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

// Like InitializeMultisig, but does not require the Rent sysvar to be provided.
type InitializeMultisig2 struct {
	// The number of signers (M) required to validate this multisignature account.
	M *uint8

	// [0] = [WRITE] account
	// ··········· The multisignature account to initialize.
	//
	// [1] = [SIGNER] signers
	// ··········· The signer accounts, must equal to N where 1 <= N <= 11.
	Accounts base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
	Signers  base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

func (obj *InitializeMultisig2) SetAccounts(accounts []*base.AccountMeta) error {
	obj.Accounts, obj.Signers = base.AccountMetaSlice(accounts).SplitFrom(1)
	return nil
}

func (slice InitializeMultisig2) GetAccounts() (accounts []*base.AccountMeta) {
	accounts = append(accounts, slice.Accounts...)
	accounts = append(accounts, slice.Signers...)
	return
}

// NewInitializeMultisig2InstructionBuilder creates a new `InitializeMultisig2` instruction builder.
func NewInitializeMultisig2InstructionBuilder() *InitializeMultisig2 {
	nd := &InitializeMultisig2{
		Accounts: make(base.AccountMetaSlice, 1),
		Signers:  make(base.AccountMetaSlice, 0),
	}
	return nd
}

// SetM sets the "m" parameter.
// The number of signers (M) required to validate this multisignature account.
func (inst *InitializeMultisig2) SetM(m uint8) *InitializeMultisig2 {
	inst.M = &m
	return inst
}

// SetAccount sets the "account" account.
// The multisignature account to initialize.
func (inst *InitializeMultisig2) SetAccount(account base.PublicKey) *InitializeMultisig2 {
	inst.Accounts[0] = base.Meta(account).WRITE()
	return inst
}

// GetAccount gets the "account" account.
// The multisignature account to initialize.
func (inst *InitializeMultisig2) GetAccount() *base.AccountMeta {
	return inst.Accounts[0]
}

// AddSigners adds the "signers" accounts.
// The signer accounts, must equal to N where 1 <= N <= 11.
func (inst *InitializeMultisig2) AddSigners(signers ...base.PublicKey) *InitializeMultisig2 {
	for _, signer := range signers {
		inst.Signers = append(inst.Signers, base.Meta(signer).SIGNER())
	}
	return inst
}

func (inst InitializeMultisig2) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint8(Instruction_InitializeMultisig2),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst InitializeMultisig2) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *InitializeMultisig2) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.M == nil {
			return errors.New("M parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.Accounts[0] == nil {
			return errors.New("accounts.Account is not set")
		}
		if len(inst.Signers) == 0 {
			return fmt.Errorf("accounts.Signers is not set")
		}
		if len(inst.Signers) > MAX_SIGNERS {
			return fmt.Errorf("too many signers; got %v, but max is 11", len(inst.Signers))
		}
	}
	return nil
}

func (obj InitializeMultisig2) MarshalWithEncoder(encoder *base.Encoder) (err error) {
	// Serialize `M` param:
	err = encoder.Encode(obj.M)
	if err != nil {
		return err
	}
	return nil
}

// NewInitializeMultisig2Instruction declares a new InitializeMultisig2 instruction with the provided parameters and accounts.
func NewInitializeMultisig2Instruction(
	// Parameters:
	m uint8,
	// Accounts:
	account base.PublicKey,
	signers []base.PublicKey,
) *InitializeMultisig2 {
	return NewInitializeMultisig2InstructionBuilder().
		SetM(m).
		SetAccount(account).
		AddSigners(signers...)
}
