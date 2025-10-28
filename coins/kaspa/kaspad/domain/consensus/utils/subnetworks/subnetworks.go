package subnetworks

import "github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/model/externalapi"

var (
	// SubnetworkIDNative is the default subnetwork ID which is used for transactions without related payload data
	SubnetworkIDNative = externalapi.DomainSubnetworkID{}

	// SubnetworkIDCoinbase is the subnetwork ID which is used for the coinbase transaction
	SubnetworkIDCoinbase = externalapi.DomainSubnetworkID{1}

	// SubnetworkIDRegistry is the subnetwork ID which is used for adding new sub networks to the registry
	SubnetworkIDRegistry = externalapi.DomainSubnetworkID{2}
)

// IsBuiltIn returns true if the subnetwork is a built in subnetwork, which
// means all nodes, including partial nodes, must validate it, and its transactions
// always use 0 gas.
func IsBuiltIn(id externalapi.DomainSubnetworkID) bool {
	return id == SubnetworkIDCoinbase || id == SubnetworkIDRegistry
}

// IsBuiltInOrNative returns true if the subnetwork is the native or a built in subnetwork,
// see IsBuiltIn for further details
func IsBuiltInOrNative(id externalapi.DomainSubnetworkID) bool {
	return id == SubnetworkIDNative || IsBuiltIn(id)
}
