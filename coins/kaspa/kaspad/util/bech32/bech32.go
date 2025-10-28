// Copyright (c) 2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package bech32

import (
	"fmt"
	"github.com/pkg/errors"
	"strings"
)

const charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"
const checksumLength = 8

// For use in convertBits. Represents a number of bits to convert to or from and whether
// to add padding.
type conversionType struct {
	fromBits uint8
	toBits   uint8
	pad      bool
}

// Conversion types to use in convertBits.
var fiveToEightBits = conversionType{fromBits: 5, toBits: 8, pad: false}
var eightToFiveBits = conversionType{fromBits: 8, toBits: 5, pad: true}

var generator = []int64{0x98f2bc8e61, 0x79b76d99e2, 0xf33e5fb3c4, 0xae2eabe2a8, 0x1e4f43e470}

// Encode prepends the version byte, converts to uint5, and encodes to Bech32.
func Encode(prefix string, payload []byte, version byte) string {
	data := make([]byte, len(payload)+1)
	data[0] = version
	copy(data[1:], payload)

	converted := convertBits(data, eightToFiveBits)

	return encode(prefix, converted)
}

// Decode decodes a string that was encoded with Encode.
func Decode(encoded string) (string, []byte, byte, error) {
	prefix, decoded, err := decode(encoded)
	if err != nil {
		return "", nil, 0, err
	}

	converted := convertBits(decoded, fiveToEightBits)
	version := converted[0]
	payload := converted[1:]

	return prefix, payload, version, nil
}

// Decode decodes a Bech32 encoded string, returning the prefix
// and the data part excluding the checksum.
func decode(encoded string) (string, []byte, error) {
	// The minimum allowed length for a Bech32 string is 10 characters,
	// since it needs a non-empty prefix, a separator, and an 8 character
	// checksum.
	if len(encoded) < checksumLength+2 {
		return "", nil, errors.Errorf("invalid bech32 string length %d",
			len(encoded))
	}
	// Only	ASCII characters between 33 and 126 are allowed.
	for i := 0; i < len(encoded); i++ {
		if encoded[i] < 33 || encoded[i] > 126 {
			return "", nil, errors.Errorf("invalid character in "+
				"string: '%c'", encoded[i])
		}
	}

	// The characters must be either all lowercase or all uppercase.
	lower := strings.ToLower(encoded)
	upper := strings.ToUpper(encoded)
	if encoded != lower && encoded != upper {
		return "", nil, errors.Errorf("string not all lowercase or all " +
			"uppercase")
	}

	// We'll work with the lowercase string from now on.
	encoded = lower

	// The string is invalid if the last ':' is non-existent, it is the
	// first character of the string (no human-readable part) or one of the
	// last 8 characters of the string (since checksum cannot contain ':'),
	// or if the string is more than 90 characters in total.
	colonIndex := strings.LastIndexByte(encoded, ':')
	if colonIndex < 1 || colonIndex+checksumLength+1 > len(encoded) {
		return "", nil, errors.Errorf("invalid index of ':'")
	}

	// The prefix part is everything before the last ':'.
	prefix := encoded[:colonIndex]
	data := encoded[colonIndex+1:]

	// Each character corresponds to the byte with value of the index in
	// 'charset'.
	decoded, err := decodeFromBase32(data)
	if err != nil {
		return "", nil, errors.Errorf("failed converting data to bytes: "+
			"%s", err)
	}

	if !verifyChecksum(prefix, decoded) {
		checksum := encoded[len(encoded)-checksumLength:]
		expected := encodeToBase32(calculateChecksum(prefix,
			decoded[:len(decoded)-checksumLength]))

		return "", nil, errors.Errorf("checksum failed. Expected %s, got %s",
			expected, checksum)
	}

	// We exclude the last 8 bytes, which is the checksum.
	return prefix, decoded[:len(decoded)-checksumLength], nil
}

// Encode encodes a byte slice into a bech32 string with the
// prefix. Note that the bytes must each encode 5 bits (base32).
func encode(prefix string, data []byte) string {
	// Calculate the checksum of the data and append it at the end.
	checksum := calculateChecksum(prefix, data)
	combined := append(data, checksum...)

	// The resulting bech32 string is the concatenation of the prefix, the
	// separator ':', data and checksum. Everything after the separator is
	// represented using the specified charset.
	base32String := encodeToBase32(combined)

	return fmt.Sprintf("%s:%s", prefix, base32String)
}

// decodeFromBase32 converts each character in the string 'chars' to the value of the
// index of the correspoding character in 'charset'.
func decodeFromBase32(base32String string) ([]byte, error) {
	decoded := make([]byte, 0, len(base32String))
	for i := 0; i < len(base32String); i++ {
		index := strings.IndexByte(charset, base32String[i])
		if index < 0 {
			return nil, errors.Errorf("invalid character not part of "+
				"charset: %c", base32String[i])
		}
		decoded = append(decoded, byte(index))
	}
	return decoded, nil
}

// Converts the byte slice 'data' to a string where each byte in 'data'
// encodes the index of a character in 'charset'.
// IMPORTANT: this function expects the data to be in uint5 format.
// CAUTION: for legacy reasons, in case of an error this function returns
// an empty string instead of an error.
func encodeToBase32(data []byte) string {
	result := make([]byte, 0, len(data))
	for _, b := range data {
		if int(b) >= len(charset) {
			return ""
		}
		result = append(result, charset[b])
	}
	return string(result)
}

// convertBits converts a byte slice where each byte is encoding fromBits bits,
// to a byte slice where each byte is encoding toBits bits.
func convertBits(data []byte, conversionType conversionType) []byte {
	// The final bytes, each byte encoding toBits bits.
	var regrouped []byte

	// Keep track of the next byte we create and how many bits we have
	// added to it out of the toBits goal.
	nextByte := byte(0)
	filledBits := uint8(0)

	for _, b := range data {
		// Discard unused bits.
		b = b << (8 - conversionType.fromBits)

		// How many bits remaining to extract from the input data.
		remainingFromBits := conversionType.fromBits
		for remainingFromBits > 0 {
			// How many bits remaining to be added to the next byte.
			remainingToBits := conversionType.toBits - filledBits

			// The number of bytes to next extract is the minimum of
			// remainingFromBits and remainingToBits.
			toExtract := remainingFromBits
			if remainingToBits < toExtract {
				toExtract = remainingToBits
			}

			// Add the next bits to nextByte, shifting the already
			// added bits to the left.
			nextByte = (nextByte << toExtract) | (b >> (8 - toExtract))

			// Discard the bits we just extracted and get ready for
			// next iteration.
			b = b << toExtract
			remainingFromBits -= toExtract
			filledBits += toExtract

			// If the nextByte is completely filled, we add it to
			// our regrouped bytes and start on the next byte.
			if filledBits == conversionType.toBits {
				regrouped = append(regrouped, nextByte)
				filledBits = 0
				nextByte = 0
			}
		}
	}

	// We pad any unfinished group if specified.
	if conversionType.pad && filledBits > 0 {
		nextByte = nextByte << (conversionType.toBits - filledBits)
		regrouped = append(regrouped, nextByte)
		filledBits = 0
		nextByte = 0
	}

	return regrouped
}

// The checksum is a 40 bits BCH codes defined over GF(2^5).
// It ensures the detection of up to 6 errors in the address and 8 in a row.
// Combined with the length check, this provides very strong guarantee against errors.
// For more details please refer to the Bech32 Address Serialization section
// of the spec.
func calculateChecksum(prefix string, payload []byte) []byte {
	prefixLower5Bits := prefixToUint5Array(prefix)
	payloadInts := ints(payload)
	templateZeroes := []int64{0, 0, 0, 0, 0, 0, 0, 0}

	// prefixLower5Bits + 0 + payloadInts + templateZeroes
	concat := append(prefixLower5Bits, 0)
	concat = append(concat, payloadInts...)
	concat = append(concat, templateZeroes...)

	polyModResult := polyMod(concat)
	var res []byte
	for i := 0; i < checksumLength; i++ {
		res = append(res, byte((polyModResult>>uint(5*(checksumLength-1-i)))&31))
	}

	return res
}

// For more details please refer to the Bech32 Address Serialization section
// of the spec.
func verifyChecksum(prefix string, payload []byte) bool {
	prefixLower5Bits := prefixToUint5Array(prefix)
	payloadInts := ints(payload)

	// prefixLower5Bits + 0 + payloadInts
	dataToVerify := append(prefixLower5Bits, 0)
	dataToVerify = append(dataToVerify, payloadInts...)

	return polyMod(dataToVerify) == 0
}

func prefixToUint5Array(prefix string) []int64 {
	prefixLower5Bits := make([]int64, len(prefix))
	for i := 0; i < len(prefix); i++ {
		char := prefix[i]
		charLower5Bits := int64(char & 31)
		prefixLower5Bits[i] = charLower5Bits
	}

	return prefixLower5Bits
}

func ints(payload []byte) []int64 {
	payloadInts := make([]int64, len(payload))
	for i, b := range payload {
		payloadInts[i] = int64(b)
	}

	return payloadInts
}

// For more details please refer to the Bech32 Address Serialization section
// of the spec.
func polyMod(values []int64) int64 {
	checksum := int64(1)
	for _, value := range values {
		topBits := checksum >> 35
		checksum = ((checksum & 0x07ffffffff) << 5) ^ value
		for i := 0; i < len(generator); i++ {
			if ((topBits >> uint(i)) & 1) == 1 {
				checksum ^= generator[i]
			}
		}
	}

	return checksum ^ 1
}
