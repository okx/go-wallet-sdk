package crypto

import (
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/internal/util"
	"strings"
)

// PrivateKeyVariant represents the type of private key
type PrivateKeyVariant string

const (
	PrivateKeyVariantEd25519   PrivateKeyVariant = "ed25519"
	PrivateKeyVariantSecp256k1 PrivateKeyVariant = "secp256k1"
)

// AIP80Prefixes contains the AIP-80 compliant prefixes for each private key type
var AIP80Prefixes = map[PrivateKeyVariant]string{
	PrivateKeyVariantEd25519:   "ed25519-priv-",
	PrivateKeyVariantSecp256k1: "secp256k1-priv-",
}

// FormatPrivateKey formats a hex input to an AIP-80 compliant string
func FormatPrivateKey(privateKey any, keyType PrivateKeyVariant) (formattedString string, err error) {
	aip80Prefix := AIP80Prefixes[keyType]

	var hexStr string
	switch v := privateKey.(type) {
	case string:
		// Remove the prefix if it exists
		if strings.HasPrefix(v, aip80Prefix) {
			parts := strings.Split(v, "-")
			v = parts[2]
		}

		// If it's already a string, just ensure it's properly formatted
		var strBytes, err = util.ParseHex(v)
		if err != nil {
			return "", err
		}

		// Reformat to have 0x prefix
		hexStr = util.BytesToHex(strBytes)
	case []byte:
		hexStr = util.BytesToHex(v)
	default:
		return "", fmt.Errorf("unsupported private key type: must be string or []byte")
	}

	return fmt.Sprintf("%s%s", aip80Prefix, hexStr), nil
}

// ParseHexInput parses a hex input that may be bytes, hex string, or an AIP-80 compliant string to bytes.
//
// You may optionally pass in a boolean to strictly enforce AIP-80 compliance.
func ParsePrivateKey(value any, keyType PrivateKeyVariant, strict ...bool) (bytes []byte, err error) {
	aip80Prefix := AIP80Prefixes[keyType]

	// Get the first boolean if it exists, otherwise nil
	var strictness *bool = nil
	if len(strict) > 1 {
		return nil, fmt.Errorf("strictness must be a single boolean")
	} else if len(strict) == 1 {
		strictness = &strict[0]
	}

	switch v := value.(type) {
	case string:
		if (strictness == nil || !*strictness) && !strings.HasPrefix(v, aip80Prefix) {
			bytes, err := util.ParseHex(v)
			if err != nil {
				return nil, err
			}

			// If strictness is not explicitly false, warn about non-AIP-80 compliance
			if strictness == nil {
				//fmt.Printf("[Aptos SDK] It is recommended that private keys are AIP-80 compliant (https://github.com/aptos-foundation/AIPs/blob/main/aips/aip-80.md). You can fix the private key by formatting it with crypto.FormatPrivateKey\n")
			}

			return bytes, nil
		} else if strings.HasPrefix(v, aip80Prefix) {
			// Parse for AIP-80 compliant String input
			parts := strings.Split(v, "-")
			return util.ParseHex(parts[2])
		}
		return nil, fmt.Errorf("invalid hex string input while parsing private key. Must be AIP-80 compliant string")
	case []byte:
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported private key type: must be string or []byte")
	}
}
