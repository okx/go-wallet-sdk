package utxo

import (
	"bytes"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/model/externalapi"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/utils/serialization"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/utils/transactionid"
	"io"

	"github.com/pkg/errors"
)

// SerializeUTXO returns the byte-slice representation for given UTXOEntry-outpoint pair
func SerializeUTXO(entry externalapi.UTXOEntry, outpoint *externalapi.DomainOutpoint) ([]byte, error) {
	w := &bytes.Buffer{}

	err := serializeOutpoint(w, outpoint)
	if err != nil {
		return nil, err
	}

	err = serializeUTXOEntry(w, entry)
	if err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

// DeserializeUTXO deserializes the given byte slice to UTXOEntry-outpoint pair
func DeserializeUTXO(utxoBytes []byte) (entry externalapi.UTXOEntry, outpoint *externalapi.DomainOutpoint, err error) {
	r := bytes.NewReader(utxoBytes)
	outpoint, err = deserializeOutpoint(r)
	if err != nil {
		return nil, nil, err
	}

	entry, err = deserializeUTXOEntry(r)
	if err != nil {
		return nil, nil, err
	}

	return entry, outpoint, nil
}

func serializeOutpoint(w io.Writer, outpoint *externalapi.DomainOutpoint) error {
	_, err := w.Write(outpoint.TransactionID.ByteSlice())
	if err != nil {
		return err
	}

	err = serialization.WriteElement(w, outpoint.Index)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func deserializeOutpoint(r io.Reader) (*externalapi.DomainOutpoint, error) {
	transactionIDBytes := make([]byte, externalapi.DomainHashSize)
	_, err := io.ReadFull(r, transactionIDBytes)
	if err != nil {
		return nil, err
	}

	transactionID, err := transactionid.FromBytes(transactionIDBytes)
	if err != nil {
		return nil, err
	}

	var index uint32
	err = serialization.ReadElement(r, &index)
	if err != nil {
		return nil, err
	}

	return &externalapi.DomainOutpoint{
		TransactionID: *transactionID,
		Index:         index,
	}, nil
}

func serializeUTXOEntry(w io.Writer, entry externalapi.UTXOEntry) error {
	err := serialization.WriteElements(w, entry.BlockDAAScore(), entry.Amount(), entry.IsCoinbase())
	if err != nil {
		return err
	}
	err = serialization.WriteElement(w, entry.ScriptPublicKey().Version)
	if err != nil {
		return err
	}
	count := uint64(len(entry.ScriptPublicKey().Script))
	err = serialization.WriteElement(w, count)
	if err != nil {
		return err
	}

	_, err = w.Write(entry.ScriptPublicKey().Script)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func deserializeUTXOEntry(r io.Reader) (externalapi.UTXOEntry, error) {
	var blockDAAScore uint64
	var amount uint64
	var isCoinbase bool
	err := serialization.ReadElements(r, &blockDAAScore, &amount, &isCoinbase)
	if err != nil {
		return nil, err
	}

	var version uint16
	err = serialization.ReadElement(r, &version)
	if err != nil {
		return nil, err
	}

	var scriptPubKeyLen uint64
	err = serialization.ReadElement(r, &scriptPubKeyLen)
	if err != nil {
		return nil, err
	}

	scriptPubKeyScript := make([]byte, scriptPubKeyLen)
	_, err = io.ReadFull(r, scriptPubKeyScript)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	scriptPubKey := externalapi.ScriptPublicKey{scriptPubKeyScript, version}

	return NewUTXOEntry(amount, &scriptPubKey, isCoinbase, blockDAAScore), nil
}
