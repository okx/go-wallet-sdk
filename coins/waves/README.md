# waves-sdk
Waves SDK is used to interact with the Waves blockchain, it contains various functions that can be used for web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/waves
```

## Usage
### New address
```golang
pubKeyHash1, _ := crypto.NewPublicKeyFromBase58("2wySdbAsXi1bfAfMBKC1NcyyJemUWLM4R5ECwXJiADUx")
got, err := GetAddress(MainNetScheme, pubKeyHash1)
```

###  Transfer 
```golang
p1 := "tMUA9XRwPTiUXCTmEvU6kFkqTFKxSpaAFvQwyAT29GR"
senderPublicKey, err := crypto.NewPublicKeyFromBase58(p1)
address, err := types.NewAddressFromString(a2)
waves := types.NewOptionalAssetWaves()
tx := NewUnsignedTransferWithSig(senderPublicKey, waves, waves, 1655401735758, 2000000,
    200000, types.NewRecipientFromAddress(address), []byte("attachment"))
// sign the tx
secretKey, err := crypto.NewSecretKeyFromBase58(s1)
if err := SignTransferWithSig(tx, secretKey); err != nil {
    return
}
idBytes, err := tx.ID.MarshalJSON()
```

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/waves/LICENSE>) licensed, see package or folder for the respective license.
