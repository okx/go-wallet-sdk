package filecoin

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/dchest/blake2b"
	"github.com/fxamacker/cbor"
	"math/big"
)

type Message struct {
	Version    uint64  `json:"Version"`
	To         string  `json:"To"`
	From       string  `json:"From"`
	Nonce      uint64  `json:"Nonce"`
	Value      *BigInt `json:"Value"`
	GasLimit   int64   `json:"GasLimit"`
	GasFeeCap  *BigInt `json:"GasFeeCap"`
	GasPremium *BigInt `json:"GasPremium"`
	Method     uint64  `json:"Method"`
	Params     []byte  `json:"Params" cbor:"omitempty"`
}

type SignedMessage struct {
	Message   *Message  `json:"Message"`
	Signature Signature `json:"Signature"`
}

type Signature struct {
	Type byte
	Data []byte
}

type BigInt big.Int

func (bn BigInt) MarshalJSON() ([]byte, error) {
	b := big.Int(bn)
	return json.Marshal(b.String())
}

func (bn *BigInt) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	bnInt, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return errors.New("not number")
	}
	*bn = BigInt(*bnInt)
	return nil
}
func (bn BigInt) Bytes() []byte {
	b := big.Int(bn)
	return b.Bytes()
}

func (m *Message) Serialize() []byte {

	i := []interface{}{
		0,
		AddressToBytes(m.To),
		AddressToBytes(m.From),
		m.Nonce,
		append([]byte{0}, m.Value.Bytes()...),
		m.GasLimit,
		append([]byte{0}, m.GasFeeCap.Bytes()...),
		append([]byte{0}, m.GasPremium.Bytes()...),
		m.Method,
		append([]byte{}, m.Params...),
	}
	bytes, err := cbor.Marshal(i, cbor.CanonicalEncOptions())
	if err != nil {
		return []byte{}
	}
	return bytes
}

func (m *Message) Hash() []byte {

	bytes := m.Serialize()
	h, _ := blake2b.New(&blake2b.Config{Size: uint8(32)})
	h.Write(bytes)
	sum := h.Sum(nil)
	prefix := []byte{0x01, 0x71, 0xa0, 0xe4, 0x02, 0x20}
	cid := append(prefix, sum...)

	h, _ = blake2b.New(&blake2b.Config{Size: uint8(32)})
	h.Write(cid)
	sum = h.Sum(nil)

	return sum
}

func (sm *SignedMessage) Serialize() []byte {
	m := sm.Message
	sig := sm.Signature
	i := []interface{}{
		[]interface{}{
			0,
			AddressToBytes(m.To),
			AddressToBytes(m.From),
			m.Nonce,
			append([]byte{0}, m.Value.Bytes()...),
			m.GasLimit,
			append([]byte{0}, m.GasFeeCap.Bytes()...),
			append([]byte{0}, m.GasPremium.Bytes()...),
			m.Method,
			append([]byte{}, m.Params...),
		},
		append([]byte{SECP256K1}, sig.Data...), // Signature
	}
	bytes, err := cbor.Marshal(i, cbor.CanonicalEncOptions())
	if err != nil {
		return []byte{}
	}

	return bytes
}

func (sm *SignedMessage) Hash() string {
	bytes := sm.Serialize()
	h, _ := blake2b.New(&blake2b.Config{Size: uint8(32)})
	h.Write(bytes)
	sum := h.Sum(nil)
	prefix := []byte{0x01, 0x71, 0xa0, 0xe4, 0x02, 0x20}
	cid := append(prefix, sum...)
	return "b" + AddressEncoding.WithPadding(-1).EncodeToString(cid)
}

func CalTxHash(tx string) (string, error) {
	var msg SignedMessage
	err := json.Unmarshal([]byte(tx), &msg)
	if err != nil {
		return "", err
	}
	return msg.Hash(), nil
}
func NewTx(from, to string, nonce, method, gasLimit int, value, gasFeeCap, gasPremium *big.Int) *Message {
	bvalue := BigInt(*value)
	bgasFeeCap := BigInt(*gasFeeCap)
	bgasPremium := BigInt(*gasPremium)
	return &Message{
		Version:    0,
		To:         to,
		From:       from,
		Nonce:      uint64(nonce),
		Value:      &bvalue,
		GasLimit:   int64(gasLimit),
		GasFeeCap:  &bgasFeeCap,
		GasPremium: &bgasPremium,
		Method:     uint64(method),
		Params:     []byte{},
	}
}

func SignedTx(message *Message, signHex string) (string, error) {
	signData, err := hex.DecodeString(signHex)
	if err != nil {
		return "", err
	}
	V := signData[0]
	R := signData[1:33]
	S := signData[33:65]
	signature := append(R, S...)
	signature = append(signature, V-27)

	tx := &SignedMessage{
		Message: message,
		Signature: Signature{
			Type: SECP256K1,
			Data: signature,
		},
	}

	bytes, err := json.Marshal(tx)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
