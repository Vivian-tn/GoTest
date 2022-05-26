package statsd

import (
	"fmt"
	"os"
	"time"

	"git.in.zhihu.com/go/base/telemetry/statsd"
	"git.in.zhihu.com/go/base/zae"
	"git.in.zhihu.com/go/logrus"
)

func NewHook() *Hook {
	return &Hook{}
}

type Hook struct{}

func (h *Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *Hook) Fire(entry *logrus.Entry) error {
	defer func() {
		if rval := recover(); rval != nil {
			logMessage(fmt.Sprint(rval))
		}
	}()

	switch entry.Level {
	case logrus.FatalLevel, logrus.PanicLevel:
		statsd.TelemetryLog(logrus.FatalLevel.String())
	case logrus.ErrorLevel:
		statsd.TelemetryLog(logrus.ErrorLevel.String())
	case logrus.WarnLevel:
		statsd.TelemetryLog(logrus.WarnLevel.String())
	case logrus.InfoLevel:
		statsd.TelemetryLog(logrus.InfoLevel.String())
	case logrus.DebugLevel, logrus.TraceLevel:
		statsd.TelemetryLog(logrus.DebugLevel.String())
	}
	return nil
}

func logMessage(message string) {
	_, _ = fmt.Fprintf(os.Stderr, "[E %s - - %s:0] %s\n", time.Now().Format("2006-01-02 15:04:05.000"), zae.Hostname(), message)
}
