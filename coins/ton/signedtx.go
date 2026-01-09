package ton

import (
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/okx/go-wallet-sdk/coins/ton/ton/wallet"
)

type SignedTx struct {
	Code    string `json:"code"`
	Data    string `json:"data"`
	Tx      string `json:"tx"`
	Hash    string `json:"txHash,omitempty"`
	Address string `json:"address"`
}

func NewSignedTx(address string) *SignedTx {
	return &SignedTx{
		Address: address,
	}
}

func (s *SignedTx) FillInit(w *wallet.Wallet, seqno uint32) error {
	if seqno != 0 {
		return nil
	}
	stateInit, err := wallet.GetStateInit(w.PublicKey(), w.GetVersionConfig(), w.GetSubwalletID())
	if err != nil {
		return err
	}

	s.Code = base64.StdEncoding.EncodeToString(stateInit.Code.ToBOC())
	s.Data = base64.StdEncoding.EncodeToString(stateInit.Data.ToBOC())
	return nil
}

func (s *SignedTx) FillTx(tx string, withHash bool) error {
	s.Tx = tx
	if !withHash {
		return nil
	}
	var err error
	s.Hash, err = CalNormMsgHash(tx)
	return err
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
