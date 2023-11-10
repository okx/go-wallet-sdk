package zksync

import (
	"encoding/json"
	"github.com/okx/go-wallet-sdk/coins/zksync/core"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestCreateWithdrawTx(t *testing.T) {
	accountId := uint32(1291712)
	addr := "0x0bc4b0c3483084bb71614e114968c1a0ae588888"
	amount := big.NewInt(12312124)
	fee := big.NewInt(10000)
	feeToken := RinkebyUSDC
	nonce := uint32(20)
	validFrom := uint64(0)
	validUntil := uint64(10000000000000000)
	tx := CreateWithdrawTx(accountId, addr, amount, fee, feeToken, nonce, validFrom, validUntil)
	txBytes, err := json.Marshal(tx)
	require.NoError(t, err)
	expected := "{\"type\":\"Withdraw\",\"accountId\":1291712,\"from\":\"0x0bc4b0c3483084bb71614e114968c1a0ae588888\",\"to\":\"0x0bc4b0c3483084bb71614e114968c1a0ae588888\",\"token\":2,\"amount\":12312124,\"fee\":\"10000\",\"nonce\":20,\"signature\":null,\"validFrom\":0,\"validUntil\":10000000000000000}"
	require.Equal(t, expected, string(txBytes))
}

func TestCreateTransferTx(t *testing.T) {
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
	txBytes, err := json.Marshal(tx)
	require.NoError(t, err)
	expected := "{\"type\":\"Transfer\",\"accountId\":1291712,\"from\":\"0x0bc4b0c3483084bb71614e114968c1a0ae588888\",\"to\":\"0x0bc4b0c3483084bb71614e114968c1a0ae588888\",\"token\":2,\"amount\":12312124,\"fee\":\"10000\",\"nonce\":20,\"signature\":null,\"validFrom\":0,\"validUntil\":10000000000000000}"
	require.Equal(t, expected, string(txBytes))
}

func TestCreateTransferWithFeeTokenTx(t *testing.T) {
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
	txBytes, err := json.Marshal(txs)
	require.NoError(t, err)
	expected := "[{\"type\":\"Transfer\",\"accountId\":1291712,\"from\":\"0x0bc4b0c3483084bb71614e114968c1a0ae588888\",\"to\":\"0x0e81575BF66e79915A22c614e2046d360e40a3f9\",\"token\":2,\"amount\":12312124,\"fee\":\"0\",\"nonce\":18,\"signature\":null,\"validFrom\":0,\"validUntil\":10000000000000000},{\"type\":\"Transfer\",\"accountId\":1291712,\"from\":\"0x0bc4b0c3483084bb71614e114968c1a0ae588888\",\"to\":\"0x0bc4b0c3483084bb71614e114968c1a0ae588888\",\"token\":1,\"amount\":0,\"fee\":\"10000\",\"nonce\":19,\"signature\":null,\"validFrom\":0,\"validUntil\":10000000000000000}]"
	require.Equal(t, expected, string(txBytes))
}

func TestCreateChangePubKeyTx(t *testing.T) {
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
	bytes, err := json.Marshal(tx)
	require.NoError(t, err)
	expected := "{\"type\":\"ChangePubKey\",\"accountId\":1291712,\"account\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"newPkHash\":\"sync:4d9badffbce878e14d60cbd57a90d07a88c7028f\",\"feeToken\":2,\"fee\":\"100000000000000\",\"nonce\":2,\"signature\":null,\"ethAuthData\":{\"type\":\"Onchain\"},\"validFrom\":0,\"validUntil\":4294967295}"
	require.Equal(t, expected, string(bytes))
}

func TestCreateMintNFTTx(t *testing.T) {
	ethPrivKeyHex := "559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"
	accountId := uint32(1291712)
	address, _ := GetAddress(ethPrivKeyHex)
	token := RinkebyUSDC
	fee := big.NewInt(100000000000000)
	nonce := uint32(29)
	NFTContentHash := "1"
	tx := CreateMintNFTTx(accountId, address, address, NFTContentHash, token, fee, nonce)
	bytes, err := json.Marshal(tx)
	require.NoError(t, err)
	expected := "{\"type\":\"MintNFT\",\"creatorId\":1291712,\"creatorAddress\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"contentHash\":\"0x0000000000000000000000000000000000000000000000000000000000000001\",\"recipient\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"fee\":\"100000000000000\",\"feeToken\":2,\"nonce\":29,\"signature\":null}"
	require.Equal(t, expected, string(bytes))
}

func TestCreateTransferNFTTx(t *testing.T) {
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
	require.NoError(t, err)
	bytes, err := json.Marshal(transfers)
	require.NoError(t, err)
	expected := "[{\"type\":\"Transfer\",\"accountId\":1291712,\"from\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"to\":\"0x0e81575BF66e79915A22c614e2046d360e40a3f9\",\"token\":113561,\"amount\":1,\"fee\":\"0\",\"nonce\":16,\"signature\":null,\"validFrom\":0,\"validUntil\":10000000000000000},{\"type\":\"Transfer\",\"accountId\":1291712,\"from\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"to\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"token\":2,\"amount\":0,\"fee\":\"100000000000000\",\"nonce\":17,\"signature\":null,\"validFrom\":0,\"validUntil\":10000000000000000}]"
	require.Equal(t, expected, string(bytes))
}

func TestCreateWithdrawNFTTx(t *testing.T) {
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
	bytes, err := json.Marshal(tx)
	require.NoError(t, err)
	expected := "{\"type\":\"WithdrawNFT\",\"accountId\":1291712,\"from\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"to\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"token\":113561,\"feeToken\":2,\"fee\":\"100000000000000\",\"nonce\":16,\"signature\":null,\"validFrom\":0,\"validUntil\":10000000000000000}"
	require.Equal(t, expected, string(bytes))

}
