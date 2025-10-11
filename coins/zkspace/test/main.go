package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/emresenyuva/go-wallet-sdk/coins/zkspace"
	"github.com/emresenyuva/go-wallet-sdk/coins/zksync/core"
	"math/big"
	"strings"
)

func main() {
	//transfer()
	//changePubkey()
}

func transfer() {
	const l1PrivateKey = "0x559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"

	l1PrivateKeyBytes, _ := hex.DecodeString(l1PrivateKey[2:])
	ethSigner, _ := core.NewOkEthSignerFromPrivBytes(l1PrivateKeyBytes)
	zkSigner, _ := zkspace.NewZkSignerFromEthSigner(ethSigner, core.ChainIdMainnet)
	fmt.Println(zkSigner.GetPublicKeyHash())

	from := ethSigner.GetAddress()
	const nonce = 4
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
	fmt.Println(message)

	ethSignature, _ := ethSigner.SignMessage([]byte(message))
	fmt.Println(hex.EncodeToString(ethSignature))
	ethereumSignature := &core.EthSignature{
		Type:      "EthereumSignature",
		Signature: "0x" + hex.EncodeToString(ethSignature),
	}

	// prepare for l2 tx data
	transferTx := &zkspace.Transfer{
		Type:       "Transfer",
		AccountId:  accountId,
		From:       from,
		To:         to,
		TokenId:    tokenId,
		Amount:     amount,
		FeeTokenId: feeTokenId,
		Fee:        fee,
		ChainId:    chainId,
		Nonce:      nonce,
	}

	signature, _ := zkSigner.SignTransfer(transferTx)
	transferTx.Signature = signature
	fmt.Println(transferTx.Signature)

	type SignedTransaction struct {
		Transaction       *zkspace.Transfer  `json:"tx"`
		EthereumSignature *core.EthSignature `json:"signature"`
	}

	signedTransaction := &SignedTransaction{
		Transaction:       transferTx,
		EthereumSignature: ethereumSignature,
	}
	fmt.Println(transferTx, ethereumSignature)
	signedTxJson, _ := json.Marshal(signedTransaction)
	fmt.Println(string(signedTxJson))
}

func changePubkey() {
	const l1PrivateKey = "0x559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"
	l1PrivateKeyBytes, _ := hex.DecodeString(l1PrivateKey[2:])

	ethSigner, _ := core.NewOkEthSignerFromPrivBytes(l1PrivateKeyBytes)
	zkSigner, _ := zkspace.NewZkSignerFromEthSigner(ethSigner, core.ChainIdMainnet)
	fmt.Println(zkSigner.GetPublicKeyHash())

	// prepare for l1 signature
	const nonce = 5
	const accountId = 11573

	nonceStr := "0x" + fmt.Sprintf("%08x", nonce)

	fmt.Println(nonceStr)

	accountIdStr := "0x" + fmt.Sprintf("%08x", accountId)
	fmt.Println(accountIdStr)

	message := fmt.Sprintf("Register ZKSwap pubkey:\n\n%s\nnonce: %s\naccount id: %s\n\nOnly sign this message for a trusted client!",
		zkSigner.GetPublicKeyHash()[5:], nonceStr, accountIdStr)

	fmt.Println(message)

	ethSignature, _ := ethSigner.SignMessage([]byte(message))
	fmt.Println("0x" + hex.EncodeToString(ethSignature))
	from := ethSigner.GetAddress()

	// prepare for l2 tx data
	changepubkey := &zkspace.ChangePubKey{
		Type:         "ChangePubKey",
		AccountId:    accountId,
		Account:      from,
		NewPkHash:    zkSigner.GetPublicKeyHash(),
		Nonce:        nonce,
		EthSignature: "0x" + hex.EncodeToString(ethSignature),
	}
	signedTxJson, _ := json.Marshal(changepubkey)
	fmt.Println(string(signedTxJson))
}
