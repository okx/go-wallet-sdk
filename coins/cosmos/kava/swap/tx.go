package swap

import (
	"encoding/json"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/cosmos/types"
	"io"
)

// MsgDeposit represents a message for depositing liquidity into a pool
type MsgDeposit struct {
	// depositor represents the address to deposit funds from
	Depositor string `protobuf:"bytes,1,opt,name=depositor,proto3" json:"depositor,omitempty"`
	// token_a represents one token of deposit pair
	TokenA types.Coin `protobuf:"bytes,2,opt,name=token_a,json=tokenA,proto3" json:"token_a"`
	// token_b represents one token of deposit pair
	TokenB types.Coin `protobuf:"bytes,3,opt,name=token_b,json=tokenB,proto3" json:"token_b"`
	// slippage represents the max decimal percentage price change
	Slippage types.Dec `protobuf:"bytes,4,opt,name=slippage,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"slippage"`
	// deadline represents the unix timestamp to complete the deposit by
	Deadline int64 `protobuf:"varint,5,opt,name=deadline,proto3" json:"deadline,string,omitempty"`
}

// MsgWithdraw represents a message for withdrawing liquidity from a pool
type MsgWithdraw struct {
	// from represents the address we are withdrawing for
	From string `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	// shares represents the amount of shares to withdraw
	Shares types.Int `protobuf:"bytes,2,opt,name=shares,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"shares"`
	// min_token_a represents the minimum a token to withdraw
	MinTokenA types.Coin `protobuf:"bytes,3,opt,name=min_token_a,json=minTokenA,proto3" json:"min_token_a"`
	// min_token_a represents the minimum a token to withdraw
	MinTokenB types.Coin `protobuf:"bytes,4,opt,name=min_token_b,json=minTokenB,proto3" json:"min_token_b"`
	// deadline represents the unix timestamp to complete the withdraw by
	Deadline int64 `protobuf:"varint,5,opt,name=deadline,proto3" json:"deadline,string,omitempty"`
}

// MsgSwapExactForTokens represents a message for trading exact coinA for coinB
type MsgSwapExactForTokens struct {
	// represents the address swaping the tokens
	Requester string `protobuf:"bytes,1,opt,name=requester,proto3" json:"requester,omitempty"`
	// exact_token_a represents the exact amount to swap for token_b
	ExactTokenA types.Coin `protobuf:"bytes,2,opt,name=exact_token_a,json=exactTokenA,proto3" json:"exact_token_a"`
	// token_b represents the desired token_b to swap for
	TokenB types.Coin `protobuf:"bytes,3,opt,name=token_b,json=tokenB,proto3" json:"token_b"`
	// slippage represents the maximum change in token_b allowed
	Slippage types.Dec `protobuf:"bytes,4,opt,name=slippage,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"slippage"`
	// deadline represents the unix timestamp to complete the swap by
	Deadline int64 `protobuf:"varint,5,opt,name=deadline,proto3" json:"deadline,string,omitempty"`
}

// MsgSwapForExactTokens represents a message for trading coinA for an exact
// coinB
type MsgSwapForExactTokens struct {
	// represents the address swaping the tokens
	Requester string `protobuf:"bytes,1,opt,name=requester,proto3" json:"requester,omitempty"`
	// token_a represents the desired token_a to swap for
	TokenA types.Coin `protobuf:"bytes,2,opt,name=token_a,json=tokenA,proto3" json:"token_a"`
	// exact_token_b represents the exact token b amount to swap for token a
	ExactTokenB types.Coin `protobuf:"bytes,3,opt,name=exact_token_b,json=exactTokenB,proto3" json:"exact_token_b"`
	// slippage represents the maximum change in token_a allowed
	Slippage types.Dec `protobuf:"bytes,4,opt,name=slippage,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"slippage"`
	// deadline represents the unix timestamp to complete the swap by
	Deadline int64 `protobuf:"varint,5,opt,name=deadline,proto3" json:"deadline,string,omitempty"`
}

func (m *MsgDeposit) MessageName() string {
	return "kava.swap.v1beta1.MsgDeposit"
}

func (m *MsgWithdraw) MessageName() string {
	return "kava.swap.v1beta1.MsgWithdraw"
}

func (m *MsgSwapExactForTokens) MessageName() string {
	return "kava.swap.v1beta1.MsgSwapExactForTokens"
}

func (m *MsgSwapForExactTokens) MessageName() string {
	return "kava.swap.v1beta1.MsgSwapForExactTokens"
}

func (m *MsgDeposit) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgDeposit) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgDeposit) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Deadline != 0 {
		i = types.EncodeVarintTx(dAtA, i, uint64(m.Deadline))
		i--
		dAtA[i] = 0x28
	}
	{
		size := m.Slippage.Size()
		i -= size
		if _, err := m.Slippage.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = types.EncodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size, err := m.TokenB.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = types.EncodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size, err := m.TokenA.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = types.EncodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Depositor) > 0 {
		i -= len(m.Depositor)
		copy(dAtA[i:], m.Depositor)
		i = types.EncodeVarintTx(dAtA, i, uint64(len(m.Depositor)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgWithdraw) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgWithdraw) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgWithdraw) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Deadline != 0 {
		i = types.EncodeVarintTx(dAtA, i, uint64(m.Deadline))
		i--
		dAtA[i] = 0x28
	}
	{
		size, err := m.MinTokenB.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = types.EncodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size, err := m.MinTokenA.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = types.EncodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.Shares.Size()
		i -= size
		if _, err := m.Shares.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = types.EncodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.From) > 0 {
		i -= len(m.From)
		copy(dAtA[i:], m.From)
		i = types.EncodeVarintTx(dAtA, i, uint64(len(m.From)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgSwapExactForTokens) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSwapExactForTokens) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSwapExactForTokens) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Deadline != 0 {
		i = types.EncodeVarintTx(dAtA, i, uint64(m.Deadline))
		i--
		dAtA[i] = 0x28
	}
	{
		size := m.Slippage.Size()
		i -= size
		if _, err := m.Slippage.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = types.EncodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size, err := m.TokenB.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = types.EncodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size, err := m.ExactTokenA.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = types.EncodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Requester) > 0 {
		i -= len(m.Requester)
		copy(dAtA[i:], m.Requester)
		i = types.EncodeVarintTx(dAtA, i, uint64(len(m.Requester)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgSwapForExactTokens) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgSwapForExactTokens) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgSwapForExactTokens) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Deadline != 0 {
		i = types.EncodeVarintTx(dAtA, i, uint64(m.Deadline))
		i--
		dAtA[i] = 0x28
	}
	{
		size := m.Slippage.Size()
		i -= size
		if _, err := m.Slippage.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = types.EncodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x22
	{
		size, err := m.ExactTokenB.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = types.EncodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size, err := m.TokenA.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = types.EncodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Requester) > 0 {
		i -= len(m.Requester)
		copy(dAtA[i:], m.Requester)
		i = types.EncodeVarintTx(dAtA, i, uint64(len(m.Requester)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgDeposit) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Depositor)
	if l > 0 {
		n += 1 + l + types.SovTx(uint64(l))
	}
	l = m.TokenA.Size()
	n += 1 + l + types.SovTx(uint64(l))
	l = m.TokenB.Size()
	n += 1 + l + types.SovTx(uint64(l))
	l = m.Slippage.Size()
	n += 1 + l + types.SovTx(uint64(l))
	if m.Deadline != 0 {
		n += 1 + types.SovTx(uint64(m.Deadline))
	}
	return n
}

func (m *MsgWithdraw) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.From)
	if l > 0 {
		n += 1 + l + types.SovTx(uint64(l))
	}
	l = m.Shares.Size()
	n += 1 + l + types.SovTx(uint64(l))
	l = m.MinTokenA.Size()
	n += 1 + l + types.SovTx(uint64(l))
	l = m.MinTokenB.Size()
	n += 1 + l + types.SovTx(uint64(l))
	if m.Deadline != 0 {
		n += 1 + types.SovTx(uint64(m.Deadline))
	}
	return n
}

func (m *MsgSwapExactForTokens) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Requester)
	if l > 0 {
		n += 1 + l + types.SovTx(uint64(l))
	}
	l = m.ExactTokenA.Size()
	n += 1 + l + types.SovTx(uint64(l))
	l = m.TokenB.Size()
	n += 1 + l + types.SovTx(uint64(l))
	l = m.Slippage.Size()
	n += 1 + l + types.SovTx(uint64(l))
	if m.Deadline != 0 {
		n += 1 + types.SovTx(uint64(m.Deadline))
	}
	return n
}

func (m *MsgSwapForExactTokens) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Requester)
	if l > 0 {
		n += 1 + l + types.SovTx(uint64(l))
	}
	l = m.TokenA.Size()
	n += 1 + l + types.SovTx(uint64(l))
	l = m.ExactTokenB.Size()
	n += 1 + l + types.SovTx(uint64(l))
	l = m.Slippage.Size()
	n += 1 + l + types.SovTx(uint64(l))
	if m.Deadline != 0 {
		n += 1 + types.SovTx(uint64(m.Deadline))
	}
	return n
}

func (m *MsgDeposit) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return types.ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgDeposit: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgDeposit: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Depositor", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return types.ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return types.ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return types.ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Depositor = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenA", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return types.ErrIntOverflowTx
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
				return types.ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return types.ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TokenA.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenB", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return types.ErrIntOverflowTx
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
				return types.ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return types.ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TokenB.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Slippage", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return types.ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return types.ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return types.ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Slippage.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Deadline", wireType)
			}
			m.Deadline = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return types.ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Deadline |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := types.SkipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return types.ErrInvalidLengthTx
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

func (m *MsgWithdraw) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return types.ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgWithdraw: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgWithdraw: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field From", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return types.ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return types.ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return types.ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.From = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Shares", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return types.ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return types.ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return types.ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Shares.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinTokenA", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return types.ErrIntOverflowTx
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
				return types.ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return types.ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinTokenA.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinTokenB", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return types.ErrIntOverflowTx
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
				return types.ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return types.ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinTokenB.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Deadline", wireType)
			}
			m.Deadline = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return types.ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Deadline |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := types.SkipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return types.ErrInvalidLengthTx
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

func (m *MsgSwapExactForTokens) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return types.ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgSwapExactForTokens: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSwapExactForTokens: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Requester", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return types.ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return types.ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return types.ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Requester = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExactTokenA", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return types.ErrIntOverflowTx
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
				return types.ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return types.ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ExactTokenA.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenB", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return types.ErrIntOverflowTx
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
				return types.ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return types.ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TokenB.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Slippage", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return types.ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return types.ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return types.ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Slippage.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Deadline", wireType)
			}
			m.Deadline = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return types.ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Deadline |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := types.SkipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return types.ErrInvalidLengthTx
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

func (m *MsgSwapForExactTokens) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return types.ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgSwapForExactTokens: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgSwapForExactTokens: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Requester", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return types.ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return types.ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return types.ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Requester = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TokenA", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return types.ErrIntOverflowTx
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
				return types.ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return types.ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TokenA.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExactTokenB", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return types.ErrIntOverflowTx
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
				return types.ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return types.ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ExactTokenB.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Slippage", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return types.ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return types.ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return types.ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Slippage.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Deadline", wireType)
			}
			m.Deadline = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return types.ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Deadline |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := types.SkipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return types.ErrInvalidLengthTx
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

func init() {
	types.RegistryMessageConverter("swap/MsgDeposit", func(s string) *types.Any {
		msg := MsgDeposit{}
		err := json.Unmarshal([]byte(s), &msg)
		if err != nil {
			panic(err)
		}
		any, err := types.NewAnyWithValue(&msg)
		if err != nil {
			panic(err)
		}
		return any
	})

	types.RegistryMessageConverter("swap/MsgWithdraw", func(s string) *types.Any {
		msg := MsgWithdraw{}
		err := json.Unmarshal([]byte(s), &msg)
		if err != nil {
			panic(err)
		}
		any, err := types.NewAnyWithValue(&msg)
		if err != nil {
			panic(err)
		}
		return any
	})

	types.RegistryMessageConverter("swap/MsgSwapExactForTokens", func(s string) *types.Any {
		msg := MsgSwapExactForTokens{}
		err := json.Unmarshal([]byte(s), &msg)
		if err != nil {
			panic(err)
		}
		any, err := types.NewAnyWithValue(&msg)
		if err != nil {
			panic(err)
		}
		return any
	})

	types.RegistryMessageConverter("swap/MsgSwapForExactTokens", func(s string) *types.Any {
		msg := MsgSwapForExactTokens{}
		err := json.Unmarshal([]byte(s), &msg)
		if err != nil {
			panic(err)
		}
		any, err := types.NewAnyWithValue(&msg)
		if err != nil {
			panic(err)
		}
		return any
	})
}
