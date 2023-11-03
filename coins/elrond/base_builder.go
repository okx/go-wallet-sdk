package elrond

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"math/big"
)

type baseBuilder struct {
	args []string
	err  error
}

func (builder *baseBuilder) addBytes(bytes []byte) {
	if len(bytes) == 0 {
		bytes = []byte{0}
	}
	builder.args = append(builder.args, hex.EncodeToString(bytes))
}

func (builder *baseBuilder) checkAddress(address string) bool {
	return ValidateAddress(address)
}

func (builder *baseBuilder) addArgHexString(hexed string) {
	if builder.err != nil {
		return
	}
	_, err := hex.DecodeString(hexed)
	if err != nil {
		builder.err = fmt.Errorf("%w in builder.ArgHexString for string %s", err, hexed)
		return
	}
	builder.args = append(builder.args, hexed)
}

func (builder *baseBuilder) addArgAddress(address string) {
	if builder.err != nil {
		return
	}
	ret := builder.checkAddress(address)
	if !ret {
		builder.err = fmt.Errorf("checkAddress failed")
		return
	}
	_, bytes, _ := bech32.DecodeToBase256(address)
	builder.addBytes(bytes)
}

func (builder *baseBuilder) addArgBigInt(value *big.Int) {
	if builder.err != nil {
		return
	}
	if value == nil {
		builder.err = fmt.Errorf("addArgBigInt value is null")
		return
	}
	builder.addBytes(value.Bytes())
}

func (builder *baseBuilder) addArgInt64(value int64) {
	if builder.err != nil {
		return
	}
	b := big.NewInt(value)
	builder.addBytes(b.Bytes())
}

func (builder *baseBuilder) addArgBytes(bytes []byte) {
	if builder.err != nil {
		return
	}
	if len(bytes) == 0 {
		builder.err = fmt.Errorf("addArgBytes bytes is empty")
	}
	builder.addBytes(bytes)
}
