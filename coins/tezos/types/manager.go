// Copyright (c) 2020-2023 Blockwatch Data Inc.
// Author: alex@blockwatch.cc
package types

import (
	"bytes"
	"strconv"
)

// Manager contains fields common for all manager operations
type Manager struct {
	Source       Address `json:"source"`
	Fee          N       `json:"fee"`
	Counter      N       `json:"counter"`
	GasLimit     N       `json:"gas_limit"`
	StorageLimit N       `json:"storage_limit"`
}

func (o *Manager) WithSource(addr Address) {
	o.Source = addr
}

func (o Manager) Limits() Limits {
	return Limits{
		Fee:          o.Fee.Int64(),
		GasLimit:     o.GasLimit.Int64(),
		StorageLimit: o.StorageLimit.Int64(),
	}
}

func (o *Manager) WithLimits(limits Limits) {
	o.Fee.SetInt64(limits.Fee)
	o.GasLimit.SetInt64(limits.GasLimit)
	o.StorageLimit.SetInt64(limits.StorageLimit)
}

func (o Manager) EncodeBuffer(buf *bytes.Buffer, _ *Params) error {
	buf.Write(o.Source.Bytes())
	o.Fee.EncodeBuffer(buf)
	o.Counter.EncodeBuffer(buf)
	o.GasLimit.EncodeBuffer(buf)
	o.StorageLimit.EncodeBuffer(buf)
	return nil
}

func (o Manager) EncodeJSON(buf *bytes.Buffer) error {
	buf.WriteString(`"source":`)
	buf.WriteString(strconv.Quote(o.Source.String()))
	buf.WriteString(`,"fee":`)
	buf.WriteString(strconv.Quote(o.Fee.String()))
	buf.WriteString(`,"counter":`)
	buf.WriteString(strconv.Quote(o.Counter.String()))
	buf.WriteString(`,"gas_limit":`)
	buf.WriteString(strconv.Quote(o.GasLimit.String()))
	buf.WriteString(`,"storage_limit":`)
	buf.WriteString(strconv.Quote(o.StorageLimit.String()))
	return nil
}

func (o *Manager) WithCounter(c int64) {
	o.Counter.SetInt64(c)
}

func (o Manager) GetCounter() int64 {
	return o.Counter.Int64()
}
