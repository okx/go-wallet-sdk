# nervos-sdk
Nervos SDK is used to interact with the Nervos blockchain, it contains various functions can be used to web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/nervos
```

## Usage
### New Address
```go
	privateKeyHex := "559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"
	address, err := GenerateAddressByPrivateKey("ckb",privateKeyHex)
	if err != nil {
		// todo
	}
```

###  Transfer
```go
	builder := NewTestnetTxBuild()
	// add inputs
	if err := builder.AddInput("0xcaf2cfb17eb961f54e22f8ced8656aa152f64f53e3db35b99705ca6b3822b5be", 0, 0); err != nil {
		// todo
	}
	// add outputs
	if err := builder.AddOutput(ckbTest2Address, 100*OneCKBShannon); err != nil {
		// todo
	}
	if err := builder.AddOutput(ckbTest1Address, 9895*OneCKBShannon); err != nil {
		// todo
	}
	tx, err := builder.Build()
	if err != nil {
		// todo
	}
	// sign
	if err := builder.SignByPrivateKey(ckbTest1PrivateKey); err != nil {
		// todo
	}
```

## Credits  This project includes code adapted from the following sources:
- [ckb-sdk-go](https://github.com/nervosnetwork/ckb-sdk-go) - Nervos Go SDK

If you are the original author and would like credit adjusted, please contact us.

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/nervos/LICENSE>) licensed, see package or folder for the respective license.
