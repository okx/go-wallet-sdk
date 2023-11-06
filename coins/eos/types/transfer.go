package types

// NewTransfer creates a new Transfer action object.
func NewTransfer(from, to string, quantity Asset, memo string) *Action {
	return &Action{
		Account: AN("eosio.token"),
		Name:    ActN("transfer"),
		Authorization: []PermissionLevel{
			{Actor: AN(from), Permission: PN("active")},
		},
		ActionData: NewActionData(Transfer{
			From:     AN(from),
			To:       AN(to),
			Quantity: quantity,
			Memo:     memo,
		}),
	}
}

// Transfer represents the `transfer` struct on `eosio.token` contract.
type Transfer struct {
	From     AccountName `json:"from"`
	To       AccountName `json:"to"`
	Quantity Asset       `json:"quantity"`
	Memo     string      `json:"memo"`
}

// NewContractTransfer creates a new ContractTransfer action object.
func NewContractTransfer(name, from, to string, quantity Asset, memo string) *Action {
	return &Action{
		Account: AN(name),
		Name:    ActN("transfer"),
		Authorization: []PermissionLevel{
			{Actor: AN(from), Permission: PN("active")},
		},
		ActionData: NewActionData(ContractTransfer{
			From:     AN(from),
			To:       AN(to),
			Quantity: quantity,
			Memo:     memo,
		}),
	}
}

// ContractTransfer represents the `transfer` struct on `eosio.token` contract.
type ContractTransfer struct {
	From     AccountName `json:"from"`
	To       AccountName `json:"to"`
	Quantity Asset       `json:"quantity"`
	Memo     string      `json:"memo"`
}
