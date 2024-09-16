package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"logger"
)

type ByteBuffer struct {
	buf    *bytes.Buffer
	offset int
	mark   int
}

var log = logger.Init("../logs", "buffer.log")

func (self *ByteBuffer) GetBuffer() *bytes.Buffer {
	return self.buf
}
func NewByteBuffer() (buffer ByteBuffer) {
	buffer = ByteBuffer{buf: &bytes.Buffer{}, offset: 0, mark: 0}
	return buffer
}
func NewByteBufferByBuf(buf *bytes.Buffer) (buffer ByteBuffer) {
	buffer = ByteBuffer{buf: buf, offset: 0, mark: 0}
	return buffer
}
func NewByteBufferByArr(data []byte) (buffer ByteBuffer) {
	buffer = ByteBuffer{buf: &bytes.Buffer{}, offset: 0, mark: 0}
	buffer.buf.Write(data)
	return buffer
}
func NewByteBufferByString(str string) (buffer ByteBuffer) {
	buffer = ByteBuffer{buf: &bytes.Buffer{}, offset: 0, mark: 0}
	buffer.buf.Write([]byte(str))
	return buffer
}

func (self *ByteBuffer) Mark() {
	self.mark = self.offset
}
func (self *ByteBuffer) RestMark() {
	self.offset = self.mark
	self.mark = 0
}

func (self *ByteBuffer) ReadByte() (b byte) {
	self.checkOffset()
	var out byte
	read := bytes.NewReader(self.buf.Bytes()[self.offset:])
	err := binary.Read(read, binary.LittleEndian, &out)
	if err != nil {
		log.Error(fmt.Sprintf("read byte error %s", err))
		return 0
	}
	self.offset++
	return out
}

func (self *ByteBuffer) ReadUint16() (u uint16, er error) {
	self.checkOffset()
	var out uint16
	read := bytes.NewReader(self.buf.Bytes()[self.offset:])
	err := binary.Read(read, binary.LittleEndian, &out)
	if err != nil {
		//log.Error(fmt.Sprintf("read uint16 error %s", err))
		return 0, err
	}
	self.offset += 2
	return out, nil
}

func (self *ByteBuffer) ReadInt16() (u int16) {
	self.checkOffset()
	var out int16
	read := bytes.NewReader(self.buf.Bytes()[self.offset:])
	err := binary.Read(read, binary.LittleEndian, &out)
	if err != nil {
		log.Error(fmt.Sprintf("read int16 error %s", err))
		return 0
	}
	self.offset += 2
	return out
}

func (self *ByteBuffer) ReadUint32() (u uint32) {
	self.checkOffset()
	var out uint32
	read := bytes.NewReader(self.buf.Bytes()[self.offset:])
	err := binary.Read(read, binary.LittleEndian, &out)
	if err != nil {
		log.Error(fmt.Sprintf("read uint32 error %s", err))
		return 0
	}
	self.offset += 4
	return out
}

func (self *ByteBuffer) ReadInt32() (u int32) {
	self.checkOffset()
	var out int32
	read := bytes.NewReader(self.buf.Bytes()[self.offset:])
	err := binary.Read(read, binary.LittleEndian, &out)
	if err != nil {
		log.Error(fmt.Sprintf("read int error %s", err))
		return 0
	}
	self.offset += 4
	return out
}

func (self *ByteBuffer) ReadUint64() (u uint64) {
	self.checkOffset()
	var out uint64
	read := bytes.NewReader(self.buf.Bytes()[self.offset:])
	err := binary.Read(read, binary.LittleEndian, &out)
	if err != nil {
		log.Error(fmt.Sprintf("read uint64 error %s", err))
		return 0
	}
	self.offset += 8
	return out
}

func (self *ByteBuffer) ReadInt64() (u int64) {
	self.checkOffset()
	var out int64
	read := bytes.NewReader(self.buf.Bytes()[self.offset:])
	err := binary.Read(read, binary.LittleEndian, &out)
	if err != nil {
		log.Error(fmt.Sprintf("read int64 error %s", err))
		return 0
	}
	self.offset += 8
	return out
}

func (self *ByteBuffer) ReadBytes(len int) (bs []byte) {
	self.checkOffset()
	out := make([]byte, len)
	read := bytes.NewReader(self.buf.Bytes()[self.offset:])
	err := binary.Read(read, binary.LittleEndian, &out)
	if err != nil {
		log.Error(fmt.Sprintf("read Bytes error %s, redLen:%d offset:%d, buf len:%d, cap:%d", err, len, self.offset, self.buf.Len(), self.buf.Cap()))
		return out
	}
	self.offset += len
	return out
}
func (self *ByteBuffer) ReadAllByte() (bs []byte) {
	self.checkOffset()
	out := make([]byte, self.buf.Len()-self.offset)
	read := bytes.NewReader(self.buf.Bytes()[self.offset:])
	err := binary.Read(read, binary.LittleEndian, &out)
	if err != nil {
		log.Error(fmt.Sprintf("read AllByte error %s", err))
		return out
	}
	self.offset += self.buf.Len()
	return out
}

func (self *ByteBuffer) WriteByte(b byte) (len int, success bool) {
	err := self.buf.WriteByte(b)
	if err != nil {
		log.Error(fmt.Sprintf("write ByteBuffer error:%s", err))
		return -1, false
	}
	return 1, true
}

func (self *ByteBuffer) WriteUint16(b uint16) (len int, success bool) {
	err := binary.Write(self.buf, binary.LittleEndian, b)
	if err != nil {
		log.Error(fmt.Sprintf("write ByteBuffer error:%s", err))
		return -1, false
	}
	return 2, true
}
func (self *ByteBuffer) WriteInt16(b int16) (len int, success bool) {
	err := binary.Write(self.buf, binary.LittleEndian, b)
	if err != nil {
		log.Error(fmt.Sprintf("write ByteBuffer error:%s", err))
		return -1, false
	}
	return 2, true
}
func (self *ByteBuffer) WriteUint32(b uint32) (len int, success bool) {
	err := binary.Write(self.buf, binary.LittleEndian, b)
	if err != nil {
		log.Error(fmt.Sprintf("write ByteBuffer error:%s", err))
		return -1, false
	}
	return 4, true
}
func (self *ByteBuffer) WriteInt32(b int32) (len int, success bool) {
	err := binary.Write(self.buf, binary.LittleEndian, b)
	if err != nil {
		log.Error(fmt.Sprintf("write ByteBuffer error:%s", err))
		return -1, false
	}
	return 4, true
}

func (self *ByteBuffer) WriteInt64(b int64) (len int, success bool) {
	err := binary.Write(self.buf, binary.LittleEndian, b)
	if err != nil {
		log.Error(fmt.Sprintf("write ByteBuffer error:%s", err))
		return -1, false
	}
	return 8, true
}

func (self *ByteBuffer) WriteUInt64(b uint64) (len int, success bool) {
	err := binary.Write(self.buf, binary.LittleEndian, b)
	if err != nil {
		log.Error(fmt.Sprintf("write ByteBuffer error:%s", err))
		return -1, false
	}
	return 8, true
}
func (self *ByteBuffer) WriteBytes(b []byte) (l int, success bool) {
	err := binary.Write(self.buf, binary.LittleEndian, b)
	if err != nil {
		log.Error(fmt.Sprintf("write ByteBuffer error:%s", err))
		return -1, false
	}
	return len(b), true
}

func (self *ByteBuffer) GetBytes() (bs []byte) {
	return self.buf.Bytes()[self.offset:self.buf.Len()]
}

func (self *ByteBuffer) checkOffset() {
	if self.offset > self.buf.Len() {
		log.Info(fmt.Sprintf("reset offset 0 , self.offset:%d > self.buf.Len():%d  cap:%d bufPtr:%p", self.offset, self.buf.Len(), self.buf.Cap(), self.buf))
		self.offset = 0
	}
}

func (self *ByteBuffer) Len() int {
	return self.buf.Len() - self.offset
}
