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
	"encoding/binary"

	"github.com/okx/go-wallet-sdk/coins/solana/base"
)

type Mint struct {
	// Optional authority used to mint new tokens. The mint authority may only be provided during
	// mint creation. If no mint authority is present then the mint has a fixed supply and no
	// further tokens may be minted.
	MintAuthority *base.PublicKey `bin:"optional"`

	// Total supply of tokens.
	Supply uint64

	// Number of base 10 digits to the right of the decimal place.
	Decimals uint8

	// Is `true` if this structure has been initialized
	IsInitialized bool

	// Optional authority to freeze token accounts.
	FreezeAuthority *base.PublicKey `bin:"optional"`
}

func (mint Mint) MarshalWithEncoder(encoder *base.Encoder) (err error) {
	{
		if mint.MintAuthority == nil {
			err = encoder.WriteUint32(0, binary.LittleEndian)
			if err != nil {
				return err
			}
			empty := base.PublicKey{}
			err = encoder.WriteBytes(empty[:], false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteUint32(1, binary.LittleEndian)
			if err != nil {
				return err
			}
			err = encoder.WriteBytes(mint.MintAuthority[:], false)
			if err != nil {
				return err
			}
		}
	}
	err = encoder.WriteUint64(mint.Supply, binary.LittleEndian)
	if err != nil {
		return err
	}
	err = encoder.WriteUint8(mint.Decimals)
	if err != nil {
		return err
	}
	err = encoder.WriteBool(mint.IsInitialized)
	if err != nil {
		return err
	}
	{
		if mint.FreezeAuthority == nil {
			err = encoder.WriteUint32(0, binary.LittleEndian)
			if err != nil {
				return err
			}
			empty := base.PublicKey{}
			err = encoder.WriteBytes(empty[:], false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteUint32(1, binary.LittleEndian)
			if err != nil {
				return err
			}
			err = encoder.WriteBytes(mint.FreezeAuthority[:], false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type Account struct {
	// The mint associated with this account
	Mint base.PublicKey

	// The owner of this account.
	Owner base.PublicKey

	// The amount of tokens this account holds.
	Amount uint64

	// If `delegate` is `Some` then `delegated_amount` represents
	// the amount authorized by the delegate
	Delegate *base.PublicKey `bin:"optional"`

	// The account's state
	State AccountState

	// If is_some, this is a native token, and the value logs the rent-exempt reserve. An Account
	// is required to be rent-exempt, so the value is used by the Processor to ensure that wrapped
	// SOL accounts do not drop below this threshold.
	IsNative *uint64 `bin:"optional"`

	// The amount delegated
	DelegatedAmount uint64

	// Optional authority to close the account.
	CloseAuthority *base.PublicKey `bin:"optional"`
}

func (mint Account) MarshalWithEncoder(encoder *base.Encoder) (err error) {
	{
		err = encoder.WriteBytes(mint.Mint[:], false)
		if err != nil {
			return err
		}
	}
	{
		err = encoder.WriteBytes(mint.Owner[:], false)
		if err != nil {
			return err
		}
	}
	{
		err = encoder.WriteUint64(mint.Amount, binary.LittleEndian)
		if err != nil {
			return err
		}
	}
	{
		if mint.Delegate == nil {
			err = encoder.WriteUint32(0, binary.LittleEndian)
			if err != nil {
				return err
			}
			empty := base.PublicKey{}
			err = encoder.WriteBytes(empty[:], false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteUint32(1, binary.LittleEndian)
			if err != nil {
				return err
			}
			err = encoder.WriteBytes(mint.Delegate[:], false)
			if err != nil {
				return err
			}
		}
	}
	err = encoder.WriteUint8(uint8(mint.State))
	if err != nil {
		return err
	}
	{
		if mint.IsNative == nil {
			err = encoder.WriteUint32(0, binary.LittleEndian)
			if err != nil {
				return err
			}
			err = encoder.WriteUint64(0, binary.LittleEndian)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteUint32(1, binary.LittleEndian)
			if err != nil {
				return err
			}
			err = encoder.WriteUint64(*mint.IsNative, binary.LittleEndian)
			if err != nil {
				return err
			}
		}
	}
	{
		err = encoder.WriteUint64(mint.DelegatedAmount, binary.LittleEndian)
		if err != nil {
			return err
		}
	}
	{
		if mint.CloseAuthority == nil {
			err = encoder.WriteUint32(0, binary.LittleEndian)
			if err != nil {
				return err
			}
			empty := base.PublicKey{}
			err = encoder.WriteBytes(empty[:], false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteUint32(1, binary.LittleEndian)
			if err != nil {
				return err
			}
			err = encoder.WriteBytes(mint.CloseAuthority[:], false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type Multisig struct {
	// Number of signers required
	M uint8
	// Number of valid signers
	N uint8
	// Is `true` if this structure has been initialized
	IsInitialized bool
	// Signer public keys
	Signers [MAX_SIGNERS]base.PublicKey
}
