package server

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"logger"
	"protobufMsg"
	"strconv"
	"testing"
	"time"
)

func TestByteServer(t *testing.T) {
	logger.Init("../logs", "test.log")
	server := NewServer(2001)
	server.Start()
}

func TestConnecter(t *testing.T) {
	logger := logger.SystemLogger
	//logger.Init("../logs", "test1.log")
	for i := 0; i < 5000; i++ {
		time.Sleep(3 * time.Millisecond)
		go func(i int) {
			codeProto := PackageFactory{}
			sc := CreateConnect("127.0.0.1:2001", &codeProto)
			req := protobufMsg.CsLogin{
				Name: "hello Name:" + strconv.Itoa(i),
				Male: i%2 == 1,
			}
			for {
				flag, responsePack := sc.SendMsgData(int32(protobufMsg.CMD_Login), &req, 1000)
				if flag {
					response := protobufMsg.ScLogin{}
					proto.Unmarshal(responsePack.Body, &response)
					logger.Info(fmt.Sprintf("receive msg:%s ", response.String()))
				}
				time.Sleep(1 * time.Second)
			}
		}(i)
	}

	time.Sleep(10240 * time.Second)
}
