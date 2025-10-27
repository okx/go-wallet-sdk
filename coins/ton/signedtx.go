package ton

import (
	"encoding/json"
	"errors"
)

type SignedTx struct {
	Code     string `json:"code"`
	Data     string `json:"data"`
	Tx       string `json:"tx"`
	Hash     string `json:"txHash,omitempty"`
	NormHash string `json:"normHash,omitempty"`
	Address  string `json:"address"`
}

func (s *SignedTx) FillTxOnly(tx string) {
	s.Tx = tx
}

func (s *SignedTx) FillTx(tx string) {
	s.Tx = tx
	s.Hash, _ = CalTxHash(tx)
	s.NormHash, _ = CalNormMsgHash(tx)
}

func (s *SignedTx) Str() (string, error) {
	if s == nil {
		return "", errors.New("invalid tx")
	}
	r, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(r), err
}
