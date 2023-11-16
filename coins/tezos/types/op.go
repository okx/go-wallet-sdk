/**
The MIT License (MIT) Copyright (c) 2021-2023 Blockwatch Data Inc.
*/

package types

import (
	"bytes"
	"fmt"
)

type OpType byte

// enums are allocated in chronological order
const (
	OpTypeInvalid                      OpType = iota
	OpTypeActivateAccount                     // 1
	OpTypeDoubleBakingEvidence                // 2
	OpTypeDoubleEndorsementEvidence           // 3
	OpTypeSeedNonceRevelation                 // 4
	OpTypeTransaction                         // 5
	OpTypeOrigination                         // 6
	OpTypeDelegation                          // 7
	OpTypeReveal                              // 8
	OpTypeEndorsement                         // 9
	OpTypeProposals                           // 10
	OpTypeBallot                              // 11
	OpTypeFailingNoop                         // 12 v009
	OpTypeEndorsementWithSlot                 // 13 v009
	OpTypeRegisterConstant                    // 14 v011
	OpTypePreendorsement                      // 15 v012
	OpTypeDoublePreendorsementEvidence        // 16 v012
	OpTypeSetDepositsLimit                    // 17 v012
	OpTypeToruOrigination                     // 18 v013
	OpTypeToruSubmitBatch                     // 19 v013
	OpTypeToruCommit                          // 20 v013
	OpTypeToruReturnBond                      // 21 v013
	OpTypeToruFinalizeCommitment              // 22 v013
	OpTypeToruRemoveCommitment                // 23 v013
	OpTypeToruRejection                       // 24 v013
	OpTypeToruDispatchTickets                 // 25 v013
	OpTypeTransferTicket                      // 26 v013
	OpTypeScruOriginate                       // 27 v013
	OpTypeScruAddMessages                     // 28 v013
	OpTypeScruCement                          // 29 v013
	OpTypeScruPublish                         // 30 v013
)

const (
	EmmyBlockWatermark                byte = 1 // deprecated
	EmmyEndorsementWatermark               = 2 // deprecated
	OperationWatermark                     = 3
	TenderbakeBlockWatermark               = 11
	TenderbakePreendorsementWatermark      = 12
	TenderbakeEndorsementWatermark         = 13
)

var (
	opTypeStrings = map[OpType]string{
		OpTypeInvalid:         "",
		OpTypeActivateAccount: "activate_account",
		OpTypeTransaction:     "transaction",
		OpTypeOrigination:     "origination",
		OpTypeDelegation:      "delegation",
		OpTypeReveal:          "reveal",
	}
	// before Babylon v005
	opTagV0 = map[OpType]byte{
		OpTypeEndorsement:               0,
		OpTypeSeedNonceRevelation:       1,
		OpTypeDoubleEndorsementEvidence: 2,
		OpTypeDoubleBakingEvidence:      3,
		OpTypeActivateAccount:           4,
		OpTypeProposals:                 5,
		OpTypeBallot:                    6,
		OpTypeReveal:                    7,
		OpTypeTransaction:               8,
		OpTypeOrigination:               9,
		OpTypeDelegation:                10,
	}
	// Babylon v005 until Hangzhou v011
	opTagV1 = map[OpType]byte{
		OpTypeEndorsement:               0,
		OpTypeSeedNonceRevelation:       1,
		OpTypeDoubleEndorsementEvidence: 2,
		OpTypeDoubleBakingEvidence:      3,
		OpTypeActivateAccount:           4,
		OpTypeProposals:                 5,
		OpTypeBallot:                    6,
		OpTypeReveal:                    107, // v005
		OpTypeTransaction:               108, // v005
		OpTypeOrigination:               109, // v005
		OpTypeDelegation:                110, // v005
		OpTypeEndorsementWithSlot:       10,  // v009
		OpTypeFailingNoop:               17,  // v009
		OpTypeRegisterConstant:          111, // v011
	}
	// Ithaca v012 and up
	opTagV2 = map[OpType]byte{
		OpTypeSeedNonceRevelation:          1,
		OpTypeDoubleEndorsementEvidence:    2,
		OpTypeDoubleBakingEvidence:         3,
		OpTypeActivateAccount:              4,
		OpTypeProposals:                    5,
		OpTypeBallot:                       6,
		OpTypeReveal:                       107, // v005
		OpTypeTransaction:                  108, // v005
		OpTypeOrigination:                  109, // v005
		OpTypeDelegation:                   110, // v005
		OpTypeFailingNoop:                  17,  // v009
		OpTypeRegisterConstant:             111, // v011
		OpTypePreendorsement:               20,  // v012
		OpTypeEndorsement:                  21,  // v012
		OpTypeDoublePreendorsementEvidence: 7,   // v012
		OpTypeSetDepositsLimit:             112, // v012
		OpTypeToruOrigination:              150, // v013
		OpTypeToruSubmitBatch:              151, // v013
		OpTypeToruCommit:                   152, // v013
		OpTypeToruReturnBond:               153, // v013
		OpTypeToruFinalizeCommitment:       154, // v013
		OpTypeToruRemoveCommitment:         155, // v013
		OpTypeToruRejection:                156, // v013
		OpTypeToruDispatchTickets:          157, // v013
		OpTypeTransferTicket:               158, // v013
		OpTypeScruOriginate:                200, // v013
		OpTypeScruAddMessages:              201, // v013
		OpTypeScruCement:                   202, // v013
		OpTypeScruPublish:                  203, // v013
	}
)

func (t OpType) String() string {
	return opTypeStrings[t]
}

func (t OpType) TagVersion(ver int) byte {
	var (
		tag byte
		ok  bool
	)
	switch ver {
	case 0:
		tag, ok = opTagV0[t]
	case 1:
		tag, ok = opTagV1[t]
	default:
		tag, ok = opTagV2[t]
	}
	if !ok {
		return 255
	}
	return tag
}

// Op is a container used to collect, serialize and sign Tezos operations.
type Op struct {
	Branch    BlockHash    `json:"branch"`    // used for TTL handling
	Contents  []Operation  `json:"contents"`  // non-zero list of transactions
	Signature Signature    `json:"signature"` // added during the lifecycle
	ChainId   *ChainIdHash `json:"-"`         // optional, used for remote signing only
	TTL       int64        `json:"-"`         // optional, specify TTL in blocks
	Params    *Params      `json:"-"`         // optional, define protocol to encode for
	Source    Address      `json:"-"`         // optional, used as manager/sender
}

type Operation interface {
	Kind() OpType
	Limits() Limits
	GetCounter() int64
	WithSource(Address)
	WithCounter(int64)
	WithLimits(Limits)
	EncodeBuffer(buf *bytes.Buffer, p *Params) error
}

func NewOp() *Op {
	return &Op{
		Params: DefaultParams,
		TTL:    DefaultParams.MaxOperationsTTL - 2, // Ithaca recommendation
	}
}

func NewJakartanetOp() *Op {
	return &Op{
		Params: JakartanetParams,
		TTL:    JakartanetParams.MaxOperationsTTL - 2, // Ithaca recommendation
	}
}

func (o *Op) Bytes() []byte {
	if len(o.Contents) == 0 || !o.Branch.IsValid() {
		return nil
	}
	p := o.Params
	if p == nil {
		p = DefaultParams
	}
	buf := bytes.NewBuffer(nil)
	buf.Write(o.Branch.Bytes())
	for _, v := range o.Contents {
		_ = v.EncodeBuffer(buf, p)
	}
	if o.Contents[0].Kind() != OpTypeEndorsementWithSlot {
		if o.Signature.IsValid() {
			buf.Write(o.Signature.Data) // raw, without type (!)
		}
	}
	return buf.Bytes()
}

// WithTransfer adds a simple value transfer transaction to the contents list.
func (o *Op) WithTransfer(to Address, amount int64) *Op {
	o.Contents = append(o.Contents, &Transaction{
		Manager: Manager{
			Source:  o.Source,
			Counter: 0,
		},
		Amount:      N(amount),
		Destination: to,
	})
	return o
}

// WithDelegation adds a delegation transaction to the contents list.
func (o *Op) WithDelegation(to Address) *Op {
	o.Contents = append(o.Contents, &Delegation{
		Manager: Manager{
			Source:  o.Source,
			Counter: 0,
		},
		Delegate: to,
	})
	return o
}

// WithUndelegation adds a delegation transaction that resets the callers baker to null
// to the contents list.
func (o *Op) WithUndelegation() *Op {
	o.Contents = append(o.Contents, &Delegation{
		Manager: Manager{
			Source:  o.Source,
			Counter: 0,
		},
	})
	return o
}

// WithSource sets the source for all manager operations to addr.
func (o *Op) WithSource(addr Address) *Op {
	for _, v := range o.Contents {
		v.WithSource(addr)
	}
	o.Source = addr
	return o
}

// WithParams defines the protocol and other chain configuration params for which
// the operation will be encoded. If unset, defaults to tezos.DefaultParams.
func (o *Op) WithParams(p *Params) *Op {
	o.Params = p
	return o
}

// WithContentsFront adds a Tezos operation to the front of the contents list.
func (o *Op) WithContentsFront(op Operation) *Op {
	o.Contents = append([]Operation{op}, o.Contents...)
	return o
}

func (o *Op) WithSignature(sig Signature) *Op {
	o.Signature = sig
	return o
}

func (o *Op) WithChainId(id ChainIdHash) *Op {
	clone := id.Clone()
	o.ChainId = &clone
	return o
}

func (o *Op) WithTTL(n int64) *Op {
	if n > o.Params.MaxOperationsTTL {
		n = o.Params.MaxOperationsTTL - 2 // Ithaca adjusted
	} else if n < 0 {
		n = 1
	}
	o.TTL = n
	return o
}

func (o Op) Limits() Limits {
	var l Limits
	for _, v := range o.Contents {
		l = l.Add(v.Limits())
	}
	return l
}

func (o Op) NeedCounter() bool {
	for _, v := range o.Contents {
		if v.GetCounter() == 0 {
			return true
		}
	}
	return false
}

func (o *Op) WithBranch(hash BlockHash) *Op {
	o.Branch = hash
	return o
}

func (o *Op) Digest() []byte {
	d := Digest(o.WatermarkedBytes())
	return d[:]
}

func (o *Op) WatermarkedBytes() []byte {
	if len(o.Contents) == 0 || !o.Branch.IsValid() {
		return nil
	}
	p := o.Params
	if p == nil {
		p = DefaultParams
	}
	buf := bytes.NewBuffer(nil)
	switch o.Contents[0].Kind() {
	case OpTypeEndorsement, OpTypeEndorsementWithSlot:
		if p.OperationTagsVersion < 2 {
			buf.WriteByte(EmmyEndorsementWatermark)
		} else {
			buf.WriteByte(TenderbakeEndorsementWatermark)
		}
		if o.ChainId != nil {
			buf.Write(o.ChainId.Bytes())
		}
	case OpTypePreendorsement:
		buf.WriteByte(TenderbakePreendorsementWatermark)
		if o.ChainId != nil {
			buf.Write(o.ChainId.Bytes())
		}
	default:
		buf.WriteByte(OperationWatermark)
	}
	buf.Write(o.Branch.Bytes())
	for _, v := range o.Contents {
		_ = v.EncodeBuffer(buf, p)
	}
	return buf.Bytes()
}

func (o *Op) Sign(key PrivateKey) error {
	if !o.Branch.IsValid() {
		return fmt.Errorf("tezos: missing branch")
	}
	if len(o.Contents) == 0 {
		return fmt.Errorf("tezos: empty operation contents")
	}
	sig, err := key.Sign(o.Digest())
	if err != nil {
		return err
	}
	o.Signature = sig
	return nil
}
