package helium

import (
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	to               = "13ECKNq99BqN3dHhqXRYsdUHAPCEnfFBJsYVh5aqVSsYB35M3wS"
	from             = "13Lqwnbh427csevUveZF9n3ra1LnVYQug31RFeENaYgXuK2s8UC"
	amount    uint64 = 120
	fee       uint64 = 35000
	nonce     uint64 = 2
	private          = "f5e029dd6cca805047ca64e131c0a6cf3bf45c7ad03a7a1e7681963c9b1f3043"
	tokenType        = "hnt"
	isMax            = true
)

func Test_CreateAddress(t *testing.T) {
	address := NewAddress(private)
	require.Equal(t, address, from)
}

func Test_Sign(t *testing.T) {
	signTx, err := Sign(private, from, to, amount, fee, nonce, tokenType, isMax)
	if err != nil {
		// todo
	}
	require.NoError(t, err)
	expected := "wgGUAQohATRzO7mymsXF5mphcGit6S+VtjKx/IRuIrvpOSqEWSLMEicKIQElWnMrrLtN3iwWLGgC3fPx3D8hzAR7R/GzQTaseWoJxBB4IAEYuJECIAIqQPFldfkKANOAus8bNMsfiSBYysh+SZXoQHB2BlRX5sLaDB4V+awcrzu99dXo9Guq4gwZMpq1AYX8A4b5Qzq95Q8="
	require.Equal(t, expected, signTx)
}
