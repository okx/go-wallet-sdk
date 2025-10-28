package serialization

import (
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/model/externalapi"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/util/binaryserializer"
	"github.com/pkg/errors"
	"io"
)

// errNoEncodingForType signifies that there's no encoding for the given type.
var errNoEncodingForType = errors.New("there's no encoding for this type")

var errMalformed = errors.New("errMalformed")

// WriteElement writes the little endian representation of element to w.
func WriteElement(w io.Writer, element interface{}) error {
	// Attempt to write the element based on the concrete type via fast
	// type assertions first.
	switch e := element.(type) {
	case []byte:
		err := WriteElement(w, uint64(len(e)))
		if err != nil {
			return err
		}
		_, err = w.Write(e)
		if err != nil {
			return err
		}
		return nil
	case int16:
		err := binaryserializer.PutUint16(w, uint16(e))
		if err != nil {
			return err
		}
		return nil
	case uint16:
		err := binaryserializer.PutUint16(w, e)
		if err != nil {
			return err
		}
		return nil
	case int32:
		err := binaryserializer.PutUint32(w, uint32(e))
		if err != nil {
			return err
		}
		return nil

	case uint32:
		err := binaryserializer.PutUint32(w, e)
		if err != nil {
			return err
		}
		return nil

	case int64:
		err := binaryserializer.PutUint64(w, uint64(e))
		if err != nil {
			return err
		}
		return nil

	case uint64:
		err := binaryserializer.PutUint64(w, e)
		if err != nil {
			return err
		}
		return nil

	case uint8:
		err := binaryserializer.PutUint8(w, e)
		if err != nil {
			return err
		}
		return nil

	case bool:
		var err error
		if e {
			err = binaryserializer.PutUint8(w, 0x01)
		} else {
			err = binaryserializer.PutUint8(w, 0x00)
		}
		if err != nil {
			return err
		}
		return nil

	case externalapi.DomainHash:
		_, err := w.Write(e.ByteSlice())
		if err != nil {
			return err
		}
		return nil

	case *externalapi.DomainHash:
		_, err := w.Write(e.ByteSlice())
		if err != nil {
			return err
		}
		return nil

	case externalapi.DomainTransactionID:
		_, err := w.Write(e.ByteSlice())
		if err != nil {
			return err
		}
		return nil

	case externalapi.DomainSubnetworkID:
		_, err := w.Write(e[:])
		if err != nil {
			return err
		}
		return nil

	case *externalapi.DomainSubnetworkID:
		_, err := w.Write(e[:])
		if err != nil {
			return err
		}
		return nil
	}

	return errors.Wrapf(errNoEncodingForType, "couldn't find a way to write type %T", element)
}

// WriteElements writes multiple items to w. It is equivalent to multiple
// calls to writeElement.
func WriteElements(w io.Writer, elements ...interface{}) error {
	for _, element := range elements {
		err := WriteElement(w, element)
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadElement reads the next sequence of bytes from r using little endian
// depending on the concrete type of element pointed to.
func ReadElement(r io.Reader, element interface{}) error {
	// Attempt to read the element based on the concrete type via fast
	// type assertions first.
	switch e := element.(type) {
	case *int16:
		rv, err := binaryserializer.Uint16(r)
		if err != nil {
			return err
		}
		*e = int16(rv)
		return nil

	case *uint16:
		rv, err := binaryserializer.Uint16(r)
		if err != nil {
			return err
		}
		*e = rv
		return nil
	case *int32:
		rv, err := binaryserializer.Uint32(r)
		if err != nil {
			return err
		}
		*e = int32(rv)
		return nil

	case *uint32:
		rv, err := binaryserializer.Uint32(r)
		if err != nil {
			return err
		}
		*e = rv
		return nil

	case *int64:
		rv, err := binaryserializer.Uint64(r)
		if err != nil {
			return err
		}
		*e = int64(rv)
		return nil

	case *uint64:
		rv, err := binaryserializer.Uint64(r)
		if err != nil {
			return err
		}
		*e = rv
		return nil

	case *uint8:
		rv, err := binaryserializer.Uint8(r)
		if err != nil {
			return err
		}
		*e = rv
		return nil

	case *bool:
		rv, err := binaryserializer.Uint8(r)
		if err != nil {
			return err
		}
		if rv == 0x00 {
			*e = false
		} else if rv == 0x01 {
			*e = true
		} else {
			return errors.Wrapf(errMalformed, "in order to keep serialization canonical, true has to"+
				" always be 0x01")
		}
		return nil
	}

	return errors.Wrapf(errNoEncodingForType, "couldn't find a way to read type %T", element)
}

// ReadElements reads multiple items from r. It is equivalent to multiple
// calls to ReadElement.
func ReadElements(r io.Reader, elements ...interface{}) error {
	for _, element := range elements {
		err := ReadElement(r, element)
		if err != nil {
			return err
		}
	}
	return nil
}

// IsMalformedError returns whether the error indicates a malformed data source
func IsMalformedError(err error) bool {
	return errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) || errors.Is(err, errMalformed)
}
