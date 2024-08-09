/**
Author： https://github.com/xssnick/tonutils-go
*/

package wallet

import (
	"context"
	"errors"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/ton/tlb"
	"github.com/okx/go-wallet-sdk/coins/ton/tvm/cell"
	"time"
)

const _VenomV3CodeHex = "B5EE9C720101010100710000DEFF0020DD2082014C97BA218201339CBAB19F71B0ED44D0D31FD31F31D70BFFE304E0A4F2608308D71820D31FD31FD31FF82313BBF263ED44D0D31FD31FD3FFD15132BAF2A15144BAF2A204F901541055F910F2A3F8009320D74A96D307D402FB00E8D101A4C8CB1FCB1FCBFFC9ED54"

type SpecVenomV3 struct {
	SpecRegular
	SpecSeqno
}

func (s *SpecVenomV3) BuildMessage(ctx context.Context, messages []*Message) (*cell.Cell, error) {
	if len(messages) > 4 {
		return nil, errors.New("for this type of wallet max 4 messages can be sent in the same time")
	}

	if s.seqnoFetcher == nil {
		return nil, errors.New("unable to get seq")
	}
	seq, err := s.seqnoFetcher(ctx, s.wallet.subwallet)
	if err != nil {
		return nil, err
	}
	seqno := uint64(seq)
	unix := uint64(s.expireAt)
	if unix == 0 {
		unix = uint64(timeNow().Add(600 * time.Second).UTC().Unix())
	}
	payload := cell.BeginCell().MustStoreUInt(uint64(s.wallet.subwallet), 32).
		MustStoreUInt(unix, 32).
		MustStoreUInt(seqno, 32)

	for i, message := range messages {
		intMsg, err := tlb.ToCell(message.InternalMessage)
		if err != nil {
			return nil, fmt.Errorf("failed to convert internal message %d to cell: %w", i, err)
		}

		payload.MustStoreUInt(uint64(message.Mode), 8).MustStoreRef(intMsg)
	}

	sign := payload.EndCell().Sign(s.wallet.key)
	msg := cell.BeginCell().MustStoreSlice(sign, 512).MustStoreBuilder(payload).EndCell()

	return msg, nil
}
