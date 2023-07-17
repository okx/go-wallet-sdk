package polkadot

func GetEra(height uint64, calPeriod uint64) []byte {
	if calPeriod == 0 {
		calPeriod = 64
	}
	phase := height % calPeriod
	index := uint64(6)
	trailingZero := index - 1

	var encoded uint64
	if trailingZero > 1 {
		encoded = trailingZero
	} else {
		encoded = 1
	}

	if trailingZero < 15 {
		encoded = trailingZero
	} else {
		encoded = 15
	}
	encoded += phase / 1 << 4
	first := byte(encoded >> 8)
	second := byte(encoded & 0xff)
	return []byte{second, first}
}
