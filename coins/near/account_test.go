package near

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAccount(t *testing.T) {
	addr, pub, prv, err := NewAccount()
	t.Logf("addr: %s\npub: %s\nprev: %s\nerr: %v", addr, pub, prv, err)
	assert.NoError(t, err)
	addr2, err := PublicKeyToAddress(pub)
	assert.NoError(t, err)
	assert.Equal(t, addr, addr2)
	addr2, err = PrivateKeyToAddr(prv)
	assert.NoError(t, err)
	assert.Equal(t, addr, addr2)
	pub2, err := PrivateKeyToPublicKey(prv)
	assert.NoError(t, err)
	assert.Equal(t, pub, pub2)
	isAddr := ValidateAddress(addr)
	assert.True(t, isAddr)
}

func TestValidateAddress(t *testing.T) {
	addr, _, _, err := NewAccount()
	assert.NoError(t, err)
	isAddr := ValidateAddress(addr)
	assert.True(t, isAddr)
	assert.True(t, ValidateAddress("58064be4ab6a0097b6c794f5cf1983ef36c60ea82c17e8488107433f6386b5ba"))
	assert.True(t, ValidateAddress("abc"))
	assert.True(t, ValidateAddress("123"))
	assert.True(t, ValidateAddress("abc.near"))
}
