// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kv/kvserver/storagepb/lease_status.proto

package storagepb

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import roachpb "github.com/cockroachdb/cockroach/pkg/roachpb"
import hlc "github.com/cockroachdb/cockroach/pkg/util/hlc"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type LeaseState int32

const (
	// ERROR indicates that the lease can't be used or acquired.
	LeaseState_ERROR LeaseState = 0
	// VALID indicates that the lease can be used.
	LeaseState_VALID LeaseState = 1
	// STASIS indicates that the lease has not expired, but can't be
	// used because it is close to expiration (a stasis period at the
	// end of each lease is one of the ways we handle clock
	// uncertainty). A lease in STASIS may become VALID for the same
	// leaseholder after a successful RequestLease (for expiration-based
	// leases) or Heartbeat (for epoch-based leases). A lease may not
	// change hands while it is in stasis; would-be acquirers must wait
	// for the stasis period to expire.
	//
	// The point of the stasis period is to prevent reads on the old leaseholder
	// (the one whose stasis we're talking about) from missing to see writes
	// performed under the next lease (held by someone else) when these writes
	// should fall in the uncertainty window. Even without the stasis, writes
	// performed by the new leaseholder are guaranteed to have higher timestamps
	// than any reads served by the old leaseholder. However, a read at timestamp
	// T needs to observe all writes at timestamps [T, T+maxOffset] and so,
	// without the stasis, only the new leaseholder might have some of these
	// writes. In other words, without the stasis, a new leaseholder with a fast
	// clock could start performing writes ordered in real time before the old
	// leaseholder considers its lease to have expired.
	LeaseState_STASIS LeaseState = 2
	// EXPIRED indicates that the lease can't be used. An expired lease
	// may become VALID for the same leaseholder on RequestLease or
	// Heartbeat, or it may be replaced by a new leaseholder with a
	// RequestLease (for expiration-based leases) or
	// IncrementEpoch+RequestLease (for epoch-based leases).
	LeaseState_EXPIRED LeaseState = 3
	// PROSCRIBED indicates that the lease's proposed timestamp is
	// earlier than allowed. This is used to detect node restarts: a
	// node that has restarted will see its former incarnation's leases
	// as PROSCRIBED so it will renew them before using them. Note that
	// the PROSCRIBED state is only visible to the leaseholder; other
	// nodes will see this as a VALID lease.
	LeaseState_PROSCRIBED LeaseState = 4
)

var LeaseState_name = map[int32]string{
	0: "ERROR",
	1: "VALID",
	2: "STASIS",
	3: "EXPIRED",
	4: "PROSCRIBED",
}
var LeaseState_value = map[string]int32{
	"ERROR":      0,
	"VALID":      1,
	"STASIS":     2,
	"EXPIRED":    3,
	"PROSCRIBED": 4,
}

func (x LeaseState) String() string {
	return proto.EnumName(LeaseState_name, int32(x))
}
func (LeaseState) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_lease_status_cf4e3e39987effff, []int{0}
}

// LeaseStatus holds the lease state, the timestamp at which the state
// is accurate, the lease and optionally the liveness if the lease is
// epoch-based.
type LeaseStatus struct {
	// Lease which this status describes.
	Lease roachpb.Lease `protobuf:"bytes,1,opt,name=lease,proto3" json:"lease"`
	// Timestamp that the lease was evaluated at.
	Timestamp hlc.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp"`
	// State of the lease at timestamp.
	State LeaseState `protobuf:"varint,3,opt,name=state,proto3,enum=cockroach.kv.kvserver.storagepb.LeaseState" json:"state,omitempty"`
	// Liveness if this is an epoch-based lease.
	Liveness Liveness `protobuf:"bytes,4,opt,name=liveness,proto3" json:"liveness"`
}

func (m *LeaseStatus) Reset()         { *m = LeaseStatus{} }
func (m *LeaseStatus) String() string { return proto.CompactTextString(m) }
func (*LeaseStatus) ProtoMessage()    {}
func (*LeaseStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_lease_status_cf4e3e39987effff, []int{0}
}
func (m *LeaseStatus) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LeaseStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (dst *LeaseStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LeaseStatus.Merge(dst, src)
}
func (m *LeaseStatus) XXX_Size() int {
	return m.Size()
}
func (m *LeaseStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_LeaseStatus.DiscardUnknown(m)
}

var xxx_messageInfo_LeaseStatus proto.InternalMessageInfo

func init() {
	proto.RegisterType((*LeaseStatus)(nil), "cockroach.kv.kvserver.storagepb.LeaseStatus")
	proto.RegisterEnum("cockroach.kv.kvserver.storagepb.LeaseState", LeaseState_name, LeaseState_value)
}
func (m *LeaseStatus) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LeaseStatus) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintLeaseStatus(dAtA, i, uint64(m.Lease.Size()))
	n1, err := m.Lease.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	dAtA[i] = 0x12
	i++
	i = encodeVarintLeaseStatus(dAtA, i, uint64(m.Timestamp.Size()))
	n2, err := m.Timestamp.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n2
	if m.State != 0 {
		dAtA[i] = 0x18
		i++
		i = encodeVarintLeaseStatus(dAtA, i, uint64(m.State))
	}
	dAtA[i] = 0x22
	i++
	i = encodeVarintLeaseStatus(dAtA, i, uint64(m.Liveness.Size()))
	n3, err := m.Liveness.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n3
	return i, nil
}

func encodeVarintLeaseStatus(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *LeaseStatus) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Lease.Size()
	n += 1 + l + sovLeaseStatus(uint64(l))
	l = m.Timestamp.Size()
	n += 1 + l + sovLeaseStatus(uint64(l))
	if m.State != 0 {
		n += 1 + sovLeaseStatus(uint64(m.State))
	}
	l = m.Liveness.Size()
	n += 1 + l + sovLeaseStatus(uint64(l))
	return n
}

func sovLeaseStatus(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozLeaseStatus(x uint64) (n int) {
	return sovLeaseStatus(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *LeaseStatus) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLeaseStatus
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: LeaseStatus: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LeaseStatus: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Lease", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLeaseStatus
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthLeaseStatus
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Lease.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLeaseStatus
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthLeaseStatus
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Timestamp.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
			}
			m.State = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLeaseStatus
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.State |= (LeaseState(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Liveness", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLeaseStatus
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthLeaseStatus
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Liveness.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLeaseStatus(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthLeaseStatus
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
func skipLeaseStatus(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowLeaseStatus
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
					return 0, ErrIntOverflowLeaseStatus
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowLeaseStatus
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthLeaseStatus
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowLeaseStatus
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipLeaseStatus(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthLeaseStatus = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowLeaseStatus   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("kv/kvserver/storagepb/lease_status.proto", fileDescriptor_lease_status_cf4e3e39987effff)
}

var fileDescriptor_lease_status_cf4e3e39987effff = []byte{
	// 365 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x91, 0xcf, 0x6a, 0xea, 0x40,
	0x14, 0x87, 0x33, 0xfe, 0xbb, 0xd7, 0x11, 0x24, 0x0c, 0x5d, 0x04, 0xa1, 0xa3, 0x94, 0x2e, 0x6c,
	0x85, 0x09, 0xd8, 0xbe, 0x40, 0xac, 0x59, 0x04, 0x05, 0x65, 0x22, 0xa5, 0x74, 0x53, 0x62, 0x3a,
	0xa8, 0x24, 0x36, 0x21, 0x33, 0xe6, 0x39, 0xba, 0xea, 0x33, 0xb9, 0x74, 0xe9, 0xaa, 0xb4, 0xf1,
	0x45, 0x4a, 0x26, 0x31, 0x76, 0x53, 0xdc, 0x1d, 0xe6, 0x9c, 0xef, 0x77, 0x3e, 0xce, 0xc0, 0xae,
	0x17, 0xeb, 0x5e, 0xcc, 0x59, 0x14, 0xb3, 0x48, 0xe7, 0x22, 0x88, 0x9c, 0x05, 0x0b, 0xe7, 0xba,
	0xcf, 0x1c, 0xce, 0x5e, 0xb8, 0x70, 0xc4, 0x86, 0x93, 0x30, 0x0a, 0x44, 0x80, 0xda, 0x6e, 0xe0,
	0x7a, 0x51, 0xe0, 0xb8, 0x4b, 0xe2, 0xc5, 0xe4, 0xc8, 0x90, 0x82, 0x69, 0x21, 0xd9, 0x0c, 0xe7,
	0xfa, 0xab, 0x23, 0x9c, 0x0c, 0x6a, 0x5d, 0xff, 0x11, 0xbf, 0x8a, 0xd9, 0x1b, 0xe3, 0x79, 0x74,
	0x4b, 0xdb, 0x88, 0x95, 0xaf, 0x2f, 0x7d, 0x57, 0x17, 0xab, 0x35, 0xe3, 0xc2, 0x59, 0x87, 0x79,
	0xe7, 0x62, 0x11, 0x2c, 0x02, 0x59, 0xea, 0x69, 0x95, 0xbd, 0x5e, 0x7d, 0x94, 0x60, 0x63, 0x9c,
	0x1a, 0xda, 0x52, 0x10, 0xdd, 0xc3, 0xaa, 0x14, 0xd6, 0x40, 0x07, 0x74, 0x1b, 0x7d, 0x8d, 0x9c,
	0x54, 0x73, 0x27, 0x22, 0xc7, 0x07, 0x95, 0xed, 0x67, 0x5b, 0xa1, 0xd9, 0x30, 0x32, 0x60, 0xbd,
	0x58, 0xa7, 0x95, 0x24, 0x79, 0xf9, 0x8b, 0x4c, 0x9d, 0xc8, 0xd2, 0x77, 0xc9, 0xec, 0x38, 0x94,
	0xe3, 0x27, 0x0a, 0x19, 0xb0, 0x9a, 0xde, 0x88, 0x69, 0xe5, 0x0e, 0xe8, 0x36, 0xfb, 0x3d, 0x72,
	0xe6, 0x46, 0xa4, 0xb0, 0x66, 0x34, 0x23, 0xd1, 0x08, 0xfe, 0x3f, 0x5e, 0x43, 0xab, 0x48, 0x89,
	0x9b, 0xf3, 0x29, 0x39, 0x90, 0x0b, 0x15, 0x01, 0xb7, 0x23, 0x08, 0x4f, 0x1b, 0x50, 0x1d, 0x56,
	0x4d, 0x4a, 0x27, 0x54, 0x55, 0xd2, 0xf2, 0xd1, 0x18, 0x5b, 0x43, 0x15, 0x20, 0x08, 0x6b, 0xf6,
	0xcc, 0xb0, 0x2d, 0x5b, 0x2d, 0xa1, 0x06, 0xfc, 0x67, 0x3e, 0x4d, 0x2d, 0x6a, 0x0e, 0xd5, 0x32,
	0x6a, 0x42, 0x38, 0xa5, 0x13, 0xfb, 0x81, 0x5a, 0x03, 0x73, 0xa8, 0x56, 0x06, 0xbd, 0xed, 0x37,
	0x56, 0xb6, 0x09, 0x06, 0xbb, 0x04, 0x83, 0x7d, 0x82, 0xc1, 0x57, 0x82, 0xc1, 0xfb, 0x01, 0x2b,
	0xbb, 0x03, 0x56, 0xf6, 0x07, 0xac, 0x3c, 0xd7, 0x0b, 0xa5, 0x79, 0x4d, 0xfe, 0xcc, 0xdd, 0x4f,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x57, 0x4d, 0x79, 0x36, 0x50, 0x02, 0x00, 0x00,
}
