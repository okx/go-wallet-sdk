package aptos

import (
	ed255192 "crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/aptos/aptos_types"
	"github.com/okx/go-wallet-sdk/coins/aptos/common"
	"github.com/okx/go-wallet-sdk/coins/aptos/serde"
	"github.com/okx/go-wallet-sdk/coins/aptos/transaction_builder"
	v2 "github.com/okx/go-wallet-sdk/coins/aptos/v2"
	bcs2 "github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/crypto"
	"github.com/okx/go-wallet-sdk/crypto/ed25519"
	"github.com/okx/go-wallet-sdk/util"
	"math/big"
	"regexp"
	"strconv"
	"strings"
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

func NewAddress(seedHex string, shortEnable bool) string {
	seedHex = StripAptosPrivateKeyPrefix(seedHex)
	publicKey, _ := ed25519.PublicKeyFromSeed(seedHex)
	publicKey = append(publicKey, 0x0)
	address := "0x" + hex.EncodeToString(common.Sha256Hash(publicKey))
	if shortEnable {
		return ShortenAddress(address)
	} else {
		return address
	}
}

////////////////////////// API for TEE ////////////////////////

type AptosTXParam struct {
	From                    string `json:"from"`
	SequenceNumber          uint64 `json:"sequenceNumber"`
	MaxGasAmount            uint64 `json:"maxGasAmount"`
	GasUnitPrice            uint64 `json:"gasUnitPrice"`
	ExpirationTimestampSecs uint64 `json:"expirationTimestampSecs"`
	ChainId                 uint8  `json:"chainId"`
	Payload                 string `json:"payload"`
}

func NewTxFromParam(txParam *AptosTXParam) (tx *v2.RawTransaction, err error) {
	txBytes, err := util.DecodeHexStringErr(txParam.Payload)
	if err != nil {
		return nil, err
	}
	der := bcs2.NewDeserializer(txBytes)
	payload := &v2.TransactionPayload{}
	payload.UnmarshalBCS(der)
	err = der.Error()
	if err != nil {
		return
	}
	return MakeRawTransactionV2(txParam.From, txParam.SequenceNumber, txParam.MaxGasAmount, txParam.GasUnitPrice, txParam.ExpirationTimestampSecs,
		txParam.ChainId, payload)
}

// todo
func NewTxFromRaw(rawTx string, isMultiAgent bool) (tx interface{}, err error) {
	txBytes, err := util.DecodeHexStringErr(rawTx)
	if err != nil {
		return nil, err
	}
	der := bcs2.NewDeserializer(txBytes)
	if isMultiAgent {
		var rawTx = &v2.RawTransactionWithData{}
		rawTx.UnmarshalBCS(der)
		err = der.Error()
		return rawTx, err
	} else {
		var rawTx = &v2.RawTransaction{}
		rawTx.UnmarshalBCS(der)
		err = der.Error()
		return rawTx, err
	}
}

func GetSigningData(rawTxnImpl interface{}) (data []byte, err error) {
	switch rawTx := rawTxnImpl.(type) {
	case *v2.RawTransactionWithData:
		return rawTx.SigningMessage()
	case *v2.RawTransaction:
		return rawTx.SigningMessage()
	default:
		return nil, errors.New("unknown transaction type")
	}
}

// func AddSignature(tx Transaction, sig []byte) (data TxData, err error)
func NewAddressFromPubkey(pubkeyHex string) (addr string, err error) {
	return NewPubKeyAddress(pubkeyHex, false)
}

////////////////////////// API for TEE  END ////////////////////////

func NewPubKeyAddress(pubHex string, shortEnable bool) (string, error) {
	if strings.HasPrefix(pubHex, "0x") {
		pubHex = pubHex[2:]
	}
	pub, err := hex.DecodeString(pubHex)
	if err != nil {
		return "", err
	}
	if len(pub) != 32 {
		return "", errors.New("invalid public key")
	}
	publicKey := ed255192.PublicKey(pub)
	publicKey = append(publicKey, 0x0)
	address := "0x" + hex.EncodeToString(common.Sha256Hash(publicKey))
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
	address := "0x" + hex.EncodeToString(common.Sha256Hash(pubKey))
	if shortEnable {
		return ShortenAddress(address), nil
	} else {
		return address, nil
	}
}

func ValidateAddress(address string, shortEnable bool) bool {
	re1, _ := regexp.Compile("^0x[\\dA-Fa-f]{62,64}$")
	re2, _ := regexp.Compile("^[\\dA-Fa-f]{64}$")
	return re1.Match([]byte(address)) || re2.Match([]byte(address))
}

func ValidateContractAddress(address string) bool {
	contractReg, err := regexp.Compile("^(0[x|X])?[\\dA-Fa-f]+::.+::.+")
	if err != nil {
		panic(err)
	}
	if contractReg.Match([]byte(address)) {
		return true
	}
	contractReg2, err := regexp.Compile("^(0[x|X])?[\\dA-Fa-f]*$")
	if err != nil {
		panic(err)
	}
	return contractReg2.Match([]byte(address))
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

func MakeRawTransactionV2(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64,
	expirationTimestampSecs uint64, chainId uint8, payload *v2.TransactionPayload) (*v2.RawTransaction, error) {
	rawTxn := v2.RawTransaction{}
	// addr, err := aptos_types.FromHex(ExpandAddress(from))
	addr := &v2.AccountAddress{}
	if err := addr.ParseStringRelaxed(from); err != nil {
		return nil, err
	}
	rawTxn.Sender = *addr
	rawTxn.SequenceNumber = sequenceNumber
	rawTxn.MaxGasAmount = maxGasAmount
	rawTxn.GasUnitPrice = gasUnitPrice
	rawTxn.ExpirationTimestampSeconds = expirationTimestampSecs
	rawTxn.ChainId = chainId
	rawTxn.Payload = *payload
	return &rawTxn, nil
}

func MakeMultiAgentTransactionV2(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64,
	expirationTimestampSecs uint64, chainId uint8, payload *v2.TransactionPayload, withFeePayer bool,
	feePayer string, additionalSigners []string) (*v2.RawTransactionWithData, error) {
	addr := &v2.AccountAddress{}
	if err := addr.ParseStringRelaxed(from); err != nil {
		return nil, err
	}
	rawTxn := &v2.RawTransaction{}
	rawTxn.Sender = *addr
	rawTxn.SequenceNumber = sequenceNumber
	rawTxn.MaxGasAmount = maxGasAmount
	rawTxn.GasUnitPrice = gasUnitPrice
	rawTxn.ExpirationTimestampSeconds = expirationTimestampSecs
	rawTxn.ChainId = chainId
	rawTxn.Payload = *payload
	secondarySigners := make([]v2.AccountAddress, 0)
	for i := 0; i < len(additionalSigners); i++ {
		var signer = &v2.AccountAddress{}
		if err := signer.ParseStringRelaxed(additionalSigners[i]); err != nil {
			return nil, err
		}
		secondarySigners = append(secondarySigners, *signer)
	}

	if withFeePayer {
		var feePayerAddr = &v2.AccountAddress{}
		if feePayer == "" {
			feePayerAddr = &v2.AccountZero
		} else {
			if err := feePayerAddr.ParseStringRelaxed(feePayer); err != nil {
				return nil, err
			}
		}
		return &v2.RawTransactionWithData{
			Variant: v2.MultiAgentWithFeePayerRawTransactionWithDataVariant,
			Inner: &v2.MultiAgentWithFeePayerRawTransactionWithData{
				RawTxn:           rawTxn,
				SecondarySigners: secondarySigners,
				FeePayer:         feePayerAddr,
			},
		}, nil
	} else if len(additionalSigners) > 0 {
		return &v2.RawTransactionWithData{
			Variant: v2.MultiAgentRawTransactionWithDataVariant,
			Inner: &v2.MultiAgentRawTransactionWithData{
				RawTxn:           rawTxn,
				SecondarySigners: secondarySigners,
			},
		}, nil
	} else {
		return nil, errors.New("no additional signers provided")
	}
}

func Transfer(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64, expirationTimestampSecs uint64, chainId uint8,
	to string, amount uint64, seedHex string) (string, error) {
	seedHex = StripAptosPrivateKeyPrefix(seedHex)
	payload, err := TransferPayload(to, amount)
	if err != nil {
		return "", err
	}
	return BuildSignedTransaction(from, sequenceNumber, maxGasAmount, gasUnitPrice, expirationTimestampSecs, chainId, payload, seedHex)
}

func TransferWithFeePayer(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64, expirationTimestampSecs uint64, chainId uint8,
	to string, amount uint64, seedHex string, feePayer string) (*TxWithAuth, error) {
	multiTx, err := BuildTransferWithFeePayerTx(from, sequenceNumber, maxGasAmount, gasUnitPrice, expirationTimestampSecs,
		chainId, to, amount, feePayer)
	if err != nil {
		return nil, err
	}
	return SignTxV2(multiTx, seedHex)
}

func BuildTransferWithFeePayerTx(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64, expirationTimestampSecs uint64, chainId uint8,
	to string, amount uint64, feePayer string) (*v2.RawTransactionWithData, error) {
	toAddr, err := parseAccountAddress(to)
	if err != nil {
		return nil, err
	}
	payload, err := v2.CoinTransferPayload(&v2.AptosCoinTypeTag, *toAddr, amount)
	if err != nil {
		return nil, err
	}
	if feePayer == "" {
		feePayer = v2.AccountZero.String()
	}
	multiTx, err := BuildMultiAgentTx(from, sequenceNumber, maxGasAmount, gasUnitPrice, expirationTimestampSecs, chainId,
		v2.TransactionPayload{
			Payload: payload,
		}, feePayer, []string{})
	if err != nil {
		return nil, err
	}
	return multiTx, nil
}

func SignAsFeePayer(rawTxn, seedHex string, feePayer string) (string, string, error) {
	seedHex = StripAptosPrivateKeyPrefix(seedHex)
	var feePayerAddr *v2.AccountAddress
	var err error
	if feePayer == "" {
		addr := NewAddress(seedHex, false)
		feePayerAddr, err = parseAccountAddress(addr)
		if err != nil {
			return "", "", err
		}
	} else {
		feePayerAddr, err = parseAccountAddress(feePayer)
		if err != nil {
			return "", "", err
		}
	}

	rawTx, auth, err := signMultiAgentTx(rawTxn, seedHex, feePayerAddr)
	if err != nil {
		return "", "", err
	}
	feePayerBytes, err := bcs2.Serialize(auth)
	if err != nil {
		return "", "", err
	}
	data, err := bcs2.Serialize(rawTx)
	if err != nil {
		return "", "", err
	}
	return hex.EncodeToString(feePayerBytes), hex.EncodeToString(data), nil
}

func SignSimpleTx(rawTxn, seedHex string) (string, error) {
	seedHex = StripAptosPrivateKeyPrefix(seedHex)
	rawTx, rawTxWithData, err := deserializeSimpleTx(rawTxn)
	if err != nil {
		return "", err
	}
	signer := &crypto.Ed25519PrivateKey{}
	if err = signer.FromHex(seedHex); err != nil {
		return "", err
	}
	var auth *crypto.AccountAuthenticator
	if rawTx != nil {
		auth, err = rawTx.Sign(signer)
		if err != nil {
			return "", err
		}
		txAuth, err := v2.NewTransactionAuthenticator(auth)
		if err != nil {
			return "", err
		}
		signedTx := &v2.SignedTransaction{
			Transaction:   rawTx,
			Authenticator: txAuth,
		}
		txHash, err := signedTx.Hash()
		if err != nil {
			return "", err
		}
		signedTxBytes, err := bcs2.Serialize(signedTx)
		if err != nil {
			return "", err
		}
		authBytes, err := bcs2.Serialize(auth)
		if err != nil {
			return "", err
		}
		return common.FormatSimpleTxRespTxResp(txHash, hex.EncodeToString(signedTxBytes), hex.EncodeToString(authBytes))
	} else if rawTxWithData != nil {
		auth, err = rawTxWithData.Sign(signer)
		if err != nil {
			return "", err
		}
		authBytes, err := bcs2.Serialize(auth)
		if err != nil {
			return "", err
		}
		return common.FormatSignMultiAgentTxResp(rawTxn, hex.EncodeToString(authBytes))
	} else {
		return "", errors.New("invalid rawTx")
	}
}

func deserializeSimpleTx(rawTxn string) (*v2.RawTransaction, *v2.RawTransactionWithData, error) {
	rawTxnBytes := util.RemoveZeroHex(rawTxn)
	rawTx := &v2.RawTransaction{}
	des := bcs2.NewDeserializer(rawTxnBytes)
	rawTx.UnmarshalBCS(des)
	err := des.Error()
	if err != nil {
		return nil, nil, err
	}
	feepayerPresent := des.Bool()
	if feepayerPresent {
		feePayer := &v2.AccountAddress{}
		feePayer.UnmarshalBCS(des)
		return nil, &v2.RawTransactionWithData{
			Variant: v2.MultiAgentWithFeePayerRawTransactionWithDataVariant,
			Inner: &v2.MultiAgentWithFeePayerRawTransactionWithData{
				RawTxn:   rawTx,
				FeePayer: feePayer,
			},
		}, nil
	}
	return rawTx, nil, nil

}

func SignMultiAgentTx(rawTxn, seedHex string) (string, error) {
	seedHex = StripAptosPrivateKeyPrefix(seedHex)
	_, auth, err := signMultiAgentTx(rawTxn, seedHex, nil)
	if err != nil {
		return "", err
	}
	authBytes, err := bcs2.Serialize(auth)
	if err != nil {
		return "", err
	}
	return common.FormatSignMultiAgentTxResp(rawTxn, hex.EncodeToString(authBytes))
}

func ToFinalMultiAgentTx(txWithAuth *TxWithAuth) (string, error) {
	var signedTx *v2.SignedTransaction
	var ok bool
	if txWithAuth.FeePayerAuth != nil {
		signedTx, ok = txWithAuth.RawTxn.ToFeePayerSignedTransaction(txWithAuth.SenderAuth, txWithAuth.FeePayerAuth, txWithAuth.AdditionalAuths)
	} else {
		signedTx, ok = txWithAuth.RawTxn.ToMultiAgentSignedTransaction(txWithAuth.SenderAuth, txWithAuth.AdditionalAuths)
	}
	if !ok {
		return "", errors.New("failed to convert raw transaction to signed transaction")
	}
	signedTxBytes, err := bcs2.Serialize(signedTx)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(signedTxBytes), nil
}

func signMultiAgentTx(rawTxn, seedHex string, feePayerAddr *v2.AccountAddress) (multiTx *v2.RawTransactionWithData, auth *crypto.AccountAuthenticator, err error) {
	seedHex = StripAptosPrivateKeyPrefix(seedHex)
	rawTxnBytes := util.RemoveZeroHex(rawTxn)
	var rawTx = &v2.RawTransactionWithData{}
	des := bcs2.NewDeserializer(rawTxnBytes)
	rawTx.UnmarshalTypeScriptBCS(des)
	if des.Error() != nil {
		return nil, nil, des.Error()
	}
	if feePayerAddr != nil {
		ok := rawTx.SetFeePayer(*feePayerAddr)
		if !ok {
			return nil, nil, errors.New("failed to set fee payer address")
		}
	}
	signer := &crypto.Ed25519PrivateKey{}
	if err = signer.FromHex(seedHex); err != nil {
		return
	}
	auth, err = rawTx.Sign(signer)
	if err != nil {
		return
	}
	return rawTx, auth, nil
}

func TransferCoins(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64, expirationTimestampSecs uint64, chainId uint8,
	to string, amount uint64, seedHex, tyArg string) (string, error) {
	seedHex = StripAptosPrivateKeyPrefix(seedHex)
	payload, err := CoinTransferPayloadV2(to, amount, tyArg)
	if err != nil {
		return "", err
	}
	return BuildSignedTransaction(from, sequenceNumber, maxGasAmount, gasUnitPrice, expirationTimestampSecs, chainId, payload, seedHex)
}

func V2CoinTransferPayload(to string, amount uint64, tyArg string) (*v2.EntryFunction, error) {
	coinTy, err := parseCoinType(tyArg)
	if err != nil {
		return nil, err
	}
	toAddr, err := parseAccountAddress(to)
	if err != nil {
		return nil, err
	}
	return v2.CoinTransferPayload(coinTy, *toAddr, amount)
}

func TransferCoinsWithPayer(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64, expirationTimestampSecs uint64, chainId uint8,
	to string, amount uint64, seedHex, tyArg string, feePayer string) (*TxWithAuth, error) {
	seedHex = StripAptosPrivateKeyPrefix(seedHex)
	coinTy, err := parseCoinType(tyArg)
	if err != nil {
		return nil, err
	}
	toAddr, err := parseAccountAddress(to)
	if err != nil {
		return nil, err
	}
	payload, err := v2.CoinTransferPayload(coinTy, *toAddr, amount)
	if err != nil {
		return nil, err
	}
	if feePayer == "" {
		feePayer = v2.AccountZero.String()
	}
	multiTx, err := BuildMultiAgentTx(from, sequenceNumber, maxGasAmount, gasUnitPrice, expirationTimestampSecs, chainId,
		v2.TransactionPayload{
			Payload: payload,
		}, feePayer, []string{})
	if err != nil {
		return nil, err
	}
	return SignTxV2(multiTx, seedHex)
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
	seedHex = StripAptosPrivateKeyPrefix(seedHex)
	rawTxn, err := MakeRawTransaction(from, sequenceNumber, maxGasAmount, gasUnitPrice, expirationTimestampSecs, chainId, payload)
	if err != nil {
		return "", err
	}
	return SignRawTransaction(rawTxn, seedHex)
}

func BuildSignedTransactionV2(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64, expirationTimestampSecs uint64, chainId uint8,
	payload *v2.TransactionPayload, seedHex string) (string, error) {
	seedHex = StripAptosPrivateKeyPrefix(seedHex)
	rawTxn, err := MakeRawTransactionV2(from, sequenceNumber, maxGasAmount, gasUnitPrice, expirationTimestampSecs, chainId, payload)
	if err != nil {
		return "", err
	}
	return SignRawTransactionV2(rawTxn, seedHex)
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
func GetTransactionHash(hexStr string) string {
	prefix := common.Sha256Hash([]byte("APTOS::Transaction"))
	return common.ComputeTransactionHash(prefix, hexStr)
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
	// 0x3::moon_coin::MoonCoin  address (hex) + module + struct
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

func CoinTransferPayloadV2(to string, amount uint64, tyArg string) (aptos_types.TransactionPayload, error) {
	bscAddress, err := aptos_types.BcsSerializeFixedBytes(aptos_types.BytesFromHex(ExpandAddress(to)))
	if err != nil {
		return nil, err
	}
	bscAmount, err := aptos_types.BcsSerializeUint64(amount)
	if err != nil {
		return nil, err
	}
	// 0x3::moon_coin::MoonCoin  address (hex) + module + struct
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
	//0x1::aptos_account::transfer_coins
	scriptFunction := aptos_types.ScriptFunction{
		Module:   aptos_types.ModuleId{Address: *aptos_types.CORE_CODE_ADDRESS, Name: "aptos_account"},
		Function: "transfer_coins",
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

func OfferNFTTokenPayload(receiver string, creator string, collectionName string, tokenName string, propertyVersion uint64, amount uint64) aptos_types.TransactionPayload {
	// moduleAddress := make([]byte, 31)
	// moduleAddress = append(moduleAddress, 0x3)
	moduleAddress, _ := aptos_types.FromHex("0x3")
	bscReceiver, _ := aptos_types.BcsSerializeFixedBytes(aptos_types.BytesFromHex(ExpandAddress(receiver)))
	bscCreator, _ := aptos_types.BcsSerializeFixedBytes(aptos_types.BytesFromHex(ExpandAddress(creator)))
	bscCollectName, _ := aptos_types.BcsSerializeStr(collectionName)
	bscTokenName, _ := aptos_types.BcsSerializeStr(tokenName)
	bscPropertyVersion, _ := aptos_types.BcsSerializeUint64(propertyVersion)
	bscAmount, _ := aptos_types.BcsSerializeUint64(amount)

	scriptFunction := aptos_types.ScriptFunction{
		Module:   aptos_types.ModuleId{Address: *moduleAddress, Name: "token_transfers"},
		Function: "offer_script",
		TyArgs:   []aptos_types.TypeTag{},
		Args:     [][]byte{bscReceiver, bscCreator, bscCollectName, bscTokenName, bscPropertyVersion, bscAmount},
	}
	return &aptos_types.TransactionPayloadEntryFunction{Value: scriptFunction}
}

func ClaimNFTTokenPayload(sender string, creator string, collectionName string, tokenName string, propertyVersion uint64) aptos_types.TransactionPayload {
	// moduleAddress := make([]byte, 31)
	// moduleAddress = append(moduleAddress, 0x3)
	moduleAddress, _ := aptos_types.FromHex("0x3")
	bscSender, _ := aptos_types.BcsSerializeFixedBytes(aptos_types.BytesFromHex(ExpandAddress(sender)))
	bscCreator, _ := aptos_types.BcsSerializeFixedBytes(aptos_types.BytesFromHex(ExpandAddress(creator)))
	bscCollectName, _ := aptos_types.BcsSerializeStr(collectionName)
	bscTokenName, _ := aptos_types.BcsSerializeStr(tokenName)
	bscPropertyVersion, _ := aptos_types.BcsSerializeUint64(propertyVersion)

	scriptFunction := aptos_types.ScriptFunction{
		Module:   aptos_types.ModuleId{Address: *moduleAddress, Name: "token_transfers"},
		Function: "claim_script",
		TyArgs:   []aptos_types.TypeTag{},
		Args:     [][]byte{bscSender, bscCreator, bscCollectName, bscTokenName, bscPropertyVersion},
	}
	return &aptos_types.TransactionPayloadEntryFunction{Value: scriptFunction}
}

func SignRawTransaction(rawTxn *aptos_types.RawTransaction, seedHex string) (string, error) {
	seedHex = StripAptosPrivateKeyPrefix(seedHex)

	publicKeyBs, err := ed25519.PublicKeyFromSeed(seedHex)
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
	ed25519Authenticator := aptos_types.TransactionAuthenticatorEd25519{PublicKey: aptos_types.Ed25519PublicKey(publicKeyBs), Signature: signature}
	signedTransaction := aptos_types.SignedTransaction{RawTxn: *rawTxn, Authenticator: &ed25519Authenticator}
	txBytes, err := signedTransaction.BcsSerialize()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(txBytes), nil
}

func SignRawTransactionV2(rawTxn *v2.RawTransaction, seedHex string) (string, error) {
	seedHex = StripAptosPrivateKeyPrefix(seedHex)
	signer := &crypto.Ed25519PrivateKey{}
	if err := signer.FromBytes(util.RemoveZeroHex(seedHex)); err != nil {
		return "", err
	}
	signedTx, err := rawTxn.SignedTransaction(signer)
	if err != nil {
		return "", err
	}
	ser := &bcs2.Serializer{}
	signedTx.MarshalBCS(ser)
	if ser.Error() != nil {
		return "", ser.Error()
	}
	return hex.EncodeToString(ser.ToBytes()), nil
}

func SimulateTransaction(rawTxn *aptos_types.RawTransaction, seedHex string) (string, error) {
	seedHex = StripAptosPrivateKeyPrefix(seedHex)
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
	p1, _ := aptos_types.FromHex(ExpandAddress(parts[0]))
	p2 := aptos_types.Identifier(parts[1])
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

func ConvertArgs(args []interface{}, argTypes []aptos_types.MoveType) ([][]byte, error) {
	if len(args) != len(argTypes) {
		return nil, fmt.Errorf("types and values size not match")
	}
	array := make([][]byte, 0)
	for i := range args {
		moveType := argTypes[i]
		moveValue := args[i]
		if moveType == "0x1::object::Object" || (strings.Contains(moveType, "<") && moveType[:strings.Index(moveType, "<")] == "0x1::object::Object") {
			moveType = "address"
		}
		switch moveType {
		case "address":
			op, ok := moveValue.(string)
			if !ok {
				return nil, fmt.Errorf("unknown argument for address")
			}
			bytes, _ := aptos_types.BcsSerializeFixedBytes(aptos_types.BytesFromHex(ExpandAddress(op)))
			array = append(array, bytes)
		case "u64":
			ai, err := Interface2U64(moveValue)
			if err != nil {
				return nil, fmt.Errorf("unknown argument for u64")
			}
			bytes, _ := aptos_types.BcsSerializeUint64(ai)
			array = append(array, bytes)
		case "bool":
			op := fmt.Sprintf("%v", moveValue)
			b, err := strconv.ParseBool(op)
			if err != nil {
				return nil, fmt.Errorf("unknown argument for bool")
			}
			bytes, _ := aptos_types.BcsSerializeBool(b)
			array = append(array, bytes)
		case "u8":
			op := fmt.Sprintf("%v", moveValue)
			ai, err := strconv.ParseUint(op, 0, 8)
			if err != nil {
				return nil, fmt.Errorf("unknown argument for u8")
			}
			bytes, _ := aptos_types.BcsSerializeU8(uint8(ai))
			array = append(array, bytes)
		case "u128":
			ii, err := Interface2U128(moveValue)
			if err != nil {
				return nil, fmt.Errorf("unknown argument for u128")
			}
			bytes, _ := aptos_types.BcsSerializeU128(*ii)
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
				bytes, _ := aptos_types.BcsSerializeBytes(inputBytes)
				array = append(array, bytes)
			case string:
				op, _ := moveValue.(string)
				v := aptos_types.BytesFromHex(op)
				bytes, _ := aptos_types.BcsSerializeBytes(v)
				array = append(array, bytes)
			default:
				return nil, fmt.Errorf("unknown argument for vector<u8>")
			}
		case "vector<u64>":
			switch moveValue.(type) {
			case []interface{}:
				vArray, _ := moveValue.([]interface{})
				targetBytes := make([]byte, 0)
				// 数组长度
				bytes, _ := aptos_types.BcsSerializeLen(uint64(len(vArray)))
				targetBytes = append(targetBytes, bytes...)
				// 序列化每一项
				for _, e := range vArray {
					v, err := Interface2U64(e)
					if err != nil {
						return nil, fmt.Errorf("unknown argument for u64")
					}
					bytes, _ = aptos_types.BcsSerializeUint64(v)
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
				// 数组长度
				bytes, _ := aptos_types.BcsSerializeLen(uint64(len(vArray)))
				targetBytes = append(targetBytes, bytes...)
				// 序列化每一项
				for _, e := range vArray {
					v, err := Interface2U128(e)
					if err != nil {
						return nil, fmt.Errorf("unknown argument for u128")
					}
					bytes, _ = aptos_types.BcsSerializeU128(*v)
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
				bytes, _ := aptos_types.BcsSerializeLen(uint64(len(vArray)))
				targetBytes = append(targetBytes, bytes...)
				for _, e := range vArray {
					v := fmt.Sprintf("%v", e)
					vv, err := strconv.ParseBool(v)
					if err != nil {
						return nil, fmt.Errorf("unknown argument for bool")
					}
					bytes, _ = aptos_types.BcsSerializeBool(vv)
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
				bytes, _ := aptos_types.BcsSerializeLen(uint64(len(vArray)))
				targetBytes = append(targetBytes, bytes...)

				for _, e := range vArray {
					v, ok := e.(string)
					if !ok {
						return nil, fmt.Errorf("unknown argument for address")
					}
					bytes, _ = aptos_types.BcsSerializeFixedBytes(aptos_types.BytesFromHex(ExpandAddress(v)))
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
				bytes, _ := aptos_types.BcsSerializeLen(uint64(len(vArray)))
				targetBytes = append(targetBytes, bytes...)

				for _, e := range vArray {
					v, ok := e.(string)
					if !ok {
						return nil, fmt.Errorf("unknown argument for string")
					}
					bytes, _ = aptos_types.BcsSerializeStr(v)
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
			bytes, _ := aptos_types.BcsSerializeStr(op)
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

func BuildScriptPayload(payload string) (*v2.TransactionPayload, error) {
	var param ScriptParam
	err := json.Unmarshal([]byte(payload), &param)
	if err != nil {
		return nil, err
	}
	s, err := parseScriptParam(&param)
	if err != nil {
		return nil, err
	}
	return &v2.TransactionPayload{
		Payload: s,
	}, nil
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
	typeArguments := entryFunction.TypeArguments
	if len(typeArguments) == 0 {
		typeArguments = entryFunction.TypeArgumentsV2
	}

	arguments := entryFunction.Arguments
	if len(arguments) == 0 {
		arguments = entryFunction.FunctionArgumentsV2
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
	typeArgsEntryFunction := typeArguments
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
	for _, tagString := range typeArguments {
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
	args, err := transaction_builder.ToBCSArgs(typeArgABIs, arguments)
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

func PayloadFromJsonAndAbiV2(payload string, abi string, options ...any) (*v2.TransactionPayload, error) {
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
	typeArguments := entryFunction.TypeArguments
	if len(typeArguments) == 0 {
		typeArguments = entryFunction.TypeArgumentsV2
	}

	arguments := entryFunction.Arguments
	if len(arguments) == 0 {
		arguments = entryFunction.FunctionArgumentsV2
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
	var moduleAddress v2.AccountAddress
	err = moduleAddress.ParseStringRelaxed(funcParts[0])
	if err != nil {
		return nil, err
	}
	ef, err := v2.EntryFunctionFromAbi(&funcAbi.MoveFunction, moduleAddress, funcParts[1], funcParts[2],
		stringArrToAnyArr(typeArguments), arguments, options...)
	if err != nil {
		return nil, err
	}

	return &v2.TransactionPayload{
		Payload: ef,
	}, nil
}

func PayloadFromSerializedHex(serializedHex string) (*v2.TransactionPayload, error) {
	data := util.RemoveZeroHex(serializedHex)
	if len(data) == 0 {
		return nil, errors.New("serialized hex empty")
	}
	pay := &v2.TransactionPayload{}
	deserializer := bcs2.NewDeserializer(data)

	pay.UnmarshalBCS(deserializer)
	if deserializer.Error() != nil {
		return nil, errors.New("unmarshalBCS payload fail")
	}

	return pay, nil
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
	txBytes, _ := signedTransaction.BcsSerialize()
	return hex.EncodeToString(txBytes), nil
}

// Staking related
func AddStake(from string, sequenceNumber uint64, maxGasAmount uint64, gasUnitPrice uint64, expirationTimestampSecs uint64, chainId uint8,
	poolAddress string, amount uint64, seedHex string) (string, error) {
	seedHex = StripAptosPrivateKeyPrefix(seedHex)
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
	seedHex = StripAptosPrivateKeyPrefix(seedHex)
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
	seedHex = StripAptosPrivateKeyPrefix(seedHex)
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
	seedHex = StripAptosPrivateKeyPrefix(seedHex)
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
func VerifyMessage(publicKey, message, signature string) error {
	if message == "" || signature == "" {
		return errors.New("invalid params")
	}

	pubBytes, err := hex.DecodeString(strings.TrimPrefix(publicKey, "0x"))
	if err != nil {
		return err
	}
	if len(pubBytes) != 32 {
		return errors.New("invalid public key")
	}

	sigBytes, err := hex.DecodeString(strings.TrimPrefix(signature, "0x"))
	if err != nil {
		return err
	}

	if !ed255192.Verify(pubBytes, []byte(message), sigBytes) {
		return errors.New("verify failed")
	}
	return nil
}
