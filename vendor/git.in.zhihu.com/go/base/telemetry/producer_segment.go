package telemetry

import (
	"context"
	"fmt"
	"time"

	"git.in.zhihu.com/go/base/telemetry/internal/halo"
	"git.in.zhihu.com/go/base/telemetry/log"
	"git.in.zhihu.com/go/base/zae"
	"github.com/opentracing/opentracing-go"
	spanlog "github.com/opentracing/opentracing-go/log"
)

type ProducerSystem = string

const (
	ProducerKafka  ProducerSystem = "Kafka"
	ProducerPulsar ProducerSystem = "Pulsar"
)

func StartProducerSegment(ctx context.Context, ps *ProducerSegment) (*ProducerSegment, context.Context, error) {
	openSpan, ctx := StartChildSpanWithContext(ctx, "")
	ps.openSpan = openSpan
	ps.haloSpan = &halo.ClientSpan{
		Client:  globalHaloClient,
		Service: zae.Service(),
		Method:  MethodFromContext(ctx),
	}
	ps.start = time.Now()
	return ps, ctx, nil
}

type ProducerSegment struct {
	openSpan opentracing.Span
	haloSpan *halo.ClientSpan
	start    time.Time
	System   ProducerSystem
	Topic    string
	Sync     bool
}

func (s *ProducerSegment) Name() string {
	return fmt.Sprintf("%s.%s/send", s.System, s.Topic)
}

func (s *ProducerSegment) End(ctx context.Context, input Error) {
	if s.System == "" {
		s.System = "unknown"
	}
	if s.Topic == "" {
		s.Topic = "unknown"
	}

	elapsed := time.Since(s.start)

	{
		s.haloSpan.TargetService = fmt.Sprintf("%s_%s", s.System, s.Topic)
		if s.Sync {
			s.haloSpan.TargetMethod = "sync_send"
		} else {
			s.haloSpan.TargetMethod = "async_send"
		}

		s.haloSpan.End(elapsed, input)
	}

	{
		if s.openSpan != nil {
			s.openSpan.SetOperationName(s.Name())

			s.openSpan.SetTag("span.kind", "producer")
			s.openSpan.SetTag("messaging.system", s.System)
			s.openSpan.SetTag("messaging.destination", s.Topic)
			s.openSpan.SetTag("messaging.sync", s.Sync)

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
				"messaging.system":      s.System,
				"messaging.destination": s.Topic,
				"elapsed":               elapsed.String(),
				"error.class":           input.Class(),
			}
			log.WithFields(ctx, fields).WithError(input).Error(input.Error())
		}
	}
}
