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
	expected := "{\"sigHashList\":[\"76ea9481dfc53bf1397b09877368ab5de829157beffced747dd0b38d9de5ecd1\",\"b6dd4824f6afa369a4b29b1fe59f8fc8d21309d864db0e283d9b7266b55160ab\",\"3b431bf23223a1f77101cdb6cb375d74a9674a6f4936810d7dc9307a5853bb12\",\"f70e35aef542ba44583341cd43f1aab2beaecf5044036f4db3cbbd1fe64063ab\"],\"commitTx\":\"02000000000104b5215a023a50176369969d886fb32a40c0b883862ab750cd061ff339dda63a4500000000171600145c005c5532ce810ddf20f9d1d939631b47089ecdfdffffffd40825b8dca2dda833e9f653da0c2930611078099c959459eea92a9f86a4c8220000000000fdffffff8789f89f3e2e4e5015765b1b1382ad3aa634d2092785bcd5965699c25e206f3c00000000210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2ffdffffff26bd8a346a51065b121a33830fbe2c7f2d3f8ddbc34318deb7e2a0dd48fa09aa0400000000fdffffff0550030000000000002251206ff0ac47ccff79fc3eaab0cd0047c28dead95cd35c6c695dfe33010b8807d16c3c03000000000000225120845a93ad3f2f36750672201709a48e6ad458cc0a42455f0786cf3bbbe42a6d183803000000000000225120be60aa4826e2e3a3245158c0e7b36543ed7ead2ed40a541c4583b80d4b3762003803000000000000225120e7ff49e9dee3ddaf3a811f12954a9c66cc98bf01c4eccb1ec093acf04ee2d1ff8462110000000000225120b7ee7f83a6a7fdb513040856c56778aa3abea9a451e0c9bb012f22a77ed99b2101210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f01210357bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f000000000000\",\"revealTxs\":[\"0200000000010115e22cee9da0b8a8954193106dccad8891dd44a4ac320f026a2cc78ad28e3a6e0000000000fdffffff012202000000000000225120b7ee7f83a6a7fdb513040856c56778aa3abea9a451e0c9bb012f22a77ed99b210340cb7567a1989a4d83a411d7cbfaa8ccef109812d055c0f00fc61b97cb62d81f9f28ae35635f0d5626bce39031ec035b30b0f0c8e45fc477c13fa819cb41843b237a2057bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2fac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d3800347b2270223a226272632d3230222c226f70223a226d696e74222c227469636b223a2278637662222c22616d74223a22313030227d6821c057bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f00000000\",\"0200000000010115e22cee9da0b8a8954193106dccad8891dd44a4ac320f026a2cc78ad28e3a6e0100000000fdffffff0122020000000000001976a9145c005c5532ce810ddf20f9d1d939631b47089ecd88ac03405aa61ba432140fa1ce9a2daf18d8fcd14eea0ddde3d83fb5f070a42bf44e7e456826c45cf1f0eddd31efa2cbd96c8d9d1076f0e51e690e02812dfec5b3a39e78792057bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2fac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d3800337b2270223a226272632d3230222c226f70223a226d696e74222c227469636b223a2278637662222c22616d74223a223130227d6821c157bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f00000000\",\"0200000000010115e22cee9da0b8a8954193106dccad8891dd44a4ac320f026a2cc78ad28e3a6e0200000000fdffffff0122020000000000001600145c005c5532ce810ddf20f9d1d939631b47089ecd0340da5e459b952b5dc4e5e8631c87e1d47d5bce8fba77dd48ea68de9eedff9988ce1cf50d715009d4a5789f98d7fbecf04f5f6d54eaa08704a67154fc17d35f18847c2057bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2fac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d3800367b2270223a226272632d3230222c226f70223a226d696e74222c227469636b223a2278637662222c22616d74223a223130303030227d6821c057bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f00000000\",\"0200000000010115e22cee9da0b8a8954193106dccad8891dd44a4ac320f026a2cc78ad28e3a6e0300000000fdffffff01220200000000000017a914ef05515a0595d15eaf90d9f62fb85873a6d8c0b48703400db47a127505525dfe0929db68a3ed20af511575d95b01485a43d8e1da2cae8f271c53414f17fc672e93351fed6c4104b6305f6056a534c2b25f9aedbd547f6d782057bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2fac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d3800327b2270223a226272632d3230222c226f70223a226d696e74222c227469636b223a2278637662222c22616d74223a2231227d6821c057bbb2d4a9cb8a2357633f201b9c518c2795ded682b7913c6beef3fe23bd6d2f00000000\"],\"commitTxFee\":1178,\"revealTxFees\":[302,282,278,278],\"commitAddrs\":[\"tb1pdlc2c37vlaulc042krxsq37z3h4djhxnt3kxjh07xvqshzq869kqz5sgrc\",\"tb1ps3df8tfl9um82pnjyqtsnfywdt293nq2gfz47puxeuamhep2d5vq0jujz6\",\"tb1phes25jpxut36xfz3trqw0vm9g0khatfw6s99g8z9swuq6jehvgqqdsrvg2\",\"tb1pull5n6w7u0w67w5pruff2j5uvmxf30cpcnkvk8kqjwk0qnhz68ls68tklf\"]}"
	require.Equal(t, expected, string(rb))

}
