package externalapi

import (
	"bytes"
	"encoding/hex"

	"github.com/pkg/errors"
)

// DomainHashSize of array used to store hashes.
const DomainHashSize = 32

// DomainHash is the domain representation of a Hash
type DomainHash struct {
	hashArray [DomainHashSize]byte
}

// NewZeroHash returns a DomainHash that represents the zero value (0x000000...000)
func NewZeroHash() *DomainHash {
	return &DomainHash{hashArray: [32]byte{}}
}

// NewDomainHashFromByteArray constructs a new DomainHash out of a byte array
func NewDomainHashFromByteArray(hashBytes *[DomainHashSize]byte) *DomainHash {
	return &DomainHash{
		hashArray: *hashBytes,
	}
}

// NewDomainHashFromByteSlice constructs a new DomainHash out of a byte slice.
// Returns an error if the length of the byte slice is not exactly `DomainHashSize`
func NewDomainHashFromByteSlice(hashBytes []byte) (*DomainHash, error) {
	if len(hashBytes) != DomainHashSize {
		return nil, errors.Errorf("invalid hash size. Want: %d, got: %d",
			DomainHashSize, len(hashBytes))
	}
	domainHash := DomainHash{
		hashArray: [DomainHashSize]byte{},
	}
	copy(domainHash.hashArray[:], hashBytes)
	return &domainHash, nil
}

// NewDomainHashFromString constructs a new DomainHash out of a hex-encoded string.
// Returns an error if the length of the string is not exactly `DomainHashSize * 2`
func NewDomainHashFromString(hashString string) (*DomainHash, error) {
	expectedLength := DomainHashSize * 2
	// Return error if hash string is too long.
	if len(hashString) != expectedLength {
		return nil, errors.Errorf("hash string length is %d, while it should be be %d",
			len(hashString), expectedLength)
	}

	hashBytes, err := hex.DecodeString(hashString)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return NewDomainHashFromByteSlice(hashBytes)
}

// String returns the Hash as the hexadecimal string of the hash.
func (hash DomainHash) String() string {
	return hex.EncodeToString(hash.hashArray[:])
}

// ByteArray returns the bytes in this hash represented as a byte array.
// The hash bytes are cloned, therefore it is safe to modify the resulting array.
func (hash *DomainHash) ByteArray() *[DomainHashSize]byte {
	arrayClone := hash.hashArray
	return &arrayClone
}

// ByteSlice returns the bytes in this hash represented as a byte slice.
// The hash bytes are cloned, therefore it is safe to modify the resulting slice.
func (hash *DomainHash) ByteSlice() []byte {
	return hash.ByteArray()[:]
}

// If this doesn't compile, it means the type definition has been changed, so it's
// an indication to update Equal and Clone accordingly.
var _ DomainHash = DomainHash{hashArray: [DomainHashSize]byte{}}

// Equal returns whether hash equals to other
func (hash *DomainHash) Equal(other *DomainHash) bool {
	if hash == nil || other == nil {
		return hash == other
	}

	return hash.hashArray == other.hashArray
}

// Less returns true if hash is less than other
func (hash *DomainHash) Less(other *DomainHash) bool {
	return bytes.Compare(hash.hashArray[:], other.hashArray[:]) < 0
}

// LessOrEqual returns true if hash is smaller or equal to other
func (hash *DomainHash) LessOrEqual(other *DomainHash) bool {
	return bytes.Compare(hash.hashArray[:], other.hashArray[:]) <= 0
}

// CloneHashes returns a clone of the given hashes slice.
// Note: since DomainHash is a read-only type, the clone is shallow
func CloneHashes(hashes []*DomainHash) []*DomainHash {
	clone := make([]*DomainHash, len(hashes))
	copy(clone, hashes)
	return clone
}

// HashesEqual returns whether the given hash slices are equal.
func HashesEqual(a, b []*DomainHash) bool {
	if len(a) != len(b) {
		return false
	}

	for i, hash := range a {
		if !hash.Equal(b[i]) {
			return false
		}
	}
	return true
}
