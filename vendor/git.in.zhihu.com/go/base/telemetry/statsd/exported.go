package statsd

import (
	"time"
)

var std Client

func init() {
	var err error
	std, err = New("")
	if err != nil {
		panic(err)
	}
}

func Gauge(name string, value float64) {
	std.Gauge(name, value)
}

func Increment(name string) {
	std.Increment(name)
}

func Count(name string, value int64) {
	std.Count(name, value)
}

func Timing(name string, value time.Duration) {
	std.Timing(name, value)
}

func TimeInMilliseconds(name string, value float64) {
	std.TimeInMilliseconds(name, value)
}
