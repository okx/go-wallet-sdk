package oasis

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewAddress(t *testing.T) {
	privateKeyHex := "a30a45ef8c019d22b7e8d18f11197677bff80ff4d2f23ab9ac14bdbac32c86e7baf40754ed3843e0464f814c3c605d8c36500cfb6892e2bd441839102f4200ed"
	address, err := NewAddress(privateKeyHex)
	require.NoError(t, err)
	require.Equal(t, "oasis1qzqrq9m2m7yfhpjk2x8ga2z3gg4fzcq2eqrexz40", address)
}
