package aptos

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/okx/go-wallet-sdk/coins/aptos/serde"
	"github.com/okx/go-wallet-sdk/coins/aptos/types"
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
	address := "0x" + hex.EncodeToString(types.Sha256Hash(publicKey))
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
	address := "0x" + hex.EncodeToString(types.Sha256Hash(pubKey))
	if shortEnable {
		return ShortenAddress(address), nil
	} else {
		return address, nil
	}
}

// ValidateAddress hex 32bytes
func ValidateAddress(address string, shortEnable bool) bool {
	re1, _ := regexp.Compile("^0x[\\dA-Fa-f]{62,64}$")
	re2, _ := regexp.Compile("^[\\dA-Fa-f]{64}$")
	return re1.Match([]byte(address)) || re2.Match([]byte(address))
}

func MakeRawTransaction(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64,
	expirationTimestampSecs uint64, chainId uint8, payload types.TransactionPayload) *types.RawTransaction {
	rawTxn := types.RawTransaction{}
	rawTxn.Sender = types.BytesFromHex(ExpandAddress(from))
	rawTxn.SequenceNumber = sequenceNumber
	rawTxn.MaxGasAmount = maxGasAmount
	rawTxn.GasUnitPrice = gasUnitPrice
	rawTxn.ExpirationTimestampSecs = expirationTimestampSecs
	rawTxn.ChainId = types.ChainId(chainId)
	rawTxn.Payload = payload
	return &rawTxn
}

func Transfer(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64, expirationTimestampSecs uint64, chainId uint8,
	to string, amount uint64, seedHex string) (string, error) {
	payload, err := TransferPayload(to, amount)
	if err != nil {
		return "", err
	}
	return BuildSignedTransaction(from, sequenceNumber, maxGasAmount, gasUnitPrice, expirationTimestampSecs, chainId, payload, seedHex)
}

func TransferPayload(to string, amount uint64) (types.TransactionPayload, error) {
	moduleAddress := make([]byte, 31)
	moduleAddress = append(moduleAddress, 0x1)
	bscAddress, err := types.BcsSerializeFixedBytes(types.BytesFromHex(ExpandAddress(to)))
	if err != nil {
		return nil, err
	}
	bscAmount, err := types.BcsSerializeUint64(amount)
	if err != nil {
		return nil, err
	}
	scriptFunction := types.ScriptFunction{
		Module:   types.ModuleId{Address: moduleAddress, Name: "aptos_account"},
		Function: "transfer",
		TyArgs:   []types.TypeTag{},
		Args:     [][]byte{bscAddress, bscAmount},
	}
	return &types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
}

func BuildSignedTransaction(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64, expirationTimestampSecs uint64, chainId uint8,
	payload types.TransactionPayload, seedHex string) (string, error) {
	rawTxn := MakeRawTransaction(from, sequenceNumber, maxGasAmount, gasUnitPrice, expirationTimestampSecs, chainId, payload)
	return SignRawTransaction(rawTxn, seedHex)
}

// pub enum Transaction {
// /// Transaction submitted by the user. e.g: P2P payment transaction, publishing module
// /// transaction, etc.
// /// TODO: We need to rename SignedTransaction to SignedUserTransaction, as well as all the other
// ///       transaction types we had in our codebase.
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
	prefix := types.Sha256Hash([]byte("APTOS::Transaction"))
	bcsBytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return "", err
	}
	message := make([]byte, 0)
	message = append(message, prefix...)
	message = append(message, 0x0)
	message = append(message, bcsBytes...)
	return "0x" + hex.EncodeToString(types.Sha256Hash(message)), nil
}

func CoinTransferPayload(to string, amount uint64, tyArg string) (types.TransactionPayload, error) {
	moduleAddress := make([]byte, 31)
	moduleAddress = append(moduleAddress, 0x1)
	bscAddress, err := types.BcsSerializeFixedBytes(types.BytesFromHex(ExpandAddress(to)))
	if err != nil {
		return nil, err
	}
	bscAmount, err := types.BcsSerializeUint64(amount)
	if err != nil {
		return nil, err
	}

	// 0x3::moon_coin::MoonCoin
	parts := strings.Split(tyArg, "::")
	contractAddr := types.BytesFromHex(ExpandAddress(parts[0]))
	tyArgs := make([]types.TypeTag, 0)
	t1 := types.TypeTag__Struct{
		Value: types.StructTag{
			Address:    contractAddr,
			Module:     types.Identifier(parts[1]),
			Name:       types.Identifier(parts[2]),
			TypeParams: []types.TypeTag{},
		},
	}
	tyArgs = append(tyArgs, &t1)

	// 0x1::coin transfer
	scriptFunction := types.ScriptFunction{
		Module:   types.ModuleId{Address: moduleAddress, Name: "coin"},
		Function: "transfer",
		TyArgs:   tyArgs,
		Args:     [][]byte{bscAddress, bscAmount},
	}
	return &types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
}

func CoinRegisterPayload(tyArg string) types.TransactionPayload {
	moduleAddress := make([]byte, 31)
	moduleAddress = append(moduleAddress, 0x1)

	parts := strings.Split(tyArg, "::")
	contractAddr := types.BytesFromHex(ExpandAddress(parts[0]))
	tyArgs := make([]types.TypeTag, 0)
	t1 := types.TypeTag__Struct{
		Value: types.StructTag{
			Address:    contractAddr,
			Module:     types.Identifier(parts[1]),
			Name:       types.Identifier(parts[2]),
			TypeParams: []types.TypeTag{},
		},
	}
	tyArgs = append(tyArgs, &t1)

	scriptFunction := types.ScriptFunction{
		Module:   types.ModuleId{Address: moduleAddress, Name: "managed_coin"},
		Function: "register",
		TyArgs:   tyArgs,
		Args:     [][]byte{},
	}
	return &types.TransactionPayloadEntryFunction{Value: scriptFunction}
}

func CoinMintPayload(receiveAddress string, amount uint64, tyArg string) (types.TransactionPayload, error) {
	moduleAddress := make([]byte, 31)
	moduleAddress = append(moduleAddress, 0x1)

	bscAddress, err := types.BcsSerializeFixedBytes(types.BytesFromHex(receiveAddress))
	if err != nil {
		return nil, err
	}
	bscAmount, err := types.BcsSerializeUint64(amount)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(tyArg, "::")
	contractAddr := types.BytesFromHex(ExpandAddress(parts[0]))
	tyArgs := make([]types.TypeTag, 0)
	t1 := types.TypeTag__Struct{
		Value: types.StructTag{
			Address:    contractAddr,
			Module:     types.Identifier(parts[1]),
			Name:       types.Identifier(parts[2]),
			TypeParams: []types.TypeTag{},
		},
	}
	tyArgs = append(tyArgs, &t1)

	scriptFunction := types.ScriptFunction{
		Module:   types.ModuleId{Address: moduleAddress, Name: "managed_coin"},
		Function: "mint",
		TyArgs:   tyArgs,
		Args:     [][]byte{bscAddress, bscAmount},
	}
	return &types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
}

func CoinBurnPayload(amount uint64, tyArg string) (types.TransactionPayload, error) {
	moduleAddress := make([]byte, 31)
	moduleAddress = append(moduleAddress, 0x1)

	bscAmount, err := types.BcsSerializeUint64(amount)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(tyArg, "::")
	contractAddr := types.BytesFromHex(ExpandAddress(parts[0]))
	tyArgs := make([]types.TypeTag, 0)
	t1 := types.TypeTag__Struct{
		Value: types.StructTag{
			Address:    contractAddr,
			Module:     types.Identifier(parts[1]),
			Name:       types.Identifier(parts[2]),
			TypeParams: []types.TypeTag{},
		},
	}
	tyArgs = append(tyArgs, &t1)

	scriptFunction := types.ScriptFunction{
		Module:   types.ModuleId{Address: moduleAddress, Name: "managed_coin"},
		Function: "burn",
		TyArgs:   tyArgs,
		Args:     [][]byte{bscAmount},
	}
	return &types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
}

func OfferNFTTokenPayload(receiver string, creator string, collectionName string, tokenName string, propertyVersion uint64, amount uint64) (types.TransactionPayload, error) {
	moduleAddress := make([]byte, 31)
	moduleAddress = append(moduleAddress, 0x3)

	bscReceiver, err := types.BcsSerializeFixedBytes(types.BytesFromHex(ExpandAddress(receiver)))
	if err != nil {
		return nil, err
	}
	bscCreator, err := types.BcsSerializeFixedBytes(types.BytesFromHex(ExpandAddress(creator)))
	if err != nil {
		return nil, err
	}
	bscCollectName, err := types.BcsSerializeStr(collectionName)
	if err != nil {
		return nil, err
	}
	bscTokenName, err := types.BcsSerializeStr(tokenName)
	if err != nil {
		return nil, err
	}
	bscPropertyVersion, err := types.BcsSerializeUint64(propertyVersion)
	if err != nil {
		return nil, err
	}
	bscAmount, err := types.BcsSerializeUint64(amount)
	if err != nil {
		return nil, err
	}

	scriptFunction := types.ScriptFunction{
		Module:   types.ModuleId{Address: moduleAddress, Name: "token_transfers"},
		Function: "offer_script",
		TyArgs:   []types.TypeTag{},
		Args:     [][]byte{bscReceiver, bscCreator, bscCollectName, bscTokenName, bscPropertyVersion, bscAmount},
	}
	return &types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
}

func ClaimNFTTokenPayload(sender string, creator string, collectionName string, tokenName string, propertyVersion uint64) (types.TransactionPayload, error) {
	moduleAddress := make([]byte, 31)
	moduleAddress = append(moduleAddress, 0x3)

	bscSender, err := types.BcsSerializeFixedBytes(types.BytesFromHex(ExpandAddress(sender)))
	if err != nil {
		return nil, err
	}
	bscCreator, err := types.BcsSerializeFixedBytes(types.BytesFromHex(ExpandAddress(creator)))
	if err != nil {
		return nil, err
	}
	bscCollectName, err := types.BcsSerializeStr(collectionName)
	if err != nil {
		return nil, err
	}
	bscTokenName, err := types.BcsSerializeStr(tokenName)
	if err != nil {
		return nil, err
	}
	bscPropertyVersion, err := types.BcsSerializeUint64(propertyVersion)
	if err != nil {
		return nil, err
	}

	scriptFunction := types.ScriptFunction{
		Module:   types.ModuleId{Address: moduleAddress, Name: "token_transfers"},
		Function: "claim_script",
		TyArgs:   []types.TypeTag{},
		Args:     [][]byte{bscSender, bscCreator, bscCollectName, bscTokenName, bscPropertyVersion},
	}
	return &types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
}

func SignRawTransaction(rawTxn *types.RawTransaction, seedHex string) (string, error) {
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

	ed25519Authenticator := types.TransactionAuthenticator__Ed25519{PublicKey: types.Ed25519PublicKey(publicKey), Signature: signature}
	signedTransaction := types.SignedTransaction{RawTxn: *rawTxn, Authenticator: &ed25519Authenticator}
	txBytes, err := signedTransaction.BcsSerialize()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txBytes), nil
}

func SimulateTransaction(rawTxn *types.RawTransaction, seedHex string) (string, error) {
	publicKey, err := ed25519.PublicKeyFromSeed(seedHex)
	if err != nil {
		return "", err
	}
	signature := make([]byte, 64)
	ed25519Authenticator := types.TransactionAuthenticator__Ed25519{PublicKey: types.Ed25519PublicKey(publicKey), Signature: signature}
	signedTransaction := types.SignedTransaction{RawTxn: *rawTxn, Authenticator: &ed25519Authenticator}
	txBytes, err := signedTransaction.BcsSerialize()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txBytes), nil
}

func parseTypeArguments(data string) *types.TypeTag__Struct {
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

	typeTags := make([]types.TypeTag, 0)
	for _, s := range right {
		if len(s) > 0 {
			temp := parseTypeArguments(s)
			typeTags = append(typeTags, temp)
		}
	}

	parts := strings.Split(left, "::")
	// module address
	p1 := types.BytesFromHex(ExpandAddress(parts[0]))
	// module name
	p2 := types.Identifier(parts[1])
	// struct name
	p3 := types.Identifier(parts[2])

	return &types.TypeTag__Struct{
		Value: types.StructTag{
			Address:    p1,
			Module:     p2,
			Name:       p3,
			TypeParams: typeTags,
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

func ConvertArgs(args []interface{}, argTypes []MoveType) ([][]byte, error) {
	if len(args) != len(argTypes) {
		return nil, fmt.Errorf("types and values size not match")
	}
	array := make([][]byte, 0)
	for i := range args {
		moveType := argTypes[i]
		moveValue := args[i]
		switch moveType {
		case "address":
			op, ok := moveValue.(string)
			if !ok {
				return nil, fmt.Errorf("unknown argument for address")
			}
			bytes, err := types.BcsSerializeFixedBytes(types.BytesFromHex(ExpandAddress(op)))
			if err != nil {
				return nil, err
			}
			array = append(array, bytes)
		case "u64":
			ai, err := Interface2U64(moveValue)
			if err != nil {
				return nil, fmt.Errorf("unknown argument for u64")
			}
			bytes, err := types.BcsSerializeUint64(ai)
			if err != nil {
				return nil, err
			}
			array = append(array, bytes)
		case "bool":
			op := fmt.Sprintf("%v", moveValue)
			b, err := strconv.ParseBool(op)
			if err != nil {
				return nil, fmt.Errorf("unknown argument for bool")
			}
			bytes, err := types.BcsSerializeBool(b)
			if err != nil {
				return nil, err
			}
			array = append(array, bytes)
		case "u8":
			op := fmt.Sprintf("%v", moveValue)
			ai, err := strconv.ParseUint(op, 0, 8)
			if err != nil {
				return nil, fmt.Errorf("unknown argument for u8")
			}
			bytes, err := types.BcsSerializeU8(uint8(ai))
			if err != nil {
				return nil, err
			}
			array = append(array, bytes)
		case "u128":
			ii, err := Interface2U128(moveValue)
			if err != nil {
				return nil, fmt.Errorf("unknown argument for u128")
			}
			bytes, err := types.BcsSerializeU128(*ii)
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
						return nil, fmt.Errorf("unknown argument for u8")
					}
					inputBytes = append(inputBytes, uint8(vv))
				}
				bytes, err := types.BcsSerializeBytes(inputBytes)
				if err != nil {
					return nil, err
				}
				array = append(array, bytes)
			case string:
				op, _ := moveValue.(string)
				v := types.BytesFromHex(op)
				bytes, err := types.BcsSerializeBytes(v)
				if err != nil {
					return nil, err
				}
				array = append(array, bytes)
			default:
				return nil, fmt.Errorf("unknown argument for vector<u8>")
			}
		case "vector<u64>":
			switch moveValue.(type) {
			case []interface{}:
				vArray, _ := moveValue.([]interface{})
				targetBytes := make([]byte, 0)
				bytes, err := types.BcsSerializeLen(uint64(len(vArray)))
				if err != nil {
					return nil, err
				}
				targetBytes = append(targetBytes, bytes...)
				for _, e := range vArray {
					v, err := Interface2U64(e)
					if err != nil {
						return nil, fmt.Errorf("unknown argument for u64")
					}
					bytes, err = types.BcsSerializeUint64(v)
					if err != nil {
						return nil, err
					}
					targetBytes = append(targetBytes, bytes...)
				}
				array = append(array, targetBytes)
			default:
				return nil, fmt.Errorf("unknown argument for vector<u64>")
			}
		case "vector<u128>":
			switch moveValue.(type) {
			case []interface{}:
				vArray, _ := moveValue.([]interface{})
				targetBytes := make([]byte, 0)
				bytes, err := types.BcsSerializeLen(uint64(len(vArray)))
				if err != nil {
					return nil, err
				}
				targetBytes = append(targetBytes, bytes...)
				for _, e := range vArray {
					v, err := Interface2U128(e)
					if err != nil {
						return nil, fmt.Errorf("unknown argument for u128")
					}
					bytes, err = types.BcsSerializeU128(*v)
					if err != nil {
						return nil, err
					}
					targetBytes = append(targetBytes, bytes...)
				}
				array = append(array, targetBytes)
			default:
				return nil, fmt.Errorf("unknown argument for vector<u128>")
			}
		case "vector<bool>":
			switch moveValue.(type) {
			case []interface{}:
				vArray, _ := moveValue.([]interface{})
				targetBytes := make([]byte, 0)
				bytes, err := types.BcsSerializeLen(uint64(len(vArray)))
				if err != nil {
					return nil, err
				}
				targetBytes = append(targetBytes, bytes...)
				for _, e := range vArray {
					v := fmt.Sprintf("%v", e)
					vv, err := strconv.ParseBool(v)
					if err != nil {
						return nil, fmt.Errorf("unknown argument for bool")
					}
					bytes, err = types.BcsSerializeBool(vv)
					if err != nil {
						return nil, err
					}
					targetBytes = append(targetBytes, bytes...)
				}
				array = append(array, targetBytes)
			default:
				return nil, fmt.Errorf("unknown argument for vector<bool>")
			}
		case "vector<address>":
			switch moveValue.(type) {
			case []interface{}:
				vArray, _ := moveValue.([]interface{})
				targetBytes := make([]byte, 0)
				bytes, err := types.BcsSerializeLen(uint64(len(vArray)))
				if err != nil {
					return nil, err
				}
				targetBytes = append(targetBytes, bytes...)

				for _, e := range vArray {
					v, ok := e.(string)
					if !ok {
						return nil, fmt.Errorf("unknown argument for address")
					}
					bytes, err = types.BcsSerializeFixedBytes(types.BytesFromHex(ExpandAddress(v)))
					if err != nil {
						return nil, err
					}
					targetBytes = append(targetBytes, bytes...)
				}
				array = append(array, targetBytes)
			default:
				return nil, fmt.Errorf("unknown argument for vector<address>")
			}
		case "vector<0x1::string::String>":
			switch moveValue.(type) {
			case []interface{}:
				vArray, _ := moveValue.([]interface{})
				targetBytes := make([]byte, 0)
				bytes, err := types.BcsSerializeLen(uint64(len(vArray)))
				if err != nil {
					return nil, err
				}
				targetBytes = append(targetBytes, bytes...)

				for _, e := range vArray {
					v, ok := e.(string)
					if !ok {
						return nil, fmt.Errorf("unknown argument for string")
					}
					bytes, err = types.BcsSerializeStr(v)
					if err != nil {
						return nil, err
					}
					targetBytes = append(targetBytes, bytes...)
				}
				array = append(array, targetBytes)
			default:
				return nil, fmt.Errorf("unknown argument for vector<string>")
			}
		case "0x1::string::String":
			op, ok := moveValue.(string)
			if !ok {
				return nil, fmt.Errorf("unknown argument for string")
			}
			bytes, err := types.BcsSerializeStr(op)
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

func PayloadFromJsonAndAbi(payload string, abi string) (types.TransactionPayload, error) {
	moveModules := make([]MoveModuleBytecode, 0)
	err := json.Unmarshal([]byte(abi), &moveModules)
	if err != nil {
		return nil, err
	}

	entryFunction := EntryFunctionPayload{}
	err = json.Unmarshal([]byte(payload), &entryFunction)
	if err != nil {
		return nil, err
	}

	funcParts := strings.Split(entryFunction.Function, "::")
	// 0x43417434fd869edee76cca2a4d2301e528a1551b1d719b75c350c3c97d15b8b9::scripts::swap
	for _, m := range moveModules {
		moveModuleAddress := ExpandAddress(funcParts[0])
		if ExpandAddress(m.Abi.Address) == moveModuleAddress && m.Abi.Name == funcParts[1] {
			for _, e := range m.Abi.ExposedFunctions {
				if e.IsEntry && e.Name == funcParts[2] {
					ma := types.BytesFromHex(moveModuleAddress)
					mn := types.Identifier(funcParts[1])
					fn := types.Identifier(funcParts[2])

					tyArgs := make([]types.TypeTag, 0)
					for _, ta := range entryFunction.TypeArguments {
						if len(ta) > 0 {
							tt := parseTypeArguments(ta)
							tyArgs = append(tyArgs, tt)
						}
					}

					args, err := ConvertArgs(entryFunction.Arguments, filterArgumentsTypes(e.Params))
					if err != nil {
						return nil, err
					}
					scriptFunction := types.ScriptFunction{
						Module:   types.ModuleId{Address: ma, Name: mn},
						Function: fn,
						TyArgs:   tyArgs,
						Args:     args,
					}
					return &types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("can not find move function from abi, %s", entryFunction.Function)
}

func filterArgumentsTypes(types []MoveType) []MoveType {
	array := make([]MoveType, 0)
	for _, t := range types {
		if !strings.Contains(t, "signer") {
			array = append(array, t)
		}
	}
	return array
}

func GetSigningHash(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64, expirationTimestampSecs uint64, chainId uint8,
	to string, amount uint64) (string, error) {

	payload, err := TransferPayload(to, amount)
	if err != nil {
		return "", err
	}
	rawTxn := MakeRawTransaction(from, sequenceNumber, maxGasAmount, gasUnitPrice, expirationTimestampSecs, chainId, payload)

	rawTxHash, err := rawTxn.GetSigningMessage()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(rawTxHash), err
}

func GetRawTxHash(rawTxn *types.RawTransaction) (string, error) {
	rawTxHash, err := rawTxn.GetSigningMessage()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(rawTxHash), err
}

func SignedTx(rawTxn *types.RawTransaction, signDataHex string, pubKey string) (string, error) {
	pb, err := hex.DecodeString(pubKey)
	if err != nil {
		return "", err
	}

	signData, err := hex.DecodeString(signDataHex)
	if err != nil {
		return "", err
	}

	ed25519Authenticator := types.TransactionAuthenticator__Ed25519{PublicKey: types.Ed25519PublicKey(pb), Signature: signData}
	signedTransaction := types.SignedTransaction{RawTxn: *rawTxn, Authenticator: &ed25519Authenticator}
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

func AddStakePayload(poolAddress string, amount uint64) (types.TransactionPayload, error) {
	moduleAddress := make([]byte, 31)
	moduleAddress = append(moduleAddress, 0x1)
	bscAddress, err := types.BcsSerializeFixedBytes(types.BytesFromHex(ExpandAddress(poolAddress)))
	if err != nil {
		return nil, err
	}
	bscAmount, err := types.BcsSerializeUint64(amount)
	if err != nil {
		return nil, err
	}
	scriptFunction := types.ScriptFunction{
		Module:   types.ModuleId{Address: moduleAddress, Name: "delegation_pool"},
		Function: "add_stake",
		TyArgs:   []types.TypeTag{},
		Args:     [][]byte{bscAddress, bscAmount},
	}
	return &types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
}

func Unlock(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64, expirationTimestampSecs uint64, chainId uint8,
	poolAddress string, amount uint64, seedHex string) (string, error) {
	payload, err := UnlockPayload(poolAddress, amount)
	if err != nil {
		return "", err
	}
	return BuildSignedTransaction(from, sequenceNumber, maxGasAmount, gasUnitPrice, expirationTimestampSecs, chainId, payload, seedHex)
}

func UnlockPayload(poolAddress string, amount uint64) (types.TransactionPayload, error) {
	moduleAddress := make([]byte, 31)
	moduleAddress = append(moduleAddress, 0x1)
	bscAddress, err := types.BcsSerializeFixedBytes(types.BytesFromHex(ExpandAddress(poolAddress)))
	if err != nil {
		return nil, err
	}
	bscAmount, err := types.BcsSerializeUint64(amount)
	if err != nil {
		return nil, err
	}
	scriptFunction := types.ScriptFunction{
		Module:   types.ModuleId{Address: moduleAddress, Name: "delegation_pool"},
		Function: "unlock",
		TyArgs:   []types.TypeTag{},
		Args:     [][]byte{bscAddress, bscAmount},
	}
	return &types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
}

func ReactivateStake(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64, expirationTimestampSecs uint64, chainId uint8,
	poolAddress string, amount uint64, seedHex string) (string, error) {
	payload, err := ReactivateStakePayload(poolAddress, amount)
	if err != nil {
		return "", err
	}
	return BuildSignedTransaction(from, sequenceNumber, maxGasAmount, gasUnitPrice, expirationTimestampSecs, chainId, payload, seedHex)
}

func ReactivateStakePayload(poolAddress string, amount uint64) (types.TransactionPayload, error) {
	moduleAddress := make([]byte, 31)
	moduleAddress = append(moduleAddress, 0x1)
	bscAddress, err := types.BcsSerializeFixedBytes(types.BytesFromHex(ExpandAddress(poolAddress)))
	if err != nil {
		return nil, err
	}
	bscAmount, err := types.BcsSerializeUint64(amount)
	if err != nil {
		return nil, err
	}
	scriptFunction := types.ScriptFunction{
		Module:   types.ModuleId{Address: moduleAddress, Name: "delegation_pool"},
		Function: "reactivate_stake",
		TyArgs:   []types.TypeTag{},
		Args:     [][]byte{bscAddress, bscAmount},
	}
	return &types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
}

func Withdraw(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64, expirationTimestampSecs uint64, chainId uint8,
	poolAddress string, amount uint64, seedHex string) (string, error) {
	payload, err := WithdrawPayload(poolAddress, amount)
	if err != nil {
		return "", err
	}
	return BuildSignedTransaction(from, sequenceNumber, maxGasAmount, gasUnitPrice, expirationTimestampSecs, chainId, payload, seedHex)
}

func WithdrawPayload(poolAddress string, amount uint64) (types.TransactionPayload, error) {
	moduleAddress := make([]byte, 31)
	moduleAddress = append(moduleAddress, 0x1)
	bscAddress, err := types.BcsSerializeFixedBytes(types.BytesFromHex(ExpandAddress(poolAddress)))
	if err != nil {
		return nil, err
	}
	bscAmount, err := types.BcsSerializeUint64(amount)
	if err != nil {
		return nil, err
	}
	scriptFunction := types.ScriptFunction{
		Module:   types.ModuleId{Address: moduleAddress, Name: "delegation_pool"},
		Function: "withdraw",
		TyArgs:   []types.TypeTag{},
		Args:     [][]byte{bscAddress, bscAmount},
	}
	return &types.TransactionPayloadEntryFunction{Value: scriptFunction}, nil
}
