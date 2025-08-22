package types

import (
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	defaultMetadata = "0x2ebb2ccac5e027a87fa0e2e5f656a3a4238d6a48d93ec9b610d570fc0aa0df12"
	defaultStore    = "0x8a9d57692a9d4deb1680eaf107b83c152436e10f7bb521143fa403fa95ef76a"
	defaultOwner    = "0xc67545d6f3d36ed01efc9b28cbfd0c1ae326d5d262dd077a29539bcee0edce9e"
)

func TestGenerateEd25519Account(t *testing.T) {
	message := []byte{0x12, 0x34}
	account, err := NewEd25519Account()
	assert.NoError(t, err)
	output, err := account.Sign(message)
	assert.NoError(t, err)
	assert.Equal(t, crypto.AccountAuthenticatorEd25519, output.Variant)
	assert.True(t, output.Auth.Verify(message))
}

func TestNewAccountFromSigner(t *testing.T) {
	message := []byte{0x12, 0x34}
	key, err := crypto.GenerateEd25519PrivateKey()
	assert.NoError(t, err)

	account, err := NewAccountFromSigner(key)
	assert.NoError(t, err)
	output, err := account.Sign(message)
	assert.NoError(t, err)
	assert.Equal(t, crypto.AccountAuthenticatorEd25519, output.Variant)
	assert.True(t, output.Auth.Verify(message))

	authKey := key.AuthKey()
	assert.Equal(t, authKey[:], account.Address[:])
}

func TestNewAccountFromSignerWithAddress(t *testing.T) {
	message := []byte{0x12, 0x34}
	key, err := crypto.GenerateEd25519PrivateKey()
	assert.NoError(t, err)

	authenticationKey := crypto.AuthenticationKey{}

	account, err := NewAccountFromSigner(key, authenticationKey)
	assert.NoError(t, err)
	output, err := account.Sign(message)
	assert.NoError(t, err)
	assert.Equal(t, crypto.AccountAuthenticatorEd25519, output.Variant)
	assert.True(t, output.Auth.Verify(message))

	assert.Equal(t, AccountZero, account.Address)
}

func TestNewAccountFromSignerWithAddressMulti(t *testing.T) {
	key, err := crypto.GenerateEd25519PrivateKey()
	assert.NoError(t, err)

	authenticationKey := crypto.AuthenticationKey{}

	_, err = NewAccountFromSigner(key, authenticationKey, authenticationKey)
	assert.Error(t, err)
}
