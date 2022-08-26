package telemetry

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"git.in.zhihu.com/go/base/telemetry/internal/halo"
	"git.in.zhihu.com/go/base/telemetry/log"
	"git.in.zhihu.com/go/base/telemetry/statsd"
	"git.in.zhihu.com/go/base/zae"
	"github.com/opentracing/opentracing-go"
	spanlog "github.com/opentracing/opentracing-go/log"
)

var ErrInvalidSystemOrMethod = errors.New("invalid system or method")

var (
	alphanum = regexp.MustCompile(`^[a-zA-Z0-9-_ ]+$`)
)

const (
	serverSegmentKey ctxKey = "server-segment-key"
)

type TransactionSystem = string

const (
	TransactionHTTP   TransactionSystem = "HTTP"
	TransactiongRPC   TransactionSystem = "gRPC"
	TransactionTZone  TransactionSystem = "TZone"
	TransactionWorker TransactionSystem = "Worker"
	TransactionExec   TransactionSystem = "Exec"
)

var (
	disableTransactionInfoLog      bool
	enableTelemetryRecordArguments bool
)

func init() {
	if os.Getenv("DISABLE_TRANSACTION_INFO_LOG") == "1" {
		disableTransactionInfoLog = true
	}
	if os.Getenv("ENABLE_TELEMETRY_RECORD_ARGUMENTS") == "1" {
		enableTelemetryRecordArguments = true
	}
}

func IsEnableTelemetryRecordArguments() bool {
	return enableTelemetryRecordArguments
}

func TransactionFromContext(ctx context.Context) *Transaction {
	val := ctx.Value(serverSegmentKey)
	if txn, ok := val.(*Transaction); ok {
		return txn
	}
	return nil
}

func MethodFromContext(ctx context.Context) (method string) {
	if txn := TransactionFromContext(ctx); txn != nil {
		method = txn.Method
	}
	if method != "" {
		return method
	}
	return "unknown"
}

func StartTransaction(ctx context.Context, txn *Transaction, opts ...opentracing.StartSpanOption) (*Transaction, context.Context, error) {
	txn.System, txn.Method = strings.TrimSpace(txn.System), strings.TrimSpace(txn.Method)
	if !alphanum.MatchString(txn.System) || (len(txn.Method) > 0 && !alphanum.MatchString(txn.Method)) {
		return nil, ctx, ErrInvalidSystemOrMethod
	}

	openSpan, ctx := StartRootSpanWithContext(ctx, "", opts...)
	txn.openSpan = openSpan
	txn.start = time.Now()
	txn.haloSpan = &halo.ServerSpan{
		Client:  globalHaloClient,
		Service: zae.Service(),
	}
	ctx = context.WithValue(ctx, serverSegmentKey, txn)
	return txn, ctx, nil
}

type Transaction struct {
	openSpan   opentracing.Span
	haloSpan   *halo.ServerSpan
	start      time.Time
	Request    *http.Request
	StatusCode int
	System     TransactionSystem
	Method     string
	Arguments  Arguments
	input      Error
}

func (s *Transaction) Name() string {
	return fmt.Sprintf("%s.%s", s.System, s.Method)
}

func (s *Transaction) SetRequest(req *http.Request) {
	s.Request = req
}

func (s *Transaction) SetStatusCode(code int) {
	s.StatusCode = code
}

func (s *Transaction) HTTPURL() string {
	if s.Request != nil {
		return s.Request.URL.String()
	}
	return ""
}

func (s *Transaction) HTTPMethod() string {
	if s.Request != nil {
		if s.Request.Method != "" {
			return s.Request.Method
		}
		// Golang's http package states that when a client's Request has
		// an empty string for Method, the method is GET.
		return "GET"
	}
	return ""
}

func (s *Transaction) SetError(input Error) {
	s.input = input
}

func (s *Transaction) Error() Error {
	return s.input
}

func (s *Transaction) End(ctx context.Context, input Error) {
	if s.input != nil {
		input = s.input
	}
	if input == nil && s.StatusCode >= http.StatusInternalServerError {
		input = NewClassErr(http.StatusText(s.StatusCode))
	}

	if s.System == "" {
		s.System = "unknown"
	}
	if s.Method == "" {
		s.Method = "unknown"
	}

	s.Method = statsd.Node(s.Method)
	s.Arguments.Truncate()

	elapsed := time.Since(s.start)

	{
		s.haloSpan.Method = s.Method

		s.haloSpan.End(elapsed, input)
	}

	{
		s.openSpan.SetOperationName(s.Name())

		s.openSpan.SetTag("span.kind", "server")
		s.openSpan.SetTag("server.system", s.System)
		s.openSpan.SetTag("server.method", s.Method)
		if enableTelemetryRecordArguments {
			s.openSpan.SetTag("server.arguments", s.Arguments)
		}

		if s.Request != nil {
			s.openSpan.SetTag("http.url", s.HTTPURL())
			s.openSpan.SetTag("http.method", s.HTTPMethod())
			s.openSpan.SetTag("http.proto", s.Request.Proto)
			s.openSpan.SetTag("http.request_content_length", s.Request.ContentLength)
			s.openSpan.SetTag("http.status_code", s.StatusCode)
		}

		if input != nil {
			s.openSpan.SetTag("error", true)
			s.openSpan.LogFields(spanlog.String("message", input.Error()))
		}

		s.openSpan.Finish()
	}

	{
		fields := log.Fields{
			"server.system":    s.System,
			"server.method":    s.Method,
			"elapsed":          elapsed.String(),
		}
		if enableTelemetryRecordArguments {
			fields["server.arguments"] = s.Arguments.ToMap()
		}
		if s.Request != nil {
			fields["http.url"] = s.HTTPURL()
			fields["http.method"] = s.HTTPMethod()
			fields["http.status_code"] = s.StatusCode
			fields["http.request_remote_addr"] = s.Request.RemoteAddr
			fields["http.request_content_length"] = s.Request.ContentLength
			if s.System == TransactionHTTP {
				appVersion := s.Request.Header.Get("X-APP-Version")
				if appVersion != "" {
					fields["http.app_version"] = appVersion
				}
				apiVersion := s.Request.Header.Get("X-API-Version")
				if apiVersion != "" {
					fields["http.api_version"] = apiVersion
				}
			}
		}

		if input == nil {
			if disableTransactionInfoLog {
				return
			}

			log.WithFields(ctx, fields).Info("finished")

			if s.System == TransactionHTTP && (s.StatusCode >= 400 && s.StatusCode < 500) {
				input = NewClassErr(http.StatusText(s.StatusCode))
				fields["error.class"] = input.Class()
				if s.Request != nil {
					fields["sentry.http_request"] = s.Request
				}
				log.WithFields(ctx, fields).WithError(input).Warnf(input.Error())
			}
		} else {
			fields["error.class"] = input.Class()
			if s.Request != nil {
				fields["sentry.http_request"] = s.Request
			}
			log.WithFields(ctx, fields).WithError(input).Error(input.Error())
		}
	}
}
