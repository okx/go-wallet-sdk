package aptos

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	v2 "github.com/okx/go-wallet-sdk/coins/aptos/v2"
	bcs2 "github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/crypto"
	"github.com/okx/go-wallet-sdk/util"
	"math/big"
	"strconv"
	"strings"
)

// Script A Move script as compiled code as a transaction
type ScriptParam struct {
	Code     string   `json:"code"`     // The compiled script bytes, hex
	ArgTypes []string `json:"argTypes"` // The types of the arguments
	Args     []string `json:"args"`     // The arguments
}

// ScriptArgument a Move script argument, which encodes its type with it
type ScriptArgument struct {
	Variant v2.ScriptArgumentVariant `json:"variant"` // The type of the argument
	Value   string                   `json:"value"`   // The value of the argument
}

func (sa *ScriptArgument) parseScriptArgument() (v2.ScriptArgument, error) {
	switch sa.Variant {
	case v2.ScriptArgumentU8:
		i, err := strconv.ParseUint(sa.Value, 10, 8)
		if err != nil {
			return v2.ScriptArgument{}, err
		}
		return v2.ScriptArgument{
			Variant: sa.Variant,
			Value:   i,
		}, nil
	case v2.ScriptArgumentU64:
		i, err := strconv.ParseUint(sa.Value, 10, 64)
		if err != nil {
			return v2.ScriptArgument{}, err
		}
		return v2.ScriptArgument{
			Variant: sa.Variant,
			Value:   i,
		}, nil
	case v2.ScriptArgumentU128:
		return v2.ScriptArgument{
			Variant: sa.Variant,
			Value:   ConvertToBigInt(sa.Value),
		}, nil
	case v2.ScriptArgumentAddress:
		addr := &v2.AccountAddress{}
		if err := addr.ParseStringRelaxed(sa.Value); err != nil {
			return v2.ScriptArgument{}, err
		}
		return v2.ScriptArgument{
			Variant: sa.Variant,
			Value:   *addr,
		}, nil
	case v2.ScriptArgumentU8Vector:
		if strings.HasPrefix(sa.Value, "0x") {
			sa.Value = strings.TrimPrefix(sa.Value, "0x")
		}
		bs, err := hex.DecodeString(sa.Value)
		if err != nil {
			return v2.ScriptArgument{}, err
		}
		return v2.ScriptArgument{
			Variant: sa.Variant,
			Value:   bs,
		}, nil
	case v2.ScriptArgumentBool:
		if strings.ToLower(sa.Value) == "true" {
			return v2.ScriptArgument{
				Variant: sa.Variant,
				Value:   true,
			}, nil
		} else {
			return v2.ScriptArgument{
				Variant: sa.Variant,
				Value:   false,
			}, nil
		}
	case v2.ScriptArgumentU16:
		i, err := strconv.ParseUint(sa.Value, 10, 16)
		if err != nil {
			return v2.ScriptArgument{}, err
		}
		return v2.ScriptArgument{
			Variant: sa.Variant,
			Value:   i,
		}, nil
	case v2.ScriptArgumentU32:
		i, err := strconv.ParseUint(sa.Value, 10, 32)
		if err != nil {
			return v2.ScriptArgument{}, err
		}
		return v2.ScriptArgument{
			Variant: sa.Variant,
			Value:   i,
		}, nil
	case v2.ScriptArgumentU256:
		return v2.ScriptArgument{
			Variant: sa.Variant,
			Value:   ConvertToBigInt(sa.Value),
		}, nil
	case v2.ScriptArgumentSerialized:
		if strings.HasPrefix(sa.Value, "0x") {
			sa.Value = strings.TrimPrefix(sa.Value, "0x")
		}
		bs, err := hex.DecodeString(sa.Value)
		if err != nil {
			return v2.ScriptArgument{}, err
		}
		return v2.ScriptArgument{
			Variant: sa.Variant,
			Value:   bs,
		}, nil
	default:
		return v2.ScriptArgument{}, errors.New("unsupported script variant")
	}
}

func parseScriptParam(param *ScriptParam) (*v2.Script, error) {
	code := util.RemoveZeroHex(param.Code)
	argTypes := make([]v2.TypeTag, 0)
	for i := 0; i < len(param.ArgTypes); i++ {
		argType, err := v2.ParseTypeTag(param.ArgTypes[i])
		if err != nil {
			return nil, err
		}
		argTypes = append(argTypes, *argType)
	}
	args := make([]v2.ScriptArgument, 0)
	for i := 0; i < len(param.Args); i++ {
		var arg v2.ScriptArgument
		data, err := util.DecodeHexStringErr(param.Args[i])
		if err != nil {
			return nil, err
		}
		arg.UnmarshalBCS(bcs2.NewDeserializer(data))
		args = append(args, arg)
	}
	return &v2.Script{
		Code:     code,
		ArgTypes: argTypes,
		Args:     args,
	}, nil
}

func parseTypeTag(tt string) (v2.TypeTagImpl, error) {
	switch tt {
	case "address":
		return &v2.AddressTag{}, nil
	case "signer":
		return &v2.SignerTag{}, nil
	case "bool":
		return &v2.BoolTag{}, nil
	case "u8":
		return &v2.U8Tag{}, nil
	case "u16":
		return &v2.U16Tag{}, nil
	case "u32":
		return &v2.U32Tag{}, nil
	case "u64":
		return &v2.U64Tag{}, nil
	case "u128":
		return &v2.U128Tag{}, nil
	case "u256":
		return &v2.U256Tag{}, nil
	default:
		return nil, fmt.Errorf("unknown type: %s", tt)
	}
}

func parseAccountAddress(addr string) (*v2.AccountAddress, error) {
	acc := &v2.AccountAddress{}
	if err := acc.ParseStringRelaxed(addr); err != nil {
		return nil, err
	}
	return acc, nil
}

func SignTxV2(rawTxnImpl *v2.RawTransactionWithData, seedHex string) (*TxWithAuth, error) {
	seedHex = StripAptosPrivateKeyPrefix(seedHex)
	signer := &crypto.Ed25519PrivateKey{}
	if err := signer.FromHex(seedHex); err != nil {
		return nil, err
	}
	auth, err := rawTxnImpl.Sign(signer)
	if err != nil {
		return nil, err
	}
	return &TxWithAuth{
		RawTxn:     rawTxnImpl,
		SenderAuth: auth,
	}, nil
}

type TxWithAuth struct {
	RawTxn          *v2.RawTransactionWithData    `json:"rawTxn"`
	SenderAuth      *crypto.AccountAuthenticator  `json:"senderAuth"`
	FeePayerAuth    *crypto.AccountAuthenticator  `json:"feePayerAuth"`
	AdditionalAuths []crypto.AccountAuthenticator `json:"additionalAuths"`
}

func (self *TxWithAuth) GetRawTxnHex() (string, error) {
	data, err := bcs2.SerializeSingle(func(ser *bcs2.Serializer) {
		self.RawTxn.MarshalTypeScriptBCS(ser)
	})
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}

func (self *TxWithAuth) GetSenderAuthHex() (string, error) {
	data, err := bcs2.Serialize(self.SenderAuth)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}

func (self *TxWithAuth) GetFeePayerAuthHex() (string, error) {
	data, err := bcs2.Serialize(self.FeePayerAuth)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}

func (self *TxWithAuth) GetAdditionalAuthsHex() ([]string, error) {
	var res []string
	for i := 0; i < len(self.AdditionalAuths); i++ {
		data, err := bcs2.Serialize(&self.AdditionalAuths[i])
		if err != nil {
			return nil, err
		}
		res = append(res, hex.EncodeToString(data))
	}
	return res, nil
}

type TxWithAuthJson struct {
	RawTxn          string   `json:"rawTxn"`
	SenderAuth      string   `json:"senderAuth"`
	FeePayerAuth    string   `json:"feePayerAuth"`
	AdditionalAuths []string `json:"additionalAuths"`
}

func (self *TxWithAuth) MarshalJson() (string, error) {
	var txWithAuthJson TxWithAuthJson
	if self.RawTxn != nil {
		data, err := bcs2.SerializeSingle(func(ser *bcs2.Serializer) {
			self.RawTxn.MarshalTypeScriptBCS(ser)
		})
		if err != nil {
			return "", err
		}
		txWithAuthJson.RawTxn = hex.EncodeToString(data)
	}
	if self.SenderAuth != nil {
		data, err := bcs2.Serialize(self.SenderAuth)
		if err != nil {
			return "", err
		}
		txWithAuthJson.SenderAuth = hex.EncodeToString(data)
	}
	if self.FeePayerAuth != nil {
		data, err := bcs2.Serialize(self.FeePayerAuth)
		if err != nil {
			return "", err
		}
		txWithAuthJson.FeePayerAuth = hex.EncodeToString(data)
	}
	for i := 0; i < len(self.AdditionalAuths); i++ {
		data, err := bcs2.Serialize(&self.AdditionalAuths[i])
		if err != nil {
			return "", err
		}
		txWithAuthJson.AdditionalAuths = append(txWithAuthJson.AdditionalAuths, hex.EncodeToString(data))
	}
	data, _ := json.Marshal(txWithAuthJson)
	return string(data), nil
}

func (self *TxWithAuth) UnmarshalJson(dataJson string) error {
	var txJson TxWithAuthJson
	if err := json.Unmarshal([]byte(dataJson), &txJson); err != nil {
		return err
	}
	if txJson.RawTxn != "" {
		rawTx := &v2.RawTransactionWithData{}
		des := bcs2.NewDeserializer(util.RemoveZeroHex(txJson.RawTxn))
		rawTx.UnmarshalTypeScriptBCS(des)
		if err := des.Error(); err != nil {
			return err
		}
		self.RawTxn = rawTx
	}

	if txJson.SenderAuth != "" {
		var senderAuth crypto.AccountAuthenticator
		if err := bcs2.Deserialize(&senderAuth, util.RemoveZeroHex(txJson.SenderAuth)); err != nil {
			return err
		}
		self.SenderAuth = &senderAuth
	}

	if txJson.FeePayerAuth != "" {
		var feePayerAuth crypto.AccountAuthenticator
		if err := bcs2.Deserialize(&feePayerAuth, util.RemoveZeroHex(txJson.FeePayerAuth)); err != nil {
			return err
		}
		self.FeePayerAuth = &feePayerAuth
	}
	if len(txJson.AdditionalAuths) > 0 {
		additionalAuths := make([]crypto.AccountAuthenticator, 0)
		for i := 0; i < len(txJson.AdditionalAuths); i++ {
			var auth crypto.AccountAuthenticator
			if err := bcs2.Deserialize(&auth, util.RemoveZeroHex(txJson.AdditionalAuths[i])); err != nil {
				return err
			}
			additionalAuths = append(additionalAuths, auth)
		}
		self.AdditionalAuths = additionalAuths
	}
	return nil
}

func BuildMultiAgentTx(sender string, sequenceNumber uint64, maxGasAmount uint64,
	gasUnitPrice uint64, expirationTimestampSecs uint64, chainId uint8, payload v2.TransactionPayload,
	feePayer string, additionalSigners []string) (rawTxnImpl *v2.RawTransactionWithData, err error) {
	senderAddr, err := parseAccountAddress(sender)
	if err != nil {
		return nil, err
	}
	rawTxn := &v2.RawTransaction{
		Sender:                     *senderAddr,
		SequenceNumber:             sequenceNumber,
		Payload:                    payload,
		MaxGasAmount:               maxGasAmount,
		GasUnitPrice:               gasUnitPrice,
		ExpirationTimestampSeconds: expirationTimestampSecs,
		ChainId:                    chainId,
	}
	rawAdditionalSigners := make([]v2.AccountAddress, 0)
	for i := 0; i < len(additionalSigners); i++ {
		parsedAddr, err := parseAccountAddress(additionalSigners[i])
		if err != nil {
			return nil, err
		}
		rawAdditionalSigners = append(rawAdditionalSigners, *parsedAddr)
	}
	if feePayer != "" {
		feePayerAddr, err := parseAccountAddress(feePayer)
		if err != nil {
			return nil, err
		}
		return &v2.RawTransactionWithData{
			Variant: v2.MultiAgentWithFeePayerRawTransactionWithDataVariant,
			Inner: &v2.MultiAgentWithFeePayerRawTransactionWithData{
				RawTxn:           rawTxn,
				FeePayer:         feePayerAddr,
				SecondarySigners: rawAdditionalSigners,
			},
		}, nil
	} else {
		return &v2.RawTransactionWithData{
			Variant: v2.MultiAgentRawTransactionWithDataVariant,
			Inner: &v2.MultiAgentRawTransactionWithData{
				RawTxn:           rawTxn,
				SecondarySigners: rawAdditionalSigners,
			},
		}, nil
	}
}

func parseCoinType(tyArg string) (*v2.TypeTag, error) {
	// 0x3::moon_coin::MoonCoin  address (hex) + module + struct
	parts := strings.Split(tyArg, "::")
	contractAddr := &v2.AccountAddress{}
	if err := contractAddr.ParseStringRelaxed(parts[0]); err != nil {
		return nil, err
	}
	return &v2.TypeTag{
		Value: &v2.StructTag{
			Address:    *contractAddr,
			Module:     parts[1],
			Name:       parts[2],
			TypeParams: []v2.TypeTag{},
		},
	}, nil
}

func stringArrToAnyArr(arr []string) (res []any) {
	for _, v := range arr {
		res = append(res, v)
	}
	return
}

const (
	PrivateKeyEd25519Prefix    = "ed25519-priv-"
	PrivateKeyVariantSecp256k1 = "secp256k1-priv-"
)

func StripAptosPrivateKeyPrefix(privateKey string) string {
	if strings.HasPrefix(privateKey, PrivateKeyVariantSecp256k1) {
		return privateKey[len(PrivateKeyVariantSecp256k1):]
	}
	if strings.HasPrefix(privateKey, PrivateKeyEd25519Prefix) {
		return privateKey[len(PrivateKeyEd25519Prefix):]
	} else {
		return privateKey
	}
}

// DeserializeMultiAgentTransaction deserializes a BCS encoded MultiAgent transaction
func DeserializeMultiAgentTransaction(data []byte) (*v2.RawTransactionWithData, error) {
	des := bcs2.NewDeserializer(data)
	rawTxnWithData := &v2.RawTransactionWithData{}

	rawTxn := &v2.RawTransaction{}
	rawTxn.UnmarshalBCS(des)

	if des.Error() != nil {
		return nil, des.Error()
	}
	secondarySigners := bcs2.DeserializeSequence[v2.AccountAddress](des)

	hasFeePayer := des.Bool()
	if hasFeePayer {
		feePayer := &v2.AccountAddress{}
		feePayer.UnmarshalBCS(des)
		rawTxnWithData.Variant = v2.MultiAgentWithFeePayerRawTransactionWithDataVariant
		rawTxnWithData.Inner = &v2.MultiAgentWithFeePayerRawTransactionWithData{
			RawTxn:           rawTxn,
			SecondarySigners: secondarySigners,
			FeePayer:         feePayer,
		}
	} else {
		rawTxnWithData.Variant = v2.MultiAgentRawTransactionWithDataVariant
		rawTxnWithData.Inner = &v2.MultiAgentRawTransactionWithData{
			RawTxn:           rawTxn,
			SecondarySigners: secondarySigners,
		}
	}
	return rawTxnWithData, nil
}

func ConvertToBigInt(v string) *big.Int {
	b := new(big.Int)
	b.SetString(v, 10)
	return b
}

func DecodeHexStringErr(hexString string) ([]byte, error) {
	return hex.DecodeString(util.RemoveHexPrefix(hexString))
}
