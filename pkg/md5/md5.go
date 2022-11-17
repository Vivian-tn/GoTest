package md5

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"strconv"
)

type MD5 struct {
	hasher hash.Hash
}

func NewMD5() *MD5 {
	return &MD5{
		hasher: md5.New(),
	}
}

func (m *MD5) WriteString(text string) {
	m.hasher.Write(String2Bytes(text))
}

func (m *MD5) WriteInt64(text int64) {
	m.hasher.Write(String2Bytes(strconv.FormatInt(text, 10)))
}

func (m *MD5) WriteInt32(text int32) {
	m.hasher.Write(String2Bytes(strconv.FormatInt(int64(text), 10)))
}

func (m *MD5) WriteInt(text int) {
	m.hasher.Write(String2Bytes(strconv.FormatInt(int64(text), 10)))
}

func (m *MD5) WriteBool(text bool) {
	m.hasher.Write(String2Bytes(strconv.FormatBool(text)))
}

func (m *MD5) WriteFloat64(text float64) {
	m.hasher.Write(String2Bytes(strconv.FormatFloat(text, 'E', -1, 64)))
}

func (m *MD5) HexDigest() string {
	return hex.EncodeToString(m.hasher.Sum(nil))
}
