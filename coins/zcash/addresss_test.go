package zcash

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateAddress(t *testing.T) {
	assert.True(t, ValidateAddress("t3VgMqDL15ZcyQDeqBsBW3W6rzfftrWP2yB"))
}
