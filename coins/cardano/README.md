# cardano-sdk
Cardano SDK is used to interact with the Cardano blockchain, it contains various functions that can be used for web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/cardano
```

## Usage
### New Address
```go
	prvKeyHex := "your_private_key_hex"
	address, err := NewAddressFromPrvKey(prvKeyHex)
	if err != nil {
		// todo
	}
```

### New Address from Public Key
```go
	pubKeyHex := "your_public_key_hex"
	address, err := NewAddressFromPubKey(pubKeyHex)
	if err != nil {
		// todo
	}
```

### Sign Transaction
```go
	txData := "your_tx_data"
	signedTx, err := SignTx(txData)
	if err != nil {
		// todo
	}
```

### Sign and Verify Message
```go
	message := "your_message"
	prvKeyHex := "your_private_key_hex"
	signature, err := SignMessage(message, prvKeyHex)
	if err != nil {
		// todo
	}

	pubKeyHex := "your_public_key_hex"
	ok := VerifyMessage(message, signature, pubKeyHex)
```

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/cardano/LICENSE>) licensed, see package or folder for the respective license.
