package utils

import (
	"fmt"
	"goserver/common/logger"
	"testing"
)

func TestByteBuffer(t *testing.T) {
	logger.Init("/logs", "test.log")
	buf := NewByteBuffer()
	buf.WriteByte(1)
	buf.WriteInt32(12321)
	buf.WriteInt16(123)
	buf.WriteInt64(21)
	buf.WriteInt64(4353)
	buf.WriteInt64(3432)
	fmt.Println(buf)
	fmt.Println("readByte:", buf.ReadByte(), "   ", buf)
	fmt.Println("ReadInt32:", buf.ReadInt32(), "   ", buf)
	fmt.Println("readInt16:", buf.ReadInt16(), "   ", buf)
	fmt.Println("readInt64:", buf.ReadInt64(), "   ", buf)
	fmt.Println("getBytes:", buf.GetBytes(), "   ", buf)

	for i := 0; i < 1000000; i++ {
		buf.WriteByte(byte(i))
		buf.buf.ReadByte()
	}
	fmt.Println("buf len:", buf.buf.Len())
}
