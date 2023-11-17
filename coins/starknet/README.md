# @okxweb3/coin-starknet
Starknet SDK is used to interact with the Starknet blockchain, it contains various functions can be used to web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/starknet
```

## Usage
### New Address
```go
	curve := SC()
	starkPub, _ := GetPubKey(curve, "0x01651242558d251b0daa72cdf11feb1713e47eb88fb55d0978a2625445a771ac")
	calculateAddress, err := CalculateContractAddressFromHash(starkPub)
	if err != nil {
		// todo
	}
	preAddress := BigToHexWithPadding(calculateAddress)
```

###  Transfer
```go
	curve := SC()
	contractAddr := ETHBridge
	from := "0x076a18ceb1638b364b2bccd7652b3d024b0192b6cd97932d7a25638cd0c38cc3"
	maxFee := big.NewInt(1864315586779310)
	nonce := big.NewInt(2)
	functionName := "initiate_withdraw"
	calldata := []string{"0x62e206b4ddd402056d881ded58c0bd87193d2913", "0x38d7ea4c68000"}
	tx, err := CreateSignedContractTx(curve, contractAddr, from, functionName, calldata, nonce, maxFee, MAINNET_ID, "0x01651242558d251b0daa72cdf11feb1713e47eb88fb55d0978a2625445a771ac")
	if err != nil {
		// todo
	}
```

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/starknet/LICENSE>) licensed, see package or folder for the respective license.
