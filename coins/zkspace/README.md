# zkspace-sdk
Zkspace SDK is used to interact with the Zkspace blockchain, it contains various functions that can be used for web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/zkspace
```

## Usage
### New address
```golang
pri := "559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"
addr, err := GetAddress(pri)
```


###  Change Pubkey
```golang
l1PrivateKeyBytes, _ := hex.DecodeString(l1PrivateKey[2:])
ethSigner, _ := core.NewOkEthSignerFromPrivBytes(l1PrivateKeyBytes)
zkSigner, _ := NewZkSignerFromEthSigner(ethSigner, core.ChainIdMainnet)
const nonce = 6
const accountId = 11573
newPkHash := zkSigner.GetPublicKeyHash()
nonceStr := "0x" + fmt.Sprintf("%08x", nonce)
accountIdStr := "0x" + fmt.Sprintf("%08x", accountId)
message := fmt.Sprintf("Register ZKSwap pubkey:\n\n%s\nnonce: %s\naccount id: %s\n\nOnly sign this message for a trusted client!",
zkSigner.GetPublicKeyHash()[5:], nonceStr, accountIdStr)
ethSignature, err := ethSigner.SignMessage([]byte(message))
from := ethSigner.GetAddress()
tx := CreateChangePubKeyTx(accountId, from, newPkHash, nonce, ethSignature)
```



###  Transfer 
```golang
l1PrivateKeyBytes, _ := hex.DecodeString(l1PrivateKey[2:])
ethSigner, _ := core.NewOkEthSignerFromPrivBytes(l1PrivateKeyBytes)
zkSigner, _ := NewZkSignerFromEthSigner(ethSigner, core.ChainIdMainnet)
from := ethSigner.GetAddress()
const nonce = 7
const accountId = 11573
const chainId = 13
const to = "0x21dceed765c30b2abea933a161479aea4702e433"
const tokenId = 1
const tokenSymbol = "ZKS"
const decimals = 18
token := &core.Token{
    Id:       1,
    Symbol:   tokenSymbol,
    Decimals: decimals,
}
amount, _ := big.NewInt(0).SetString("5000000000000000000", 10)
readableAmount := token.ToDecimalString(amount)
// calculate fee
const feeUSDT = "0.5"
const feeTokenId = 1
const feeTokenSymbol = "ZKS"
const feeDecimals = 18
const price = "0.0593863548182511"
feeToken := &core.Token{
    Id:       1,
    Symbol:   feeTokenSymbol,
    Decimals: feeDecimals,
}
fee, _ := big.NewInt(0).SetString("8410000000000000000", 10)
readableFee := feeToken.ToDecimalString(fee)

// prepare for l1 signature
message := fmt.Sprintf("Transfer %s %s\nTo: %s\nChain Id: %d\nNonce: %d\nFee: %s %s\nAccount Id: %d",
readableAmount, tokenSymbol, strings.ToLower(to), chainId, nonce, readableFee, feeTokenSymbol, accountId)

ethSignature, err := ethSigner.SignMessage([]byte(message)) // 0x549dd4788ef9abb59240d6ee0952e789df02b98890f89abf30987291b89270a73b363ddc69e9da9165cba1e7e95d23576372bd38761c4e713473d336638fd55e1b
ethereumSignature := &core.EthSignature{
    Type:      "EthereumSignature",
    Signature: "0x" + hex.EncodeToString(ethSignature),
}
tx := CreateTransferTx(accountId, from, to, tokenId, amount, feeTokenId, fee, chainId, nonce)
signature, err := zkSigner.SignTransfer(&tx)
tx.Signature = signature
transferTx := CreateSignTransferTx(&tx)
signedTransaction := &SignedTransaction{
    Transaction:       transferTx,
    EthereumSignature: ethereumSignature,
} 
```
 

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/zkspace/LICENSE>) licensed, see package or folder for the respective license.
