package hashes

import (
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/model/externalapi"
	"github.com/pkg/errors"
	"golang.org/x/crypto/sha3"
	"hash"
)

// HashWriter is used to incrementally hash data without concatenating all of the data to a single buffer
// it exposes an io.Writer api and a Finalize function to get the resulting hash.
// The used hash function is blake2b.
// This can only be created via one of the domain separated constructors
type HashWriter struct {
	hash.Hash
}

// InfallibleWrite is just like write but doesn't return anything
func (h HashWriter) InfallibleWrite(p []byte) {
	// This write can never return an error, this is part of the hash.Hash interface contract.
	_, err := h.Write(p)
	if err != nil {
		panic(errors.Wrap(err, "this should never happen. hash.Hash interface promises to not return errors."))
	}
}

// Finalize returns the resulting hash
func (h HashWriter) Finalize() *externalapi.DomainHash {
	var sum [externalapi.DomainHashSize]byte
	// This should prevent `Sum` for allocating an output buffer, by using the DomainHash buffer. we still copy because we don't want to rely on that.
	copy(sum[:], h.Sum(sum[:0]))
	return externalapi.NewDomainHashFromByteArray(&sum)
}

// ShakeHashWriter is exactly the same as HashWriter but for CShake256
type ShakeHashWriter struct {
	sha3.ShakeHash
}

// InfallibleWrite is just like write but doesn't return anything
func (h *ShakeHashWriter) InfallibleWrite(p []byte) {
	// This write can never return an error, this is part of the hash.Hash interface contract.
	_, err := h.Write(p)
	if err != nil {
		panic(errors.Wrap(err, "this should never happen. sha3.ShakeHash interface promises to not return errors."))
	}
}

// Finalize returns the resulting hash
func (h *ShakeHashWriter) Finalize() *externalapi.DomainHash {
	var sum [externalapi.DomainHashSize]byte
	// This should prevent `Sum` for allocating an output buffer, by using the DomainHash buffer. we still copy because we don't want to rely on that.
	_, err := h.Read(sum[:])
	if err != nil {
		panic(errors.Wrap(err, "this should never happen. sha3.ShakeHash interface promises to not return errors."))
	}
	h.ShakeHash = nil // prevent double reading as it will return a different hash
	return externalapi.NewDomainHashFromByteArray(&sum)
}
