package protobufMsg

import (
	"fmt"
	"github.com/golang/protobuf/proto"
)

func CreateProtoRequestMessage(cmd int32) (msg proto.Message) {
	switch cmd {

	case int32(100):
		return &CsLogin{}

	case int32(10000):
		return &CsServer2Server{}

	default:
		fmt.Println("CreateProtoRequestMessage not found cmdId:", cmd)
		return nil

	}
}

func CreateProtoResponseMessage(cmd int32) (msg proto.Message) {
	switch cmd {

	case int32(100):
		return &ScLogin{}

	case int32(10000):
		return &ScServer2Server{}

	default:
		fmt.Println("CreateProtoResponseMessage not found cmdId:", cmd)
		return nil

	}
}
