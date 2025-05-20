package sei

import (
	"github.com/okx/go-wallet-sdk/coins/cosmos"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestNewAddress(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	address, err := NewAddress(privateKeyHex)
	require.Nil(t, err)
	expected := "sei145q0tcdur4tcx2ya5cphqx96e54yflfyd7jmd4"
	require.Equal(t, expected, address)

	ret := ValidateAddress(address)
	require.True(t, ret)
}

// txHash: https://sei.explorers.guru/transaction/08D076CFE1903AB697974E7CB756F5E9A3344FF07220892CE2A75EBB29494435
// https://rest.atlantic-2.seinetwork.io/cosmos/auth/v1beta1/accounts/sei1s95zvpxwxcr0ykdkj3ymscrevdam7wvs24dk57
// curl -X POST -d '{"tx_bytes":"CosBCogBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEmgKKnNlaTFzOTV6dnB4d3hjcjB5a2RrajN5bXNjcmV2ZGFtN3d2czI0ZGs1NxIqc2VpMXVyZGVkZWVqMGZkNGtzbHpuM3VxNnM4bW5kaDh3dDd1c2s2YTR6Gg4KBHVzZWkSBjEwMDAwMBJmClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECWzs64TTLQ3sGP88eAUtzXoHtGHUYauDmYZWgBLyYUesSBAoCCAEYARISCgwKBHVzZWkSBDEwMDAQoI0GGkAPt3BsqAL807wgpPtKQdF8mYsPwM52HjRaScsc27rIdh30d6JxWnu9Zy1Tm9funsAYIOtStq7GKTfekctaIRK/","mode":"BROADCAST_MODE_SYNC"}' https://rest.atlantic-2.seinetwork.io/cosmos/tx/v1beta1/txs
func TestTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	address, err := NewAddress(privateKeyHex)
	require.Nil(t, err)
	expected := "sei145q0tcdur4tcx2ya5cphqx96e54yflfyd7jmd4"
	require.Equal(t, expected, address)

	param := cosmos.TransferParam{}
	param.FromAddress = address
	param.ToAddress = address
	param.Denom = "usei"
	param.Amount = "100000"
	param.CommonParam.ChainId = "atlantic-2"
	param.CommonParam.Sequence = 1
	param.CommonParam.AccountNumber = 4050874
	param.CommonParam.FeeDenom = "usei"
	param.CommonParam.FeeAmount = "1000"
	param.CommonParam.GasLimit = 100000
	param.CommonParam.Memo = ""
	param.CommonParam.TimeoutHeight = 0
	signedTx, err := cosmos.Transfer(param, privateKeyHex)
	require.Nil(t, err)
	expected = "CosBCogBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEmgKKnNlaTE0NXEwdGNkdXI0dGN4MnlhNWNwaHF4OTZlNTR5ZmxmeWQ3am1kNBIqc2VpMTQ1cTB0Y2R1cjR0Y3gyeWE1Y3BocXg5NmU1NHlmbGZ5ZDdqbWQ0Gg4KBHVzZWkSBjEwMDAwMBJmClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDEFPp7wKV0zS2uyLiDMcX6xoWpUb2klcsiDC0vBTBNnYSBAoCCAEYARISCgwKBHVzZWkSBDEwMDAQoI0GGkDJivGlkzrZffAXInJAHY2QBOgD8QraVrDhLyX5yAjr+QLdwloCq5maBrniL3PmBrRUrPxFNC8Wt5QklQ70YYDF"
	require.Equal(t, expected, signedTx)
}

// https://sei.explorers.guru/transaction/E07194819858ED6C8BF355AF55A1F57E6346C45E3A9C5CDEE2DFFDEA3992BE5B
// curl -X POST -d '{"tx_bytes":"CrsBCrgBCikvaWJjLmFwcGxpY2F0aW9ucy50cmFuc2Zlci52MS5Nc2dUcmFuc2ZlchKKAQoIdHJhbnNmZXISCWNoYW5uZWwtNxoMCgR1c2VpEgQxMDAwIipzZWkxczk1enZweHd4Y3IweWtka2ozeW1zY3JldmRhbTd3dnMyNGRrNTcqLWF4ZWxhcjFydnM1eHBoNGwzcHgyZWZ5bnFzdGh1czhwNnI0ZXh5cjZrcXZkZDIAOIDs+ImekvirFxJmClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECWzs64TTLQ3sGP88eAUtzXoHtGHUYauDmYZWgBLyYUesSBAoCCAEYBBISCgwKBHVzZWkSBDEwMDAQoI0GGkB2KRZkV1uDATVA9Tk7rpqCeSu6b6j1POSrOKQLX2Vijxoj5hpIRm32sKykpf+O+ED7QiZwzF24+WTowx9eAMwb","mode":"BROADCAST_MODE_SYNC"}' https://rest.atlantic-2.seinetwork.io/cosmos/tx/v1beta1/txs
func TestIbcTransfer(t *testing.T) {
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	p := cosmos.IbcTransferParam{}
	p.CommonParam.ChainId = "atlantic-2"
	p.CommonParam.Sequence = 4
	p.CommonParam.AccountNumber = 4050874
	p.CommonParam.FeeDenom = "usei"
	p.CommonParam.FeeAmount = "1000"
	p.CommonParam.GasLimit = 100000
	p.CommonParam.Memo = ""
	p.CommonParam.TimeoutHeight = 0
	p.FromAddress = "sei1s95zvpxwxcr0ykdkj3ymscrevdam7wvs24dk57"
	p.ToAddress = "axelar1rvs5xph4l3px2efynqsthus8p6r4exyr6kqvdd"
	p.Denom = "usei"
	p.Amount = "1000"
	p.SourcePort = "transfer"
	p.SourceChannel = "channel-7"
	p.TimeOutInSeconds = uint64(time.Now().UnixMilli()/1000) + 300
	signedIBCTx, err := cosmos.IbcTransfer(p, privateKeyHex)
	require.Nil(t, err)
	t.Log("signedIBCTx : ", signedIBCTx)
}

// https://www.seiscan.app/atlantic-2/txs/B564BE5BFA5ED5096287F72F11778E3BE83F24CCD6074B28BACD5BBF0BFC9A3C
// curl -X POST -d '{"tx_bytes":"Cs0BCsoBCiQvY29zbXdhc20ud2FzbS52MS5Nc2dFeGVjdXRlQ29udHJhY3QSoQEKKnNlaTFzOTV6dnB4d3hjcjB5a2RrajN5bXNjcmV2ZGFtN3d2czI0ZGs1NxI+c2VpMTY1bDVuY2gzMDlyY2d2NzY1anlkZHdmMzV0YW5lM3QyMDllZjM1NXpmNzUzazI2czc2NHNxa242cTkaKHsidXBkYXRlX3ByaWNlIjp7ImJhaWxfb25fZXJyb3IiOmZhbHNlfX0qCQoEdXNlaRIBMRJoClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiECWzs64TTLQ3sGP88eAUtzXoHtGHUYauDmYZWgBLyYUesSBAoCCAEYEhIUCg4KBHVzZWkSBjI1MDAwMBCAiXoaQFhyD8Iv42M5OeJJe0TvHU3PN5XmHYSVXx8MsIjdX1A+BW5et4Sqa3vQoF0pjDzEzGc93dFllESWJZF7QtywyYs=","mode":"BROADCAST_MODE_SYNC"}' https://rest.atlantic-2.seinetwork.io/cosmos/tx/v1beta1/txs
func TestSignMessage(t *testing.T) {
	// sei1s95zvpxwxcr0ykdkj3ymscrevdam7wvs24dk57
	privateKeyHex := "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"
	data := "{\n" +
		"\"chain_id\": \"atlantic-2\",\n" +
		"\"account_number\": \"4050874\",\n" +
		"\"sequence\": \"18\",\n" +
		"\"fee\": {\n " +
		"	\"gas\": \"2000000\",\n" +
		"	\"amount\": [\n    " +
		"		{\n    " +
		"			\"denom\": \"usei\",\n" +
		"			\"amount\": \"250000\"\n" +
		"       }\n" +
		"    ]\n" +
		"},\n" +
		"\"msgs\": [\n" +
		" {\n" +
		"	\"type\": \"cosmwasm.wasm.v1.MsgExecuteContract\",\n " +
		"	\"value\": {\n" +
		"		\"sender\": \"sei1s95zvpxwxcr0ykdkj3ymscrevdam7wvs24dk57\",\n " +
		"		\"contract\": \"sei165l5nch309rcgv765jyddwf35tane3t209ef355zf753k26s764sqkn6q9\",\n " +
		"		\"funds\": [\n" +
		"	  		{\n" +
		"	 			\"amount\": \"1\",\n" +
		"   			\"denom\": \"usei\"\n" +
		"     		}\n" +
		" 		],\n" +
		"		\"msg\": {\n" +
		"			\"update_price\": {\n" +
		"				\"bail_on_error\": false\n" +
		"	    	}\n" +
		"		}\n" +
		" 	}\n" +
		" }\n" +
		"],\n " +
		"\"memo\": \"\"\n" +
		"}"
	tt, _, err := cosmos.SignMessage(data, privateKeyHex)
	require.Nil(t, err)
	t.Log(tt)
}
