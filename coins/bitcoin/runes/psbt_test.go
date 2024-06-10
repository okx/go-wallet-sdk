package bitcoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/psbt"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/okx/go-wallet-sdk/coins/bitcoin"
	"github.com/stretchr/testify/assert"
	"testing"
)

// https://mempool.space/testnet/tx/766d32e560fca19f22ee3cf526edd00c5fbfd2f9833b5ff51fb64d3564321833
func TestMergePsbt(t *testing.T) {
	sellerPsbt := "cHNidP8BAN0CAAAAApV3TG/w54SD9fdtU6W5j2cBOlpdj33WuIUbrckXp3w6AQAAAAD/////3Hae3avlOA0ThY1rsXFiaFygFIX3WBHy1cNXJe7cVKwBAAAAAP////8DZQAAAAAAAAAiUSDZM1dVZgQG1pEcxQ2hKvQFtwfyJx6IV+z8hSN6zbZwRxAnAAAAAAAAIlEg2TNXVWYEBtaRHMUNoSr0BbcH8iceiFfs/IUjes22cEeWmRoAAAAAACJRINkzV1VmBAbWkRzFDaEq9AW3B/InHohX7PyFI3rNtnBHAAAAAAABASulwRoAAAAAACJRINkzV1VmBAbWkRzFDaEq9AW3B/InHohX7PyFI3rNtnBHAQMEAQAAAAETQA+Y1MmQQtNh33YYoUrZ5l276lwWmwGUtLijAespxz17HQpf3XfA5BE5x9ZJ8o5o+x+KAP1Zkdmwo1KznWj88lcBFyAplEsuIf4iLmoDGOkWHBC0Cn7HA2v7Rl0gCiCXK31XngAAAAAA"
	sp, err := psbt.NewFromRawBytes(bytes.NewReader([]byte(sellerPsbt)), true)
	assert.NoError(t, err)
	buyPsbt := "cHNidP8BALICAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD/////3Hae3avlOA0ThY1rsXFiaFygFIX3WBHy1cNXJe7cVKwBAAAAAP////8CAAAAAAAAAAAiUSDBJU2tQKrqEh5T4M7AmYY2L2zUIieW0ZSMlHkbNSgctRAnAAAAAAAAIlEg2TNXVWYEBtaRHMUNoSr0BbcH8iceiFfs/IUjes22cEcAAAAAAAEBKwAAAAAAAAAAIlEgwSVNrUCq6hIeU+DOwJmGNi9s1CInltGUjJR5GzUoHLUAAQErZQAAAAAAAAAiUSDZM1dVZgQG1pEcxQ2hKvQFtwfyJx6IV+z8hSN6zbZwRwEDBIMAAAABE0EyqMsv2PszUAkidLXerXjgrP1m/gUoTdbfK1xcaxuONjWKReHZH9FCal3gcQg0uYik6ukg4N7gyLwT9Ptb+szagwEXICmUSy4h/iIuagMY6RYcELQKfscDa/tGXSAKIJcrfVeeAAAA"
	bp, err := psbt.NewFromRawBytes(bytes.NewReader([]byte(buyPsbt)), true)
	assert.NoError(t, err)
	sp.Inputs[1] = bp.Inputs[1]
	for k, _ := range sp.Inputs {
		err = psbt.Finalize(sp, k)
		assert.NoError(t, err)
	}
	err = psbt.MaybeFinalizeAll(sp)
	assert.NoError(t, err)

	buyerSignedTx, err := psbt.Extract(sp)
	assert.NoError(t, err)
	var buf bytes.Buffer
	err = buyerSignedTx.Serialize(&buf)
	assert.NoError(t, err)
	fmt.Println("vsize", bitcoin.GetTxVirtualSize(btcutil.NewTx(buyerSignedTx)))
	fmt.Println(hex.EncodeToString(buf.Bytes()))
	assert.Equal(t, "0200000000010295774c6ff0e78483f5f76d53a5b98f67013a5a5d8f7dd6b8851badc917a77c3a0100000000ffffffffdc769eddabe5380d13858d6bb17162685ca01485f75811f2d5c35725eedc54ac0100000000ffffffff036500000000000000225120d9335755660406d6911cc50da12af405b707f2271e8857ecfc85237acdb670471027000000000000225120d9335755660406d6911cc50da12af405b707f2271e8857ecfc85237acdb6704796991a0000000000225120d9335755660406d6911cc50da12af405b707f2271e8857ecfc85237acdb6704701400f98d4c99042d361df7618a14ad9e65dbbea5c169b0194b4b8a301eb29c73d7b1d0a5fdd77c0e41139c7d649f28e68fb1f8a00fd5991d9b0a352b39d68fcf257014132a8cb2fd8fb3350092274b5dead78e0acfd66fe05284dd6df2b5c5c6b1b8e36358a45e1d91fd1426a5de0710834b988a4eae920e0dee0c8bc13f4fb5bfaccda8300000000", hex.EncodeToString(buf.Bytes()))
}

func TestPsbt(t *testing.T) {
	network := &chaincfg.TestNet3Params
	// seller
	txInput := &TxInput{
		TxId:       "ac54dcee2557c3d5f21158f78514a05c686271b16b8d85130d38e5abdd9e76dc",
		VOut:       uint32(1),
		Amount:     int64(101),
		Address:    "tb1pmye4w4txqsrddyguc5x6z2h5qkms0u38r6y90m8us53h4ndkwprst34fnw",
		PrivateKey: "cSWVEyJPTXLcNdEyAKzngz3diBXXEhAZHUyURzv2JsuUohopZkdE",
	}

	txOutput := &TxOutput{
		Address: "tb1pmye4w4txqsrddyguc5x6z2h5qkms0u38r6y90m8us53h4ndkwprst34fnw",
		Amount:  int64(10000),
	}

	sellerPsbt, err := GenerateRunesSignedListingPSBTBase64(txInput, txOutput, network)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("sellerPsbt", sellerPsbt)
	// buyer
	var inputs []*TxInput
	inputs = append(inputs, &TxInput{
		TxId:       "3a7ca717c9ad1b85b8d67d8f5d5a3a01678fb9a5536df7f58384e7f06f4c7795",
		VOut:       1,
		Amount:     1753509,
		Address:    "tb1pmye4w4txqsrddyguc5x6z2h5qkms0u38r6y90m8us53h4ndkwprst34fnw",
		PrivateKey: "cSWVEyJPTXLcNdEyAKzngz3diBXXEhAZHUyURzv2JsuUohopZkdE",
	})

	// seller input
	inputs = append(inputs, txInput)
	var outputs []*TxOutput
	outputs = append(outputs, &TxOutput{
		Address: "tb1pmye4w4txqsrddyguc5x6z2h5qkms0u38r6y90m8us53h4ndkwprst34fnw",
		Amount:  int64(101),
	})
	// seller output
	outputs = append(outputs, txOutput)
	outputs = append(outputs, &TxOutput{
		Address: "tb1pmye4w4txqsrddyguc5x6z2h5qkms0u38r6y90m8us53h4ndkwprst34fnw",
		Amount:  int64(0),
	})
	var sellerPsbts []string
	sellerPsbts = append(sellerPsbts, sellerPsbt)

	fee, buyerTx, err := GenerateRunesSignedBuyingTx(inputs, outputs, 1, 1, sellerPsbts, network)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fee)
	t.Log(buyerTx)
	assert.Equal(t, "cHNidP8BAN0CAAAAApV3TG/w54SD9fdtU6W5j2cBOlpdj33WuIUbrckXp3w6AQAAAAD/////3Hae3avlOA0ThY1rsXFiaFygFIX3WBHy1cNXJe7cVKwBAAAAAP////8DZQAAAAAAAAAiUSDZM1dVZgQG1pEcxQ2hKvQFtwfyJx6IV+z8hSN6zbZwRxAnAAAAAAAAIlEg2TNXVWYEBtaRHMUNoSr0BbcH8iceiFfs/IUjes22cEeWmRoAAAAAACJRINkzV1VmBAbWkRzFDaEq9AW3B/InHohX7PyFI3rNtnBHAAAAAAABASulwRoAAAAAACJRINkzV1VmBAbWkRzFDaEq9AW3B/InHohX7PyFI3rNtnBHAQMEAQAAAAETQA+Y1MmQQtNh33YYoUrZ5l276lwWmwGUtLijAespxz17HQpf3XfA5BE5x9ZJ8o5o+x+KAP1Zkdmwo1KznWj88lcBFyAplEsuIf4iLmoDGOkWHBC0Cn7HA2v7Rl0gCiCXK31XngABAStlAAAAAAAAACJRINkzV1VmBAbWkRzFDaEq9AW3B/InHohX7PyFI3rNtnBHAQMEgwAAAAETQDKoyy/Y+zNQCSJ0td6teOCs/Wb+BShN1t8rXFxrG442NYpF4dkf0UJqXeBxCDS5iKTq6SDg3uDIvBP0+1v6zNoBFyAplEsuIf4iLmoDGOkWHBC0Cn7HA2v7Rl0gCiCXK31XngAAAAA=",
		buyerTx)

}

func TestPsbtBatch(t *testing.T) {
	network := &chaincfg.TestNet3Params
	// seller
	txInput := &TxInput{
		TxId:       "7bc2683ddfa47d48c58be91c01ef3a5dad4cee975259c5130b44a87b014882fb",
		VOut:       uint32(0),
		Amount:     int64(1000),
		Address:    "tb1ppfc0mx9j3070zqleu257zt46ch2v9f9n9urkhlg7n7pswcmpqq0qt3pswx",
		PrivateKey: "cSWVEyJPTXLcNdEyAKzngz3diBXXEhAZHUyURzv2JsuUohopZkdE",
	}

	txOutput := &TxOutput{
		Address: "tb1ppfc0mx9j3070zqleu257zt46ch2v9f9n9urkhlg7n7pswcmpqq0qt3pswx",
		Amount:  int64(2000), // price
	}

	sellerPsbt, err := GenerateRunesSignedListingPSBTBase64(txInput, txOutput, network)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sellerPsbt)

	txInput2 := &TxInput{
		TxId:       "6abc1613438645b04435ac887e6e450f6ca57c3648e2091e968bc20a12e94a5e",
		VOut:       uint32(1),
		Amount:     int64(1000),
		Address:    "tb1ppfc0mx9j3070zqleu257zt46ch2v9f9n9urkhlg7n7pswcmpqq0qt3pswx",
		PrivateKey: "cSWVEyJPTXLcNdEyAKzngz3diBXXEhAZHUyURzv2JsuUohopZkdE",
	}

	txOutput2 := &TxOutput{
		Address: "tb1ppfc0mx9j3070zqleu257zt46ch2v9f9n9urkhlg7n7pswcmpqq0qt3pswx",
		Amount:  int64(2000), // price
	}
	sellerPsbt2, err := GenerateRunesSignedListingPSBTBase64(txInput2, txOutput2, network)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(sellerPsbt2)

	// buyer
	var inputs []*TxInput
	inputs = append(inputs, &TxInput{
		TxId:       "2ff5e2a38cb897ff6bf7708f15405bbed5b672308b404fb37e63b3b24b976b11",
		VOut:       2,
		Amount:     1351041,
		Address:    "tb1ppfc0mx9j3070zqleu257zt46ch2v9f9n9urkhlg7n7pswcmpqq0qt3pswx",
		PrivateKey: "cSWVEyJPTXLcNdEyAKzngz3diBXXEhAZHUyURzv2JsuUohopZkdE",
	})

	// seller input
	inputs = append(inputs, txInput)
	inputs = append(inputs, txInput2)

	var outputs []*TxOutput
	outputs = append(outputs, &TxOutput{
		Address: "tb1ppfc0mx9j3070zqleu257zt46ch2v9f9n9urkhlg7n7pswcmpqq0qt3pswx",
		Amount:  int64(547),
	})
	// seller output
	outputs = append(outputs, txOutput)
	outputs = append(outputs, txOutput2)

	outputs = append(outputs, &TxOutput{
		Address: "tb1ppfc0mx9j3070zqleu257zt46ch2v9f9n9urkhlg7n7pswcmpqq0qt3pswx",
		Amount:  int64(1),
	})

	var sellerPsbts []string
	sellerPsbts = append(sellerPsbts, sellerPsbt)
	sellerPsbts = append(sellerPsbts, sellerPsbt2)

	fee, buyerTx, err := GenerateRunesSignedBuyingTx(inputs, outputs, 1, 1, sellerPsbts, network)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fee)
	t.Log(buyerTx)
	assert.Equal(t, int64(355), fee)

	assert.Equal(t, "cHNidP8BAP0xAQIAAAADEWuXS7KzY36zT0CLMHK21b5bQBWPcPdr/5e4jKPi9S8CAAAAAP/////7gkgBe6hECxPFWVKX7kytXTrvARzpi8VIfaTfPWjCewAAAAAA/////15K6RIKwouWHgniSDZ8pWwPRW5+iKw1RLBFhkMTFrxqAQAAAAD/////BCMCAAAAAAAAIlEgCnD9mLKL/PED+eKp4S66xdTCpLMvB2v9Hp+DB2NhAB7QBwAAAAAAACJRIApw/Ziyi/zxA/niqeEuusXUwqSzLwdr/R6fgwdjYQAe0AcAAAAAAAAiUSAKcP2Ysov88QP54qnhLrrF1MKksy8Ha/0en4MHY2EAHiqSFAAAAAAAIlEgCnD9mLKL/PED+eKp4S66xdTCpLMvB2v9Hp+DB2NhAB4AAAAAAAEBK4GdFAAAAAAAIlEgCnD9mLKL/PED+eKp4S66xdTCpLMvB2v9Hp+DB2NhAB4BAwQBAAAAARNA+JSJqU8WbcA+FnXWHzI/pGLFdpCjwe5E2qCXXFQtU1XhQ5pFm22uRZKK3q7FGpQMpPMUJ6R4ZlsY8VinQrSa+gEXICmUSy4h/iIuagMY6RYcELQKfscDa/tGXSAKIJcrfVeeAAEBK+gDAAAAAAAAIlEgCnD9mLKL/PED+eKp4S66xdTCpLMvB2v9Hp+DB2NhAB4BAwSDAAAAARNAVgPw5EzeHCe9cLrDBt2KGF1G+cmKXpo5bxyIhBjOWFg18nSQBzV1syhU20uNyY7KVy1bqCdT/vVOlZOQ6Cu3NAEXICmUSy4h/iIuagMY6RYcELQKfscDa/tGXSAKIJcrfVeeAAEBK+gDAAAAAAAAIlEgCnD9mLKL/PED+eKp4S66xdTCpLMvB2v9Hp+DB2NhAB4BAwSDAAAAARNApU2osI9D7LsV7pImSDoy+xD7xLVnAZxvuXDNUev8mAswTIGfej+wNDb7lXQtROoVuTHI1GPfBuAy6LvPv39PsAEXICmUSy4h/iIuagMY6RYcELQKfscDa/tGXSAKIJcrfVeeAAAAAAA=", buyerTx)

}

func TestMergePsbtBatch(t *testing.T) {
	sellerPsbt := "cHNidP8BALICAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD/////+4JIAXuoRAsTxVlSl+5MrV067wEc6YvFSH2k3z1ownsAAAAAAP////8CAAAAAAAAAAAiUSAqfZ5VuWBDpp6/Ye7m7xjbUgvubGLNU95k/nksx3r2itAHAAAAAAAAIlEgCnD9mLKL/PED+eKp4S66xdTCpLMvB2v9Hp+DB2NhAB4AAAAAAAEBKwAAAAAAAAAAIlEgKn2eVblgQ6aev2Hu5u8Y21IL7mxizVPeZP55LMd69ooAAQEr6AMAAAAAAAAiUSAKcP2Ysov88QP54qnhLrrF1MKksy8Ha/0en4MHY2EAHgEDBIMAAAABE0FSvoLbzAWBUqq4YIsO7Az4GUz6l9jvmrxfDWi4a/WjBH2DXBUNDCURUlqpKkEE5/ubZu1CiHMWqMDi2NQDy6xmgwEXIJDh9zu/sGbDuLp0rwfERDXJ6wi99BwPO4YO2JQKt5GKAAAA"
	sp, err := psbt.NewFromRawBytes(bytes.NewReader([]byte(sellerPsbt)), true)
	assert.NoError(t, err)
	sellerPsbt2 := "cHNidP8BALICAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAD/////XkrpEgrCi5YeCeJINnylbA9Fbn6IrDVEsEWGQxMWvGoBAAAAAP////8CAAAAAAAAAAAiUSAqfZ5VuWBDpp6/Ye7m7xjbUgvubGLNU95k/nksx3r2itAHAAAAAAAAIlEgCnD9mLKL/PED+eKp4S66xdTCpLMvB2v9Hp+DB2NhAB4AAAAAAAEBKwAAAAAAAAAAIlEgKn2eVblgQ6aev2Hu5u8Y21IL7mxizVPeZP55LMd69ooAAQEr6AMAAAAAAAAiUSAKcP2Ysov88QP54qnhLrrF1MKksy8Ha/0en4MHY2EAHgEDBIMAAAABE0F4/MRQHw85k4q4Yir0ivRNVyzCJwFn3b0EqYGhtdvaTok7vBhg16NFCqojfRuNjNUyyQRdOGu6/8fkbX4qmKrbgwEXIJDh9zu/sGbDuLp0rwfERDXJ6wi99BwPO4YO2JQKt5GKAAAA"
	sp2, err := psbt.NewFromRawBytes(bytes.NewReader([]byte(sellerPsbt2)), true)
	assert.NoError(t, err)

	buyPsbt := "cHNidP8BAP0xAQIAAAADEWuXS7KzY36zT0CLMHK21b5bQBWPcPdr/5e4jKPi9S8CAAAAAP/////7gkgBe6hECxPFWVKX7kytXTrvARzpi8VIfaTfPWjCewAAAAAA/////15K6RIKwouWHgniSDZ8pWwPRW5+iKw1RLBFhkMTFrxqAQAAAAD/////BCMCAAAAAAAAIlEgCnD9mLKL/PED+eKp4S66xdTCpLMvB2v9Hp+DB2NhAB7QBwAAAAAAACJRIApw/Ziyi/zxA/niqeEuusXUwqSzLwdr/R6fgwdjYQAe0AcAAAAAAAAiUSAKcP2Ysov88QP54qnhLrrF1MKksy8Ha/0en4MHY2EAHiqSFAAAAAAAIlEgCnD9mLKL/PED+eKp4S66xdTCpLMvB2v9Hp+DB2NhAB4AAAAAAAEBK4GdFAAAAAAAIlEgCnD9mLKL/PED+eKp4S66xdTCpLMvB2v9Hp+DB2NhAB4BAwQBAAAAARNAijpaJm3+Aagvgzhbq0qDSeiEODAthWCyVrlZXypLzfU2+tYIUDdjfCxAY7tq7suYwaqwrvZEwgO+2BOhUPaoDAEXIJDh9zu/sGbDuLp0rwfERDXJ6wi99BwPO4YO2JQKt5GKAAEBK+gDAAAAAAAAIlEgCnD9mLKL/PED+eKp4S66xdTCpLMvB2v9Hp+DB2NhAB4BAwSDAAAAARNBUr6C28wFgVKquGCLDuwM+BlM+pfY75q8Xw1ouGv1owR9g1wVDQwlEVJaqSpBBOf7m2btQohzFqjA4tjUA8usZoMBFyCQ4fc7v7Bmw7i6dK8HxEQ1yesIvfQcDzuGDtiUCreRigABASvoAwAAAAAAACJRIApw/Ziyi/zxA/niqeEuusXUwqSzLwdr/R6fgwdjYQAeAQMEgwAAAAETQXj8xFAfDzmTirhiKvSK9E1XLMInAWfdvQSpgaG129pOiTu8GGDXo0UKqiN9G42M1TLJBF04a7r/x+RtfiqYqtuDARcgkOH3O7+wZsO4unSvB8RENcnrCL30HA87hg7YlAq3kYoAAAAAAA=="
	bp, err := psbt.NewFromRawBytes(bytes.NewReader([]byte(buyPsbt)), true)
	assert.NoError(t, err)
	bp.Inputs[1] = sp.Inputs[1]
	bp.Inputs[2] = sp2.Inputs[1]

	for k, _ := range bp.Inputs {
		err = psbt.Finalize(bp, k)
		assert.NoError(t, err)
	}
	err = psbt.MaybeFinalizeAll(bp)
	assert.NoError(t, err)

	buyerSignedTx, err := psbt.Extract(bp)
	assert.NoError(t, err)
	var buf bytes.Buffer
	err = buyerSignedTx.Serialize(&buf)
	assert.NoError(t, err)
	fmt.Println("vsize", bitcoin.GetTxVirtualSize(btcutil.NewTx(buyerSignedTx)))
	fmt.Println(hex.EncodeToString(buf.Bytes()))
	assert.Equal(t, "02000000000103116b974bb2b3637eb34f408b3072b6d5be5b40158f70f76bff97b88ca3e2f52f0200000000fffffffffb8248017ba8440b13c5595297ee4cad5d3aef011ce98bc5487da4df3d68c27b0000000000ffffffff5e4ae9120ac28b961e09e248367ca56c0f456e7e88ac3544b04586431316bc6a0100000000ffffffff0423020000000000002251200a70fd98b28bfcf103f9e2a9e12ebac5d4c2a4b32f076bfd1e9f83076361001ed0070000000000002251200a70fd98b28bfcf103f9e2a9e12ebac5d4c2a4b32f076bfd1e9f83076361001ed0070000000000002251200a70fd98b28bfcf103f9e2a9e12ebac5d4c2a4b32f076bfd1e9f83076361001e2a921400000000002251200a70fd98b28bfcf103f9e2a9e12ebac5d4c2a4b32f076bfd1e9f83076361001e01408a3a5a266dfe01a82f83385bab4a8349e88438302d8560b256b9595f2a4bcdf536fad6085037637c2c4063bb6aeecb98c1aab0aef644c203bed813a150f6a80c014152be82dbcc058152aab8608b0eec0cf8194cfa97d8ef9abc5f0d68b86bf5a3047d835c150d0c2511525aa92a4104e7fb9b66ed42887316a8c0e2d8d403cbac6683014178fcc4501f0f39938ab8622af48af44d572cc2270167ddbd04a981a1b5dbda4e893bbc1860d7a3450aaa237d1b8d8cd532c9045d386bbaffc7e46d7e2a98aadb8300000000", hex.EncodeToString(buf.Bytes()))
	assert.Equal(t, int64(356), bitcoin.GetTxVirtualSize(btcutil.NewTx(buyerSignedTx)))
}
