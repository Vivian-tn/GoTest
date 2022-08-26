package general

import "git.in.zhihu.com/go/base/tzone/internal/decoder/protocol"

type generalStructDecoder struct {
}

func (decoder *generalStructDecoder) Decode(val interface{}, iter protocol.Iterator) {
	*val.(*Struct) = readStruct(iter).(Struct)
}

func readStruct(iter protocol.Iterator) interface{} {
	generalStruct := Struct{}
	iter.ReadStructHeader()
	for {
		fieldType, fieldId := iter.ReadStructField()
		if fieldType == protocol.TypeStop {
			return generalStruct
		}
		generalReader := generalReaderOf(fieldType)
		generalStruct[fieldId] = generalReader(iter)
	}
}
