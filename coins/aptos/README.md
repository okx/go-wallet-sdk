# @okxweb3/aptos-sdk
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

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/aptos/LICENSE>) licensed, see package or folder for the respective license.
