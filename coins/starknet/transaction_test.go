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
	from := "0x6c3c93eeb1643740a80a338b9346c0c9a06177bfcc098a6d86e353532090ae4"
	to := "0x026e9E8c411056B64B2D044EBCb39FC810D652Cfbe694326651d796BB078320b"
	amount := big.NewInt(1700000000000000)
	maxFee := big.NewInt(14000000000000)
	nonce := big.NewInt(1)

	tx, err := CreateTransferTx(curve, contractAddr, from, to, amount, nonce, maxFee, MAINNET_ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := SignTx(curve, tx, "//todo please replace your key"); err != nil {
		t.Fatal(err)
	}

	request := tx.GetTxRequest()
	b, _ := json.Marshal(request)
	t.Logf(string(b))

	assert.Equal(t, []byte(`{"type":"INVOKE_FUNCTION","sender_address":"0x6c3c93eeb1643740a80a338b9346c0c9a06177bfcc098a6d86e353532090ae4","calldata":["1","2087021424722619777119509474943472645767659996348769578120564519014510906823","232670485425082704932579856502088130646006032362877466777181098476241604910","0","3","3","1100073131459501680801927467743186870973801404098697873181544877894944698891","1700000000000000","0"],"max_fee":"0xcbba106e000","signature":["2046726132177223766402968530695733227169557314964382677711945876686369864286","1790217669415917429705589389033333595673730480755481198246314553747591090991"],"version":"0x1","nonce":"0x1"}`), b)
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
	t.Log(GetPubKeyPoint(sc, "//todo please replace your key"))
	t.Log(SignMsg(sc, "0xb0a391057a8c2ce9a6e8799f2609da2012970a513a700960e68f05c5c0cc26", "//todo please replace your key"))
}

func TestCreateContractTx(t *testing.T) {
	curve := SC()
	contractAddr := ETHBridge
	from := "0x06c3c93eeb1643740a80a338b9346c0c9a06177bfcc098a6d86e353532090ae4"
	maxFee := big.NewInt(1864315586779310)
	nonce := big.NewInt(2)
	functionName := "initiate_withdraw"
	calldata := []string{"0x62e206b4ddd402056d881ded58c0bd87193d2913", "0x38d7ea4c68000"}

	tx, err := CreateSignedContractTx(curve, contractAddr, from, functionName, calldata, nonce, maxFee, MAINNET_ID, "//todo please replace your key")
	if err != nil {
		t.Fatal(err)
	}

	request := tx.GetTxRequest()
	b, _ := json.Marshal(request)
	t.Logf(string(b))
	assert.Equal(t, `{"type":"INVOKE_FUNCTION","sender_address":"0x06c3c93eeb1643740a80a338b9346c0c9a06177bfcc098a6d86e353532090ae4","calldata":["1","3256441166037631918262930812410838598500200462657642943867372734773841898370","403823062618199777388530751713272716715733872218085068081490028803159187238","0","3","3","564521648175006025532572708057195208089056127251","1000000000000000","0"],"max_fee":"0x69f95cc4c98ae","signature":["847473586541842316388942211795213889856494548988837959760160024500693390782","1348638286841361823893095410439312197628401195006599391371680622656774652575"],"version":"0x1","nonce":"0x2"}`, string(b))
}
