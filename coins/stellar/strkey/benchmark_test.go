package strkey_test

import (
	"github.com/okx/go-wallet-sdk/coins/stellar/strkey"
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkDecode_accountID(b *testing.B) {
	accountID, err := strkey.Encode(strkey.VersionByteAccountID, make([]byte, 32))
	assert.Nil(b, err)
	for i := 0; i < b.N; i++ {
		_, _ = strkey.Decode(strkey.VersionByteAccountID, accountID)
	}
}

func BenchmarkEncode_accountID(b *testing.B) {
	accountID := make([]byte, 32)
	for i := 0; i < b.N; i++ {
		_, _ = strkey.Encode(strkey.VersionByteAccountID, accountID)
	}
}
