package stacks

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateAddress(t *testing.T) {
	privatekey := "598d99970d04be67e8b41ddd5c5453487eeab5345ea1638c9a2849dee377f2a301"
	pubkey, err := GetPublicKey(privatekey)
	if err != nil {
		t.Fatal(err)
	}
	address, err := GetAddressFromPublicKey(pubkey)
	if err != nil {
		t.Fatal(err)
	}

	expected := "SP1QCZZWWXT5CADKWGEPGG6F4RM0BDH3NTTNM86ZG"
	assert.Equal(t, expected, address)
}

func TestPubKeyfromPrivKeyWithPKCompressed(t *testing.T) {
	// Compressed
	pub, err := GetPublicKey("598d99970d04be67e8b41ddd5c5453487eeab5345ea1638c9a2849dee377f2a301")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, pub, "032e615bd2b300081af80d3b8449168c6c2d6ae9478ed1c820233f1ba6fef85eef")

	// UnCompressed
	pub2, err := GetPublicKey("598d99970d04be67e8b41ddd5c5453487eeab5345ea1638c9a2849dee377f2a3")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, pub2, "042e615bd2b300081af80d3b8449168c6c2d6ae9478ed1c820233f1ba6fef85eef8f31549af2f43622d6397f135608f49242d3057830bfa74423443db7701e717f")
}

func TestValidAddress(t *testing.T) {
	address1 := "SP1A6RRGGQ5DJM9FWRPYQRPHPFBNN1VKPGRB02581"
	assert.Equal(t, ValidAddress(address1), true)

	address2 := "1A6RRGGQ5DJM9FWRPYQRPHPFBNN1VKPGRB02581"
	assert.Equal(t, ValidAddress(address2), false)
}

func TestNewAddress(t *testing.T) {
	pri1 := "598d99970d04be67e8b41ddd5c5453487eeab5345ea1638c9a2849dee377f2a301"
	address1, err := NewAddress(pri1)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(address1)
	assert.Equal(t, address1, "SP1QCZZWWXT5CADKWGEPGG6F4RM0BDH3NTTNM86ZG")
}
