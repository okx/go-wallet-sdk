package generated

import "github.com/okx/go-wallet-sdk/coins/waves/types"

type SignedTransaction struct {
	// Types that are assignable to Transaction:
	//	*SignedTransaction_WavesTransaction
	//	*SignedTransaction_EthereumTransaction
	Transaction isSignedTransaction_Transaction `protobuf_oneof:"transaction"`
	Proofs      [][]byte                        `protobuf:"bytes,2,rep,name=proofs,proto3" json:"proofs,omitempty"`
}

type isSignedTransaction_Transaction interface {
	isSignedTransaction_Transaction()
}

type SignedTransaction_WavesTransaction struct {
	WavesTransaction *Transaction `protobuf:"bytes,1,opt,name=waves_transaction,json=wavesTransaction,proto3,oneof"`
}

type SignedTransaction_EthereumTransaction struct {
	EthereumTransaction []byte `protobuf:"bytes,3,opt,name=ethereum_transaction,json=ethereumTransaction,proto3,oneof"`
}

func (*SignedTransaction_WavesTransaction) isSignedTransaction_Transaction() {}

func (*SignedTransaction_EthereumTransaction) isSignedTransaction_Transaction() {}

type Transaction struct {
	ChainId         int32         `protobuf:"varint,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	SenderPublicKey []byte        `protobuf:"bytes,2,opt,name=sender_public_key,json=senderPublicKey,proto3" json:"sender_public_key,omitempty"`
	Fee             *types.Amount `protobuf:"bytes,3,opt,name=fee,proto3" json:"fee,omitempty"`
	Timestamp       int64         `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Version         int32         `protobuf:"varint,5,opt,name=version,proto3" json:"version,omitempty"`
	// Types that are assignable to Data:
	//	*Transaction_Genesis
	//	*Transaction_Payment
	//	*Transaction_Issue
	//	*Transaction_Transfer
	//	*Transaction_Reissue
	//	*Transaction_Burn
	//	*Transaction_Exchange
	//	*Transaction_Lease
	//	*Transaction_LeaseCancel
	//	*Transaction_CreateAlias
	//	*Transaction_MassTransfer
	//	*Transaction_DataTransaction
	//	*Transaction_SetScript
	//	*Transaction_SponsorFee
	//	*Transaction_SetAssetScript
	//	*Transaction_InvokeScript
	//	*Transaction_UpdateAssetInfo
	//	*Transaction_InvokeExpression
	Data isTransaction_Data `protobuf_oneof:"data"`
}

type isTransaction_Data interface {
	isTransaction_Data()
}

type Transaction_Genesis struct {
	Genesis *types.GenesisTransactionData `protobuf:"bytes,101,opt,name=genesis,proto3,oneof"`
}

type Transaction_Payment struct {
	Payment *types.PaymentTransactionData `protobuf:"bytes,102,opt,name=payment,proto3,oneof"`
}

type Transaction_Issue struct {
	Issue *types.IssueTransactionData `protobuf:"bytes,103,opt,name=issue,proto3,oneof"`
}

type Transaction_Transfer struct {
	Transfer *types.TransferTransactionData `protobuf:"bytes,104,opt,name=transfer,proto3,oneof"`
}

type Transaction_Reissue struct {
	Reissue *types.ReissueTransactionData `protobuf:"bytes,105,opt,name=reissue,proto3,oneof"`
}

type Transaction_Burn struct {
	Burn *types.BurnTransactionData `protobuf:"bytes,106,opt,name=burn,proto3,oneof"`
}

type Transaction_Exchange struct {
	Exchange *types.ExchangeTransactionData `protobuf:"bytes,107,opt,name=exchange,proto3,oneof"`
}

type Transaction_Lease struct {
	Lease *types.LeaseTransactionData `protobuf:"bytes,108,opt,name=lease,proto3,oneof"`
}

type Transaction_LeaseCancel struct {
	LeaseCancel *types.LeaseCancelTransactionData `protobuf:"bytes,109,opt,name=lease_cancel,json=leaseCancel,proto3,oneof"`
}

type Transaction_CreateAlias struct {
	CreateAlias *types.CreateAliasTransactionData `protobuf:"bytes,110,opt,name=create_alias,json=createAlias,proto3,oneof"`
}

type Transaction_MassTransfer struct {
	MassTransfer *types.MassTransferTransactionData `protobuf:"bytes,111,opt,name=mass_transfer,json=massTransfer,proto3,oneof"`
}

type Transaction_DataTransaction struct {
	DataTransaction *types.DataTransactionData `protobuf:"bytes,112,opt,name=data_transaction,json=dataTransaction,proto3,oneof"`
}

type Transaction_SetScript struct {
	SetScript *types.SetScriptTransactionData `protobuf:"bytes,113,opt,name=set_script,json=setScript,proto3,oneof"`
}

type Transaction_SponsorFee struct {
	SponsorFee *types.SponsorFeeTransactionData `protobuf:"bytes,114,opt,name=sponsor_fee,json=sponsorFee,proto3,oneof"`
}

type Transaction_SetAssetScript struct {
	SetAssetScript *types.SetAssetScriptTransactionData `protobuf:"bytes,115,opt,name=set_asset_script,json=setAssetScript,proto3,oneof"`
}

type Transaction_InvokeScript struct {
	InvokeScript *types.InvokeScriptTransactionData `protobuf:"bytes,116,opt,name=invoke_script,json=invokeScript,proto3,oneof"`
}

type Transaction_UpdateAssetInfo struct {
	UpdateAssetInfo *types.UpdateAssetInfoTransactionData `protobuf:"bytes,117,opt,name=update_asset_info,json=updateAssetInfo,proto3,oneof"`
}

type Transaction_InvokeExpression struct {
	InvokeExpression *types.InvokeExpressionTransactionData `protobuf:"bytes,119,opt,name=invoke_expression,json=invokeExpression,proto3,oneof"`
}

func (*Transaction_Genesis) isTransaction_Data() {}

func (*Transaction_Payment) isTransaction_Data() {}

func (*Transaction_Issue) isTransaction_Data() {}

func (*Transaction_Transfer) isTransaction_Data() {}

func (*Transaction_Reissue) isTransaction_Data() {}

func (*Transaction_Burn) isTransaction_Data() {}

func (*Transaction_Exchange) isTransaction_Data() {}

func (*Transaction_Lease) isTransaction_Data() {}

func (*Transaction_LeaseCancel) isTransaction_Data() {}

func (*Transaction_CreateAlias) isTransaction_Data() {}

func (*Transaction_MassTransfer) isTransaction_Data() {}

func (*Transaction_DataTransaction) isTransaction_Data() {}

func (*Transaction_SetScript) isTransaction_Data() {}

func (*Transaction_SponsorFee) isTransaction_Data() {}

func (*Transaction_SetAssetScript) isTransaction_Data() {}

func (*Transaction_InvokeScript) isTransaction_Data() {}

func (*Transaction_UpdateAssetInfo) isTransaction_Data() {}

func (*Transaction_InvokeExpression) isTransaction_Data() {}
