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

// Consumes a stored nonce, replacing it with a successor
type AdvanceNonceAccount struct {

	// [0] = [WRITE] NonceAccount
	// ··········· Nonce account
	//
	// [1] = [] $(SysVarRecentBlockHashesPubkey)
	// ··········· RecentBlockhashes sysvar
	//
	// [2] = [SIGNER] NonceAuthorityAccount
	// ··········· Nonce authority
	base.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewAdvanceNonceAccountInstructionBuilder creates a new `AdvanceNonceAccount` instruction builder.
func NewAdvanceNonceAccountInstructionBuilder() *AdvanceNonceAccount {
	nd := &AdvanceNonceAccount{
		AccountMetaSlice: make(base.AccountMetaSlice, 3),
	}
	nd.AccountMetaSlice[1] = base.Meta(base.SysVarRecentBlockHashesPubkey)
	return nd
}

// Nonce account
func (inst *AdvanceNonceAccount) SetNonceAccount(nonceAccount base.PublicKey) *AdvanceNonceAccount {
	inst.AccountMetaSlice[0] = base.Meta(nonceAccount).WRITE()
	return inst
}

func (inst *AdvanceNonceAccount) GetNonceAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[0]
}

// RecentBlockhashes sysvar
func (inst *AdvanceNonceAccount) SetSysVarRecentBlockHashesPubkeyAccount(SysVarRecentBlockHashesPubkey base.PublicKey) *AdvanceNonceAccount {
	inst.AccountMetaSlice[1] = base.Meta(SysVarRecentBlockHashesPubkey)
	return inst
}

func (inst *AdvanceNonceAccount) GetSysVarRecentBlockHashesPubkeyAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[1]
}

// Nonce authority
func (inst *AdvanceNonceAccount) SetNonceAuthorityAccount(nonceAuthorityAccount base.PublicKey) *AdvanceNonceAccount {
	inst.AccountMetaSlice[2] = base.Meta(nonceAuthorityAccount).SIGNER()
	return inst
}

func (inst *AdvanceNonceAccount) GetNonceAuthorityAccount() *base.AccountMeta {
	return inst.AccountMetaSlice[2]
}

func (inst AdvanceNonceAccount) Build() *Instruction {
	return &Instruction{BaseVariant: base.BaseVariant{
		Impl:   inst,
		TypeID: base.TypeIDFromUint32(Instruction_AdvanceNonceAccount, binary.LittleEndian),
	}}
}

func (inst AdvanceNonceAccount) MarshalWithEncoder(encoder *base.Encoder) error {
	return nil
}

// NewAdvanceNonceAccountInstruction declares a new AdvanceNonceAccount instruction with the provided parameters and accounts.
func NewAdvanceNonceAccountInstruction(
	// Accounts:
	nonceAccount base.PublicKey,
	SysVarRecentBlockHashesPubkey base.PublicKey,
	nonceAuthorityAccount base.PublicKey) *AdvanceNonceAccount {
	return NewAdvanceNonceAccountInstructionBuilder().
		SetNonceAccount(nonceAccount).
		SetSysVarRecentBlockHashesPubkeyAccount(SysVarRecentBlockHashesPubkey).
		SetNonceAuthorityAccount(nonceAuthorityAccount)
}
