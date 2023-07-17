package filecoin

import (
	"github.com/okx/go-wallet-sdk/util"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestNewTx(t *testing.T) {

	from := "f1b3cn3x3uhyuqcgtjewjj3vpychvgztu3oxk36ta"
	to := from
	nonce := 1
	value := big.NewInt(100)
	gasLimit := 608335
	gasFeeCap := big.NewInt(1643831112)
	gasPremium := big.NewInt(99707)
	method := 0

	message := NewTx(from, to, nonce, method, gasLimit, value, gasFeeCap, gasPremium)
	h := util.EncodeHexWith0x(message.Hash())
	t.Logf(h)
	assert.Equal(t, h, "0x69f8861931b0d568c238eee7615a22957c561bb34e0f31cf4885d981ac4bdc02")
}
