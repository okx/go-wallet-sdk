package externalapi

import "encoding/hex"

// DomainSubnetworkIDSize is the size of the array used to store subnetwork IDs.
const DomainSubnetworkIDSize = 20

// DomainSubnetworkID is the domain representation of a Subnetwork ID
type DomainSubnetworkID [DomainSubnetworkIDSize]byte

// String stringifies a subnetwork ID.
func (id DomainSubnetworkID) String() string {
	return hex.EncodeToString(id[:])
}

// Clone returns a clone of DomainSubnetworkID
func (id *DomainSubnetworkID) Clone() *DomainSubnetworkID {
	idClone := *id
	return &idClone
}

// If this doesn't compile, it means the type definition has been changed, so it's
// an indication to update Equal and Clone accordingly.
var _ DomainSubnetworkID = [DomainSubnetworkIDSize]byte{}

// Equal returns whether id equals to other
func (id *DomainSubnetworkID) Equal(other *DomainSubnetworkID) bool {
	if id == nil || other == nil {
		return id == other
	}

	return *id == *other
}
