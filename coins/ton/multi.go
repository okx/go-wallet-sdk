package ton

import (
	"errors"
	"github.com/okx/go-wallet-sdk/coins/ton/address"
	"math/big"
	"time"
)

const (
	maxUnix            = 10_000_000_000
	mainNet            = "-239"
	TransferTimeoutSec = 600
)

var (
	ErrInvalidMultiRequest = errors.New("invalid multi request content")
)

type Msg struct {
	Address   string `json:"address"`
	Amount    string `json:"amount"`
	Payload   string `json:"payload"`
	StateInit string `json:"stateInit"`
}

func (s *Msg) Check() error {
	if s == nil {
		return ErrInvalidMultiRequest
	}
	addr, err := address.ParseAddr(s.Address)
	if err != nil {
		return err
	}
	if addr.IsTestnetOnly() {
		return ErrInvalidMultiRequest
	}
	b, ok := new(big.Int).SetString(s.Amount, 10)
	if !ok || b.Cmp(big.NewInt(0)) < 0 {
		return ErrInvalidMultiRequest
	}
	return nil
}

type MultiRequest struct {
	Messages   []*Msg `json:"messages"`
	ValidUntil int64  `json:"valid_until"`
	From       string `json:"from"`
	Network    string `json:"network"`
}

func (s *MultiRequest) Check() error {
	if s == nil || len(s.Messages) > 4 || s.ValidUntil < 0 {
		return ErrInvalidMultiRequest
	}
	if mainNet != s.Network {
		return ErrInvalidMultiRequest
	}
	addr, err := address.ParseAddr(s.From)
	if err != nil {
		return err
	}
	if addr.IsTestnetOnly() {
		return ErrInvalidMultiRequest
	}
	for _, v := range s.Messages {
		if err := v.Check(); err != nil {
			return err
		}
	}
	if s.ValidUntil == 0 {
		s.ValidUntil = time.Now().Unix() + TransferTimeoutSec
	}
	if s.ValidUntil > maxUnix {
		s.ValidUntil = s.ValidUntil / 1000
	}
	return nil
}
