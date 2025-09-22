package solana

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/okx/go-wallet-sdk/crypto/base58"
)

func NewAddressFromPubkey(pubkey []byte) (addr string, err error) {
	if len(pubkey) != ed25519.PublicKeySize {
		return "", errors.New("invalid pubkey")
	}
	return base58.Encode(pubkey), nil
}

// SolanaTxParams represents the parameters needed to construct a Solana transaction
type SolanaTxParams struct {
	FeePayer        string        `json:"feePayer"`        // The account that will pay for the transaction
	RecentBlockHash string        `json:"recentBlockHash"` // Recent block hash for transaction validity
	Instructions    []Instruction `json:"instructions"`    // List of instructions to execute
	LookupTables    []LookupTable `json:"lookupTables"`    // Address lookup tables for account compression
}

type Instruction struct {
	ProgramId string       `json:"programId"`
	Keys      []AccountKey `json:"keys"`
	Data      string       `json:"data"`
}

type AccountKey struct {
	Pubkey     string `json:"pubkey"`
	IsSigner   bool   `json:"isSigner"`
	IsWritable bool   `json:"isWritable"`
}

type LookupTable struct {
	TableAccount string   `json:"tableAccount"`
	AddressList  []string `json:"addressList"`
}

func ValidateTxParams(params *SolanaTxParams) error {
	if params.FeePayer == "" {
		return errors.New("feePayer cannot be empty")
	}
	if params.RecentBlockHash == "" {
		return errors.New("recentBlockHash cannot be empty")
	}
	if len(params.Instructions) == 0 {
		return errors.New("instructions cannot be empty")
	}

	return nil
}

func NewTxFromParams(txParams SolanaTxParams) (tx types.Transaction, err error) {
	if err := ValidateTxParams(&txParams); err != nil {
		return tx, err
	}

	// handle instructions
	var ixs []types.Instruction
	for _, ixParam := range txParams.Instructions {
		var accounts []types.AccountMeta
		for _, keyParam := range ixParam.Keys {
			accounts = append(accounts, types.AccountMeta{
				PubKey:     common.PublicKeyFromString(keyParam.Pubkey),
				IsSigner:   keyParam.IsSigner,
				IsWritable: keyParam.IsWritable,
			})
		}

		ix := types.Instruction{
			Accounts:  accounts,
			ProgramID: common.PublicKeyFromString(ixParam.ProgramId),
			Data:      base58.Decode(ixParam.Data),
		}
		ixs = append(ixs, ix)
	}

	// handle lookuptables
	var lookupTables []types.AddressLookupTableAccount
	for _, lookupTableParam := range txParams.LookupTables {
		var addresses []common.PublicKey
		for _, addressParam := range lookupTableParam.AddressList {
			addresses = append(addresses, common.PublicKeyFromString(addressParam))
		}

		lookupTables = append(lookupTables, types.AddressLookupTableAccount{
			Key:       common.PublicKeyFromString(lookupTableParam.TableAccount),
			Addresses: addresses,
		})
	}

	tx, err = types.NewTransaction(
		types.NewTransactionParam{
			Message: types.NewMessage(
				types.NewMessageParam{
					FeePayer:                   common.PublicKeyFromString(txParams.FeePayer),
					Instructions:               ixs,
					RecentBlockhash:            txParams.RecentBlockHash,
					AddressLookupTableAccounts: lookupTables,
				},
			),
		},
	)

	return tx, err
}

func NewTxFromRaw(rawTx string, encoding string) (tx types.Transaction, err error) {
	if rawTx == "" {
		return tx, errors.New("raw transaction cannot be empty")
	}

	var rawBytes []byte
	if encoding == "base64" {
		rawBytes, err = base64.StdEncoding.DecodeString(rawTx)
		if err != nil {
			return tx, errors.New("invalid base64 encoding")
		}
	} else if encoding == "base58" || encoding == "" {
		rawBytes = base58.Decode(rawTx)
		if len(rawBytes) == 0 {
			return tx, errors.New("invalid base58 encoding")
		}
	} else {
		return tx, errors.New("invalid encoding")
	}

	tx, err = types.TransactionDeserialize(rawBytes)
	if err != nil {
		return tx, fmt.Errorf("failed to deserialize transaction: %w", err)
	}

	return tx, nil
}

func GetSigningData(tx types.Transaction) (data []byte, err error) {
	return tx.Message.Serialize()
}

type TxData struct {
	RawTx string
	TxId  string
}

func AddSignature(tx types.Transaction, sig []byte, encoding string) (data TxData, err error) {
	err = tx.AddSignature(sig)
	if err != nil {
		return TxData{}, err
	}

	rawBytes, err := tx.Serialize()
	if err != nil {
		return TxData{}, err
	}

	var rawTx string
	if encoding == "base64" {
		rawTx = base64.StdEncoding.EncodeToString(rawBytes)
	} else if encoding == "base58" || encoding == "" {
		rawTx = base58.Encode(rawBytes)
	} else {
		return TxData{}, errors.New("invalid encoding")
	}

	txHash := base58.Encode(tx.Signatures[0])

	return TxData{
		RawTx: rawTx,
		TxId:  txHash,
	}, nil
}
