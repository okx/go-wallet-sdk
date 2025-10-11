package filecoin

import (
	"encoding/json"
	"github.com/emresenyuva/go-wallet-sdk/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestNewPrivateKey(t *testing.T) {
	privateKey := NewPrivateKey()
	t.Log("privateKey : ", privateKey)
}

func TestGetPublicKey(t *testing.T) {
	privateKeyHex := "0x7c6207c56b6aa5ed4345c5f662816408e273cbdf64e2f01d54ced0125d6172c2"
	publicKeyHex, err := GetPublicKey(privateKeyHex)
	require.Nil(t, err)
	expected := "0x04e4f0c46b3dd02bf1579c848decf9b4c8d8e92cb1583f9b866b7c59f2d0ccc7fb87e96d2be53b72fd2cee4a83f08af5bc36fa0e927e4de31e3d424e4cc1e17f69"
	require.Equal(t, expected, publicKeyHex)
}

func TestGetAddressByPublicKey(t *testing.T) {
	publicKeyHex := "0x04e4f0c46b3dd02bf1579c848decf9b4c8d8e92cb1583f9b866b7c59f2d0ccc7fb87e96d2be53b72fd2cee4a83f08af5bc36fa0e927e4de31e3d424e4cc1e17f69"
	address, err := GetAddressByPublicKey(publicKeyHex, MainnetPrefix)
	require.Nil(t, err)
	expected := "f1bh3d2y6xxugpg3ygzxnjhcrs5ffxh5nvqmanbia"
	require.Equal(t, expected, address)
}

func TestGetAddressByPrivateKey(t *testing.T) {
	privateKeyHex := "0x7c6207c56b6aa5ed4345c5f662816408e273cbdf64e2f01d54ced0125d6172c2"
	address, err := GetAddressByPrivateKey(privateKeyHex, MainnetPrefix)
	require.Nil(t, err)
	expected := "f1bh3d2y6xxugpg3ygzxnjhcrs5ffxh5nvqmanbia"
	require.Equal(t, expected, address)
}

func TestAddressToBytes(t *testing.T) {
	addr := "f1bh3d2y6xxugpg3ygzxnjhcrs5ffxh5nvqmanbia"
	addrBytes := AddressToBytes(addr)
	expected := "0x0109f63d63d7bd0cf36f06cdda938a32e94b73f5b5"
	require.Equal(t, expected, util.EncodeHexWith0x(addrBytes))
}

func TestSignTx(t *testing.T) {
	from := "f1bh3d2y6xxugpg3ygzxnjhcrs5ffxh5nvqmanbia"
	to := "f1fvs2fjqr6ozk477zkwzvermhledmfkswt34cmhi"
	nonce := 0
	value := big.NewInt(20000000000000000)
	gasLimit := 210000
	gasFeeCap := big.NewInt(9455791480)
	gasPremium := big.NewInt(120242)
	method := 0
	message := NewTx(from, to, nonce, method, gasLimit, value, gasFeeCap, gasPremium)
	privateKeyHex := "0x7c6207c56b6aa5ed4345c5f662816408e273cbdf64e2f01d54ced0125d6172c2"
	tx, err := SignTx(message, privateKeyHex)
	require.Nil(t, err)
	txBytes, err := json.Marshal(tx)
	require.Nil(t, err)
	expected := "{\"Message\":{\"Version\":0,\"To\":\"f1fvs2fjqr6ozk477zkwzvermhledmfkswt34cmhi\",\"From\":\"f1bh3d2y6xxugpg3ygzxnjhcrs5ffxh5nvqmanbia\",\"Nonce\":0,\"Value\":\"20000000000000000\",\"GasLimit\":210000,\"GasFeeCap\":\"9455791480\",\"GasPremium\":\"120242\",\"Method\":0,\"Params\":\"\"},\"Signature\":{\"Type\":1,\"Data\":\"/0tQo2pRIYeSax/nt/+Jvdovz2CQwctKvthqsQRHfDd3B6h69K0ayW6z9CriMRX93USZv0uPXDs5SxCNkCJpgwE=\"}}"
	assert.Equal(t, expected, string(txBytes))
}
