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
	"testing"
)

func TestBtcTx(t *testing.T) {

	txBuild := NewTxBuild(1, &chaincfg.TestNet3Params)
	txBuild.AddInput("0b2c23f5c2e6326c90cfa1d3925b0d83f4b08035ca6af8fd8f606385dfbc5822", 1, "7214b52a4821690bac8a3139f36e15ab2f78c396f51d33f2749943332c083039", "")
	txBuild.AddOutput("mvNnCR7EJS4aUReLEw2sL2ZtTZh8CAP8Gp", 53000)
	txBuild.AddOutput("mvNnCR7EJS4aUReLEw2sL2ZtTZh8CAP8Gp", 10000)
	txHex := txBuild.SingleBuild()
	assert.Equal(t, "01000000012258bcdf8563608ffdf86aca3580b0f4830d5b92d3a1cf906c32e6c2f5232c0b010000006a47304402200b73518a4c2ed2f85afe3d8074e3169d2a661841cbef1127a1e77c651a7cba2102207d1e8e814e620f1d2db27adb07b01b1ae426d328c2732d27db852dd7ef4fb6f80121022bc0ca1d6aea1c1e523bfcb33f46131bd1a3240aa04f71c34b1a177cfd5ff933ffffffff0208cf0000000000001976a914a2fe215e4789e607401a4bf85358cbbfae13a97e88ac10270000000000001976a914a2fe215e4789e607401a4bf85358cbbfae13a97e88ac00000000", txHex)

}

func TestMultiAddress(t *testing.T) {
	var pubKeys = []string{"022bc0ca1d6aea1c1e523bfcb33f46131bd1a3240aa04f71c34b1a177cfd5ff933", "035dc63727e7719824978161cdd94609db5235537bc8339a07b6838a6075f02530", "033eeee979afb70450d2aebb17ace1b170a96199b495cdf3dd0631eb96aa21e6a8"}
	redeemScript, err := GetRedeemScript(pubKeys, 2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hex.EncodeToString(redeemScript))
	multiAddress, err := GenerateMultiAddress(redeemScript, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "2N6DPSDtyXxUdJdACE1eHQ71z8vEVhDJZKF", multiAddress)
}

/*
*
2N6DPSDtyXxUdJdACE1eHQ71z8vEVhDJZKF
5221022bc0ca1d6aea1c1e523bfcb33f46131bd1a3240aa04f71c34b1a177cfd5ff93321035dc63727e7719824978161cdd94609db5235537bc8339a07b6838a6075f0253021033eeee979afb70450d2aebb17ace1b170a96199b495cdf3dd0631eb96aa21e6a853ae

7214b52a4821690bac8a3139f36e15ab2f78c396f51d33f2749943332c083039 022bc0ca1d6aea1c1e523bfcb33f46131bd1a3240aa04f71c34b1a177cfd5ff933
5dfce364a4e9020d1bc187c9c14060e1a2f8815b3b0ceb40f45e7e39eb122103 035dc63727e7719824978161cdd94609db5235537bc8339a07b6838a6075f02530
11a1d32084a3a5f58a05398600da06ec9da4c4bfe61f5e4110a1d94cfec68dbf 033eeee979afb70450d2aebb17ace1b170a96199b495cdf3dd0631eb96aa21e6a8
*/
func TestMultiTx(t *testing.T) {
	txBuild := NewTxBuild(1, &chaincfg.TestNet3Params)
	txBuild.AddInput("f9c199cb3f43c0a1cd1b84f9912c77e3a62381cfe350ecc15a49c9bbd2633377", 0, "7214b52a4821690bac8a3139f36e15ab2f78c396f51d33f2749943332c083039", "5221022bc0ca1d6aea1c1e523bfcb33f46131bd1a3240aa04f71c34b1a177cfd5ff93321035dc63727e7719824978161cdd94609db5235537bc8339a07b6838a6075f0253021033eeee979afb70450d2aebb17ace1b170a96199b495cdf3dd0631eb96aa21e6a853ae")
	txBuild.AddOutput("2NCVwfBKJ3zQBz1bimvKHC4kW7XHwbtQvF7", 2698100000)
	firstHex := txBuild.SingleBuild()
	assert.Equal(t, "0100000001773363d2bbc9495ac1ec50e3cf8123a6e3772c91f9841bcda1c0433fcb99c1f900000000b400473044022056c8b555ae624eb4d399827d22cad3d24bf1e520710383c5a838224db30d783102202550c560bbd813d95da6af59274fb50c429e9447c7084a52da2ce108beb6bf8f014c695221022bc0ca1d6aea1c1e523bfcb33f46131bd1a3240aa04f71c34b1a177cfd5ff93321035dc63727e7719824978161cdd94609db5235537bc8339a07b6838a6075f0253021033eeee979afb70450d2aebb17ace1b170a96199b495cdf3dd0631eb96aa21e6a853aeffffffff0120bdd1a00000000017a914d332ff8b6514f0b5fb07f090c5022d3d1c6c4bdf8700000000", firstHex)

	tx := NewTxFromHex(firstHex)
	var priKeyList []string
	priKeyList = append(priKeyList, "5dfce364a4e9020d1bc187c9c14060e1a2f8815b3b0ceb40f45e7e39eb122103")
	secondHex := MultiSignBuild(tx, priKeyList)
	assert.Equal(t, "0100000001773363d2bbc9495ac1ec50e3cf8123a6e3772c91f9841bcda1c0433fcb99c1f900000000fc00473044022056c8b555ae624eb4d399827d22cad3d24bf1e520710383c5a838224db30d783102202550c560bbd813d95da6af59274fb50c429e9447c7084a52da2ce108beb6bf8f01473044022055ae37dcee60c77f038359daa54188d1716f978690ac1e28cd57e12aa297be9e02202bda548506a5ee4620b296467535d3bfd31568c501b5434c7125b202fc718c64014c695221022bc0ca1d6aea1c1e523bfcb33f46131bd1a3240aa04f71c34b1a177cfd5ff93321035dc63727e7719824978161cdd94609db5235537bc8339a07b6838a6075f0253021033eeee979afb70450d2aebb17ace1b170a96199b495cdf3dd0631eb96aa21e6a853aeffffffff0120bdd1a00000000017a914d332ff8b6514f0b5fb07f090c5022d3d1c6c4bdf8700000000", secondHex)
}

func TestAddress(t *testing.T) {
	privateBytes, err := hex.DecodeString("7214b52a4821690bac8a3139f36e15ab2f78c396f51d33f2749943332c083039")
	if err != nil {
		t.Fatal(err)
	}
	_, publicKey := btcec.PrivKeyFromBytes(privateBytes)

	pubKey := publicKey.SerializeCompressed()
	addr, err := btcutil.NewAddressPubKey(pubKey, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}
	addresEncode := addr.EncodeAddress()
	assert.Equal(t, "mvNnCR7EJS4aUReLEw2sL2ZtTZh8CAP8Gp", addresEncode)
}

func TestUnsingedTx(t *testing.T) {
	txBuild := NewTxBuild(1, &chaincfg.TestNet3Params)
	txBuild.AddInput("0b2c23f5c2e6326c90cfa1d3925b0d83f4b08035ca6af8fd8f606385dfbc5822", 1, "", "")
	txBuild.AddOutput("mvNnCR7EJS4aUReLEw2sL2ZtTZh8CAP8Gp", 53000)
	txBuild.AddOutput("mvNnCR7EJS4aUReLEw2sL2ZtTZh8CAP8Gp", 10000)

	pubKeyMap := make(map[int]string)
	pubKeyMap[0] = "022bc0ca1d6aea1c1e523bfcb33f46131bd1a3240aa04f71c34b1a177cfd5ff933"

	txHex, hashes, err := txBuild.UnSignedTx(pubKeyMap)
	if err != nil {
		t.Fatal(err)
	}

	signatureMap := make(map[int]string)
	for i, h := range hashes {
		privateBytes, _ := hex.DecodeString("7214b52a4821690bac8a3139f36e15ab2f78c396f51d33f2749943332c083039")
		prvKey, _ := btcec.PrivKeyFromBytes(privateBytes)
		sign := ecdsa.Sign(prvKey, util.RemoveZeroHex(h))
		signatureMap[i] = hex.EncodeToString(sign.Serialize())
	}
	txHex, err = SignTx(txHex, pubKeyMap, signatureMap)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "01000000012258bcdf8563608ffdf86aca3580b0f4830d5b92d3a1cf906c32e6c2f5232c0b010000006a47304402200b73518a4c2ed2f85afe3d8074e3169d2a661841cbef1127a1e77c651a7cba2102207d1e8e814e620f1d2db27adb07b01b1ae426d328c2732d27db852dd7ef4fb6f80121022bc0ca1d6aea1c1e523bfcb33f46131bd1a3240aa04f71c34b1a177cfd5ff933ffffffff0208cf0000000000001976a914a2fe215e4789e607401a4bf85358cbbfae13a97e88ac10270000000000001976a914a2fe215e4789e607401a4bf85358cbbfae13a97e88ac00000000", txHex)
}

func TestGenerateAddress(t *testing.T) {
	address, err := GenerateAddress("022bc0ca1d6aea1c1e523bfcb33f46131bd1a3240aa04f71c34b1a177cfd5ff933", &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(address)
	assert.Equal(t, "1FrpuN2FVQdKhKAiXN4VW7MZba6RMevpkR", address)

}
