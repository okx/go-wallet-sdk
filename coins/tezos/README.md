# tezos-sdk
Tezos SDK is used to interact with the Tezos blockchain, it contains various functions can be used to web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/tezos
```

## Usage
### New address
```golang
addr, err := GetAddress("edpkucde3WUTR2s6KgDBwvR7NiezGyHNj1aGz6WrJg6SeZWHNjDA8N")
fmt.Println(addr)
```


###  Transfer 
```golang
var amount int64 = 6000000
var fee int64 = 10000 // 0.010000 XTZ
var counter int64 = 339709
opt := NewCallOptions("BL74GqeaJ8tdFuBR2RhsGXET7MNonprQ49BBreZHE9yn9x85hJP", counter, false)
privateKey, err := types.ParsePrivateKey(p1)
require.NoError(t, err)
tx, err := NewJakartanetTransaction(n1, n3, amount, opt)
require.NoError(t, err)
err = BuildTransaction(tx, fee, privateKey.Public(), opt)
require.NoError(t, err)
rawTx, err := SignTransaction(tx, p1, opt)
require.NoError(t, err)
expected := "33b684e3912522308951aea7e274f0f97a920d9ea268de31c2ca842cba8edd5a6c00cb15c8cb2ebe15662ad5697e139eabf3e0f1aea6904efedd1480bd3fe0d403809bee0200008e1c63b65a34abf66f88b0314549ca3295004eb700cfce6d27ca0feac5877bd24a7080c52a6c89c3378f3d45642d9e6729386e8dc1bff7b4041e8e2255d01f8ab0634bba2823314406844c026dd9b9dd9d3f989708"
require.Equal(t, expected, hex.EncodeToString(rawTx))
```

###  New jakartanet delegation transaction 
```golang
var to string = "tz1foXHgRzdYdaLgX6XhpZGxbBv42LZ6ubvE"
var fee int64 = 10000
var counter int64 = 331345
opt := NewCallOptions("BL7kQbhcCsMYB954n94XcSZmS1oTYcg8J7ut2wj6iZpL3fRBdM3", counter, false)
privateKey, err := types.ParsePrivateKey(p2)
tx, err := NewJakartanetDelegationTransaction(n2, to, opt)
err = BuildTransaction(tx, fee, privateKey.Public(), opt)
rawTx, err := SignTransaction(tx, p2, opt)
expected := "3548bdffd1ce6eb49efda7ebfaa9518a5868870fadf4fc0b45906412496c24376e00de6e07b12d524f72641623528d0de4da3a4cbce3904ed29c1480bd3fe0d403ff00dd2e214620a9ceaf0c38da92f9a56954f81e5ee006e44a704a0914a1c9b5adcc99fc5a238da3a1b7649a08fb5b89d66cc04c494a67fa9a40da6c1cf91ee652bf510dc6403cb654d6a0a849ed08b06a9280b9fd02"
```

###  New jakartanet undelegation transaction 
```golang
var fee int64 = 10000 // 0.010000 XTZ
var counter int64 = 339708
opt := NewCallOptions("BL7kQbhcCsMYB954n94XcSZmS1oTYcg8J7ut2wj6iZpL3fRBdM3", counter, false)
privateKey, err := types.ParsePrivateKey(p2)
tx, err := NewJakartanetUnDelegationTransaction(n2, opt)
err = BuildTransaction(tx, fee, privateKey.Public(), opt)
rawTx, err := SignTransaction(tx, p2, opt)
```

## Credits  This project includes code adapted from the following sources:
- [tzgo](https://github.com/trilitech/tzgo) - Tezos Go SDK

If you are the original author and would like credit adjusted, please contact us.

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/aptos/LICENSE>) licensed, see package or folder for the respective license.
