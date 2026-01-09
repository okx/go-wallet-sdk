package example

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignTx(t *testing.T) {
	mnemonic := "limb alter vapor lava clown pigeon exist pulp ride dry wage middle battle tell suspect pigeon want thrive sugar smoke merit tower curve local"

	hdPath := GetDerivedPath(0)
	derivePrivateKey, err := GetDerivedPrivateKey(mnemonic, hdPath)
	assert.NoError(t, err)

	newAddress := GetNewAddress(derivePrivateKey)
	assert.Equal(t, "0xd5cb882a2ace84806c0554c247f8135d161378b4", newAddress)

	valid := ValidAddress(newAddress)
	assert.True(t, valid)

	txJson := `{
		"txType":2,
		"chainId":"11155111",
		"nonce":"1",
		"to":"0x31c514837ee0f6062eaffb0882d764170a178004",
		"value":"21000",
		"gasLimit":"21000",
		"gasPrice":"66799178286",
		"maxFeePerGas":"20000000000",
		"maxPriorityFeePerGas":"1500000000"
	}`
	signedTx, err := SignTransaction(txJson, derivePrivateKey)
	assert.NoError(t, err)
	assert.Equal(t, `0x02f87083aa36a7018459682f008504a817c8008252089431c514837ee0f6062eaffb0882d764170a17800482520880c001a0b7588bed05e60cd5edbef6ac9cde46bc5807d4d3d538fbd6f4a6081b161bc8e5a049d5992df046a5c97c77cf6a6a23982b9a51af240b22a755ca1089addcc3607e`, signedTx)

	hash, err := CalTxHash(signedTx)
	assert.NoError(t, err)
	assert.Equal(t, "6673e25ced49eb2160d676db950837ab3280955b68b4b2eea05f124ec0ed6942", hash)
}
