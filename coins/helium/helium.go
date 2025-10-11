package helium

import (
	"encoding/base64"
	"github.com/emresenyuva/go-wallet-sdk/coins/helium/keypair"
	"github.com/emresenyuva/go-wallet-sdk/coins/helium/transactions"
)

func Sign(private string, from, to string, amount, fee, nonce uint64, tokenType string, isMax bool) (string, error) {
	fromAble := keypair.NewAddressable(from)
	tmpSig := make([]byte, 64)
	toAmount := map[string]uint64{to: amount}
	v2 := transactions.NewPaymentV2Tx(fromAble, toAmount, fee, nonce, tokenType, isMax, tmpSig)
	transaction, err := v2.BuildTransaction(true)
	if err != nil {
		return "", err
	}
	// 1 for edd25519
	kp := keypair.NewKeypairFromHex(1, private)
	sig, err := kp.Sign(transaction)
	if err != nil {
		return "", err
	}
	v2.SetSignature(sig)
	ser, err := v2.Serialize()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ser), nil
}

func NewAddress(private string) string {
	kp := keypair.NewKeypairFromHex(1, private)
	return kp.CreateAddressable().GetAddress()

}
