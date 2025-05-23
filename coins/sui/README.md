# sui-sdk
Sui SDK is used to interact with the Sui blockchain, it contains various functions that can be used for web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/sui
```

## Usage
### New address
```golang
b, err := base64.StdEncoding.DecodeString("uemYAwkvsf/a7q2DdoMKNHWP7DlDhLLmgUh6coTtp94=")
addr := NewAddress(hex.EncodeToString(b[0:32]))
```


###  Stake 
```golang
data, err := BuildStakeTx("0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406", "0x72169c90b7ea87f8101285c849c09cacced9968f83aa30786dad546bb94c78ab",
		[]*SuiObjectRef{{ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Digest: "AMGM65x2qTfM4kfPjbv7Aqpap6MBiVVa4W8hrakgvPjB", Version: 1978816}},
		1000000000, 0, 9644512, 820)
require.NoError(t, err)
fmt.Println(base64.StdEncoding.EncodeToString(data))
key := "b9e99803092fb1ffdaeead8376830a34758fec394384b2e681487a7284eda7de"
tx, err := SignTransaction(base64.StdEncoding.EncodeToString(data), key)
```

###  Withdraw 
```golang
data, err := BuildWithdrawStakeTx("0x19a8a5cadebbe8f73364cb1ce4d9aa2d0ffb62a3098c7ff93bd0426960f56406",
    []*SuiObjectRef{{ObjectId: "0x26c26d20be986dbd99d33ff1b3a0bb16437039d3de85b4cb9d56a3f57066ef54", Digest: "AihBh2VjG96NDTCw1HvZj8TtWEpmYZUh9rb1D92oQ7Ak", Version: 5656730}},
    &SuiObjectRef{Digest: "CkmUVCkHFWyjH27Zg5xTd5xbUQZt1BReQnvtq2zeT6zW", Version: 5656730, ObjectId: "0x194acb4ec803ef63f15331efa9e701b4a334cf417fa15432d736d90978ce43e4"}, 0, 9534000, 820)
key := "b9e99803092fb1ffdaeead8376830a34758fec394384b2e681487a7284eda7de"
tx, err := SignTransaction(base64.StdEncoding.EncodeToString(data), key)
```

###  Sign and Verify Message 
```golang
	b, err := base64.StdEncoding.DecodeString(seed)
    seedHex := hex.EncodeToString(b[0:32])
    pubKey, err := ed25519.PublicKeyFromSeed(seedHex)

    message := "im from okx"
    signature, err := SignMessage(message, seedHex)

    err = VerifyMessage(message, base64.StdEncoding.EncodeToString(pubKey), signature)

    hash, err := hex.DecodeString("ddb521e9f8756257e16cbb657feb022ba4c270939990e3bf0194e1330be44082")
    sign, err := base64.StdEncoding.DecodeString(signature)
    err = VerifySign(pubKey, sign, hash)
```

###  Transfer amount
```golang
suiObjects := []*SuiObjectRef{{
    Digest:   "ESUg3nLfPmcMK2vf8kAyyX967w1whtgv8dk6pZhNHh6N",
    ObjectId: "0x0cd3e81f2130b922a25f113f017d8f76cae8c2f9d7ebed690e56754a0b3a5784",
    Version:  60784,
    }, {
    Digest:   "A3Nk1uPDmLLaYgBDPM235ZhTHBJVhXUD74mRVn9zZx4Z",
    ObjectId: "0x1d892361074249073f82613ad08b387bdec26185ea013d379f5ae0bb2a611ebc",
    Version:  60785,
    }, {
    Digest:   "2bacji3hre1MZiatXVjqQ1yVZXXfq5rw7yPyHwSH9CPj",
    ObjectId: "0x8565ab3ac7072abd8f7a0d4e81974a0c8669defaac9340dda7083e62542cd2f9",
    Version:  60783,
    }}
b, err := base64.StdEncoding.DecodeString("uemYAwkvsf/a7q2DdoMKNHWP7DlDhLLmgUh6coTtp94=")
addr := NewAddress(hex.EncodeToString(b[0:32]))
amount:=1
pay := &PaySuiRequest{suiObjects, amount, 0}
raw, err := json.Marshal(pay)
res, err := Execute(&Request{Data: string(raw)}, addr, recipient, gasBudget, gasPrice, hex.EncodeToString(b[0:32]))

```

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/sui/LICENSE>) licensed, see package or folder for the respective license.
