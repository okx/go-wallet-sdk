package bitcoin

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTaprootAddress(t *testing.T) {
	publicKey := "04c9ad860eaf65d7769ff3b262ec30ce025422a4eb3cbd84e203b60025d3770b7d6a7e9b158bf7db80b936f7d54138be12b1b4bdcfb74e2b028de7b0231a62178b"
	script := "20c9ad860eaf65d7769ff3b262ec30ce025422a4eb3cbd84e203b60025d3770b7dac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d38000c68656c6c6f2c20776f726c6468"
	r, err := NewTaprootAddress(script, &chaincfg.MainNetParams, publicKey)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, "c0c9ad860eaf65d7769ff3b262ec30ce025422a4eb3cbd84e203b60025d3770b7d", r.ControlBlockWitness)
	assert.Equal(t, "bc1ps2mt0gnqechh9eezqhj50n09lyl06e02kz2h5zgu5je56rus245stgkez5", r.TaprootAddress)
}
