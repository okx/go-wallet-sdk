package subnetworks

import (
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/model/externalapi"
	"github.com/pkg/errors"
)

// FromBytes creates a DomainSubnetworkID from the given byte slice
func FromBytes(subnetworkIDBytes []byte) (*externalapi.DomainSubnetworkID, error) {
	if len(subnetworkIDBytes) != externalapi.DomainSubnetworkIDSize {
		return nil, errors.Errorf("invalid hash size. Want: %d, got: %d",
			externalapi.DomainSubnetworkIDSize, len(subnetworkIDBytes))
	}
	var domainSubnetworkID externalapi.DomainSubnetworkID
	copy(domainSubnetworkID[:], subnetworkIDBytes)
	return &domainSubnetworkID, nil
}
