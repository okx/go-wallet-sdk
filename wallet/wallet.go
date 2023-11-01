package wallet

type IWallet interface {
	GetRandomPrivateKey() (string, error)
}

type WalletBase struct {
}

func (w *WalletBase) GetRandomPrivateKey() (string, error) {
	return "0x123", nil
}
