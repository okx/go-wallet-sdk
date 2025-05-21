# near-sdk
Near SDK is used to interact with the Near blockchain, it contains various functions can be used to web3 wallet.

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

### Calculate Tx Hash
```go
	signedTx := `QAAAAGQ3Mzg4OGEyNjE5Yzc3NjE3MzVmMjNjNzk4NTM2MTQ1ZGZhODdmOTMwNmI1ZjIxMjc1ZWI0YjFhN2JhOTcxYjkA1ziIomGcd2FzXyPHmFNhRd+of5MGtfISdetLGnupcbnjWwQAAAAAAEAAAAA4OWY5Nzc1ODU5ZWQzNDY3OGVhNDhlOWExYWViMjAyY2Q0YzI5ZGNlMTViZTA2NTJiOWY1MGUyMmEwYjY3ZWY3r4iB+lQhXiP818JF0LPDjkAFNvOeVJ/lAoe14WgEF6cBAAAAAwAA4ntBSX/LsDkAAAAAAAAAJv0PcmRmmTopCCBHfD2GNR3IKgmLzEL0K70jwXkjwXqbESEFCVaymK9VP/o9bFoPYeU+AFW92TyPy1fssMHaDQ==`
	hash, err := CalTxHash(signedTx, true)
 	if err != nil {
		// todo
	}
```

### Sign Message
```go
	nonce := make([]byte, 32)
	nonce[31] = 1
	payload := serialize.NewSignMessagePayload("hello world", nonce, "", "")
	privateKey := "790e2778e0bfdae3da6419ef68c2451e80449de81e7bed9150b1cbc72b56a219d25cfdae0f9832e98bbdc87f3a156bb765cd9964e00878bf66da74591537e0a9"
	bs, err := payload.Serialize()
 	if err != nil {
		// todo
	}

	payload = serialize.NewSignMessagePayload("hello world", nonce, "", "1")
	bs, err = payload.Serialize()
 	if err != nil {
		// todo
	}

	s, err := SignMessage(payload, privateKey)
 	if err != nil {
		// todo
	}
```

## Credits  This project includes code adapted from the following sources:
- [wallet-srv](https://github.com/cnmars/wallet-srv) - Near Go SDK

If you are the original author and would like credit adjusted, please contact us.

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/near/LICENSE>) licensed, see package or folder for the respective license.
