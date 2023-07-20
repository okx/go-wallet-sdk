package sui

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/okx/go-wallet-sdk/crypto/base58"
	"math/big"
	"strings"
)

func WriteUint64(buf *bytes.Buffer, value uint64) error {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, value)
	_, err := buf.Write(b)
	return err
}

func WriteUint32(buf *bytes.Buffer, value uint32) error {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, value)
	_, err := buf.Write(b)
	return err
}

func WriteString(buf *bytes.Buffer, value string) error {
	if err := WriteLen(buf, len(value)); err != nil {
		return err
	}
	_, err := buf.WriteString(value)
	return err
}

func WriteBigInt(buf *bytes.Buffer, value *big.Int, bits int) error {
	if bits < 1 || bits%8 != 0 || value.Sign() < 0 || value.BitLen() > bits {
		return errors.New("invalid bits withtou multiples of 8")
	}
	b := make([]byte, bits/8)
	for i := 0; i < len(b); i += 8 {
		lo := value.Uint64()
		binary.LittleEndian.PutUint64(b[i:i+8], lo)
		value = value.Rsh(value, 64)
	}
	_, err := buf.Write(b)
	return err
}

func WriteBool(buf *bytes.Buffer, value bool) error {
	if value {
		return buf.WriteByte(1)
	} else {
		return buf.WriteByte(0)
	}
}

func DecodeHexString(s string) ([]byte, error) {
	if len(s) == 0 {
		return nil, nil
	}
	index := 0
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		index = 2
	}
	if len(s)%2 == 0 {
		return hex.DecodeString(s[index:])
	}
	return hex.DecodeString("0" + s[index:])
}

func WriteUint16(buf *bytes.Buffer, value uint16) error {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, value)
	_, err := buf.Write(b)
	return err
}

func WriteUint8(buf *bytes.Buffer, value uint8) error {
	return buf.WriteByte(byte(value))
}

func UlebEncode(num uint64) []byte {
	len, arr := 0, make([]byte, 0)
	if num == 0 {
		arr = append(arr, byte(0))
		return arr
	}
	for num > 0 {
		arr = append(arr, byte(num&0x7f))
		if num = num >> 7; num > 0 {
			arr[len] |= 0x80
		}
		len += 1
	}
	return arr
}

func WriteLen(buf *bytes.Buffer, l int) error {
	_, err := buf.Write(UlebEncode(uint64(l)))
	return err
}

func WriteHash(buf *bytes.Buffer, hash string) error {
	b := base58.Decode(hash)
	if err := WriteLen(buf, len(b)); err != nil {
		return err
	}
	_, err := buf.Write(b)
	return err
}

func WriteAddress(buf *bytes.Buffer, address string) error {
	b, err := DecodeHexString(address)
	if err != nil {
		return err
	}
	if len(b) > 32 {
		return errors.New("invalid address")
	}
	if len(b) < 32 {
		bb := make([]byte, 32)
		copy(bb[32-len(b):], b)
		b = bb
	}
	buf.Write(b)
	return nil
}
