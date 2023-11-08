package oasis

import (
	"testing"
)

func TestNewAddress(t *testing.T) {
	privateKeyHex := "d10a45ef8c019d22b7e8d18f77297677bff80ff4d2f23ab9ac14bdbac32c86e7baf40754ed3843e0464f814c3c605d8c36500cfb6892e2bd441839102f4200ed"
	address, err := NewAddress(privateKeyHex)
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Logf(address)
}
