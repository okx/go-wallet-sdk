package crypto

import (
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/internal/util"
)

//region SingleSigner

// SingleSigner is a wrapper around different types of MessageSigners to allow for many types of keys
//
// Implements:
//   - [Signer]
//   - [MessageSigner]
type SingleSigner struct {
	Signer MessageSigner // Signer is the actual signer or private key
}

// NewSingleSigner creates a new [SingleSigner] with the given [MessageSigner]
func NewSingleSigner(input MessageSigner) *SingleSigner {
	return &SingleSigner{Signer: input}
}

// SignMessage similar, but doesn't implement [MessageSigner] so there's no circular usage
func (key *SingleSigner) SignMessage(msg []byte) (Signature, error) {
	signature, err := key.Signer.SignMessage(msg)
	if err != nil {
		return nil, err
	}

	return &AnySignature{
		Variant:   key.SignatureVariant(),
		Signature: signature,
	}, nil
}

func (key *SingleSigner) SignatureVariant() AnySignatureVariant {
	sigType := AnySignatureVariantEd25519
	switch key.Signer.(type) {
	case *Ed25519PrivateKey:
		sigType = AnySignatureVariantEd25519
	}
	return sigType
}

func (key *SingleSigner) EmptySignature() *AnySignature {
	return &AnySignature{
		Variant:   key.SignatureVariant(),
		Signature: key.Signer.EmptySignature(),
	}
}

// region SingleSigner Signer implementation

// Sign signs a transaction and returns an associated [AccountAuthenticator]
//
// Implements:
//   - [Signer]
func (key *SingleSigner) Sign(msg []byte) (authenticator *AccountAuthenticator, err error) {
	signature, err := key.SignMessage(msg)
	if err != nil {
		return nil, err
	}

	auth := &SingleKeyAuthenticator{}
	auth.PubKey = key.PubKey().(*AnyPublicKey)
	auth.Sig = signature.(*AnySignature)
	return &AccountAuthenticator{Variant: AccountAuthenticatorSingleSender, Auth: auth}, nil
}

// SimulationAuthenticator creates a new [AccountAuthenticator] for simulation purposes
//
// Implements:
//   - [Signer]
func (key *SingleSigner) SimulationAuthenticator() *AccountAuthenticator {
	return &AccountAuthenticator{
		Variant: AccountAuthenticatorSingleSender,
		Auth: &SingleKeyAuthenticator{
			PubKey: key.PubKey().(*AnyPublicKey),
			Sig:    key.EmptySignature(),
		},
	}
}

// AuthKey gives the [AuthenticationKey] associated with the [Signer]
//
// Implements:
//   - [Signer]
func (key *SingleSigner) AuthKey() *AuthenticationKey {
	out := &AuthenticationKey{}
	out.FromPublicKey(key.PubKey())
	return out
}

// PubKey Retrieve the [PublicKey] for [Signature] verification
//
// Implements:
//   - [Signer]
func (key *SingleSigner) PubKey() PublicKey {
	innerPubKey := key.Signer.VerifyingKey()
	keyType := AnyPublicKeyVariantEd25519
	switch key.Signer.(type) {
	case *Ed25519PrivateKey:
		keyType = AnyPublicKeyVariantEd25519
	}
	return &AnyPublicKey{
		Variant: keyType,
		PubKey:  innerPubKey,
	}
}

//endregion
//endregion

//region AnyPublicKey

// AnyPublicKeyVariant is an enum ID for the public key used in AnyPublicKey
type AnyPublicKeyVariant uint32

const (
	AnyPublicKeyVariantEd25519   AnyPublicKeyVariant = 0 // AnyPublicKeyVariantEd25519 is the variant for [Ed25519PublicKey]
	AnyPublicKeyVariantSecp256k1 AnyPublicKeyVariant = 1 // AnyPublicKeyVariantSecp256k1 is the variant for [Secp256k1PublicKey]
)

// AnyPublicKey is used by SingleSigner and MultiKey to allow for using different keys with the same structs
// Implements [VerifyingKey], [PublicKey], [CryptoMaterial], [bcs.Marshaler], [bcs.Unmarshaler], [bcs.Struct]
type AnyPublicKey struct {
	Variant AnyPublicKeyVariant // Variant is the type of public key
	PubKey  VerifyingKey        // PubKey is the actual public key
}

// ToAnyPublicKey converts a [VerifyingKey] to an [AnyPublicKey]
func ToAnyPublicKey(key VerifyingKey) (*AnyPublicKey, error) {
	out := &AnyPublicKey{}
	switch key.(type) {
	case *Ed25519PublicKey:
		out.Variant = AnyPublicKeyVariantEd25519
	case *AnyPublicKey:
		// Passthrough for conversion
		return key.(*AnyPublicKey), nil
	default:
		return nil, fmt.Errorf("unknown public key type: %T", key)
	}
	out.PubKey = key
	return out, nil
}

//region AnyPublicKey VerifyingKey implementation

// Verify verifies the signature against the message
//
// Implements:
//   - [VerifyingKey]
func (key *AnyPublicKey) Verify(msg []byte, sig Signature) bool {
	switch sig := sig.(type) {
	case *AnySignature:
		return key.PubKey.Verify(msg, sig.Signature)
	default:
		return false
	}
}

//endregion

//region AnyPublicKey PublicKey implementation

// AuthKey converts the public key to an authentication key
//
// Implements:
//   - [PublicKey]
func (key *AnyPublicKey) AuthKey() *AuthenticationKey {
	out := &AuthenticationKey{}
	out.FromPublicKey(key)
	return out
}

// Scheme returns the scheme for the public key
//
// Implements:
//   - [PublicKey]
func (key *AnyPublicKey) Scheme() uint8 {
	return SingleKeyScheme
}

//endregion

//region AnyPublicKey CryptoMaterial implementation

// Bytes returns the raw bytes of the [AnyPublicKey]
//
// Implements:
//   - [CryptoMaterial]
func (key *AnyPublicKey) Bytes() []byte {
	val, _ := bcs.Serialize(key)
	return val
}

// FromBytes sets the [AnyPublicKey] to the given bytes
//
// Implements:
//   - [CryptoMaterial]
func (key *AnyPublicKey) FromBytes(bytes []byte) (err error) {
	return bcs.Deserialize(key, bytes)
}

// ToHex returns the hex string representation of the [AnyPublicKey], with a leading 0x
//
// Implements:
//   - [CryptoMaterial]
func (key *AnyPublicKey) ToHex() string {
	return util.BytesToHex(key.Bytes())
}

// FromHex sets the [AnyPublicKey] to the bytes represented by the hex string, with or without a leading 0x
//
// Implements:
//   - [CryptoMaterial]
func (key *AnyPublicKey) FromHex(hexStr string) (err error) {
	bytes, err := util.ParseHex(hexStr)
	if err != nil {
		return err
	}
	return key.FromBytes(bytes)
}

//endregion

//region AnyPublicKey bcs.Struct implementation

// MarshalBCS serializes the [AnyPublicKey] to bytes
//
// Implements:
//   - [bcs.Marshaler]
func (key *AnyPublicKey) MarshalBCS(ser *bcs.Serializer) {
	ser.Uleb128(uint32(key.Variant))
	ser.Struct(key.PubKey)
}

// UnmarshalBCS deserializes the [AnyPublicKey] from bytes
//
// Implements:
//   - [bcs.Unmarshaler]
func (key *AnyPublicKey) UnmarshalBCS(des *bcs.Deserializer) {
	key.Variant = AnyPublicKeyVariant(des.Uleb128())
	switch key.Variant {
	case AnyPublicKeyVariantEd25519:
		key.PubKey = &Ed25519PublicKey{}
	default:
		des.SetError(fmt.Errorf("unknown public key variant: %d", key.Variant))
		return
	}
	des.Struct(key.PubKey)
}

//endregion
//endregion

//region AnySignature

// AnySignatureVariant is an enum ID for the signature used in AnySignature
type AnySignatureVariant uint32

const (
	AnySignatureVariantEd25519   AnySignatureVariant = 0 // AnySignatureVariantEd25519 is the variant for [Ed25519Signature]
	AnySignatureVariantSecp256k1 AnySignatureVariant = 1 // AnySignatureVariantSecp256k1 is the variant for [Secp256k1Signature]
)

// AnySignature is a wrapper around signatures signed with SingleSigner and verified with AnyPublicKey
//
// Implements:
//   - [Signature]
//   - [CryptoMaterial]
//   - [bcs.Marshaler]
//   - [bcs.Unmarshaler]
//   - [bcs.Struct]
type AnySignature struct {
	Variant   AnySignatureVariant
	Signature Signature
}

// region AnySignature CryptoMaterial implementation

// Bytes returns the raw bytes of the [AnySignature]
//
// Implements:
//   - [CryptoMaterial]
func (e *AnySignature) Bytes() []byte {
	val, _ := bcs.Serialize(e)
	return val
}

// FromBytes sets the [AnySignature] to the given bytes
//
// Implements:
//   - [CryptoMaterial]
func (e *AnySignature) FromBytes(bytes []byte) (err error) {
	return bcs.Deserialize(e, bytes)
}

// ToHex returns the hex string representation of the [AnySignature], with a leading 0x
//
// Implements:
//   - [CryptoMaterial]
func (e *AnySignature) ToHex() string {
	return util.BytesToHex(e.Bytes())
}

// FromHex sets the [AnySignature] to the bytes represented by the hex string, with or without a leading 0x
//
// Implements:
//   - [CryptoMaterial]
func (e *AnySignature) FromHex(hexStr string) (err error) {
	bytes, err := util.ParseHex(hexStr)
	if err != nil {
		return err
	}
	return e.FromBytes(bytes)
}

//endregion

//region AnySignature bcs.Struct implementation

// MarshalBCS serializes the [AnySignature] to bytes
//
// Implements:
//   - [bcs.Marshaler]
func (e *AnySignature) MarshalBCS(ser *bcs.Serializer) {
	ser.Uleb128(uint32(e.Variant))
	ser.Struct(e.Signature)
}

// UnmarshalBCS deserializes the [AnySignature] from bytes
//
// Implements:
//   - [bcs.Unmarshaler]
func (e *AnySignature) UnmarshalBCS(des *bcs.Deserializer) {
	e.Variant = AnySignatureVariant(des.Uleb128())
	switch e.Variant {
	case AnySignatureVariantEd25519:
		e.Signature = &Ed25519Signature{}
	default:
		des.SetError(fmt.Errorf("unknown signature variant: %d", e.Variant))
		return
	}
	des.Struct(e.Signature)
}

//endregion
//endregion

//region SingleKeyAuthenticator

// SingleKeyAuthenticator is an authenticator for a [SingleSigner]
//
// Implements:
//   - [AccountAuthenticatorImpl]
//   - [bcs.Marshaler]
//   - [bcs.Unmarshaler]
//   - [bcs.Struct]
type SingleKeyAuthenticator struct {
	PubKey *AnyPublicKey
	Sig    *AnySignature
}

//region SingleKeyAuthenticator AccountAuthenticatorImpl implementation

// PublicKey returns the public key of the authenticator
//
// Implements:
//   - [AccountAuthenticatorImpl]
func (ea *SingleKeyAuthenticator) PublicKey() PublicKey {
	return ea.PubKey
}

// Signature returns the signature of the authenticator
//
// Implements:
//   - [AccountAuthenticatorImpl]
func (ea *SingleKeyAuthenticator) Signature() Signature {
	return ea.Sig
}

// Verify verifies the signature against the message
//
// Implements:
//   - [AccountAuthenticatorImpl]
func (ea *SingleKeyAuthenticator) Verify(msg []byte) bool {
	return ea.PubKey.Verify(msg, ea.Sig)
}

//endregion

//region SingleKeyAuthenticator bcs.Struct implementation

// MarshalBCS serializes the authenticator to bytes
//
// Implements:
//   - [bcs.Marshaler]
func (ea *SingleKeyAuthenticator) MarshalBCS(ser *bcs.Serializer) {
	ser.Struct(ea.PublicKey())
	ser.Struct(ea.Signature())
}

// UnmarshalBCS deserializes the authenticator from bytes
//
// Implements:
//   - [bcs.Unmarshaler]
func (ea *SingleKeyAuthenticator) UnmarshalBCS(des *bcs.Deserializer) {
	ea.PubKey = &AnyPublicKey{}
	des.Struct(ea.PubKey)
	err := des.Error()
	if err != nil {
		return
	}
	ea.Sig = &AnySignature{}
	des.Struct(ea.Sig)
}

//endregion
//endregion
