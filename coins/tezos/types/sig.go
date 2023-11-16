/**
The MIT License (MIT) Copyright (c) 2021-2023 Blockwatch Data Inc.
*/

package types

// SignatureType represents the type of a Tezos signature.
type SignatureType byte

const (
	SignatureTypeEd25519 SignatureType = iota
	SignatureTypeSecp256k1
	SignatureTypeP256
	SignatureTypeGeneric
	SignatureTypeInvalid
)

func (t SignatureType) Len() int {
	if t.IsValid() {
		return 64
	}
	return 0
}

func (t SignatureType) IsValid() bool {
	return t < SignatureTypeInvalid
}

func (t SignatureType) PrefixBytes() []byte {
	switch t {
	case SignatureTypeEd25519:
		return ED25519_SIGNATURE_ID
	case SignatureTypeSecp256k1:
		return SECP256K1_SIGNATURE_ID
	case SignatureTypeP256:
		return P256_SIGNATURE_ID
	case SignatureTypeGeneric:
		return GENERIC_SIGNATURE_ID
	default:
		return nil
	}
}

// Signature represents a typed Tezos signature.
type Signature struct {
	Type SignatureType
	Data []byte
}

func (s Signature) IsValid() bool {
	return s.Type.IsValid() && s.Type.Len() == len(s.Data)
}

func (s Signature) String() string {
	if !s.IsValid() {
		return ""
	}
	return CheckEncode(s.Data, s.Type.PrefixBytes())
}

func (s Signature) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}
