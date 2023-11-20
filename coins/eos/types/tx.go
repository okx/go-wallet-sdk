package types

import (
	"bytes"
	"compress/flate"
	"compress/zlib"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/eoscanada/eos-go/ecc"
	"time"
)

// TxOptions represents options you want to pass to the transaction
// you're sending.
type TxOptions struct {
	ChainID          Checksum256 // If specified, we won't hit the API to fetch it
	HeadBlockID      Checksum256 // If provided, don't hit API to fetch it.  This allows offline transaction signing.
	MaxNetUsageWords uint32
	DelaySecs        uint32
	MaxCPUUsageMS    uint8 // If you want to override the CPU usage (in counts of 1024)
	//ExtraKCPUUsage uint32 // If you want to *add* some CPU usage to the estimated amount (in counts of 1024)
	Compress   CompressionType
	Expiration time.Duration
}

// TransactionHeader Tx is an EOS transaction.
type TransactionHeader struct {
	Expiration     JSONTime `json:"expiration"`
	RefBlockNum    uint16   `json:"ref_block_num"`
	RefBlockPrefix uint32   `json:"ref_block_prefix"`

	MaxNetUsageWords Varuint32 `json:"max_net_usage_words"`
	MaxCPUUsageMS    uint8     `json:"max_cpu_usage_ms"`
	DelaySec         Varuint32 `json:"delay_sec"` // number of secs to delay, making it cancellable for that duration
}

type Extension struct {
	Type uint16
	Data []byte
}

type Transaction struct {
	TransactionHeader

	ContextFreeActions []*Action    `json:"context_free_actions"`
	Actions            []*Action    `json:"actions"`
	Extensions         []*Extension `json:"transaction_extensions"`
}

func (tx *Transaction) Fill(headBlockID Checksum256, delaySecs, maxNetUsageWords uint32,
	maxCPUUsageMS uint8, expiration time.Duration) {
	tx.setRefBlock(headBlockID)

	if tx.ContextFreeActions == nil {
		tx.ContextFreeActions = make([]*Action, 0, 0)
	}
	if tx.Extensions == nil {
		tx.Extensions = make([]*Extension, 0, 0)
	}

	tx.MaxNetUsageWords = Varuint32(maxNetUsageWords)
	tx.MaxCPUUsageMS = maxCPUUsageMS
	tx.DelaySec = Varuint32(delaySecs)

	tx.SetExpiration(30 * time.Second)
	if expiration > 30*time.Second && expiration < 3*time.Minute {
		tx.SetExpiration(expiration)
	}
}

func (tx *Transaction) setRefBlock(blockID []byte) {
	if len(blockID) == 0 {
		return
	}
	tx.RefBlockNum = uint16(binary.BigEndian.Uint32(blockID[:4]))
	tx.RefBlockPrefix = binary.LittleEndian.Uint32(blockID[8:16])
}

// SetExpiration sets the expiration of the transaction.
func (tx *Transaction) SetExpiration(in time.Duration) {
	tx.Expiration = JSONTime{time.Now().UTC().Add(in)}
}

type SignedTransaction struct {
	*Transaction

	Signatures      []ecc.Signature `json:"signatures"`
	ContextFreeData [][]byte        `json:"context_free_data"`

	packed *PackedTransaction
}

func NewSignedTransaction(tx *Transaction) *SignedTransaction {
	return &SignedTransaction{
		Transaction:     tx,
		Signatures:      make([]ecc.Signature, 0),
		ContextFreeData: make([][]byte, 0),
	}
}

func (s *SignedTransaction) String() string {
	data, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func (s *SignedTransaction) PackedTransactionAndCFD() ([]byte, []byte, error) {
	rawtrx, err := MarshalBinary(s.Transaction)
	if err != nil {
		return nil, nil, err
	}

	var rawcfd []byte
	if len(s.ContextFreeData) > 0 {
		rawcfd, err = MarshalBinary(s.ContextFreeData)
		if err != nil {
			return nil, nil, err
		}
	}

	return rawtrx, rawcfd, nil
}

type CompressionType uint8

const (
	CompressionNone = CompressionType(iota)
	CompressionZlib
)

func (c CompressionType) String() string {
	switch c {
	case CompressionNone:
		return "none"
	case CompressionZlib:
		return "zlib"
	default:
		return ""
	}
}

func (s *SignedTransaction) Pack(compression CompressionType) (*PackedTransaction, error) {
	rawtrx, rawcfd, err := s.PackedTransactionAndCFD()
	if err != nil {
		return nil, err
	}

	switch compression {
	case CompressionZlib:
		var trx bytes.Buffer
		var cfd bytes.Buffer

		// Compress Trx
		writer, _ := zlib.NewWriterLevel(&trx, flate.BestCompression) // can only fail if invalid `level`..
		writer.Write(rawtrx)                                          // ignore error, could only bust memory
		err = writer.Close()
		if err != nil {
			return nil, fmt.Errorf("tx writer close %s", err)
		}
		rawtrx = trx.Bytes()

		// Compress ContextFreeData
		writer, _ = zlib.NewWriterLevel(&cfd, flate.BestCompression) // can only fail if invalid `level`..
		writer.Write(rawcfd)                                         // ignore errors, memory errors only
		err = writer.Close()
		if err != nil {
			return nil, fmt.Errorf("cfd writer close %s", err)
		}
		rawcfd = cfd.Bytes()
	}

	packed := &PackedTransaction{
		Signatures:            s.Signatures,
		Compression:           compression,
		PackedContextFreeData: rawcfd,
		PackedTransaction:     rawtrx,
	}

	return packed, nil
}

// PackedTransaction represents a fully packed transaction, with
// signatures, and all. They circulate like that on the P2P net, and
// that's how they are stored.
type PackedTransaction struct {
	Signatures            []ecc.Signature `json:"signatures"`
	Compression           CompressionType `json:"compression"` // in C++, it's an enum, not sure how it Binary-marshals..
	PackedContextFreeData HexBytes        `json:"packed_context_free_data"`
	PackedTransaction     HexBytes        `json:"packed_trx"`
}

func (p *PackedTransaction) ID() (Checksum256, error) {
	h := sha256.New()
	_, _ = h.Write(p.PackedTransaction)
	return h.Sum(nil), nil
}
