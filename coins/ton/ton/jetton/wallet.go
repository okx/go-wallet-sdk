/**
Author： https://github.com/xssnick/tonutils-go
*/

package jetton

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/ton/address"
	"github.com/okx/go-wallet-sdk/coins/ton/tlb"
	"github.com/okx/go-wallet-sdk/coins/ton/tvm/cell"
)

type TransferPayload struct {
	_                   tlb.Magic        `tlb:"#0f8a7ea5"`
	QueryID             uint64           `tlb:"## 64"`
	Amount              tlb.Coins        `tlb:"."`
	Destination         *address.Address `tlb:"addr"`
	ResponseDestination *address.Address `tlb:"addr"`
	CustomPayload       *cell.Cell       `tlb:"maybe ^"`
	ForwardTONAmount    tlb.Coins        `tlb:"."`
	ForwardPayload      *cell.Cell       `tlb:"either . ^"`
}

type BurnPayload struct {
	_                   tlb.Magic        `tlb:"#595f07bc"`
	QueryID             uint64           `tlb:"## 64"`
	Amount              tlb.Coins        `tlb:"."`
	ResponseDestination *address.Address `tlb:"addr"`
	CustomPayload       *cell.Cell       `tlb:"maybe ^"`
}

type WalletClient struct {
	master *Client
	addr   *address.Address
}

func (c *WalletClient) Address() *address.Address {
	return c.addr
}

// Deprecated: use BuildTransferPayloadV2
func (c *WalletClient) BuildTransferPayload(to *address.Address, amountCoins, amountForwardTON tlb.Coins, payloadForward *cell.Cell) (*cell.Cell, error) {
	return c.BuildTransferPayloadV2(to, to, amountCoins, amountForwardTON, payloadForward, nil, 0)
}

func (c *WalletClient) BuildTransferPayloadV2(to, responseTo *address.Address, amountCoins, amountForwardTON tlb.Coins, payloadForward, customPayload *cell.Cell, rndU64 uint64) (*cell.Cell, error) {
	if payloadForward == nil {
		payloadForward = cell.BeginCell().EndCell()
	}

	rnd := rndU64
	if rnd == 0 {
		buf := make([]byte, 8)
		if _, err := rand.Read(buf); err != nil {
			return nil, err
		}
		rnd = binary.LittleEndian.Uint64(buf)
	}

	body, err := tlb.ToCell(TransferPayload{
		QueryID:             rnd,
		Amount:              amountCoins,
		Destination:         to,
		ResponseDestination: responseTo,
		CustomPayload:       customPayload,
		ForwardTONAmount:    amountForwardTON,
		ForwardPayload:      payloadForward,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to convert TransferPayload to cell: %w", err)
	}

	return body, nil
}

func (c *WalletClient) BuildBurnPayload(amountCoins tlb.Coins, notifyAddr *address.Address) (*cell.Cell, error) {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}
	rnd := binary.LittleEndian.Uint64(buf)

	body, err := tlb.ToCell(BurnPayload{
		QueryID:             rnd,
		Amount:              amountCoins,
		ResponseDestination: notifyAddr,
		CustomPayload:       nil,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to convert BurnPayload to cell: %w", err)
	}

	return body, nil
}
