# filecoin-sdk
Filecoin SDK is used to interact with the Filecoin blockchain, it contains various functions that can be used for web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/filecoin
```

## Usage
### New Address
```golang
	privateKeyHex := "0x7c6207c56b6aa5ed4345c5f662816408e273cbdf64e2f01d54ced0125d6172c2"
	address, err := GetAddressByPrivateKey(privateKeyHex, MainnetPrefix)
	if err != nil {
		// todo
	}
```

###  Transfer 
```golang
	from := "f1bh3d2y6xxugpg3ygzxnjhcrs5ffxh5nvqmanbia"
	to := "f1fvs2fjqr6ozk477zkwzvermhledmfkswt34cmhi"
	nonce := 0
	value := big.NewInt(20000000000000000)
	gasLimit := 210000
	gasFeeCap := big.NewInt(9455791480)
	gasPremium := big.NewInt(120242)
	method := 0
	message := NewTx(from, to, nonce, method, gasLimit, value, gasFeeCap, gasPremium)
	privateKeyHex := "0x7c6207c56b6aa5ed4345c5f662816408e273cbdf64e2f01d54ced0125d6172c2"
	tx, err := SignTx(message, privateKeyHex)
	if err != nil {
		// todo
	}
    hash, err := CalTxHash(tx)
    if err != nil {
        // todo
    }
```

## Credits  This project includes code adapted from the following sources:
- [lotus](https://github.com/filecoin-project/lotus) - Filecoin Go SDK

If you are the original author and would like credit adjusted, please contact us.

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/filecoin/LICENSE>) licensed, see package or folder for the respective license.
