package tezos

import (
	"encoding/hex"
	"encoding/json"
	"gitlab.okg.com/wallet-sign-core/go-parent-sdk/coin/tezos/types"
	"testing"
)

const (
	p1 = "edskRqZzHUtUb2sLMQpM6fpZX1qNfbDZ8jSG7QHNYVy5mPNE6J15wmdMSLjomdoPRQWYPDufdV8xxrG73FiyL2oNGw9kxshAnw"
	p2 = "edskS2nYyy3vDYvVsQjg6WxDhVTGNLRQ6NYxpb11wQJrc4sqDRJgrhx1k49ZqBG99Wyci4GXPjSRv4YxJX4VozpKAL8VKgGztU"
	p3 = "edskS76q1W5Gf3KqekuaTh5t95hAsm4DML7S4BTizAucn7wjwioMNVKjy5HvWcsMxys6jRuWodnmcPVmuhcefbGa3rQq5JXZ3p"
	p4 = "edskRfQwFULfLgdNoWxsdbvSvdYduNQuF8xa1hYCzFZhYWhvVCjiYnAXo6A8Yt1J2G1VMM5eUKRjvXo3K33qRkc8mTLh7MZyJt"
	n1 = "tz1e9qxfTpvqub3rj1nryeRVshcA1tca8Tsq"
	n2 = "tz1fv8Wj8jyxb4wp9M6dfAk7irLhh8T62Uze"
	n3 = "tz1YbSetwCa2EzqgFQ6Amii4VHd87g92z8mK"
	n4 = "tz1VXXwzPS1dREuz4adwZHsJSbDv7urjvMfX"
)

func TestNewTransaction(t *testing.T) {
	var amount int64 = 6000000
	var fee int64 = 10000      // 0.010000 XTZ
	var counter int64 = 339709 // 与 operation 有关
	opt := NewCallOptions("BL74GqeaJ8tdFuBR2RhsGXET7MNonprQ49BBreZHE9yn9x85hJP", counter, false)
	privateKey, err := types.ParsePrivateKey(p1)
	if err != nil {
		t.Error(err)
		return
	}
	tx, err := NewJakartanetTransaction(n1, n3, amount, opt)
	if err != nil {
		t.Error(err)
		return
	}
	err = BuildTransaction(tx, fee, privateKey.Public(), opt)
	if err != nil {
		t.Error(err)
		return
	}
	rawTx, err := SignTransaction(tx, p1, opt)
	if err != nil {
		t.Error(err)
		return
	}

	rawTxByte, _ := json.Marshal(tx)
	t.Log(string(rawTxByte))
	t.Log(rawTx)
	t.Log(hex.EncodeToString(rawTx))
}

func TestNewJakartanetDelegationTransaction(t *testing.T) {
	var to string = "tz1foXHgRzdYdaLgX6XhpZGxbBv42LZ6ubvE" // 找一个面包师
	var fee int64 = 10000                                  // 0.010000 XTZ
	var counter int64 = 331345                             // 与 operation 有关
	opt := NewCallOptions("BL7kQbhcCsMYB954n94XcSZmS1oTYcg8J7ut2wj6iZpL3fRBdM3", counter, false)
	privateKey, err := types.ParsePrivateKey(p2)
	if err != nil {
		t.Error(err)
		return
	}
	tx, err := NewJakartanetDelegationTransaction(n2, to, opt)
	if err != nil {
		t.Error(err)
		return
	}
	err = BuildTransaction(tx, fee, privateKey.Public(), opt)
	if err != nil {
		t.Error(err)
		return
	}
	rawTx, err := SignTransaction(tx, p2, opt)
	if err != nil {
		t.Error(err)
		return
	}

	rawTxByte, _ := json.Marshal(tx)
	t.Log(string(rawTxByte))
	t.Log(rawTx)
	t.Log(hex.EncodeToString(rawTx))
}

func TestNewJakartanetUnDelegationTransaction(t *testing.T) {
	var fee int64 = 10000      // 0.010000 XTZ
	var counter int64 = 339708 // 与 operation 有关

	opt := NewCallOptions("BL7kQbhcCsMYB954n94XcSZmS1oTYcg8J7ut2wj6iZpL3fRBdM3", counter, false)
	privateKey, err := types.ParsePrivateKey(p2)
	if err != nil {
		t.Error(err)
		return
	}
	tx, err := NewJakartanetUnDelegationTransaction(n2, opt)
	if err != nil {
		t.Error(err)
		return
	}
	err = BuildTransaction(tx, fee, privateKey.Public(), opt)
	if err != nil {
		t.Error(err)
		return
	}
	rawTx, err := SignTransaction(tx, p2, opt)
	if err != nil {
		t.Error(err)
		return
	}

	rawTxByte, _ := json.Marshal(tx)
	t.Log(string(rawTxByte))
	t.Log(rawTx)
	t.Log(hex.EncodeToString(rawTx))
}
