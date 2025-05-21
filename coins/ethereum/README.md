# ethereum-sdk
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

### Sign message
```golang

    prv := "49c0722d56d6bac802bdf5c480a17c870d1d18bc4355d8344aa05390eb778280"
    prvB, err := hex.DecodeString(prv)
    assert.NoError(t, err)
    msg := "0x49c0722d56d6bac802bdf5c480a17c870d1d18bc4355d8344aa05390eb778280"
    prvKey, pub := btcec.PrivKeyFromBytes(prvB)
    sig, err := SignEthTypeMessage(msg, prvKey, true)   //true means using the format "\x19Ethereum Signed Message:\n%d%s"
	if err != nil {
	// todo
	}

```
### Verify Signed message
```golang
    sig:="d87758593e0b89f8a2deef5e053ce484fe971a75124bf5d89d6f4d4f586604120d0110d03c91260fec9ec917354caae50c1744d246e30ff48def277d7d9aec831b"
    msg:="0x49c0722d56d6bac802bdf5c480a17c870d1d18bc4355d8344aa05390eb778280"
    addr:="0xd74c65ad81aa8537327e9ba943011a8cec7a7b6b"
    err := VerifySignMsg(sig, msg, addr, true) //true means using the format "\x19Ethereum Signed Message:\n%d%s"
    if err != nil {
        // todo
    }

```

### EIP712
```golang
    var typedData TypedData
    str := `{"domain":{"name":"AuthTransfer","chainId":1,"verifyingContract":"0x1243C09717e4441341472c4b142B8ac0B71F7672"},"message":{"details":[{"token":"0x0000000000000000000000000000000000000000","expiration":1853395200}],"spenders":["0x1B256B89462710a6b459540B999AbE5771d45A6e"],"nonce":0},"primaryType":"Permits","types":{"EIP712Domain":[{"name":"name","type":"string"},{"name":"chainId","type":"uint256"},{"name":"verifyingContract","type":"address"}],"Permits":[{"name":"details","type":"PermitDetails[]"},{"name":"spenders","type":"address[]"},{"name":"nonce","type":"uint256"}],"PermitDetails":[{"name":"token","type":"address"},{"name":"expiration","type":"uint256"}]}}`
    err := json.Unmarshal([]byte(str), &typedData)
    assert.NoError(t, err)
    hash, str2, err := TypedDataAndHash(typedData)
	if err != nil {
	// todo
	}
```

### Dynamic Fee Tx
```golang
	tx := NewEthDynamicFeeTx(big.NewInt(int64(11155111)),
        16,
        big.NewInt(int64(420000)),
        big.NewInt(int64(20000000000)),
        420000,
        big.NewInt(int64(1234)), "2de4898dd458d6dce097e29026d446300e3815fa", "", AccessList{})
    p, _ := hex.DecodeString("5dfce364a4e9020d1bc187c9c14060e1a2f8815b3b0ceb40f45e7e39eb122103")
    prvKey, _ := btcec.PrivKeyFromBytes(p)
    txStr, err := tx.SignTransaction(prvKey)
    if err != nil {
        // todo
    }
```

## Credits  This project includes code adapted from the following sources:  
- [go-ethereum](https://github.com/ethereum/go-ethereum) - Ethereum Go SDK

If you are the original author and would like credit adjusted, please contact us.

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/ethereum/LICENSE>) licensed, see package or folder for the respective license.
