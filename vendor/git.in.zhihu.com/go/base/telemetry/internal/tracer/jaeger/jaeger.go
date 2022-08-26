package jaeger

import (
	"fmt"
	"net"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/zipkin"
)

// OptimalUDPPayloadSize defines the optimal payload size for a UDP datagram, 1432 bytes
// is optimal for regular networks with an MTU of 1500 so datagrams don't get
// fragmented. It's generally recommended not to fragment UDP datagrams as losing
// a single fragment will cause the entire datagram to be lost.
const OptimalUDPPayloadSize = 1432

type Config struct {
	Service   string
	AgentHost string
}

func New(config *Config) (opentracing.Tracer, error) {
	var sampler jaeger.Sampler

	samplingServer := fmt.Sprintf("%s:%d", config.AgentHost, jaeger.DefaultSamplingServerPort)
	if checkAvailable(samplingServer) {
		sampler = jaeger.NewRemotelyControlledSampler(
			config.Service,
			jaeger.SamplerOptions.SamplingStrategyFetcher(
				newStrategyFetcher(fmt.Sprintf("http://%s/sampling", samplingServer)),
			),
		)
	} else {
		sampler, _ = jaeger.NewProbabilisticSampler(0.001)
	}

	sender, err := jaeger.NewUDPTransport(
		fmt.Sprintf("%s:%d", config.AgentHost, jaeger.DefaultUDPSpanServerPort),
		OptimalUDPPayloadSize,
	)
	if err != nil {
		return nil, err
	}
	reporter := jaeger.NewRemoteReporter(sender)

	// https://github.com/jaegertracing/jaeger-client-go/blob/v2.25.0/zipkin/README.md#zipkin-compatibility-features
	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	injector := jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, zipkinPropagator)
	extractor := jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, zipkinPropagator)
	zipkinSharedRPCSpan := jaeger.TracerOptions.ZipkinSharedRPCSpan(true)

	tracer, _ := jaeger.NewTracer(config.Service, sampler, reporter, injector, extractor, zipkinSharedRPCSpan)
	return tracer, nil
}

func checkAvailable(address string) bool {
	conn, err := net.DialTimeout("tcp", address, time.Second)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
