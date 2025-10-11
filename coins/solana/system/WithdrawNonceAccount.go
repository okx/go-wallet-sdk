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

// Withdraw funds from a nonce account
type WithdrawNonceAccount struct {
	// The u64 parameter is the lamports to withdraw, which must leave the account balance above the rent exempt reserve or at zero.
	Lamports *uint64

	// [0] = [WRITE] NonceAccount
	// ··········· Nonce account
	//
	// [1] = [WRITE] RecipientAccount
	// ··········· Recipient account
	//
	// [2] = [] $(SysVarRecentBlockHashesPubkey)
	// ··········· RecentBlockhashes sysvar
	//
	// [3] = [] $(SysVarRentPubkey)
	// ··········· Rent sysvar
	//
	// [4] = [SIGNER] NonceAuthorityAccount
	// ··········· Nonce authority
	base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewWithdrawNonceAccountInstructionBuilder creates a new `WithdrawNonceAccount` instruction builder.
func NewWithdrawNonceAccountInstructionBuilder() *WithdrawNonceAccount {
	nd := &WithdrawNonceAccount{
		AccountMetaSlice: make(base.AccountMetaSlice, 5),
	}
	nd.AccountMetaSlice[2] = base.Meta(base.SysVarRecentBlockHashesPubkey)
	nd.AccountMetaSlice[3] = base.Meta(base.SysVarRentPubkey)
	return nd
}

// The u64 parameter is the lamports to withdraw, which must leave the account balance above the rent exempt reserve or at zero.
func (inst *WithdrawNonceAccount) SetLamports(lamports uint64) *WithdrawNonceAccount {
	inst.Lamports = &lamports
	return inst
}

// Nonce account
func (inst *WithdrawNonceAccount) SetNonceAccount(nonceAccount base.PublicKey) *WithdrawNonceAccount {
	inst.AccountMetaSlice[0] = base.Meta(nonceAccount).WRITE()
	return inst
}

func (inst *WithdrawNonceAccount) GetNonceAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[0]
}

// Recipient account
func (inst *WithdrawNonceAccount) SetRecipientAccount(recipientAccount base.PublicKey) *WithdrawNonceAccount {
	inst.AccountMetaSlice[1] = base.Meta(recipientAccount).WRITE()
	return inst
}

func (inst *WithdrawNonceAccount) GetRecipientAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[1]
}

// RecentBlockhashes sysvar
func (inst *WithdrawNonceAccount) SetSysVarRecentBlockHashesPubkeyAccount(SysVarRecentBlockHashesPubkey base.PublicKey) *WithdrawNonceAccount {
	inst.AccountMetaSlice[2] = base.Meta(SysVarRecentBlockHashesPubkey)
	return inst
}

func (inst *WithdrawNonceAccount) GetSysVarRecentBlockHashesPubkeyAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[2]
}

// Rent sysvar
func (inst *WithdrawNonceAccount) SetSysVarRentPubkeyAccount(SysVarRentPubkey base.PublicKey) *WithdrawNonceAccount {
	inst.AccountMetaSlice[3] = base.Meta(SysVarRentPubkey)
	return inst
}

func (inst *WithdrawNonceAccount) GetSysVarRentPubkeyAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[3]
}

// Nonce authority
func (inst *WithdrawNonceAccount) SetNonceAuthorityAccount(nonceAuthorityAccount base.PublicKey) *WithdrawNonceAccount {
	inst.AccountMetaSlice[4] = base.Meta(nonceAuthorityAccount).SIGNER()
	return inst
}

func (inst *WithdrawNonceAccount) GetNonceAuthorityAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[4]
}

func (inst WithdrawNonceAccount) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint32(Instruction_WithdrawNonceAccount, binary.LittleEndian),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst WithdrawNonceAccount) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *WithdrawNonceAccount) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Lamports == nil {
			return errors.New("Lamports parameter is not set")
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

func (inst WithdrawNonceAccount) MarshalWithEncoder(encoder *base.Encoder) error {
	// Serialize `Lamports` param:
	{
		err := encoder.Encode(*inst.Lamports)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewWithdrawNonceAccountInstruction declares a new WithdrawNonceAccount instruction with the provided parameters and accounts.
func NewWithdrawNonceAccountInstruction(
	// Parameters:
	lamports uint64,
	// Accounts:
	nonceAccount base.PublicKey,
	recipientAccount base.PublicKey,
	nonceAuthorityAccount base.PublicKey) *WithdrawNonceAccount {
	return NewWithdrawNonceAccountInstructionBuilder().
		SetLamports(lamports).
		SetNonceAccount(nonceAccount).
		SetRecipientAccount(recipientAccount).
		SetSysVarRecentBlockHashesPubkeyAccount(base.SysVarRecentBlockHashesPubkey).
		SetSysVarRentPubkeyAccount(base.SysVarRentPubkey).
		SetNonceAuthorityAccount(nonceAuthorityAccount)
}
