package aptos

import "testing"

func TestAptosWallet_GetRandomPrivateKey(t *testing.T) {
	w := &AptosWallet{}
	t.Log(w.GetRandomPrivateKey())
}
