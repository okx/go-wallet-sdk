package ton

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/ton/ton/wallet"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressStrings(t *testing.T) {
	a := "UQBjHRX_2tP47XFqKcoec4n9gbj9p4BH69Z4Plh-0qYTP9UI"
	all, err := AddressStrings(a)
	assert.NoError(t, err)
	fmt.Println(all)
	assert.Equal(t, 3, len(all))
	s, err := json.Marshal(all)
	fmt.Println(string(s), err)
	assert.Equal(t, "0:631d15ffdad3f8ed716a29ca1e7389fd81b8fda78047ebd6783e587ed2a6133f", all[0].Addr)
	assert.Equal(t, "EQBjHRX_2tP47XFqKcoec4n9gbj9p4BH69Z4Plh-0qYTP4jN", all[1].Addr)
	assert.Equal(t, "UQBjHRX_2tP47XFqKcoec4n9gbj9p4BH69Z4Plh-0qYTP9UI", all[2].Addr)
}

func TestGetStateInit(t *testing.T) {
	seedHex := "45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40"
	seed, _ := hex.DecodeString(seedHex)
	pubKey := ed25519.NewKeyFromSeed(seed).Public().(ed25519.PublicKey)
	stateInit, err := wallet.GetStateInit(pubKey, wallet.V4R2, wallet.DefaultSubwallet)
	assert.NoError(t, err)
	fmt.Println("Data", base64.StdEncoding.EncodeToString(stateInit.Data.ToBOC()))
	fmt.Println("Code", base64.StdEncoding.EncodeToString(stateInit.Code.ToBOC()))
	data := "te6cckEBAQEAKwAAUQAAAAApqaMXDC88bau0oGAOzK6Hrqo5JCBC+aV2qo3KAeG0Gc8X16JA0rBAuw=="
	code := "te6cckECFAEAAtQAART/APSkE/S88sgLAQIBIAIDAgFIBAUE+PKDCNcYINMf0x/THwL4I7vyZO1E0NMf0x/T//QE0VFDuvKhUVG68qIF+QFUEGT5EPKj+AAkpMjLH1JAyx9SMMv/UhD0AMntVPgPAdMHIcAAn2xRkyDXSpbTB9QC+wDoMOAhwAHjACHAAuMAAcADkTDjDQOkyMsfEssfy/8GBwgJAubQAdDTAyFxsJJfBOAi10nBIJJfBOAC0x8hghBwbHVnvSKCEGRzdHK9sJJfBeAD+kAwIPpEAcjKB8v/ydDtRNCBAUDXIfQEMFyBAQj0Cm+hMbOSXwfgBdM/yCWCEHBsdWe6kjgw4w0DghBkc3RyupJfBuMNCgsCASAMDQBu0gf6ANTUIvkABcjKBxXL/8nQd3SAGMjLBcsCIs8WUAX6AhTLaxLMzMlz+wDIQBSBAQj0UfKnAgBwgQEI1xj6ANM/yFQgR4EBCPRR8qeCEG5vdGVwdIAYyMsFywJQBs8WUAT6AhTLahLLH8s/yXP7AAIAbIEBCNcY+gDTPzBSJIEBCPRZ8qeCEGRzdHJwdIAYyMsFywJQBc8WUAP6AhPLassfEss/yXP7AAAK9ADJ7VQAeAH6APQEMPgnbyIwUAqhIb7y4FCCEHBsdWeDHrFwgBhQBMsFJs8WWPoCGfQAy2kXyx9SYMs/IMmAQPsABgCKUASBAQj0WTDtRNCBAUDXIMgBzxb0AMntVAFysI4jghBkc3Rygx6xcIAYUAXLBVADzxYj+gITy2rLH8s/yYBA+wCSXwPiAgEgDg8AWb0kK29qJoQICga5D6AhhHDUCAhHpJN9KZEM5pA+n/mDeBKAG3gQFImHFZ8xhAIBWBARABG4yX7UTQ1wsfgAPbKd+1E0IEBQNch9AQwAsjKB8v/ydABgQEI9ApvoTGACASASEwAZrc52omhAIGuQ64X/wAAZrx32omhAEGuQ64WPwGb/qfE="
	assert.Equal(t, data, base64.StdEncoding.EncodeToString(stateInit.Data.ToBOC()))
	assert.Equal(t, code, base64.StdEncoding.EncodeToString(stateInit.Code.ToBOC()))
}

func TestFromSeedV4R2(t *testing.T) {
	prv, err := FromSeedV4R2("good word gossip learn giggle nose bar silk crawl fold hire exercise bulk game rebel hello indicate lunar indoor scrap flip silent orbit twice", "")
	assert.NoError(t, err)
	fmt.Println(hex.EncodeToString(prv[0:32]))
	address, err := NewAddress(prv.Seed())
	assert.Nil(t, err)
	fmt.Println(address)
	//assert.Equal(t, "UQBjHRX_2tP47XFqKcoec4n9gbj9p4BH69Z4Plh-0qYTP9UI", address)
}

func TestNewAddress(t *testing.T) {
	seedHex := "45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40"
	seed, _ := hex.DecodeString(seedHex)
	address, err := NewAddress(seed)
	fmt.Println(address)
	assert.Nil(t, err)
	assert.Equal(t, "UQC8hsclj77EPhJCHG3VLor0zlv1J7wfIWMuH-hov7SbgIIM", address)
}

func TestValidateAddress(t *testing.T) {
	address := "UQD7w9qG8Cq0PgX0hnp5nVpHPeDYL0QlhcLtjFMmna43sMxz"
	isValid := ValidateAddress(address)
	assert.True(t, isValid)

	address2 := "UQD7w9qG8Cq0PgX0hnp5nVpHPeDYL0QlhcLtjFMmna43sMxy"
	isValid2 := ValidateAddress(address2)
	assert.False(t, isValid2)
}

func TestVenomNewAddress(t *testing.T) {
	seedHex := "45d3bd794c5bc6ed91ae41c93c0baed679935703dfac72c48d27f8321b8d3a40"
	seed, _ := hex.DecodeString(seedHex)
	address, err := VenomNewAddress(seed)
	assert.Nil(t, err)
	assert.Equal(t, "0:a47d62625b54a73d662f76372326e223fa2c4041df66d44552c3c9418ba09479", address)
}
