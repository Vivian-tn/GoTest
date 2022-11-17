package md5

import (
	"unsafe"
)

// String converts slice to string without copy.
// Use at your own risk.
func Bytes2String(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return *(*string)(unsafe.Pointer(&b))
}

// Slice converts string to slice without copy.
// Use at your own risk.
func String2Bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
