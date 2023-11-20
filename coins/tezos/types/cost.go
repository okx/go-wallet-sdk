/**
The MIT License (MIT) Copyright (c) 2021-2023 Blockwatch Data Inc.
*/

package types

// Limits represents all resource limits defined for an operation in Tezos.
type Limits struct {
	Fee          int64
	GasLimit     int64
	StorageLimit int64
}

// Add adds two limits z = x + y and returns the sum z without changing any of the inputs.
func (x Limits) Add(y Limits) Limits {
	x.Fee += y.Fee
	x.GasLimit += y.GasLimit
	x.StorageLimit += y.StorageLimit
	return x
}
