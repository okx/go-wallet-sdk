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

// Approves a delegate.  A delegate is given the authority over tokens on
// behalf of the source account's owner.
type Approve struct {
	// The amount of tokens the delegate is approved for.
	Amount *uint64

	// [0] = [WRITE] source
	// ··········· The source account.
	//
	// [1] = [] delegate
	// ··········· The delegate.
	//
	// [2] = [] owner
	// ··········· The source account owner.
	//
	// [3...] = [SIGNER] signers
	// ··········· M signer accounts.
	Accounts base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
	Signers  base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func (obj *Approve) SetAccounts(accounts []*base.AccountMeta) error {
	obj.Accounts, obj.Signers = base.AccountMetaSlice(accounts).SplitFrom(3)
	return nil
}

func (slice Approve) GetAccounts() (accounts []*base.AccountMeta) {
	accounts = append(accounts, slice.Accounts...)
	accounts = append(accounts, slice.Signers...)
	return
}

// NewApproveInstructionBuilder creates a new `Approve` instruction builder.
func NewApproveInstructionBuilder() *Approve {
	nd := &Approve{
		Accounts: make(base.AccountMetaSlice, 3),
		Signers:  make(base.AccountMetaSlice, 0),
	}
	return nd
}

// SetAmount sets the "amount" parameter.
// The amount of tokens the delegate is approved for.
func (inst *Approve) SetAmount(amount uint64) *Approve {
	inst.Amount = &amount
	return inst
}

// SetSourceAccount sets the "source" account.
// The source account.
func (inst *Approve) SetSourceAccount(source base.PublicKey) *Approve {
	inst.Accounts[0] = base.Meta(source).WRITE()
	return inst
}

// GetSourceAccount gets the "source" account.
// The source account.
func (inst *Approve) GetSourceAccount() *base.AccountMeta {
	return inst.Accounts[0]
}

// SetDelegateAccount sets the "delegate" account.
// The delegate.
func (inst *Approve) SetDelegateAccount(delegate base.PublicKey) *Approve {
	inst.Accounts[1] = base.Meta(delegate)
	return inst
}

// GetDelegateAccount gets the "delegate" account.
// The delegate.
func (inst *Approve) GetDelegateAccount() *base.AccountMeta {
	return inst.Accounts[1]
}

// SetOwnerAccount sets the "owner" account.
// The source account owner.
func (inst *Approve) SetOwnerAccount(owner base.PublicKey, multisigSigners ...base.PublicKey) *Approve {
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
// The source account owner.
func (inst *Approve) GetOwnerAccount() *base.AccountMeta {
	return inst.Accounts[2]
}

func (inst Approve) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint8(Instruction_Approve),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst Approve) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *Approve) Validate() error {
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
			return errors.New("accounts.Delegate is not set")
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

func (obj Approve) MarshalWithEncoder(encoder *base.Encoder) (err error) {
	// Serialize `Amount` param:
	err = encoder.Encode(obj.Amount)
	if err != nil {
		return err
	}
	return nil
}

// NewApproveInstruction declares a new Approve instruction with the provided parameters and accounts.
func NewApproveInstruction(
	// Parameters:
	amount uint64,
	// Accounts:
	source base.PublicKey,
	delegate base.PublicKey,
	owner base.PublicKey,
	multisigSigners []base.PublicKey,
) *Approve {
	return NewApproveInstructionBuilder().
		SetAmount(amount).
		SetSourceAccount(source).
		SetDelegateAccount(delegate).
		SetOwnerAccount(owner, multisigSigners...)
}
