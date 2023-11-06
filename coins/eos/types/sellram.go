package types

// NewSellRAM will sell at current market price a given number of
// bytes of RAM.
func NewSellRAM(account string, bytes uint64) *Action {
	a := &Action{
		Account: AN("eosio"),
		Name:    ActN("sellram"),
		Authorization: []PermissionLevel{
			{Actor: AN(account), Permission: PermissionName("active")},
		},
		ActionData: NewActionData(SellRAM{
			Account: AN(account),
			Bytes:   bytes,
		}),
	}
	return a
}

// SellRAM represents the `eosio.system::sellram` action.
type SellRAM struct {
	Account AccountName `json:"account"`
	Bytes   uint64      `json:"bytes"`
}
