package xdr

import (
	"bytes"
	"encoding"
	"errors"
	"fmt"
	xdr "github.com/okx/go-wallet-sdk/coins/stellar/xdr3"
	"io"
)

var ErrMaxDecodingDepthReached = errors.New("maximum decoding depth reached")

type xdrType interface {
	xdrType()
}

type decoderFrom interface {
	DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error)
}

type TrustLineFlags int32

const (
	TrustLineFlagsAuthorizedFlag                      TrustLineFlags = 1
	TrustLineFlagsAuthorizedToMaintainLiabilitiesFlag TrustLineFlags = 2
	TrustLineFlagsTrustlineClawbackEnabledFlag        TrustLineFlags = 4
)

var trustLineFlagsMap = map[int32]string{
	1: "TrustLineFlagsAuthorizedFlag",
	2: "TrustLineFlagsAuthorizedToMaintainLiabilitiesFlag",
	4: "TrustLineFlagsTrustlineClawbackEnabledFlag",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for TrustLineFlags
func (e TrustLineFlags) ValidEnum(v int32) bool {
	_, ok := trustLineFlagsMap[v]
	return ok
}

// String returns the name of `e`
func (e TrustLineFlags) String() string {
	name, _ := trustLineFlagsMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e TrustLineFlags) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := trustLineFlagsMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid TrustLineFlags enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*TrustLineFlags)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *TrustLineFlags) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding TrustLineFlags: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding TrustLineFlags: %w", err)
	}
	if _, ok := trustLineFlagsMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid TrustLineFlags enum value", v)
	}
	*e = TrustLineFlags(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TrustLineFlags) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TrustLineFlags) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TrustLineFlags)(nil)
	_ encoding.BinaryUnmarshaler = (*TrustLineFlags)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s TrustLineFlags) xdrType() {}

var _ xdrType = (*TrustLineFlags)(nil)

// MaskTrustlineFlags is an XDR Const defines as:
//
//	const MASK_TRUSTLINE_FLAGS = 1;
const MaskTrustlineFlags = 1

// MaskTrustlineFlagsV13 is an XDR Const defines as:
//
//	const MASK_TRUSTLINE_FLAGS_V13 = 3;
const MaskTrustlineFlagsV13 = 3

// MaskTrustlineFlagsV17 is an XDR Const defines as:
//
//	const MASK_TRUSTLINE_FLAGS_V17 = 7;
const MaskTrustlineFlagsV17 = 7

// LiquidityPoolType is an XDR Enum defines as:
//
//	enum LiquidityPoolType
//	 {
//	     LIQUIDITY_POOL_CONSTANT_PRODUCT = 0
//	 };
type LiquidityPoolType int32

const (
	LiquidityPoolTypeLiquidityPoolConstantProduct LiquidityPoolType = 0
)

var liquidityPoolTypeMap = map[int32]string{
	0: "LiquidityPoolTypeLiquidityPoolConstantProduct",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for LiquidityPoolType
func (e LiquidityPoolType) ValidEnum(v int32) bool {
	_, ok := liquidityPoolTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e LiquidityPoolType) String() string {
	name, _ := liquidityPoolTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e LiquidityPoolType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := liquidityPoolTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid LiquidityPoolType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*LiquidityPoolType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *LiquidityPoolType) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding LiquidityPoolType: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding LiquidityPoolType: %w", err)
	}
	if _, ok := liquidityPoolTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid LiquidityPoolType enum value", v)
	}
	*e = LiquidityPoolType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LiquidityPoolType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LiquidityPoolType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LiquidityPoolType)(nil)
	_ encoding.BinaryUnmarshaler = (*LiquidityPoolType)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s LiquidityPoolType) xdrType() {}

var _ xdrType = (*LiquidityPoolType)(nil)

// TrustLineAsset is an XDR Union defines as:
//
//	union TrustLineAsset switch (AssetType type)
//	 {
//	 case ASSET_TYPE_NATIVE: // Not credit
//	     void;
//
//	 case ASSET_TYPE_CREDIT_ALPHANUM4:
//	     AlphaNum4 alphaNum4;
//
//	 case ASSET_TYPE_CREDIT_ALPHANUM12:
//	     AlphaNum12 alphaNum12;
//
//	 case ASSET_TYPE_POOL_SHARE:
//	     PoolID liquidityPoolID;
//
//	     // add other asset types here in the future
//	 };
type TrustLineAsset struct {
	Type            AssetType
	AlphaNum4       *AlphaNum4
	AlphaNum12      *AlphaNum12
	LiquidityPoolId *PoolId
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TrustLineAsset) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TrustLineAsset
func (u TrustLineAsset) ArmForSwitch(sw int32) (string, bool) {
	switch AssetType(sw) {
	case AssetTypeAssetTypeNative:
		return "", true
	case AssetTypeAssetTypeCreditAlphanum4:
		return "AlphaNum4", true
	case AssetTypeAssetTypeCreditAlphanum12:
		return "AlphaNum12", true
	case AssetTypeAssetTypePoolShare:
		return "LiquidityPoolId", true
	}
	return "-", false
}

// NewTrustLineAsset creates a new  TrustLineAsset.
func NewTrustLineAsset(aType AssetType, value interface{}) (result TrustLineAsset, err error) {
	result.Type = aType
	switch AssetType(aType) {
	case AssetTypeAssetTypeNative:
		// void
	case AssetTypeAssetTypeCreditAlphanum4:
		tv, ok := value.(AlphaNum4)
		if !ok {
			err = errors.New("invalid value, must be AlphaNum4")
			return
		}
		result.AlphaNum4 = &tv
	case AssetTypeAssetTypeCreditAlphanum12:
		tv, ok := value.(AlphaNum12)
		if !ok {
			err = errors.New("invalid value, must be AlphaNum12")
			return
		}
		result.AlphaNum12 = &tv
	case AssetTypeAssetTypePoolShare:
		tv, ok := value.(PoolId)
		if !ok {
			err = errors.New("invalid value, must be PoolId")
			return
		}
		result.LiquidityPoolId = &tv
	}
	return
}

// MustAlphaNum4 retrieves the AlphaNum4 value from the union,
// panicing if the value is not set.
func (u TrustLineAsset) MustAlphaNum4() AlphaNum4 {
	val, ok := u.GetAlphaNum4()

	if !ok {
		panic("arm AlphaNum4 is not set")
	}

	return val
}

// GetAlphaNum4 retrieves the AlphaNum4 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TrustLineAsset) GetAlphaNum4() (result AlphaNum4, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AlphaNum4" {
		result = *u.AlphaNum4
		ok = true
	}

	return
}

// MustAlphaNum12 retrieves the AlphaNum12 value from the union,
// panicing if the value is not set.
func (u TrustLineAsset) MustAlphaNum12() AlphaNum12 {
	val, ok := u.GetAlphaNum12()

	if !ok {
		panic("arm AlphaNum12 is not set")
	}

	return val
}

// GetAlphaNum12 retrieves the AlphaNum12 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TrustLineAsset) GetAlphaNum12() (result AlphaNum12, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AlphaNum12" {
		result = *u.AlphaNum12
		ok = true
	}

	return
}

// MustLiquidityPoolId retrieves the LiquidityPoolId value from the union,
// panicing if the value is not set.
func (u TrustLineAsset) MustLiquidityPoolId() PoolId {
	val, ok := u.GetLiquidityPoolId()

	if !ok {
		panic("arm LiquidityPoolId is not set")
	}

	return val
}

// GetLiquidityPoolId retrieves the LiquidityPoolId value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TrustLineAsset) GetLiquidityPoolId() (result PoolId, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "LiquidityPoolId" {
		result = *u.LiquidityPoolId
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u TrustLineAsset) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch AssetType(u.Type) {
	case AssetTypeAssetTypeNative:
		// Void
		return nil
	case AssetTypeAssetTypeCreditAlphanum4:
		if err = (*u.AlphaNum4).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case AssetTypeAssetTypeCreditAlphanum12:
		if err = (*u.AlphaNum12).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case AssetTypeAssetTypePoolShare:
		if err = (*u.LiquidityPoolId).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (AssetType) switch value '%d' is not valid for union TrustLineAsset", u.Type)
}

var _ decoderFrom = (*TrustLineAsset)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *TrustLineAsset) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding TrustLineAsset: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AssetType: %w", err)
	}
	switch AssetType(u.Type) {
	case AssetTypeAssetTypeNative:
		// Void
		return n, nil
	case AssetTypeAssetTypeCreditAlphanum4:
		u.AlphaNum4 = new(AlphaNum4)
		nTmp, err = (*u.AlphaNum4).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AlphaNum4: %w", err)
		}
		return n, nil
	case AssetTypeAssetTypeCreditAlphanum12:
		u.AlphaNum12 = new(AlphaNum12)
		nTmp, err = (*u.AlphaNum12).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AlphaNum12: %w", err)
		}
		return n, nil
	case AssetTypeAssetTypePoolShare:
		u.LiquidityPoolId = new(PoolId)
		nTmp, err = (*u.LiquidityPoolId).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding PoolId: %w", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union TrustLineAsset has invalid Type (AssetType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TrustLineAsset) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TrustLineAsset) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TrustLineAsset)(nil)
	_ encoding.BinaryUnmarshaler = (*TrustLineAsset)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s TrustLineAsset) xdrType() {}

var _ xdrType = (*TrustLineAsset)(nil)

// TrustLineEntryExtensionV2Ext is an XDR NestedUnion defines as:
//
//	union switch (int v)
//	     {
//	     case 0:
//	         void;
//	     }
type TrustLineEntryExtensionV2Ext struct {
	V int32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TrustLineEntryExtensionV2Ext) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TrustLineEntryExtensionV2Ext
func (u TrustLineEntryExtensionV2Ext) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	}
	return "-", false
}

// NewTrustLineEntryExtensionV2Ext creates a new  TrustLineEntryExtensionV2Ext.
func NewTrustLineEntryExtensionV2Ext(v int32, value interface{}) (result TrustLineEntryExtensionV2Ext, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u TrustLineEntryExtensionV2Ext) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union TrustLineEntryExtensionV2Ext", u.V)
}

var _ decoderFrom = (*TrustLineEntryExtensionV2Ext)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *TrustLineEntryExtensionV2Ext) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding TrustLineEntryExtensionV2Ext: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %w", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union TrustLineEntryExtensionV2Ext has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TrustLineEntryExtensionV2Ext) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TrustLineEntryExtensionV2Ext) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TrustLineEntryExtensionV2Ext)(nil)
	_ encoding.BinaryUnmarshaler = (*TrustLineEntryExtensionV2Ext)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s TrustLineEntryExtensionV2Ext) xdrType() {}

var _ xdrType = (*TrustLineEntryExtensionV2Ext)(nil)

// TrustLineEntryExtensionV2 is an XDR Struct defines as:
//
//	struct TrustLineEntryExtensionV2
//	 {
//	     int32 liquidityPoolUseCount;
//
//	     union switch (int v)
//	     {
//	     case 0:
//	         void;
//	     }
//	     ext;
//	 };
type TrustLineEntryExtensionV2 struct {
	LiquidityPoolUseCount Int32
	Ext                   TrustLineEntryExtensionV2Ext
}

// EncodeTo encodes this value using the Encoder.
func (s *TrustLineEntryExtensionV2) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.LiquidityPoolUseCount.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*TrustLineEntryExtensionV2)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TrustLineEntryExtensionV2) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding TrustLineEntryExtensionV2: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.LiquidityPoolUseCount.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int32: %w", err)
	}
	nTmp, err = s.Ext.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TrustLineEntryExtensionV2Ext: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TrustLineEntryExtensionV2) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TrustLineEntryExtensionV2) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TrustLineEntryExtensionV2)(nil)
	_ encoding.BinaryUnmarshaler = (*TrustLineEntryExtensionV2)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s TrustLineEntryExtensionV2) xdrType() {}

var _ xdrType = (*TrustLineEntryExtensionV2)(nil)

// TrustLineEntryV1Ext is an XDR NestedUnion defines as:
//
//	union switch (int v)
//	             {
//	             case 0:
//	                 void;
//	             case 2:
//	                 TrustLineEntryExtensionV2 v2;
//	             }
type TrustLineEntryV1Ext struct {
	V  int32
	V2 *TrustLineEntryExtensionV2
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TrustLineEntryV1Ext) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TrustLineEntryV1Ext
func (u TrustLineEntryV1Ext) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	case 2:
		return "V2", true
	}
	return "-", false
}

// NewTrustLineEntryV1Ext creates a new  TrustLineEntryV1Ext.
func NewTrustLineEntryV1Ext(v int32, value interface{}) (result TrustLineEntryV1Ext, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	case 2:
		tv, ok := value.(TrustLineEntryExtensionV2)
		if !ok {
			err = errors.New("invalid value, must be TrustLineEntryExtensionV2")
			return
		}
		result.V2 = &tv
	}
	return
}

// MustV2 retrieves the V2 value from the union,
// panicing if the value is not set.
func (u TrustLineEntryV1Ext) MustV2() TrustLineEntryExtensionV2 {
	val, ok := u.GetV2()

	if !ok {
		panic("arm V2 is not set")
	}

	return val
}

// GetV2 retrieves the V2 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TrustLineEntryV1Ext) GetV2() (result TrustLineEntryExtensionV2, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "V2" {
		result = *u.V2
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u TrustLineEntryV1Ext) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	case 2:
		if err = (*u.V2).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union TrustLineEntryV1Ext", u.V)
}

var _ decoderFrom = (*TrustLineEntryV1Ext)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *TrustLineEntryV1Ext) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding TrustLineEntryV1Ext: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %w", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	case 2:
		u.V2 = new(TrustLineEntryExtensionV2)
		nTmp, err = (*u.V2).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding TrustLineEntryExtensionV2: %w", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union TrustLineEntryV1Ext has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TrustLineEntryV1Ext) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TrustLineEntryV1Ext) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TrustLineEntryV1Ext)(nil)
	_ encoding.BinaryUnmarshaler = (*TrustLineEntryV1Ext)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s TrustLineEntryV1Ext) xdrType() {}

var _ xdrType = (*TrustLineEntryV1Ext)(nil)

// TrustLineEntryV1 is an XDR NestedStruct defines as:
//
//	struct
//	         {
//	             Liabilities liabilities;
//
//	             union switch (int v)
//	             {
//	             case 0:
//	                 void;
//	             case 2:
//	                 TrustLineEntryExtensionV2 v2;
//	             }
//	             ext;
//	         }
type TrustLineEntryV1 struct {
	Liabilities Liabilities
	Ext         TrustLineEntryV1Ext
}

// EncodeTo encodes this value using the Encoder.
func (s *TrustLineEntryV1) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Liabilities.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*TrustLineEntryV1)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TrustLineEntryV1) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding TrustLineEntryV1: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.Liabilities.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Liabilities: %w", err)
	}
	nTmp, err = s.Ext.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TrustLineEntryV1Ext: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TrustLineEntryV1) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TrustLineEntryV1) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TrustLineEntryV1)(nil)
	_ encoding.BinaryUnmarshaler = (*TrustLineEntryV1)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s TrustLineEntryV1) xdrType() {}

var _ xdrType = (*TrustLineEntryV1)(nil)

type AssetType int32

const (
	AssetTypeAssetTypeNative           AssetType = 0
	AssetTypeAssetTypeCreditAlphanum4  AssetType = 1
	AssetTypeAssetTypeCreditAlphanum12 AssetType = 2
	AssetTypeAssetTypePoolShare        AssetType = 3
)

var assetTypeMap = map[int32]string{
	0: "AssetTypeAssetTypeNative",
	1: "AssetTypeAssetTypeCreditAlphanum4",
	2: "AssetTypeAssetTypeCreditAlphanum12",
	3: "AssetTypeAssetTypePoolShare",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for AssetType
func (e AssetType) ValidEnum(v int32) bool {
	_, ok := assetTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e AssetType) String() string {
	name, _ := assetTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e AssetType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := assetTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid AssetType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*AssetType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *AssetType) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding AssetType: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding AssetType: %w", err)
	}
	if _, ok := assetTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid AssetType enum value", v)
	}
	*e = AssetType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AssetType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AssetType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AssetType)(nil)
	_ encoding.BinaryUnmarshaler = (*AssetType)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s AssetType) xdrType() {}

var _ xdrType = (*AssetType)(nil)

// AssetCode is an XDR Union defines as:
//
//	union AssetCode switch (AssetType type)
//	 {
//	 case ASSET_TYPE_CREDIT_ALPHANUM4:
//	     AssetCode4 assetCode4;
//
//	 case ASSET_TYPE_CREDIT_ALPHANUM12:
//	     AssetCode12 assetCode12;
//
//	     // add other asset types here in the future
//	 };
type AssetCode struct {
	Type        AssetType
	AssetCode4  *AssetCode4
	AssetCode12 *AssetCode12
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AssetCode) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of AssetCode
func (u AssetCode) ArmForSwitch(sw int32) (string, bool) {
	switch AssetType(sw) {
	case AssetTypeAssetTypeCreditAlphanum4:
		return "AssetCode4", true
	case AssetTypeAssetTypeCreditAlphanum12:
		return "AssetCode12", true
	}
	return "-", false
}

// NewAssetCode creates a new  AssetCode.
func NewAssetCode(aType AssetType, value interface{}) (result AssetCode, err error) {
	result.Type = aType
	switch AssetType(aType) {
	case AssetTypeAssetTypeCreditAlphanum4:
		tv, ok := value.(AssetCode4)
		if !ok {
			err = errors.New("invalid value, must be AssetCode4")
			return
		}
		result.AssetCode4 = &tv
	case AssetTypeAssetTypeCreditAlphanum12:
		tv, ok := value.(AssetCode12)
		if !ok {
			err = errors.New("invalid value, must be AssetCode12")
			return
		}
		result.AssetCode12 = &tv
	}
	return
}

// MustAssetCode4 retrieves the AssetCode4 value from the union,
// panicing if the value is not set.
func (u AssetCode) MustAssetCode4() AssetCode4 {
	val, ok := u.GetAssetCode4()

	if !ok {
		panic("arm AssetCode4 is not set")
	}

	return val
}

// GetAssetCode4 retrieves the AssetCode4 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u AssetCode) GetAssetCode4() (result AssetCode4, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AssetCode4" {
		result = *u.AssetCode4
		ok = true
	}

	return
}

// MustAssetCode12 retrieves the AssetCode12 value from the union,
// panicing if the value is not set.
func (u AssetCode) MustAssetCode12() AssetCode12 {
	val, ok := u.GetAssetCode12()

	if !ok {
		panic("arm AssetCode12 is not set")
	}

	return val
}

// GetAssetCode12 retrieves the AssetCode12 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u AssetCode) GetAssetCode12() (result AssetCode12, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AssetCode12" {
		result = *u.AssetCode12
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u AssetCode) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch AssetType(u.Type) {
	case AssetTypeAssetTypeCreditAlphanum4:
		if err = (*u.AssetCode4).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case AssetTypeAssetTypeCreditAlphanum12:
		if err = (*u.AssetCode12).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (AssetType) switch value '%d' is not valid for union AssetCode", u.Type)
}

var _ decoderFrom = (*AssetCode)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *AssetCode) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding AssetCode: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AssetType: %w", err)
	}
	switch AssetType(u.Type) {
	case AssetTypeAssetTypeCreditAlphanum4:
		u.AssetCode4 = new(AssetCode4)
		nTmp, err = (*u.AssetCode4).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AssetCode4: %w", err)
		}
		return n, nil
	case AssetTypeAssetTypeCreditAlphanum12:
		u.AssetCode12 = new(AssetCode12)
		nTmp, err = (*u.AssetCode12).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AssetCode12: %w", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union AssetCode has invalid Type (AssetType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AssetCode) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AssetCode) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AssetCode)(nil)
	_ encoding.BinaryUnmarshaler = (*AssetCode)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s AssetCode) xdrType() {}

var _ xdrType = (*AssetCode)(nil)

// AlphaNum4 is an XDR Struct defines as:
//
//	struct AlphaNum4
//	 {
//	     AssetCode4 assetCode;
//	     AccountID issuer;
//	 };
type AlphaNum4 struct {
	AssetCode AssetCode4
	Issuer    AccountId
}

// EncodeTo encodes this value using the Encoder.
func (s *AlphaNum4) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.AssetCode.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Issuer.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*AlphaNum4)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *AlphaNum4) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding AlphaNum4: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.AssetCode.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AssetCode4: %w", err)
	}
	nTmp, err = s.Issuer.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AlphaNum4) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AlphaNum4) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AlphaNum4)(nil)
	_ encoding.BinaryUnmarshaler = (*AlphaNum4)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s AlphaNum4) xdrType() {}

var _ xdrType = (*AlphaNum4)(nil)

// AlphaNum12 is an XDR Struct defines as:
//
//	struct AlphaNum12
//	 {
//	     AssetCode12 assetCode;
//	     AccountID issuer;
//	 };
type AlphaNum12 struct {
	AssetCode AssetCode12
	Issuer    AccountId
}

// EncodeTo encodes this value using the Encoder.
func (s *AlphaNum12) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.AssetCode.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Issuer.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*AlphaNum12)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *AlphaNum12) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding AlphaNum12: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.AssetCode.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AssetCode12: %w", err)
	}
	nTmp, err = s.Issuer.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AlphaNum12) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AlphaNum12) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AlphaNum12)(nil)
	_ encoding.BinaryUnmarshaler = (*AlphaNum12)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s AlphaNum12) xdrType() {}

var _ xdrType = (*AlphaNum12)(nil)

// Asset is an XDR Union defines as:
//
//	union Asset switch (AssetType type)
//	 {
//	 case ASSET_TYPE_NATIVE: // Not credit
//	     void;
//
//	 case ASSET_TYPE_CREDIT_ALPHANUM4:
//	     AlphaNum4 alphaNum4;
//
//	 case ASSET_TYPE_CREDIT_ALPHANUM12:
//	     AlphaNum12 alphaNum12;
//
//	     // add other asset types here in the future
//	 };
type Asset struct {
	Type       AssetType
	AlphaNum4  *AlphaNum4
	AlphaNum12 *AlphaNum12
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u Asset) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of Asset
func (u Asset) ArmForSwitch(sw int32) (string, bool) {
	switch AssetType(sw) {
	case AssetTypeAssetTypeNative:
		return "", true
	case AssetTypeAssetTypeCreditAlphanum4:
		return "AlphaNum4", true
	case AssetTypeAssetTypeCreditAlphanum12:
		return "AlphaNum12", true
	}
	return "-", false
}

// NewAsset creates a new  Asset.
func NewAsset(aType AssetType, value interface{}) (result Asset, err error) {
	result.Type = aType
	switch AssetType(aType) {
	case AssetTypeAssetTypeNative:
		// void
	case AssetTypeAssetTypeCreditAlphanum4:
		tv, ok := value.(AlphaNum4)
		if !ok {
			err = errors.New("invalid value, must be AlphaNum4")
			return
		}
		result.AlphaNum4 = &tv
	case AssetTypeAssetTypeCreditAlphanum12:
		tv, ok := value.(AlphaNum12)
		if !ok {
			err = errors.New("invalid value, must be AlphaNum12")
			return
		}
		result.AlphaNum12 = &tv
	}
	return
}

// MustAlphaNum4 retrieves the AlphaNum4 value from the union,
// panicing if the value is not set.
func (u Asset) MustAlphaNum4() AlphaNum4 {
	val, ok := u.GetAlphaNum4()

	if !ok {
		panic("arm AlphaNum4 is not set")
	}

	return val
}

// GetAlphaNum4 retrieves the AlphaNum4 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Asset) GetAlphaNum4() (result AlphaNum4, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AlphaNum4" {
		result = *u.AlphaNum4
		ok = true
	}

	return
}

// MustAlphaNum12 retrieves the AlphaNum12 value from the union,
// panicing if the value is not set.
func (u Asset) MustAlphaNum12() AlphaNum12 {
	val, ok := u.GetAlphaNum12()

	if !ok {
		panic("arm AlphaNum12 is not set")
	}

	return val
}

// GetAlphaNum12 retrieves the AlphaNum12 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Asset) GetAlphaNum12() (result AlphaNum12, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AlphaNum12" {
		result = *u.AlphaNum12
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u Asset) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch AssetType(u.Type) {
	case AssetTypeAssetTypeNative:
		// Void
		return nil
	case AssetTypeAssetTypeCreditAlphanum4:
		if err = (*u.AlphaNum4).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case AssetTypeAssetTypeCreditAlphanum12:
		if err = (*u.AlphaNum12).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (AssetType) switch value '%d' is not valid for union Asset", u.Type)
}

var _ decoderFrom = (*Asset)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *Asset) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding Asset: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AssetType: %w", err)
	}
	switch AssetType(u.Type) {
	case AssetTypeAssetTypeNative:
		// Void
		return n, nil
	case AssetTypeAssetTypeCreditAlphanum4:
		u.AlphaNum4 = new(AlphaNum4)
		nTmp, err = (*u.AlphaNum4).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AlphaNum4: %w", err)
		}
		return n, nil
	case AssetTypeAssetTypeCreditAlphanum12:
		u.AlphaNum12 = new(AlphaNum12)
		nTmp, err = (*u.AlphaNum12).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AlphaNum12: %w", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union Asset has invalid Type (AssetType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Asset) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Asset) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Asset)(nil)
	_ encoding.BinaryUnmarshaler = (*Asset)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s Asset) xdrType() {}

var _ xdrType = (*Asset)(nil)

type Hash [32]byte

// XDRMaxSize implements the Sized interface for Hash
func (e Hash) XDRMaxSize() int {
	return 32
}

// EncodeTo encodes this value using the Encoder.
func (s *Hash) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeFixedOpaque(s[:]); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Hash)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Hash) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding Hash: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = d.DecodeFixedOpaqueInplace(s[:])
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Hash: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Hash) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Hash) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Hash)(nil)
	_ encoding.BinaryUnmarshaler = (*Hash)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s Hash) xdrType() {}

var _ xdrType = (*Hash)(nil)

// Uint256 is an XDR Typedef defines as:
//
//	typedef opaque uint256[32];
type Uint256 [32]byte

// XDRMaxSize implements the Sized interface for Uint256
func (e Uint256) XDRMaxSize() int {
	return 32
}

// EncodeTo encodes this value using the Encoder.
func (s *Uint256) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeFixedOpaque(s[:]); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Uint256)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Uint256) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding Uint256: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = d.DecodeFixedOpaqueInplace(s[:])
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint256: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Uint256) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Uint256) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Uint256)(nil)
	_ encoding.BinaryUnmarshaler = (*Uint256)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s Uint256) xdrType() {}

var _ xdrType = (*Uint256)(nil)

// Uint32 is an XDR Typedef defines as:
//
//	typedef unsigned int uint32;
type Uint32 uint32

// EncodeTo encodes this value using the Encoder.
func (s Uint32) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeUint(uint32(s)); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Uint32)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Uint32) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding Uint32: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	var v uint32
	v, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Unsigned int: %w", err)
	}
	*s = Uint32(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Uint32) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Uint32) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Uint32)(nil)
	_ encoding.BinaryUnmarshaler = (*Uint32)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s Uint32) xdrType() {}

var _ xdrType = (*Uint32)(nil)

// Int32 is an XDR Typedef defines as:
//
//	typedef int int32;
type Int32 int32

// EncodeTo encodes this value using the Encoder.
func (s Int32) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(s)); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Int32)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Int32) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding Int32: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	var v int32
	v, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %w", err)
	}
	*s = Int32(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Int32) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Int32) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Int32)(nil)
	_ encoding.BinaryUnmarshaler = (*Int32)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s Int32) xdrType() {}

var _ xdrType = (*Int32)(nil)

// Uint64 is an XDR Typedef defines as:
//
//	typedef unsigned hyper uint64;
type Uint64 uint64

// EncodeTo encodes this value using the Encoder.
func (s Uint64) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeUhyper(uint64(s)); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Uint64)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Uint64) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding Uint64: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	var v uint64
	v, nTmp, err = d.DecodeUhyper()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Unsigned hyper: %w", err)
	}
	*s = Uint64(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Uint64) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Uint64) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Uint64)(nil)
	_ encoding.BinaryUnmarshaler = (*Uint64)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s Uint64) xdrType() {}

var _ xdrType = (*Uint64)(nil)

// Int64 is an XDR Typedef defines as:
//
//	typedef hyper int64;
type Int64 int64

// EncodeTo encodes this value using the Encoder.
func (s Int64) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeHyper(int64(s)); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Int64)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Int64) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding Int64: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	var v int64
	v, nTmp, err = d.DecodeHyper()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Hyper: %w", err)
	}
	*s = Int64(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Int64) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Int64) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Int64)(nil)
	_ encoding.BinaryUnmarshaler = (*Int64)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s Int64) xdrType() {}

var _ xdrType = (*Int64)(nil)

// TimePoint is an XDR Typedef defines as:
//
//	typedef uint64 TimePoint;
type TimePoint Uint64

// EncodeTo encodes this value using the Encoder.
func (s TimePoint) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = Uint64(s).EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*TimePoint)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TimePoint) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding TimePoint: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = (*Uint64)(s).DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TimePoint) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TimePoint) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TimePoint)(nil)
	_ encoding.BinaryUnmarshaler = (*TimePoint)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s TimePoint) xdrType() {}

var _ xdrType = (*TimePoint)(nil)

// Duration is an XDR Typedef defines as:
//
//	typedef uint64 Duration;
type Duration Uint64

// EncodeTo encodes this value using the Encoder.
func (s Duration) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = Uint64(s).EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Duration)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Duration) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding Duration: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = (*Uint64)(s).DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Duration) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Duration) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Duration)(nil)
	_ encoding.BinaryUnmarshaler = (*Duration)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s Duration) xdrType() {}

var _ xdrType = (*Duration)(nil)

// ExtensionPoint is an XDR Union defines as:
//
//	union ExtensionPoint switch (int v)
//	 {
//	 case 0:
//	     void;
//	 };
type ExtensionPoint struct {
	V int32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ExtensionPoint) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ExtensionPoint
func (u ExtensionPoint) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	}
	return "-", false
}

// NewExtensionPoint creates a new  ExtensionPoint.
func NewExtensionPoint(v int32, value interface{}) (result ExtensionPoint, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u ExtensionPoint) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union ExtensionPoint", u.V)
}

var _ decoderFrom = (*ExtensionPoint)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *ExtensionPoint) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding ExtensionPoint: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %w", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union ExtensionPoint has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ExtensionPoint) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ExtensionPoint) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ExtensionPoint)(nil)
	_ encoding.BinaryUnmarshaler = (*ExtensionPoint)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s ExtensionPoint) xdrType() {}

var _ xdrType = (*ExtensionPoint)(nil)

// PoolId is an XDR Typedef defines as:
//
//	typedef Hash PoolID;
type PoolId Hash

// EncodeTo encodes this value using the Encoder.
func (s *PoolId) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = (*Hash)(s).EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*PoolId)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *PoolId) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding PoolId: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = (*Hash)(s).DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Hash: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PoolId) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PoolId) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PoolId)(nil)
	_ encoding.BinaryUnmarshaler = (*PoolId)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s PoolId) xdrType() {}

var _ xdrType = (*PoolId)(nil)

// Liabilities is an XDR Struct defines as:
//
//	struct Liabilities
//	 {
//	     int64 buying;
//	     int64 selling;
//	 };
type Liabilities struct {
	Buying  Int64
	Selling Int64
}

// EncodeTo encodes this value using the Encoder.
func (s *Liabilities) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Buying.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Selling.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Liabilities)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Liabilities) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding Liabilities: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.Buying.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %w", err)
	}
	nTmp, err = s.Selling.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Liabilities) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Liabilities) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Liabilities)(nil)
	_ encoding.BinaryUnmarshaler = (*Liabilities)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s Liabilities) xdrType() {}

var _ xdrType = (*Liabilities)(nil)

// AssetCode4 is an XDR Typedef defines as:
//
//	typedef opaque AssetCode4[4];
type AssetCode4 [4]byte

// XDRMaxSize implements the Sized interface for AssetCode4
func (e AssetCode4) XDRMaxSize() int {
	return 4
}

// EncodeTo encodes this value using the Encoder.
func (s *AssetCode4) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeFixedOpaque(s[:]); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*AssetCode4)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *AssetCode4) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding AssetCode4: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = d.DecodeFixedOpaqueInplace(s[:])
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AssetCode4: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AssetCode4) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AssetCode4) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AssetCode4)(nil)
	_ encoding.BinaryUnmarshaler = (*AssetCode4)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s AssetCode4) xdrType() {}

var _ xdrType = (*AssetCode4)(nil)

// AssetCode12 is an XDR Typedef defines as:
//
//	typedef opaque AssetCode12[12];
type AssetCode12 [12]byte

// XDRMaxSize implements the Sized interface for AssetCode12
func (e AssetCode12) XDRMaxSize() int {
	return 12
}

// EncodeTo encodes this value using the Encoder.
func (s *AssetCode12) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeFixedOpaque(s[:]); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*AssetCode12)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *AssetCode12) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding AssetCode12: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = d.DecodeFixedOpaqueInplace(s[:])
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AssetCode12: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AssetCode12) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AssetCode12) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AssetCode12)(nil)
	_ encoding.BinaryUnmarshaler = (*AssetCode12)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s AssetCode12) xdrType() {}

var _ xdrType = (*AssetCode12)(nil)

// AccountId is an XDR Typedef defines as:
//
//	typedef PublicKey AccountID;
type AccountId PublicKey

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u AccountId) SwitchFieldName() string {
	return PublicKey(u).SwitchFieldName()
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PublicKey
func (u AccountId) ArmForSwitch(sw int32) (string, bool) {
	return PublicKey(u).ArmForSwitch(sw)
}

// NewAccountId creates a new  AccountId.
func NewAccountId(aType PublicKeyType, value interface{}) (result AccountId, err error) {
	u, err := NewPublicKey(aType, value)
	result = AccountId(u)
	return
}

// MustEd25519 retrieves the Ed25519 value from the union,
// panicing if the value is not set.
func (u AccountId) MustEd25519() Uint256 {
	return PublicKey(u).MustEd25519()
}

// GetEd25519 retrieves the Ed25519 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u AccountId) GetEd25519() (result Uint256, ok bool) {
	return PublicKey(u).GetEd25519()
}

// EncodeTo encodes this value using the Encoder.
func (s AccountId) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = PublicKey(s).EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*AccountId)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *AccountId) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding AccountId: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = (*PublicKey)(s).DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding PublicKey: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AccountId) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AccountId) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AccountId)(nil)
	_ encoding.BinaryUnmarshaler = (*AccountId)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s AccountId) xdrType() {}

var _ xdrType = (*AccountId)(nil)

type PublicKeyType int32

const (
	PublicKeyTypePublicKeyTypeEd25519 PublicKeyType = 0
)

var publicKeyTypeMap = map[int32]string{
	0: "PublicKeyTypePublicKeyTypeEd25519",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for PublicKeyType
func (e PublicKeyType) ValidEnum(v int32) bool {
	_, ok := publicKeyTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e PublicKeyType) String() string {
	name, _ := publicKeyTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e PublicKeyType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := publicKeyTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid PublicKeyType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*PublicKeyType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *PublicKeyType) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding PublicKeyType: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding PublicKeyType: %w", err)
	}
	if _, ok := publicKeyTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid PublicKeyType enum value", v)
	}
	*e = PublicKeyType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PublicKeyType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PublicKeyType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PublicKeyType)(nil)
	_ encoding.BinaryUnmarshaler = (*PublicKeyType)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s PublicKeyType) xdrType() {}

var _ xdrType = (*PublicKeyType)(nil)

// SignerKeyType is an XDR Enum defines as:
//
//	enum SignerKeyType
//	 {
//	     SIGNER_KEY_TYPE_ED25519 = KEY_TYPE_ED25519,
//	     SIGNER_KEY_TYPE_PRE_AUTH_TX = KEY_TYPE_PRE_AUTH_TX,
//	     SIGNER_KEY_TYPE_HASH_X = KEY_TYPE_HASH_X,
//	     SIGNER_KEY_TYPE_ED25519_SIGNED_PAYLOAD = KEY_TYPE_ED25519_SIGNED_PAYLOAD
//	 };
type SignerKeyType int32

const (
	SignerKeyTypeSignerKeyTypeEd25519              SignerKeyType = 0
	SignerKeyTypeSignerKeyTypePreAuthTx            SignerKeyType = 1
	SignerKeyTypeSignerKeyTypeHashX                SignerKeyType = 2
	SignerKeyTypeSignerKeyTypeEd25519SignedPayload SignerKeyType = 3
)

var signerKeyTypeMap = map[int32]string{
	0: "SignerKeyTypeSignerKeyTypeEd25519",
	1: "SignerKeyTypeSignerKeyTypePreAuthTx",
	2: "SignerKeyTypeSignerKeyTypeHashX",
	3: "SignerKeyTypeSignerKeyTypeEd25519SignedPayload",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for SignerKeyType
func (e SignerKeyType) ValidEnum(v int32) bool {
	_, ok := signerKeyTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e SignerKeyType) String() string {
	name, _ := signerKeyTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e SignerKeyType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := signerKeyTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid SignerKeyType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*SignerKeyType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *SignerKeyType) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding SignerKeyType: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding SignerKeyType: %w", err)
	}
	if _, ok := signerKeyTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid SignerKeyType enum value", v)
	}
	*e = SignerKeyType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s SignerKeyType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SignerKeyType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*SignerKeyType)(nil)
	_ encoding.BinaryUnmarshaler = (*SignerKeyType)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s SignerKeyType) xdrType() {}

var _ xdrType = (*SignerKeyType)(nil)

// PublicKey is an XDR Union defines as:
//
//	union PublicKey switch (PublicKeyType type)
//	 {
//	 case PUBLIC_KEY_TYPE_ED25519:
//	     uint256 ed25519;
//	 };
type PublicKey struct {
	Type    PublicKeyType
	Ed25519 *Uint256
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u PublicKey) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of PublicKey
func (u PublicKey) ArmForSwitch(sw int32) (string, bool) {
	switch PublicKeyType(sw) {
	case PublicKeyTypePublicKeyTypeEd25519:
		return "Ed25519", true
	}
	return "-", false
}

// NewPublicKey creates a new  PublicKey.
func NewPublicKey(aType PublicKeyType, value interface{}) (result PublicKey, err error) {
	result.Type = aType
	switch PublicKeyType(aType) {
	case PublicKeyTypePublicKeyTypeEd25519:
		tv, ok := value.(Uint256)
		if !ok {
			err = errors.New("invalid value, must be Uint256")
			return
		}
		result.Ed25519 = &tv
	}
	return
}

// MustEd25519 retrieves the Ed25519 value from the union,
// panicing if the value is not set.
func (u PublicKey) MustEd25519() Uint256 {
	val, ok := u.GetEd25519()

	if !ok {
		panic("arm Ed25519 is not set")
	}

	return val
}

// GetEd25519 retrieves the Ed25519 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u PublicKey) GetEd25519() (result Uint256, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Ed25519" {
		result = *u.Ed25519
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u PublicKey) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch PublicKeyType(u.Type) {
	case PublicKeyTypePublicKeyTypeEd25519:
		if err = (*u.Ed25519).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (PublicKeyType) switch value '%d' is not valid for union PublicKey", u.Type)
}

var _ decoderFrom = (*PublicKey)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *PublicKey) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding PublicKey: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding PublicKeyType: %w", err)
	}
	switch PublicKeyType(u.Type) {
	case PublicKeyTypePublicKeyTypeEd25519:
		u.Ed25519 = new(Uint256)
		nTmp, err = (*u.Ed25519).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint256: %w", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union PublicKey has invalid Type (PublicKeyType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PublicKey) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PublicKey) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PublicKey)(nil)
	_ encoding.BinaryUnmarshaler = (*PublicKey)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s PublicKey) xdrType() {}

var _ xdrType = (*PublicKey)(nil)

// SignerKeyEd25519SignedPayload is an XDR NestedStruct defines as:
//
//	struct
//	     {
//	         /* Public key that must sign the payload. */
//	         uint256 ed25519;
//	         /* Payload to be raw signed by ed25519. */
//	         opaque payload<64>;
//	     }
type SignerKeyEd25519SignedPayload struct {
	Ed25519 Uint256
	Payload []byte `xdrmaxsize:"64"`
}

// EncodeTo encodes this value using the Encoder.
func (s *SignerKeyEd25519SignedPayload) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Ed25519.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeOpaque(s.Payload[:]); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*SignerKeyEd25519SignedPayload)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *SignerKeyEd25519SignedPayload) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding SignerKeyEd25519SignedPayload: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.Ed25519.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint256: %w", err)
	}
	s.Payload, nTmp, err = d.DecodeOpaque(64)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Payload: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s SignerKeyEd25519SignedPayload) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SignerKeyEd25519SignedPayload) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*SignerKeyEd25519SignedPayload)(nil)
	_ encoding.BinaryUnmarshaler = (*SignerKeyEd25519SignedPayload)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s SignerKeyEd25519SignedPayload) xdrType() {}

var _ xdrType = (*SignerKeyEd25519SignedPayload)(nil)

// SignerKey is an XDR Union defines as:
//
//	union SignerKey switch (SignerKeyType type)
//	 {
//	 case SIGNER_KEY_TYPE_ED25519:
//	     uint256 ed25519;
//	 case SIGNER_KEY_TYPE_PRE_AUTH_TX:
//	     /* SHA-256 Hash of TransactionSignaturePayload structure */
//	     uint256 preAuthTx;
//	 case SIGNER_KEY_TYPE_HASH_X:
//	     /* Hash of random 256 bit preimage X */
//	     uint256 hashX;
//	 case SIGNER_KEY_TYPE_ED25519_SIGNED_PAYLOAD:
//	     struct
//	     {
//	         /* Public key that must sign the payload. */
//	         uint256 ed25519;
//	         /* Payload to be raw signed by ed25519. */
//	         opaque payload<64>;
//	     } ed25519SignedPayload;
//	 };
type SignerKey struct {
	Type                 SignerKeyType
	Ed25519              *Uint256
	PreAuthTx            *Uint256
	HashX                *Uint256
	Ed25519SignedPayload *SignerKeyEd25519SignedPayload
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u SignerKey) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of SignerKey
func (u SignerKey) ArmForSwitch(sw int32) (string, bool) {
	switch SignerKeyType(sw) {
	case SignerKeyTypeSignerKeyTypeEd25519:
		return "Ed25519", true
	case SignerKeyTypeSignerKeyTypePreAuthTx:
		return "PreAuthTx", true
	case SignerKeyTypeSignerKeyTypeHashX:
		return "HashX", true
	case SignerKeyTypeSignerKeyTypeEd25519SignedPayload:
		return "Ed25519SignedPayload", true
	}
	return "-", false
}

// NewSignerKey creates a new  SignerKey.
func NewSignerKey(aType SignerKeyType, value interface{}) (result SignerKey, err error) {
	result.Type = aType
	switch SignerKeyType(aType) {
	case SignerKeyTypeSignerKeyTypeEd25519:
		tv, ok := value.(Uint256)
		if !ok {
			err = errors.New("invalid value, must be Uint256")
			return
		}
		result.Ed25519 = &tv
	case SignerKeyTypeSignerKeyTypePreAuthTx:
		tv, ok := value.(Uint256)
		if !ok {
			err = errors.New("invalid value, must be Uint256")
			return
		}
		result.PreAuthTx = &tv
	case SignerKeyTypeSignerKeyTypeHashX:
		tv, ok := value.(Uint256)
		if !ok {
			err = errors.New("invalid value, must be Uint256")
			return
		}
		result.HashX = &tv
	case SignerKeyTypeSignerKeyTypeEd25519SignedPayload:
		tv, ok := value.(SignerKeyEd25519SignedPayload)
		if !ok {
			err = errors.New("invalid value, must be SignerKeyEd25519SignedPayload")
			return
		}
		result.Ed25519SignedPayload = &tv
	}
	return
}

// MustEd25519 retrieves the Ed25519 value from the union,
// panicing if the value is not set.
func (u SignerKey) MustEd25519() Uint256 {
	val, ok := u.GetEd25519()

	if !ok {
		panic("arm Ed25519 is not set")
	}

	return val
}

// GetEd25519 retrieves the Ed25519 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u SignerKey) GetEd25519() (result Uint256, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Ed25519" {
		result = *u.Ed25519
		ok = true
	}

	return
}

// MustPreAuthTx retrieves the PreAuthTx value from the union,
// panicing if the value is not set.
func (u SignerKey) MustPreAuthTx() Uint256 {
	val, ok := u.GetPreAuthTx()

	if !ok {
		panic("arm PreAuthTx is not set")
	}

	return val
}

// GetPreAuthTx retrieves the PreAuthTx value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u SignerKey) GetPreAuthTx() (result Uint256, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "PreAuthTx" {
		result = *u.PreAuthTx
		ok = true
	}

	return
}

// MustHashX retrieves the HashX value from the union,
// panicing if the value is not set.
func (u SignerKey) MustHashX() Uint256 {
	val, ok := u.GetHashX()

	if !ok {
		panic("arm HashX is not set")
	}

	return val
}

// GetHashX retrieves the HashX value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u SignerKey) GetHashX() (result Uint256, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "HashX" {
		result = *u.HashX
		ok = true
	}

	return
}

// MustEd25519SignedPayload retrieves the Ed25519SignedPayload value from the union,
// panicing if the value is not set.
func (u SignerKey) MustEd25519SignedPayload() SignerKeyEd25519SignedPayload {
	val, ok := u.GetEd25519SignedPayload()

	if !ok {
		panic("arm Ed25519SignedPayload is not set")
	}

	return val
}

// GetEd25519SignedPayload retrieves the Ed25519SignedPayload value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u SignerKey) GetEd25519SignedPayload() (result SignerKeyEd25519SignedPayload, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Ed25519SignedPayload" {
		result = *u.Ed25519SignedPayload
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u SignerKey) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch SignerKeyType(u.Type) {
	case SignerKeyTypeSignerKeyTypeEd25519:
		if err = (*u.Ed25519).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case SignerKeyTypeSignerKeyTypePreAuthTx:
		if err = (*u.PreAuthTx).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case SignerKeyTypeSignerKeyTypeHashX:
		if err = (*u.HashX).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case SignerKeyTypeSignerKeyTypeEd25519SignedPayload:
		if err = (*u.Ed25519SignedPayload).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (SignerKeyType) switch value '%d' is not valid for union SignerKey", u.Type)
}

var _ decoderFrom = (*SignerKey)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *SignerKey) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding SignerKey: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SignerKeyType: %w", err)
	}
	switch SignerKeyType(u.Type) {
	case SignerKeyTypeSignerKeyTypeEd25519:
		u.Ed25519 = new(Uint256)
		nTmp, err = (*u.Ed25519).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint256: %w", err)
		}
		return n, nil
	case SignerKeyTypeSignerKeyTypePreAuthTx:
		u.PreAuthTx = new(Uint256)
		nTmp, err = (*u.PreAuthTx).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint256: %w", err)
		}
		return n, nil
	case SignerKeyTypeSignerKeyTypeHashX:
		u.HashX = new(Uint256)
		nTmp, err = (*u.HashX).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint256: %w", err)
		}
		return n, nil
	case SignerKeyTypeSignerKeyTypeEd25519SignedPayload:
		u.Ed25519SignedPayload = new(SignerKeyEd25519SignedPayload)
		nTmp, err = (*u.Ed25519SignedPayload).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding SignerKeyEd25519SignedPayload: %w", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union SignerKey has invalid Type (SignerKeyType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s SignerKey) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SignerKey) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*SignerKey)(nil)
	_ encoding.BinaryUnmarshaler = (*SignerKey)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s SignerKey) xdrType() {}

var _ xdrType = (*SignerKey)(nil)

// Signature is an XDR Typedef defines as:
//
//	typedef opaque Signature<64>;
type Signature []byte

// XDRMaxSize implements the Sized interface for Signature
func (e Signature) XDRMaxSize() int {
	return 64
}

// EncodeTo encodes this value using the Encoder.
func (s Signature) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeOpaque(s[:]); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Signature)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Signature) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding Signature: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	(*s), nTmp, err = d.DecodeOpaque(64)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Signature: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Signature) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Signature) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Signature)(nil)
	_ encoding.BinaryUnmarshaler = (*Signature)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s Signature) xdrType() {}

var _ xdrType = (*Signature)(nil)

// SignatureHint is an XDR Typedef defines as:
//
//	typedef opaque SignatureHint[4];
type SignatureHint [4]byte

// XDRMaxSize implements the Sized interface for SignatureHint
func (e SignatureHint) XDRMaxSize() int {
	return 4
}

// EncodeTo encodes this value using the Encoder.
func (s *SignatureHint) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeFixedOpaque(s[:]); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*SignatureHint)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *SignatureHint) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding SignatureHint: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = d.DecodeFixedOpaqueInplace(s[:])
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SignatureHint: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s SignatureHint) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SignatureHint) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*SignatureHint)(nil)
	_ encoding.BinaryUnmarshaler = (*SignatureHint)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s SignatureHint) xdrType() {}

var _ xdrType = (*SignatureHint)(nil)

type ChangeTrustAsset struct {
	Type       AssetType
	AlphaNum4  *AlphaNum4
	AlphaNum12 *AlphaNum12
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u ChangeTrustAsset) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of ChangeTrustAsset
func (u ChangeTrustAsset) ArmForSwitch(sw int32) (string, bool) {
	switch AssetType(sw) {
	case AssetTypeAssetTypeNative:
		return "", true
	case AssetTypeAssetTypeCreditAlphanum4:
		return "AlphaNum4", true
	case AssetTypeAssetTypeCreditAlphanum12:
		return "AlphaNum12", true
	case AssetTypeAssetTypePoolShare:
		return "LiquidityPool", true
	}
	return "-", false
}

// NewChangeTrustAsset creates a new  ChangeTrustAsset.
func NewChangeTrustAsset(aType AssetType, value interface{}) (result ChangeTrustAsset, err error) {
	result.Type = aType
	switch AssetType(aType) {
	case AssetTypeAssetTypeNative:
		// void
	case AssetTypeAssetTypeCreditAlphanum4:
		tv, ok := value.(AlphaNum4)
		if !ok {
			err = errors.New("invalid value, must be AlphaNum4")
			return
		}
		result.AlphaNum4 = &tv
	case AssetTypeAssetTypeCreditAlphanum12:
		tv, ok := value.(AlphaNum12)
		if !ok {
			err = errors.New("invalid value, must be AlphaNum12")
			return
		}
		result.AlphaNum12 = &tv
	}
	return
}

// MustAlphaNum4 retrieves the AlphaNum4 value from the union,
// panicing if the value is not set.
func (u ChangeTrustAsset) MustAlphaNum4() AlphaNum4 {
	val, ok := u.GetAlphaNum4()

	if !ok {
		panic("arm AlphaNum4 is not set")
	}

	return val
}

// GetAlphaNum4 retrieves the AlphaNum4 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ChangeTrustAsset) GetAlphaNum4() (result AlphaNum4, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AlphaNum4" {
		result = *u.AlphaNum4
		ok = true
	}

	return
}

// MustAlphaNum12 retrieves the AlphaNum12 value from the union,
// panicing if the value is not set.
func (u ChangeTrustAsset) MustAlphaNum12() AlphaNum12 {
	val, ok := u.GetAlphaNum12()

	if !ok {
		panic("arm AlphaNum12 is not set")
	}

	return val
}

// GetAlphaNum12 retrieves the AlphaNum12 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u ChangeTrustAsset) GetAlphaNum12() (result AlphaNum12, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AlphaNum12" {
		result = *u.AlphaNum12
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u ChangeTrustAsset) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch AssetType(u.Type) {
	case AssetTypeAssetTypeNative:
		// Void
		return nil
	case AssetTypeAssetTypeCreditAlphanum4:
		if err = (*u.AlphaNum4).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case AssetTypeAssetTypeCreditAlphanum12:
		if err = (*u.AlphaNum12).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (AssetType) switch value '%d' is not valid for union ChangeTrustAsset", u.Type)
}

var _ decoderFrom = (*ChangeTrustAsset)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *ChangeTrustAsset) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding ChangeTrustAsset: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AssetType: %w", err)
	}
	switch AssetType(u.Type) {
	case AssetTypeAssetTypeNative:
		// Void
		return n, nil
	case AssetTypeAssetTypeCreditAlphanum4:
		u.AlphaNum4 = new(AlphaNum4)
		nTmp, err = (*u.AlphaNum4).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AlphaNum4: %w", err)
		}
		return n, nil
	case AssetTypeAssetTypeCreditAlphanum12:
		u.AlphaNum12 = new(AlphaNum12)
		nTmp, err = (*u.AlphaNum12).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AlphaNum12: %w", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union ChangeTrustAsset has invalid Type (AssetType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ChangeTrustAsset) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ChangeTrustAsset) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ChangeTrustAsset)(nil)
	_ encoding.BinaryUnmarshaler = (*ChangeTrustAsset)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s ChangeTrustAsset) xdrType() {}

var _ xdrType = (*ChangeTrustAsset)(nil)

// ChangeTrustOp is an XDR Struct defines as:
//
//	struct ChangeTrustOp
//	 {
//	     ChangeTrustAsset line;
//
//	     // if limit is set to 0, deletes the trust line
//	     int64 limit;
//	 };
type ChangeTrustOp struct {
	Line  ChangeTrustAsset
	Limit Int64
}

// EncodeTo encodes this value using the Encoder.
func (s *ChangeTrustOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Line.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Limit.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*ChangeTrustOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *ChangeTrustOp) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding ChangeTrustOp: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.Line.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding ChangeTrustAsset: %w", err)
	}
	nTmp, err = s.Limit.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s ChangeTrustOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *ChangeTrustOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*ChangeTrustOp)(nil)
	_ encoding.BinaryUnmarshaler = (*ChangeTrustOp)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s ChangeTrustOp) xdrType() {}

var _ xdrType = (*ChangeTrustOp)(nil)

// AllowTrustOp is an XDR Struct defines as:
//
//	struct AllowTrustOp
//	 {
//	     AccountID trustor;
//	     AssetCode asset;
//
//	     // One of 0, AUTHORIZED_FLAG, or AUTHORIZED_TO_MAINTAIN_LIABILITIES_FLAG
//	     uint32 authorize;
//	 };
type AllowTrustOp struct {
	Trustor   AccountId
	Asset     AssetCode
	Authorize Uint32
}

// EncodeTo encodes this value using the Encoder.
func (s *AllowTrustOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Trustor.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Asset.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Authorize.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*AllowTrustOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *AllowTrustOp) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding AllowTrustOp: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.Trustor.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %w", err)
	}
	nTmp, err = s.Asset.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AssetCode: %w", err)
	}
	nTmp, err = s.Authorize.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s AllowTrustOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *AllowTrustOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*AllowTrustOp)(nil)
	_ encoding.BinaryUnmarshaler = (*AllowTrustOp)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s AllowTrustOp) xdrType() {}

var _ xdrType = (*AllowTrustOp)(nil)

type LedgerKey struct {
	Type      LedgerEntryType
	Account   *LedgerKeyAccount
	TrustLine *LedgerKeyTrustLine
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerKey) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerKey
func (u LedgerKey) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerEntryType(sw) {
	case LedgerEntryTypeAccount:
		return "Account", true
	case LedgerEntryTypeTrustline:
		return "TrustLine", true
	case LedgerEntryTypeOffer:
		return "Offer", true
	case LedgerEntryTypeData:
		return "Data", true
	}
	return "-", false
}

// NewLedgerKey creates a new  LedgerKey.
func NewLedgerKey(aType LedgerEntryType, value interface{}) (result LedgerKey, err error) {
	result.Type = aType
	switch LedgerEntryType(aType) {
	case LedgerEntryTypeAccount:
		tv, ok := value.(LedgerKeyAccount)
		if !ok {
			err = errors.New("invalid value, must be LedgerKeyAccount")
			return
		}
		result.Account = &tv
	case LedgerEntryTypeTrustline:
		tv, ok := value.(LedgerKeyTrustLine)
		if !ok {
			err = errors.New("invalid value, must be LedgerKeyTrustLine")
			return
		}
		result.TrustLine = &tv
	}
	return
}

// MustAccount retrieves the Account value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustAccount() LedgerKeyAccount {
	val, ok := u.GetAccount()

	if !ok {
		panic("arm Account is not set")
	}

	return val
}

// GetAccount retrieves the Account value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetAccount() (result LedgerKeyAccount, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Account" {
		result = *u.Account
		ok = true
	}

	return
}

// MustTrustLine retrieves the TrustLine value from the union,
// panicing if the value is not set.
func (u LedgerKey) MustTrustLine() LedgerKeyTrustLine {
	val, ok := u.GetTrustLine()

	if !ok {
		panic("arm TrustLine is not set")
	}

	return val
}

// GetTrustLine retrieves the TrustLine value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerKey) GetTrustLine() (result LedgerKeyTrustLine, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "TrustLine" {
		result = *u.TrustLine
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u LedgerKey) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch LedgerEntryType(u.Type) {
	case LedgerEntryTypeAccount:
		if err = (*u.Account).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case LedgerEntryTypeTrustline:
		if err = (*u.TrustLine).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (LedgerEntryType) switch value '%d' is not valid for union LedgerKey", u.Type)
}

var _ decoderFrom = (*LedgerKey)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *LedgerKey) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding LedgerKey: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerEntryType: %w", err)
	}
	switch LedgerEntryType(u.Type) {
	case LedgerEntryTypeAccount:
		u.Account = new(LedgerKeyAccount)
		nTmp, err = (*u.Account).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LedgerKeyAccount: %w", err)
		}
		return n, nil
	case LedgerEntryTypeTrustline:
		u.TrustLine = new(LedgerKeyTrustLine)
		nTmp, err = (*u.TrustLine).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LedgerKeyTrustLine: %w", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union LedgerKey has invalid Type (LedgerEntryType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerKey) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerKey) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerKey)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerKey)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s LedgerKey) xdrType() {}

var _ xdrType = (*LedgerKey)(nil)

type LedgerEntry struct {
	LastModifiedLedgerSeq Uint32
	Data                  LedgerEntryData
	Ext                   LedgerEntryExt
}

// EncodeTo encodes this value using the Encoder.
func (s *LedgerEntry) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.LastModifiedLedgerSeq.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Data.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*LedgerEntry)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LedgerEntry) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding LedgerEntry: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.LastModifiedLedgerSeq.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %w", err)
	}
	nTmp, err = s.Data.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerEntryData: %w", err)
	}
	nTmp, err = s.Ext.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerEntryExt: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerEntry) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerEntry) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerEntry)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerEntry)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s LedgerEntry) xdrType() {}

var _ xdrType = (*LedgerEntry)(nil)

type LedgerEntryData struct {
	Type      LedgerEntryType
	TrustLine *TrustLineEntry
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerEntryData) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerEntryData
func (u LedgerEntryData) ArmForSwitch(sw int32) (string, bool) {
	switch LedgerEntryType(sw) {
	case LedgerEntryTypeAccount:
		return "Account", true
	case LedgerEntryTypeTrustline:
		return "TrustLine", true
	case LedgerEntryTypeOffer:
		return "Offer", true
	case LedgerEntryTypeData:
		return "Data", true
	}
	return "-", false
}

// NewLedgerEntryData creates a new  LedgerEntryData.
func NewLedgerEntryData(aType LedgerEntryType, value interface{}) (result LedgerEntryData, err error) {
	result.Type = aType
	switch LedgerEntryType(aType) {
	case LedgerEntryTypeTrustline:
		tv, ok := value.(TrustLineEntry)
		if !ok {
			err = errors.New("invalid value, must be TrustLineEntry")
			return
		}
		result.TrustLine = &tv
	}
	return
}

// MustTrustLine retrieves the TrustLine value from the union,
// panicing if the value is not set.
func (u LedgerEntryData) MustTrustLine() TrustLineEntry {
	val, ok := u.GetTrustLine()

	if !ok {
		panic("arm TrustLine is not set")
	}

	return val
}

// GetTrustLine retrieves the TrustLine value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryData) GetTrustLine() (result TrustLineEntry, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "TrustLine" {
		result = *u.TrustLine
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u LedgerEntryData) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch LedgerEntryType(u.Type) {
	case LedgerEntryTypeTrustline:
		if err = (*u.TrustLine).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (LedgerEntryType) switch value '%d' is not valid for union LedgerEntryData", u.Type)
}

var _ decoderFrom = (*LedgerEntryData)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *LedgerEntryData) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding LedgerEntryData: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerEntryType: %w", err)
	}
	switch LedgerEntryType(u.Type) {
	case LedgerEntryTypeTrustline:
		u.TrustLine = new(TrustLineEntry)
		nTmp, err = (*u.TrustLine).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding TrustLineEntry: %w", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union LedgerEntryData has invalid Type (LedgerEntryType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerEntryData) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerEntryData) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerEntryData)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerEntryData)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s LedgerEntryData) xdrType() {}

var _ xdrType = (*LedgerEntryData)(nil)

type LedgerKeyTrustLine struct {
	AccountId AccountId
	Asset     TrustLineAsset
}

// EncodeTo encodes this value using the Encoder.
func (s *LedgerKeyTrustLine) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.AccountId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Asset.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*LedgerKeyTrustLine)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LedgerKeyTrustLine) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding LedgerKeyTrustLine: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.AccountId.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %w", err)
	}
	nTmp, err = s.Asset.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TrustLineAsset: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerKeyTrustLine) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerKeyTrustLine) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerKeyTrustLine)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerKeyTrustLine)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s LedgerKeyTrustLine) xdrType() {}

var _ xdrType = (*LedgerKeyTrustLine)(nil)

type LedgerEntryType int32

const (
	LedgerEntryTypeAccount   LedgerEntryType = 0
	LedgerEntryTypeTrustline LedgerEntryType = 1
	LedgerEntryTypeOffer     LedgerEntryType = 2
	LedgerEntryTypeData      LedgerEntryType = 3
)

var ledgerEntryTypeMap = map[int32]string{
	0: "LedgerEntryTypeAccount",
	1: "LedgerEntryTypeTrustline",
	2: "LedgerEntryTypeOffer",
	3: "LedgerEntryTypeData",
	4: "LedgerEntryTypeClaimableBalance",
	5: "LedgerEntryTypeLiquidityPool",
	6: "LedgerEntryTypeContractData",
	7: "LedgerEntryTypeContractCode",
	8: "LedgerEntryTypeConfigSetting",
	9: "LedgerEntryTypeTtl",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for LedgerEntryType
func (e LedgerEntryType) ValidEnum(v int32) bool {
	_, ok := ledgerEntryTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e LedgerEntryType) String() string {
	name, _ := ledgerEntryTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e LedgerEntryType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := ledgerEntryTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid LedgerEntryType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*LedgerEntryType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *LedgerEntryType) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding LedgerEntryType: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding LedgerEntryType: %w", err)
	}
	if _, ok := ledgerEntryTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid LedgerEntryType enum value", v)
	}
	*e = LedgerEntryType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerEntryType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerEntryType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerEntryType)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerEntryType)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s LedgerEntryType) xdrType() {}

var _ xdrType = (*LedgerEntryType)(nil)

type LedgerKeyAccount struct {
	AccountId AccountId
}

// EncodeTo encodes this value using the Encoder.
func (s *LedgerKeyAccount) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.AccountId.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*LedgerKeyAccount)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LedgerKeyAccount) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding LedgerKeyAccount: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.AccountId.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerKeyAccount) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerKeyAccount) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerKeyAccount)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerKeyAccount)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s LedgerKeyAccount) xdrType() {}

var _ xdrType = (*LedgerKeyAccount)(nil)

type OperationType int32

const (
	OperationTypePayment                  OperationType = 1
	OperationTypePathPaymentStrictReceive OperationType = 2
	OperationTypeManageSellOffer          OperationType = 3
	OperationTypeCreatePassiveSellOffer   OperationType = 4
	OperationTypeSetOptions               OperationType = 5
	OperationTypeChangeTrust              OperationType = 6
	OperationTypeAllowTrust               OperationType = 7
	OperationTypeAccountMerge             OperationType = 8
	OperationTypeInflation                OperationType = 9
)

var operationTypeMap = map[int32]string{
	0:  "OperationTypeCreateAccount",
	1:  "OperationTypePayment",
	2:  "OperationTypePathPaymentStrictReceive",
	3:  "OperationTypeManageSellOffer",
	4:  "OperationTypeCreatePassiveSellOffer",
	5:  "OperationTypeSetOptions",
	6:  "OperationTypeChangeTrust",
	7:  "OperationTypeAllowTrust",
	8:  "OperationTypeAccountMerge",
	9:  "OperationTypeInflation",
	10: "OperationTypeManageData",
	11: "OperationTypeBumpSequence",
	12: "OperationTypeManageBuyOffer",
	13: "OperationTypePathPaymentStrictSend",
	14: "OperationTypeCreateClaimableBalance",
	15: "OperationTypeClaimClaimableBalance",
	16: "OperationTypeBeginSponsoringFutureReserves",
	17: "OperationTypeEndSponsoringFutureReserves",
	18: "OperationTypeRevokeSponsorship",
	19: "OperationTypeClawback",
	20: "OperationTypeClawbackClaimableBalance",
	21: "OperationTypeSetTrustLineFlags",
	22: "OperationTypeLiquidityPoolDeposit",
	23: "OperationTypeLiquidityPoolWithdraw",
	24: "OperationTypeInvokeHostFunction",
	25: "OperationTypeExtendFootprintTtl",
	26: "OperationTypeRestoreFootprint",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for OperationType
func (e OperationType) ValidEnum(v int32) bool {
	_, ok := operationTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e OperationType) String() string {
	name, _ := operationTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e OperationType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := operationTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid OperationType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*OperationType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *OperationType) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding OperationType: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding OperationType: %w", err)
	}
	if _, ok := operationTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid OperationType enum value", v)
	}
	*e = OperationType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s OperationType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *OperationType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*OperationType)(nil)
	_ encoding.BinaryUnmarshaler = (*OperationType)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s OperationType) xdrType() {}

var _ xdrType = (*OperationType)(nil)

type TrustLineEntry struct {
	AccountId AccountId
	Asset     TrustLineAsset
	Balance   Int64
	Limit     Int64
	Flags     Uint32
	Ext       TrustLineEntryExt
}

// EncodeTo encodes this value using the Encoder.
func (s *TrustLineEntry) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.AccountId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Asset.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Balance.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Limit.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Flags.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*TrustLineEntry)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TrustLineEntry) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding TrustLineEntry: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.AccountId.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding AccountId: %w", err)
	}
	nTmp, err = s.Asset.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TrustLineAsset: %w", err)
	}
	nTmp, err = s.Balance.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %w", err)
	}
	nTmp, err = s.Limit.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %w", err)
	}
	nTmp, err = s.Flags.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %w", err)
	}
	nTmp, err = s.Ext.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TrustLineEntryExt: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TrustLineEntry) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TrustLineEntry) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TrustLineEntry)(nil)
	_ encoding.BinaryUnmarshaler = (*TrustLineEntry)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s TrustLineEntry) xdrType() {}

var _ xdrType = (*TrustLineEntry)(nil)

type TrustLineEntryExt struct {
	V  int32
	V1 *TrustLineEntryV1
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TrustLineEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TrustLineEntryExt
func (u TrustLineEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	case 1:
		return "V1", true
	}
	return "-", false
}

// NewTrustLineEntryExt creates a new  TrustLineEntryExt.
func NewTrustLineEntryExt(v int32, value interface{}) (result TrustLineEntryExt, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	case 1:
		tv, ok := value.(TrustLineEntryV1)
		if !ok {
			err = errors.New("invalid value, must be TrustLineEntryV1")
			return
		}
		result.V1 = &tv
	}
	return
}

// MustV1 retrieves the V1 value from the union,
// panicing if the value is not set.
func (u TrustLineEntryExt) MustV1() TrustLineEntryV1 {
	val, ok := u.GetV1()

	if !ok {
		panic("arm V1 is not set")
	}

	return val
}

// GetV1 retrieves the V1 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TrustLineEntryExt) GetV1() (result TrustLineEntryV1, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "V1" {
		result = *u.V1
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u TrustLineEntryExt) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	case 1:
		if err = (*u.V1).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union TrustLineEntryExt", u.V)
}

var _ decoderFrom = (*TrustLineEntryExt)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *TrustLineEntryExt) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding TrustLineEntryExt: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %w", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	case 1:
		u.V1 = new(TrustLineEntryV1)
		nTmp, err = (*u.V1).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding TrustLineEntryV1: %w", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union TrustLineEntryExt has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TrustLineEntryExt) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TrustLineEntryExt) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TrustLineEntryExt)(nil)
	_ encoding.BinaryUnmarshaler = (*TrustLineEntryExt)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s TrustLineEntryExt) xdrType() {}

var _ xdrType = (*TrustLineEntryExt)(nil)

type LedgerEntryExt struct {
	V  int32
	V1 *LedgerEntryExtensionV1
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerEntryExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerEntryExt
func (u LedgerEntryExt) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	case 1:
		return "V1", true
	}
	return "-", false
}

// NewLedgerEntryExt creates a new  LedgerEntryExt.
func NewLedgerEntryExt(v int32, value interface{}) (result LedgerEntryExt, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	case 1:
		tv, ok := value.(LedgerEntryExtensionV1)
		if !ok {
			err = errors.New("invalid value, must be LedgerEntryExtensionV1")
			return
		}
		result.V1 = &tv
	}
	return
}

// MustV1 retrieves the V1 value from the union,
// panicing if the value is not set.
func (u LedgerEntryExt) MustV1() LedgerEntryExtensionV1 {
	val, ok := u.GetV1()

	if !ok {
		panic("arm V1 is not set")
	}

	return val
}

// GetV1 retrieves the V1 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u LedgerEntryExt) GetV1() (result LedgerEntryExtensionV1, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.V))

	if armName == "V1" {
		result = *u.V1
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u LedgerEntryExt) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	case 1:
		if err = (*u.V1).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union LedgerEntryExt", u.V)
}

var _ decoderFrom = (*LedgerEntryExt)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *LedgerEntryExt) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding LedgerEntryExt: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %w", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	case 1:
		u.V1 = new(LedgerEntryExtensionV1)
		nTmp, err = (*u.V1).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LedgerEntryExtensionV1: %w", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union LedgerEntryExt has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerEntryExt) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerEntryExt) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerEntryExt)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerEntryExt)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s LedgerEntryExt) xdrType() {}

var _ xdrType = (*LedgerEntryExt)(nil)

type SponsorshipDescriptor = *AccountId

type LedgerEntryExtensionV1 struct {
	SponsoringId SponsorshipDescriptor
	Ext          LedgerEntryExtensionV1Ext
}

// EncodeTo encodes this value using the Encoder.
func (s *LedgerEntryExtensionV1) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeBool(s.SponsoringId != nil); err != nil {
		return err
	}
	if s.SponsoringId != nil {
		if err = (*s.SponsoringId).EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*LedgerEntryExtensionV1)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LedgerEntryExtensionV1) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding LedgerEntryExtensionV1: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	var b bool
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SponsorshipDescriptor: %w", err)
	}
	s.SponsoringId = nil
	if b {
		s.SponsoringId = new(AccountId)
		nTmp, err = s.SponsoringId.DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding SponsorshipDescriptor: %w", err)
		}
	}
	nTmp, err = s.Ext.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerEntryExtensionV1Ext: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerEntryExtensionV1) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerEntryExtensionV1) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerEntryExtensionV1)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerEntryExtensionV1)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s LedgerEntryExtensionV1) xdrType() {}

var _ xdrType = (*LedgerEntryExtensionV1)(nil)

type LedgerEntryExtensionV1Ext struct {
	V int32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u LedgerEntryExtensionV1Ext) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of LedgerEntryExtensionV1Ext
func (u LedgerEntryExtensionV1Ext) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	}
	return "-", false
}

// NewLedgerEntryExtensionV1Ext creates a new  LedgerEntryExtensionV1Ext.
func NewLedgerEntryExtensionV1Ext(v int32, value interface{}) (result LedgerEntryExtensionV1Ext, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u LedgerEntryExtensionV1Ext) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union LedgerEntryExtensionV1Ext", u.V)
}

var _ decoderFrom = (*LedgerEntryExtensionV1Ext)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *LedgerEntryExtensionV1Ext) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding LedgerEntryExtensionV1Ext: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %w", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union LedgerEntryExtensionV1Ext has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerEntryExtensionV1Ext) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerEntryExtensionV1Ext) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerEntryExtensionV1Ext)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerEntryExtensionV1Ext)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s LedgerEntryExtensionV1Ext) xdrType() {}

var _ xdrType = (*LedgerEntryExtensionV1Ext)(nil)

type String32 string

// XDRMaxSize implements the Sized interface for String32
func (e String32) XDRMaxSize() int {
	return 32
}

// EncodeTo encodes this value using the Encoder.
func (s String32) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeString(string(s)); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*String32)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *String32) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding String32: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	var v string
	v, nTmp, err = d.DecodeString(32)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding String32: %w", err)
	}
	*s = String32(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s String32) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *String32) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*String32)(nil)
	_ encoding.BinaryUnmarshaler = (*String32)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s String32) xdrType() {}

var _ xdrType = (*String32)(nil)

// String64 is an XDR Typedef defines as:
//
//	typedef string string64<64>;
type String64 string

// XDRMaxSize implements the Sized interface for String64
func (e String64) XDRMaxSize() int {
	return 64
}

// EncodeTo encodes this value using the Encoder.
func (s String64) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeString(string(s)); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*String64)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *String64) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding String64: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	var v string
	v, nTmp, err = d.DecodeString(64)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding String64: %w", err)
	}
	*s = String64(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s String64) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *String64) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*String64)(nil)
	_ encoding.BinaryUnmarshaler = (*String64)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s String64) xdrType() {}

var _ xdrType = (*String64)(nil)

type Transaction struct {
	SourceAccount MuxedAccount
	Fee           Uint32
	SeqNum        SequenceNumber
	Cond          Preconditions
	Memo          Memo
	Operations    []Operation `xdrmaxsize:"100"`
	Ext           TransactionExt
}

// EncodeTo encodes this value using the Encoder.
func (s *Transaction) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.SourceAccount.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Fee.EncodeTo(e); err != nil {
		return err
	}
	if err = s.SeqNum.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Cond.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Memo.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.Operations))); err != nil {
		return err
	}
	for i := 0; i < len(s.Operations); i++ {
		if err = s.Operations[i].EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Transaction)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Transaction) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding Transaction: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.SourceAccount.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding MuxedAccount: %w", err)
	}
	nTmp, err = s.Fee.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %w", err)
	}
	nTmp, err = s.SeqNum.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SequenceNumber: %w", err)
	}
	nTmp, err = s.Cond.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Preconditions: %w", err)
	}
	nTmp, err = s.Memo.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Memo: %w", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Operation: %w", err)
	}
	if l > 100 {
		return n, fmt.Errorf("decoding Operation: data size (%d) exceeds size limit (100)", l)
	}
	s.Operations = nil
	if l > 0 {
		if il, ok := d.InputLen(); ok && uint(il) < uint(l) {
			return n, fmt.Errorf("decoding Operation: length (%d) exceeds remaining input length (%d)", l, il)
		}
		s.Operations = make([]Operation, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Operations[i].DecodeFrom(d, maxDepth)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding Operation: %w", err)
			}
		}
	}
	nTmp, err = s.Ext.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionExt: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Transaction) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Transaction) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Transaction)(nil)
	_ encoding.BinaryUnmarshaler = (*Transaction)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s Transaction) xdrType() {}

var _ xdrType = (*Transaction)(nil)

type PreconditionType int32

const (
	PreconditionTypePrecondNone PreconditionType = 0
	PreconditionTypePrecondTime PreconditionType = 1
	PreconditionTypePrecondV2   PreconditionType = 2
)

var preconditionTypeMap = map[int32]string{
	0: "PreconditionTypePrecondNone",
	1: "PreconditionTypePrecondTime",
	2: "PreconditionTypePrecondV2",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for PreconditionType
func (e PreconditionType) ValidEnum(v int32) bool {
	_, ok := preconditionTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e PreconditionType) String() string {
	name, _ := preconditionTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e PreconditionType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := preconditionTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid PreconditionType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*PreconditionType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *PreconditionType) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding PreconditionType: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding PreconditionType: %w", err)
	}
	if _, ok := preconditionTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid PreconditionType enum value", v)
	}
	*e = PreconditionType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PreconditionType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PreconditionType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PreconditionType)(nil)
	_ encoding.BinaryUnmarshaler = (*PreconditionType)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s PreconditionType) xdrType() {}

var _ xdrType = (*PreconditionType)(nil)

type MuxedAccount struct {
	Type     CryptoKeyType
	Ed25519  *Uint256
	Med25519 *MuxedAccountMed25519
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u MuxedAccount) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of MuxedAccount
func (u MuxedAccount) ArmForSwitch(sw int32) (string, bool) {
	switch CryptoKeyType(sw) {
	case CryptoKeyTypeKeyTypeEd25519:
		return "Ed25519", true
	case CryptoKeyTypeKeyTypeMuxedEd25519:
		return "Med25519", true
	}
	return "-", false
}

// NewMuxedAccount creates a new  MuxedAccount.
func NewMuxedAccount(aType CryptoKeyType, value interface{}) (result MuxedAccount, err error) {
	result.Type = aType
	switch CryptoKeyType(aType) {
	case CryptoKeyTypeKeyTypeEd25519:
		tv, ok := value.(Uint256)
		if !ok {
			err = errors.New("invalid value, must be Uint256")
			return
		}
		result.Ed25519 = &tv
	case CryptoKeyTypeKeyTypeMuxedEd25519:
		tv, ok := value.(MuxedAccountMed25519)
		if !ok {
			err = errors.New("invalid value, must be MuxedAccountMed25519")
			return
		}
		result.Med25519 = &tv
	}
	return
}

// MustEd25519 retrieves the Ed25519 value from the union,
// panicing if the value is not set.
func (u MuxedAccount) MustEd25519() Uint256 {
	val, ok := u.GetEd25519()

	if !ok {
		panic("arm Ed25519 is not set")
	}

	return val
}

// GetEd25519 retrieves the Ed25519 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u MuxedAccount) GetEd25519() (result Uint256, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Ed25519" {
		result = *u.Ed25519
		ok = true
	}

	return
}

// MustMed25519 retrieves the Med25519 value from the union,
// panicing if the value is not set.
func (u MuxedAccount) MustMed25519() MuxedAccountMed25519 {
	val, ok := u.GetMed25519()

	if !ok {
		panic("arm Med25519 is not set")
	}

	return val
}

// GetMed25519 retrieves the Med25519 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u MuxedAccount) GetMed25519() (result MuxedAccountMed25519, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Med25519" {
		result = *u.Med25519
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u MuxedAccount) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch CryptoKeyType(u.Type) {
	case CryptoKeyTypeKeyTypeEd25519:
		if err = (*u.Ed25519).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case CryptoKeyTypeKeyTypeMuxedEd25519:
		if err = (*u.Med25519).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (CryptoKeyType) switch value '%d' is not valid for union MuxedAccount", u.Type)
}

var _ decoderFrom = (*MuxedAccount)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *MuxedAccount) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding MuxedAccount: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding CryptoKeyType: %w", err)
	}
	switch CryptoKeyType(u.Type) {
	case CryptoKeyTypeKeyTypeEd25519:
		u.Ed25519 = new(Uint256)
		nTmp, err = (*u.Ed25519).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint256: %w", err)
		}
		return n, nil
	case CryptoKeyTypeKeyTypeMuxedEd25519:
		u.Med25519 = new(MuxedAccountMed25519)
		nTmp, err = (*u.Med25519).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding MuxedAccountMed25519: %w", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union MuxedAccount has invalid Type (CryptoKeyType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s MuxedAccount) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *MuxedAccount) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*MuxedAccount)(nil)
	_ encoding.BinaryUnmarshaler = (*MuxedAccount)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s MuxedAccount) xdrType() {}

var _ xdrType = (*MuxedAccount)(nil)

type SequenceNumber Int64

// EncodeTo encodes this value using the Encoder.
func (s SequenceNumber) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = Int64(s).EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*SequenceNumber)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *SequenceNumber) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding SequenceNumber: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = (*Int64)(s).DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s SequenceNumber) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *SequenceNumber) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*SequenceNumber)(nil)
	_ encoding.BinaryUnmarshaler = (*SequenceNumber)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s SequenceNumber) xdrType() {}

var _ xdrType = (*SequenceNumber)(nil)

type Preconditions struct {
	Type       PreconditionType
	TimeBounds *TimeBounds
	V2         *PreconditionsV2
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u Preconditions) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of Preconditions
func (u Preconditions) ArmForSwitch(sw int32) (string, bool) {
	switch PreconditionType(sw) {
	case PreconditionTypePrecondNone:
		return "", true
	case PreconditionTypePrecondTime:
		return "TimeBounds", true
	case PreconditionTypePrecondV2:
		return "V2", true
	}
	return "-", false
}

// NewPreconditions creates a new  Preconditions.
func NewPreconditions(aType PreconditionType, value interface{}) (result Preconditions, err error) {
	result.Type = aType
	switch PreconditionType(aType) {
	case PreconditionTypePrecondNone:
		// void
	case PreconditionTypePrecondTime:
		tv, ok := value.(TimeBounds)
		if !ok {
			err = errors.New("invalid value, must be TimeBounds")
			return
		}
		result.TimeBounds = &tv
	case PreconditionTypePrecondV2:
		tv, ok := value.(PreconditionsV2)
		if !ok {
			err = errors.New("invalid value, must be PreconditionsV2")
			return
		}
		result.V2 = &tv
	}
	return
}

// MustTimeBounds retrieves the TimeBounds value from the union,
// panicing if the value is not set.
func (u Preconditions) MustTimeBounds() TimeBounds {
	val, ok := u.GetTimeBounds()

	if !ok {
		panic("arm TimeBounds is not set")
	}

	return val
}

// GetTimeBounds retrieves the TimeBounds value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Preconditions) GetTimeBounds() (result TimeBounds, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "TimeBounds" {
		result = *u.TimeBounds
		ok = true
	}

	return
}

// MustV2 retrieves the V2 value from the union,
// panicing if the value is not set.
func (u Preconditions) MustV2() PreconditionsV2 {
	val, ok := u.GetV2()

	if !ok {
		panic("arm V2 is not set")
	}

	return val
}

// GetV2 retrieves the V2 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Preconditions) GetV2() (result PreconditionsV2, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "V2" {
		result = *u.V2
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u Preconditions) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch PreconditionType(u.Type) {
	case PreconditionTypePrecondNone:
		// Void
		return nil
	case PreconditionTypePrecondTime:
		if err = (*u.TimeBounds).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case PreconditionTypePrecondV2:
		if err = (*u.V2).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (PreconditionType) switch value '%d' is not valid for union Preconditions", u.Type)
}

var _ decoderFrom = (*Preconditions)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *Preconditions) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding Preconditions: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding PreconditionType: %w", err)
	}
	switch PreconditionType(u.Type) {
	case PreconditionTypePrecondNone:
		// Void
		return n, nil
	case PreconditionTypePrecondTime:
		u.TimeBounds = new(TimeBounds)
		nTmp, err = (*u.TimeBounds).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding TimeBounds: %w", err)
		}
		return n, nil
	case PreconditionTypePrecondV2:
		u.V2 = new(PreconditionsV2)
		nTmp, err = (*u.V2).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding PreconditionsV2: %w", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union Preconditions has invalid Type (PreconditionType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Preconditions) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Preconditions) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Preconditions)(nil)
	_ encoding.BinaryUnmarshaler = (*Preconditions)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s Preconditions) xdrType() {}

var _ xdrType = (*Preconditions)(nil)

type TimeBounds struct {
	MinTime TimePoint
	MaxTime TimePoint
}

// EncodeTo encodes this value using the Encoder.
func (s *TimeBounds) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.MinTime.EncodeTo(e); err != nil {
		return err
	}
	if err = s.MaxTime.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*TimeBounds)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TimeBounds) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding TimeBounds: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.MinTime.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TimePoint: %w", err)
	}
	nTmp, err = s.MaxTime.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TimePoint: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TimeBounds) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TimeBounds) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TimeBounds)(nil)
	_ encoding.BinaryUnmarshaler = (*TimeBounds)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s TimeBounds) xdrType() {}

var _ xdrType = (*TimeBounds)(nil)

type PreconditionsV2 struct {
	TimeBounds      *TimeBounds
	LedgerBounds    *LedgerBounds
	MinSeqNum       *SequenceNumber
	MinSeqAge       Duration
	MinSeqLedgerGap Uint32
	ExtraSigners    []SignerKey `xdrmaxsize:"2"`
}

// EncodeTo encodes this value using the Encoder.
func (s *PreconditionsV2) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeBool(s.TimeBounds != nil); err != nil {
		return err
	}
	if s.TimeBounds != nil {
		if err = (*s.TimeBounds).EncodeTo(e); err != nil {
			return err
		}
	}
	if _, err = e.EncodeBool(s.LedgerBounds != nil); err != nil {
		return err
	}
	if s.LedgerBounds != nil {
		if err = (*s.LedgerBounds).EncodeTo(e); err != nil {
			return err
		}
	}
	if _, err = e.EncodeBool(s.MinSeqNum != nil); err != nil {
		return err
	}
	if s.MinSeqNum != nil {
		if err = (*s.MinSeqNum).EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.MinSeqAge.EncodeTo(e); err != nil {
		return err
	}
	if err = s.MinSeqLedgerGap.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.ExtraSigners))); err != nil {
		return err
	}
	for i := 0; i < len(s.ExtraSigners); i++ {
		if err = s.ExtraSigners[i].EncodeTo(e); err != nil {
			return err
		}
	}
	return nil
}

var _ decoderFrom = (*PreconditionsV2)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *PreconditionsV2) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding PreconditionsV2: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	var b bool
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TimeBounds: %w", err)
	}
	s.TimeBounds = nil
	if b {
		s.TimeBounds = new(TimeBounds)
		nTmp, err = s.TimeBounds.DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding TimeBounds: %w", err)
		}
	}
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding LedgerBounds: %w", err)
	}
	s.LedgerBounds = nil
	if b {
		s.LedgerBounds = new(LedgerBounds)
		nTmp, err = s.LedgerBounds.DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding LedgerBounds: %w", err)
		}
	}
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SequenceNumber: %w", err)
	}
	s.MinSeqNum = nil
	if b {
		s.MinSeqNum = new(SequenceNumber)
		nTmp, err = s.MinSeqNum.DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding SequenceNumber: %w", err)
		}
	}
	nTmp, err = s.MinSeqAge.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Duration: %w", err)
	}
	nTmp, err = s.MinSeqLedgerGap.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %w", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SignerKey: %w", err)
	}
	if l > 2 {
		return n, fmt.Errorf("decoding SignerKey: data size (%d) exceeds size limit (2)", l)
	}
	s.ExtraSigners = nil
	if l > 0 {
		if il, ok := d.InputLen(); ok && uint(il) < uint(l) {
			return n, fmt.Errorf("decoding SignerKey: length (%d) exceeds remaining input length (%d)", l, il)
		}
		s.ExtraSigners = make([]SignerKey, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.ExtraSigners[i].DecodeFrom(d, maxDepth)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding SignerKey: %w", err)
			}
		}
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PreconditionsV2) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PreconditionsV2) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PreconditionsV2)(nil)
	_ encoding.BinaryUnmarshaler = (*PreconditionsV2)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s PreconditionsV2) xdrType() {}

var _ xdrType = (*PreconditionsV2)(nil)

type MemoType int32

const (
	MemoTypeMemoNone   MemoType = 0
	MemoTypeMemoText   MemoType = 1
	MemoTypeMemoId     MemoType = 2
	MemoTypeMemoHash   MemoType = 3
	MemoTypeMemoReturn MemoType = 4
)

var memoTypeMap = map[int32]string{
	0: "MemoTypeMemoNone",
	1: "MemoTypeMemoText",
	2: "MemoTypeMemoId",
	3: "MemoTypeMemoHash",
	4: "MemoTypeMemoReturn",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for MemoType
func (e MemoType) ValidEnum(v int32) bool {
	_, ok := memoTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e MemoType) String() string {
	name, _ := memoTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e MemoType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := memoTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid MemoType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*MemoType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *MemoType) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding MemoType: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding MemoType: %w", err)
	}
	if _, ok := memoTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid MemoType enum value", v)
	}
	*e = MemoType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s MemoType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *MemoType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*MemoType)(nil)
	_ encoding.BinaryUnmarshaler = (*MemoType)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s MemoType) xdrType() {}

var _ xdrType = (*MemoType)(nil)

type Memo struct {
	Type    MemoType
	Text    *string `xdrmaxsize:"28"`
	Id      *Uint64
	Hash    *Hash
	RetHash *Hash
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u Memo) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of Memo
func (u Memo) ArmForSwitch(sw int32) (string, bool) {
	switch MemoType(sw) {
	case MemoTypeMemoNone:
		return "", true
	case MemoTypeMemoText:
		return "Text", true
	case MemoTypeMemoId:
		return "Id", true
	case MemoTypeMemoHash:
		return "Hash", true
	case MemoTypeMemoReturn:
		return "RetHash", true
	}
	return "-", false
}

// NewMemo creates a new  Memo.
func NewMemo(aType MemoType, value interface{}) (result Memo, err error) {
	result.Type = aType
	switch MemoType(aType) {
	case MemoTypeMemoNone:
		// void
	case MemoTypeMemoText:
		tv, ok := value.(string)
		if !ok {
			err = errors.New("invalid value, must be string")
			return
		}
		result.Text = &tv
	case MemoTypeMemoId:
		tv, ok := value.(Uint64)
		if !ok {
			err = errors.New("invalid value, must be Uint64")
			return
		}
		result.Id = &tv
	case MemoTypeMemoHash:
		tv, ok := value.(Hash)
		if !ok {
			err = errors.New("invalid value, must be Hash")
			return
		}
		result.Hash = &tv
	case MemoTypeMemoReturn:
		tv, ok := value.(Hash)
		if !ok {
			err = errors.New("invalid value, must be Hash")
			return
		}
		result.RetHash = &tv
	}
	return
}

// MustText retrieves the Text value from the union,
// panicing if the value is not set.
func (u Memo) MustText() string {
	val, ok := u.GetText()

	if !ok {
		panic("arm Text is not set")
	}

	return val
}

// GetText retrieves the Text value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Memo) GetText() (result string, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Text" {
		result = *u.Text
		ok = true
	}

	return
}

// MustId retrieves the Id value from the union,
// panicing if the value is not set.
func (u Memo) MustId() Uint64 {
	val, ok := u.GetId()

	if !ok {
		panic("arm Id is not set")
	}

	return val
}

// GetId retrieves the Id value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Memo) GetId() (result Uint64, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Id" {
		result = *u.Id
		ok = true
	}

	return
}

// MustHash retrieves the Hash value from the union,
// panicing if the value is not set.
func (u Memo) MustHash() Hash {
	val, ok := u.GetHash()

	if !ok {
		panic("arm Hash is not set")
	}

	return val
}

// GetHash retrieves the Hash value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Memo) GetHash() (result Hash, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Hash" {
		result = *u.Hash
		ok = true
	}

	return
}

// MustRetHash retrieves the RetHash value from the union,
// panicing if the value is not set.
func (u Memo) MustRetHash() Hash {
	val, ok := u.GetRetHash()

	if !ok {
		panic("arm RetHash is not set")
	}

	return val
}

// GetRetHash retrieves the RetHash value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u Memo) GetRetHash() (result Hash, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "RetHash" {
		result = *u.RetHash
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u Memo) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch MemoType(u.Type) {
	case MemoTypeMemoNone:
		// Void
		return nil
	case MemoTypeMemoText:
		if _, err = e.EncodeString(string((*u.Text))); err != nil {
			return err
		}
		return nil
	case MemoTypeMemoId:
		if err = (*u.Id).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case MemoTypeMemoHash:
		if err = (*u.Hash).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case MemoTypeMemoReturn:
		if err = (*u.RetHash).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (MemoType) switch value '%d' is not valid for union Memo", u.Type)
}

var _ decoderFrom = (*Memo)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *Memo) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding Memo: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding MemoType: %w", err)
	}
	switch MemoType(u.Type) {
	case MemoTypeMemoNone:
		// Void
		return n, nil
	case MemoTypeMemoText:
		u.Text = new(string)
		(*u.Text), nTmp, err = d.DecodeString(28)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Text: %w", err)
		}
		return n, nil
	case MemoTypeMemoId:
		u.Id = new(Uint64)
		nTmp, err = (*u.Id).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Uint64: %w", err)
		}
		return n, nil
	case MemoTypeMemoHash:
		u.Hash = new(Hash)
		nTmp, err = (*u.Hash).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Hash: %w", err)
		}
		return n, nil
	case MemoTypeMemoReturn:
		u.RetHash = new(Hash)
		nTmp, err = (*u.RetHash).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Hash: %w", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union Memo has invalid Type (MemoType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Memo) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Memo) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Memo)(nil)
	_ encoding.BinaryUnmarshaler = (*Memo)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s Memo) xdrType() {}

var _ xdrType = (*Memo)(nil)

type Operation struct {
	SourceAccount *MuxedAccount
	Body          OperationBody
}

// EncodeTo encodes this value using the Encoder.
func (s *Operation) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeBool(s.SourceAccount != nil); err != nil {
		return err
	}
	if s.SourceAccount != nil {
		if err = (*s.SourceAccount).EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.Body.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Operation)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Operation) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding Operation: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	var b bool
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding MuxedAccount: %w", err)
	}
	s.SourceAccount = nil
	if b {
		s.SourceAccount = new(MuxedAccount)
		nTmp, err = s.SourceAccount.DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding MuxedAccount: %w", err)
		}
	}
	nTmp, err = s.Body.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding OperationBody: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Operation) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Operation) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Operation)(nil)
	_ encoding.BinaryUnmarshaler = (*Operation)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s Operation) xdrType() {}

var _ xdrType = (*Operation)(nil)

type TransactionExt struct {
	V int32
	//SorobanData *SorobanTransactionData
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionExt
func (u TransactionExt) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	case 1:
		return "SorobanData", true
	}
	return "-", false
}

// NewTransactionExt creates a new  TransactionExt.
func NewTransactionExt(v int32, value interface{}) (result TransactionExt, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	case 1:
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u TransactionExt) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union TransactionExt", u.V)
}

var _ decoderFrom = (*TransactionExt)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *TransactionExt) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding TransactionExt: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %w", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	case 1:
	}
	return n, fmt.Errorf("union TransactionExt has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionExt) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionExt) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionExt)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionExt)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s TransactionExt) xdrType() {}

var _ xdrType = (*TransactionExt)(nil)

type CryptoKeyType int32

const (
	CryptoKeyTypeKeyTypeEd25519              CryptoKeyType = 0
	CryptoKeyTypeKeyTypePreAuthTx            CryptoKeyType = 1
	CryptoKeyTypeKeyTypeHashX                CryptoKeyType = 2
	CryptoKeyTypeKeyTypeEd25519SignedPayload CryptoKeyType = 3
	CryptoKeyTypeKeyTypeMuxedEd25519         CryptoKeyType = 256
)

var cryptoKeyTypeMap = map[int32]string{
	0:   "CryptoKeyTypeKeyTypeEd25519",
	1:   "CryptoKeyTypeKeyTypePreAuthTx",
	2:   "CryptoKeyTypeKeyTypeHashX",
	3:   "CryptoKeyTypeKeyTypeEd25519SignedPayload",
	256: "CryptoKeyTypeKeyTypeMuxedEd25519",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for CryptoKeyType
func (e CryptoKeyType) ValidEnum(v int32) bool {
	_, ok := cryptoKeyTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e CryptoKeyType) String() string {
	name, _ := cryptoKeyTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e CryptoKeyType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := cryptoKeyTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid CryptoKeyType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*CryptoKeyType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *CryptoKeyType) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding CryptoKeyType: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding CryptoKeyType: %w", err)
	}
	if _, ok := cryptoKeyTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid CryptoKeyType enum value", v)
	}
	*e = CryptoKeyType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s CryptoKeyType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *CryptoKeyType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*CryptoKeyType)(nil)
	_ encoding.BinaryUnmarshaler = (*CryptoKeyType)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s CryptoKeyType) xdrType() {}

var _ xdrType = (*CryptoKeyType)(nil)

type MuxedAccountMed25519 struct {
	Id      Uint64
	Ed25519 Uint256
}

// EncodeTo encodes this value using the Encoder.
func (s *MuxedAccountMed25519) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Id.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ed25519.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*MuxedAccountMed25519)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *MuxedAccountMed25519) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding MuxedAccountMed25519: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.Id.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %w", err)
	}
	nTmp, err = s.Ed25519.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint256: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s MuxedAccountMed25519) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *MuxedAccountMed25519) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*MuxedAccountMed25519)(nil)
	_ encoding.BinaryUnmarshaler = (*MuxedAccountMed25519)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s MuxedAccountMed25519) xdrType() {}

var _ xdrType = (*MuxedAccountMed25519)(nil)

type LedgerBounds struct {
	MinLedger Uint32
	MaxLedger Uint32
}

// EncodeTo encodes this value using the Encoder.
func (s *LedgerBounds) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.MinLedger.EncodeTo(e); err != nil {
		return err
	}
	if err = s.MaxLedger.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*LedgerBounds)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *LedgerBounds) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding LedgerBounds: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.MinLedger.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %w", err)
	}
	nTmp, err = s.MaxLedger.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s LedgerBounds) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *LedgerBounds) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*LedgerBounds)(nil)
	_ encoding.BinaryUnmarshaler = (*LedgerBounds)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s LedgerBounds) xdrType() {}

var _ xdrType = (*LedgerBounds)(nil)

type OperationBody struct {
	Type          OperationType
	PaymentOp     *PaymentOp
	ChangeTrustOp *ChangeTrustOp
	AllowTrustOp  *AllowTrustOp
	Destination   *MuxedAccount
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u OperationBody) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of OperationBody
func (u OperationBody) ArmForSwitch(sw int32) (string, bool) {
	switch OperationType(sw) {
	case OperationTypePayment:
		return "PaymentOp", true
	case OperationTypePathPaymentStrictReceive:
		return "PathPaymentStrictReceiveOp", true
	case OperationTypeManageSellOffer:
		return "ManageSellOfferOp", true
	case OperationTypeCreatePassiveSellOffer:
		return "CreatePassiveSellOfferOp", true
	case OperationTypeSetOptions:
		return "SetOptionsOp", true
	case OperationTypeChangeTrust:
		return "ChangeTrustOp", true
	case OperationTypeAllowTrust:
		return "AllowTrustOp", true
	case OperationTypeAccountMerge:
		return "Destination", true
	case OperationTypeInflation:
		return "", true
	}
	return "-", false
}

// NewOperationBody creates a new  OperationBody.
func NewOperationBody(aType OperationType, value interface{}) (result OperationBody, err error) {
	result.Type = aType
	switch OperationType(aType) {
	case OperationTypePayment:
		tv, ok := value.(PaymentOp)
		if !ok {
			err = errors.New("invalid value, must be PaymentOp")
			return
		}
		result.PaymentOp = &tv
	case OperationTypeChangeTrust:
		tv, ok := value.(ChangeTrustOp)
		if !ok {
			err = errors.New("invalid value, must be ChangeTrustOp")
			return
		}
		result.ChangeTrustOp = &tv
	case OperationTypeAllowTrust:
		tv, ok := value.(AllowTrustOp)
		if !ok {
			err = errors.New("invalid value, must be AllowTrustOp")
			return
		}
		result.AllowTrustOp = &tv
	case OperationTypeAccountMerge:
		tv, ok := value.(MuxedAccount)
		if !ok {
			err = errors.New("invalid value, must be MuxedAccount")
			return
		}
		result.Destination = &tv
	}
	return
}

// MustPaymentOp retrieves the PaymentOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustPaymentOp() PaymentOp {
	val, ok := u.GetPaymentOp()

	if !ok {
		panic("arm PaymentOp is not set")
	}

	return val
}

// GetPaymentOp retrieves the PaymentOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetPaymentOp() (result PaymentOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "PaymentOp" {
		result = *u.PaymentOp
		ok = true
	}

	return
}

// MustChangeTrustOp retrieves the ChangeTrustOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustChangeTrustOp() ChangeTrustOp {
	val, ok := u.GetChangeTrustOp()

	if !ok {
		panic("arm ChangeTrustOp is not set")
	}

	return val
}

// GetChangeTrustOp retrieves the ChangeTrustOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetChangeTrustOp() (result ChangeTrustOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "ChangeTrustOp" {
		result = *u.ChangeTrustOp
		ok = true
	}

	return
}

// MustAllowTrustOp retrieves the AllowTrustOp value from the union,
// panicing if the value is not set.
func (u OperationBody) MustAllowTrustOp() AllowTrustOp {
	val, ok := u.GetAllowTrustOp()

	if !ok {
		panic("arm AllowTrustOp is not set")
	}

	return val
}

// GetAllowTrustOp retrieves the AllowTrustOp value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetAllowTrustOp() (result AllowTrustOp, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "AllowTrustOp" {
		result = *u.AllowTrustOp
		ok = true
	}

	return
}

// MustDestination retrieves the Destination value from the union,
// panicing if the value is not set.
func (u OperationBody) MustDestination() MuxedAccount {
	val, ok := u.GetDestination()

	if !ok {
		panic("arm Destination is not set")
	}

	return val
}

// GetDestination retrieves the Destination value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u OperationBody) GetDestination() (result MuxedAccount, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Destination" {
		result = *u.Destination
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u OperationBody) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch OperationType(u.Type) {
	case OperationTypePayment:
		if err = (*u.PaymentOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeChangeTrust:
		if err = (*u.ChangeTrustOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeAllowTrust:
		if err = (*u.AllowTrustOp).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeAccountMerge:
		if err = (*u.Destination).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case OperationTypeInflation:
		// Void
		return nil
	}
	return fmt.Errorf("Type (OperationType) switch value '%d' is not valid for union OperationBody", u.Type)
}

var _ decoderFrom = (*OperationBody)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *OperationBody) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding OperationBody: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding OperationType: %w", err)
	}
	switch OperationType(u.Type) {
	case OperationTypePayment:
		u.PaymentOp = new(PaymentOp)
		nTmp, err = (*u.PaymentOp).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding PaymentOp: %w", err)
		}
		return n, nil
	case OperationTypeChangeTrust:
		u.ChangeTrustOp = new(ChangeTrustOp)
		nTmp, err = (*u.ChangeTrustOp).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding ChangeTrustOp: %w", err)
		}
		return n, nil
	case OperationTypeAllowTrust:
		u.AllowTrustOp = new(AllowTrustOp)
		nTmp, err = (*u.AllowTrustOp).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding AllowTrustOp: %w", err)
		}
		return n, nil
	case OperationTypeAccountMerge:
		u.Destination = new(MuxedAccount)
		nTmp, err = (*u.Destination).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding MuxedAccount: %w", err)
		}
		return n, nil
	case OperationTypeInflation:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union OperationBody has invalid Type (OperationType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s OperationBody) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *OperationBody) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*OperationBody)(nil)
	_ encoding.BinaryUnmarshaler = (*OperationBody)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s OperationBody) xdrType() {}

var _ xdrType = (*OperationBody)(nil)

type PaymentOp struct {
	Destination MuxedAccount
	Asset       Asset
	Amount      Int64
}

// EncodeTo encodes this value using the Encoder.
func (s *PaymentOp) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Destination.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Asset.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Amount.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*PaymentOp)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *PaymentOp) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding PaymentOp: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.Destination.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding MuxedAccount: %w", err)
	}
	nTmp, err = s.Asset.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Asset: %w", err)
	}
	nTmp, err = s.Amount.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s PaymentOp) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *PaymentOp) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*PaymentOp)(nil)
	_ encoding.BinaryUnmarshaler = (*PaymentOp)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s PaymentOp) xdrType() {}

var _ xdrType = (*PaymentOp)(nil)

type DecoratedSignature struct {
	Hint      SignatureHint
	Signature Signature
}

// EncodeTo encodes this value using the Encoder.
func (s *DecoratedSignature) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Hint.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Signature.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*DecoratedSignature)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *DecoratedSignature) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding DecoratedSignature: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.Hint.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SignatureHint: %w", err)
	}
	nTmp, err = s.Signature.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Signature: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s DecoratedSignature) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *DecoratedSignature) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*DecoratedSignature)(nil)
	_ encoding.BinaryUnmarshaler = (*DecoratedSignature)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s DecoratedSignature) xdrType() {}

var _ xdrType = (*DecoratedSignature)(nil)

type TransactionEnvelope struct {
	Type    EnvelopeType
	V0      *TransactionV0Envelope
	V1      *TransactionV1Envelope
	FeeBump *FeeBumpTransactionEnvelope
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionEnvelope) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionEnvelope
func (u TransactionEnvelope) ArmForSwitch(sw int32) (string, bool) {
	switch EnvelopeType(sw) {
	case EnvelopeTypeEnvelopeTypeTxV0:
		return "V0", true
	case EnvelopeTypeEnvelopeTypeTx:
		return "V1", true
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		return "FeeBump", true
	}
	return "-", false
}

// NewTransactionEnvelope creates a new  TransactionEnvelope.
func NewTransactionEnvelope(aType EnvelopeType, value interface{}) (result TransactionEnvelope, err error) {
	result.Type = aType
	switch EnvelopeType(aType) {
	case EnvelopeTypeEnvelopeTypeTxV0:
		tv, ok := value.(TransactionV0Envelope)
		if !ok {
			err = errors.New("invalid value, must be TransactionV0Envelope")
			return
		}
		result.V0 = &tv
	case EnvelopeTypeEnvelopeTypeTx:
		tv, ok := value.(TransactionV1Envelope)
		if !ok {
			err = errors.New("invalid value, must be TransactionV1Envelope")
			return
		}
		result.V1 = &tv
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		tv, ok := value.(FeeBumpTransactionEnvelope)
		if !ok {
			err = errors.New("invalid value, must be FeeBumpTransactionEnvelope")
			return
		}
		result.FeeBump = &tv
	}
	return
}

// MustV0 retrieves the V0 value from the union,
// panicing if the value is not set.
func (u TransactionEnvelope) MustV0() TransactionV0Envelope {
	val, ok := u.GetV0()

	if !ok {
		panic("arm V0 is not set")
	}

	return val
}

// GetV0 retrieves the V0 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TransactionEnvelope) GetV0() (result TransactionV0Envelope, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "V0" {
		result = *u.V0
		ok = true
	}

	return
}

// MustV1 retrieves the V1 value from the union,
// panicing if the value is not set.
func (u TransactionEnvelope) MustV1() TransactionV1Envelope {
	val, ok := u.GetV1()

	if !ok {
		panic("arm V1 is not set")
	}

	return val
}

// GetV1 retrieves the V1 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TransactionEnvelope) GetV1() (result TransactionV1Envelope, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "V1" {
		result = *u.V1
		ok = true
	}

	return
}

// MustFeeBump retrieves the FeeBump value from the union,
// panicing if the value is not set.
func (u TransactionEnvelope) MustFeeBump() FeeBumpTransactionEnvelope {
	val, ok := u.GetFeeBump()

	if !ok {
		panic("arm FeeBump is not set")
	}

	return val
}

// GetFeeBump retrieves the FeeBump value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TransactionEnvelope) GetFeeBump() (result FeeBumpTransactionEnvelope, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "FeeBump" {
		result = *u.FeeBump
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u TransactionEnvelope) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch EnvelopeType(u.Type) {
	case EnvelopeTypeEnvelopeTypeTxV0:
		if err = (*u.V0).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case EnvelopeTypeEnvelopeTypeTx:
		if err = (*u.V1).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		if err = (*u.FeeBump).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (EnvelopeType) switch value '%d' is not valid for union TransactionEnvelope", u.Type)
}

var _ decoderFrom = (*TransactionEnvelope)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *TransactionEnvelope) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding TransactionEnvelope: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding EnvelopeType: %w", err)
	}
	switch EnvelopeType(u.Type) {
	case EnvelopeTypeEnvelopeTypeTxV0:
		u.V0 = new(TransactionV0Envelope)
		nTmp, err = (*u.V0).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding TransactionV0Envelope: %w", err)
		}
		return n, nil
	case EnvelopeTypeEnvelopeTypeTx:
		u.V1 = new(TransactionV1Envelope)
		nTmp, err = (*u.V1).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding TransactionV1Envelope: %w", err)
		}
		return n, nil
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		u.FeeBump = new(FeeBumpTransactionEnvelope)
		nTmp, err = (*u.FeeBump).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding FeeBumpTransactionEnvelope: %w", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union TransactionEnvelope has invalid Type (EnvelopeType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionEnvelope) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionEnvelope) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionEnvelope)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionEnvelope)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s TransactionEnvelope) xdrType() {}

var _ xdrType = (*TransactionEnvelope)(nil)

type TransactionV0Envelope struct {
	Tx         TransactionV0
	Signatures []DecoratedSignature `xdrmaxsize:"20"`
}

// EncodeTo encodes this value using the Encoder.
func (s *TransactionV0Envelope) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Tx.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.Signatures))); err != nil {
		return err
	}
	for i := 0; i < len(s.Signatures); i++ {
		if err = s.Signatures[i].EncodeTo(e); err != nil {
			return err
		}
	}
	return nil
}

var _ decoderFrom = (*TransactionV0Envelope)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TransactionV0Envelope) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding TransactionV0Envelope: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.Tx.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionV0: %w", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding DecoratedSignature: %w", err)
	}
	if l > 20 {
		return n, fmt.Errorf("decoding DecoratedSignature: data size (%d) exceeds size limit (20)", l)
	}
	s.Signatures = nil
	if l > 0 {
		if il, ok := d.InputLen(); ok && uint(il) < uint(l) {
			return n, fmt.Errorf("decoding DecoratedSignature: length (%d) exceeds remaining input length (%d)", l, il)
		}
		s.Signatures = make([]DecoratedSignature, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Signatures[i].DecodeFrom(d, maxDepth)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding DecoratedSignature: %w", err)
			}
		}
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionV0Envelope) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionV0Envelope) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionV0Envelope)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionV0Envelope)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s TransactionV0Envelope) xdrType() {}

var _ xdrType = (*TransactionV0Envelope)(nil)

type TransactionV1Envelope struct {
	Tx         Transaction
	Signatures []DecoratedSignature `xdrmaxsize:"20"`
}

// EncodeTo encodes this value using the Encoder.
func (s *TransactionV1Envelope) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Tx.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.Signatures))); err != nil {
		return err
	}
	for i := 0; i < len(s.Signatures); i++ {
		if err = s.Signatures[i].EncodeTo(e); err != nil {
			return err
		}
	}
	return nil
}

var _ decoderFrom = (*TransactionV1Envelope)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TransactionV1Envelope) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding TransactionV1Envelope: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.Tx.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Transaction: %w", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding DecoratedSignature: %w", err)
	}
	if l > 20 {
		return n, fmt.Errorf("decoding DecoratedSignature: data size (%d) exceeds size limit (20)", l)
	}
	s.Signatures = nil
	if l > 0 {
		if il, ok := d.InputLen(); ok && uint(il) < uint(l) {
			return n, fmt.Errorf("decoding DecoratedSignature: length (%d) exceeds remaining input length (%d)", l, il)
		}
		s.Signatures = make([]DecoratedSignature, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Signatures[i].DecodeFrom(d, maxDepth)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding DecoratedSignature: %w", err)
			}
		}
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionV1Envelope) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionV1Envelope) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionV1Envelope)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionV1Envelope)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s TransactionV1Envelope) xdrType() {}

var _ xdrType = (*TransactionV1Envelope)(nil)

type EnvelopeType int32

const (
	EnvelopeTypeEnvelopeTypeTxV0                 EnvelopeType = 0
	EnvelopeTypeEnvelopeTypeScp                  EnvelopeType = 1
	EnvelopeTypeEnvelopeTypeTx                   EnvelopeType = 2
	EnvelopeTypeEnvelopeTypeAuth                 EnvelopeType = 3
	EnvelopeTypeEnvelopeTypeScpvalue             EnvelopeType = 4
	EnvelopeTypeEnvelopeTypeTxFeeBump            EnvelopeType = 5
	EnvelopeTypeEnvelopeTypeOpId                 EnvelopeType = 6
	EnvelopeTypeEnvelopeTypePoolRevokeOpId       EnvelopeType = 7
	EnvelopeTypeEnvelopeTypeContractId           EnvelopeType = 8
	EnvelopeTypeEnvelopeTypeSorobanAuthorization EnvelopeType = 9
)

var envelopeTypeMap = map[int32]string{
	0: "EnvelopeTypeEnvelopeTypeTxV0",
	1: "EnvelopeTypeEnvelopeTypeScp",
	2: "EnvelopeTypeEnvelopeTypeTx",
	3: "EnvelopeTypeEnvelopeTypeAuth",
	4: "EnvelopeTypeEnvelopeTypeScpvalue",
	5: "EnvelopeTypeEnvelopeTypeTxFeeBump",
	6: "EnvelopeTypeEnvelopeTypeOpId",
	7: "EnvelopeTypeEnvelopeTypePoolRevokeOpId",
	8: "EnvelopeTypeEnvelopeTypeContractId",
	9: "EnvelopeTypeEnvelopeTypeSorobanAuthorization",
}

// ValidEnum validates a proposed value for this enum.  Implements
// the Enum interface for EnvelopeType
func (e EnvelopeType) ValidEnum(v int32) bool {
	_, ok := envelopeTypeMap[v]
	return ok
}

// String returns the name of `e`
func (e EnvelopeType) String() string {
	name, _ := envelopeTypeMap[int32(e)]
	return name
}

// EncodeTo encodes this value using the Encoder.
func (e EnvelopeType) EncodeTo(enc *xdr.Encoder) error {
	if _, ok := envelopeTypeMap[int32(e)]; !ok {
		return fmt.Errorf("'%d' is not a valid EnvelopeType enum value", e)
	}
	_, err := enc.EncodeInt(int32(e))
	return err
}

var _ decoderFrom = (*EnvelopeType)(nil)

// DecodeFrom decodes this value using the Decoder.
func (e *EnvelopeType) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding EnvelopeType: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	v, n, err := d.DecodeInt()
	if err != nil {
		return n, fmt.Errorf("decoding EnvelopeType: %w", err)
	}
	if _, ok := envelopeTypeMap[v]; !ok {
		return n, fmt.Errorf("'%d' is not a valid EnvelopeType enum value", v)
	}
	*e = EnvelopeType(v)
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s EnvelopeType) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *EnvelopeType) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*EnvelopeType)(nil)
	_ encoding.BinaryUnmarshaler = (*EnvelopeType)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s EnvelopeType) xdrType() {}

var _ xdrType = (*EnvelopeType)(nil)

type FeeBumpTransactionEnvelope struct {
	Tx         FeeBumpTransaction
	Signatures []DecoratedSignature `xdrmaxsize:"20"`
}

// EncodeTo encodes this value using the Encoder.
func (s *FeeBumpTransactionEnvelope) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Tx.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.Signatures))); err != nil {
		return err
	}
	for i := 0; i < len(s.Signatures); i++ {
		if err = s.Signatures[i].EncodeTo(e); err != nil {
			return err
		}
	}
	return nil
}

var _ decoderFrom = (*FeeBumpTransactionEnvelope)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *FeeBumpTransactionEnvelope) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding FeeBumpTransactionEnvelope: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.Tx.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding FeeBumpTransaction: %w", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding DecoratedSignature: %w", err)
	}
	if l > 20 {
		return n, fmt.Errorf("decoding DecoratedSignature: data size (%d) exceeds size limit (20)", l)
	}
	s.Signatures = nil
	if l > 0 {
		if il, ok := d.InputLen(); ok && uint(il) < uint(l) {
			return n, fmt.Errorf("decoding DecoratedSignature: length (%d) exceeds remaining input length (%d)", l, il)
		}
		s.Signatures = make([]DecoratedSignature, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Signatures[i].DecodeFrom(d, maxDepth)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding DecoratedSignature: %w", err)
			}
		}
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s FeeBumpTransactionEnvelope) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *FeeBumpTransactionEnvelope) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*FeeBumpTransactionEnvelope)(nil)
	_ encoding.BinaryUnmarshaler = (*FeeBumpTransactionEnvelope)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s FeeBumpTransactionEnvelope) xdrType() {}

var _ xdrType = (*FeeBumpTransactionEnvelope)(nil)

type TransactionV0 struct {
	SourceAccountEd25519 Uint256
	Fee                  Uint32
	SeqNum               SequenceNumber
	TimeBounds           *TimeBounds
	Memo                 Memo
	Operations           []Operation `xdrmaxsize:"100"`
	Ext                  TransactionV0Ext
}

// EncodeTo encodes this value using the Encoder.
func (s *TransactionV0) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.SourceAccountEd25519.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Fee.EncodeTo(e); err != nil {
		return err
	}
	if err = s.SeqNum.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeBool(s.TimeBounds != nil); err != nil {
		return err
	}
	if s.TimeBounds != nil {
		if err = (*s.TimeBounds).EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.Memo.EncodeTo(e); err != nil {
		return err
	}
	if _, err = e.EncodeUint(uint32(len(s.Operations))); err != nil {
		return err
	}
	for i := 0; i < len(s.Operations); i++ {
		if err = s.Operations[i].EncodeTo(e); err != nil {
			return err
		}
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*TransactionV0)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TransactionV0) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding TransactionV0: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.SourceAccountEd25519.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint256: %w", err)
	}
	nTmp, err = s.Fee.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %w", err)
	}
	nTmp, err = s.SeqNum.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SequenceNumber: %w", err)
	}
	var b bool
	b, nTmp, err = d.DecodeBool()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TimeBounds: %w", err)
	}
	s.TimeBounds = nil
	if b {
		s.TimeBounds = new(TimeBounds)
		nTmp, err = s.TimeBounds.DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding TimeBounds: %w", err)
		}
	}
	nTmp, err = s.Memo.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Memo: %w", err)
	}
	var l uint32
	l, nTmp, err = d.DecodeUint()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Operation: %w", err)
	}
	if l > 100 {
		return n, fmt.Errorf("decoding Operation: data size (%d) exceeds size limit (100)", l)
	}
	s.Operations = nil
	if l > 0 {
		if il, ok := d.InputLen(); ok && uint(il) < uint(l) {
			return n, fmt.Errorf("decoding Operation: length (%d) exceeds remaining input length (%d)", l, il)
		}
		s.Operations = make([]Operation, l)
		for i := uint32(0); i < l; i++ {
			nTmp, err = s.Operations[i].DecodeFrom(d, maxDepth)
			n += nTmp
			if err != nil {
				return n, fmt.Errorf("decoding Operation: %w", err)
			}
		}
	}
	nTmp, err = s.Ext.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionV0Ext: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionV0) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionV0) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionV0)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionV0)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s TransactionV0) xdrType() {}

var _ xdrType = (*TransactionV0)(nil)

type FeeBumpTransaction struct {
	FeeSource MuxedAccount
	Fee       Int64
	InnerTx   FeeBumpTransactionInnerTx
	Ext       FeeBumpTransactionExt
}

// EncodeTo encodes this value using the Encoder.
func (s *FeeBumpTransaction) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.FeeSource.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Fee.EncodeTo(e); err != nil {
		return err
	}
	if err = s.InnerTx.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Ext.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*FeeBumpTransaction)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *FeeBumpTransaction) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding FeeBumpTransaction: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.FeeSource.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding MuxedAccount: %w", err)
	}
	nTmp, err = s.Fee.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %w", err)
	}
	nTmp, err = s.InnerTx.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding FeeBumpTransactionInnerTx: %w", err)
	}
	nTmp, err = s.Ext.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding FeeBumpTransactionExt: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s FeeBumpTransaction) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *FeeBumpTransaction) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*FeeBumpTransaction)(nil)
	_ encoding.BinaryUnmarshaler = (*FeeBumpTransaction)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s FeeBumpTransaction) xdrType() {}

var _ xdrType = (*FeeBumpTransaction)(nil)

type TransactionV0Ext struct {
	V int32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionV0Ext) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionV0Ext
func (u TransactionV0Ext) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	}
	return "-", false
}

// NewTransactionV0Ext creates a new  TransactionV0Ext.
func NewTransactionV0Ext(v int32, value interface{}) (result TransactionV0Ext, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u TransactionV0Ext) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union TransactionV0Ext", u.V)
}

var _ decoderFrom = (*TransactionV0Ext)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *TransactionV0Ext) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding TransactionV0Ext: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %w", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union TransactionV0Ext has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionV0Ext) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionV0Ext) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionV0Ext)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionV0Ext)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s TransactionV0Ext) xdrType() {}

var _ xdrType = (*TransactionV0Ext)(nil)

type FeeBumpTransactionInnerTx struct {
	Type EnvelopeType
	V1   *TransactionV1Envelope
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u FeeBumpTransactionInnerTx) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of FeeBumpTransactionInnerTx
func (u FeeBumpTransactionInnerTx) ArmForSwitch(sw int32) (string, bool) {
	switch EnvelopeType(sw) {
	case EnvelopeTypeEnvelopeTypeTx:
		return "V1", true
	}
	return "-", false
}

// NewFeeBumpTransactionInnerTx creates a new  FeeBumpTransactionInnerTx.
func NewFeeBumpTransactionInnerTx(aType EnvelopeType, value interface{}) (result FeeBumpTransactionInnerTx, err error) {
	result.Type = aType
	switch EnvelopeType(aType) {
	case EnvelopeTypeEnvelopeTypeTx:
		tv, ok := value.(TransactionV1Envelope)
		if !ok {
			err = errors.New("invalid value, must be TransactionV1Envelope")
			return
		}
		result.V1 = &tv
	}
	return
}

// MustV1 retrieves the V1 value from the union,
// panicing if the value is not set.
func (u FeeBumpTransactionInnerTx) MustV1() TransactionV1Envelope {
	val, ok := u.GetV1()

	if !ok {
		panic("arm V1 is not set")
	}

	return val
}

// GetV1 retrieves the V1 value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u FeeBumpTransactionInnerTx) GetV1() (result TransactionV1Envelope, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "V1" {
		result = *u.V1
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u FeeBumpTransactionInnerTx) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch EnvelopeType(u.Type) {
	case EnvelopeTypeEnvelopeTypeTx:
		if err = (*u.V1).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (EnvelopeType) switch value '%d' is not valid for union FeeBumpTransactionInnerTx", u.Type)
}

var _ decoderFrom = (*FeeBumpTransactionInnerTx)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *FeeBumpTransactionInnerTx) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding FeeBumpTransactionInnerTx: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding EnvelopeType: %w", err)
	}
	switch EnvelopeType(u.Type) {
	case EnvelopeTypeEnvelopeTypeTx:
		u.V1 = new(TransactionV1Envelope)
		nTmp, err = (*u.V1).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding TransactionV1Envelope: %w", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union FeeBumpTransactionInnerTx has invalid Type (EnvelopeType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s FeeBumpTransactionInnerTx) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *FeeBumpTransactionInnerTx) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*FeeBumpTransactionInnerTx)(nil)
	_ encoding.BinaryUnmarshaler = (*FeeBumpTransactionInnerTx)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s FeeBumpTransactionInnerTx) xdrType() {}

var _ xdrType = (*FeeBumpTransactionInnerTx)(nil)

type FeeBumpTransactionExt struct {
	V int32
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u FeeBumpTransactionExt) SwitchFieldName() string {
	return "V"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of FeeBumpTransactionExt
func (u FeeBumpTransactionExt) ArmForSwitch(sw int32) (string, bool) {
	switch int32(sw) {
	case 0:
		return "", true
	}
	return "-", false
}

// NewFeeBumpTransactionExt creates a new  FeeBumpTransactionExt.
func NewFeeBumpTransactionExt(v int32, value interface{}) (result FeeBumpTransactionExt, err error) {
	result.V = v
	switch int32(v) {
	case 0:
		// void
	}
	return
}

// EncodeTo encodes this value using the Encoder.
func (u FeeBumpTransactionExt) EncodeTo(e *xdr.Encoder) error {
	var err error
	if _, err = e.EncodeInt(int32(u.V)); err != nil {
		return err
	}
	switch int32(u.V) {
	case 0:
		// Void
		return nil
	}
	return fmt.Errorf("V (int32) switch value '%d' is not valid for union FeeBumpTransactionExt", u.V)
}

var _ decoderFrom = (*FeeBumpTransactionExt)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *FeeBumpTransactionExt) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding FeeBumpTransactionExt: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	u.V, nTmp, err = d.DecodeInt()
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int: %w", err)
	}
	switch int32(u.V) {
	case 0:
		// Void
		return n, nil
	}
	return n, fmt.Errorf("union FeeBumpTransactionExt has invalid V (int32) switch value '%d'", u.V)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s FeeBumpTransactionExt) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *FeeBumpTransactionExt) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*FeeBumpTransactionExt)(nil)
	_ encoding.BinaryUnmarshaler = (*FeeBumpTransactionExt)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s FeeBumpTransactionExt) xdrType() {}

var _ xdrType = (*FeeBumpTransactionExt)(nil)

type Signer struct {
	Key    SignerKey
	Weight Uint32
}

// EncodeTo encodes this value using the Encoder.
func (s *Signer) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Key.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Weight.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Signer)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Signer) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding Signer: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.Key.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding SignerKey: %w", err)
	}
	nTmp, err = s.Weight.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint32: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Signer) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Signer) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Signer)(nil)
	_ encoding.BinaryUnmarshaler = (*Signer)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s Signer) xdrType() {}

var _ xdrType = (*Signer)(nil)

// Marshal writes an xdr element `v` into `w`.
func Marshal(w io.Writer, v interface{}) (int, error) {
	if _, ok := v.(xdrType); ok {
		if bm, ok := v.(encoding.BinaryMarshaler); ok {
			b, err := bm.MarshalBinary()
			if err != nil {
				return 0, err
			}
			return w.Write(b)
		}
	}
	// delegate to xdr package's Marshal
	return xdr.Marshal(w, v)
}

// Unmarshal reads an xdr element from `r` into `v`.
func Unmarshal(r io.Reader, v interface{}) (int, error) {
	return UnmarshalWithOptions(r, v, xdr.DefaultDecodeOptions)
}

// UnmarshalWithOptions works like Unmarshal but uses decoding options.
func UnmarshalWithOptions(r io.Reader, v interface{}, options xdr.DecodeOptions) (int, error) {
	if decodable, ok := v.(decoderFrom); ok {
		d := xdr.NewDecoderWithOptions(r, options)
		return decodable.DecodeFrom(d, options.MaxDepth)
	}
	// delegate to xdr package's Unmarshal
	return xdr.UnmarshalWithOptions(r, v, options)
}

type TransactionSignaturePayloadTaggedTransaction struct {
	Type    EnvelopeType
	Tx      *Transaction
	FeeBump *FeeBumpTransaction
}

// SwitchFieldName returns the field name in which this union's
// discriminant is stored
func (u TransactionSignaturePayloadTaggedTransaction) SwitchFieldName() string {
	return "Type"
}

// ArmForSwitch returns which field name should be used for storing
// the value for an instance of TransactionSignaturePayloadTaggedTransaction
func (u TransactionSignaturePayloadTaggedTransaction) ArmForSwitch(sw int32) (string, bool) {
	switch EnvelopeType(sw) {
	case EnvelopeTypeEnvelopeTypeTx:
		return "Tx", true
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		return "FeeBump", true
	}
	return "-", false
}

// NewTransactionSignaturePayloadTaggedTransaction creates a new  TransactionSignaturePayloadTaggedTransaction.
func NewTransactionSignaturePayloadTaggedTransaction(aType EnvelopeType, value interface{}) (result TransactionSignaturePayloadTaggedTransaction, err error) {
	result.Type = aType
	switch EnvelopeType(aType) {
	case EnvelopeTypeEnvelopeTypeTx:
		tv, ok := value.(Transaction)
		if !ok {
			err = errors.New("invalid value, must be Transaction")
			return
		}
		result.Tx = &tv
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		tv, ok := value.(FeeBumpTransaction)
		if !ok {
			err = errors.New("invalid value, must be FeeBumpTransaction")
			return
		}
		result.FeeBump = &tv
	}
	return
}

// MustTx retrieves the Tx value from the union,
// panicing if the value is not set.
func (u TransactionSignaturePayloadTaggedTransaction) MustTx() Transaction {
	val, ok := u.GetTx()

	if !ok {
		panic("arm Tx is not set")
	}

	return val
}

// GetTx retrieves the Tx value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TransactionSignaturePayloadTaggedTransaction) GetTx() (result Transaction, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "Tx" {
		result = *u.Tx
		ok = true
	}

	return
}

// MustFeeBump retrieves the FeeBump value from the union,
// panicing if the value is not set.
func (u TransactionSignaturePayloadTaggedTransaction) MustFeeBump() FeeBumpTransaction {
	val, ok := u.GetFeeBump()

	if !ok {
		panic("arm FeeBump is not set")
	}

	return val
}

// GetFeeBump retrieves the FeeBump value from the union,
// returning ok if the union's switch indicated the value is valid.
func (u TransactionSignaturePayloadTaggedTransaction) GetFeeBump() (result FeeBumpTransaction, ok bool) {
	armName, _ := u.ArmForSwitch(int32(u.Type))

	if armName == "FeeBump" {
		result = *u.FeeBump
		ok = true
	}

	return
}

// EncodeTo encodes this value using the Encoder.
func (u TransactionSignaturePayloadTaggedTransaction) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = u.Type.EncodeTo(e); err != nil {
		return err
	}
	switch EnvelopeType(u.Type) {
	case EnvelopeTypeEnvelopeTypeTx:
		if err = (*u.Tx).EncodeTo(e); err != nil {
			return err
		}
		return nil
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		if err = (*u.FeeBump).EncodeTo(e); err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Type (EnvelopeType) switch value '%d' is not valid for union TransactionSignaturePayloadTaggedTransaction", u.Type)
}

var _ decoderFrom = (*TransactionSignaturePayloadTaggedTransaction)(nil)

// DecodeFrom decodes this value using the Decoder.
func (u *TransactionSignaturePayloadTaggedTransaction) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding TransactionSignaturePayloadTaggedTransaction: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = u.Type.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding EnvelopeType: %w", err)
	}
	switch EnvelopeType(u.Type) {
	case EnvelopeTypeEnvelopeTypeTx:
		u.Tx = new(Transaction)
		nTmp, err = (*u.Tx).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding Transaction: %w", err)
		}
		return n, nil
	case EnvelopeTypeEnvelopeTypeTxFeeBump:
		u.FeeBump = new(FeeBumpTransaction)
		nTmp, err = (*u.FeeBump).DecodeFrom(d, maxDepth)
		n += nTmp
		if err != nil {
			return n, fmt.Errorf("decoding FeeBumpTransaction: %w", err)
		}
		return n, nil
	}
	return n, fmt.Errorf("union TransactionSignaturePayloadTaggedTransaction has invalid Type (EnvelopeType) switch value '%d'", u.Type)
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionSignaturePayloadTaggedTransaction) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionSignaturePayloadTaggedTransaction) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionSignaturePayloadTaggedTransaction)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionSignaturePayloadTaggedTransaction)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s TransactionSignaturePayloadTaggedTransaction) xdrType() {}

var _ xdrType = (*TransactionSignaturePayloadTaggedTransaction)(nil)

type TransactionSignaturePayload struct {
	NetworkId         Hash
	TaggedTransaction TransactionSignaturePayloadTaggedTransaction
}

// EncodeTo encodes this value using the Encoder.
func (s *TransactionSignaturePayload) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.NetworkId.EncodeTo(e); err != nil {
		return err
	}
	if err = s.TaggedTransaction.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*TransactionSignaturePayload)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *TransactionSignaturePayload) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding TransactionSignaturePayload: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.NetworkId.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Hash: %w", err)
	}
	nTmp, err = s.TaggedTransaction.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding TransactionSignaturePayloadTaggedTransaction: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s TransactionSignaturePayload) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *TransactionSignaturePayload) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*TransactionSignaturePayload)(nil)
	_ encoding.BinaryUnmarshaler = (*TransactionSignaturePayload)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s TransactionSignaturePayload) xdrType() {}

var _ xdrType = (*TransactionSignaturePayload)(nil)

type Int128Parts struct {
	Hi Int64
	Lo Uint64
}

// EncodeTo encodes this value using the Encoder.
func (s *Int128Parts) EncodeTo(e *xdr.Encoder) error {
	var err error
	if err = s.Hi.EncodeTo(e); err != nil {
		return err
	}
	if err = s.Lo.EncodeTo(e); err != nil {
		return err
	}
	return nil
}

var _ decoderFrom = (*Int128Parts)(nil)

// DecodeFrom decodes this value using the Decoder.
func (s *Int128Parts) DecodeFrom(d *xdr.Decoder, maxDepth uint) (int, error) {
	if maxDepth == 0 {
		return 0, fmt.Errorf("decoding Int128Parts: %w", ErrMaxDecodingDepthReached)
	}
	maxDepth -= 1
	var err error
	var n, nTmp int
	nTmp, err = s.Hi.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Int64: %w", err)
	}
	nTmp, err = s.Lo.DecodeFrom(d, maxDepth)
	n += nTmp
	if err != nil {
		return n, fmt.Errorf("decoding Uint64: %w", err)
	}
	return n, nil
}

// MarshalBinary implements encoding.BinaryMarshaler.
func (s Int128Parts) MarshalBinary() ([]byte, error) {
	b := bytes.Buffer{}
	e := xdr.NewEncoder(&b)
	err := s.EncodeTo(e)
	return b.Bytes(), err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (s *Int128Parts) UnmarshalBinary(inp []byte) error {
	r := bytes.NewReader(inp)
	o := xdr.DefaultDecodeOptions
	o.MaxInputLen = len(inp)
	d := xdr.NewDecoderWithOptions(r, o)
	_, err := s.DecodeFrom(d, o.MaxDepth)
	return err
}

var (
	_ encoding.BinaryMarshaler   = (*Int128Parts)(nil)
	_ encoding.BinaryUnmarshaler = (*Int128Parts)(nil)
)

// xdrType signals that this type represents XDR values defined by this package.
func (s Int128Parts) xdrType() {}

var _ xdrType = (*Int128Parts)(nil)
