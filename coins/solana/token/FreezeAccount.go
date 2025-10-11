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

// Freeze an Initialized account using the Mint's freeze_authority (if set).
type FreezeAccount struct {

	// [0] = [WRITE] account
	// ··········· The account to freeze.
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

func (obj *FreezeAccount) SetAccounts(accounts []*base.AccountMeta) error {
	obj.Accounts, obj.Signers = base.AccountMetaSlice(accounts).SplitFrom(3)
	return nil
}

func (slice FreezeAccount) GetAccounts() (accounts []*base.AccountMeta) {
	accounts = append(accounts, slice.Accounts...)
	accounts = append(accounts, slice.Signers...)
	return
}

// NewFreezeAccountInstructionBuilder creates a new `FreezeAccount` instruction builder.
func NewFreezeAccountInstructionBuilder() *FreezeAccount {
	nd := &FreezeAccount{
		Accounts: make(base.AccountMetaSlice, 3),
		Signers:  make(base.AccountMetaSlice, 0),
	}
	return nd
}

// SetAccount sets the "account" account.
// The account to freeze.
func (inst *FreezeAccount) SetAccount(account base.PublicKey) *FreezeAccount {
	inst.Accounts[0] = base.Meta(account).WRITE()
	return inst
}

// GetAccount gets the "account" account.
// The account to freeze.
func (inst *FreezeAccount) GetAccount() *base.AccountMeta {
	return inst.Accounts[0]
}

// SetMintAccount sets the "mint" account.
// The token mint.
func (inst *FreezeAccount) SetMintAccount(mint base.PublicKey) *FreezeAccount {
	inst.Accounts[1] = base.Meta(mint)
	return inst
}

// GetMintAccount gets the "mint" account.
// The token mint.
func (inst *FreezeAccount) GetMintAccount() *base.AccountMeta {
	return inst.Accounts[1]
}

// SetAuthorityAccount sets the "authority" account.
// The mint freeze authority.
func (inst *FreezeAccount) SetAuthorityAccount(authority base.PublicKey, multisigSigners ...base.PublicKey) *FreezeAccount {
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
func (inst *FreezeAccount) GetAuthorityAccount() *base.AccountMeta {
	return inst.Accounts[2]
}

func (inst FreezeAccount) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint8(Instruction_FreezeAccount),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst FreezeAccount) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *FreezeAccount) Validate() error {
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

func (obj FreezeAccount) MarshalWithEncoder(encoder *base.Encoder) (err error) {
	return nil
}

// NewFreezeAccountInstruction declares a new FreezeAccount instruction with the provided parameters and accounts.
func NewFreezeAccountInstruction(
	// Accounts:
	account base.PublicKey,
	mint base.PublicKey,
	authority base.PublicKey,
	multisigSigners []base.PublicKey,
) *FreezeAccount {
	return NewFreezeAccountInstructionBuilder().
		SetAccount(account).
		SetMintAccount(mint).
		SetAuthorityAccount(authority, multisigSigners...)
}
