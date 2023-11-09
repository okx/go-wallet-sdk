# @okxweb3/coin-cosmos
Cosmos SDK is used to interact with the Aptos blockchain, it contains various functions can be used to web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/cosmos
```

## Usage

```golang
	// address
	pri := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	address, err := cosmos.NewAddress(pri, "cosmos", false)
	if err != nil {
		// todo
		fmt.Println(err)
	}
	fmt.Println(address)

	// transfer
	pk, err := hex.DecodeString(pri)
	k, _ := btcec.PrivKeyFromBytes(pk)
	param := cosmos.TransferParam{}
	param.FromAddress = "cosmos145q0tcdur4tcx2ya5cphqx96e54yflfyqjrdt5"
	param.ToAddress = "cosmos1jun53r4ycc8g2v6tffp4cmxjjhv6y7ntat62wn"
	param.Demon = "uatom"
	param.Amount = "10000"
	param.CommonParam.ChainId = "cosmoshub-4"
	param.CommonParam.Sequence = 0
	param.CommonParam.AccountNumber = 623151
	param.CommonParam.FeeDemon = "uatom"
	param.CommonParam.FeeAmount = "10"
	param.CommonParam.GasLimit = 100
	param.CommonParam.Memo = "memo"
	param.CommonParam.TimeoutHeight = 0
	doc, err := cosmos.GetRawTransaction(param, hex.EncodeToString(k.PubKey().SerializeCompressed()))
	signature, err := cosmos.SignRawTransaction(doc, k)
	signedTransaction, err := cosmos.GetSignedTransaction(doc, signature)
	if err != nil {
		// todo
		fmt.Println(err)
	}
	fmt.Println(signedTransaction)
```

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/aptos/LICENSE>) licensed, see package or folder for the respective license.
