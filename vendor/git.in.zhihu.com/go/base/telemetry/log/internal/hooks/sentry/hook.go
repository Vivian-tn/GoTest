package sentry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"git.in.zhihu.com/go/base/telemetry/internal/halo"
	"git.in.zhihu.com/go/base/zae"
	"git.in.zhihu.com/go/logrus"
	"github.com/getsentry/raven-go"
)

var (
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
)

var (
	severityMap = map[logrus.Level]raven.Severity{
		logrus.TraceLevel: raven.DEBUG,
		logrus.DebugLevel: raven.DEBUG,
		logrus.InfoLevel:  raven.INFO,
		logrus.WarnLevel:  raven.WARNING,
		logrus.ErrorLevel: raven.ERROR,
		logrus.FatalLevel: raven.FATAL,
		logrus.PanicLevel: raven.FATAL,
	}
)

func NewHook(dsn, release string, level logrus.Level) (logrus.Hook, error) {
	client, err := raven.New(dsn)
	if err != nil {
		return nil, err
	}
	client.SetRelease(release)

	return &Hook{client: client, level: level}, nil
}

type Hook struct {
	client *raven.Client
	level  logrus.Level
}

func (h *Hook) Levels() []logrus.Level {
	var ret []logrus.Level
	for _, level := range logrus.AllLevels {
		if h.level >= level {
			ret = append(ret, level)
		}
	}
	return ret
}

func (h *Hook) Fire(entry *logrus.Entry) error {
	defer func() {
		if rval := recover(); rval != nil {
			logMessage(fmt.Sprint(rval))
		}
	}()

	packet := raven.NewPacket(entry.Message)
	packet.Level = severityMap[entry.Level]
	packet.Timestamp = raven.Timestamp(entry.Time)
	packet.Tags = raven.Tags{
		{
			Key:   "region",
			Value: zae.Region(),
		},
		{
			Key:   "environment",
			Value: zae.Environment(),
		},
		{
			Key:   "service",
			Value: zae.Service(),
		},
	}

	var err Error
	if er, ok := entry.Data[logrus.ErrorKey].(halo.Error); ok {
		delete(entry.Data, logrus.ErrorKey)

		err = Error{message: entry.Message, frames: er.Frames()}
	} else {
		err = Error{message: entry.Message, frames: halo.GetFrames()}
	}
	packet.Interfaces = append(packet.Interfaces, raven.NewException(err, err.Stacktrace()))
	packet.Culprit = err.Error()

	// log 会使用自己的 getRuntimeCaller
	// if len(err.frames) > 0 {
	// 	entry.Caller = err.Caller()
	// }

	extra := raven.Extra{}
	extra["hostname"] = zae.Hostname()

	if req, ok := entry.Data["sentry.http_request"].(*http.Request); ok {
		delete(entry.Data, "sentry.http_request")

		extra["http.proto"] = req.Proto
		extra["http.request_remote_addr"] = req.RemoteAddr
		extra["http.request_content_length"] = req.ContentLength
		extra["http.request_headers"] = headerToMap(req.Header)

		if req.Body != nil {
			buffer := bufferPool.Get().(*bytes.Buffer)
			defer bufferPool.Put(buffer)
			buffer.Reset()

			if _, err := buffer.ReadFrom(io.LimitReader(req.Body, 512)); err == nil {
				_ = req.Body.Close()

				if req.Header.Get("Content-Type") != "application/x-thrift" {
					extra["http.request_content"] = buffer.String()
				}
			}
		}
	}

	if resp, ok := entry.Data["sentry.http_response"].(*http.Response); ok {
		delete(entry.Data, "sentry.http_response")

		extra["http.response_remote_addr"] = resp.Request.RemoteAddr
		extra["http.response_headers"] = headerToMap(resp.Header)
		extra["http.response_content_length"] = resp.ContentLength
		if resp.ContentLength > 0 && resp.ContentLength < 512 {
			// 拦截的 http 请求的响应，业务代码里是需要继续读里面的数据的
			// 用 buffer pool 不能确定释放获取的 buffer 释放时间，所以这里用不了
			buffer := new(bytes.Buffer)

			if _, err := buffer.ReadFrom(io.LimitReader(resp.Body, 512)); err == nil {
				_ = resp.Body.Close()

				resp.Body = ioutil.NopCloser(buffer)
				extra["http.response_content"] = buffer.String()
			}
		}
	}

	packet.Extra = deepCopyToExtra(entry.Data, extra)

	if zae.IsDebug() {
		data, err := json.MarshalIndent(packet, "", "  ")
		if err != nil {
			return err
		}
		log.Printf("[DEBUG:Sentry] %s", string(data))
	}

	_, _ = h.client.Capture(packet, nil)
	return nil
}

func headerToMap(header http.Header) map[string]string {
	data := make(map[string]string)
	for k, v := range header {
		data[k] = strings.Join(v, ",")
	}
	return data
}

// deepcopy avoid concurrent access in raven & formater
func deepCopyToExtra(fields logrus.Fields, extra raven.Extra) raven.Extra {
	body, err := json.Marshal(fields)
	if err != nil {
		return extra
	}
	_ = json.Unmarshal(body, &extra)
	return extra
}

func logMessage(message string) {
	_, _ = fmt.Fprintf(os.Stderr, "[E %s - - %s:0] %s\n", time.Now().Format("2006-01-02 15:04:05.000"), zae.Hostname(), message)
}
