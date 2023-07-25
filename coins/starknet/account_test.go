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
	starkPub, _ := GetPubKey(curve, "0x01651242558d251b0daa72cdf11feb1713e47eb88fb55d0978a2625445a771ac")
	fmt.Println(starkPub)
	calculateAddress, err := CalculateContractAddressFromHash(starkPub)
	if err != nil {
		t.Fatal(err)
	}
	preAddress := BigToHexWithPadding(calculateAddress)
	assert.Equal(t, "0x027850700bb0c1a9fe7c4dc7c253548e40f4b4fcc4d36f68551a557b19c0b3a2", preAddress)
	fmt.Println(preAddress)
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
	assert.Equal(t, `{"type":"DEPLOY_ACCOUNT","contract_address_salt":"0x4c3eb6ed976748ba6038adc996bc7efbae3915e71cb5ad9af3ae839aa5fe28e","constructor_calldata":["1449178161945088530446351771646113898511736767359683664273252560520029776866","215307247182100370520050591091822763712463273430149262739280891880522753123","2","2155411470851976624741272041507081168444512578889852119281177422738732606094","0"],"class_hash":"0x25ec026985a3bf9d0cc1fe17326b245dfdc3ff89b8fde106542a3ea56c5a918","max_fee":"0x7157cb0e14a0","version":"0x1","nonce":"0x0","signature":["2937863057704564206936475349638492456044380740625858385484870946711142421031","2715437793949255530037338139693287898654578972743579624372647390009359438003"]}`, string(jsonReq))
	fmt.Println(string(jsonReq))
}
