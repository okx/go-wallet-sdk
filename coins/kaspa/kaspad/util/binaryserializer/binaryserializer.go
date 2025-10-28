package binaryserializer

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

// maxItems is the number of buffers to keep in the free
// list to use for binary serialization and deserialization.
const maxItems = 1024

// Borrow returns a byte slice from the free list with a length of 8. A new
// buffer is allocated if there are not any available on the free list.
func Borrow() []byte {
	var buf []byte
	select {
	case buf = <-binaryFreeList:
	default:
		buf = make([]byte, 8)
	}
	return buf[:8]
}

// Return puts the provided byte slice back on the free list. The buffer MUST
// have been obtained via the Borrow function and therefore have a cap of 8.
func Return(buf []byte) {
	select {
	case binaryFreeList <- buf:
	default:
		// Let it go to the garbage collector.
	}
}

// Uint8 reads a single byte from the provided reader using a buffer from the
// free list and returns it as a uint8.
func Uint8(r io.Reader) (uint8, error) {
	buf := Borrow()[:1]
	if _, err := io.ReadFull(r, buf); err != nil {
		Return(buf)
		return 0, errors.WithStack(err)
	}
	rv := buf[0]
	Return(buf)
	return rv, nil
}

// Uint16 reads two bytes from the provided reader using a buffer from the
// free list, converts it to a number using the provided byte order, and returns
// the resulting uint16.
func Uint16(r io.Reader) (uint16, error) {
	buf := Borrow()[:2]
	if _, err := io.ReadFull(r, buf); err != nil {
		Return(buf)
		return 0, errors.WithStack(err)
	}
	rv := binary.LittleEndian.Uint16(buf)
	Return(buf)
	return rv, nil
}

// Uint32 reads four bytes from the provided reader using a buffer from the
// free list, converts it to a number using the provided byte order, and returns
// the resulting uint32.
func Uint32(r io.Reader) (uint32, error) {
	buf := Borrow()[:4]
	if _, err := io.ReadFull(r, buf); err != nil {
		Return(buf)
		return 0, errors.WithStack(err)
	}
	rv := binary.LittleEndian.Uint32(buf)
	Return(buf)
	return rv, nil
}

// Uint64 reads eight bytes from the provided reader using a buffer from the
// free list, converts it to a number using the provided byte order, and returns
// the resulting uint64.
func Uint64(r io.Reader) (uint64, error) {
	buf := Borrow()[:8]
	if _, err := io.ReadFull(r, buf); err != nil {
		Return(buf)
		return 0, errors.WithStack(err)
	}
	rv := binary.LittleEndian.Uint64(buf)
	Return(buf)
	return rv, nil
}

// PutUint8 copies the provided uint8 into a buffer from the free list and
// writes the resulting byte to the given writer.
func PutUint8(w io.Writer, val uint8) error {
	buf := Borrow()[:1]
	buf[0] = val
	_, err := w.Write(buf)
	Return(buf)
	return errors.WithStack(err)
}

// PutUint16 serializes the provided uint16 using the given byte order into a
// buffer from the free list and writes the resulting two bytes to the given
// writer.
func PutUint16(w io.Writer, val uint16) error {
	buf := Borrow()[:2]
	binary.LittleEndian.PutUint16(buf, val)
	_, err := w.Write(buf)
	Return(buf)
	return errors.WithStack(err)
}

// PutUint32 serializes the provided uint32 using the given byte order into a
// buffer from the free list and writes the resulting four bytes to the given
// writer.
func PutUint32(w io.Writer, val uint32) error {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], val)
	_, err := w.Write(buf[:])
	return errors.WithStack(err)
}

// PutUint64 serializes the provided uint64 using the given byte order into a
// buffer from the free list and writes the resulting eight bytes to the given
// writer.
func PutUint64(w io.Writer, val uint64) error {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], val)
	_, err := w.Write(buf[:])
	return errors.WithStack(err)
}

// binaryFreeList provides a free list of buffers to use for serializing and
// deserializing primitive integer values to and from io.Readers and io.Writers.
//
// It defines a concurrent safe free list of byte slices (up to the
// maximum number defined by the maxItems constant) that have a
// cap of 8 (thus it supports up to a uint64). It is used to provide temporary
// buffers for serializing and deserializing primitive numbers to and from their
// binary encoding in order to greatly reduce the number of allocations
// required.
//
// For convenience, functions are provided for each of the primitive unsigned
// integers that automatically obtain a buffer from the free list, perform the
// necessary binary conversion, read from or write to the given io.Reader or
// io.Writer, and return the buffer to the free list.
var binaryFreeList = make(chan []byte, maxItems)
