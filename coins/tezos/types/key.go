/**
The MIT License (MIT) Copyright (c) 2021-2023 Blockwatch Data Inc.
*/

package types

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"golang.org/x/crypto/blake2b"
	"strings"
)

// PassphraseFunc is a callback used to obtain a passphrase for decrypting a private key
type PassphraseFunc func() ([]byte, error)

type KeyType byte

const (
	KeyTypeEd25519 KeyType = iota
	KeyTypeSecp256k1
	KeyTypeP256
	KeyTypeBls12_381
	KeyTypeInvalid
)

var (
	// ErrUnknownKeyType describes an error where a type for a
	// public key is undefined.
	ErrUnknownKeyType = errors.New("tezos: unknown key type")
	// ErrPassphrase is returned when a required passphrase is missing
	ErrPassphrase = errors.New("tezos: passphrase required")
)

func (t KeyType) AddressType() AddressType {
	switch t {
	case KeyTypeEd25519:
		return AddressTypeEd25519
	case KeyTypeSecp256k1:
		return AddressTypeSecp256k1
	case KeyTypeP256:
		return AddressTypeP256
	case KeyTypeBls12_381:
		return AddressTypeBls12_381
	default:
		return AddressTypeInvalid
	}
}

func (t KeyType) IsValid() bool {
	return t < KeyTypeInvalid
}

func (t KeyType) Curve() elliptic.Curve {
	switch t {
	case KeyTypeSecp256k1:
		return secp256k1.S256()
	case KeyTypeP256:
		return elliptic.P256()
	default:
		return nil
	}
}

func (t KeyType) SkHashType() HashType {
	switch t {
	case KeyTypeEd25519:
		return HashTypeSkEd25519
	case KeyTypeSecp256k1:
		return HashTypeSkSecp256k1
	case KeyTypeP256:
		return HashTypeSkP256
	case KeyTypeBls12_381:
		return HashTypeSkBls12_381
	default:
		return HashTypeInvalid
	}
}

func (t KeyType) PkHashType() HashType {
	switch t {
	case KeyTypeEd25519:
		return HashTypePkEd25519
	case KeyTypeSecp256k1:
		return HashTypePkSecp256k1
	case KeyTypeP256:
		return HashTypePkP256
	case KeyTypeBls12_381:
		return HashTypePkBls12_381
	default:
		return HashTypeInvalid
	}
}

func (t KeyType) SkPrefixBytes() []byte {
	switch t {
	case KeyTypeEd25519:
		return ED25519_SEED_ID
	case KeyTypeSecp256k1:
		return SECP256K1_SECRET_KEY_ID
	case KeyTypeP256:
		return P256_SECRET_KEY_ID
	case KeyTypeBls12_381:
		return BLS12_381_SECRET_KEY_ID
	default:
		return nil
	}
}

func (t KeyType) PkPrefixBytes() []byte {
	switch t {
	case KeyTypeEd25519:
		return ED25519_PUBLIC_KEY_ID
	case KeyTypeSecp256k1:
		return SECP256K1_PUBLIC_KEY_ID
	case KeyTypeP256:
		return P256_PUBLIC_KEY_ID
	case KeyTypeBls12_381:
		return BLS12_381_PUBLIC_KEY_ID
	default:
		return nil
	}
}

func (t KeyType) Tag() byte {
	switch t {
	case KeyTypeEd25519:
		return 0
	case KeyTypeSecp256k1:
		return 1
	case KeyTypeP256:
		return 2
	case KeyTypeBls12_381:
		return 3
	default:
		return 255
	}
}

// Key represents a public key on the Tezos blockchain.
type Key struct {
	Type KeyType
	Data []byte
}

func (k Key) Address() Address {
	return Address{
		Type: k.Type.AddressType(),
		Hash: k.Hash(),
	}
}

func (k Key) Hash() []byte {
	h, _ := blake2b.New(20, nil)
	h.Write(k.Data)
	return h.Sum(nil)
}

func (k Key) IsValid() bool {
	return k.Type.IsValid() && k.Type.PkHashType().Len() == len(k.Data)
}

func (k Key) String() string {
	if !k.IsValid() {
		return ""
	}
	return CheckEncode(k.Data, k.Type.PkPrefixBytes())
}

func (k Key) Bytes() []byte {
	if !k.Type.IsValid() {
		return nil
	}
	return append([]byte{k.Type.Tag()}, k.Data...)
}

func ParseKey(s string) (Key, error) {
	k := Key{}
	if len(s) == 0 {
		return k, nil
	}
	decoded, version, err := CheckDecode(s, 4, nil)
	if err != nil {
		if err == ErrChecksum {
			return k, ErrChecksum
		}
		return k, fmt.Errorf("tezos: unknown format for key %s: %w", s, err)
	}
	switch true {
	case bytes.Equal(version, ED25519_PUBLIC_KEY_ID):
		k.Type = KeyTypeEd25519
	case bytes.Equal(version, SECP256K1_PUBLIC_KEY_ID):
		k.Type = KeyTypeSecp256k1
	case bytes.Equal(version, P256_PUBLIC_KEY_ID):
		k.Type = KeyTypeP256
	case bytes.Equal(version, BLS12_381_PUBLIC_KEY_ID):
		k.Type = KeyTypeBls12_381
	default:
		return k, fmt.Errorf("tezos: unknown version %x for key %s", version, s)
	}
	if l := len(decoded); l != k.Type.PkHashType().Len() {
		return k, fmt.Errorf("tezos: invalid length %d for key data", l)
	}
	k.Data = decoded
	return k, nil
}

// PrivateKey represents a typed private key used for signing messages.
type PrivateKey struct {
	Type KeyType
	Data []byte
}

func (k PrivateKey) IsValid() bool {
	return k.Type.IsValid() && k.Type.SkHashType().Len() == len(k.Data)
}

func (k PrivateKey) String() string {
	var buf []byte
	switch k.Type {
	case KeyTypeEd25519:
		buf = ed25519.PrivateKey(k.Data).Seed()
	case KeyTypeSecp256k1, KeyTypeP256, KeyTypeBls12_381:
		buf = k.Data
	default:
		return ""
	}
	return CheckEncode(buf, k.Type.SkPrefixBytes())
}

func (k PrivateKey) Address() Address {
	return k.Public().Address()
}

func (k PrivateKey) MarshalText() ([]byte, error) {
	return []byte(k.String()), nil
}

func (k *PrivateKey) UnmarshalText(data []byte) error {
	key, err := ParsePrivateKey(string(data))
	if err != nil {
		return err
	}
	*k = key
	return nil
}

// GenerateKey creates a random private key.
func GenerateKey(typ KeyType) (PrivateKey, error) {
	key := PrivateKey{
		Type: typ,
	}
	switch typ {
	case KeyTypeEd25519:
		_, sk, err := ed25519.GenerateKey(nil)
		if err != nil {
			return key, err
		}
		key.Data = []byte(sk)
	case KeyTypeSecp256k1, KeyTypeP256:
		curve := typ.Curve()
		ecKey, err := ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			return key, err
		}
		key.Data = make([]byte, typ.SkHashType().Len())
		ecKey.D.FillBytes(key.Data)
	case KeyTypeBls12_381:
	}
	return key, nil
}

// Public returns the public key associated with the private key.
func (k PrivateKey) Public() Key {
	pk := Key{
		Type: k.Type,
	}
	switch k.Type {
	case KeyTypeEd25519:
		pk.Data = []byte(ed25519.PrivateKey(k.Data).Public().(ed25519.PublicKey))
	case KeyTypeSecp256k1, KeyTypeP256:
		curve := k.Type.Curve()
		ecKey, err := ecPrivateKeyFromBytes(k.Data, curve)
		if err != nil {
			pk.Type = KeyTypeInvalid
			return pk
		}
		pk.Data = elliptic.MarshalCompressed(curve, ecKey.PublicKey.X, ecKey.PublicKey.Y)
	case KeyTypeBls12_381:
	}
	return pk
}

func IsEncryptedKey(s string) bool {
	for _, prefix := range []string{
		ED25519_ENCRYPTED_SEED_PREFIX,
		SECP256K1_ENCRYPTED_SECRET_KEY_PREFIX,
		P256_ENCRYPTED_SECRET_KEY_PREFIX,
		BLS12_381_ENCRYPTED_SECRET_KEY_PREFIX,
	} {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

// Sign signs the digest (hash) of a message with the private key.
func (k PrivateKey) Sign(hash []byte) (Signature, error) {
	switch k.Type {
	case KeyTypeEd25519:
		return Signature{
			Type: SignatureTypeEd25519,
			Data: ed25519.Sign(ed25519.PrivateKey(k.Data), hash),
		}, nil
	case KeyTypeSecp256k1, KeyTypeP256:
		curve := k.Type.Curve()
		sig := Signature{
			Type: SignatureTypeSecp256k1,
		}
		if k.Type == KeyTypeP256 {
			sig.Type = SignatureTypeP256
		}
		ecKey, err := ecPrivateKeyFromBytes(k.Data, curve)
		if err != nil {
			return sig, err
		}
		sig.Data, err = ecSign(ecKey, hash)
		return sig, err
	case KeyTypeBls12_381:
		// TODO
		return Signature{}, ErrUnknownKeyType
	default:
		return Signature{}, ErrUnknownKeyType
	}
}

func ParsePrivateKey(s string) (PrivateKey, error) {
	return ParseEncryptedPrivateKey(s, nil)
}

// ParseEncryptedPrivateKey attempts to parse and optionally decrypt a Tezos private key.
func ParseEncryptedPrivateKey(s string, fn PassphraseFunc) (k PrivateKey, err error) {
	var (
		prefixLen     int = 4
		shouldDecrypt bool
	)
	if IsEncryptedKey(s) {
		prefixLen = 5
		shouldDecrypt = true
	}

	// decode base58, version length differs between encrypted and non-encrypted keys
	decoded, version, err := CheckDecode(s, prefixLen, nil)
	if err != nil {
		if err == ErrChecksum {
			err = ErrChecksumMismatch
			return
		}
		err = fmt.Errorf("tezos: unknown format for private key %s: %w", s, err)
		return
	}

	// decrypt if necessary
	if shouldDecrypt {
		decoded, err = decryptPrivateKey(decoded, fn)
		if err != nil {
			return
		}
		switch true {
		case bytes.Equal(version, ED25519_ENCRYPTED_SEED_ID):
			version = ED25519_SEED_ID
		case bytes.Equal(version, SECP256K1_ENCRYPTED_SECRET_KEY_ID):
			version = SECP256K1_SECRET_KEY_ID
		case bytes.Equal(version, P256_ENCRYPTED_SECRET_KEY_ID):
			version = P256_SECRET_KEY_ID
		case bytes.Equal(version, BLS12_381_ENCRYPTED_SECRET_KEY_ID):
			version = BLS12_381_SECRET_KEY_ID
		}
	}

	// detect type
	switch true {
	case bytes.Equal(version, ED25519_SEED_ID):
		if l := len(decoded); l != ed25519.SeedSize {
			return k, fmt.Errorf("tezos: invalid ed25519 seed length: %d", l)
		}
		k.Type = KeyTypeEd25519
		decoded = []byte(ed25519.NewKeyFromSeed(decoded))
	case bytes.Equal(version, ED25519_SECRET_KEY_ID):
		k.Type = KeyTypeEd25519
	case bytes.Equal(version, SECP256K1_SECRET_KEY_ID):
		k.Type = KeyTypeSecp256k1
	case bytes.Equal(version, P256_SECRET_KEY_ID):
		k.Type = KeyTypeP256
	case bytes.Equal(version, BLS12_381_SECRET_KEY_ID):
		k.Type = KeyTypeBls12_381
	default:
		err = fmt.Errorf("tezos: unknown version %x for private key %s", version, s)
		return
	}
	if l := len(decoded); l != k.Type.SkHashType().Len() {
		return k, fmt.Errorf("tezos: invalid length %d for private key data", l)
	}
	k.Data = decoded
	return
}
