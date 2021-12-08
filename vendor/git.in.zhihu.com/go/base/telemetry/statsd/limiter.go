package statsd

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/juju/ratelimit"
)

func newBucketLimiter(rate int64) *bucketLimiter {
	bucket := ratelimit.NewBucket(time.Second/time.Duration(rate), rate)
	// clean bucket
	bucket.Take(rate)
	return &bucketLimiter{bucket: bucket}
}

type bucketLimiter struct {
	bucket *ratelimit.Bucket
	active atomic.Value
}

func (t *bucketLimiter) Available() bool {
	t.active.Store(time.Now())
	return t.bucket.TakeAvailable(1) > 0
}

func (t *bucketLimiter) ActiveAt() time.Time {
	return t.active.Load().(time.Time)
}

func newLimiter(rate int64, gcDuration time.Duration) *limiter {
	l := &limiter{
		rate:    rate,
		ticker:  time.NewTicker(gcDuration),
		closing: make(chan struct{}),
	}
	go l.gcLoop()
	return l
}

type limiter struct {
	buckets sync.Map
	rate    int64
	ticker  *time.Ticker
	closing chan struct{}
}

func (l *limiter) Available(bucket string) bool {
	var bl *bucketLimiter
	value, ok := l.buckets.Load(bucket)
	if !ok {
		bl = newBucketLimiter(l.rate)
		l.buckets.Store(bucket, bl)
	} else {
		bl = value.(*bucketLimiter)
	}
	return bl.Available()
}

func (l *limiter) gcLoop() {
	for {
		select {
		case <-l.ticker.C:
			now := time.Now()
			l.buckets.Range(func(key, value interface{}) bool {
				bucket := key.(string)
				bl := value.(*bucketLimiter)
				if now.Sub(bl.ActiveAt()) > time.Minute {
					l.buckets.Delete(bucket)
				}
				return true
			})
		case <-l.closing:
			return
		}
	}
}

func (l *limiter) Close() {
	l.ticker.Stop()
	close(l.closing)
}
