package log

import (
	"context"
	"io"

	"git.in.zhihu.com/go/base/telemetry/log/internal/hooks/sentry"
	"git.in.zhihu.com/go/logrus"
)

// fork from https://github.com/sirupsen/logrus/blob/v1.6.0/exported.go

var std *ZLogger

type Level = logrus.Level
type Fields = logrus.Fields
type Formatter = logrus.Formatter

const (
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel = logrus.FatalLevel
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel = logrus.ErrorLevel
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel = logrus.WarnLevel
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel = logrus.InfoLevel
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel = logrus.DebugLevel
)

func SetIncludePaths(p []string) {
	sentry.SetIncludePaths(p)
}

func StandardLogger() *ZLogger {
	return std
}

// SetLevel sets the standard logger level
func SetLevel(level Level) {
	std.SetLevel(level)
}

// SetFilelineNumberLevel sets the standard logger to include the calling depending on the entry's level
// It will be ignored is ReportCaller is true
func SetFileLineNumberLevel(level Level) {
	std.SetFileLineNumberLevel(level)
}

// SetSentryLevel sets the minimum level to report to sentry for standard logger
func SetSentryLevel(level Level) {
	std.SetSentryLevel(level)
}

// SetReportCaller sets whether the standard logger will include the calling
// method as a field.
func SetReportCaller(include bool) {
	std.SetReportCaller(include)
}

// SetOutput sets one output for standard logger
func SetOutput(output io.Writer) {
	std.SetOutput(output)
}

// SetOutputs sets multiple outputs for standard logger
func SetOutputs(outputs ...io.Writer) {
	std.SetOutputs(outputs...)
}

// SetFormatter sets formatter for standard logger
func SetFormatter(formatter logrus.Formatter) {
	std.SetFormatter(formatter)
}

// GetLevel returns the standard logger level.
func GetLevel() Level {
	return std.GetLevel()
}

func GetFileLineNumberLevel() Level {
	return std.GetFileLineNumberLevel()
}

func GetSentryLevel() Level {
	return std.GetSentryLevel()
}

// IsLevelEnabled checks if the log level of the standard logger is greater than the level param
func IsLevelEnabled(level Level) bool {
	return std.IsLevelEnabled(level)
}

// WithFields creates an entry from the standard logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
func WithFields(ctx context.Context, fields Fields) *logrus.Entry {
	return std.WithFields(ctx, fields)
}

func WithField(ctx context.Context, key string, value interface{}) *logrus.Entry {
	return std.WithField(ctx, key, value)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(ctx context.Context, args ...interface{}) {
	std.Debug(ctx, args...)
}

// Info logs a message at level Info on the standard logger.
func Info(ctx context.Context, args ...interface{}) {
	std.Info(ctx, args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(ctx context.Context, args ...interface{}) {
	std.Warn(ctx, args...)
}

// Error logs a message at level Error on the standard logger.
func Error(ctx context.Context, args ...interface{}) {
	std.Error(ctx, args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(ctx context.Context, args ...interface{}) {
	std.Fatal(ctx, args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(ctx context.Context, format string, args ...interface{}) {
	std.Debugf(ctx, format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(ctx context.Context, format string, args ...interface{}) {
	std.Infof(ctx, format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(ctx context.Context, format string, args ...interface{}) {
	std.Warnf(ctx, format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(ctx context.Context, format string, args ...interface{}) {
	std.Errorf(ctx, format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalf(ctx context.Context, format string, args ...interface{}) {
	std.Fatalf(ctx, format, args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(ctx context.Context, args ...interface{}) {
	std.Debugln(ctx, args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(ctx context.Context, args ...interface{}) {
	std.Infoln(ctx, args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(ctx context.Context, args ...interface{}) {
	std.Warnln(ctx, args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(ctx context.Context, args ...interface{}) {
	std.Errorln(ctx, args...)
}

// Fatalln logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalln(ctx context.Context, args ...interface{}) {
	std.Fatalln(ctx, args...)
}

type FieldName = string

const (
	LoggerNameField FieldName = "logger"
)

func (zl *ZLogger) defaultFields() Fields {
	return Fields{
		LoggerNameField: zl.name,
	}
}

func (zl *ZLogger) GetName() string {
	return zl.name
}

// SetLevel sets the logger level
func (zl *ZLogger) SetLevel(level Level) {
	zl.Logger.SetLevel(level)
}

// SetFilelineNumberLevel sets the logger to report caller when entry's level is greater than FilelineNumberLevel，it will be ignored when ReportCaller is true
func (zl *ZLogger) SetFileLineNumberLevel(level Level) {
	zl.Logger.SetFilelineNumberLevel(level)
}

// SetSentryLevel sets the minimum level to report to sentry for the logger
func (zl *ZLogger) SetSentryLevel(level Level) {
	zl.mu.Lock()
	defer zl.mu.Unlock()
	zl.sentryLevel = level
	reloadHooks(zl)
}

// SetReportCaller sets the logger to report caller while logging
func (zl *ZLogger) SetReportCaller(include bool) {
	zl.Logger.SetReportCaller(include)
}

// SetOutput sets one output for the logger
func (zl *ZLogger) SetOutput(output io.Writer) {
	zl.Logger.SetOutput(output)
}

// SetOutputs sets multiple outputs for the logger
func (zl *ZLogger) SetOutputs(outputs ...io.Writer) {
	zl.Logger.SetOutputs(outputs...)
}

// SetFormatter sets formatter for the logger
func (zl *ZLogger) SetFormatter(formatter logrus.Formatter) {
	zl.Logger.SetFormatter(formatter)
}

func (zl *ZLogger) GetLevel() Level {
	return zl.Logger.GetLevel()
}

func (zl *ZLogger) GetFileLineNumberLevel() Level {
	return zl.Logger.FilelineNumberLevel
}

func (zl *ZLogger) GetSentryLevel() Level {
	return zl.sentryLevel
}

func (zl *ZLogger) IsLevelEnabled(level Level) bool {
	return zl.Logger.IsLevelEnabled(level)
}

func (zl *ZLogger) WithFields(ctx context.Context, fields Fields) *logrus.Entry {
	targets := zl.defaultFields()
	for key, value := range fields {
		targets[key] = value
	}
	return zl.Logger.WithContext(ctx).WithFields(targets)
}

func (zl *ZLogger) WithField(ctx context.Context, key string, value interface{}) *logrus.Entry {
	return zl.WithFields(ctx, Fields{key: value})
}

// Debug logs a message at level Debug on the zlogger.
func (zl *ZLogger) Debug(ctx context.Context, args ...interface{}) {
	zl.Logger.WithFields(zl.defaultFields()).WithContext(ctx).Debug(args...)
}

// Info logs a message at level Info on the zlogger.
func (zl *ZLogger) Info(ctx context.Context, args ...interface{}) {
	zl.Logger.WithFields(zl.defaultFields()).WithContext(ctx).Info(args...)
}

// Warn logs a message at level Warn on the zlogger.
func (zl *ZLogger) Warn(ctx context.Context, args ...interface{}) {
	zl.Logger.WithFields(zl.defaultFields()).WithContext(ctx).Warn(args...)
}

// Error logs a message at level Error on the zlogger.
func (zl *ZLogger) Error(ctx context.Context, args ...interface{}) {
	zl.Logger.WithFields(zl.defaultFields()).WithContext(ctx).Error(args...)
}

// Fatal logs a message at level Fatal on the zlogger then the process will exit with status set to 1.
func (zl *ZLogger) Fatal(ctx context.Context, args ...interface{}) {
	zl.Logger.WithFields(zl.defaultFields()).WithContext(ctx).Fatal(args...)
}

// Debugf logs a message at level Debug on the zlogger.
func (zl *ZLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	zl.Logger.WithFields(zl.defaultFields()).WithContext(ctx).Debugf(format, args...)
}

// Infof logs a message at level Info on the zlogger.
func (zl *ZLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	zl.Logger.WithFields(zl.defaultFields()).WithContext(ctx).Infof(format, args...)
}

// Warnf logs a message at level Warn on the zlogger.
func (zl *ZLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	zl.Logger.WithFields(zl.defaultFields()).WithContext(ctx).Warnf(format, args...)
}

// Errorf logs a message at level Error on the zlogger.
func (zl *ZLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	zl.Logger.WithFields(zl.defaultFields()).WithContext(ctx).Errorf(format, args...)
}

// Fatalf logs a message at level Fatal on the zlogger then the process will exit with status set to 1.
func (zl *ZLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	zl.Logger.WithFields(zl.defaultFields()).WithContext(ctx).Fatalf(format, args...)
}

// Debugln logs a message at level Debug on the zlogger.
func (zl *ZLogger) Debugln(ctx context.Context, args ...interface{}) {
	zl.Logger.WithFields(zl.defaultFields()).WithContext(ctx).Debugln(args...)
}

// Infoln logs a message at level Info on the zlogger.
func (zl *ZLogger) Infoln(ctx context.Context, args ...interface{}) {
	zl.Logger.WithFields(zl.defaultFields()).WithContext(ctx).Infoln(args...)
}

// Warnln logs a message at level Warn on the zlogger.
func (zl *ZLogger) Warnln(ctx context.Context, args ...interface{}) {
	zl.Logger.WithFields(zl.defaultFields()).WithContext(ctx).Warnln(args...)
}

// Errorln logs a message at level Error on the zlogger.
func (zl *ZLogger) Errorln(ctx context.Context, args ...interface{}) {
	zl.Logger.WithFields(zl.defaultFields()).WithContext(ctx).Errorln(args...)
}

// Fatalln logs a message at level Fatal on the zlogger then the process will exit with status set to 1.
func (zl *ZLogger) Fatalln(ctx context.Context, args ...interface{}) {
	zl.Logger.WithFields(zl.defaultFields()).WithContext(ctx).Fatalln(args...)
}
