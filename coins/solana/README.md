# solana-sdk
Solana SDK is used to interact with the Solana blockchain, it contains various functions that can be used for web3 wallet.

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
    txHash, err:= CalTxHash(tx)
	if err != nil {
		// todo
	}

```

### Sign and Verify Message
```go
    // base58 message
    sig, err := SignMessage("2nCvHtAjwgpHHuaRMHcq3atYxyLV1oYh2tzUA6N83Xxr3sVEebEPJuY2oAb6ZwfRCYbWkHRkvw1dfsTFmpvjq3T5", "87PYrKY7ewJ25qaivxFzQ4g3fYH2ZT1CuRePJo9jCyEydJQMoVkxtS6pyAbKKBjSTxXT3PVGST3BpTpxvtEGMMQQMbbqeJAgzkF5TMNLkovkcEE7ZPm1qq6S9Ros4ZExAyckimPi8wfQW8rHhmMn9PnNaXS2bv4HJeHXXjEvzn2Ezi3CWbNQRvJs695KKtFfhGTqoabp9URM")
    if err != nil {
        // todo
    }
    err := VerifySignedMessage("2uWejjxZtzuqLrQeCH4gwh3C5TNn2rhHTdvC26dWzKfM", "87PYrKY7ewJ25qaivxFzQ4g3fYH2ZT1CuRePJo9jCyEydJQMoVkxtS6pyAbKKBjSTxXT3PVGST3BpTpxvtEGMMQQMbbqeJAgzkF5TMNLkovkcEE7ZPm1qq6S9Ros4ZExAyckimPi8wfQW8rHhmMn9PnNaXS2bv4HJeHXXjEvzn2Ezi3CWbNQRvJs695KKtFfhGTqoabp9URM", sig)
    if err != nil {
        // todo
    }

    // utf-8 message
    sig, err := SignUtf8Message("2nCvHtAjwgpHHuaRMHcq3atYxyLV1oYh2tzUA6N83Xxr3sVEebEPJuY2oAb6ZwfRCYbWkHRkvw1dfsTFmpvjq3T5", "this is a message to be signed by solana")
    if err != nil {
        // todo
    }
    err := VerifySignedUtf8Message("2uWejjxZtzuqLrQeCH4gwh3C5TNn2rhHTdvC26dWzKfM", "this is a message to be signed by solana", sig)
    if err != nil {
        // todo
    }
```

## Credits  This project includes code adapted from the following sources:  
- [solana-go](https://github.com/gagliardetto/solana-go) - Solana Go SDK

If you are the original author and would like credit adjusted, please contact us.

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/solana/LICENSE>) licensed, see package or folder for the respective license.
