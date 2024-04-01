package aptos_types

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFromHex(t *testing.T) {
	t.Log(CORE_CODE_ADDRESS)
	addr, err := FromHex("0x1")
	require.NoError(t, err)
	require.True(t, *addr == *CORE_CODE_ADDRESS)
}

func TestAccountAddress_BcsSerialize(t *testing.T) {
	b, err := CORE_CODE_ADDRESS.BcsSerialize()
	require.NoError(t, err)
	m := map[string][]byte{}
	m["Address"] = b
	jm, _ := json.Marshal(m)
	expected := "{\"Address\":\"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAE=\"}"
	require.Equal(t, expected, string(jm))
}
