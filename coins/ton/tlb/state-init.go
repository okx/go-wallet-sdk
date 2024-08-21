/**
Author： https://github.com/xssnick/tonutils-go
*/

package tlb

import (
	"github.com/okx/go-wallet-sdk/coins/ton/tvm/cell"
)

type TickTock struct {
	Tick bool `tlb:"bool"`
	Tock bool `tlb:"bool"`
}

type StateInit struct {
	Depth    *uint64          `tlb:"maybe ## 5"`
	TickTock *TickTock        `tlb:"maybe ."`
	Code     *cell.Cell       `tlb:"maybe ^"`
	Data     *cell.Cell       `tlb:"maybe ^"`
	Lib      *cell.Dictionary `tlb:"dict 256"`
}
