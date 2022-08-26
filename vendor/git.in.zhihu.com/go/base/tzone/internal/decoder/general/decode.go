package general

import (
	"reflect"

	"git.in.zhihu.com/go/base/tzone/internal/decoder/protocol"
)

type ValDecoder interface {
	Decode(val interface{}, iter protocol.Iterator)
}

func DecoderOf(valType reflect.Type) ValDecoder {
	switch valType {
	case reflect.TypeOf((*List)(nil)):
		return &generalListDecoder{}
	case reflect.TypeOf((*Map)(nil)):
		return &generalMapDecoder{}
	case reflect.TypeOf((*Struct)(nil)):
		return &generalStructDecoder{}
	case reflect.TypeOf((*Message)(nil)):
		return &messageDecoder{}
	case reflect.TypeOf((*protocol.MessageHeader)(nil)):
		return &messageHeaderDecoder{}
	}
	return nil
}

func generalReaderOf(ttype protocol.TType) func(iter protocol.Iterator) interface{} {
	switch ttype {
	case protocol.TypeBool:
		return readBool
	case protocol.TypeI08:
		return readInt8
	case protocol.TypeI16:
		return readInt16
	case protocol.TypeI32:
		return readInt32
	case protocol.TypeI64:
		return readInt64
	case protocol.TypeString:
		return readString
	case protocol.TypeDouble:
		return readFloat64
	case protocol.TypeList:
		return readList
	case protocol.TypeMap:
		return readMap
	case protocol.TypeStruct:
		return readStruct
	default:
		panic("unsupported type")
	}
}

func readFloat64(iter protocol.Iterator) interface{} {
	return iter.ReadFloat64()
}

func readBool(iter protocol.Iterator) interface{} {
	return iter.ReadBool()
}

func readInt8(iter protocol.Iterator) interface{} {
	return iter.ReadInt8()
}

func readInt16(iter protocol.Iterator) interface{} {
	return iter.ReadInt16()
}

func readInt32(iter protocol.Iterator) interface{} {
	return iter.ReadInt32()
}

func readInt64(iter protocol.Iterator) interface{} {
	return iter.ReadInt64()
}

func readString(iter protocol.Iterator) interface{} {
	return iter.ReadString()
}
