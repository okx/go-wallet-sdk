package types

func NewTokenTransferAction(contract, from, to, amount, memo, token string, precision uint8) *Action {
	//1. new token Asset
	var symbol = Symbol{Precision: precision, Symbol: token}
	asset, err := NewFixedSymbolAssetFromString(symbol, amount)
	if err != nil {
		return nil
	}

	//2. new token Action
	return NewTokenTransfer(contract, from, to, asset, memo)
}

func NewTokenTransfer(contract, from, to string, quantity Asset, memo string) *Action {
	return &Action{
		Account: AccountName(contract),
		Name:    ActionName("transfer"),
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
