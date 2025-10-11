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

// Burns tokens by removing them from an account.  `BurnChecked` does not
// support accounts associated with the native mint, use `CloseAccount`
// instead.
//
// This instruction differs from Burn in that the decimals value is checked
// by the caller. This may be useful when creating transactions offline or
// within a hardware wallet.
type BurnChecked struct {
	// The amount of tokens to burn.
	Amount *uint64

	// Expected number of base 10 digits to the right of the decimal place.
	Decimals *uint8

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

func (obj *BurnChecked) SetAccounts(accounts []*base.AccountMeta) error {
	obj.Accounts, obj.Signers = base.AccountMetaSlice(accounts).SplitFrom(3)
	return nil
}

func (slice BurnChecked) GetAccounts() (accounts []*base.AccountMeta) {
	accounts = append(accounts, slice.Accounts...)
	accounts = append(accounts, slice.Signers...)
	return
}

// NewBurnCheckedInstructionBuilder creates a new `BurnChecked` instruction builder.
func NewBurnCheckedInstructionBuilder() *BurnChecked {
	nd := &BurnChecked{
		Accounts: make(base.AccountMetaSlice, 3),
		Signers:  make(base.AccountMetaSlice, 0),
	}
	return nd
}

// SetAmount sets the "amount" parameter.
// The amount of tokens to burn.
func (inst *BurnChecked) SetAmount(amount uint64) *BurnChecked {
	inst.Amount = &amount
	return inst
}

// SetDecimals sets the "decimals" parameter.
// Expected number of base 10 digits to the right of the decimal place.
func (inst *BurnChecked) SetDecimals(decimals uint8) *BurnChecked {
	inst.Decimals = &decimals
	return inst
}

// SetSourceAccount sets the "source" account.
// The account to burn from.
func (inst *BurnChecked) SetSourceAccount(source base.PublicKey) *BurnChecked {
	inst.Accounts[0] = base.Meta(source).WRITE()
	return inst
}

// GetSourceAccount gets the "source" account.
// The account to burn from.
func (inst *BurnChecked) GetSourceAccount() *base.AccountMeta {
	return inst.Accounts[0]
}

// SetMintAccount sets the "mint" account.
// The token mint.
func (inst *BurnChecked) SetMintAccount(mint base.PublicKey) *BurnChecked {
	inst.Accounts[1] = base.Meta(mint).WRITE()
	return inst
}

// GetMintAccount gets the "mint" account.
// The token mint.
func (inst *BurnChecked) GetMintAccount() *base.AccountMeta {
	return inst.Accounts[1]
}

// SetOwnerAccount sets the "owner" account.
// The account's owner/delegate.
func (inst *BurnChecked) SetOwnerAccount(owner base.PublicKey, multisigSigners ...base.PublicKey) *BurnChecked {
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
func (inst *BurnChecked) GetOwnerAccount() *base.AccountMeta {
	return inst.Accounts[2]
}

func (inst BurnChecked) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint8(Instruction_BurnChecked),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst BurnChecked) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *BurnChecked) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Amount == nil {
			return errors.New("Amount parameter is not set")
		}
		if inst.Decimals == nil {
			return errors.New("Decimals parameter is not set")
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

func (obj BurnChecked) MarshalWithEncoder(encoder *base.Encoder) (err error) {
	// Serialize `Amount` param:
	err = encoder.Encode(obj.Amount)
	if err != nil {
		return err
	}
	// Serialize `Decimals` param:
	err = encoder.Encode(obj.Decimals)
	if err != nil {
		return err
	}
	return nil
}

// NewBurnCheckedInstruction declares a new BurnChecked instruction with the provided parameters and accounts.
func NewBurnCheckedInstruction(
	// Parameters:
	amount uint64,
	decimals uint8,
	// Accounts:
	source base.PublicKey,
	mint base.PublicKey,
	owner base.PublicKey,
	multisigSigners []base.PublicKey,
) *BurnChecked {
	return NewBurnCheckedInstructionBuilder().
		SetAmount(amount).
		SetDecimals(decimals).
		SetSourceAccount(source).
		SetMintAccount(mint).
		SetOwnerAccount(owner, multisigSigners...)
}
