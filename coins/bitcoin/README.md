# @okxweb3/coin-bitcoin
Bitcoin SDK is used to interact with the Aptos blockchain, it contains various functions can be used to web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/bitcoin
```

## Usage

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
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/aptos/LICENSE>) licensed, see package or folder for the respective license.
