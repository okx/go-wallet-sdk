# bitcoin-sdk
Bitcoin SDK is used to interact with the Bitcoin Mainnet or Testnet, it contains various functions that can be used for web3 wallet.
The SDK not only supports Bitcoin, it also supports the following chains:

- BTC
- BSV
- DOGE
- LTC
- TBTC

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/bitcoin
```

## Usage

### Supported Functions

```golang
// PubKeyToAddr
// GetAddressByPublicKey
// GenerateAddress
// SignTx
// // pbst
// GenerateSignedListingPSBTBase64
```

### New Address
```golang
	// address
	network := &chaincfg.TestNet3Params
	pubKeyHex := "0357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f"
	publicKey, err := hex.DecodeString(pubKeyHex)
	p2pkh, err := bitcoin.PubKeyToAddr(publicKey, bitcoin.LEGACY, network)
	if err != nil {
		// todo
		fmt.Println(err)
	}
	fmt.Println(p2pkh)
```
### Transfer
```golang
	// transfer btc
	txBuild := bitcoin.NewTxBuild(1, &chaincfg.TestNet3Params)
	txBuild.AddInput("0b2c23f5c2e6326c90cfa1d3925b0d83f4b08035ca6af8fd8f606385dfbc5822", 1, "", "", "", 0)
	txBuild.AddOutput("mvNnCR7EJS4aUReLEw2sL2ZtTZh8CAP8Gp", 53000)
	txBuild.AddOutput("mvNnCR7EJS4aUReLEw2sL2ZtTZh8CAP8Gp", 10000)
	pubKeyMap := make(map[int]string)
	pubKeyMap[0] = "022bc0ca1d6aea1c1e523bfcb33f46131bd1a3240aa04f71c34b1a177cfd5ff933"
	txHex, hashes, err := txBuild.UnSignedTx(pubKeyMap)
	signatureMap := make(map[int]string)
	for i, h := range hashes {
		privateBytes, _ := hex.DecodeString("1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37")
		prvKey, _ := btcec.PrivKeyFromBytes(privateBytes)
		sign := ecdsa.Sign(prvKey, util.RemoveZeroHex(h))
		signatureMap[i] = hex.EncodeToString(sign.Serialize())
	}
	txHex, err = bitcoin.SignTx(txHex, pubKeyMap, signatureMap)
	if err != nil {
		// todo
		fmt.Println(err)
	}
	fmt.Println(txHex)
```

### PSBT
```golang
	// psbt
	var inputs []*bitcoin.TxInput
	inputs = append(inputs, &bitcoin.TxInput{
		TxId:              "46e3ce050474e6da80760a2a0b062836ff13e2a42962dc1c9b17b8f962444206",
		VOut:              uint32(0),
		Sequence:          1,
		Amount:            int64(546),
		Address:           "2NF33rckfiQTiE5Guk5ufUdwms8PgmtnEdc",
		PrivateKey:        "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
		MasterFingerprint: 0xF23F9FD2,
		DerivationPath:    "m/44'/0'/0'/0/0",
		PublicKey:         "0357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f",
	})

	var outputs []*bitcoin.TxOutput
	outputs = append(outputs, &bitcoin.TxOutput{
		Address: "2NF33rckfiQTiE5Guk5ufUdwms8PgmtnEdc",
		Amount:  int64(100000),
	})
	psbtHex, err := bitcoin.GenerateUnsignedPSBTHex(inputs, outputs, network)
	if err != nil {
		// todo
		fmt.Println(err)
	}
	fmt.Println(psbtHex)
```

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/bitcoin/LICENSE>) licensed, see package or folder for the respective license.
