package nervos

import (
	"encoding/hex"
	"errors"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/okx/go-wallet-sdk/coins/nervos/crypto"
	"github.com/okx/go-wallet-sdk/coins/nervos/types"
	"strings"
)

type Mode string
type Type string

const (
	Mainnet Mode = "ckb"
	Testnet Mode = "ckt"

	Short       Type = "Short"
	FullBech32  Type = "FullBech32"
	FullBech32m Type = "FullBech32m"

	TYPE_FULL_WITH_BECH32M    = "00"
	ShortFormat               = "01"
	CodeHashIndexSingleSig    = "00"
	CodeHashIndexMultisigSig  = "01"
	CodeHashIndexAnyoneCanPay = "02"

	MAINNET_ACP_CODE_HASH    = "0xd369597ff47f29fbc0d47d2e3775370d1250b85140c670e4718af712983a2354"
	TESTNET_ACP_CODE_HASH    = "0x3419a1c09eb2567f6552ee7a8ecffd64155cffe0f1796e6e61ec088d740c1356"
	MAINNET_CHEQUE_CODE_HASH = "0xe4d4ecc6e5f9a059bf2f7a82cca292083aebc0c421566a52484fe2ec51a9fb0c"
	TESTNET_CHEQUE_CODE_HASH = "0x60d5f39efce409c587cb9ea359cefdead650ca128f0bd9cb3855348f98c70d5b"

	AnyoneCanPayCodeHashOnLina   = "0xd369597ff47f29fbc0d47d2e3775370d1250b85140c670e4718af712983a2354"
	AnyoneCanPayCodeHashOnAggron = "0x3419a1c09eb2567f6552ee7a8ecffd64155cffe0f1796e6e61ec088d740c1356"
)

type AddressGenerateResult struct {
	Address    string
	LockArgs   string
	PrivateKey string
}

func GenerateAddress() (*AddressGenerateResult, error) {
	return GenerateBech32mFullAddress(Mainnet)
}

// GenerateTestnetAddress generates a testnet address.
func GenerateTestnetAddress() (*AddressGenerateResult, error) {
	return GenerateBech32mFullAddress(Testnet)
}

// GenerateBech32mFullAddress generates a bech32m address with a full script.
func GenerateBech32mFullAddress(mode Mode) (*AddressGenerateResult, error) {
	key, err := crypto.RandomNew()
	if err != nil {
		return nil, err
	}

	address, err := GenerateBech32mFullAddressByPublicKey(mode, key.PubKey())
	if err != nil {
		return nil, err
	}
	pubKey, err := crypto.Blake160(key.PubKey())
	if err != nil {
		return nil, err
	}

	return &AddressGenerateResult{
		Address:    address,
		LockArgs:   types.Encode(pubKey),
		PrivateKey: types.Encode(key.Bytes()),
	}, err
}

func GenerateBech32mFullAddressByPublicKey(mode Mode, publicKey []byte) (string, error) {
	pubKey, err := crypto.Blake160(publicKey)
	if err != nil {
		return "", err
	}

	script := &types.Script{
		CodeHash: types.HexToHash(types.SECP256K1_BLAKE160_SIGHASH_ALL_TYPE_HASH),
		HashType: types.HashTypeType,
		Args:     types.FromHex(hex.EncodeToString(pubKey)),
	}

	address, err := ConvertScriptToBech32mFullAddress(mode, script)
	if err != nil {
		return "", err
	}

	return address, err
}

func GenerateAddressByPrivateKey(mode string, privateKey string) (string, error) {
	key, err := crypto.HexToKey(privateKey)
	if err != nil {
		return "", err
	}
	if mode != "ckt" {
		mode = "ckb"
	}
	return GenerateBech32mFullAddressByPublicKey(Mode(mode), key.PubKey())
}

func ConvertScriptToAddress(mode Mode, script *types.Script) (string, error) {
	return ConvertScriptToBech32mFullAddress(mode, script)
}

// ConvertScriptToBech32mFullAddress converts a script to a bech32m full address.
func ConvertScriptToBech32mFullAddress(mode Mode, script *types.Script) (string, error) {
	hashType, err := types.SerializeHashType(script.HashType)
	if err != nil {
		return "", err
	}
	// https://github.com/nervosnetwork/rfcs/blob/master/rfcs/0021-ckb-address-format/0021-ckb-address-format.md
	// Payload: type(00) | code hash | hash type | args
	payload := TYPE_FULL_WITH_BECH32M
	payload += script.CodeHash.Hex()[2:]
	payload += hashType

	payload += hex.EncodeToString(script.Args)

	dataPart, err := bech32.ConvertBits(types.FromHex(payload), 8, 5, true)
	if err != nil {
		return "", err
	}
	return bech32.EncodeM(string(mode), dataPart)
}

type ParsedAddress struct {
	Mode   Mode
	Type   Type
	Script *types.Script
}

func Parse(address string) (*ParsedAddress, error) {
	encoding, hrp, decoded, err := crypto.Bech32Decode(address)
	if err != nil {
		return nil, err
	}
	data, err := bech32.ConvertBits(decoded, 5, 8, false)
	if err != nil {
		return nil, err
	}
	payload := hex.EncodeToString(data)

	var addressType Type
	var script types.Script
	if strings.HasPrefix(payload, "01") {
		if encoding != crypto.BECH32 {
			return nil, errors.New("payload header 0x01 should have encoding BECH32")
		}
		addressType = Short
		if CodeHashIndexSingleSig == payload[2:4] {
			if len(payload) != 44 {
				return nil, errors.New("payload bytes length of secp256k1-sighash-all " +
					"short address should be 22")
			}
			script = types.Script{
				CodeHash: types.HexToHash(types.SECP256K1_BLAKE160_SIGHASH_ALL_TYPE_HASH),
				HashType: types.HashTypeType,
				Args:     types.Hex2Bytes(payload[4:]),
			}
		} else if CodeHashIndexAnyoneCanPay == payload[2:4] {
			if len(payload) < 44 || len(payload) > 48 {
				return nil, errors.New("payload bytes length of acp short address should between 22-24")
			}
			script = types.Script{
				HashType: types.HashTypeType,
				Args:     types.Hex2Bytes(payload[4:]),
			}
			if hrp == (string)(Testnet) {
				script.CodeHash = types.HexToHash(AnyoneCanPayCodeHashOnAggron)
			} else {
				script.CodeHash = types.HexToHash(AnyoneCanPayCodeHashOnLina)
			}
		} else if CodeHashIndexMultisigSig == payload[2:4] {
			if len(payload) != 44 {
				return nil, errors.New("payload bytes length of secp256k1-multisig-all " +
					"short address should be 22")
			}
			script = types.Script{
				CodeHash: types.HexToHash(types.SECP256K1_BLAKE160_MULTISIG_ALL_TYPE_HASH),
				HashType: types.HashTypeType,
				Args:     types.Hex2Bytes(payload[4:]),
			}
		} else {
			return nil, errors.New("unknown code hash index " + payload[2:4])
		}
	} else if strings.HasPrefix(payload, "02") {
		if encoding != crypto.BECH32 {
			return nil, errors.New("payload header 0x02 should have encoding BECH32")
		}
		addressType = FullBech32
		script = types.Script{
			CodeHash: types.HexToHash(payload[2:66]),
			HashType: types.HashTypeData,
			Args:     types.Hex2Bytes(payload[66:]),
		}
	} else if strings.HasPrefix(payload, "04") {
		if encoding != crypto.BECH32 {
			return nil, errors.New("payload header 0x04 should have encoding BECH32")
		}
		addressType = FullBech32
		script = types.Script{
			CodeHash: types.HexToHash(payload[2:66]),
			HashType: types.HashTypeType,
			Args:     types.Hex2Bytes(payload[66:]),
		}
	} else if strings.HasPrefix(payload, "00") {
		if encoding != crypto.BECH32M {
			return nil, errors.New("payload header 0x00 should have encoding BECH32")
		}
		addressType = FullBech32m
		script = types.Script{
			CodeHash: types.HexToHash(payload[2:66]),
			Args:     types.Hex2Bytes(payload[68:]),
		}

		hashType, err := types.DeserializeHashType(payload[66:68])
		if err != nil {
			return nil, err
		}

		script.HashType = hashType

	} else {
		return nil, errors.New("address type error:" + payload[:2])
	}

	result := &ParsedAddress{
		Mode:   Mode(hrp),
		Type:   addressType,
		Script: &script,
	}
	return result, nil
}

func ValidateAddress(address string) bool {
	_, err := Parse(address)
	return err == nil
}
