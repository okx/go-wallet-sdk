package ton

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/okx/go-wallet-sdk/coins/ton/ton/wallet"
	"github.com/okx/go-wallet-sdk/coins/ton/tvm/cell"
	"testing"

	"github.com/stretchr/testify/assert"
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
	seed, _ := hex.DecodeString("11c01440658fce38bfc0a6029a1d5bfc8a34842d836bf00b521016d4aa6adcc1")
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
	assert.Equal(t, `{"address":"UQCjQa2xs4l01wvQnrWl46By9vMuy9cGycLoc7pgt8sy-1jO","body":"te6cckECAwEAAQ8AAZzxZne0gGdYWqMdaGMB/6/YnaQYDmnT/jo6KX9d1qUwcA/XgePUZaPUvA46tanb30NfVIiilog/m7RblGHEU7cKKamjF2ci1CIAAAC0AAMBAdNiAAioWoxZMTVqjEz8xEP8QSW4AyorIq+/8UCfgJNM0gMPoFz7tgAAAAAAAAAAAAAAAAAAD4p+pQAAAGo5eW95OYloCADvO5kConGyoByJOKUjz+JOcYR6rramIAAe1Ep3rA5wnBAsG4EDAgCdJZOFYYAcheIbhLa4Ep+9OcNZXQKKpcDAPO/gohpxCiCQKayPQ24LHKQ85mXdACjQa2xs4l01wvQnrWl46By9vMuy9cGycLoc7pgt8sy+0Ea2hE8=","init_data":"","init_code":"","ignore_chksig":true}`, tt.Str())
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
	seed, _ := hex.DecodeString("fc81e6f42150458f53d8c42551a8ab91978a55d0e22b1fd890b85139086b93f8")
	pubKey := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
	info, err := GetWalletInformation(seed, pubKey, wallet.V4R2)
	assert.NoError(t, err)
	initCode := "te6cckECFAEAAtQAART/APSkE/S88sgLAQIBIAIDAgFIBAUE+PKDCNcYINMf0x/THwL4I7vyZO1E0NMf0x/T//QE0VFDuvKhUVG68qIF+QFUEGT5EPKj+AAkpMjLH1JAyx9SMMv/UhD0AMntVPgPAdMHIcAAn2xRkyDXSpbTB9QC+wDoMOAhwAHjACHAAuMAAcADkTDjDQOkyMsfEssfy/8GBwgJAubQAdDTAyFxsJJfBOAi10nBIJJfBOAC0x8hghBwbHVnvSKCEGRzdHK9sJJfBeAD+kAwIPpEAcjKB8v/ydDtRNCBAUDXIfQEMFyBAQj0Cm+hMbOSXwfgBdM/yCWCEHBsdWe6kjgw4w0DghBkc3RyupJfBuMNCgsCASAMDQBu0gf6ANTUIvkABcjKBxXL/8nQd3SAGMjLBcsCIs8WUAX6AhTLaxLMzMlz+wDIQBSBAQj0UfKnAgBwgQEI1xj6ANM/yFQgR4EBCPRR8qeCEG5vdGVwdIAYyMsFywJQBs8WUAT6AhTLahLLH8s/yXP7AAIAbIEBCNcY+gDTPzBSJIEBCPRZ8qeCEGRzdHJwdIAYyMsFywJQBc8WUAP6AhPLassfEss/yXP7AAAK9ADJ7VQAeAH6APQEMPgnbyIwUAqhIb7y4FCCEHBsdWeDHrFwgBhQBMsFJs8WWPoCGfQAy2kXyx9SYMs/IMmAQPsABgCKUASBAQj0WTDtRNCBAUDXIMgBzxb0AMntVAFysI4jghBkc3Rygx6xcIAYUAXLBVADzxYj+gITy2rLH8s/yYBA+wCSXwPiAgEgDg8AWb0kK29qJoQICga5D6AhhHDUCAhHpJN9KZEM5pA+n/mDeBKAG3gQFImHFZ8xhAIBWBARABG4yX7UTQ1wsfgAPbKd+1E0IEBQNch9AQwAsjKB8v/ydABgQEI9ApvoTGACASASEwAZrc52omhAIGuQ64X/wAAZrx32omhAEGuQ64WPwGb/qfE="
	initData := `te6cckEBAQEAKwAAUQAAAAApqaMXSEu9odGrPom4O0I2hanUy+eii4uWu77lHQVTtH5NtC9As9YYUg==`
	walletStateInit := "te6cckECFgEAAwQAAgE0AQIBFP8A9KQT9LzyyAsDAFEAAAAAKamjF0hLvaHRqz6JuDtCNoWp1MvnoouLlru+5R0FU7R+TbQvQAIBIAQFAgFIBgcE+PKDCNcYINMf0x/THwL4I7vyZO1E0NMf0x/T//QE0VFDuvKhUVG68qIF+QFUEGT5EPKj+AAkpMjLH1JAyx9SMMv/UhD0AMntVPgPAdMHIcAAn2xRkyDXSpbTB9QC+wDoMOAhwAHjACHAAuMAAcADkTDjDQOkyMsfEssfy/8ICQoLAubQAdDTAyFxsJJfBOAi10nBIJJfBOAC0x8hghBwbHVnvSKCEGRzdHK9sJJfBeAD+kAwIPpEAcjKB8v/ydDtRNCBAUDXIfQEMFyBAQj0Cm+hMbOSXwfgBdM/yCWCEHBsdWe6kjgw4w0DghBkc3RyupJfBuMNDA0CASAODwBu0gf6ANTUIvkABcjKBxXL/8nQd3SAGMjLBcsCIs8WUAX6AhTLaxLMzMlz+wDIQBSBAQj0UfKnAgBwgQEI1xj6ANM/yFQgR4EBCPRR8qeCEG5vdGVwdIAYyMsFywJQBs8WUAT6AhTLahLLH8s/yXP7AAIAbIEBCNcY+gDTPzBSJIEBCPRZ8qeCEGRzdHJwdIAYyMsFywJQBc8WUAP6AhPLassfEss/yXP7AAAK9ADJ7VQAeAH6APQEMPgnbyIwUAqhIb7y4FCCEHBsdWeDHrFwgBhQBMsFJs8WWPoCGfQAy2kXyx9SYMs/IMmAQPsABgCKUASBAQj0WTDtRNCBAUDXIMgBzxb0AMntVAFysI4jghBkc3Rygx6xcIAYUAXLBVADzxYj+gITy2rLH8s/yYBA+wCSXwPiAgEgEBEAWb0kK29qJoQICga5D6AhhHDUCAhHpJN9KZEM5pA+n/mDeBKAG3gQFImHFZ8xhAIBWBITABG4yX7UTQ1wsfgAPbKd+1E0IEBQNch9AQwAsjKB8v/ydABgQEI9ApvoTGACASAUFQAZrc52omhAIGuQ64X/wAAZrx32omhAEGuQ64WPwL8I0H4="

	assert.Equal(t, walletStateInit, info.WalletStateInit)
	assert.Equal(t, initData, info.InitData)
	assert.Equal(t, initCode, info.InitCode)
}

func TestSignProof(t *testing.T) {
	seed, err := hex.DecodeString("fc81e6f42150458f53d8c42551a8ab91978a55d0e22b1fd890b85139086b93f8")
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
	expect := "tkrepSlC/7RMdGmtqNQ/OKRkxtdzvwF6GYBP6sgKWI/9mP0KcqyUwpFHAEGF6xeNOrwwoxIce8KJuEHLxuIgDw=="
	assert.Equal(t, r, expect)
	assert.NoError(t, VerifySignProof(addr, pub, sign, proof))
}

func TestVerify(t *testing.T) {
	seed, err := hex.DecodeString("fc81e6f42150458f53d8c42551a8ab91978a55d0e22b1fd890b85139086b93f8")
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
	expect := "tkrepSlC/7RMdGmtqNQ/OKRkxtdzvwF6GYBP6sgKWI/9mP0KcqyUwpFHAEGF6xeNOrwwoxIce8KJuEHLxuIgDw=="
	assert.Equal(t, r, expect)
	assert.NoError(t, VerifySignProofStr(addr, hex.EncodeToString(pub), r, proof))
}

func TestSignMultiTransferWithExtraFlags(t *testing.T) {
	r := &MultiRequest{
		Messages: []*Msg{
			{
				Address:    "EQAHxgkOidVwrDMsGmmf3m5y8752Z55rf-iu3ktN6xR8rx99",
				Amount:     "100000000",
				Payload:    "te6cckEBBAEAowABoXNpbnR///8RaO9fgAAAABKqAtOfR8Sau3JMYfCRQcYJyYHW5oJjdeDCNaWEmoKHmGv4Lh2UoSvCAix84WLSLK4EX6vLsv/SLu0FqZ8DrPpAoAECCg7DyG0DAgMAAACEQgAD4wSHROq4VhmWDTTP7zc5ed87M881v/RXbyWm9Yo+V6Hc1lAAAAAAAAAAAAAAAAAAAAAAAAB2NCB3cmFwIHY1SQmJyA==",
				ExtraFlags: "3",
			},
		},
		From:       "UQCd7tJX0EgL0y7kRy8HtMQCCiNinFPdS722ksDGRKK_pS8V",
		ValidUntil: 1719309177,
		Network:    "-239",
	}

	nonce := uint32(17)
	seed, _ := hex.DecodeString("87bfafe77b75a8bfbb95ace9d997798956b01c0a31a269daa83a73e2122e6fe7")
	assert.NoError(t, r.Check())
	s, err := SignMultiTransfer(seed, nil, nonce, r, false, wallet.V4R2)
	assert.NoError(t, err)

	expectedTx := `te6cckECBgEAAU8AAeGIATvdpK+gkBemXciOXg9piAQURsU4p7qXe20lgYyJRX9KBZL8tmyZZ17VHmu1fHQ2EZtGaMFs2PFG/dvA6hZWt5XreHGq+OxqNs3V2fcx8se9QHJ8DQK1Y0kdZsX/gsYE2GlNTRi7M9SbyAAAAIgAHAEBamIAA+MEh0TquFYZlg00z+83OXnfOzPPNb/0V28lpvWKPlegL68IAEDAAAAAAAAAAAAAAAABAgGhc2ludH///xFo71+AAAAAEqoC059HxJq7ckxh8JFBxgnJgdbmgmN14MI1pYSagoeYa/guHZShK8ICLHzhYtIsrgRfq8uy/9Iu7QWpnwOs+kCgAwIKDsPIbQMEBQAAAIRCAAPjBIdE6rhWGZYNNM/vNzl53zszzzW/9FdvJab1ij5XodzWUAAAAAAAAAAAAAAAAAAAAAAAAHY0IHdyYXAgdjVn1oM6`
	expectedTxHash := `2bcc71ab22a9aa3580cf46636f6c9c0d08a12710696ec7dabf7821861cdf9eaa`
	expectedNormHash := `6d882d17d483a0330ec2a4c8ca5c84d1d6ee4a76ca900b7823abf43dff441c62`
	assert.Equal(t, expectedTx, s.Tx)
	assert.Equal(t, expectedTxHash, s.Hash)
	assert.Equal(t, expectedNormHash, s.NormHash)

}
