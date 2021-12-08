package halo

import (
	"time"

	"git.in.zhihu.com/go/base/telemetry/statsd"
)

type ClientSpan struct {
	Client        statsd.Client
	SlowThreshold time.Duration
	Service       string
	Method        string
	TargetService string
	TargetMethod  string
}

func (s *ClientSpan) End(elapsed time.Duration, input Error) {
	s.Client.Increment(s.Count())

	if input != nil {
		s.Client.Increment(s.Error(input))
	}

	if s.SlowThreshold != 0 && elapsed > s.SlowThreshold {
		s.Client.Increment(s.SlowLog())
	}

	s.Client.Timing(s.Timing(), elapsed)
}

func (s *ClientSpan) Timing() string {
	return s.withSuffix("request_time")
}

func (s *ClientSpan) Count() string {
	return s.withSuffix("count")
}

func (s *ClientSpan) Error(err Error) string {
	return s.withSuffix("error", statsd.Node(err.Class()), "count")
}

func (s *ClientSpan) SlowLog() string {
	return s.withSuffix("slow", "count")
}

func (s *ClientSpan) withSuffix(suffix ...string) string {
	return statsd.Join(
		statsd.Node(s.Service),
		"_all",
		"client",
		statsd.Node(s.Method),
		statsd.Node(s.TargetService),
		"_all",
		statsd.Node(s.TargetMethod),
		statsd.Join(suffix...),
	)
}
