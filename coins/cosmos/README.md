# cosmos-sdk
Cosmos SDK is used to interact with the Cosmos blockchain, it contains various functions that can be used for web3 wallet.
The SDK not only supports Atom, it also supports the following chains:
- Atom
- Axelar
- Cronos
- Evmos
- Iris
- Juno
- Kava
- Kujira
- Okc
- Osmos
- Secret
- Sei
- Stargaze
- Terra
- Tia

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/cosmos
```

## Usage

### Supported Functions

```golang
// NewAddress
// GetAddressByPublicKey
// ValidateAddress
// GetRawTransaction
// SignRawTransaction
// SignMessage
// Transfer
// IbcTransfer
```
### address
```golang
	// address
	hrp := "cosmos"
	pri := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	address, err := cosmos.NewAddress(pri, hrp, false)
	if err != nil {
		// todo
		fmt.Println(err)
	}
	fmt.Println(address)
	p, _ := hex.DecodeString(pri)
	_, pub := btcec.PrivKeyFromBytes(p)
	a, err := cosmos.GetAddressByPublicKey(hex.EncodeToString(pub.SerializeCompressed()), hrp)
	if err != nil {
		// todo
		fmt.Println(err)
	}
	fmt.Println(a == address)

	// ValidateAddress
	valid := cosmos.ValidateAddress(address, hrp)
	fmt.Println(valid)
```
### transfer
```golang
	// GetRawTransaction and SignRawTransaction
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
### SignMessage
```golang
	// SignMessage
	data := "{\n  \"chain_id\": \"cosmoshub-4\",\n  \"account_number\": \"584406\",\n  \"sequence\": \"1\",\n  \"fee\": {\n    \"gas\": \"250000\",\n    \"amount\": [\n      {\n        \"denom\": \"uatom\",\n        \"amount\": \"0\"\n      }\n    ]\n  },\n  \"msgs\": [\n    {\n      \"type\": \"atom/gamm/swap-exact-amount-in\",\n      \"value\": {\n        \"sender\": \"cosmos145q0tcdur4tcx2ya5cphqx96e54yflfyqjrdt5\",\n        \"routes\": [\n          {\n            \"poolId\": \"722\",\n            \"tokenOutDenom\": \"ibc/6AE98883D4D5D5FF9E50D7130F1305DA2FFA0C652D1DD9C123657C6B4EB2DF8A\"\n          }\n        ],\n        \"tokenIn\": {\n          \"denom\": \"uatom\",\n          \"amount\": \"10000\"\n        },\n        \"tokenOutMinAmount\": \"3854154180813018\"\n      }\n    }\n  ],\n  \"memo\": \"\"\n}"
	signedMessage, _, err := cosmos.SignMessage(data, pri)
	if err != nil {
		// todo
		fmt.Println(err)
	}
	fmt.Println(signedMessage)
```
### IbcTransfer
```golang
	// IbcTransfer
	paramIbc := cosmos.IbcTransferParam{}
	paramIbc.CommonParam.ChainId = "evmos_9001-2"
	paramIbc.CommonParam.Sequence = 1
	paramIbc.CommonParam.AccountNumber = 2091572
	paramIbc.CommonParam.FeeDemon = "aevmos"
	paramIbc.CommonParam.FeeAmount = "4000000000000000"
	paramIbc.CommonParam.GasLimit = 200000
	paramIbc.CommonParam.Memo = ""
	paramIbc.CommonParam.TimeoutHeight = 0
	paramIbc.FromAddress = "evmos1yc4q6svsl9xy9g2gplgnlpxwhnzr3y73wfs0xh"
	paramIbc.ToAddress = "cosmos1rvs5xph4l3px2efynqsthus8p6r4exyr7ckyxv"
	paramIbc.Demon = "aevmos"
	paramIbc.Amount = "10000000000000000"
	paramIbc.SourcePort = "transfer"
	paramIbc.SourceChannel = "channel-3"
	paramIbc.TimeOutInSeconds = uint64(time.Now().UnixMilli()/1000) + 300
	signedIBCTx, err := cosmos.IbcTransferAction(paramIbc, pri, true)
	if err != nil {
		// todo
		fmt.Println(err)
	}
	fmt.Println(signedIBCTx)
```

## Credits  This project includes code adapted from the following sources:  
- [cosmos-sdk](https://github.com/cosmos/cosmos-sdk) - Cosmos Go SDK

If you are the original author and would like credit adjusted, please contact us.

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/cosmos/LICENSE>) licensed, see package or folder for the respective license.
