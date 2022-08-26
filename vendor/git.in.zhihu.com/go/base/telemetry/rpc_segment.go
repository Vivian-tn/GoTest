package telemetry

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	spanlog "github.com/opentracing/opentracing-go/log"

	"git.in.zhihu.com/go/base/telemetry/internal/halo"
	"git.in.zhihu.com/go/base/telemetry/log"
	"git.in.zhihu.com/go/base/zae"
)

// https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/trace/semantic_conventions/rpc.md

type RPCSystem = string

const (
	RPCgRPC  RPCSystem = "gRPC"
	RPCTZone RPCSystem = "TZone"
)

const (
	rpcSegmentKey ctxKey = "rpc-segment-key"
)

func RPCSegmentFromContext(ctx context.Context) *RPCSegment {
	val := ctx.Value(rpcSegmentKey)
	if rs, ok := val.(*RPCSegment); ok {
		return rs
	}
	return nil
}

func StartRPCSegment(ctx context.Context, rs *RPCSegment) (*RPCSegment, context.Context, error) {
	openSpan, ctx := StartChildSpanWithContext(ctx, "")
	rs.openSpan = openSpan
	rs.start = time.Now()
	rs.haloSpan = &halo.ClientSpan{
		Client:        globalHaloClient,
		SlowThreshold: rs.SlowThreshold,
		Service:       zae.Service(),
		Method:        MethodFromContext(ctx),
	}
	ctx = context.WithValue(ctx, rpcSegmentKey, rs)
	return rs, ctx, nil
}

type RPCSegment struct {
	openSpan        opentracing.Span
	haloSpan        *halo.ClientSpan
	start           time.Time
	SlowThreshold   time.Duration
	System          RPCSystem
	TargetService   string
	TargetMethod    string
	// RPC 请求的参数, 仅在 IsEnableTelemetryRecordArguments() 为 true 时有可用值
	TargetArguments Arguments
}

func (s *RPCSegment) Name() string {
	return fmt.Sprintf("%s.%s/%s", s.System, s.TargetService, s.TargetMethod)
}

func (s *RPCSegment) OutboundHeaders() http.Header {
	if s.openSpan == nil {
		return nil
	}
	return InjectHTTPHeaders(s.openSpan.Context())
}

func (s *RPCSegment) End(ctx context.Context, input Error) {
	if s.System == "" {
		s.System = "unknown"
	}
	if s.TargetService == "" {
		s.TargetService = "unknown"
	}
	if s.TargetMethod == "" {
		s.TargetMethod = "unknown"
	}

	s.TargetArguments.Truncate()
	elapsed := time.Since(s.start)

	{
		s.haloSpan.TargetService = s.TargetService
		s.haloSpan.TargetMethod = s.TargetMethod

		s.haloSpan.End(elapsed, input)
	}

	{
		if s.openSpan != nil {
			s.openSpan.SetOperationName(s.Name())

			s.openSpan.SetTag("span.kind", "client")
			s.openSpan.SetTag("rpc.system", s.System)
			s.openSpan.SetTag("rpc.service", s.TargetService)
			s.openSpan.SetTag("rpc.method", s.TargetMethod)
			if enableTelemetryRecordArguments {
				s.openSpan.SetTag("rpc.arguments", s.TargetArguments)
			}

			if input != nil {
				s.openSpan.SetTag("error", true)
				s.openSpan.LogFields(spanlog.String("message", input.Error()))
			}

			s.openSpan.Finish()
		}
	}

	{
		if input != nil {
			fields := log.Fields{
				"rpc.system":    s.System,
				"rpc.service":   s.TargetService,
				"rpc.method":    s.TargetMethod,
				"elapsed":       elapsed.String(),
				"error.class":   input.Class(),
			}
			if enableTelemetryRecordArguments {
				fields["rpc.arguments"] = s.TargetArguments.ToMap()
			}
			log.WithFields(ctx, fields).WithError(input).Error(input.Error())
		}
	}
}
