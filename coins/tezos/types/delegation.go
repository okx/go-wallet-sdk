// Copyright (c) 2020-2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc
//

package types

import (
	"bytes"
)

type Delegation struct {
	Manager
	Delegate Address `json:"delegate"`
}

func (o Delegation) Kind() OpType {
	return OpTypeDelegation
}

func (o Delegation) EncodeBuffer(buf *bytes.Buffer, p *Params) error {
	buf.WriteByte(o.Kind().TagVersion(p.OperationTagsVersion))
	o.Manager.EncodeBuffer(buf, p)
	if o.Delegate.IsValid() {
		buf.WriteByte(0xff)
		buf.Write(o.Delegate.Bytes())
	} else {
		buf.WriteByte(0x0)
	}
	return nil
}
