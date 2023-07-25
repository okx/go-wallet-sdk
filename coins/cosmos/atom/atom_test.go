package atom

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/okx/go-wallet-sdk/coins/cosmos"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestAtom(t *testing.T) {
	pk, err := hex.DecodeString("1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37")
	if err != nil {
		t.Fatal(err)
	}
	k, _ := btcec.PrivKeyFromBytes(pk)
	address, err := NewAddress(k)
	if err != nil {
		t.Fatal(err)
	}
	if address != "cosmos145q0tcdur4tcx2ya5cphqx96e54yflfyqjrdt5" {
		t.Fatal("NewAddress failed")
	}

	ret := ValidateAddress(address)
	if !ret {
		t.Fatal("ValidateAddress failed")
	}

	chainId := "cosmoshub-4"
	from := "cosmos145q0tcdur4tcx2ya5cphqx96e54yflfyqjrdt5"
	to := "cosmos1jun53r4ycc8g2v6tffp4cmxjjhv6y7ntat62wn"
	demon := "uatom"
	memo := "memo"
	amount := big.NewInt(10000)
	sequence := 0
	accountNumber := 623151
	feeAmount := big.NewInt(10)
	gasLimit := 100
	hexStr, err := SignStart(chainId, from, to, demon, memo, amount, 0, uint64(sequence), uint64(accountNumber), feeAmount, uint64(gasLimit), k)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hexStr)
	if hexStr != "0a97010a8e010a1c2f636f736d6f732e62616e6b2e763162657461312e4d736753656e64126e0a2d636f736d6f733134357130746364757234746378327961356370687178393665353479666c6679716a72647435122d636f736d6f73316a756e3533723479636338673276367466667034636d786a6a68763679376e7461743632776e1a0e0a057561746f6d1205313030303012046d656d6f12610a4e0a460a1f2f636f736d6f732e63727970746f2e736563703235366b312e5075624b657912230a21031053e9ef0295d334b6bb22e20cc717eb1a16a546f692572c8830b4bc14c1367612040a020801120f0a0b0a057561746f6d1202313010641a0b636f736d6f736875622d3420af8426" {
		t.Fatal("make transaction failed, hexStr: ", hexStr)
	}

	signStr, err := Sign(hexStr, k)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(signStr)
	if signStr != "57fa782b2982e9119d285a2cbbb2e4a8d8c08a1c7b419bad2cf4a6d219046f2c6d9d42424417b23d551e35cc60f461f5caaf5739334b59a9e9458719e8920296" {
		t.Fatal("sgin transaction failed, signStr: ", signStr)
	}

	trans, err := SignEnd(hexStr, signStr)
	if err != nil {
		t.Fatal(err)
	}
	if trans != "CpcBCo4BChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEm4KLWNvc21vczE0NXEwdGNkdXI0dGN4MnlhNWNwaHF4OTZlNTR5ZmxmeXFqcmR0NRItY29zbW9zMWp1bjUzcjR5Y2M4ZzJ2NnRmZnA0Y214ampodjZ5N250YXQ2MnduGg4KBXVhdG9tEgUxMDAwMBIEbWVtbxJhCk4KRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDEFPp7wKV0zS2uyLiDMcX6xoWpUb2klcsiDC0vBTBNnYSBAoCCAESDwoLCgV1YXRvbRICMTAQZBpAV/p4KymC6RGdKFosu7LkqNjAihx7QZutLPSm0hkEbyxtnUJCRBeyPVUeNcxg9GH1yq9XOTNLWanpRYcZ6JIClg==" {
		t.Fatal("SignEnd failed, trans: ", trans)
	}
}

func TestAtom2(t *testing.T) {
	pk, err := hex.DecodeString("1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37")
	if err != nil {
		t.Fatal(err)
	}
	k, _ := btcec.PrivKeyFromBytes(pk)

	param := cosmos.TransferParam{}
	param.FromAddress = "cosmos145q0tcdur4tcx2ya5cphqx96e54yflfyqjrdt5"
	param.ToAddress = "cosmos1jun53r4ycc8g2v6tffp4cmxjjhv6y7ntat62wn"
	param.Demon = "uatom"
	param.Amount = "10000"
	param.CommonParam.ChainId = "cosmoshub-4"
	param.CommonParam.Sequence = 0
	param.CommonParam.AccountNumber = 623151
	param.CommonParam.FeeDemon = "uatom"
	param.CommonParam.FeeAmount = "10"
	param.CommonParam.GasLimit = 100
	param.CommonParam.Memo = "memo"
	param.CommonParam.TimeoutHeight = 0

	doc, err := cosmos.GetRawTransaction(param, hex.EncodeToString(k.PubKey().SerializeCompressed()))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(doc)

	signature, err := cosmos.SignRawTransaction(doc, k)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(signature)

	result, err := cosmos.GetSignedTransaction(doc, signature)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)

	if result != "CpcBCo4BChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEm4KLWNvc21vczE0NXEwdGNkdXI0dGN4MnlhNWNwaHF4OTZlNTR5ZmxmeXFqcmR0NRItY29zbW9zMWp1bjUzcjR5Y2M4ZzJ2NnRmZnA0Y214ampodjZ5N250YXQ2MnduGg4KBXVhdG9tEgUxMDAwMBIEbWVtbxJhCk4KRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDEFPp7wKV0zS2uyLiDMcX6xoWpUb2klcsiDC0vBTBNnYSBAoCCAESDwoLCgV1YXRvbRICMTAQZBpAV/p4KymC6RGdKFosu7LkqNjAihx7QZutLPSm0hkEbyxtnUJCRBeyPVUeNcxg9GH1yq9XOTNLWanpRYcZ6JIClg==" {
		t.Fatal("Sign failed")
	}
}

// Check account details
// https://api.cosmos.network/cosmos/auth/v1beta1/accounts/cosmos1rvs5xph4l3px2efynqsthus8p6r4exyr7ckyxv
// curl -X POST -d '{"tx_bytes":"CpABCo0BChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEm0KLWNvc21vczFydnM1eHBoNGwzcHgyZWZ5bnFzdGh1czhwNnI0ZXh5cjdja3l4dhItY29zbW9zMXJ2czV4cGg0bDNweDJlZnlucXN0aHVzOHA2cjRleHlyN2NreXh2Gg0KBXVhdG9tEgQxMDAwEmYKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQOcJMA96W11QpNEacdGblBLXYYIw5nd27SBSxlh+Pc6UxIECgIIfxgCEhIKDAoFdWF0b20SAzEzMBCgjQYaQA04G6nhx6Zo8uYBHKhyw46t7RkvxLwDO0XfkRG3hVRRDmCg6xl+61FhXe3x2A/temH/hGsIt1bjs37vcDQAgg4=","mode":"BROADCAST_MODE_SYNC"}' https://api.cosmos.network/cosmos/tx/v1beta1/txs
func TestAtom3(t *testing.T) {
	pk, err := hex.DecodeString("1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37")
	if err != nil {
		t.Fatal(err)
	}
	k, _ := btcec.PrivKeyFromBytes(pk)
	address, err := NewAddress(k)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(address)

	param := cosmos.TransferParam{}
	param.FromAddress = "cosmos1rvs5xph4l3px2efynqsthus8p6r4exyr7ckyxv"
	param.ToAddress = "cosmos1rvs5xph4l3px2efynqsthus8p6r4exyr7ckyxv"
	param.Demon = "uatom"
	param.Amount = "1000"
	param.CommonParam.ChainId = "cosmoshub-4"
	param.CommonParam.Sequence = 2
	param.CommonParam.AccountNumber = 1225716
	param.CommonParam.FeeDemon = "uatom"
	param.CommonParam.FeeAmount = "130"
	param.CommonParam.GasLimit = 100000
	param.CommonParam.Memo = ""

	doc, err := cosmos.GetRawJsonTransaction(param)
	if err != nil {
		t.Fatal(err)
	}

	signature, err := cosmos.SignRawJsonTransaction(doc, k)
	if err != nil {
		t.Fatal(err)
	}

	publicKey := hex.EncodeToString(k.PubKey().SerializeCompressed())
	result, err := cosmos.GetSignedJsonTransaction(doc, publicKey, signature)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
	assert.Equal(t, result, "CpABCo0BChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEm0KLWNvc21vczFydnM1eHBoNGwzcHgyZWZ5bnFzdGh1czhwNnI0ZXh5cjdja3l4dhItY29zbW9zMXJ2czV4cGg0bDNweDJlZnlucXN0aHVzOHA2cjRleHlyN2NreXh2Gg0KBXVhdG9tEgQxMDAwEmYKUApGCh8vY29zbW9zLmNyeXB0by5zZWNwMjU2azEuUHViS2V5EiMKIQMQU+nvApXTNLa7IuIMxxfrGhalRvaSVyyIMLS8FME2dhIECgIIfxgCEhIKDAoFdWF0b20SAzEzMBCgjQYaQNbZuALq96PyACNflFUnnS5kd/fCZniLitsmRhCR092UHIJYGEUugnG89Be+v6BbWe0E0jPRqPfw36thibn6ix0=")
}
