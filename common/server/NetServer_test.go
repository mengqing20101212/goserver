package server

import (
	"bytes"
	"fmt"
	"goserver/common/logger"
	"testing"
	"time"
)

func TestByteServer(t *testing.T) {
	logger.Init("../logs", "test.log")
	server := NewServer(2001)
	server.Start()
}

func TestConnecter(t *testing.T) {
	codeProto := PackageFactory{}
	for i := 0; i < 1000; i++ {
		go func() {
			sc := CreateConnect("127.0.0.1:2001")
			logger.Init("../logs", "test1.log")
			for {
				pack := CreatePackage(int32(i), int32(i), uint32(time.Now().Unix()), uint16(i), bytes.NewBufferString(fmt.Sprintf("hello word %d", i)).Bytes())
				sc.SendMsg(codeProto.Encode(pack))
				logger.Info(fmt.Sprintf("send msg:%s", pack))
				time.Sleep(20 * time.Millisecond)
			}
		}()
	}

	time.Sleep(10240 * time.Second)
}
