# starknet-sdk
Starknet SDK is used to interact with the Starknet blockchain, it contains various functions that can be used for web3 wallet.

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

### Sign and Verify Message
```go
	sc := SC()
	pk, err := GetPubKeyPoint(sc, "5d54413cc091c4c25584706c2eca3bdfd9119b9313eb81f457afd263c52eabd")
	if err != nil {
		// todo
	}
	
	sig, err := SignMsg(sc, "0xb0a391057a8c2ce9a6e8799f2609da2012970a513a700960e68f05c5c0cc26", "5d54413cc091c4c25584706c2eca3bdfd9119b9313eb81f457afd263c52eabd")
	if err != nil {
		// todo
	}
	
    ok := VerifyMsgSign("025dd2ddf7155286341c632b5ee092b52267cc09c73a079756393c79baf5d5b8", "0xb0a391057a8c2ce9a6e8799f2609da2012970a513a700960e68f05c5c0cc26", sig)

```

### Calculate Tx Hash
```go
    txHash, err := GetTxHash("{\"type\":\"INVOKE_FUNCTION\",\"sender_address\":\"0x0179aa76deab144ef996ddda6b37f9fb259c291f7b79f4e0fca63e64228a53f5\",\"calldata\":[\"1\",\"2087021424722619777119509474943472645767659996348769578120564519014510906823\",\"232670485425082704932579856502088130646006032362877466777181098476241604910\",\"0\",\"3\",\"3\",\"2101208752900774171800778204657581671583980985264470748488502665721568772719\",\"6805000000000000\",\"0\"],\"max_fee\":\"0x1fd512913000\",\"signature\":[\"1239534042151864320196515971505820081878747415448571244676980360270654138666\",\"1089426459092570934939188723465380278018225693915907648771183733616479070762\"],\"version\":\"0x1\",\"nonce\":\"0xa2758\"}")
    if err != nil {
        // todo
	}
}

```

## Credits  This project includes code adapted from the following sources:  
- [starknet.go](https://github.com/NethermindEth/starknet.go) - Starknet Go SDK

If you are the original author and would like credit adjusted, please contact us.

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/starknet/LICENSE>) licensed, see package or folder for the respective license.
