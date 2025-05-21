package ethereum

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/okx/go-wallet-sdk/util"
	"golang.org/x/crypto/sha3"
	"math/big"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// TypedDataAndHash is a helper function that calculates a hash for typed data conforming to EIP-712.
// This hash can then be safely used to calculate a signature.
//
// See https://eips.ethereum.org/EIPS/eip-712 for the full specification.
//
// This gives context to the signed typed data and prevents signing of transactions.
func TypedDataAndHash(typedData TypedData) ([]byte, string, error) {
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, "", err
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, "", err
	}
	rawData := fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash))
	hash := sha3.NewLegacyKeccak256()
	hash.Write([]byte(rawData))
	h := hash.Sum(nil)
	return h, rawData, nil
}

func EIP712Hash(typedData TypedData) string {
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return ""
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return ""
	}

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	hash := keccak256(rawData)
	return "0x" + hex.EncodeToString(hash)
}

func EIP712ParamsHash(typedData TypedData) string {
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return ""
	}
	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return ""
	}
	//rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	res := map[string]string{}
	res["domainSeparator"] = hex.EncodeToString(domainSeparator)
	res["typedDataHash"] = hex.EncodeToString(typedDataHash)
	jsonValue, err := json.Marshal(res)
	if err != nil {
		return ""
	}
	return string(jsonValue)
}

func EIP712HashWithTypedDataHash(domainSeparatorHex, typedDataHashHex string) string {
	domainSeparator := util.RemoveZeroHex(domainSeparatorHex)
	typedDataHash := util.RemoveZeroHex(typedDataHashHex)
	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	hash := keccak256(rawData)
	return "0x" + hex.EncodeToString(hash)
}

func keccak256(data []byte) []byte {
	s256 := sha3.NewLegacyKeccak256()
	s256.Write(data)
	txId := s256.Sum(nil)
	txIdHash := hex.EncodeToString(txId)
	res, _ := hex.DecodeString(txIdHash)
	return res
}

type decError struct{ msg string }

// Errors
var (
	ErrEmptyString   = &decError{"empty hex string"}
	ErrSyntax        = &decError{"invalid hex string"}
	ErrMissingPrefix = &decError{"hex string without 0x prefix"}
	ErrOddLength     = &decError{"hex string of odd length"}
	ErrUint64Range   = &decError{"hex number > 64 bits"}
)

func (err decError) Error() string { return err.msg }

func has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

func mapError(err error) error {
	if err, ok := err.(*strconv.NumError); ok {
		switch err.Err {
		case strconv.ErrRange:
			return ErrUint64Range
		case strconv.ErrSyntax:
			return ErrSyntax
		}
	}
	if _, ok := err.(hex.InvalidByteError); ok {
		return ErrSyntax
	}
	if err == hex.ErrLength {
		return ErrOddLength
	}
	return err
}

func Decode(input string) ([]byte, error) {
	if len(input) == 0 {
		return nil, ErrEmptyString
	}
	if !has0xPrefix(input) {
		return nil, ErrMissingPrefix
	}
	b, err := hex.DecodeString(input[2:])
	if err != nil {
		err = mapError(err)
	}
	return b, err
}

/////////// Address

// Address represents the 20 byte address of an Ethereum account.
type Address [20]byte

// BytesToAddress returns Address with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-20:]
	}
	copy(a[20-len(b):], b)
}

// BigToAddress returns Address with byte values of b.
// If b is larger than len(h), b will be cropped from the left.
func BigToAddress(b *big.Int) Address { return BytesToAddress(b.Bytes()) }

// FromHex returns the bytes represented by the hexadecimal string s.
// s may be prefixed with "0x".
func FromHex(s string) []byte {
	if has0xPrefix(s) {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}

// Hex2Bytes returns the bytes represented by the hexadecimal string str.
func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

// HexToAddress returns Address with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToAddress(s string) Address { return BytesToAddress(FromHex(s)) }

// IsHexAddress verifies whether a string can represent a valid hex-encoded
// Ethereum address or not.
func IsHexAddress(s string) bool {
	if has0xPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*20 && isHex(s)
}

// isHex validates whether each byte is valid hexadecimal string.
func isHex(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}

// Bytes gets the string representation of the underlying address.
func (a Address) Bytes() []byte { return a[:] }

// Hex returns an EIP55-compliant hex string representation of the address.
func (a Address) Hex() string {
	return string(a.checksumHex())
}

func (a Address) hex() []byte {
	var buf [len(a)*2 + 2]byte
	copy(buf[:2], "0x")
	hex.Encode(buf[2:], a[:])
	return buf[:]
}

func (a *Address) checksumHex() []byte {
	buf := a.hex()
	// compute checksum
	sha := sha3.NewLegacyKeccak256()
	sha.Write(buf[2:])
	hash := sha.Sum(nil)
	for i := 2; i < len(buf); i++ {
		hashByte := hash[(i-2)/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if buf[i] > '9' && hashByte > 7 {
			buf[i] -= 32
		}
	}
	return buf[:]
}

// String implements fmt.Stringer.
func (a Address) String() string {
	return a.Hex()
}

type TypedData struct {
	Types       Types            `json:"types"`
	PrimaryType string           `json:"primaryType"`
	Domain      TypedDataDomain  `json:"domain"`
	Message     TypedDataMessage `json:"message"`
}

type Type struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (t *Type) isArray() bool {
	return strings.HasSuffix(t.Type, "[]")
}

// typeName returns the canonical name of the type. If the type is 'Person[]', then
// this method returns 'Person'
func (t *Type) typeName() string {
	if strings.HasSuffix(t.Type, "[]") {
		return strings.TrimSuffix(t.Type, "[]")
	}
	return t.Type
}

func (t *Type) isReferenceType() bool {
	if len(t.Type) == 0 {
		return false
	}
	// Reference types must have a leading uppercase character
	return unicode.IsUpper([]rune(t.Type)[0])
}

type Types map[string][]Type

type TypePriority struct {
	Type  string
	Value uint
}

type TypedDataMessage = map[string]interface{}

type TypedDataDomain struct {
	Name              string   `json:"name"`
	Version           string   `json:"version"`
	ChainId           *big.Int `json:"chainId"`
	VerifyingContract string   `json:"verifyingContract"`
	Salt              string   `json:"salt"`
}

var typedDataReferenceTypeRegexp = regexp.MustCompile(`^[A-Z](\w*)(\[\])?$`)

// HashStruct generates a keccak256 hash of the encoding of the provided data
func (typedData *TypedData) HashStruct(primaryType string, data TypedDataMessage) ([]byte, error) {
	encodedData, err := typedData.EncodeData(primaryType, data, 1)
	if err != nil {
		return nil, err
	}
	return keccak256(encodedData), nil
}

// Dependencies returns an array of custom types ordered by their hierarchical reference tree
func (typedData *TypedData) Dependencies(primaryType string, found []string) []string {
	primaryType = strings.Split(primaryType, "[")[0]

	if contains(found, primaryType) {
		return found
	}
	if typedData.Types[primaryType] == nil {
		return found
	}
	found = append(found, primaryType)
	for _, field := range typedData.Types[primaryType] {
		for _, dep := range typedData.Dependencies(field.Type, found) {
			if !contains(found, dep) {
				found = append(found, dep)
			}
		}
	}
	return found
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// EncodeType generates the following encoding:
// `name ‖ "(" ‖ member₁ ‖ "," ‖ member₂ ‖ "," ‖ … ‖ memberₙ ")"`
//
// each member is written as `type ‖ " " ‖ name` encodings cascade down and are sorted by name
func (typedData *TypedData) EncodeType(primaryType string) hexutil.Bytes {
	// Get dependencies primary first, then alphabetical
	deps := typedData.Dependencies(primaryType, []string{})
	if len(deps) > 0 {
		slicedDeps := deps[1:]
		sort.Strings(slicedDeps)
		deps = append([]string{primaryType}, slicedDeps...)
	}

	// Format as a string with fields
	var buffer bytes.Buffer
	for _, dep := range deps {
		buffer.WriteString(dep)
		buffer.WriteString("(")
		for _, obj := range typedData.Types[dep] {
			buffer.WriteString(obj.Type)
			buffer.WriteString(" ")
			buffer.WriteString(obj.Name)
			buffer.WriteString(",")
		}
		buffer.Truncate(buffer.Len() - 1)
		buffer.WriteString(")")
	}
	return buffer.Bytes()
}

// TypeHash creates the keccak256 hash  of the data
func (typedData *TypedData) TypeHash(primaryType string) []byte {
	return keccak256(typedData.EncodeType(primaryType))
}

// EncodeData generates the following encoding:
// `enc(value₁) ‖ enc(value₂) ‖ … ‖ enc(valueₙ)`
//
// each encoded member is 32-byte long
func (typedData *TypedData) EncodeData(primaryType string, data map[string]interface{}, depth int) ([]byte, error) {
	if err := typedData.validate(); err != nil {
		return nil, err
	}

	buffer := bytes.Buffer{}

	// Verify extra data
	if exp, got := len(typedData.Types[primaryType]), len(data); exp < got {
		return nil, fmt.Errorf("there is extra data provided in the message (%d < %d)", exp, got)
	}

	// Add typehash
	buffer.Write(typedData.TypeHash(primaryType))
	// Add field contents. Structs and arrays have special handlers.
	for _, field := range typedData.Types[primaryType] {
		encType := field.Type
		encValue := data[field.Name]
		if encType[len(encType)-1:] == "]" {
			encodedData, err := typedData.encodeArrayValue(encValue, encType, depth)
			if err != nil {
				return nil, err
			}
			buffer.Write(encodedData)
		} else if typedData.Types[field.Type] != nil {
			mapValue, ok := encValue.(map[string]interface{})
			if !ok {
				return nil, dataMismatchError(encType, encValue)
			}
			encodedData, err := typedData.EncodeData(field.Type, mapValue, depth+1)
			if err != nil {
				return nil, err
			}
			buffer.Write(keccak256(encodedData))
		} else {
			byteValue, err := typedData.EncodePrimitiveValue(encType, encValue, depth)
			if err != nil {
				return nil, err
			}
			buffer.Write(byteValue)
		}
	}

	return buffer.Bytes(), nil
}

func (typedData *TypedData) encodeArrayValue(encValue interface{}, encType string, depth int) (hexutil.Bytes, error) {
	arrayValue, err := convertDataToSlice(encValue)
	if err != nil {
		return nil, dataMismatchError(encType, encValue)
	}

	arrayBuffer := new(bytes.Buffer)
	parsedType := strings.Split(encType, "[")[0]
	for _, item := range arrayValue {
		if reflect.TypeOf(item).Kind() == reflect.Slice ||
			reflect.TypeOf(item).Kind() == reflect.Array {
			encodedData, err := typedData.encodeArrayValue(item, parsedType, depth+1)
			if err != nil {
				return nil, err
			}
			arrayBuffer.Write(encodedData)
		} else {
			if typedData.Types[parsedType] != nil {
				mapValue, ok := item.(map[string]interface{})
				if !ok {
					return nil, dataMismatchError(parsedType, item)
				}
				encodedData, err := typedData.EncodeData(parsedType, mapValue, depth+1)
				if err != nil {
					return nil, err
				}
				digest := keccak256(encodedData)
				arrayBuffer.Write(digest)
			} else {
				bytesValue, err := typedData.EncodePrimitiveValue(parsedType, item, depth)
				if err != nil {
					return nil, err
				}
				arrayBuffer.Write(bytesValue)
			}
		}
	}
	return keccak256(arrayBuffer.Bytes()), nil
}

func convertDataToSlice(encValue interface{}) ([]interface{}, error) {
	var outEncValue []interface{}
	rv := reflect.ValueOf(encValue)
	if rv.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			outEncValue = append(outEncValue, rv.Index(i).Interface())
		}
	} else {
		return outEncValue, fmt.Errorf("provided data '%v' is not slice", encValue)
	}
	return outEncValue, nil
}

// Attempt to parse bytes in different formats: byte array, hex string, hexutil.Bytes.
func parseBytes(encType interface{}) ([]byte, bool) {
	switch v := encType.(type) {
	case []byte:
		return v, true
	case string:
		bytes, err := Decode(v)
		if err != nil {
			return nil, false
		}
		return bytes, true
	default:
		return nil, false
	}
}

func parseInteger(encType string, encValue interface{}) (*big.Int, error) {
	var (
		length int
		signed = strings.HasPrefix(encType, "int")
		b      *big.Int
	)
	if encType == "int" || encType == "uint" {
		length = 256
	} else {
		lengthStr := ""
		if strings.HasPrefix(encType, "uint") {
			lengthStr = strings.TrimPrefix(encType, "uint")
		} else {
			lengthStr = strings.TrimPrefix(encType, "int")
		}
		atoiSize, err := strconv.Atoi(lengthStr)
		if err != nil {
			return nil, fmt.Errorf("invalid size on integer: %v", lengthStr)
		}
		length = atoiSize
	}
	switch v := encValue.(type) {
	case *big.Int:
		b = (*big.Int)(v)
	case string:
		var hexIntValue big.Int
		if err := hexIntValue.UnmarshalText([]byte(v)); err != nil {
			return nil, err
		}
		b = (*big.Int)(&hexIntValue)
	case float64:
		// JSON parses non-strings as float64. Fail if we cannot
		// convert it losslessly
		if float64(int64(v)) == v {
			b = big.NewInt(int64(v))
		} else {
			return nil, fmt.Errorf("invalid float value %v for type %v", v, encType)
		}
	}
	if b == nil {
		return nil, fmt.Errorf("invalid integer value %v/%v for type %v", encValue, reflect.TypeOf(encValue), encType)
	}
	if b.BitLen() > length {
		return nil, fmt.Errorf("integer larger than '%v'", encType)
	}
	if !signed && b.Sign() == -1 {
		return nil, fmt.Errorf("invalid negative value for unsigned type %v", encType)
	}
	return b, nil
}

// EncodePrimitiveValue deals with the primitive values found
// while searching through the typed data
func (typedData *TypedData) EncodePrimitiveValue(encType string, encValue interface{}, depth int) ([]byte, error) {
	switch encType {
	case "address":
		stringValue, ok := encValue.(string)
		if !ok || !IsHexAddress(stringValue) {
			return nil, dataMismatchError(encType, encValue)
		}
		retval := make([]byte, 32)
		copy(retval[12:], HexToAddress(stringValue).Bytes())
		return retval, nil
	case "bool":
		boolValue, ok := encValue.(bool)
		if !ok {
			return nil, dataMismatchError(encType, encValue)
		}
		if boolValue {
			big1 := big.NewInt(1)
			return PaddedBigBytes(big1, 32), nil
		}
		big0 := big.NewInt(0)
		return PaddedBigBytes(big0, 32), nil
	case "string":
		strVal, ok := encValue.(string)
		if !ok {
			return nil, dataMismatchError(encType, encValue)
		}
		return keccak256([]byte(strVal)), nil
	case "bytes":
		bytesValue, ok := parseBytes(encValue)
		if !ok {
			return nil, dataMismatchError(encType, encValue)
		}
		return keccak256(bytesValue), nil
	}
	if strings.HasPrefix(encType, "bytes") {
		lengthStr := strings.TrimPrefix(encType, "bytes")
		length, err := strconv.Atoi(lengthStr)
		if err != nil {
			return nil, fmt.Errorf("invalid size on bytes: %v", lengthStr)
		}
		if length < 0 || length > 32 {
			return nil, fmt.Errorf("invalid size on bytes: %d", length)
		}
		if byteValue, ok := parseBytes(encValue); !ok || len(byteValue) != length {
			return nil, dataMismatchError(encType, encValue)
		} else {
			// Right-pad the bits
			dst := make([]byte, 32)
			copy(dst, byteValue)
			return dst, nil
		}
	}
	if strings.HasPrefix(encType, "int") || strings.HasPrefix(encType, "uint") {
		b, err := parseInteger(encType, encValue)
		if err != nil {
			return nil, err
		}
		return U256Bytes(b), nil
	}
	return nil, fmt.Errorf("unrecognized type '%s'", encType)
}

// BigPow returns a ** b as a big integer.
func BigPow(a, b int64) *big.Int {
	r := big.NewInt(a)
	return r.Exp(r, big.NewInt(b), nil)
}

// Various big integer limit values.
var (
	tt256   = BigPow(2, 256)
	tt256m1 = new(big.Int).Sub(tt256, big.NewInt(1))
)

// U256 encodes as a 256 bit two's complement number. This operation is destructive.
func U256(x *big.Int) *big.Int {
	return x.And(x, tt256m1)
}

// U256Bytes converts a big Int into a 256bit EVM number.
// This operation is destructive.
func U256Bytes(n *big.Int) []byte {
	return PaddedBigBytes(U256(n), 32)
}

// PaddedBigBytes encodes a big integer as a big-endian byte slice. The length
// of the slice is at least n bytes.
func PaddedBigBytes(bigint *big.Int, n int) []byte {
	if bigint.BitLen()/8 >= n {
		return bigint.Bytes()
	}
	ret := make([]byte, n)
	ReadBits(bigint, ret)
	return ret
}

const (
	// number of bits in a big.Word
	wordBits = 32 << (uint64(^big.Word(0)) >> 63)
	// number of bytes in a big.Word
	wordBytes = wordBits / 8
)

// ReadBits encodes the absolute value of bigint as big-endian bytes. Callers must ensure
// that buf has enough space. If buf is too short the result will be incomplete.
func ReadBits(bigint *big.Int, buf []byte) {
	i := len(buf)
	for _, d := range bigint.Bits() {
		for j := 0; j < wordBytes && i > 0; j++ {
			i--
			buf[i] = byte(d)
			d >>= 8
		}
	}
}

// dataMismatchError generates an error for a mismatch between
// the provided type and data
func dataMismatchError(encType string, encValue interface{}) error {
	return fmt.Errorf("provided data '%v' doesn't match type '%s'", encValue, encType)
}

// validate makes sure the types are sound
func (typedData *TypedData) validate() error {
	if err := typedData.Types.validate(); err != nil {
		return err
	}
	if err := typedData.Domain.validate(); err != nil {
		return err
	}
	return nil
}

// Map generates a map version of the typed data
func (typedData *TypedData) Map() map[string]interface{} {
	dataMap := map[string]interface{}{
		"types":       typedData.Types,
		"domain":      typedData.Domain.Map(),
		"primaryType": typedData.PrimaryType,
		"message":     typedData.Message,
	}
	return dataMap
}

// Format returns a representation of typedData, which can be easily displayed by a user-interface
// without in-depth knowledge about 712 rules
func (typedData *TypedData) Format() ([]*NameValueType, error) {
	domain, err := typedData.formatData("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, err
	}
	ptype, err := typedData.formatData(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, err
	}
	var nvts []*NameValueType
	nvts = append(nvts, &NameValueType{
		Name:  "EIP712Domain",
		Value: domain,
		Typ:   "domain",
	})
	nvts = append(nvts, &NameValueType{
		Name:  typedData.PrimaryType,
		Value: ptype,
		Typ:   "primary type",
	})
	return nvts, nil
}

func (typedData *TypedData) formatData(primaryType string, data map[string]interface{}) ([]*NameValueType, error) {
	var output []*NameValueType

	// Add field contents. Structs and arrays have special handlers.
	for _, field := range typedData.Types[primaryType] {
		encName := field.Name
		encValue := data[encName]
		item := &NameValueType{
			Name: encName,
			Typ:  field.Type,
		}
		if field.isArray() {
			arrayValue, _ := encValue.([]interface{})
			parsedType := field.typeName()
			for _, v := range arrayValue {
				if typedData.Types[parsedType] != nil {
					mapValue, _ := v.(map[string]interface{})
					mapOutput, err := typedData.formatData(parsedType, mapValue)
					if err != nil {
						return nil, err
					}
					item.Value = mapOutput
				} else {
					primitiveOutput, err := formatPrimitiveValue(field.Type, encValue)
					if err != nil {
						return nil, err
					}
					item.Value = primitiveOutput
				}
			}
		} else if typedData.Types[field.Type] != nil {
			if mapValue, ok := encValue.(map[string]interface{}); ok {
				mapOutput, err := typedData.formatData(field.Type, mapValue)
				if err != nil {
					return nil, err
				}
				item.Value = mapOutput
			} else {
				item.Value = "<nil>"
			}
		} else {
			primitiveOutput, err := formatPrimitiveValue(field.Type, encValue)
			if err != nil {
				return nil, err
			}
			item.Value = primitiveOutput
		}
		output = append(output, item)
	}
	return output, nil
}

func formatPrimitiveValue(encType string, encValue interface{}) (string, error) {
	switch encType {
	case "address":
		if stringValue, ok := encValue.(string); !ok {
			return "", fmt.Errorf("could not format value %v as address", encValue)
		} else {
			return HexToAddress(stringValue).String(), nil
		}
	case "bool":
		if boolValue, ok := encValue.(bool); !ok {
			return "", fmt.Errorf("could not format value %v as bool", encValue)
		} else {
			return fmt.Sprintf("%t", boolValue), nil
		}
	case "bytes", "string":
		return fmt.Sprintf("%s", encValue), nil
	}
	if strings.HasPrefix(encType, "bytes") {
		return fmt.Sprintf("%s", encValue), nil

	}
	if strings.HasPrefix(encType, "uint") || strings.HasPrefix(encType, "int") {
		if b, err := parseInteger(encType, encValue); err != nil {
			return "", err
		} else {
			return fmt.Sprintf("%d (0x%x)", b, b), nil
		}
	}
	return "", fmt.Errorf("unhandled type %v", encType)
}

// NameValueType is a very simple struct with Name, Value and Type. It's meant for simple
// json structures used to communicate signing-info about typed data with the UI
type NameValueType struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
	Typ   string      `json:"type"`
}

// Pprint returns a pretty-printed version of nvt
func (nvt *NameValueType) Pprint(depth int) string {
	output := bytes.Buffer{}
	output.WriteString(strings.Repeat("\u00a0", depth*2))
	output.WriteString(fmt.Sprintf("%s [%s]: ", nvt.Name, nvt.Typ))
	if nvts, ok := nvt.Value.([]*NameValueType); ok {
		output.WriteString("\n")
		for _, next := range nvts {
			sublevel := next.Pprint(depth + 1)
			output.WriteString(sublevel)
		}
	} else {
		if nvt.Value != nil {
			output.WriteString(fmt.Sprintf("%q\n", nvt.Value))
		} else {
			output.WriteString("\n")
		}
	}
	return output.String()
}

// Validate checks if the types object is conformant to the specs
func (t Types) validate() error {
	for typeKey, typeArr := range t {
		if len(typeKey) == 0 {
			return fmt.Errorf("empty type key")
		}
		for i, typeObj := range typeArr {
			if len(typeObj.Type) == 0 {
				return fmt.Errorf("type %q:%d: empty Type", typeKey, i)
			}
			if len(typeObj.Name) == 0 {
				return fmt.Errorf("type %q:%d: empty Name", typeKey, i)
			}
			if typeKey == typeObj.Type {
				return fmt.Errorf("type %q cannot reference itself", typeObj.Type)
			}
			if typeObj.isReferenceType() {
				if _, exist := t[typeObj.typeName()]; !exist {
					return fmt.Errorf("reference type %q is undefined", typeObj.Type)
				}
				if !typedDataReferenceTypeRegexp.MatchString(typeObj.Type) {
					return fmt.Errorf("unknown reference type %q", typeObj.Type)
				}
			} else if !isPrimitiveTypeValid(typeObj.Type) {
				return fmt.Errorf("unknown type %q", typeObj.Type)
			}
		}
	}
	return nil
}

// Checks if the primitive value is valid
func isPrimitiveTypeValid(primitiveType string) bool {
	if primitiveType == "address" ||
		primitiveType == "address[]" ||
		primitiveType == "bool" ||
		primitiveType == "bool[]" ||
		primitiveType == "string" ||
		primitiveType == "string[]" {
		return true
	}
	if primitiveType == "bytes" ||
		primitiveType == "bytes[]" ||
		primitiveType == "bytes1" ||
		primitiveType == "bytes1[]" ||
		primitiveType == "bytes2" ||
		primitiveType == "bytes2[]" ||
		primitiveType == "bytes3" ||
		primitiveType == "bytes3[]" ||
		primitiveType == "bytes4" ||
		primitiveType == "bytes4[]" ||
		primitiveType == "bytes5" ||
		primitiveType == "bytes5[]" ||
		primitiveType == "bytes6" ||
		primitiveType == "bytes6[]" ||
		primitiveType == "bytes7" ||
		primitiveType == "bytes7[]" ||
		primitiveType == "bytes8" ||
		primitiveType == "bytes8[]" ||
		primitiveType == "bytes9" ||
		primitiveType == "bytes9[]" ||
		primitiveType == "bytes10" ||
		primitiveType == "bytes10[]" ||
		primitiveType == "bytes11" ||
		primitiveType == "bytes11[]" ||
		primitiveType == "bytes12" ||
		primitiveType == "bytes12[]" ||
		primitiveType == "bytes13" ||
		primitiveType == "bytes13[]" ||
		primitiveType == "bytes14" ||
		primitiveType == "bytes14[]" ||
		primitiveType == "bytes15" ||
		primitiveType == "bytes15[]" ||
		primitiveType == "bytes16" ||
		primitiveType == "bytes16[]" ||
		primitiveType == "bytes17" ||
		primitiveType == "bytes17[]" ||
		primitiveType == "bytes18" ||
		primitiveType == "bytes18[]" ||
		primitiveType == "bytes19" ||
		primitiveType == "bytes19[]" ||
		primitiveType == "bytes20" ||
		primitiveType == "bytes20[]" ||
		primitiveType == "bytes21" ||
		primitiveType == "bytes21[]" ||
		primitiveType == "bytes22" ||
		primitiveType == "bytes22[]" ||
		primitiveType == "bytes23" ||
		primitiveType == "bytes23[]" ||
		primitiveType == "bytes24" ||
		primitiveType == "bytes24[]" ||
		primitiveType == "bytes25" ||
		primitiveType == "bytes25[]" ||
		primitiveType == "bytes26" ||
		primitiveType == "bytes26[]" ||
		primitiveType == "bytes27" ||
		primitiveType == "bytes27[]" ||
		primitiveType == "bytes28" ||
		primitiveType == "bytes28[]" ||
		primitiveType == "bytes29" ||
		primitiveType == "bytes29[]" ||
		primitiveType == "bytes30" ||
		primitiveType == "bytes30[]" ||
		primitiveType == "bytes31" ||
		primitiveType == "bytes31[]" ||
		primitiveType == "bytes32" ||
		primitiveType == "bytes32[]" {
		return true
	}
	if primitiveType == "int" ||
		primitiveType == "int[]" ||
		primitiveType == "int8" ||
		primitiveType == "int8[]" ||
		primitiveType == "int16" ||
		primitiveType == "int16[]" ||
		primitiveType == "int32" ||
		primitiveType == "int32[]" ||
		primitiveType == "int64" ||
		primitiveType == "int64[]" ||
		primitiveType == "int128" ||
		primitiveType == "int128[]" ||
		primitiveType == "int256" ||
		primitiveType == "int256[]" {
		return true
	}
	if primitiveType == "uint" ||
		primitiveType == "uint[]" ||
		primitiveType == "uint8" ||
		primitiveType == "uint8[]" ||
		primitiveType == "uint16" ||
		primitiveType == "uint16[]" ||
		primitiveType == "uint32" ||
		primitiveType == "uint32[]" ||
		primitiveType == "uint64" ||
		primitiveType == "uint64[]" ||
		primitiveType == "uint128" ||
		primitiveType == "uint128[]" ||
		primitiveType == "uint256" ||
		primitiveType == "uint256[]" {
		return true
	}
	return false
}

// validate checks if the given domain is valid, i.e. contains at least
// the minimum viable keys and values
func (domain *TypedDataDomain) validate() error {
	if domain.ChainId == nil && len(domain.Name) == 0 && len(domain.Version) == 0 && len(domain.VerifyingContract) == 0 && len(domain.Salt) == 0 {
		return errors.New("domain is undefined")
	}
	return nil
}

// Map is a helper function to generate a map version of the domain
func (domain *TypedDataDomain) Map() map[string]interface{} {
	dataMap := map[string]interface{}{}

	if domain.ChainId != nil {
		dataMap["chainId"] = domain.ChainId
	}

	if len(domain.Name) > 0 {
		dataMap["name"] = domain.Name
	}

	if len(domain.Version) > 0 {
		dataMap["version"] = domain.Version
	}

	if len(domain.VerifyingContract) > 0 {
		dataMap["verifyingContract"] = domain.VerifyingContract
	}

	if len(domain.Salt) > 0 {
		dataMap["salt"] = domain.Salt
	}
	return dataMap
}
