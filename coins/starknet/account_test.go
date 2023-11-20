package starknet

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestNewKeyPair(t *testing.T) {
	curve := SC()

	privateKey, publicKey, err := NewKeyPair(curve)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("private key : %s", privateKey)
	t.Logf("public key : %s", publicKey)

	bn, err := HexToBN(publicKey)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("public key bn : %s", bn.String())
}

func TestCalculateContractAddressFromHash(t *testing.T) {
	curve := SC()
	starkPub, _ := GetPubKey(curve, "0x01651242558d251b0daa72cdf11feb1713e47eb88fb55d0978a2625445a771ac")
	calculateAddress, err := CalculateContractAddressFromHash(starkPub)
	if err != nil {
		t.Fatal(err)
	}
	preAddress := BigToHexWithPadding(calculateAddress)
	assert.Equal(t, "0x076a18ceb1638b364b2bccd7652b3d024b0192b6cd97932d7a25638cd0c38cc3", preAddress)
}

func TestCreateSignedDeployAccountTx(t *testing.T) {
	curve := SC()
	pri := "0x01651242558d251b0daa72cdf11feb1713e47eb88fb55d0978a2625445a771ac"
	starkPub, err := GetPubKey(curve, pri)
	if err != nil {
		t.Fatal(err)
	}
	nonce := big.NewInt(0)
	maxFee := big.NewInt(124621882791072)
	tx, err := CreateSignedDeployAccountTx(curve, starkPub, nonce, maxFee, MAINNET_ID, pri)
	if err != nil {
		t.Fatal(err)
	}
	req := tx.GetDeployAccountReq()
	jsonReq, err := json.Marshal(req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, `{"type":"DEPLOY_ACCOUNT","contract_address_salt":"0x4c3eb6ed976748ba6038adc996bc7efbae3915e71cb5ad9af3ae839aa5fe28e","constructor_calldata":["1374167106255892599010711965180388247554893597343032596700351269194389035468","215307247182100370520050591091822763712463273430149262739280891880522753123","2","2155411470851976624741272041507081168444512578889852119281177422738732606094","0"],"class_hash":"0x3530cc4759d78042f1b543bf797f5f3d647cde0388c33734cf91b7f7b9314a9","max_fee":"0x7157cb0e14a0","version":"0x1","nonce":"0x0","signature":["1807287688609244063980416760800578370903099344795965990726567528439781006612","1207725556014061784011398195262770566842867382342357956710841026815216690815"]}`, string(jsonReq))
}

func TestValidAddress(t *testing.T) {
	address := "0x06c3c93eeb1643740a80a338b9346c0c9a06177bfcc098a6d86e353532090ae4"
	assert.Equal(t, true, ValidAddress(address))

	address2 := "0x1127aeb6f4cc7fcfaec0f82722bef78d23acd172d350969c32545e36e0aa4d0b65"
	assert.Equal(t, false, ValidAddress(address2))

	address3 := "6c3c93eeb1643740a80a338b9346c0c9a06177bfcc098a6d86e353532090ae4"
	assert.Equal(t, true, ValidAddress(address3))
}
