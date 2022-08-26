package telemetry

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"sync"

	"git.in.zhihu.com/go/base/telemetry/log"
	"git.in.zhihu.com/go/base/telemetry/sentry"
)

var (
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
)

type responseRecorder struct {
	http.ResponseWriter
	StatusCode int
}

func newResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
	}
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func NewRequestSniffer(source io.ReadCloser, buffer *bytes.Buffer, limit int) *RequestSniffer {
	if source == nil {
		return nil
	}
	return &RequestSniffer{source: source, buffer: buffer, limit: limit}
}

type RequestSniffer struct {
	ctx        context.Context
	source     io.ReadCloser
	buffer     *bytes.Buffer
	limit      int
	bufferRead int
	bufferSize int
	sniffing   bool
	lastErr    error
}

func (s *RequestSniffer) Read(p []byte) (int, error) {
	if s.bufferSize > s.bufferRead {
		bn := copy(p, s.buffer.Bytes()[s.bufferRead:s.bufferSize])
		s.bufferRead += bn
		return bn, s.lastErr
	}

	sn, sErr := s.source.Read(p)
	if sn > 0 && s.sniffing {
		s.lastErr = sErr
		// 🚨 🚨 🚨 不要滥用这个功能，很危险 🚨 🚨 🚨
		// buffer 剩余空间不足时，只能保证第一次读的有效，重复读的数据超过 limit 的部分会错位，参见单元测试
		// 所以要保证 sniffer 单次的数据不超过 limit 长度，不要传长度大于 limit 的 slice 进来
		rbn := s.limit - s.buffer.Len()
		if sn > rbn {
			// 尽量写一点到 buffer 里，可能不全
			if rbn > 0 {
				_, _ = s.buffer.Write(p[:rbn])
			}
		} else {
			if wn, wErr := s.buffer.Write(p[:sn]); wErr != nil {
				return wn, wErr
			}
		}
	}
	return sn, sErr
}

func (s *RequestSniffer) Start() {
	s.reset(true)
}

func (s *RequestSniffer) Stop() {
	s.reset(false)
}

func (s *RequestSniffer) reset(snif bool) {
	s.sniffing = snif
	s.bufferRead = 0
	s.bufferSize = s.buffer.Len()
}

func (s *RequestSniffer) Bytes() []byte {
	return s.buffer.Bytes()
}

func (s *RequestSniffer) Close() error {
	return s.source.Close()
}

func Middleware(system TransactionSystem, f func(r *http.Request) string) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			sentry.Recover(func() {
				var method string
				if f != nil {
					method = f(r)
				}

				if system == TransactionHTTP && method == "" {
					// 兼容 go/box Stats Middleware 处理，忽略业务未定义的路由
					next.ServeHTTP(w, r)
					return
				}

				buffer := bufferPool.Get().(*bytes.Buffer)
				defer bufferPool.Put(buffer)
				buffer.Reset()

				txn, ctx, err := StartTransaction(r.Context(), &Transaction{
					System: system,
					Method: method,
				}, ExtractHTTPHeaders(r.Header))
				if err != nil {
					w.WriteHeader(http.StatusTooManyRequests)
					log.Error(r.Context(), "start transaction failed %w", err)
					return
				}

				r = r.WithContext(ctx)
				r.Body = NewRequestSniffer(r.Body, buffer, 256)
				txn.SetRequest(r)

				wr := newResponseRecorder(w)

				sentry.Recover(func() {
					next.ServeHTTP(wr, r)
					err = txn.Error()
				}, func(e error) {
					err = e
				})
				if err != nil {
					wr.WriteHeader(http.StatusInternalServerError)
					txn.End(ctx, WrapErrWithUnknownClass(err))
					return
				}

				txn.SetStatusCode(wr.StatusCode)
				txn.End(ctx, nil)
			}, func(err error) {
				w.WriteHeader(http.StatusInternalServerError)
				log.Error(r.Context(), err)
			})
		}
		return http.HandlerFunc(fn)
	}
}
