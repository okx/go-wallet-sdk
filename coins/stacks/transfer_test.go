package stacks

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestTransfer(t *testing.T) {
	result, err := Transfer("598d99970d04be67e8b41ddd5c5453487eeab5345ea1638c9a2849dee377f2a3", "SP2P58SJY1XH6GX4W3YGEPZ2058DD3JHBPJ8W843Q", "20", big.NewInt(3000), big.NewInt(8), big.NewInt(200))
	require.NoError(t, err)
	var s TransactionRes
	err = json.Unmarshal([]byte(result), &s)
	require.NoError(t, err)
	expected := "00000000010400971dae66fc76f7b2ba992e026d2f58578e32ee48000000000000000800000000000000c801007dc19e1e4eb9f376656b48520f86c3e59b1de6ae02462a92356e544ced9d19f6568b25f5301f9f8346997486ed7763a3ffddb6795947d7bf45558da47de3d790030200000000000516ac54665e0f6268749c1fa0eb7c402a1ad1ca2bb40000000000000bb832300000000000000000000000000000000000000000000000000000000000000000"
	require.Equal(t, expected, s.TxSerialize)
}

func TestTransfer_index_compressedAddress(t *testing.T) {
	result, err := Transfer("598d99970d04be67e8b41ddd5c5453487eeab5345ea1638c9a2849dee377f2a301", "SP2P58SJY1XH6GX4W3YGEPZ2058DD3JHBPJ8W843Q", "20", big.NewInt(3000), big.NewInt(8), big.NewInt(200))
	require.NoError(t, err)
	var s TransactionRes
	err = json.Unmarshal([]byte(result), &s)
	require.NoError(t, err)
	expected := "000000000104006ecfff9cee8ac5367c83ad0819e4c500b6c475d6000000000000000800000000000000c80001d9e056da510cc2f247ce93b72ff82d1820b21751e91e07c63662b4ee6fae9578545c1576c39f68f654eb268f1ffafb62360f056ccc5efae6aa9a9e995a0a7aeb030200000000000516ac54665e0f6268749c1fa0eb7c402a1ad1ca2bb40000000000000bb832300000000000000000000000000000000000000000000000000000000000000000"
	require.Equal(t, expected, s.TxSerialize)
}
