# How to sign message and verify signed message?

We currently support message signing and verifying for utxo and evm chains.

chains list

|       | Chain Index | Chain Name | Coin/Short Name | Chain type | Prefix for UTXO | Remark  |
|-------|-------------|------------|-----------|------------|-----------------|---------|
|1 | 0 | BTC | btc | BTC        | "Bitcoin"                   |  |  
|2 | 1 | Ethereum | eth | EVM        |                 |  |  
|3 | 2 | Litecoin | LTC | BTC        | "Litecoin"      |  |  
|4 | 3 | Dogecoin | DOGE | BTC        | "Dogecoin"      |  |  
|5 | 5 | Dash | DASH | BTC        | "DarkCoin"      |  |  
|6 | 10 | Optimism | Optimism | EVM        |                 |  |  
|7 | 14 | Flare | Flr | EVM        |                 |  |  
|8 | 25 | Cronos | Cronos | EVM        |                 |  |  
|9 | 56 | BNB Smart Chain | bsc | EVM        |                 |  |  
|10 | 61 | Ethereum Classic | ETC | EVM        |                 |  |  
|11 | 66 | OKTC | okt | EVM        |                 | Only evm addresses are supported  |
|12 | 100 | Gnosis | XDAI | EVM        |                 |  |  
|13 | 133 | Zcash | ZEC | BTC        | "Zcash"         |  |  
|14 | 137 | Polygon | matic | EVM        |                 |  |  
|15 | 145 | Bitcoin Cash | BCH | BTC        | "Bitcoin"       |  |  
|16 | 169 | Manta Pacific | Manta | EVM        |                 |  |  
|17 | 204 | opBNB | op_bnb | EVM        |                 |  |  
|18 | 236 | BitcoinSV | BSV | BTC        | "Bitcoin"       |  |  
|19 | 250 | Fantom | ftm | EVM        |                 |  |  
|20 | 288 | Boba | BOBA | EVM        |                 |  |  
|21 | 314 | Filecoin | FIL | EVM        |                 |Only evm addresses are supported  |
|22 | 321 | KCC | KCC | EVM        |                 |  |  
|23 | 324 | zkSync Era | ZKSync2 | EVM        |                 |  |  
|24 | 369 | PulseChain | PLS |            |                 |   |
|25 | 408 | Omega Network | Omega | EVM        |                 |
|26 | 648 | Endurance | ACE | EVM        |                 |  |  
|27 | 1030 | Conflux | cfx | EVM        |                 |  |  
|28 | 1088 | Metis | METIS | EVM        |                 |  |  
|29 | 1101 | Polygon zkEVM | PolygonZK | EVM        |                 |  |  
|30 | 1111 | Wemix 3.0 | wemix | EVM        |                 |  |  
|31 | 1116 | Core | core | EVM        |                 |  |  
|32 | 2020 | Ronin | RONIN | EVM       |                 |   |  |  
|33 | 2222 | Kava EVM | EVM_KAVA | EVM        |                 |  |  
|34 | 42161 | Arbitrum One | Arbitrum | EVM        |                 |  |  
|35 | 42170 | Arbitrum Nova | Nova | EVM        |                 |  |  
|36 | 42220 | Celo | Celo | EVM        |                 |  |  
|37 | 43114 | Avalanche C | avax | EVM        |                 |  |  
|38 | 59144 | Linea | linea_eth | EVM        |                 |  |  
|39 | 7000 | ZetaChain Mainnet | ZETACHAIN_MAINNET | EVM        |                 |  |  
|40 | 8217 | Klaytn | Klay | EVM        |                 |  |  
|41 | 81457 | Blast | Blast | EVM        |                 |  |  
|42 | 10001 | EthereumPoW | ETHW | EVM        |                 |  |  
|43 | 4200 | Merlin Chain | Merlin_Chain | EVM        |                 |  |  
|44 | 534352 | Scroll | Scroll_eth | EVM        |                 |  |  
|45 | 11155111 | Sepolia | SEPOLIA_ETH | EVM        |                 |  |  
|46 | 13371 | Immutable zkEVM | Immutable_zkEVM | EVM        |                 |  |  
|51 | 1101 | Polygon zkEVM | PolygonZK | EVM        |                 |  |  
|52 | 1111 | Wemix 3.0 | wemix | EVM        |                 |  |  
|53 | 324 | zkSync Era | ZKSync2 | EVM        |                 |  |  
|54 | 13371 | Immutable zkEVM | Immutable_zkEVM | EVM        |                 |  |  
|47 | 1313161554 | Aurora | Aurora | EVM        |                 |  |  
|48 | 1284 | Moonbeam | GLMR | EVM        |                 |  |  
|49 | 1285 | Moonriver | MOVR | EVM        |                 |  |  
|50 | 1666600000 | Harmony | one | EVM        |                 |  |  
|51 | 513100 | DIS CHAIN | DIS | EVM        |                 |  |  
|52 | 5000 | Mantle | Mantle | EVM        |                 |  |  
|53 | 8453 | Base | base_eth | EVM        |                 |  |  
|54 | 11235 | HAQQ Network | HAQQ | EVM        |                 |  |  
|55 | 70000038 | BTC Testnet | TBTC | BTC        | "Bitcoin"       |        |        

# Evm sign message and verify signed message

## go.mod

```golang

module example

go 1.19

require github.com/okx/go -wallet-sdk/coins/ethereum v0.0.4

require (
github.com/btcsuite/btcd/btcec/v2 v2.3.2 // indirect
github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
github.com/ethereum/go -ethereum v1.12.2 // indirect
github.com/go -stack/stack v1.8.1 // indirect
github.com/holiman/uint256 v1.2.3 // indirect
github.com/okx/go -wallet-sdk/crypto v0.0.1 // indirect
github.com/okx/go-wallet-sdk/util v0.0.1 // indirect
golang.org/x/crypto v0.12.0 // indirect
golang.org/x/sys v0.11.0 // indirect
)


```

## code

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

## output

>
> sign msg success.msg:`im from okx`
>
,sign:`ae3363a92811dcaa19676490423caa46c07140f51e047b662433ef1940afd4ff1db880f4a632f7d776b603f038ed894e2d695c9291754807d7451daff9a5026f1c`
> ,address:`0xd74c65ad81aa8537327e9ba943011a8cec7a7b6b`
> verify signed msg success.
>
>
>

#

## go.mod

```go
module example

go 1.19

require github.com/okx/go -wallet-sdk/coins/ethereum v0.0.4

require (
github.com/btcsuite/btcd/btcec/v2 v2.3.2
github.com/okx/go -wallet-sdk/coins/bitcoin v0.0.3
)

require (
github.com/bits-and-blooms/bitset v1.7.0 // indirect
github.com/btcsuite/btcd v0.23.4 // indirect
github.com/btcsuite/btcd/btcutil v1.1.3 // indirect
github.com/btcsuite/btcd/btcutil/psbt v1.1.8 // indirect
github.com/btcsuite/btcd/chaincfg/chainhash v1.0.2 // indirect
github.com/btcsuite/btclog v0.0.0-20170628155309-84c8d2346e9f // indirect
github.com/btcsuite/btcutil v1.0.2 // indirect
github.com/consensys/bavard v0.1.13 // indirect
github.com/consensys/gnark-crypto v0.12.1 // indirect
github.com/crate-crypto/go -kzg-4844 v0.3.0 // indirect
github.com/decred/dcrd/crypto/blake256 v1.0.0 // indirect
github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
github.com/ethereum/c-kzg-4844 v0.3.1 // indirect
github.com/ethereum/go -ethereum v1.13.4 // indirect
github.com/go -stack/stack v1.8.1 // indirect
github.com/holiman/uint256 v1.2.3 // indirect
github.com/mmcloughlin/addchain v0.4.0 // indirect
github.com/okx/go -wallet-sdk/crypto v0.0.1 // indirect
github.com/okx/go -wallet-sdk/util v0.0.1 // indirect
github.com/supranational/blst v0.3.11 // indirect
golang.org/x/crypto v0.14.0 // indirect
golang.org/x/sync v0.3.0 // indirect
golang.org/x/sys v0.13.0 // indirect
rsc.io/tmplfunc v0.0.3 // indirect
)



# BTC sign message and verify signed message

```

## code

```go
package main

import (
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/bitcoin"
)

func main() {
	wif := "L5M17m37LiS4vABs6MzVY7VG4BoR9T8QQxMLEDyyB8b2caYKsdgp"
	msg := "im from okx"
	sign, err := bitcoin.SignMessage(wif, "", msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	address := "1F4skTsARmoJ3HQK94th7rzbteo4Rdkic9"
	pubKeyHex := "02802fdc092f06cad59b2b737020030bde0ca0d2f30d7644b51fb5062606ab1b65"
	fmt.Println(fmt.Sprintf("sign msg success.msg:`%s`,sign:`%s`,address:`%s`", msg, sign, address))
	err = bitcoin.VerifyMessage(sign, "", msg, pubKeyHex, address, "", bitcoin.GetBTCMainNetParams())
	if err != nil {
		fmt.Println("VerifySignMsg failed.")
		return
	}
	fmt.Println("verify signed msg success.")
}

```

### output

>
> sign msg success.msg:`im from okx`
> ,sign:`H5lIZuepl9PgjHqmhSNS8gcQvHdjwqHMWdsw+xpcYOpoZGeKgw+KrIwYlP46LcCE+SaUYemLojvCR0gUh+rR1Co=`
> ,address:`1F4skTsARmoJ3HQK94th7rzbteo4Rdkic9`
> verify signed msg success.
>
> 