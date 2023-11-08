package crypto

import "github.com/btcsuite/btcd/btcutil/bech32"

type Encoding uint

const (
	BECH32 Encoding = iota
	BECH32M
)

const (
	charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"
)

var gen = []int{0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3}

func Bech32Decode(bech string) (Encoding, string, []byte, error) {
	hrp, decoded, err := bech32.DecodeNoLimit(bech)
	if err != nil {
		return 0, "", nil, err
	}

	ints := make([]int, len(decoded))
	for i := 0; i < len(decoded); i++ {
		ints[i] = int(decoded[i])
	}

	polymod := append(bech32HrpExpand(hrp), ints...)
	i := bech32Polymod(polymod)
	var encoding Encoding
	if i == 1 {
		encoding = BECH32
	} else {
		encoding = BECH32M
	}

	return encoding, hrp, decoded, nil
}

func bech32Polymod(values []int) int {
	chk := 1
	for _, v := range values {
		b := chk >> 25
		chk = (chk&0x1ffffff)<<5 ^ v
		for i := 0; i < 5; i++ {
			if (b>>uint(i))&1 == 1 {
				chk ^= gen[i]
			}
		}
	}
	return chk
}

func bech32HrpExpand(hrp string) []int {
	v := make([]int, 0, len(hrp)*2+1)
	for i := 0; i < len(hrp); i++ {
		v = append(v, int(hrp[i]>>5))
	}
	v = append(v, 0)
	for i := 0; i < len(hrp); i++ {
		v = append(v, int(hrp[i]&31))
	}
	return v
}
