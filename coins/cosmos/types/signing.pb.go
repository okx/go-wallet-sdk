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

import (
	fmt "fmt"
	io "io"
	math_bits "math/bits"
)

// SignMode represents a signing mode with its own security guarantees.
type SignMode int32

const (
	// SIGN_MODE_UNSPECIFIED specifies an unknown signing mode and will be
	// rejected
	SignMode_SIGN_MODE_UNSPECIFIED SignMode = 0
	// SIGN_MODE_DIRECT specifies a signing mode which uses SignDoc and is
	// verified with raw bytes from Tx
	SignMode_SIGN_MODE_DIRECT SignMode = 1
	// SIGN_MODE_TEXTUAL is a future signing mode that will verify some
	// human-readable textual representation on top of the binary representation
	// from SIGN_MODE_DIRECT
	SignMode_SIGN_MODE_TEXTUAL SignMode = 2
	// SIGN_MODE_LEGACY_AMINO_JSON is a backwards compatibility mode which uses
	// Amino JSON and will be removed in the future
	SignMode_SIGN_MODE_LEGACY_AMINO_JSON SignMode = 127
)

var SignMode_name = map[int32]string{
	0:   "SIGN_MODE_UNSPECIFIED",
	1:   "SIGN_MODE_DIRECT",
	2:   "SIGN_MODE_TEXTUAL",
	127: "SIGN_MODE_LEGACY_AMINO_JSON",
}

var SignMode_value = map[string]int32{
	"SIGN_MODE_UNSPECIFIED":       0,
	"SIGN_MODE_DIRECT":            1,
	"SIGN_MODE_TEXTUAL":           2,
	"SIGN_MODE_LEGACY_AMINO_JSON": 127,
}

// SignatureDescriptor is a convenience type which represents the full data for
// a signature including the public key of the signer, signing modes and the
// signature itself. It is primarily used for coordinating signatures between
// clients.
type SignatureDescriptor struct {
	// public_key is the public key of the signer
	PublicKey *Any                      `protobuf:"bytes,1,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	Data      *SignatureDescriptor_Data `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	// sequence is the sequence of the account, which describes the
	// number of committed transactions signed by a given address. It is used to prevent
	// replay attacks.
	Sequence uint64 `protobuf:"varint,3,opt,name=sequence,proto3" json:"sequence,omitempty"`
}

// Data represents signature data
type SignatureDescriptor_Data struct {
	// sum is the oneof that specifies whether this represents single or multi-signature data
	//
	// Types that are valid to be assigned to Sum:
	//	*SignatureDescriptor_Data_Single_
	//	*SignatureDescriptor_Data_Multi_
	Sum isSignatureDescriptor_Data_Sum `protobuf_oneof:"sum"`
}

type isSignatureDescriptor_Data_Sum interface {
	isSignatureDescriptor_Data_Sum()
	MarshalTo([]byte) (int, error)
	Size() int
}

type SignatureDescriptor_Data_Single_ struct {
	Single *SignatureDescriptor_Data_Single `protobuf:"bytes,1,opt,name=single,proto3,oneof" json:"single,omitempty"`
}
type SignatureDescriptor_Data_Multi_ struct {
	Multi *SignatureDescriptor_Data_Multi `protobuf:"bytes,2,opt,name=multi,proto3,oneof" json:"multi,omitempty"`
}

func (*SignatureDescriptor_Data_Single_) isSignatureDescriptor_Data_Sum() {}
func (*SignatureDescriptor_Data_Multi_) isSignatureDescriptor_Data_Sum()  {}

// Single is the signature data for a single signer
type SignatureDescriptor_Data_Single struct {
	// mode is the signing mode of the single signer
	Mode SignMode `protobuf:"varint,1,opt,name=mode,proto3,enum=cosmos.tx.signing.v1beta1.SignMode" json:"mode,omitempty"`
	// signature is the raw signature bytes
	Signature []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
}

// Multi is the signature data for a multisig public key
type SignatureDescriptor_Data_Multi struct {
	// bitarray specifies which keys within the multisig are signing
	Bitarray *CompactBitArray `protobuf:"bytes,1,opt,name=bitarray,proto3" json:"bitarray,omitempty"`
	// signatures is the signatures of the multi-signature
	Signatures []*SignatureDescriptor_Data `protobuf:"bytes,2,rep,name=signatures,proto3" json:"signatures,omitempty"`
}

// SignatureDescriptors wraps multiple SignatureDescriptor's.
type SignatureDescriptors struct {
	// signatures are the signature descriptors
	Signatures []*SignatureDescriptor `protobuf:"bytes,1,rep,name=signatures,proto3" json:"signatures,omitempty"`
}

func (m *SignatureDescriptors) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SignatureDescriptors) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SignatureDescriptors) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Signatures) > 0 {
		for iNdEx := len(m.Signatures) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Signatures[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSigning(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *SignatureDescriptor) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SignatureDescriptor) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SignatureDescriptor) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Sequence != 0 {
		i = encodeVarintSigning(dAtA, i, uint64(m.Sequence))
		i--
		dAtA[i] = 0x18
	}
	if m.Data != nil {
		{
			size, err := m.Data.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSigning(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.PublicKey != nil {
		{
			size, err := m.PublicKey.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSigning(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *SignatureDescriptor_Data) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SignatureDescriptor_Data) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SignatureDescriptor_Data) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Sum != nil {
		{
			size := m.Sum.Size()
			i -= size
			if _, err := m.Sum.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
		}
	}
	return len(dAtA) - i, nil
}

func (m *SignatureDescriptor_Data_Single_) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SignatureDescriptor_Data_Single_) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Single != nil {
		{
			size, err := m.Single.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSigning(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}
func (m *SignatureDescriptor_Data_Multi_) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SignatureDescriptor_Data_Multi_) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	if m.Multi != nil {
		{
			size, err := m.Multi.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSigning(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	return len(dAtA) - i, nil
}
func (m *SignatureDescriptor_Data_Single) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SignatureDescriptor_Data_Single) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SignatureDescriptor_Data_Single) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Signature) > 0 {
		i -= len(m.Signature)
		copy(dAtA[i:], m.Signature)
		i = encodeVarintSigning(dAtA, i, uint64(len(m.Signature)))
		i--
		dAtA[i] = 0x12
	}
	if m.Mode != 0 {
		i = encodeVarintSigning(dAtA, i, uint64(m.Mode))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *SignatureDescriptor_Data_Multi) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SignatureDescriptor_Data_Multi) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *SignatureDescriptor_Data_Multi) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Signatures) > 0 {
		for iNdEx := len(m.Signatures) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Signatures[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintSigning(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Bitarray != nil {
		{
			size, err := m.Bitarray.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintSigning(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintSigning(dAtA []byte, offset int, v uint64) int {
	offset -= sovSigning(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *SignatureDescriptors) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Signatures) > 0 {
		for _, e := range m.Signatures {
			l = e.Size()
			n += 1 + l + sovSigning(uint64(l))
		}
	}
	return n
}

func (m *SignatureDescriptor) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PublicKey != nil {
		l = m.PublicKey.Size()
		n += 1 + l + sovSigning(uint64(l))
	}
	if m.Data != nil {
		l = m.Data.Size()
		n += 1 + l + sovSigning(uint64(l))
	}
	if m.Sequence != 0 {
		n += 1 + sovSigning(uint64(m.Sequence))
	}
	return n
}

func (m *SignatureDescriptor_Data) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Sum != nil {
		n += m.Sum.Size()
	}
	return n
}

func (m *SignatureDescriptor_Data_Single_) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Single != nil {
		l = m.Single.Size()
		n += 1 + l + sovSigning(uint64(l))
	}
	return n
}
func (m *SignatureDescriptor_Data_Multi_) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Multi != nil {
		l = m.Multi.Size()
		n += 1 + l + sovSigning(uint64(l))
	}
	return n
}
func (m *SignatureDescriptor_Data_Single) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Mode != 0 {
		n += 1 + sovSigning(uint64(m.Mode))
	}
	l = len(m.Signature)
	if l > 0 {
		n += 1 + l + sovSigning(uint64(l))
	}
	return n
}

func (m *SignatureDescriptor_Data_Multi) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Bitarray != nil {
		l = m.Bitarray.Size()
		n += 1 + l + sovSigning(uint64(l))
	}
	if len(m.Signatures) > 0 {
		for _, e := range m.Signatures {
			l = e.Size()
			n += 1 + l + sovSigning(uint64(l))
		}
	}
	return n
}

func sovSigning(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}

func (m *SignatureDescriptors) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSigning
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SignatureDescriptors: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SignatureDescriptors: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signatures", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSigning
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSigning
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSigning
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signatures = append(m.Signatures, &SignatureDescriptor{})
			if err := m.Signatures[len(m.Signatures)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSigning(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSigning
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SignatureDescriptor) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSigning
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SignatureDescriptor: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SignatureDescriptor: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PublicKey", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSigning
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSigning
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSigning
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.PublicKey == nil {
				m.PublicKey = &Any{}
			}
			if err := m.PublicKey.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Data", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSigning
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSigning
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSigning
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Data == nil {
				m.Data = &SignatureDescriptor_Data{}
			}
			if err := m.Data.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sequence", wireType)
			}
			m.Sequence = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSigning
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Sequence |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipSigning(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSigning
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SignatureDescriptor_Data) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSigning
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Data: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Data: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Single", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSigning
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSigning
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSigning
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &SignatureDescriptor_Data_Single{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Sum = &SignatureDescriptor_Data_Single_{v}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Multi", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSigning
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSigning
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSigning
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &SignatureDescriptor_Data_Multi{}
			if err := v.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Sum = &SignatureDescriptor_Data_Multi_{v}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSigning(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSigning
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SignatureDescriptor_Data_Single) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSigning
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Single: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Single: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Mode", wireType)
			}
			m.Mode = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSigning
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Mode |= SignMode(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signature", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSigning
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthSigning
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthSigning
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signature = append(m.Signature[:0], dAtA[iNdEx:postIndex]...)
			if m.Signature == nil {
				m.Signature = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSigning(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSigning
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *SignatureDescriptor_Data_Multi) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSigning
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Multi: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Multi: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bitarray", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSigning
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSigning
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSigning
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Bitarray == nil {
				m.Bitarray = &CompactBitArray{}
			}
			if err := m.Bitarray.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signatures", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSigning
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthSigning
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthSigning
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signatures = append(m.Signatures, &SignatureDescriptor_Data{})
			if err := m.Signatures[len(m.Signatures)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSigning(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSigning
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipSigning(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSigning
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSigning
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowSigning
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthSigning
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSigning
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSigning
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSigning        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSigning          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSigning = fmt.Errorf("proto: unexpected end of group")
)
