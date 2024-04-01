package transaction_builder

import (
	"errors"
	"github.com/okx/go-wallet-sdk/coins/aptos/aptos_types"
	"github.com/okx/go-wallet-sdk/coins/aptos/bcs"
)

func ToBCSArgs(abiArgs []aptos_types.ArgumentABI, args []any) ([][]byte, error) {
	if len(abiArgs) != len(args) {
		return nil, errors.New("wrong number of args provided")
	}
	res := make([][]byte, 0)
	for i, arg := range args {
		serializer := bcs.NewSerializer()
		err := serializeArg(arg, abiArgs[i].TypeTag, serializer)
		if err != nil {
			return nil, err
		}
		res = append(res, serializer.GetBytes())
	}
	return res, nil
}
