# @okxweb3/coin-aptos
Aptos SDK is used to interact with the Aptos blockchain, it contains various functions can be used to web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get github.com/okx/go-wallet-sdk/coins/aptos
```

## Usage

### Generate private key

```golang
package main

import (
	"github.com/okx/go-wallet-sdk/coins/aptos"
)

func main() {
	wallet := AptosWallet{}
	wallet.getRandomPrivateKey()
}
```

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/LICENSE>) licensed, see package or folder for the respective license.
