package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// unpackUint32BigEndian 函数用于从字节切片中按大端字节序解包出一个无符号 32 位整数
func UnpackUint32BigEndian(d []byte) uint32 {
	// 检查字节切片长度是否足够解包为无符号 32 位整数（需要 4 字节）
	if len(d) < 4 {
		return 0
	}
	// 创建一个字节读取器，从字节切片的前 4 个字节读取数据
	buffer := bytes.NewReader(d[:4])
	var result uint32
	// 使用 binary.Read 函数按大端字节序从字节读取器中读取一个 uint32 类型的值到 result 变量
	err := binary.Read(buffer, binary.BigEndian, &result)
	if err != nil {
		return 0
	}
	return result
}
func UnpackUint32LittleEndian(d []byte) uint32 {
	// 检查字节切片的长度是否足够
	if len(d) < 4 {
		fmt.Println("输入的字节切片长度不足 4 字节，无法解包为 uint32")
		return 0
	}
	// 创建一个字节读取器
	buffer := bytes.NewReader(d[:4])
	var result uint32
	// 使用 binary.Read 函数按照小端字节序读取一个 uint32 类型的值
	err := binary.Read(buffer, binary.LittleEndian, &result)
	if err != nil {
		fmt.Println("eerr", err)
		return 0
	}
	return result
}
func UnpackUint64LittleEndian(d []byte) uint64 {
	// 检查字节切片长度是否足够解包为无符号 64 位整数（需要 8 字节）
	if len(d) < 8 {
		return 0
	}
	// 创建一个字节读取器，从字节切片的前 8 个字节读取数据
	buffer := bytes.NewReader(d[:8])
	var result uint64
	// 使用 binary.Read 函数按小端字节序从字节读取器中读取一个 uint64 类型的值到 result 变量
	err := binary.Read(buffer, binary.LittleEndian, &result)
	if err != nil {
		return 0
	}
	return result
}
