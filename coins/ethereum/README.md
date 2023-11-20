# @okxweb3/coin-ethereum
Ethereum SDK is used to interact with the Ethereum blockchain or Evm blockchains, it contains various functions can be used to web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/ethereum
```

## Usage
### New Address
```golang
	p, _ := hex.DecodeString("559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a")
	prvKey, _ := btcec.PrivKeyFromBytes(p)
	address := GetNewAddress(prvKey.PubKey())
	addr, err := PubKeyToAddr(prvKey.PubKey().SerializeUncompressed())
	if err != nil {
		// todo
	}
```

###  Transfer 
```golang
	transaction := NewEthTransaction(
		big.NewInt(int64(00)),
		big.NewInt(int64(420000)),
		big.NewInt(int64(200000000000)),
		big.NewInt(int64(100000000000000000)),
		"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728", "0x",
	)
	hash, raw, _ := transaction.GetSigningHash(big.NewInt(int64(10)))
	tx, err := transaction.SignTransaction(big.NewInt(int64(10)), prvKey)
	if err != nil {
		// todo
	}
```

### Transfer Token
```golang
	transfer, _ := token.Transfer("0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728", big.NewInt(int64(100000000000000000)))
	transaction := NewEthTransaction(
		big.NewInt(int64(00)),
		big.NewInt(int64(420000)),
		big.NewInt(int64(200000000000)),
		big.NewInt(int64(0)),
		"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728", hex.EncodeToString(transfer),
	)
	tx, err := transaction.SignTransaction(big.NewInt(int64(10)), prvKey)
	if err != nil {
		// todo
	}
```

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/ethereum/LICENSE>) licensed, see package or folder for the respective license.
