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

// Burns tokens by removing them from an account.  `Burn` does not support
// accounts associated with the native mint, use `CloseAccount` instead.
type Burn struct {
	// The amount of tokens to burn.
	Amount *uint64

	// [0] = [WRITE] source
	// ··········· The account to burn from.
	//
	// [1] = [WRITE] mint
	// ··········· The token mint.
	//
	// [2] = [] owner
	// ··········· The account's owner/delegate.
	//
	// [3...] = [SIGNER] signers
	// ··········· M signer accounts.
	Accounts base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
	Signers  base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func (obj *Burn) SetAccounts(accounts []*base.AccountMeta) error {
	obj.Accounts, obj.Signers = base.AccountMetaSlice(accounts).SplitFrom(3)
	return nil
}

func (slice Burn) GetAccounts() (accounts []*base.AccountMeta) {
	accounts = append(accounts, slice.Accounts...)
	accounts = append(accounts, slice.Signers...)
	return
}

// NewBurnInstructionBuilder creates a new `Burn` instruction builder.
func NewBurnInstructionBuilder() *Burn {
	nd := &Burn{
		Accounts: make(base.AccountMetaSlice, 3),
		Signers:  make(base.AccountMetaSlice, 0),
	}
	return nd
}

// SetAmount sets the "amount" parameter.
// The amount of tokens to burn.
func (inst *Burn) SetAmount(amount uint64) *Burn {
	inst.Amount = &amount
	return inst
}

// SetSourceAccount sets the "source" account.
// The account to burn from.
func (inst *Burn) SetSourceAccount(source base.PublicKey) *Burn {
	inst.Accounts[0] = base.Meta(source).WRITE()
	return inst
}

// GetSourceAccount gets the "source" account.
// The account to burn from.
func (inst *Burn) GetSourceAccount() *base.AccountMeta {
	return inst.Accounts[0]
}

// SetMintAccount sets the "mint" account.
// The token mint.
func (inst *Burn) SetMintAccount(mint base.PublicKey) *Burn {
	inst.Accounts[1] = base.Meta(mint).WRITE()
	return inst
}

// GetMintAccount gets the "mint" account.
// The token mint.
func (inst *Burn) GetMintAccount() *base.AccountMeta {
	return inst.Accounts[1]
}

// SetOwnerAccount sets the "owner" account.
// The account's owner/delegate.
func (inst *Burn) SetOwnerAccount(owner base.PublicKey, multisigSigners ...base.PublicKey) *Burn {
	inst.Accounts[2] = base.Meta(owner)
	if len(multisigSigners) == 0 {
		inst.Accounts[2].SIGNER()
	}
	for _, signer := range multisigSigners {
		inst.Signers = append(inst.Signers, base.Meta(signer).SIGNER())
	}
	return inst
}

// GetOwnerAccount gets the "owner" account.
// The account's owner/delegate.
func (inst *Burn) GetOwnerAccount() *base.AccountMeta {
	return inst.Accounts[2]
}

func (inst Burn) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint8(Instruction_Burn),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst Burn) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *Burn) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Amount == nil {
			return errors.New("Amount parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.Accounts[0] == nil {
			return errors.New("accounts.Source is not set")
		}
		if inst.Accounts[1] == nil {
			return errors.New("accounts.Mint is not set")
		}
		if inst.Accounts[2] == nil {
			return errors.New("accounts.Owner is not set")
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

func (obj Burn) MarshalWithEncoder(encoder *base.Encoder) (err error) {
	// Serialize `Amount` param:
	err = encoder.Encode(obj.Amount)
	if err != nil {
		return err
	}
	return nil
}

// NewBurnInstruction declares a new Burn instruction with the provided parameters and accounts.
func NewBurnInstruction(
	// Parameters:
	amount uint64,
	// Accounts:
	source base.PublicKey,
	mint base.PublicKey,
	owner base.PublicKey,
	multisigSigners []base.PublicKey,
) *Burn {
	return NewBurnInstructionBuilder().
		SetAmount(amount).
		SetSourceAccount(source).
		SetMintAccount(mint).
		SetOwnerAccount(owner, multisigSigners...)
}
