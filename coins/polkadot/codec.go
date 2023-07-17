package polkadot

import (
	"encoding/binary"
	"encoding/hex"
	"math/big"
)

const (
	modeBits                  = 2
	singleMode           byte = 0
	twoByteMode          byte = 1
	fourByteMode         byte = 2
	bigIntMode           byte = 3
	singleModeMaxValue        = 63
	twoByteModeMaxValue       = 16383
	fourByteModeMaxValue      = 1073741823
)

var modeToNumOfBytes = map[byte]uint{
	singleMode:   1,
	twoByteMode:  2,
	fourByteMode: 4,
}

func CompactLength(data uint32) int {
	if data >= 0 && data <= singleModeMaxValue {
		return 1
	} else if data <= twoByteModeMaxValue {
		return 2
	} else if data <= fourByteModeMaxValue {
		return 4
	} else {
		return 5
	}
}

func ExtendLEBytes(input []byte, length int) []byte {
	diff := length - len(input)
	if diff == 0 {
		return input
	}
	for i := 0; i < diff; i++ {
		input = append(input, 0)
	}
	return input
}

func uint32ToLittleEndianBytes(data uint32) []byte {
	tmp := [4]byte{}
	binary.LittleEndian.PutUint32(tmp[:], data)
	return tmp[:]
}

func removeExtraLEBytes(input []byte) []byte {
	index := len(input)
	for {
		if input[index-1] != 0 {
			break
		} else {
			index--
		}
	}
	return input[:index]
}

func BytesToCompactBytes(bytes []byte) (res []byte) {
	lenOfBytes := len(bytes)
	if lenOfBytes > 4 {
		zeroByte := len(bytes) - 4
		zeroByte = zeroByte << modeBits
		zeroByte |= int(bigIntMode)

		res = []byte{byte(zeroByte)}
		res = append(res, bytes...)
	} else {
		mode := fourByteMode

		switch lenOfBytes {
		case 1:
			mode = singleMode
		case 2:
			mode = twoByteMode
		}

		var nextRepl byte
		for i := range bytes {
			repl := bytes[i] & 192
			repl = repl >> 6
			bytes[i] = bytes[i] << modeBits
			if i != 0 {
				bytes[i] |= nextRepl
			}
			nextRepl = repl
		}
		if nextRepl != 0 {
			bytes = append(bytes, nextRepl)
		}
		bytes[0] |= mode
		bytes = ExtendLEBytes(bytes, int(modeToNumOfBytes[mode]))

		res = bytes
	}
	return
}

func getLeadingZeros(data uint64) int {
	return 64 - big.NewInt(int64(data)).BitLen()
}

func getPrefix(data uint64) string {
	return hex.EncodeToString([]byte{byte(3 + (((8 - getLeadingZeros(data)/8) - 4) << 2))})
}

func uint64ToLittleEndianArrayWithoutLeadingZeros(data uint64) string {
	tmp := [8]byte{}
	binary.LittleEndian.PutUint64(tmp[:], data)
	return hex.EncodeToString(removeExtraLEBytes(tmp[:]))
}

func Encode(data uint64) string {
	if data > fourByteModeMaxValue {
		return getPrefix(data) + uint64ToLittleEndianArrayWithoutLeadingZeros(data)
	}
	bytes := uint32ToLittleEndianBytes(uint32(data))
	bytes = removeExtraLEBytes(bytes)
	compactLength := CompactLength(uint32(data))
	length := len(bytes)
	if length < compactLength {
		for i := 0; i < compactLength-length; i++ {
			bytes = append(bytes, 0)
		}
	}

	ret := BytesToCompactBytes(bytes)
	if compactLength == 5 {
		ret[0] = 0x03
	}
	return hex.EncodeToString(ret)
}
