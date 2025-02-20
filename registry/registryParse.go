package registry

import (
	"bytes"
	"encoding/binary"
	"math"
	"time"
)

type RegistryBlock struct {
	Buffer []byte
	Offset int
	Parent *RegistryBlock
}

func NewRegistryBlock(buffer []byte, offset int, parent *RegistryBlock) *RegistryBlock {
	return &RegistryBlock{
		Buffer: buffer,
		Offset: offset,
		Parent: parent,
	}
}

// UnpackBinary 从相对偏移量开始提取指定长度的二进制数据
func (u *RegistryBlock) UnpackBinary(offset, length int) []byte {
	start := u.Offset + offset
	end := start + length
	return u.Buffer[start:end]
}

// UnpackWord 从相对偏移量开始提取小端字节序的 2 字节无符号整数
func (u *RegistryBlock) UnpackWord(offset int) uint16 {
	var result uint16
	buf := bytes.NewReader(u.Buffer[u.Offset+offset:])
	binary.Read(buf, binary.LittleEndian, &result)
	return result
}

// UnpackDword 从相对偏移量开始提取小端字节序的 4 字节无符号整数
func (u *RegistryBlock) UnpackDword(offset int) uint32 {
	var result uint32
	buf := bytes.NewReader(u.Buffer[u.Offset+offset:])
	binary.Read(buf, binary.LittleEndian, &result)
	return result
}

// UnpackInt 从相对偏移量开始提取小端字节序的 4 字节有符号整数
func (u *RegistryBlock) UnpackInt(offset int) int32 {
	var result int32
	buf := bytes.NewReader(u.Buffer[u.Offset+offset:])
	binary.Read(buf, binary.LittleEndian, &result)
	return result
}

// UnpackQword 从相对偏移量开始提取小端字节序的 8 字节无符号整数
func (u *RegistryBlock) UnpackQword(offset int) uint64 {
	var result uint64
	buf := bytes.NewReader(u.Buffer[u.Offset+offset:])
	binary.Read(buf, binary.LittleEndian, &result)
	return result
}

// UnpackString 从相对偏移量开始提取指定长度的字节字符串
func (u *RegistryBlock) UnpackString(offset, length int) []byte {
	start := u.Offset + offset
	end := start + length
	return u.Buffer[start:end]
}

// AbsoluteOffset 计算相对偏移量对应的绝对偏移量
func (u *RegistryBlock) AbsoluteOffset(offset int) int {
	return u.Offset + offset
}

type REGFBlock struct {
	RegistryBlock
}

func NewREGFBlock(buffer []byte, offset int, parent *RegistryBlock) *REGFBlock {
	reg := RegistryBlock{
		Buffer: buffer,
		Offset: offset,
		Parent: parent,
	}
	ID := reg.UnpackDword(0)
	if ID != 0x66676572 {
		return nil
	}
	return &REGFBlock{
		RegistryBlock: reg,
	}
}

func (u *REGFBlock) FirstKey() *NKRecord {
	first_hbin := u.Hbins()
	key_offset := first_hbin.AbsoluteOffset(int(u.UnpackDword(0x24)))
	d := NewHBINCell(u.Buffer, key_offset, &first_hbin.RegistryBlock)
	return NewNKRecord(u.Buffer, d.Data_offset(), &first_hbin.RegistryBlock)
}
func (u *REGFBlock) Hbins() *HBINBlock {
	return NewHBINBlock(u.Buffer, u.first_hbin_offset(), &u.RegistryBlock)
}
func (u *REGFBlock) first_hbin_offset() int {
	return 0x1000
}

type HBINBlock struct {
	RegistryBlock
	reloffset_next_hbin uint32
	offset_next_hbin    uint32
}

func NewHBINBlock(buffer []byte, offset int, parent *RegistryBlock) *HBINBlock {
	reg := RegistryBlock{
		Buffer: buffer,
		Offset: offset,
		Parent: parent,
	}
	ID := reg.UnpackDword(0)
	if ID != 0x6E696268 {
		panic("not hbin")
	}
	reloffset_next_hbin := reg.UnpackDword(0x8)
	offset_next_hbin := reloffset_next_hbin + uint32(offset)
	return &HBINBlock{
		RegistryBlock:       reg,
		reloffset_next_hbin: reloffset_next_hbin,
		offset_next_hbin:    offset_next_hbin,
	}

}
func (u *HBINBlock) First_hbin() *HBINBlock {
	reloffset_from_first_hbin := u.UnpackQword(0x4)
	return NewHBINBlock(u.Buffer, u.Offset-int(reloffset_from_first_hbin), u.Parent)
}

type HBINCell struct {
	RegistryBlock
	size int32
}

func NewHBINCell(buffer []byte, offset int, parent *RegistryBlock) *HBINCell {
	reg := RegistryBlock{
		Buffer: buffer,
		Offset: offset,
		Parent: parent,
	}
	size := reg.UnpackInt(0)
	return &HBINCell{
		RegistryBlock: reg,
		size:          size,
	}

}
func (u *HBINCell) Data_offset() int {
	return u.Offset + 0x4
}
func (u *HBINCell) Data_id() []byte {
	return u.UnpackString(0x4, 2)
}
func (u *HBINCell) abs_offset_from_hbin_offset(offset uint32) int {
	h := u.Parent
	for h.Parent != nil && h.Parent.Offset != 0 {
		h = h.Parent
	}
	reloffset_from_first_hbin := h.UnpackDword(0x4)
	return NewHBINBlock(u.Buffer, h.Offset-int(reloffset_from_first_hbin), h.Parent).Offset + int(offset)
}
func (u *HBINCell) Raw_data() []byte {
	return u.Buffer[u.Data_offset() : u.Data_offset()+int(u.size)]
}
func (u *HBINCell) Child() *NKRecord {
	if u.size > 0 {
		return nil
	}
	id := u.Data_id()
	switch string(id) {
	case "vk":
		vk := NewVKRecord(u.Buffer, u.Data_offset(), &u.RegistryBlock)
		return &NKRecord{Record: Record{RegistryBlock: vk.RegistryBlock}}
	case "nk":
		nk := NewNKRecord(u.Buffer, u.Data_offset(), &u.RegistryBlock)
		return &NKRecord{Record: Record{RegistryBlock: nk.RegistryBlock}}
	case "lf":
		lf := NewLFRecord(u.Buffer, u.Data_offset(), &u.RegistryBlock)
		return &NKRecord{Record: Record{RegistryBlock: lf.RegistryBlock}}
	case "lh":
		lh := NewLHRecord(u.Buffer, u.Data_offset(), &u.RegistryBlock)
		return &NKRecord{Record: Record{RegistryBlock: lh.RegistryBlock}}
	case "ri":
		ri := NewRIRecord(u.Buffer, u.Data_offset(), &u.RegistryBlock)
		return &NKRecord{Record: Record{RegistryBlock: ri.RegistryBlock}}
	case "li":
		li := NewLIRecord(u.Buffer, u.Data_offset(), &u.RegistryBlock)
		return &NKRecord{Record: Record{RegistryBlock: li.RegistryBlock}}
	case "sk":
		sk := NewSKRecord(u.Buffer, u.Data_offset(), &u.RegistryBlock)
		return &NKRecord{Record: Record{RegistryBlock: sk.RegistryBlock}}
	case "db":
		db := NewDBRecord(u.Buffer, u.Data_offset(), &u.RegistryBlock)
		return &NKRecord{Record: Record{RegistryBlock: db.RegistryBlock}}
	default:
		data := NewDataRecord(u.Buffer, u.Data_offset(), &u.RegistryBlock)
		return &NKRecord{Record: Record{RegistryBlock: data.RegistryBlock}}
	}
}

type ValuesList struct {
	HBINCell
	number uint32
}

func NewValuesList(buffer []byte, offset int, parent *RegistryBlock, number uint32) *ValuesList {
	hc := NewHBINCell(buffer, offset, parent)
	return &ValuesList{
		HBINCell: *hc,
		number:   number,
	}
}
func (u *ValuesList) Values() []*VKRecord {
	result := make([]*VKRecord, 0)
	value_item := 0x0
	for i := 0; i < int(u.number); i++ {
		value_offset := u.abs_offset_from_hbin_offset(u.UnpackDword(value_item))
		d := NewHBINCell(u.Buffer, value_offset, &u.HBINCell.RegistryBlock)
		v := NewVKRecord(u.Buffer, d.Data_offset(), &u.HBINCell.RegistryBlock)
		value_item += 4
		result = append(result, v)
	}
	return result
}

// ParseTimestamp 用于解析时间戳
func ParseTimestamp(ticks int64, resolution int64, epoch time.Time) time.Time {
	// Go 语言的 time 包支持纳秒精度，这里我们将其转换为微秒
	datetimeResolution := int64(1e6)

	// 将自纪元以来的刻度转换为自纪元以来的微秒
	us := int64(math.Round(float64(ticks*datetimeResolution) / float64(resolution)))

	// 转换为 time.Time
	return epoch.Add(time.Duration(us) * time.Microsecond)
}

// ParseWindowsTimestamp 解析 Windows 时间戳
func ParseWindowsTimestamp(qword int64) time.Time {
	// 定义 Windows 时间戳的纪元
	epoch := time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC)
	// 每秒钟的刻度数
	resolution := int64(1e7)
	return ParseTimestamp(qword, resolution, epoch)
}
