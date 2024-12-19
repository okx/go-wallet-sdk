package terra

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"math/big"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/okx/go-wallet-sdk/coins/cosmos/tx"
	"github.com/okx/go-wallet-sdk/coins/cosmos/types"
	"golang.org/x/crypto/ripemd160"
)

type TransactionInput struct {
	ChainId       string
	Memo          string
	Sequence      uint64
	AccountNumber uint64
	Fee           types.Coins
	GasLimit      uint64

	SendArray     []types.MsgSend
	SwapArray     []MsgSwap
	ContractArray []MsgExecuteContract
}

func (m *TransactionInput) AppendFeeCoin(demon string, amount *big.Int) {
	feeCoin := types.NewCoin(demon, types.NewIntFromBigInt(amount))
	if m.Fee == nil {
		m.Fee = types.NewCoins(feeCoin)
	} else {
		m.Fee = append(m.Fee, feeCoin)
	}
}

func (m *TransactionInput) AppendSendMsg(from string, to string, coins *types.Coins) {
	msg := types.MsgSend{}
	msg.FromAddress = from
	msg.ToAddress = to
	msg.Amount = *coins
	if m.SendArray == nil {
		m.SendArray = make([]types.MsgSend, 0)
	}
	m.SendArray = append(m.SendArray, msg)
}

func (m *TransactionInput) AppendSwapMsg(trader string, askDemon string, demon string, amount *big.Int) {
	msg := MsgSwap{}
	msg.Trader = trader
	msg.AskDenom = askDemon
	msg.OfferCoin = types.NewCoin(demon, types.NewIntFromBigInt(amount))
	if m.SwapArray == nil {
		m.SwapArray = make([]MsgSwap, 0)
	}
	m.SwapArray = append(m.SwapArray, msg)
}

func (m *TransactionInput) AppendContractMsg(sender string, contract string, executeMsg string, coins *types.Coins) {
	msg := MsgExecuteContract{}
	msg.Sender = sender
	msg.Contract = contract
	msg.ExecuteMsg = []byte(executeMsg)
	if coins != nil {
		msg.Coins = *coins
	}
	if m.ContractArray == nil {
		m.ContractArray = make([]MsgExecuteContract, 0)
	}
	m.ContractArray = append(m.ContractArray, msg)
}

func NewAddress(privateKeyHex string) (string, error) {
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", err
	}
	_, publicKey := btcec.PrivKeyFromBytes(privateKeyBytes)
	sha := sha256.Sum256(publicKey.SerializeCompressed())
	hasherRIPEMD160 := ripemd160.New()
	hasherRIPEMD160.Write(sha[:])
	address, err := bech32.EncodeFromBase256(HRP, hasherRIPEMD160.Sum(nil))
	if err != nil {
		return "", err
	}
	return address, nil
}

func ValidateAddress(address string) bool {
	hrp, _, err := bech32.DecodeToBase256(address)
	return err == nil && hrp == HRP
}

func NewTransaction(input TransactionInput, privateKeyHex string) string {
	privateKeyBytes, _ := hex.DecodeString(privateKeyHex)
	_, publicKey := btcec.PrivKeyFromBytes(privateKeyBytes)
	messages := make([]*types.Any, 0)
	if input.SendArray != nil {
		for _, msg := range input.SendArray {
			anyMsg, _ := types.NewAnyWithValue(&msg)
			messages = append(messages, anyMsg)
		}
	}
	if input.SwapArray != nil {
		for _, msg := range input.SwapArray {
			anyMsg, _ := types.NewAnyWithValue(&msg)
			messages = append(messages, anyMsg)
		}
	}
	if input.ContractArray != nil {
		for _, msg := range input.ContractArray {
			anyMsg, _ := types.NewAnyWithValue(&msg)
			messages = append(messages, anyMsg)
		}
	}
	body := tx.TxBody{Messages: messages, Memo: input.Memo, TimeoutHeight: 0}
	pubkey := types.PubKey{Key: publicKey.SerializeCompressed()}
	anyPubkey, _ := types.NewAnyWithValue(&pubkey)
	single := tx.ModeInfo_Single{Mode: types.SignMode_SIGN_MODE_DIRECT}
	single_ := tx.ModeInfo_Single_{Single: &single}
	modeInfo := tx.ModeInfo{Sum: &single_}
	signerInfo := make([]*tx.SignerInfo, 0)
	signerInfo = append(signerInfo, &tx.SignerInfo{PublicKey: anyPubkey, ModeInfo: &modeInfo, Sequence: input.Sequence})
	fee := tx.Fee{Amount: input.Fee, GasLimit: input.GasLimit}
	authInfo := tx.AuthInfo{SignerInfos: signerInfo, Fee: &fee}
	bodyBytes, _ := body.Marshal()
	authInfoBytes, _ := authInfo.Marshal()
	signDoc := tx.SignDoc{BodyBytes: bodyBytes, AuthInfoBytes: authInfoBytes, ChainId: input.ChainId, AccountNumber: input.AccountNumber}
	signDocBtyes, _ := signDoc.Marshal()
	return hex.EncodeToString(signDocBtyes)
}

func NewTransactionWithTypeUrl(input TransactionInput, privateKeyHex string) string {
	privateKeyBytes, _ := hex.DecodeString(privateKeyHex)
	_, publicKey := btcec.PrivKeyFromBytes(privateKeyBytes)
	messages := make([]*types.Any, 0)
	if input.SendArray != nil {
		for _, msg := range input.SendArray {
			anyMsg, _ := types.NewAnyWithValue(&msg)
			messages = append(messages, anyMsg)
		}
	}
	if input.SwapArray != nil {
		for _, msg := range input.SwapArray {
			anyMsg, _ := types.NewAnyWithValue(&msg)
			messages = append(messages, anyMsg)
		}
	}
	if input.ContractArray != nil {
		for _, msg := range input.ContractArray {
			anyMsg, _ := types.NewAnyWithValueAndName(&msg)
			messages = append(messages, anyMsg)
		}
	}
	body := tx.TxBody{Messages: messages, Memo: input.Memo, TimeoutHeight: 0}
	pubkey := types.PubKey{Key: publicKey.SerializeCompressed()}
	anyPubkey, _ := types.NewAnyWithValue(&pubkey)
	single := tx.ModeInfo_Single{Mode: types.SignMode_SIGN_MODE_DIRECT}
	single_ := tx.ModeInfo_Single_{Single: &single}
	modeInfo := tx.ModeInfo{Sum: &single_}
	signerInfo := make([]*tx.SignerInfo, 0)
	signerInfo = append(signerInfo, &tx.SignerInfo{PublicKey: anyPubkey, ModeInfo: &modeInfo, Sequence: input.Sequence})
	fee := tx.Fee{Amount: input.Fee, GasLimit: input.GasLimit}
	authInfo := tx.AuthInfo{SignerInfos: signerInfo, Fee: &fee}
	bodyBytes, _ := body.Marshal()
	authInfoBytes, _ := authInfo.Marshal()
	signDoc := tx.SignDoc{BodyBytes: bodyBytes, AuthInfoBytes: authInfoBytes, ChainId: input.ChainId, AccountNumber: input.AccountNumber}
	signDocBtyes, _ := signDoc.Marshal()
	return hex.EncodeToString(signDocBtyes)
}

func Sign(rawHex string, privateKeyHex string) string {
	privateKeyBytes, _ := hex.DecodeString(privateKeyHex)
	privateKey, _ := btcec.PrivKeyFromBytes(privateKeyBytes)

	signDocBtyes, _ := hex.DecodeString(rawHex)
	hash := sha256.Sum256(signDocBtyes)
	b := ecdsa.SignCompact(privateKey, hash[:], false)
	return hex.EncodeToString(b[1:])
}

func SignEnd(rawHex string, signHex string) string {
	signDocBytes, _ := hex.DecodeString(rawHex)
	var signDoc tx.SignDoc
	signDoc.Unmarshal(signDocBytes)

	signBytes, _ := hex.DecodeString(signHex)
	signatures := make([][]byte, 0)
	signatures = append(signatures, signBytes)

	trans := tx.TxRaw{BodyBytes: signDoc.BodyBytes, AuthInfoBytes: signDoc.AuthInfoBytes, Signatures: signatures}
	transBytes, _ := trans.Marshal()
	return base64.StdEncoding.EncodeToString(transBytes)
}

func GetRawTxHex(input TransactionInput, pub string) string {
	pb, _ := hex.DecodeString(pub)
	messages := make([]*types.Any, 0)
	if input.ContractArray != nil {
		for _, msg := range input.ContractArray {
			anyMsg, _ := types.NewAnyWithValueAndName(&msg)
			messages = append(messages, anyMsg)
		}
	}
	body := tx.TxBody{Messages: messages, Memo: input.Memo, TimeoutHeight: 0}
	pubkey := types.PubKey{Key: pb}
	anyPubkey, _ := types.NewAnyWithValue(&pubkey)
	single := tx.ModeInfo_Single{Mode: types.SignMode_SIGN_MODE_DIRECT}
	single_ := tx.ModeInfo_Single_{Single: &single}
	modeInfo := tx.ModeInfo{Sum: &single_}
	signerInfo := make([]*tx.SignerInfo, 0)
	signerInfo = append(signerInfo, &tx.SignerInfo{PublicKey: anyPubkey, ModeInfo: &modeInfo, Sequence: input.Sequence})
	fee := tx.Fee{Amount: input.Fee, GasLimit: input.GasLimit}
	authInfo := tx.AuthInfo{SignerInfos: signerInfo, Fee: &fee}
	bodyBytes, _ := body.Marshal()
	authInfoBytes, _ := authInfo.Marshal()
	signDoc := tx.SignDoc{BodyBytes: bodyBytes, AuthInfoBytes: authInfoBytes, ChainId: input.ChainId, AccountNumber: input.AccountNumber}
	signDocBtyes, _ := signDoc.Marshal()
	return hex.EncodeToString(signDocBtyes)
}
