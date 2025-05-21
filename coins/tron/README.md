# tron-sdk
Tron SDK is used to interact with the Tron blockchain, it contains various functions can be used to web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/tron
```

## Usage
### New address
```golang
pubKeyHex := "0357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f"
publicKey, _ := hex.DecodeString(pubKeyHex)
pub,_:=btcec.ParsePubKey(publicKey)
addr:=GetAddress(pub)
fmt.Println(addr)
```


###  Transfer 
```golang
currentTime := time.Now()
k1 := make([]byte, 8)
binary.BigEndian.PutUint64(k1, 47102802)
k2, _ := hex.DecodeString("0000000002cebb52bb1c53a37236902bac251e302a4541452b6df63f594562b9")
d2, _ := newTransfer(
    "TSAaoJuxBUxSqU7JGxzTH3gx237PTJxfwV",
    "TWYrgz7RDP2NpumQRPY1jBmPKLWVSnrzWZ",
    10000000,
    hex.EncodeToString(k1[6:8]),
    hex.EncodeToString(k2[8:16]),
    currentTime.UnixMilli()+3600*1000,
    currentTime.UnixMilli())

```

## Credits  This project includes code adapted from the following sources:
- [gotron-sdk](https://github.com/fbsobreira/gotron-sdk) - Tron Go SDK

If you are the original author and would like credit adjusted, please contact us.

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/aptos/LICENSE>) licensed, see package or folder for the respective license.
