package types

import "github.com/eoscanada/eos-go/ecc"

// NewNewAccount returns a `newaccount` action that lives on the
// `eosio.system` contract.
func NewNewAccount(creator, newAccount string, publicKey ecc.PublicKey) *Action {
	return &Action{
		Account: "eosio",
		Name:    "newaccount",
		Authorization: []PermissionLevel{
			{Actor: AN(creator), Permission: "active"},
		},
		ActionData: NewActionData(NewAccount{
			Creator: AN(creator),
			Name:    AN(newAccount),
			Owner: Authority{
				Threshold: 1,
				Keys: []KeyWeight{
					{
						PublicKey: publicKey,
						Weight:    1,
					},
				},
				Accounts: []PermissionLevelWeight{},
			},
			Active: Authority{
				Threshold: 1,
				Keys: []KeyWeight{
					{
						PublicKey: publicKey,
						Weight:    1,
					},
				},
				Accounts: []PermissionLevelWeight{},
			},
		}),
	}
}

// NewAccount represents a `newaccount` action on the `eosio.system`
// contract. It is one of the rare ones to be hard-coded into the
// blockchain.
type NewAccount struct {
	Creator AccountName `json:"creator"`
	Name    AccountName `json:"name"`
	Owner   Authority   `json:"owner"`
	Active  Authority   `json:"active"`
}
