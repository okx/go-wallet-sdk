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

package token

import (
	"errors"
	"fmt"

	"github.com/emresenyuva/go-wallet-sdk/coins/solana/base"
)

// Sets a new authority of a mint or account.
type SetAuthority struct {
	// The type of authority to update.
	AuthorityType *AuthorityType

	// The new authority.
	NewAuthority *base.PublicKey `bin:"optional"`

	// [0] = [WRITE] subject
	// ··········· The mint or account to change the authority of.
	//
	// [1] = [] authority
	// ··········· The current authority of the mint or account.
	//
	// [2...] = [SIGNER] signers
	// ··········· M signer accounts.
	Accounts base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
	Signers  base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func (obj *SetAuthority) SetAccounts(accounts []*base.AccountMeta) error {
	obj.Accounts, obj.Signers = base.AccountMetaSlice(accounts).SplitFrom(2)
	return nil
}

func (slice SetAuthority) GetAccounts() (accounts []*base.AccountMeta) {
	accounts = append(accounts, slice.Accounts...)
	accounts = append(accounts, slice.Signers...)
	return
}

// NewSetAuthorityInstructionBuilder creates a new `SetAuthority` instruction builder.
func NewSetAuthorityInstructionBuilder() *SetAuthority {
	nd := &SetAuthority{
		Accounts: make(base.AccountMetaSlice, 2),
		Signers:  make(base.AccountMetaSlice, 0),
	}
	return nd
}

// SetAuthorityType sets the "authority_type" parameter.
// The type of authority to update.
func (inst *SetAuthority) SetAuthorityType(authority_type AuthorityType) *SetAuthority {
	inst.AuthorityType = &authority_type
	return inst
}

// SetNewAuthority sets the "new_authority" parameter.
// The new authority.
func (inst *SetAuthority) SetNewAuthority(new_authority base.PublicKey) *SetAuthority {
	inst.NewAuthority = &new_authority
	return inst
}

// SetSubjectAccount sets the "subject" account.
// The mint or account to change the authority of.
func (inst *SetAuthority) SetSubjectAccount(subject base.PublicKey) *SetAuthority {
	inst.Accounts[0] = base.Meta(subject).WRITE()
	return inst
}

// GetSubjectAccount gets the "subject" account.
// The mint or account to change the authority of.
func (inst *SetAuthority) GetSubjectAccount() *base.AccountMeta {
	return inst.Accounts[0]
}

// SetAuthorityAccount sets the "authority" account.
// The current authority of the mint or account.
func (inst *SetAuthority) SetAuthorityAccount(authority base.PublicKey, multisigSigners ...base.PublicKey) *SetAuthority {
	inst.Accounts[1] = base.Meta(authority)
	if len(multisigSigners) == 0 {
		inst.Accounts[1].SIGNER()
	}
	for _, signer := range multisigSigners {
		inst.Signers = append(inst.Signers, base.Meta(signer).SIGNER())
	}
	return inst
}

// GetAuthorityAccount gets the "authority" account.
// The current authority of the mint or account.
func (inst *SetAuthority) GetAuthorityAccount() *base.AccountMeta {
	return inst.Accounts[1]
}

func (inst SetAuthority) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint8(Instruction_SetAuthority),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst SetAuthority) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *SetAuthority) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.AuthorityType == nil {
			return errors.New("AuthorityType parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.Accounts[0] == nil {
			return errors.New("accounts.Subject is not set")
		}
		if inst.Accounts[1] == nil {
			return errors.New("accounts.Authority is not set")
		}
		if !inst.Accounts[1].IsSigner && len(inst.Signers) == 0 {
			return fmt.Errorf("accounts.Signers is not set")
		}
		if len(inst.Signers) > MAX_SIGNERS {
			return fmt.Errorf("too many signers; got %v, but max is 11", len(inst.Signers))
		}
	}
	return nil
}

func (obj SetAuthority) MarshalWithEncoder(encoder *base.Encoder) (err error) {
	// Serialize `AuthorityType` param:
	err = encoder.Encode(obj.AuthorityType)
	if err != nil {
		return err
	}
	// Serialize `NewAuthority` param (optional):
	{
		if obj.NewAuthority == nil {
			err = encoder.WriteBool(false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteBool(true)
			if err != nil {
				return err
			}
			err = encoder.Encode(obj.NewAuthority)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// NewSetAuthorityInstruction declares a new SetAuthority instruction with the provided parameters and accounts.
func NewSetAuthorityInstruction(
	// Parameters:
	authority_type AuthorityType,
	new_authority base.PublicKey,
	// Accounts:
	subject base.PublicKey,
	authority base.PublicKey,
	multisigSigners []base.PublicKey,
) *SetAuthority {
	return NewSetAuthorityInstructionBuilder().
		SetAuthorityType(authority_type).
		SetNewAuthority(new_authority).
		SetSubjectAccount(subject).
		SetAuthorityAccount(authority, multisigSigners...)
}
