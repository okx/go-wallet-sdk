package types

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// Denominations can be 3 ~ 128 characters long and support letters, followed by either
	// a letter, a number or a separator ('/').
	ibcReDnmString = `[a-zA-Z][a-zA-Z0-9/-]{2,127}`
	ibcReDecAmt    = `[[:digit:]]+(?:\.[[:digit:]]+)?|\.[[:digit:]]+`
	ibcReSpc       = `[[:space:]]*`
	ibcReDnm       *regexp.Regexp
	ibcReDecCoin   *regexp.Regexp
)
var ibcCoinDenomRegex = DefaultCoinDenomRegex

func init() {
	SetIBCCoinDenomRegex(DefaultIBCCoinDenomRegex)
}

// DefaultCoinDenomRegex returns the default regex string
func DefaultIBCCoinDenomRegex() string {
	return ibcReDnmString
}

func SetIBCCoinDenomRegex(reFn func() string) {
	ibcCoinDenomRegex = reFn

	ibcReDnm = regexp.MustCompile(fmt.Sprintf(`^%s$`, ibcCoinDenomRegex()))
	ibcReDecCoin = regexp.MustCompile(fmt.Sprintf(`^(%s)%s(%s)$`, ibcReDecAmt, ibcReSpc, ibcCoinDenomRegex()))
}

func validIBCCoinDenom(denom string) bool {
	return ibcReDnm.MatchString(denom)
}

func IBCParseDecCoin(coinStr string) (coin DecCoin, err error) {
	coinStr = strings.TrimSpace(coinStr)

	matches := ibcReDecCoin.FindStringSubmatch(coinStr)
	if matches == nil {
		return DecCoin{}, fmt.Errorf("invalid decimal coin expression: %s", coinStr)
	}

	amountStr, denomStr := matches[1], matches[2]

	amount, err := NewDecFromStr(amountStr)
	if err != nil {
		return DecCoin{}, fmt.Errorf("failed to parse decimal coin amount: %s", amountStr)
	}

	if err := ValidateDenom(denomStr); err != nil {
		return DecCoin{}, fmt.Errorf("invalid denom cannot contain upper case characters or spaces: %s", err)
	}

	return NewDecCoinFromDec(denomStr, amount), nil
}

const DefaultBondDenom = "okt"

// ValidateDenom validates a denomination string returning an error if it is
// invalid.
func ValidateDenom(denom string) error {
	// TODO ,height
	if denom == DefaultBondDenom {
		return nil
	}
	if !reDnm.MatchString(denom) && !rePoolTokenDnm.MatchString(denom) {
		if strings.HasPrefix(denom, "ibc") {
			if validIBCCoinDenom(denom) {
				return nil
			}
		}
		return fmt.Errorf("invalid denom: %s", denom)
	}
	return nil
}
