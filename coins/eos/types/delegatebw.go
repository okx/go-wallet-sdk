package types

// NewDelegateBW returns a `delegatebw` action that lives on the
// `eosio.system` contract.
func NewDelegateBW(from, receiver string, stakeCPU, stakeNet Asset, transfer bool) *Action {
	return &Action{
		Account: "eosio",
		Name:    "delegatebw",
		Authorization: []PermissionLevel{
			{Actor: AN(from), Permission: "active"},
		},
		ActionData: NewActionData(DelegateBW{
			From:     AN(from),
			Receiver: AN(receiver),
			StakeNet: stakeNet,
			StakeCPU: stakeCPU,
			Transfer: transfer,
		}),
	}
}

// DelegateBW represents the `eosio.system::delegatebw` action.
type DelegateBW struct {
	From     AccountName `json:"from"`
	Receiver AccountName `json:"receiver"`
	StakeNet Asset       `json:"stake_net_quantity"`
	StakeCPU Asset       `json:"stake_cpu_quantity"`
	Transfer bool        `json:"transfer"`
}
