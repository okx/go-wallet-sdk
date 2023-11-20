# @okxweb3/near-sdk
Kaspa SDK is used to interact with the Kaspa blockchain, it contains various functions can be used to web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/near
```

## Usage
### New Address
```go
	addr, _, err := NewAccount()
	if err != nil {
		// todo
	}
```

###  Transfer 
```go
	privateKey := "b9ec4d26ab5bec8df4314a9e3b8fc3f9c96d410b4cd13caa675018dcfc7916cceefbba85caaa14cb87b83314d5b86895f2d4b7633e29012e65bfb037c885c804"
	val := 0.222
	to := "ggasii.testnet"
	blockHash := "D7CPxgTXyRKYTSYuwAiRwDJH5RdHz7DwPt4EViptAW4L"
	nonce := int64(1)
	addr, err := PrivateKeyToAddr(privateKey)
    if err != nil {
		// todo
	}
	publicKeyHex, err := PrivateKeyToPublicKeyHex(privateKey)
    if err != nil {
		// todo
	}
	tx, err := CreateTransaction(addr, to, publicKeyHex, blockHash, nonce)
    if err != nil {
		// todo
	}
	amount := decimal.NewFromFloat(val).Shift(24)
	ta, err := serialize.CreateTransfer(amount.String())
    if err != nil {
		// todo
	}
	tx.SetAction(ta)
	txData, err := tx.Serialize()
    if err != nil {
		// todo
	}
	txBase58 := base58.Encode(txData)
	sig, err := SignTransaction(txBase58, privateKey)
    if err != nil {
		// todo
	}
	stx, err := CreateSignedTransaction(tx, sig)
```

###  Transfer Token
```go
	privateKey := "b9ec4d26ab5bec8df4314a9e3b8fc3f9c96d410b4cd13caa675018dcfc7916cceefbba85caaa14cb87b83314d5b86895f2d4b7633e29012e65bfb037c885c804"
	val := 0.222
	to := "ft.examples.testnet"
	blockHash := "D7CPxgTXyRKYTSYuwAiRwDJH5RdHz7DwPt4EViptAW4L"
	nonce := int64(1)
	argsStr := `{"account_id": "serhii.testnet"}`
	gas := big.NewInt(1)
	addr := "ggasii.testnet"
	publicKeyHex, err := PrivateKeyToPublicKeyHex(privateKey)
	require.NoError(t, err)
	tx, err := CreateTransaction(addr, to, publicKeyHex, blockHash, nonce)
	amount := decimal.NewFromFloat(val).Shift(24).BigInt()
	ta, err := serialize.CreateFunctionCall("storage_balance_of", []byte(argsStr), gas, amount)
	tx.SetAction(ta)
	txData, err := tx.Serialize()
	txHash := base58.Encode(txData)
	sig, err := SignTransaction(txHash, privateKey)
	stx, err := CreateSignedTransaction(tx, sig)
	if err != nil {
		// todo
	}
```

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/near/LICENSE>) licensed, see package or folder for the respective license.
