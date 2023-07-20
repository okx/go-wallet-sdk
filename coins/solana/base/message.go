// Copyright 2021 github.com/gagliardetto
// This file has been modified by github.com/gagliardetto
//
// Copyright 2020 dfuse Platform Inc.
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

package base

import (
	"encoding/base64"
	"fmt"
)

type Instruction interface {
	ProgramID() PublicKey     // the programID the instruction acts on
	Accounts() []*AccountMeta // returns the list of accounts the instructions requires
	Data() ([]byte, error)    // the binary encoded instructions
}

type MessageVersion int

const (
	MessageVersionLegacy MessageVersion = 0 // default
	MessageVersionV0     MessageVersion = 1 // v0
)

type MessageAddressTableLookupSlice []MessageAddressTableLookup

// NumLookups returns the number of accounts in all the MessageAddressTableLookupSlice
func (lookups MessageAddressTableLookupSlice) NumLookups() int {
	count := 0
	for _, lookup := range lookups {
		// TODO: check if this is correct.
		count += len(lookup.ReadonlyIndexes)
		count += len(lookup.WritableIndexes)
	}
	return count
}

// GetTableIDs returns the list of all address table IDs.
func (lookups MessageAddressTableLookupSlice) GetTableIDs() PublicKeySlice {
	if lookups == nil {
		return nil
	}
	ids := make(PublicKeySlice, 0)
	for _, lookup := range lookups {
		ids.UniqueAppend(lookup.AccountKey)
	}
	return ids
}

type MessageAddressTableLookup struct {
	AccountKey      PublicKey // The account key of the address table.
	WritableIndexes []uint8
	ReadonlyIndexes []uint8
}

type Message struct {
	version MessageVersion
	// List of base-58 encoded public keys used by the transaction,
	// including by the instructions and for signatures.
	// The first `message.header.numRequiredSignatures` public keys must sign the transaction.
	AccountKeys []PublicKey `json:"accountKeys"`

	// Details the account types and signatures required by the transaction.
	Header MessageHeader `json:"header"`

	// A base-58 encoded hash of a recent block in the ledger used to
	// prevent transaction duplication and to give transactions lifetimes.
	RecentBlockhash Hash `json:"recentBlockhash"`

	// List of program instructions that will be executed in sequence
	// and committed in one atomic transaction if all succeed.
	Instructions []CompiledInstruction `json:"instructions"`

	// List of address table lookups used to load additional accounts for this transaction.
	addressTableLookups MessageAddressTableLookupSlice

	// The actual tables that contain the list of account pubkeys.
	// NOTE: you need to fetch these from the chain, and then call `SetAddressTables`
	// before you use this transaction -- otherwise, you will get a panic.
	addressTables map[PublicKey][]PublicKey
}

func (mx *Message) MarshalBinary() ([]byte, error) {
	buf := []byte{
		mx.Header.NumRequiredSignatures,
		mx.Header.NumReadonlySignedAccounts,
		mx.Header.NumReadonlyUnsignedAccounts,
	}

	EncodeCompactU16Length(&buf, len(mx.AccountKeys))
	for _, key := range mx.AccountKeys {
		buf = append(buf, key[:]...)
	}

	buf = append(buf, mx.RecentBlockhash[:]...)

	EncodeCompactU16Length(&buf, len(mx.Instructions))
	for _, instruction := range mx.Instructions {
		buf = append(buf, byte(instruction.ProgramIDIndex))
		EncodeCompactU16Length(&buf, len(instruction.Accounts))
		for _, accountIdx := range instruction.Accounts {
			buf = append(buf, byte(accountIdx))
		}

		EncodeCompactU16Length(&buf, len(instruction.Data))
		buf = append(buf, instruction.Data...)
	}
	return buf, nil
}

// Signers returns the pubkeys of all accounts that are signers.
func (m *Message) Signers() PublicKeySlice {
	out := make(PublicKeySlice, 0, len(m.AccountKeys))
	for _, a := range m.AccountKeys {
		if m.IsSigner(a) {
			out = append(out, a)
		}
	}
	return out
}

// Writable returns the pubkeys of all accounts that are writable.
func (m *Message) Writable() (out PublicKeySlice) {
	for _, a := range m.AccountKeys {
		if m.IsWritable(a) {
			out = append(out, a)
		}
	}
	return out
}

func (m *Message) IsSigner(account PublicKey) bool {
	for idx, acc := range m.AccountKeys {
		if acc.Equals(account) {
			return idx < int(m.Header.NumRequiredSignatures)
		}
	}
	return false
}

func (m *Message) IsWritable(account PublicKey) bool {
	index := 0
	found := false
	for idx, acc := range m.AccountKeys {
		if acc.Equals(account) {
			found = true
			index = idx
		}
	}
	if !found {
		return false
	}
	h := m.Header
	return (index < int(h.NumRequiredSignatures-h.NumReadonlySignedAccounts)) ||
		((index >= int(h.NumRequiredSignatures)) && (index < len(m.AccountKeys)-int(h.NumReadonlyUnsignedAccounts)))
}

func (m *Message) signerKeys() []PublicKey {
	return m.AccountKeys[0:m.Header.NumRequiredSignatures]
}

type MessageHeader struct {
	// The total number of signatures required to make the transaction valid.
	// The signatures must match the first `numRequiredSignatures` of `message.account_keys`.
	NumRequiredSignatures uint8 `json:"numRequiredSignatures"`

	// The last numReadonlySignedAccounts of the signed keys are read-only accounts.
	// Programs may process multiple transactions that load read-only accounts within
	// a single PoH entry, but are not permitted to credit or debit lamports or modify
	// account data.
	// Transactions targeting the same read-write account are evaluated sequentially.
	NumReadonlySignedAccounts uint8 `json:"numReadonlySignedAccounts"`

	// The last `numReadonlyUnsignedAccounts` of the unsigned keys are read-only accounts.
	NumReadonlyUnsignedAccounts uint8 `json:"numReadonlyUnsignedAccounts"`
}

func (mx *Message) UnmarshalWithDecoder(decoder *Decoder) (err error) {
	// peek first byte to determine if this is a legacy or v0 message
	versionNum, err := decoder.Peek(1)
	if err != nil {
		return err
	}
	// TODO: is this the right way to determine if this is a legacy or v0 message?
	if versionNum[0] < 127 {
		mx.version = MessageVersionLegacy
	} else {
		mx.version = MessageVersionV0
	}
	switch mx.version {
	case MessageVersionV0:
		return mx.UnmarshalV0(decoder)
	case MessageVersionLegacy:
		return mx.UnmarshalLegacy(decoder)
	default:
		return fmt.Errorf("invalid message version: %d", mx.version)
	}
}

func (mx *Message) UnmarshalBase64(b64 string) error {
	b, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}
	return mx.UnmarshalWithDecoder(NewBinDecoder(b))
}

func (mx *Message) UnmarshalV0(decoder *Decoder) (err error) {
	version, err := decoder.ReadByte()
	if err != nil {
		return fmt.Errorf("failed to read message version: %w", err)
	}
	// TODO: check version
	mx.version = MessageVersion(version - 127)

	// The middle of the message is the same as the legacy message:
	err = mx.UnmarshalLegacy(decoder)
	if err != nil {
		return err
	}

	// Read address table lookups length:
	addressTableLookupsLen, err := decoder.ReadByte()
	if err != nil {
		return fmt.Errorf("failed to read address table lookups length: %w", err)
	}
	if addressTableLookupsLen > 0 {
		mx.addressTableLookups = make([]MessageAddressTableLookup, addressTableLookupsLen)
		for i := 0; i < int(addressTableLookupsLen); i++ {
			// read account pubkey
			_, err = decoder.Read(mx.addressTableLookups[i].AccountKey[:])
			if err != nil {
				return fmt.Errorf("failed to read account pubkey: %w", err)
			}

			// read writable indexes
			writableIndexesLen, err := decoder.ReadCompactU16Length()
			if err != nil {
				return fmt.Errorf("failed to read writable indexes length: %w", err)
			}
			mx.addressTableLookups[i].WritableIndexes = make([]byte, writableIndexesLen)
			_, err = decoder.Read(mx.addressTableLookups[i].WritableIndexes)
			if err != nil {
				return fmt.Errorf("failed to read writable indexes: %w", err)
			}

			// read readonly indexes
			readonlyIndexesLen, err := decoder.ReadCompactU16Length()
			if err != nil {
				return fmt.Errorf("failed to read readonly indexes length: %w", err)
			}
			mx.addressTableLookups[i].ReadonlyIndexes = make([]byte, readonlyIndexesLen)
			_, err = decoder.Read(mx.addressTableLookups[i].ReadonlyIndexes)
			if err != nil {
				return fmt.Errorf("failed to read readonly indexes: %w", err)
			}
		}
	}
	return nil
}

func (mx *Message) UnmarshalLegacy(decoder *Decoder) (err error) {
	{
		mx.Header.NumRequiredSignatures, err = decoder.ReadUint8()
		if err != nil {
			return fmt.Errorf("unable to decode mx.Header.NumRequiredSignatures: %w", err)
		}
		mx.Header.NumReadonlySignedAccounts, err = decoder.ReadUint8()
		if err != nil {
			return fmt.Errorf("unable to decode mx.Header.NumReadonlySignedAccounts: %w", err)
		}
		mx.Header.NumReadonlyUnsignedAccounts, err = decoder.ReadUint8()
		if err != nil {
			return fmt.Errorf("unable to decode mx.Header.NumReadonlyUnsignedAccounts: %w", err)
		}
	}
	{
		numAccountKeys, err := decoder.ReadCompactU16()
		if err != nil {
			return fmt.Errorf("unable to decode numAccountKeys: %w", err)
		}
		mx.AccountKeys = make([]PublicKey, numAccountKeys)
		for i := 0; i < numAccountKeys; i++ {
			_, err := decoder.Read(mx.AccountKeys[i][:])
			if err != nil {
				return fmt.Errorf("unable to decode mx.AccountKeys[%d]: %w", i, err)
			}
		}
	}
	{
		_, err := decoder.Read(mx.RecentBlockhash[:])
		if err != nil {
			return fmt.Errorf("unable to decode mx.RecentBlockhash: %w", err)
		}
	}
	{
		numInstructions, err := decoder.ReadCompactU16()
		if err != nil {
			return fmt.Errorf("unable to decode numInstructions: %w", err)
		}
		mx.Instructions = make([]CompiledInstruction, numInstructions)
		for instructionIndex := 0; instructionIndex < numInstructions; instructionIndex++ {
			programIDIndex, err := decoder.ReadUint8()
			if err != nil {
				return fmt.Errorf("unable to decode mx.Instructions[%d].ProgramIDIndex: %w", instructionIndex, err)
			}
			mx.Instructions[instructionIndex].ProgramIDIndex = uint16(programIDIndex)

			{
				numAccounts, err := decoder.ReadCompactU16()
				if err != nil {
					return fmt.Errorf("unable to decode numAccounts for ix[%d]: %w", instructionIndex, err)
				}
				mx.Instructions[instructionIndex].Accounts = make([]uint16, numAccounts)
				for i := 0; i < numAccounts; i++ {
					accountIndex, err := decoder.ReadUint8()
					if err != nil {
						return fmt.Errorf("unable to decode accountIndex for ix[%d].Accounts[%d]: %w", instructionIndex, i, err)
					}
					mx.Instructions[instructionIndex].Accounts[i] = uint16(accountIndex)
				}
			}
			{
				dataLen, err := decoder.ReadCompactU16()
				if err != nil {
					return fmt.Errorf("unable to decode dataLen for ix[%d]: %w", instructionIndex, err)
				}
				dataBytes, err := decoder.ReadNBytes(dataLen)
				if err != nil {
					return fmt.Errorf("unable to decode dataBytes for ix[%d]: %w", instructionIndex, err)
				}
				mx.Instructions[instructionIndex].Data = dataBytes
			}
		}
	}

	return nil
}
