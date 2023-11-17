# @okxweb3/coin-solana
Solana SDK is used to interact with the Solana blockchain, it contains various functions can be used to web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/solana
```

## Usage
### New Address
```go
    pk, _ := base.NewRandomPrivateKey()
	address, err := NewAddress(hex.EncodeToString(pk.Bytes()))
    if err != nil {
        // todo
    }
```

###  Transfer
```go
	fromPrivate, _ := base.PrivateKeyFromBase58("tzyJiBd5PzFPFfVnnfVx14rsfC8FKW8idpJwNhH6FxzZAdhgBp4CrDxcUW9D89f5k3W6WhVnybbAw7RRB2HPxnt")
	to := "7NRmECq1R4tCtXNvmvDAuXmii3vN1J9DRZWhMCuuUnkM"
	hash := "Cfudd6AiXTzPYrmEBGNFsHgaNKJ3xrrsGCT39avLkoiu"
	// FZNZLT5diWHooSBjcng9qitykwcL9v3RiNrpC3fp9PU1
	from := fromPrivate.PublicKey().String()
	rawTransaction := NewRawTransaction(hash, from)
	rawTransaction.AppendTransferInstruction(1000000000, from, to)
	rawTransaction.AppendSigner(hex.EncodeToString(fromPrivate.Bytes()))
	tx, err := rawTransaction.Sign(true)
	if err != nil {
		// todo
	}
```

### Transfer Token
```go
	hash := "H6TNM3fDg5wTYT4eiv2PnGdd1555a45FEJtxVLtzv9dJ"
	fromPrivate, _ := base.PrivateKeyFromBase58("tzyJiBd5PzFPFfVnnfVx14rsfC8FKW8idpJwNhH6FxzZAdhgBp4CrDxcUW9D89f5k3W6WhVnybbAw7RRB2HPxnt")
	from := fromPrivate.PublicKey().String()
	to := "7NRmECq1R4tCtXNvmvDAuXmii3vN1J9DRZWhMCuuUnkM"
	mint := "4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU"
	fromAssociated, _, _ := base.FindAssociatedTokenAddress(base.MustPublicKeyFromBase58(from), base.MustPublicKeyFromBase58(mint))
	toAssociated, _, _ := base.FindAssociatedTokenAddress(base.MustPublicKeyFromBase58(to), base.MustPublicKeyFromBase58(mint))
	rawTransaction := NewRawTransaction(hash, from)
	rawTransaction.AppendAssociatedTokenAccountCreateInstruction(from, to, mint)
	rawTransaction.AppendTokenTransferInstruction(1000000, fromAssociated.String(), toAssociated.String(), from)
	rawTransaction.AppendSigner(hex.EncodeToString(fromPrivate.Bytes()))
	tx, err := rawTransaction.Sign(true)
	if err != nil {
		// todo
	}
```

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/solana/LICENSE>) licensed, see package or folder for the respective license.
