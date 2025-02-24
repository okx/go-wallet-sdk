package v2

import (
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
)

// ModuleId the identifier for a module e.g. 0x1::coin
type ModuleId struct {
	Address AccountAddress
	Name    string
}

func (mod *ModuleId) MarshalBCS(bcs *bcs.Serializer) {
	mod.Address.MarshalBCS(bcs)
	bcs.WriteString(mod.Name)
}
func (mod *ModuleId) UnmarshalBCS(bcs *bcs.Deserializer) {
	mod.Address.UnmarshalBCS(bcs)
	mod.Name = bcs.ReadString()
}
