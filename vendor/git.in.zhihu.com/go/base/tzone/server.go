package tzone

import (
	"context"
	"net/http"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"git.in.zhihu.com/go/base/telemetry"
	"git.in.zhihu.com/go/base/tzone/internal/decoder"
	"git.in.zhihu.com/go/base/tzone/internal/errors"
	"git.in.zhihu.com/go/base/tzone/internal/middleware"
)

type TProcessor = thrift.TProcessor
type Caller = middleware.Caller

func CallerFromContext(ctx context.Context) *Caller {
	return middleware.CallerFromContext(ctx)
}

func NewServer(services map[string]TProcessor) *Server {
	s := &Server{
		processor:       thrift.NewTMultiplexedProcessor(),
		protocolFactory: thrift.NewTBinaryProtocolFactoryDefault(),
	}
	for serviceName, serviceProcessor := range services {
		s.processor.RegisterProcessor(serviceName, serviceProcessor)
	}

	s.Use(telemetry.Middleware(telemetry.TransactionTZone, nil))
	s.Use(middleware.Auth)
	s.Use(middleware.InjectCaller)

	return s
}

type Server struct {
	httpServer      *http.Server
	processor       *thrift.TMultiplexedProcessor
	protocolFactory *thrift.TBinaryProtocolFactory
	middlewares     []func(http.Handler) http.Handler
}

func (s *Server) Use(middlewares ...func(http.Handler) http.Handler) *Server {
	s.middlewares = append(s.middlewares, middlewares...)
	return s
}

func (s *Server) Chain(h http.Handler) http.Handler {
	for i := range s.middlewares {
		h = s.middlewares[len(s.middlewares)-1-i](h)
	}
	return h
}

func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	txn := telemetry.TransactionFromContext(r.Context())

	inputTransport := thrift.NewStreamTransportR(r.Body)
	inputProtocol := s.protocolFactory.GetProtocol(inputTransport)

	outputTransport := thrift.NewTMemoryBuffer()
	outputProtocol := s.protocolFactory.GetProtocol(outputTransport)

	sniffer := r.Body.(*telemetry.RequestSniffer)
	sniffer.Start()

	method, err := decoder.ReadMethod(sniffer)
	if err != nil {
		txn.SetError(errors.WrapError(err))
		return
	}
	txn.Method = method

	sniffer.Stop()

	sniffer.Start()
	defer func() {
		sniffer.Stop()

		if method, args, err := decoder.DumpBody(sniffer.Bytes()); err == nil {
			if txn.Method == "" {
				txn.Method = method
			}
			txn.Arguments = args
		} else {
			if txn.Method == "" {
				txn.Method = r.Header.Get("X-ZONE-API")
			}
			txn.Arguments = telemetry.LargeArguments()
		}
	}()

	_, err = s.processor.Process(r.Context(), inputProtocol, outputProtocol)
	if err != nil {
		txn.SetError(errors.WrapError(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/x-thrift")
	_, _ = outputTransport.WriteTo(w)
}

func (s *Server) Run(addr string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/check_health", func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("zhi~"))
	})
	mux.Handle("/", s.Chain(http.HandlerFunc(s.Handler)))

	s.httpServer = &http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Close() error {
	return s.httpServer.Close()
}
