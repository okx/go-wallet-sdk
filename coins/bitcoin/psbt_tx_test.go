package bitcoin

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require.Nil(t, err)
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
	require.Nil(t, err)
	t.Log(fee)

	buyerTx, err := GenerateSignedBuyingTx(inputs, outputs, sellerPsbt, network)
	require.Nil(t, err)
	t.Log(buyerTx)
}

func TestGenerateUnsignedPSBTHex(t *testing.T) {
	network := &chaincfg.TestNet3Params
	var inputs []*TxInput
	inputs = append(inputs, &TxInput{
		TxId:              "46e3ce050474e6da80760a2a0b062836ff13e2a42962dc1c9b17b8f962444206",
		VOut:              uint32(0),
		Sequence:          1,
		Amount:            int64(546),
		Address:           "2NF33rckfiQTiE5Guk5ufUdwms8PgmtnEdc",
		PrivateKey:        "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
		MasterFingerprint: 0xF23F9FD2,
		DerivationPath:    "m/44'/0'/0'/0/0",
		PublicKey:         "0357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f",
	})

	var outputs []*TxOutput
	outputs = append(outputs, &TxOutput{
		Address: "2NF33rckfiQTiE5Guk5ufUdwms8PgmtnEdc",
		Amount:  int64(100000),
	})
	psbtHex, err := GenerateUnsignedPSBTHex(inputs, outputs, network)
	require.Nil(t, err)
	t.Log(psbtHex)
}

func TestExtractTxFromSignedPSBT(t *testing.T) {
	psbtHex := "70736274ff0100fd90010200000004870fa29a7da1acff1cd4fb274fd15904ff1c867ad41d309577d4c8268ad0b9250000000000ffffffff1558fd0c79199219e27ce50e07a84c4b01d7563e5c53f9e6550d7c4450aa596d0000000000ffffffff06424462f9b8179b1cdc6229a4e213ff3628060b2a0a7680dae6740405cee3460000000000ffffffffa21ba51db540d68c0feaf3fb958058e1f2f123194f9238d9b2c86e04106c69d10000000000ffffffff06400d0300000000001600145c005c5532ce810ddf20f9d1d939631b47089ecd400d0300000000001976a9145c005c5532ce810ddf20f9d1d939631b47089ecd88aca08601000000000017a914ef05515a0595d15eaf90d9f62fb85873a6d8c0b487e4c2030000000000225120b7ee7f83a6a7fdb513040856c56778aa3abea9a451e0c9bb012f22a77ed99b21e803000000000000225120b7ee7f83a6a7fdb513040856c56778aa3abea9a451e0c9bb012f22a77ed99b21e803000000000000225120b7ee7f83a6a7fdb513040856c56778aa3abea9a451e0c9bb012f22a77ed99b21000000000001011f08ce0300000000001600145c005c5532ce810ddf20f9d1d939631b47089ecd01086c02483045022100a1d12dee8d87d2f8a12ff43f656a6b52183fa5ce4ffd1ab349b978d4dc5e68620220060d8c6d20ea34d3b2f744624d9f027c9020cb80cfb9babe015ebd70db0a927a01210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f000100bf02000000010a6b13715c8effde51dac60d572358005a589cd80413a88e0912e4c6d275abbe010000006a473044022019e34aa16cf55eb9c7a8627f61bcd671525a3818a23ab8a78af13c35121ea3c8022055a5bfb3e8486f6e83707660f1fca3da06f140f449902a63900625f43fadf10501210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2fffffffff019c9e0700000000001976a9145c005c5532ce810ddf20f9d1d939631b47089ecd88ac0000000001076b483045022100bd9b8c17d68efed18f0882bdb77db303a0a547864305e32ed7a9a951b650caa90220131c361e5c27652a3a05603306a87d8f6e117b78fdb1082db23d8960eb6214bf01210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f0001012b2202000000000000225120b7ee7f83a6a7fdb513040856c56778aa3abea9a451e0c9bb012f22a77ed99b210108430141f24c018bc95e051c33e4659cacad365db8f3afbaf61ee163e3e1bf1d419baaeb681f681c75a545a19d4ade0b972e226448015d9cbdaee121f4148b5bee9d27068300010120900100000000000017a914ef05515a0595d15eaf90d9f62fb85873a6d8c0b4870107171600145c005c5532ce810ddf20f9d1d939631b47089ecd01086c02483045022100bb251cc4a4db4eab3352d54541a03d20d5067e8261b6f7ba8a20a7d955dfafde022078be1dd187ff61934177a9245872f4a90beef32ec40b69f75d9c50c32053d97101210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f00000000000000"
	txHex, err := ExtractTxFromSignedPSBT(psbtHex)
	require.Nil(t, err)
	assert.Equal(t, "02000000000104870fa29a7da1acff1cd4fb274fd15904ff1c867ad41d309577d4c8268ad0b9250000000000ffffffff1558fd0c79199219e27ce50e07a84c4b01d7563e5c53f9e6550d7c4450aa596d000000006b483045022100bd9b8c17d68efed18f0882bdb77db303a0a547864305e32ed7a9a951b650caa90220131c361e5c27652a3a05603306a87d8f6e117b78fdb1082db23d8960eb6214bf01210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2fffffffff06424462f9b8179b1cdc6229a4e213ff3628060b2a0a7680dae6740405cee3460000000000ffffffffa21ba51db540d68c0feaf3fb958058e1f2f123194f9238d9b2c86e04106c69d100000000171600145c005c5532ce810ddf20f9d1d939631b47089ecdffffffff06400d0300000000001600145c005c5532ce810ddf20f9d1d939631b47089ecd400d0300000000001976a9145c005c5532ce810ddf20f9d1d939631b47089ecd88aca08601000000000017a914ef05515a0595d15eaf90d9f62fb85873a6d8c0b487e4c2030000000000225120b7ee7f83a6a7fdb513040856c56778aa3abea9a451e0c9bb012f22a77ed99b21e803000000000000225120b7ee7f83a6a7fdb513040856c56778aa3abea9a451e0c9bb012f22a77ed99b21e803000000000000225120b7ee7f83a6a7fdb513040856c56778aa3abea9a451e0c9bb012f22a77ed99b2102483045022100a1d12dee8d87d2f8a12ff43f656a6b52183fa5ce4ffd1ab349b978d4dc5e68620220060d8c6d20ea34d3b2f744624d9f027c9020cb80cfb9babe015ebd70db0a927a01210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f000141f24c018bc95e051c33e4659cacad365db8f3afbaf61ee163e3e1bf1d419baaeb681f681c75a545a19d4ade0b972e226448015d9cbdaee121f4148b5bee9d27068302483045022100bb251cc4a4db4eab3352d54541a03d20d5067e8261b6f7ba8a20a7d955dfafde022078be1dd187ff61934177a9245872f4a90beef32ec40b69f75d9c50c32053d97101210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f00000000", txHex)
}
