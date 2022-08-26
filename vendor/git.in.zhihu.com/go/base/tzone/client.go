package tzone

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
	"git.in.zhihu.com/go/base/internal/diplomat"
	"git.in.zhihu.com/go/base/telemetry"
	"git.in.zhihu.com/go/base/tzone/internal/decoder"
	"git.in.zhihu.com/go/base/tzone/internal/errors"
	"git.in.zhihu.com/go/base/zae"
)

func NewClient(serviceName string, opts ...Option) *Client {
	c := &Client{
		timeout:     500 * time.Millisecond,
		serviceName: serviceName,
		discovery:   diplomat.Discover(),
	}
	for _, opt := range opts {
		opt(c)
	}

	if c.targetName == "" && c.hostPort == "" {
		panic("client: either targetName or HostPort option must be specified.")
	}
	if c.targetName == "" {
		c.targetName = "localhost"
	}

	if c.hostPort != "" {
		c.discovery.RegisterLocal(c.targetName, c.hostPort)
	}

	c.httpClient = &http.Client{
		Timeout: c.timeout,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          10240,
			MaxIdleConnsPerHost:   1024,
			IdleConnTimeout:       1 * time.Second, // 过长的存活时间链接可能会被 server 关掉
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	return c
}

type Client struct {
	targetName      string
	timeout         time.Duration
	serviceName     string
	discovery       *diplomat.Diplomat
	hostPort        string
	httpClient      *http.Client
	thriftTransport *thrift.TTransport
}

func (c *Client) Call(ctx context.Context, method string, args, result thrift.TStruct) error {
	var targetArgs telemetry.Arguments
	if telemetry.IsEnableTelemetryRecordArguments() {
		targetArgs = decoder.DumpArgumentsTStruct(args)
	}
	rs, ctx, err := telemetry.StartRPCSegment(ctx, &telemetry.RPCSegment{
		System:          telemetry.RPCTZone,
		TargetService:   c.targetName,
		TargetMethod:    fmt.Sprintf("%s_%s", c.serviceName, method),
		TargetArguments: targetArgs,
	})
	if err != nil {
		return err
	}

	entry, err := c.discovery.Select(ctx, c.targetName, diplomat.Roundrobin, false)
	if err != nil {
		return err
	}
	defer func() {
		if _, ok := err.(thrift.TTransportException); ok {
			c.discovery.Discard(entry)
		}
	}()

	transport, err := thrift.NewTHttpClientWithOptions(fmt.Sprintf("http://%s", entry.Address()), thrift.THttpClientOptions{
		Client: c.httpClient,
	})
	if err != nil {
		return err
	}
	defer func() { _ = transport.Close() }()

	tclient := transport.(*thrift.THttpClient)

	headers := http.Header{
		"X-ZONE-API":           []string{fmt.Sprintf("%s.%s", c.serviceName, method)},
		"X-ZONE-ORIGIN":        []string{zae.Service()},
		"X-ZONE-ORIGIN-APP":    []string{zae.App()},
		"X-ZONE-ORIGIN-UNIT":   []string{zae.Service()},
		"X-ZONE-ORIGIN-TOKEN":  []string{zae.AppToken()},
		"X-ZONE-ORIGIN-REGION": []string{zae.Region()},
	}
	for key, values := range rs.OutboundHeaders() {
		for _, value := range values {
			headers.Add(key, value)
		}
	}
	for key, values := range headers {
		for _, value := range values {
			tclient.SetHeader(key, value)
		}
	}

	inputProtocol := thrift.NewTBinaryProtocolTransport(thrift.NewTBufferedTransport(transport, 4*1024))
	outputProtocol := thrift.NewTMultiplexedProtocol(inputProtocol, c.serviceName)

	sclient := thrift.NewTStandardClient(inputProtocol, outputProtocol)

	err = sclient.Call(ctx, method, args, result)
	if err != nil {
		rs.End(ctx, errors.WrapError(err))
		return err
	}

	rs.End(ctx, nil)
	return nil
}
