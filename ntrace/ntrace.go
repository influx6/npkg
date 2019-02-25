package ntrace

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gokit/npkg"

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

// Trace defines a function which giving provided context, traceName will call provided function
// with the trace created and finished. Note if tracing is disabled, then a nil span is provided
// to passed in function.
//
// WARNING: The do function may receive a nil span, if no span was found within context, has
// this is used to indicate if tracing was enabled.
func Trace(ctx context.Context, traceName string, do func(context.Context, opentracing.Span)) {
	var span opentracing.Span
	if ctx, span = NewSpanFromContext(ctx, traceName); span != nil {
		defer span.Finish()
	}

	do(ctx, span)
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
// Service Name in Jaeger Query - edit accordingly to the Kubernetes service name (ask us if you don't know)
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
