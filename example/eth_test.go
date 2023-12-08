package example

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// testcode https://sepolia.etherscan.io/address/0x2354aa76877c6043bf30119eb23bf3fdd02c4808
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
				"chainId":"11155111",
				"txType":2,
				"nonce":"1",
				"isToken":false,
				"to":"0x31c514837ee0f6062eaffb0882d764170a178004",
				"value":"21000",
				"gasLimit":"21000",
				"gasPrice":"66799178286",
				"maxFeePerGas":"20000000000",
				"maxPriorityFeePerGas":"1500000000"
			}`
	//02a501a1622ecdbdca2ff9ae36dfcf58603006e8fd5ddd4809e8b8b9b8a4cf9f8b
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

	// Verify address
	valid := ValidAddress(newAddress)
	fmt.Println("verify address isValid:", valid)

	// Sign a transaction of type 1
	txJson := `{
				"chainId":"11155111",
				"txType":2,
				"nonce":"1",
				"isToken":false,
				"to":"0x31c514837ee0f6062eaffb0882d764170a178004",
				"value":"21000",
				"gasLimit":"21000",
				"gasPrice":"66799178286",
				"maxFeePerGas":"20000000000",
				"maxPriorityFeePerGas":"1500000000"
			}`
	//02a501a1622ecdbdca2ff9ae36dfcf58603006e8fd5ddd4809e8b8b9b8a4cf9f8b
	signedTx, err := SignTransaction(txJson, derivePrivateKey)
	assert.NoError(t, err)
	fmt.Println("signed tx1:", signedTx)
	assert.Equal(t, signedTx, "{\"hash\":\"0x6673e25ced49eb2160d676db950837ab3280955b68b4b2eea05f124ec0ed6942\",\"hex\":\"0x02f87083aa36a7018459682f008504a817c8008252089431c514837ee0f6062eaffb0882d764170a17800482520880c001a0b7588bed05e60cd5edbef6ac9cde46bc5807d4d3d538fbd6f4a6081b161bc8e5a049d5992df046a5c97c77cf6a6a23982b9a51af240b22a755ca1089addcc3607e\"}")

	// Sign a transaction of type 1
	txJson2 := `{
				"chainId":"11155111",
				"txType":2,
				"nonce":"1",
				"isToken":false,
				"to":"0x31c514837ee0f6062eaffb0882d764170a178004",
				"value":"21000",
				"gasLimit":"21000",
				"gasPrice":"66799178286",
				"maxFeePerGas":"20000000000",
				"maxPriorityFeePerGas":"1500000000"
			}`
	//02a501a1622ecdbdca2ff9ae36dfcf58603006e8fd5ddd4809e8b8b9b8a4cf9f8b
	signedTx2, err := SignTransaction(txJson2, derivePrivateKey)
	assert.NoError(t, err)
	fmt.Println("signed tx2:", signedTx2)
	assert.Equal(t, signedTx2, "{\"hash\":\"0x6673e25ced49eb2160d676db950837ab3280955b68b4b2eea05f124ec0ed6942\",\"hex\":\"0x02f87083aa36a7018459682f008504a817c8008252089431c514837ee0f6062eaffb0882d764170a17800482520880c001a0b7588bed05e60cd5edbef6ac9cde46bc5807d4d3d538fbd6f4a6081b161bc8e5a049d5992df046a5c97c77cf6a6a23982b9a51af240b22a755ca1089addcc3607e\"}")

}
