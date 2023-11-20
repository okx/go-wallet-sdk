// Copyright (c) 2020-2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc
//

package types

import "bytes"

// Transaction represents "transaction" operation
type Transaction struct {
	Manager
	Amount      N           `json:"amount"`
	Destination Address     `json:"destination"`
	Parameters  *Parameters `json:"parameters,omitempty"`
}

func (o Transaction) Kind() OpType {
	return OpTypeTransaction
}

func (o Transaction) EncodeBuffer(buf *bytes.Buffer, p *Params) error {
	buf.WriteByte(o.Kind().TagVersion(p.OperationTagsVersion))
	o.Manager.EncodeBuffer(buf, p)
	o.Amount.EncodeBuffer(buf)
	buf.Write(o.Destination.Bytes22())
	buf.WriteByte(0x0)
	return nil
}
