package near

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"github.com/okx/go-wallet-sdk/crypto/base58"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAccount(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	t.Log(hex.EncodeToString(privateKey))
	if err != nil {
		t.Fatal(err)
	}
	address := base58.Encode(publicKey)
	t.Logf(address)
}

func TestValidateAddress(t *testing.T) {
	addr, _, err := NewAccount()
	if err != nil {
		t.Fatal(err)
	}
	isAddr := ValidateAddress(addr)
	t.Logf("%s is address? %t", addr, isAddr)
	assert.Equal(t, isAddr, true)
}
