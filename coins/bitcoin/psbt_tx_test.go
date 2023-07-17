package bitcoin

import (
	"github.com/btcsuite/btcd/chaincfg"
	"testing"
)

func TestPsbt(t *testing.T) {
	network := &chaincfg.TestNet3Params
	// seller
	txInput := &TxInput{
		TxId:       "46e3ce050474e6da80760a2a0b062836ff13e2a42962dc1c9b17b8f962444206",
		VOut:       uint32(0),
		Amount:     int64(546),
		Address:    "tb1pklh8lqax5l7m2ycypptv2emc4gata2dy28svnwcp9u32wlkenvsspcvhsr",
		PrivateKey: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
	}

	txOutput := &TxOutput{
		Address: "2NF33rckfiQTiE5Guk5ufUdwms8PgmtnEdc",
		Amount:  int64(100000),
	}

	sellerPsbt, err := GenerateSignedListingPSBTBase64(txInput, txOutput, network)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sellerPsbt)

	// buyer
	var inputs []*TxInput
	inputs = append(inputs, &TxInput{
		TxId:       "25b9d08a26c8d47795301dd47a861cff0459d14f27fbd41cffaca17d9aa20f87",
		VOut:       uint32(0),
		Amount:     int64(249352),
		Address:    "tb1qtsq9c4fje6qsmheql8gajwtrrdrs38kdzeersc",
		PrivateKey: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
	})
	inputs = append(inputs, &TxInput{
		TxId:           "6d59aa50447c0d55e6f9535c3e56d7014b4ca8070ee57ce2199219790cfd5815",
		VOut:           uint32(0),
		Amount:         int64(499356),
		Address:        "mouQtmBWDS7JnT65Grj2tPzdSmGKJgRMhE",
		PrivateKey:     "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
		NonWitnessUtxo: "02000000010a6b13715c8effde51dac60d572358005a589cd80413a88e0912e4c6d275abbe010000006a473044022019e34aa16cf55eb9c7a8627f61bcd671525a3818a23ab8a78af13c35121ea3c8022055a5bfb3e8486f6e83707660f1fca3da06f140f449902a63900625f43fadf10501210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2fffffffff019c9e0700000000001976a9145c005c5532ce810ddf20f9d1d939631b47089ecd88ac00000000",
	})
	// seller input
	inputs = append(inputs, txInput)
	inputs = append(inputs, &TxInput{
		TxId:       "d1696c10046ec8b2d938924f1923f1f2e1588095fbf3ea0f8cd640b51da51ba2",
		VOut:       uint32(0),
		Amount:     int64(400),
		Address:    "2NF33rckfiQTiE5Guk5ufUdwms8PgmtnEdc",
		PrivateKey: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
	})

	var outputs []*TxOutput
	outputs = append(outputs, &TxOutput{
		Address: "tb1qtsq9c4fje6qsmheql8gajwtrrdrs38kdzeersc",
		Amount:  int64(200000),
	})
	outputs = append(outputs, &TxOutput{
		Address: "mouQtmBWDS7JnT65Grj2tPzdSmGKJgRMhE",
		Amount:  int64(200000),
	})
	// seller output
	outputs = append(outputs, txOutput)
	outputs = append(outputs, &TxOutput{
		Address: "tb1pklh8lqax5l7m2ycypptv2emc4gata2dy28svnwcp9u32wlkenvsspcvhsr",
		Amount:  int64(246500),
	})
	outputs = append(outputs, &TxOutput{
		Address: "tb1pklh8lqax5l7m2ycypptv2emc4gata2dy28svnwcp9u32wlkenvsspcvhsr",
		Amount:  int64(1000),
	})
	outputs = append(outputs, &TxOutput{
		Address: "tb1pklh8lqax5l7m2ycypptv2emc4gata2dy28svnwcp9u32wlkenvsspcvhsr",
		Amount:  int64(1000),
	})

	fee, err := CalcFee(inputs, outputs, sellerPsbt, 2, network)
	if err != nil {
		return
	}
	t.Log(fee)

	buyerTx, err := GenerateSignedBuyingTx(inputs, outputs, sellerPsbt, network)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(buyerTx)
}
