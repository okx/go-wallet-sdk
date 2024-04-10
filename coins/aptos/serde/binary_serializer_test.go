package serde

import (
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestBinarySerializer_SerializeU64(t *testing.T) {
	bs := NewBinarySerializer(500)
	err := bs.SerializeU64(1700312272)
	require.NoError(t, err)
	h := hex.EncodeToString(bs.GetBytes())
	t.Log(h)
	require.Equal(t, "d0b4586500000000", h)
}

func TestBinarySerializer_SerializeU128(t *testing.T) {
	bs := NewBinarySerializer(500)
	s := "fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe"
	b, ok := new(big.Int).SetString(s[32:], 16)
	require.True(t, ok)
	u, err := FromBig(b)
	require.NoError(t, err)
	err = bs.SerializeU128(*u)
	require.NoError(t, err)
	h := hex.EncodeToString(bs.GetBytes())
	t.Log(h)
	require.Equal(t, "feffffffffffffffffffffffffffffff", h)
}

func TestBinarySerializer_SerializeU256(t *testing.T) {
	bs := NewBinarySerializer(500)
	b, ok := new(big.Int).SetString("1ffffffffffffffffffffffffffffffefffffffffffffffffffffffffffffffe", 16)
	//b, ok := new(big.Int).SetString("1", 16)
	require.True(t, ok)
	u256, err := BigInt2U256(b)
	require.NoError(t, err)
	err = bs.SerializeU256(*u256)
	require.NoError(t, err)
	h := hex.EncodeToString(bs.GetBytes())
	t.Log(h)
	require.Equal(t, "fefffffffffffffffffffffffffffffffeffffffffffffffffffffffffffff1f", h)
}
