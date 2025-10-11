/**
MIT License

Copyright (c) 2018 WavesPlatform

*/

package types

import (
	"github.com/emresenyuva/go-wallet-sdk/coins/waves/crypto"
	"strings"
)

const (
	jsonNull       = "null"
	WavesAssetName = "WAVES"
)

// OptionalAsset represents an optional asset identification
type OptionalAsset struct {
	Present bool
	ID      crypto.Digest
}

// MarshalJSON writes OptionalAsset as a JSON string Value
func (a OptionalAsset) MarshalJSON() ([]byte, error) {
	if a.Present {
		return a.ID.MarshalJSON()
	}
	return []byte(jsonNull), nil
}

func (a OptionalAsset) BinarySize() int {
	s := 1
	if a.Present {
		s += crypto.DigestSize
	}
	return s
}

// MarshalBinary marshals the optional asset to its binary representation.
func (a OptionalAsset) MarshalBinary() ([]byte, error) {
	buf := make([]byte, a.BinarySize())
	PutBool(buf, a.Present)
	if a.Present {
		copy(buf[1:], a.ID[:])
	}
	return buf, nil
}

// NewOptionalAssetFromString creates an OptionalAsset structure from its string representation.
func NewOptionalAssetFromString(s string) (*OptionalAsset, error) {
	switch strings.ToUpper(s) {
	case WavesAssetName, "":
		return &OptionalAsset{Present: false}, nil
	default:
		d, err := crypto.NewDigestFromBase58(s)
		if err != nil {
			return nil, err
		}
		return NewOptionalAssetFromDigest(d), nil
	}
}

// NewOptionalAssetFromBytes parses bytes as crypto.Digest and returns OptionalAsset.
func NewOptionalAssetFromBytes(b []byte) (*OptionalAsset, error) {
	d, err := crypto.NewDigestFromBytes(b)
	if err != nil {
		return nil, err
	}
	return NewOptionalAssetFromDigest(d), nil
}

func NewOptionalAsset(present bool, id crypto.Digest) OptionalAsset {
	return OptionalAsset{Present: present, ID: id}
}

func NewOptionalAssetFromDigest(d crypto.Digest) *OptionalAsset {
	return &OptionalAsset{Present: true, ID: d}
}

func NewOptionalAssetWaves() OptionalAsset {
	return OptionalAsset{Present: false}
}
