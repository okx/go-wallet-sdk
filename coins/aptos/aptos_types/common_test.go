package aptos_types

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBcsSerializeFixedBytes(t *testing.T) {
	addr := BytesFromHex("0x1")
	b, _ := BcsSerializeFixedBytes(addr)
	t.Log(b)
	m := map[string][]byte{}
	m["Address"] = b
	jm, _ := json.Marshal(m)
	require.Equal(t, "{\"Address\":null}", string(jm))
}
