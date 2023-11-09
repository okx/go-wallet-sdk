package crypto

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewPublicKeyFromBytes(t *testing.T) {
	pkStr := "J3pmMgPHKhaTdi74UENsEXfmetxjCGkYdqWW3rphowYa"
	pk1, err := NewPublicKeyFromBase58(pkStr)
	require.NoError(t, err)
	pk2, err := NewPublicKeyFromBytes(pk1.Bytes())
	require.NoError(t, err)
	require.Equal(t, pkStr, pk2.String())
}

func TestSign(t *testing.T) {
	var privateKey = "6QhEoSnJ12QDgeEAt3HYkPDBiYe15BArgSKWrV3DUctG"
	sk, err := NewSecretKeyFromBase58(privateKey)
	require.NoError(t, err)
	data := []byte("hello")
	got, err := Sign(sk, data)
	require.NoError(t, err)
	t.Log(got)
}
