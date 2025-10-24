package transactionhelper

import (
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/model/externalapi"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/utils/subnetworks"
)

// NewSubnetworkTransaction returns a new trsnactions in the specified subnetwork with specified gas and payload
func NewSubnetworkTransaction(version uint16, inputs []*externalapi.DomainTransactionInput,
	outputs []*externalapi.DomainTransactionOutput, subnetworkID *externalapi.DomainSubnetworkID,
	gas uint64, payload []byte) *externalapi.DomainTransaction {

	return &externalapi.DomainTransaction{
		Version:      version,
		Inputs:       inputs,
		Outputs:      outputs,
		LockTime:     0,
		SubnetworkID: *subnetworkID,
		Gas:          gas,
		Payload:      payload,
		Fee:          0,
		Mass:         0,
	}
}

// NewNativeTransaction returns a new native transaction
func NewNativeTransaction(version uint16, inputs []*externalapi.DomainTransactionInput,
	outputs []*externalapi.DomainTransactionOutput) *externalapi.DomainTransaction {
	return &externalapi.DomainTransaction{
		Version:      version,
		Inputs:       inputs,
		Outputs:      outputs,
		LockTime:     0,
		SubnetworkID: subnetworks.SubnetworkIDNative,
		Gas:          0,
		Payload:      []byte{},
		Fee:          0,
		Mass:         0,
	}
}
