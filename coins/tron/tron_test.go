package tron

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"math/big"
	"testing"
	"time"

	"github.com/okx/go-wallet-sdk/coins/tron/token"
	"github.com/okx/go-wallet-sdk/util/abi"
)

func TestTron_Address(t *testing.T) {
	trxAddress := "TNrEPvnnX7Hwj1z6tb1aTXpMad7z4BxoNW"
	ret := ValidateAddress(trxAddress)
	fmt.Println(ret)

	ah, _ := GetAddressHash("TGpKmWjRRQLuMn2G2PX5yCWJ9HfVsawJjY")
	fmt.Println(hex.EncodeToString(ah))
}

// curl -X POST  http://3.225.171.164:8090/wallet/createtransaction -d
// '{
//     "to_address": "41e9d79cc47518930bc322d9bf7cddd260a0260a8d",
//     "owner_address": "41D1E7A6BC354106CB410E65FF8B181C600FF14292",
//     "amount": 1000
// }

// curl -X POST https://3.225.171.164:8090/wallet/triggersmartcontract -d '{
// 	"contract_address":"419E62BE7F4F103C36507CB2A753418791B1CDC182",
// 	"function_selector":"transfer(address,uint256)",
// 	"parameter":"00000000000000000000004115208EF33A926919ED270E2FA61367B2DA3753DA0000000000000000000000000000000000000000000000000000000000000032",
// 	"fee_limit":100000000,
// 	"call_value":0,
// 	"owner_address":"41977C20977F412C2A1AA4EF3D49FEE5EC4C31CDFB"
// 	}'

func TestTron_Decode(t *testing.T) {
	// txStr := "0a86010a02607d2208aa7fdc0d42355da640e8cf96a5a72e5a68080112640a2d747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e5472616e73666572436f6e747261637412330a1541415d595646691ea4f6f02ab0cacaa6a57c70d81e121541f7c90c365ce1c5a2e6d510b4b7a016204d566d4d18e0e0d50d70fda4b7a3a72e"
	// txStr := "0ad4010a025fc222080142e695bcb5a9924080b1f4a4a72e5aae01081f12a9010a31747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e54726967676572536d617274436f6e747261637412740a1541415d595646691ea4f6f02ab0cacaa6a57c70d81e121541a614f803b6fd780986a42c78ec9c7f77e6ded13c2244a9059cbb000000000000000000000041f7c90c365ce1c5a2e6d510b4b7a016204d566d4d0000000000000000000000000000000000000000000000000000000000970fe070949195a3a72e90018094ebdc03"
	txStr := "0aab010a2000000000025263fdceee4a70744223c4ab3458f3480c002c5ec500c760f3dd7f222000000000025263fdceee4a70744223c4ab3458f3480c002c5ec500c760f3dd7f5a65080112610a2d747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e5472616e73666572436f6e747261637412300a15412fb81e30aa59250b4ddd5b312b0f80ac6205ee8612154179309abcff2cf531070ca9222a1f72c4a513687418011241a1d54fd4d2f3a38c6b4ba00d29b07b46c5fb1ba4547ea2df547b24298fa1f7a928cc7b4cf185c35d4a65afa79f1b1b913747fcd2a046552ecee06806d9b8a6f800"
	tran, _ := ParseTxStr(txStr)

	fmt.Println("tran: ", tran)

	dataToSign, _ := SignStart(txStr)
	fmt.Println("hash: ", dataToSign)

	pk, _ := btcec.NewPrivateKey()
	signStr, _ := Sign(dataToSign, pk)
	fmt.Println("sign: ", signStr)

	data, _ := SignEnd(txStr, signStr)
	fmt.Println("signed: ", data)

	b, _ := token.Transfer("2ed5dd8a98aea00ae32517742ea5289761b2710e", big.NewInt(50000000000))
	c, _ := token.Abi20.PackParams("transfer", "2ed5dd8a98aea00ae32517742ea5289761b2710e", big.NewInt(50000000000))
	fmt.Println("abi: ", hex.EncodeToString(b))
	fmt.Println("params: ", hex.EncodeToString(c))

	k1 := make([]byte, 8)
	binary.BigEndian.PutUint64(k1, 39184438)
	k2, _ := hex.DecodeString("000000000255e836afa34ffc90bb059cac301f849134253d3c505509306b091b")

	currentTime := time.Now()
	a, _ := abi.ParseBig256("1000000")
	d1, _ := newTRC20TokenTransfer(
		"TWhevFCRWEMAu9gqJ2Wymba3QbvKaBR3z4",
		"TEjxQjU3CxkFrSDcPfHwZXSuPpCpdQ27NJ",
		"TLa2f6VPqDgRE67v1736s7bJ8Ray5wYjU7",
		a,
		10000000,
		hex.EncodeToString(k1[6:8]),
		hex.EncodeToString(k2[8:16]),
		currentTime.UnixMilli()+3600*1000,
		currentTime.UnixMilli())
	fmt.Println("d1: ", d1)

	d2, _ := newTransfer(
		"TEjxQjU3CxkFrSDcPfHwZXSuPpCpdQ27NJ",
		"TWhevFCRWEMAu9gqJ2Wymba3QbvKaBR3z4",
		10000000,
		hex.EncodeToString(k1[6:8]),
		hex.EncodeToString(k2[8:16]),
		currentTime.UnixMilli()+3600*1000,
		currentTime.UnixMilli())
	fmt.Println("d2: ", d2)
}

func TestTron_Sign(t *testing.T) {
	currentTime := time.Now()
	k1 := make([]byte, 8)
	binary.BigEndian.PutUint64(k1, 47102802)
	k2, _ := hex.DecodeString("0000000002cebb52bb1c53a37236902bac251e302a4541452b6df63f594562b9")
	d2, _ := newTransfer(
		"TSAaoJuxBUxSqU7JGxzTH3gx237PTJxfwV",
		"TWYrgz7RDP2NpumQRPY1jBmPKLWVSnrzWZ",
		10000000,
		hex.EncodeToString(k1[6:8]),
		hex.EncodeToString(k2[8:16]),
		currentTime.UnixMilli()+3600*1000,
		currentTime.UnixMilli())

	fmt.Println(d2)
}
