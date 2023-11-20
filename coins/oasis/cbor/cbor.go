// Package cbor provides helpers for encoding and decoding canonical CBOR.
//
// Using this package will produce canonical encodings which can be used
// in cryptographic contexts like signing as the same message is guaranteed
// to always have the same serialization.
package cbor

import (
	"github.com/okx/go-wallet-sdk/crypto/cbor"
	"io"
)

// RawMessage is a raw encoded CBOR value. It implements Marshaler and
// Unmarshaler interfaces and can be used to delay CBOR decoding or
// precompute a CBOR encoding.
type RawMessage = cbor.RawMessage

// Marshaler is the interface implemented by types that can marshal themselves
// into valid CBOR.
type Marshaler = cbor.Marshaler

// Unmarshaler is the interface implemented by types that wish to unmarshal
// CBOR data themselves.  The input is a valid CBOR value. UnmarshalCBOR
// must copy the CBOR data if it needs to use it after returning.
type Unmarshaler = cbor.Unmarshaler

var (
	encOptions = cbor.EncOptions{
		Sort:          cbor.SortCanonical,
		ShortestFloat: cbor.ShortestFloat16,
		NaNConvert:    cbor.NaNConvert7e00,
		InfConvert:    cbor.InfConvertFloat16,
		IndefLength:   cbor.IndefLengthForbidden,
		Time:          cbor.TimeUnix,
		TagsMd:        cbor.TagsForbidden,
	}

	// decOptions are decoding options for UNTRUSTED inputs (used by default).
	decOptions = cbor.DecOptions{
		DupMapKey:         cbor.DupMapKeyEnforcedAPF,
		IndefLength:       cbor.IndefLengthForbidden,
		TagsMd:            cbor.TagsForbidden,
		ExtraReturnErrors: cbor.ExtraDecErrorUnknownField,
		MaxArrayElements:  10_000_000, // Usually limited by blob size limits anyway.
		MaxMapPairs:       10_000_000, // Usually limited by blob size limits anyway.
	}

	// decOptionsTrusted are decoding options for TRUSTED inputs. They are only used when explicitly
	// requested by using the UnmarshalTrusted method.
	decOptionsTrusted = cbor.DecOptions{
		MaxArrayElements: 2147483647, // Maximum allowed.
		MaxMapPairs:      2147483647, // Maximum allowed.
	}

	encMode        cbor.EncMode
	decMode        cbor.DecMode
	decModeTrusted cbor.DecMode
)

func init() {
	var err error
	if encMode, err = encOptions.EncMode(); err != nil {
		panic(err)
	}
	if decMode, err = decOptions.DecMode(); err != nil {
		panic(err)
	}
	if decModeTrusted, err = decOptionsTrusted.DecMode(); err != nil {
		panic(err)
	}
}

// FixSliceForSerde will convert `nil` to `[]byte` to work around serde
// brain damage.
func FixSliceForSerde(b []byte) []byte {
	if b != nil {
		return b
	}
	return []byte{}
}

// Marshal serializes a given type into a CBOR byte vector.
func Marshal(src interface{}) []byte {
	b, err := encMode.Marshal(src)
	if err != nil {
		panic("common/cbor: failed to marshal: " + err.Error())
	}
	return b
}

// Unmarshal deserializes a CBOR byte vector into a given type.
func Unmarshal(data []byte, dst interface{}) error {
	if data == nil {
		return nil
	}

	return decMode.Unmarshal(data, dst)
}

// UnmarshalTrusted deserializes a CBOR byte vector into a given type.
//
// This method MUST ONLY BE USED FOR TRUSTED INPUTS as it relaxes some decoding restrictions.
func UnmarshalTrusted(data []byte, dst interface{}) error {
	if data == nil {
		return nil
	}

	return decModeTrusted.Unmarshal(data, dst)
}

// MustUnmarshal deserializes a CBOR byte vector into a given type.
// Panics if unmarshal fails.
func MustUnmarshal(data []byte, dst interface{}) {
	if err := Unmarshal(data, dst); err != nil {
		panic(err)
	}
}

// NewEncoder creates a new CBOR encoder.
func NewEncoder(w io.Writer) *cbor.Encoder {
	return encMode.NewEncoder(w)
}

// NewDecoder creates a new CBOR decoder.
func NewDecoder(r io.Reader) *cbor.Decoder {
	return decMode.NewDecoder(r)
}
