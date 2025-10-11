/*
*
MIT License

Copyright (c) 2018 WavesPlatform
*/
package types

import (
	"bytes"
	"encoding/json"
	"github.com/emresenyuva/go-wallet-sdk/crypto/base58"
)

type TransactionType byte

// All transaction types supported.
const (
	GenesisTransaction          TransactionType = iota + 1 // 1 - Genesis transaction
	PaymentTransaction                                     // 2 - Payment transaction
	IssueTransaction                                       // 3 - Issue transaction
	TransferTransaction                                    // 4 - Transfer transaction
	ReissueTransaction                                     // 5 - Reissue transaction
	BurnTransaction                                        // 6 - Burn transaction
	ExchangeTransaction                                    // 7 - Exchange transaction
	LeaseTransaction                                       // 8 - Lease transaction
	LeaseCancelTransaction                                 // 9 - LeaseCancel transaction
	CreateAliasTransaction                                 // 10 - CreateAlias transaction
	MassTransferTransaction                                // 11 - MassTransfer transaction
	DataTransaction                                        // 12 - Data transaction
	SetScriptTransaction                                   // 13 - SetScript transaction
	SponsorshipTransaction                                 // 14 - Sponsorship transaction
	SetAssetScriptTransaction                              // 15 - SetAssetScript transaction
	InvokeScriptTransaction                                // 16 - InvokeScript transaction
	UpdateAssetInfoTransaction                             // 17 - UpdateAssetInfoTransaction
	EthereumMetamaskTransaction                            // 18 - EthereumMetamaskTransaction is a transaction which is received from metamask
	InvokeExpressionTransaction                            // 19 - InvokeExpressionTransaction
)

type Attachment []byte

func (a Attachment) MarshalJSON() ([]byte, error) {
	return json.Marshal(base58.Encode(a))
}

func (a Attachment) Size() int {
	return len(a)
}

func (a Attachment) Bytes() ([]byte, error) {
	return a, nil
}

func NewAttachmentFromBase58(s string) Attachment {
	return base58.Decode(s)
}

type Scheme byte

// B58Bytes represents bytes as Base58 string in JSON
type B58Bytes []byte

// String represents underlying bytes as Base58 string
func (b B58Bytes) String() string {
	return base58.Encode(b)
}

// MarshalJSON writes B58Bytes Value as JSON string
func (b B58Bytes) MarshalJSON() ([]byte, error) {
	return ToBase58JSON(b), nil
}

func ToBase58JSON(b []byte) []byte {
	s := base58.Encode(b)
	var sb bytes.Buffer
	sb.Grow(2 + len(s))
	sb.WriteRune('"')
	sb.WriteString(s)
	sb.WriteRune('"')
	return sb.Bytes()
}
