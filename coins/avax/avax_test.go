package avax

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestSign(t *testing.T) {
	pk, _ := hex.DecodeString("d27a851e2ffe50d81d639a5bc17ccb488b1441307fea7636e264b9da0ce577a1")
	pk2, _ := hex.DecodeString("201ae7b3cf9dce9fec7fa6f9b48fbb4390279ed039331fce4d7509e52c15d1ec")
	_, b := btcec.PrivKeyFromBytes(pk)
	_, b2 := btcec.PrivKeyFromBytes(pk2)
	addr, err := NewAddress(CHAINID_X, HRP_FUJI, b)
	if err != nil {
		t.Fatal(err)
	}
	addr2, err := NewAddress(CHAINID_X, HRP_FUJI, b2)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "X-fuji185n87yeprljt4376qqsdtltqs53rtqj8mz8u20", addr)
	assert.Equal(t, "X-fuji1qetw42gr9l8t626rrlzxglf8km532m9c9mvdz9", addr2)

	_, chainId, hrp, err := ParseAddress(addr)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "X", chainId)
	assert.Equal(t, "fuji", hrp)

	var inputs []TransferInput
	var outputs []TransferOutPut

	c := math.Pow10(9)
	inputs = append(inputs, TransferInput{TxId: "sJNJVJQzmjyrAoPfshkDhKNf55jNSNW7NXK8SygGdNrst2waA", Index: 0, Amount: uint64(2 * c), AssetId: ASSET_AVAX_FUJI, PrivateKey: "bf77591baae00a9b2826ae63d6668fe5c1cd934fcaf5c99946af9d55457533ce"})
	outputs = append(outputs, TransferOutPut{Address: "X-fuji1xqq48uejmydn95dwmvk4ge7rs9mj60nlx94dst", AssetId: ASSET_AVAX_FUJI, Value: uint64(c)})
	outputs = append(outputs, TransferOutPut{Address: "X-fuji1asep0ygju0g2trqq2pvpez736gngthh29lkazf", AssetId: ASSET_AVAX_FUJI, Value: uint64(0.99 * c)})

	ret, err := NewTransferTransaction(NETWORK_FUJI, BLOCKCHAIN_FUJI, &inputs, &outputs)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "111111111489TZ9qvxFYbETdE4vQ9b9ksjoUzgnoH6hoqcdW9G7FUn2StUWnXXMsVnR1L88SJXygzJrsge43HLKVHNjBAaDfqDmpZg8Up6uVnzhtx2cXsbGonMRCFPDE64LThuXWKXYkh88zdVr7gDLyW4489BTwa3rpj81Hk1VuZWWdK9E82XHZgzyQJoF2jUabXP55PRVf3pq4VqEVJeJsBAX5fgjUzsAwVh3cBHt4RD1vK6zxGYZRLHgg6Acdsr6a3GTHfuycCkUUX4HFAdym9K6SoRjndHKufGWkukcDdjGWiKSV5wX6wEMvmKioZ9jrpPgULNpHxcBgUtR1pcna2wXRiaJE5HQgYs5jHLi4X7XEFvc4DJEro5yzhfb4RZkA1wKbt3Ar6zL28g5U4NwxZcejT9qiioaWjbgFtDT8UdeEAgT1ehmvyqNvSqSwBTeLGgTR1qG3C3QPeNeJPFMmBqrUBoxTmej5dx5kv1jtgw8i6WJyrRnBQ4TvE9YAJBgUL8nd5N2L6", ret)
}
