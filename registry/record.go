package registry

import (
	"bytes"
	"fmt"
	"slices"
	"strings"

	"github.com/OblivionTime/go-registry/utils"
)

type Record struct {
	RegistryBlock
}

func NewRecord(buffer []byte, offset int, parent *RegistryBlock) *Record {
	reg := RegistryBlock{
		Buffer: buffer,
		Offset: offset,
		Parent: parent,
	}
	return &Record{
		RegistryBlock: reg,
	}
}
func (u *Record) abs_offset_from_hbin_offset(offset uint32) int {
	h := u.Parent
	for h.Parent != nil && h.Parent.Offset != 0 {
		h = h.Parent
	}
	reloffset_from_first_hbin := h.UnpackDword(0x4)
	return NewHBINBlock(u.Buffer, h.Offset-int(reloffset_from_first_hbin), h.Parent).Offset + int(offset)
}
func (u *Record) Large_data(length int) []byte {
	off := u.abs_offset_from_hbin_offset(u.UnpackDword(0x8))
	cell := NewHBINCell(u.Buffer, off, &u.RegistryBlock)
	dbi := NewDBIndirectBlock(u.Buffer, cell.Data_offset(), &u.RegistryBlock)
	return dbi.Large_data(length)
}

type VKRecord struct {
	Record
}

func NewVKRecord(buffer []byte, offset int, parent *RegistryBlock) *VKRecord {
	reg := RegistryBlock{
		Buffer: buffer,
		Offset: offset,
		Parent: parent,
	}
	id := reg.UnpackString(0x0, 2)
	if string(id) != "vk" {
		fmt.Println("not vk")
		return nil
	}
	return &VKRecord{
		Record: Record{
			RegistryBlock: reg,
		},
	}
}
func (u *VKRecord) Data_type_ori() int {
	return u.data_type()
}
func (u *VKRecord) Data_type_str() string {
	data_type := u.data_type()
	switch data_type {
	case RegSZ:
		return "RegSZ"
	case RegExpandSZ:
		return "RegExpandSZ"
	case RegBin:
		return "RegBin"
	case RegDWord:
		return "RegDWord"
	case RegMultiSZ:
		return "RegMultiSZ"
	case RegQWord:
		return "RegQWord"
	case RegNone:
		return "RegNone"
	case RegBigEndian:
		return "RegBigEndian"
	case RegLink:
		return "RegLink"
	case RegResourceList:
		return "RegResourceList"
	case RegFullResourceDescriptor:
		return "RegFullResourceDescriptor"
	case RegResourceRequirementsList:
		return "RegResourceRequirementsList"
	case RegFileTime:
		return "RegFileTime"
	case RegUint8:
		return "RegUint8"
	case RegInt16:
		return "RegInt16"
	case RegUint16:
		return "RegUint16"
	case RegInt32:
		return "RegInt32"
	case RegUint32:
		return "RegUint32"
	case RegInt64:
		return "RegInt64"
	case RegUint64:
		return "RegUint64"
	case RegFloat:
		return "RegFloat"
	case RegDouble:
		return "RegDouble"
	case RegUnicodeChar:
		return "RegUnicodeChar"
	case RegBoolean:
		return "RegBoolean"
	case RegUnicodeString:
		return "RegUnicodeString"
	case RegCompositeValue:
		return "RegCompositeValue"
	case RegDateTimeOffset:
		return "RegDateTimeOffset"
	case RegTimeSpan:
		return "RegTimeSpan"
	case RegGUID:
		return "RegGUID"
	case RegUnk111:
		return "RegUnk111"
	case RegUnk112:
		return "RegUnk112"
	case RegUnk113:
		return "RegUnk113"
	case RegBytesArray:
		return "RegBytesArray"
	case RegInt16Array:
		return "RegInt16Array"
	case RegUint16Array:
		return "RegUint16Array"
	case RegInt32Array:
		return "RegInt32Array"
	case RegUInt32Array:
		return "RegUInt32Array"
	case RegInt64Array:
		return "RegInt64Array"
	case RegUInt64Array:
		return "RegUInt64Array"
	case RegFloatArray:
		return "RegFloatArray"
	case RegDoubleArray:
		return "RegDoubleArray"
	case RegUnicodeCharArray:
		return "RegUnicodeCharArray"
	case RegBooleanArray:
		return "RegBooleanArray"
	case RegUnicodeStringArray:
		return "RegUnicodeStringArray"
	default:
		return fmt.Sprintf("Unknown type: 0x%X", data_type)
	}
}
func (u *VKRecord) data_type() int {
	return int(u.UnpackDword(0xC)) & DEVPROP_MASK_TYPE
}
func (u *VKRecord) raw_data_length() int {
	return int(u.UnpackDword(0x4))
}
func (u *VKRecord) data_offset() int {
	if u.raw_data_length() < 5 || u.raw_data_length() >= 0x80000000 {
		return u.Offset + 0x8

	} else {
		return u.abs_offset_from_hbin_offset(u.UnpackDword(0x8))
	}
}
func (u *VKRecord) Has_name() bool {
	return u.UnpackWord(0x2) != 0
}
func (u *VKRecord) has_ascii_name() bool {
	return u.UnpackWord(0x10)&1 == 1
}
func (u *VKRecord) Name() string {
	if !u.Has_name() {
		return ""
	}
	name_length := u.UnpackWord(0x2)
	unpacked_string := u.UnpackString(0x14, int(name_length))
	if u.has_ascii_name() {
		return utils.DecodeWindows1252(unpacked_string)
	}
	return utils.DecodeUTF16(unpacked_string)
}
func (u *VKRecord) Data(overrun int) interface{} {
	data_type := u.data_type()
	data_length := u.raw_data_length()
	d := u.raw_data(overrun)
	if data_type == RegSZ || data_type == RegExpandSZ {
		if overrun > 0 {
			//decode_utf16le() only returns the first string, but if we explicitly
			//ask for overrun, let's make a best-effort to decode as much as possible.
			return utils.DecodeUTF16(d)
		} else {
			return utils.DecodeUTF16LE(d)
		}
	} else if data_type == RegBin || data_type == RegNone {
		return d
	} else if data_type == RegDWord {
		return utils.UnpackUint32LittleEndian(d)
	} else if data_type == RegMultiSZ {
		return strings.Split(utils.DecodeUTF16(d), "\x00")
	} else if data_type == RegQWord {
		return utils.UnpackUint64LittleEndian(d)
	} else if data_type == RegBigEndian {
		return utils.UnpackUint32BigEndian(d)
	} else if data_type == RegLink || data_type == RegResourceList || data_type == RegFullResourceDescriptor || data_type == RegResourceRequirementsList {
		return d
	} else if slices.Contains(tt, data_type) {
		d = d[0 : len(d)-8]            //remove timestamp from end
		comp_type := data_type & 0xEFF // Apply mask for composite types
		return ParseAppDataCompositeValue(comp_type, d, len(d))
	} else if data_type == RegFileTime {
		return ParseWindowsTimestamp(int64(utils.UnpackUint64LittleEndian(d)))
	} else if data_length < 5 || data_length >= 0x80000000 {
		return utils.UnpackUint32LittleEndian(d)
	} else {
		return nil
	}
}
func (u *VKRecord) raw_data(overrun int) []byte {
	data_type := u.data_type()
	data_length := u.raw_data_length()
	data_offset := u.data_offset()
	var ret []byte
	if data_type == RegSZ || data_type == RegExpandSZ {
		if data_length >= 0x80000000 {
			// data is contained in the data_offset field
			ret = u.Buffer[data_offset : data_offset+0x4]
		} else if 0x3fd8 < data_length && data_length < 0x80000000 {
			d := NewHBINCell(u.Buffer, data_offset, &u.RegistryBlock)
			if bytes.Equal(d.Data_id(), []byte("db")) {
				ret = d.Child().Large_data(data_length + overrun)
			} else {
				ret = d.Raw_data()[:data_length+overrun]
			}
		} else {
			d := NewHBINCell(u.Buffer, data_offset, &u.RegistryBlock)
			data_offset = d.Data_offset()
			ret = u.Buffer[data_offset : data_offset+data_length]
		}
	} else if data_type == RegBin || data_type == RegNone || slices.Contains(tt, data_type) {
		if data_length >= 0x80000000 {
			ret = []byte("")
		} else if 0x3fd8 < data_length && data_length < 0x80000000 {
			d := NewHBINCell(u.Buffer, data_offset, &u.RegistryBlock)
			if bytes.Equal(d.Data_id(), []byte("db")) {
				ret = d.Child().Large_data(data_length + overrun)
			} else {
				ret = d.Raw_data()[:data_length+overrun]
			}
		} else {
			ret = u.Buffer[data_offset+4 : data_offset+4+data_length+overrun]
		}
	} else if data_type == RegDWord {
		ret = u.UnpackBinary(0x8, 0x4)
	} else if data_type == RegMultiSZ {
		if data_length >= 0x80000000 {
			data_length -= 0x80000000
			ret = u.Buffer[data_offset : data_offset+data_length+overrun]
		} else if 0x3fd8 < data_length && data_length < 0x80000000 {
			d := NewHBINCell(u.Buffer, data_offset, &u.RegistryBlock)
			if bytes.Equal(d.Data_id(), []byte("db")) {
				ret = d.Child().Large_data(data_length + overrun)
			} else {
				ret = d.Raw_data()[:data_length+overrun]
			}
		} else {
			ret = u.Buffer[data_offset+4 : data_offset+4+data_length+overrun]
		}
	} else if data_type == RegQWord {
		d := NewHBINCell(u.Buffer, data_offset, &u.RegistryBlock)
		data_offset = d.Data_offset()
		ret = u.Buffer[data_offset : data_offset+0x8]
	} else if data_type == RegBigEndian {
		d := NewHBINCell(u.Buffer, data_offset, &u.RegistryBlock)
		data_offset = d.Data_offset()
		ret = u.Buffer[data_offset : data_offset+4]
	} else if data_type == RegLink || data_type == RegResourceList || data_type == RegFullResourceDescriptor || data_type == RegResourceRequirementsList {
		if data_length >= 0x80000000 {
			data_length -= 0x80000000
			ret = u.Buffer[data_offset : data_offset+data_length]
		} else if 0x3fd8 < data_length && data_length < 0x80000000 {
			d := NewHBINCell(u.Buffer, data_offset, &u.RegistryBlock)
			if bytes.Equal(d.Data_id(), []byte("db")) {
				ret = d.Child().Large_data(data_length)
			} else {
				ret = d.Raw_data()[:data_length]
			}
		} else {
			ret = u.Buffer[data_offset+4 : data_offset+4+data_length]
		}
	} else if data_type == RegFileTime {
		ret = u.Buffer[data_offset+4 : data_offset+4+data_length]
	} else if data_length < 5 || data_length >= 0x80000000 {
		ret = u.UnpackBinary(0x8, 4)
	} else {
		if data_length >= 0x80000000 {
			data_length -= 0x80000000
			ret = u.Buffer[data_offset : data_offset+data_length]
		} else if 0x3fd8 < data_length && data_length < 0x80000000 {
			d := NewHBINCell(u.Buffer, data_offset, &u.RegistryBlock)
			if bytes.Equal(d.Data_id(), []byte("db")) {
				ret = d.Child().Large_data(data_length)
			} else {
				ret = d.Raw_data()[:data_length]
			}
		} else {
			ret = u.Buffer[data_offset+4 : data_offset+4+data_length]
		}
	}
	return ret
}

type NKRecord struct {
	Record
}

func NewNKRecord(buffer []byte, offset int, parent *RegistryBlock) *NKRecord {
	reg := RegistryBlock{
		Buffer: buffer,
		Offset: offset,
		Parent: parent,
	}
	id := reg.UnpackString(0x0, 2)
	if string(id) != "nk" {
		// panic("not nk")
		fmt.Println("not nk")
		return nil
	}
	return &NKRecord{
		Record: Record{
			RegistryBlock: reg,
		},
	}
}

func (u *NKRecord) Subkey_number() uint32 {
	number := u.UnpackDword(0x14)
	if number == 0xFFFFFFFF {
		return 0
	}
	return number
}
func (u *NKRecord) Subkey_List() *NKRecord {
	subkey_list_offset := u.abs_offset_from_hbin_offset(u.UnpackDword(0x1C))
	d := NewHBINCell(u.Buffer, subkey_list_offset, &u.RegistryBlock)
	id := d.Data_id()
	switch string(id) {
	case "lf":
		lf := NewLFRecord(u.Buffer, d.Data_offset(), &u.RegistryBlock)
		return &NKRecord{Record: Record{RegistryBlock: lf.RegistryBlock}}
	case "lh":
		lh := NewLHRecord(u.Buffer, d.Data_offset(), &u.RegistryBlock)
		return &NKRecord{Record: Record{RegistryBlock: lh.RegistryBlock}}
	case "ri":
		ri := NewRIRecord(u.Buffer, d.Data_offset(), &u.RegistryBlock)
		return &NKRecord{Record: Record{RegistryBlock: ri.RegistryBlock}}
	case "li":
		li := NewLIRecord(u.Buffer, d.Data_offset(), &u.RegistryBlock)
		return &NKRecord{Record: Record{RegistryBlock: li.RegistryBlock}}
	}
	return nil
}

func (u *NKRecord) Keys() []*NKRecord {
	key_index := 0x4
	result := make([]*NKRecord, 0)
	for i := 0; i < u._keys_len(); i++ {
		offset := u.UnpackDword(key_index)
		key_offset := u.abs_offset_from_hbin_offset(offset)
		d := NewHBINCell(u.Buffer, key_offset, &u.RegistryBlock)
		result = append(result, NewNKRecord(u.Buffer, d.Data_offset(), &u.RegistryBlock))
		key_index += 8
	}
	return result

}
func (u *NKRecord) _keys_len() int {
	return int(u.UnpackWord(0x2))
}
func (u *NKRecord) has_ascii_name() bool {
	return u.UnpackWord(0x2)&0x0020 > 0
}
func (u *NKRecord) name() string {
	name_length := u.UnpackWord(0x48)
	unpacked_string := u.UnpackString(0x4C, int(name_length))
	if u.has_ascii_name() {
		return utils.DecodeWindows1252(unpacked_string)
	}
	return utils.DecodeUTF16(unpacked_string)
}
func (u *NKRecord) is_root() bool {
	return u.UnpackWord(0x2)&0x0004 > 0
}
func (u *NKRecord) Path() string {
	name := []string{u.name()}
	offsets := []int{u.Offset}
	p := u
	for p.has_parent_key() {
		p = p.parent_key()
		for _, offset := range offsets {
			if p.Offset == offset {
				name = append(name, "[path cycle]")
				break
			}
		}
		name = append(name, p.name())
		offsets = append(offsets, p.Offset)
	}
	return utils.JoinReversedWithBackslash(name)
}
func (u *NKRecord) parent_key() *NKRecord {
	offset := u.abs_offset_from_hbin_offset(u.UnpackDword(0x10))
	d := NewHBINCell(u.Buffer, offset, u.Parent)
	return NewNKRecord(u.Buffer, d.Data_offset(), u.Parent)
}
func (u *NKRecord) has_parent_key() bool {
	if u.is_root() {
		return false
	}
	return u.parent_key() != nil
}
func (u *NKRecord) values_number() uint32 {
	num := u.UnpackDword(0x24)
	if num == 0xFFFFFFFF {
		return 0
	}
	return num
}
func (u *NKRecord) Values_list() *ValuesList {
	if u.values_number() == 0 {
		return nil
	}
	values_list_offset := u.abs_offset_from_hbin_offset(u.UnpackDword(0x28))
	d := NewHBINCell(u.Buffer, values_list_offset, &u.RegistryBlock)
	return NewValuesList(u.Buffer, d.Data_offset(), &u.RegistryBlock, u.values_number())
}

type LFRecord struct {
	Record
}

func NewLFRecord(buffer []byte, offset int, parent *RegistryBlock) *LFRecord {
	reg := RegistryBlock{
		Buffer: buffer,
		Offset: offset,
		Parent: parent,
	}

	return &LFRecord{
		Record: Record{
			RegistryBlock: reg,
		},
	}
}

type LHRecord struct {
	Record
}

func NewLHRecord(buffer []byte, offset int, parent *RegistryBlock) *LHRecord {
	reg := RegistryBlock{
		Buffer: buffer,
		Offset: offset,
		Parent: parent,
	}

	return &LHRecord{
		Record: Record{
			RegistryBlock: reg,
		},
	}
}

type LIRecord struct {
	Record
}

func NewLIRecord(buffer []byte, offset int, parent *RegistryBlock) *LIRecord {
	reg := RegistryBlock{
		Buffer: buffer,
		Offset: offset,
		Parent: parent,
	}

	return &LIRecord{
		Record: Record{
			RegistryBlock: reg,
		},
	}
}

type RIRecord struct {
	Record
}

func NewRIRecord(buffer []byte, offset int, parent *RegistryBlock) *RIRecord {
	reg := RegistryBlock{
		Buffer: buffer,
		Offset: offset,
		Parent: parent,
	}

	return &RIRecord{
		Record: Record{
			RegistryBlock: reg,
		},
	}
}

type SKRecord struct {
	Record
}

func NewSKRecord(buffer []byte, offset int, parent *RegistryBlock) *SKRecord {
	reg := RegistryBlock{
		Buffer: buffer,
		Offset: offset,
		Parent: parent,
	}

	return &SKRecord{
		Record: Record{
			RegistryBlock: reg,
		},
	}
}

type DBRecord struct {
	Record
}

func NewDBRecord(buffer []byte, offset int, parent *RegistryBlock) *DBRecord {
	reg := RegistryBlock{
		Buffer: buffer,
		Offset: offset,
		Parent: parent,
	}

	return &DBRecord{
		Record: Record{
			RegistryBlock: reg,
		},
	}
}

type DataRecord struct {
	Record
}

func NewDataRecord(buffer []byte, offset int, parent *RegistryBlock) *DataRecord {
	reg := RegistryBlock{
		Buffer: buffer,
		Offset: offset,
		Parent: parent,
	}

	return &DataRecord{
		Record: Record{
			RegistryBlock: reg,
		},
	}
}

type DBIndirectBlock struct {
	Record
}

func NewDBIndirectBlock(buffer []byte, offset int, parent *RegistryBlock) *DBIndirectBlock {
	reg := RegistryBlock{
		Buffer: buffer,
		Offset: offset,
		Parent: parent,
	}

	return &DBIndirectBlock{
		Record: Record{
			RegistryBlock: reg,
		},
	}
}
func (u *DBIndirectBlock) Large_data(length int) []byte {
	b := make([]byte, 0)
	count := 0
	for length > 0 {
		off := u.abs_offset_from_hbin_offset(u.UnpackDword(4 * count))
		cell := NewHBINCell(u.Buffer, off, &u.RegistryBlock)
		size := slices.Min([]int{length, int(cell.size)})
		b = append(b, NewHBINCell(u.Buffer, off, &u.RegistryBlock).Raw_data()[:size]...)
		count += 1
		length -= size
	}
	return b
}
