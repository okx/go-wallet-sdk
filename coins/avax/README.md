# avax-sdk
Avax SDK is used to interact with the Avalanche blockchain, it contains various functions that can be used for web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/avax
```

## Usage
### New Address
```golang
	// address
	pk, _ := hex.DecodeString("d27a851e2ffe50d81d639a5bc17ccb488b1441307fea7636e264b9da0ce577a1")
	_, b := btcec.PrivKeyFromBytes(pk)
	addr, err := avax.NewAddress("X", "fuji", b)
	if err != nil {
		// todo
		fmt.Println(err)
	}

```
###  Transfer
```golang
	// transfer
	var inputs []avax.TransferInput
	var outputs []avax.TransferOutPut

	c := math.Pow10(9)
	inputs = append(inputs, avax.TransferInput{TxId: "sJNJVJQzmjyrAoPfshkDhKNf55jNSNW7NXK8SygGdNrst2waA", Index: 0, Amount: uint64(2 * c), AssetId: "U8iRqJoiJm8xZHAacmvYyZVwqQx6uDNtQeP3CQ6fcgQk3JqnK", PrivateKey: "bf77591baae00a9b2826ae63d6668fe5c1cd934fcaf5c99946af9d55457533ce"})
	outputs = append(outputs, avax.TransferOutPut{Address: "X-fuji1xqq48uejmydn95dwmvk4ge7rs9mj60nlx94dst", AssetId: "U8iRqJoiJm8xZHAacmvYyZVwqQx6uDNtQeP3CQ6fcgQk3JqnK", Value: uint64(c)})
	outputs = append(outputs, avax.TransferOutPut{Address: "X-fuji1asep0ygju0g2trqq2pvpez736gngthh29lkazf", AssetId: "U8iRqJoiJm8xZHAacmvYyZVwqQx6uDNtQeP3CQ6fcgQk3JqnK", Value: uint64(0.99 * c)})

	ret, err := avax.NewTransferTransaction(5, "2JVSBoinj9C2J33VntvzYtVJNZdN2NKiwwKjcumHUWEb5DbBrm", &inputs, &outputs)
	if err != nil {
		// todo
		fmt.Println(err)
	}
	fmt.Println(ret)
```

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/avax/LICENSE>) licensed, see package or folder for the respective license.
