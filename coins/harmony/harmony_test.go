package harmony

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/okx/go-wallet-sdk/coins/ethereum"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestNewAddress(t *testing.T) {
	seedHex := "3332391f0d09af398394bbd192bb0305e130f136da8db8da82e980b9aa849681"
	addr, err := NewAddress(seedHex, true)
	require.Nil(t, err)
	expected := "0xfd01ba8507367c8a689913ea92a1526dd3893fc1"
	require.Equal(t, expected, addr)

	bytes, err := hex.DecodeString(addr[2:])
	require.Nil(t, err)
	bech32Address, err := bech32.EncodeFromBase256(HRP, bytes)
	require.Nil(t, err)
	expected = "one1l5qm4pg8xe7g56yez04f9g2jdhfcj07p4xcn0u"
	require.Equal(t, expected, bech32Address)
	hrp, addrBytes, err := bech32.DecodeToBase256(bech32Address)
	require.Nil(t, err)
	require.Equal(t, addr, "0x"+hex.EncodeToString(addrBytes))
	require.Equal(t, HRP, hrp)
}

func TestTransfer(t *testing.T) {
	p, _ := hex.DecodeString("3332391f0d09af398394bbd192bb0305e130f136da8db8da82e980b9aa849681")
	prvKey, _ := btcec.PrivKeyFromBytes(p)
	transaction := ethereum.NewEthTransaction(
		big.NewInt(int64(00)),
		big.NewInt(int64(420000)),
		big.NewInt(int64(200000000000)),
		big.NewInt(int64(100000000000)),
		"fd01ba8507367c8a689913ea92a1526dd3893fc1", "",
	)
	signedTx, err := Transfer(transaction, big.NewInt(int64(1666700000)), prvKey)
	require.Nil(t, err)
	expected := "0xf86e80852e90edd000830668a094fd01ba8507367c8a689913ea92a1526dd3893fc185174876e8008084c6afa5e4a07c3d623039edfa98327437170a3e087213dadc29d134b7151f37ff92267a04fba06ea511af335fba57bca3f4d900651459f1a292208d56088751148aa1f795bb72"
	require.Equal(t, expected, signedTx)
}
