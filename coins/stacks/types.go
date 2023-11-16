package stacks

import (
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
)

type Message interface {
	GetType() int
}

type ClarityValue interface {
	Message
}

type Payload interface {
	getPayloadType() int
	getRecipient() *StandardPrincipalCV
	getAmount() *big.Int
	getMemo() *Memo
	getContractAddress() *Address
	getContractName() *LengthPrefixedString
	getFunctionName() *LengthPrefixedString
	getFunctionArgs() []ClarityValue
}

type PostConditionInterface interface {
	Message
	getConditionType() int
}

type DecodeBtcAddressBean struct {
	hashMode int
	data     []byte
}

type SignedContractCallOptions struct {
	ContractAddress         string                   `json:"contractAddress"`
	ContractName            string                   `json:"contractName"`
	FunctionName            string                   `json:"functionName"`
	FunctionArgs            []ClarityValue           `json:"functionArgs"`
	SendKey                 string                   `json:"sendKey"`
	ValidateWithAbi         bool                     `json:"validateWithAbi"`
	Fee                     big.Int                  `json:"fee"`
	Nonce                   big.Int                  `json:"nonce"`
	AnchorMode              int                      `json:"anchorMode"`
	PostConditionMode       int                      `json:"postConditionMode"`
	PostConditions          []PostConditionInterface `json:"postConditions"`
	SerializePostConditions []string                 `json:"serializePostConditions"`
}

type StacksTransaction struct {
	Version           int                   `json:"version"`
	ChainId           int64                 `json:"chainId"`
	Auth              StandardAuthorization `json:"auth"`
	AnchorMode        int                   `json:"anchorMode"`
	Payload           Payload               `json:"payload"`
	PostConditionMode int                   `json:"postConditionMode"`
	PostConditions    LPList                `json:"postConditions"`
}

type SignedTokenTransferOptions struct {
	Recipient string  `json:"recipient"`
	Amount    big.Int `json:"amount"`
	Fee       big.Int `json:"fee"`
	Nonce     big.Int `json:"nonce"`
	Memo      string  `json:"memo"`
	SenderKey string  `json:"senderKey"`
}

type SingleSigSpendingCondition struct {
	HashMode    uint64           `json:"hashMode"`
	Signer      string           `json:"signer"`
	Nonce       big.Int          `json:"nonce"`
	Fee         big.Int          `json:"fee"`
	KeyEncoding uint64           `json:"keyEncoding"`
	Signature   MessageSignature `json:"signature"`
}

type StacksPrivateKey struct {
	Data       []byte
	Compressed bool
}

type PostCondition struct {
	StacksMessage
	ConditionType int
	// Principal     PostConditionPrincipal
	Principal     PostConditionPrincipalInterface
	ConditionCode int
}

func (p PostCondition) getConditionType() int {
	return p.ConditionType
}

func (p PostCondition) GetType() int {
	// return p.Type
	return p.ConditionType
}

type AssetInfo struct {
	type_        int
	address      Address
	contractName LengthPrefixedString
	assetName    LengthPrefixedString
}

type PostConditionPrincipal struct {
	Type    int
	Prefix  int
	Address Address
}

func (p PostConditionPrincipal) getPrefix() int {
	return p.Prefix
}

type ContractPrincipal struct {
	PostConditionPrincipal
	contractName LengthPrefixedString
}

func (c ContractPrincipal) getPrefix() int {
	return c.Prefix
}

type PostConditionPrincipalInterface interface {
	getPrefix() int
}

type STXPostCondition struct {
	PostCondition
	amount *big.Int
}

func (s STXPostCondition) getConditionType() int {
	return s.ConditionType
}

func (s STXPostCondition) GetType() int {
	return s.Type
}

type FungiblePostCondition struct {
	PostCondition
	assetInfo AssetInfo
	amount    *big.Int
}

func (f FungiblePostCondition) getConditionType() int {
	return f.ConditionType
}

func (f FungiblePostCondition) GetType() int {
	return f.Type
}

type LPList struct {
	Type              int                      `json:"type"`
	LengthPrefixBytes int                      `json:"lengthPrefixBytes"`
	Values            []PostConditionInterface `json:"values"`
}

type StacksMessage struct {
	Type int `json:"type"`
}

type StandardAuthorization struct {
	AuthType                 uint64                      `json:"authType"`
	SpendingCondition        *SingleSigSpendingCondition `json:"spendingCondition"`
	SponsorSpendingCondition *SingleSigSpendingCondition `json:"sponsorSpendingCondition"`
}

type StacksTransferSig struct {
	Version   uint8  `json:"version"`
	Nonce     uint64 `json:"nonce"`
	Recipient string `json:"recipient"`
	Amount    uint64 `json:"amount"`
	Fee       uint64 `json:"fee"`
	Memo      string `json:"memo"`
}

type LegacyAddress struct {
	Bytes []byte
	P2sh  bool
}

type MessageSignature struct {
	Type_ uint64 `json:"type"`
	Data  string `json:"data"`
}

type StacksPublicKey struct {
	Type_ uint64 `json:"type"`
	Data  string `json:"data"`
}

type Signer struct {
	Type_   uint64 `json:"type"`
	Version uint64 `json:"version"`
	Hash160 string `json:"hash160"`
}

type NextSignature struct {
	nextSig     MessageSignature
	nextSigHash string
}

type TransactionSigner struct {
	transaction   *StacksTransaction
	sigHash       string
	originDone    bool
	checkOversign bool
	checkOverlap  bool
}

func (signer *TransactionSigner) signOrigin(privateKey *StacksPrivateKey) error {
	if signer.checkOverlap && signer.originDone {
		panic("Cannot sign origin after sponsor key")
	}
	//signer.transaction.auth.spendingCondition.signature.data
	nextSig, err := nextSignature(signer.sigHash, 4, &signer.transaction.Auth.SpendingCondition.Fee, &signer.transaction.Auth.SpendingCondition.Nonce, privateKey)
	if err != nil {
		return err
	}
	signer.sigHash = nextSig.nextSigHash
	signer.transaction.Auth.SpendingCondition.Signature = nextSig.nextSig
	return nil
}

type PrincipalCV struct {
	Address Address `json:"address"`
	Type_   uint32  `json:"Type_"`
}

type TupleData struct {
	hashBytes *BufferCV
	version   *BufferCV
}

type BufferCV struct {
	Type_  int    `json:"type"`
	Buffer []byte `json:"buffer"`
}

func (b BufferCV) GetType() int {
	return b.Type_
}

type TupleCV struct {
	Type int                     `json:"type"`
	Data map[string]ClarityValue `json:"data"`
}

func (t TupleCV) GetType() int {
	return t.Type
}

type KeyValuePair struct {
	Name  string
	Value ClarityValue
}

// A function that sorts keys in lexicographic order
func sortByKey(pairs []KeyValuePair) {
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Name < pairs[j].Name
	})
}

type SomeCV struct {
	Type_ int          `json:"type"`
	Value ClarityValue `json:"value"`
}

func (s SomeCV) GetType() int {
	return s.Type_
}

type ResponseCV struct {
	Type_ int          `json:"type"`
	Value ClarityValue `json:"value"`
}

func (r ResponseCV) GetType() int {
	return r.Type_
}

type ListCV struct {
	Type_ int            `json:"type"`
	List  []ClarityValue `json:"list"`
}

func (l ListCV) GetType() int {
	return l.Type_
}

type StringCV struct {
	Type_ int    `json:"type"`
	Data  string `json:"data"`
}

func (s StringCV) GetType() int {
	return s.Type_
}

type NoneCV struct {
	Type_ int
}

func (n NoneCV) GetType() int {
	return n.Type_
}

type UintCV struct {
	Type_ int      `json:"type"`
	Value *big.Int `json:"value"`
}

func (u UintCV) GetType() int {
	return u.Type_
}

func NewUintCV(value *big.Int) *UintCV {
	MaxU128, _ := new(big.Int).SetString("ffffffffffffffffffffffffffffffff", 16)
	MinU128, _ := new(big.Int).SetString("0", 10)
	if value.Cmp(MaxU128) > 0 || value.Cmp(MinU128) < 0 {
		panic("Cannot construct clarity UintCV from the value")
	}
	return &UintCV{Uint, value}
}

type IntCV struct {
	Type_ int      `json:"type"`
	Value *big.Int `json:"value"`
}

func (i IntCV) GetType() int {
	return i.Type_
}

func NewIntCV(value *big.Int) *IntCV {
	MaxI128, _ := new(big.Int).SetString("7fffffffffffffffffffffffffffffff", 16)
	MinI128, _ := new(big.Int).SetString("-170141183460469231731687303715884105728", 10)
	if value.Cmp(MaxI128) > 0 || value.Cmp(MinI128) < 0 {
		panic("Cannot construct clarity IntCV from the value")
	}
	return &IntCV{Int, value}
}

type BooleanCV struct {
	Type_ int `json:"type"`
}

func (b BooleanCV) GetType() int {
	return b.Type_
}

type LengthPrefixedString struct {
	Content           string `json:"content"`
	LengthPrefixBytes int    `json:"lengthPrefixBytes"`
	MaxLengthBytes    int    `json:"maxLengthBytes"`
	Type              int    `json:"type"`
}

func (l LengthPrefixedString) GetType() int {
	return l.Type
}

type Address struct {
	Type_   uint32 `json:"type"`
	Version uint64 `json:"version"`
	Hash160 string `json:"hash160"`
}

func (a Address) GetType() int {
	return int(a.Type_)
}

type StandardPrincipalCV struct {
	Address *Address `json:"address"`
	Type_   int      `json:"type"`
}

func (s StandardPrincipalCV) GetType() int {
	return s.Type_
}

type ContractPrincipalCV struct {
	Type_        int                  `json:"type"`
	Address      Address              `json:"address"`
	ContractName LengthPrefixedString `json:"contractName"`
}

func (c ContractPrincipalCV) GetType() int {
	return c.Type_
}

type Memo struct {
	Type_   uint32 `json:"type"`
	Content string `json:"content"`
}

func (m Memo) GetType() int {
	return int(m.Type_)
}

type TokenTransferPayload struct {
	Type_       int                  `json:"type"`
	PayloadType int                  `json:"payloadType"`
	Recipient   *StandardPrincipalCV `json:"recipient"`
	Amount      big.Int              `json:"amount"`
	Memo        *Memo                `json:"memo"`
}

func (p *TokenTransferPayload) getContractAddress() *Address {
	panic("not have Address")
}

func (p *TokenTransferPayload) getContractName() *LengthPrefixedString {
	panic("no have ContractName")
}

func (p *TokenTransferPayload) getFunctionName() *LengthPrefixedString {
	panic("no have FunctionName")
}

func (p *TokenTransferPayload) getFunctionArgs() []ClarityValue {
	panic("no have FunctionArgs")
}

func (p *TokenTransferPayload) getAmount() *big.Int {
	return &p.Amount
}

func (p *TokenTransferPayload) getMemo() *Memo {
	return p.Memo
}

type ContractCallPayload struct {
	type_           uint32
	payloadType     uint32
	contractAddress *Address
	contractName    *LengthPrefixedString
	functionName    *LengthPrefixedString
	functionArgs    []ClarityValue
}

func (c ContractCallPayload) getContractAddress() *Address {
	return c.contractAddress
}

func (c ContractCallPayload) getContractName() *LengthPrefixedString {
	return c.contractName
}

func (c ContractCallPayload) getFunctionName() *LengthPrefixedString {
	return c.functionName
}

func (c ContractCallPayload) getFunctionArgs() []ClarityValue {
	return c.functionArgs
}

func (c ContractCallPayload) getAmount() *big.Int {
	panic("not have Amount")
}

func (c ContractCallPayload) getMemo() *Memo {
	panic("not have Memo")
}

func (c ContractCallPayload) getRecipient() *StandardPrincipalCV {
	panic("not have Recipient")
}

func (c ContractCallPayload) getPayloadType() int {
	return int(c.payloadType)
}

func (p *TokenTransferPayload) getPayloadType() int {
	return p.PayloadType
}

func (p *TokenTransferPayload) getRecipient() *StandardPrincipalCV {
	return p.Recipient
}

func createLPString(content string, lengthPrefixBytes, maxLengthBytes *int) *LengthPrefixedString {
	if lengthPrefixBytes == nil {
		lengthPrefixBytes = new(int)
		*lengthPrefixBytes = 1
	}
	if maxLengthBytes == nil {
		maxLengthBytes = new(int)
		*maxLengthBytes = 128
	}
	if len(content) > 128 {
		panic("...")
	}

	return &LengthPrefixedString{
		Content:           content,
		LengthPrefixBytes: *lengthPrefixBytes,
		MaxLengthBytes:    *maxLengthBytes,
		Type:              2,
	}
}

func createLPList(lengthPrefixBytes *int, list []PostConditionInterface) *LPList {
	if lengthPrefixBytes == nil {
		defaultLengthPrefixBytes := 4
		lengthPrefixBytes = &defaultLengthPrefixBytes
	}

	lpList := &LPList{}
	lpList.Type = 7
	lpList.LengthPrefixBytes = *lengthPrefixBytes
	lpList.Values = list
	return lpList
}

func createSingleSigSpendingCondition(addressHashMode uint64, pubKey string, nonce big.Int, fee big.Int) (*SingleSigSpendingCondition, error) {
	stacksPublicKey := &StacksPublicKey{}
	stacksPublicKey.Type_ = 6
	stacksPublicKey.Data = pubKey
	signer, err := addressFromPublicKeys(0, addressHashMode, 1, []StacksPublicKey{*stacksPublicKey})
	if err != nil {
		return nil, err
	}
	flag := isCompressed(*stacksPublicKey)
	keyEncoding := 0
	if !flag {
		keyEncoding = 1
	}
	spendingCondition := &SingleSigSpendingCondition{}
	spendingCondition.HashMode = addressHashMode
	spendingCondition.Signer = signer.Hash160
	spendingCondition.Nonce = nonce
	spendingCondition.Fee = fee
	spendingCondition.KeyEncoding = uint64(keyEncoding)
	emptyMessagesignature := MessageSignature{}
	emptyMessagesignature.Type_ = 9
	emptyMessagesignature.Data = "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
	spendingCondition.Signature = emptyMessagesignature
	return spendingCondition, nil

}

func createSignedTokenTransferOptions(recipient string, amount big.Int, fee big.Int, nonce big.Int, memo string, senderKey string) *SignedTokenTransferOptions {
	tx := SignedTokenTransferOptions{
		Recipient: recipient,
		Amount:    amount,
		Fee:       fee,
		Nonce:     nonce,
		Memo:      memo,
		SenderKey: senderKey,
	}
	return &tx
}

func CreateTokenTransferPayload(recipientCV string, amount big.Int, memo string) (*TokenTransferPayload, error) {
	if len(memo) == 0 {
		memo = " "
	}
	memoObj, err := createMemoString(memo)
	if err != nil {
		return nil, err
	}
	standardPrincipalCV := NewStandardPrincipalCV(recipientCV)
	payload := &TokenTransferPayload{}
	payload.Type_ = 8
	payload.PayloadType = 0
	payload.Recipient = standardPrincipalCV
	payload.Amount = amount
	payload.Memo = memoObj
	return payload, nil
}

func NewStandardPrincipalCV(address string) *StandardPrincipalCV {
	addrObj := createAddress(address)
	standardPrincipalCV := &StandardPrincipalCV{}
	standardPrincipalCV.Address = addrObj
	standardPrincipalCV.Type_ = PrincipalStandard
	return standardPrincipalCV
}

func NewContractPrincipalCV(principal string) (*ContractPrincipalCV, error) {
	if !strings.Contains(principal, ".") {
		return nil, fmt.Errorf("not contract address")
	}
	splitPrincipal := strings.Split(principal, ".")
	address := splitPrincipal[0]
	contractName := splitPrincipal[1]

	addr := createAddress(address)
	lengthPrefixedContractName := createLPString(contractName, nil, nil)
	if len(lengthPrefixedContractName.Content) >= 128 {
		return nil, fmt.Errorf("contract name must be less than 128 bytes")
	}
	return &ContractPrincipalCV{
		PrincipalContract,
		*addr,
		*lengthPrefixedContractName,
	}, nil
}

func createAddress(address string) *Address {
	addObj := &Address{}
	strs, err := c32addressDecode(address)
	if err != nil {
		panic(err)
	}
	addObj.Type_ = 0
	version, err := strconv.ParseUint(strs[0], 10, 64)
	if err != nil {
		panic(err)
	}
	addObj.Version = version
	addObj.Hash160 = strs[1]
	return addObj
}

func createMemoString(memo string) (*Memo, error) {
	if len([]rune(memo)) > 34 {
		return nil, fmt.Errorf("stacks memo length max is 34")
	}
	var memoObj_ = &Memo{
		Type_:   3,
		Content: memo,
	}
	return memoObj_, nil
}

type TransactionRes struct {
	TxId        string `json:"txId"`
	TxSerialize string `json:"txSerialize"`
}

func NewUntilBurnBlockHeight(untilBurnBlockHeight *big.Int) ClarityValue {
	if untilBurnBlockHeight == nil {
		return &NoneCV{OptionalNone}
	} else {
		return &SomeCV{OptionalSome, NewUintCV(untilBurnBlockHeight)}
	}
}

type BytesReader struct {
	source   []byte
	consumed int
}

func NewBytesReader(arr []byte) *BytesReader {
	return &BytesReader{
		source:   arr,
		consumed: 0,
	}
}

func (br *BytesReader) ReadBytes(length int) []byte {
	view := br.source[br.consumed : br.consumed+length]
	br.consumed += length
	return view
}

func (br *BytesReader) ReadUInt8() uint8 {
	b := br.ReadBytes(1)
	return b[0]
}

func (br *BytesReader) readUInt32BE() uint32 {
	bytes := br.ReadBytes(4)
	return uint32(bytes[0])<<24 | uint32(bytes[1])<<16 | uint32(bytes[2])<<8 | uint32(bytes[3])
}
