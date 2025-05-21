package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromHex(t *testing.T) {
	key := Ed25519PrivateKey{}
	err := key.FromHex("5f7bf6af2e3d40d2a2a7e6d0f8fa8047abb40d98a41a23cd404ff11620a5209d91b59738acda2ac81183fe78d91f24b1765b0ef5b0146e858e44bc69a3c5d009")
	assert.NoError(t, err)
	assert.Equal(t, "5f7bf6af2e3d40d2a2a7e6d0f8fa8047abb40d98a41a23cd404ff11620a5209d91b59738acda2ac81183fe78d91f24b1765b0ef5b0146e858e44bc69a3c5d009", key.ToHex())
}
