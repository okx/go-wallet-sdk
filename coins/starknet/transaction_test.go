package starknet

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
	"testing"
)

func TestJsonRpc(t *testing.T) {
	//create transaction
	curve := SC()
	contractAddr := ETH
	from := "0x027850700bb0c1a9fe7c4dc7c253548e40f4b4fcc4d36f68551a557b19c0b3a2"
	to := "0x6c3c93eeb1643740a80a338b9346c0c9a06177bfcc098a6d86e353532090ae4"
	amount := big.NewInt(1700000000000000)
	maxFee := big.NewInt(14000000000000)
	nonce := big.NewInt(1)

	tx, err := CreateTransferTx(curve, contractAddr, from, to, amount, nonce, maxFee, MAINNET_ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := SignTx(curve, tx, "0x01651242558d251b0daa72cdf11feb1713e47eb88fb55d0978a2625445a771ac"); err != nil {
		t.Fatal(err)
	}

	request := tx.GetTxRequest()
	b, _ := json.Marshal(request)
	t.Logf(string(b))

	assert.Equal(t, `{"type":"INVOKE_FUNCTION","sender_address":"0x027850700bb0c1a9fe7c4dc7c253548e40f4b4fcc4d36f68551a557b19c0b3a2","calldata":["1","2087021424722619777119509474943472645767659996348769578120564519014510906823","232670485425082704932579856502088130646006032362877466777181098476241604910","0","3","3","3059801216421328751596122112822479687228268238922911799033205908290402847460","1700000000000000","0"],"max_fee":"0xcbba106e000","signature":["3316003862596213685187927938396555541641565477986172535417248434450520075036","207497642185058112223994813923795045316971441670508760190438453143208354315"],"version":"0x1","nonce":"0x1"}`, string(b))
}

func getNonce(address string, net string) (*big.Int, error) {

	nonceRequest := struct {
		ContractAddress    string   `json:"contract_address,omitempty"`
		EntryPointSelector string   `json:"entry_point_selector,omitempty"`
		Calldata           []string `json:"calldata"`
		Signature          []string `json:"signature"`
	}{
		ContractAddress:    address,
		EntryPointSelector: BigToHex(GetSelectorFromName("get_nonce")),
		Calldata:           []string{},
		Signature:          []string{},
	}

	url := fmt.Sprintf("%s/feeder_gateway/call_contract", net)
	contentType := "appliaction/json"
	bodyBytes, _ := json.Marshal(nonceRequest)

	fmt.Printf(string(bodyBytes))

	client := http.Client{}
	response, err := client.Post(url, contentType, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, err
	}

	respBodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	nonceResp := &struct {
		Result []string `json:"result"`
	}{}
	err = json.Unmarshal(respBodyBytes, nonceResp)
	if err != nil {
		return nil, err
	}

	nonce := HexToBN(nonceResp.Result[0])
	return nonce, nil
}

func TestStarkEx(t *testing.T) {
	sc := SC()
	t.Log(GetPubKeyPoint(sc, "0x01651242558d251b0daa72cdf11feb1713e47eb88fb55d0978a2625445a771ac"))
	t.Log(SignMsg(sc, "0xb0a391057a8c2ce9a6e8799f2609da2012970a513a700960e68f05c5c0cc26", "0x01651242558d251b0daa72cdf11feb1713e47eb88fb55d0978a2625445a771ac"))
}

func TestCreateContractTx(t *testing.T) {
	curve := SC()
	contractAddr := ETHBridge
	from := "0x027850700bb0c1a9fe7c4dc7c253548e40f4b4fcc4d36f68551a557b19c0b3a2"
	maxFee := big.NewInt(1864315586779310)
	nonce := big.NewInt(2)
	functionName := "initiate_withdraw"
	calldata := []string{"0x62e206b4ddd402056d881ded58c0bd87193d2913", "0x38d7ea4c68000"}

	tx, err := CreateSignedContractTx(curve, contractAddr, from, functionName, calldata, nonce, maxFee, MAINNET_ID, "0x01651242558d251b0daa72cdf11feb1713e47eb88fb55d0978a2625445a771ac")
	if err != nil {
		t.Fatal(err)
	}

	request := tx.GetTxRequest()
	b, _ := json.Marshal(request)
	t.Logf(string(b))
	assert.Equal(t, `{"type":"INVOKE_FUNCTION","sender_address":"0x027850700bb0c1a9fe7c4dc7c253548e40f4b4fcc4d36f68551a557b19c0b3a2","calldata":["1","3256441166037631918262930812410838598500200462657642943867372734773841898370","403823062618199777388530751713272716715733872218085068081490028803159187238","0","3","3","564521648175006025532572708057195208089056127251","1000000000000000","0"],"max_fee":"0x69f95cc4c98ae","signature":["205010620035651583898118207596984197687754972756456862120522896476547105292","1154716553731371905121479884741907505965615780937256245747265484914503769789"],"version":"0x1","nonce":"0x2"}`, string(b))
}
