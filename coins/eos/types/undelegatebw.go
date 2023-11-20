package types

// NewUndelegateBW returns a `undelegatebw` action that lives on the
// `eosio.system` contract.
func NewUndelegateBW(from, receiver string, unstakeCPU, unstakeNet Asset) *Action {
	return &Action{
		Account: AN("eosio"),
		Name:    ActN("undelegatebw"),
		Authorization: []PermissionLevel{
			{Actor: AN(from), Permission: PN("active")},
		},
		ActionData: NewActionData(UndelegateBW{
			From:       AN(from),
			Receiver:   AN(receiver),
			UnstakeNet: unstakeNet,
			UnstakeCPU: unstakeCPU,
		}),
	}
}

// UndelegateBW represents the `eosio.system::undelegatebw` action.
type UndelegateBW struct {
	From       AccountName `json:"from"`
	Receiver   AccountName `json:"receiver"`
	UnstakeNet Asset       `json:"unstake_net_quantity"`
	UnstakeCPU Asset       `json:"unstake_cpu_quantity"`
}
