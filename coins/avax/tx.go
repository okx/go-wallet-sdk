package avax

import (
	"encoding/hex"
	"sort"
	"strings"
)

const (
	BASETX          = 0
	SECPINPUTID     = 5
	SECPOUTPUTID    = 7
	SECPCREDENTIAL  = 9
	CHAINID_X       = "X"
	HRP_FUJI        = "fuji"
	NETWORK_FUJI    = 5
	BLOCKCHAIN_FUJI = "2JVSBoinj9C2J33VntvzYtVJNZdN2NKiwwKjcumHUWEb5DbBrm"
	ASSET_AVAX_FUJI = "U8iRqJoiJm8xZHAacmvYyZVwqQx6uDNtQeP3CQ6fcgQk3JqnK"
)

type Credential struct {
	Signature    []byte
	CredentialId uint32
}

type outputSlices []*TransferableOutput
type inputSlices []*TransferableInput

type Transaction struct {
	Codecid      uint32
	TypeId       uint32
	NetworkID    uint32
	BlockchainID []byte
	Outs         outputSlices
	Ins          inputSlices
	Memo         [4]byte
	Credentials  []Credential
}

func (t *Transaction) AddInput(txId string, index uint32, amount uint64, assetId string, privateKey string) error {
	txIdBytes, err := CheckDecodeWithCheckSumLast(txId)
	if err != nil {
		return err
	}
	assetIdBytes, err := CheckDecodeWithCheckSumLast(assetId)
	if err != nil {
		return err
	}
	t.Ins = append(t.Ins, &TransferableInput{txIdBytes, index, assetIdBytes, TxInput{SECPINPUTID, amount, 0}, privateKey})
	return nil
}

func (t *Transaction) AddOutput(address string, value uint64, assetId string) error {
	assetIdBytes, err := CheckDecodeWithCheckSumLast(assetId)
	if err != nil {
		return err
	}
	converted, _, _, err := ParseAddress(address)
	if err != nil {
		return err
	}
	t.Outs = append(t.Outs, &TransferableOutput{assetIdBytes, TxOutput{SECPOUTPUTID, value, 0, 1, converted}})
	return nil
}

func (s outputSlices) Len() int {
	return len(s)
}

func (s outputSlices) Less(i, j int) bool {
	sl := NewSerializer()
	s[i].SerializeToBytes(sl)
	hi := hex.EncodeToString(sl.Body[0:sl.Offset])

	sl = NewSerializer()
	s[j].SerializeToBytes(sl)
	hj := hex.EncodeToString(sl.Body[0:sl.Offset])
	return strings.Compare(hi, hj) <= 0
}

func (s outputSlices) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s inputSlices) Len() int {
	return len(s)
}

func (s inputSlices) Less(i, j int) bool {
	sl := NewSerializer()

	sl.WriteFixedBytes(s[i].TxID)
	sl.WriteInt(s[i].UtxoIndex)
	hi := hex.EncodeToString(sl.Payload())

	sl = NewSerializer()
	sl.WriteFixedBytes(s[j].TxID)
	sl.WriteInt(s[j].UtxoIndex)
	hj := hex.EncodeToString(sl.Payload())
	return strings.Compare(hi, hj) <= 0
}

func (s inputSlices) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (t *Transaction) sortInputAndOutput() {
	sort.Sort(t.Ins)
	sort.Sort(t.Outs)
}

func (t *Transaction) SerializeToBytes(sl *Serializer) error {
	sl.WriteShort(uint16(t.Codecid), false)
	sl.WriteInt(t.TypeId)
	sl.WriteInt(t.NetworkID)
	sl.WriteFixedBytes(t.BlockchainID)
	sl.WriteInt(uint32(len(t.Outs)))
	for _, out := range t.Outs {
		if err := out.SerializeToBytes(sl); err != nil {
			return err
		}
	}

	sl.WriteInt(uint32(len(t.Ins)))
	for _, input := range t.Ins {
		if err := input.SerializeToBytes(sl); err != nil {
			return err
		}
	}

	sl.WriteBytes(t.Memo[:])

	if t.Credentials != nil {
		sl.WriteInt(uint32(len(t.Credentials)))
		for _, credential := range t.Credentials {
			sl.WriteInt(credential.CredentialId)
			sl.WriteInt(1)
			sl.WriteFixedBytes(credential.Signature)
		}
	}
	return sl.err
}

type TransferableOutput struct {
	AssetID []byte
	Output  TxOutput
}

func (t *TransferableOutput) SerializeToBytes(sl *Serializer) error {
	sl.WriteFixedBytes(t.AssetID)
	if sl.err != nil {
		return sl.err
	}
	t.Output.SerializeToBytes(sl)
	return sl.err
}

type TxOutput struct {
	TypeID     uint32
	Amount     uint64
	Locktime   uint64
	Threshold  uint32
	PubKeyHash []byte
}

func (t *TxOutput) SerializeToBytes(sl *Serializer) {
	sl.WriteInt(t.TypeID)
	sl.WriteLong(t.Amount)
	sl.WriteLong(t.Locktime)
	sl.WriteInt(t.Threshold)
	sl.WriteInt(1)
	sl.WriteFixedBytes(t.PubKeyHash)
}

type TransferableInput struct {
	TxID       []byte
	UtxoIndex  uint32
	AssetID    []byte
	Input      TxInput
	PrivateKey string
}

func (t *TransferableInput) SerializeToBytes(sl *Serializer) error {
	sl.WriteFixedBytes(t.TxID)
	sl.WriteInt(t.UtxoIndex)
	sl.WriteFixedBytes(t.AssetID)
	t.Input.SerializeToBytes(sl)
	return sl.err
}

type TxInput struct {
	TypeID       uint32
	Amount       uint64
	AddressIndic uint32
}

func (t *TxInput) SerializeToBytes(sl *Serializer) {
	sl.WriteInt(t.TypeID)
	sl.WriteLong(t.Amount)
	sl.WriteInt(1)
	sl.WriteInt(t.AddressIndic)
}

type TransferInput struct {
	TxId       string
	Index      uint32
	Amount     uint64
	AssetId    string
	PrivateKey string
}

type TransferOutPut struct {
	Address string
	Value   uint64
	AssetId string
}
