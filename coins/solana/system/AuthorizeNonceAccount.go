// Copyright 2021 github.com/gagliardetto
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package system

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/okx/go-wallet-sdk/coins/solana/base"
)

// Change the entity authorized to execute nonce instructions on the account
type AuthorizeNonceAccount struct {
	// The Pubkey parameter identifies the entity to authorize.
	Authorized *base.PublicKey

	// [0] = [WRITE] NonceAccount
	// ··········· Nonce account
	//
	// [1] = [SIGNER] NonceAuthorityAccount
	// ··········· Nonce authority
	base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewAuthorizeNonceAccountInstructionBuilder creates a new `AuthorizeNonceAccount` instruction builder.
func NewAuthorizeNonceAccountInstructionBuilder() *AuthorizeNonceAccount {
	nd := &AuthorizeNonceAccount{
		AccountMetaSlice: make(base.AccountMetaSlice, 2),
	}
	return nd
}

// The Pubkey parameter identifies the entity to authorize.
func (inst *AuthorizeNonceAccount) SetAuthorized(authorized base.PublicKey) *AuthorizeNonceAccount {
	inst.Authorized = &authorized
	return inst
}

// Nonce account
func (inst *AuthorizeNonceAccount) SetNonceAccount(nonceAccount base.PublicKey) *AuthorizeNonceAccount {
	inst.AccountMetaSlice[0] = base.Meta(nonceAccount).WRITE()
	return inst
}

func (inst *AuthorizeNonceAccount) GetNonceAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[0]
}

// Nonce authority
func (inst *AuthorizeNonceAccount) SetNonceAuthorityAccount(nonceAuthorityAccount base.PublicKey) *AuthorizeNonceAccount {
	inst.AccountMetaSlice[1] = base.Meta(nonceAuthorityAccount).SIGNER()
	return inst
}

func (inst *AuthorizeNonceAccount) GetNonceAuthorityAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[1]
}

func (inst AuthorizeNonceAccount) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint32(Instruction_AuthorizeNonceAccount, binary.LittleEndian),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst AuthorizeNonceAccount) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *AuthorizeNonceAccount) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Authorized == nil {
			return errors.New("Authorized parameter is not set")
		}
	}

	// Check whether all accounts are set:
	for accIndex, acc := range inst.AccountMetaSlice {
		if acc == nil {
			return fmt.Errorf("ins.AccountMetaSlice[%v] is not set", accIndex)
		}
	}
	return nil
}

func (inst AuthorizeNonceAccount) MarshalWithEncoder(encoder *base.Encoder) error {
	// Serialize `Authorized` param:
	{
		err := encoder.Encode(*inst.Authorized)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewAuthorizeNonceAccountInstruction declares a new AuthorizeNonceAccount instruction with the provided parameters and accounts.
func NewAuthorizeNonceAccountInstruction(
	// Parameters:
	authorized base.PublicKey,
	// Accounts:
	nonceAccount base.PublicKey,
	nonceAuthorityAccount base.PublicKey) *AuthorizeNonceAccount {
	return NewAuthorizeNonceAccountInstructionBuilder().
		SetAuthorized(authorized).
		SetNonceAccount(nonceAccount).
		SetNonceAuthorityAccount(nonceAuthorityAccount)
}
