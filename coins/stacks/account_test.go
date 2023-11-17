package stacks

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateAddress(t *testing.T) {
	priKey := "598d99970d04be67e8b41ddd5c5453487eeab5345ea1638c9a2849dee377f2a301"
	pubKey, err := GetPublicKey(priKey)
	require.NoError(t, err)
	address, err := GetAddressFromPublicKey(pubKey)
	require.NoError(t, err)
	expected := "SP1QCZZWWXT5CADKWGEPGG6F4RM0BDH3NTTNM86ZG"
	require.Equal(t, expected, address)
}

func TestPubKeyfromPrivKeyWithPKCompressed(t *testing.T) {
	// Compressed
	pub, err := GetPublicKey("598d99970d04be67e8b41ddd5c5453487eeab5345ea1638c9a2849dee377f2a301")
	require.NoError(t, err)
	require.Equal(t, "032e615bd2b300081af80d3b8449168c6c2d6ae9478ed1c820233f1ba6fef85eef", pub)
	// UnCompressed
	pub2, err := GetPublicKey("598d99970d04be67e8b41ddd5c5453487eeab5345ea1638c9a2849dee377f2a3")
	require.NoError(t, err)
	require.Equal(t, "042e615bd2b300081af80d3b8449168c6c2d6ae9478ed1c820233f1ba6fef85eef8f31549af2f43622d6397f135608f49242d3057830bfa74423443db7701e717f", pub2)
}

func TestValidAddress(t *testing.T) {
	address1 := "SP1A6RRGGQ5DJM9FWRPYQRPHPFBNN1VKPGRB02581"
	assert.Equal(t, ValidAddress(address1), true)

	address2 := "1A6RRGGQ5DJM9FWRPYQRPHPFBNN1VKPGRB02581"
	assert.Equal(t, ValidAddress(address2), false)
}

func TestNewAddress(t *testing.T) {
	priKey := "598d99970d04be67e8b41ddd5c5453487eeab5345ea1638c9a2849dee377f2a301"
	addr, err := NewAddress(priKey)
	require.NoError(t, err)
	assert.Equal(t, "SP1QCZZWWXT5CADKWGEPGG6F4RM0BDH3NTTNM86ZG", addr)
}
