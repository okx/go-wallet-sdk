package crypto

import (
	"crypto/hmac"
	"crypto/sha512"
	"fmt"

	"filippo.io/edwards25519"
	"github.com/okx/go-wallet-sdk/coins/cardano/ed25519"
)

// Derive derives a children XPrv using BIP32-Ed25519
func (xsk XPrvKey) Derive(index uint32) XPrvKey {
	xpriv := xsk[:64]
	chainCode := xsk[64:]
	zmac := hmac.New(sha512.New, chainCode)
	ccmac := hmac.New(sha512.New, chainCode)

	sindex := serializeIndex(index)
	if isHardenedDerivation(index) {
		zmac.Write([]byte{0x0})
		zmac.Write(xpriv)
		zmac.Write(sindex)
		ccmac.Write([]byte{0x1})
		ccmac.Write(xpriv)
		ccmac.Write(sindex)
	} else {
		pub := ed25519.PublicKeyFrom(ed25519.ExtendedPrivateKey(xpriv))
		zmac.Write([]byte{0x2})
		zmac.Write(pub)
		zmac.Write(sindex)
		ccmac.Write([]byte{0x3})
		ccmac.Write(pub)
		ccmac.Write(sindex)
	}
	z := zmac.Sum(nil)
	zl := z[:32]
	zr := z[32:64]

	kl := add28Mul8(xsk[:32], zl)
	kr := addMod256(xsk[32:64], zr)

	cc := ccmac.Sum(nil)
	cc = cc[32:]

	cxsk := make([]byte, 96)
	copy(cxsk[:32], kl)
	copy(cxsk[32:64], kr)
	copy(cxsk[64:], cc)

	return cxsk
}

// Derive derives a children XPub using BIP32-Ed25519
func (xvk XPubKey) Derive(index uint32) (XPubKey, error) {
	pub := []byte(xvk[:32])
	chainCode := []byte(xvk[32:64])
	zmac := hmac.New(sha512.New, chainCode)
	ccmac := hmac.New(sha512.New, chainCode)

	sindex := serializeIndex(index)
	if isHardenedDerivation(index) {
		return XPubKey{}, fmt.Errorf("expected soft derivation")
	}

	zmac.Write([]byte{0x2})
	zmac.Write(pub)
	zmac.Write(sindex)
	ccmac.Write([]byte{0x3})
	ccmac.Write(pub)
	ccmac.Write(sindex)

	z := zmac.Sum(nil)
	zl := z[:32]

	// PKc = PK + 8ZL*B
	// C = A + B
	var (
		zero [32]byte
		a    edwards25519.Point
		b    edwards25519.Point
		c    edwards25519.Point
		s    edwards25519.Scalar
	)

	zl8 := add28Mul8(zero[:], zl)
	_, err := s.SetCanonicalBytes(zl8)
	if err != nil {
		return nil, err
	}
	b.ScalarBaseMult(&s)
	_, err = a.SetBytes(pub)
	if err != nil {
		return nil, err
	}

	c.Add(&a, &b)

	cxvk := make([]byte, 64)
	cc := ccmac.Sum(nil)
	cc = cc[32:64]

	copy(cxvk[:32], c.Bytes())
	copy(cxvk[32:64], cc)

	return XPubKey(cxvk), nil
}

func serializeIndex(index uint32) []byte {
	return []byte{byte(index), byte(index >> 8), byte(index >> 16), byte(index >> 24)}
}

func isHardenedDerivation(index uint32) bool {
	return index >= 0x80000000
}

func add28Mul8(x, y []byte) []byte {
	out := make([]byte, 32)
	var carry uint16

	for i, xi := range x[:28] {
		r := uint16(xi) + ((uint16(y[i])) << 3) + carry
		out[i] = byte(r & 0xff)
		carry = r >> 8
	}
	for i, xi := range x[28:32] {
		r := uint16(xi) + carry
		out[i+28] = byte(r & 0xff)
		carry = r >> 8
	}

	return out
}

func addMod256(x, y []byte) []byte {
	out := make([]byte, 32)
	var carry uint16

	for i, xi := range x[:32] {
		r := uint16(xi) + uint16(y[i]) + carry
		out[i] = byte(r)
		carry = r >> 8

	}

	return out
}
