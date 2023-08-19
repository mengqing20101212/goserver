package utils

import (
	"errors"
	"fmt"
	"sync"
	"time"
	"unsafe"
)

const (
	WriteTypeSmall byte = iota //小端模式
	WriteTypeBig               //大端模式
)

/*
RingBuffer 是一个环形队列，数据满了抛异常, 支持单生产者 单消费者使用
*/
type RingBuffer struct {
	data     []byte       //数据数组
	readPos  int          //读指针
	writePos int          //写指针
	makePos  int          //标识
	capacity int          //最大容量
	lock     sync.RWMutex //读写锁
	model    byte
}

func (ringBuf *RingBuffer) WriteByte(b byte) bool {
	ringBuf.checkCanWrite(1, time.Duration(-1))
	ringBuf.writeVal(uint64(b), 1)
	return true
}
func (ringBuf *RingBuffer) WriteByteWhiteTimeOut(b byte, timeout int) bool {
	ringBuf.checkCanWrite(1, time.Duration(timeout))
	ringBuf.writeVal(uint64(b), 1)
	return true
}
func (ringBuf *RingBuffer) ReadByte() byte {
	ringBuf.checkCanRead(1, time.Duration(-1))
	val := readVal[byte](ringBuf, 1)
	return val
}

func (ringBuf *RingBuffer) ReadByteWhiteTimeOut(timeout int) byte {
	ringBuf.checkCanRead(1, time.Duration(timeout))
	val := readVal[byte](ringBuf, 1)
	return val
}

func (ringBuf *RingBuffer) WriteUint16(val uint16) bool {
	ringBuf.checkCanWrite(2, time.Duration(-1))
	ringBuf.writeVal(uint64(val), 2)
	return true
}

func (ringBuf *RingBuffer) ReadUint16() uint16 {
	ringBuf.checkCanRead(2, time.Duration(-1))
	val := readVal[uint16](ringBuf, 2)
	return val
}

func (ringBuf *RingBuffer) WriteUint16WhiteTimeOut(val uint16, timeout int) bool {
	ringBuf.checkCanWrite(2, time.Duration(timeout))
	ringBuf.writeVal(uint64(val), 2)
	return true
}
func (ringBuf *RingBuffer) ReadUint16WhiteTimeOut(timeout int) uint16 {
	ringBuf.checkCanRead(2, time.Duration(timeout))
	val := readVal[uint16](ringBuf, 2)
	return val
}

func (ringBuf *RingBuffer) writeVal(val uint64, len int) {
	bs := make([]byte, 4)
	for i := 0; i < 4; i++ {
		b := byte(val & 0xff)
		bs[i] = b
		val >>= 8
	}
	pos := ringBuf.writePos
	if ringBuf.model == WriteTypeBig {
		for i := 0; i < len; i++ {
			ringBuf.data[pos] = bs[len-1-i]
			pos = ringBuf.writePosAutoincrement(pos)
		}
	} else if ringBuf.model == WriteTypeSmall {
		for i := 0; i < len; i++ {
			ringBuf.data[pos] = bs[i]
			pos = ringBuf.writePosAutoincrement(pos)
		}
	}
	ringBuf.writePos = pos
}

func (ringBuf *RingBuffer) Rest() {
	ringBuf.lock.Lock()
	defer ringBuf.lock.Unlock()
	ringBuf.readPos = 0
	ringBuf.writePos = 0
	ringBuf.makePos = -1
}

func (ringBuf *RingBuffer) checkCanWrite(writeLen int, timeout time.Duration) {
	canWriteLen := ringBuf.canWriteLen()
	if canWriteLen < writeLen {
		if timeout != -1 {
			time.Sleep(time.Millisecond * timeout)
			canWriteLen = ringBuf.canWriteLen()
		}
		if canWriteLen < writeLen {
			panic(errors.New(fmt.Sprintf("not write more data [%s] , write len:%d", ringBuf.toString(), writeLen)))
		}
	}
}

func (ringBuf *RingBuffer) MakeMask() {
	ringBuf.makePos = ringBuf.readPos
}
func (ringBuf *RingBuffer) RestMask() {
	ringBuf.readPos = ringBuf.makePos
	ringBuf.makePos = -1
}

func (ringBuf *RingBuffer) toString() string {
	return fmt.Sprintf(" ptr:%p, readPos:%d, writePos:%d, capacity:%d , makePos:%d dataPtr:%p", ringBuf, ringBuf.readPos, ringBuf.writePos, ringBuf.capacity, ringBuf.makePos, &ringBuf.data)
}

func (ringBuf *RingBuffer) canWriteLen() int {
	ringBuf.lock.Lock()
	defer ringBuf.lock.Unlock()
	if ringBuf.writePos >= ringBuf.readPos {
		return ringBuf.capacity - ringBuf.writePos + ringBuf.readPos
	}
	return ringBuf.readPos - ringBuf.writePos
}

func (ringBuf *RingBuffer) canReadLen() int {
	ringBuf.lock.RLock()
	defer ringBuf.lock.RUnlock()
	if ringBuf.writePos >= ringBuf.readPos {
		return ringBuf.writePos - ringBuf.readPos
	}
	return ringBuf.capacity - ringBuf.readPos + ringBuf.writePos
}

func (ringBuf *RingBuffer) writePosAutoincrement(pos int) int {
	pos++
	if ringBuf.capacity <= pos {
		pos -= ringBuf.capacity
	}
	return pos
}

func (ringBuf *RingBuffer) checkCanRead(readLen int, timeout time.Duration) {
	canReadLen := ringBuf.canReadLen()
	if canReadLen < readLen {
		if timeout > 0 {
			time.Sleep(time.Millisecond * timeout)
			canReadLen = ringBuf.canReadLen()
		}
		if canReadLen < readLen {
			panic(errors.New(fmt.Sprintf("not read more data [%s] , write len:%d", ringBuf.toString(), readLen)))
		}
	}
}

func (ringBuf *RingBuffer) readPosAutoincrement(pos int) int {
	if pos >= ringBuf.capacity {
		pos -= ringBuf.capacity
	}
	return pos
}

func readVal[T byte | uint16 | int16 | int32 | int | uint32 | int64 | uint64 | float32 | float64](ringBuf *RingBuffer, len int) T {
	var res uint64
	bs := make([]byte, 4)
	pos := ringBuf.readPos
	readPos := pos
	if ringBuf.model == WriteTypeBig {
		for i := 0; i < len; i++ {
			readPos = ringBuf.readPosAutoincrement(pos + i)
			bs[i] = ringBuf.data[readPos]
		}
	} else {
		j := 0
		for i := 4; i > 0; i-- {
			readPos = ringBuf.readPosAutoincrement(pos + i - 1)
			bs[j] = ringBuf.data[readPos]
			j++
		}

	}
	result := T(res)
	typeLen := unsafe.Sizeof(result)
	if typeLen >= 1 {
		res |= uint64(bs[0])
	}
	if typeLen >= 2 {
		res <<= 8
		res |= uint64(bs[1])
	}
	if typeLen >= 3 {
		res <<= 8
		res |= uint64(bs[2])
	}
	if typeLen >= 4 {
		res <<= 8
		res |= uint64(bs[3])
	}
	result = T(res)
	ringBuf.readPos = readPos + 1
	return result
}

func NewRingBuffer(capacity int, model byte) *RingBuffer {
	ringBuf := RingBuffer{data: make([]byte, capacity), readPos: 0, writePos: 0, capacity: capacity, makePos: -1, model: model}
	return &ringBuf
}

func CallNewCapatity(len int) int {
	highestBitPosition := -1
	tempLen := len
	for tempLen > 0 {
		tempLen >>= 1
		highestBitPosition += 1
	}
	if len < 1024 {
		return 1 << (highestBitPosition + 1)
	}
	tempLen = (1 << highestBitPosition) + 1024
	for tempLen < len {
		tempLen += 1024
	}
	return tempLen
}
