package statsd

import (
	"bytes"
	"strings"
	"sync"
)

var (
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
)

// Verify check the bucket valid or not.
// The key should be checked follow this rule:
// https://github.com/statsd/statsd/blob/b7f1d9daf3abc363542c3f1bc7b750a40c243b8b/backends/graphite.js#L168
func Verify(rs string) bool {
	if rs == "" {
		return false
	}
	for i, r := range rs {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || (r == '-' || r == '_'):
		case r == '.':
			switch {
			case i == 0:
				return false
			case i+1 == len(rs):
				return false
			case rs[i+1] == '.':
				return false
			default:
			}
		default:
			return false
		}
	}
	return true
}

// Join joins any number of bucket elements into a single bucket
func Join(elms ...string) (bucket string) {
	buffer := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buffer)
	buffer.Reset()

	for i, elm := range elms {
		elm = strings.Trim(elm, ".")
		if elm == "" {
			elm = "!"
		}
		if i > 0 {
			buffer.WriteByte('.')
		}
		buffer.WriteString(elm)
	}
	return buffer.String()
}

func Node(rs string) string {
	buffer := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buffer)
	buffer.Reset()

	var insertUnderline bool
	for _, r := range rs {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' {
			if insertUnderline {
				buffer.WriteByte('_')
			}
			buffer.WriteRune(r)
			insertUnderline = false
		} else {
			insertUnderline = true
		}
	}
	return buffer.String()
}
