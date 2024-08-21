/**
Author： https://github.com/xssnick/tonutils-go
*/

package wallet

import (
	"context"
	"github.com/okx/go-wallet-sdk/coins/ton/tvm/cell"
)

type RegularBuilder interface {
	BuildMessage(ctx context.Context, isInitialized bool /*_ *ton.BlockIDExt,*/, messages []*Message) (*cell.Cell, error)
}

type SpecRegular struct {
	wallet *Wallet

	expireAt int64
}

func (s *SpecRegular) SetExpireAt(expireAt int64) {
	s.expireAt = expireAt
}

type SpecSeqno struct {
	// Instead of calling contract 'seqno' method,
	// this function wil be used (if not nil) to get seqno for new transaction.
	// You may use it to set seqno according to your own logic,
	// for example for additional idempotency,
	// if build message is not enough, or for additional security
	seqnoFetcher func(ctx context.Context, subWallet uint32) (uint32, error)
}

// Deprecated: Use SetSeqnoFetcher
func (s *SpecSeqno) SetCustomSeqnoFetcher(fetcher func() uint32) {
	s.seqnoFetcher = func(ctx context.Context, subWallet uint32) (uint32, error) {
		return fetcher(), nil
	}
}

func (s *SpecSeqno) SetSeqnoFetcher(fetcher func(ctx context.Context, subWallet uint32) (uint32, error)) {
	s.seqnoFetcher = fetcher
}

type SpecQuery struct {
	// Instead of generating random query id with message ttl,
	// this function wil be used (if not nil) to get query id for new transaction.
	// You may use it to set query id according to your own logic,
	// for example for additional idempotency,
	// if build message is not enough, or for additional security
	//
	// Do not set ttl to high if you are sending many messages,
	// unexpired executed messages will be cached in contract,
	// and it may become too expensive to make transactions.
	customQueryIDFetcher func() (ttl uint32, randPart uint32)
}

func (s *SpecQuery) SetCustomQueryIDFetcher(fetcher func() (ttl uint32, randPart uint32)) {
	s.customQueryIDFetcher = fetcher
}
