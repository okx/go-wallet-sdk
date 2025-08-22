package types

import (
	"encoding/hex"
	"errors"
	"github.com/okx/go-wallet-sdk/coins/aptos/common"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/crypto"
	"github.com/okx/go-wallet-sdk/crypto/base58"
	"strings"
)

// Account represents an on-chain account, with an associated signer, which must be a [crypto.Signer]
//
// Implements:
//   - [crypto.Signer]
type Account struct {
	Address AccountAddress
	Signer  crypto.Signer
}

// NewAccountFromSigner creates an account from a [crypto.Signer] with an optional [crypto.AuthenticationKey]
func NewAccountFromSigner(signer crypto.Signer, authKey ...crypto.AuthenticationKey) (*Account, error) {
	out := &Account{}
	if len(authKey) == 1 {
		copy(out.Address[:], authKey[0][:])
	} else if len(authKey) > 1 {
		// Throw error
		return nil, errors.New("must only provide one auth key")
	} else {
		copy(out.Address[:], signer.AuthKey()[:])
	}
	out.Signer = signer
	return out, nil
}

// NewEd25519Account creates an account with a new random Ed25519 private key
func NewEd25519Account() (*Account, error) {
	privateKey, err := crypto.GenerateEd25519PrivateKey()
	if err != nil {
		return nil, err
	}
	return NewAccountFromSigner(privateKey)
}

// NewEd25519SingleSignerAccount creates a new random Ed25519 account
func NewEd25519SingleSignerAccount() (*Account, error) {
	privateKey, err := crypto.GenerateEd25519PrivateKey()
	if err != nil {
		return nil, err
	}
	signer := &crypto.SingleSigner{Signer: privateKey}
	return NewAccountFromSigner(signer)
}

// Sign signs a message, returning an appropriate authenticator for the signer
func (account *Account) Sign(message []byte) (authenticator *crypto.AccountAuthenticator, err error) {
	return account.Signer.Sign(message)
}

// SignMessage signs a message and returns the raw signature without a public key for verification
func (account *Account) SignMessage(message []byte) (signature crypto.Signature, err error) {
	return account.Signer.SignMessage(message)
}

// SimulationAuthenticator creates a new authenticator for simulation purposes
func (account *Account) SimulationAuthenticator() *crypto.AccountAuthenticator {
	return account.Signer.SimulationAuthenticator()
}

// PubKey retrieves the public key for signature verification
func (account *Account) PubKey() crypto.PublicKey {
	return account.Signer.PubKey()
}

// AuthKey retrieves the authentication key associated with the signer
func (account *Account) AuthKey() *crypto.AuthenticationKey {
	return account.Signer.AuthKey()
}

// AccountAddress retrieves the account address
func (account *Account) AccountAddress() AccountAddress {
	return account.Address
}

// ErrAddressTooShort is returned when an AccountAddress is too short
var ErrAddressTooShort = errors.New("AccountAddress too short")

// ErrAddressTooLong is returned when an AccountAddress is too long
var ErrAddressTooLong = errors.New("AccountAddress too long")

// ParseStringRelaxed parses a string into an AccountAddress
// TODO: add strict mode checking
func (aa *AccountAddress) ParseStringRelaxed(x string) error {
	if strings.HasPrefix(x, "0x") || common.IsHexString(x) {
		x = x[2:]
		if len(x) < 1 {
			return ErrAddressTooShort
		}
		if len(x) > 64 {
			return ErrAddressTooLong
		}
		if len(x)%2 != 0 {
			x = "0" + x
		}
		bytes, err := hex.DecodeString(x)
		if err != nil {
			return err
		}
		// zero-prefix/right-align what bytes we got
		copy((*aa)[32-len(bytes):], bytes)
	} else {
		if len(x) < 30 {
			return ErrAddressTooShort
		}

		base58Bytes := base58.Decode(x)
		copy((*aa)[:], base58Bytes)
	}
	return nil
}

// ParseStringWithPrefixRelaxed parses a string into an AccountAddress
func (aa *AccountAddress) ParseStringWithPrefixRelaxed(x string) error {
	if !strings.HasPrefix(x, "0x") {
		return ErrAddressTooShort
	}
	x = x[2:]
	if len(x) < 1 {
		return ErrAddressTooShort
	}
	if len(x) > 64 {
		return ErrAddressTooLong
	}
	if len(x)%2 != 0 {
		x = "0" + x
	}
	bytes, err := hex.DecodeString(x)
	if err != nil {
		return err
	}
	// zero-prefix/right-align what bytes we got
	copy((*aa)[32-len(bytes):], bytes)

	return nil
}
