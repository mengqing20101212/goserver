package server

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"goserver/common/logger"
	"goserver/protobuf"
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
	logger.Init("../logs", "test1.log")
	for i := 0; i < 5000; i++ {
		time.Sleep(3 * time.Millisecond)
		go func(i int) {
			codeProto := PackageFactory{}
			sc := CreateConnect("127.0.0.1:2001", &codeProto)
			req := protobuf.CsLogin{
				Name: "hello Name:" + strconv.Itoa(i),
				Male: i%2 == 1,
			}
			for {
				flag, responsePack := sc.SendMsgData(protobuf.CMD_cmd_login, &req)
				if flag {
					response := protobuf.ScLogin{}
					proto.Unmarshal(responsePack.body, &response)
					logger.Info(fmt.Sprintf("receive msg:%s ", response.String()))
				}
				time.Sleep(1 * time.Second)
			}
		}(i)
	}

	time.Sleep(10240 * time.Second)
}
