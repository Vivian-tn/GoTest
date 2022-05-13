package log

import (
	"bytes"
	"fmt"
	"strings"

	"git.in.zhihu.com/go/logrus"
)

var (
	defaultFormatter = &customFormatter{}
)

type customFormatter struct {
}

// Example:
// 	[E0627 19:43:21.629] messages	error=error message
//
func (*customFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	buffer := bytes.NewBufferString("[")
	buffer.WriteString(string(strings.ToUpper(entry.Level.String())[0]))
	buffer.WriteString(entry.Time.Format("0102 15:04:05.000"))
	buffer.WriteString("] ")
	buffer.WriteString(entry.Message)
	buffer.WriteByte('\t')
	for key, value := range entry.Data {
		buffer.WriteString(key)
		buffer.WriteByte('=')
		buffer.WriteString(fmt.Sprint(value))
		buffer.WriteByte(' ')
	}
	buffer.WriteByte('\n')
	return buffer.Bytes(), nil
}
