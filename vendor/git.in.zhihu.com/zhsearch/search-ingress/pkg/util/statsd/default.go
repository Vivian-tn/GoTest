package statsd

var DefaultClient = New()

func NewTiming() Timing {
	return DefaultClient.NewTiming()
}

func RecordTime(bucket string, millisecond int64) {
	DefaultClient.RecordTime(bucket, millisecond)
}

func Time(bucket string, timing Timing) {
	DefaultClient.Timing(bucket, timing)
}

func RecordCount(bucket string, count int) {
	DefaultClient.RecordCount(bucket, count)
}

func Count(bucket string, count int) {
	DefaultClient.Count(bucket, count)
}

func Gauge(bucket string, value interface{}) {
	DefaultClient.Gauge(bucket, value)
}

func Increment(bucket string) {
	DefaultClient.Increment(bucket)
}

func Close() {
	// default client 不支持关闭
}
