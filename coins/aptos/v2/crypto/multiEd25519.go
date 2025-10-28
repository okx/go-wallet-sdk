package crypto

import (
	"crypto/ed25519"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/internal/util"
)

// region MultiEd25519PublicKey

// MultiEd25519PublicKey is the public key for off-chain multi-sig on Aptos with Ed25519 keys
//
// Implements:
//   - [VerifyingKey]
//   - [PublicKey]
//   - [CryptoMaterial]
//   - [bcs.Marshaler]
//   - [bcs.Unmarshaler]
//   - [bcs.Struct]
type MultiEd25519PublicKey struct {
	// PubKeys is the list of all public keys associated with the off-chain multi-sig
	PubKeys []*Ed25519PublicKey
	// SignaturesRequired is the number of signatures required to pass verification
	SignaturesRequired uint8
}

// region MultiEd25519PublicKey VerifyingKey implementation

// Verify verifies the signature against the message
//
// # This function will return true if the number of verified signatures is greater than or equal to the number of required signatures
//
// Implements:
//   - [VerifyingKey]
func (key *MultiEd25519PublicKey) Verify(msg []byte, signature Signature) bool {
	switch sig := signature.(type) {
	case *MultiEd25519Signature:
		verified := uint8(0)
		// TODO: Verify with bitmap, and check that this works properly
		for i, pubKey := range key.PubKeys {
			if pubKey.Verify(msg, sig.Signatures[i]) {
				verified++
			}
		}

		return verified >= key.SignaturesRequired
	default:
		return false
	}
}

// endregion

// region MultiEd25519PublicKey PublicKey implementation

// AuthKey converts the public key to an authentication key
//
// Implements:
//
//   - [PublicKey]
func (key *MultiEd25519PublicKey) AuthKey() *AuthenticationKey {
	out := &AuthenticationKey{}
	out.FromPublicKey(key)
	return out
}

// Scheme returns the scheme for the public key
//
// Implements:
//   - [PublicKey]
func (key *MultiEd25519PublicKey) Scheme() uint8 {
	return MultiEd25519Scheme
}

// endregion

// region MultiEd25519PublicKey CryptoMaterial implementation

// Bytes serializes the public key to bytes
//
// Implements:
//   - [CryptoMaterial]
func (key *MultiEd25519PublicKey) Bytes() []byte {
	keyBytes := make([]byte, len(key.PubKeys)*ed25519.PublicKeySize+1)
	for i, publicKey := range key.PubKeys {
		start := i * ed25519.PublicKeySize
		end := start + ed25519.PublicKeySize
		copy(keyBytes[start:end], publicKey.Bytes())
	}
	keyBytes[len(keyBytes)-1] = key.SignaturesRequired
	return keyBytes
}

// FromBytes deserializes the public key from bytes
//
// Returns an error if deserialization fails due to invalid keys.
//
// Implements:
//   - [CryptoMaterial]
func (key *MultiEd25519PublicKey) FromBytes(bytes []byte) error {
	keyBytesLength := len(bytes)
	numKeys := keyBytesLength / ed25519.PublicKeySize
	signaturesRequired := bytes[keyBytesLength-1]

	pubKeys := make([]*Ed25519PublicKey, numKeys)
	for i := range numKeys {
		start := i * ed25519.PublicKeySize
		end := start + ed25519.PublicKeySize
		pubKeys[i] = &Ed25519PublicKey{}
		err := pubKeys[i].FromBytes(bytes[start:end])
		if err != nil {
			return fmt.Errorf("failed to deserialize multi ed25519 public key sub key %d: %w", i, err)
		}
	}

	key.SignaturesRequired = signaturesRequired
	key.PubKeys = pubKeys
	return nil
}

// ToHex serializes the public key to a hex string
//
// Implements:
//   - [CryptoMaterial]
func (key *MultiEd25519PublicKey) ToHex() string {
	return util.BytesToHex(key.Bytes())
}

// FromHex deserializes the public key from a hex string
//
// Returns an error if deserialization fails due to invalid keys.
//
// Implements:
//   - [CryptoMaterial]
func (key *MultiEd25519PublicKey) FromHex(hexStr string) error {
	bytes, err := util.ParseHex(hexStr)
	if err != nil {
		return err
	}
	return key.FromBytes(bytes)
}

// endregion

// region MultiEd25519PublicKey bcs.Struct implementation

// MarshalBCS serializes the public key to bytes
//
// Implements:
//   - [bcs.Marshaler]
func (key *MultiEd25519PublicKey) MarshalBCS(ser *bcs.Serializer) {
	ser.WriteBytes(key.Bytes())
}

// UnmarshalBCS deserializes the public key from bytes
//
// Returns an error if deserialization fails due to invalid keys or not enough bytes.
//
// Implements:
//   - [bcs.Unmarshaler]
func (key *MultiEd25519PublicKey) UnmarshalBCS(des *bcs.Deserializer) {
	keyBytes := des.ReadBytes()
	if des.Error() != nil {
		return
	}
	err := key.FromBytes(keyBytes)
	if err != nil {
		des.SetError(err)
	}
}

// endregion
// endregion

// region MultiEd25519Authenticator

// MultiEd25519Authenticator is an authenticator for a MultiEd25519Signature
//
// Implements:
//   - [AccountAuthenticatorImpl]
//   - [bcs.Marshaler]
//   - [bcs.Unmarshaler]
//   - [bcs.Struct]
type MultiEd25519Authenticator struct {
	PubKey *MultiEd25519PublicKey
	Sig    *MultiEd25519Signature
}

// region MultiEd25519Authenticator AccountAuthenticatorImpl implementation

// PublicKey returns the public key associated with the authenticator
//
// Implements:
//   - [AccountAuthenticatorImpl]
func (ea *MultiEd25519Authenticator) PublicKey() PublicKey {
	return ea.PubKey
}

// Signature returns the signature associated with the authenticator
//
// Implements:
//   - [AccountAuthenticatorImpl]
func (ea *MultiEd25519Authenticator) Signature() Signature {
	return ea.Sig
}

// Verify verifies the signature against the message
//
// Implements:
//   - [AccountAuthenticatorImpl]
func (ea *MultiEd25519Authenticator) Verify(msg []byte) bool {
	return ea.PubKey.Verify(msg, ea.Sig)
}

// endregion

// region MultiEd25519Authenticator bcs.Struct implementation

// MarshalBCS serializes the authenticator to bytes
//
// Implements:
//   - [bcs.Marshaler]
func (ea *MultiEd25519Authenticator) MarshalBCS(ser *bcs.Serializer) {
	ser.Struct(ea.PublicKey())
	ser.Struct(ea.Signature())
}

// UnmarshalBCS deserializes the authenticator from bytes
//
// Returns an error if deserialization fails due to invalid keys or not enough bytes.
//
// Implements:
//   - [bcs.Unmarshaler]
func (ea *MultiEd25519Authenticator) UnmarshalBCS(des *bcs.Deserializer) {
	ea.PubKey = &MultiEd25519PublicKey{}
	des.Struct(ea.PubKey)
	err := des.Error()
	if err != nil {
		return
	}
	ea.Sig = &MultiEd25519Signature{}
	des.Struct(ea.Sig)
}

// endregion
// endregion

// region MultiEd25519Signature

// MultiEd25519BitmapLen is number of bytes in the bitmap representing who signed the transaction
const MultiEd25519BitmapLen = 4

// MultiEd25519Signature is a signature for off-chain multi-sig
//
// Implements:
//   - [Signature]
//   - [CryptoMaterial]
//   - [bcs.Marshaler]
//   - [bcs.Unmarshaler]
//   - [bcs.Struct]
type MultiEd25519Signature struct {
	Signatures []*Ed25519Signature
	Bitmap     [MultiEd25519BitmapLen]byte
}

// region MultiEd25519Signature CryptoMaterial implementation

// Bytes serializes the signature to bytes
//
// Implements:
//   - [CryptoMaterial]
func (e *MultiEd25519Signature) Bytes() []byte {
	// This is a weird one, we need to serialize in set bytes
	sigBytes := make([]byte, len(e.Signatures)*ed25519.SignatureSize+MultiEd25519BitmapLen)
	for i, signature := range e.Signatures {
		start := i * ed25519.SignatureSize
		end := start + ed25519.SignatureSize
		copy(sigBytes[start:end], signature.Bytes())
	}
	copy(sigBytes[len(sigBytes)-MultiEd25519BitmapLen:], e.Bitmap[:])
	return sigBytes
}

// FromBytes deserializes the signature from bytes
//
// Returns an error if deserialization fails due to invalid keys or not enough bytes.
//
// Implements:
//   - [CryptoMaterial]
func (e *MultiEd25519Signature) FromBytes(bytes []byte) error {
	signatures := make([]*Ed25519Signature, len(bytes)/ed25519.SignatureSize)
	for i := 0; (i+1)*ed25519.SignatureSize < len(bytes); i++ {
		start := i * ed25519.SignatureSize
		end := start + ed25519.SignatureSize
		signatures[i] = &Ed25519Signature{}
		err := signatures[i].FromBytes(bytes[start:end])
		if err != nil {
			return fmt.Errorf("failed to deserialize multi ed25519 signature sub signature %d: %w", i, err)
		}
	}
	copy(e.Bitmap[:], bytes[len(bytes)-MultiEd25519BitmapLen:])
	e.Signatures = signatures
	return nil
}

// ToHex serializes the signature to a hex string
//
// Implements:
//   - [CryptoMaterial]
func (e *MultiEd25519Signature) ToHex() string {
	return util.BytesToHex(e.Bytes())
}

// FromHex deserializes the signature from a hex string
//
// Returns an error if deserialization fails due to invalid keys.
//
// Implements:
//   - [CryptoMaterial]
func (e *MultiEd25519Signature) FromHex(hexStr string) error {
	bytes, err := util.ParseHex(hexStr)
	if err != nil {
		return err
	}
	return e.FromBytes(bytes)
}

// endregion

// region MultiEd25519Signature bcs.Struct implementation

// MarshalBCS serializes the signature to bytes
//
// Implements:
//   - [bcs.Marshaler]
func (e *MultiEd25519Signature) MarshalBCS(ser *bcs.Serializer) {
	ser.WriteBytes(e.Bytes())
}

// UnmarshalBCS deserializes the signature from bytes
//
// Returns an error if deserialization fails due to invalid keys or not enough bytes.
//
// Implements
//   - [bcs.Unmarshaler]
func (e *MultiEd25519Signature) UnmarshalBCS(des *bcs.Deserializer) {
	bytes := des.ReadBytes()
	err := e.FromBytes(bytes)
	if err != nil {
		des.SetError(err)
	}
}

// endregion
// endregion
