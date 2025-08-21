package crypto

import (
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/internal/util"
)

//region AuthenticationKey

// DeriveScheme is the key type for deriving the AuthenticationKey.  It is used in a SHA3-256 hash.
//
// Types:
//   - [Ed25519Scheme]
//   - [MultiEd25519Scheme]
//   - [SingleKeyScheme]
//   - [MultiKeyScheme]
//   - [DeriveObjectScheme]
//   - [NamedObjectScheme]
//   - [ResourceAccountScheme]
type DeriveScheme = uint8

// Seeds for deriving addresses from addresses
const (
	Ed25519Scheme         DeriveScheme = 0   // Ed25519Scheme is the default scheme for deriving the AuthenticationKey
	MultiEd25519Scheme    DeriveScheme = 1   // MultiEd25519Scheme is the scheme for deriving the AuthenticationKey for Multi-ed25519 accounts
	SingleKeyScheme       DeriveScheme = 2   // SingleKeyScheme is the scheme for deriving the AuthenticationKey for single-key accounts
	MultiKeyScheme        DeriveScheme = 3   // MultiKeyScheme is the scheme for deriving the AuthenticationKey for multi-key accounts
	DeriveObjectScheme    DeriveScheme = 252 // DeriveObjectScheme is the scheme for deriving the AuthenticationKey for objects, used to create new object addresses
	NamedObjectScheme     DeriveScheme = 254 // NamedObjectScheme is the scheme for deriving the AuthenticationKey for named objects, used to create new named object addresses
	ResourceAccountScheme DeriveScheme = 255 // ResourceAccountScheme is the scheme for deriving the AuthenticationKey for resource accounts, used to create new resource account addresses
)

// AuthenticationKeyLength is the length of a SHA3-256 Hash
const AuthenticationKeyLength = 32

// AuthenticationKey a hash representing the method for authorizing an account
//
// Implements:
//   - [CryptoMaterial]
//   - [bcs.Marshaler]
//   - [bcs.Unmarshaler]
//   - [bcs.Struct]
type AuthenticationKey [AuthenticationKeyLength]byte

// FromPublicKey for private / public key pairs, the [AuthenticationKey] is derived from the [PublicKey] directly
func (ak *AuthenticationKey) FromPublicKey(publicKey PublicKey) {
	ak.FromBytesAndScheme(publicKey.Bytes(), publicKey.Scheme())
}

// FromBytesAndScheme derives the [AuthenticationKey] directly from the SHA3-256 hash of the combined array
func (ak *AuthenticationKey) FromBytesAndScheme(bytes []byte, scheme DeriveScheme) {
	authBytes := util.Sha3256Hash([][]byte{
		bytes,
		{scheme},
	})
	copy((*ak)[:], authBytes)
}

//region AuthenticationKey CryptoMaterial

// Bytes returns the raw bytes of the [AuthenticationKey]
//
// Implements:
//   - [CryptoMaterial]
func (ak *AuthenticationKey) Bytes() []byte {
	return ak[:]
}

// FromBytes sets the [AuthenticationKey] to the given bytes
//
// Implements:
//   - [CryptoMaterial]
func (ak *AuthenticationKey) FromBytes(bytes []byte) (err error) {
	if len(bytes) != AuthenticationKeyLength {
		return fmt.Errorf("invalid authentication key, not 32 bytes")
	}
	copy((*ak)[:], bytes)
	return nil
}

// ToHex returns the hex string representation of the [AuthenticationKey], with a leading 0x
//
// Implements:
//   - [CryptoMaterial]
func (ak *AuthenticationKey) ToHex() string {
	return util.BytesToHex(ak[:])
}

// FromHex sets the [AuthenticationKey] to the bytes represented by the hex string, with or without a leading 0x
//
// Implements:
//   - [CryptoMaterial]
func (ak *AuthenticationKey) FromHex(hexStr string) (err error) {
	bytes, err := util.ParseHex(hexStr)
	if err != nil {
		return err
	}
	return ak.FromBytes(bytes)
}

//endregion

//region AuthenticationKey bcs.Struct

// MarshalBCS serializes the [AuthenticationKey] to BCS bytes
//
// Implements:
//   - [bcs.Marshaler]
func (ak *AuthenticationKey) MarshalBCS(ser *bcs.Serializer) {
	ser.Uleb128(AuthenticationKeyLength)
	ser.FixedBytes(ak[:])
}

// UnmarshalBCS deserializes the [AuthenticationKey] from BCS bytes
//
// Implements:
//   - [bcs.Unmarshaler]
func (ak *AuthenticationKey) UnmarshalBCS(des *bcs.Deserializer) {
	length := des.Uleb128()
	if length != AuthenticationKeyLength {
		des.SetError(fmt.Errorf("authentication key has wrong length %d", length))
		return
	}
	des.ReadFixedBytesInto(ak[:])
}

//endregion
//endregion
