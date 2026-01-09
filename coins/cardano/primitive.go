package cardano

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"

	"github.com/okx/go-wallet-sdk/crypto/cbor"
)

type Network byte

const (
	Testnet Network = 0
	Mainnet Network = 1
)

// String implements Stringer.
func (n Network) String() string {
	if n == Mainnet {
		return "mainnet"
	} else {
		return "testnet"
	}
}

type BigNum uint64

// Coin represents the Cardano Native Token, in Lovelace.
type Coin BigNum

// Value is a bundle of transferable Cardano Native Tokens.
type Value struct {
	Coin       Coin
	MultiAsset *MultiAsset
}

// NewValue returns a new only-coin Value.
func NewValue(coin Coin) *Value {
	return &Value{Coin: coin, MultiAsset: NewMultiAsset()}
}

// NewValueWithAssets returns a new MultiAsset Value.
func NewValueWithAssets(coin Coin, assets *MultiAsset) *Value {
	return &Value{Coin: coin, MultiAsset: assets}
}

// OnlyCoin returns true if the Value only holds coins.
func (v *Value) OnlyCoin() bool {
	return v.MultiAsset == nil || len(v.MultiAsset.m) == 0
}

// IsZero returns true if the Value is zero.
func (v *Value) IsZero() bool {
	for _, assets := range v.MultiAsset.m {
		for _, value := range assets.m {
			if value != 0 {
				return false
			}
		}
	}
	return v.Coin == 0
}

// Add computes the addition of two Values and returns the result.
func (v *Value) Add(rhs *Value) *Value {
	coin := v.Coin + rhs.Coin
	result := NewValue(coin)

	for _, ma := range []*MultiAsset{v.MultiAsset, rhs.MultiAsset} {
		for policy, assets := range ma.m {
			for assetName, value := range assets.m {
				current := result.MultiAsset
				if _, policyExists := current.m[policy]; policyExists {
					if _, assetExists := result.MultiAsset.m[policy].m[assetName]; assetExists {
						current.m[policy].m[assetName] += value
					} else {
						current.m[policy].m[assetName] = value
					}
				} else {
					current.m[policy] = &Assets{
						m: map[cbor.ByteString]BigNum{
							assetName: value,
						},
					}
				}
			}
		}
	}

	return result
}

// Sub computes the substracion of two Values and returns the result.
func (v *Value) Sub(rhs *Value) *Value {
	var coin Coin
	if v.Coin > rhs.Coin {
		coin = v.Coin - rhs.Coin
	}

	result := NewValue(coin)
	for policy, assets := range v.MultiAsset.m {
		reAssets := NewAssets()
		for assetName, value := range assets.m {
			reAssets.m[assetName] = value
		}
		result.MultiAsset.m[policy] = reAssets
	}

	current := result.MultiAsset
	for policy, assets := range rhs.MultiAsset.m {
		for assetName, value := range assets.m {
			if _, policyExists := current.m[policy]; policyExists {
				if _, assetExists := result.MultiAsset.m[policy].m[assetName]; assetExists {
					lastValue := current.m[policy].m[assetName]
					if lastValue > value {
						current.m[policy].m[assetName] -= value
					} else {
						delete(current.m[policy].m, assetName)
					}
				}
			}
		}
		// Check if the policy exists before checking if it's empty
		if val, ok := current.m[policy]; ok {
			if len(val.m) == 0 {
				delete(current.m, policy)
			}
		}
	}

	return result
}

// Compares two Values and returns
//
//	-1 if v < rhs
//	 0 if v == rhs
//	 1 if v > rhs
//	 2 if not comparable
func (v *Value) Cmp(rhs *Value) int {
	lrZero := v.Sub(rhs).IsZero()
	rlZero := rhs.Sub(v).IsZero()

	if !lrZero && !rlZero {
		return 2
	} else if lrZero && !rlZero {
		return -1
	} else if !lrZero && rlZero {
		return 1
	} else {
		return 0
	}
}

// MarshalCBOR implements cbor.Marshaler.
func (v *Value) MarshalCBOR() ([]byte, error) {
	if v.OnlyCoin() {
		return cborEnc.Marshal(v.Coin)
	} else {
		return cborEnc.Marshal([]interface{}{v.Coin, v.MultiAsset})
	}
}

// UnmarshalCBOR implements cbor.Unmarshaler.
func (v *Value) UnmarshalCBOR(data []byte) error {
	type arrayValue struct {
		_          struct{} `cbor:"_,toarray"`
		Coin       Coin
		MultiAsset *MultiAsset
	}

	var coin Coin
	// Try to unmarshal value into a coin
	err := cborDec.Unmarshal(data, &coin)
	if err != nil {
		var arrValue arrayValue
		err := cborDec.Unmarshal(data, &arrValue)
		if err != nil {
			return err
		}
		v.Coin = arrValue.Coin
		v.MultiAsset = arrValue.MultiAsset
	} else {
		v.Coin = coin
		v.MultiAsset = NewMultiAsset()
	}

	return nil
}

// PolicyID is the native token policy id.
type PolicyID struct {
	bs cbor.ByteString
}

// NewPolicyID returns a new PolicyID using a native script.
func NewPolicyID(script NativeScript) (PolicyID, error) {
	scriptHash, err := script.Hash()
	if err != nil {
		return PolicyID{}, err
	}
	return PolicyID{bs: cbor.NewByteString(scriptHash)}, nil
}

// NewPolicyIDFromHash returns a new PolicyID using a script hash.
func NewPolicyIDFromHash(scriptHash Hash28) PolicyID {
	return PolicyID{bs: cbor.NewByteString(scriptHash)}
}

// NewPolicyIDFromHex returns a new PolicyID using a script hash hex.
func NewPolicyIDFromHex(scriptHashHex string) (PolicyID, error) {
	scriptHash, err := hex.DecodeString(scriptHashHex)
	if err != nil {
		return PolicyID{}, err
	}
	return PolicyID{bs: cbor.NewByteString(scriptHash)}, nil
}

// Bytes returns underlying script hash.
func (p *PolicyID) Bytes() []byte {
	return p.bs.Bytes()
}

// String implements Stringer.
func (p *PolicyID) String() string {
	return p.bs.String()
}

// AssetName represents an Asset name.
type AssetName struct {
	bs cbor.ByteString
}

// NewAssetName returns a new AssetName.
func NewAssetName(name string) AssetName {
	return AssetName{bs: cbor.NewByteString([]byte(name))}
}

// NewAssetNameFromHex returns a new AssetName using a hex encoded name.
func NewAssetNameFromHex(hexName string) (AssetName, error) {
	bytes, err := hex.DecodeString(hexName)
	if err != nil {
		return AssetName{}, err
	}
	return AssetName{bs: cbor.NewByteString(bytes)}, nil
}

// Bytes returns the underlying name bytes.
func (an *AssetName) Bytes() []byte {
	return an.bs.Bytes()
}

// String implements Stringer.
func (an AssetName) String() string {
	return string(an.bs.Bytes())
}

// Assets repressents a set of Cardano Native Tokens.
type Assets struct {
	m map[cbor.ByteString]BigNum
}

// NewAssets returns a new empty Assets.
func NewAssets() *Assets {
	return &Assets{m: make(map[cbor.ByteString]BigNum)}
}

// Set sets the value of a given Asset in Assets.
func (a *Assets) Set(name AssetName, val BigNum) *Assets {
	a.m[name.bs] = val
	return a
}

// Get returns the value of a given Asset in Assets.
func (a *Assets) Get(name AssetName) BigNum {
	return a.m[name.bs]
}

// Keys returns all the AssetNames stored in Assets.
func (a *Assets) Keys() []AssetName {
	assetNames := []AssetName{}
	for k := range a.m {
		assetNames = append(assetNames, AssetName{bs: k})
	}
	return assetNames
}

// MarshalCBOR implements cbor.Marshaler
func (a *Assets) MarshalCBOR() ([]byte, error) {
	return cborEnc.Marshal(a.m)
}

// UnmarshalCBOR implements cbor.Unmarshaler.
func (a *Assets) UnmarshalCBOR(data []byte) error {
	return cborDec.Unmarshal(data, &a.m)
}

// MultiAsset is a bundle of Assets indexed by Policy.
type MultiAsset struct {
	m map[cbor.ByteString]*Assets
}

// NewMultiAsset returns a new empty MultiAsset.
func NewMultiAsset() *MultiAsset {
	return &MultiAsset{m: make(map[cbor.ByteString]*Assets)}
}

// Set sets the Assets of a given Policy in MultiAsset.
func (ma *MultiAsset) Set(policyID PolicyID, assets *Assets) *MultiAsset {
	ma.m[policyID.bs] = assets
	return ma
}

// Get returns the Assets of a given Policy in MultiAsset.
func (ma *MultiAsset) Get(policyID PolicyID) *Assets {
	return ma.m[policyID.bs]
}

// Keys returns all the Policies stored in MultiAsset.
func (ma *MultiAsset) Keys() []PolicyID {
	policyIDs := []PolicyID{}
	for id := range ma.m {
		policyIDs = append(policyIDs, PolicyID{bs: id})
	}
	return policyIDs
}

// String implements Stringer.
func (ma MultiAsset) String() string {
	vMap := map[string]uint64{}
	for _, pool := range ma.Keys() {
		for _, assets := range ma.Get(pool).Keys() {
			vMap[assets.String()] = uint64(ma.Get(pool).Get(assets))
		}
	}
	return fmt.Sprintf("%+v", vMap)
}

func (ma *MultiAsset) numPIDs() uint {
	return uint(len(ma.m))
}

func (ma *MultiAsset) numAssets() uint {
	var num uint
	for _, assets := range ma.m {
		num += uint(len(assets.m))
	}
	return num
}

func (ma *MultiAsset) assetsLength() uint {
	var sum uint
	for _, assets := range ma.m {
		for assetName := range assets.m {
			sum += uint(len(assetName.Bytes()))
		}
	}
	return sum
}

// MarshalCBOR implements cbor.Marshaler
func (ma *MultiAsset) MarshalCBOR() ([]byte, error) {
	return cborEnc.Marshal(ma.m)
}

// UnmarshalCBOR implements cbor.Unmarshaler.
func (ma *MultiAsset) UnmarshalCBOR(data []byte) error {
	return cborDec.Unmarshal(data, &ma.m)
}

// MintAssets represents a set of Cardano Native Tokens to be minted.
type MintAssets struct {
	m map[cbor.ByteString]*big.Int
}

// NewMintAssets returns a new empty MintAssets.
func NewMintAssets() *MintAssets {
	return &MintAssets{m: make(map[cbor.ByteString]*big.Int)}
}

// Set sets the value of a given Asset in MintAssets.
func (a *MintAssets) Set(name AssetName, val *big.Int) *MintAssets {
	a.m[name.bs] = val
	return a
}

// Get returns the value of a given Asset in MintAssets.
func (a *MintAssets) Get(name AssetName) *big.Int {
	return a.m[name.bs]
}

// Keys returns all the AssetNames stored in MintAssets.
func (a *MintAssets) Keys() []AssetName {
	assetNames := []AssetName{}
	for k := range a.m {
		assetNames = append(assetNames, AssetName{bs: k})
	}
	return assetNames
}

// MarshalCBOR implements cbor.Marshaler
func (a *MintAssets) MarshalCBOR() ([]byte, error) {
	return cborEnc.Marshal(a.m)
}

// UnmarshalCBOR implements cbor.Unmarshaler.
func (a *MintAssets) UnmarshalCBOR(data []byte) error {
	return cborDec.Unmarshal(data, &a.m)
}

// Mint is a bundle of MintAssets indexed by Policy.
type Mint struct {
	m map[cbor.ByteString]*MintAssets
}

// NewMint returns a new empty Mint.
func NewMint() *Mint {
	return &Mint{m: make(map[cbor.ByteString]*MintAssets)}
}

// Set sets the MintAssets of a given Policy in Mint.
func (m *Mint) Set(policyID PolicyID, assets *MintAssets) *Mint {
	m.m[policyID.bs] = assets
	return m
}

// Get returns the MintAssets of a given Policy in Mint.
func (m *Mint) Get(policyID PolicyID) *MintAssets {
	return m.m[policyID.bs]
}

// Keys returns all the Policies stored in Mint.
func (m *Mint) Keys() []PolicyID {
	policyIDs := []PolicyID{}
	for id := range m.m {
		policyIDs = append(policyIDs, PolicyID{bs: id})
	}
	return policyIDs
}

// MultiAsset returns a new MultiAsset created from Mint.
func (m *Mint) MultiAsset() *MultiAsset {
	ma := NewMultiAsset()
	for policy, mintAssets := range m.m {
		assets := NewAssets()
		for assetName, value := range mintAssets.m {
			posVal := value.Abs(value)
			if posVal.IsUint64() {
				assets.m[assetName] = BigNum(posVal.Uint64())
			} else {
				panic("MintAsset value cannot be represented as a uint64")
			}
		}
		ma.m[policy] = assets
	}
	return ma
}

func (ma *Mint) numPIDs() uint {
	return uint(len(ma.m))
}

func (ma *Mint) numAssets() uint {
	var num uint
	for _, assets := range ma.m {
		num += uint(len(assets.m))
	}
	return num
}

func (ma *Mint) assetsLength() uint {
	var sum uint
	for _, assets := range ma.m {
		for assetName := range assets.m {
			sum += uint(len(assetName.Bytes()))
		}
	}
	return sum
}

// MarshalCBOR implements cbor.Marshaler
func (ma *Mint) MarshalCBOR() ([]byte, error) {
	return cborEnc.Marshal(ma.m)
}

// UnmarshalCBOR implements cbor.Unmarshaler.
func (ma *Mint) UnmarshalCBOR(data []byte) error {
	return cborDec.Unmarshal(data, &ma.m)
}

type AddrKeyHash = Hash28

type PoolKeyHash = Hash28

type Hash28 []byte

// NewHash28 returns a new Hash28 from a hex encoded string.
func NewHash28(h string) (Hash28, error) {
	hash := make([]byte, 28)
	b, err := hex.DecodeString(h)
	if err != nil {
		return hash, err
	}
	copy(hash[:], b)
	return hash, nil
}

// String returns the hex encoding representation of a Hash28.
func (h Hash28) String() string {
	return hex.EncodeToString(h[:])
}

type Hash32 []byte

// NewHash32 returns a new Hash32 from a hex encoded string.
func NewHash32(h string) (Hash32, error) {
	hash := make([]byte, 32)
	b, err := hex.DecodeString(h)
	if err != nil {
		return hash, err
	}
	copy(hash[:], b)
	return hash, nil
}

// String returns the hex encoding representation of a Hash32
func (h Hash32) String() string {
	return hex.EncodeToString(h[:])
}

type Uint64 *uint64

func NewUint64(u uint64) Uint64 {
	return Uint64(&u)
}

type String *string

func NewString(s string) String {
	return String(&s)
}

type UnitInterval = Rational

type Rational struct {
	_ struct{} `cbor:",toarray"`
	P uint64
	Q uint64
}

// MarshalCBOR implements cbor.Marshaler
func (r *Rational) MarshalCBOR() ([]byte, error) {
	type rational Rational

	// Register tag 30 for rational numbers
	tags, err := r.tagSet(rational{})
	if err != nil {
		return nil, err
	}

	em, err := cbor.CanonicalEncOptions().EncModeWithTags(tags)
	if err != nil {
		return nil, err
	}

	return em.Marshal(rational(*r))
}

// UnmarshalCBOR implements cbor.Unmarshaler
func (r *Rational) UnmarshalCBOR(data []byte) error {
	type rational Rational

	// Register tag 30 for rational numbers
	tags, err := r.tagSet(rational{})
	if err != nil {
		return err
	}

	dm, err := cbor.DecOptions{}.DecModeWithTags(tags)
	if err != nil {
		return err
	}

	var rr rational
	if err := dm.Unmarshal(data, &rr); err != nil {
		return err
	}
	r.P = rr.P
	r.Q = rr.Q

	return nil
}

func (r *Rational) tagSet(contentType interface{}) (cbor.TagSet, error) {
	tags := cbor.NewTagSet()
	err := tags.Add(
		cbor.TagOptions{EncTag: cbor.EncTagRequired, DecTag: cbor.DecTagRequired},
		reflect.TypeOf(contentType),
		30)
	if err != nil {
		return nil, err
	}

	return tags, nil
}
