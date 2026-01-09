package cairo1

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/okx/go-wallet-sdk/coins/starknet"
	"math/big"
	"testing"
)

func TestCalculateContractAddressFromHash(t *testing.T) {
	sc := starknet.SC()
	privateKey := "0x043c269ab7bce618cf53c990db2f56c9223952da735f30c064fcd96a681084b9"
	publicKey, err := starknet.GetPubKey(sc, privateKey)
	assert.NoError(t, err)

	result := []string{
		"0x05a3e31056b5db28c67e7ca7b0140eea4c746ee174a5c797714ef6f052bd6a0b",
		"0x04daf5f47959134734ea09d4d71f3aa6189d9b66b607a5ab6d73322b81026a7c",
		"0x0722808348404c439d75c3de9e027d3b40113460a976fad14f41cd02b7b52cf8",
	}
	for i := 0; i < len(starknet.ArgentClassHashCairo1); i++ {
		calculateAddress, err := CalculateContractAddressFromHash(publicKey, starknet.ArgentClassHashCairo1[i])
		assert.NoError(t, err)
		preAddress := starknet.BigToHexWithPadding(calculateAddress)
		assert.Equal(t, result[i], preAddress)
	}

}

func TestCreateSignedDeployAccountTx(t *testing.T) {
	curve := starknet.SC()
	pri := "0x01a820094d8a382db8f6b78a84eb70d909f4344b98a497c52e78de19e049f2da"

	starkPub, err := starknet.GetPubKey(curve, pri)
	assert.NoError(t, err)
	nonce := big.NewInt(0)
	maxFee := big.NewInt(217231761102532)

	tx, err := CreateDeployAccountTx(starkPub, starknet.OKXAccountClassHashCairo1, nonce, maxFee, starknet.MAINNET_ID)
	assert.NoError(t, err)

	err = starknet.SignTx(curve, tx, pri)
	assert.NoError(t, err)

	req := tx.GetDeployAccountReq()
	jsonReq, err := json.Marshal(req)
	assert.NoError(t, err)
	assert.Equal(t, "0x5a510a37b6e99f684104172b35187626ddf5551ec3dccedc5e4e52e8fe90bf6", starknet.BigToHex(tx.TransactionHash))
	assert.Equal(t, `{"type":"DEPLOY_ACCOUNT","contract_address_salt":"0x72ff9867ba607f204042c328cde87ddefe405b830e6515563fbe3ced9342109","constructor_calldata":["3250953912622366112757261009292357488752618517996924917314048668568370028809","0"],"class_hash":"0x1c0bb51e2ce73dc007601a1e7725453627254016c28f118251a71bbb0507fcb","max_fee":"0xc59235f456c4","version":"0x1","nonce":"0x0","signature":["2766731716724855679795418940476970091597219837260156056026575342854911781563","529998056881645898391433349383052287083635576224272450659107534975419683759"]}`, string(jsonReq))
	// https://starkscan.co/tx/0x05a510a37b6e99f684104172b35187626ddf5551ec3dccedc5e4e52e8fe90bf6
}
