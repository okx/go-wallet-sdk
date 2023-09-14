package avax

import (
	"encoding/binary"
	"errors"
)

func NewSerializer() *Serializer {
	return &Serializer{MaxLen: 256 * 1024, Body: make([]byte, 0, 128)}
}

const (
	ShortLen = 2
	IntLen   = 4
	LongLen  = 8
)

var (
	errBadLen       = errors.New("insufficient input length ")
	errNegOffset    = errors.New("negative offset")
	errInvalidInput = errors.New("input conflict")
)

type Serializer struct {
	err error

	MaxLen int
	Body   []byte
	Offset int
}

func (c *Serializer) Payload() []byte {
	if c.Offset > 0 && c.Body != nil {
		return c.Body[0:c.Offset]
	}
	return make([]byte, 0)
}

func (c *Serializer) Check(bytes int) {
	switch {
	case c.Offset < 0:
		c.Add(errNegOffset)
	case bytes < 0:
		c.Add(errInvalidInput)
	case len(c.Body)-c.Offset < bytes:
		c.Add(errBadLen)
	}
}

func (c *Serializer) AddCap(bytes int) {
	neededSize := bytes + c.Offset
	switch {
	case neededSize <= len(c.Body):
		return
	case neededSize > c.MaxLen:
		c.err = errBadLen
		return
	case neededSize <= cap(c.Body):
		c.Body = c.Body[:neededSize]
		return
	default:
		c.Body = append(c.Body[:cap(c.Body)], make([]byte, neededSize-cap(c.Body))...)
	}
}

func (c *Serializer) Errored() bool { return c.err != nil }

func (c *Serializer) Add(errors ...error) {
	if c.err == nil {
		for _, err := range errors {
			if err != nil {
				c.err = err
				break
			}
		}
	}
}

func (c *Serializer) WriteShort(val uint16, bigEndian bool) {
	c.AddCap(ShortLen)
	if c.Errored() {
		return
	}

	if bigEndian {
		binary.BigEndian.PutUint16(c.Body[c.Offset:], val)
	} else {
		binary.LittleEndian.PutUint16(c.Body[c.Offset:], val)
	}
	c.Offset += ShortLen
}

func (c *Serializer) ReadShort() uint16 {
	c.Check(ShortLen)
	if c.Errored() {
		return 0
	}

	val := binary.BigEndian.Uint16(c.Body[c.Offset:])
	c.Offset += ShortLen
	return val
}

func (c *Serializer) WriteInt(val uint32) {
	c.AddCap(IntLen)
	if c.Errored() {
		return
	}

	binary.BigEndian.PutUint32(c.Body[c.Offset:], val)
	c.Offset += IntLen
}

func (c *Serializer) ReadInt() uint32 {
	c.Check(IntLen)
	if c.Errored() {
		return 0
	}

	val := binary.BigEndian.Uint32(c.Body[c.Offset:])
	c.Offset += IntLen
	return val
}

func (c *Serializer) WriteLong(val uint64) {
	c.AddCap(LongLen)
	if c.Errored() {
		return
	}

	binary.BigEndian.PutUint64(c.Body[c.Offset:], val)
	c.Offset += LongLen
}

func (c *Serializer) ReadLong() uint64 {
	c.Check(LongLen)
	if c.Errored() {
		return 0
	}

	val := binary.BigEndian.Uint64(c.Body[c.Offset:])
	c.Offset += LongLen
	return val
}

func (c *Serializer) WriteFixedBytes(bytes []byte) {
	c.AddCap(len(bytes))
	if c.Errored() {
		return
	}

	copy(c.Body[c.Offset:], bytes)
	c.Offset += len(bytes)
}

func (c *Serializer) ReadFixedBytes(size int) []byte {
	c.Check(size)
	if c.Errored() {
		return nil
	}

	bytes := c.Body[c.Offset : c.Offset+size]
	c.Offset += size
	return bytes
}

func (c *Serializer) WriteBytes(bytes []byte) {
	c.WriteInt(uint32(len(bytes)))
	c.WriteFixedBytes(bytes)
}

func (c *Serializer) ReadBytes() []byte {
	size := c.ReadInt()
	return c.ReadFixedBytes(int(size))
}
