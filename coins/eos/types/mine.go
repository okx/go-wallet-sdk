package types

// NewMine creates a new Mine action object. account example: "m.federation"
func NewMine(account, miner string, nonce string) *Action {
	return &Action{
		Account: AN(account),
		Name:    ActN("mine"),
		Authorization: []PermissionLevel{
			{Actor: AN(miner), Permission: PN("active")},
		},
		ActionData: NewActionData(Mine{
			Miner: AN(miner),
			Nonce: nonce,
		}),
	}
}

type Mine struct {
	Miner AccountName `json:"miner"`
	Nonce string      `json:"nonce"`
}
