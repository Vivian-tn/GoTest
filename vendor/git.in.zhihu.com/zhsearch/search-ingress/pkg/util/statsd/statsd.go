package statsd

import (
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"git.in.zhihu.com/zhsearch/search-ingress/pkg/log"
	"github.com/juju/ratelimit"
	"gopkg.in/alexcesaro/statsd.v2"
)

func newCountAggreMetric(count int64) *aggreMetric {
	aMetric := &aggreMetric{}
	aMetric.AddCount(count)
	return aMetric
}

func newGaugeAggreMetric(value interface{}) *aggreMetric {
	aMetric := &aggreMetric{}
	aMetric.SetGauge(value)
	return aMetric
}

type aggreMetric struct {
	value atomic.Value
	count int64
}

func (a *aggreMetric) AddCount(count int64) {
	atomic.AddInt64(&a.count, count)
}

func (a *aggreMetric) Count() int64 {
	return atomic.LoadInt64(&a.count)
}

func (a *aggreMetric) SetGauge(value interface{}) *aggreMetric {
	a.value.Store(value)
	return a
}

func (a *aggreMetric) Gauge() interface{} {
	return a.value.Load()
}

type Client struct {
	cli       *statsd.Client
	buffer    sync.Map
	flusher   *time.Ticker
	limiter   *ratelimit.Bucket
	closing   chan struct{}
	debugMode bool
}

type Timing statsd.Timing

func (t Timing) Duration() time.Duration {
	return (statsd.Timing)(t).Duration()
}

// Deprecated: use Instance instead.
var New = Instance

var (
	instance       *Client
	instanceLocker sync.Mutex
)

// Instance get the global statsd Client instance.
// Client's create should be lazy-loading.
func Instance() *Client {
	instanceLocker.Lock()
	if instance == nil {
		instance = NewWithOptions(statsd.Address("status:8126"), statsd.MaxPacketSize(63000))
	}
	instanceLocker.Unlock()
	return instance
}

const (
	flushDuration       = 100 * time.Millisecond
	dropWarningDuration = 5 * time.Second
	timingLimit         = 5000
)

func NewWithOptions(opts ...statsd.Option) *Client {
	cli, err := statsd.New(opts...)
	if err != nil {
		log.Error(err)
	}

	client := &Client{
		cli:     cli,
		buffer:  sync.Map{},
		flusher: time.NewTicker(flushDuration),
		// 200 微秒生产一个令牌，即 QPS 5k
		limiter:   ratelimit.NewBucket(time.Second/timingLimit, timingLimit),
		debugMode: os.Getenv("DEBUG") != "",
		closing:   make(chan struct{}),
	}

	go client.flushLoop()

	return client
}

func (c *Client) Increment(bucket string) {
	c.Count(bucket, 1)
}

func (c *Client) Count(bucket string, count int) {
	// Carbon 默认的 aggregationMethod 是 avg，知乎内部只有给 \.count$ 的指标才配置成 sum，所以这里 count 操作统一添加后缀。
	// 参考：http://wiki.in.zhihu.com/pages/viewpage.action?pageId=6195516
	if !strings.HasSuffix(bucket, ".count") {
		bucket += ".count"
	}

	c.RecordCount(bucket, count)
}

func (c *Client) RecordCount(bucket string, count int) {
	if c.debugMode {
		c.debug(bucket, count)
	}

	c.aggreCount(bucket, int64(count))
}

func (c *Client) Gauge(bucket string, value interface{}) {
	if c.debugMode {
		c.debug(bucket, value)
	}

	c.aggreGauge(bucket, value)
}

func (c *Client) NewTiming() Timing {
	return (Timing)(c.cli.NewTiming())
}

func (c *Client) Timing(bucket string, timing Timing) {
	// millisecond 使用 float64 保证精度不受损
	millisecond := timing.Duration().Seconds() * 1000
	if c.debugMode {
		c.debug(bucket, millisecond)
	}
	c.limitTiming(bucket, millisecond)
}

func (c *Client) RecordTime(bucket string, millisecond int64) {
	// 跟 Timing 方法数据类型不一样，保持向前兼容
	if c.debugMode {
		c.debug(bucket, millisecond)
	}
	c.limitTiming(bucket, millisecond)
}

func (c *Client) limitTiming(bucket string, value interface{}) {
	if c.limiter.TakeAvailable(1) > 0 {
		c.cli.Timing(bucket, value)
	}
}

func (c *Client) aggreGauge(bucket string, value interface{}) {
	metric, ok := c.buffer.Load(bucket)
	if !ok {
		c.buffer.Store(bucket, newGaugeAggreMetric(value))
	} else {
		metric.(*aggreMetric).SetGauge(value)
	}
}

func (c *Client) aggreCount(bucket string, count int64) {
	aMetric, ok := c.buffer.Load(bucket)
	if !ok {
		c.buffer.Store(bucket, newCountAggreMetric(count))
	} else {
		aMetric.(*aggreMetric).AddCount(count)
	}
}

func (c *Client) flush() (dropped int) {
	c.buffer.Range(func(key, value interface{}) bool {
		bucket := key.(string)
		if !ValidBucket(bucket) {
			dropped++
			return true
		}

		aMetric := value.(*aggreMetric)
		// count 跟 gauge 的 bucket 可能是同一个
		if count := aMetric.Count(); count > 0 {
			c.cli.Count(bucket, count)
		}
		if gauge := aMetric.Gauge(); gauge != nil {
			c.cli.Gauge(bucket, gauge)
		}
		c.buffer.Delete(bucket)
		return true
	})
	return dropped
}

func (c *Client) flushLoop() {
	var flushTimes uint64 = 0
	var totalDropped = 0
	for {
		select {
		case <-c.flusher.C:
			dropped := c.flush()
			// print warning message every few seconds.
			totalDropped += dropped
			flushTimes++
			if flushTimes%(uint64(dropWarningDuration/flushDuration)) == 0 {
				if totalDropped > 0 {
					log.Warnf("[STATSD] dropped %d invalid metrics", totalDropped)
				}
				totalDropped = 0
			}
		case <-c.closing:
			return
		}
	}
}

func (c *Client) debug(bucket string, value interface{}) {
	log.Debugf("[STATSD] bucket:%s,value:%v\n", bucket, value)
}

func (c *Client) Close() {
	c.flush()
	close(c.closing)
	c.flusher.Stop()
	c.cli.Close()
}
