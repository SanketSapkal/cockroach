// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sql/execinfrapb/processors_table_stats.proto

package execinfrapb

/*
	Beware! This package name must not be changed, even though it doesn't match
	the Go package name, because it defines the Protobuf message names which
	can't be changed without breaking backward compatibility.
*/

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import github_com_cockroachdb_cockroach_pkg_sql_catalog_descpb "github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"

import encoding_binary "encoding/binary"

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

type SketchType int32

const (
	// This is the github.com/axiomhq/hyperloglog binary format (as of commit
	// 730eea1) for a sketch with precision 14. Values are encoded using their key
	// encoding, except integers which are encoded in 8 bytes (little-endian).
	SketchType_HLL_PLUS_PLUS_V1 SketchType = 0
)

var SketchType_name = map[int32]string{
	0: "HLL_PLUS_PLUS_V1",
}
var SketchType_value = map[string]int32{
	"HLL_PLUS_PLUS_V1": 0,
}

func (x SketchType) Enum() *SketchType {
	p := new(SketchType)
	*p = x
	return p
}
func (x SketchType) String() string {
	return proto.EnumName(SketchType_name, int32(x))
}
func (x *SketchType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(SketchType_value, data, "SketchType")
	if err != nil {
		return err
	}
	*x = SketchType(value)
	return nil
}
func (SketchType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_processors_table_stats_82c970dea3e562e9, []int{0}
}

// SketchSpec contains the specification for a generated statistic.
type SketchSpec struct {
	SketchType SketchType `protobuf:"varint,1,opt,name=sketch_type,json=sketchType,enum=cockroach.sql.distsqlrun.SketchType" json:"sketch_type"`
	// Each value is an index identifying a column in the input stream.
	// TODO(radu): currently only one column is supported.
	Columns []uint32 `protobuf:"varint,2,rep,name=columns" json:"columns,omitempty"`
	// If set, we generate a histogram for the first column in the sketch.
	GenerateHistogram bool `protobuf:"varint,3,opt,name=generate_histogram,json=generateHistogram" json:"generate_histogram"`
	// Controls the maximum number of buckets in the histogram.
	// Only used by the SampleAggregator.
	HistogramMaxBuckets uint32 `protobuf:"varint,4,opt,name=histogram_max_buckets,json=histogramMaxBuckets" json:"histogram_max_buckets"`
	// Only used by the SampleAggregator.
	StatName string `protobuf:"bytes,5,opt,name=stat_name,json=statName" json:"stat_name"`
}

func (m *SketchSpec) Reset()         { *m = SketchSpec{} }
func (m *SketchSpec) String() string { return proto.CompactTextString(m) }
func (*SketchSpec) ProtoMessage()    {}
func (*SketchSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_processors_table_stats_82c970dea3e562e9, []int{0}
}
func (m *SketchSpec) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SketchSpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (dst *SketchSpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SketchSpec.Merge(dst, src)
}
func (m *SketchSpec) XXX_Size() int {
	return m.Size()
}
func (m *SketchSpec) XXX_DiscardUnknown() {
	xxx_messageInfo_SketchSpec.DiscardUnknown(m)
}

var xxx_messageInfo_SketchSpec proto.InternalMessageInfo

// SamplerSpec is the specification of a "sampler" processor which
// returns a sample (random subset) of the input columns and computes
// cardinality estimation sketches on sets of columns.
//
// The sampler is configured with a sample size and sets of columns
// for the sketches. It produces one row with global statistics, one
// row with sketch information for each sketch plus at most
// sample_size sampled rows.
//
// For each column with an inverted index, a sketch and sample reservoir are
// created. Each of these produces one sketch row and at most sample_size
// sampled rows from the inverted index keys.
//
// The following method is used to do reservoir sampling: we generate a
// "rank" for each row, which is just a random, uniformly distributed
// 64-bit value. The rows with the smallest <sample_size> ranks are selected.
// This method is chosen because it allows to combine sample sets very easily.
//
// The internal schema of the processor is formed of three column groups:
//   1. sampled row columns:
//       - columns that map 1-1 to the columns in the input (same
//         schema as the input). Note that columns unused in a histogram are
//         set to NULL.
//       - an INT column with the "rank" of the row; this is a random value
//         associated with the row (necessary for combining sample sets).
//   2. sketch columns:
//       - an INT column indicating the sketch index
//         (0 to len(sketches) - 1).
//       - an INT column indicating the number of rows processed
//       - an INT column indicating the number of rows with NULL values
//         on all columns of the sketch.
//       - a BYTES column with the binary sketch data (format
//         dependent on the sketch type).
//   3. inverted columns:
//       - an INT column identifying the column index for this inverted sample
//       - a BYTE column of the inverted index key.
//
// There are four row types produced:
//   1. sample rows, using column group #1.
//   2. sketch rows, using column group #2.
//   3. inverted sample rows, using column group #3 and the rank column from #1.
//   4. inverted sketch rows, using column group #2 and first column from #3.
//
// Rows have NULLs on either all the sampled row columns or on all the
// sketch columns.
type SamplerSpec struct {
	Sketches         []SketchSpec `protobuf:"bytes,1,rep,name=sketches" json:"sketches"`
	InvertedSketches []SketchSpec `protobuf:"bytes,4,rep,name=inverted_sketches,json=invertedSketches" json:"inverted_sketches"`
	SampleSize       uint32       `protobuf:"varint,2,opt,name=sample_size,json=sampleSize" json:"sample_size"`
	// Setting this value enables throttling; this is the fraction of time that
	// the sampler processors will be idle when the recent CPU usage is high. The
	// throttling is adaptive so the actual idle fraction will depend on CPU
	// usage; this value is a ceiling.
	//
	// Currently, this field is set only for automatic statistics based on the
	// value of the cluster setting
	// sql.stats.automatic_collection.max_fraction_idle.
	MaxFractionIdle float64 `protobuf:"fixed64,3,opt,name=max_fraction_idle,json=maxFractionIdle" json:"max_fraction_idle"`
}

func (m *SamplerSpec) Reset()         { *m = SamplerSpec{} }
func (m *SamplerSpec) String() string { return proto.CompactTextString(m) }
func (*SamplerSpec) ProtoMessage()    {}
func (*SamplerSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_processors_table_stats_82c970dea3e562e9, []int{1}
}
func (m *SamplerSpec) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SamplerSpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (dst *SamplerSpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SamplerSpec.Merge(dst, src)
}
func (m *SamplerSpec) XXX_Size() int {
	return m.Size()
}
func (m *SamplerSpec) XXX_DiscardUnknown() {
	xxx_messageInfo_SamplerSpec.DiscardUnknown(m)
}

var xxx_messageInfo_SamplerSpec proto.InternalMessageInfo

// SampleAggregatorSpec is the specification of a processor that aggregates the
// results from multiple sampler processors and writes out the statistics to
// system.table_statistics.
//
// The input schema it expects matches the output schema of a sampler spec (see
// the comment for SamplerSpec for all the details):
//  1. sampled row columns:
//    - sampled columns
//    - row rank
//  2. sketch columns:
//    - sketch index
//    - number of rows processed
//    - number of rows encountered with NULL values on all columns of the sketch
//    - binary sketch data
//  3. inverted columns:
//    - column index for inverted sample
//    - sample column
type SampleAggregatorSpec struct {
	Sketches         []SketchSpec `protobuf:"bytes,1,rep,name=sketches" json:"sketches"`
	InvertedSketches []SketchSpec `protobuf:"bytes,8,rep,name=inverted_sketches,json=invertedSketches" json:"inverted_sketches"`
	// The processor merges reservoir sample sets into a single
	// sample set of this size. This must match the sample size
	// used for each Sampler.
	SampleSize uint32 `protobuf:"varint,2,opt,name=sample_size,json=sampleSize" json:"sample_size"`
	// The i-th value indicates the ColumnID of the i-th sampled row column.
	// These are necessary for writing out the statistic data.
	SampledColumnIDs []github_com_cockroachdb_cockroach_pkg_sql_catalog_descpb.ColumnID `protobuf:"varint,3,rep,name=sampled_column_ids,json=sampledColumnIds,casttype=github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb.ColumnID" json:"sampled_column_ids,omitempty"`
	TableID          github_com_cockroachdb_cockroach_pkg_sql_catalog_descpb.ID         `protobuf:"varint,4,opt,name=table_id,json=tableId,casttype=github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb.ID" json:"table_id"`
	// JobID is the id of the CREATE STATISTICS job.
	JobID int64 `protobuf:"varint,6,opt,name=job_id,json=jobId" json:"job_id"`
	// The total number of rows expected in the table based on previous runs of
	// CREATE STATISTICS. Used for progress reporting. If rows expected is 0,
	// reported progress is 0 until the very end.
	RowsExpected uint64 `protobuf:"varint,7,opt,name=rows_expected,json=rowsExpected" json:"rows_expected"`
}

func (m *SampleAggregatorSpec) Reset()         { *m = SampleAggregatorSpec{} }
func (m *SampleAggregatorSpec) String() string { return proto.CompactTextString(m) }
func (*SampleAggregatorSpec) ProtoMessage()    {}
func (*SampleAggregatorSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_processors_table_stats_82c970dea3e562e9, []int{2}
}
func (m *SampleAggregatorSpec) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *SampleAggregatorSpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (dst *SampleAggregatorSpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SampleAggregatorSpec.Merge(dst, src)
}
func (m *SampleAggregatorSpec) XXX_Size() int {
	return m.Size()
}
func (m *SampleAggregatorSpec) XXX_DiscardUnknown() {
	xxx_messageInfo_SampleAggregatorSpec.DiscardUnknown(m)
}

var xxx_messageInfo_SampleAggregatorSpec proto.InternalMessageInfo

func init() {
	proto.RegisterType((*SketchSpec)(nil), "cockroach.sql.distsqlrun.SketchSpec")
	proto.RegisterType((*SamplerSpec)(nil), "cockroach.sql.distsqlrun.SamplerSpec")
	proto.RegisterType((*SampleAggregatorSpec)(nil), "cockroach.sql.distsqlrun.SampleAggregatorSpec")
	proto.RegisterEnum("cockroach.sql.distsqlrun.SketchType", SketchType_name, SketchType_value)
}
func (m *SketchSpec) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SketchSpec) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0x8
	i++
	i = encodeVarintProcessorsTableStats(dAtA, i, uint64(m.SketchType))
	if len(m.Columns) > 0 {
		for _, num := range m.Columns {
			dAtA[i] = 0x10
			i++
			i = encodeVarintProcessorsTableStats(dAtA, i, uint64(num))
		}
	}
	dAtA[i] = 0x18
	i++
	if m.GenerateHistogram {
		dAtA[i] = 1
	} else {
		dAtA[i] = 0
	}
	i++
	dAtA[i] = 0x20
	i++
	i = encodeVarintProcessorsTableStats(dAtA, i, uint64(m.HistogramMaxBuckets))
	dAtA[i] = 0x2a
	i++
	i = encodeVarintProcessorsTableStats(dAtA, i, uint64(len(m.StatName)))
	i += copy(dAtA[i:], m.StatName)
	return i, nil
}

func (m *SamplerSpec) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SamplerSpec) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Sketches) > 0 {
		for _, msg := range m.Sketches {
			dAtA[i] = 0xa
			i++
			i = encodeVarintProcessorsTableStats(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	dAtA[i] = 0x10
	i++
	i = encodeVarintProcessorsTableStats(dAtA, i, uint64(m.SampleSize))
	dAtA[i] = 0x19
	i++
	encoding_binary.LittleEndian.PutUint64(dAtA[i:], uint64(math.Float64bits(float64(m.MaxFractionIdle))))
	i += 8
	if len(m.InvertedSketches) > 0 {
		for _, msg := range m.InvertedSketches {
			dAtA[i] = 0x22
			i++
			i = encodeVarintProcessorsTableStats(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func (m *SampleAggregatorSpec) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *SampleAggregatorSpec) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Sketches) > 0 {
		for _, msg := range m.Sketches {
			dAtA[i] = 0xa
			i++
			i = encodeVarintProcessorsTableStats(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	dAtA[i] = 0x10
	i++
	i = encodeVarintProcessorsTableStats(dAtA, i, uint64(m.SampleSize))
	if len(m.SampledColumnIDs) > 0 {
		for _, num := range m.SampledColumnIDs {
			dAtA[i] = 0x18
			i++
			i = encodeVarintProcessorsTableStats(dAtA, i, uint64(num))
		}
	}
	dAtA[i] = 0x20
	i++
	i = encodeVarintProcessorsTableStats(dAtA, i, uint64(m.TableID))
	dAtA[i] = 0x30
	i++
	i = encodeVarintProcessorsTableStats(dAtA, i, uint64(m.JobID))
	dAtA[i] = 0x38
	i++
	i = encodeVarintProcessorsTableStats(dAtA, i, uint64(m.RowsExpected))
	if len(m.InvertedSketches) > 0 {
		for _, msg := range m.InvertedSketches {
			dAtA[i] = 0x42
			i++
			i = encodeVarintProcessorsTableStats(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func encodeVarintProcessorsTableStats(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *SketchSpec) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	n += 1 + sovProcessorsTableStats(uint64(m.SketchType))
	if len(m.Columns) > 0 {
		for _, e := range m.Columns {
			n += 1 + sovProcessorsTableStats(uint64(e))
		}
	}
	n += 2
	n += 1 + sovProcessorsTableStats(uint64(m.HistogramMaxBuckets))
	l = len(m.StatName)
	n += 1 + l + sovProcessorsTableStats(uint64(l))
	return n
}

func (m *SamplerSpec) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Sketches) > 0 {
		for _, e := range m.Sketches {
			l = e.Size()
			n += 1 + l + sovProcessorsTableStats(uint64(l))
		}
	}
	n += 1 + sovProcessorsTableStats(uint64(m.SampleSize))
	n += 9
	if len(m.InvertedSketches) > 0 {
		for _, e := range m.InvertedSketches {
			l = e.Size()
			n += 1 + l + sovProcessorsTableStats(uint64(l))
		}
	}
	return n
}

func (m *SampleAggregatorSpec) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Sketches) > 0 {
		for _, e := range m.Sketches {
			l = e.Size()
			n += 1 + l + sovProcessorsTableStats(uint64(l))
		}
	}
	n += 1 + sovProcessorsTableStats(uint64(m.SampleSize))
	if len(m.SampledColumnIDs) > 0 {
		for _, e := range m.SampledColumnIDs {
			n += 1 + sovProcessorsTableStats(uint64(e))
		}
	}
	n += 1 + sovProcessorsTableStats(uint64(m.TableID))
	n += 1 + sovProcessorsTableStats(uint64(m.JobID))
	n += 1 + sovProcessorsTableStats(uint64(m.RowsExpected))
	if len(m.InvertedSketches) > 0 {
		for _, e := range m.InvertedSketches {
			l = e.Size()
			n += 1 + l + sovProcessorsTableStats(uint64(l))
		}
	}
	return n
}

func sovProcessorsTableStats(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozProcessorsTableStats(x uint64) (n int) {
	return sovProcessorsTableStats(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *SketchSpec) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProcessorsTableStats
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
			return fmt.Errorf("proto: SketchSpec: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SketchSpec: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SketchType", wireType)
			}
			m.SketchType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessorsTableStats
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SketchType |= (SketchType(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType == 0 {
				var v uint32
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowProcessorsTableStats
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (uint32(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.Columns = append(m.Columns, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowProcessorsTableStats
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthProcessorsTableStats
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.Columns) == 0 {
					m.Columns = make([]uint32, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v uint32
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowProcessorsTableStats
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (uint32(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.Columns = append(m.Columns, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field Columns", wireType)
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GenerateHistogram", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessorsTableStats
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.GenerateHistogram = bool(v != 0)
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field HistogramMaxBuckets", wireType)
			}
			m.HistogramMaxBuckets = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessorsTableStats
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.HistogramMaxBuckets |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StatName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessorsTableStats
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthProcessorsTableStats
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.StatName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProcessorsTableStats(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthProcessorsTableStats
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
func (m *SamplerSpec) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProcessorsTableStats
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
			return fmt.Errorf("proto: SamplerSpec: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SamplerSpec: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sketches", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessorsTableStats
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
				return ErrInvalidLengthProcessorsTableStats
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sketches = append(m.Sketches, SketchSpec{})
			if err := m.Sketches[len(m.Sketches)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SampleSize", wireType)
			}
			m.SampleSize = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessorsTableStats
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SampleSize |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 1 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxFractionIdle", wireType)
			}
			var v uint64
			if (iNdEx + 8) > l {
				return io.ErrUnexpectedEOF
			}
			v = uint64(encoding_binary.LittleEndian.Uint64(dAtA[iNdEx:]))
			iNdEx += 8
			m.MaxFractionIdle = float64(math.Float64frombits(v))
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InvertedSketches", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessorsTableStats
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
				return ErrInvalidLengthProcessorsTableStats
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.InvertedSketches = append(m.InvertedSketches, SketchSpec{})
			if err := m.InvertedSketches[len(m.InvertedSketches)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProcessorsTableStats(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthProcessorsTableStats
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
func (m *SampleAggregatorSpec) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProcessorsTableStats
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
			return fmt.Errorf("proto: SampleAggregatorSpec: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SampleAggregatorSpec: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sketches", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessorsTableStats
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
				return ErrInvalidLengthProcessorsTableStats
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Sketches = append(m.Sketches, SketchSpec{})
			if err := m.Sketches[len(m.Sketches)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SampleSize", wireType)
			}
			m.SampleSize = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessorsTableStats
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SampleSize |= (uint32(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType == 0 {
				var v github_com_cockroachdb_cockroach_pkg_sql_catalog_descpb.ColumnID
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowProcessorsTableStats
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					v |= (github_com_cockroachdb_cockroach_pkg_sql_catalog_descpb.ColumnID(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				m.SampledColumnIDs = append(m.SampledColumnIDs, v)
			} else if wireType == 2 {
				var packedLen int
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return ErrIntOverflowProcessorsTableStats
					}
					if iNdEx >= l {
						return io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					packedLen |= (int(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				if packedLen < 0 {
					return ErrInvalidLengthProcessorsTableStats
				}
				postIndex := iNdEx + packedLen
				if postIndex > l {
					return io.ErrUnexpectedEOF
				}
				var elementCount int
				var count int
				for _, integer := range dAtA {
					if integer < 128 {
						count++
					}
				}
				elementCount = count
				if elementCount != 0 && len(m.SampledColumnIDs) == 0 {
					m.SampledColumnIDs = make([]github_com_cockroachdb_cockroach_pkg_sql_catalog_descpb.ColumnID, 0, elementCount)
				}
				for iNdEx < postIndex {
					var v github_com_cockroachdb_cockroach_pkg_sql_catalog_descpb.ColumnID
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowProcessorsTableStats
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						v |= (github_com_cockroachdb_cockroach_pkg_sql_catalog_descpb.ColumnID(b) & 0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					m.SampledColumnIDs = append(m.SampledColumnIDs, v)
				}
			} else {
				return fmt.Errorf("proto: wrong wireType = %d for field SampledColumnIDs", wireType)
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TableID", wireType)
			}
			m.TableID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessorsTableStats
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TableID |= (github_com_cockroachdb_cockroach_pkg_sql_catalog_descpb.ID(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field JobID", wireType)
			}
			m.JobID = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessorsTableStats
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.JobID |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RowsExpected", wireType)
			}
			m.RowsExpected = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessorsTableStats
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RowsExpected |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InvertedSketches", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProcessorsTableStats
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
				return ErrInvalidLengthProcessorsTableStats
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.InvertedSketches = append(m.InvertedSketches, SketchSpec{})
			if err := m.InvertedSketches[len(m.InvertedSketches)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProcessorsTableStats(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthProcessorsTableStats
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
func skipProcessorsTableStats(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowProcessorsTableStats
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
					return 0, ErrIntOverflowProcessorsTableStats
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
					return 0, ErrIntOverflowProcessorsTableStats
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
				return 0, ErrInvalidLengthProcessorsTableStats
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowProcessorsTableStats
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
				next, err := skipProcessorsTableStats(dAtA[start:])
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
	ErrInvalidLengthProcessorsTableStats = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowProcessorsTableStats   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("sql/execinfrapb/processors_table_stats.proto", fileDescriptor_processors_table_stats_82c970dea3e562e9)
}

var fileDescriptor_processors_table_stats_82c970dea3e562e9 = []byte{
	// 653 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x54, 0x4b, 0x6f, 0xd3, 0x40,
	0x10, 0x8e, 0xf3, 0x68, 0xd2, 0x0d, 0x81, 0xd4, 0x14, 0xc9, 0xea, 0xc1, 0x31, 0xa1, 0x48, 0x06,
	0x41, 0xcc, 0xe3, 0x82, 0x38, 0x41, 0x08, 0x55, 0x53, 0x5a, 0x84, 0x92, 0xf2, 0x10, 0x17, 0x6b,
	0xbd, 0x3b, 0x75, 0xdc, 0xd8, 0x5e, 0x77, 0x77, 0x03, 0x69, 0xef, 0xc0, 0x95, 0x9f, 0xc1, 0x4f,
	0xe9, 0xb1, 0xc7, 0x9e, 0x22, 0x48, 0xff, 0x45, 0x4f, 0xc8, 0x8f, 0x84, 0x22, 0x40, 0x08, 0x38,
	0x70, 0xb1, 0x76, 0xe7, 0x9b, 0xf9, 0xe6, 0xf3, 0xa7, 0x99, 0x45, 0x37, 0xc4, 0x9e, 0x6f, 0xc1,
	0x18, 0x88, 0x17, 0xee, 0x70, 0x1c, 0x39, 0x56, 0xc4, 0x19, 0x01, 0x21, 0x18, 0x17, 0xb6, 0xc4,
	0x8e, 0x0f, 0xb6, 0x90, 0x58, 0x8a, 0x56, 0xc4, 0x99, 0x64, 0xaa, 0x46, 0x18, 0x19, 0x72, 0x86,
	0xc9, 0xa0, 0x25, 0xf6, 0xfc, 0x16, 0xf5, 0x84, 0x14, 0x7b, 0x3e, 0x1f, 0x85, 0x2b, 0x57, 0x62,
	0x1e, 0x82, 0x25, 0xf6, 0x99, 0x6b, 0x51, 0x10, 0x24, 0x72, 0x2c, 0x21, 0xf9, 0x88, 0xc8, 0x11,
	0x07, 0x9a, 0x96, 0xaf, 0x2c, 0xbb, 0xcc, 0x65, 0xc9, 0xd1, 0x8a, 0x4f, 0x69, 0xb4, 0xf9, 0x2e,
	0x8f, 0x50, 0x7f, 0x08, 0x92, 0x0c, 0xfa, 0x11, 0x10, 0xf5, 0x09, 0xaa, 0x8a, 0xe4, 0x66, 0xcb,
	0xfd, 0x08, 0x34, 0xc5, 0x50, 0xcc, 0xf3, 0x77, 0x56, 0x5b, 0xbf, 0xea, 0xdc, 0x4a, 0x4b, 0xb7,
	0xf7, 0x23, 0x68, 0x17, 0x0f, 0x27, 0x8d, 0x5c, 0x0f, 0x89, 0x79, 0x44, 0xd5, 0x50, 0x99, 0x30,
	0x7f, 0x14, 0x84, 0x42, 0xcb, 0x1b, 0x05, 0xb3, 0xd6, 0x9b, 0x5d, 0xd5, 0xbb, 0x48, 0x75, 0x21,
	0x04, 0x8e, 0x25, 0xd8, 0x03, 0x4f, 0x48, 0xe6, 0x72, 0x1c, 0x68, 0x05, 0x43, 0x31, 0x2b, 0x19,
	0xcf, 0xd2, 0x0c, 0x5f, 0x9f, 0xc1, 0xea, 0x3d, 0x74, 0x69, 0x9e, 0x6b, 0x07, 0x78, 0x6c, 0x3b,
	0x23, 0x32, 0x04, 0x29, 0xb4, 0xa2, 0xa1, 0x98, 0xb5, 0xac, 0xee, 0xe2, 0x3c, 0x65, 0x0b, 0x8f,
	0xdb, 0x69, 0x82, 0x7a, 0x19, 0x2d, 0xc6, 0x46, 0xda, 0x21, 0x0e, 0x40, 0x2b, 0x19, 0x8a, 0xb9,
	0x98, 0x65, 0x57, 0xe2, 0xf0, 0x53, 0x1c, 0x40, 0xf3, 0x43, 0x1e, 0x55, 0xfb, 0x38, 0x88, 0x7c,
	0xe0, 0x89, 0x11, 0x6b, 0xa8, 0x92, 0xfe, 0x09, 0x08, 0x4d, 0x31, 0x0a, 0x66, 0xf5, 0xf7, 0x2e,
	0xc4, 0x75, 0x73, 0xde, 0xac, 0x56, 0xbd, 0x8a, 0xaa, 0x22, 0xa1, 0xb5, 0x85, 0x77, 0x00, 0x5a,
	0xfe, 0x8c, 0x54, 0x94, 0x02, 0x7d, 0xef, 0x00, 0xd4, 0x5b, 0x68, 0x29, 0xfe, 0xa3, 0x1d, 0x8e,
	0x89, 0xf4, 0x58, 0x68, 0x7b, 0xd4, 0x87, 0xc4, 0x0f, 0x25, 0x4b, 0xbe, 0x10, 0xe0, 0xf1, 0x5a,
	0x86, 0x76, 0xa9, 0x0f, 0xea, 0x4b, 0xb4, 0xe4, 0x85, 0x6f, 0x80, 0x4b, 0xa0, 0xf6, 0x5c, 0x69,
	0xf1, 0x8f, 0x95, 0xd6, 0x67, 0x24, 0xfd, 0x8c, 0xa3, 0xf9, 0xa9, 0x88, 0x96, 0x53, 0x27, 0x1e,
	0xba, 0x2e, 0x07, 0x17, 0x4b, 0xf6, 0x5f, 0x2c, 0x79, 0xaf, 0x20, 0x35, 0xbd, 0x52, 0x3b, 0x9d,
	0x1b, 0xdb, 0xa3, 0x42, 0x2b, 0xc4, 0x93, 0xd4, 0x7e, 0x35, 0x9d, 0x34, 0xea, 0xa9, 0x4a, 0xfa,
	0x28, 0x01, 0xbb, 0x1d, 0x71, 0x3a, 0x69, 0x3c, 0x70, 0x3d, 0x39, 0x18, 0x39, 0x2d, 0xc2, 0x02,
	0x6b, 0xae, 0x8d, 0x3a, 0xdf, 0xce, 0x56, 0x34, 0x74, 0xad, 0x1f, 0xd7, 0xa5, 0x35, 0x23, 0xe9,
	0xd5, 0xc5, 0x77, 0xac, 0x54, 0xa8, 0x03, 0x54, 0x49, 0x97, 0xd1, 0xa3, 0xd9, 0xa8, 0x6d, 0xc5,
	0x62, 0xa7, 0x93, 0x46, 0x79, 0x3b, 0x8e, 0x77, 0x3b, 0xa7, 0x93, 0xc6, 0xfd, 0xbf, 0x6d, 0xdc,
	0xed, 0xf4, 0xca, 0x09, 0x7d, 0x97, 0xaa, 0xab, 0x68, 0x61, 0x97, 0x39, 0x71, 0x9f, 0x05, 0x43,
	0x31, 0x0b, 0xed, 0x5a, 0xd6, 0xa7, 0xb4, 0xc1, 0x9c, 0x6e, 0xa7, 0x57, 0xda, 0x65, 0x4e, 0x97,
	0xaa, 0xd7, 0x50, 0x8d, 0xb3, 0xb7, 0xc2, 0x86, 0x71, 0x04, 0x44, 0x02, 0xd5, 0xca, 0x86, 0x62,
	0x16, 0x33, 0x07, 0xcf, 0xc5, 0xd0, 0xe3, 0x0c, 0xf9, 0xf9, 0x90, 0x54, 0xfe, 0x7d, 0x48, 0x36,
	0x8a, 0x95, 0x52, 0x7d, 0xe1, 0x7a, 0x73, 0xf6, 0x76, 0x24, 0xeb, 0xbe, 0x8c, 0xea, 0xeb, 0x9b,
	0x9b, 0xf6, 0xb3, 0xcd, 0xe7, 0xfd, 0xf4, 0xf3, 0xe2, 0x76, 0x3d, 0xd7, 0xbe, 0x79, 0xf8, 0x45,
	0xcf, 0x1d, 0x4e, 0x75, 0xe5, 0x68, 0xaa, 0x2b, 0xc7, 0x53, 0x5d, 0xf9, 0x3c, 0xd5, 0x95, 0x8f,
	0x27, 0x7a, 0xee, 0xe8, 0x44, 0xcf, 0x1d, 0x9f, 0xe8, 0xb9, 0xd7, 0xd5, 0x33, 0x2f, 0xe0, 0xd7,
	0x00, 0x00, 0x00, 0xff, 0xff, 0xf4, 0x92, 0x79, 0x38, 0x13, 0x05, 0x00, 0x00,
}
