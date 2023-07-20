package ethereum

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/okx/go-wallet-sdk/coins/ethereum/token"
	"github.com/stretchr/testify/assert"
)

func TestEth(t *testing.T) {
	p, _ := hex.DecodeString("//todo please replace your key")
	prvKey, _ := btcec.PrivKeyFromBytes(p)
	address := GetNewAddress(prvKey.PubKey())
	assert.Equal(t, "0x2de4898dd458d6dce097e29026d446300e3815fa", address)

	transaction := NewEthTransaction(
		big.NewInt(int64(00)),
		big.NewInt(int64(420000)),
		big.NewInt(int64(200000000000)),
		big.NewInt(int64(100000000000000000)),
		"2de4898dd458d6dce097e29026d446300e3815fa", "0x",
	)

	hash, raw, _ := transaction.GetSigningHash(big.NewInt(int64(10)))
	assert.Equal(t, "07f7adf7bf9efaf9442f792f3c7cd36b4505ec114c63493effa30b10b72d23e5", hash)
	assert.Equal(t, "ed80852e90edd000830668a0942de4898dd458d6dce097e29026d446300e3815fa88016345785d8a0000800a8080", raw)

	tx, err := transaction.SignTransaction(big.NewInt(int64(10)), prvKey)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "0xf86d80852e90edd000830668a0942de4898dd458d6dce097e29026d446300e3815fa88016345785d8a00008038a0c115f04d7dd555746085b01a4c8add779353f7033560b8eb61b244fe37772138a00c333c52c61c2e1ad2af83343b5f1b396a8b0fb60252ff8529b30f2e90967580", tx)

	b, _ := hex.DecodeString("2e3390fc71f35035b2ec378cced62632ef19c8d54b6b2f447e1f809c3d11ed0e")
	d, _ := SignAsRecoverable(b, prvKey)
	signature := d.ToHex()
	assert.Equal(t, "94f48c1dc793960d4c0e0ea0b34b95a8975d8d254edad8e25ab92a085914b37f30c943419fb1188244a8ffbdb7312003f8365ac381df5cd673a14cda28cf9b4f1b", signature)
}

func TestEthToken(t *testing.T) {
	p, _ := hex.DecodeString("//todo please replace your key")
	prvKey, _ := btcec.PrivKeyFromBytes(p)

	transfer, _ := token.Transfer("2de4898dd458d6dce097e29026d446300e3815fa", big.NewInt(int64(100000000000000000)))
	transaction := NewEthTransaction(
		big.NewInt(int64(00)),
		big.NewInt(int64(420000)),
		big.NewInt(int64(200000000000)),
		big.NewInt(int64(0)),
		"2ca70e7d0c396c36e8b9d206d988607a013483cf", hex.EncodeToString(transfer),
	)
	tx, err := transaction.SignTransaction(big.NewInt(int64(10)), prvKey)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "0xf8aa80852e90edd000830668a0942ca70e7d0c396c36e8b9d206d988607a013483cf80b844a9059cbb0000000000000000000000002de4898dd458d6dce097e29026d446300e3815fa000000000000000000000000000000000000000000000000016345785d8a000037a0afe573c296b30e9c4ef664ec64f63c48112e84167cdb8ef1ea567efc651f63c6a0283d2cc14a342464d2f62e83166bb9bf8d7ecdbe02d9b518eac02d458e158ccc", tx)
}

func TestEth2(t *testing.T) {
	p, _ := hex.DecodeString("//todo please replace your key")
	prvKey, _ := btcec.PrivKeyFromBytes(p)

	chainId := big.NewInt(int64(10))
	transaction := NewEthTransaction(
		big.NewInt(int64(00)),
		big.NewInt(int64(420000)),
		big.NewInt(int64(200000000000)),
		big.NewInt(int64(100000000000000000)),
		"2de4898dd458d6dce097e29026d446300e3815fa", "",
	)

	unSignedHex, _ := transaction.UnSignedTx(chainId)
	unSignedBytes, _ := hex.DecodeString(unSignedHex)
	sig, err := SignMessage(unSignedBytes, prvKey)
	if err != nil {
		t.Fatal(err)
	}
	tx, _ := transaction.SignedTx(chainId, sig)
	assert.Equal(t, "0xf86d80852e90edd000830668a0942de4898dd458d6dce097e29026d446300e3815fa88016345785d8a00008038a0c115f04d7dd555746085b01a4c8add779353f7033560b8eb61b244fe37772138a00c333c52c61c2e1ad2af83343b5f1b396a8b0fb60252ff8529b30f2e90967580", tx)
}

func TestEth3(t *testing.T) {
	t1, _ := NewTransactionFromRaw("0xf86f6a8506fc23ac0082520894e1d4fd72a48af968d80f6d9ef161d57bb9293837880de0b6b3a764000080830150f5a0b64b1eb1c2f41b95dac35c4751f4070ca8c185b9a94ea2d44454f47ca7944a23a004a98084513c3233ec962d9cfdf2a45d06473fa51554a0b3d99939e1ed387ed7")
	t2, _ := NewTransactionFromRaw("0xf869698506fc23ac0082520894af133678d4188ddbfd13655cf12e8e15f28fdecb8203e880830150f6a0b64b1eb1c2f41b95dac35c4751f4070ca8c185b9a94ea2d44454f47ca7944a23a0522afd4359e208cf86c036c4272bf65d7fdddc73f33006f845c1c94fa826befc")
	fmt.Println(t1)
	fmt.Println(t2)
}
