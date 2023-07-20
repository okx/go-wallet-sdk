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

	"github.com/okx/go-wallet-sdk/coins/solana/base"
)

// Like InitializeAccount, but the owner pubkey is passed via instruction data
// rather than the accounts list. This variant may be preferable when using
// Cross Program Invocation from an instruction that does not need the owner's
// `AccountInfo` otherwise.
type InitializeAccount2 struct {
	// The new account's owner/multisignature.
	Owner *base.PublicKey

	// [0] = [WRITE] account
	// ··········· The account to initialize.
	//
	// [1] = [] mint
	// ··········· The mint this account will be associated with.
	//
	// [2] = [] $(SysVarRentPubkey)
	// ··········· Rent sysvar.
	base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewInitializeAccount2InstructionBuilder creates a new `InitializeAccount2` instruction builder.
func NewInitializeAccount2InstructionBuilder() *InitializeAccount2 {
	nd := &InitializeAccount2{
		AccountMetaSlice: make(base.AccountMetaSlice, 3),
	}
	nd.AccountMetaSlice[2] = base.Meta(base.SysVarRentPubkey)
	return nd
}

// SetOwner sets the "owner" parameter.
// The new account's owner/multisignature.
func (inst *InitializeAccount2) SetOwner(owner base.PublicKey) *InitializeAccount2 {
	inst.Owner = &owner
	return inst
}

// SetAccount sets the "account" account.
// The account to initialize.
func (inst *InitializeAccount2) SetAccount(account base.PublicKey) *InitializeAccount2 {
	inst.AccountMetaSlice[0] = base.Meta(account).WRITE()
	return inst
}

// GetAccount gets the "account" account.
// The account to initialize.
func (inst *InitializeAccount2) GetAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[0]
}

// SetMintAccount sets the "mint" account.
// The mint this account will be associated with.
func (inst *InitializeAccount2) SetMintAccount(mint base.PublicKey) *InitializeAccount2 {
	inst.AccountMetaSlice[1] = base.Meta(mint)
	return inst
}

// GetMintAccount gets the "mint" account.
// The mint this account will be associated with.
func (inst *InitializeAccount2) GetMintAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[1]
}

// SetSysVarRentPubkeyAccount sets the "$(SysVarRentPubkey)" account.
// Rent sysvar.
func (inst *InitializeAccount2) SetSysVarRentPubkeyAccount(SysVarRentPubkey base.PublicKey) *InitializeAccount2 {
	inst.AccountMetaSlice[2] = base.Meta(SysVarRentPubkey)
	return inst
}

// GetSysVarRentPubkeyAccount gets the "$(SysVarRentPubkey)" account.
// Rent sysvar.
func (inst *InitializeAccount2) GetSysVarRentPubkeyAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[2]
}

func (inst InitializeAccount2) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint8(Instruction_InitializeAccount2),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst InitializeAccount2) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *InitializeAccount2) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Owner == nil {
			return errors.New("Owner parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Account is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Mint is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.SysVarRentPubkey is not set")
		}
	}
	return nil
}

func (obj InitializeAccount2) MarshalWithEncoder(encoder *base.Encoder) (err error) {
	// Serialize `Owner` param:
	err = encoder.Encode(obj.Owner)
	if err != nil {
		return err
	}
	return nil
}

// NewInitializeAccount2Instruction declares a new InitializeAccount2 instruction with the provided parameters and accounts.
func NewInitializeAccount2Instruction(
	// Parameters:
	owner base.PublicKey,
	// Accounts:
	account base.PublicKey,
	mint base.PublicKey) *InitializeAccount2 {
	return NewInitializeAccount2InstructionBuilder().
		SetOwner(owner).
		SetAccount(account).
		SetMintAccount(mint).
		SetSysVarRentPubkeyAccount(base.SysVarRentPubkey)
}
