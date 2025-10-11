package zksync

import (
	"encoding/hex"
	"encoding/json"
	"github.com/emresenyuva/go-wallet-sdk/coins/zksync/core"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestAccount(t *testing.T) {

	privKeyBytes, _ := hex.DecodeString("559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a")
	ethSigner, err := core.NewOkEthSignerFromPrivBytes(privKeyBytes)
	require.NoError(t, err)
	require.Equal(t, "0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728", ethSigner.GetAddress())
	zkSigner, err := core.NewZkSignerFromEthSigner(ethSigner, core.ChainIdRinkeby)
	require.NoError(t, err)
	require.Equal(t, "sync:4d9badffbce878e14d60cbd57a90d07a88c7028f", zkSigner.GetPublicKeyHash())
}

func TestGetPubKeyHash(t *testing.T) {
	ethPrivKeyHex := "559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"
	pubKeyHash, err := GetPubKeyHash(ethPrivKeyHex, int(core.ChainIdRinkeby))
	require.NoError(t, err)
	require.Equal(t, "sync:4d9badffbce878e14d60cbd57a90d07a88c7028f", pubKeyHash)
}

func TestGetAddress(t *testing.T) {
	ethPrivKeyHex := "559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"
	address, err := GetAddress(ethPrivKeyHex)
	require.NoError(t, err)
	require.Equal(t, "0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728", address)
}

func TestSignChangePubKey(t *testing.T) {
	ethPrivKeyHex := "559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"
	accountId := uint32(1292869)
	address, _ := GetAddress(ethPrivKeyHex)
	pubKeyHash, _ := GetPubKeyHash(ethPrivKeyHex, int(core.ChainIdRinkeby))
	fee := big.NewInt(100000000000000)
	nonce := uint32(0)
	validFrom := uint64(0)
	validUntil := uint64(4294967295)
	tx := CreateChangePubKeyTx(accountId, address, pubKeyHash, RinkebyUSDC, fee, nonce, validFrom, validUntil)
	signedTx, err := SignChangePubKey(tx, ethPrivKeyHex, int(core.ChainIdRinkeby))
	require.NoError(t, err)
	txBytes, err := json.Marshal(signedTx)
	require.NoError(t, err)
	expected := "{\"tx\":{\"type\":\"ChangePubKey\",\"accountId\":1292869,\"account\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"newPkHash\":\"sync:4d9badffbce878e14d60cbd57a90d07a88c7028f\",\"feeToken\":2,\"fee\":\"100000000000000\",\"nonce\":0,\"signature\":{\"pubKey\":\"e570ffa4c84b298bac4b881d3570ade9a709e57df3d597413d82f89b83172c23\",\"signature\":\"3bbd7bbdad8f7b4b63d57342f8f8a1ccc31044acd6b355e45ee166f5e3148081e7546a6b5f608b77340c334911d787d8d6141867458abad8cded05f0176de801\"},\"ethAuthData\":{\"type\":\"ECDSA\",\"ethSignature\":\"0xa2e314dfae7e74d56c7d084e2366da7a0e77b1e542e4a64690a8739bc40c5bd27fce1d4390cf7ca29f81625d0ac09f2654e573a9ca36163e07125ce2f6b2496a1c\",\"batchHash\":\"0x0000000000000000000000000000000000000000000000000000000000000000\"},\"validFrom\":0,\"validUntil\":4294967295},\"signature\":null}"
	require.Equal(t, expected, string(txBytes))
}

func TestSignWithdraw(t *testing.T) {
	ethPrivKeyHex := "559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"
	accountId := uint32(1291712)
	amount := big.NewInt(1000000)
	fee := big.NewInt(1000000)
	token := RinkebyUSDC
	nonce := uint32(38)
	validFrom := uint64(0)
	validUntil := uint64(4294967295)
	addr, _ := GetAddress(ethPrivKeyHex)
	tx := CreateWithdrawTx(accountId, addr, amount, fee, token, nonce, validFrom, validUntil)

	signedTx, err := SignWithdraw(tx, token, ethPrivKeyHex, int(core.ChainIdRinkeby))
	require.NoError(t, err)
	txBytes, err := json.Marshal(signedTx)
	require.NoError(t, err)
	expected := "{\"tx\":{\"type\":\"Withdraw\",\"accountId\":1291712,\"from\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"to\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"token\":2,\"amount\":1000000,\"fee\":\"1000000\",\"nonce\":38,\"signature\":{\"pubKey\":\"e570ffa4c84b298bac4b881d3570ade9a709e57df3d597413d82f89b83172c23\",\"signature\":\"40fc841aaaf9576d8fc2e1cb6db2c431cf1c0586b517dcadda4a7adbe1997197cadb3fc4f18e28537a4b5c5df12ae41215aef765e01fdc8054b441ea9b908c03\"},\"validFrom\":0,\"validUntil\":4294967295},\"signature\":{\"type\":\"EthereumSignature\",\"signature\":\"0xa823f58932a0e1346c0a78222816571fae152366c8e1b35ffaf7d21b18f47b901b6be16c334a6f873ac5091f2e99f53472cfd8abfb84a116a2d337274bf6073c1b\"}}"
	require.Equal(t, expected, string(txBytes))
}

func TestSignTransfer(t *testing.T) {
	ethPrivKeyHex := "559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"

	accountId := uint32(1291712)
	to := "0x5007B4AD1ca695Bb3f9ef69D32f8F54691be4D14"
	amount := big.NewInt(10000)
	fee := big.NewInt(10000000000000000)
	token := ETH
	nonce := uint32(43)
	validFrom := uint64(0)
	validUntil := uint64(4294967295)
	from, _ := GetAddress(ethPrivKeyHex)
	tx := CreateTransferTx(accountId, from, to, amount, fee, token, nonce, validFrom, validUntil)

	signedTx, err := SignTransfer(tx, ethPrivKeyHex, int(core.ChainIdRinkeby))
	require.NoError(t, err)

	txBytes, err := json.Marshal(signedTx)
	require.NoError(t, err)
	expected := "{\"tx\":{\"type\":\"Transfer\",\"accountId\":1291712,\"from\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"to\":\"0x5007B4AD1ca695Bb3f9ef69D32f8F54691be4D14\",\"token\":0,\"amount\":10000,\"fee\":\"10000000000000000\",\"nonce\":43,\"signature\":{\"pubKey\":\"e570ffa4c84b298bac4b881d3570ade9a709e57df3d597413d82f89b83172c23\",\"signature\":\"014ccbee836e9dadaa59fd3d421bcc512926d3d3cab8b7aef83f29bc8b761f95144d7ab1889238b86073603f1778575795d0b48f06ccfd1c745ad56be3ebfa04\"},\"validFrom\":0,\"validUntil\":4294967295},\"signature\":{\"type\":\"EthereumSignature\",\"signature\":\"0x267554e6eeecd64317f9b241d1ebd1142036f57291c3596836cd14288a750f3319ccceddd6b2c5f9c23df58337f2661b50e5f9e6ec97b19d61d98dda6d4e5f291c\"}}"
	require.Equal(t, expected, string(txBytes))
}

func TestSignTransferWithFeeToken(t *testing.T) {
	ethPrivKeyHex := "559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"

	accountId := uint32(1291712)
	from := "0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728"
	to := "0x0e81575BF66e79915A22c614e2046d360e40a3f9"
	amount := big.NewInt(12312124)
	fee := big.NewInt(10000)
	feeToken := RinkebyUSDC
	nonce := uint32(35)
	validFrom := uint64(0)
	validUntil := uint64(10000000000000000)
	txs := CreateTransferWithFeeTokenTx(accountId, from, to, amount, feeToken, fee, RinkebyUSDC, nonce, validFrom, validUntil)

	signedTx, err := SignBatchTransfer(txs, ethPrivKeyHex, int(core.ChainIdRinkeby))
	require.NoError(t, err)

	txBytes, err := json.Marshal(signedTx)
	require.NoError(t, err)
	expected := "{\"txs\":[{\"tx\":{\"type\":\"Transfer\",\"accountId\":1291712,\"from\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"to\":\"0x0e81575BF66e79915A22c614e2046d360e40a3f9\",\"token\":2,\"amount\":12312124,\"fee\":\"0\",\"nonce\":35,\"signature\":{\"pubKey\":\"e570ffa4c84b298bac4b881d3570ade9a709e57df3d597413d82f89b83172c23\",\"signature\":\"12d7f9b7df5811c21715c2a2ebe84ab37cdb64b4049ca72d383c75dda136c5895a49d0d1c8060d3e87838e084282d30b62852c8fc2f4634259516471ccf36304\"},\"validFrom\":0,\"validUntil\":10000000000000000},\"signature\":null},{\"tx\":{\"type\":\"Transfer\",\"accountId\":1291712,\"from\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"to\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"token\":2,\"amount\":0,\"fee\":\"10000\",\"nonce\":36,\"signature\":{\"pubKey\":\"e570ffa4c84b298bac4b881d3570ade9a709e57df3d597413d82f89b83172c23\",\"signature\":\"9fc3debfbf77130ba2fea92ea50e3966dfd07d8e98074772ee1eb869e3621b2d05d309d3a639ef0e587f5d51ad68680988d4eab02b92acd6cdcaa9e5d42c2c04\"},\"validFrom\":0,\"validUntil\":10000000000000000},\"signature\":null}],\"signature\":{\"type\":\"EthereumSignature\",\"signature\":\"0xb79827822f3fdb86558b14b658aa17a4de89d9b5481121d8a12c4e53e1094a825ebb320f3af8cbb75d13deeca3fa1a5bb76f3fd24fce8f78ea00d785ff5c75961b\"}}"
	require.Equal(t, expected, string(txBytes))
}

func TestSignBatchTransfer(t *testing.T) {
	ethPrivKeyHex := "559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"

	toAddress := "0x0e81575BF66e79915A22c614e2046d360e40a3f9"

	accountId := uint32(1291712)
	address, _ := GetAddress(ethPrivKeyHex)
	fee := big.NewInt(10000)
	feeToken := RinkebyUSDC
	nonce := uint32(37)
	nftSymbol := "NFT-113561"
	validFrom := uint64(0)
	validUntil := uint64(10000000000000000)
	txs, err := CreateTransferNFTTx(accountId, address, toAddress, nftSymbol, feeToken, fee, nonce, validFrom, validUntil)
	require.NoError(t, err)

	signedTx, err := SignBatchTransfer(txs, ethPrivKeyHex, int(core.ChainIdRinkeby))
	require.NoError(t, err)

	txBytes, err := json.Marshal(signedTx)
	require.NoError(t, err)
	expected := "{\"txs\":[{\"tx\":{\"type\":\"Transfer\",\"accountId\":1291712,\"from\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"to\":\"0x0e81575BF66e79915A22c614e2046d360e40a3f9\",\"token\":113561,\"amount\":1,\"fee\":\"0\",\"nonce\":37,\"signature\":{\"pubKey\":\"e570ffa4c84b298bac4b881d3570ade9a709e57df3d597413d82f89b83172c23\",\"signature\":\"513b3e484f51a01d20893d8cfd8ffe971a3a7e064beba67be18ea9025d7e9c1825f6e6378fd6c6e43be8e70c59877ee8c7da917cceee4d256d0c5474c7b92002\"},\"validFrom\":0,\"validUntil\":10000000000000000},\"signature\":null},{\"tx\":{\"type\":\"Transfer\",\"accountId\":1291712,\"from\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"to\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"token\":2,\"amount\":0,\"fee\":\"10000\",\"nonce\":38,\"signature\":{\"pubKey\":\"e570ffa4c84b298bac4b881d3570ade9a709e57df3d597413d82f89b83172c23\",\"signature\":\"85bc29d7f6cdda01e87a49b1cd8356660df070983a09606557058feb3c13ab23a5e72246b9d9069ba5a68d9f24e4be7f2b1ba381b134f7dddbd3da3002a24504\"},\"validFrom\":0,\"validUntil\":10000000000000000},\"signature\":null}],\"signature\":{\"type\":\"EthereumSignature\",\"signature\":\"0x9125ab76b999e5cd44946468b179006172d336020a9df11e6e5b56f782dbd7187f6625ebbd6f8dff1f859d9c3fc99e6ac91f713c2eb7620dc5c81af44bb35fcd1c\"}}"
	require.Equal(t, expected, string(txBytes))
}

func TestSignMintNFT(t *testing.T) {
	ethPrivKeyHex := "559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"
	accountId := uint32(1291712)
	address, _ := GetAddress(ethPrivKeyHex)
	fee := big.NewInt(10000)
	feeToken := RinkebyUSDC
	nonce := uint32(45)
	NFTContentHash := "bafybeieqlewrb6pkogtzvhah5ujz4tbfxihgdcezvduw64mthf77i7akru"
	tx := CreateMintNFTTx(accountId, address, address, NFTContentHash, feeToken, fee, nonce)

	signedTx, err := SignMintNFT(tx, feeToken, ethPrivKeyHex, int(core.ChainIdRinkeby))
	require.NoError(t, err)

	txBytes, err := json.Marshal(signedTx)
	require.NoError(t, err)
	expected := "{\"tx\":{\"type\":\"MintNFT\",\"creatorId\":1291712,\"creatorAddress\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"contentHash\":\"0x0000000000000000000000000000000000000000000000000000000000000baf\",\"recipient\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"fee\":\"10000\",\"feeToken\":2,\"nonce\":45,\"signature\":{\"pubKey\":\"e570ffa4c84b298bac4b881d3570ade9a709e57df3d597413d82f89b83172c23\",\"signature\":\"c03b4dba313e57814267f6da7e48fdb5e3e6621a5144322853eca0e0dccd3c99d0a08fd76f9ef0f9ca13ea790351fea26eb7b4b7e56b76c50ea98bf2c3d9eb00\"}},\"signature\":{\"type\":\"EthereumSignature\",\"signature\":\"0x3ecd4f964e1e106fc4498393a1ac45aac8699b5a544a3e117f21f3f23c0506f32e27d653b62fe88ee9de067de9bf3ccd36893180adc289773699792d33a9bcd81b\"}}"
	require.Equal(t, expected, string(txBytes))
}

func TestSignWithdrawNFT(t *testing.T) {
	ethPrivKeyHex := "559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a"
	accountId := uint32(1291712)
	address, _ := GetAddress(ethPrivKeyHex)
	feeToken := RinkebyUSDC
	fee := big.NewInt(1000000)
	nonce := uint32(41)
	nftId := uint32(113561)
	validFrom := uint64(0)
	validUntil := uint64(10000000000000000)
	tx := CreateWithdrawNFTTx(accountId, address, nftId, feeToken, fee, nonce, validFrom, validUntil)

	signedTx, err := SignWithdrawNFT(tx, feeToken, ethPrivKeyHex, int(core.ChainIdRinkeby))
	require.NoError(t, err)
	txBytes, err := json.Marshal(signedTx)
	require.NoError(t, err)
	expected := "{\"tx\":{\"type\":\"WithdrawNFT\",\"accountId\":1291712,\"from\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"to\":\"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728\",\"token\":113561,\"feeToken\":2,\"fee\":\"1000000\",\"nonce\":41,\"signature\":{\"pubKey\":\"e570ffa4c84b298bac4b881d3570ade9a709e57df3d597413d82f89b83172c23\",\"signature\":\"80044f6fa10417c033529e7bdaea1777b62b0ab2cce34516575ee1a90e61aea738235427a3638f026ef8546936814a1755910f84a3d57c370166c5f78a562f02\"},\"validFrom\":0,\"validUntil\":10000000000000000},\"signature\":{\"type\":\"EthereumSignature\",\"signature\":\"0x5f40f081ced32b78926f7acc2834cc8e642f4b33a3df3a206fd630af4d6e9edb5df3bdf1df8792833cda3447e0ffed0270c445324e5a0ddd6ec10842bae31e6c1c\"}}"
	require.Equal(t, expected, string(txBytes))
}
