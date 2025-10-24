package crypto

import "github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"

// region NoAuthenticator

// NoAuthenticator represents a verifiable signature with it's accompanied public key
//
// Implements:
//
//   - [AccountAuthenticatorImpl]
//   - [bcs.Marshaler]
//   - [bcs.Unmarshaler]
//   - [bcs.Struct]
type NoAuthenticator struct{}

// region NoAuthenticator AccountAuthenticatorImpl implementation

// PublicKey returns the [PublicKey] of the authenticator
//
// Implements:
//   - [AccountAuthenticatorImpl]
func (ea *NoAuthenticator) PublicKey() PublicKey {
	return nil
}

// Signature returns the [Signature] of the authenticator
//
// Implements:
//   - [AccountAuthenticatorImpl]
func (ea *NoAuthenticator) Signature() Signature {
	return nil
}

// Verify returns true if the authenticator can be cryptographically verified
//
// Implements:
//   - [AccountAuthenticatorImpl]
func (ea *NoAuthenticator) Verify([]byte) bool {
	return false
}

// endregion

// region NoAuthenticator bcs.Struct implementation

// MarshalBCS serializes the [NoAuthenticator] to BCS bytes
//
// Implements:
//   - [bcs.Marshaler]
func (ea *NoAuthenticator) MarshalBCS(*bcs.Serializer) {
	// TODO: Double check nothing is needed here
}

// UnmarshalBCS deserializes the [NoAuthenticator] from BCS bytes
//
// Sets [bcs.Deserializer.Error] if it fails to read the required bytes.
//
// Implements:
//   - [bcs.Unmarshaler]
func (ea *NoAuthenticator) UnmarshalBCS(*bcs.Deserializer) {
	// TODO: Double check nothing is needed here
}

func NoAccountAuthenticator() *AccountAuthenticator {
	return &AccountAuthenticator{
		Variant: AccountAuthenticatorNone,
		Auth:    &NoAuthenticator{},
	}
}

// endregion
// endregion
