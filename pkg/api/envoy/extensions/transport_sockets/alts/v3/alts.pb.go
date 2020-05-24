// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: envoy/extensions/transport_sockets/alts/v3/alts.proto

package envoy_extensions_transport_sockets_alts_v3

import (
	fmt "fmt"
	_ "github.com/cncf/udpa/go/udpa/annotations"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Configuration for ALTS transport socket. This provides Google's ALTS protocol to Envoy.
// https://cloud.google.com/security/encryption-in-transit/application-layer-transport-security/
type Alts struct {
	// The location of a handshaker service, this is usually 169.254.169.254:8080
	// on GCE.
	HandshakerService string `protobuf:"bytes,1,opt,name=handshaker_service,json=handshakerService,proto3" json:"handshaker_service,omitempty"`
	// The acceptable service accounts from peer, peers not in the list will be rejected in the
	// handshake validation step. If empty, no validation will be performed.
	PeerServiceAccounts  []string `protobuf:"bytes,2,rep,name=peer_service_accounts,json=peerServiceAccounts,proto3" json:"peer_service_accounts,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Alts) Reset()         { *m = Alts{} }
func (m *Alts) String() string { return proto.CompactTextString(m) }
func (*Alts) ProtoMessage()    {}
func (*Alts) Descriptor() ([]byte, []int) {
	return fileDescriptor_3b780b78c099503e, []int{0}
}
func (m *Alts) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Alts) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Alts.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Alts) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Alts.Merge(m, src)
}
func (m *Alts) XXX_Size() int {
	return m.Size()
}
func (m *Alts) XXX_DiscardUnknown() {
	xxx_messageInfo_Alts.DiscardUnknown(m)
}

var xxx_messageInfo_Alts proto.InternalMessageInfo

func (m *Alts) GetHandshakerService() string {
	if m != nil {
		return m.HandshakerService
	}
	return ""
}

func (m *Alts) GetPeerServiceAccounts() []string {
	if m != nil {
		return m.PeerServiceAccounts
	}
	return nil
}

func init() {
	proto.RegisterType((*Alts)(nil), "envoy.extensions.transport_sockets.alts.v3.Alts")
}

func init() {
	proto.RegisterFile("envoy/extensions/transport_sockets/alts/v3/alts.proto", fileDescriptor_3b780b78c099503e)
}

var fileDescriptor_3b780b78c099503e = []byte{
	// 308 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x90, 0xb1, 0x4a, 0x03, 0x31,
	0x1c, 0xc6, 0x49, 0x2d, 0x95, 0x66, 0xd2, 0x13, 0xb1, 0x14, 0x3c, 0xaa, 0x53, 0x71, 0x48, 0xb0,
	0xc5, 0x22, 0x6e, 0xbd, 0x27, 0x28, 0x75, 0x74, 0x28, 0x7f, 0xaf, 0xb1, 0x0d, 0x3d, 0x92, 0x90,
	0xfc, 0x1b, 0xda, 0xcd, 0xd1, 0x67, 0xf0, 0x11, 0x7c, 0x04, 0x77, 0xa1, 0xa3, 0xbe, 0x81, 0xf4,
	0x31, 0x9c, 0xe4, 0x72, 0x27, 0x37, 0xdc, 0xe2, 0x94, 0xc0, 0xf7, 0xfb, 0xbe, 0xe4, 0xfb, 0xe8,
	0x8d, 0x50, 0x5e, 0x6f, 0xb9, 0xd8, 0xa0, 0x50, 0x4e, 0x6a, 0xe5, 0x38, 0x5a, 0x50, 0xce, 0x68,
	0x8b, 0x33, 0xa7, 0xd3, 0x95, 0x40, 0xc7, 0x21, 0x43, 0xc7, 0xfd, 0x30, 0x9c, 0xcc, 0x58, 0x8d,
	0x3a, 0xba, 0x0a, 0x36, 0x56, 0xd9, 0x58, 0xcd, 0xc6, 0x02, 0xee, 0x87, 0xdd, 0xf3, 0xf5, 0xdc,
	0x00, 0x07, 0xa5, 0x34, 0x02, 0x86, 0x27, 0x1c, 0x02, 0xae, 0xcb, 0xa8, 0xee, 0x45, 0x4d, 0xf6,
	0xc2, 0xe6, 0x99, 0x52, 0x2d, 0x4a, 0xe4, 0xcc, 0x43, 0x26, 0xe7, 0x80, 0x82, 0xff, 0x5d, 0x0a,
	0xe1, 0xf2, 0x8d, 0xd0, 0xe6, 0x38, 0x43, 0x17, 0x8d, 0x68, 0xb4, 0x04, 0x35, 0x77, 0x4b, 0x58,
	0x09, 0x3b, 0x73, 0xc2, 0x7a, 0x99, 0x8a, 0x0e, 0xe9, 0x91, 0x7e, 0x3b, 0x39, 0xfc, 0x49, 0x9a,
	0xb6, 0xd1, 0x23, 0xd3, 0xe3, 0x0a, 0xb9, 0x2f, 0x88, 0x68, 0x40, 0x4f, 0x8d, 0xa8, 0x1c, 0x33,
	0x48, 0x53, 0xbd, 0x56, 0xe8, 0x3a, 0x8d, 0xde, 0x41, 0xbf, 0x3d, 0x3d, 0xc9, 0xc5, 0x92, 0x1d,
	0x97, 0xd2, 0xdd, 0xe8, 0xf5, 0xe3, 0x25, 0xbe, 0xa6, 0xbc, 0x98, 0x20, 0xd5, 0xea, 0x49, 0x2e,
	0x6a, 0xf5, 0xcb, 0xf6, 0x03, 0xc8, 0xcc, 0x12, 0x58, 0xfe, 0xc7, 0xe4, 0x61, 0xb7, 0x8f, 0xc9,
	0xe7, 0x3e, 0x26, 0xdf, 0xfb, 0x98, 0xbc, 0x3f, 0xef, 0xbe, 0x5a, 0x8d, 0x23, 0x42, 0x6f, 0xa5,
	0x66, 0x21, 0xc9, 0x58, 0xbd, 0xd9, 0xb2, 0xff, 0xef, 0x9a, 0xb4, 0xf3, 0xd4, 0x49, 0xbe, 0xc3,
	0x84, 0x3c, 0xb6, 0xc2, 0x20, 0xc3, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x36, 0x52, 0xd2, 0x51,
	0xd0, 0x01, 0x00, 0x00,
}

func (m *Alts) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Alts) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Alts) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.PeerServiceAccounts) > 0 {
		for iNdEx := len(m.PeerServiceAccounts) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.PeerServiceAccounts[iNdEx])
			copy(dAtA[i:], m.PeerServiceAccounts[iNdEx])
			i = encodeVarintAlts(dAtA, i, uint64(len(m.PeerServiceAccounts[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.HandshakerService) > 0 {
		i -= len(m.HandshakerService)
		copy(dAtA[i:], m.HandshakerService)
		i = encodeVarintAlts(dAtA, i, uint64(len(m.HandshakerService)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintAlts(dAtA []byte, offset int, v uint64) int {
	offset -= sovAlts(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Alts) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.HandshakerService)
	if l > 0 {
		n += 1 + l + sovAlts(uint64(l))
	}
	if len(m.PeerServiceAccounts) > 0 {
		for _, s := range m.PeerServiceAccounts {
			l = len(s)
			n += 1 + l + sovAlts(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovAlts(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozAlts(x uint64) (n int) {
	return sovAlts(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Alts) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAlts
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
			return fmt.Errorf("proto: Alts: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Alts: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HandshakerService", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAlts
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
				return ErrInvalidLengthAlts
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAlts
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.HandshakerService = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PeerServiceAccounts", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAlts
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
				return ErrInvalidLengthAlts
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAlts
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PeerServiceAccounts = append(m.PeerServiceAccounts, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAlts(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAlts
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthAlts
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
func skipAlts(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAlts
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
					return 0, ErrIntOverflowAlts
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
					return 0, ErrIntOverflowAlts
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
				return 0, ErrInvalidLengthAlts
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupAlts
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthAlts
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthAlts        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAlts          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupAlts = fmt.Errorf("proto: unexpected end of group")
)
