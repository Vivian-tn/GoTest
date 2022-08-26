package log

import (
	"io"
	"os"
	"sync"

	"git.in.zhihu.com/go/base/telemetry/log/internal/hooks/sentry"
	"git.in.zhihu.com/go/base/telemetry/log/internal/hooks/statsd"
	"git.in.zhihu.com/go/base/zae"
	"git.in.zhihu.com/go/logrus"
)

var (
	manager logManager
	managerOnce = sync.Once{}
)

func getManager() *logManager {
	managerOnce.Do(func() {
		manager = logManager{
			loggerStore: sync.Map{},
		}
	})
	return &manager
}

func init() {
	std = newDefaultLogger()
	addHooks(std)
	// add default logger to manager
	getManager().add(std.GetName(), std)
}

type ZLogger struct {
	Logger      *logrus.Logger
	mu          sync.Mutex
	name        string
	sentryLevel Level
}

const empty = ""

var (
	defaultName        = "default" // standard logger name
	defaultOut         = os.Stderr
	defaultFormatter   = NewDefaultDynamicFormatter()
	defaultLevel       = InfoLevel
	defaultSentryLevel = ErrorLevel
)

func New(name string) *ZLogger {
	logger := GetLogger(name)
	if logger != nil {
		return logger
	}

	return newLogger(name)
}

func GetLogger(name string) *ZLogger {
	// use standard logger if there's no options
	if name == empty || name == defaultName {
		return StandardLogger()
	}

	if logger := getManager().getLogger(name); logger != nil {
		return logger
	}

	return nil
}

func newLogger(name string) *ZLogger {
	logger := newDefaultLogger()
	logger.name = name
	addHooks(logger)

	// add logger to the manager
	getManager().add(name, logger)
	return logger
}

func newDefaultLogger() *ZLogger {
	return &ZLogger{
		Logger: &logrus.Logger{
			Outs:         []io.Writer{defaultOut},
			Formatter:    defaultFormatter,
			Hooks:        make(logrus.LevelHooks, 0),
			ReportCaller: false,
			Level:        defaultLevel,
		},
		name:        defaultName,
		sentryLevel: defaultSentryLevel,
	}
}

// GetLoggers return the loggers owned by the log manager
func GetLoggers() []*ZLogger {
	var rets []*ZLogger
	getManager().loggerStore.Range(func(k, v interface{}) bool {
		rets = append(rets, v.(*ZLogger))
		return true
	})
	return rets
}

func addHooks(logger *ZLogger) {
	dsn, release := zae.Sentry()
	hook, err := sentry.NewHook(dsn, release, logger.sentryLevel)
	if err != nil {
		logger.Logger.WithError(err).Error("add sentry hook failed")
	} else {
		logger.Logger.AddHook(hook)
	}

	logger.Logger.AddHook(statsd.NewHook())
}

func reloadHooks(logger *ZLogger) {
	hooks := make(logrus.LevelHooks, 0)

	dsn, release := zae.Sentry()
	hook, err := sentry.NewHook(dsn, release, logger.sentryLevel)
	if err != nil {
		logger.Logger.WithError(err).Error("reload sentry hook failed")
	} else {
		hooks.Add(hook)
	}

	hooks.Add(statsd.NewHook())

	logger.Logger.ReplaceHooks(hooks)
}
