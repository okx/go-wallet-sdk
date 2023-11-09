package types

import (
	"encoding/binary"
	"github.com/okx/go-wallet-sdk/coins/waves/crypto"
)

const (
	// https://sourcegraph.com/github.com/wavesplatform/gowaves@master/-/blob/pkg/proto/transactions.go
	transferLen = crypto.PublicKeySize + 1 + 1 + 8 + 8 + 8 + 2
)

type Transfer struct {
	SenderPK    crypto.PublicKey `json:"senderPublicKey"`
	AmountAsset OptionalAsset    `json:"assetId"`
	FeeAsset    OptionalAsset    `json:"feeAssetId"`
	Timestamp   uint64           `json:"timestamp,omitempty"`
	Amount      uint64           `json:"amount"`
	Fee         uint64           `json:"fee"`
	Recipient   Recipient        `json:"recipient"`
	Attachment  Attachment       `json:"attachment,omitempty"`
}

func (tr *Transfer) marshalBinary() ([]byte, error) {
	p := 0
	aal := 0
	if tr.AmountAsset.Present {
		aal += crypto.DigestSize
	}
	fal := 0
	if tr.FeeAsset.Present {
		fal += crypto.DigestSize
	}
	rb, err := tr.Recipient.MarshalBinary()
	if err != nil {
		return nil, err
	}
	rl := len(rb)
	att := tr.Attachment
	atl := att.Size()
	buf := make([]byte, transferLen+aal+fal+atl+rl)
	copy(buf[p:], tr.SenderPK[:])
	p += crypto.PublicKeySize
	aab, err := tr.AmountAsset.MarshalBinary()
	if err != nil {
		return nil, err
	}
	copy(buf[p:], aab)
	p += 1 + aal
	fab, err := tr.FeeAsset.MarshalBinary()
	if err != nil {
		return nil, err
	}
	copy(buf[p:], fab)
	p += 1 + fal
	binary.BigEndian.PutUint64(buf[p:], tr.Timestamp)
	p += 8
	binary.BigEndian.PutUint64(buf[p:], tr.Amount)
	p += 8
	binary.BigEndian.PutUint64(buf[p:], tr.Fee)
	p += 8
	copy(buf[p:], rb)
	p += rl
	attBytes, err := att.Bytes()
	if err != nil {
		return nil, err
	}
	PutBytesWithUInt16Len(buf[p:], attBytes)
	return buf, nil
}

// TransferWithSig transaction to transfer any token from one account to another. Version 1.
type TransferWithSig struct {
	Type      TransactionType   `json:"type"`
	Version   byte              `json:"version,omitempty"`
	ID        *crypto.Digest    `json:"id,omitempty"`
	Signature *crypto.Signature `json:"signature,omitempty"`
	Transfer
}

func (tx *TransferWithSig) Sign(scheme Scheme, secretKey crypto.SecretKey) error {
	b, err := MarshalTxBody(scheme, tx)
	if err != nil {
		return err
	}
	s, err := crypto.Sign(secretKey, b)
	if err != nil {
		return err
	}
	tx.Signature = &s
	d, err := crypto.FastHash(b)
	if err != nil {
		return err
	}
	tx.ID = &d
	return nil
}

func (tx *TransferWithSig) BodyMarshalBinary() ([]byte, error) {
	b, err := tx.Transfer.marshalBinary()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 1+len(b))
	buf[0] = byte(tx.Type)
	copy(buf[1:], b)
	return buf, nil
}

func MarshalTxBody(scheme Scheme, tx *TransferWithSig) ([]byte, error) {
	return tx.BodyMarshalBinary()
}
