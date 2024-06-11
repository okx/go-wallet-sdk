package bitcoin

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
)

func TestSrc20Inscribe(t *testing.T) {
	network := &chaincfg.TestNet3Params

	commitTxPrevOutputList := make([]*PrevOutput, 0)
	commitTxPrevOutputList = append(commitTxPrevOutputList, &PrevOutput{
		TxId:       "c865cd4dc206ccdaf1cff0fad4f0272f2075af5c975c670debbf8d56045391ad",
		VOut:       3,
		Amount:     202000,
		Address:    "2NF33rckfiQTiE5Guk5ufUdwms8PgmtnEdc",
		PrivateKey: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
	})

	inscriptionData := &InscriptionData{
		ContentType: "stamp:",
		Body:        []byte(`{"p":"src-20","op":"deploy","tick":"coder","max":"21000000","lim":"1000","dec":"8"}`),
		RevealAddr:  "2NF33rckfiQTiE5Guk5ufUdwms8PgmtnEdc",
	}

	request := &Src20InscriptionRequest{
		CommitTxPrevOutputList: commitTxPrevOutputList,
		CommitFeeRate:          100,
		RevealOutValue:         790,
		InscriptionData:        inscriptionData,
		Address:                "2NF33rckfiQTiE5Guk5ufUdwms8PgmtnEdc",
	}

	txs, _ := Src20Inscribe(network, request)

	expected := "{\"commitTx\":\"02000000000101ad915304568dbfeb0d675c975caf75202f27f0d4faf0cff1dacc06c24dcd65c803000000171600145c005c5532ce810ddf20f9d1d939631b47089ecdfdffffff04160300000000000017a914ef05515a0595d15eaf90d9f62fb85873a6d8c0b4871603000000000000695121034a54cfbca897d6e5bd94c8b03e0524e9849b8d5f19ac6eb79ec78ea402271d002102651491c55c5a27dc6838d312ca9e9350ae2cbdc02f4903bf0fcbf87ffc9096002102020202020202020202020202020202020202020202020202020202020202020253ae160300000000000069512103a964c52310e9976582c01d9705c7308949173d7e571df1e244ceb348b54a850021024a4637e826e37fb67470f97bcd954a0b5a4e20ef37d16f5b5d64cbc58081b8002102020202020202020202020202020202020202020202020202020202020202020253ae5ad101000000000017a914ef05515a0595d15eaf90d9f62fb85873a6d8c0b48702473044022021bb8f5c0bf9e6e12f31e7e9ae99d777a541fed6310c86cf63af5e5ce24e54f502204c99f6f0d5793ee3ccd206874d6430de3e10076e95631d4f94c4a40f666b76a901210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f00000000\",\"revealTxs\":[],\"commitTxFee\":80500,\"revealTxFees\":[],\"commitAddrs\":[]}"
	txsBytes, _ := json.Marshal(txs)
	assert.Equal(t, len(txs.RevealTxs), 0)
	assert.Equal(t, txs.CommitTxFee, int64(80500))
	assert.Equal(t, expected, string(txsBytes))
}
