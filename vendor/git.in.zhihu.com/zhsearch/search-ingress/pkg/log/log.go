package log

import (
	"context"

	zlog "git.in.zhihu.com/go/base/telemetry/log"
	"git.in.zhihu.com/go/logrus"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/getsentry/raven-go"
)

type Fields = logrus.Fields

func SetDebug() {
	zlog.SetLevel(zlog.DebugLevel)
}

// 已知的正常错误不用打到 sentry 上
func isExceptionUnknown(err error) bool {
	switch e := err.(type) {
	case thrift.TTransportException:
		tid := e.TypeId()
		if tid == thrift.TIMED_OUT {
			return false
		}
	}
	return true
}

// WithError creates an entry from the standard logger and adds an error to it, using the value defined in ErrorKey as key.
// A error message to sentry is fired.
func WithError(err error) *logrus.Entry {
	if isExceptionUnknown(err) {
		raven.CaptureError(err, nil)
	}
	return zlog.WithFields(context.Background(), zlog.Fields{logrus.ErrorKey: err})
}

// WithField creates an entry from the standard logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// on the Entry it returns.
func WithField(key string, value interface{}) *logrus.Entry {
	return zlog.WithFields(context.Background(), zlog.Fields{key: value})
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// on the Entry it returns.
func WithFields(fields logrus.Fields) *logrus.Entry {
	return zlog.WithFields(context.Background(), fields)
}

// Debug logs a message at level Debug on the standard zlog.
func Debug(args ...interface{}) {
	zlog.Debug(context.Background(), args...)
}

// Info logs a message at level Info on the standard zlog.
func Info(args ...interface{}) {
	zlog.Info(context.Background(), args...)
}

// Warn logs a message at level Warn on the standard zlog.
func Warn(args ...interface{}) {
	zlog.Warn(context.Background(), args...)
}

// Warning logs a message at level Warn on the standard zlog.
func Warning(args ...interface{}) {
	zlog.Warn(context.Background(), args...)
}

// Error logs a message at level Error on the standard zlog.
func Error(args ...interface{}) {
	zlog.Error(context.Background(), args...)
}

// Fatal logs a message at level Fatal on the standard zlog.
func Fatal(args ...interface{}) {
	zlog.Fatal(context.Background(), args...)
}

// Debugf logs a message at level Debug on the standard zlog.
func Debugf(format string, args ...interface{}) {
	zlog.Debugf(context.Background(), format, args...)
}

// Infof logs a message at level Info on the standard zlog.
func Infof(format string, args ...interface{}) {
	zlog.Infof(context.Background(), format, args...)
}

// Warnf logs a message at level Warn on the standard zlog.
func Warnf(format string, args ...interface{}) {
	zlog.Warnf(context.Background(), format, args...)
}

// Warningf logs a message at level Warn on the standard zlog.
func Warningf(format string, args ...interface{}) {
	zlog.Warnf(context.Background(), format, args...)
}

// Errorf logs a message at level Error on the standard zlog.
func Errorf(format string, args ...interface{}) {
	zlog.Errorf(context.Background(), format, args...)
}

// Fatalf logs a message at level Fatal on the standard zlog.
func Fatalf(format string, args ...interface{}) {
	zlog.Fatalf(context.Background(), format, args...)
}

// Debugln logs a message at level Debug on the standard zlog.
func Debugln(args ...interface{}) {
	zlog.Debugln(context.Background(), args...)
}

// Infoln logs a message at level Info on the standard zlog.
func Infoln(args ...interface{}) {
	zlog.Infoln(context.Background(), args...)
}

// Warnln logs a message at level Warn on the standard zlog.
func Warnln(args ...interface{}) {
	zlog.Warnln(context.Background(), args...)
}

// Warningln logs a message at level Warn on the standard zlog.
func Warningln(args ...interface{}) {
	zlog.Warnln(context.Background(), args...)
}

// Errorln logs a message at level Error on the standard zlog.
func Errorln(args ...interface{}) {
	zlog.Errorln(context.Background(), args...)
}

// Fatalln logs a message at level Fatal on the standard zlog.
func Fatalln(args ...interface{}) {
	zlog.Fatalln(context.Background(), args...)
}

func WithCtx(ctx context.Context) *logrus.Entry {
	return zlog.WithFields(ctx, Fields{})
}
