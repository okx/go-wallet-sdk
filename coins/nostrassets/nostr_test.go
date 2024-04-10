package nostrassets

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVerify(t *testing.T) {
	hash, err := hex.DecodeString("385eb020a83cb7e547659922b6c092a55e88c5127d9448370d1e55221aaeb5dd")
	if err != nil {
		t.Fatal(err)
	}
	pub, err := hex.DecodeString("14ccbe1d4a55fe23628576a7f04637f647fd6b86d362f983f4ebd7b95d47796f")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("pub", pub)
	sig, err := hex.DecodeString("622b442f220c83debe0a1a89fb142dc121bbd5b413befe1e62566351705f4f46e6de9b1726cb8372b5714af87aa8f4c13ec7b9fcd7d7a558f36acdf466a08a4d")
	if err != nil {
		t.Fatal(err)
	}
	s, err := schnorr.ParseSignature(sig)
	if err != nil {
		t.Fatal(err)
	}
	p, err := schnorr.ParsePubKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	if !s.Verify(hash, p) {
		t.Fatal("invalid sign")
	}
}

func TestSignEvent(t *testing.T) {
	prvBech := "nsec1hvwfx5ytjck8a7c2xsyys5ut930hhfkyfe2l2guf4gfj5t7n2gdqxvh70y"
	event := &Event{Kind: 1, CreatedAt: 1000, Content: "hello", Tags: [][]string{}}
	rr, err := SignEvent(prvBech, event)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.Marshal(rr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(b))
	assert.Equal(t, rr.Pubkey, "14ccbe1d4a55fe23628576a7f04637f647fd6b86d362f983f4ebd7b95d47796f")
	assert.Equal(t, rr.Id, "385eb020a83cb7e547659922b6c092a55e88c5127d9448370d1e55221aaeb5dd")
	assert.Equal(t, VerifyEvent(rr), true)
}

func TestGetEventHash(t *testing.T) {
	event := &Event{Kind: 1, CreatedAt: 1000, Pubkey: "14ccbe1d4a55fe23628576a7f04637f647fd6b86d362f983f4ebd7b95d47796f", Content: "hello", Tags: [][]string{}}
	rr, err := GetEventHash(event)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, rr, "385eb020a83cb7e547659922b6c092a55e88c5127d9448370d1e55221aaeb5dd")
}

func TestEncrypt(t *testing.T) {
	r, err := Encrypt("nsec1gfvzgfpwxquwqfhhe04klc5fedhle7k3l224cvvvz942ruhn907qsmz8sf", "0x8a0523d045d09c30765029af9307d570cb0d969e4b9400c08887c23250626eea", "hello", "RK8MdoXLJLe0R4DbYAd2oQ==")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "gReE/VquUuD+O0nbodbolQ==?iv=RK8MdoXLJLe0R4DbYAd2oQ==", r)
	rr, err := Decrypt("nsec1gfvzgfpwxquwqfhhe04klc5fedhle7k3l224cvvvz942ruhn907qsmz8sf", "0x8a0523d045d09c30765029af9307d570cb0d969e4b9400c08887c23250626eea", r)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "hello", rr)
}

func TestNpubEncode2(t *testing.T) {
	s, b, err := bech32.DecodeToBase256("nsec1hvwfx5ytjck8a7c2xsyys5ut930hhfkyfe2l2guf4gfj5t7n2gdqxvh70y")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s, hex.EncodeToString(b))
}

func TestNpubEncode(t *testing.T) {
	prvBech := "nsec1hvwfx5ytjck8a7c2xsyys5ut930hhfkyfe2l2guf4gfj5t7n2gdqxvh70y"
	pub, err := GetPublicKey(prvBech)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "14ccbe1d4a55fe23628576a7f04637f647fd6b86d362f983f4ebd7b95d47796f", pub)
	res, err := NpubEncode(pub)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "npub1znxtu8222hlzxc59w6nlq33h7erl66ux6d30nql5a0tmjh2809hstw0d22", res)
	addr, err := AddressFromPrvKey(prvBech)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "npub1znxtu8222hlzxc59w6nlq33h7erl66ux6d30nql5a0tmjh2809hstw0d22", addr)
}
