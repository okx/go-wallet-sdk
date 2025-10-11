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

// Close an account by transferring all its SOL to the destination account.
// Non-native accounts may only be closed if its token amount is zero.
type CloseAccount struct {

	// [0] = [WRITE] account
	// ··········· The account to close.
	//
	// [1] = [WRITE] destination
	// ··········· The destination account.
	//
	// [2] = [] owner
	// ··········· The account's owner.
	//
	// [3...] = [SIGNER] signers
	// ··········· M signer accounts.
	Accounts base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
	Signers  base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func (obj *CloseAccount) SetAccounts(accounts []*base.AccountMeta) error {
	obj.Accounts, obj.Signers = base.AccountMetaSlice(accounts).SplitFrom(3)
	return nil
}

func (slice CloseAccount) GetAccounts() (accounts []*base.AccountMeta) {
	accounts = append(accounts, slice.Accounts...)
	accounts = append(accounts, slice.Signers...)
	return
}

// NewCloseAccountInstructionBuilder creates a new `CloseAccount` instruction builder.
func NewCloseAccountInstructionBuilder() *CloseAccount {
	nd := &CloseAccount{
		Accounts: make(base.AccountMetaSlice, 3),
		Signers:  make(base.AccountMetaSlice, 0),
	}
	return nd
}

// SetAccount sets the "account" account.
// The account to close.
func (inst *CloseAccount) SetAccount(account base.PublicKey) *CloseAccount {
	inst.Accounts[0] = base.Meta(account).WRITE()
	return inst
}

// GetAccount gets the "account" account.
// The account to close.
func (inst *CloseAccount) GetAccount() *base.AccountMeta {
	return inst.Accounts[0]
}

// SetDestinationAccount sets the "destination" account.
// The destination account.
func (inst *CloseAccount) SetDestinationAccount(destination base.PublicKey) *CloseAccount {
	inst.Accounts[1] = base.Meta(destination).WRITE()
	return inst
}

// GetDestinationAccount gets the "destination" account.
// The destination account.
func (inst *CloseAccount) GetDestinationAccount() *base.AccountMeta {
	return inst.Accounts[1]
}

// SetOwnerAccount sets the "owner" account.
// The account's owner.
func (inst *CloseAccount) SetOwnerAccount(owner base.PublicKey, multisigSigners ...base.PublicKey) *CloseAccount {
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
// The account's owner.
func (inst *CloseAccount) GetOwnerAccount() *base.AccountMeta {
	return inst.Accounts[2]
}

func (inst CloseAccount) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint8(Instruction_CloseAccount),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst CloseAccount) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *CloseAccount) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if inst.Accounts[0] == nil {
			return errors.New("accounts.Account is not set")
		}
		if inst.Accounts[1] == nil {
			return errors.New("accounts.Destination is not set")
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

func (obj CloseAccount) MarshalWithEncoder(encoder *base.Encoder) (err error) {
	return nil
}

// NewCloseAccountInstruction declares a new CloseAccount instruction with the provided parameters and accounts.
func NewCloseAccountInstruction(
	// Accounts:
	account base.PublicKey,
	destination base.PublicKey,
	owner base.PublicKey,
	multisigSigners []base.PublicKey,
) *CloseAccount {
	return NewCloseAccountInstructionBuilder().
		SetAccount(account).
		SetDestinationAccount(destination).
		SetOwnerAccount(owner, multisigSigners...)
}
