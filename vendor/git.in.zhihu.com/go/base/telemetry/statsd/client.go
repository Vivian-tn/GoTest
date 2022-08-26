package statsd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/DataDog/datadog-go/statsd"

	"git.in.zhihu.com/go/base/zae"
)

const (
	MetadataApp = "app"
)

type Client interface {
	Gauge(name string, value float64)
	Increment(name string)
	Count(name string, value int64)
	Timing(name string, value time.Duration)
	TimeInMilliseconds(name string, value float64)
	Close() error
}

func New(prefix string, addrs ...string) (Client, error) {
	return NewWithOptions(prefix, WithAddrs(addrs...))
}

func NewWithOptions(prefix string, opts ...option) (Client, error) {
	o := &options{
		metadata: true,
	}
	for _, opt := range opts {
		opt(o)
	}

	var addr string
	switch len(o.addrs) {
	case 0:
		if defaultAddr, err := zae.DiscoveryOne(zae.ResourceStatsd, "default", ""); err == nil {
			addr = defaultAddr.Host
		} else {
			addr = "status:8126"
		}
	case 1:
		addr = o.addrs[0]
	default:
		return nil, fmt.Errorf("invalid addr %v", o.addrs)
	}
	if prefix != "" && !Verify(prefix) {
		return nil, fmt.Errorf("invalid prefix %v", prefix)
	}

	if o.metadata {
		if o.globalTags == nil {
			o.globalTags = make(map[string]string)
		}
		o.globalTags[MetadataApp] = zae.App()
	}

	tags := make([]string, 0, len(o.globalTags))
	for k, v := range o.globalTags {
		if strings.ContainsAny(k, ":,") {
			return nil, fmt.Errorf("invalid tag key %v", k)
		}
		if strings.ContainsAny(v, ":,") {
			return nil, fmt.Errorf("invalid tag value %v", v)
		}
		tags = append(tags, k+":"+v)
	}

	var client *statsd.Client
	var err error

	sopts := []statsd.Option{
		statsd.WithClientSideAggregation(),
		statsd.WithoutTelemetry(),
		statsd.WithTags(tags),
	}
	if len(prefix) > 0 {
		sopts = append(sopts, statsd.WithNamespace(prefix))
	}

	client, err = statsd.New(addr, sopts...)
	if err != nil {
		return nil, err
	}

	name := Node(prefix)
	if name == "" {
		name = "default"
	}

	wrap := &wrapClient{
		name:    name,
		client:  client,
		metrics: new(Metrics),
		// 频控 5K QPS，与原 go/box 实现保持一致
		limiter: newLimiter(5000, time.Minute),
		closing: make(chan struct{}),
	}

	go wrap.reporter()

	return wrap, nil
}

type Metrics struct {
	Total   uint64
	Dropped uint64
}

type wrapClient struct {
	name    string
	client  *statsd.Client
	limiter *limiter
	metrics *Metrics
	closing chan struct{}
}

func (c *wrapClient) reporter() {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			TelemetryCount(Join("statsd_stats", c.name, "total"), int64(atomic.SwapUint64(&c.metrics.Total, 0)))
			TelemetryCount(Join("statsd_stats", c.name, "dropped"), int64(atomic.SwapUint64(&c.metrics.Dropped, 0)))
		case <-c.closing:
			return
		}
	}
}

func (c *wrapClient) Gauge(name string, value float64) {
	atomic.AddUint64(&c.metrics.Total, 1)

	var err error
	defer func() {
		if err != nil {
			atomic.AddUint64(&c.metrics.Dropped, 1)
			logError(err)
		}
	}()

	if !Verify(name) {
		err = fmt.Errorf("invalid metric name: %s", name)
		return
	}
	if strings.HasSuffix(name, ".count") {
		err = fmt.Errorf("invalid gauge name with count suffix: %s", name)
		return
	}

	if zae.IsDebug() {
		log.Printf("[DEBUG:Statsd] %s %v\n", name, value)
	}

	err = c.client.Gauge(name, value, nil, 1)
}

func (c *wrapClient) Increment(name string) {
	c.Count(name, 1)
}

func (c *wrapClient) Count(name string, value int64) {
	atomic.AddUint64(&c.metrics.Total, 1)

	var err error
	defer func() {
		if err != nil {
			atomic.AddUint64(&c.metrics.Dropped, 1)
			logError(err)
		}
	}()

	if !Verify(name) {
		err = fmt.Errorf("invalid metric name: %s", name)
		return
	}
	if !strings.HasSuffix(name, ".count") {
		name += ".count"
	}

	if zae.IsDebug() {
		log.Printf("[DEBUG:Statsd] %s %v\n", name, value)
	}

	err = c.client.Count(name, value, nil, 1)
}

func (c *wrapClient) Timing(name string, value time.Duration) {
	c.TimeInMilliseconds(name, value.Seconds()*1000)
}

func (c *wrapClient) TimeInMilliseconds(name string, value float64) {
	atomic.AddUint64(&c.metrics.Total, 1)

	var err error
	defer func() {
		if err != nil {
			atomic.AddUint64(&c.metrics.Dropped, 1)
			logError(err)
		}
	}()

	if !Verify(name) {
		err = fmt.Errorf("invalid metric name: %s", name)
		return
	}
	if strings.HasSuffix(name, ".count") {
		err = fmt.Errorf("invalid timing name with count suffix: %s", name)
		return
	}
	if !c.limiter.Available(name) {
		return
	}

	if zae.IsDebug() {
		log.Printf("[DEBUG:Statsd] %s %v\n", name, value)
	}

	err = c.client.TimeInMilliseconds(name, value, nil, 1)
}

func (c *wrapClient) Close() error {
	close(c.closing)
	c.limiter.Close()
	return c.client.Close()
}

func logError(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "[E %s - - %s:0] %s\n", time.Now().Format("2006-01-02 15:04:05.000"), zae.Hostname(), err.Error())
}
