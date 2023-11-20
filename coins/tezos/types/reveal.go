// Copyright (c) 2020-2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc
package types

import (
	"bytes"
)

type Reveal struct {
	Manager
	PublicKey Key `json:"public_key"`
}

func (o Reveal) Kind() OpType {
	return OpTypeReveal
}

func (o Reveal) EncodeBuffer(buf *bytes.Buffer, p *Params) error {
	buf.WriteByte(o.Kind().TagVersion(p.OperationTagsVersion))
	o.Manager.EncodeBuffer(buf, p)
	buf.Write(o.PublicKey.Bytes())
	return nil
}
