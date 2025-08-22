package util

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestFloat64ToBigInt(t *testing.T) {
	a, err := Float64ToBigInt(1.0)
	assert.Nil(t, err)
	assert.Equal(t, a.Int64(), int64(1))
	a, err = Float64ToBigInt(1.1)
	assert.Nil(t, err)
	assert.Equal(t, a.Int64(), int64(1))

	_, err = Float64ToBigInt(-1.1)
	assert.NotNil(t, err)
}

func TestIntToU32(t *testing.T) {
	a, err := IntToU32(math.MaxUint32)
	assert.Nil(t, err)
	assert.Equal(t, uint32(a), uint32(math.MaxUint32))

	_, err = IntToU32(-math.MaxUint32)
	assert.NotNil(t, err)
}

func TestFloat64ToU16(t *testing.T) {
	_, err := Float64ToU16(math.MaxUint16 + 1)
	assert.NotNil(t, err)
	_, err = Float64ToU32(math.MaxUint32 + 1)
	assert.NotNil(t, err)
}
