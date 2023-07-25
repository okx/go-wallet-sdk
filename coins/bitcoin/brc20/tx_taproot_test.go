package brc20

import (
	"fmt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReveal(t *testing.T) {
	builder := NewTxBuildV1(&chaincfg.TestNet3Params)

	contentType := "text/plain;charset=utf-8"
	body := []byte(fmt.Sprintf(`{"p":"brc-20","op":"%s","tick":"%s","amt":"%s"}`, "transfer", "ordi", "1"))

	inscription := NewInscription(contentType, body)
	builder.AddInput("9f9ff5acc7b3966ccfc6acc77027209d62aab34e563a09180c58ef7296fca74b",
		1,
		"1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37",
		"tb1pmwus5lpxnnet6wcyqtevls07y7u8h5wun7q7p9jglk707y2czfnsdlqqjw",
		"1600",
		inscription,
	)
	builder.AddOutput("tb1pp6v2zc4dfxrx0c6xmh340u9w958w2mklyfhz5ufrf7t8m6wunj2q4uvfj0", "546")
	builder.AddOutput("tb1pmwus5lpxnnet6wcyqtevls07y7u8h5wun7q7p9jglk707y2czfnsdlqqjw", "754")
	tx, _ := builder.Build()
	assert.Equal(t, "010000000001014ba7fc9672ef580c18093a564eb3aa629d202770c7acc6cf6c96b3c7acf59f9f0100000000ffffffff0222020000000000002251200e98a162ad498667e346dde357f0ae2d0ee56edf226e2a71234f967de9dc9c94f202000000000000225120dbb90a7c269cf2bd3b0402f2cfc1fe27b87bd1dc9f81e09648fdbcff115812670340b914b2133850ed6c2b96adaac45d141d18a3a789b3a6002cff57f6854779f90769ac530a5b25d3eee2702066e7d15f8771dfeaf76effc97ade5271211bf55f337c201053e9ef0295d334b6bb22e20cc717eb1a16a546f692572c8830b4bc14c13676ac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d3800367b2270223a226272632d3230222c226f70223a227472616e73666572222c227469636b223a226f726469222c22616d74223a2231227d6821c11053e9ef0295d334b6bb22e20cc717eb1a16a546f692572c8830b4bc14c1367600000000", tx)
}
