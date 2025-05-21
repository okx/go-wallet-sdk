package filecoin

import (
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
	signHex := "1b886880e6e39c9e842dbb16161cf5d3b469af30a1da0b6b5037cdb00e17691ee834c4eb46d2ffadbeef856a6d1b478835b627f9609c60d5e9905861215ed57d2b"
	signedMessage, err := SignedTx(message, signHex)
	assert.NoError(t, err)
	assert.Equal(t, `{"Message":{"Version":0,"To":"f1b3cn3x3uhyuqcgtjewjj3vpychvgztu3oxk36ta","From":"f1b3cn3x3uhyuqcgtjewjj3vpychvgztu3oxk36ta","Nonce":1,"Value":"100","GasLimit":608335,"GasFeeCap":"1643831112","GasPremium":"99707","Method":0,"Params":""},"Signature":{"Type":1,"Data":"iGiA5uOcnoQtuxYWHPXTtGmvMKHaC2tQN82wDhdpHug0xOtG0v+tvu+Fam0bR4g1tif5YJxg1emQWGEhXtV9KwA="}}`, signedMessage)

}

func TestCalTxHash(t *testing.T) {
	tx := `{"Message":{"Nonce":347341,"Version":0,"Value":"1239021584920000000000","To":"f410fodaz6xmnyfehneslmvptejn7mzaykgvqst6e7gq","From":"f1ys5qqiciehcml3sp764ymbbytfn3qoar5fo3iwy","Method":3844450837,"GasPremium":"264158","GasLimit":3075406,"GasFeeCap":"884744"},"Signature":{"Type":1,"Data":"iGiA5uOcnoQtuxYWHPXTtGmvMKHaC2tQN82wDhdpHug0xOtG0v+tvu+Fam0bR4g1tif5YJxg1emQWGEhXtV9KwA="}}`
	res, err := CalTxHash(tx)
	assert.NoError(t, err)
	assert.Equal(t, "bafy2bzaceccvft3let663mpexkgi5db22fcguonetb2vg7i5bazcyh2gjvz34", res)
}
