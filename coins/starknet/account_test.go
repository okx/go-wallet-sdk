package starknet

import (
	"encoding/json"
	"fmt"
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

	bn := HexToBN(publicKey)
	t.Logf("public key bn : %s", bn.String())
}

func TestCalculateContractAddressFromHash(t *testing.T) {
	curve := SC()
	starkPub, _ := GetPubKey(curve, "//todo please replace your key")
	fmt.Println(starkPub)
	calculateAddress, err := CalculateContractAddressFromHash(starkPub)
	if err != nil {
		t.Fatal(err)
	}
	preAddress := BigToHexWithPadding(calculateAddress)
	assert.Equal(t, "0x06c3c93eeb1643740a80a338b9346c0c9a06177bfcc098a6d86e353532090ae4", preAddress)
	fmt.Println(preAddress)
}

func TestCreateSignedDeployAccountTx(t *testing.T) {
	curve := SC()
	pri := "//todo please replace your key"

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
	assert.Equal(t, `{"type":"DEPLOY_ACCOUNT","contract_address_salt":"0x2f4a65ecea5351f49f181841bdddcdf62f600d0e4864755699386d42dd17e37","constructor_calldata":["1374167106255892599010711965180388247554893597343032596700351269194389035468","215307247182100370520050591091822763712463273430149262739280891880522753123","2","1336884626863307009745693974738944585680195300936188147148938838915943595575","0"],"class_hash":"0x3530cc4759d78042f1b543bf797f5f3d647cde0388c33734cf91b7f7b9314a9","max_fee":"0x7157cb0e14a0","version":"0x1","nonce":"0x0","signature":["1743576707672350586938093874140587768903567601625974071199004868774070770998","2517494932084439140630351310818252639109372374885508507240315248980355503830"]}`, string(jsonReq))
	fmt.Println(string(jsonReq))
}
