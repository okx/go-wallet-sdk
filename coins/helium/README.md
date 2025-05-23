# helium-sdk
Helium SDK is used to interact with the Helium blockchain, it contains various functions that can be used for web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/helium
```

## Usage
### New Address
```go
	seedHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	addr = NewAddress(seedHex)
```

###  Transfer 
```go
	to               = "13ECKNq99BqN3dHhqXRYsdUHAPCEnfFBJsYVh5aqVSsYB35M3wS"
	from             = "13Lqwnbh427csevUveZF9n3ra1LnVYQug31RFeENaYgXuK2s8UC"
	amount    uint64 = 120
	fee       uint64 = 35000
	nonce     uint64 = 2
	private          = "f5e029dd6cca805047ca64e131c0a6cf3bf45c7ad03a7a1e7681963c9b1f3043"
	tokenType        = "hnt"
	isMax            = true
	signTx, err := Sign(private, from, to, amount, fee, nonce, tokenType, isMax)
	if err != nil {
		// todo
	}
```

## Credits  This project includes code adapted from the following sources:  
- [block_sign](https://github.com/hecodev007/block_sign/tree/main/flynn/helium-go) - Helium Go SDK

If you are the original author and would like credit adjusted, please contact us.

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/helium/LICENSE>) licensed, see package or folder for the respective license.
