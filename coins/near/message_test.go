package near

import (
	"encoding/hex"
	"github.com/okx/go-wallet-sdk/coins/near/serialize"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSignMessage(t *testing.T) {
	nonce := make([]byte, 32)
	nonce[31] = 1
	payload := serialize.NewSignMessagePayload("hello world", nonce, "", "")
	privateKey := "790e2778e0bfdae3da6419ef68c2451e80449de81e7bed9150b1cbc72b56a219d25cfdae0f9832e98bbdc87f3a156bb765cd9964e00878bf66da74591537e0a9"
	bs, err := payload.Serialize()
	assert.Nil(t, err)
	expectData := "9d0100800b00000068656c6c6f20776f726c6400000000000000000000000000000000000000000000000000000000000000010000000000"
	assert.Equal(t, expectData, hex.EncodeToString(bs))

	payload = serialize.NewSignMessagePayload("hello world", nonce, "", "1")
	bs, err = payload.Serialize()
	assert.Nil(t, err)
	expectData = "9d0100800b00000068656c6c6f20776f726c64000000000000000000000000000000000000000000000000000000000000000100000000010100000031"
	assert.Equal(t, expectData, hex.EncodeToString(bs))

	s, err := SignMessage(payload, privateKey)
	assert.Nil(t, err)
	assert.Equal(t, "3GwbyivHZuPxAD1oWe3yLJGjUyDcXAJ1fx2RhuawA5b75VbcBi1iZVz6HPbCe3dGae1buWQYHvASfTQwgLY9FzS4", s)

}
