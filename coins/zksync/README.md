# zksync-sdk
Zksync SDK is used to interact with the Zksync blockchain, it contains various functions can be used to web3 wallet.

## Installation

### go get

To obtain the latest version, simply require the project using :

```shell
go get -u github.com/okx/go-wallet-sdk/coins/zksync
```

## Usage
### New address
```golang
privKeyBytes, _ := hex.DecodeString("559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a")
ethSigner, err := core.NewOkEthSignerFromPrivBytes(privKeyBytes)
zkSigner, err := core.NewZkSignerFromEthSigner(ethSigner, core.ChainIdRinkeby)

```
###  Create withdraw transaction
```golang
accountId := uint32(1291712)
addr := "0x0bc4b0c3483084bb71614e114968c1a0ae588888"
amount := big.NewInt(12312124)
fee := big.NewInt(10000)
feeToken := RinkebyUSDC
nonce := uint32(20)
validFrom := uint64(0)
validUntil := uint64(10000000000000000)
tx := CreateWithdrawTx(accountId, addr, amount, fee, feeToken, nonce, validFrom, validUntil)

signedTx, err := SignWithdraw(tx, token, ethPrivKeyHex, int(core.ChainIdRinkeby))
```

###  Create transfer transaction
```golang 
accountId := uint32(1291712)
from := "0x0bc4b0c3483084bb71614e114968c1a0ae588888"
addr := "0x0bc4b0c3483084bb71614e114968c1a0ae588888"
amount := big.NewInt(12312124)
fee := big.NewInt(10000)
feeToken := RinkebyUSDC
nonce := uint32(20)
validFrom := uint64(0)
validUntil := uint64(10000000000000000)
tx := CreateTransferTx(accountId, from, addr, amount, fee, feeToken, nonce, validFrom, validUntil)

signedTx, err := SignTransfer(tx, ethPrivKeyHex, int(core.ChainIdRinkeby))
```

###  Create transfer with fee token transaction
```golang 
accountId := uint32(1291712)
from := "0x0bc4b0c3483084bb71614e114968c1a0ae588888"
to := "0x0e81575BF66e79915A22c614e2046d360e40a3f9"
amount := big.NewInt(12312124)
fee := big.NewInt(10000)
feeToken := RinkebyUSDC
nonce := uint32(18)
validFrom := uint64(0)
validUntil := uint64(10000000000000000)
txs := CreateTransferWithFeeTokenTx(accountId, from, to, amount, feeToken, fee, RinkebyUSDT, nonce, validFrom, validUntil)

signedTx, err := SignBatchTransfer(txs, ethPrivKeyHex, int(core.ChainIdRinkeby))
```

###  Create change pubkey transaction
```golang 
ethPrivKeyHex := "559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"
accountId := uint32(1291712)
address, _ := GetAddress(ethPrivKeyHex)
pubKeyHash, _ := GetPubKeyHash(ethPrivKeyHex, int(core.ChainIdRinkeby))
token := RinkebyUSDC
fee := big.NewInt(100000000000000)
nonce := uint32(2)
validFrom := uint64(0)
validUntil := uint64(4294967295)
tx := CreateChangePubKeyTx(accountId, address, pubKeyHash, token, fee, nonce, validFrom, validUntil)

signedTx, err := SignChangePubKey(tx, ethPrivKeyHex, int(core.ChainIdRinkeby))
```

###  Create mint nft transaction
```golang 
ethPrivKeyHex := "559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"
accountId := uint32(1291712)
address, _ := GetAddress(ethPrivKeyHex)
token := RinkebyUSDC
fee := big.NewInt(100000000000000)
nonce := uint32(29)
NFTContentHash := "1"
tx := CreateMintNFTTx(accountId, address, address, NFTContentHash, token, fee, nonce)

signedTx, err := SignMintNFT(tx, feeToken, ethPrivKeyHex, int(core.ChainIdRinkeby))
```

###  Create transfer nft transaction
```golang 
ethPrivKeyHex := "559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"
toAddress := "0x0e81575BF66e79915A22c614e2046d360e40a3f9"
accountId := uint32(1291712)
address, _ := GetAddress(ethPrivKeyHex)
token := RinkebyUSDC
fee := big.NewInt(100000000000000)
nonce := uint32(16)
nftSymbol := "NFT-113561"
validFrom := uint64(0)
validUntil := uint64(10000000000000000)
transfers, err := CreateTransferNFTTx(accountId, address, toAddress, nftSymbol, token, fee, nonce, validFrom, validUntil)

signedTx, err := SignBatchTransfer(txs, ethPrivKeyHex, int(core.ChainIdRinkeby))
```

###  Create withdraw nft transaction
```golang
ethPrivKeyHex := "559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"
accountId := uint32(1291712)
address, _ := GetAddress(ethPrivKeyHex)
token := RinkebyUSDC
fee := big.NewInt(100000000000000)
nonce := uint32(16)
nftId := uint32(113561)
validFrom := uint64(0)
validUntil := uint64(10000000000000000)
tx := CreateWithdrawNFTTx(accountId, address, nftId, token, fee, nonce, validFrom, validUntil)

signedTx, err := SignWithdrawNFT(tx, feeToken, ethPrivKeyHex, int(core.ChainIdRinkeby))
```

## Credits  This project includes code adapted from the following sources:  
- [zksync-go](https://github.com/zksync-sdk/zksync-go) - ZkSync Go SDK

If you are the original author and would like credit adjusted, please contact us.

## License
Most packages or folder are [MIT](<https://github.com/okx/go-wallet-sdk/blob/main/coins/aptos/LICENSE>) licensed, see package or folder for the respective license.
