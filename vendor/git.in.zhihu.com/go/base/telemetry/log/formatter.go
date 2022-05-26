package log

/*
DynamicFormatter:
BenchmarkInfoLogWithoutFilelineNoTest-4    	  129972	      7701 ns/op
BenchmarkInfoLogWithFilelineNoTest-4       	  111541	     10741 ns/op
BenchmarkErrorLogWithoutFilelineNoTest-4   	   40750	     28723 ns/op
BenchmarkErrorLogWithFilelineNoTest-4      	   37250	     32306 ns/op

DefaultFormatter:
BenchmarkInfoLogWithoutFilelineNoTest-4    	  227107	      5074 ns/op
BenchmarkInfoLogWithFilelineNoTest-4       	  158587	      8202 ns/op
BenchmarkErrorLogWithoutFilelineNoTest-4   	   39020	     29119 ns/op
BenchmarkErrorLogWithFilelineNoTest-4      	   32292	     34377 ns/op
*/

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"git.in.zhihu.com/go/base/zae"
	"git.in.zhihu.com/go/logrus"
)

const (
	traceIDKey = "X-B3-Traceid" // equal with telemetry
)

var (
	processID = strconv.Itoa(os.Getpid())
)

// 日志格式和采样规则 http://wiki.in.zhihu.com/pages/viewpage.action?pageId=69050370

var (
	severityMap = map[logrus.Level]byte{
		logrus.TraceLevel: 'D',
		logrus.DebugLevel: 'D',
		logrus.InfoLevel:  'I',
		logrus.WarnLevel:  'W',
		logrus.ErrorLevel: 'E',
		logrus.FatalLevel: 'F',
		logrus.PanicLevel: 'F',
	}
)

type DefaultFormatter struct{}

func (f *DefaultFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Buffer.WriteByte('[')

	entry.Buffer.WriteByte(severityMap[entry.Level])

	entry.Buffer.WriteByte(' ')

	entry.Buffer.WriteString(entry.Time.Format("2006-01-02 15:04:05.000"))

	entry.Buffer.WriteByte(' ')

	var traceID string
	if entry.Context != nil {
		traceID, _ = entry.Context.Value(traceIDKey).(string)
	}
	if traceID != "" {
		entry.Buffer.WriteString(traceID)
	} else {
		entry.Buffer.WriteByte('-')
	}

	entry.Buffer.WriteByte(' ')

	if entry.Caller != nil {
		entry.Buffer.WriteString(filepath.Base(entry.Caller.File))
		entry.Buffer.WriteByte(':')
		entry.Buffer.WriteString(strconv.Itoa(entry.Caller.Line))
	} else {
		entry.Buffer.WriteByte('-')
	}

	entry.Buffer.WriteByte(' ')

	entry.Buffer.WriteString(zae.Hostname())
	entry.Buffer.WriteByte(':')
	entry.Buffer.WriteString(processID)

	entry.Buffer.WriteString("] ")

	if len(entry.Data) > 0 {
		data := make(Fields, len(entry.Data))
		for k, v := range entry.Data {
			switch v := v.(type) {
			case error:
				// Otherwise errors are ignored by `encoding/json`
				// https://github.com/sirupsen/logrus/issues/137
				data[k] = v.Error()
			default:
				data[k] = v
			}
		}
		if extraData, err := json.Marshal(data); err == nil {
			entry.Buffer.WriteString("[")
			entry.Buffer.Write(extraData)
			entry.Buffer.WriteString("] ")
		}
	}

	entry.Buffer.WriteString(entry.Message)

	entry.Buffer.WriteByte('\n')

	return entry.Buffer.Bytes(), nil
}

const defaultPattern = "[{Level} {Time} {TraceID} {FileLine} {Hostname}:{ProcessID}] [{Data}] {Message}"

type dynamicFormatter struct {
	pattern string
}

func NewDynamicFormatter(pattern string) *dynamicFormatter {
	if pattern == "" {
		return NewDefaultDynamicFormatter()
	}
	return &dynamicFormatter{pattern}
}

func NewDefaultDynamicFormatter() *dynamicFormatter {
	return &dynamicFormatter{defaultPattern}
}

func (f *dynamicFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	traceID := "-"
	if entry.Context != nil {
		traceID, _ = entry.Context.Value(traceIDKey).(string)
	}
	if traceID == "" {
		traceID = "-"
	}

	var fileline string
	if entry.Caller != nil {
		fileline = filepath.Base(entry.Caller.File) + ":" + strconv.Itoa(entry.Caller.Line)
	} else {
		fileline = "-"
	}

	hostname := zae.Hostname()

	fields := ""
	if len(entry.Data) > 0 {
		data := make(logrus.Fields, len(entry.Data))
		for k, v := range entry.Data {
			switch v := v.(type) {
			case error:
				data[k] = v.Error()
			default:
				data[k] = v
			}
		}
		if extraData, err := json.Marshal(data); err == nil {
			fields = string(extraData)
		}
	}

	str := templateParse(f.pattern,
		"{Level}", string(severityMap[entry.Level]),
		"{Time}", entry.Time.Format("2006-01-02 15:04:05.000"),
		"{TraceID}", traceID,
		"{FileLine}", fileline,
		"{Hostname}", hostname,
		"{ProcessID}", processID,
		"{Data}", fields,
		"{Message}", entry.Message,
	)
	entry.Buffer.WriteString(str)
	entry.Buffer.WriteByte('\n')
	return entry.Buffer.Bytes(), nil
}

func templateParse(format string, args ...string) string {
	r := strings.NewReplacer(args...)
	return r.Replace(format)
}
