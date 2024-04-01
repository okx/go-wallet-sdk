package aptos_types

import (
	"github.com/okx/go-wallet-sdk/coins/aptos/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/serde"
)

type AccountAuthenticator interface {
	Serialize(serializer serde.Serializer) error
	BcsSerialize() ([]byte, error)
}
type AccountAuthenticatorEd25519 struct {
	PublicKey Ed25519PublicKey
	Signature Ed25519Signature
}

func (o *AccountAuthenticatorEd25519) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	err := serializer.SerializeVariantIndex(0)
	if err != nil {
		return err
	}
	if err := o.PublicKey.Serialize(serializer); err != nil {
		return err
	}
	if err := o.Signature.Serialize(serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (o *AccountAuthenticatorEd25519) BcsSerialize() ([]byte, error) {
	if o == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := o.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}

type TransactionAuthenticator interface {
	Serialize(serializer serde.Serializer) error
	BcsSerialize() ([]byte, error)
}

type TransactionAuthenticatorEd25519 struct {
	PublicKey Ed25519PublicKey
	Signature Ed25519Signature
}

func (obj *TransactionAuthenticatorEd25519) Serialize(serializer serde.Serializer) error {
	if err := serializer.IncreaseContainerDepth(); err != nil {
		return err
	}
	err := serializer.SerializeVariantIndex(0)
	if err != nil {
		return err
	}
	if err := obj.PublicKey.Serialize(serializer); err != nil {
		return err
	}
	if err := obj.Signature.Serialize(serializer); err != nil {
		return err
	}
	serializer.DecreaseContainerDepth()
	return nil
}

func (obj *TransactionAuthenticatorEd25519) BcsSerialize() ([]byte, error) {
	if obj == nil {
		return nil, ErrNullObject
	}
	serializer := bcs.NewSerializer()
	if err := obj.Serialize(serializer); err != nil {
		return nil, err
	}
	return serializer.GetBytes(), nil
}
