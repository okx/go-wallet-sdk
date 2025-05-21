# polkdot-sdk
Polkdot SDK is used to interact with the Polkdot blockchain, it contains various functions can be used to web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/polkdot
```

## Usage
### New Address
```go
	priKey, _ := hex.DecodeString("45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40")
	p := ed25519.NewKeyFromSeed(priKey)
	publicKey := p.Public().(ed25519.PublicKey)
	address, _ := PubKeyToAddress(publicKey, PolkadotPrefix)
```

###  Transfer
```go
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
	if err != nil {
		// todo
	}
    txHash, err := CalTxHash(signed)
    hash, err 
```

## Credits  This project includes code adapted from the following sources:  
- [go-owcdrivers](https://github.com/blocktree/go-owcdrivers/tree/master/polkadotTransaction) - Polkadot Go SDK

If you are the original author and would like credit adjusted, please contact us.

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/polkdot/LICENSE>) licensed, see package or folder for the respective license.
