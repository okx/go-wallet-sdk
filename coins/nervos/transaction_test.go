package nervos

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	ckbTest1Address    = "ckt1qzda0cr08m85hc8jlnfp3zer7xulejywt49kt2rr0vthywaa50xwsq08pk6ldw7944vqmvulq555739qnlpap8sglxd9q"
	ckbTest2Address    = "ckt1qzda0cr08m85hc8jlnfp3zer7xulejywt49kt2rr0vthywaa50xwsqg4w6hxs0zvlh6kfrwfjfleq8qpaw2r7pcx24f6u"
	ckbTest1PrivateKey = "0171ecab8a308cd26fef99efb7ea02fa17ec9c210d8e9f6e32543694a6623ece"
	ckbTest2PrivateKey = "7fde15d42081b384af0d6fde3f575d3b34cdc068a6e5660127d77716fa27ef0b"
)

func TestNewTransaction(t *testing.T) {
	builder := NewTestnetTxBuild()
	// add inputs
	err := builder.AddInputWithPrivateKey("0x41c28858a53fd6e6a15e0df0c557bd7f2eba38b400f83bb64ce9d3c96914a5be",
		1, 0, ckbTest1PrivateKey)
	require.NoError(t, err)
	// add outputs
	err = builder.AddOutput(ckbTest2Address, 100*OneCKBShannon)
	require.NoError(t, err)
	err = builder.AddOutput(ckbTest1Address, 9685*OneCKBShannon)
	require.NoError(t, err)
	// build
	tx, err := builder.Build()
	require.NoError(t, err)
	// sign
	err = builder.Sign()
	require.NoError(t, err)
	// serialize
	txHash, err := tx.ComputeHash()
	require.NoError(t, err)
	expected := "0xed7a3a7f4ccd29a5fda103dd9e0abaa8435ecf95fa64bccaec7a7b2d994c844b"
	require.Equal(t, expected, txHash.Hex())
}

func TestNewTransactionWithTx(t *testing.T) {
	builder := NewTestnetTxBuild()
	// add inputs
	if err := builder.AddInput("0xcaf2cfb17eb961f54e22f8ced8656aa152f64f53e3db35b99705ca6b3822b5be", 0, 0); err != nil {
		t.Error(err)
		return
	}
	// add outputs
	if err := builder.AddOutput(ckbTest2Address, 100*OneCKBShannon); err != nil {
		t.Error(err)
		return
	}
	if err := builder.AddOutput(ckbTest1Address, 9895*OneCKBShannon); err != nil {
		t.Error(err)
		return
	}
	tx, err := builder.Build()
	if err != nil {
		t.Error(err)
		return
	}
	// sign
	if err := builder.SignByPrivateKey(ckbTest1PrivateKey); err != nil {
		t.Error(err)
		return
	}
	fmt.Println(tx.ComputeHash())
	// serialize
	fmt.Println(builder.DumpTx())
}
