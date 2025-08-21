package crypto

import (
	"crypto/ed25519"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/internal/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testEd25519PrivateKey = "ed25519-priv-0xc5338cd251c22daa8c9c9cc94f498cc8a5c7e1d2e75287a5dda91096fe64efa5"
const testEd25519PrivateKeyHex = "0xc5338cd251c22daa8c9c9cc94f498cc8a5c7e1d2e75287a5dda91096fe64efa5"
const testEd25519PublicKey = "0xde19e5d1880cac87d57484ce9ed2e84cf0f9599f12e7cc3a52e4e7657a763f2c"
const testEd25519Address = "0x978c213990c4833df71548df7ce49d54c759d6b6d932de22b24d56060b7af2aa"
const testEd25519Message = "0x68656c6c6f20776f726c64"
const testEd25519Signature = "0x9e653d56a09247570bb174a389e85b9226abd5c403ea6c504b386626a145158cd4efd66fc5e071c0e19538a96a05ddbda24d3c51e1e6a9dacc6bb1ce775cce07"

func TestEd25519Keys(t *testing.T) {
	testEd25519PrivateKeyBytes := []byte{0xc5, 0x33, 0x8c, 0xd2, 0x51, 0xc2, 0x2d, 0xaa, 0x8c, 0x9c, 0x9c, 0xc9, 0x4f, 0x49, 0x8c, 0xc8, 0xa5, 0xc7, 0xe1, 0xd2, 0xe7, 0x52, 0x87, 0xa5, 0xdd, 0xa9, 0x10, 0x96, 0xfe, 0x64, 0xef, 0xa5}

	// First ensure bytes and hex are the same
	readBytes, err := util.ParseHex(testEd25519PrivateKeyHex)
	assert.NoError(t, err)
	assert.Equal(t, testEd25519PrivateKeyBytes, readBytes)

	// Either bytes or hex should work
	privateKey := &Ed25519PrivateKey{}
	err = privateKey.FromHex(testEd25519PrivateKey)
	assert.NoError(t, err)
	privateKey2 := &Ed25519PrivateKey{}
	err = privateKey2.FromBytes(testEd25519PrivateKeyBytes)
	assert.NoError(t, err)
	assert.Equal(t, privateKey, privateKey2)

	// The outputs should match as well
	assert.Equal(t, privateKey.Bytes(), testEd25519PrivateKeyBytes)
	assert.Equal(t, privateKey.ToHex(), testEd25519PrivateKeyHex)
	formattedString, err := privateKey.ToAIP80()
	assert.NoError(t, err)
	assert.Equal(t, formattedString, testEd25519PrivateKey)

	// Auth key should match
	assert.Equal(t, testEd25519Address, privateKey.AuthKey().ToHex())

	// Test signature
	message, err := util.ParseHex(testEd25519Message)
	assert.NoError(t, err)
	authenticator, err := privateKey.Sign(message)
	assert.NoError(t, err)

	// Check public keys
	publicKey := authenticator.PubKey()
	assert.Equal(t, testEd25519PublicKey, privateKey.PubKey().ToHex())
	assert.Equal(t, testEd25519PublicKey, publicKey.ToHex())

	// Check signature
	expectedSignature, err := util.ParseHex(testEd25519Signature)
	assert.NoError(t, err)
	assert.Equal(t, expectedSignature, authenticator.Signature().Bytes())

	// Verify signature with the key and the authenticator directly
	assert.True(t, authenticator.Verify(message))
	assert.True(t, publicKey.Verify(message, authenticator.Signature()))

	// Verify serialization of public key
	publicKeyBytes, err := bcs.Serialize(publicKey)
	assert.NoError(t, err)
	expectedPublicKeyBytes, err := util.ParseHex(testEd25519PublicKey)
	assert.NoError(t, err)
	// Need to prepend the length
	expectedBcsPublicKeyBytes := []byte{ed25519.PublicKeySize}
	expectedBcsPublicKeyBytes = append(expectedBcsPublicKeyBytes, expectedPublicKeyBytes[:]...)
	assert.Equal(t, expectedBcsPublicKeyBytes, publicKeyBytes)

	publicKey2 := &Ed25519PublicKey{}
	err = bcs.Deserialize(publicKey2, publicKeyBytes)
	assert.NoError(t, err)
	assert.Equal(t, publicKey, publicKey2)

	// Check from bytes and from hex
	publicKey3 := &Ed25519PublicKey{}
	err = publicKey3.FromHex(testEd25519PublicKey)
	assert.NoError(t, err)
	publicKey4 := &Ed25519PublicKey{}
	err = publicKey4.FromBytes(expectedPublicKeyBytes)
	assert.NoError(t, err)

	assert.Equal(t, publicKey, publicKey3)
	assert.Equal(t, publicKey, publicKey4)

	// Test serialization and deserialization of authenticator
	authenticatorBytes, err := bcs.Serialize(authenticator)
	assert.NoError(t, err)
	authenticator2 := &AccountAuthenticator{}
	err = bcs.Deserialize(authenticator2, authenticatorBytes)
	assert.NoError(t, err)
	assert.Equal(t, authenticator, authenticator2)
}

func TestEd25519PrivateKeyWrongLength(t *testing.T) {
	privateKey := &Ed25519PrivateKey{}
	err := privateKey.FromBytes([]byte{0x01})
	assert.Error(t, err)
}

func TestEd25519PublicKeyWrongLength(t *testing.T) {
	key := &Ed25519PublicKey{}
	err := key.FromBytes([]byte{0x01})
	assert.Error(t, err)
}

func TestEd25519SignatureWrongLength(t *testing.T) {
	sig := &Ed25519Signature{}
	err := sig.FromBytes([]byte{0x01})
	assert.Error(t, err)
}
