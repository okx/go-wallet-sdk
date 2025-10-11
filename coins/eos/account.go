package eos

import (
	"github.com/emresenyuva/go-wallet-sdk/coins/eos/types"
	"github.com/eoscanada/eos-go/ecc"
)

// NewAccountTransaction creates a new account
func NewAccountTransaction(creator, newAccount string, pubKey ecc.PublicKey, buyRAMAmount, cpuStake, netStake types.Asset,
	doTransfer bool, opts *types.TxOptions) *types.Transaction {
	if len(creator) > 12 || len(newAccount) > 12 {
		return nil
	}
	var actions []*types.Action
	actions = append(actions, types.NewNewAccount(creator, newAccount, pubKey))
	actions = append(actions, types.NewDelegateBW(creator, newAccount, cpuStake, netStake, doTransfer))
	actions = append(actions, types.NewBuyRAM(creator, newAccount, buyRAMAmount))

	return NewTransaction(actions, opts)
}

func GenerateKeyPair() (privKey, pubKey string) {
	privateKey, _ := ecc.NewRandomPrivateKey()
	return privateKey.String(), privateKey.PublicKey().String()
}
