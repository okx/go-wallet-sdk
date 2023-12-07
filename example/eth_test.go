package example

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExample(t *testing.T) {
	// get menmonic
	mnemonic, err := GenerateMnemonic()
	assert.NoError(t, err)
	fmt.Println(mnemonic)
	// get derived key
	hdPath := GetDerivedPath(0)
	derivePrivateKey, err := GetDerivedPrivateKey(mnemonic, hdPath)
	assert.NoError(t, err)
	fmt.Println("generate derived private key:", derivePrivateKey, ",derived path: ", hdPath)

	// get new address
	newAddress := GetNewAddress(derivePrivateKey)
	fmt.Println("generate new address:", newAddress)

	// Verify address
	valid := ValidAddress(newAddress)
	fmt.Println("verify address isValid:", valid)

	// Sign a transaction
	txJson := `{
				"chainId":"1",
				"txType":2,
				"nonce":"244",
				"isToken":false,
				"to":"0x31c514837ee0f6062eaffb0882d764170a178004",
				"value":"1000000000000000",
				"gasLimit":"21000",
				"maxFeePerGas":"20000000000",
				"maxPriorityFeePerGas":"1500000000"
			}`

	signedTx, err := SignTransaction(txJson, derivePrivateKey)
	assert.NoError(t, err)
	fmt.Println("signed tx:", signedTx)
}

func TestSignTx(t *testing.T) {
	// get menmonic
	mnemonic := "limb alter vapor lava clown pigeon exist pulp ride dry wage middle battle tell suspect pigeon want thrive sugar smoke merit tower curve local"
	fmt.Println(mnemonic)
	// get derived key
	hdPath := GetDerivedPath(0)
	derivePrivateKey, err := GetDerivedPrivateKey(mnemonic, hdPath)
	assert.NoError(t, err)
	fmt.Println("generate derived private key:", derivePrivateKey, ",derived path: ", hdPath)

	// get new address
	newAddress := GetNewAddress(derivePrivateKey)
	fmt.Println("generate new address:", newAddress)
	assert.Equal(t, "0xd5cb882a2ace84806c0554c247f8135d161378b4", newAddress)
}
