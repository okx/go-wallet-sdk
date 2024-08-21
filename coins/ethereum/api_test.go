package ethereum

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestDecodeTx(t *testing.T) {
	r, err := DecodeTx("0xf86c81f48504a817c8008252089431c514837ee0f6062eaffb0882d764170a17800487038d7ea4c680008025a08d2cce871494c89cc0057b75ebe7ba5fcb0ca12376f8b7c3b5e1ee0bce77c2a2a028bdf0d6d9e52f50b2eecb6c21aeb288e4727478b2f588c9d0787e845cc95144")
	require.Nil(t, err)
	expected := `{"chainId":1,"from":"0x05d132975D8EfCD67262980C54f9030319C91Af0","gasLimit":21000,"gasPrice":20000000000,"nonce":244,"r":63855278322598500466154057791001702486226441534857140050433689054814601200290,"s":18428110250072635139397296879143057669157476180829444102104539880782742245700,"to":"0x31C514837ee0f6062EafFb0882D764170A178004","txType":0,"v":37,"value":1000000000000000}`
	require.Equal(t, expected, r)
}

func TestMessageHash(t *testing.T) {
	r := MessageHash("data")
	expected := "0x8edd100985b029cc35de22b18d970ad826ca8948fc18e6f4783f4728af80b113"
	require.Equal(t, expected, r)
}

func TestEcRecover(t *testing.T) {
	validRes, err := EcRecover("25064ca2f492d0a8a99801e57e5f30fc6c69335e02487d652dc98448145866556007ddc34a0ccdce592176f022e05d3e83b83a039d97aae86c1c7839cb44221e1b", "data", true)
	assert.NoError(t, err)
	expected := "0xdcab8e02b4d06d0a07ddb1dfa6e2c94cf2da2356"
	require.Equal(t, expected, validRes)

	res, err := EcRecover("0xb715378a9d1cce098c27399f40e408fe1ac314aac8ced9704905f14d0d6840c7027fdfffb800958ba826ecbdcc411571af9a348e2a78166b3f8518ce77b35c701b", "0x214f333f99d572b10721be1024700ba551b1b18ecd9072c2975cc82da63cc631", false)
	assert.NoError(t, err)
	expected = "0xb8cf89fa8f3a0ddcda3d8fdb9006859628665ef4"
	require.Equal(t, expected, res)
}

func TestGetAddress(t *testing.T) {
	validRes := GetAddress("0x04dcb83211f9cbc7f846595d93613f395fa410ca9bc1ae6ac4cbc5c63c66fc83c599dcdf22033d8a833d28a7f88da8c0d4d25ded358623068f4f2a07bdcb8d6d2c")
	require.Equal(t, "0x769ecd386d4ad9d7f7aea69f674e39efe7ea0720", validRes)
}

func TestValidateAddress(t *testing.T) {
	require.True(t, ValidateAddress("0x769ecd386d4ad9d7f7aea69f674e39efe7ea0720"))
	require.False(t, ValidateAddress("0x769ecd386d4ad9d7f7aea69f674e139ef3e7ea0720"))
}

func TestGenerateRawTransactionWithSignature(t *testing.T) {
	unsignedRawTx := "eb808505d21dba0082520894917ce722ac0b323c038e141eec7bd6bf5a00c6ac87017056cdfb974a80388080"
	r := "6fd922e7bbd9796cc49827f0b9645adfc34a16017e316e8c0b3062270b1a15bc"
	s := "685a42b37c414aa0fca332f87c1f566f59065e88f3efdaac74f2739f5788d615"
	v := "93"
	signedTx, err := GenerateRawTransactionWithSignature(0, "56", unsignedRawTx, r, s, v)
	require.Nil(t, err)
	expected := "0xf86c808505d21dba0082520894917ce722ac0b323c038e141eec7bd6bf5a00c6ac87017056cdfb974a808193a06fd922e7bbd9796cc49827f0b9645adfc34a16017e316e8c0b3062270b1a15bca0685a42b37c414aa0fca332f87c1f566f59065e88f3efdaac74f2739f5788d615"
	assert.Equal(t, expected, signedTx)
}

func TestGenerateTxWithJSON(t *testing.T) {
	raw := `{"to":"","value":"404993102026570","fee":"100000", "data":"0x123123123123","gasPrice":"20000000000","gasLimit":"21000"}`
	tx, err := GenerateTxWithJSON(raw, big.NewInt(1), false)
	require.Nil(t, err)
	require.Equal(t, "dd808504a817c8008252088087017056cdfb974a86123123123123018080", tx.Tx)
	raw = `{
				"txType":2,
				"nonce":"244",
				"to":"0x31c514837ee0f6062eaffb0882d764170a178004",
				"value":"1000000000000000",
				"gasLimit":"21000",
				"maxFeePerGas":"20000000000",
				"maxPriorityFeePerGas":"1500000000"
			}`
	tx, err = GenerateTxWithJSON(raw, big.NewInt(1), false)
	require.Nil(t, err)
	require.Equal(t, "0x02f30181f48459682f008504a817c8008252089431c514837ee0f6062eaffb0882d764170a17800487038d7ea4c6800080c0808080", tx.Tx)
	require.Equal(t, "0xbf498de3f0bc5fff92798a26f91f0a44842cf4471504de23e2eec80d53d65ad7", tx.Hash)
}
