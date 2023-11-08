package cbor

import (
	"errors"
	"github.com/okx/go-wallet-sdk/crypto/cbor"
	"math"
)

const invalidVersion = math.MaxUint16

var (
	// ErrInvalidVersion is the error returned when a versioned
	// serialized blob is either missing, or has an invalid version.
	ErrInvalidVersion = errors.New("cbor: missing or invalid version")

	decOptionsVersioned = decOptions

	decModeVersioned cbor.DecMode
)

// Versioned is a generic versioned serializable data structure.
type Versioned struct {
	V uint16 `json:"v"`
}

// GetVersion returns the version of a versioned serializable data
// structure, if any.
func GetVersion(data []byte) (uint16, error) {
	vblob := Versioned{
		V: invalidVersion,
	}
	if err := decModeVersioned.Unmarshal(data, &vblob); err != nil {
		return 0, err
	}
	if vblob.V == invalidVersion {
		return 0, ErrInvalidVersion
	}
	return vblob.V, nil
}

// NewVersioned creates a new Versioned structure with the specified version.
func NewVersioned(v uint16) Versioned {
	if v == invalidVersion {
		panic("cbor: invalid version specified")
	}
	return Versioned{V: v}
}

func init() {
	// Use the untrusted decode options, but ignore unknown fields.
	decOptionsVersioned.ExtraReturnErrors = cbor.ExtraDecErrorNone

	var err error
	if decModeVersioned, err = decOptionsVersioned.DecMode(); err != nil {
		panic(err)
	}
}
