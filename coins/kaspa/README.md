# kaspa-sdk
Kaspa SDK is used to interact with the Kaspa blockchain, it contains various functions that can be used for web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/kaspa
```

## Usage
### New Address
```go
	privateKeyHex := "559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"
	address, err := NewAddressWithNetParams(privateKeyHex, dagconfig.DevnetParams)
	if err != nil {
		// todo
	}
```

###  Transfer 
```go
	var txInputs []*TxInput
	txInputs = append(txInputs, &TxInput{
		TxId:       "120c5410cc4512f29da50a8befc67c1cfbf7bb4f594ef91c14741150d8dadd24",
		Index:      0,
		Address:    "kaspa:qrcnkrtrjptghtrntvyqkqafj06f9tamn0pnqvelmt2vmz68yp4gqj5lnal2h",
		Amount:     "900000",
		PrivateKey: "b827bb46d921bde498a535999d7554071045f02e4fdfdebb10b08583f1c6afbe",
	})
	txData := &TxData{
		TxInputs:      txInputs,
		ToAddress:     "kaspa:qqvxjssnw024e93vykhzd8d7t6dua2sx8ak4mj7xm8s9370yevxcv0jgl2xfj", // 443642da97444e52af9eb35e3d32d6270f47d255854b63299b29f21c1ded4c7c
		Amount:        "100000",
		Fee:           "10000",
		ChangeAddress: "kaspa:qrcnkrtrjptghtrntvyqkqafj06f9tamn0pnqvelmt2vmz68yp4gqj5lnal2h",
		MinOutput:     "546",
	}

	signedTx, err := Transfer(txData)
	if err != nil {
		// todo
	}
```

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/kaspa/LICENSE>) licensed, see package or folder for the respective license.
