package statsd

import (
	"strings"
)

// ValidBucket check the bucket valid or not.
// The key should be checked follow this rule:
// https://github.com/statsd/statsd/blob/b7f1d9daf3abc363542c3f1bc7b750a40c243b8b/backends/graphite.js#L168
func ValidBucket(bucket string) bool {
	rs := ([]rune)(bucket)
	for i, r := range rs {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '-' || r == '_':
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

const anyBadChar = "!"

// Join joins any number of bucket elements into a single bucket
func Join(elms ...string) (bucket string) {
	lazyBuf := elms
	lazyCreate := false
	for i, elm := range elms {
		elm = strings.Trim(elm, ".")
		if elm == "" {
			elm = anyBadChar
		}
		// check if changed elm length
		if len(elm) != len(elms[i]) && !lazyCreate {
			lazyCreate = true
			lazyBuf = make([]string, len(elms))
			copy(lazyBuf, elms[:i])
		}
		lazyBuf[i] = elm
	}
	return strings.Join(lazyBuf, ".")
}

func Node(name string) string {
	rebuild := false
	var rs []byte
	var n = len(name)
	for i := 0; i < n; i++ {
		r := name[i]
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '-' || r == '_':
		default:
			r = '_'
			if !rebuild {
				rebuild = true
				rs = make([]byte, len(name))
				copy(rs, name[:i])
			}
		}
		if rebuild {
			rs[i] = r
		}
	}
	if rebuild {
		return string(rs)
	}
	return name
}
