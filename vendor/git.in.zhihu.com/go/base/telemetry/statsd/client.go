package statsd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"git.in.zhihu.com/go/base/zae"
	"github.com/DataDog/datadog-go/statsd"
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
	var addr string
	switch len(addrs) {
	case 0:
		if defaultAddr, err := zae.DiscoveryOne(zae.ResourceStatsd, "default", ""); err == nil {
			addr = defaultAddr.Host
		} else {
			addr = "status:8126"
		}
	case 1:
		addr = addrs[0]
	default:
		return nil, fmt.Errorf("invalid addr %v", addrs)
	}
	if prefix != "" && !Verify(prefix) {
		return nil, fmt.Errorf("invalid prefix %v", prefix)
	}

	var client *statsd.Client
	var err error

	opts := []statsd.Option{
		statsd.WithClientSideAggregation(),
		statsd.WithoutTelemetry(),
	}
	if len(prefix) > 0 {
		opts = append(opts, statsd.WithNamespace(prefix))
	}

	client, err = statsd.New(addr, opts...)
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
		limiter: newLimiter(100, time.Minute),
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
