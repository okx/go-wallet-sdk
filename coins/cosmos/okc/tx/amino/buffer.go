package amino

import (
	"bytes"
	"sync"
)

type Buffer struct {
	bytes.Buffer
}

func NewBuffer(buf []byte) *Buffer {
	return &Buffer{*bytes.NewBuffer(buf)}
}

func NewBufferString(s string) *Buffer {
	return &Buffer{*bytes.NewBufferString(s)}
}

func (b *Buffer) BytesCopy() []byte {
	bz := b.Bytes()
	length := len(bz)

	if length == 0 {
		return nil
	} else {
		ret := make([]byte, length)
		copy(ret, bz)
		return ret
	}
}

var bufferPool = &sync.Pool{
	New: func() interface{} {
		return NewBuffer(nil)
	},
}

// GetBuffer returns a new bytes.Buffer from the pool.
// you must call PutBuffer on the buffer when you are done with it.
func GetBuffer() *Buffer {
	return bufferPool.Get().(*Buffer)
}

// PutBuffer returns a bytes.Buffer to the pool.
func PutBuffer(b *Buffer) {
	b.Reset()
	bufferPool.Put(b)
}

func GetBytesBufferCopy(buf *bytes.Buffer) []byte {
	if buf == nil {
		return nil
	}
	bz := buf.Bytes()
	length := len(bz)
	if length == 0 {
		return nil
	} else {
		ret := make([]byte, length)
		copy(ret, bz)
		return ret
	}
}
