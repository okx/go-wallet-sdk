/**
Author： https://github.com/hecodev007/block_sign
*/

package utils

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"github.com/emresenyuva/go-wallet-sdk/crypto/base58"
	"github.com/shopspring/decimal"
	"math/big"
	"strconv"
)

func DoubleSha256(data []byte) []byte {
	d := sha256.Sum256(data)
	dd := sha256.Sum256(d[:])
	return dd[:]

}

func CoinToFloat(val *big.Int, decimalnum int32) float64 {
	bigDecimal := decimal.NewFromBigInt(val, -decimalnum)
	fCoin, _ := bigDecimal.Truncate(decimalnum).Float64()
	return fCoin
}

// FloatToCoin coverts float64 to coin.
func FloatToCoin(val float64, decimalnum int32) *big.Int {
	bigDecimal := decimal.NewFromFloat(val).Truncate(decimalnum)
	coin := decimal.NewFromBigInt(GetTokenCoinNum(decimalnum), 0)
	bigDecimal = bigDecimal.Mul(coin)
	result, _ := new(big.Int).SetString(bigDecimal.String(), 10)
	return result
}

func GetTokenCoinNum(decimal int32) *big.Int {
	switch decimal {
	case 0:
		return big.NewInt(1e0)
	case 1:
		return big.NewInt(1e1)
	case 2:
		return big.NewInt(1e2)
	case 3:
		return big.NewInt(1e3)
	case 4:
		return big.NewInt(1e4)
	case 5:
		return big.NewInt(1e5)
	case 6:
		return big.NewInt(1e6)
	case 7:
		return big.NewInt(1e7)
	case 8:
		return big.NewInt(1e8)
	case 9:
		return big.NewInt(1e9)
	case 10:
		return big.NewInt(1e10)
	case 11:
		return big.NewInt(1e11)
	case 12:
		return big.NewInt(1e12)
	case 13:
		return big.NewInt(1e13)
	case 14:
		return big.NewInt(1e14)
	case 15:
		return big.NewInt(1e15)
	case 16:
		return big.NewInt(1e16)
	case 17:
		return big.NewInt(1e17)
	case 18:
		return big.NewInt(1e18)
	}
	return big.NewInt(1e18)
}

func MathSub(temp float64, temp2 float64, de int32) float64 {
	d := decimal.NewFromFloat(temp)
	d2 := decimal.NewFromFloat(temp2)

	d3 := d.Sub(d2).Truncate(de)

	reulsat, _ := strconv.ParseFloat(d3.String(), 64)
	return reulsat
}

func MathAdd(temp float64, temp2 float64, de int32) float64 {
	d := decimal.NewFromFloat(temp)
	d2 := decimal.NewFromFloat(temp2)

	d3 := d.Add(d2).Truncate(de)

	reulsat, _ := strconv.ParseFloat(d3.String(), 64)
	return reulsat
}

func ValidHeliumAddress(address string) error {
	if address == "" {
		return errors.New("address is null")
	}
	data := base58.Decode(address)
	checkSum1 := data[len(data)-4:]
	payload := data[:len(data)-4]
	checkSum2 := DoubleSha256(payload)[:4]
	if !bytes.Equal(checkSum1, checkSum2) {
		return errors.New("valid checksum error")
	}
	return nil
}
