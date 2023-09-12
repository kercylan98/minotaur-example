package packet

import "google.golang.org/protobuf/proto"

type Reader func(message proto.Message) proto.Message

func (slf Reader) ReadTo(message proto.Message) proto.Message {
	return slf(message)
}

func (slf Reader) Read(message proto.Message) proto.Message {
	return slf(message)
}
