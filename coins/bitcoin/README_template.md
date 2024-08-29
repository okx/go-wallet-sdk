# bitcoin-sdk
Bitcoin SDK is used to interact with the Bitcoin Mainnet or Testnet, it contains various functions can be used to web3 wallet.
The SDK not only support Bitcoin, it also supports following chains:

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

## Supported Functions

* [1. PubKeyToAddr](#1-PubKeyToAddr)
* [2. SignTx](#2-SignTx)
* [3. PSBT](#3-PSBT)
  * [3.1 GenerateUnsignedPSBTHex](#31-GenerateUnsignedPSBTHex)

### 1. PubKeyToAddr
根据公钥、地址类型和网络类型计算出来bitcoin地址
* Parameters:
    1. **publicKey**: `[]byte`, 必须是长度为65或33的字节数组
    2. **addrType**: `string`,  地址类型，目前支持的地址类型有`LEGACY`,`SEGWIT_NATIVE`,`SEGWIT_NESTED`和`SEGWIT_NESTED`，共4中类型的地址。
    3. **network**: `*chaincfg.Params`, bitcoin链参数，`chaincfg.MainNetParams`或者`chaincfg.TestNet3Params`等
* Returns:
    1. `string`,  根据输入参数的不同生成不同的bitcoin地址
    2. `error`, 
* example
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

### 2. SignTx
签名交易
* Parameters://方法输入参数介绍，格式为，参数名：类型，描述
    1. **raw**: `string`, 待签名的btc交易，`&wire.MsgTx`类型的序列化结果。
    2. **pubKeyMap**: `map[int]string`,  公钥列表。
    3. **signatureMap**: `map[int]string`, 签名列表。
* Returns:
    1. `string`,  返回签好名的交易，hex编码的字符串
    2. `error`,
* example
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

### 3. PSBT

构建psbt交易

#### 3.1 GenerateUnsignedPSBTHex
生成未签名的psbt交易

* Parameters://方法输入参数介绍，格式为，参数名：类型，描述
    1. **ins**: `[]*TxInput`,  交易输入参数。
        1. **TxId**: `string`, utxo交易ID
        2. **VOut**: `uint32`, utxo交易输出index
        3. **Sequence**: `uint32`,
        4. **Amount**:    `int64`, utxo交易输出的BTC数量
        5. **Address**:  `string`, 
        6. **PrivateKey**: `string`, base58编码的私钥，wif
        7. **MasterFingerprint**: `uint32`,
        8. **DerivationPath**:   `string`,
        9. **PublicKey**: `string`,
    2. **outs**: `[]*TxOutput`,  交易输出参数。
    3. **network**: `*chaincfg.Params`, bitcoin网络参数。
* Returns:
    1. `string`,  返回签好名的交易，hex编码的字符串
    2. `error`,
* example
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
