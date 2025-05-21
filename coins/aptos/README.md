# aptos-sdk
Aptos SDK is used to interact with the Aptos blockchain, it contains various functions can be used to web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/aptos
```

## Usage
### New Address
```golang
	// address
	addr, err := aptos.NewAddress("1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37", false)
	if err != nil {
		// todo
		fmt.Println(err)
	}

```
###  Transfer 
```golang
	// transfer
	from := addr
	to := "0xedc4410aa38b512e3173fcd1e119abb13872d6928dce0842664ad6ada1ccd28"
	amount := 1000
	sequenceNumber := 1
	maxGasAmount := 10000
	gasUnitPrice := 100
	expirationTimestampSecs := time.Now().Unix() + 300
	chainId := 2
	seedHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	data, err := aptos.Transfer(from, uint64(sequenceNumber), uint64(maxGasAmount), uint64(gasUnitPrice), uint64(expirationTimestampSecs), uint8(chainId),
		to, uint64(amount), seedHex)
	if err != nil {
		// todo
		fmt.Println(err)
	}
	fmt.Println(data)
```


### TransferV2
```golang
	from := "0xd1028d1c19e05b737a5ff9e2bfddee4821d329f1b1efd9e21c002aea04b83862"
	to := "0x00ca226de86c2da6716aaeddfddc2a16c76d35c67a0da2148c408d2ea1e5ad38"
	amount := 100
	sequenceNumber := 15
	maxGasAmount := 200000
	gasUnitPrice := 100
	expirationTimestampSecs := 1722564321
	chainId := 1
	seedHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	payload, err := CoinTransferPayloadV2(to, uint64(amount), "0xf22bede237a07e121b56d91a491eb7bcdfd1f5907926a9e58338f964a01b17fa::asset::USDT")

    if err != nil {
        // todo
        fmt.Println(err)
    }
	data, err := BuildSignedTransaction(from, uint64(sequenceNumber), uint64(maxGasAmount), uint64(gasUnitPrice), uint64(expirationTimestampSecs), uint8(chainId), payload, seedHex)
    if err != nil {
        // todo
        fmt.Println(err)
    }
```

### TransferCoin
```golang
	from := "0xd1028d1c19e05b737a5ff9e2bfddee4821d329f1b1efd9e21c002aea04b83862"
	to := "0x00ca226de86c2da6716aaeddfddc2a16c76d35c67a0da2148c408d2ea1e5ad38"
	amount := uint64(100)
	sequenceNumber := uint64(15)
	maxGasAmount := uint64(200000)
	gasUnitPrice := uint64(100)
	expirationTimestampSecs := uint64(1722564321)
	chainId := uint8(1)
	seedHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	tyArgs := "0xf22bede237a07e121b56d91a491eb7bcdfd1f5907926a9e58338f964a01b17fa::asset::USDT"

	data, err := TransferCoins(from, sequenceNumber, maxGasAmount, gasUnitPrice, expirationTimestampSecs, chainId, to, amount, seedHex, tyArgs)
    if err != nil {
        // todo
        fmt.Println(err)
    }

```

## Credits  This project includes code adapted from the following sources:  
- [aptos-go-sdk](https://github.com/aptos-labs/aptos-go-sdk) - Aptos Go SDK

If you are the original author and would like credit adjusted, please contact us.

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/aptos/LICENSE>) licensed, see package or folder for the respective license.
