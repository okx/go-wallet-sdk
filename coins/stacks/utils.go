package stacks

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/emresenyuva/go-wallet-sdk/crypto/base58"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/crypto/ripemd160"
)

func decodeBtcAddress(btcAddress string) (*DecodeBtcAddressBean, error) {
	hashMode, err := getAddressHashMode(btcAddress)
	if err != nil {
		return nil, err
	}

	legacyAddress, err := FromBase58(btcAddress)
	if err != nil {
		return nil, err
	}

	return &DecodeBtcAddressBean{
		hashMode: hashMode,
		data:     legacyAddress.Bytes,
	}, nil
}

func getAddressHashMode(btcAddress string) (int, error) {
	if strings.HasPrefix(btcAddress, "bc1") || strings.HasPrefix(btcAddress, "tb1") {
		return 0, errors.New("segwit addresses are currently not supported")
	} else {
		legacyAddress, err := FromBase58(btcAddress)
		if err != nil {
			return 0, err
		}

		var version int
		if legacyAddress.P2sh {
			version = 5
		} else {
			version = 0
		}
		switch version {
		case 0:
			return SerializeP2PKH, nil
		case 111:
			return SerializeP2PKH, nil
		case 5:
			return SerializeP2SH, nil
		case 196:
			return SerializeP2SH, nil
		default:
			return 0, errors.New("getAddressHashMode error")
		}
	}
}

func intToHexString(data int, length *int) string {
	a := big.NewInt(int64(data))
	sb := fmt.Sprintf("%x", a)
	if length == nil {
		defaultLength := 8
		length = &defaultLength
	}
	s := strings.Repeat("0", *length*2-len(sb)) + sb
	return s
}

func intToBytes(value int, numBytes int) []byte {
	b := make([]byte, numBytes)
	for i := numBytes - 1; i >= 0; i-- {
		b[i] = byte(value & 0xff)
		value >>= 8
	}
	return b
}

func getBytesByLength(l int64, length int) []byte {
	newInt := big.NewInt(l)
	s := newInt.Text(16)
	if len(s) == 1 {
		s = "0" + s
	}
	for len(s) < length {
		s = "0" + s
	}
	decoded, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return decoded
}

func getBytes(i int64, length int) []byte {
	newInt := big.NewInt(i)
	s := newInt.Text(16)
	if len(s) == 1 {
		s = "0" + s
	}
	for len(s) < length {
		s = "0" + s
	}
	decoded, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return decoded
}

func sliceByteBuffer(bufferArray *bytes.Buffer) []byte {
	position := bufferArray.Len()
	b := make([]byte, position)
	bufferArray.Read(b)
	return b
}

func fromHexString(data string) []byte {
	if data == "" {
		return []byte{}
	}
	if strings.HasPrefix(data, "0x") {
		data = data[2:]
	}
	if len(data)%2 != 0 {
		data = "0" + data
	}
	result, err := hex.DecodeString(data)
	if err != nil {
		panic(err)
	}
	return result
}

func isCompressed(key StacksPublicKey) bool {
	return !strings.HasPrefix(key.Data, "04")
}

func hashP2PKH(data string) string {
	return hex.EncodeToString(hash160(data))
}

func hash160(pubKey string) []byte {
	decode, err := hex.DecodeString(pubKey)
	if err != nil {
		panic(err)
	}
	hash1 := sha256.Sum256(decode)
	hasher := ripemd160.New()
	hasher.Write(hash1[:])
	hashBytes := hasher.Sum(nil)
	return hashBytes
}

func c32addressDecode(address string) (strs []string, err error) {
	if len([]rune(address)) <= 5 {
		return nil, fmt.Errorf("invalid c32 address: invalid length")
	}
	if address[0] != 'S' {
		return nil, fmt.Errorf(`invalid c32 address: must start with "S"`)
	}
	return c32checkDecode(address[1:])
}

func c32checkDecode(address string) (strs []string, err error) {
	c32data := c32normalize(address)
	dataHex, err := c32decode(c32data[1:])
	if err != nil {
		return nil, err
	}
	versionChar := c32data[:1]
	version := strings.Index(c32, versionChar)
	checksum := dataHex[len(dataHex)-8:]

	versionHex := fmt.Sprintf("%x", version)
	if len(versionHex) == 1 {
		versionHex = "0" + versionHex
	}

	cs := versionHex + dataHex[:len(dataHex)-8]

	if r, err := c32checksum(cs); err != nil || r != checksum {
		return nil, fmt.Errorf("invalid c32check string: checksum mismatch")
	}

	strs = []string{fmt.Sprintf("%d", version), dataHex[:len(dataHex)-8]}
	return strs, nil
}

func c32checksum(data string) (checksum string, err error) {
	dataBytes, err := hex.DecodeString(data)
	if err != nil {
		return "", err
	}
	hash1 := sha256.Sum256(dataBytes)
	hash2 := sha256.Sum256(hash1[:])
	result := make([]byte, 4)
	copy(result, hash2[:4])
	checksum = hex.EncodeToString(result)
	return checksum, nil
}

func isC32(addr string) bool {
	pattern := "^[" + strings.Join([]string{c32}, "") + "]*$"
	match, err := regexp.MatchString(pattern, addr)
	if err != nil {
		return false
	}
	return match
}

func c32decode(addr string) (hexStr string, err error) {
	c32input := c32normalize(addr)
	// must result in a c32 string
	if !isC32(addr) {
		return "", errors.New("not a c32-encoded string")
	}

	var sb strings.Builder
	if strings.HasPrefix(c32input, "0") {
		chars := []rune(c32input)
		for _, aChar := range chars {
			if aChar != '0' {
				break
			}
			sb.WriteRune(aChar)
		}
	}

	var res []string
	carry := 0
	carryBits := 0
	for i := len(c32input) - 1; i >= 0; i-- {
		if carryBits == 4 {
			res = append([]string{constHex[carry : carry+1]}, res...)
			carryBits = 0
			carry = 0
		}
		currentCode := strings.Index(c32, c32input[i:i+1]) << carryBits
		currentValue := currentCode + carry
		y := currentValue % 16
		currentHexDigit := constHex[y : y+1]
		carryBits += 1
		carry = currentValue >> 4
		if carry > 1<<carryBits {
			return "", errors.New("panic error in decoding")
		}
		res = append([]string{currentHexDigit}, res...)
	}
	// one last carry
	res = append([]string{constHex[carry : carry+1]}, res...)

	if len(res)%2 == 1 {
		res = append([]string{"0"}, res...)
	}

	hexLeadingZeros := 0
	for _, value := range res {
		if value != "0" {
			break
		} else {
			hexLeadingZeros++
		}
	}

	res = res[hexLeadingZeros-(hexLeadingZeros%2):]
	hexStr = strings.Join(res, "")

	s := strings.TrimSpace(sb.String())
	if len(s) > 0 {
		for i := 0; i < len(s); i++ {
			hexStr = "00" + hexStr
		}
	}
	return hexStr, nil
}

func c32normalize(data string) (normalizedData string) {
	s := strings.ToUpper(data)
	s = strings.ReplaceAll(s, "O", "0")
	s = strings.ReplaceAll(s, "L", "1")
	s = strings.ReplaceAll(s, "I", "1")
	return s
}

func toArrayLike(bi *big.Int, length int) []byte {
	s := bi.Text(16)
	length *= 2
	sb := strings.Builder{}
	for i := 0; i < length-len(s); i++ {
		sb.WriteRune('0')
	}
	sb.WriteString(s)
	result, err := hex.DecodeString(sb.String())
	if err != nil {
		panic(err)
	}
	return result
}

func FromBase58(b58 string) (*LegacyAddress, error) {
	dataBytes, version, err := base58.CheckDecode(b58)
	if err != nil {
		return nil, err
	}
	version = version & 0xFF
	if version == 0 {
		return &LegacyAddress{
			Bytes: dataBytes,
			P2sh:  false,
		}, nil
	} else if version == 5 {
		return &LegacyAddress{
			Bytes: dataBytes,
			P2sh:  true,
		}, nil
	}

	return nil, nil
}

func c32encode(inputHex string) string {
	if len(inputHex)%2 != 0 {
		inputHex = "0" + inputHex
	}
	inputHex = strings.ToLower(inputHex)

	s := "0123456789abcdef"

	var res []string
	carry := 0
	for i := len(inputHex) - 1; i >= 0; i-- {
		if carry < 4 {
			currentCode := strings.Index(s, inputHex[i:i+1]) >> carry

			nextCode := 0
			if i != 0 {
				nextCode = strings.Index(s, inputHex[i-1:i])
			}
			nextBits := 1 + carry
			nextLowBits := (nextCode % (1 << nextBits)) << (5 - nextBits)

			sum := currentCode + nextLowBits
			curC32Digit := c32[sum : sum+1]
			carry = nextBits
			res = append([]string{curC32Digit}, res...)
		} else {
			carry = 0
		}
	}
	C32leadingZeros := 0
	for i := 0; i < len(res); i++ {
		if res[i] != "0" {
			break
		} else {
			C32leadingZeros++
		}
	}

	res = res[C32leadingZeros:]

	inputBytes, err := hex.DecodeString(inputHex)
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile(`^\x00*`)
	zeroPrefix := re.FindString(string(inputBytes))
	numLeadingZeroBytesInHex := len(zeroPrefix)

	for i := 0; i < numLeadingZeroBytesInHex; i++ {
		res = append([]string{string(c32[0])}, res...)
	}

	return strings.Join(res, "")
}

func leftPadHex(hex string) string {
	if len(hex)%2 == 0 {
		return hex
	} else {
		return "0" + hex
	}
}

func txidFromData(serialized []byte) string {
	hash := sha512.Sum512_256(serialized)
	return hex.EncodeToString(hash[:])
}

func hexToBytes(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic("invalid hex in source file: " + s)
	}
	return b
}

func hexToInt(hex string) int {
	value, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		panic(err)
	}
	return int(value)
}

func writeUInt32BE(destination []byte, value uint32, offset int) {
	destination[offset+3] = byte(value)
	value >>= 8
	destination[offset+2] = byte(value)
	value >>= 8
	destination[offset+1] = byte(value)
	value >>= 8
	destination[offset] = byte(value)
}

func asciiToBytes(str string) []byte {
	byteArray := make([]byte, len(str))
	for i, char := range str {
		byteArray[i] = byte(char & 0xff)
	}
	return byteArray
}
