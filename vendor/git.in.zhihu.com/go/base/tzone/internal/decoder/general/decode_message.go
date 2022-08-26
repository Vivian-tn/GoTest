package general

import "git.in.zhihu.com/go/base/tzone/internal/decoder/protocol"

type messageDecoder struct {
}

func (decoder *messageDecoder) Decode(val interface{}, iter protocol.Iterator) {
	*val.(*Message) = Message{
		MessageHeader: iter.ReadMessageHeader(),
		Arguments:     readStruct(iter).(Struct),
	}
}

type messageHeaderDecoder struct {
}

func (decoder *messageHeaderDecoder) Decode(val interface{}, iter protocol.Iterator) {
	*val.(*protocol.MessageHeader) = iter.ReadMessageHeader()
}
