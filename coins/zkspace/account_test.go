package zkspace

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetAddress(t *testing.T) {
	pri := "559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"
	addr, err := GetAddress(pri)
	require.NoError(t, err)
	require.Equal(t, "0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728", addr)
}

func TestGetPubKeyHash(t *testing.T) {
	pri := "559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"
	pub, err := GetPubKeyHash(pri, 1)
	require.NoError(t, err)
	require.Equal(t, "sync:89497052061f2e34e3c11f5afdb65df454c0d7b6", pub)
}
