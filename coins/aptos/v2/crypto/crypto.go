package crypto

import "github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"

// Signer a generic interface for any kind of signing
type Signer interface {
	// Sign signs a transaction and returns an associated [AccountAuthenticator]
	Sign(msg []byte) (authenticator *AccountAuthenticator, err error)

	// SignMessage signs a message and returns the raw [Signature] without a [PublicKey] for verification
	SignMessage(msg []byte) (signature Signature, err error)

	// SimulationAuthenticator creates a new [AccountAuthenticator] for simulation purposes
	SimulationAuthenticator() *AccountAuthenticator

	// AuthKey gives the [AuthenticationKey] associated with the [Signer]
	AuthKey() *AuthenticationKey

	// PubKey Retrieve the [PublicKey] for [Signature] verification
	PubKey() PublicKey
}

// MessageSigner a generic interface for a signing private key, a private key isn't always a signer, see SingleSender
//
// This is not BCS serializable, because this doesn't go on-chain.  An example is [Secp256k1PrivateKey]
type MessageSigner interface {
	// SignMessage signs a message and returns the raw [Signature] without a [VerifyingKey]
	SignMessage(msg []byte) (signature Signature, err error)

	// EmptySignature creates an empty signature for use in simulation
	EmptySignature() Signature

	// VerifyingKey Retrieve the [VerifyingKey] for signature verification.
	VerifyingKey() VerifyingKey
}

// PublicKey is an interface for a public key that can be used to verify transactions in a TransactionAuthenticator
type PublicKey interface {
	VerifyingKey

	// AuthKey gives the [AuthenticationKey] associated with the [PublicKey]
	AuthKey() *AuthenticationKey

	// Scheme The [DeriveScheme] used for address derivation
	Scheme() DeriveScheme
}

// VerifyingKey a generic interface for a public key associated with the private key, but it cannot necessarily stand on
// its own as a [PublicKey] for authentication on Aptos.  An example is [Secp256k1PublicKey].  All [PublicKey]s are also
// VerifyingKeys.
type VerifyingKey interface {
	bcs.Struct
	CryptoMaterial

	// Verify verifies a message with the public key
	Verify(msg []byte, sig Signature) bool
}

// Signature is an identifier for a serializable [Signature] for on-chain representation
type Signature interface {
	bcs.Struct
	CryptoMaterial
}

// CryptoMaterial is a set of functions for serializing and deserializing a key to and from bytes and hex
// This mirrors the trait in Rust
type CryptoMaterial interface {
	// Bytes outputs the raw byte representation of the [CryptoMaterial]
	Bytes() []byte

	// FromBytes loads the [CryptoMaterial] from the raw bytes
	FromBytes([]byte) error

	// ToHex outputs the hex representation of the [CryptoMaterial] with a leading `0x`
	ToHex() string

	// FromHex parses the hex representation of the [CryptoMaterial] with or without a leading `0x`
	FromHex(string) error
}
