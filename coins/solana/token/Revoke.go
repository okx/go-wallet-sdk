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

// Revokes the delegate's authority.
type Revoke struct {

	// [0] = [WRITE] source
	// ··········· The source account.
	//
	// [1] = [] owner
	// ··········· The source account's owner.
	//
	// [2...] = [SIGNER] signers
	// ··········· M signer accounts.
	Accounts base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
	Signers  base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func (obj *Revoke) SetAccounts(accounts []*base.AccountMeta) error {
	obj.Accounts, obj.Signers = base.AccountMetaSlice(accounts).SplitFrom(2)
	return nil
}

func (slice Revoke) GetAccounts() (accounts []*base.AccountMeta) {
	accounts = append(accounts, slice.Accounts...)
	accounts = append(accounts, slice.Signers...)
	return
}

// NewRevokeInstructionBuilder creates a new `Revoke` instruction builder.
func NewRevokeInstructionBuilder() *Revoke {
	nd := &Revoke{
		Accounts: make(base.AccountMetaSlice, 2),
		Signers:  make(base.AccountMetaSlice, 0),
	}
	return nd
}

// SetSourceAccount sets the "source" account.
// The source account.
func (inst *Revoke) SetSourceAccount(source base.PublicKey) *Revoke {
	inst.Accounts[0] = base.Meta(source).WRITE()
	return inst
}

// GetSourceAccount gets the "source" account.
// The source account.
func (inst *Revoke) GetSourceAccount() *base.AccountMeta {
	return inst.Accounts[0]
}

// SetOwnerAccount sets the "owner" account.
// The source account's owner.
func (inst *Revoke) SetOwnerAccount(owner base.PublicKey, multisigSigners ...base.PublicKey) *Revoke {
	inst.Accounts[1] = base.Meta(owner)
	if len(multisigSigners) == 0 {
		inst.Accounts[1].SIGNER()
	}
	for _, signer := range multisigSigners {
		inst.Signers = append(inst.Signers, base.Meta(signer).SIGNER())
	}
	return inst
}

// GetOwnerAccount gets the "owner" account.
// The source account's owner.
func (inst *Revoke) GetOwnerAccount() *base.AccountMeta {
	return inst.Accounts[1]
}

func (inst Revoke) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint8(Instruction_Revoke),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst Revoke) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *Revoke) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if inst.Accounts[0] == nil {
			return errors.New("accounts.Source is not set")
		}
		if inst.Accounts[1] == nil {
			return errors.New("accounts.Owner is not set")
		}
		if !inst.Accounts[1].IsSigner && len(inst.Signers) == 0 {
			return fmt.Errorf("accounts.Signers is not set")
		}
		if len(inst.Signers) > MAX_SIGNERS {
			return fmt.Errorf("too many signers; got %v, but max is 11", len(inst.Signers))
		}
	}
	return nil
}

func (obj Revoke) MarshalWithEncoder(encoder *base.Encoder) (err error) {
	return nil
}

// NewRevokeInstruction declares a new Revoke instruction with the provided parameters and accounts.
func NewRevokeInstruction(
	// Accounts:
	source base.PublicKey,
	owner base.PublicKey,
	multisigSigners []base.PublicKey,
) *Revoke {
	return NewRevokeInstructionBuilder().
		SetSourceAccount(source).
		SetOwnerAccount(owner, multisigSigners...)
}
