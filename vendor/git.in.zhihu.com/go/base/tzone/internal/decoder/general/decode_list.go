package general

import "git.in.zhihu.com/go/base/tzone/internal/decoder/protocol"

type generalListDecoder struct {
}

func (decoder *generalListDecoder) Decode(val interface{}, iter protocol.Iterator) {
	*val.(*List) = readList(iter).(List)
}

func readList(iter protocol.Iterator) interface{} {
	elemType, length := iter.ReadListHeader()
	generalReader := generalReaderOf(elemType)
	var generalList List
	for i := 0; i < length; i++ {
		generalList = append(generalList, generalReader(iter))
	}
	return generalList
}
