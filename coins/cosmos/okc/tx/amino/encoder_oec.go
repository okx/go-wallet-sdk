package amino

import (
	"bytes"
	"time"
)

func EncodeByteSliceWithKeyToBuffer(w *bytes.Buffer, bz []byte, key ...byte) error {
	_, err := w.Write(key)
	if err != nil {
		return err
	}
	return EncodeByteSliceToBuffer(w, bz)
}

func EncodeStringWithKeyToBuffer(w *bytes.Buffer, s string, key ...byte) (err error) {
	_, err = w.Write(key)
	if err != nil {
		return
	}
	err = EncodeStringToBuffer(w, s)
	return
}

func EncodeUvarintWithKeyToBuffer(w *bytes.Buffer, u uint64, key ...byte) (err error) {
	_, err = w.Write(key)
	if err != nil {
		return
	}
	err = EncodeUvarintToBuffer(w, u)
	return
}

func EncodeBoolWithKeyToBuffer(w *bytes.Buffer, b bool, key ...byte) error {
	_, err := w.Write(key)
	if err != nil {
		return err
	}
	return EncodeBoolToBuffer(w, b)
}

func EncodeTimeWithKeyToBuffer(w *bytes.Buffer, t time.Time, key ...byte) error {
	_, err := w.Write(key)
	if err != nil {
		return err
	}
	return EncodeTimeToBuffer(w, t)
}
