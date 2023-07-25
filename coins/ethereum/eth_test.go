package ethereum

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"math/big"
	"testing"

	"github.com/okx/go-wallet-sdk/coins/ethereum/token"
	"github.com/stretchr/testify/assert"
)

func TestEth(t *testing.T) {
	p, _ := hex.DecodeString("559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a")
	prvKey, _ := btcec.PrivKeyFromBytes(p)
	address := GetNewAddress(prvKey.PubKey())
	assert.Equal(t, "0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728", address)

	transaction := NewEthTransaction(
		big.NewInt(int64(00)),
		big.NewInt(int64(420000)),
		big.NewInt(int64(200000000000)),
		big.NewInt(int64(100000000000000000)),
		"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728", "0x",
	)

	hash, raw, _ := transaction.GetSigningHash(big.NewInt(int64(10)))
	assert.Equal(t, "790f2b826ad9dfa7f2a53ec68e37ea51dc58652ecfde812da37c96a1069fcdbb", hash)
	assert.Equal(t, "ed80852e90edd000830668a0941ca96f8cfe7276bb053b25e57188f1b5ec6a472888016345785d8a0000800a8080", raw)

	tx, err := transaction.SignTransaction(big.NewInt(int64(10)), prvKey)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "0xf86d80852e90edd000830668a0941ca96f8cfe7276bb053b25e57188f1b5ec6a472888016345785d8a00008037a0afd10738449dd9ab4f95b6f49244dc076ae5f1251397c7f010ba529edecf8517a03eb5492b35278b2636870843550040edb60f6b1026bff42ee5a803c6de1b0e04", tx)

	b, _ := hex.DecodeString("2e3390fc71f35035b2ec378cced62632ef19c8d54b6b2f447e1f809c3d11ed0e")
	d, _ := SignAsRecoverable(b, prvKey)
	signature := d.ToHex()
	assert.Equal(t, "32466d55329625198458901517ccae23f0162fc42b333f770e8e59ab62d3d40e6c2e85072ad2fd4273d8be86af5e005b1c9df39bf3f2014897347ec81ce6bc7f1b", signature)
}

func TestEthToken(t *testing.T) {
	p, _ := hex.DecodeString("559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a")
	prvKey, _ := btcec.PrivKeyFromBytes(p)

	transfer, _ := token.Transfer("0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728", big.NewInt(int64(100000000000000000)))
	transaction := NewEthTransaction(
		big.NewInt(int64(00)),
		big.NewInt(int64(420000)),
		big.NewInt(int64(200000000000)),
		big.NewInt(int64(0)),
		"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728", hex.EncodeToString(transfer),
	)
	tx, err := transaction.SignTransaction(big.NewInt(int64(10)), prvKey)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "0xf8aa80852e90edd000830668a0941ca96f8cfe7276bb053b25e57188f1b5ec6a472880b844a9059cbb0000000000000000000000001ca96f8cfe7276bb053b25e57188f1b5ec6a4728000000000000000000000000000000000000000000000000016345785d8a000038a0ad7d69a4eeb889a2bdd82e2c62d4063467936350f7d3cc466aa513e7abcbb077a071b5b06e8253352f3e1aed57a6db4fdf5113b00c961fedecf7fbde96c94cb66f", tx)
}

func TestEth2(t *testing.T) {
	p, _ := hex.DecodeString("559376194bb4c9a9dfb33fde4a2ab15daa8a899a3f43dee787046f57d5f7b10a")
	prvKey, _ := btcec.PrivKeyFromBytes(p)

	chainId := big.NewInt(int64(10))
	transaction := NewEthTransaction(
		big.NewInt(int64(00)),
		big.NewInt(int64(420000)),
		big.NewInt(int64(200000000000)),
		big.NewInt(int64(100000000000000000)),
		"0x1ca96f8cfe7276bb053b25e57188f1b5ec6a4728", "",
	)

	unSignedHex, _ := transaction.UnSignedTx(chainId)
	unSignedBytes, _ := hex.DecodeString(unSignedHex)
	sig, err := SignMessage(unSignedBytes, prvKey)
	if err != nil {
		t.Fatal(err)
	}
	tx, _ := transaction.SignedTx(chainId, sig)
	assert.Equal(t, "0xf86d80852e90edd000830668a0941ca96f8cfe7276bb053b25e57188f1b5ec6a472888016345785d8a00008037a0afd10738449dd9ab4f95b6f49244dc076ae5f1251397c7f010ba529edecf8517a03eb5492b35278b2636870843550040edb60f6b1026bff42ee5a803c6de1b0e04", tx)
}

func TestEth3(t *testing.T) {
	t1, _ := NewTransactionFromRaw("0xf86f6a8506fc23ac0082520894e1d4fd72a48af968d80f6d9ef161d57bb9293837880de0b6b3a764000080830150f5a0b64b1eb1c2f41b95dac35c4751f4070ca8c185b9a94ea2d44454f47ca7944a23a004a98084513c3233ec962d9cfdf2a45d06473fa51554a0b3d99939e1ed387ed7")
	t2, _ := NewTransactionFromRaw("0xf869698506fc23ac0082520894af133678d4188ddbfd13655cf12e8e15f28fdecb8203e880830150f6a0b64b1eb1c2f41b95dac35c4751f4070ca8c185b9a94ea2d44454f47ca7944a23a0522afd4359e208cf86c036c4272bf65d7fdddc73f33006f845c1c94fa826befc")
	fmt.Println(t1)
	fmt.Println(t2)
}
