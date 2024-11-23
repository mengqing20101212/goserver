package utils

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestCallNewCapatity(t *testing.T) {
	type args struct {
		len int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CallNewCapatity(tt.args.len); got != tt.want {
				t.Errorf("CallNewCapatity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRingBuffer(t *testing.T) {
	type args struct {
		capacity int
		model    byte
	}
	tests := []struct {
		name string
		args args
		want *RingBuffer
	}{
		{
			args: args{capacity: 1024, model: WriteTypeBig},
			name: "TestNewRingBuffer",
			want: NewRingBuffer(1024, WriteTypeBig),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := NewRingBuffer(tt.args.capacity, tt.args.model)
			got.WriteUint16(12345)
			got.MakeMask()
			fmt.Println(got.ReadUint16())
			got.RestMask()
			fmt.Println(got.ReadUint16())
			defer func() {
				if r := recover(); r != nil {
					for got.canReadLen() > 0 {
						fmt.Println("read data:", got.ReadUint16())
					}
				}
			}()

			testFloat(got)
			testString(got)
			/*go func() {
				for i := 0; i < 1000; i++ {
					got.WriteUint16WhiteTimeOut(uint16(i), 100)
				}
				time.Sleep(2 * time.Second)
			}()
			go func() {
				for true {
					fmt.Println("read data:", got.ReadUint16WhiteTimeOut(4000))
				}

			}()*/
			/*for i := 0; i < 1000; i++ {
				got.WriteUint16(uint16(i))
				fmt.Println("read data:", got.ReadUint16())
			}*/
			time.Sleep(5 * time.Second)
			fmt.Println(got)

			/*if got := NewRingBuffer(tt.args.capacity, tt.args.model); !reflect.DeepEqual(got, tt.want) {
				for i := 0; i < 1000; i++ {
					got.WriteUint16(uint16(i))
				}
				t.Errorf("NewRingBuffer() = %v, want %v", got, tt.want)
			}*/
		})
	}
}

func testString(buffer *RingBuffer) {
	buffer.WriteString("12312, dadwa === !! 我不是中国人  \r\n \\m")
	fmt.Println(buffer.ReadString())
}

func testFloat(buffer *RingBuffer) {
	buffer.WriteFloat32(13.211)
	fmt.Println(buffer.ReadFloat32())
	buffer.WriteFloat64(12.111)
	fmt.Println(buffer.ReadFloat64())
	buffer.WriteInt16(-12)
	fmt.Println(buffer.ReadInt16())
	buffer.WriteUint64(20)
	fmt.Println(buffer.ReadUint64())
	buffer.WriteInt64(-21)
	fmt.Println(buffer.ReadInt64())
}

func TestRingBuffer_MakeMask(t *testing.T) {
	type fields struct {
		data     []byte
		readPos  int
		writePos int
		makePos  int
		capacity int
		lock     sync.RWMutex
		model    byte
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ringBuf := &RingBuffer{
				data:     tt.fields.data,
				readPos:  tt.fields.readPos,
				writePos: tt.fields.writePos,
				makePos:  tt.fields.makePos,
				capacity: tt.fields.capacity,
				model:    tt.fields.model,
			}
			ringBuf.MakeMask()
		})
	}
}

func TestRingBuffer_RestMask(t *testing.T) {
	type fields struct {
		data     []byte
		readPos  int
		writePos int
		makePos  int
		capacity int
		lock     sync.RWMutex
		model    byte
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ringBuf := &RingBuffer{
				data:     tt.fields.data,
				readPos:  tt.fields.readPos,
				writePos: tt.fields.writePos,
				makePos:  tt.fields.makePos,
				capacity: tt.fields.capacity,
				model:    tt.fields.model,
			}
			ringBuf.RestMask()
		})
	}
}

func TestRingBuffer_WriteByte(t *testing.T) {
	type fields struct {
		data     []byte
		readPos  int
		writePos int
		makePos  int
		capacity int
		lock     sync.RWMutex
		model    byte
	}
	type args struct {
		b byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ringBuf := &RingBuffer{
				data:     tt.fields.data,
				readPos:  tt.fields.readPos,
				writePos: tt.fields.writePos,
				makePos:  tt.fields.makePos,
				capacity: tt.fields.capacity,

				model: tt.fields.model,
			}
			if got := ringBuf.WriteByte(tt.args.b); got != tt.want {
				t.Errorf("WriteByte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRingBuffer_WriteUint16(t *testing.T) {
	type fields struct {
		data     []byte
		readPos  int
		writePos int
		makePos  int
		capacity int
		lock     sync.RWMutex
		model    byte
	}
	type args struct {
		val uint16
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ringBuf := &RingBuffer{
				data:     tt.fields.data,
				readPos:  tt.fields.readPos,
				writePos: tt.fields.writePos,
				makePos:  tt.fields.makePos,
				capacity: tt.fields.capacity,

				model: tt.fields.model,
			}
			if got := ringBuf.WriteUint16(tt.args.val); got != tt.want {
				t.Errorf("WriteUint16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRingBuffer_WriteUint16WhiteTimeOut(t *testing.T) {
	type fields struct {
		data     []byte
		readPos  int
		writePos int
		makePos  int
		capacity int
		lock     sync.RWMutex
		model    byte
	}
	type args struct {
		val     uint16
		timeout int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ringBuf := &RingBuffer{
				data:     tt.fields.data,
				readPos:  tt.fields.readPos,
				writePos: tt.fields.writePos,
				makePos:  tt.fields.makePos,
				capacity: tt.fields.capacity,

				model: tt.fields.model,
			}
			if got := ringBuf.WriteUint16WhiteTimeOut(tt.args.val, tt.args.timeout); got != tt.want {
				t.Errorf("WriteUint16WhiteTimeOut() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRingBuffer_checkScalingUp(t *testing.T) {
	type fields struct {
		data     []byte
		readPos  int
		writePos int
		makePos  int
		capacity int
		lock     sync.RWMutex
		model    byte
	}
	type args struct {
		writeLen int
		timeout  time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ringBuf := &RingBuffer{
				data:     tt.fields.data,
				readPos:  tt.fields.readPos,
				writePos: tt.fields.writePos,
				makePos:  tt.fields.makePos,
				capacity: tt.fields.capacity,

				model: tt.fields.model,
			}
			ringBuf.checkCanWrite(tt.args.writeLen, tt.args.timeout)
		})
	}
}

func TestRingBuffer_toString(t *testing.T) {
	type fields struct {
		data     []byte
		readPos  int
		writePos int
		makePos  int
		capacity int
		lock     sync.RWMutex
		model    byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ringBuf := &RingBuffer{
				data:     tt.fields.data,
				readPos:  tt.fields.readPos,
				writePos: tt.fields.writePos,
				makePos:  tt.fields.makePos,
				capacity: tt.fields.capacity,

				model: tt.fields.model,
			}
			if got := ringBuf.toString(); got != tt.want {
				t.Errorf("toString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRingBuffer_writeVal(t *testing.T) {
	type fields struct {
		data     []byte
		readPos  int
		writePos int
		makePos  int
		capacity int
		lock     sync.RWMutex
		model    byte
	}
	type args struct {
		val uint64
		len int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ringBuf := &RingBuffer{
				data:     tt.fields.data,
				readPos:  tt.fields.readPos,
				writePos: tt.fields.writePos,
				makePos:  tt.fields.makePos,
				capacity: tt.fields.capacity,

				model: tt.fields.model,
			}
			ringBuf.writeVal(tt.args.val, tt.args.len)
		})
	}
}
