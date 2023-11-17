package eos

import (
	"fmt"
	"github.com/eoscanada/eos-go/ecc"
	"github.com/okx/go-wallet-sdk/coins/eos/types"
)

// NewTransaction creates a new EOS Transaction object, ready to sign.
func NewTransaction(actions []*types.Action, opts *types.TxOptions) *types.Transaction {
	if opts == nil {
		opts = &types.TxOptions{}
	}
	// note: need HeadBlockID, ChainID
	tx := &types.Transaction{Actions: actions}
	tx.Fill(opts.HeadBlockID, opts.DelaySecs, opts.MaxNetUsageWords, opts.MaxCPUUsageMS, opts.Expiration)
	return tx
}

func NewTransactionWithParams(from, to, memo string, quantity types.Asset, opts *types.TxOptions) *types.Transaction {
	if len(from) > 12 || len(to) > 12 {
		return nil
	}
	return NewTransaction([]*types.Action{types.NewTransfer(from, to, quantity, memo)}, opts)
}

func NewContractTransaction(name, from, to, memo string, quantity types.Asset, opts *types.TxOptions) *types.Transaction {
	if len(from) > 12 || len(to) > 12 || len(name) > 12 {
		return nil
	}
	return NewTransaction([]*types.Action{types.NewContractTransfer(name, from, to, quantity, memo)}, opts)
}

func NewBuyRamTransaction(from, to string, quantity types.Asset, opts *types.TxOptions) *types.Transaction {
	if len(from) > 12 || len(to) > 12 {
		return nil
	}
	return NewTransaction([]*types.Action{types.NewBuyRAM(from, to, quantity)}, opts)
}

// NewBuyRAMBytesTransaction creates a new EOS BuyRAMBytes transaction.
func NewBuyRAMBytesTransaction(from, to string, bytes uint64, opts *types.TxOptions) *types.Transaction {
	if len(from) > 12 || len(to) > 12 {
		return nil
	}
	return NewTransaction([]*types.Action{types.NewBuyRAMBytes(from, to, uint32(bytes))}, opts)
}

// NewDelegateBWTransaction creates a new EOS DelegateBW transaction.
func NewDelegateBWTransaction(from, to string, stakeCPU, stakeNet types.Asset, doTransfer bool,
	opts *types.TxOptions) *types.Transaction {
	if len(from) > 12 || len(to) > 12 {
		return nil
	}
	return NewTransaction([]*types.Action{types.NewDelegateBW(from, to, stakeCPU, stakeNet, doTransfer)}, opts)
}

// NewSellRAMTransaction creates a new EOS SellRAM transaction.
func NewSellRAMTransaction(account string, bytes uint64, opts *types.TxOptions) *types.Transaction {
	if len(account) > 12 {
		return nil
	}
	return NewTransaction([]*types.Action{types.NewSellRAM(account, bytes)}, opts)
}

// NewUndelegateBWTransaction creates a new EOS UndelegateBW transaction.
func NewUndelegateBWTransaction(from, to string, unstakeCPU, unstakeNet types.Asset, opts *types.TxOptions) *types.Transaction {
	if len(from) > 12 || len(to) > 12 {
		return nil
	}
	return NewTransaction([]*types.Action{types.NewUndelegateBW(from, to, unstakeCPU, unstakeNet)}, opts)
}

// SignTransactionWithWIFs signs a transaction with the given WIFs.
func SignTransactionWithWIFs(wifs []string, tx *types.Transaction, chainID types.Checksum256,
	compression types.CompressionType) (*types.SignedTransaction, *types.PackedTransaction, error) {
	stx := types.NewSignedTransaction(tx)

	signer, err := NewSignerFromWIFs(wifs)
	if err != nil {
		return nil, nil, err
	}

	requiredKeys := make([]ecc.PublicKey, len(wifs))
	for i, wif := range wifs {
		key, err := ecc.NewPrivateKey(wif)
		if err != nil {
			continue
		}
		requiredKeys[i] = key.PublicKey()
	}

	signedTx, err := signer.Sign(stx, chainID, requiredKeys...)
	if err != nil {
		return nil, nil, fmt.Errorf("signing through wallet: %w", err)
	}

	packed, err := signedTx.Pack(compression)
	if err != nil {
		return nil, nil, fmt.Errorf("packing transaction: %w", err)
	}
	return signedTx, packed, nil
}

// SignTransaction signs a transaction with the given WIF.
func SignTransaction(wifKey string, tx *types.Transaction, chainID types.Checksum256,
	compression types.CompressionType) (*types.SignedTransaction, *types.PackedTransaction, error) {
	return SignTransactionWithWIFs([]string{wifKey}, tx, chainID, compression)
}
