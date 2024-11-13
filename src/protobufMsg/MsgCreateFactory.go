package protobufMsg

import (
	"fmt"
	"github.com/golang/protobuf/proto"
)

func CreateProtoRequestMessage(cmd int32) (msg proto.Message) {
	switch cmd {

	case int32(CMD_login):
		return &CsLogin{}

	default:
		fmt.Println("CreateProtoRequestMessage not found cmdId:", cmd)
		return nil

	}
}

func CreateProtoResponseMessage(cmd int32) (msg proto.Message) {
	switch cmd {

	case int32(CMD_login):
		return &ScLogin{}

	default:
		fmt.Println("CreateProtoResponseMessage not found cmdId:", cmd)
		return nil

	}
}
