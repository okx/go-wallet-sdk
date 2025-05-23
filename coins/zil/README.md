# zil-sdk
Zil SDK is used to interact with the Zil blockchain, it contains various functions that can be used for web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/zil
```

## Usage
### New address
```golang
privHex := "c0dc46b9f9d6ef1c88dff3f1e82adc61cb11d77ab76a8d66338f14c2711cb4d8"
address, err := GetAddressFromPrivateKey(privHex)
```

###  New  bech32 address
```golang
addr, err := FromBech32Addr("zil1h6j9d76cp997r3lenwmdzkzdemry9v9su5ddz8")
```


###  Sign transaction
```golang
privateKey := "c0dc46b9f9d6ef1c88dff3f1e82adc61cb11d77ab76a8d66338f14c2722cb4d8"
to := "zil1fwh4ltdguhde9s7nysnp33d5wye6uqpugufkz7"
gasPrice := "2000000000"
amount := big.NewInt(10000000000)
gasLimit := big.NewInt(50)
nonce := 2
chainId := 333
tx := CreateTransferTransaction(to, gasPrice, amount, gasLimit, nonce, chainId)
err := SignTransaction(privateKey, tx)
payload := tx.ToTransactionPayload()
```

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/zil/LICENSE>) licensed, see package or folder for the respective license.
