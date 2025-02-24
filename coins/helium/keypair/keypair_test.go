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

func TestKeypairSign(t *testing.T) {
	priv, err := hex.DecodeString("8b4b7063a3722dfb301739569feb1532834887d014b26b3b77729dff1f2c77e9")
	require.NoError(t, err)
	kp.SetPrivateKey(priv)
	signature, err := kp.Sign([]byte("test ed25519 signature"))
	require.NoError(t, err)
	require.Equal(t, "ef156333ca010c9789aba37f838e5dc944d6bad3da01d874c95042bb7858f12281e8fa2e79276698ef193208809dc8ab64881f6ebdde7c0993ab47ba2792b304", hex.EncodeToString(signature))
}
