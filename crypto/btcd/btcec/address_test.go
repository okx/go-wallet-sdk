package btcec

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnsupportedWitnessProgLenError(t *testing.T) {
	r := UnsupportedWitnessProgLenError(1)
	assert.Equal(t, "unsupported witness program length: 1", r.Error())
}
