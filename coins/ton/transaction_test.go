package ton

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	expireAt = int64(1719482102)
)

func TestTransfer(t *testing.T) {
	seed, _ := hex.DecodeString("45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40")
	address, err := NewAddress(seed)
	fmt.Println(address)
	to := "UQC6QJ31Bv_hjmsoaUjRmpZYqj9NXbBbvufCNycnc0gjReqR"
	amount := "100000000"
	comment := ""
	seqno := uint32(0)
	pubKey := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
	signedTx, err := Transfer(seed, pubKey, to, amount, comment, seqno, expireAt, 3, false)
	assert.Nil(t, err)
	t.Log(signedTx.Tx)
	expect := "te6cckECFwEAA6wAA+OIAXkNjksffYh8JIQ426pdFemct+pPeD5Cxlw/0NF/aTcAEZGlf6R72mNUWwgiDUiI+hyrC1r1udyECFNcwKRQmie+LUYK0fpwG9QeRdxkHQMEokNS3VflTzPLCJu4WiykrYAlNTRi7M+m3sAAAAAAAHABAgMBFP8A9KQT9LzyyAsEAFEAAAAAKamjFwwvPG2rtKBgDsyuh66qOSQgQvmldqqNygHhtBnPF9eiQABoQgBdIE76g3/wxzWUNKRozUssVR+mrtgt33PhG5OTuaQRoqAvrwgAAAAAAAAAAAAAAAAAAAIBIAUGAgFIBwgE+PKDCNcYINMf0x/THwL4I7vyZO1E0NMf0x/T//QE0VFDuvKhUVG68qIF+QFUEGT5EPKj+AAkpMjLH1JAyx9SMMv/UhD0AMntVPgPAdMHIcAAn2xRkyDXSpbTB9QC+wDoMOAhwAHjACHAAuMAAcADkTDjDQOkyMsfEssfy/8JCgsMAubQAdDTAyFxsJJfBOAi10nBIJJfBOAC0x8hghBwbHVnvSKCEGRzdHK9sJJfBeAD+kAwIPpEAcjKB8v/ydDtRNCBAUDXIfQEMFyBAQj0Cm+hMbOSXwfgBdM/yCWCEHBsdWe6kjgw4w0DghBkc3RyupJfBuMNDQ4CASAPEABu0gf6ANTUIvkABcjKBxXL/8nQd3SAGMjLBcsCIs8WUAX6AhTLaxLMzMlz+wDIQBSBAQj0UfKnAgBwgQEI1xj6ANM/yFQgR4EBCPRR8qeCEG5vdGVwdIAYyMsFywJQBs8WUAT6AhTLahLLH8s/yXP7AAIAbIEBCNcY+gDTPzBSJIEBCPRZ8qeCEGRzdHJwdIAYyMsFywJQBc8WUAP6AhPLassfEss/yXP7AAAK9ADJ7VQAeAH6APQEMPgnbyIwUAqhIb7y4FCCEHBsdWeDHrFwgBhQBMsFJs8WWPoCGfQAy2kXyx9SYMs/IMmAQPsABgCKUASBAQj0WTDtRNCBAUDXIMgBzxb0AMntVAFysI4jghBkc3Rygx6xcIAYUAXLBVADzxYj+gITy2rLH8s/yYBA+wCSXwPiAgEgERIAWb0kK29qJoQICga5D6AhhHDUCAhHpJN9KZEM5pA+n/mDeBKAG3gQFImHFZ8xhAIBWBMUABG4yX7UTQ1wsfgAPbKd+1E0IEBQNch9AQwAsjKB8v/ydABgQEI9ApvoTGACASAVFgAZrc52omhAIGuQ64X/wAAZrx32omhAEGuQ64WPwGwOSsk="
	assert.Equal(t, expect, signedTx.Tx)
	t.Log(signedTx.Hash)
	assert.Equal(t, "0ddcbf78f63bdd2bc6d11ff5bd79213d337b748ea4c536cba3b1cce57b21b7b3", signedTx.Hash)
}

func TestTransfer2(t *testing.T) {
	seed, _ := hex.DecodeString("961e76f2e06f689b14002118a96d4a075fdb14f4e82aa3f150a5bd519aa077e2")
	address, err := NewAddress(seed)
	fmt.Println(address)
	to := "UQC6QJ31Bv_hjmsoaUjRmpZYqj9NXbBbvufCNycnc0gjReqR"
	amount := "10000"
	comment := ""
	seqno := uint32(0)
	pubKey := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
	signedTx, err := Transfer(seed, pubKey, to, amount, comment, seqno, expireAt, 3, false)
	assert.Nil(t, err)
	t.Log(signedTx.Tx)
	expect := "te6cckECFwEAA6oAA+OIAKR4w+Ra92rxd7NkYekkw/qWVVV/gMAUcSSMvEJ3Qn/GEYnpj8dYQpnru9Xoofh2mmV9dy+EO78HGlvBQ9tX6Ri5avIuLbrMqDN/mkSe4bNLSP9zJNEMciXdjELNqaNNYuGFNTRi7M+m3sAAAAAAAHABAgMBFP8A9KQT9LzyyAsEAFEAAAAAKamjF9Nb8UAZZyMmw8R9+HT3ASp834+WtVDTcVEDVTRvrfJYQABkQgBdIE76g3/wxzWUNKRozUssVR+mrtgt33PhG5OTuaQRopE4gAAAAAAAAAAAAAAAAAACASAFBgIBSAcIBPjygwjXGCDTH9Mf0x8C+CO78mTtRNDTH9Mf0//0BNFRQ7ryoVFRuvKiBfkBVBBk+RDyo/gAJKTIyx9SQMsfUjDL/1IQ9ADJ7VT4DwHTByHAAJ9sUZMg10qW0wfUAvsA6DDgIcAB4wAhwALjAAHAA5Ew4w0DpMjLHxLLH8v/CQoLDALm0AHQ0wMhcbCSXwTgItdJwSCSXwTgAtMfIYIQcGx1Z70ighBkc3RyvbCSXwXgA/pAMCD6RAHIygfL/8nQ7UTQgQFA1yH0BDBcgQEI9ApvoTGzkl8H4AXTP8glghBwbHVnupI4MOMNA4IQZHN0crqSXwbjDQ0OAgEgDxAAbtIH+gDU1CL5AAXIygcVy//J0Hd0gBjIywXLAiLPFlAF+gIUy2sSzMzJc/sAyEAUgQEI9FHypwIAcIEBCNcY+gDTP8hUIEeBAQj0UfKnghBub3RlcHSAGMjLBcsCUAbPFlAE+gIUy2oSyx/LP8lz+wACAGyBAQjXGPoA0z8wUiSBAQj0WfKnghBkc3RycHSAGMjLBcsCUAXPFlAD+gITy2rLHxLLP8lz+wAACvQAye1UAHgB+gD0BDD4J28iMFAKoSG+8uBQghBwbHVngx6xcIAYUATLBSbPFlj6Ahn0AMtpF8sfUmDLPyDJgED7AAYAilAEgQEI9Fkw7UTQgQFA1yDIAc8W9ADJ7VQBcrCOI4IQZHN0coMesXCAGFAFywVQA88WI/oCE8tqyx/LP8mAQPsAkl8D4gIBIBESAFm9JCtvaiaECAoGuQ+gIYRw1AgIR6STfSmRDOaQPp/5g3gSgBt4EBSJhxWfMYQCAVgTFAARuMl+1E0NcLH4AD2ynftRNCBAUDXIfQEMALIygfL/8nQAYEBCPQKb6ExgAgEgFRYAGa3OdqJoQCBrkOuF/8AAGa8d9qJoQBBrkOuFj8BWioQb"
	assert.Equal(t, expect, signedTx.Tx)
	t.Log(signedTx.Hash)
	assert.Equal(t, "43b1fe500eb878b62a712ce7d2591fcf3d3a10e20346543465d24ae7edea3f2f", signedTx.Hash)
}

func TestTransferJetton(t *testing.T) {
	seed, _ := hex.DecodeString("45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40")
	fromJettonAccount := "UQD7w9qG8Cq0PgX0hnp5nVpHPeDYL0QlhcLtjFMmna43sMxz"
	to := "UQC27fdnAFQcQDaXDrR89OKx-lW_Zyxuzcy5CjfPrS9A6vZf"
	amount := "1"
	seqno := uint32(0) // start 0
	pubKey := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
	signedTx, err := TransferJetton(seed, pubKey, fromJettonAccount, to, amount, 9, seqno, "", "1000000", "jetton test", 1, 1, false)
	assert.Nil(t, err)
	t.Log(signedTx)
	t.Log(signedTx.Hash)
	expect := "te6ccgECGAEABBIAA+OIAXkNjksffYh8JIQ426pdFemct+pPeD5Cxlw/0NF/aTcAEYDYxqIkc7gHj/wspWbNel478bGOM2usOFS7+tyy0aw7vqC26VK77CwqmwFrbp6HdZRq9VDEWDKH8o3cHj4kUeEFNTRi4AAAACAAAAAAAHABAgMBFP8A9KQT9LzyyAsEAFEAAAAAKamjFwwvPG2rtKBgDsyuh66qOSQgQvmldqqNygHhtBnPF9eiQAFoQgB94e1DeBVaHwL6Qz08zq0jnvBsF6ISwuF2ximTTtcb2CAX14QAAAAAAAAAAAAAAAAAAQUCASAGBwDGD4p+pQAAAAAAAAABEBgBbdvuzgCoOIBtLh1o+enFY/Srfs5Y3ZuZchRvn1pegdUALyGxyWPvsQ+EkIcbdUuivTOW/UnvB8hYy4f6Gi/tJuAGHoSAAAAAAGpldHRvbiB0ZXN0AgFICAkE+PKDCNcYINMf0x/THwL4I7vyZO1E0NMf0x/T//QE0VFDuvKhUVG68qIF+QFUEGT5EPKj+AAkpMjLH1JAyx9SMMv/UhD0AMntVPgPAdMHIcAAn2xRkyDXSpbTB9QC+wDoMOAhwAHjACHAAuMAAcADkTDjDQOkyMsfEssfy/8KCwwNAubQAdDTAyFxsJJfBOAi10nBIJJfBOAC0x8hghBwbHVnvSKCEGRzdHK9sJJfBeAD+kAwIPpEAcjKB8v/ydDtRNCBAUDXIfQEMFyBAQj0Cm+hMbOSXwfgBdM/yCWCEHBsdWe6kjgw4w0DghBkc3RyupJfBuMNDg8CASAQEQBu0gf6ANTUIvkABcjKBxXL/8nQd3SAGMjLBcsCIs8WUAX6AhTLaxLMzMlz+wDIQBSBAQj0UfKnAgBwgQEI1xj6ANM/yFQgR4EBCPRR8qeCEG5vdGVwdIAYyMsFywJQBs8WUAT6AhTLahLLH8s/yXP7AAIAbIEBCNcY+gDTPzBSJIEBCPRZ8qeCEGRzdHJwdIAYyMsFywJQBc8WUAP6AhPLassfEss/yXP7AAAK9ADJ7VQAeAH6APQEMPgnbyIwUAqhIb7y4FCCEHBsdWeDHrFwgBhQBMsFJs8WWPoCGfQAy2kXyx9SYMs/IMmAQPsABgCKUASBAQj0WTDtRNCBAUDXIMgBzxb0AMntVAFysI4jghBkc3Rygx6xcIAYUAXLBVADzxYj+gITy2rLH8s/yYBA+wCSXwPiAgEgEhMAWb0kK29qJoQICga5D6AhhHDUCAhHpJN9KZEM5pA+n/mDeBKAG3gQFImHFZ8xhAIBWBQVABG4yX7UTQ1wsfgAPbKd+1E0IEBQNch9AQwAsjKB8v/ydABgQEI9ApvoTGACASAWFwAZrc52omhAIGuQ64X/wAAZrx32omhAEGuQ64WPwA=="
	assert.Equal(t, expect, signedTx.Tx)
	assert.Equal(t, "abf0e9f9ab9a4e2b737580dea00091917aee63f52f1ec0d05b429d2dbd1f53a8", signedTx.Hash)
}

func TestCalTxHash(t *testing.T) {
	boc := "te6cckEBAgEArAAB34gB5RWsNJCcjosmHIR/8ivtI4/c6mpOliSwHD8/60xbNFoGXpoySnQ674GpDG2c4ercTebNjg+ARoFRVYZdNR2thDnpSC4LSQjLgesn7oH4/G9pNDb4L406HlVa+u9UEDVYCU1NGLsvVLYgAAAAqBwBAG5iAFhQVN5tJeFlyhBONxrLX6juythJMQtNZS8dw8/MvsxPnMS0AAAAAAAAAAAAAAAAAAAAAAAAtOdiNg=="
	hash, err := CalTxHash(boc)
	assert.Nil(t, err)
	//0q7RQCpKCrkSGLp2EW9Vk9b5VjlLMf3LT9/PE+Dnwxc=
	t.Log(hash)
	assert.Equal(t, "d2aed1402a4a0ab91218ba76116f5593d6f956394b31fdcb4fdfcf13e0e7c317", hash)
}
