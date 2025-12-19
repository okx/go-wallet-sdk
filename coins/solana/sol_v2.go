package solana

import (
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/okx/go-wallet-sdk/crypto/base58"
	"github.com/okx/go-wallet-sdk/util"
)

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

type TxData struct {
	RawTx string `json:"rawTx"`
	TxId  string `json:"txId"`
}

func NewAddressFromPubkey(pubkey []byte) (addr string, err error) {
	if len(pubkey) != ed25519.PublicKeySize {
		return "", errors.New("invalid pubkey")
	}
	return base58.Encode(pubkey), nil
}

// NewPublicKeyFromBase58 decodes a base58 encoded string to a public key ensuring the length is correct
func NewPublicKeyFromBase58(s string) (common.PublicKey, error) {
	b, err := util.DecodeBase58(s)
	if err != nil {
		return common.PublicKey{}, err
	}
	if len(b) != common.PublicKeyLength {
		return common.PublicKey{}, errors.New("invalid public key length")
	}
	var publicKey common.PublicKey
	copy(publicKey[:], b)
	return publicKey, nil
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
		return types.Transaction{}, err
	}

	// handle instructions
	ixs := make([]types.Instruction, 0, len(txParams.Instructions))
	for _, ixParam := range txParams.Instructions {
		ix, err := buildInstruction(ixParam)
		if err != nil {
			return types.Transaction{}, err
		}
		ixs = append(ixs, ix)
	}

	// handle lookupTables
	lookupTables := make([]types.AddressLookupTableAccount, 0, len(txParams.LookupTables))
	for _, lookupTableParam := range txParams.LookupTables {
		lookupTable, err := buildLookupTable(lookupTableParam)
		if err != nil {
			return types.Transaction{}, err
		}
		lookupTables = append(lookupTables, lookupTable)
	}

	feePayer, err := NewPublicKeyFromBase58(txParams.FeePayer)
	if err != nil {
		return types.Transaction{}, fmt.Errorf("failed to decode fee payer: %w", err)
	}

	return types.NewTransaction(
		types.NewTransactionParam{
			Message: types.NewMessage(
				types.NewMessageParam{
					FeePayer:                   feePayer,
					Instructions:               ixs,
					RecentBlockhash:            txParams.RecentBlockHash,
					AddressLookupTableAccounts: lookupTables,
				},
			),
		},
	)
}

func NewTxFromRaw(rawTx string, encoding string) (tx types.Transaction, err error) {
	if rawTx == "" {
		return types.Transaction{}, errors.New("raw transaction cannot be empty")
	}

	var rawBytes []byte
	if encoding == "base64" {
		rawBytes, err = base64.StdEncoding.DecodeString(rawTx)
		if err != nil {
			return tx, errors.New("invalid base64 encoding")
		}
	} else if encoding == "base58" || encoding == "" {
		rawBytes, err = util.DecodeBase58(rawTx)
		if err != nil {
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

func AddSignature(tx types.Transaction, sig []byte, encoding string) (data TxData, err error) {
	if err := tx.AddSignature(sig); err != nil {
		return TxData{}, fmt.Errorf("failed to add signature: %w", err)
	}

	rawBytes, err := tx.Serialize()
	if err != nil {
		return TxData{}, fmt.Errorf("failed to serialize transaction: %w", err)
	}

	var rawTx string
	if encoding == "base64" {
		rawTx = base64.StdEncoding.EncodeToString(rawBytes)
	} else if encoding == "base58" || encoding == "" {
		rawTx = base58.Encode(rawBytes)
	} else {
		return TxData{}, errors.New("invalid encoding")
	}

	txId := base58.Encode(tx.Signatures[0])
	return TxData{
		RawTx: rawTx,
		TxId:  txId,
	}, nil
}

func buildInstruction(ixParam Instruction) (types.Instruction, error) {
	accounts := make([]types.AccountMeta, 0, len(ixParam.Keys))
	for _, k := range ixParam.Keys {
		pubKey, err := util.DecodeBase58(k.Pubkey)
		if err != nil {
			return types.Instruction{}, fmt.Errorf("failed to decode pubkey: %w", err)
		}
		accounts = append(accounts, types.AccountMeta{
			PubKey:     common.PublicKeyFromBytes(pubKey),
			IsSigner:   k.IsSigner,
			IsWritable: k.IsWritable,
		})
	}

	programId, err := util.DecodeBase58(ixParam.ProgramId)
	if err != nil {
		return types.Instruction{}, fmt.Errorf("failed to decode program id: %w", err)
	}

	data, err := util.DecodeBase58(ixParam.Data)
	if err != nil {
		return types.Instruction{}, fmt.Errorf("failed to decode data: %w", err)
	}

	return types.Instruction{
		Accounts:  accounts,
		ProgramID: common.PublicKeyFromBytes(programId),
		Data:      data,
	}, nil
}

func buildLookupTable(lookupTableParam LookupTable) (types.AddressLookupTableAccount, error) {
	addresses := make([]common.PublicKey, 0, len(lookupTableParam.AddressList))
	for _, addressParam := range lookupTableParam.AddressList {
		address, err := util.DecodeBase58(addressParam)
		if err != nil {
			return types.AddressLookupTableAccount{}, fmt.Errorf("failed to decode address: %w", err)
		}
		addresses = append(addresses, common.PublicKeyFromBytes(address))
	}

	tableAccount, err := util.DecodeBase58(lookupTableParam.TableAccount)
	if err != nil {
		return types.AddressLookupTableAccount{}, fmt.Errorf("failed to decode table account: %w", err)
	}

	return types.AddressLookupTableAccount{
		Key:       common.PublicKeyFromBytes(tableAccount),
		Addresses: addresses,
	}, nil
}
