package amino

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

// ParseProtoPosAndTypeMustOneByte Parse field number and type from one byte,
// if original field number and type encode to multiple bytes, you should not use this function.
func ParseProtoPosAndTypeMustOneByte(data byte) (pos int, pb3Type Typ3, err error) {
	if data&0x80 == 0x80 {
		err = errors.New("func ParseProtoPosAndTypeMustOneBytevarint can't parse more than one byte")
		return
	}
	pb3Type = Typ3(data & 0x07)
	pos = int(data) >> 3
	return
}

func EncodeProtoPosAndTypeMustOneByte(pos int, typ Typ3) (byte, error) {
	// 1 1111 111
	if pos > 15 {
		return 0, fmt.Errorf("pos must be less than 16")
	}
	if typ > 7 {
		return 0, fmt.Errorf("typ must be less than 8")
	}
	data := byte(pos)
	data <<= 3
	data |= byte(typ)
	return data, nil
}

// StrToBytes is meant to make a zero allocation conversion
// from string -> []byte to speed up operations, it is not meant
// to be used generally
func StrToBytes(s string) []byte {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	var b []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	hdr.Cap = stringHeader.Len
	hdr.Len = stringHeader.Len
	hdr.Data = stringHeader.Data
	return b
}

// BytesToStr is meant to make a zero allocation conversion
// from []byte -> string to speed up operations, it is not meant
// to be used generally, but for a specific pattern to delete keys
// from a map.
func BytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func GetBinaryBareFromBinaryLengthPrefixed(bz []byte) ([]byte, error) {
	if len(bz) == 0 {
		return nil, errors.New("cannot be empty bytes")
	}

	// Read byte-length prefix.
	u64, n := binary.Uvarint(bz)
	if n < 0 {
		return nil, fmt.Errorf("Error reading msg byte-length prefix: got code %v", n)
	}
	if u64 > uint64(len(bz)-n) {
		return nil, fmt.Errorf("Not enough bytes to read, want %v more bytes but only have %v",
			u64, len(bz)-n)
	} else if u64 < uint64(len(bz)-n) {
		return nil, fmt.Errorf("Bytes left over, should read %v more bytes but have %v",
			u64, len(bz)-n)
	}
	return bz[n:], nil
}

const is64Bit = strconv.IntSize == 64

func UnmarshalBigIntBase10(bz []byte) (*big.Int, error) {
	ret := new(big.Int)
	if len(bz) < 19 {
		i, err := strconv.ParseInt(BytesToStr(bz), 10, 0)
		if err == nil {
			ret.SetInt64(i)
			return ret, nil
		}
	}

	err := ret.UnmarshalText(bz)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func MarshalBigIntToText(bi *big.Int) (string, error) {
	if bi == nil {
		// copy from big.Int.MarshalText
		return "<nil>", nil
	}
	si := bi.Sign()
	words := bi.Bits()

	if si == 0 {
		return "0", nil
	}

	var num uint64

	if is64Bit && len(words) == 1 {
		num = uint64(words[0])
	} else if !is64Bit && len(words) < 3 {
		num = bi.Uint64()
	} else {
		t, err := bi.MarshalText()
		return BytesToStr(t), err
	}

	if si > 0 {
		return strconv.FormatUint(num, 10), nil
	} else {
		if num <= uint64(math.MaxInt64)+1 {
			return strconv.FormatInt(-int64(num), 10), nil
		}
	}

	t, err := bi.MarshalText()
	return BytesToStr(t), err
}

func HexEncodeToString(src []byte) string {
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return BytesToStr(dst)
}

const hextableUpper = "0123456789ABCDEF"

func HexEncodeToStringUpper(src []byte) string {
	dst := make([]byte, hex.EncodedLen(len(src)))
	j := 0
	for _, v := range src {
		dst[j] = hextableUpper[v>>4]
		dst[j+1] = hextableUpper[v&0x0f]
		j += 2
	}
	return BytesToStr(dst)
}

func TimeSize(t time.Time) int {
	var size = 0
	s := t.Unix()
	// skip if default/zero value:
	if s != 0 {
		size += 1 + UvarintSize(uint64(s))
	}
	ns := int32(t.Nanosecond()) // this int64 -> int32 cast is safe (nanos are in [0, 999999999])
	// skip if default/zero value:
	if ns != 0 {
		// do not encode if nanos exceed allowed interval
		size += 1 + UvarintSize(uint64(ns))
	}

	return size
}

func calcUintNum(n uint64) int {
	c := 1
	n1 := n

	for ; n1 >= 100; n = n1 {
		n1 = n / 100
		c += 2
	}
	if n1 >= 10 {
		c++
	}
	return c
}

var divisor = new(big.Int).SetUint64(10000000000000000000)

type twoBigInts struct {
	A, B big.Int
}

var bigIntPool = &sync.Pool{
	New: func() interface{} {
		return new(twoBigInts)
	},
}

func CalcBigIntTextSize(bi *big.Int) int {
	if bi == nil {
		return 5 // "<nil>"
	}
	si := bi.Sign()
	words := bi.Bits()

	if si == 0 {
		return 1 // "0"
	}

	signCount := 0
	if si < 0 {
		signCount = 1
	}

	var num uint64

	if is64Bit && len(words) == 1 {
		num = uint64(words[0])
	} else if !is64Bit && len(words) < 3 {
		num = bi.Uint64()
	} else {
		wordCountOfUint64 := 1
		if !is64Bit {
			wordCountOfUint64 = 2
		}
		twoBi := bigIntPool.Get().(*twoBigInts)
		bi2 := twoBi.A.Set(bi)

		c := 0
		for len(bi2.Bits()) > wordCountOfUint64 {
			bi2.QuoRem(bi2, divisor, &twoBi.B)
			c += 1
		}
		c = calcUintNum(bi2.Uint64()) + c*19 + signCount
		bigIntPool.Put(twoBi)
		return c
	}

	return calcUintNum(num) + signCount
}
