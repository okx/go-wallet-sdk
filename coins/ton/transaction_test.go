package ton

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/okx/go-wallet-sdk/coins/ton/ton/wallet"
)

var (
	expireAt           = int64(1719482102)
	testSeedHex        = "961e76f2e06f689b14002118a96d4a075fdb14f4e82aa3f150a5bd519aa077e9"
	testToAddress      = "UQC6QJ31Bv_hjmsoaUjRmpZYqj9NXbBbvufCNycnc0gjReqR"
	testAmount         = "100000000"
	testComment        = ""
	testMode           = uint8(3)
	testWalletAddress  = "UQCY_awZ5qi9Od-AyWm7XNKIdS1raUMoCvk0YsLCkTqfpnEb"
	testCodeSeqno0     = "te6cckECFAEAAtQAART/APSkE/S88sgLAQIBIAIDAgFIBAUE+PKDCNcYINMf0x/THwL4I7vyZO1E0NMf0x/T//QE0VFDuvKhUVG68qIF+QFUEGT5EPKj+AAkpMjLH1JAyx9SMMv/UhD0AMntVPgPAdMHIcAAn2xRkyDXSpbTB9QC+wDoMOAhwAHjACHAAuMAAcADkTDjDQOkyMsfEssfy/8GBwgJAubQAdDTAyFxsJJfBOAi10nBIJJfBOAC0x8hghBwbHVnvSKCEGRzdHK9sJJfBeAD+kAwIPpEAcjKB8v/ydDtRNCBAUDXIfQEMFyBAQj0Cm+hMbOSXwfgBdM/yCWCEHBsdWe6kjgw4w0DghBkc3RyupJfBuMNCgsCASAMDQBu0gf6ANTUIvkABcjKBxXL/8nQd3SAGMjLBcsCIs8WUAX6AhTLaxLMzMlz+wDIQBSBAQj0UfKnAgBwgQEI1xj6ANM/yFQgR4EBCPRR8qeCEG5vdGVwdIAYyMsFywJQBs8WUAT6AhTLahLLH8s/yXP7AAIAbIEBCNcY+gDTPzBSJIEBCPRZ8qeCEGRzdHJwdIAYyMsFywJQBc8WUAP6AhPLassfEss/yXP7AAAK9ADJ7VQAeAH6APQEMPgnbyIwUAqhIb7y4FCCEHBsdWeDHrFwgBhQBMsFJs8WWPoCGfQAy2kXyx9SYMs/IMmAQPsABgCKUASBAQj0WTDtRNCBAUDXIMgBzxb0AMntVAFysI4jghBkc3Rygx6xcIAYUAXLBVADzxYj+gITy2rLH8s/yYBA+wCSXwPiAgEgDg8AWb0kK29qJoQICga5D6AhhHDUCAhHpJN9KZEM5pA+n/mDeBKAG3gQFImHFZ8xhAIBWBARABG4yX7UTQ1wsfgAPbKd+1E0IEBQNch9AQwAsjKB8v/ydABgQEI9ApvoTGACASASEwAZrc52omhAIGuQ64X/wAAZrx32omhAEGuQ64WPwGb/qfE="
	testDataSeqno0     = "te6cckEBAQEAKwAAUQAAAAApqaMXESxRf3avxJCx2HaBkMoZwfu935e57ya0WDPIwb7UQMxAY1OUUA=="
	testNormHashSeqno0 = "d49d90840e8c92a28d699761841e1548a1f4cceca22ab7d271a0966963eb09fa"
	testNormHashSeqno5 = "ff0b11b768837424267e0c1a5b659b15efa18f08c8a0de0d0799572a6b1331ab"
)

func getTestSeed() []byte {
	seed, _ := hex.DecodeString(testSeedHex)
	return seed
}

func getTestPubKey() ed25519.PublicKey {
	seed := getTestSeed()
	return ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
}

func getTestTransferParams(seqno uint32, simulate bool) *TransferParams {
	return &TransferParams{
		Seed:     getTestSeed(),
		PubKey:   getTestPubKey(),
		Seqno:    seqno,
		ExpireAt: expireAt,
		Simulate: simulate,
		Version:  wallet.V4R2,
		To:       testToAddress,
		Amount:   testAmount,
		Comment:  testComment,
		Mode:     testMode,
		IsToken:  false,
	}
}

func TestTransfer(t *testing.T) {
	//Data te6cckEBAQEAKwAAUQAAAAApqaMXESxRf3avxJCx2HaBkMoZwfu935e57ya0WDPIwb7UQMxAY1OUUA==
	//Code te6cckECFAEAAtQAART/APSkE/S88sgLAQIBIAIDAgFIBAUE+PKDCNcYINMf0x/THwL4I7vyZO1E0NMf0x/T//QE0VFDuvKhUVG68qIF+QFUEGT5EPKj+AAkpMjLH1JAyx9SMMv/UhD0AMntVPgPAdMHIcAAn2xRkyDXSpbTB9QC+wDoMOAhwAHjACHAAuMAAcADkTDjDQOkyMsfEssfy/8GBwgJAubQAdDTAyFxsJJfBOAi10nBIJJfBOAC0x8hghBwbHVnvSKCEGRzdHK9sJJfBeAD+kAwIPpEAcjKB8v/ydDtRNCBAUDXIfQEMFyBAQj0Cm+hMbOSXwfgBdM/yCWCEHBsdWe6kjgw4w0DghBkc3RyupJfBuMNCgsCASAMDQBu0gf6ANTUIvkABcjKBxXL/8nQd3SAGMjLBcsCIs8WUAX6AhTLaxLMzMlz+wDIQBSBAQj0UfKnAgBwgQEI1xj6ANM/yFQgR4EBCPRR8qeCEG5vdGVwdIAYyMsFywJQBs8WUAT6AhTLahLLH8s/yXP7AAIAbIEBCNcY+gDTPzBSJIEBCPRZ8qeCEGRzdHJwdIAYyMsFywJQBc8WUAP6AhPLassfEss/yXP7AAAK9ADJ7VQAeAH6APQEMPgnbyIwUAqhIb7y4FCCEHBsdWeDHrFwgBhQBMsFJs8WWPoCGfQAy2kXyx9SYMs/IMmAQPsABgCKUASBAQj0WTDtRNCBAUDXIMgBzxb0AMntVAFysI4jghBkc3Rygx6xcIAYUAXLBVADzxYj+gITy2rLH8s/yYBA+wCSXwPiAgEgDg8AWb0kK29qJoQICga5D6AhhHDUCAhHpJN9KZEM5pA+n/mDeBKAG3gQFImHFZ8xhAIBWBARABG4yX7UTQ1wsfgAPbKd+1E0IEBQNch9AQwAsjKB8v/ydABgQEI9ApvoTGACASASEwAZrc52omhAIGuQ64X/wAAZrx32omhAEGuQ64WPwGb/qfE=
	seed, _ := hex.DecodeString("961e76f2e06f689b14002118a96d4a075fdb14f4e82aa3f150a5bd519aa077e9")
	to := "UQC6QJ31Bv_hjmsoaUjRmpZYqj9NXbBbvufCNycnc0gjReqR"
	amount := "100000000"
	comment := ""
	seqno := uint32(0)
	pubKey := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
	signedTx, err := Transfer(seed, pubKey, to, amount, comment, seqno, expireAt, 3, false, wallet.V4R2)
	assert.Nil(t, err)
	expect := "te6cckECFwEAA6wAA+OIATH7WDPNUXpzvwGS03a5pRDqWtbShlAV8mjFhYUidT9MEYyu6cZSXcE0TVy6ZQSUC8FA6Q5UQLAzYT0q9TCi1WEDvUWxNF7hVnAhFhhds+sryS4VzMWAKMMbYIxpiIrbu2HFNTRi7M+m3sAAAAAAAHABAgMBFP8A9KQT9LzyyAsEAFEAAAAAKamjFxEsUX92r8SQsdh2gZDKGcH7vd+Xue8mtFgzyMG+1EDMQABoQgBdIE76g3/wxzWUNKRozUssVR+mrtgt33PhG5OTuaQRoqAvrwgAAAAAAAAAAAAAAAAAAAIBIAUGAgFIBwgE+PKDCNcYINMf0x/THwL4I7vyZO1E0NMf0x/T//QE0VFDuvKhUVG68qIF+QFUEGT5EPKj+AAkpMjLH1JAyx9SMMv/UhD0AMntVPgPAdMHIcAAn2xRkyDXSpbTB9QC+wDoMOAhwAHjACHAAuMAAcADkTDjDQOkyMsfEssfy/8JCgsMAubQAdDTAyFxsJJfBOAi10nBIJJfBOAC0x8hghBwbHVnvSKCEGRzdHK9sJJfBeAD+kAwIPpEAcjKB8v/ydDtRNCBAUDXIfQEMFyBAQj0Cm+hMbOSXwfgBdM/yCWCEHBsdWe6kjgw4w0DghBkc3RyupJfBuMNDQ4CASAPEABu0gf6ANTUIvkABcjKBxXL/8nQd3SAGMjLBcsCIs8WUAX6AhTLaxLMzMlz+wDIQBSBAQj0UfKnAgBwgQEI1xj6ANM/yFQgR4EBCPRR8qeCEG5vdGVwdIAYyMsFywJQBs8WUAT6AhTLahLLH8s/yXP7AAIAbIEBCNcY+gDTPzBSJIEBCPRZ8qeCEGRzdHJwdIAYyMsFywJQBc8WUAP6AhPLassfEss/yXP7AAAK9ADJ7VQAeAH6APQEMPgnbyIwUAqhIb7y4FCCEHBsdWeDHrFwgBhQBMsFJs8WWPoCGfQAy2kXyx9SYMs/IMmAQPsABgCKUASBAQj0WTDtRNCBAUDXIMgBzxb0AMntVAFysI4jghBkc3Rygx6xcIAYUAXLBVADzxYj+gITy2rLH8s/yYBA+wCSXwPiAgEgERIAWb0kK29qJoQICga5D6AhhHDUCAhHpJN9KZEM5pA+n/mDeBKAG3gQFImHFZ8xhAIBWBMUABG4yX7UTQ1wsfgAPbKd+1E0IEBQNch9AQwAsjKB8v/ydABgQEI9ApvoTGACASAVFgAZrc52omhAIGuQ64X/wAAZrx32omhAEGuQ64WPwCaI0B4="
	assert.Equal(t, expect, signedTx.Tx)
	assert.Equal(t, "d49d90840e8c92a28d699761841e1548a1f4cceca22ab7d271a0966963eb09fa", signedTx.Hash)
}

func TestTransfer2(t *testing.T) {
	seed, _ := hex.DecodeString("961e76f2e06f689b14002118a96d4a075fdb14f4e82aa3f150a5bd519aa077e2")
	to := "UQC6QJ31Bv_hjmsoaUjRmpZYqj9NXbBbvufCNycnc0gjReqR"
	amount := "10000"
	comment := ""
	seqno := uint32(0)
	pubKey := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
	signedTx, err := Transfer(seed, pubKey, to, amount, comment, seqno, expireAt, 3, false, wallet.V4R2)
	assert.Nil(t, err)
	expect := "te6cckECFwEAA6oAA+OIAKR4w+Ra92rxd7NkYekkw/qWVVV/gMAUcSSMvEJ3Qn/GEYnpj8dYQpnru9Xoofh2mmV9dy+EO78HGlvBQ9tX6Ri5avIuLbrMqDN/mkSe4bNLSP9zJNEMciXdjELNqaNNYuGFNTRi7M+m3sAAAAAAAHABAgMBFP8A9KQT9LzyyAsEAFEAAAAAKamjF9Nb8UAZZyMmw8R9+HT3ASp834+WtVDTcVEDVTRvrfJYQABkQgBdIE76g3/wxzWUNKRozUssVR+mrtgt33PhG5OTuaQRopE4gAAAAAAAAAAAAAAAAAACASAFBgIBSAcIBPjygwjXGCDTH9Mf0x8C+CO78mTtRNDTH9Mf0//0BNFRQ7ryoVFRuvKiBfkBVBBk+RDyo/gAJKTIyx9SQMsfUjDL/1IQ9ADJ7VT4DwHTByHAAJ9sUZMg10qW0wfUAvsA6DDgIcAB4wAhwALjAAHAA5Ew4w0DpMjLHxLLH8v/CQoLDALm0AHQ0wMhcbCSXwTgItdJwSCSXwTgAtMfIYIQcGx1Z70ighBkc3RyvbCSXwXgA/pAMCD6RAHIygfL/8nQ7UTQgQFA1yH0BDBcgQEI9ApvoTGzkl8H4AXTP8glghBwbHVnupI4MOMNA4IQZHN0crqSXwbjDQ0OAgEgDxAAbtIH+gDU1CL5AAXIygcVy//J0Hd0gBjIywXLAiLPFlAF+gIUy2sSzMzJc/sAyEAUgQEI9FHypwIAcIEBCNcY+gDTP8hUIEeBAQj0UfKnghBub3RlcHSAGMjLBcsCUAbPFlAE+gIUy2oSyx/LP8lz+wACAGyBAQjXGPoA0z8wUiSBAQj0WfKnghBkc3RycHSAGMjLBcsCUAXPFlAD+gITy2rLHxLLP8lz+wAACvQAye1UAHgB+gD0BDD4J28iMFAKoSG+8uBQghBwbHVngx6xcIAYUATLBSbPFlj6Ahn0AMtpF8sfUmDLPyDJgED7AAYAilAEgQEI9Fkw7UTQgQFA1yDIAc8W9ADJ7VQBcrCOI4IQZHN0coMesXCAGFAFywVQA88WI/oCE8tqyx/LP8mAQPsAkl8D4gIBIBESAFm9JCtvaiaECAoGuQ+gIYRw1AgIR6STfSmRDOaQPp/5g3gSgBt4EBSJhxWfMYQCAVgTFAARuMl+1E0NcLH4AD2ynftRNCBAUDXIfQEMALIygfL/8nQAYEBCPQKb6ExgAgEgFRYAGa3OdqJoQCBrkOuF/8AAGa8d9qJoQBBrkOuFj8BWioQb"
	assert.Equal(t, expect, signedTx.Tx)
	assert.Equal(t, "f763badbe9235b8e1c1dc22b86a9ac82c8edad955ab5a58c54f53f97108f1de2", signedTx.Hash)
}

func TestTransferJetton(t *testing.T) {
	seed, _ := hex.DecodeString("961e76f2e06f689b14002118a96d4a075fdb14f4e82aa3f150a5bd519aa077e9")
	fromJettonAccount := "UQD7w9qG8Cq0PgX0hnp5nVpHPeDYL0QlhcLtjFMmna43sMxz"
	to := "UQC27fdnAFQcQDaXDrR89OKx-lW_Zyxuzcy5CjfPrS9A6vZf"
	amount := "1"
	seqno := uint32(0) // start 0
	pubKey := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
	signedTx, err := TransferJetton(seed, pubKey, fromJettonAccount, to, amount, 9, seqno, "", "1000000", "", "", "jetton test", 1, 1, false, wallet.V4R2)
	assert.Nil(t, err)
	expect := "te6ccgECGAEABBIAA+OIATH7WDPNUXpzvwGS03a5pRDqWtbShlAV8mjFhYUidT9MEY5pNbous8IBc0KISEiDna5zVywCAYWSTZjB8R61mR26xYovo9igzt7uphefSCwvTWGfb8HvOJDzV/ztx3jg7YBlNTRi4AAAACAAAAAAAHABAgMBFP8A9KQT9LzyyAsEAFEAAAAAKamjFxEsUX92r8SQsdh2gZDKGcH7vd+Xue8mtFgzyMG+1EDMQAFoQgB94e1DeBVaHwL6Qz08zq0jnvBsF6ISwuF2ximTTtcb2CAX14QAAAAAAAAAAAAAAAAAAQUCASAGBwDGD4p+pQAAAAAAAAABEBgBbdvuzgCoOIBtLh1o+enFY/Srfs5Y3ZuZchRvn1pegdUAJj9rBnmqL0534DJabtc0oh1LWtpQygK+TRiwsKROp+mGHoSAAAAAAGpldHRvbiB0ZXN0AgFICAkE+PKDCNcYINMf0x/THwL4I7vyZO1E0NMf0x/T//QE0VFDuvKhUVG68qIF+QFUEGT5EPKj+AAkpMjLH1JAyx9SMMv/UhD0AMntVPgPAdMHIcAAn2xRkyDXSpbTB9QC+wDoMOAhwAHjACHAAuMAAcADkTDjDQOkyMsfEssfy/8KCwwNAubQAdDTAyFxsJJfBOAi10nBIJJfBOAC0x8hghBwbHVnvSKCEGRzdHK9sJJfBeAD+kAwIPpEAcjKB8v/ydDtRNCBAUDXIfQEMFyBAQj0Cm+hMbOSXwfgBdM/yCWCEHBsdWe6kjgw4w0DghBkc3RyupJfBuMNDg8CASAQEQBu0gf6ANTUIvkABcjKBxXL/8nQd3SAGMjLBcsCIs8WUAX6AhTLaxLMzMlz+wDIQBSBAQj0UfKnAgBwgQEI1xj6ANM/yFQgR4EBCPRR8qeCEG5vdGVwdIAYyMsFywJQBs8WUAT6AhTLahLLH8s/yXP7AAIAbIEBCNcY+gDTPzBSJIEBCPRZ8qeCEGRzdHJwdIAYyMsFywJQBc8WUAP6AhPLassfEss/yXP7AAAK9ADJ7VQAeAH6APQEMPgnbyIwUAqhIb7y4FCCEHBsdWeDHrFwgBhQBMsFJs8WWPoCGfQAy2kXyx9SYMs/IMmAQPsABgCKUASBAQj0WTDtRNCBAUDXIMgBzxb0AMntVAFysI4jghBkc3Rygx6xcIAYUAXLBVADzxYj+gITy2rLH8s/yYBA+wCSXwPiAgEgEhMAWb0kK29qJoQICga5D6AhhHDUCAhHpJN9KZEM5pA+n/mDeBKAG3gQFImHFZ8xhAIBWBQVABG4yX7UTQ1wsfgAPbKd+1E0IEBQNch9AQwAsjKB8v/ydABgQEI9ApvoTGACASAWFwAZrc52omhAIGuQ64X/wAAZrx32omhAEGuQ64WPwA=="
	assert.Equal(t, expect, signedTx.Tx)
	assert.Equal(t, "857990b895733665f6d6f1478c4adf40dbbf0b24e5f0406e093f03547aef6601", signedTx.Hash)
}

func TestTransferMintlessJetton(t *testing.T) {
	seed, _ := hex.DecodeString("961e76f2e06f689b14002118a96d4a075fdb14f4e82aa3f150a5bd519aa077e9")
	fromJettonAccount := "UQD7w9qG8Cq0PgX0hnp5nVpHPeDYL0QlhcLtjFMmna43sMxz"
	to := "UQC27fdnAFQcQDaXDrR89OKx-lW_Zyxuzcy5CjfPrS9A6vZf"
	amount := "0"
	seqno := uint32(473) // start 0
	pubKey := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
	now := int64(1739482102)
	payload := "te6ccgECLwEABBIAAQgN9gLWAQlGA6+1FWXC4ss/wvDOFwMk2bVM97AUEWqaUhh63uWfQ26nAB4CIgWBcAIDBChIAQEZG2ZqtEYGAq27TvzHdGuGrhhKoICBU+Zg9Xq/qRMHGAAdIgEgBQYiASAHCChIAQEV0tdPcZG01smq0thhsmqf9ZzE0QqpP3c+ERvuHF1JDgAbKEgBAf3dO8qdKoPys7AWvavs1wMNWCOq5XashXaRopmksx/LABsiASAJCiIBIAsMKEgBAWP0xUs9JBrfQRl1FkF2tIfIDYpwLdf3fXqMi6BqxNtmABoiASANDihIAQFOErI5E7ld/nTAgHXdGI74UH8kxIaFyAkH42P54tEC9QAYIgEgDxAoSAEBrF16Czdlg18FB467CrR6Ucwxb8H+Z1e4qDeFWbkz1WEAFyhIAQEHi8q3zQEZfXYwkOhpSzAmR7QK0WmkucDeHV316c8c7QAWIgEgERIoSAEBR2lihAaD6vvRhL/JPmzaQggIKPhOYPFPz+gEDcJwGDEAFSIBIBMUIgEgFRYoSAEBGqLTCRJq5U9xF/0wg6m1Ofz5ajWar7G7OgFmuQNTEo8AEyhIAQHDhe4RDpUKfUvoxXabHqwgzzkeMRsUTXWathAf3eDi+gATIgEgFxgiASAZGihIAQE8uhwFLVCmW7tRL8CsXSPEyW4rSTSrAm3sjMnIORDhfAARIgEgGxwoSAEBYnXkuRcSaZ5EPsepRvY/G1DkCaHSzmWpyNWIW44xbrsAECIBIB0eKEgBAdKemJi3F0hYrXNqx18xFogtSbKwkITtmdyo3Z8iHIa9ABAiASAfIChIAQHe8C8s8ieQWjP9id65wykzkW78REHvWHVr12YzrDAGmQANIgEgISIoSAEBnDka+TCGgaCBt4hXmvS+P+zM6kh/mW16kZioyI5GZDkADChIAQG5/hysJjrmhyOYQIDn1fJm9RlaW/xAL06YgIqRG4V13wAMIgEgIyQoSAEBI+SIgvG4g8IIRDr0RlXbXGJmoxUkPY4mxyxeZe5EbMsACiIBICUmIgEgJygoSAEBma4khVS3h9wAGj5xLOuHvZdOTJhLxz9NuBbu6FLfUTIACShIAQEHUHwoLt9sVpcrmO+6JQ0PYjugfJ5GHOaWCKrzpPvNuQAHIgEgKSoiASArLChIAQHSAWkfu8dj0LoUIVlQtxsFatOZlM/EFv3N5bvK3IwnGQAFKEgBAR8z/6zOlECNrW9yAoHqhYCTTG5sN7xdM3AFFlWWyC4WAAUiASAtLihIAQEMtuRH10MNAMK1Nhd4P/9D/C/3KxQ6aWqXvx9q648OEgABAF66mxs4l01wvQnrWl46By9vMuy9cGycLoc7pgt8sy+0O5rKAAAAZuDdgAAAZvVa0A=="
	stateInit := "te6ccgEBAwEAjwACATQBAghCAg7xnhv0Dyukkvyqw4buylm/aCejhQcI2fzZrbaDq8M2AMoAgBRoNbY2cS6a4XoT1rS8dA5e3mXZeuDZOF0Od0wW+WZfcAPpn0MdzkzH7w8jwjgGMZfR3Y2FqlpEArXYKCy3B42gyr7UVZcLiyz/C8M4XAyTZtUz3sBQRappSGHre5Z9DbqcAQ=="

	signedTx, err := TransferJetton(seed, pubKey, fromJettonAccount, to, amount, 9, seqno, "100000000", "1", payload, stateInit, "", now, 1, false, wallet.V4R2)
	assert.Nil(t, err)
	expect := "te6ccgECNAEABZ4AAeGIATH7WDPNUXpzvwGS03a5pRDqWtbShlAV8mjFhYUidT9MBAbRy1grq1YaCD4Drxt/kYeqEqY1UvouoLpgZiX017UoQXpcfoSBXN8baOEyPzaJzO6uXCSo1Boj4T1SvStn0FFNTRi7PXMfsAAADsgAHAEDaUIAfeHtQ3gVWh8C+kM9PM6tI57wbBeiEsLhdsYpk07XG9ggL68IAAAAAAAAAAAAAAAAAAI2AgMECEICDvGeG/QPK6SS/KrDhu7KWb9oJ6OFBwjZ/NmttoOrwzYAygCAFGg1tjZxLprhehPWtLx0Dl7eZdl64Nk4XQ53TBb5Zl9wA+mfQx3OTMfvDyPCOAYxl9HdjYWqWkQCtdgoLLcHjaDKvtRVlwuLLP8LwzhcDJNm1TPewFBFqmlIYet7ln0NupwBAaIPin6lAAAAAAAAAAEIAW3b7s4AqDiAbS4daPnpxWP0q37OWN2bmXIUb59aXoHVACY/awZ5qi9Od+AyWm7XNKIdS1raUMoCvk0YsLCkTqfpogIFAQgN9gLWBglGA6+1FWXC4ss/wvDOFwMk2bVM97AUEWqaUhh63uWfQ26nAB4HIgWBcAIICShIAQEZG2ZqtEYGAq27TvzHdGuGrhhKoICBU+Zg9Xq/qRMHGAAdIgEgCgsiASAMDShIAQEV0tdPcZG01smq0thhsmqf9ZzE0QqpP3c+ERvuHF1JDgAbKEgBAf3dO8qdKoPys7AWvavs1wMNWCOq5XashXaRopmksx/LABsiASAODyIBIBARKEgBAWP0xUs9JBrfQRl1FkF2tIfIDYpwLdf3fXqMi6BqxNtmABoiASASEyhIAQFOErI5E7ld/nTAgHXdGI74UH8kxIaFyAkH42P54tEC9QAYIgEgFBUoSAEBrF16Czdlg18FB467CrR6Ucwxb8H+Z1e4qDeFWbkz1WEAFyhIAQEHi8q3zQEZfXYwkOhpSzAmR7QK0WmkucDeHV316c8c7QAWIgEgFhcoSAEBR2lihAaD6vvRhL/JPmzaQggIKPhOYPFPz+gEDcJwGDEAFSIBIBgZIgEgGhsoSAEBGqLTCRJq5U9xF/0wg6m1Ofz5ajWar7G7OgFmuQNTEo8AEyhIAQHDhe4RDpUKfUvoxXabHqwgzzkeMRsUTXWathAf3eDi+gATIgEgHB0iASAeHyhIAQE8uhwFLVCmW7tRL8CsXSPEyW4rSTSrAm3sjMnIORDhfAARIgEgICEoSAEBYnXkuRcSaZ5EPsepRvY/G1DkCaHSzmWpyNWIW44xbrsAECIBICIjKEgBAdKemJi3F0hYrXNqx18xFogtSbKwkITtmdyo3Z8iHIa9ABAiASAkJShIAQHe8C8s8ieQWjP9id65wykzkW78REHvWHVr12YzrDAGmQANIgEgJicoSAEBnDka+TCGgaCBt4hXmvS+P+zM6kh/mW16kZioyI5GZDkADChIAQG5/hysJjrmhyOYQIDn1fJm9RlaW/xAL06YgIqRG4V13wAMIgEgKCkoSAEBI+SIgvG4g8IIRDr0RlXbXGJmoxUkPY4mxyxeZe5EbMsACiIBICorIgEgLC0oSAEBma4khVS3h9wAGj5xLOuHvZdOTJhLxz9NuBbu6FLfUTIACShIAQEHUHwoLt9sVpcrmO+6JQ0PYjugfJ5GHOaWCKrzpPvNuQAHIgEgLi8iASAwMShIAQHSAWkfu8dj0LoUIVlQtxsFatOZlM/EFv3N5bvK3IwnGQAFKEgBAR8z/6zOlECNrW9yAoHqhYCTTG5sN7xdM3AFFlWWyC4WAAUiASAyMyhIAQEMtuRH10MNAMK1Nhd4P/9D/C/3KxQ6aWqXvx9q648OEgABAF66mxs4l01wvQnrWl46By9vMuy9cGycLoc7pgt8sy+0O5rKAAAAZuDdgAAAZvVa0A=="
	assert.Equal(t, expect, signedTx.Tx)
	assert.Equal(t, "e49673c3318cce940172df53edfee75ce3873906f69a4cfc209065211aa865e1", signedTx.Hash)
}

func TestCalTxHash(t *testing.T) {
	boc := "te6cckECCgEAAcMAAeGIAFj4Q0XnjCz5PaCmffX7LoD3fW0u23rLm2iS7OiPTb++Aa4Ryjk9tmvjfp0/T7qN/6XH2pvGPXWVA49KyUFwBTMaJQPaH5K7c+JF2EtW6pi+YAmaa9z1jty0XFk0ik09ECFNTRi7QzjRyAAAAAgAHAEEgWIAL3x2MA2Hkw8ANLbT8LnXO1PbxcFPjtQmAqnkxUPN9hegcPCv0AAAAAAAAAAAAAAAAACrze8TAAAAAAAAAADwAgMEBQID0AgGBwCAz5OLCwJUe6jC2zBOxTIAbNEbhl5Uyb9RzPXc23q5rvUKIeaWu+U+GA3X2x/2JMYzN9aG8+wWaDj1tKgJqRMxDACAJomsdXOD9V1dherAPhvpta+dSdW5lZadUxeQsZ9cut8rP2o5Oz9WYSaCtL6T5zmXuBRnEb+j11levBeOUxcZAABDgAKjUGeyXpSFlLMK+sUSa+N+F055IaK07JfoUwXa+lbLSAEBIAgBASAJAE6ABuBXrm5cCldCCudzpcXEdP7Hd9o3iQ7TZAdYv544RXiqAlQL5AAAToAa18Cpqpt0vH1t6e1WvLApomKv5aCwyokOInEYhKrZfKoFloLwADrVB6o="
	hash, err := CalTxHash(boc)
	assert.Nil(t, err)
	assert.Equal(t, "2dc368fb446f6e2251bc408f673e40d730a77cf9632d70235f365936042a99da", hash)

	normHash, err := CalNormMsgHash(boc)
	assert.Nil(t, err)
	assert.Equal(t, "716b7272ac30e7afb0819a52646da0445c009386db2a4f2e64164468afbd7291", normHash)

	boc = "te6ccgECNAEABX8AAZwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAKamjFwAAASwAAACzAAMBA2lCAAa9gOKN9IIU6LsyeemRMBIm9e1iSSXxOXoCB7iEYDlqIC+vCAAAAAAAAAAAAAAAAAACNgIDBAhCAg7xnhv0Dyukkvyqw4buylm/aCejhQcI2fzZrbaDq8M2AMoAgBRoNbY2cS6a4XoT1rS8dA5e3mXZeuDZOF0Od0wW+WZfcAPpn0MdzkzH7w8jwjgGMZfR3Y2FqlpEArXYKCy3B42gyr7UVZcLiyz/C8M4XAyTZtUz3sBQRappSGHre5Z9DbqcAQGqD4p+pQAAAAAAACcQQ7msoAgBJwKmyORh2SwlHeZ3XdOVRwU45FMWjK+T9PRYw7yswQkAKNBrbGziXTXC9CetaXjoHL28y7L1wbJwuhzumC3yzL7iAgUBCA32AtYGCUYDr7UVZcLiyz/C8M4XAyTZtUz3sBQRappSGHre5Z9DbqcAHgciBYFwAggJKEgBARkbZmq0RgYCrbtO/Md0a4auGEqggIFT5mD1er+pEwcYAB0iASAKCyIBIAwNKEgBARXS109xkbTWyarS2GGyap/1nMTRCqk/dz4RG+4cXUkOABsoSAEB/d07yp0qg/KzsBa9q+zXAw1YI6rldqyFdpGimaSzH8sAGyIBIA4PIgEgEBEoSAEBY/TFSz0kGt9BGXUWQXa0h8gNinAt1/d9eoyLoGrE22YAGiIBIBITKEgBAU4SsjkTuV3+dMCAdd0YjvhQfyTEhoXICQfjY/ni0QL1ABgiASAUFShIAQGsXXoLN2WDXwUHjrsKtHpRzDFvwf5nV7ioN4VZuTPVYQAXKEgBAQeLyrfNARl9djCQ6GlLMCZHtArRaaS5wN4dXfXpzxztABYiASAWFyhIAQFHaWKEBoPq+9GEv8k+bNpCCAgo+E5g8U/P6AQNwnAYMQAVIgEgGBkiASAaGyhIAQEaotMJEmrlT3EX/TCDqbU5/PlqNZqvsbs6AWa5A1MSjwATKEgBAcOF7hEOlQp9S+jFdpserCDPOR4xGxRNdZq2EB/d4OL6ABMiASAcHSIBIB4fKEgBATy6HAUtUKZbu1EvwKxdI8TJbitJNKsCbeyMycg5EOF8ABEiASAgIShIAQFideS5FxJpnkQ+x6lG9j8bUOQJodLOZanI1YhbjjFuuwAQIgEgIiMoSAEB0p6YmLcXSFitc2rHXzEWiC1JsrCQhO2Z3KjdnyIchr0AECIBICQlKEgBAd7wLyzyJ5BaM/2J3rnDKTORbvxEQe9YdWvXZjOsMAaZAA0iASAmJyhIAQGcORr5MIaBoIG3iFea9L4/7MzqSH+ZbXqRmKjIjkZkOQAMKEgBAbn+HKwmOuaHI5hAgOfV8mb1GVpb/EAvTpiAipEbhXXfAAwiASAoKShIAQEj5IiC8biDwghEOvRGVdtcYmajFSQ9jibHLF5l7kRsywAKIgEgKisiASAsLShIAQGZriSFVLeH3AAaPnEs64e9l05MmEvHP024Fu7oUt9RMgAJKEgBAQdQfCgu32xWlyuY77olDQ9iO6B8nkYc5pYIqvOk+825AAciASAuLyIBIDAxKEgBAdIBaR+7x2PQuhQhWVC3GwVq05mUz8QW/c3lu8rcjCcZAAUoSAEBHzP/rM6UQI2tb3ICgeqFgJNMbmw3vF0zcAUWVZbILhYABSIBIDIzKEgBAQy25EfXQw0AwrU2F3g//0P8L/crFDppape/H2rrjw4SAAEAXrqbGziXTXC9CetaXjoHL28y7L1wbJwuhzumC3yzL7Q7msoAAABm4N2AAABm9VrQ"
	hash, err = CalTxHash(boc)
	assert.Nil(t, err)
	assert.Equal(t, "f775cbee53479fb134ba1ae5a152bca07c8f9601d3196a0baa219a85effd4dde", hash)

	normHash, err = CalNormMsgHash(boc)
	assert.Error(t, err)
	assert.Equal(t, "", normHash)
}

func TestBuildTransferSigningHash(t *testing.T) {
	seed, _ := hex.DecodeString("961e76f2e06f689b14002118a96d4a075fdb14f4e82aa3f150a5bd519aa077e9")
	pubKey := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
	to := "UQC6QJ31Bv_hjmsoaUjRmpZYqj9NXbBbvufCNycnc0gjReqR"
	amount := "100000000"
	comment := ""
	seqno := uint32(0)

	params := &TransferParams{
		Seed:     seed,
		PubKey:   pubKey,
		Seqno:    seqno,
		ExpireAt: expireAt,
		Simulate: false,
		Version:  wallet.V4R2,
		To:       to,
		Amount:   amount,
		Comment:  comment,
		Mode:     3,
		IsToken:  false,
	}

	builder, err := NewTonTransferBuilder(params)
	assert.Nil(t, err)

	// Test error case: no messages
	hash, err := builder.BuildTransferSigningHash()
	assert.NotNil(t, err)
	assert.Nil(t, hash)
	assert.Equal(t, "no messages", err.Error())

	// Build a message first
	_, err = builder.BuildTonMessage()
	assert.Nil(t, err)

	// Test success case
	hash, err = builder.BuildTransferSigningHash()
	assert.Nil(t, err)
	assert.NotNil(t, hash)
	assert.Equal(t, 32, len(hash))    // SHA256 hash is 32 bytes
	assert.NotNil(t, builder.payload) // payload should be set
}

func TestBuildTransferWithSignature(t *testing.T) {
	seed, _ := hex.DecodeString("961e76f2e06f689b14002118a96d4a075fdb14f4e82aa3f150a5bd519aa077e9")
	pubKey := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
	to := "UQC6QJ31Bv_hjmsoaUjRmpZYqj9NXbBbvufCNycnc0gjReqR"
	amount := "100000000"
	comment := ""
	seqno := uint32(0)

	params := &TransferParams{
		Seed:     seed,
		PubKey:   pubKey,
		Seqno:    seqno,
		ExpireAt: expireAt,
		Simulate: false,
		Version:  wallet.V4R2,
		To:       to,
		Amount:   amount,
		Comment:  comment,
		Mode:     3,
		IsToken:  false,
	}

	builder, err := NewTonTransferBuilder(params)
	assert.Nil(t, err)

	// Test error case: payload is not set
	signature := make([]byte, 64) // 64 bytes for ed25519 signature
	msg, err := builder.BuildTransferWithSignature(signature)
	assert.NotNil(t, err)
	assert.Nil(t, msg)
	assert.Equal(t, "payload is not set", err.Error())

	// Build message and signing hash first
	_, err = builder.BuildTonMessage()
	assert.Nil(t, err)
	_, err = builder.BuildTransferSigningHash()
	assert.Nil(t, err)

	// Test success case with valid signature
	msg, err = builder.BuildTransferWithSignature(signature)
	assert.Nil(t, err)
	assert.NotNil(t, msg)
	assert.NotNil(t, builder.externalMessage) // externalMessage should be set
	assert.Equal(t, msg, builder.externalMessage)
}

func TestBuildSignedTx(t *testing.T) {
	// Test error case: external message is not set
	params := getTestTransferParams(0, false)
	builder, err := NewTonTransferBuilder(params)
	assert.Nil(t, err)

	signedTx, err := builder.BuildSignedTx(true)
	assert.NotNil(t, err)
	assert.Nil(t, signedTx)
	assert.Equal(t, "external message is not set", err.Error())

	// Table-driven tests for various combinations
	tests := []struct {
		name             string
		seqno            uint32
		simulate         bool
		useBOCWithFlags  bool
		expectedAddress  string
		expectedCode     string
		expectedData     string
		expectedNormHash string
		expectedTx       string
	}{
		{
			name:             "seqno=0, simulate=false, useBOCWithFlags=true",
			seqno:            0,
			simulate:         false,
			useBOCWithFlags:  true,
			expectedAddress:  testWalletAddress,
			expectedCode:     testCodeSeqno0,
			expectedData:     testDataSeqno0,
			expectedNormHash: testNormHashSeqno0,
			expectedTx:       "te6cckECFwEAA6wAA+OIATH7WDPNUXpzvwGS03a5pRDqWtbShlAV8mjFhYUidT9MEYyu6cZSXcE0TVy6ZQSUC8FA6Q5UQLAzYT0q9TCi1WEDvUWxNF7hVnAhFhhds+sryS4VzMWAKMMbYIxpiIrbu2HFNTRi7M+m3sAAAAAAAHABAgMBFP8A9KQT9LzyyAsEAFEAAAAAKamjFxEsUX92r8SQsdh2gZDKGcH7vd+Xue8mtFgzyMG+1EDMQABoQgBdIE76g3/wxzWUNKRozUssVR+mrtgt33PhG5OTuaQRoqAvrwgAAAAAAAAAAAAAAAAAAAIBIAUGAgFIBwgE+PKDCNcYINMf0x/THwL4I7vyZO1E0NMf0x/T//QE0VFDuvKhUVG68qIF+QFUEGT5EPKj+AAkpMjLH1JAyx9SMMv/UhD0AMntVPgPAdMHIcAAn2xRkyDXSpbTB9QC+wDoMOAhwAHjACHAAuMAAcADkTDjDQOkyMsfEssfy/8JCgsMAubQAdDTAyFxsJJfBOAi10nBIJJfBOAC0x8hghBwbHVnvSKCEGRzdHK9sJJfBeAD+kAwIPpEAcjKB8v/ydDtRNCBAUDXIfQEMFyBAQj0Cm+hMbOSXwfgBdM/yCWCEHBsdWe6kjgw4w0DghBkc3RyupJfBuMNDQ4CASAPEABu0gf6ANTUIvkABcjKBxXL/8nQd3SAGMjLBcsCIs8WUAX6AhTLaxLMzMlz+wDIQBSBAQj0UfKnAgBwgQEI1xj6ANM/yFQgR4EBCPRR8qeCEG5vdGVwdIAYyMsFywJQBs8WUAT6AhTLahLLH8s/yXP7AAIAbIEBCNcY+gDTPzBSJIEBCPRZ8qeCEGRzdHJwdIAYyMsFywJQBc8WUAP6AhPLassfEss/yXP7AAAK9ADJ7VQAeAH6APQEMPgnbyIwUAqhIb7y4FCCEHBsdWeDHrFwgBhQBMsFJs8WWPoCGfQAy2kXyx9SYMs/IMmAQPsABgCKUASBAQj0WTDtRNCBAUDXIMgBzxb0AMntVAFysI4jghBkc3Rygx6xcIAYUAXLBVADzxYj+gITy2rLH8s/yYBA+wCSXwPiAgEgERIAWb0kK29qJoQICga5D6AhhHDUCAhHpJN9KZEM5pA+n/mDeBKAG3gQFImHFZ8xhAIBWBMUABG4yX7UTQ1wsfgAPbKd+1E0IEBQNch9AQwAsjKB8v/ydABgQEI9ApvoTGACASAVFgAZrc52omhAIGuQ64X/wAAZrx32omhAEGuQ64WPwCaI0B4=",
		},
		{
			name:             "seqno=5, simulate=false, useBOCWithFlags=true",
			seqno:            5,
			simulate:         false,
			useBOCWithFlags:  true,
			expectedAddress:  testWalletAddress,
			expectedCode:     "",
			expectedData:     "",
			expectedNormHash: testNormHashSeqno5,
			expectedTx:       "te6cckEBAgEAqgAB4YgBMftYM81RenO/AZLTdrmlEOpa1tKGUBXyaMWFhSJ1P0wH7m9RYi42M605i4bf0B17bJm4yYskIj/MPjaCoYJw9DzO/ZFO1QK2r/oLWudZIVlgqXWO1YOd/5LhqI4MrEOAcU1NGLsz6bewAAAAKAAcAQBoQgBdIE76g3/wxzWUNKRozUssVR+mrtgt33PhG5OTuaQRoqAvrwgAAAAAAAAAAAAAAAAAAPEs5fk=",
		},
		{
			name:             "seqno=0, simulate=true, useBOCWithFlags=true",
			seqno:            0,
			simulate:         true,
			useBOCWithFlags:  true,
			expectedAddress:  testWalletAddress,
			expectedCode:     testCodeSeqno0,
			expectedData:     testDataSeqno0,
			expectedNormHash: "",
			expectedTx:       "te6cckEBAgEAhwABnGV3TjKS7gmiauXTKCSgXgoHSHKiBYGbCelXqYUWqwgd6i2JovcKs4EIsMLtn1leSXCuZiwBRhjbBGNMRFbd2w4pqaMXZn029gAAAAAAAwEAaEIAXSBO+oN/8Mc1lDSkaM1LLFUfpq7YLd9z4RuTk7mkEaKgL68IAAAAAAAAAAAAAAAAAAAj6dYr",
		},
		{
			name:             "seqno=10, simulate=true, useBOCWithFlags=true",
			seqno:            10,
			simulate:         true,
			useBOCWithFlags:  true,
			expectedAddress:  testWalletAddress,
			expectedCode:     "",
			expectedData:     "",
			expectedNormHash: "",
			expectedTx:       "te6cckEBAgEAhwABnNaRiInLDjFGCZ9+A5itbau22guNqoV3+nqnxGYKzdAvPYf+7YodsCdfCKLXUG6LwMa17la+jYDsnksvG8PqLQ4pqaMXZn029gAAAAoAAwEAaEIAXSBO+oN/8Mc1lDSkaM1LLFUfpq7YLd9z4RuTk7mkEaKgL68IAAAAAAAAAAAAAAAAAADpbKiY",
		},
		{
			name:             "seqno=0, simulate=false, useBOCWithFlags=false",
			seqno:            0,
			simulate:         false,
			useBOCWithFlags:  false,
			expectedAddress:  testWalletAddress,
			expectedCode:     testCodeSeqno0,
			expectedData:     testDataSeqno0,
			expectedNormHash: testNormHashSeqno0,
			expectedTx:       "te6ccgECFwEAA6wAA+OIATH7WDPNUXpzvwGS03a5pRDqWtbShlAV8mjFhYUidT9MEYyu6cZSXcE0TVy6ZQSUC8FA6Q5UQLAzYT0q9TCi1WEDvUWxNF7hVnAhFhhds+sryS4VzMWAKMMbYIxpiIrbu2HFNTRi7M+m3sAAAAAAAHABAgMBFP8A9KQT9LzyyAsEAFEAAAAAKamjFxEsUX92r8SQsdh2gZDKGcH7vd+Xue8mtFgzyMG+1EDMQABoQgBdIE76g3/wxzWUNKRozUssVR+mrtgt33PhG5OTuaQRoqAvrwgAAAAAAAAAAAAAAAAAAAIBIAUGAgFIBwgE+PKDCNcYINMf0x/THwL4I7vyZO1E0NMf0x/T//QE0VFDuvKhUVG68qIF+QFUEGT5EPKj+AAkpMjLH1JAyx9SMMv/UhD0AMntVPgPAdMHIcAAn2xRkyDXSpbTB9QC+wDoMOAhwAHjACHAAuMAAcADkTDjDQOkyMsfEssfy/8JCgsMAubQAdDTAyFxsJJfBOAi10nBIJJfBOAC0x8hghBwbHVnvSKCEGRzdHK9sJJfBeAD+kAwIPpEAcjKB8v/ydDtRNCBAUDXIfQEMFyBAQj0Cm+hMbOSXwfgBdM/yCWCEHBsdWe6kjgw4w0DghBkc3RyupJfBuMNDQ4CASAPEABu0gf6ANTUIvkABcjKBxXL/8nQd3SAGMjLBcsCIs8WUAX6AhTLaxLMzMlz+wDIQBSBAQj0UfKnAgBwgQEI1xj6ANM/yFQgR4EBCPRR8qeCEG5vdGVwdIAYyMsFywJQBs8WUAT6AhTLahLLH8s/yXP7AAIAbIEBCNcY+gDTPzBSJIEBCPRZ8qeCEGRzdHJwdIAYyMsFywJQBc8WUAP6AhPLassfEss/yXP7AAAK9ADJ7VQAeAH6APQEMPgnbyIwUAqhIb7y4FCCEHBsdWeDHrFwgBhQBMsFJs8WWPoCGfQAy2kXyx9SYMs/IMmAQPsABgCKUASBAQj0WTDtRNCBAUDXIMgBzxb0AMntVAFysI4jghBkc3Rygx6xcIAYUAXLBVADzxYj+gITy2rLH8s/yYBA+wCSXwPiAgEgERIAWb0kK29qJoQICga5D6AhhHDUCAhHpJN9KZEM5pA+n/mDeBKAG3gQFImHFZ8xhAIBWBMUABG4yX7UTQ1wsfgAPbKd+1E0IEBQNch9AQwAsjKB8v/ydABgQEI9ApvoTGACASAVFgAZrc52omhAIGuQ64X/wAAZrx32omhAEGuQ64WPwA==",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := getTestTransferParams(tt.seqno, tt.simulate)
			builder, err := NewTonTransferBuilder(params)
			assert.Nil(t, err)

			// Build external message
			_, err = builder.BuildTonMessage()
			assert.Nil(t, err)
			_, err = builder.BuildTransferDirect()
			assert.Nil(t, err)

			// Build signed tx
			signedTx, err := builder.BuildSignedTx(tt.useBOCWithFlags)
			assert.Nil(t, err)
			assert.NotNil(t, signedTx)

			// Verify address
			assert.Equal(t, tt.expectedAddress, signedTx.Address)

			// Verify Code (only when seqno = 0)
			assert.Equal(t, tt.expectedCode, signedTx.Code)

			// Verify Data (only when seqno = 0)
			assert.Equal(t, tt.expectedData, signedTx.Data)

			// Verify NormHash (only when simulate = false)
			assert.Equal(t, tt.expectedNormHash, signedTx.Hash)

			// Verify Tx
			assert.Equal(t, tt.expectedTx, signedTx.Tx)
		})
	}

	// Test BOC flags difference
	t.Run("BOC flags produce different encodings", func(t *testing.T) {
		params := getTestTransferParams(0, false)
		builder1, _ := NewTonTransferBuilder(params)
		builder1.BuildTonMessage()
		builder1.BuildTransferDirect()
		signedTx1, _ := builder1.BuildSignedTx(false)

		builder2, _ := NewTonTransferBuilder(params)
		builder2.BuildTonMessage()
		builder2.BuildTransferDirect()
		signedTx2, _ := builder2.BuildSignedTx(true)

		assert.NotEqual(t, signedTx1.Tx, signedTx2.Tx)
	})
}
