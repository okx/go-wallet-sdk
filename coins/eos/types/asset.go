// Copyright 2018 EOS Canada <alex@eoscanada.com>
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package types

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// EOSSymbol represents the standard EOS symbol on the chain.  It's
// here just to speed up things.
var EOSSymbol = Symbol{Precision: 4, Symbol: "EOS"}

// REXSymbol represents the standard REX symbol on the chain.  It's
// here just to speed up things.
var REXSymbol = Symbol{Precision: 4, Symbol: "REX"}

// TNTSymbol represents the standard EOSIO Testnet symbol on the testnet chain.
// Temporary Network Token (TNT) is the native token of the EOSIO Testnet.
// It's here just to speed up things.
var TNTSymbol = Symbol{Precision: 4, Symbol: "TNT"}

// WAXSymbol represents the standard WAX symbol on the chain.
var WAXSymbol = Symbol{Precision: 8, Symbol: "WAX"}

type Asset struct {
	Amount int64
	Symbol
}

func (a Asset) String() string {
	amt := a.Amount
	if amt < 0 {
		amt = -amt
	}

	precisionDigitCount := int(a.Symbol.Precision)
	dotAndPrecisionDigitCount := precisionDigitCount + 1

	strInt := strconv.FormatInt(int64(amt), 10)
	if len(strInt) < dotAndPrecisionDigitCount {
		// prepend `0` for the difference:
		strInt = strings.Repeat("0", dotAndPrecisionDigitCount-len(strInt)) + strInt
	}

	result := strInt
	if a.Symbol.Precision > 0 {
		result = strInt[:len(strInt)-precisionDigitCount] + "." + strInt[len(strInt)-precisionDigitCount:]
	}

	if a.Amount < 0 {
		result = "-" + result
	}

	return fmt.Sprintf("%s %s", result, a.Symbol.Symbol)
}

func (a Asset) MarshalJSON() (data []byte, err error) {
	return json.Marshal(a.String())
}

// Symbol NOTE: there's also a new ExtendedSymbol (which includes the contract (as AccountName) on which it is)
type Symbol struct {
	Precision uint8
	Symbol    string

	// Caching of symbol code if it was computed once
	symbolCode uint64
}

type SymbolCode uint64

func (s Symbol) MustSymbolCode() (SymbolCode, error) {
	symbolCode, err := StringToSymbolCode(s.Symbol)
	if err != nil {
		return symbolCode, fmt.Errorf("invalid symbol code " + s.Symbol)
	}

	return symbolCode, nil
}

func (s Symbol) SymbolCode() (SymbolCode, error) {
	if s.symbolCode != 0 {
		return SymbolCode(s.symbolCode), nil
	}

	symbolCode, err := StringToSymbolCode(s.Symbol)
	if err != nil {
		return 0, err
	}

	return SymbolCode(symbolCode), nil
}

func (s Symbol) ToUint64() (uint64, error) {
	symbolCode, err := s.SymbolCode()
	if err != nil {
		return 0, fmt.Errorf("symbol %s is not a valid symbol code: %w", s.Symbol, err)
	}

	return uint64(symbolCode)<<8 | uint64(s.Precision), nil
}

func (s Symbol) String() string {
	return fmt.Sprintf("%d,%s", s.Precision, s.Symbol)
}

func (s Symbol) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (sc SymbolCode) String() string {
	builder := strings.Builder{}

	symbolCode := uint64(sc)
	for i := 0; i < 7; i++ {
		if symbolCode == 0 {
			return builder.String()
		}

		builder.WriteByte(byte(symbolCode & 0xFF))
		symbolCode >>= 8
	}

	return builder.String()
}

func (sc SymbolCode) MarshalJSON() (data []byte, err error) {
	return []byte(`"` + sc.String() + `"`), nil
}

func NewWAXAsset(amount int64) Asset {
	return Asset{Amount: amount, Symbol: WAXSymbol}
}

func NewEOSAsset(amount int64) Asset {
	return Asset{Amount: amount, Symbol: EOSSymbol}
}

func NewEOSAssetFromString(input string) (Asset, error) {
	return NewFixedSymbolAssetFromString(EOSSymbol, input)
}

func NewAssetFromString(input string, symbol Symbol) (Asset, error) {
	return NewFixedSymbolAssetFromString(symbol, input)
}

func StringToSymbolCode(str string) (SymbolCode, error) {
	if len(str) > 7 {
		return 0, fmt.Errorf("string is too long to be a valid symbol_code")
	}

	var symbolCode uint64
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] < 'A' || str[i] > 'Z' {
			return 0, fmt.Errorf("only uppercase letters allowed in symbol_code string")
		}

		symbolCode <<= 8
		symbolCode = symbolCode | uint64(str[i])
	}

	return SymbolCode(symbolCode), nil
}

func splitAsset(input string) (integralPart, decimalPart, symbolPart string, err error) {
	input = strings.Trim(input, " ")
	if len(input) == 0 {
		return "", "", "", fmt.Errorf("input cannot be empty")
	}

	parts := strings.Split(input, " ")
	if len(parts) >= 1 {
		integralPart, decimalPart, err = splitAssetAmount(parts[0])
		if err != nil {
			return
		}
	}

	if len(parts) == 2 {
		symbolPart = parts[1]
		if len(symbolPart) > 7 {
			return "", "", "", fmt.Errorf("invalid asset %q, symbol should have less than 7 characters", input)
		}
	}

	if len(parts) > 2 {
		return "", "", "", fmt.Errorf("invalid asset %q, expecting an amount alone or an amount and a currency symbol", input)
	}

	return
}

func splitAssetAmount(input string) (integralPart, decimalPart string, err error) {
	parts := strings.Split(input, ".")
	switch len(parts) {
	case 1:
		integralPart = parts[0]
	case 2:
		integralPart = parts[0]
		decimalPart = parts[1]

		if len(decimalPart) > math.MaxUint8 {
			err = fmt.Errorf("invalid asset amount precision %q, should have less than %d characters", input, math.MaxUint8)

		}
	default:
		return "", "", fmt.Errorf("invalid asset amount %q, expected amount to have at most a single dot", input)
	}

	return
}

func NewFixedSymbolAssetFromString(symbol Symbol, input string) (out Asset, err error) {
	integralPart, decimalPart, symbolPart, err := splitAsset(input)
	if err != nil {
		return out, err
	}

	symbolCode, err := symbol.MustSymbolCode()
	if err != nil {
		return out, err
	}
	precision := symbol.Precision

	if len(decimalPart) > int(precision) {
		return out, fmt.Errorf("symbol %s precision mismatch: expected %d, got %d", symbol, precision, len(decimalPart))
	}

	if symbolPart != "" && symbolPart != symbolCode.String() {
		return out, fmt.Errorf("symbol %s code mismatch: expected %s, got %s", symbol, symbolCode, symbolPart)
	}

	if len(decimalPart) < int(precision) {
		decimalPart += strings.Repeat("0", int(precision)-len(decimalPart))
	}

	val, err := strconv.ParseInt(integralPart+decimalPart, 10, 64)
	if err != nil {
		return out, err
	}

	return Asset{
		Amount: val,
		Symbol: Symbol{Precision: precision, Symbol: symbolCode.String()},
	}, nil
}
