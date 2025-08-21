package crypto

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/internal/util"
	"io"
)

// Ed25519PrivateKey represents an Ed25519Private key
//
// Implements:
//   - [Signer]
//   - [MessageSigner]
//   - [CryptoMaterial]
//   - [bcs.Marshaler]
//   - [bcs.Unmarshaler]
//   - [bcs.Struct]
type Ed25519PrivateKey struct {
	Inner ed25519.PrivateKey // Inner is the actual private key
}

// GenerateEd25519PrivateKey generates a random [Ed25519PrivateKey]
//
// An [io.Reader] can be provided for randomness, otherwise the default randomness source is from [ed25519.GenerateKey].
// The [io.Reader] must provide 32 bytes of input.
//
// Returns an error if the key generation fails.
func GenerateEd25519PrivateKey(rand ...io.Reader) (privateKey *Ed25519PrivateKey, err error) {
	var priv ed25519.PrivateKey
	if len(rand) > 0 {
		_, priv, err = ed25519.GenerateKey(rand[0])
	} else {
		_, priv, err = ed25519.GenerateKey(nil)
	}
	if err != nil {
		return nil, err
	}
	return &Ed25519PrivateKey{priv}, nil
}

//region Ed25519PrivateKey Signer Implementation

// Sign signs a message and returns an [AccountAuthenticator] with the [Ed25519Signature] and [Ed25519PublicKey]
//
// Never returns an error.
//
// Implements:
//   - [Signer]
func (key *Ed25519PrivateKey) Sign(msg []byte) (authenticator *AccountAuthenticator, err error) {
	// Can't error
	signature, _ := key.SignMessage(msg)
	publicKeyBytes := key.PubKey().Bytes()

	return &AccountAuthenticator{
		Variant: AccountAuthenticatorEd25519,
		Auth: &Ed25519Authenticator{
			PubKey: &Ed25519PublicKey{Inner: publicKeyBytes},
			Sig:    signature.(*Ed25519Signature),
		},
	}, nil
}

// SimulationAuthenticator creates a new [AccountAuthenticator] for simulation purposes
//
// Implements:
//   - [Signer]
func (key *Ed25519PrivateKey) SimulationAuthenticator() *AccountAuthenticator {
	return &AccountAuthenticator{
		Variant: AccountAuthenticatorEd25519,
		Auth: &Ed25519Authenticator{
			PubKey: key.PubKey().(*Ed25519PublicKey),
			Sig:    &Ed25519Signature{},
		},
	}
}

// PubKey returns the [Ed25519PublicKey] associated with the [Ed25519PrivateKey]
//
// Implements:
//   - [Signer]
func (key *Ed25519PrivateKey) PubKey() PublicKey {
	pubKey := key.Inner.Public()
	return &Ed25519PublicKey{
		pubKey.(ed25519.PublicKey),
	}
}

// AuthKey returns the [AuthenticationKey] associated with the [Ed25519PrivateKey] for a [Ed25519Scheme].
//
// Implements:
//   - [Signer]
func (key *Ed25519PrivateKey) AuthKey() *AuthenticationKey {
	out := &AuthenticationKey{}
	out.FromPublicKey(key.PubKey())
	return out
}

//endregion

//region Ed25519PrivateKey MessageSigner Implementation

// SignMessage signs a message and returns the raw [Signature] without a [VerifyingKey] for verification
//
// Never returns an error.
//
// Implements:
//   - [MessageSigner]
func (key *Ed25519PrivateKey) SignMessage(msg []byte) (sig Signature, err error) {
	sigBytes := ed25519.Sign(key.Inner, msg)
	return &Ed25519Signature{Inner: [64]byte(sigBytes)}, nil
}

// EmptySignature creates an empty signature for use in simulation
//
// Implements:
//   - [MessageSigner]
func (key *Ed25519PrivateKey) EmptySignature() Signature {
	return &Ed25519Signature{}
}

// VerifyingKey returns the [Ed25519PublicKey] associated with the [Ed25519PrivateKey]
//
// Implements:
//   - [MessageSigner]
func (key *Ed25519PrivateKey) VerifyingKey() VerifyingKey {
	return key.PubKey()
}

//endregion

//region Ed25519PrivateKey CryptoMaterial Implementation

// Bytes returns the raw bytes of the [Ed25519PrivateKey]
//
// Implements:
//   - [CryptoMaterial]
func (key *Ed25519PrivateKey) Bytes() []byte {
	return key.Inner.Seed()
}

// FromBytes sets the [Ed25519PrivateKey] to the given bytes
//
// Returns an error if the bytes length is not [ed25519.SeedSize].
//
// Implements:
//   - [CryptoMaterial]
func (key *Ed25519PrivateKey) FromBytes(bytes []byte) (err error) {
	bytes, err = ParsePrivateKey(bytes, PrivateKeyVariantEd25519, false)
	if err != nil {
		return err
	}
	if len(bytes) != ed25519.SeedSize {
		return fmt.Errorf("invalid ed25519 private key size %d", len(bytes))
	}
	key.Inner = ed25519.NewKeyFromSeed(bytes)
	return nil
}

// ToHex returns the hex string representation of the [Ed25519PrivateKey], with a leading 0x
//
// Implements:
//   - [CryptoMaterial]
func (key *Ed25519PrivateKey) ToHex() string {
	return util.BytesToHex(key.Bytes())
}

// ToAIP80 formats the private key to AIP-80 compliant string
func (key *Ed25519PrivateKey) ToAIP80() (formattedString string, err error) {
	return FormatPrivateKey(key.ToHex(), PrivateKeyVariantEd25519)
}

// FromHex sets the [Ed25519PrivateKey] to the bytes represented by the hex string, with or without a leading 0x
//
// Errors if the hex string is not valid, or if the bytes length is not [ed25519.SeedSize].
//
// Implements:
//   - [CryptoMaterial]
func (key *Ed25519PrivateKey) FromHex(hexStr string) (err error) {
	bytes, err := ParsePrivateKey(hexStr, PrivateKeyVariantEd25519)
	if err != nil {
		return err
	}
	return key.FromBytes(bytes)
}

//endregion

//endregion

//region Ed25519PublicKey

// Ed25519PublicKey is a Ed25519PublicKey which can be used to verify signatures
//
// Implements:
//   - [VerifyingKey]
//   - [PublicKey]
//   - [CryptoMaterial]
//   - [bcs.Marshaler]
//   - [bcs.Unmarshaler]
//   - [bcs.Struct]
type Ed25519PublicKey struct {
	Inner ed25519.PublicKey // Inner is the actual public key
}

//region Ed25519PublicKey VerifyingKey implementation

// Verify verifies a message with the public key and [Signature]
//
// Returns false if the signature is not [Ed25519Signature], or if the verification fails.
//
// Implements:
//   - [VerifyingKey]
func (key *Ed25519PublicKey) Verify(msg []byte, sig Signature) bool {
	switch sig := sig.(type) {
	case *Ed25519Signature:
		return ed25519.Verify(key.Inner, msg, sig.Bytes())
	default:
		return false
	}
}

//endregion

//region Ed25519PublicKey PublicKey implementation

// AuthKey returns the [AuthenticationKey] associated with the [Ed25519PublicKey] for a [Ed25519Scheme].
//
// Implements:
//   - [PublicKey]
func (key *Ed25519PublicKey) AuthKey() *AuthenticationKey {
	out := &AuthenticationKey{}
	out.FromPublicKey(key)
	return out
}

// Scheme returns the [Ed25519Scheme] for the [Ed25519PublicKey]
//
// Implements:
//   - [PublicKey]
func (key *Ed25519PublicKey) Scheme() uint8 {
	return Ed25519Scheme
}

//endregion

//region Ed25519PublicKey CryptoMaterial implementation

// Bytes returns the raw bytes of the [Ed25519PublicKey]
//
// Implements:
//   - [CryptoMaterial]
func (key *Ed25519PublicKey) Bytes() []byte {
	return key.Inner[:]
}

// FromBytes sets the [Ed25519PublicKey] to the given bytes
//
// Returns an error if the bytes length is not [ed25519.PublicKeySize].
//
// Implements:
//   - [CryptoMaterial]
func (key *Ed25519PublicKey) FromBytes(bytes []byte) (err error) {
	if len(bytes) != ed25519.PublicKeySize {
		return errors.New("invalid ed25519 public key size")
	}
	key.Inner = bytes
	return nil
}

// ToHex returns the hex string representation of the [Ed25519PublicKey] with a leading 0x
//
// Implements:
//   - [CryptoMaterial]
func (key *Ed25519PublicKey) ToHex() string {
	return util.BytesToHex(key.Bytes())
}

// FromHex sets the [Ed25519PublicKey] to the bytes represented by the hex string, with or without a leading 0x
//
// Errors if the hex string is not valid, or if the bytes length is not [ed25519.PublicKeySize].
//
// Implements:
//   - [CryptoMaterial]
func (key *Ed25519PublicKey) FromHex(hexStr string) (err error) {
	bytes, err := util.ParseHex(hexStr)
	if err != nil {
		return err
	}
	return key.FromBytes(bytes)
}

//endregion

//region Ed25519PublicKey bcs.Struct implementation

// MarshalBCS serializes the [Ed25519PublicKey] to BCS bytes
//
// Implements:
//   - [bcs.Marshaler]
func (key *Ed25519PublicKey) MarshalBCS(ser *bcs.Serializer) {
	ser.WriteBytes(key.Inner)
}

// UnmarshalBCS deserializes the [Ed25519PublicKey] from BCS bytes
//
// Sets [bcs.Deserializer.Error] if the bytes length is not [ed25519.PublicKeySize], or if it fails to read the required bytes.
//
// Implements:
//   - [bcs.Unmarshaler]
func (key *Ed25519PublicKey) UnmarshalBCS(des *bcs.Deserializer) {
	kb := des.ReadBytes()
	if des.Error() != nil {
		return
	}
	err := key.FromBytes(kb)
	if err != nil {
		des.SetError(err)
		return
	}
}

//endregion
//endregion

//region Ed25519Authenticator

// Ed25519Authenticator represents a verifiable signature with it's accompanied public key
//
// Implements:
//
//   - [AccountAuthenticatorImpl]
//   - [bcs.Marshaler]
//   - [bcs.Unmarshaler]
//   - [bcs.Struct]
type Ed25519Authenticator struct {
	PubKey *Ed25519PublicKey // PubKey is the public key
	Sig    *Ed25519Signature // Sig is the signature
}

//region Ed25519Authenticator AccountAuthenticatorImpl implementation

// PublicKey returns the [PublicKey] of the authenticator
//
// Implements:
//   - [AccountAuthenticatorImpl]
func (ea *Ed25519Authenticator) PublicKey() PublicKey {
	return ea.PubKey
}

// Signature returns the [Signature] of the authenticator
//
// Implements:
//   - [AccountAuthenticatorImpl]
func (ea *Ed25519Authenticator) Signature() Signature {
	return ea.Sig
}

// Verify returns true if the authenticator can be cryptographically verified
//
// Implements:
//   - [AccountAuthenticatorImpl]
func (ea *Ed25519Authenticator) Verify(msg []byte) bool {
	return ea.PubKey.Verify(msg, ea.Sig)
}

//endregion

//region Ed25519Authenticator bcs.Struct implementation

// MarshalBCS serializes the [Ed25519Authenticator] to BCS bytes
//
// Implements:
//   - [bcs.Marshaler]
func (ea *Ed25519Authenticator) MarshalBCS(ser *bcs.Serializer) {
	ser.Struct(ea.PublicKey())
	ser.Struct(ea.Signature())
}

// UnmarshalBCS deserializes the [Ed25519Authenticator] from BCS bytes
//
// Sets [bcs.Deserializer.Error] if it fails to read the required bytes.
//
// Implements:
//   - [bcs.Unmarshaler]
func (ea *Ed25519Authenticator) UnmarshalBCS(des *bcs.Deserializer) {
	ea.PubKey = &Ed25519PublicKey{}
	des.Struct(ea.PubKey)
	if des.Error() != nil {
		return
	}
	ea.Sig = &Ed25519Signature{}
	des.Struct(ea.Sig)
}

//endregion
//endregion

//region Ed25519Signature

// Ed25519Signature a wrapper for serialization of Ed25519 signatures
//
// Implements:
//   - [Signature]
//   - [CryptoMaterial]
//   - [bcs.Marshaler]
//   - [bcs.Unmarshaler]
//   - [bcs.Struct]
type Ed25519Signature struct {
	Inner [ed25519.SignatureSize]byte // Inner is the actual signature
}

//region Ed25519Signature CryptoMaterial implementation

// Bytes returns the raw bytes of the [Ed25519Signature]
//
// Implements:
//   - [CryptoMaterial]
func (e *Ed25519Signature) Bytes() []byte {
	return e.Inner[:]
}

// FromBytes sets the [Ed25519Signature] to the given bytes
//
// Returns an error if the bytes length is not [ed25519.SignatureSize].
//
// Implements:
//   - [CryptoMaterial]
func (e *Ed25519Signature) FromBytes(bytes []byte) (err error) {
	if len(bytes) != ed25519.SignatureSize {
		return errors.New("invalid ed25519 signature size")
	}
	copy(e.Inner[:], bytes)
	return nil
}

// ToHex returns the hex string representation of the [Ed25519Signature], with a leading 0x
//
// Implements:
//   - [CryptoMaterial]
func (e *Ed25519Signature) ToHex() string {
	return util.BytesToHex(e.Bytes())
}

// FromHex sets the [Ed25519Signature] to the bytes represented by the hex string, with or without a leading 0x
//
// Errors if the hex string is not valid, or if the bytes length is not [ed25519.SignatureSize].
//
// Implements:
//   - [CryptoMaterial]
func (e *Ed25519Signature) FromHex(hexStr string) (err error) {
	bytes, err := util.ParseHex(hexStr)
	if err != nil {
		return err
	}
	return e.FromBytes(bytes)
}

//endregion

//region Ed25519Signature bcs.Struct implementation

// MarshalBCS serializes the [Ed25519Signature] to BCS bytes
//
// Implements:
//   - [bcs.Marshaler]
func (e *Ed25519Signature) MarshalBCS(ser *bcs.Serializer) {
	ser.WriteBytes(e.Bytes())
}

// UnmarshalBCS deserializes the [Ed25519Signature] from BCS bytes
//
// Sets [bcs.Deserializer.Error] if it fails to read the required [ed25519.SignatureSize] bytes.
//
// Implements:
//   - [bcs.Unmarshaler]
func (e *Ed25519Signature) UnmarshalBCS(des *bcs.Deserializer) {
	bytes := des.ReadBytes()
	if des.Error() != nil {
		return
	}
	err := e.FromBytes(bytes)
	if err != nil {
		des.SetError(err)
	}
}
