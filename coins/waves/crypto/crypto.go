/*
*
MIT License

Copyright (c) 2018 WavesPlatform
*/
package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"github.com/agl/ed25519/edwards25519"
	"github.com/emresenyuva/go-wallet-sdk/crypto/base58"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/sha3"
)

const (
	DigestSize     = 32
	SignatureSize  = 64
	PublicKeySize  = 32
	SecretKeySize  = 32
	PrivateKeySize = 64
)

var (
	prefix = []byte{
		0xfe, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	}
)

type Digest [DigestSize]byte

func (d Digest) MarshalJSON() ([]byte, error) {
	return ToBase58JSON(d[:]), nil
}

type Signature [SignatureSize]byte

func (s Signature) String() string {
	return base58.Encode(s[:])
}

func (s Signature) MarshalJSON() ([]byte, error) {
	return ToBase58JSON(s[:]), nil
}

type PublicKey [PublicKeySize]byte

func (k *PublicKey) Bytes() []byte {
	return k[:]
}

func (k PublicKey) String() string {
	return base58.Encode(k[:])
}

func (k PublicKey) MarshalJSON() ([]byte, error) {
	return ToBase58JSON(k[:]), nil
}

type SecretKey [SecretKeySize]byte

func (k SecretKey) Bytes() []byte {
	out := make([]byte, len(k))
	copy(out, k[:])
	return out
}

func (k SecretKey) String() string {
	return base58.Encode(k[:])
}

func ToBase58JSON(b []byte) []byte {
	s := base58.Encode(b)
	var sb bytes.Buffer
	sb.Grow(2 + len(s))
	sb.WriteRune('"')
	sb.WriteString(s)
	sb.WriteRune('"')
	return sb.Bytes()
}

func array32FromBase58(s, name string) ([32]byte, error) {
	var r [32]byte
	b := base58.Decode(s)
	if l := len(b); l != 32 {
		return r, fmt.Errorf("incorrect %s length %d, expected %d", name, l, 32)
	}
	copy(r[:], b[:32])
	return r, nil
}

func NewDigestFromBase58(s string) (Digest, error) {
	return array32FromBase58(s, "Digest")
}

func NewDigestFromBytes(b []byte) (Digest, error) {
	if len(b) != DigestSize {
		return Digest{}, errors.New("invalid digest len")
	}
	var r Digest
	copy(r[:], b)
	return r, nil
}

func NewSecretKeyFromBase58(s string) (SecretKey, error) {
	return array32FromBase58(s, "SecretKey")
}

func NewPublicKeyFromBase58(s string) (PublicKey, error) {
	return array32FromBase58(s, "PublicKey")
}

func NewPublicKeyFromBytes(b []byte) (PublicKey, error) {
	var pk PublicKey
	if l := len(b); l < PublicKeySize {
		return pk, fmt.Errorf("insufficient array length %d, expected atleast %d", l, PublicKeySize)
	}
	copy(pk[:], b[:PublicKeySize])
	return pk, nil
}

func GenerateSecretKey(seed []byte) SecretKey {
	var sk SecretKey
	copy(sk[:], seed[:SecretKeySize])
	sk[0] &= 248
	sk[31] &= 127
	sk[31] |= 64
	return sk
}

// GeneratePublicKey generates a public key from a secret key.
func GeneratePublicKey(sk SecretKey) PublicKey {
	reader := bytes.NewReader(sk[:])
	publicKey, _, _ := GenerateWavesKey(reader)
	var pk PublicKey
	copy(pk[:], publicKey[:])
	return pk
}

func GenerateKeyPair(seed []byte) (SecretKey, PublicKey, error) {
	var sk SecretKey
	var pk PublicKey
	h := sha256.New()
	if _, err := h.Write(seed); err != nil {
		return sk, pk, err
	}
	digest := h.Sum(nil)
	sk = GenerateSecretKey(digest)
	pk = GeneratePublicKey(sk)
	return sk, pk, nil
}

func SecureHash(data []byte) (Digest, error) {
	var d Digest
	fh, err := blake2b.New256(nil)
	if err != nil {
		return d, err
	}
	if _, err := fh.Write(data); err != nil {
		return d, err
	}
	fh.Sum(d[:0])
	h := sha3.NewLegacyKeccak256()
	if _, err := h.Write(d[:DigestSize]); err != nil {
		return d, err
	}
	h.Sum(d[:0])
	return d, nil
}

func FastHash(data []byte) (Digest, error) {
	var d Digest
	h, err := blake2b.New256(nil)
	if err != nil {
		return d, err
	}
	if _, err := h.Write(data); err != nil {
		return d, err
	}
	h.Sum(d[:0])
	return d, nil
}

func Sign(secretKey SecretKey, data []byte) (Signature, error) {
	var sig Signature
	var A edwards25519.ExtendedGroupElement
	var hBytes [32]byte
	copy(hBytes[:], secretKey[:])
	var pk1 [32]byte
	A.ToBytes(&pk1)

	var wideBytes [64]byte
	copy(wideBytes[:], secretKey[:])
	wideBytes[0] &= 248
	wideBytes[31] &= 63
	wideBytes[31] |= 64
	var out [32]byte
	edwards25519.ScReduce(&out, &wideBytes)
	edwards25519.GeScalarMultBase(&A, &hBytes)
	var pkb [32]byte
	A.ToBytes(&pkb)

	sf := pkb[31] & 0x80

	random := make([]byte, sha512.Size)
	if _, err := rand.Read(random); err != nil {
		return sig, err
	}

	md := make([]byte, 0, sha512.Size)
	h := sha512.New()
	if _, err := h.Write(prefix); err != nil {
		return sig, err
	}
	if _, err := h.Write(out[:]); err != nil {
		return sig, err
	}
	if _, err := h.Write(data); err != nil {
		return sig, err
	}
	if _, err := h.Write(random); err != nil {
		return sig, err
	}
	md = h.Sum(md)

	var wideBytes2 [64]byte
	copy(wideBytes2[:], md[:])
	var out2 [32]byte
	edwards25519.ScReduce(&out2, &wideBytes2)

	var R edwards25519.ExtendedGroupElement
	edwards25519.GeScalarMultBase(&R, &out2)
	var Rb [32]byte
	R.ToBytes(&Rb)

	hd := make([]byte, 0, sha512.Size)
	h.Reset()
	if _, err := h.Write(Rb[:]); err != nil {
		return sig, err
	}
	if _, err := h.Write(pkb[:]); err != nil {
		return sig, err
	}
	if _, err := h.Write(data); err != nil {
		return sig, err
	}
	hd = h.Sum(hd)

	var wideBytes3 [64]byte
	copy(wideBytes3[:], hd[:])
	var out3 [32]byte
	edwards25519.ScReduce(&out3, &wideBytes3)

	var Sb [32]byte
	edwards25519.ScMulAdd(&Sb, &out3, &out, &out2)

	copy(sig[:DigestSize], Rb[:])
	copy(sig[DigestSize:], Sb[:])

	sig[63] &= 0x7f
	sig[63] |= sf
	return sig, nil
}
