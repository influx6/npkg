package ntrace

import (
	"context"
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-client-go/zipkin"
)

// Trace defines a function
func Trace(ctx context.Context, traceName string, do func(context.Context, opentracing.Span)) {

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
