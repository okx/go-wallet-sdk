package bitcoin

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	// MaxVarIntPayload is the maximum payload size for a variable length integer.
	MaxVarIntPayload = 9
)

// errNonCanonicalVarInt is the common format string used for non-canonically
// encoded variable length integer errors.
var errNonCanonicalVarInt = "non-canonical varint %x - discriminant %x must " +
	"encode a value greater than %x"

// MessageError describes an issue with a message.
// An example of some potential issues are messages from the wrong bitcoin
// network, invalid commands, mismatched checksums, and exceeding max payloads.
//
// This provides a mechanism for the caller to type assert the error to
// differentiate between general io errors such as io.EOF and issues that
// resulted from malformed messages.
type MessageError struct {
	Func        string // Function name
	Description string // Human readable description of the issue
}

// Error satisfies the error interface and prints human-readable errors.
func (e *MessageError) Error() string {
	if e.Func != "" {
		return fmt.Sprintf("%v: %v", e.Func, e.Description)
	}
	return e.Description
}

// messageError creates an error for the given function and description.
func messageError(f string, desc string) *MessageError {
	return &MessageError{Func: f, Description: desc}
}

var (
	// littleEndian is a convenience variable since binary.LittleEndian is
	// quite long.
	littleEndian = binary.LittleEndian

	// bigEndian is a convenience variable since binary.BigEndian is quite
	// long.
	bigEndian = binary.BigEndian
)

func ReadVarIntBuf(r io.Reader, pver uint32, buf []byte) (uint64, error) {
	if _, err := io.ReadFull(r, buf[:1]); err != nil {
		return 0, err
	}
	discriminant := buf[0]

	var rv uint64
	switch discriminant {
	case 0xff:
		if _, err := io.ReadFull(r, buf); err != nil {
			return 0, err
		}
		rv = littleEndian.Uint64(buf)

		// The encoding is not canonical if the value could have been
		// encoded using fewer bytes.
		min := uint64(0x100000000)
		if rv < min {
			return 0, messageError("ReadVarInt", fmt.Sprintf(
				errNonCanonicalVarInt, rv, discriminant, min))
		}

	case 0xfe:
		if _, err := io.ReadFull(r, buf[:4]); err != nil {
			return 0, err
		}
		rv = uint64(littleEndian.Uint32(buf[:4]))

		// The encoding is not canonical if the value could have been
		// encoded using fewer bytes.
		min := uint64(0x10000)
		if rv < min {
			return 0, messageError("ReadVarInt", fmt.Sprintf(
				errNonCanonicalVarInt, rv, discriminant, min))
		}

	case 0xfd:
		if _, err := io.ReadFull(r, buf[:2]); err != nil {
			return 0, err
		}
		rv = uint64(littleEndian.Uint16(buf[:2]))

		// The encoding is not canonical if the value could have been
		// encoded using fewer bytes.
		min := uint64(0xfd)
		if rv < min {
			return 0, messageError("ReadVarInt", fmt.Sprintf(
				errNonCanonicalVarInt, rv, discriminant, min))
		}

	default:
		rv = uint64(discriminant)
	}

	return rv, nil
}
