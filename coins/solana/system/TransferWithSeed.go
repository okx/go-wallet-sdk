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
	"errors"
	"fmt"

	"github.com/emresenyuva/go-wallet-sdk/coins/solana/base"
)

// Transfer lamports from a derived address
type TransferWithSeed struct {
	// Amount to transfer
	Lamports *uint64

	// Seed to use to derive the funding account address
	FromSeed *string

	// Owner to use to derive the funding account address
	FromOwner *base.PublicKey

	// [0] = [WRITE] FundingAccount
	// ··········· Funding account
	//
	// [1] = [SIGNER] BaseForFundingAccount
	// ··········· Base for funding account
	//
	// [2] = [WRITE] RecipientAccount
	// ··········· Recipient account
	base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewTransferWithSeedInstructionBuilder creates a new `TransferWithSeed` instruction builder.
func NewTransferWithSeedInstructionBuilder() *TransferWithSeed {
	nd := &TransferWithSeed{
		AccountMetaSlice: make(base.AccountMetaSlice, 3),
	}
	return nd
}

// Amount to transfer
func (inst *TransferWithSeed) SetLamports(lamports uint64) *TransferWithSeed {
	inst.Lamports = &lamports
	return inst
}

// Seed to use to derive the funding account address
func (inst *TransferWithSeed) SetFromSeed(from_seed string) *TransferWithSeed {
	inst.FromSeed = &from_seed
	return inst
}

// Owner to use to derive the funding account address
func (inst *TransferWithSeed) SetFromOwner(from_owner base.PublicKey) *TransferWithSeed {
	inst.FromOwner = &from_owner
	return inst
}

// Funding account
func (inst *TransferWithSeed) SetFundingAccount(fundingAccount base.PublicKey) *TransferWithSeed {
	inst.AccountMetaSlice[0] = base.Meta(fundingAccount).WRITE()
	return inst
}

func (inst *TransferWithSeed) GetFundingAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[0]
}

// Base for funding account
func (inst *TransferWithSeed) SetBaseForFundingAccount(baseForFundingAccount base.PublicKey) *TransferWithSeed {
	inst.AccountMetaSlice[1] = base.Meta(baseForFundingAccount).SIGNER()
	return inst
}

func (inst *TransferWithSeed) GetBaseForFundingAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[1]
}

// Recipient account
func (inst *TransferWithSeed) SetRecipientAccount(recipientAccount base.PublicKey) *TransferWithSeed {
	inst.AccountMetaSlice[2] = base.Meta(recipientAccount).WRITE()
	return inst
}

func (inst *TransferWithSeed) GetRecipientAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[2]
}

func (inst TransferWithSeed) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint32(Instruction_TransferWithSeed, binary.LittleEndian),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst TransferWithSeed) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *TransferWithSeed) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Lamports == nil {
			return errors.New("Lamports parameter is not set")
		}
		if inst.FromSeed == nil {
			return errors.New("FromSeed parameter is not set")
		}
		if inst.FromOwner == nil {
			return errors.New("FromOwner parameter is not set")
		}
	}

	// Check whether all accounts are set:
	for accIndex, acc := range inst.AccountMetaSlice {
		if acc == nil {
			return fmt.Errorf("ins.AccountMetaSlice[%v] is not set", accIndex)
		}
	}
	return nil
}

func (inst TransferWithSeed) MarshalWithEncoder(encoder *base.Encoder) error {
	// Serialize `Lamports` param:
	{
		err := encoder.Encode(*inst.Lamports)
		if err != nil {
			return err
		}
	}
	// Serialize `FromSeed` param:
	{
		err := encoder.WriteRustString(*inst.FromSeed)
		if err != nil {
			return err
		}
	}
	// Serialize `FromOwner` param:
	{
		err := encoder.Encode(*inst.FromOwner)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewTransferWithSeedInstruction declares a new TransferWithSeed instruction with the provided parameters and accounts.
func NewTransferWithSeedInstruction(
	// Parameters:
	lamports uint64,
	from_seed string,
	from_owner base.PublicKey,
	// Accounts:
	fundingAccount base.PublicKey,
	baseForFundingAccount base.PublicKey,
	recipientAccount base.PublicKey) *TransferWithSeed {
	return NewTransferWithSeedInstructionBuilder().
		SetLamports(lamports).
		SetFromSeed(from_seed).
		SetFromOwner(from_owner).
		SetFundingAccount(fundingAccount).
		SetBaseForFundingAccount(baseForFundingAccount).
		SetRecipientAccount(recipientAccount)
}
