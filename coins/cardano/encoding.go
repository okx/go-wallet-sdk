package cardano

import (
	"fmt"

	"github.com/okx/go-wallet-sdk/crypto/cbor"
)

var cborEnc, _ = cbor.CanonicalEncOptions().EncMode()
var cborDec, _ = cbor.DecOptions{MapKeyByteString: cbor.MapKeyByteStringWrap}.DecMode()

func getTypeFromCBORArray(data []byte) (uint64, error) {
	raw := []interface{}{}
	if err := cborDec.Unmarshal(data, &raw); err != nil {
		return 0, err
	}

	if len(raw) == 0 {
		return 0, fmt.Errorf("empty CBOR array")
	}

	t, ok := raw[0].(uint64)
	if !ok {
		return 0, fmt.Errorf("invalid Type")
	}

	return t, nil
}
