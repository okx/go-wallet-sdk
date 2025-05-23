# harmony-sdk
Harmony SDK is used to interact with the Harmony blockchain, it contains various functions that can be used for web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/harmony
```

## Usage
### New Address
```go
	seedHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	addr, err := NewAddress(seedHex, true)
	if err != nil {
		// todo
	}
```

###  Transfer 
```go
	p, _ := hex.DecodeString("1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37")
	prvKey, _ := btcec.PrivKeyFromBytes(p)
	transaction := ethereum.NewEthTransaction(
		big.NewInt(int64(00)),
		big.NewInt(int64(420000)),
		big.NewInt(int64(200000000000)),
		big.NewInt(int64(100000000000)),
		"97e2728c08bd0bfba631929e10bceaec8fc5c961", "",
	)
	signedTx, err := Transfer(transaction, big.NewInt(int64(1666700000)), prvKey)
	if err != nil {
		// todo
	}
```

### Verify Signed message
```go
    p, _ := hex.DecodeString("1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37")
    prvKey, _ := btcec.PrivKeyFromBytes(p)
	
    bech32Address, err := GetAddress(prvKey.PubKey())
    ethAddress := ethereum.GetNewAddress(prvKey.PubKey())

    msg := "hello world"

    signature, err := ethereum.SignEthTypeMessage(msg, prvKey, true)

    assert.Nil(t, VerifySignMsg(signature, msg, bech32Address, true))
    if err != nil {
        // todo
    }
    assert.Nil(t, VerifySignMsg(signature, msg, ethAddress, true))
    if err != nil {
		// todo
	}

```

###  Validate Address
```go
    valid := ValidateAddress("0xfd01ba8507367c8a689913ea92a1526dd3893fc1")
    validOne := ValidateAddress("one1l5qm4pg8xe7g56yez04f9g2jdhfcj07p4xcn0u")
```


## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/harmony/LICENSE>) licensed, see package or folder for the respective license.
