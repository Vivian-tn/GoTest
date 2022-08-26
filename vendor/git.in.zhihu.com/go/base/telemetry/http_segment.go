package telemetry

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"git.in.zhihu.com/go/base/telemetry/internal/halo"
	"git.in.zhihu.com/go/base/telemetry/log"
	"git.in.zhihu.com/go/base/telemetry/statsd"
	"git.in.zhihu.com/go/base/zae"
	"github.com/opentracing/opentracing-go"
	spanlog "github.com/opentracing/opentracing-go/log"
)

var ErrInvalidRequest = errors.New("must set *http.Request in StartHTTPSegment")

func StartHTTPSegment(ctx context.Context, req *http.Request) (*HTTPSegment, context.Context, error) {
	if req == nil {
		return nil, nil, ErrInvalidRequest
	}
	openSpan, ctx := StartChildSpanWithContext(ctx, "")
	es := &HTTPSegment{
		openSpan: openSpan,
		haloSpan: &halo.ClientSpan{
			Client:  globalHaloClient,
			Service: zae.Service(),
			Method:  MethodFromContext(ctx),
		},
		start:   time.Now(),
		Request: req,
	}
	if req.Header == nil {
		req.Header = make(http.Header)
	}
	for key, values := range es.OutboundHeaders() {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	return es, ctx, nil
}

type HTTPSegment struct {
	openSpan opentracing.Span
	haloSpan *halo.ClientSpan
	start    time.Time
	Request  *http.Request
	Response *http.Response
}

func (s *HTTPSegment) Name() string {
	return fmt.Sprintf("http.%s/%s", s.Service(), s.Method())
}

func (s *HTTPSegment) OutboundHeaders() http.Header {
	if s.openSpan == nil {
		return nil
	}
	return InjectHTTPHeaders(s.openSpan.Context())
}

func (s *HTTPSegment) URL() *url.URL {
	r := s.Request
	if s.Response != nil && s.Response.Request != nil {
		r = s.Response.Request
	}
	if r.Host != "" {
		r.URL.Host = r.Host
	}
	return r.URL
}

func (s *HTTPSegment) Service() string {
	service := s.Request.Header.Get("x-telemtry-service")
	if service != "" {
		return service
	}
	if host := s.URL().Host; host != "" {
		return statsd.Node(host)
	}
	return "unknown"
}

func (s *HTTPSegment) Method() string {
	if s.Request.Method != "" {
		return s.Request.Method
	}
	// Golang's http package states that when a client's Request has
	// an empty string for Method, the method is GET.
	return "GET"
}

func (s *HTTPSegment) End(ctx context.Context, input Error) {
	if input == nil && s.Response != nil && s.Response.StatusCode >= http.StatusBadRequest {
		input = NewClassErr(http.StatusText(s.Response.StatusCode))
	}

	elapsed := time.Since(s.start)

	{
		s.haloSpan.TargetService = fmt.Sprintf("HTTP_%s", s.Service())
		s.haloSpan.TargetMethod = s.Method()

		s.haloSpan.End(elapsed, input)
	}

	{
		if s.openSpan != nil {
			s.openSpan.SetOperationName(s.Name())

			s.openSpan.SetTag("span.kind", "client")
			s.openSpan.SetTag("http.method", s.Method())
			s.openSpan.SetTag("http.url", s.URL().String())
			s.openSpan.SetTag("http.proto", s.Request.Proto)
			s.openSpan.SetTag("http.request_content_length", s.Request.ContentLength)

			if s.Response != nil {
				s.openSpan.SetTag("http.status_code", s.Response.StatusCode)
				s.openSpan.SetTag("http.response_content_length", s.Response.ContentLength)
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
				"http.url":            s.URL().String(),
				"http.method":         s.Method(),
				"elapsed":             elapsed.String(),
				"error.class":         input.Class(),
				"sentry.http_request": s.Request,
			}
			if s.Response != nil {
				fields["http.status_code"] = s.Response.StatusCode
				fields["sentry.http_response"] = s.Response
			}
			log.WithFields(ctx, fields).WithError(input).Error(input.Error())
		}
	}
}
