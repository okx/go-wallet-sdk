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

package system

import (
	"encoding/binary"

	"github.com/okx/go-wallet-sdk/coins/solana/base"
)

// Create a new account
type CreateAccount struct {
	// Number of lamports to transfer to the new account
	Lamports *uint64

	// Number of bytes of memory to allocate
	Space *uint64

	// Address of program that will own the new account
	Owner *base.PublicKey

	// [0] = [WRITE, SIGNER] FundingAccount
	// ··········· Funding account
	//
	// [1] = [WRITE, SIGNER] NewAccount
	// ··········· New account
	base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewCreateAccountInstructionBuilder creates a new `CreateAccount` instruction builder.
func NewCreateAccountInstructionBuilder() *CreateAccount {
	nd := &CreateAccount{
		AccountMetaSlice: make(base.AccountMetaSlice, 2),
	}
	return nd
}

// Number of lamports to transfer to the new account
func (inst *CreateAccount) SetLamports(lamports uint64) *CreateAccount {
	inst.Lamports = &lamports
	return inst
}

// Number of bytes of memory to allocate
func (inst *CreateAccount) SetSpace(space uint64) *CreateAccount {
	inst.Space = &space
	return inst
}

// Address of program that will own the new account
func (inst *CreateAccount) SetOwner(owner base.PublicKey) *CreateAccount {
	inst.Owner = &owner
	return inst
}

// Funding account
func (inst *CreateAccount) SetFundingAccount(fundingAccount base.PublicKey) *CreateAccount {
	inst.AccountMetaSlice[0] = base.Meta(fundingAccount).WRITE().SIGNER()
	return inst
}

func (inst *CreateAccount) GetFundingAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[0]
}

// New account
func (inst *CreateAccount) SetNewAccount(newAccount base.PublicKey) *CreateAccount {
	inst.AccountMetaSlice[1] = base.Meta(newAccount).WRITE().SIGNER()
	return inst
}

func (inst *CreateAccount) GetNewAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[1]
}

func (inst CreateAccount) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint32(Instruction_CreateAccount, binary.LittleEndian),
	}}
}

func (inst CreateAccount) MarshalWithEncoder(encoder *base.Encoder) error {
	// Serialize `Lamports` param:
	{
		err := encoder.Encode(*inst.Lamports)
		if err != nil {
			return err
		}
	}
	// Serialize `Space` param:
	{
		err := encoder.Encode(*inst.Space)
		if err != nil {
			return err
		}
	}
	// Serialize `Owner` param:
	{
		err := encoder.Encode(*inst.Owner)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewCreateAccountInstruction declares a new CreateAccount instruction with the provided parameters and accounts.
func NewCreateAccountInstruction(
	// Parameters:
	lamports uint64,
	space uint64,
	owner base.PublicKey,
	// Accounts:
	fundingAccount base.PublicKey,
	newAccount base.PublicKey) *CreateAccount {
	return NewCreateAccountInstructionBuilder().
		SetLamports(lamports).
		SetSpace(space).
		SetOwner(owner).
		SetFundingAccount(fundingAccount).
		SetNewAccount(newAccount)
}
