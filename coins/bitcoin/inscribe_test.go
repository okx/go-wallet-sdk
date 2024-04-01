package bitcoin

import (
	"encoding/json"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInscribe(t *testing.T) {
	network := &chaincfg.TestNet3Params

	commitTxPrevOutputList := make([]*PrevOutput, 0)
	commitTxPrevOutputList = append(commitTxPrevOutputList, &PrevOutput{
		TxId:       "453aa6dd39f31f06cd50b72a8683b8c0402ab36f889d96696317503a025a21b5",
		VOut:       0,
		Amount:     546,
		Address:    "2NF33rckfiQTiE5Guk5ufUdwms8PgmtnEdc",
		PrivateKey: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
	})
	commitTxPrevOutputList = append(commitTxPrevOutputList, &PrevOutput{
		TxId:       "22c8a4869f2aa9ee5994959c0978106130290cda53f6e933a8dda2dcb82508d4",
		VOut:       0,
		Amount:     546,
		Address:    "tb1qtsq9c4fje6qsmheql8gajwtrrdrs38kdzeersc",
		PrivateKey: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
	})
	commitTxPrevOutputList = append(commitTxPrevOutputList, &PrevOutput{
		TxId:       "3c6f205ec2995696d5bc852709d234a63aad82131b5b7615504e2e3e9ff88987",
		VOut:       0,
		Amount:     546,
		Address:    "mouQtmBWDS7JnT65Grj2tPzdSmGKJgRMhE",
		PrivateKey: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
	})
	commitTxPrevOutputList = append(commitTxPrevOutputList, &PrevOutput{
		TxId:       "aa09fa48dda0e2b7de1843c3db8d3f2d7f2cbe0f83331a125b06516a348abd26",
		VOut:       4,
		Amount:     1142196,
		Address:    "tb1pklh8lqax5l7m2ycypptv2emc4gata2dy28svnwcp9u32wlkenvsspcvhsr",
		PrivateKey: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
	})

	inscriptionDataList := make([]InscriptionData, 0)
	inscriptionDataList = append(inscriptionDataList, InscriptionData{
		ContentType: "text/plain;charset=utf-8",
		Body:        []byte(`{"p":"brc-20","op":"mint","tick":"xcvb","amt":"100"}`),
		RevealAddr:  "tb1pklh8lqax5l7m2ycypptv2emc4gata2dy28svnwcp9u32wlkenvsspcvhsr",
	})
	inscriptionDataList = append(inscriptionDataList, InscriptionData{
		ContentType: "text/plain;charset=utf-8",
		Body:        []byte(`{"p":"brc-20","op":"mint","tick":"xcvb","amt":"10"}`),
		RevealAddr:  "mouQtmBWDS7JnT65Grj2tPzdSmGKJgRMhE",
	})
	inscriptionDataList = append(inscriptionDataList, InscriptionData{
		ContentType: "text/plain;charset=utf-8",
		Body:        []byte(`{"p":"brc-20","op":"mint","tick":"xcvb","amt":"10000"}`),
		RevealAddr:  "tb1qtsq9c4fje6qsmheql8gajwtrrdrs38kdzeersc",
	})
	inscriptionDataList = append(inscriptionDataList, InscriptionData{
		ContentType: "text/plain;charset=utf-8",
		Body:        []byte(`{"p":"brc-20","op":"mint","tick":"xcvb","amt":"1"}`),
		RevealAddr:  "2NF33rckfiQTiE5Guk5ufUdwms8PgmtnEdc",
	})

	request := &InscriptionRequest{
		CommitTxPrevOutputList: commitTxPrevOutputList,
		CommitFeeRate:          2,
		RevealFeeRate:          2,
		RevealOutValue:         546,
		InscriptionDataList:    inscriptionDataList,
		ChangeAddress:          "tb1pklh8lqax5l7m2ycypptv2emc4gata2dy28svnwcp9u32wlkenvsspcvhsr",
	}

	txs, err := Inscribe(network, request)
	if err != nil {
		t.Fatal(err)
	}
	txsBytes, err := json.Marshal(txs)
	if err != nil {
		t.Fatal(err)
	}
	expected := `{"commitTx":"02000000000104b5215a023a50176369969d886fb32a40c0b883862ab750cd061ff339dda63a4500000000171600145c005c5532ce810ddf20f9d1d939631b47089ecdfdffffffd40825b8dca2dda833e9f653da0c2930611078099c959459eea92a9f86a4c8220000000000fdffffff8789f89f3e2e4e5015765b1b1382ad3aa634d2092785bcd5965699c25e206f3c000000006b483045022100f754ad06bad6452f96ca89fcde5f8fb5d66f5add8ea95c0d3c28ef5209a7a58d022045259c123ed509acdf625fa41c449776f4e0a049bc9ab1f0f2c136ef9896a9ce01210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2ffdffffff26bd8a346a51065b121a33830fbe2c7f2d3f8ddbc34318deb7e2a0dd48fa09aa0400000000fdffffff0550030000000000002251206ff0ac47ccff79fc3eaab0cd0047c28dead95cd35c6c695dfe33010b8807d16c3c03000000000000225120845a93ad3f2f36750672201709a48e6ad458cc0a42455f0786cf3bbbe42a6d183803000000000000225120be60aa4826e2e3a3245158c0e7b36543ed7ead2ed40a541c4583b80d4b3762003803000000000000225120e7ff49e9dee3ddaf3a811f12954a9c66cc98bf01c4eccb1ec093acf04ee2d1ff8062110000000000225120b7ee7f83a6a7fdb513040856c56778aa3abea9a451e0c9bb012f22a77ed99b210247304402207589b3e41b82547801a3613efbd3edb1438576679f211ee104e30e02732e42a702200341c77095a196fb7e4c20eb446fb7e9ab6ab4d02609eb72a26708cf7b453daa01210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f02483045022100f56201b18bd33472e19f4564a84c819b08af6fddab55ca408112f12cabc849be0220343f5c7f391b5cb69ef90a41d513bf1b05f24f4c6701ad4fb2db3e9d1a64ab4c01210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f000140dca26614fd80c47dee4eb2db0c776c2e888c453d51afc2fc01d28b1e2903d1f444e0717d6c1d4110d3631e1c480a1b20314315c72929945606acdcb01309910b00000000","revealTxs":["02000000000101a4a801d4e06cf7d6e3d376686edb048e26cede46bf248f94ddb290dfe9d426640000000000fdffffff012202000000000000225120b7ee7f83a6a7fdb513040856c56778aa3abea9a451e0c9bb012f22a77ed99b21034061d734a5a91aacb5a257a74e73ed6ed99d81918e3be4f917cb1532b7087175d807cb6f7a7e6518ed1db087318e9ed536071d4fe3e86590abd01bb1335484cccf7a2057bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2fac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d3800347b2270223a226272632d3230222c226f70223a226d696e74222c227469636b223a2278637662222c22616d74223a22313030227d6821c057bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f00000000","02000000000101a4a801d4e06cf7d6e3d376686edb048e26cede46bf248f94ddb290dfe9d426640100000000fdffffff0122020000000000001976a9145c005c5532ce810ddf20f9d1d939631b47089ecd88ac03403876a2ef916ebc497912941b6bee621389a734a8d88eb639bde717bb614f0aa42e511745506a70ee5f4693cfbd17df015e49cbd53fbbec37cf15f3b9b5cd7dfe792057bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2fac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d3800337b2270223a226272632d3230222c226f70223a226d696e74222c227469636b223a2278637662222c22616d74223a223130227d6821c157bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f00000000","02000000000101a4a801d4e06cf7d6e3d376686edb048e26cede46bf248f94ddb290dfe9d426640200000000fdffffff0122020000000000001600145c005c5532ce810ddf20f9d1d939631b47089ecd0340f3af92405a2fbb5105cad1a9c498432ff9e69097801b3997d149d962fca8a84f68c4c3a7e46723c4f503c46b62aac751ade35a3de8882247d329c01e98863ff87c2057bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2fac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d3800367b2270223a226272632d3230222c226f70223a226d696e74222c227469636b223a2278637662222c22616d74223a223130303030227d6821c057bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f00000000","02000000000101a4a801d4e06cf7d6e3d376686edb048e26cede46bf248f94ddb290dfe9d426640300000000fdffffff01220200000000000017a914ef05515a0595d15eaf90d9f62fb85873a6d8c0b487034097ee19b8f9a51bc32cd8cacb16e8abf8a12119a58d0c591f1072286034e4ff7622b271d8a1d7b87ba1b958bfe7eaa23d2923a8b32793e23f3be594a565b8e3cf782057bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2fac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d3800327b2270223a226272632d3230222c226f70223a226d696e74222c227469636b223a2278637662222c22616d74223a2231227d6821c057bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f00000000"],"commitTxFee":1182,"revealTxFees":[302,282,278,278],"commitAddrs":["tb1pdlc2c37vlaulc042krxsq37z3h4djhxnt3kxjh07xvqshzq869kqz5sgrc","tb1ps3df8tfl9um82pnjyqtsnfywdt293nq2gfz47puxeuamhep2d5vq0jujz6","tb1phes25jpxut36xfz3trqw0vm9g0khatfw6s99g8z9swuq6jehvgqqdsrvg2","tb1pull5n6w7u0w67w5pruff2j5uvmxf30cpcnkvk8kqjwk0qnhz68ls68tklf"]}`
	require.Equal(t, expected, string(txsBytes))
}

func TestInscribeForMPCUnsigned(t *testing.T) {
	network := &chaincfg.TestNet3Params

	commitTxPrevOutputList := make([]*PrevOutput, 0)
	commitTxPrevOutputList = append(commitTxPrevOutputList, &PrevOutput{
		TxId:       "453aa6dd39f31f06cd50b72a8683b8c0402ab36f889d96696317503a025a21b5",
		VOut:       0,
		Amount:     546,
		Address:    "2NF33rckfiQTiE5Guk5ufUdwms8PgmtnEdc",
		PrivateKey: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
		PublicKey:  "0357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f",
	})
	commitTxPrevOutputList = append(commitTxPrevOutputList, &PrevOutput{
		TxId:       "22c8a4869f2aa9ee5994959c0978106130290cda53f6e933a8dda2dcb82508d4",
		VOut:       0,
		Amount:     546,
		Address:    "tb1qtsq9c4fje6qsmheql8gajwtrrdrs38kdzeersc",
		PrivateKey: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
		PublicKey:  "0357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f",
	})
	commitTxPrevOutputList = append(commitTxPrevOutputList, &PrevOutput{
		TxId:       "3c6f205ec2995696d5bc852709d234a63aad82131b5b7615504e2e3e9ff88987",
		VOut:       0,
		Amount:     546,
		Address:    "mouQtmBWDS7JnT65Grj2tPzdSmGKJgRMhE",
		PrivateKey: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
		PublicKey:  "0357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f",
	})
	commitTxPrevOutputList = append(commitTxPrevOutputList, &PrevOutput{
		TxId:       "aa09fa48dda0e2b7de1843c3db8d3f2d7f2cbe0f83331a125b06516a348abd26",
		VOut:       4,
		Amount:     1142196,
		Address:    "tb1pklh8lqax5l7m2ycypptv2emc4gata2dy28svnwcp9u32wlkenvsspcvhsr",
		PrivateKey: "cPnvkvUYyHcSSS26iD1dkrJdV7k1RoUqJLhn3CYxpo398PdLVE22",
		PublicKey:  "0357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f",
	})

	inscriptionDataList := make([]InscriptionData, 0)
	inscriptionDataList = append(inscriptionDataList, InscriptionData{
		ContentType: "text/plain;charset=utf-8",
		Body:        []byte(`{"p":"brc-20","op":"mint","tick":"xcvb","amt":"100"}`),
		RevealAddr:  "tb1pklh8lqax5l7m2ycypptv2emc4gata2dy28svnwcp9u32wlkenvsspcvhsr",
	})
	inscriptionDataList = append(inscriptionDataList, InscriptionData{
		ContentType: "text/plain;charset=utf-8",
		Body:        []byte(`{"p":"brc-20","op":"mint","tick":"xcvb","amt":"10"}`),
		RevealAddr:  "mouQtmBWDS7JnT65Grj2tPzdSmGKJgRMhE",
	})
	inscriptionDataList = append(inscriptionDataList, InscriptionData{
		ContentType: "text/plain;charset=utf-8",
		Body:        []byte(`{"p":"brc-20","op":"mint","tick":"xcvb","amt":"10000"}`),
		RevealAddr:  "tb1qtsq9c4fje6qsmheql8gajwtrrdrs38kdzeersc",
	})
	inscriptionDataList = append(inscriptionDataList, InscriptionData{
		ContentType: "text/plain;charset=utf-8",
		Body:        []byte(`{"p":"brc-20","op":"mint","tick":"xcvb","amt":"1"}`),
		RevealAddr:  "2NF33rckfiQTiE5Guk5ufUdwms8PgmtnEdc",
	})

	request := &InscriptionRequest{
		CommitTxPrevOutputList: commitTxPrevOutputList,
		CommitFeeRate:          2,
		RevealFeeRate:          2,
		RevealOutValue:         546,
		InscriptionDataList:    inscriptionDataList,
		ChangeAddress:          "tb1pklh8lqax5l7m2ycypptv2emc4gata2dy28svnwcp9u32wlkenvsspcvhsr",
	}
	res, err := InscribeForMPCUnsigned(request, network, nil, nil)
	require.NoError(t, err)
	rb, err := json.Marshal(res)
	require.NoError(t, err)
	expected := `{"sigHashList":["89d176c6dd56cf7ac84c2b0136098c7394cdcb29318c8513092150af7f0ef685","a22c61c3fdead3e958364786ffc796daaeeb918ca1033b8dc7228e8180a5859b","13c56286442af478c8b89b8d313f54b98fdf9ee0ddd0429b025c718913f92c96","a1cf51c368086658d473c0f8045b7fd5bf90178f7e4ce8926ec1b1e7d629b419"],"commitTx":"02000000000104b5215a023a50176369969d886fb32a40c0b883862ab750cd061ff339dda63a4500000000171600145c005c5532ce810ddf20f9d1d939631b47089ecdfdffffffd40825b8dca2dda833e9f653da0c2930611078099c959459eea92a9f86a4c8220000000000fdffffff8789f89f3e2e4e5015765b1b1382ad3aa634d2092785bcd5965699c25e206f3c00000000210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2ffdffffff26bd8a346a51065b121a33830fbe2c7f2d3f8ddbc34318deb7e2a0dd48fa09aa0400000000fdffffff0550030000000000002251206ff0ac47ccff79fc3eaab0cd0047c28dead95cd35c6c695dfe33010b8807d16c3c03000000000000225120845a93ad3f2f36750672201709a48e6ad458cc0a42455f0786cf3bbbe42a6d183803000000000000225120be60aa4826e2e3a3245158c0e7b36543ed7ead2ed40a541c4583b80d4b3762003803000000000000225120e7ff49e9dee3ddaf3a811f12954a9c66cc98bf01c4eccb1ec093acf04ee2d1ff8262110000000000225120b7ee7f83a6a7fdb513040856c56778aa3abea9a451e0c9bb012f22a77ed99b2101210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f01210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f000000000000","revealTxs":["020000000001015c3a8f2abcd39b0e4a1fcf9fff905e17ed130fccd81a079271eb3f28e127a7e80000000000fdffffff012202000000000000225120b7ee7f83a6a7fdb513040856c56778aa3abea9a451e0c9bb012f22a77ed99b2103407d77a1c8dee85e59b2446f707e2e37aac600ce45cbb2ceb90554c8a391540de0c6df415177c65ab279b87abb29ca39fe6e07cae73fb9f726674e66412fd9b3bf7a2057bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2fac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d3800347b2270223a226272632d3230222c226f70223a226d696e74222c227469636b223a2278637662222c22616d74223a22313030227d6821c057bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f00000000","020000000001015c3a8f2abcd39b0e4a1fcf9fff905e17ed130fccd81a079271eb3f28e127a7e80100000000fdffffff0122020000000000001976a9145c005c5532ce810ddf20f9d1d939631b47089ecd88ac0340e8d1b62dd426a98abe501dabf83969767d44a2c3542acf358a66d3dbbf5f6f8fa2144183fdead4e4a3e972cb522de94bfd12af4d940ac7694b90757d5651055a792057bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2fac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d3800337b2270223a226272632d3230222c226f70223a226d696e74222c227469636b223a2278637662222c22616d74223a223130227d6821c157bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f00000000","020000000001015c3a8f2abcd39b0e4a1fcf9fff905e17ed130fccd81a079271eb3f28e127a7e80200000000fdffffff0122020000000000001600145c005c5532ce810ddf20f9d1d939631b47089ecd0340ae4d6c59687a723c69a011253855f047481c309d084e783f80a2ea1df16190db8ce16598da992416b678b7b3626379184939f1cea45421e77b0a74b8fbecf23a7c2057bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2fac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d3800367b2270223a226272632d3230222c226f70223a226d696e74222c227469636b223a2278637662222c22616d74223a223130303030227d6821c057bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f00000000","020000000001015c3a8f2abcd39b0e4a1fcf9fff905e17ed130fccd81a079271eb3f28e127a7e80300000000fdffffff01220200000000000017a914ef05515a0595d15eaf90d9f62fb85873a6d8c0b4870340ab4d04bbf1e15eb488229f074713de28cd0798cc4ce570bb0022106c97c2ba5fa8f80d15603d70a54470ba05887b05b01acbaca7b4ee5deaf6fd51846e95cfce782057bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2fac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d3800327b2270223a226272632d3230222c226f70223a226d696e74222c227469636b223a2278637662222c22616d74223a2231227d6821c057bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f00000000"],"commitTxFee":1180,"revealTxFees":[302,282,278,278],"commitAddrs":["tb1pdlc2c37vlaulc042krxsq37z3h4djhxnt3kxjh07xvqshzq869kqz5sgrc","tb1ps3df8tfl9um82pnjyqtsnfywdt293nq2gfz47puxeuamhep2d5vq0jujz6","tb1phes25jpxut36xfz3trqw0vm9g0khatfw6s99g8z9swuq6jehvgqqdsrvg2","tb1pull5n6w7u0w67w5pruff2j5uvmxf30cpcnkvk8kqjwk0qnhz68ls68tklf"]}`
	require.Equal(t, expected, string(rb))

}
