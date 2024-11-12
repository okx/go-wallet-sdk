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
	from := "0x076a18ceb1638b364b2bccd7652b3d024b0192b6cd97932d7a25638cd0c38cc3"
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

	assert.Equal(t, `{"type":"INVOKE_FUNCTION","sender_address":"0x076a18ceb1638b364b2bccd7652b3d024b0192b6cd97932d7a25638cd0c38cc3","calldata":["1","2087021424722619777119509474943472645767659996348769578120564519014510906823","232670485425082704932579856502088130646006032362877466777181098476241604910","0","3","3","3059801216421328751596122112822479687228268238922911799033205908290402847460","1700000000000000","0"],"max_fee":"0xcbba106e000","signature":["3401628654028065778841119644550530891055048890957579557528273384328460537572","2931185323066124925754585519252992698493396142856989191034518059513125992525"],"version":"0x1","nonce":"0x1"}`, string(b))
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

	nonce, err := HexToBN(nonceResp.Result[0])
	if err != nil {
		return nil, err
	}
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
	from := "0x076a18ceb1638b364b2bccd7652b3d024b0192b6cd97932d7a25638cd0c38cc3"
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
	assert.Equal(t, `{"type":"INVOKE_FUNCTION","sender_address":"0x076a18ceb1638b364b2bccd7652b3d024b0192b6cd97932d7a25638cd0c38cc3","calldata":["1","3256441166037631918262930812410838598500200462657642943867372734773841898370","403823062618199777388530751713272716715733872218085068081490028803159187238","0","3","3","564521648175006025532572708057195208089056127251","1000000000000000","0"],"max_fee":"0x69f95cc4c98ae","signature":["997530001645826245217076867597411200842987068991821492017648518099258243988","1328037646647848873512528341661941339565391129262876084553664272191049071268"],"version":"0x1","nonce":"0x2"}`, string(b))
}

func TestCreateMutiContracTx(t *testing.T) {
	curve := SC()
	contractAddr := ETH
	functionName := "0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e"
	from := "0x04f46d2B784A75d85163364930c941116664f272c8b96D70491dB228B1d20daa"

	to1 := "0x026e9E8c411056B64B2D044EBCb39FC810D652Cfbe694326651d796BB078320b"
	to2 := "0x004eb36472e15019967568f5D09eAF985e4CaC8Cce3CD6c1930841442270A582"
	maxFee := big.NewInt(1864315586779310)
	nonce := big.NewInt(3)

	txs := []Calls{
		{
			ContractAddress: contractAddr,
			Entrypoint:      functionName,
			Calldata:        []string{to1, "0x38d7ea4c68000", "0"},
		}, {
			ContractAddress: contractAddr,
			Entrypoint:      functionName,
			Calldata:        []string{to2, "0x38d7ea4c68000", "0"},
		},
	}

	tx, err := CreateSignedMultiContractTx(curve, from, txs, nonce, maxFee, GOERLI_ID, "0x0603c85d20500520d4c653352ff6c524f358afeab7e41a511c73733e49c3075e")
	if err != nil {
		t.Fatal(err)
	}
	request := tx.GetTxRequest()
	b, _ := json.Marshal(request)
	t.Logf(string(b))
	assert.Equal(t, "{\"type\":\"INVOKE_FUNCTION\",\"sender_address\":\"0x04f46d2b784a75d85163364930c941116664f272c8b96d70491db228b1d20daa\",\"calldata\":[\"2\",\"2087021424722619777119509474943472645767659996348769578120564519014510906823\",\"232670485425082704932579856502088130646006032362877466777181098476241604910\",\"0\",\"3\",\"2087021424722619777119509474943472645767659996348769578120564519014510906823\",\"232670485425082704932579856502088130646006032362877466777181098476241604910\",\"3\",\"3\",\"6\",\"1100073131459501680801927467743186870973801404098697873181544877894944698891\",\"1000000000000000\",\"0\",\"139052191741745800914583662380955970011320307658408865736700552519382771074\",\"1000000000000000\",\"0\"],\"max_fee\":\"0x69f95cc4c98ae\",\"signature\":[\"760447045916449177622851549765390799754187040144222481889821032357924640812\",\"2927285097903583350907810664360564794371093466717502011984475989141773484812\"],\"version\":\"0x1\",\"nonce\":\"0x3\"}", string(b))
}

func TestCreateSignedUpgradeTx(t *testing.T) {
	curve := SC()

	pri := "0x6ee700c6032c2b3032c548744032a2f04efedeb44434860375b64de8eaeaca8"
	from := "0x078e4b47a9039490e384daf45fdd907c563ba47f076755824e9dc4bfd9a090a4"
	maxFee := big.NewInt(101360058727033)
	nonce := big.NewInt(0)

	tx, err := CreateSignedUpgradeTx(curve, from, nonce, maxFee, MAINNET_ID, pri)
	if err != nil {
		t.Fatal(err)
	}
	request := tx.GetOldTxRequest()
	bytes, _ := json.Marshal(request)
	t.Logf(string(bytes))
	assert.Equal(t, `{"type":"INVOKE_FUNCTION","contract_address":"0x78e4b47a9039490e384daf45fdd907c563ba47f076755824e9dc4bfd9a090a4","calldata":["1","3417601786212868110109684890481242235888924928954829564265128058049878069412","429286934060636239444256046255241512105662385954349596568652644383873724621","0","1","1","1449178161945088530446351771646113898511736767359683664273252560520029776866","0"],"entry_point_selector":"0x15d40a3d6ca2ac30f4031e42be28da9b056fef9bb7357ac5e85627ee876e5ad","max_fee":"0x5c2fba4b7a79","signature":["3357525698737071998813962926394033071725708642828196663017857258566474124153","2490443692580058414860821492828517520325641579361425628600243520322083164202"]}`, string(bytes))
}

func TestGetTxHash(t *testing.T) {
	txStr := `{"type":"INVOKE_FUNCTION","sender_address":"0x06c3c93eeb1643740a80a338b9346c0c9a06177bfcc098a6d86e353532090ae4","calldata":["1","3256441166037631918262930812410838598500200462657642943867372734773841898370","403823062618199777388530751713272716715733872218085068081490028803159187238","0","3","3","564521648175006025532572708057195208089056127251","1000000000000000","0"],"max_fee":"0x69f95cc4c98ae","signature":["847473586541842316388942211795213889856494548988837959760160024500693390782","1348638286841361823893095410439312197628401195006599391371680622656774652575"],"version":"0x1","nonce":"0x2"}`
	deployStr := `{"type":"DEPLOY_ACCOUNT","contract_address_salt":"0x2f4a65ecea5351f49f181841bdddcdf62f600d0e4864755699386d42dd17e37","constructor_calldata":["1374167106255892599010711965180388247554893597343032596700351269194389035468","215307247182100370520050591091822763712463273430149262739280891880522753123","2","1336884626863307009745693974738944585680195300936188147148938838915943595575","0"],"class_hash":"0x3530cc4759d78042f1b543bf797f5f3d647cde0388c33734cf91b7f7b9314a9","max_fee":"0x7157cb0e14a0","version":"0x1","nonce":"0x0","signature":["1743576707672350586938093874140587768903567601625974071199004868774070770998","2517494932084439140630351310818252639109372374885508507240315248980355503830"]}
`

	fmt.Println(GetTxHash(txStr))
	fmt.Println(GetTxHash(deployStr))
}

func TestVerifyMessageHash(t *testing.T) {
	curve := SC()
	pub := "0x346262ffa4ec2f40feb9ae81e416af7cca9fcfa8871f1f9169e6dccd63aa667"
	pubX, pubY := curve.XToPubKey(pub)

	messageHash := "0x45514f85d4e7e2d3db3aac059a5d937f6c5d0f61f87ba25fa138c038248ce7a"
	sigR := "03da226bd3985e75d344ecc653967fdf13647773a537845ae8cbd9c62e3e5208"
	sigS := "04c6be1b89be6f86ad3cdb13ac5b55c05df9b1076a0adce73c058225a652c17a"

	b := curve.Verify(HexToBig(messageHash), HexToBig(sigR), HexToBig(sigS), pubX, pubY)
	assert.True(t, b)
}
