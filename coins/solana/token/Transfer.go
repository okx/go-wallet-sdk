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

// Transfers tokens from one account to another either directly or via a
// delegate.  If this account is associated with the native mint then equal
// amounts of SOL and Tokens will be transferred to the destination
// account.
type Transfer struct {
	// The amount of tokens to transfer.
	Amount *uint64

	// [0] = [WRITE] source
	// ··········· The source account.
	//
	// [1] = [WRITE] destination
	// ··········· The destination account.
	//
	// [2] = [] owner
	// ··········· The source account owner/delegate.
	//
	// [3...] = [SIGNER] signers
	// ··········· M signer accounts.
	Accounts base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
	Signers  base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func (obj *Transfer) SetAccounts(accounts []*base.AccountMeta) error {
	obj.Accounts, obj.Signers = base.AccountMetaSlice(accounts).SplitFrom(3)
	return nil
}

func (slice Transfer) GetAccounts() (accounts []*base.AccountMeta) {
	accounts = append(accounts, slice.Accounts...)
	accounts = append(accounts, slice.Signers...)
	return
}

// NewTransferInstructionBuilder creates a new `Transfer` instruction builder.
func NewTransferInstructionBuilder() *Transfer {
	nd := &Transfer{
		Accounts: make(base.AccountMetaSlice, 3),
		Signers:  make(base.AccountMetaSlice, 0),
	}
	return nd
}

// SetAmount sets the "amount" parameter.
// The amount of tokens to transfer.
func (inst *Transfer) SetAmount(amount uint64) *Transfer {
	inst.Amount = &amount
	return inst
}

// SetSourceAccount sets the "source" account.
// The source account.
func (inst *Transfer) SetSourceAccount(source base.PublicKey) *Transfer {
	inst.Accounts[0] = base.Meta(source).WRITE()
	return inst
}

// GetSourceAccount gets the "source" account.
// The source account.
func (inst *Transfer) GetSourceAccount() *base.AccountMeta {
	return inst.Accounts[0]
}

// SetDestinationAccount sets the "destination" account.
// The destination account.
func (inst *Transfer) SetDestinationAccount(destination base.PublicKey) *Transfer {
	inst.Accounts[1] = base.Meta(destination).WRITE()
	return inst
}

// GetDestinationAccount gets the "destination" account.
// The destination account.
func (inst *Transfer) GetDestinationAccount() *base.AccountMeta {
	return inst.Accounts[1]
}

// SetOwnerAccount sets the "owner" account.
// The source account owner/delegate.
func (inst *Transfer) SetOwnerAccount(owner base.PublicKey, multisigSigners ...base.PublicKey) *Transfer {
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
// The source account owner/delegate.
func (inst *Transfer) GetOwnerAccount() *base.AccountMeta {
	return inst.Accounts[2]
}

func (inst Transfer) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint8(Instruction_Transfer),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst Transfer) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *Transfer) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Amount == nil {
			return errors.New("Amount parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.Accounts[0] == nil {
			return fmt.Errorf("accounts.Source is not set")
		}
		if inst.Accounts[1] == nil {
			return fmt.Errorf("accounts.Destination is not set")
		}
		if inst.Accounts[2] == nil {
			return fmt.Errorf("accounts.Owner is not set")
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

func (obj Transfer) MarshalWithEncoder(encoder *base.Encoder) (err error) {
	// Serialize `Amount` param:
	err = encoder.Encode(obj.Amount)
	if err != nil {
		return err
	}
	return nil
}

// NewTransferInstruction declares a new Transfer instruction with the provided parameters and accounts.
func NewTransferInstruction(
	// Parameters:
	amount uint64,
	// Accounts:
	source base.PublicKey,
	destination base.PublicKey,
	owner base.PublicKey,
	multisigSigners []base.PublicKey) *Transfer {
	return NewTransferInstructionBuilder().
		SetAmount(amount).
		SetSourceAccount(source).
		SetDestinationAccount(destination).
		SetOwnerAccount(owner, multisigSigners...)
}
