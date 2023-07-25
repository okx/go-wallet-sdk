package polkadot

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestAddress(t *testing.T) {
	priKey, _ := hex.DecodeString("45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40")
	p := ed25519.NewKeyFromSeed(priKey)
	publicKey := p.Public().(ed25519.PublicKey)
	fmt.Println("publicKey: ", hex.EncodeToString(publicKey))

	address, _ := PubKeyToAddress(publicKey, PolkadotPrefix)
	fmt.Println("address: ", address)

	validateAddress := ValidateAddress(address)
	fmt.Println(validateAddress)

	key, _ := AddressToPublicKey("1GycBnSYfhVWN8yRxmLxJ6CXyyZEHgVJyiF8MWZDsKTLfhs")
	fmt.Println(key)
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

	signed, _ := SignTx(tx, Transfer, "45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40")
	fmt.Println(signed)
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

	signed, _ := SignTx(tx, TransferAll, "45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40")
	fmt.Println(signed)
}
