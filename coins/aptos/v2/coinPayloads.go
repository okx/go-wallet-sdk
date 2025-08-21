package v2

import "github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"

// CoinTransferPayload builds an EntryFunction payload for transferring coins
//
// Args:
//   - coinType is the type of coin to transfer. If none is provided, it will transfer 0x1::aptos_coin:AptosCoin
//   - dest is the destination [AccountAddress]
//   - amount is the amount of coins to transfer
func CoinTransferPayload(coinType *TypeTag, dest AccountAddress, amount uint64) (payload *EntryFunction, err error) {
	amountBytes, err := bcs.SerializeU64(amount)
	if err != nil {
		return nil, err
	}

	if coinType == nil || *coinType == AptosCoinTypeTag {
		return &EntryFunction{
			Module: ModuleId{
				Address: AccountOne,
				Name:    "aptos_account",
			},
			Function: "transfer",
			ArgTypes: []TypeTag{},
			Args: [][]byte{
				dest[:],
				amountBytes,
			},
		}, nil
	} else {
		return &EntryFunction{
			Module: ModuleId{
				Address: AccountOne,
				Name:    "aptos_account",
			},
			Function: "transfer_coins",
			ArgTypes: []TypeTag{*coinType},
			Args: [][]byte{
				dest[:],
				amountBytes,
			},
		}, nil
	}
}

// CoinBatchTransferPayload builds an EntryFunction payload for transferring coins to multiple receivers
//
// Args:
//   - coinType is the type of coin to transfer. If none is provided, it will transfer 0x1::aptos_coin:AptosCoin
//   - dests are the destination [AccountAddress]s
//   - amounts are the amount of coins to transfer per destination
func CoinBatchTransferPayload(coinType *TypeTag, dests []AccountAddress, amounts []uint64) (payload *EntryFunction, err error) {
	destBytes, err := bcs.SerializeSequenceOnly(dests)
	if err != nil {
		return nil, err
	}
	amountsBytes, err := bcs.SerializeSingle(func(ser *bcs.Serializer) {
		bcs.SerializeSequenceWithFunction(amounts, ser, func(ser *bcs.Serializer, amount uint64) {
			ser.U64(amount)
		})
	})
	if err != nil {
		return nil, err
	}

	if coinType == nil || *coinType == AptosCoinTypeTag {
		return &EntryFunction{
			Module: ModuleId{
				Address: AccountOne,
				Name:    "aptos_account",
			},
			Function: "batch_transfer",
			ArgTypes: []TypeTag{},
			Args: [][]byte{
				destBytes,
				amountsBytes,
			},
		}, nil
	} else {
		return &EntryFunction{
			Module: ModuleId{
				Address: AccountOne,
				Name:    "aptos_account",
			},
			Function: "batch_transfer_coins",
			ArgTypes: []TypeTag{*coinType},
			Args: [][]byte{
				destBytes,
				amountsBytes,
			},
		}, nil
	}
}
