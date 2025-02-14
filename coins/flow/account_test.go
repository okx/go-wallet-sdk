package flow

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	privKey, pubKey := GenerateKeyPair()
	assert.Equal(t, 64, len(privKey))
	assert.Equal(t, 128, len(pubKey))
}
func TestValidateAddress(t *testing.T) {
	assert.True(t, ValidateAddress("0xa8d1a60acba12a20"))
}
