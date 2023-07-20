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

	"github.com/okx/go-wallet-sdk/coins/solana/base"
)

// Thaw a Frozen account using the Mint's freeze_authority (if set).
type ThawAccount struct {

	// [0] = [WRITE] account
	// ··········· The account to thaw.
	//
	// [1] = [] mint
	// ··········· The token mint.
	//
	// [2] = [] authority
	// ··········· The mint freeze authority.
	//
	// [3...] = [SIGNER] signers
	// ··········· M signer accounts.
	Accounts base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
	Signers  base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func (obj *ThawAccount) SetAccounts(accounts []*base.AccountMeta) error {
	obj.Accounts, obj.Signers = base.AccountMetaSlice(accounts).SplitFrom(3)
	return nil
}

func (slice ThawAccount) GetAccounts() (accounts []*base.AccountMeta) {
	accounts = append(accounts, slice.Accounts...)
	accounts = append(accounts, slice.Signers...)
	return
}

// NewThawAccountInstructionBuilder creates a new `ThawAccount` instruction builder.
func NewThawAccountInstructionBuilder() *ThawAccount {
	nd := &ThawAccount{
		Accounts: make(base.AccountMetaSlice, 3),
		Signers:  make(base.AccountMetaSlice, 0),
	}
	return nd
}

// SetAccount sets the "account" account.
// The account to thaw.
func (inst *ThawAccount) SetAccount(account base.PublicKey) *ThawAccount {
	inst.Accounts[0] = base.Meta(account).WRITE()
	return inst
}

// GetAccount gets the "account" account.
// The account to thaw.
func (inst *ThawAccount) GetAccount() *base.AccountMeta {
	return inst.Accounts[0]
}

// SetMintAccount sets the "mint" account.
// The token mint.
func (inst *ThawAccount) SetMintAccount(mint base.PublicKey) *ThawAccount {
	inst.Accounts[1] = base.Meta(mint)
	return inst
}

// GetMintAccount gets the "mint" account.
// The token mint.
func (inst *ThawAccount) GetMintAccount() *base.AccountMeta {
	return inst.Accounts[1]
}

// SetAuthorityAccount sets the "authority" account.
// The mint freeze authority.
func (inst *ThawAccount) SetAuthorityAccount(authority base.PublicKey, multisigSigners ...base.PublicKey) *ThawAccount {
	inst.Accounts[2] = base.Meta(authority)
	if len(multisigSigners) == 0 {
		inst.Accounts[2].SIGNER()
	}
	for _, signer := range multisigSigners {
		inst.Signers = append(inst.Signers, base.Meta(signer).SIGNER())
	}
	return inst
}

// GetAuthorityAccount gets the "authority" account.
// The mint freeze authority.
func (inst *ThawAccount) GetAuthorityAccount() *base.AccountMeta {
	return inst.Accounts[2]
}

func (inst ThawAccount) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint8(Instruction_ThawAccount),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst ThawAccount) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *ThawAccount) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if inst.Accounts[0] == nil {
			return errors.New("accounts.Account is not set")
		}
		if inst.Accounts[1] == nil {
			return errors.New("accounts.Mint is not set")
		}
		if inst.Accounts[2] == nil {
			return errors.New("accounts.Authority is not set")
		}
		if !inst.Accounts[2].IsSigner && len(inst.Signers) == 0 {
			return fmt.Errorf("accounts.Signers is not set")
		}
		if len(inst.Signers) > MAX_SIGNERS {
			return fmt.Errorf("too many signers; got %v, but max is 11", len(inst.Signers))
		}
	}
	return nil
}

func (obj ThawAccount) MarshalWithEncoder(encoder *base.Encoder) (err error) {
	return nil
}

// NewThawAccountInstruction declares a new ThawAccount instruction with the provided parameters and accounts.
func NewThawAccountInstruction(
	// Accounts:
	account base.PublicKey,
	mint base.PublicKey,
	authority base.PublicKey,
	multisigSigners []base.PublicKey,
) *ThawAccount {
	return NewThawAccountInstructionBuilder().
		SetAccount(account).
		SetMintAccount(mint).
		SetAuthorityAccount(authority, multisigSigners...)
}
