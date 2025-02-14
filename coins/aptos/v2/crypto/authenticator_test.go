package crypto

import (
	"crypto/ed25519"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/internal/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuthenticationKey_FromPublicKey(t *testing.T) {
	// Ed25519
	_, publicKey, err := GenerateEd5519Keys()
	assert.NoError(t, err)

	authKey := AuthenticationKey{}
	authKey.FromPublicKey(&publicKey)

	hash := util.Sha3256Hash([][]byte{
		publicKey.Bytes(),
		{Ed25519Scheme},
	})

	assert.Equal(t, hash[:], authKey[:])
}

//func Test_AuthenticationKeySerialization(t *testing.T) {
//	bytesWithLength := []byte{
//		32,
//		0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef,
//		0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef,
//		0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef,
//		0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef,
//	}
//	bytes := []byte{
//		0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef,
//		0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef,
//		0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef,
//		0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef,
//	}
//	authKey := AuthenticationKey(bytes)
//	serialized, err := bcs.Serialize(&authKey)
//	assert.NoError(t, err)
//	assert.Equal(t, bytesWithLength, serialized)
//
//	newAuthKey := AuthenticationKey{}
//	err = bcs.Deserialize(&newAuthKey, serialized)
//	assert.NoError(t, err)
//	assert.Equal(t, authKey, newAuthKey)
//}

func Test_AuthenticatorSerialization(t *testing.T) {
	msg := []byte{0x01, 0x02}
	privateKey, _, err := GenerateEd5519Keys()
	assert.NoError(t, err)

	authenticator, err := privateKey.Sign(msg)
	assert.NoError(t, err)

	serialized, err := bcs.Serialize(&authenticator)
	assert.NoError(t, err)
	assert.Equal(t, uint8(AuthenticatorEd25519), serialized[0])
	assert.Len(t, serialized, 1+(1+ed25519.PublicKeySize)+(1+ed25519.SignatureSize))

	newAuthenticator := Authenticator{}
	err = bcs.Deserialize(&newAuthenticator, serialized)
	assert.NoError(t, err)
	assert.Equal(t, authenticator, newAuthenticator)
}

func Test_AuthenticatorVerification(t *testing.T) {
	msg := []byte{0x01, 0x02}
	privateKey, _, err := GenerateEd5519Keys()
	assert.NoError(t, err)

	authenticator, err := privateKey.Sign(msg)
	assert.NoError(t, err)

	assert.True(t, authenticator.Verify(msg))
}

func Test_InvalidAuthenticatorDeserialization(t *testing.T) {
	serialized := []byte{0xFF}
	newAuthenticator := Authenticator{}
	err := bcs.Deserialize(&newAuthenticator, serialized)
	assert.Error(t, err)
	serialized = []byte{0x4F}
	err = bcs.Deserialize(&newAuthenticator, serialized)
	assert.Error(t, err)
}
func Test_InvalidAuthenticationKeyDeserialization(t *testing.T) {
	serialized := []byte{0xFF}
	newAuthkey := AuthenticationKey{}
	err := bcs.Deserialize(&newAuthkey, serialized)
	assert.Error(t, err)
}
