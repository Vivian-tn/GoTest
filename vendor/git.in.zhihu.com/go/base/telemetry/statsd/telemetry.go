package statsd

import (
	"database/sql"
	"time"

	"git.in.zhihu.com/go/base/zae"
	"github.com/bmhatfield/go-runtime-metrics/collector"
)

var telemetry Client

func init() {
	var err error
	telemetry, err = New("telemetry")
	if err != nil {
		panic(err)
	}
}

func init() {
	if !zae.IsDevelopEnv() && !zae.IsCIEnv() {
		go runtimeReporter()
	}
}

func runtimeReporter() {
	c := collector.New(func(key string, value uint64) {
		// gauge 指标不能以 count 结尾
		if key == "mem.gc.count" {
			key = "mem.gc.cnt"
		}
		TelemetryGauge(Join("runtime", key), float64(value))
	})
	c.Run()
}

func TelemetryGauge(name string, value float64) {
	telemetry.Gauge(Join(zae.Service(), zae.Hostname(), "go", name), value)
}

func TelemetryCount(name string, value int64) {
	telemetry.Count(Join(zae.Service(), zae.Hostname(), "go", name), value)
}

func TelemetryLog(level string) {
	TelemetryCount(Join("log", level, "count"), 1)
}

func TelemetryDBStats(name string, stats sql.DBStats) {
	TelemetryGauge(Join("db_stats", name, "max_open_connections"), float64(stats.MaxOpenConnections))
	TelemetryGauge(Join("db_stats", name, "open_connections"), float64(stats.OpenConnections))
	TelemetryGauge(Join("db_stats", name, "in_use"), float64(stats.InUse))
	TelemetryGauge(Join("db_stats", name, "idle"), float64(stats.Idle))
	TelemetryGauge(Join("db_stats", name, "wait_count"), float64(stats.WaitCount))
	TelemetryGauge(Join("db_stats", name, "wait_duration"), float64(stats.WaitDuration/time.Millisecond))
	TelemetryGauge(Join("db_stats", name, "max_idle_closed"), float64(stats.MaxIdleClosed))
	TelemetryGauge(Join("db_stats", name, "max_life_time_closed"), float64(stats.MaxLifetimeClosed))
}
