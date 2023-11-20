package waves

import (
	"encoding/json"
	"github.com/okx/go-wallet-sdk/coins/waves/crypto"
	"github.com/okx/go-wallet-sdk/coins/waves/types"
	"reflect"
	"testing"
)

const (
	p1 = "tMUA9XRwPTiUXCTmEvU6kFkqTFKxSpaAFvQwyAT29GR"
	s1 = "5NBbF9dHDfuJw2WC8m3Am5kJwKMXbLmN2eh4Cmqsgo5w"
	// a1 = "3Mq7eCKTgNAoEag4eQVHZYGZKRNYKmodEpM"
	a1 = "3NAmitjJqrxPqNNERfa1ZaNnN7t1FkzL26r"
	p2 = "GRcXDTsfpJZU6qUPkhjBX7dY1yKJ5mV2JJyWHWW1mUYK"
	s2 = "6QhEoSnJ12QDgeEAt3HYkPDBiYe15BArgSKWrV3DUctG"
	a2 = "3NB2pUqjoavApZeAmdsVYS84hyRGXZpeytA"
)

func TestNewTransfer(t *testing.T) {
	senderPublicKey, err := crypto.NewPublicKeyFromBase58(p1)
	if err != nil {
		t.Fatal(err)
		return
	}
	address, err := types.NewAddressFromString(a2)
	if err != nil {
		t.Fatal(err)
		return
	}
	waves := types.NewOptionalAssetWaves()
	// new transfer
	// amountAsset: WAVES, feeAsset: WAVES
	// timestamp: 1546300800, amount: 1, fee: 1
	// recipient: 3PJQXu9uQ8qoQK2f8zqXzjQYzjQ8JXhj2a
	tx := NewUnsignedTransferWithSig(senderPublicKey, waves, waves, 1655401735758, 2000000,
		200000, types.NewRecipientFromAddress(address), []byte("attachment"))
	// sign the tx
	secretKey, err := crypto.NewSecretKeyFromBase58(s1)
	if err != nil {
		t.Fatal(err)
		return
	}
	if err := SignTransferWithSig(tx, secretKey); err != nil {
		t.Fatal(err)
		return
	}
	idBytes, err := tx.ID.MarshalJSON()
	t.Log(string(idBytes))
	t.Log(tx)

	bts, err := json.Marshal(tx)
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(string(bts))
}

func TestNewUnsignedTransferWithSig(t *testing.T) {
	senderPublicKey, err := crypto.NewPublicKeyFromBase58("2wySdbAsXi1bfAfMBKC1NcyyJemUWLM4R5ECwXJiADUx")
	if err != nil {
		t.Fatal(err)
		return
	}
	address, err := types.NewAddressFromString("3NB2pUqjoavApZeAmdsVYS84hyRGXZpeytA")
	if err != nil {
		t.Fatal(err)
		return
	}
	recipient := types.NewRecipientFromAddress(address)
	amountAsset, err := types.NewOptionalAssetFromString("WAVES")
	if err != nil {
		t.Fatal(err)
		return
	}
	feeAsset, err := types.NewOptionalAssetFromString("WAVES")
	if err != nil {
		t.Fatal(err)
		return
	}
	type args struct {
		senderPK    crypto.PublicKey
		amountAsset types.OptionalAsset
		feeAsset    types.OptionalAsset
		timestamp   uint64
		amount      uint64
		fee         uint64
		recipient   types.Recipient
		attachment  types.Attachment
	}
	tests := []struct {
		name string
		args args
		want *types.TransferWithSig
	}{
		{
			name: "test",
			args: args{
				senderPK:    senderPublicKey,
				amountAsset: *amountAsset,
				feeAsset:    *feeAsset,
				timestamp:   1546300800,
				amount:      1,
				fee:         1,
				recipient:   recipient,
				attachment:  types.NewAttachmentFromBase58("test"),
			},
			want: &types.TransferWithSig{
				Type:    types.TransferTransaction,
				Version: 1,
				Transfer: types.Transfer{
					SenderPK:    senderPublicKey,
					Recipient:   recipient,
					AmountAsset: *amountAsset,
					Amount:      1,
					FeeAsset:    *feeAsset,
					Fee:         1,
					Timestamp:   1546300800,
					Attachment:  types.NewAttachmentFromBase58("test"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUnsignedTransferWithSig(tt.args.senderPK, tt.args.amountAsset, tt.args.feeAsset, tt.args.timestamp, tt.args.amount, tt.args.fee, tt.args.recipient, tt.args.attachment); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUnsignedTransferWithSig() = %v, want %v", got, tt.want)
			}
		})
	}
}
