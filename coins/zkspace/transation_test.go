package zkspace

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/emresenyuva/go-wallet-sdk/coins/zksync/core"
	"github.com/stretchr/testify/require"
	"math/big"
	"strings"
	"testing"
)

const l1PrivateKey = "0x559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"

func TestChangePubkeyTx(t *testing.T) {
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
	require.NoError(t, err)
	from := ethSigner.GetAddress()
	tx := CreateChangePubKeyTx(accountId, from, newPkHash, nonce, ethSignature)
	signedTxJson, err := json.Marshal(tx)
	require.NoError(t, err)
	expected := "{\"type\":\"ChangePubKey\",\"accountId\":11573,\"account\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"newPkHash\":\"sync:89497052061f2e34e3c11f5afdb65df454c0d7b6\",\"nonce\":6,\"ethSignature\":\"0x3460adf9665743c9bab92fcd5ab0b0ecfdf77bb0aabbbd3ce7452a0d4d23a63e2b6b9a80f990ed0a972964db19efe7364732ad749a0f398da333a9b70c7d76b71b\"}"
	require.Equal(t, expected, string(signedTxJson))
}

func TestTransferTx(t *testing.T) {
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
	require.NoError(t, err)
	ethereumSignature := &core.EthSignature{
		Type:      "EthereumSignature",
		Signature: "0x" + hex.EncodeToString(ethSignature),
	}
	tx := CreateTransferTx(accountId, from, to, tokenId, amount, feeTokenId, fee, chainId, nonce)
	signature, err := zkSigner.SignTransfer(&tx)
	require.NoError(t, err)
	tx.Signature = signature
	transferTx := CreateSignTransferTx(&tx)
	signedTransaction := &SignedTransaction{
		Transaction:       transferTx,
		EthereumSignature: ethereumSignature,
	}
	signedTxJson, err := json.Marshal(signedTransaction)
	require.NoError(t, err)
	expected := "{\"tx\":{\"type\":\"Transfer\",\"accountId\":11573,\"from\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"to\":\"0x21dceed765c30b2abea933a161479aea4702e433\",\"token\":1,\"amount\":\"5000000000000000000\",\"feeToken\":1,\"fee\":\"8410000000000000000\",\"chainId\":13,\"nonce\":7,\"signature\":{\"pubKey\":\"38e9bc8c9e1e7b019553cc53fce4dde3f71fe5d678ebfc86cf000acdf413cf2c\",\"signature\":\"58e4d4e4ebd08a0f3f7fa00a6d9330a0f3d4f8aa3b0582a35f6182155105b428d9d8b3030dbaec6a78bb7e9e95fd95a5bb7f74a25a6c25cefc60b5421a989f05\"}},\"signature\":{\"type\":\"EthereumSignature\",\"signature\":\"0x2c64162414224b1f44c349912373ba6f2611fe63fe18890922359040dad6e3d73adaf4a9fe9e6887823012368b4a056c98c32c32aa21fd68196d03db9426ccde1b\"}}"
	require.Equal(t, expected, string(signedTxJson))
}
