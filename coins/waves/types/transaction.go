package types

import "github.com/okx/go-wallet-sdk/coins/waves/crypto"

// Transaction is a set of common transaction functions.
type Transaction interface {
	// Sign transaction with given secret key.
	// It also sets transaction ID.
	Sign(scheme Scheme, sk crypto.SecretKey) error

	// MarshalBinary functions for custom binary format serialization.
	// MarshalBinary() is analogous to MarshalSignedToProtobuf() for Protobuf.
	MarshalBinary() ([]byte, error)

	// BodyMarshalBinary is analogous to MarshalToProtobuf() for Protobuf.
	BodyMarshalBinary() ([]byte, error)
}
