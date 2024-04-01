package serialize

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestCreateAddFullAccessKey(t *testing.T) {
	act, err := CreateAddFullAccessKey("58064be4ab6a0097b6c794f5cf1983ef36c60ea82c17e8488107433f6386b5ba")
	assert.NoError(t, err)
	b, err := act.Serialize()
	assert.NoError(t, err)
	assert.Equal(t, "050058064be4ab6a0097b6c794f5cf1983ef36c60ea82c17e8488107433f6386b5ba000000000000000001", hex.EncodeToString(b))
}

func TestCreateAddFunctionCallAccessKey(t *testing.T) {
	act, err := CreateAddFunctionCallAccessKey("58064be4ab6a0097b6c794f5cf1983ef36c60ea82c17e8488107433f6386b5ba", "1", "a", []string{"a"})
	assert.NoError(t, err)
	b, err := act.Serialize()
	assert.NoError(t, err)
	assert.Equal(t, "050058064be4ab6a0097b6c794f5cf1983ef36c60ea82c17e8488107433f6386b5ba00000000000000000001010000000000000000000000000000000100000061010000000100000061", hex.EncodeToString(b))
	act, err = CreateAddFunctionCallAccessKey("58064be4ab6a0097b6c794f5cf1983ef36c60ea82c17e8488107433f6386b5ba", "1", "a", []string{})
	assert.NoError(t, err)
	b, err = act.Serialize()
	assert.NoError(t, err)
	assert.Equal(t, "050058064be4ab6a0097b6c794f5cf1983ef36c60ea82c17e8488107433f6386b5ba0000000000000000000101000000000000000000000000000000010000006100000000", hex.EncodeToString(b))
	act, err = CreateAddFunctionCallAccessKey("58064be4ab6a0097b6c794f5cf1983ef36c60ea82c17e8488107433f6386b5ba", "", "a", []string{})
	assert.NoError(t, err)
	b, err = act.Serialize()
	assert.NoError(t, err)
	assert.Equal(t, "050058064be4ab6a0097b6c794f5cf1983ef36c60ea82c17e8488107433f6386b5ba00000000000000000000010000006100000000", hex.EncodeToString(b))
}

func TestCreateFunctionCall(t *testing.T) {
	data := []byte(`{"amount":"1000000000000000000","receiver_id":"316e10e0e93bef0927f4b0bc48849759a42c218b0e81a39ccb8eb15f048b00e8"}`)
	act, err := CreateFunctionCall("ft_transfer", data, big.NewInt(10000000000), big.NewInt(1))
	assert.NoError(t, err)
	b, err := act.Serialize()
	assert.NoError(t, err)
	assert.Equal(t, "020b00000066745f7472616e73666572710000007b22616d6f756e74223a2231303030303030303030303030303030303030222c2272656365697665725f6964223a2233313665313065306539336265663039323766346230626334383834393735396134326332313862306538316133396363623865623135663034386230306538227d00e40b540200000001000000000000000000000000000000", hex.EncodeToString(b))
}

func TestCreateStake(t *testing.T) {
	act, err := CreateStake("58064be4ab6a0097b6c794f5cf1983ef36c60ea82c17e8488107433f6386b5ba", "10000000000")
	assert.NoError(t, err)
	b, err := act.Serialize()
	assert.NoError(t, err)
	assert.Equal(t, "0400e40b540200000000000000000000000058064be4ab6a0097b6c794f5cf1983ef36c60ea82c17e8488107433f6386b5ba", hex.EncodeToString(b))
}

func TestCreateDeleteKey(t *testing.T) {
	act, err := CreateDeleteKey("58064be4ab6a0097b6c794f5cf1983ef36c60ea82c17e8488107433f6386b5ba")
	assert.NoError(t, err)
	b, err := act.Serialize()
	assert.NoError(t, err)
	assert.Equal(t, "060058064be4ab6a0097b6c794f5cf1983ef36c60ea82c17e8488107433f6386b5ba", hex.EncodeToString(b))
}

func TestCreateDeleteAccount(t *testing.T) {
	act, err := CreateDeleteAccount("58064be4ab6a0097b6c794f5cf1983ef36c60ea82c17e8488107433f6386b5ba")
	assert.NoError(t, err)
	b, err := act.Serialize()
	assert.NoError(t, err)
	assert.Equal(t, "074000000035383036346265346162366130303937623663373934663563663139383365663336633630656138326331376538343838313037343333663633383662356261", hex.EncodeToString(b))
}

func TestCreateDeployContract(t *testing.T) {
	act, err := CreateDeployContract([]byte("58064be4ab6a0097b6c794f5cf1983ef36c60ea82c17e8488107433f6386b5ba"))
	assert.NoError(t, err)
	b, err := act.Serialize()
	assert.NoError(t, err)
	assert.Equal(t, "014000000035383036346265346162366130303937623663373934663563663139383365663336633630656138326331376538343838313037343333663633383662356261", hex.EncodeToString(b))
}
