package zil

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetAddressFromPrivateKey(t *testing.T) {
	privHex := "c0dc46b9f9d6ef1c88dff3f1e82adc61cb11d77ab76a8d66338f14c2711cb4d8"
	address, err := GetAddressFromPrivateKey(privHex)
	require.NoError(t, err)
	require.Equal(t, "zil1uxcaatglxsgm9rluagx0yc65cuzj7nmhaw7l4u", address)
}

func TestGetPublicKeyFromPrivateKey(t *testing.T) {
	privKeyHex := "c0dc46b9f9d6ef1c88dff3f1e82adc61cb11d77ab76a8d66338f14c2711cb4d8"
	publicKey, err := GetPublicKeyFromPrivateKey(privKeyHex)
	require.NoError(t, err)
	require.Equal(t, "0374d4a96ac2f5f87ff6cb1687d330badea4988615ade53100653c85d69b1f40e9", publicKey)

}

func TestFromBech32Addr(t *testing.T) {
	addr, err := FromBech32Addr("zil1h6j9d76cp997r3lenwmdzkzdemry9v9su5ddz8")
	require.NoError(t, err)
	require.Equal(t, "bea456fb58094be1c7f99bb6d1584dcec642b0b0", addr)
}

func TestToBech32Address(t *testing.T) {
	address, err := ToBech32Address("4BAF5faDA8e5Db92C3d3242618c5B47133AE003C")
	require.NoError(t, err)
	require.Equal(t, "zil1fwh4ltdguhde9s7nysnp33d5wye6uqpugufkz7", address)
}
