/*
*
MIT License

Copyright (c) 2018 WavesPlatform
*/
package types

type Order_PriceMode int32

const (
	Order_DEFAULT        Order_PriceMode = 0
	Order_FIXED_DECIMALS Order_PriceMode = 1
	Order_ASSET_DECIMALS Order_PriceMode = 2
)

type Order_Side int32

const (
	Order_BUY  Order_Side = 0
	Order_SELL Order_Side = 1
)

type Order struct {
	ChainId          int32           `protobuf:"varint,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	MatcherPublicKey []byte          `protobuf:"bytes,3,opt,name=matcher_public_key,json=matcherPublicKey,proto3" json:"matcher_public_key,omitempty"`
	AssetPair        *AssetPair      `protobuf:"bytes,4,opt,name=asset_pair,json=assetPair,proto3" json:"asset_pair,omitempty"`
	OrderSide        Order_Side      `protobuf:"varint,5,opt,name=order_side,json=orderSide,proto3,enum=waves.Order_Side" json:"order_side,omitempty"`
	Amount           int64           `protobuf:"varint,6,opt,name=amount,proto3" json:"amount,omitempty"`
	Price            int64           `protobuf:"varint,7,opt,name=price,proto3" json:"price,omitempty"`
	Timestamp        int64           `protobuf:"varint,8,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Expiration       int64           `protobuf:"varint,9,opt,name=expiration,proto3" json:"expiration,omitempty"`
	MatcherFee       *Amount         `protobuf:"bytes,10,opt,name=matcher_fee,json=matcherFee,proto3" json:"matcher_fee,omitempty"`
	Version          int32           `protobuf:"varint,11,opt,name=version,proto3" json:"version,omitempty"`
	Proofs           [][]byte        `protobuf:"bytes,12,rep,name=proofs,proto3" json:"proofs,omitempty"`
	PriceMode        Order_PriceMode `protobuf:"varint,14,opt,name=price_mode,json=priceMode,proto3,enum=waves.Order_PriceMode" json:"price_mode,omitempty"`
	// Types that are assignable to Sender:
	//	*Order_SenderPublicKey
	//	*Order_Eip712Signature
	Sender isOrder_Sender `protobuf_oneof:"sender"`
}

type AssetPair struct {
	AmountAssetId []byte `protobuf:"bytes,1,opt,name=amount_asset_id,json=amountAssetId,proto3" json:"amount_asset_id,omitempty"`
	PriceAssetId  []byte `protobuf:"bytes,2,opt,name=price_asset_id,json=priceAssetId,proto3" json:"price_asset_id,omitempty"`
}

type isOrder_Sender interface {
	isOrder_Sender()
}

type Order_SenderPublicKey struct {
	SenderPublicKey []byte `protobuf:"bytes,2,opt,name=sender_public_key,json=senderPublicKey,proto3,oneof"`
}

type Order_Eip712Signature struct {
	Eip712Signature []byte `protobuf:"bytes,13,opt,name=eip712_signature,json=eip712Signature,proto3,oneof"`
}

func (*Order_SenderPublicKey) isOrder_Sender() {}

func (*Order_Eip712Signature) isOrder_Sender() {}
