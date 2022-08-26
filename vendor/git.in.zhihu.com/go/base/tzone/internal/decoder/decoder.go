package decoder

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync"

	"git.apache.org/thrift.git/lib/go/thrift"
	"git.in.zhihu.com/go/base/telemetry"
	"git.in.zhihu.com/go/base/tzone/internal/decoder/general"
	"git.in.zhihu.com/go/base/tzone/internal/decoder/protocol/binary"
)

var (
	readerPool = &sync.Pool{
		New: func() interface{} {
			return bufio.NewReaderSize(nil, 128)
		},
	}
	thriftBufferPool = &sync.Pool{
		New: func() interface{} {
			return thrift.NewTMemoryBuffer()
		},
	}
)

func UnmarshalMessage(buf []byte) (msg general.Message, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("decode error")
		}
	}()

	iter := binary.NewIterator(buf)
	decoder := general.DecoderOf(reflect.TypeOf(&msg))
	decoder.Decode(&msg, iter)
	err = iter.Error()
	return
}

func DumpBody(buf []byte) (service string, req map[string]interface{}, err error) {
	msg, err := UnmarshalMessage(buf)
	if err != nil {
		return "", nil, err
	}
	data, err := json.Marshal(msg.Arguments)
	if err != nil {
		return "", nil, err
	}
	err = json.Unmarshal(data, &req)
	if err != nil {
		return "", nil, err
	}
	return msg.MessageName, req, nil
}

func DumpArgumentsTStruct(args thrift.TStruct) (req telemetry.Arguments) {
	buffer := thriftBufferPool.Get().(*thrift.TMemoryBuffer)
	defer thriftBufferPool.Put(buffer)
	buffer.Reset()

	oprot := thrift.NewTSimpleJSONProtocol(buffer)
	_ = args.Write(oprot)
	_ = oprot.Flush()
	_ = json.Unmarshal(buffer.Bytes(), &req)
	return req
}

const (
	VersionMask = 0xffff0000
	Version1    = 0x80010000
)

func ReadMethod(r io.Reader) (string, error) {
	reader := readerPool.Get().(*bufio.Reader)
	defer readerPool.Put(reader)
	reader.Reset(r)

	b, err := reader.Peek(8)
	if err != nil {
		return "", err
	}
	versionAndMessageType := int32(uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24)
	version := int64(versionAndMessageType) & VersionMask
	if version != Version1 {
		return "", fmt.Errorf("fetch tzone method failed: unexpected version %d", version)
	}

	nameLength := uint32(b[7]) | uint32(b[6])<<8 | uint32(b[5])<<16 | uint32(b[4])<<24
	b, err = reader.Peek(8 + int(nameLength))
	if err != nil {
		return "", err
	}
	return string(b[8:]), nil
}
