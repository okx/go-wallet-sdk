package cairo1

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/okx/go-wallet-sdk/coins/starknet"
	"math/big"
	"testing"
)

func TestCreateTransferTxWithUpgradeAccount(t *testing.T) {
	curve := starknet.SC()
	contractAddr := starknet.ETH
	from := "0x03a15f3065b3d78575dc2ebeea5b49c6958b4966ad06f1edd252ac468d692034"
	to := "0x03a15f3065b3d78575dc2ebeea5b49c6958b4966ad06f1edd252ac468d692034"
	amount := big.NewInt(100)
	maxFee := big.NewInt(1864315586779310)
	nonce := big.NewInt(2)

	tx, err := CreateTransferTx(curve, contractAddr, from, to, amount, nonce, maxFee, starknet.GOERLI_ID)
	assert.NoError(t, err)
	err = starknet.SignTx(curve, tx, "0x01a820094d8a382db8f6b78a84eb70d909f4344b98a497c52e78de19e049f2da")
	assert.NoError(t, err)

	request := tx.GetTxRequest()
	b, _ := json.Marshal(request)
	assert.Equal(t, "0x3bf56f7f6670d9a691579d92306cf5611eeba622c7cdc9f63eb9ddf7f3f595c", starknet.BigToHex(tx.TransactionHash))
	assert.Equal(t, `{"type":"INVOKE_FUNCTION","sender_address":"0x03a15f3065b3d78575dc2ebeea5b49c6958b4966ad06f1edd252ac468d692034","calldata":["1","2087021424722619777119509474943472645767659996348769578120564519014510906823","232670485425082704932579856502088130646006032362877466777181098476241604910","3","1642057893870028285861767096082132425734966191365968126438019791489639981108","100","0"],"max_fee":"0x69f95cc4c98ae","signature":["451616947665820480367602057181276549901241039205427900927743104645439515241","2671458317742571126019934882108631154768632784578200130982031733775952690805"],"version":"0x1","nonce":"0x2"}`, string(b))
	// https://testnet.starkscan.co/tx/0x3bf56f7f6670d9a691579d92306cf5611eeba622c7cdc9f63eb9ddf7f3f595c
}

func TestCreateTransferTx(t *testing.T) {
	curve := starknet.SC()
	contractAddr := starknet.ETH
	from := "0x0109d6f1e3821348320b8021f58c192dafd25fd67b66b29bde5096a1ab94f26d"
	to := "0x02ba43f799e51a799d441acfa8946ef40db9143dd18f3498dab7789bae1858c7"
	amount := big.NewInt(525200000000000)
	maxFee := big.NewInt(300264264106568)
	nonce := big.NewInt(1)

	tx, err := CreateTransferTx(curve, contractAddr, from, to, amount, nonce, maxFee, starknet.MAINNET_ID)
	assert.NoError(t, err)
	err = starknet.SignTx(curve, tx, "0x01a820094d8a382db8f6b78a84eb70d909f4344b98a497c52e78de19e049f2da")
	assert.NoError(t, err)

	request := tx.GetTxRequest()
	b, _ := json.Marshal(request)
	assert.Equal(t, "0x79c13840c054ced9461ce2a194adaef3f697978a0c4ca013c393aabb0069d78", starknet.BigToHex(tx.TransactionHash))
	assert.Equal(t, `{"type":"INVOKE_FUNCTION","sender_address":"0x0109d6f1e3821348320b8021f58c192dafd25fd67b66b29bde5096a1ab94f26d","calldata":["1","2087021424722619777119509474943472645767659996348769578120564519014510906823","232670485425082704932579856502088130646006032362877466777181098476241604910","3","1233728343534402321001430329096373382066876795249705380731790971568720992455","525200000000000","0"],"max_fee":"0x11116b8cd0248","signature":["1611692104187385326457642558786597593472514394867788527145169488578467933305","3583199820395466713929179717515709636681528038843985085799961404412095062038"],"version":"0x1","nonce":"0x1"}`, string(b))
	// https://starkscan.co/tx/0x79c13840c054ced9461ce2a194adaef3f697978a0c4ca013c393aabb0069d78
}

func TestCreateCrossChain(t *testing.T) {
	curve := starknet.SC()
	contractAddr := starknet.ETHBridge
	from := "0x02953ed33d04e24dabe74473ea895fe61eb2e8472249881eb98ab2d52705b69b"
	functionName := "initiate_withdraw"

	to := "0x62e206b4dDd402056D881DED58c0bd87193d2913"
	amount := big.NewInt(100)
	calldata := []string{to, starknet.BigToHex(amount)}

	maxFee := big.NewInt(1864315586779310)
	nonce := big.NewInt(2)

	tx, err := CreateContractTx(curve, contractAddr, from, functionName, calldata, nonce, maxFee, starknet.GOERLI_ID)
	assert.NoError(t, err)
	err = starknet.SignTx(curve, tx, "0x01a820094d8a382db8f6b78a84eb70d909f4344b98a497c52e78de19e049f2da")
	assert.NoError(t, err)

	request := tx.GetTxRequest()
	b, _ := json.Marshal(request)
	assert.Equal(t, "0x5ad70fab55be1508d27cd05dc8d1e8db77a66d5b5a49f328628cb5211332bbf", starknet.BigToHex(tx.TransactionHash))
	assert.Equal(t, `{"type":"INVOKE_FUNCTION","sender_address":"0x02953ed33d04e24dabe74473ea895fe61eb2e8472249881eb98ab2d52705b69b","calldata":["1","3256441166037631918262930812410838598500200462657642943867372734773841898370","403823062618199777388530751713272716715733872218085068081490028803159187238","3","564521648175006025532572708057195208089056127251","100","0"],"max_fee":"0x69f95cc4c98ae","signature":["3023601260780906156876842196562896672203512794084297088705888080025487805957","2446954015321026967907118403498701441994409929294212502462091775444501406894"],"version":"0x1","nonce":"0x2"}`, string(b))
	// https://testnet.starkscan.co/tx/0x05ad70fab55be1508d27cd05dc8d1e8db77a66d5b5a49f328628cb5211332bbf
}

func TestCreateMutiContracTx(t *testing.T) {
	curve := starknet.SC()
	contractAddr := starknet.ETH
	from := "0x02953ed33d04e24dabe74473ea895fe61eb2e8472249881eb98ab2d52705b69b"

	to1 := "0x02953ed33d04e24dabe74473ea895fe61eb2e8472249881eb98ab2d52705b69b"
	to2 := "0x026e9E8c411056B64B2D044EBCb39FC810D652Cfbe694326651d796BB078320b"

	maxFee := big.NewInt(509872256144985)
	nonce := big.NewInt(4)

	txs := []starknet.Calls{
		{
			ContractAddress: contractAddr,
			Entrypoint:      "transfer",
			Calldata:        []string{to1, "100", "0"},
		},
		{
			ContractAddress: contractAddr,
			Entrypoint:      "transfer",
			Calldata:        []string{to2, "100", "0"},
		},
	}

	tx, err := CreateSignedMultiContractTx(curve, from, txs, nonce, maxFee, starknet.MAINNET_ID, "0x01a820094d8a382db8f6b78a84eb70d909f4344b98a497c52e78de19e049f2da")
	assert.NoError(t, err)
	request := tx.GetTxRequest()
	b, _ := json.Marshal(request)
	assert.Equal(t, "0x6a2a2493a78d60eff44942fb4643c8fe381ba23bd1668450a9a959c568c36ea", starknet.BigToHex(tx.TransactionHash))
	assert.Equal(t, `{"type":"INVOKE_FUNCTION","sender_address":"0x02953ed33d04e24dabe74473ea895fe61eb2e8472249881eb98ab2d52705b69b","calldata":["2","2087021424722619777119509474943472645767659996348769578120564519014510906823","232670485425082704932579856502088130646006032362877466777181098476241604910","3","1168319513066818777483581844199041165067911991331081399636770980953504200347","100","0","2087021424722619777119509474943472645767659996348769578120564519014510906823","232670485425082704932579856502088130646006032362877466777181098476241604910","3","1100073131459501680801927467743186870973801404098697873181544877894944698891","100","0"],"max_fee":"0x1cfb9e2b55659","signature":["2169048640450690796796724595548366497645155157111145813203434715952997830872","3565526321264969572091926062882031431251004881005167625006935384245793389602"],"version":"0x1","nonce":"0x4"}`, string(b))
	// https://starkscan.co/tx/0x6a2a2493a78d60eff44942fb4643c8fe381ba23bd1668450a9a959c568c36ea
}

func TestName(t *testing.T) {
	curve := starknet.SC()

	address := "0x0109d6f1e3821348320b8021f58c192dafd25fd67b66b29bde5096a1ab94f26d"
	msg := "0x45514f85d4e7e2d3db3aac059a5d937f6c5d0f61f87ba25fa138c038248ce7a"

	sig, err := SignMessageV1(curve, address, msg, starknet.MAINNET_ID, "0x01a820094d8a382db8f6b78a84eb70d909f4344b98a497c52e78de19e049f2da")
	assert.NoError(t, err)
	assert.Equal(t, "{\"publicKey\":\"0x072ff9867ba607f204042c328cde87ddefe405b830e6515563fbe3ced9342109\",\"publicKeyY\":\"0x035217f0237715f7c90d0e94cd57203da9164e618bf86e7ad3e17d4238b7da55\",\"signedDataR\":\"0x0023694531f7be66487610ca29abd37e2687340dd0081298016e6c427e737419\",\"signedDataS\":\"0x07ea239c69e6bf871807d31d7015759fa685a7907b9e5ee5d47ebaa8012f2980\"}", sig)
}
