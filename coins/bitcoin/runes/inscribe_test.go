package bitcoin

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestBuildOpReturnData_forRunesMain(t *testing.T) {
	data := `{"edicts":[{"block":"837557","id":"1234","amount":"100000000000000000100000000000000000","output":0}],"isDefaultOutput":true,"defaultOutput":1,"mint":true,"mintNum":1}`
	res, err := BuildOpReturnDataJson([]byte(data))
	require.NoError(t, err)
	assert.Equal(t, "6a5d0914b58f3314d2091601", hex.EncodeToString(res))
}

func TestCalculateMintTxSize(t *testing.T) {
	{
		data := `{"edicts":[{"block":"837557","id":"1234","amount":"100000000000000000100000000000000000","output":0}],"isDefaultOutput":true,"defaultOutput":1}`
		res, err := BuildOpReturnDataJson([]byte(data))
		assert.NoError(t, err)
		assert.Equal(t, "6a5d1a160100b58f33d2098080a8ec85acb5f5ac84b6baacae98a11300", hex.EncodeToString(res))
		network := &chaincfg.MainNetParams
		vsize, cost, err := CalculateMintTxFee("1PY7wFLq74G1yDkuvM9g5isWbwrZ1DiM2Y", "1PY7wFLq74G1yDkuvM9g5isWbwrZ1DiM2Y", hex.EncodeToString(res), "KwsDCBpcTQsc35Ev26LCw4QJCA73QJM5Lg8ViM7RBNdwsnEjFzri", network, 1, 1000)
		assert.NoError(t, err)
		assert.Equal(t, int64(230), vsize)
		assert.Equal(t, int64(1230), cost)
	}
	{
		data := `{"edicts":[{"block":"837557","id":"1234","amount":"100000000000000000100000000000000000","output":0}],"isDefaultOutput":true,"defaultOutput":1,"mint":true,"mintNum":1}`
		res, err := BuildOpReturnDataJson([]byte(data))
		assert.NoError(t, err)
		assert.Equal(t, "6a5d0914b58f3314d2091601", hex.EncodeToString(res))
		network := &chaincfg.MainNetParams
		vsize, cost, err := CalculateMintTxFee("1PY7wFLq74G1yDkuvM9g5isWbwrZ1DiM2Y", "1PY7wFLq74G1yDkuvM9g5isWbwrZ1DiM2Y", hex.EncodeToString(res), "KwsDCBpcTQsc35Ev26LCw4QJCA73QJM5Lg8ViM7RBNdwsnEjFzri", network, 1, 1000)
		assert.NoError(t, err)
		assert.Equal(t, int64(213), vsize)
		assert.Equal(t, int64(1213), cost)
	}
}

func TestBuildOpReturnScript(t *testing.T) {
	edict := &Edict{
		Id:     "1234",
		Block:  "837557",
		Amount: "21000",
		Output: 0,
	}
	edicts := []*Edict{edict}
	res, err := BuildOpReturnData(edicts, true, false, 0)
	assert.NoError(t, err)
	assert.Equal(t, "6a5d0c160000b58f33d20988a40100", hex.EncodeToString(res))
}

func TestBuildOpReturnScript_forRunesMain(t *testing.T) {
	edict := &Edict{
		Id:     "1234",
		Block:  "837557",
		Amount: "21000",
		Output: 0,
	}
	edicts := []*Edict{edict}
	//res, err := BuildOpReturnData(edicts, true, false, 0)
	res, err := BuildOpReturnScriptForRuneMainEdict(edicts, true, true, 1)
	assert.NoError(t, err)
	assert.Equal(t, "6a0952554e455f54455354090c0100a33480a30800", hex.EncodeToString(res))
}

func TestEncodeToVecV22(t *testing.T) {
	buf := bytes.Buffer{}
	EncodeToVecV2(big.NewInt(837557), &buf)
	bytes := []byte{181, 143, 51}
	assert.Equal(t, bytes, buf.Bytes())
}

func TestEncodeToVec(t *testing.T) {
	e := &Edict{
		Id:     "2aa16001b",
		Amount: "21000",
		Output: 0,
	}
	//res := EncodeToVec(new(big.Int).SetInt64(0))
	idB, ok := new(big.Int).SetString(e.Id, 16)
	assert.True(t, ok)
	tagBody := new(big.Int).SetInt64(0)
	payload := []int64{}
	payload = append(payload, EncodeToVec(tagBody)...)
	payload = append(payload, EncodeToVec(idB)...)
	amountB, _ := new(big.Int).SetString(e.Amount, 10)
	payload = append(payload, EncodeToVec(amountB)...)
	output := new(big.Int).SetUint64(uint64(e.Output))
	payload = append(payload, EncodeToVec(output)...)
	expected := []int64{0, 169, 207, 214, 255, 27, 128, 163, 8, 0}
	assert.Equal(t, expected, payload)
}
