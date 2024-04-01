package aptos_types

import (
	"github.com/okx/go-wallet-sdk/coins/aptos/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/serde"
)

type ChainId uint8

func (o *ChainId) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	if err := serializer.SerializeU8((uint8)(*o)); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *ChainId) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}
