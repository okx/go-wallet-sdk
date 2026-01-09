package cardano

import (
	"errors"
	"math/big"

	"github.com/okx/go-wallet-sdk/crypto/bech32"
	"github.com/okx/go-wallet-sdk/crypto/cbor"
	"golang.org/x/crypto/blake2b"
)

type AddressType byte

const (
	Base       AddressType = 0x00
	Ptr        AddressType = 0x04
	Enterprise AddressType = 0x06
)

// Address represents a Cardano address.
type Address struct {
	Network Network
	Type    AddressType
	Pointer Pointer

	Payment StakeCredential
	Stake   StakeCredential
}

// NewAddress creates an Address from a bech32 encoded string.
func NewAddress(bech string) (Address, error) {
	hrp, data, err := bech32.DecodeNoLimit(bech)
	if hrp != getHrp(Mainnet) {
		return Address{}, errors.New("invalid bech32 address")
	}
	if err != nil {
		return Address{}, err
	}
	bytes, err := bech32.ConvertBits(data, 5, 8, false)
	if err != nil {
		return Address{}, err
	}
	return NewAddressFromBytes(bytes)
}

// NewAddressFromBytes creates an Address from bytes.
func NewAddressFromBytes(bytes []byte) (Address, error) {
	addr := Address{
		Type:    AddressType(bytes[0] >> 4),
		Network: Network(bytes[0] & 0x01),
	}

	switch addr.Type {
	case Base:
		if len(bytes) != 57 {
			return addr, errors.New("base address length should be 29")
		}
		addr.Payment = StakeCredential{
			Type:    KeyCredential,
			KeyHash: bytes[1:29],
		}
		addr.Stake = StakeCredential{
			Type:    KeyCredential,
			KeyHash: bytes[29:57],
		}
	case Base + 1:
		if len(bytes) != 57 {
			return addr, errors.New("base address length should be 29")
		}
		addr.Payment = StakeCredential{
			Type:       ScriptCredential,
			ScriptHash: bytes[1:29],
		}
		addr.Stake = StakeCredential{
			Type:    KeyCredential,
			KeyHash: bytes[29:57],
		}
	case Base + 2:
		if len(bytes) != 57 {
			return addr, errors.New("base address length should be 29")
		}
		addr.Payment = StakeCredential{
			Type:    KeyCredential,
			KeyHash: bytes[1:29],
		}
		addr.Stake = StakeCredential{
			Type:       ScriptCredential,
			ScriptHash: bytes[29:57],
		}
	case Base + 3:
		if len(bytes) != 57 {
			return addr, errors.New("base address length should be 29")
		}
		addr.Payment = StakeCredential{
			Type:       ScriptCredential,
			ScriptHash: bytes[1:29],
		}
		addr.Stake = StakeCredential{
			Type:       ScriptCredential,
			ScriptHash: bytes[29:57],
		}
	case Ptr:
		if len(bytes) <= 29 {
			return addr, errors.New("enterprise address length should be greater than 29")
		}

		index := uint(29)
		slot, sn, err := decodeFromNat(bytes[29:])
		if err != nil {
			return addr, err
		}
		index += sn
		txIndex, tn, err := decodeFromNat(bytes[index:])
		if err != nil {
			return addr, err
		}
		index += tn
		certIndex, _, err := decodeFromNat(bytes[index:])
		if err != nil {
			return addr, err
		}

		addr.Payment = StakeCredential{
			Type:    KeyCredential,
			KeyHash: bytes[1:29],
		}
		addr.Pointer = Pointer{Slot: slot, TxIndex: txIndex, CertIndex: certIndex}
	case Ptr + 1:
		if len(bytes) <= 29 {
			return addr, errors.New("enterprise address length should be greater than 29")
		}

		index := uint(29)
		slot, sn, err := decodeFromNat(bytes[29:])
		if err != nil {
			return addr, err
		}
		index += sn
		txIndex, tn, err := decodeFromNat(bytes[index:])
		if err != nil {
			return addr, err
		}
		index += tn
		certIndex, _, err := decodeFromNat(bytes[index:])
		if err != nil {
			return addr, err
		}

		addr.Payment = StakeCredential{
			Type:       ScriptCredential,
			ScriptHash: bytes[1:29],
		}
		addr.Pointer = Pointer{Slot: slot, TxIndex: txIndex, CertIndex: certIndex}
	case Enterprise:
		if len(bytes) != 29 {
			return addr, errors.New("enterprise address length should be 29")
		}
		addr.Payment = StakeCredential{
			Type:    KeyCredential,
			KeyHash: bytes[1:29],
		}
	case Enterprise + 1:
		if len(bytes) != 29 {
			return addr, errors.New("enterprise address length should be 29")
		}
		addr.Payment = StakeCredential{
			Type:       ScriptCredential,
			ScriptHash: bytes[1:29],
		}
	}

	return addr, nil
}

// MarshalCBOR implements cbor.Marshaler.
func (addr *Address) MarshalCBOR() ([]byte, error) {
	em, _ := cbor.CanonicalEncOptions().EncMode()
	return em.Marshal(addr.Bytes())
}

// UnmarshalCBOR implements cbor.Unmarshaler.
func (addr *Address) UnmarshalCBOR(data []byte) error {
	bytes := []byte{}
	if err := cborDec.Unmarshal(data, &bytes); err != nil {
		return nil
	}
	decoded, err := NewAddressFromBytes(bytes)
	if err != nil {
		return err
	}

	addr.Network = decoded.Network
	addr.Type = decoded.Type
	addr.Payment = decoded.Payment
	addr.Stake = decoded.Stake
	addr.Pointer = decoded.Pointer

	return nil
}

// Bytes returns the CBOR encoding of the Address as bytes.
func (addr *Address) Bytes() []byte {
	addrBytes := []byte{byte(addr.Type<<4) | (byte(addr.Network) & 0xFF)}
	switch addr.Type {
	case Base, Base + 1, Base + 2, Base + 3:
		addrBytes = append(addrBytes, addr.Payment.Hash()...)
		addrBytes = append(addrBytes, addr.Stake.Hash()...)
	case Enterprise, Enterprise + 1:
		addrBytes = append(addrBytes, addr.Payment.Hash()...)
	case Ptr, Ptr + 1:
		addrBytes = append(addrBytes, addr.Payment.Hash()...)
		addrBytes = append(addrBytes, encodeToNat(addr.Pointer.Slot)...)
		addrBytes = append(addrBytes, encodeToNat(addr.Pointer.TxIndex)...)
		addrBytes = append(addrBytes, encodeToNat(addr.Pointer.CertIndex)...)
	}

	return addrBytes
}

// Bech32 returns the Address encoded as bech32.
func (addr *Address) Bech32() string {
	addrStr, err := bech32.EncodeFromBase256(getHrp(addr.Network), addr.Bytes())
	if err != nil {
		panic(err)
	}
	return addrStr
}

// String returns the Address encoded as bech32.
func (addr Address) String() string {
	return addr.Bech32()
}

// NewBaseAddress returns a new Base Address.
func NewBaseAddress(network Network, payment StakeCredential, stake StakeCredential) (Address, error) {
	addrType := Base
	if payment.Type == ScriptCredential && stake.Type == KeyCredential {
		addrType = Base + 1
	} else if payment.Type == KeyCredential && stake.Type == ScriptCredential {
		addrType = Base + 2
	} else if payment.Type == ScriptCredential && stake.Type == ScriptCredential {
		addrType = Base + 3
	}
	return Address{Type: addrType, Network: network, Payment: payment, Stake: stake}, nil
}

// NewEnterpriseAddress returns a new Enterprise Address.
func NewEnterpriseAddress(network Network, payment StakeCredential) (Address, error) {
	addrType := Enterprise
	if payment.Type == ScriptCredential {
		addrType = Enterprise + 1
	}
	return Address{Type: addrType, Network: network, Payment: payment}, nil
}

// Pointer is the location of the Stake Registration Certificate in the blockchain.
type Pointer struct {
	Slot      uint64
	TxIndex   uint64
	CertIndex uint64
}

// NewPointerAddress returns a new Pointer Address.
func NewPointerAddress(network Network, payment StakeCredential, ptr Pointer) (Address, error) {
	addrType := Ptr
	if payment.Type == ScriptCredential {
		addrType = Ptr + 1
	}
	return Address{Type: addrType, Network: network, Payment: payment, Pointer: ptr}, nil
}

func decodeFromNat(data []byte) (uint64, uint, error) {
	out := big.NewInt(0)
	n := uint(0)
	for _, b := range data {
		out.Lsh(out, 7)
		out.Or(out, big.NewInt(int64(b&0x7F)))
		if !out.IsUint64() {
			return 0, 0, errors.New("too big to decode (> math.MaxUint64)")
		}
		n += 1
		if b&0x80 == 0 {
			return out.Uint64(), n, nil
		}
	}
	return 0, 0, errors.New("bad nat encoding")
}

func encodeToNat(n uint64) []byte {
	out := []byte{byte(n) & 0x7F}

	n >>= 7
	for n != 0 {
		out = append(out, byte(n)|0x80)
		n >>= 7
	}

	// reverse
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}

	return out
}

func Blake224Hash(b []byte) ([]byte, error) {
	hash, err := blake2b.New(224/8, nil)
	if err != nil {
		return nil, err
	}
	_, err = hash.Write(b)
	if err != nil {
		return nil, err
	}
	return hash.Sum(nil), err
}

func getHrp(network Network) string {
	if network == Testnet {
		return "addr_test"
	} else {
		return "addr"
	}
}
