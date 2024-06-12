# How to sign message and verify signed message?

We currently support message signing and verifying for utxo and evm chains.


# Evm sign message and verify signed message

go.mod
```golang

module example

go 1.19

require github.com/okx/go-wallet-sdk/coins/ethereum v0.0.4

require (
github.com/btcsuite/btcd/btcec/v2 v2.3.2 // indirect
github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
github.com/ethereum/go-ethereum v1.12.2 // indirect
github.com/go-stack/stack v1.8.1 // indirect
github.com/holiman/uint256 v1.2.3 // indirect
github.com/okx/go-wallet-sdk/crypto v0.0.1 // indirect
github.com/okx/go-wallet-sdk/util v0.0.1 // indirect
golang.org/x/crypto v0.12.0 // indirect
golang.org/x/sys v0.11.0 // indirect
)


```

code

```golang

package main

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/okx/go-wallet-sdk/coins/ethereum"
)

func main() {
	prvKeyHex := "49c0722d56d6bac802bdf5c480a17c870d1d18bc4355d8344aa05390eb778280"
	msg := "im from okx"
	pk, err := hex.DecodeString(prvKeyHex)
	if err != nil {
		fmt.Println(err)
		return
	}
	prv, pub := btcec.PrivKeyFromBytes(pk)
	addr := ethereum.GetAddress(hex.EncodeToString(pub.SerializeUncompressed()))
	sign, err := ethereum.SignEthTypeMessage(msg, prv, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(fmt.Sprintf("sign msg success.msg:`%s`,sign:`%s`,address:`%s`", msg, sign, addr))
	err = ethereum.VerifySignMsg(sign, msg, addr, true)
	if err != nil {
		fmt.Println("VerifySignMsg failed.")
		return
	}
	fmt.Println("verify signed msg success.")
}


```

output :
>
> sign msg success.msg:`im from okx`,sign:`ae3363a92811dcaa19676490423caa46c07140f51e047b662433ef1940afd4ff1db880f4a632f7d776b603f038ed894e2d695c9291754807d7451daff9a5026f1c`,address:`0xd74c65ad81aa8537327e9ba943011a8cec7a7b6b`
> verify signed msg success.
> 