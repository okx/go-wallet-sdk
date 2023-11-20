/*
*
MIT License

Copyright (c) 2018 WavesPlatform
*/
package types

import (
	"fmt"
	"strings"
)

const (
	AliasPrefix    = "alias"
	aliasVersion   = 1
	aliasFixedSize = 4
)

// Recipient could be an Alias or an WavesAddress.
type Recipient struct {
	Address *WavesAddress
	Alias   *Alias
	len     int
}

// MarshalBinary makes bytes of the Recipient.
func (r *Recipient) MarshalBinary() ([]byte, error) {
	if r.Alias != nil {
		return r.Alias.MarshalBinary()
	}
	return r.Address[:], nil
}

// NewRecipientFromAddress creates the Recipient from given address.
func NewRecipientFromAddress(a WavesAddress) Recipient {
	return Recipient{Address: &a, len: WavesAddressSize}
}

// NewRecipientFromAlias creates a Recipient with the given Alias inside.
func NewRecipientFromAlias(a Alias) Recipient {
	return Recipient{Alias: &a, len: aliasFixedSize + len(a.Alias)}
}

func NewRecipientFromString(s string) (Recipient, error) {
	if strings.Contains(s, AliasPrefix) {
		a, err := NewAliasFromString(s)
		if err != nil {
			return Recipient{}, err
		}
		return NewRecipientFromAlias(*a), nil
	}
	a, err := NewAddressFromString(s)
	if err != nil {
		return Recipient{}, err
	}
	return NewRecipientFromAddress(a), nil
}

// Alias represents the nickname tha could be attached to the WavesAddress.
type Alias struct {
	Version byte
	Scheme  byte
	Alias   string
}

// NewAliasFromString creates an Alias from its string representation. Function does not check that the result is a valid Alias.
// String representation of an Alias should have a following format: "alias:<scheme>:<alias>". Scheme should be represented with a one-byte ASCII symbol.
func NewAliasFromString(s string) (*Alias, error) {
	ps := strings.Split(s, ":")
	if len(ps) != 3 {
		return nil, fmt.Errorf("incorrect alias string representation '%s'", s)
	}
	if ps[0] != AliasPrefix {
		return nil, fmt.Errorf("alias should start with prefix '%s'", AliasPrefix)
	}
	scheme := ps[1]
	if len(scheme) != 1 {
		return nil, fmt.Errorf("incorrect alias chainID '%s'", scheme)
	}
	a := Alias{Version: aliasVersion, Scheme: scheme[0], Alias: ps[2]}
	return &a, nil
}

// MarshalBinary converts the Alias to the slice of bytes. Just calls Bytes().
func (a *Alias) MarshalBinary() ([]byte, error) {
	return a.Bytes(), nil
}

// Bytes converts the Alias to the slice of bytes.
func (a *Alias) Bytes() []byte {
	al := len(a.Alias)
	buf := make([]byte, aliasFixedSize+al)
	buf[0] = a.Version
	buf[1] = a.Scheme
	PutStringWithUInt16Len(buf[2:], a.Alias)
	return buf
}

// MarshalJSON is a custom JSON marshalling function.
func (a Alias) MarshalJSON() ([]byte, error) {
	var sb strings.Builder
	sb.WriteRune('"')
	sb.WriteString(a.String())
	sb.WriteRune('"')
	return []byte(sb.String()), nil
}

// String converts the Alias to its 3-part string representation.
func (a Alias) String() string {
	sb := new(strings.Builder)
	sb.WriteString(AliasPrefix)
	sb.WriteRune(':')
	sb.WriteByte(a.Scheme)
	sb.WriteRune(':')
	sb.WriteString(a.Alias)
	return sb.String()
}

// MarshalJSON converts the Recipient to its JSON representation.
func (r Recipient) MarshalJSON() ([]byte, error) {
	if r.Alias != nil {
		return r.Alias.MarshalJSON()
	}
	return r.Address.MarshalJSON()
}
