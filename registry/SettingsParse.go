package registry

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/OblivionTime/go-registry/utils"

	"github.com/google/uuid"
)

// ParseAppDataCompositeValue 解析复合数据值
func ParseAppDataCompositeValue(itemType int, data []byte, dataSize int) interface{} {
	var value interface{}
	switch itemType {
	case RegUint8:
		if len(data) < 1 {
			return nil
		}
		value = data[0]
	case RegInt16:
		if len(data) < 2 {
			return nil
		}
		var val int16
		binary.Read(bytes.NewReader(data[:2]), binary.LittleEndian, &val)
		value = val
	case RegUint16:
		if len(data) < 2 {
			return nil
		}
		var val uint16
		binary.Read(bytes.NewReader(data[:2]), binary.LittleEndian, &val)
		value = val
	case RegInt32:
		if len(data) < 4 {
			return nil
		}
		var val int32
		binary.Read(bytes.NewReader(data[:4]), binary.LittleEndian, &val)
		value = val
	case RegUint32:
		if len(data) < 4 {
			return nil
		}
		var val uint32
		binary.Read(bytes.NewReader(data[:4]), binary.LittleEndian, &val)
		value = val
	case RegInt64:
		if len(data) < 8 {
			return nil
		}
		var val int64
		binary.Read(bytes.NewReader(data[:8]), binary.LittleEndian, &val)
		value = val
	case RegUint64:
		if len(data) < 8 {
			return nil
		}
		var val uint64
		binary.Read(bytes.NewReader(data[:8]), binary.LittleEndian, &val)
		value = val
	case RegFloat:
		if len(data) < 4 {
			return nil
		}
		var val float32
		binary.Read(bytes.NewReader(data[:4]), binary.LittleEndian, &val)
		value = val
	case RegDouble:
		if len(data) < 8 {
			return nil
		}
		var val float64
		binary.Read(bytes.NewReader(data[:8]), binary.LittleEndian, &val)
		value = val
	case RegUnicodeChar:
		if len(data) < 2 {
			return nil
		}
		// 这里简单假设 data 是 UTF-16LE 编码
		value = string(data[:2])
	case RegBoolean:
		if len(data) < 1 {
			return nil
		}
		value = data[0] != 0
	case RegUnicodeString:
		// 这里简单假设 data 是 UTF-16LE 编码
		value = string(data)
	case RegCompositeValue:
		// 假设 ParseAppDataCompositeStream 函数已实现
		value = ParseAppDataCompositeStream(data)
	case RegDateTimeOffset:
		if len(data) < 8 {
			return nil
		}
		var val uint64
		binary.Read(bytes.NewReader(data[:8]), binary.LittleEndian, &val)
		value = ParseWindowsTimestamp(int64(val))
	case RegTimeSpan:
		if len(data) < 8 {
			return nil
		}
		var val uint64
		binary.Read(bytes.NewReader(data[:8]), binary.LittleEndian, &val)
		value = time.Duration(val*100) * time.Nanosecond
	case RegGUID:
		// 假设 ReadGuid 函数已实现
		value = ReadGuid(data)
	case RegBytesArray:
		if len(data) < dataSize {
			return nil
		}
		value = data[:dataSize]
	case RegInt16Array:
		if len(data) < dataSize {
			return nil
		}
		count := dataSize / 2
		result := make([]int16, count)
		for i := 0; i < count; i++ {
			binary.Read(bytes.NewReader(data[i*2:(i+1)*2]), binary.LittleEndian, &result[i])
		}
		value = result
	case RegUint16Array:
		if len(data) < dataSize {
			return nil
		}
		count := dataSize / 2
		result := make([]uint16, count)
		for i := 0; i < count; i++ {
			binary.Read(bytes.NewReader(data[i*2:(i+1)*2]), binary.LittleEndian, &result[i])
		}
		value = result
	case RegInt32Array:
		if len(data) < dataSize {
			return nil
		}
		count := dataSize / 4
		result := make([]int32, count)
		for i := 0; i < count; i++ {
			binary.Read(bytes.NewReader(data[i*4:(i+1)*4]), binary.LittleEndian, &result[i])
		}
		value = result
	case RegUInt32Array:
		if len(data) < dataSize {
			return nil
		}
		count := dataSize / 4
		result := make([]uint32, count)
		for i := 0; i < count; i++ {
			binary.Read(bytes.NewReader(data[i*4:(i+1)*4]), binary.LittleEndian, &result[i])
		}
		value = result
	case RegInt64Array:
		if len(data) < dataSize {
			return nil
		}
		count := dataSize / 8
		result := make([]int64, count)
		for i := 0; i < count; i++ {
			binary.Read(bytes.NewReader(data[i*8:(i+1)*8]), binary.LittleEndian, &result[i])
		}
		value = result
	case RegUInt64Array:
		if len(data) < dataSize {
			return nil
		}
		count := dataSize / 8
		result := make([]uint64, count)
		for i := 0; i < count; i++ {
			binary.Read(bytes.NewReader(data[i*8:(i+1)*8]), binary.LittleEndian, &result[i])
		}
		value = result
	case RegFloatArray:
		if len(data) < dataSize {
			return nil
		}
		count := dataSize / 4
		result := make([]float32, count)
		for i := 0; i < count; i++ {
			binary.Read(bytes.NewReader(data[i*4:(i+1)*4]), binary.LittleEndian, &result[i])
		}
		value = result
	case RegDoubleArray:
		if len(data) < dataSize {
			return nil
		}
		count := dataSize / 8
		result := make([]float64, count)
		for i := 0; i < count; i++ {
			binary.Read(bytes.NewReader(data[i*8:(i+1)*8]), binary.LittleEndian, &result[i])
		}
		value = result
	case RegUnicodeCharArray:
		// 这里简单假设 data 是 UTF-16LE 编码
		value = string(data)
	case RegBooleanArray:
		if len(data) < dataSize {
			return nil
		}
		result := make([]bool, dataSize)
		for i := 0; i < dataSize; i++ {
			result[i] = data[i] != 0
		}
		value = result
	case RegUnicodeStringArray:
		// 假设 ReadUnicodeStringArray 函数已实现
		value = ReadUnicodeStringArray(data)
	default:
		fmt.Printf("UNKNOWN TYPE FOUND 0x%X data=%v \nPlease report to developers!\n", itemType, data)
		value = string(data)
	}
	return value
}

// ReadGuid 从字节切片中读取并返回一个 UUID
func ReadGuid(buf []byte) uuid.UUID {
	// 检查字节切片长度是否足够
	if len(buf) < 16 {
		return uuid.Nil
	}
	// 提取前 16 个字节
	guidBytes := buf[:16]

	// 处理小端字节序
	var parts [4]uint32
	binary.Read(bytes.NewReader(guidBytes[:4]), binary.LittleEndian, &parts[0])
	binary.Read(bytes.NewReader(guidBytes[4:6]), binary.LittleEndian, &parts[1])
	binary.Read(bytes.NewReader(guidBytes[6:8]), binary.LittleEndian, &parts[2])

	// 第 3 部分和第 4 部分不需要转换字节序
	var part3, part4 uint16
	binary.Read(bytes.NewReader(guidBytes[8:10]), binary.BigEndian, &part3)
	binary.Read(bytes.NewReader(guidBytes[10:12]), binary.BigEndian, &part4)

	// 创建 UUID
	var result uuid.UUID
	binary.BigEndian.PutUint32(result[:4], parts[0])
	binary.BigEndian.PutUint16(result[4:6], uint16(parts[1]))
	binary.BigEndian.PutUint16(result[6:8], uint16(parts[2]))
	binary.BigEndian.PutUint16(result[8:10], part3)
	binary.BigEndian.PutUint16(result[10:12], part4)
	copy(result[12:], guidBytes[12:16])

	return result
}

// ParseAppDataCompositeStream 解析包含 ApplicationDataCompositeData 二进制对象的缓冲区
func ParseAppDataCompositeStream(buf []byte) map[string]interface{} {
	compositeData := make(map[string]interface{})
	bufLen := len(buf)
	pos := 0
	itemPos := 0

	for pos < bufLen {
		var itemByteLen, itemType, itemNameLen uint32
		// 从 buf 中按小端字节序解包出 itemByteLen, itemType, itemNameLen
		err := binary.Read(bytes.NewReader(buf[pos:pos+12]), binary.LittleEndian, &[]interface{}{&itemByteLen, &itemType, &itemNameLen})
		if err != nil {
			fmt.Printf("解包数据时出错: %v\n", err)
			break
		}
		itemPos = pos
		pos += 12

		// 提取并解码 itemName
		itemNameBytes := buf[pos : pos+int(itemNameLen)*2]
		itemName := string(itemNameBytes) // 这里简单假设是 UTF-16LE 编码，可按需调整
		pos += (int(itemNameLen) + 1) * 2

		// 计算数据大小
		dataSize := int(itemByteLen) - 12 - (int(itemNameLen)+1)*2
		data := buf[pos : pos+dataSize]

		// 调用 ParseAppDataCompositeValue 解析数据
		value := ParseAppDataCompositeValue(int(itemType), data, dataSize)
		compositeData[itemName] = value

		pos = itemPos + int(itemByteLen)
		// 对齐到 8 字节边界
		if pos%8 != 0 {
			pos += 8 - (pos % 8)
		}
	}
	return compositeData
}

// ReadUnicodeStringArray 从字节缓冲区读取 Unicode 字符串数组
func ReadUnicodeStringArray(buf []byte) []string {
	var strings []string
	bufLen := len(buf)
	pos := 0

	for pos < bufLen {
		var itemByteLen uint32
		// 从 buf 中按小端字节序解包出 itemByteLen
		err := binary.Read(bytes.NewReader(buf[pos:pos+4]), binary.LittleEndian, &itemByteLen)
		if err != nil {
			fmt.Printf("解包字符串长度时出错: %v\n", err)
			break
		}
		pos += 4

		// 提取字符串数据
		stringData := buf[pos : pos+int(itemByteLen)]
		// 将 UTF-16 数据转换为 UTF-8 字符串
		var runes []byte
		for i := 0; i < len(stringData); i += 2 {
			r := binary.LittleEndian.Uint16(stringData[i : i+2])
			runes = append(runes, byte(r))
		}
		// 去除末尾的空字符
		for len(runes) > 0 && runes[len(runes)-1] == 0 {
			runes = runes[:len(runes)-1]
		}
		utf8Str := utils.DecodeUTF16(runes)
		strings = append(strings, utf8Str)

		pos += int(itemByteLen)
	}
	return strings
}
