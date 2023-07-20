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
	"errors"
	"fmt"
	"io"
	math_bits "math/bits"
	time "time"
)

type Duration struct {
	// Signed seconds of the span of time. Must be from -315,576,000,000
	// to +315,576,000,000 inclusive. Note: these bounds are computed from:
	// 60 sec/min * 60 min/hr * 24 hr/day * 365.25 days/year * 10000 years
	Seconds int64 `protobuf:"varint,1,opt,name=seconds,proto3" json:"seconds,omitempty"`
	// Signed fractions of a second at nanosecond resolution of the span
	// of time. Durations less than one second are represented with a 0
	// `seconds` field and a positive or negative `nanos` field. For durations
	// of one second or more, a non-zero value for the `nanos` field must be
	// of the same sign as the `seconds` field. Must be from -999,999,999
	// to +999,999,999 inclusive.
	Nanos                int32    `protobuf:"varint,2,opt,name=nanos,proto3" json:"nanos,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Duration) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Duration) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Duration) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Nanos != 0 {
		i = encodeVarintDuration(dAtA, i, uint64(m.Nanos))
		i--
		dAtA[i] = 0x10
	}
	if m.Seconds != 0 {
		i = encodeVarintDuration(dAtA, i, uint64(m.Seconds))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintDuration(dAtA []byte, offset int, v uint64) int {
	offset -= sovDuration(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Duration) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Seconds != 0 {
		n += 1 + sovDuration(uint64(m.Seconds))
	}
	if m.Nanos != 0 {
		n += 1 + sovDuration(uint64(m.Nanos))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovDuration(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozDuration(x uint64) (n int) {
	return sovDuration(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Duration) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowDuration
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
			return fmt.Errorf("proto: Duration: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Duration: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Seconds", wireType)
			}
			m.Seconds = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDuration
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Seconds |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nanos", wireType)
			}
			m.Nanos = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowDuration
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Nanos |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipDuration(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthDuration
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func skipDuration(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowDuration
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
					return 0, ErrIntOverflowDuration
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
					return 0, ErrIntOverflowDuration
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
				return 0, ErrInvalidLengthDuration
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupDuration
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthDuration
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthDuration        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowDuration          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupDuration = fmt.Errorf("proto: unexpected end of group")
)

// DurationProto converts a time.Duration to a Duration.
func DurationProto(d time.Duration) *Duration {
	nanos := d.Nanoseconds()
	secs := nanos / 1e9
	nanos -= secs * 1e9
	return &Duration{
		Seconds: secs,
		Nanos:   int32(nanos),
	}
}

const (
	// Range of a Duration in seconds, as specified in
	// google/protobuf/duration.proto. This is about 10,000 years in seconds.
	maxSeconds = int64(10000 * 365.25 * 24 * 60 * 60)
	minSeconds = -maxSeconds
)

func validateDuration(d *Duration) error {
	if d == nil {
		return errors.New("duration: nil Duration")
	}
	if d.Seconds < minSeconds || d.Seconds > maxSeconds {
		return fmt.Errorf("duration: %#v: seconds out of range", d)
	}
	if d.Nanos <= -1e9 || d.Nanos >= 1e9 {
		return fmt.Errorf("duration: %#v: nanos out of range", d)
	}
	// Seconds and Nanos must have the same sign, unless d.Nanos is zero.
	if (d.Seconds < 0 && d.Nanos > 0) || (d.Seconds > 0 && d.Nanos < 0) {
		return fmt.Errorf("duration: %#v: seconds and nanos have different signs", d)
	}
	return nil
}

func DurationFromProto(p *Duration) (time.Duration, error) {
	if err := validateDuration(p); err != nil {
		return 0, err
	}
	d := time.Duration(p.Seconds) * time.Second
	if int64(d/time.Second) != p.Seconds {
		return 0, fmt.Errorf("duration: %#v is out of range for time.Duration", p)
	}
	if p.Nanos != 0 {
		d += time.Duration(p.Nanos) * time.Nanosecond
		if (d < 0) != (p.Nanos < 0) {
			return 0, fmt.Errorf("duration: %#v is out of range for time.Duration", p)
		}
	}
	return d, nil
}

func SizeOfStdDuration(d time.Duration) int {
	dur := DurationProto(d)
	return dur.Size()
}

func StdDurationMarshal(d time.Duration) ([]byte, error) {
	size := SizeOfStdDuration(d)
	buf := make([]byte, size)
	_, err := StdDurationMarshalTo(d, buf)
	return buf, err
}

func StdDurationMarshalTo(d time.Duration, data []byte) (int, error) {
	dur := DurationProto(d)
	return dur.MarshalTo(data)
}

func StdDurationUnmarshal(d *time.Duration, data []byte) error {
	dur := &Duration{}
	if err := dur.Unmarshal(data); err != nil {
		return err
	}
	dd, err := DurationFromProto(dur)
	if err != nil {
		return err
	}
	*d = dd
	return nil
}
