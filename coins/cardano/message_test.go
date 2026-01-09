package cardano

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignMessage(t *testing.T) {
	prvKey := "40e27cee4a8b4c4ec19e6e6060d78c1e0ecf18141df615b4d6294800b905425e511a28ae1a7595b8e9465b76f6cf4981238ac8be8584e767234fdca7f15e14d7c8fc0f836e5e7599d95cea1c9ac639c01b46ecf213ef7fe0af1804d7b905425e63d2f4e0b46a2c978c961ad55d5396a0f221e49cb004a26fbe422859fd77c5d1"
	message := "1234"

	address, err := NewAddressFromPrvKey(prvKey)
	assert.NoError(t, err)

	res, err := SignMessage(prvKey, address, message)
	assert.NoError(t, err)
	expectedSignature := "845846a201276761646472657373583901e6369b50e580eeaeccf18e47e21f7995f78f362e3787d6741469d9983726764440ed33139ac55cd1b0359cad2e8fcaf8eed41f8e1cde5dc7a166686173686564f442123458409b2ae77d56f197dd1e2e6f3fbf6084d842d877b2a1a069afb20756e7b537fec79bfca6e2e2246b3727b8edf4934b12839db6d5c1f63966b86b5de05c20f94604"
	assert.Equal(t, expectedSignature, res.Signature)
	expectedKey := "a4010103272006215820a57677142d5785216774a659f7ce05556792d950702d8b2306bbb11381f22c14"
	assert.Equal(t, expectedKey, res.Key)

}

func TestVerifyMessageNoAddr(t *testing.T) {
	signature := "845846a201276761646472657373583901e6369b50e580eeaeccf18e47e21f7995f78f362e3787d6741469d9983726764440ed33139ac55cd1b0359cad2e8fcaf8eed41f8e1cde5dc7a166686173686564f442123458409b2ae77d56f197dd1e2e6f3fbf6084d842d877b2a1a069afb20756e7b537fec79bfca6e2e2246b3727b8edf4934b12839db6d5c1f63966b86b5de05c20f94604"
	key := "a4010103272006215820a57677142d5785216774a659f7ce05556792d950702d8b2306bbb11381f22c14"
	pubKey := "a57677142d5785216774a659f7ce05556792d950702d8b2306bbb11381f22c14ef3a9790fb2034d5ed38afa3171c29318167c5bdb0d45279edcce40945437998"
	//address := "addr1q8nrdx6sukqwatkv7x8y0csl0x2l0rek9cmc04n5z35anxphyemygs8dxvfe432u6xcrt89d968u478w6s0cu8x7thrsng2ztx"
	message := "1234"
	verified, err := VerifyMessageNoAddr(signature, key, pubKey, message)
	assert.NoError(t, err)
	assert.True(t, verified)

}
func TestVerifyMessage(t *testing.T) {
	signature := "845846a201276761646472657373583901e6369b50e580eeaeccf18e47e21f7995f78f362e3787d6741469d9983726764440ed33139ac55cd1b0359cad2e8fcaf8eed41f8e1cde5dc7a166686173686564f442123458409b2ae77d56f197dd1e2e6f3fbf6084d842d877b2a1a069afb20756e7b537fec79bfca6e2e2246b3727b8edf4934b12839db6d5c1f63966b86b5de05c20f94604"
	key := "a4010103272006215820a57677142d5785216774a659f7ce05556792d950702d8b2306bbb11381f22c14"
	pubKey := "a57677142d5785216774a659f7ce05556792d950702d8b2306bbb11381f22c14ef3a9790fb2034d5ed38afa3171c29318167c5bdb0d45279edcce40945437998"
	address := "addr1q8nrdx6sukqwatkv7x8y0csl0x2l0rek9cmc04n5z35anxphyemygs8dxvfe432u6xcrt89d968u478w6s0cu8x7thrsng2ztx"
	message := "1234"
	verified, err := VerifyMessage(signature, key, pubKey, address, message)
	assert.NoError(t, err)
	assert.True(t, verified)

}
