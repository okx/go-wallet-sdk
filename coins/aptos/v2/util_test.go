package v2

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSHA3_256Hash(t *testing.T) {
	input := [][]byte{{0x1}, {0x2}, {0x3}}
	expected, err := ParseHex("fd1780a6fc9ee0dab26ceb4b3941ab03e66ccd970d1db91612c66df4515b0a0a")
	assert.NoError(t, err)
	assert.Equal(t, expected, Sha3256Hash(input))
}

func TestParseHex(t *testing.T) {
	// TODO: Last case is very weird, unsure if we want to allow that
	inputs := []string{"0x012345", "012345", "0x"}
	expected := [][]byte{{0x01, 0x23, 0x45}, {0x01, 0x23, 0x45}, {}}

	for i, input := range inputs {
		val, err := ParseHex(input)
		assert.NoError(t, err)
		assert.Equal(t, expected[i], val)
	}
}
