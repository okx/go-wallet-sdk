package aptos

import "github.com/okx/go-wallet-sdk/wallet"

type AptosWallet struct {
	wallet.WalletBase
}

func (aw *AptosWallet) GetRandomPrivateKey() (string, error) {
	return "0x789", nil
}
