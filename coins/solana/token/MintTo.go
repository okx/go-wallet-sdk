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

// Mints new tokens to an account.  The native mint does not support
// minting.
type MintTo struct {
	// The amount of new tokens to mint.
	Amount *uint64

	// [0] = [WRITE] mint
	// ··········· The mint.
	//
	// [1] = [WRITE] destination
	// ··········· The account to mint tokens to.
	//
	// [2] = [] authority
	// ··········· The mint's minting authority.
	//
	// [3...] = [SIGNER] signers
	// ··········· M signer accounts.
	Accounts base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
	Signers  base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func (obj *MintTo) SetAccounts(accounts []*base.AccountMeta) error {
	obj.Accounts, obj.Signers = base.AccountMetaSlice(accounts).SplitFrom(3)
	return nil
}

func (slice MintTo) GetAccounts() (accounts []*base.AccountMeta) {
	accounts = append(accounts, slice.Accounts...)
	accounts = append(accounts, slice.Signers...)
	return
}

// NewMintToInstructionBuilder creates a new `MintTo` instruction builder.
func NewMintToInstructionBuilder() *MintTo {
	nd := &MintTo{
		Accounts: make(base.AccountMetaSlice, 3),
		Signers:  make(base.AccountMetaSlice, 0),
	}
	return nd
}

// SetAmount sets the "amount" parameter.
// The amount of new tokens to mint.
func (inst *MintTo) SetAmount(amount uint64) *MintTo {
	inst.Amount = &amount
	return inst
}

// SetMintAccount sets the "mint" account.
// The mint.
func (inst *MintTo) SetMintAccount(mint base.PublicKey) *MintTo {
	inst.Accounts[0] = base.Meta(mint).WRITE()
	return inst
}

// GetMintAccount gets the "mint" account.
// The mint.
func (inst *MintTo) GetMintAccount() *base.AccountMeta {
	return inst.Accounts[0]
}

// SetDestinationAccount sets the "destination" account.
// The account to mint tokens to.
func (inst *MintTo) SetDestinationAccount(destination base.PublicKey) *MintTo {
	inst.Accounts[1] = base.Meta(destination).WRITE()
	return inst
}

// GetDestinationAccount gets the "destination" account.
// The account to mint tokens to.
func (inst *MintTo) GetDestinationAccount() *base.AccountMeta {
	return inst.Accounts[1]
}

// SetAuthorityAccount sets the "authority" account.
// The mint's minting authority.
func (inst *MintTo) SetAuthorityAccount(authority base.PublicKey, multisigSigners ...base.PublicKey) *MintTo {
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
// The mint's minting authority.
func (inst *MintTo) GetAuthorityAccount() *base.AccountMeta {
	return inst.Accounts[2]
}

func (inst MintTo) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint8(Instruction_MintTo),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst MintTo) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *MintTo) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Amount == nil {
			return errors.New("Amount parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.Accounts[0] == nil {
			return errors.New("accounts.Mint is not set")
		}
		if inst.Accounts[1] == nil {
			return errors.New("accounts.Destination is not set")
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

func (obj MintTo) MarshalWithEncoder(encoder *base.Encoder) (err error) {
	// Serialize `Amount` param:
	err = encoder.Encode(obj.Amount)
	if err != nil {
		return err
	}
	return nil
}

// NewMintToInstruction declares a new MintTo instruction with the provided parameters and accounts.
func NewMintToInstruction(
	// Parameters:
	amount uint64,
	// Accounts:
	mint base.PublicKey,
	destination base.PublicKey,
	authority base.PublicKey,
	multisigSigners []base.PublicKey,
) *MintTo {
	return NewMintToInstructionBuilder().
		SetAmount(amount).
		SetMintAccount(mint).
		SetDestinationAccount(destination).
		SetAuthorityAccount(authority, multisigSigners...)
}
