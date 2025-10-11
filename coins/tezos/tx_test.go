package tezos

import (
	"encoding/hex"
	"github.com/emresenyuva/go-wallet-sdk/coins/tezos/types"
	"github.com/stretchr/testify/require"
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
	var fee int64 = 10000 // 0.010000 XTZ
	var counter int64 = 339709
	opt := NewCallOptions("BL74GqeaJ8tdFuBR2RhsGXET7MNonprQ49BBreZHE9yn9x85hJP", counter, false)
	privateKey, err := types.ParsePrivateKey(p1)
	require.NoError(t, err)
	tx, err := NewJakartanetTransaction(n1, n3, amount, opt)
	require.NoError(t, err)
	err = BuildTransaction(tx, fee, privateKey.Public(), opt)
	require.NoError(t, err)
	rawTx, err := SignTransaction(tx, p1, opt)
	require.NoError(t, err)
	expected := "33b684e3912522308951aea7e274f0f97a920d9ea268de31c2ca842cba8edd5a6c00cb15c8cb2ebe15662ad5697e139eabf3e0f1aea6904efedd1480bd3fe0d403809bee0200008e1c63b65a34abf66f88b0314549ca3295004eb700cfce6d27ca0feac5877bd24a7080c52a6c89c3378f3d45642d9e6729386e8dc1bff7b4041e8e2255d01f8ab0634bba2823314406844c026dd9b9dd9d3f989708"
	require.Equal(t, expected, hex.EncodeToString(rawTx))
}

func TestNewJakartanetDelegationTransaction(t *testing.T) {
	var to string = "tz1foXHgRzdYdaLgX6XhpZGxbBv42LZ6ubvE"
	var fee int64 = 10000
	var counter int64 = 331345
	opt := NewCallOptions("BL7kQbhcCsMYB954n94XcSZmS1oTYcg8J7ut2wj6iZpL3fRBdM3", counter, false)
	privateKey, err := types.ParsePrivateKey(p2)
	require.NoError(t, err)
	tx, err := NewJakartanetDelegationTransaction(n2, to, opt)
	require.NoError(t, err)
	err = BuildTransaction(tx, fee, privateKey.Public(), opt)
	require.NoError(t, err)
	rawTx, err := SignTransaction(tx, p2, opt)
	require.NoError(t, err)
	expected := "3548bdffd1ce6eb49efda7ebfaa9518a5868870fadf4fc0b45906412496c24376e00de6e07b12d524f72641623528d0de4da3a4cbce3904ed29c1480bd3fe0d403ff00dd2e214620a9ceaf0c38da92f9a56954f81e5ee006e44a704a0914a1c9b5adcc99fc5a238da3a1b7649a08fb5b89d66cc04c494a67fa9a40da6c1cf91ee652bf510dc6403cb654d6a0a849ed08b06a9280b9fd02"
	require.Equal(t, expected, hex.EncodeToString(rawTx))
}

func TestNewJakartanetUnDelegationTransaction(t *testing.T) {
	var fee int64 = 10000 // 0.010000 XTZ
	var counter int64 = 339708
	opt := NewCallOptions("BL7kQbhcCsMYB954n94XcSZmS1oTYcg8J7ut2wj6iZpL3fRBdM3", counter, false)
	privateKey, err := types.ParsePrivateKey(p2)
	require.NoError(t, err)
	tx, err := NewJakartanetUnDelegationTransaction(n2, opt)
	require.NoError(t, err)
	err = BuildTransaction(tx, fee, privateKey.Public(), opt)
	require.NoError(t, err)
	rawTx, err := SignTransaction(tx, p2, opt)
	require.NoError(t, err)
	expected := "3548bdffd1ce6eb49efda7ebfaa9518a5868870fadf4fc0b45906412496c24376e00de6e07b12d524f72641623528d0de4da3a4cbce3904efddd1480bd3fe0d40300d4399a50232fef86e0353b7b33962e9e86509de20c77d93cc32f475634fd7be352cf9f91fc9eae23c25fe47e0850f4c3945a96071b0fbbf2fc33c5dc0539720e"
	require.Equal(t, expected, hex.EncodeToString(rawTx))
}
