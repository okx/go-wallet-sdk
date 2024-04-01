package aptos

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/aptos/aptos_types"
	"github.com/okx/go-wallet-sdk/coins/aptos/transaction_builder"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/okx/go-wallet-sdk/coins/aptos/serde"
	"github.com/okx/go-wallet-sdk/crypto/ed25519"
)

func ShortenAddress(address string) string {
	re, _ := regexp.Compile("^0x0*")
	return re.ReplaceAllString(address, "0x")
}

func ExpandAddress(address string) string {
	rest := strings.TrimPrefix(address, "0x")
	if len(rest) < 64 {
		return "0x" + strings.Repeat("0", 64-len(rest)) + rest
	} else {
		return address
	}
}

func NewAddress(seedHex string, shortEnable bool) (string, error) {
	publicKey, err := ed25519.PublicKeyFromSeed(seedHex)
	if err != nil {
		return "", err
	}
	publicKey = append(publicKey, 0x0)
	address := "0x" + hex.EncodeToString(aptos_types.Sha256Hash(publicKey))
	if shortEnable {
		return ShortenAddress(address), nil
	} else {
		return address, nil
	}
}

func GetAddressByPubKey(pubKeyHex string, shortEnable bool) (string, error) {
	pubKey, err := hex.DecodeString(pubKeyHex)
	if err != nil {
		return "", err
	}

	pubKey = append(pubKey, 0x0)
	address := "0x" + hex.EncodeToString(aptos_types.Sha256Hash(pubKey))
	if shortEnable {
		return ShortenAddress(address), nil
	} else {
		return address, nil
	}
}

// hex 32bytes
func ValidateAddress(address string, shortEnable bool) bool {
	re1, _ := regexp.Compile("^0x[\\dA-Fa-f]{62,64}$")
	re2, _ := regexp.Compile("^[\\dA-Fa-f]{64}$")
	return re1.Match([]byte(address)) || re2.Match([]byte(address))
}

func MakeRawTransaction(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64,
	expirationTimestampSecs uint64, chainId uint8, payload aptos_types.TransactionPayload) (*aptos_types.RawTransaction, error) {
	rawTxn := aptos_types.RawTransaction{}
	// addr, err := aptos_types.FromHex(ExpandAddress(from))
	addr, err := aptos_types.FromHex(from)
	if err != nil {
		return nil, err
	}
	rawTxn.Sender = *addr
	rawTxn.SequenceNumber = sequenceNumber
	rawTxn.MaxGasAmount = maxGasAmount
	rawTxn.GasUnitPrice = gasUnitPrice
	rawTxn.ExpirationTimestampSecs = expirationTimestampSecs
	rawTxn.ChainId = aptos_types.ChainId(chainId)
	rawTxn.Payload = payload
	return &rawTxn, nil
}

func Transfer(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64, expirationTimestampSecs uint64, chainId uint8,
	to string, amount uint64, seedHex string) (string, error) {
	payload, err := TransferPayload(to, amount)
	if err != nil {
		return "", err
	}
	return BuildSignedTransaction(from, sequenceNumber, maxGasAmount, gasUnitPrice, expirationTimestampSecs, chainId, payload, seedHex)
}

func TransferPayload(to string, amount uint64) (aptos_types.TransactionPayload, error) {
	//moduleAddress := make([]byte, 31)
	//moduleAddress = append(moduleAddress, 0x1)
	bscAddress, err := aptos_types.BcsSerializeFixedBytes(aptos_types.BytesFromHex(ExpandAddress(to)))
	if err != nil {
		return nil, err
	}
	bscAmount, err := aptos_types.BcsSerializeUint64(amount)
	if err != nil {
		return nil, err
	}
	scriptFunction := aptos_types.ScriptFunction{
		Module:   aptos_types.ModuleId{Address: *aptos_types.CORE_CODE_ADDRESS, Name: "aptos_account"},
		Function: "transfer",
		TyArgs:   []aptos_types.TypeTag{},
		Args:     [][]byte{bscAddress, bscAmount},
	}
	return &aptos_types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
}

func BuildSignedTransaction(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64, expirationTimestampSecs uint64, chainId uint8,
	payload aptos_types.TransactionPayload, seedHex string) (string, error) {
	rawTxn, err := MakeRawTransaction(from, sequenceNumber, maxGasAmount, gasUnitPrice, expirationTimestampSecs, chainId, payload)
	if err != nil {
		return "", err
	}
	return SignRawTransaction(rawTxn, seedHex)
}

// pub enum Transaction {
// /// Transaction submitted by the user. e.g: P2P payment transaction, publishing module
// /// transaction, etc.
// /// TODO: We need to rename SignedTransaction to SignedUserTransaction, as well as all the other
// ///       transaction aptos_types we had in our codebase.
// UserTransaction(SignedTransaction),
//
// /// Transaction that applies a WriteSet to the current storage, it's applied manually via db-bootstrapper.
// GenesisTransaction(WriteSetPayload),
//
// /// Transaction to update the block metadata resource at the beginning of a block.
// BlockMetadata(BlockMetadata),
//
// /// Transaction to let the executor update the global state tree and record the root hash
// /// in the TransactionInfo
// /// The hash value inside is unique block id which can generate unique hash of state checkpoint transaction
// StateCheckpoint(HashValue),
// }
func GetTransactionHash(hexStr string) (string, error) {
	prefix := aptos_types.Sha256Hash([]byte("APTOS::Transaction"))
	bcsBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}
	message := make([]byte, 0)
	message = append(message, prefix...)
	message = append(message, 0x0)
	message = append(message, bcsBytes...)
	return "0x" + hex.EncodeToString(aptos_types.Sha256Hash(message)), nil
}

func CoinTransferPayload(to string, amount uint64, tyArg string) (aptos_types.TransactionPayload, error) {
	bscAddress, err := aptos_types.BcsSerializeFixedBytes(aptos_types.BytesFromHex(ExpandAddress(to)))
	if err != nil {
		return nil, err
	}
	bscAmount, err := aptos_types.BcsSerializeUint64(amount)
	if err != nil {
		return nil, err
	}
	// 0x3::moon_coin::MoonCoin  address（hex） + module + structure
	parts := strings.Split(tyArg, "::")
	contractAddr, err := aptos_types.FromHex(ExpandAddress(parts[0]))
	if err != nil {
		return nil, err
	}
	tyArgs := make([]aptos_types.TypeTag, 0)
	t1 := aptos_types.TypeTagStruct{
		Value: aptos_types.StructTag{
			Address:    *contractAddr,
			ModuleName: aptos_types.Identifier(parts[1]),
			Name:       aptos_types.Identifier(parts[2]),
			TypeArgs:   []aptos_types.TypeTag{},
		},
	}
	tyArgs = append(tyArgs, &t1)
	// 0x1::coin transfer
	scriptFunction := aptos_types.ScriptFunction{
		Module:   aptos_types.ModuleId{Address: *aptos_types.CORE_CODE_ADDRESS, Name: "coin"},
		Function: "transfer",
		TyArgs:   tyArgs,
		Args:     [][]byte{bscAddress, bscAmount},
	}
	return &aptos_types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
}

func CoinRegisterPayload(tyArg string) (aptos_types.TransactionPayload, error) {
	moduleAddress := make([]byte, 31)
	moduleAddress = append(moduleAddress, 0x1)

	parts := strings.Split(tyArg, "::")
	contractAddr, err := aptos_types.FromHex(ExpandAddress(parts[0]))
	if err != nil {
		return nil, err
	}
	tyArgs := make([]aptos_types.TypeTag, 0)
	t1 := aptos_types.TypeTagStruct{
		Value: aptos_types.StructTag{
			Address:    *contractAddr,
			ModuleName: aptos_types.Identifier(parts[1]),
			Name:       aptos_types.Identifier(parts[2]),
			TypeArgs:   []aptos_types.TypeTag{},
		},
	}
	tyArgs = append(tyArgs, &t1)

	scriptFunction := aptos_types.ScriptFunction{
		Module:   aptos_types.ModuleId{Address: *aptos_types.CORE_CODE_ADDRESS, Name: "managed_coin"},
		Function: "register",
		TyArgs:   tyArgs,
		Args:     [][]byte{},
	}
	return &aptos_types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
}

func CoinMintPayload(receiveAddress string, amount uint64, tyArg string) (aptos_types.TransactionPayload, error) {
	moduleAddress := make([]byte, 31)
	moduleAddress = append(moduleAddress, 0x1)

	bscAddress, _ := aptos_types.BcsSerializeFixedBytes(aptos_types.BytesFromHex(receiveAddress))
	bscAmount, _ := aptos_types.BcsSerializeUint64(amount)

	parts := strings.Split(tyArg, "::")
	contractAddr, err := aptos_types.FromHex(ExpandAddress(parts[0]))
	if err != nil {
		return nil, err
	}
	tyArgs := make([]aptos_types.TypeTag, 0)
	t1 := aptos_types.TypeTagStruct{
		Value: aptos_types.StructTag{
			Address:    *contractAddr,
			ModuleName: aptos_types.Identifier(parts[1]),
			Name:       aptos_types.Identifier(parts[2]),
			TypeArgs:   []aptos_types.TypeTag{},
		},
	}
	tyArgs = append(tyArgs, &t1)

	scriptFunction := aptos_types.ScriptFunction{
		Module:   aptos_types.ModuleId{Address: *aptos_types.CORE_CODE_ADDRESS, Name: "managed_coin"},
		Function: "mint",
		TyArgs:   tyArgs,
		Args:     [][]byte{bscAddress, bscAmount},
	}
	return &aptos_types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
}

func CoinBurnPayload(amount uint64, tyArg string) (aptos_types.TransactionPayload, error) {
	moduleAddress := make([]byte, 31)
	moduleAddress = append(moduleAddress, 0x1)

	bscAmount, err := aptos_types.BcsSerializeUint64(amount)

	parts := strings.Split(tyArg, "::")
	contractAddr, err := aptos_types.FromHex(ExpandAddress(parts[0]))
	if err != nil {
		return nil, err
	}
	tyArgs := make([]aptos_types.TypeTag, 0)
	t1 := aptos_types.TypeTagStruct{
		Value: aptos_types.StructTag{
			Address:    *contractAddr,
			ModuleName: aptos_types.Identifier(parts[1]),
			Name:       aptos_types.Identifier(parts[2]),
			TypeArgs:   []aptos_types.TypeTag{},
		},
	}
	tyArgs = append(tyArgs, &t1)

	scriptFunction := aptos_types.ScriptFunction{
		Module:   aptos_types.ModuleId{Address: *aptos_types.CORE_CODE_ADDRESS, Name: "managed_coin"},
		Function: "burn",
		TyArgs:   tyArgs,
		Args:     [][]byte{bscAmount},
	}
	return &aptos_types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
}

func OfferNFTTokenPayload(receiver string, creator string, collectionName string, tokenName string, propertyVersion uint64, amount uint64) (aptos_types.TransactionPayload, error) {
	moduleAddress, _ := aptos_types.FromHex("0x3")
	bscReceiver, err := aptos_types.BcsSerializeFixedBytes(aptos_types.BytesFromHex(ExpandAddress(receiver)))
	if err != nil {
		return nil, err
	}
	bscCreator, err := aptos_types.BcsSerializeFixedBytes(aptos_types.BytesFromHex(ExpandAddress(creator)))
	if err != nil {
		return nil, err
	}
	bscCollectName, err := aptos_types.BcsSerializeStr(collectionName)
	if err != nil {
		return nil, err
	}
	bscTokenName, err := aptos_types.BcsSerializeStr(tokenName)
	if err != nil {
		return nil, err
	}
	bscPropertyVersion, err := aptos_types.BcsSerializeUint64(propertyVersion)
	if err != nil {
		return nil, err
	}
	bscAmount, err := aptos_types.BcsSerializeUint64(amount)
	if err != nil {
		return nil, err
	}

	scriptFunction := aptos_types.ScriptFunction{
		Module:   aptos_types.ModuleId{Address: *moduleAddress, Name: "token_transfers"},
		Function: "offer_script",
		TyArgs:   []aptos_types.TypeTag{},
		Args:     [][]byte{bscReceiver, bscCreator, bscCollectName, bscTokenName, bscPropertyVersion, bscAmount},
	}
	return &aptos_types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
}

func ClaimNFTTokenPayload(sender string, creator string, collectionName string, tokenName string, propertyVersion uint64) (aptos_types.TransactionPayload, error) {
	moduleAddress, err := aptos_types.FromHex("0x3")
	if err != nil {
		return nil, err
	}
	bscSender, err := aptos_types.BcsSerializeFixedBytes(aptos_types.BytesFromHex(ExpandAddress(sender)))
	if err != nil {
		return nil, err
	}
	bscCreator, err := aptos_types.BcsSerializeFixedBytes(aptos_types.BytesFromHex(ExpandAddress(creator)))
	if err != nil {
		return nil, err
	}
	bscCollectName, err := aptos_types.BcsSerializeStr(collectionName)
	if err != nil {
		return nil, err
	}
	bscTokenName, err := aptos_types.BcsSerializeStr(tokenName)
	if err != nil {
		return nil, err
	}
	bscPropertyVersion, err := aptos_types.BcsSerializeUint64(propertyVersion)
	if err != nil {
		return nil, err
	}

	scriptFunction := aptos_types.ScriptFunction{
		Module:   aptos_types.ModuleId{Address: *moduleAddress, Name: "token_transfers"},
		Function: "claim_script",
		TyArgs:   []aptos_types.TypeTag{},
		Args:     [][]byte{bscSender, bscCreator, bscCollectName, bscTokenName, bscPropertyVersion},
	}
	return &aptos_types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
}

func SignRawTransaction(rawTxn *aptos_types.RawTransaction, seedHex string) (string, error) {
	publicKey, err := ed25519.PublicKeyFromSeed(seedHex)
	if err != nil {
		return "", err
	}

	message, err := rawTxn.GetSigningMessage()
	if err != nil {
		return "", err
	}
	signature, err := ed25519.Sign(seedHex, message)
	if err != nil {
		return "", err
	}
	ed25519Authenticator := aptos_types.TransactionAuthenticatorEd25519{PublicKey: aptos_types.Ed25519PublicKey(publicKey), Signature: signature}
	signedTransaction := aptos_types.SignedTransaction{RawTxn: *rawTxn, Authenticator: &ed25519Authenticator}
	txBytes, err := signedTransaction.BcsSerialize()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txBytes), nil
}

func SimulateTransaction(rawTxn *aptos_types.RawTransaction, seedHex string) (string, error) {
	publicKey, err := ed25519.PublicKeyFromSeed(seedHex)
	if err != nil {
		return "", err
	}

	signature := make([]byte, 64)
	ed25519Authenticator := aptos_types.TransactionAuthenticatorEd25519{PublicKey: aptos_types.Ed25519PublicKey(publicKey), Signature: signature}
	signedTransaction := aptos_types.SignedTransaction{RawTxn: *rawTxn, Authenticator: &ed25519Authenticator}
	txBytes, err := signedTransaction.BcsSerialize()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txBytes), nil
}

func parseTypeArguments(data string) *aptos_types.TypeTagStruct {
	i1 := strings.Index(data, "<")
	i2 := strings.Index(data, ">")
	var left string
	var right []string

	if i1 != -1 && i2 != -1 {
		left = data[0:i1]
		right = strings.Split(data[i1+1:i2], ",")
	} else {
		left = data[:]
		right = make([]string, 0)
	}

	typeTags := make([]aptos_types.TypeTag, 0)
	for _, s := range right {
		if len(s) > 0 {
			temp := parseTypeArguments(s)
			typeTags = append(typeTags, temp)
		}
	}

	parts := strings.Split(left, "::")
	// module address
	p1, _ := aptos_types.FromHex(ExpandAddress(parts[0]))
	// module name
	p2 := aptos_types.Identifier(parts[1])
	// struct name
	p3 := aptos_types.Identifier(parts[2])

	return &aptos_types.TypeTagStruct{
		Value: aptos_types.StructTag{
			Address:    *p1,
			ModuleName: p2,
			Name:       p3,
			TypeArgs:   typeTags,
		},
	}
}

func String2U128(str string) (*serde.Uint128, error) {
	ii := big.Int{}
	_, ret := ii.SetString(str, 0)
	if !ret {
		return nil, fmt.Errorf("unknown argument for u128")
	}
	return serde.FromBig(&ii)
}

func Interface2U64(value interface{}) (uint64, error) {
	switch value.(type) {
	case string:
		op, _ := value.(string)
		return strconv.ParseUint(op, 0, 64)
	case float64:
		op, _ := value.(float64)
		return uint64(op), nil
	default:
		return 0, fmt.Errorf("convert value to u64 fail")
	}
}

func Interface2U128(value interface{}) (*serde.Uint128, error) {
	switch value.(type) {
	case string:
		op, _ := value.(string)
		return String2U128(op)
	case float64:
		op, _ := value.(float64)
		u := serde.From64(uint64(op))
		return &u, nil
	default:
		return nil, fmt.Errorf("convert value to u128 fail")
	}
}

func ConvertArgs(args []interface{}, arg_types []aptos_types.MoveType) ([][]byte, error) {
	if len(args) != len(arg_types) {
		return nil, fmt.Errorf("aptos_types and values size not match")
	}
	array := make([][]byte, 0)
	for i := range args {
		moveType := arg_types[i]
		moveValue := args[i]
		switch moveType {
		case "address":
			op, ok := moveValue.(string)
			if !ok {
				return nil, fmt.Errorf("unknown argument for address")
			}
			bytes, err := aptos_types.BcsSerializeFixedBytes(aptos_types.BytesFromHex(ExpandAddress(op)))
			if err != nil {
				return nil, err
			}
			array = append(array, bytes)
		case "u64":
			ai, err := Interface2U64(moveValue)
			if err != nil {
				return nil, fmt.Errorf("unknown argument for u64, %w", err)
			}
			bytes, err := aptos_types.BcsSerializeUint64(ai)
			if err != nil {
				return nil, err
			}
			array = append(array, bytes)
		case "bool":
			op := fmt.Sprintf("%v", moveValue)
			b, err := strconv.ParseBool(op)
			if err != nil {
				return nil, fmt.Errorf("unknown argument for bool, %w", err)
			}
			bytes, err := aptos_types.BcsSerializeBool(b)
			if err != nil {
				return nil, err
			}
			array = append(array, bytes)
		case "u8":
			op := fmt.Sprintf("%v", moveValue)
			ai, err := strconv.ParseUint(op, 0, 8)
			if err != nil {
				return nil, fmt.Errorf("unknown argument for u8, %w", err)
			}
			bytes, err := aptos_types.BcsSerializeU8(uint8(ai))
			if err != nil {
				return nil, err
			}
			array = append(array, bytes)
		case "u128":
			ii, err := Interface2U128(moveValue)
			if err != nil {
				return nil, fmt.Errorf("unknown argument for u128, %w", err)
			}
			bytes, err := aptos_types.BcsSerializeU128(*ii)
			if err != nil {
				return nil, err
			}
			array = append(array, bytes)
		case "vector<u8>":
			switch moveValue.(type) {
			case []interface{}:
				vArray, _ := moveValue.([]interface{})
				inputBytes := make([]byte, 0)
				for _, e := range vArray {
					v := fmt.Sprintf("%v", e)
					vv, err := strconv.ParseUint(v, 0, 8)
					if err != nil {
						return nil, fmt.Errorf("unknown argument for u8, %w", err)
					}
					inputBytes = append(inputBytes, uint8(vv))
				}
				bytes, err := aptos_types.BcsSerializeBytes(inputBytes)
				if err != nil {
					return nil, err
				}
				array = append(array, bytes)
			case string:
				op, _ := moveValue.(string)
				v := aptos_types.BytesFromHex(op)
				bytes, err := aptos_types.BcsSerializeBytes(v)
				if err != nil {
					return nil, err
				}
				array = append(array, bytes)
			default:
				return nil, errors.New("unknown argument for vector<u8>")
			}
		case "vector<u64>":
			switch moveValue.(type) {
			case []interface{}:
				vArray, _ := moveValue.([]interface{})
				targetBytes := make([]byte, 0)
				bytes, err := aptos_types.BcsSerializeLen(uint64(len(vArray)))
				if err != nil {
					return nil, err
				}
				targetBytes = append(targetBytes, bytes...)
				for _, e := range vArray {
					v, err := Interface2U64(e)
					if err != nil {
						return nil, fmt.Errorf("unknown argument for u64, %w", err)
					}
					bytes, err = aptos_types.BcsSerializeUint64(v)
					if err != nil {
						return nil, err
					}
					targetBytes = append(targetBytes, bytes...)
				}
				array = append(array, targetBytes)
			default:
				return nil, errors.New("unknown argument for vector<u64>")
			}
		case "vector<u128>":
			switch moveValue.(type) {
			case []interface{}:
				vArray, _ := moveValue.([]interface{})
				targetBytes := make([]byte, 0)
				bytes, err := aptos_types.BcsSerializeLen(uint64(len(vArray)))
				if err != nil {
					return nil, err
				}
				targetBytes = append(targetBytes, bytes...)
				for _, e := range vArray {
					v, err := Interface2U128(e)
					if err != nil {
						return nil, fmt.Errorf("unknown argument for u128, %w", err)
					}
					bytes, err = aptos_types.BcsSerializeU128(*v)
					if err != nil {
						return nil, err
					}
					targetBytes = append(targetBytes, bytes...)
				}
				array = append(array, targetBytes)
			default:
				return nil, errors.New("unknown argument for vector<u128>")
			}
		case "vector<bool>":
			switch moveValue.(type) {
			case []interface{}:
				vArray, _ := moveValue.([]interface{})
				targetBytes := make([]byte, 0)
				bytes, err := aptos_types.BcsSerializeLen(uint64(len(vArray)))
				if err != nil {
					return nil, err
				}
				targetBytes = append(targetBytes, bytes...)
				for _, e := range vArray {
					v := fmt.Sprintf("%v", e)
					vv, err := strconv.ParseBool(v)
					if err != nil {
						return nil, fmt.Errorf("unknown argument for bool, %w", err)
					}
					bytes, err = aptos_types.BcsSerializeBool(vv)
					if err != nil {
						return nil, err
					}
					targetBytes = append(targetBytes, bytes...)
				}
				array = append(array, targetBytes)
			default:
				return nil, errors.New("unknown argument for vector<bool>")
			}
		case "vector<address>":
			switch moveValue.(type) {
			case []interface{}:
				vArray, _ := moveValue.([]interface{})
				targetBytes := make([]byte, 0)
				bytes, err := aptos_types.BcsSerializeLen(uint64(len(vArray)))
				if err != nil {
					return nil, err
				}
				targetBytes = append(targetBytes, bytes...)

				for _, e := range vArray {
					v, ok := e.(string)
					if !ok {
						return nil, errors.New("unknown argument for address")
					}
					bytes, err = aptos_types.BcsSerializeFixedBytes(aptos_types.BytesFromHex(ExpandAddress(v)))
					if err != nil {
						return nil, err
					}
					targetBytes = append(targetBytes, bytes...)
				}
				array = append(array, targetBytes)
			default:
				return nil, errors.New("unknown argument for vector<address>")
			}
		case "vector<0x1::string::String>":
			switch moveValue.(type) {
			case []interface{}:
				vArray, _ := moveValue.([]interface{})
				targetBytes := make([]byte, 0)
				bytes, err := aptos_types.BcsSerializeLen(uint64(len(vArray)))
				if err != nil {
					return nil, err
				}
				targetBytes = append(targetBytes, bytes...)

				for _, e := range vArray {
					v, ok := e.(string)
					if !ok {
						return nil, errors.New("unknown argument for string")
					}
					bytes, err = aptos_types.BcsSerializeStr(v)
					if err != nil {
						return nil, err
					}
					targetBytes = append(targetBytes, bytes...)
				}
				array = append(array, targetBytes)
			default:
				return nil, errors.New("unknown argument for vector<string>")
			}
		case "0x1::string::String":
			op, ok := moveValue.(string)
			if !ok {
				return nil, errors.New("unknown argument for string")
			}
			bytes, err := aptos_types.BcsSerializeStr(op)
			if err != nil {
				return nil, err
			}
			array = append(array, bytes)
		default:
			return nil, fmt.Errorf("unkonown type, %s", moveType)
		}
	}
	return array, nil
}
func fetchABI(modules []aptos_types.MoveModuleBytecode) map[string]aptos_types.MoveFunctionFullName {
	abiMap := map[string]aptos_types.MoveFunctionFullName{}
	for _, module := range modules {
		abi := module.Abi
		for _, ef := range abi.ExposedFunctions {
			if ef.IsEntry {
				fullName := abi.Address + "::" + abi.Name + "::" + ef.Name
				abiMap[fullName] = aptos_types.MoveFunctionFullName{
					FullName:     fullName,
					MoveFunction: ef,
				}
			}
		}
	}
	return abiMap
}

func filterMoveFunctionParams(funcAbi aptos_types.MoveFunctionFullName) []string {
	res := make([]string, 0)
	for _, param := range funcAbi.Params {
		if param == "signer" || param == "&signer" {
			continue
		}
		res = append(res, param)
	}
	return res
}

func PayloadFromJsonAndAbi(payload string, abi string) (aptos_types.TransactionPayload, error) {
	moveModules := make([]aptos_types.MoveModuleBytecode, 0)
	err := json.Unmarshal([]byte(abi), &moveModules)
	if err != nil {
		return nil, err
	}

	entryFunction := aptos_types.EntryFunctionPayload{}
	err = json.Unmarshal([]byte(payload), &entryFunction)
	if err != nil {
		return nil, err
	}

	f := entryFunction.Function
	r, err := regexp.Compile("^0[xX]0*")
	if err != nil {
		return nil, err
	}
	function := r.ReplaceAllString(f, "0x")
	funcParts := strings.Split(function, "::")
	if len(funcParts) != 3 {
		return nil, errors.New("func needs to be a fully qualified function name in format <address>::<module>::<function>, e.g. 0x1::coin::transfer")
	}
	abiMap := fetchABI(moveModules)
	funcAbi, ok := abiMap[function]
	if !ok {
		return nil, errors.New("abi miss")
	}
	abiArgs := filterMoveFunctionParams(funcAbi)
	typeArgsEntryFunction := entryFunction.TypeArguments
	typeArgABIs := make([]aptos_types.ArgumentABI, 0)
	for i, abiArg := range abiArgs {
		parser, err := aptos_types.NewTypeTagParser(abiArg, typeArgsEntryFunction)
		if err != nil {
			return nil, err
		}
		typeTag, err := parser.ParseTypeTag()
		if err != nil {
			return nil, err
		}
		argAbi := aptos_types.ArgumentABI{
			Name:    strconv.Itoa(i),
			TypeTag: typeTag,
		}
		typeArgABIs = append(typeArgABIs, argAbi)
	}

	// here only support EntryFunctionABI
	// todo support TransactionScriptABI
	// argument in input data
	typeTags := []aptos_types.TypeTag{}
	for _, tagString := range entryFunction.TypeArguments {
		if tagString == "" {
			continue
		}
		parser, err := aptos_types.NewTypeTagParser(tagString, nil)
		if err != nil {
			return nil, err
		}
		tag, err := parser.ParseTypeTag()
		if err != nil {
			return nil, err
		}
		typeTags = append(typeTags, tag)
	}
	args, err := transaction_builder.ToBCSArgs(typeArgABIs, entryFunction.Arguments)
	if err != nil {
		return nil, err
	}
	ma, err := aptos_types.FromHex(funcParts[0])
	if err != nil {
		return nil, err
	}
	mn := aptos_types.Identifier(funcParts[1])
	fn := aptos_types.Identifier(funcParts[2])
	scriptFunction := aptos_types.ScriptFunction{
		Module:   aptos_types.ModuleId{Address: *ma, Name: mn},
		Function: fn,
		TyArgs:   typeTags,
		Args:     args,
	}
	return &aptos_types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
}

func GetSigningHash(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64, expirationTimestampSecs uint64, chainId uint8,
	to string, amount uint64) (string, error) {
	payload, err := TransferPayload(to, amount)
	if err != nil {
		return "", err
	}
	rawTxn, err := MakeRawTransaction(from, sequenceNumber, maxGasAmount, gasUnitPrice, expirationTimestampSecs, chainId, payload)
	if err != nil {
		return "", err
	}
	rawTxHash, err := rawTxn.GetSigningMessage()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(rawTxHash), nil
}

func GetRawTxHash(rawTxn *aptos_types.RawTransaction) (string, error) {
	rawTxHash, err := rawTxn.GetSigningMessage()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(rawTxHash), nil
}

func SignedTx(rawTxn *aptos_types.RawTransaction, signDataHex string, pubKey string) (string, error) {
	pb, err := hex.DecodeString(pubKey)
	if err != nil {
		return "", err
	}

	signData, err := hex.DecodeString(signDataHex)
	if err != nil {
		return "", err
	}

	ed25519Authenticator := aptos_types.TransactionAuthenticatorEd25519{PublicKey: aptos_types.Ed25519PublicKey(pb), Signature: signData}
	signedTransaction := aptos_types.SignedTransaction{RawTxn: *rawTxn, Authenticator: &ed25519Authenticator}
	txBytes, err := signedTransaction.BcsSerialize()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txBytes), nil
}

func AddStake(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64, expirationTimestampSecs uint64, chainId uint8,
	poolAddress string, amount uint64, seedHex string) (string, error) {
	payload, err := AddStakePayload(poolAddress, amount)
	if err != nil {
		return "", err
	}
	return BuildSignedTransaction(from, sequenceNumber, maxGasAmount, gasUnitPrice, expirationTimestampSecs, chainId, payload, seedHex)
}

func AddStakePayload(poolAddress string, amount uint64) (aptos_types.TransactionPayload, error) {
	// moduleAddress := make([]byte, 31)
	// moduleAddress = append(moduleAddress, 0x1)
	bscAddress, err := aptos_types.BcsSerializeFixedBytes(aptos_types.BytesFromHex(ExpandAddress(poolAddress)))
	if err != nil {
		return nil, err
	}
	bscAmount, err := aptos_types.BcsSerializeUint64(amount)
	if err != nil {
		return nil, err
	}
	scriptFunction := aptos_types.ScriptFunction{
		Module:   aptos_types.ModuleId{Address: *aptos_types.CORE_CODE_ADDRESS, Name: "delegation_pool"},
		Function: "add_stake",
		TyArgs:   []aptos_types.TypeTag{},
		Args:     [][]byte{bscAddress, bscAmount},
	}
	return &aptos_types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
}

func Unlock(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64, expirationTimestampSecs uint64, chainId uint8,
	poolAddress string, amount uint64, seedHex string) (string, error) {
	payload, err := UnlockPayload(poolAddress, amount)
	if err != nil {
		return "", err
	}
	return BuildSignedTransaction(from, sequenceNumber, maxGasAmount, gasUnitPrice, expirationTimestampSecs, chainId, payload, seedHex)
}

func UnlockPayload(poolAddress string, amount uint64) (aptos_types.TransactionPayload, error) {
	// moduleAddress := make([]byte, 31)
	// moduleAddress = append(moduleAddress, 0x1)
	bscAddress, err := aptos_types.BcsSerializeFixedBytes(aptos_types.BytesFromHex(ExpandAddress(poolAddress)))
	if err != nil {
		return nil, err
	}
	bscAmount, err := aptos_types.BcsSerializeUint64(amount)
	if err != nil {
		return nil, err
	}
	scriptFunction := aptos_types.ScriptFunction{
		Module:   aptos_types.ModuleId{Address: *aptos_types.CORE_CODE_ADDRESS, Name: "delegation_pool"},
		Function: "unlock",
		TyArgs:   []aptos_types.TypeTag{},
		Args:     [][]byte{bscAddress, bscAmount},
	}
	return &aptos_types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
}

func ReactivateStake(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64, expirationTimestampSecs uint64, chainId uint8,
	poolAddress string, amount uint64, seedHex string) (string, error) {
	payload, err := ReactivateStakePayload(poolAddress, amount)
	if err != nil {
		return "", err
	}
	return BuildSignedTransaction(from, sequenceNumber, maxGasAmount, gasUnitPrice, expirationTimestampSecs, chainId, payload, seedHex)
}

func ReactivateStakePayload(poolAddress string, amount uint64) (aptos_types.TransactionPayload, error) {
	moduleAddress := make([]byte, 31)
	moduleAddress = append(moduleAddress, 0x1)
	bscAddress, err := aptos_types.BcsSerializeFixedBytes(aptos_types.BytesFromHex(ExpandAddress(poolAddress)))
	if err != nil {
		return nil, err
	}
	bscAmount, err := aptos_types.BcsSerializeUint64(amount)
	if err != nil {
		return nil, err
	}
	scriptFunction := aptos_types.ScriptFunction{
		Module:   aptos_types.ModuleId{Address: *aptos_types.CORE_CODE_ADDRESS, Name: "delegation_pool"},
		Function: "reactivate_stake",
		TyArgs:   []aptos_types.TypeTag{},
		Args:     [][]byte{bscAddress, bscAmount},
	}
	return &aptos_types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
}

func Withdraw(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64, expirationTimestampSecs uint64, chainId uint8,
	poolAddress string, amount uint64, seedHex string) (string, error) {
	payload, err := WithdrawPayload(poolAddress, amount)
	if err != nil {
		return "", err
	}
	return BuildSignedTransaction(from, sequenceNumber, maxGasAmount, gasUnitPrice, expirationTimestampSecs, chainId, payload, seedHex)
}

func WithdrawPayload(poolAddress string, amount uint64) (aptos_types.TransactionPayload, error) {
	moduleAddress := make([]byte, 31)
	moduleAddress = append(moduleAddress, 0x1)
	bscAddress, err := aptos_types.BcsSerializeFixedBytes(aptos_types.BytesFromHex(ExpandAddress(poolAddress)))
	if err != nil {
		return nil, err
	}
	bscAmount, err := aptos_types.BcsSerializeUint64(amount)
	if err != nil {
		return nil, err
	}
	scriptFunction := aptos_types.ScriptFunction{
		Module:   aptos_types.ModuleId{Address: *aptos_types.CORE_CODE_ADDRESS, Name: "delegation_pool"},
		Function: "withdraw",
		TyArgs:   []aptos_types.TypeTag{},
		Args:     [][]byte{bscAddress, bscAmount},
	}
	return &aptos_types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
}

func SignMessage(priKey, message string) (string, error) {
	if len(priKey) == 0 || len(message) == 0 {
		return "", fmt.Errorf("invalid params message %s", message)
	}
	signature, err := ed25519.Sign(priKey, []byte(message))
	if err != nil {
		return "", err
	}
	return "0x" + hex.EncodeToString(signature), nil
}
