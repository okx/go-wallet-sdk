package ton

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/okx/go-wallet-sdk/coins/ton/ton/wallet"
	"github.com/okx/go-wallet-sdk/coins/ton/tvm/cell"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromBOC(t *testing.T) {
	b, err := base64.StdEncoding.DecodeString("te6ccsEBAwEATgAAHEEBMQAAAAEAAAAAAAAAAAAAAAAAAAAAQC+vCAgBAUOAGJYSqICRE9K5vqy5X2NGnMwhoOf82pCxB7OYZJNqSG1QAgAWbXlfbmZ0Lmpzb25Yc7tj")
	assert.NoError(t, err)
	c, err := cell.FromBOC(b)
	assert.NoError(t, err)
	assert.Equal(t, "ecfcb15660a415e7cbd2012f64ef9b3b4f30e749e292e388d7f28dd01919f643", hex.EncodeToString(c.Hash()))
}

func tryParseBase64(body string) ([]byte, error) {
	if b, errStd := base64.StdEncoding.DecodeString(body); errStd == nil && len(b) > 0 {
		return b, nil
	}
	if b, errRawStd := base64.RawStdEncoding.DecodeString(body); errRawStd == nil && len(b) > 0 {
		return b, nil
	}
	if b, errUrl := base64.URLEncoding.DecodeString(body); errUrl == nil && len(b) > 0 {
		return b, nil
	}
	if b, errRawUrl := base64.RawURLEncoding.DecodeString(body); errRawUrl == nil && len(b) > 0 {
		return b, nil
	}
	return nil, errors.New("invalid base64 string")
}

func TestParsebase64(t *testing.T) {
	d, err := base64.StdEncoding.DecodeString("te6cckEBAgEAigABaw+KfqUAAABqOXlveTmJaAgA7zuZAqJxsqAciTilI8/iTnGEeq62piAAHtRKd6wOcJwQLBuBAwEAnSWThWGAHIXiG4S2uBKfvTnDWV0CiqXAwDzv4KIacQogkCmsj0NuCxykPOZl3QAo0GtsbOJdNcL0J61peOgcvbzLsvXBsnC6HO6YLfLMvtDoIHzZ")
	assert.NoError(t, err)
	d1, err := tryParseBase64(base64.URLEncoding.EncodeToString(d))
	assert.NoError(t, err)
	assert.Equal(t, d1, d)
	d2, err := tryParseBase64(base64.RawStdEncoding.EncodeToString(d))
	assert.NoError(t, err)
	assert.Equal(t, d2, d)
	d3, err := tryParseBase64(base64.RawURLEncoding.EncodeToString(d))
	assert.NoError(t, err)
	assert.Equal(t, d3, d)
	d4, err := tryParseBase64(base64.URLEncoding.EncodeToString(d))
	assert.NoError(t, err)
	assert.Equal(t, d4, d)
}

func TestSignMultiTransfer(t *testing.T) {
	var r MultiRequest
	code := `{
        "messages":[{
	"address": "EQARULUYsmJq1RiZ-YiH-IJLcAZUVkVff-KBPwEmmaQGH6aC",
	"amount": "195000000",
	"payload":"te6cckEBAgEAigABaw+KfqUAAABqOXlveTmJaAgA7zuZAqJxsqAciTilI8/iTnGEeq62piAAHtRKd6wOcJwQLBuBAwEAnSWThWGAHIXiG4S2uBKfvTnDWV0CiqXAwDzv4KIacQogkCmsj0NuCxykPOZl3QAo0GtsbOJdNcL0J61peOgcvbzLsvXBsnC6HO6YLfLMvtDoIHzZ"
	}],
        "from": "0:a341adb1b38974d70bd09eb5a5e3a072f6f32ecbd706c9c2e873ba60b7cb32fb",
 "valid_until": 1730335778,
        "network": "-239"
}`
	nonce := uint32(180)
	err := json.Unmarshal([]byte(code), &r)
	assert.NoError(t, err)
	seed, err := hex.DecodeString("45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40")
	assert.NoError(t, err)
	assert.NoError(t, r.Check())
	s, err := SignMultiTransfer(seed, nil, nonce, &r, true, wallet.V4R2)
	assert.NoError(t, err)
	tt := &testSignedTx{
		Address:      s.Address,
		Body:         s.Tx,
		InitData:     s.Data,
		InitCode:     s.Code,
		IgnoreChksig: true,
	}
	assert.Equal(t, `{"address":"UQC8hsclj77EPhJCHG3VLor0zlv1J7wfIWMuH-hov7SbgIIM","body":"te6cckECAwEAAQ8AAZwEKGXfLJaboKoat0eAmT4V1Sn/2Mr5J/dB3LVr2zDN41x8bbs4LNRTW3jOtysMSpJ7nOF7suuuPCkTN9JlhyUAKamjF2ci1CIAAAC0AAMBAdNiAAioWoxZMTVqjEz8xEP8QSW4AyorIq+/8UCfgJNM0gMPoFz7tgAAAAAAAAAAAAAAAAAAD4p+pQAAAGo5eW95OYloCADvO5kConGyoByJOKUjz+JOcYR6rramIAAe1Ep3rA5wnBAsG4EDAgCdJZOFYYAcheIbhLa4Ep+9OcNZXQKKpcDAPO/gohpxCiCQKayPQ24LHKQ85mXdACjQa2xs4l01wvQnrWl46By9vMuy9cGycLoc7pgt8sy+0IKdeXA=","init_data":"","init_code":"","ignore_chksig":true}`, tt.Str())
}

type testSignedTx struct {
	Address      string `json:"address"`
	Body         string `json:"body"`
	InitData     string `json:"init_data"`
	InitCode     string `json:"init_code"`
	IgnoreChksig bool   `json:"ignore_chksig"`
}

func (t *testSignedTx) Str() string {
	if t == nil {
		return ""
	}
	j, err := json.Marshal(t)
	if err != nil {
		return ""
	}
	return string(j)
}

func TestGetAccontInfo(t *testing.T) {
	seed, _ := hex.DecodeString("45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40")
	pubKey := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
	info, err := GetWalletInformation(seed, pubKey, wallet.V4R2)
	assert.NoError(t, err)
	initCode := "te6cckECFAEAAtQAART/APSkE/S88sgLAQIBIAIDAgFIBAUE+PKDCNcYINMf0x/THwL4I7vyZO1E0NMf0x/T//QE0VFDuvKhUVG68qIF+QFUEGT5EPKj+AAkpMjLH1JAyx9SMMv/UhD0AMntVPgPAdMHIcAAn2xRkyDXSpbTB9QC+wDoMOAhwAHjACHAAuMAAcADkTDjDQOkyMsfEssfy/8GBwgJAubQAdDTAyFxsJJfBOAi10nBIJJfBOAC0x8hghBwbHVnvSKCEGRzdHK9sJJfBeAD+kAwIPpEAcjKB8v/ydDtRNCBAUDXIfQEMFyBAQj0Cm+hMbOSXwfgBdM/yCWCEHBsdWe6kjgw4w0DghBkc3RyupJfBuMNCgsCASAMDQBu0gf6ANTUIvkABcjKBxXL/8nQd3SAGMjLBcsCIs8WUAX6AhTLaxLMzMlz+wDIQBSBAQj0UfKnAgBwgQEI1xj6ANM/yFQgR4EBCPRR8qeCEG5vdGVwdIAYyMsFywJQBs8WUAT6AhTLahLLH8s/yXP7AAIAbIEBCNcY+gDTPzBSJIEBCPRZ8qeCEGRzdHJwdIAYyMsFywJQBc8WUAP6AhPLassfEss/yXP7AAAK9ADJ7VQAeAH6APQEMPgnbyIwUAqhIb7y4FCCEHBsdWeDHrFwgBhQBMsFJs8WWPoCGfQAy2kXyx9SYMs/IMmAQPsABgCKUASBAQj0WTDtRNCBAUDXIMgBzxb0AMntVAFysI4jghBkc3Rygx6xcIAYUAXLBVADzxYj+gITy2rLH8s/yYBA+wCSXwPiAgEgDg8AWb0kK29qJoQICga5D6AhhHDUCAhHpJN9KZEM5pA+n/mDeBKAG3gQFImHFZ8xhAIBWBARABG4yX7UTQ1wsfgAPbKd+1E0IEBQNch9AQwAsjKB8v/ydABgQEI9ApvoTGACASASEwAZrc52omhAIGuQ64X/wAAZrx32omhAEGuQ64WPwGb/qfE="
	initData := "te6cckEBAQEAKwAAUQAAAAApqaMXDC88bau0oGAOzK6Hrqo5JCBC+aV2qo3KAeG0Gc8X16JA0rBAuw=="
	walletStateInit := "te6cckECFgEAAwQAAgE0AQIBFP8A9KQT9LzyyAsDAFEAAAAAKamjFwwvPG2rtKBgDsyuh66qOSQgQvmldqqNygHhtBnPF9eiQAIBIAQFAgFIBgcE+PKDCNcYINMf0x/THwL4I7vyZO1E0NMf0x/T//QE0VFDuvKhUVG68qIF+QFUEGT5EPKj+AAkpMjLH1JAyx9SMMv/UhD0AMntVPgPAdMHIcAAn2xRkyDXSpbTB9QC+wDoMOAhwAHjACHAAuMAAcADkTDjDQOkyMsfEssfy/8ICQoLAubQAdDTAyFxsJJfBOAi10nBIJJfBOAC0x8hghBwbHVnvSKCEGRzdHK9sJJfBeAD+kAwIPpEAcjKB8v/ydDtRNCBAUDXIfQEMFyBAQj0Cm+hMbOSXwfgBdM/yCWCEHBsdWe6kjgw4w0DghBkc3RyupJfBuMNDA0CASAODwBu0gf6ANTUIvkABcjKBxXL/8nQd3SAGMjLBcsCIs8WUAX6AhTLaxLMzMlz+wDIQBSBAQj0UfKnAgBwgQEI1xj6ANM/yFQgR4EBCPRR8qeCEG5vdGVwdIAYyMsFywJQBs8WUAT6AhTLahLLH8s/yXP7AAIAbIEBCNcY+gDTPzBSJIEBCPRZ8qeCEGRzdHJwdIAYyMsFywJQBc8WUAP6AhPLassfEss/yXP7AAAK9ADJ7VQAeAH6APQEMPgnbyIwUAqhIb7y4FCCEHBsdWeDHrFwgBhQBMsFJs8WWPoCGfQAy2kXyx9SYMs/IMmAQPsABgCKUASBAQj0WTDtRNCBAUDXIMgBzxb0AMntVAFysI4jghBkc3Rygx6xcIAYUAXLBVADzxYj+gITy2rLH8s/yYBA+wCSXwPiAgEgEBEAWb0kK29qJoQICga5D6AhhHDUCAhHpJN9KZEM5pA+n/mDeBKAG3gQFImHFZ8xhAIBWBITABG4yX7UTQ1wsfgAPbKd+1E0IEBQNch9AQwAsjKB8v/ydABgQEI9ApvoTGACASAUFQAZrc52omhAIGuQ64X/wAAZrx32omhAEGuQ64WPwF9YoYQ="

	assert.Equal(t, walletStateInit, info.WalletStateInit)
	assert.Equal(t, initData, info.InitData)
	assert.Equal(t, initCode, info.InitCode)
}

func TestSignProof(t *testing.T) {
	seed, err := hex.DecodeString("45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40")
	assert.NoError(t, err)
	addr := "EQA3_JIJKDC0qauDUEQe2KjQj1iLwQRtrEREzmfDxbCKw9Kr"
	proof := &ProofData{
		Timestamp: 1719309177,
		Domain:    "ton.org.com",
		Payload:   "123",
	}
	r, err := SignProof(addr, seed, proof)
	assert.NoError(t, err)
	pub := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
	sign, err := base64.StdEncoding.DecodeString(r)
	assert.NoError(t, err)
	expect := "V1ImmDgpt4DtZYYeGeZz38w7J+dXtYbBf/Hl7DLcWLEad23TOexKCSTO1f+N7i3UDreGVfycaVNbOspJnr9aDw=="
	assert.Equal(t, r, expect)
	assert.NoError(t, VerifySignProof(addr, pub, sign, proof))
}

func TestVerify(t *testing.T) {
	seed, err := hex.DecodeString("45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40")
	assert.NoError(t, err)
	addr := "EQA3_JIJKDC0qauDUEQe2KjQj1iLwQRtrEREzmfDxbCKw9Kr"
	proof := &ProofData{
		Timestamp: 1719309177,
		Domain:    "ton.org.com",
		Payload:   "123",
	}
	r, err := SignProof(addr, seed, proof)
	assert.NoError(t, err)
	pub := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
	expect := "V1ImmDgpt4DtZYYeGeZz38w7J+dXtYbBf/Hl7DLcWLEad23TOexKCSTO1f+N7i3UDreGVfycaVNbOspJnr9aDw=="
	assert.Equal(t, r, expect)
	assert.NoError(t, VerifySignProofStr(addr, hex.EncodeToString(pub), r, proof))
}
