package crypto

import (
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiEd25519Keys(t *testing.T) {
	t.Parallel()
	key1, key2, publicKey := createMultiEd25519Key(t)

	message := []byte("hello world")

	signature := createMultiEd25519Signature(t, key1, key2, message)

	// Test verification of signature
	assert.True(t, publicKey.Verify(message, signature))

	// Test serialization / deserialization authenticator
	auth := &MultiEd25519Authenticator{
		PubKey: publicKey,
		Sig:    signature,
	}
	assert.True(t, auth.Verify(message))
}

func TestMultiEd25519KeySerialization(t *testing.T) {
	t.Parallel()
	key1, key2, publicKey := createMultiEd25519Key(t)

	// Test serialization / deserialization public key
	keyBytes, err := bcs.Serialize(publicKey)
	require.NoError(t, err)
	publicKeyDeserialized := &MultiEd25519PublicKey{}
	err = bcs.Deserialize(publicKeyDeserialized, keyBytes)
	require.NoError(t, err)
	assert.Equal(t, publicKey, publicKeyDeserialized)

	// Test serialization / deserialization signature
	signature := createMultiEd25519Signature(t, key1, key2, []byte("test message"))
	sigBytes, err := bcs.Serialize(signature)
	require.NoError(t, err)
	signatureDeserialized := &MultiEd25519Signature{}
	err = bcs.Deserialize(signatureDeserialized, sigBytes)
	require.NoError(t, err)
	assert.Equal(t, signature, signatureDeserialized)

	// Test serialization / deserialization authenticator
	auth := &AccountAuthenticator{
		Variant: AccountAuthenticatorMultiEd25519,
		Auth: &MultiEd25519Authenticator{
			PubKey: publicKey,
			Sig:    signature,
		},
	}
	authBytes, err := bcs.Serialize(auth)
	require.NoError(t, err)
	authDeserialized := &AccountAuthenticator{}
	err = bcs.Deserialize(authDeserialized, authBytes)
	require.NoError(t, err)
	assert.Equal(t, auth, authDeserialized)
}

func createMultiEd25519Key(t *testing.T) (
	*Ed25519PrivateKey,
	*Ed25519PrivateKey,
	*MultiEd25519PublicKey,
) {
	t.Helper()
	key1, err := GenerateEd25519PrivateKey()
	require.NoError(t, err)
	// TODO: Maybe we should have a typed function for the public keys
	pubkey1, ok := key1.PubKey().(*Ed25519PublicKey)
	require.True(t, ok)
	key2, err := GenerateEd25519PrivateKey()
	require.NoError(t, err)
	pubkey2, ok := key2.PubKey().(*Ed25519PublicKey)
	require.True(t, ok)

	publicKey := &MultiEd25519PublicKey{
		PubKeys:            []*Ed25519PublicKey{pubkey1, pubkey2},
		SignaturesRequired: 2,
	}

	return key1, key2, publicKey
}

func createMultiEd25519Signature(t *testing.T, key1 *Ed25519PrivateKey, key2 *Ed25519PrivateKey, message []byte) *MultiEd25519Signature {
	t.Helper()
	sig1, err := key1.SignMessage(message)
	require.NoError(t, err)
	sig2, err := key2.SignMessage(message)
	require.NoError(t, err)
	sig1Typed, ok := sig1.(*Ed25519Signature)
	require.True(t, ok)
	sig2Typed, ok := sig2.(*Ed25519Signature)
	require.True(t, ok)

	// TODO: This signature should be built easier, ergonomics to fix this late
	return &MultiEd25519Signature{
		Signatures: []*Ed25519Signature{
			sig1Typed,
			sig2Typed,
		},
		Bitmap: [4]byte([]byte("c0000000")),
	}
}
