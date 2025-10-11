package waves

import (
	"github.com/emresenyuva/go-wallet-sdk/coins/waves/crypto"
	"github.com/emresenyuva/go-wallet-sdk/coins/waves/types"
)

// NewUnsignedTransferWithSig creates new TransferWithSig transaction without signature and ID.
func NewUnsignedTransferWithSig(senderPK crypto.PublicKey, amountAsset, feeAsset types.OptionalAsset,
	timestamp, amount, fee uint64, recipient types.Recipient, attachment types.Attachment) *types.TransferWithSig {
	t := types.Transfer{
		SenderPK:    senderPK,
		Recipient:   recipient,
		AmountAsset: amountAsset,
		Amount:      amount,
		FeeAsset:    feeAsset,
		Fee:         fee,
		Timestamp:   timestamp,
		Attachment:  attachment,
	}
	return &types.TransferWithSig{Type: types.TransferTransaction, Version: 1, Transfer: t}
}

func SignTransferWithSig(t *types.TransferWithSig, privateKey crypto.SecretKey) error {
	scheme := types.Scheme(MainNetScheme)
	err := t.Sign(scheme, privateKey)
	if err != nil {
		return err
	}
	return nil
}
