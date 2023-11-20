package types

func NewBuyRAM(payer, receiver string, quantity Asset) *Action {
	a := &Action{
		Account: "eosio",
		Name:    "buyram",
		Authorization: []PermissionLevel{
			{Actor: AN(payer), Permission: "active"},
		},
		ActionData: NewActionData(BuyRAM{
			Payer:    AN(payer),
			Receiver: AN(receiver),
			Quantity: quantity,
		}),
	}
	return a
}

// BuyRAM represents the `eosio.system::buyram` action.
type BuyRAM struct {
	Payer    AccountName `json:"payer"`
	Receiver AccountName `json:"receiver"`
	Quantity Asset       `json:"quant"` // specified in EOS
}

// NewBuyRAMBytes will buy at current market price a given number of
// bytes of RAM, and grant them to the `receiver` account.
func NewBuyRAMBytes(payer, receiver string, bytes uint32) *Action {
	a := &Action{
		Account: AN("eosio"),
		Name:    ActN("buyrambytes"),
		Authorization: []PermissionLevel{
			{Actor: AN(payer), Permission: PermissionName("active")},
		},
		ActionData: NewActionData(BuyRAMBytes{
			Payer:    AN(payer),
			Receiver: AN(receiver),
			Bytes:    bytes,
		}),
	}
	return a
}

// BuyRAMBytes represents the `eosio.system::buyrambytes` action.
type BuyRAMBytes struct {
	Payer    AccountName `json:"payer"`
	Receiver AccountName `json:"receiver"`
	Bytes    uint32      `json:"bytes"`
}
