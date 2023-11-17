package polkadot

import (
	"crypto/ed25519"
	"encoding/hex"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddress(t *testing.T) {
	priKey, _ := hex.DecodeString("45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40")
	p := ed25519.NewKeyFromSeed(priKey)
	publicKey := p.Public().(ed25519.PublicKey)
	require.Equal(t, "0c2f3c6dabb4a0600eccae87aeaa39242042f9a576aa8dca01e1b419cf17d7a2", hex.EncodeToString(publicKey))
	address, _ := PubKeyToAddress(publicKey, PolkadotPrefix)
	require.Equal(t, "1GycBnSYfhVWN8yRxmLxJ6CXyyZEHgVJyiF8MWZDsKTLfhs", address)
	validateAddress := ValidateAddress(address)
	require.True(t, validateAddress)
	key, _ := AddressToPublicKey("1GycBnSYfhVWN8yRxmLxJ6CXyyZEHgVJyiF8MWZDsKTLfhs")
	require.Equal(t, "0c2f3c6dabb4a0600eccae87aeaa39242042f9a576aa8dca01e1b419cf17d7a2", key)
}

func TestTransfer(t *testing.T) {
	tx := TxStruct{
		From:         "1GycBnSYfhVWN8yRxmLxJ6CXyyZEHgVJyiF8MWZDsKTLfhs",
		To:           "1GycBnSYfhVWN8yRxmLxJ6CXyyZEHgVJyiF8MWZDsKTLfhs",
		Amount:       10000000000,
		Nonce:        18,
		Tip:          0,
		BlockHeight:  10672081,
		BlockHash:    "0x569e9705bdcd3cf15edb1378433148d437f585a21ad0e2691f0d8c0083021580",
		GenesisHash:  "0x91b171bb158e2d3848fa23a9f1c25182fb8e20313b2c1eb49219da7a70ce90c3",
		SpecVersion:  9220,
		TxVersion:    12,
		ModuleMethod: "0500",
		Version:      "84",
	}

	signed, err := SignTx(tx, Transfer, "45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40")
	require.NoError(t, err)
	expected := "0x410284000c2f3c6dabb4a0600eccae87aeaa39242042f9a576aa8dca01e1b419cf17d7a200823181d175794c0438f88340b8f314d1e0e1f0e7fda5b0c0375be35482468ea6284e3831ce67b622322ad984f5a1d1868e7536e4558735fc1c9050443e1c8503150148000500000c2f3c6dabb4a0600eccae87aeaa39242042f9a576aa8dca01e1b419cf17d7a20700e40b5402"
	require.Equal(t, expected, signed)
}

func TestTransferAll(t *testing.T) {
	tx := TxStruct{
		From:         "1GycBnSYfhVWN8yRxmLxJ6CXyyZEHgVJyiF8MWZDsKTLfhs",
		To:           "1GycBnSYfhVWN8yRxmLxJ6CXyyZEHgVJyiF8MWZDsKTLfhs",
		KeepAlive:    "00", // destroy the account
		Nonce:        18,
		Tip:          0,
		BlockHeight:  10672081,
		BlockHash:    "0x569e9705bdcd3cf15edb1378433148d437f585a21ad0e2691f0d8c0083021580",
		GenesisHash:  "0x91b171bb158e2d3848fa23a9f1c25182fb8e20313b2c1eb49219da7a70ce90c3",
		SpecVersion:  9220,
		TxVersion:    12,
		ModuleMethod: "0504",
		Version:      "84",
		EraHeight:    512, // 512 blocks valid
	}

	signed, err := SignTx(tx, TransferAll, "45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40")
	require.NoError(t, err)
	expected := "0x2d0284000c2f3c6dabb4a0600eccae87aeaa39242042f9a576aa8dca01e1b419cf17d7a200f30bef08367a97e17cac7b92512d109d2b43d78c3426832ec05467c2debb8fbdf3d8a8b7ef67afc92d68c716c9ddb18b141adcfca66093b39d2ecb9db7be210e151d48000504000c2f3c6dabb4a0600eccae87aeaa39242042f9a576aa8dca01e1b419cf17d7a200"
	require.Equal(t, expected, signed)
}
