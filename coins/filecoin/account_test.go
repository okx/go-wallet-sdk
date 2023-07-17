package filecoin

import (
	"encoding/json"
	"github.com/okx/go-wallet-sdk/util"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestNewPrivateKey(t *testing.T) {
	privateKey := NewPrivateKey()
	t.Logf("private key hex: %s", privateKey)
}

func TestGetPublicKey(t *testing.T) {
	privateKeyHex := "//todo please replace your key"
	publicKeyHex, err := GetPublicKey(privateKeyHex)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("public key : %s", publicKeyHex)
}

func TestGetAddressByPublicKey(t *testing.T) {
	publicKeyHex := "0x04c7d2209a4b286046cdeaf457e499a40a9a1da5d7bc6e85c05e5ac9e6af9c7a35063c8a8efaa7cc4cd294c3b76dd4b0a3f5773cc421fef44e6a99914c8c85c971"
	address, err := GetAddressByPublicKey(publicKeyHex, MainnetPrefix)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("address : %s", address)
	assert.Equal(t, address, "f12cs7ppvnhwhma3xzhkm4pavq2q47blmprcxvg6i")
}

func TestGetAddressByPrivateKey(t *testing.T) {
	privateKeyHex := "//todo please replace your key"
	address, err := GetAddressByPrivateKey(privateKeyHex, MainnetPrefix)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("address : %s", address)
}

func TestAddressToBytes(t *testing.T) {
	addr := "f12cs7ppvnhwhma3xzhkm4pavq2q47blmprcxvg6i"
	bytes := AddressToBytes(addr)
	t.Logf(util.EncodeHexWith0x(bytes))
	assert.Equal(t, util.EncodeHexWith0x(bytes), "0x01d0a5f7bead3d8ec06ef93a99c782b0d439f0ad8f")
}

func TestSignTx(t *testing.T) {
	from := "f12cs7ppvnhwhma3xzhkm4pavq2q47blmprcxvg6i"
	to := "f1izmcd3o7pqiyob5yf3q3mat3w3rf5dzrccjhhhi"
	nonce := 0
	value := big.NewInt(20000000000000000)
	gasLimit := 210000
	gasFeeCap := big.NewInt(9455791480)
	gasPremium := big.NewInt(120242)
	method := 0

	message := NewTx(from, to, nonce, method, gasLimit, value, gasFeeCap, gasPremium)

	privateKeyHex := "//todo please replace your key"
	tx, err := SignTx(message, privateKeyHex)
	if err != nil {
		t.Fatal(err)
	}
	bytes, err := json.Marshal(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(string(bytes))
	assert.Equal(t, string(bytes), "{\"Message\":{\"Version\":0,\"To\":\"f1izmcd3o7pqiyob5yf3q3mat3w3rf5dzrccjhhhi\",\"From\":\"f12cs7ppvnhwhma3xzhkm4pavq2q47blmprcxvg6i\",\"Nonce\":0,\"Value\":\"20000000000000000\",\"GasLimit\":210000,\"GasFeeCap\":\"9455791480\",\"GasPremium\":\"120242\",\"Method\":0,\"Params\":\"\"},\"Signature\":{\"Type\":1,\"Data\":\"3hsJRgdhmhLE5FyMPtVpW1DgikBKBdeW0JSt8z3iFeENJ8tO9+Yc4RNVEupGBSBVoSTz1H3Y5giFyhLLv3UbBgA=\"}}")
}
