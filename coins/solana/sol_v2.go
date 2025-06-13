package solana

import (
	"crypto/ed25519"
	"encoding/json"
	"errors"
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

type TeeTxParams struct {
	FeePayer        string        `json:"feePayer"`
	RecentBlockHash string        `json:"recentBlockHash"`
	Instructions    []Instruction `json:"instructions"`
	LookupTables    []LookupTable `json:"lookupTables"`
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

func NewTxFromJson(txJson string) (tx types.Transaction, err error) {
	var txParams TeeTxParams
	err = json.Unmarshal([]byte(txJson), &txParams)
	if err != nil {
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

func NewTxFromRaw(rawTx string) (tx types.Transaction, err error) {
	tx, err = types.TransactionDeserialize(base58.Decode(rawTx))
	return tx, err
}

func GetSigningData(tx types.Transaction) (data []byte, err error) {
	return tx.Message.Serialize()
}

type TxData struct {
	RawTx string
	TxId  string
}

func AddSignature(tx types.Transaction, sig []byte) (data TxData, err error) {
	err = tx.AddSignature(sig)
	if err != nil {
		return TxData{}, err
	}

	rawTx, err := tx.Serialize()
	if err != nil {
		return TxData{}, err
	}

	rawTxBase58 := base58.Encode(rawTx)

	txHash := base58.Encode(tx.Signatures[0])

	return TxData{
		RawTx: rawTxBase58,
		TxId:  txHash,
	}, nil
}
