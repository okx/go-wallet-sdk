package doginals

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewDogeAddr(t *testing.T) {
	prv := "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22"
	privateKeyWif, err := btcutil.DecodeWIF(prv)
	if err != nil {
		t.Fatal(err)
	}
	pub := privateKeyWif.PrivKey.PubKey().SerializeCompressed()
	pkHash := btcutil.Hash160(pub)
	addr, err := btcutil.NewAddressPubKeyHash(pkHash, &DogeMainNetParams)
	if err != nil {
		t.Fatal(err)
	}
	require.Equal(t, "DDXZ8y3AhpaLYLo4HskDcEwuTuPumYLHwr", addr.EncodeAddress())
}

func TestInscribe(t *testing.T) {
	commitTxPrevOutputList := make([]*PrevOutput, 0)

	commitTxPrevOutputList = append(commitTxPrevOutputList, &PrevOutput{
		TxId:       "adc5edd2a536c92fed35b3d75cbdbc9f11212fe3aa6b55c0ac88c289ba7c4fae",
		VOut:       2,
		Amount:     317250000,
		Address:    "DFuDR3Vn22KMnrnVCxh6YavMAJP8TCPeA2",
		PrivateKey: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
	})

	inscriptionData := &InscriptionData{
		ContentType: "text/plain;charset=utf8",
		Body:        []byte(`{"p":"drc-20","op":"mint","tick":"tril","amt":"100"}`),
		RevealAddr:  "DFuDR3Vn22KMnrnVCxh6YavMAJP8TCPeA2",
	}
	request := &InscriptionRequest{
		CommitTxPrevOutputList: commitTxPrevOutputList,
		CommitFeeRate:          100000,
		RevealFeeRate:          100000,
		RevealOutValue:         100000,
		InscriptionData:        inscriptionData,
		Address:                "DFuDR3Vn22KMnrnVCxh6YavMAJP8TCPeA2",
	}

	txs, err := Inscribe(request)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "0200000001ae4f7cba89c288acc0556baae32f21119fbcbd5cd7b335ed2fc936a5d2edc5ad020000006a4730440220179bc484c573990cee3f75e911624eb9f1c7cd591e987a43c73cd56becfd24a802204e2218fff29cb916082478174a31a4801f85743a837abc670f49b30bae79f71501210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2ffdffffff03a08601000000000017a914183d9f2cc745bc20d6bf26ce15750f09bc9e44258720897603000000001976a9145c005c5532ce810ddf20f9d1d939631b47089ecd88ac30b03a0d000000001976a91476094cb45e019a8942a4861c02f4fd766bb662e588ac00000000", txs.CommitTx)
	assert.Equal(t, 1, len(txs.RevealTxs))
	assert.Equal(t, "0200000002d1a8fbc96f2a5b8b5ecd261a5d3459577a4fcb9e0040237c822cf073a7dd5da600000000c6036f72645117746578742f706c61696e3b636861727365743d7574663800347b2270223a226472632d3230222c226f70223a226d696e74222c227469636b223a227472696c222c22616d74223a22313030227d483045022100a76b11d4febbccdcb9e0cba6ddedb7522d8ee0d571ce9a39174db034495e838b02202a12bb6c30bbe4f54ba0ceb4a64759e61b47e61d4f3763001b30973cdf65caf40129210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2fad757575757551fdffffffd1a8fbc96f2a5b8b5ecd261a5d3459577a4fcb9e0040237c822cf073a7dd5da6010000006b483045022100a7a94d711cc8914742a45d5852995e0c0e581d3dda85771a99787b9a85624a0f02203afaa646cde2c1f550e227d57c09be996a72396cdbe298c2e7d74f2ee40e689101210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2ffdffffff01a0860100000000001976a91476094cb45e019a8942a4861c02f4fd766bb662e588ac00000000", txs.RevealTxs[0])
}
