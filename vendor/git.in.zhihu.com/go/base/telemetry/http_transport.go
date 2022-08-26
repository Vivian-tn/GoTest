package telemetry

import (
	"bytes"
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptrace"
)

func init() {
	http.DefaultTransport = WrapRoundTripper(http.DefaultTransport)
}

func WrapRoundTripperWithService(original http.RoundTripper, service string) http.RoundTripper {
	return WrapRoundTripper(roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		r.Header.Set("x-telemtry-service", service)
		return original.RoundTrip(r)
	}))
}

func WrapRoundTripper(original http.RoundTripper) http.RoundTripper {
	return roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		// The specification of http.RoundTripper requires that the r is never modified.
		ctx := r.Context()
		r = cloneRequest(r)

		hs, ctx, err := StartHTTPSegment(ctx, r)
		if err != nil {
			return nil, err
		}

		buffer := bufferPool.Get().(*bytes.Buffer)
		defer bufferPool.Put(buffer)
		buffer.Reset()

		var sniffer *RequestSniffer
		if r.Body != nil {
			sniffer = NewRequestSniffer(r.Body, buffer, 256)
			sniffer.Start()
			r.Body = sniffer
		}

		var remoteAddr string
		trace := &httptrace.ClientTrace{
			GotConn: func(connInfo httptrace.GotConnInfo) {
				// FIX: https://github.com/golang/go/issues/34282
				if connInfo.Conn != nil && connInfo.Conn.RemoteAddr() != nil {
					remoteAddr = connInfo.Conn.RemoteAddr().String()
				} else {
					remoteAddr = "-"
				}

			},
		}
		r = r.WithContext(httptrace.WithClientTrace(ctx, trace))

		resp, err := original.RoundTrip(r)

		if resp != nil {
			resp.Request.RemoteAddr = remoteAddr
		}

		if sniffer != nil {
			sniffer.Stop()
		}

		hs.Response = resp

		var input Error
		if errors.Is(err, context.Canceled) {
			input = WrapErr(err, "CtxCanceled")
		} else if errors.Is(err, context.DeadlineExceeded) {
			input = WrapErr(err, "CtxDeadline")
		} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			input = WrapErr(netErr, "Timeout")
		} else if opErr, ok := err.(*net.OpError); ok {
			input = WrapErr(err, opErr.Op)
		} else {
			input = WrapErrWithUnknownClass(err)
		}
		hs.End(ctx, input)

		return resp, err
	})
}

// cloneRequest mimics implementation of https://godoc.org/github.com/google/go-github/github#BasicAuthTransport.RoundTrip
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }
