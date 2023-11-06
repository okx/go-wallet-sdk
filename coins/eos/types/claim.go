package types

// NewClaim creates a new Claim action object.
// ex. https://wax.bloks.io/transaction/ca77897065604decd0ff16f62624bcea6706abdd688fd1c7743c0630e8db8dda?tab=raw
// account example: "farmersworld"
func NewClaim(account, owner string, assertId string) *Action {
	a := &Action{
		Account: AN("farmersworld"),
		Name:    ActN("claim"),
		Authorization: []PermissionLevel{
			{Actor: AN(owner), Permission: PermissionName("active")},
		},
		ActionData: NewActionData(Claim{
			Owner:   AN(owner),
			AssetID: assertId,
		}),
	}
	return a
}

// Claim represents the `farmersworld` action.
type Claim struct {
	Owner   AccountName `json:"owner"`
	AssetID string      `json:"asset_id"`
}
