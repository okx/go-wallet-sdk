package starknet

import "encoding/json"

func VerifyMsgSign(pub, messageHash, sign string) bool {
	curve := SC()
	pubX, pubY, err := curve.XToPubKeyErr(pub)
	if err != nil {
		return false
	}
	var signRes SignRes
	if err := json.Unmarshal([]byte(sign), &signRes); err != nil {
		return false
	}
	if pubX.Cmp(HexToBig(signRes.X)) != 0 {
		return false
	}
	return curve.Verify(HexToBig(messageHash), HexToBig(signRes.R), HexToBig(signRes.S), pubX, pubY)
}
