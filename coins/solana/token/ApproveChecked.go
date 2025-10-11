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

// Approves a delegate.  A delegate is given the authority over tokens on
// behalf of the source account's owner.
//
// This instruction differs from Approve in that the token mint and
// decimals value is checked by the caller.  This may be useful when
// creating transactions offline or within a hardware wallet.
type ApproveChecked struct {
	// The amount of tokens the delegate is approved for.
	Amount *uint64

	// Expected number of base 10 digits to the right of the decimal place.
	Decimals *uint8

	// [0] = [WRITE] source
	// ··········· The source account.
	//
	// [1] = [] mint
	// ··········· The token mint.
	//
	// [2] = [] delegate
	// ··········· The delegate.
	//
	// [3] = [] owner
	// ··········· The source account owner.
	//
	// [4...] = [SIGNER] signers
	// ··········· M signer accounts.
	Accounts base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
	Signers  base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func (obj *ApproveChecked) SetAccounts(accounts []*base.AccountMeta) error {
	obj.Accounts, obj.Signers = base.AccountMetaSlice(accounts).SplitFrom(4)
	return nil
}

func (slice ApproveChecked) GetAccounts() (accounts []*base.AccountMeta) {
	accounts = append(accounts, slice.Accounts...)
	accounts = append(accounts, slice.Signers...)
	return
}

// NewApproveCheckedInstructionBuilder creates a new `ApproveChecked` instruction builder.
func NewApproveCheckedInstructionBuilder() *ApproveChecked {
	nd := &ApproveChecked{
		Accounts: make(base.AccountMetaSlice, 4),
		Signers:  make(base.AccountMetaSlice, 0),
	}
	return nd
}

// SetAmount sets the "amount" parameter.
// The amount of tokens the delegate is approved for.
func (inst *ApproveChecked) SetAmount(amount uint64) *ApproveChecked {
	inst.Amount = &amount
	return inst
}

// SetDecimals sets the "decimals" parameter.
// Expected number of base 10 digits to the right of the decimal place.
func (inst *ApproveChecked) SetDecimals(decimals uint8) *ApproveChecked {
	inst.Decimals = &decimals
	return inst
}

// SetSourceAccount sets the "source" account.
// The source account.
func (inst *ApproveChecked) SetSourceAccount(source base.PublicKey) *ApproveChecked {
	inst.Accounts[0] = base.Meta(source).WRITE()
	return inst
}

// GetSourceAccount gets the "source" account.
// The source account.
func (inst *ApproveChecked) GetSourceAccount() *base.AccountMeta {
	return inst.Accounts[0]
}

// SetMintAccount sets the "mint" account.
// The token mint.
func (inst *ApproveChecked) SetMintAccount(mint base.PublicKey) *ApproveChecked {
	inst.Accounts[1] = base.Meta(mint)
	return inst
}

// GetMintAccount gets the "mint" account.
// The token mint.
func (inst *ApproveChecked) GetMintAccount() *base.AccountMeta {
	return inst.Accounts[1]
}

// SetDelegateAccount sets the "delegate" account.
// The delegate.
func (inst *ApproveChecked) SetDelegateAccount(delegate base.PublicKey) *ApproveChecked {
	inst.Accounts[2] = base.Meta(delegate)
	return inst
}

// GetDelegateAccount gets the "delegate" account.
// The delegate.
func (inst *ApproveChecked) GetDelegateAccount() *base.AccountMeta {
	return inst.Accounts[2]
}

// SetOwnerAccount sets the "owner" account.
// The source account owner.
func (inst *ApproveChecked) SetOwnerAccount(owner base.PublicKey, multisigSigners ...base.PublicKey) *ApproveChecked {
	inst.Accounts[3] = base.Meta(owner)
	if len(multisigSigners) == 0 {
		inst.Accounts[3].SIGNER()
	}
	for _, signer := range multisigSigners {
		inst.Signers = append(inst.Signers, base.Meta(signer).SIGNER())
	}
	return inst
}

// GetOwnerAccount gets the "owner" account.
// The source account owner.
func (inst *ApproveChecked) GetOwnerAccount() *base.AccountMeta {
	return inst.Accounts[3]
}

func (inst ApproveChecked) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint8(Instruction_ApproveChecked),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst ApproveChecked) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *ApproveChecked) Validate() error {
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
			return errors.New("accounts.Delegate is not set")
		}
		if inst.Accounts[3] == nil {
			return errors.New("accounts.Owner is not set")
		}
		if !inst.Accounts[3].IsSigner && len(inst.Signers) == 0 {
			return fmt.Errorf("accounts.Signers is not set")
		}
		if len(inst.Signers) > MAX_SIGNERS {
			return fmt.Errorf("too many signers; got %v, but max is 11", len(inst.Signers))
		}
	}
	return nil
}

func (obj ApproveChecked) MarshalWithEncoder(encoder *base.Encoder) (err error) {
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

// NewApproveCheckedInstruction declares a new ApproveChecked instruction with the provided parameters and accounts.
func NewApproveCheckedInstruction(
	// Parameters:
	amount uint64,
	decimals uint8,
	// Accounts:
	source base.PublicKey,
	mint base.PublicKey,
	delegate base.PublicKey,
	owner base.PublicKey,
	multisigSigners []base.PublicKey,
) *ApproveChecked {
	return NewApproveCheckedInstructionBuilder().
		SetAmount(amount).
		SetDecimals(decimals).
		SetSourceAccount(source).
		SetMintAccount(mint).
		SetDelegateAccount(delegate).
		SetOwnerAccount(owner, multisigSigners...)
}
