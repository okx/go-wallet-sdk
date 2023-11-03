/*
*

	Copyright Cosmos-SDK Authors
	Copyright 2016 All in Bits, Inc

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	    http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/
package types

type Any struct {
	// nolint
	TypeUrl string `protobuf:"bytes,1,opt,name=type_url,json=typeUrl,proto3" json:"type_url,omitempty"`
	// Must be a valid serialized protocol buffer of the above specified type.
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`

	// nolint
	XXX_unrecognized []byte `json:"-"`
}

type Message interface {
	MessageName() string
	Marshal() (dAtA []byte, err error)
	Unmarshal(dAtA []byte) error
}

type PackError struct {
	msg string
}

func NewPackError(detail string) *PackError {
	return &PackError{detail}
}
func (p PackError) Error() string {
	return p.msg
}

// NewAnyWithValue constructs a new Any packed with the value provided or
// returns an error if that value couldn't be packed. This also caches
// the packed value so that it can be retrieved from GetCachedValue without
// unmarshaling
func NewAnyWithValue(v Message) (*Any, error) {
	if v == nil {
		return nil, NewPackError("Expecting non nil value to create a new Any")
	}

	bz, err := v.Marshal()
	if err != nil {
		return nil, err
	}

	return &Any{
		TypeUrl: "/" + v.MessageName(),
		Value:   bz,
	}, nil
}

func NewAnyWithValueAndName(v Message) (*Any, error) {
	if v == nil {
		return nil, NewPackError("Expecting non nil value to create a new Any")
	}

	bz, err := v.Marshal()
	if err != nil {
		return nil, err
	}

	return &Any{
		TypeUrl: "/cosmwasm.wasm.v1.MsgExecuteContract",
		Value:   bz,
	}, nil
}
