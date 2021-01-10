package ntrace

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/influx6/npkg"
	"github.com/influx6/npkg/nframes"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-client-go/zipkin"
)

const (
	// SpanKey provides giving key-name used to store open-tracing span into the context.
	SpanKey = contextKey("SPAN_KEY")
)

type contextKey string

// WithKV as new key-value pair into the spans baggage store.
//
// This values get propagated down to child spans from this span.
func WithKV(span opentracing.Span, key string, value string) {
	if span == nil {
		return
	}
	span.SetBaggageItem(key, value)
}

// WithTag adds a new tag into a giving span.
func WithTag(span opentracing.Span, key string, value interface{}) {
	if span == nil {
		return
	}
	span.SetTag(key, value)
}

// WithTrace returns a new context.Context and a function which can be used to finish the
// opentracing.Span attached to giving context.
//
// It's an alternative to using Trace.
func WithTrace(ctx context.Context, methodName string) (context.Context, func()) {
	var span opentracing.Span
	ctx, span = NewSpanFromContext(ctx, methodName)

	return ctx, func() {
		if span == nil {
			return
		}
		span.Finish()
	}
}

// WithTrace returns a new context.Context and a function which can be used to finish the
// opentracing.Span attached to giving context.
//
// It's an alternative to using Trace.
func WithMethodTrace(ctx context.Context) (context.Context, func()) {
	return WithTrace(ctx, nframes.GetCallerNameWith(2))
}

// GetSpanFromContext returns a OpenTracing span if available from provided context.
//
// WARNING: Second returned value can be nil if no parent span is in context.
func GetSpanFromContext(ctx context.Context) (opentracing.Span, bool) {
	if span, ok := ctx.Value(SpanKey).(opentracing.Span); ok {
		return span, true
	}
	return nil, false
}

// NewMethodSpanFromContext returns a OpenTracing span if available from provided context.
// It automatically gets the caller name using runtime.
//
// WARNING: Second returned value can be nil if no parent span is in context.
func NewMethodSpanFromContext(ctx context.Context) (context.Context, opentracing.Span) {
	return NewSpanFromContext(ctx, nframes.GetCallerNameWith(2))
}

// NewSpanFromContext returns a new OpenTracing child span which was created as a child of a
// existing span from the underline context if it exists, else returning no Span.
//
// WARNING: Second returned value can be nil if no parent span is in context.
func NewSpanFromContext(ctx context.Context, traceName string) (context.Context, opentracing.Span) {
	if span, ok := ctx.Value(SpanKey).(opentracing.Span); ok {
		var childSpan = opentracing.StartSpan(
			traceName,
			opentracing.ChildOf(span.Context()),
		)

		var newContext = context.WithValue(ctx, SpanKey, childSpan)
		return newContext, childSpan
	}
	return ctx, nil
}

// CloseOpenTracingMWSpan returns a new middleware for closing the OpenTracing span which
// then writes trace to underline tracing service.
func CloseOpenTracingMWSpan(getter npkg.Getter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if span, ok := (*r).Context().Value(SpanKey).(opentracing.Span); ok {
				// Get debug flag if enabled.
				var debug = getter.Bool(npkg.DEBUGKey)
				if debug {
					log.Printf("Closing span, path: " + r.URL.Path)
				}

				defer span.Finish()
			}
			next.ServeHTTP(w, r)
		})
	}
}

// OpenTracingSpanMW returns a middleware function able to based on tracing configuration key, enable
// and setup tracing using opentracing spans.
func OpenTracingMW(tracingKey string, getter npkg.Getter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var methodName = requestMethodParser(r)

			// Get debug flag if enabled.
			var debug = getter.Bool(npkg.DEBUGKey)

			// Get tracing flag if enabled, else skip
			var enabled = getter.Bool(tracingKey)
			if !enabled {
				next.ServeHTTP(w, r)
				return
			}

			var serverSpan opentracing.Span

			// Extracting B3 tracing context from the request.
			// This step is important to extract the actual request context
			// from outside of the applications.
			var wireContext, err = opentracing.GlobalTracer().Extract(
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(r.Header),
			)

			if err != nil {
				if debug {
					log.Printf("Attaching span fails err: " + err.Error())
					for k, h := range r.Header {
						for _, v := range h {
							log.Printf(fmt.Sprintf("Header: %s - %s", k, v))
						}
					}
				}

				// Create span as a parent without parent span.
				serverSpan = opentracing.StartSpan(
					methodName,
				)
			} else {
				// Create span as a child of parent Span from wireContext.
				serverSpan = opentracing.StartSpan(
					methodName,
					opentracing.ChildOf(wireContext),
				)
			}

			if debug {
				// Optional: Record span creation
				log.Printf("Attaching span, Starting child span: " + methodName)
			}

			if traceID := serverSpan.BaggageItem("trace_id"); traceID != "" {
				serverSpan.SetTag("trace_id", traceID)
			} else {
				var traceID = r.Header.Get("X-Request-Id")
				if traceID != "" {
					serverSpan.SetBaggageItem("trace_id", traceID)
					serverSpan.SetTag("trace_id", traceID)
				}
			}

			// We are passing the context as an item in Go context. Span is also attached
			// so that we can close the span after the request. Span needs to be finished
			// in order to report it to Jaeger collector
			var newContext = context.WithValue(r.Context(), SpanKey, serverSpan)
			next.ServeHTTP(w, r.WithContext(newContext))
		})
	}
}

// InitZipkinTracer inits the jaeger client
// Service PageName in Jaeger Query - edit accordingly to the Kubernetes service name (ask us if you don't know)
func InitZipkinTracer(srvName string, samplingServerURL, localAgentHost string) (io.Closer, error) {
	// Standard option for starting Zipkin based Tracing, please include.
	var jLogger = jaegerlog.StdLogger
	var zipkinPropagator = zipkin.NewZipkinB3HTTPHeaderPropagator()
	var injector = jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, zipkinPropagator)
	var extractor = jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, zipkinPropagator)
	var zipkinSharedRPCSpan = jaeger.TracerOptions.ZipkinSharedRPCSpan(true)

	// Create new jaeger reporter and sampler. URL must be fixed.
	var samplercfg = &jaegercfg.SamplerConfig{
		Type:              jaeger.SamplerTypeConst,
		Param:             1,
		SamplingServerURL: samplingServerURL,
	}

	var reportercfg = &jaegercfg.ReporterConfig{
		LogSpans:           true,
		LocalAgentHostPort: localAgentHost,
	}

	// Jaeger sampler and reporter
	var jMetrics = jaeger.NewNullMetrics()

	var sampler, err = samplercfg.NewSampler(srvName, jMetrics)
	if err != nil {
		return nil, err
	}

	var reporter jaeger.Reporter
	reporter, err = reportercfg.NewReporter(srvName, jMetrics, jLogger)
	if err != nil {
		return nil, err
	}

	var tracer, closer = jaeger.NewTracer(srvName, sampler, reporter, injector, extractor, zipkinSharedRPCSpan)
	opentracing.SetGlobalTracer(tracer)

	return closer, nil
}

var requestMethodParser = func(r *http.Request) string {
	var path = r.URL.Path
	var count int
	for i, c := range path {
		if count == 3 {
			return path[:i]
		}
		if c == '/' {
			count++
		}
	}
	return r.Method + path
}
