package bitcoin

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/assert"
)

func TestSignTx_Step(t *testing.T) {
	// legacy address
	txBuild := NewTxBuild(1, &chaincfg.TestNet3Params)
	txBuild.AddInput2("c44a7f98434e5e875a573339f77d36022c79c525771fa88c72fa53f3a55eeaf7", 1, "", "mouQtmBWDS7JnT65Grj2tPzdSmGKJgRMhE", 1488430)
	txBuild.AddOutput("mouQtmBWDS7JnT65Grj2tPzdSmGKJgRMhE", 1488200)
	tx, err := txBuild.Build2()
	assert.Nil(t, err)
	err = SignBuildTx(tx, []Input{
		{
			txId:    "c44a7f98434e5e875a573339f77d36022c79c525771fa88c72fa53f3a55eeaf7",
			vOut:    1,
			address: "mouQtmBWDS7JnT65Grj2tPzdSmGKJgRMhE",
			amount:  1488430,
		},
	}, map[int]string{
		0: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
	}, &chaincfg.TestNet3Params)
	assert.Nil(t, err)
	txHex, err := GetTxHex(tx)
	assert.Nil(t, err)
	assert.Equal(t, "0100000001f7ea5ea5f353fa728ca81f7725c5792c02367df73933575a875e4e43987f4ac40100000069463043021f58d5662b5215849834d0e402dd2e27f6f5ff06f132eabd2bfcf1eb0ac0cf6602202684f8d65100c6316ed9d1ea0ca3e7119fc4868677928485b0924ea381f6686501210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2fffffffff0148b51600000000001976a9145c005c5532ce810ddf20f9d1d939631b47089ecd88ac00000000", txHex)

	// segwit_nested address
	txBuild = NewTxBuild(1, &chaincfg.TestNet3Params)
	txBuild.AddInput2("ce69ca86b68708afc8484dacb7730006e7eff6d0c18b18a16a9e91abeefeb08a", 0, "", "2NF33rckfiQTiE5Guk5ufUdwms8PgmtnEdc", 2000)
	txBuild.AddOutput("2NF33rckfiQTiE5Guk5ufUdwms8PgmtnEdc", 900)
	txBuild.AddOutput("2NF33rckfiQTiE5Guk5ufUdwms8PgmtnEdc", 850)
	tx, err = txBuild.Build2()
	assert.Nil(t, err)
	err = SignBuildTx(tx, []Input{
		{
			txId:    "ce69ca86b68708afc8484dacb7730006e7eff6d0c18b18a16a9e91abeefeb08a",
			vOut:    0,
			address: "2NF33rckfiQTiE5Guk5ufUdwms8PgmtnEdc",
			amount:  2000,
		},
	}, map[int]string{
		0: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
	}, &chaincfg.TestNet3Params)
	assert.Nil(t, err)
	txHex, err = GetTxHex(tx)
	assert.Nil(t, err)
	assert.Equal(t, "010000000001018ab0feeeab919e6aa1188bc1d0f6efe7060073b7ac4d48c8af0887b686ca69ce00000000171600145c005c5532ce810ddf20f9d1d939631b47089ecdffffffff02840300000000000017a914ef05515a0595d15eaf90d9f62fb85873a6d8c0b487520300000000000017a914ef05515a0595d15eaf90d9f62fb85873a6d8c0b48702483045022100f8bc4f7e5f0a29a3f5b8a75f60a3eac2291b1e5ae8300403d355a134aa99568d02203a3eecf6ca9ae8ce20bf6ba329b0ff27290fae5d389566e84b7c625b9681752201210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f00000000", txHex)

	// segwit_native address
	txBuild = NewTxBuild(1, &chaincfg.TestNet3Params)
	txBuild.AddInput2("0bc66f18fd95ca00b6569471aa2dcd47fe45d3446fbaeec9ced228b00713fe8c", 0, "", "tb1qtsq9c4fje6qsmheql8gajwtrrdrs38kdzeersc", 200000)
	txBuild.AddOutput("tb1qtsq9c4fje6qsmheql8gajwtrrdrs38kdzeersc", 199700)
	tx, err = txBuild.Build2()
	assert.Nil(t, err)
	err = SignBuildTx(tx, []Input{
		{
			txId:    "0bc66f18fd95ca00b6569471aa2dcd47fe45d3446fbaeec9ced228b00713fe8c",
			vOut:    0,
			address: "tb1qtsq9c4fje6qsmheql8gajwtrrdrs38kdzeersc",
			amount:  200000,
		},
	}, map[int]string{
		0: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
	}, &chaincfg.TestNet3Params)
	assert.Nil(t, err)
	txHex, err = GetTxHex(tx)
	assert.Equal(t, "010000000001018cfe1307b028d2cec9eeba6f44d345fe47cd2daa719456b600ca95fd186fc60b0000000000ffffffff01140c0300000000001600145c005c5532ce810ddf20f9d1d939631b47089ecd0247304402200fcc9be29e3ab99b81f30fdf0788d883576d1a313fc809fce616813b8b2db62002207d199936c665f8e2c353c03d968f49dfcc16dacaf5877487c7d49f532180a6a101210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f00000000", txHex)

	// taproot_keypath address
	txBuild = NewTxBuild(1, &chaincfg.TestNet3Params)
	txBuild.AddInput2("3cb62c77c5c3fc032100af4cae9eeb342829cbc5b49815f8db1bb8156314a784", 0, "", "tb1pklh8lqax5l7m2ycypptv2emc4gata2dy28svnwcp9u32wlkenvsspcvhsr", 546)
	txBuild.AddOutput("tb1pklh8lqax5l7m2ycypptv2emc4gata2dy28svnwcp9u32wlkenvsspcvhsr", 300)
	tx, err = txBuild.Build2()
	assert.Nil(t, err)
	err = SignBuildTx(tx, []Input{
		{
			txId:    "3cb62c77c5c3fc032100af4cae9eeb342829cbc5b49815f8db1bb8156314a784",
			vOut:    0,
			address: "tb1pklh8lqax5l7m2ycypptv2emc4gata2dy28svnwcp9u32wlkenvsspcvhsr",
			amount:  546,
		},
	}, map[int]string{
		0: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
	}, &chaincfg.TestNet3Params)
	assert.Nil(t, err)
	txHex, err = GetTxHex(tx)
	assert.Nil(t, err)
	assert.Equal(t, "0100000000010184a7146315b81bdbf81598b4c5cb292834eb9eae4caf002103fcc3c5772cb63c0000000000ffffffff012c01000000000000225120b7ee7f83a6a7fdb513040856c56778aa3abea9a451e0c9bb012f22a77ed99b2101408399ed30c87c4dbcb7de58b9beb5448a7b7d9c4a0d8048fca51571648007c0a81cf0024d6c72d3cea13d6d3d69f341bc77a30612f3dfe10ee64874b0feded3b300000000", txHex)

	// x
	txBuild = NewTxBuild(1, &chaincfg.TestNet3Params)
	txBuild.AddInput2("914f38907e27578505e12438448143e3eb6708ccac80ec68eb5d7901b6530119", 1, "cSpPPDWzJ36N4vrLq6QSR69yuKAPZL7YvJnZ2SpxitwtPEcT4PYm", "tb1qkd8hetlagq3ylg88zupqsu9w6rqykw5mvthhr4", 48870)
	txBuild.AddInput2("a5a91933bb75e1dab85868dcab3094ecffe8013486a43710b40215e9968b4a75", 0, "cSpPPDWzJ36N4vrLq6QSR69yuKAPZL7YvJnZ2SpxitwtPEcT4PYm", "tb1qkd8hetlagq3ylg88zupqsu9w6rqykw5mvthhr4", 8000)
	txBuild.AddOutput("mws4UFRP8XE8JhweXhgyMGkVPZfCMSFgmx", 56683)
	tx, err = txBuild.Build2()
	assert.Nil(t, err)
	err = SignBuildTx(tx, []Input{
		{
			txId:    "914f38907e27578505e12438448143e3eb6708ccac80ec68eb5d7901b6530119",
			vOut:    1,
			address: "tb1qkd8hetlagq3ylg88zupqsu9w6rqykw5mvthhr4",
			amount:  48870,
		},
		{
			txId:    "a5a91933bb75e1dab85868dcab3094ecffe8013486a43710b40215e9968b4a75",
			vOut:    0,
			address: "tb1qkd8hetlagq3ylg88zupqsu9w6rqykw5mvthhr4",
			amount:  8000,
		},
	}, map[int]string{
		0: "cSpPPDWzJ36N4vrLq6QSR69yuKAPZL7YvJnZ2SpxitwtPEcT4PYm",
		1: "cSpPPDWzJ36N4vrLq6QSR69yuKAPZL7YvJnZ2SpxitwtPEcT4PYm",
	}, &chaincfg.TestNet3Params)
	assert.Nil(t, err)
	txHex, err = GetTxHex(tx)
	assert.Nil(t, err)
	t.Log(txHex)

	tx2, err := txBuild.Build()
	assert.Nil(t, err)
	txHex2, err := GetTxHex(tx2)
	assert.Nil(t, err)
	t.Log(txHex2)
}
