/*
*
MIT License

Copyright (c) 2018 WavesPlatform
*/
package types

type GenesisTransactionData struct {
	RecipientAddress []byte `protobuf:"bytes,1,opt,name=recipient_address,json=recipientAddress,proto3" json:"recipient_address,omitempty"`
	Amount           int64  `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
}

type PaymentTransactionData struct {
	RecipientAddress []byte `protobuf:"bytes,1,opt,name=recipient_address,json=recipientAddress,proto3" json:"recipient_address,omitempty"`
	Amount           int64  `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
}

type TransferTransactionData struct {
	Recipient  *Recipient `protobuf:"bytes,1,opt,name=recipient,proto3" json:"recipient,omitempty"`
	Amount     *Amount    `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount,omitempty"`
	Attachment []byte     `protobuf:"bytes,3,opt,name=attachment,proto3" json:"attachment,omitempty"`
}

type CreateAliasTransactionData struct {
	Alias string `protobuf:"bytes,1,opt,name=alias,proto3" json:"alias,omitempty"`
}

type DataTransactionData struct {
	Data []*DataTransactionData_DataEntry `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
}

type MassTransferTransactionData struct {
	AssetId    []byte                                  `protobuf:"bytes,1,opt,name=asset_id,json=assetId,proto3" json:"asset_id,omitempty"`
	Transfers  []*MassTransferTransactionData_Transfer `protobuf:"bytes,2,rep,name=transfers,proto3" json:"transfers,omitempty"`
	Attachment []byte                                  `protobuf:"bytes,3,opt,name=attachment,proto3" json:"attachment,omitempty"`
}

type LeaseTransactionData struct {
	Recipient *Recipient `protobuf:"bytes,1,opt,name=recipient,proto3" json:"recipient,omitempty"`
	Amount    int64      `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
}

type LeaseCancelTransactionData struct {
	LeaseId []byte `protobuf:"bytes,1,opt,name=lease_id,json=leaseId,proto3" json:"lease_id,omitempty"`
}

type BurnTransactionData struct {
	AssetAmount *Amount `protobuf:"bytes,1,opt,name=asset_amount,json=assetAmount,proto3" json:"asset_amount,omitempty"`
}

type IssueTransactionData struct {
	Name        string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Amount      int64  `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
	Decimals    int32  `protobuf:"varint,4,opt,name=decimals,proto3" json:"decimals,omitempty"`
	Reissuable  bool   `protobuf:"varint,5,opt,name=reissuable,proto3" json:"reissuable,omitempty"`
	Script      []byte `protobuf:"bytes,6,opt,name=script,proto3" json:"script,omitempty"`
}

type ReissueTransactionData struct {
	AssetAmount *Amount `protobuf:"bytes,1,opt,name=asset_amount,json=assetAmount,proto3" json:"asset_amount,omitempty"`
	Reissuable  bool    `protobuf:"varint,2,opt,name=reissuable,proto3" json:"reissuable,omitempty"`
}

type SetAssetScriptTransactionData struct {
	AssetId []byte `protobuf:"bytes,1,opt,name=asset_id,json=assetId,proto3" json:"asset_id,omitempty"`
	Script  []byte `protobuf:"bytes,2,opt,name=script,proto3" json:"script,omitempty"`
}

type SetScriptTransactionData struct {
	Script []byte `protobuf:"bytes,1,opt,name=script,proto3" json:"script,omitempty"`
}

type ExchangeTransactionData struct {
	Amount         int64    `protobuf:"varint,1,opt,name=amount,proto3" json:"amount,omitempty"`
	Price          int64    `protobuf:"varint,2,opt,name=price,proto3" json:"price,omitempty"`
	BuyMatcherFee  int64    `protobuf:"varint,3,opt,name=buy_matcher_fee,json=buyMatcherFee,proto3" json:"buy_matcher_fee,omitempty"`
	SellMatcherFee int64    `protobuf:"varint,4,opt,name=sell_matcher_fee,json=sellMatcherFee,proto3" json:"sell_matcher_fee,omitempty"`
	Orders         []*Order `protobuf:"bytes,5,rep,name=orders,proto3" json:"orders,omitempty"`
}

type SponsorFeeTransactionData struct {
	MinFee *Amount `protobuf:"bytes,1,opt,name=min_fee,json=minFee,proto3" json:"min_fee,omitempty"`
}

type InvokeScriptTransactionData struct {
	DApp         *Recipient `protobuf:"bytes,1,opt,name=d_app,json=dApp,proto3" json:"d_app,omitempty"`
	FunctionCall []byte     `protobuf:"bytes,2,opt,name=function_call,json=functionCall,proto3" json:"function_call,omitempty"`
	Payments     []*Amount  `protobuf:"bytes,3,rep,name=payments,proto3" json:"payments,omitempty"`
}

type UpdateAssetInfoTransactionData struct {
	AssetId     []byte `protobuf:"bytes,1,opt,name=asset_id,json=assetId,proto3" json:"asset_id,omitempty"`
	Name        string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
}

type InvokeExpressionTransactionData struct {
	Expression []byte `protobuf:"bytes,1,opt,name=expression,proto3" json:"expression,omitempty"`
}

type DataTransactionData_DataEntry struct {
	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// Types that are assignable to Value:
	//	*DataTransactionData_DataEntry_IntValue
	//	*DataTransactionData_DataEntry_BoolValue
	//	*DataTransactionData_DataEntry_BinaryValue
	//	*DataTransactionData_DataEntry_StringValue
	Value isDataTransactionData_DataEntry_Value `protobuf_oneof:"value"`
}

type isDataTransactionData_DataEntry_Value interface {
	isDataTransactionData_DataEntry_Value()
}

type DataTransactionData_DataEntry_IntValue struct {
	IntValue int64 `protobuf:"varint,10,opt,name=int_value,json=intValue,proto3,oneof"`
}

type DataTransactionData_DataEntry_BoolValue struct {
	BoolValue bool `protobuf:"varint,11,opt,name=bool_value,json=boolValue,proto3,oneof"`
}

type DataTransactionData_DataEntry_BinaryValue struct {
	BinaryValue []byte `protobuf:"bytes,12,opt,name=binary_value,json=binaryValue,proto3,oneof"`
}

type DataTransactionData_DataEntry_StringValue struct {
	StringValue string `protobuf:"bytes,13,opt,name=string_value,json=stringValue,proto3,oneof"`
}

func (*DataTransactionData_DataEntry_IntValue) isDataTransactionData_DataEntry_Value() {}

func (*DataTransactionData_DataEntry_BoolValue) isDataTransactionData_DataEntry_Value() {}

func (*DataTransactionData_DataEntry_BinaryValue) isDataTransactionData_DataEntry_Value() {}

func (*DataTransactionData_DataEntry_StringValue) isDataTransactionData_DataEntry_Value() {}

type MassTransferTransactionData_Transfer struct {
	Recipient *Recipient `protobuf:"bytes,1,opt,name=recipient,proto3" json:"recipient,omitempty"`
	Amount    int64      `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
}
