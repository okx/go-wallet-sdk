package bitcoin

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/okx/go-wallet-sdk/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// support for single private key address formats (legacy/segwit_nested/segwit_native/taproot_keypath)
func TestSignTx(t *testing.T) {
	// legacy address
	txBuild := NewTxBuild(1, &chaincfg.TestNet3Params)
	txBuild.AddInput2("c44a7f98434e5e875a573339f77d36022c79c525771fa88c72fa53f3a55eeaf7", 1, "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22", "mouQtmBWDS7JnT65Grj2tPzdSmGKJgRMhE", 1488430)
	txBuild.AddOutput("mouQtmBWDS7JnT65Grj2tPzdSmGKJgRMhE", 1488200)
	tx, err := txBuild.Build()
	assert.Nil(t, err)
	txHex, err := GetTxHex(tx)
	assert.Nil(t, err)
	assert.Equal(t, "0100000001f7ea5ea5f353fa728ca81f7725c5792c02367df73933575a875e4e43987f4ac40100000069463043021f58d5662b5215849834d0e402dd2e27f6f5ff06f132eabd2bfcf1eb0ac0cf6602202684f8d65100c6316ed9d1ea0ca3e7119fc4868677928485b0924ea381f6686501210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2fffffffff0148b51600000000001976a9145c005c5532ce810ddf20f9d1d939631b47089ecd88ac00000000", txHex)

	// segwit_nested address
	txBuild = NewTxBuild(1, &chaincfg.TestNet3Params)
	txBuild.AddInput2("ce69ca86b68708afc8484dacb7730006e7eff6d0c18b18a16a9e91abeefeb08a", 0, "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22", "2NF33rckfiQTiE5Guk5ufUdwms8PgmtnEdc", 2000)
	txBuild.AddOutput("2NF33rckfiQTiE5Guk5ufUdwms8PgmtnEdc", 900)
	txBuild.AddOutput("2NF33rckfiQTiE5Guk5ufUdwms8PgmtnEdc", 850)
	tx, err = txBuild.Build()
	assert.Nil(t, err)
	txHex, err = GetTxHex(tx)
	assert.Nil(t, err)
	assert.Equal(t, "010000000001018ab0feeeab919e6aa1188bc1d0f6efe7060073b7ac4d48c8af0887b686ca69ce00000000171600145c005c5532ce810ddf20f9d1d939631b47089ecdffffffff02840300000000000017a914ef05515a0595d15eaf90d9f62fb85873a6d8c0b487520300000000000017a914ef05515a0595d15eaf90d9f62fb85873a6d8c0b48702483045022100f8bc4f7e5f0a29a3f5b8a75f60a3eac2291b1e5ae8300403d355a134aa99568d02203a3eecf6ca9ae8ce20bf6ba329b0ff27290fae5d389566e84b7c625b9681752201210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f00000000", txHex)

	// segwit_native address
	txBuild = NewTxBuild(1, &chaincfg.TestNet3Params)
	txBuild.AddInput2("0bc66f18fd95ca00b6569471aa2dcd47fe45d3446fbaeec9ced228b00713fe8c", 0, "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22", "tb1qtsq9c4fje6qsmheql8gajwtrrdrs38kdzeersc", 200000)
	txBuild.AddOutput("tb1qtsq9c4fje6qsmheql8gajwtrrdrs38kdzeersc", 199700)
	tx, err = txBuild.Build()
	assert.Nil(t, err)
	txHex, err = GetTxHex(tx)
	assert.Equal(t, "010000000001018cfe1307b028d2cec9eeba6f44d345fe47cd2daa719456b600ca95fd186fc60b0000000000ffffffff01140c0300000000001600145c005c5532ce810ddf20f9d1d939631b47089ecd0247304402200fcc9be29e3ab99b81f30fdf0788d883576d1a313fc809fce616813b8b2db62002207d199936c665f8e2c353c03d968f49dfcc16dacaf5877487c7d49f532180a6a101210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f00000000", txHex)

	// taproot_keypath address
	txBuild = NewTxBuild(1, &chaincfg.TestNet3Params)
	txBuild.AddInput2("3cb62c77c5c3fc032100af4cae9eeb342829cbc5b49815f8db1bb8156314a784", 0, "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22", "tb1pklh8lqax5l7m2ycypptv2emc4gata2dy28svnwcp9u32wlkenvsspcvhsr", 546)
	txBuild.AddOutput("tb1pklh8lqax5l7m2ycypptv2emc4gata2dy28svnwcp9u32wlkenvsspcvhsr", 300)
	tx, err = txBuild.Build()
	assert.Nil(t, err)
	txHex, err = GetTxHex(tx)
	assert.Nil(t, err)
	assert.Equal(t, "0100000000010184a7146315b81bdbf81598b4c5cb292834eb9eae4caf002103fcc3c5772cb63c0000000000ffffffff012c01000000000000225120b7ee7f83a6a7fdb513040856c56778aa3abea9a451e0c9bb012f22a77ed99b2101408399ed30c87c4dbcb7de58b9beb5448a7b7d9c4a0d8048fca51571648007c0a81cf0024d6c72d3cea13d6d3d69f341bc77a30612f3dfe10ee64874b0feded3b300000000", txHex)
}

func TestBtcTx(t *testing.T) {
	txBuild := NewTxBuild(1, &chaincfg.TestNet3Params)
	txBuild.AddInput("0b2c23f5c2e6326c90cfa1d3925b0d83f4b08035ca6af8fd8f606385dfbc5822", 1, "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37", "", "", 0) // replace to your private key
	txBuild.AddOutput("mvNnCR7EJS4aUReLEw2sL2ZtTZh8CAP8Gp", 53000)
	txBuild.AddOutput("mvNnCR7EJS4aUReLEw2sL2ZtTZh8CAP8Gp", 10000)
	txHex, err := txBuild.SingleBuild()
	require.Nil(t, err)
	assert.Equal(t, "01000000012258bcdf8563608ffdf86aca3580b0f4830d5b92d3a1cf906c32e6c2f5232c0b010000006a473044022028022b1b92fa0a10927e5ffa26da98aba737eeed6485b92af38071349a0cf1cd02202119a5b8b33f4a186d061f08e800f7dd1d9b908048035bb42516a1706278914d0121031053e9ef0295d334b6bb22e20cc717eb1a16a546f692572c8830b4bc14c13676ffffffff0208cf0000000000001976a914a2fe215e4789e607401a4bf85358cbbfae13a97e88ac10270000000000001976a914a2fe215e4789e607401a4bf85358cbbfae13a97e88ac00000000", txHex)
}

func TestBtcScript(t *testing.T) {
	txBuild := NewTxBuild(1, &chaincfg.TestNet3Params)
	txBuild.AddInput2("02133b22fdd190519ef9b49aca9a8dfdcbab0197c77109bb829cd51e17debed1", 0, "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22", "tb1qtsq9c4fje6qsmheql8gajwtrrdrs38kdzeersc", 8000)
	txBuild.AddOutput("tb1qtsq9c4fje6qsmheql8gajwtrrdrs38kdzeersc", 6000)
	txBuild.AddOutput2("", "6a01520b0080c7f6cf9b7c858c2002", 0)
	tx, err := txBuild.Build()
	assert.Nil(t, err)
	txHex, err := GetTxHex(tx)
	assert.Equal(t, "01000000000101d1bede171ed59c82bb0971c79701abcbfd8d9aca9ab4f99e5190d1fd223b13020000000000ffffffff0270170000000000001600145c005c5532ce810ddf20f9d1d939631b47089ecd00000000000000000f6a01520b0080c7f6cf9b7c858c20020248304502210082fe8e18b707302c253f7fd9e8ffdd25204986d754675ffbedd08031b4e0708302200ff3ae7e9f1e4d4bcd9182237e8c39577eb2a0242fc8e3c6a646646fa8efe3d201210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f00000000", txHex)
}

func TestMultiAddress(t *testing.T) {
	var pubKeys = []string{"022bc0ca1d6aea1c1e523bfcb33f46131bd1a3240aa04f71c34b1a177cfd5ff933", "035dc63727e7719824978161cdd94609db5235537bc8339a07b6838a6075f02530", "033eeee979afb70450d2aebb17ace1b170a96199b495cdf3dd0631eb96aa21e6a8"}
	redeemScript, err := GetRedeemScript(pubKeys, 2)
	require.Nil(t, err)
	fmt.Println(hex.EncodeToString(redeemScript))
	multiAddress, err := GenerateMultiAddress(redeemScript, &chaincfg.TestNet3Params)
	require.Nil(t, err)
	assert.Equal(t, "2N6DPSDtyXxUdJdACE1eHQ71z8vEVhDJZKF", multiAddress)
}

func TestMultiTx(t *testing.T) {
	txBuild := NewTxBuild(1, &chaincfg.TestNet3Params)
	txBuild.AddInput("f9c199cb3f43c0a1cd1b84f9912c77e3a62381cfe350ecc15a49c9bbd2633377", 0, "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37", "5221022bc0ca1d6aea1c1e523bfcb33f46131bd1a3240aa04f71c34b1a177cfd5ff93321035dc63727e7719824978161cdd94609db5235537bc8339a07b6838a6075f0253021033eeee979afb70450d2aebb17ace1b170a96199b495cdf3dd0631eb96aa21e6a853ae", "", 0) // replace to your private key
	txBuild.AddOutput("2NCVwfBKJ3zQBz1bimvKHC4kW7XHwbtQvF7", 2698100000)
	firstHex, err := txBuild.SingleBuild()
	require.Nil(t, err)
	assert.Equal(t, "0100000001773363d2bbc9495ac1ec50e3cf8123a6e3772c91f9841bcda1c0433fcb99c1f900000000b40047304402203f4ea02bc3ec719a4c1c5a3f798f613ae80fc4a3e03c6199c78fe31808912a4d02201c419de2ac9d1e4b2d8369a9205ac9a6da8ee69cfc421e73e98f6245aa40617f014c695221022bc0ca1d6aea1c1e523bfcb33f46131bd1a3240aa04f71c34b1a177cfd5ff93321035dc63727e7719824978161cdd94609db5235537bc8339a07b6838a6075f0253021033eeee979afb70450d2aebb17ace1b170a96199b495cdf3dd0631eb96aa21e6a853aeffffffff0120bdd1a00000000017a914d332ff8b6514f0b5fb07f090c5022d3d1c6c4bdf8700000000", firstHex)

	tx, err := NewTxFromHex(firstHex)
	require.Nil(t, err)
	var priKeyList []string
	priKeyList = append(priKeyList, "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37") // replace to your private key
	secondHex, err := MultiSignBuild(tx, priKeyList)
	require.Nil(t, err)
	assert.Equal(t, "0100000001773363d2bbc9495ac1ec50e3cf8123a6e3772c91f9841bcda1c0433fcb99c1f900000000fc0047304402203f4ea02bc3ec719a4c1c5a3f798f613ae80fc4a3e03c6199c78fe31808912a4d02201c419de2ac9d1e4b2d8369a9205ac9a6da8ee69cfc421e73e98f6245aa40617f0147304402203f4ea02bc3ec719a4c1c5a3f798f613ae80fc4a3e03c6199c78fe31808912a4d02201c419de2ac9d1e4b2d8369a9205ac9a6da8ee69cfc421e73e98f6245aa40617f014c695221022bc0ca1d6aea1c1e523bfcb33f46131bd1a3240aa04f71c34b1a177cfd5ff93321035dc63727e7719824978161cdd94609db5235537bc8339a07b6838a6075f0253021033eeee979afb70450d2aebb17ace1b170a96199b495cdf3dd0631eb96aa21e6a853aeffffffff0120bdd1a00000000017a914d332ff8b6514f0b5fb07f090c5022d3d1c6c4bdf8700000000", secondHex)
}

func TestAddress(t *testing.T) {
	privateBytes, err := hex.DecodeString("1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37")
	require.Nil(t, err)
	_, publicKey := btcec.PrivKeyFromBytes(privateBytes)
	pubKey := publicKey.SerializeCompressed()
	addr, err := btcutil.NewAddressPubKey(pubKey, &chaincfg.TestNet3Params)
	require.Nil(t, err)
	addressEncoded := addr.EncodeAddress()
	assert.Equal(t, "mwHiLyYXKVbhN6zwJkenkPsydj9MBLXb1K", addressEncoded)
}

func TestSingTx(t *testing.T) {
	txBuild := NewTxBuild(1, &chaincfg.TestNet3Params)
	txBuild.AddInput("0b2c23f5c2e6326c90cfa1d3925b0d83f4b08035ca6af8fd8f606385dfbc5822", 1, "", "", "", 0)
	txBuild.AddOutput("mvNnCR7EJS4aUReLEw2sL2ZtTZh8CAP8Gp", 53000)
	txBuild.AddOutput("mvNnCR7EJS4aUReLEw2sL2ZtTZh8CAP8Gp", 10000)
	pubKeyMap := make(map[int]string)
	pubKeyMap[0] = "022bc0ca1d6aea1c1e523bfcb33f46131bd1a3240aa04f71c34b1a177cfd5ff933"
	txHex, hashes, err := txBuild.UnSignedTx(pubKeyMap)
	require.Nil(t, err)
	signatureMap := make(map[int]string)
	for i, h := range hashes {
		privateBytes, err := hex.DecodeString("1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37")
		require.Nil(t, err)
		prvKey, _ := btcec.PrivKeyFromBytes(privateBytes)
		sign := ecdsa.Sign(prvKey, util.RemoveZeroHex(h))
		signatureMap[i] = hex.EncodeToString(sign.Serialize())
	}
	txHex, err = SignTx(txHex, pubKeyMap, signatureMap)
	require.Nil(t, err)
	assert.Equal(t, "01000000012258bcdf8563608ffdf86aca3580b0f4830d5b92d3a1cf906c32e6c2f5232c0b010000006a47304402206bdac667fb3d6f1a62e0b0d1123a5caa58d8c0fd95c2a2c8cd091374960a871702204f301e6883866570ce309573e569d6a32a44386af5bf928b5f9e1dcd7e2dd0ed0121022bc0ca1d6aea1c1e523bfcb33f46131bd1a3240aa04f71c34b1a177cfd5ff933ffffffff0208cf0000000000001976a914a2fe215e4789e607401a4bf85358cbbfae13a97e88ac10270000000000001976a914a2fe215e4789e607401a4bf85358cbbfae13a97e88ac00000000", txHex)
}

func TestGenerateAddress(t *testing.T) {
	address, err := GenerateAddress("022bc0ca1d6aea1c1e523bfcb33f46131bd1a3240aa04f71c34b1a177cfd5ff933", &chaincfg.TestNet3Params)
	require.Nil(t, err)
	assert.Equal(t, "1FrpuN2FVQdKhKAiXN4VW7MZba6RMevpkR", address)

}
