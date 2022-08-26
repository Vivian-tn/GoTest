package jaeger

import (
	"context"
	"fmt"
	"net/http"

	"git.in.zhihu.com/go/base/internal/request"
	"github.com/uber/jaeger-client-go"
)

var _ jaeger.SamplingStrategyFetcher = new(strategyFetcher)

func newStrategyFetcher(serverURL string) *strategyFetcher {
	return &strategyFetcher{
		url: serverURL,
		req: request.New(),
	}
}

type strategyFetcher struct {
	url string
	req *request.Request
}

func (f *strategyFetcher) Fetch(serviceName string) ([]byte, error) {
	resp, err := f.req.Get(context.TODO(), f.url, request.Query{
		"service": serviceName,
	}, request.Headers{
		"x-telemtry-service": "jaeger-sampling-strategy",
	})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("StatusCode: %d, Body: %s", resp.StatusCode, resp.String())
	}
	return resp.ReadAll()
}
