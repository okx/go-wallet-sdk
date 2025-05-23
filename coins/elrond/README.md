# elrond-sdk
Elrond SDK is used to interact with the Elrond blockchain, it contains various functions that can be used for web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/elrond
```

## Usage
### address
```golang
	// address
    pk, _ := hex.DecodeString("27d57eb22fc218b83e9ea2da55746d9318ba6b89cfa31b797e7296bf8a66e4f1")
    address, err := AddressFromSeed(hex.EncodeToString(pk))
	if err != nil {
		// todo
		fmt.Println(err)
	}

```

###  transfer 
```golang
	// transfer
	args := ArgCreateTransaction{
		Nonce:    3,
		Value:    "10000000000000000", // decimal 18
		RcvAddr:  toAddress,
		GasPrice: 1000000000,
		GasLimit: 50000,
		ChainID:  "T",
		Version:  2,
		Options:  1,
	}
    signedTx, err := Transfer(args, pk)
	if err != nil {
		// todo
		fmt.Println(err)
	}
	fmt.Println(data)
```

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/elrond/LICENSE>) licensed, see package or folder for the respective license.
