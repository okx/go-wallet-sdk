/*
*
MIT License

Copyright (c) 2018 WavesPlatform
*/
package types

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/waves/crypto"
	"github.com/okx/go-wallet-sdk/crypto/base58"
)

const (
	WavesAddressVersion      byte = 0x01
	AddressIDSize                 = 20
	wavesAddressHeaderSize        = 2
	wavesAddressBodySize          = AddressIDSize
	wavesAddressChecksumSize      = 4

	WavesAddressSize = wavesAddressHeaderSize + wavesAddressBodySize + wavesAddressChecksumSize
)

type WavesAddress [WavesAddressSize]byte

// String produces the BASE58 string representation of the WavesAddress.
func (a WavesAddress) String() string {
	return base58.Encode(a[:])
}

// MarshalJSON is the custom JSON marshal function for the WavesAddress.
func (a WavesAddress) MarshalJSON() ([]byte, error) {
	return B58Bytes(a[:]).MarshalJSON()
}

// Valid checks that version and checksum of the WavesAddress are correct.
func (a *WavesAddress) Valid() (bool, error) {
	if a[0] != WavesAddressVersion {
		return false, fmt.Errorf("unsupported address version %d", a[0])
	}
	hb := a[:wavesAddressHeaderSize+wavesAddressBodySize]
	ec, err := addressChecksum(hb)
	if err != nil {
		return false, err
	}
	ac := a[wavesAddressHeaderSize+wavesAddressBodySize:]
	if !bytes.Equal(ec, ac) {
		return false, errors.New("invalid WavesAddress checksum")
	}
	return true, nil
}

func addressChecksum(b []byte) ([]byte, error) {
	h, err := crypto.SecureHash(b)
	if err != nil {
		return nil, err
	}
	c := make([]byte, wavesAddressChecksumSize)
	copy(c, h[:wavesAddressChecksumSize])
	return c, nil
}

func NewAddressFromPublicKeyHash(scheme byte, pubKeyHash []byte) (WavesAddress, error) {
	var addr WavesAddress
	addr[0] = WavesAddressVersion
	addr[1] = scheme
	copy(addr[wavesAddressHeaderSize:], pubKeyHash[:wavesAddressBodySize])
	checksum, err := addressChecksum(addr[:wavesAddressHeaderSize+wavesAddressBodySize])
	if err != nil {
		return addr, err
	}
	copy(addr[wavesAddressHeaderSize+wavesAddressBodySize:], checksum)
	return addr, nil
}

// NewAddressFromString creates a WavesAddress from its string representation. This function checks that the address is valid.
func NewAddressFromString(s string) (WavesAddress, error) {
	var a WavesAddress
	var err error
	b := base58.Decode(s)
	a, err = NewAddressFromBytes(b)
	if err != nil {
		return a, err
	}
	return a, nil
}

// NewAddressFromBytes creates a WavesAddress from the slice of bytes and checks that the result address is valid address.
func NewAddressFromBytes(b []byte) (WavesAddress, error) {
	var a WavesAddress
	if l := len(b); l < WavesAddressSize {
		return a, fmt.Errorf("insufficient array length %d, expected at least %d", l, WavesAddressSize)
	}
	copy(a[:], b[:WavesAddressSize])
	if ok, err := a.Valid(); !ok {
		return a, err
	}
	return a, nil
}
