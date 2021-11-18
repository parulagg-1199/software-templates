package tracing

import (
	// local packages
	"fmt"
	"io"

	// opentracing and jaeger packages
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	ot "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	config "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/zipkin"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"google.golang.org/grpc"
)

// tracer is a structure that contains service-name, jaeger host and port
type tracer struct {
	ServiceName string
	JaegerHost  string
	JeagerPort  string
}

// New is a method for getting new tracer default configuration setting
func New() tracer {
	return tracer{
		ServiceName: "default",
		JaegerHost:  "localhost",
		JeagerPort:  "6831",
	}
}

// InitJaeger is a function for intialising jaeger tracing using opentracing
func (t *tracer) InitJaeger() (ot.Tracer, io.Closer) {
	// creating new metrics factory using prometheus
	metricsFactory := prometheus.New()
	// NewZipkinB3HTTPHeaderPropagator creates a Propagator for extracting and injecting
	// Zipkin HTTP B3 headers into SpanContexts.
	propagator := zipkin.NewZipkinB3HTTPHeaderPropagator()

	// setting jeager configuration
	cfg := &config.Configuration{
		// rpc metrics enabled
		RPCMetrics: true,
		// sampler enabled
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		// reporter enabled
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: t.JaegerHost + ":" + t.JeagerPort,
		},
	}

	// using above config creating new tracer
	tracer, closer, err := cfg.New(t.ServiceName, config.ZipkinSharedRPCSpan(true), config.Injector(ot.HTTPHeaders, propagator), config.Extractor(ot.HTTPHeaders, propagator), config.Logger(jaeger.StdLogger), config.Metrics(metricsFactory))
	if err != nil {
		// any error during creation of new tracer
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	// passing new tracer
	return tracer, closer
}

// GrpcServerOptions is a function to set tracer in grpc options
func (t *tracer) GrpcServerOptions() (grpc.UnaryServerInterceptor, grpc.StreamServerInterceptor, ot.Tracer, io.Closer) {
	// creating new jaeger tracer
	tracer, closer := t.InitJaeger()

	// return unary server interceptor and stream server interceptor, tracer, closer
	return grpc_opentracing.UnaryServerInterceptor(grpc_opentracing.WithTracer(tracer)), grpc_opentracing.StreamServerInterceptor(grpc_opentracing.WithTracer(tracer)), tracer, closer
}
