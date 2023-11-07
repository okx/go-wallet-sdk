package keypair

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"
)

var kp = New(Ed25519Version)

func TestKeypair_GenerateKey(t *testing.T) {
	priv, pub := kp.GenerateKey()
	t.Log(hex.EncodeToString(priv))
	addr := kp.CreateAddress(pub)
	t.Log("addr : ", addr)
	kPair := NewKeypairFromHex(Ed25519Version, hex.EncodeToString(priv))
	a := kPair.CreateAddressable()
	require.Equal(t, addr, a.GetAddress())
}

func TestKeypair_Sign(t *testing.T) {
	priv, _ := kp.GenerateKey()
	kp.SetPrivateKey(priv)
	signedMessage, err := kp.Sign([]byte("test ed25519 signature"))
	require.NoError(t, err)
	t.Log("signedMessage :", hex.EncodeToString(signedMessage))
}
