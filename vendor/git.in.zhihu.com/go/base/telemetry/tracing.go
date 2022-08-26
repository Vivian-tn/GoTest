package telemetry

import (
	"context"
	"net/http"

	"git.in.zhihu.com/go/base/telemetry/internal/tracer/jaeger"
	"git.in.zhihu.com/go/base/telemetry/log"
	"git.in.zhihu.com/go/base/zae"
	"github.com/opentracing/opentracing-go"
)

type ctxKey string

const (
	traceIDKey = "X-B3-Traceid" // equal with log, cannot be used the ctxKey type
)

func init() {
	if zae.Service() == "" {
		panic("can't find the service name")
	}

	tracer, err := jaeger.New(&jaeger.Config{
		Service:   zae.Service(),
		AgentHost: zae.Host(),
	})
	if err != nil {
		panic(err)
	}
	opentracing.SetGlobalTracer(tracer)
}

func InjectHTTPHeaders(sc opentracing.SpanContext) http.Header {
	headers := make(http.Header)
	if err := opentracing.GlobalTracer().Inject(sc, opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(headers)); err != nil {
		log.Fatalf(context.TODO(), "all opentracing.Tracer implementations MUST support all BuiltinFormats: %s", err)
	}
	return headers
}

func ExtractHTTPHeaders(headers http.Header) opentracing.StartSpanOption {
	sc, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(headers))
	return opentracing.ChildOf(sc)
}

func StartRootSpanWithContext(ctx context.Context, op string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	if !zae.IsDevelopEnv() && !zae.IsCIEnv() && !opentracing.IsGlobalTracerRegistered() {
		panic("can't find the global tracer")
	}

	if TraceIDFromContext(ctx) != "" {
		panic("can't set root span twice in the same context")
	}

	openSpan, ctx := opentracing.StartSpanFromContext(ctx, op, opts...)

	headers := InjectHTTPHeaders(openSpan.Context())
	ctx = ContextWithTraceID(ctx, headers.Get(traceIDKey))
	return openSpan, ctx
}

func StartChildSpanWithContext(ctx context.Context, op string, opts ...opentracing.StartSpanOption) (opentracing.Span, context.Context) {
	if parentSpan := opentracing.SpanFromContext(ctx); parentSpan == nil {
		return nil, ctx
	}
	return opentracing.StartSpanFromContext(ctx, op, opts...)
}

func TraceIDFromContext(ctx context.Context) string {
	val := ctx.Value(traceIDKey)
	if sp, ok := val.(string); ok {
		return sp
	}
	return ""
}

func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

func BaggageItemFromContext(ctx context.Context, restrictedKey string) string {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return ""
	}
	return span.BaggageItem(restrictedKey)
}
