package utils

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func DecodeWindows1252(data []byte) string {
	// 创建一个 windows-1252 解码器
	decoder := charmap.Windows1252.NewDecoder()
	// 创建一个 transformer 用于解码操作
	transformer := transform.NewReader(bytes.NewReader(data), decoder)
	// 读取解码后的数据
	decodedData, err := ioutil.ReadAll(transformer)
	if err != nil {
		return ""
	}
	// 将解码后的数据转换为字符串
	return string(decodedData)
}

func DecodeUTF16(data []byte) string {
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	transformer := transform.NewReader(bytes.NewReader(data), decoder)
	decoded, err := io.ReadAll(transformer)
	if err != nil {
		fmt.Printf("解码出错: %v\n", err)
		return ""
	}
	return string(decoded)
}

// DecodeUTF16LE 尝试将字节切片按 UTF-16LE 编码解码
func DecodeUTF16LE(s []byte) string {
	// 处理连续的 \x00\x00
	if bytes.Contains(s, []byte{0x00, 0x00}) {
		index := bytes.Index(s, []byte{0x00, 0x00})
		if index > 2 {
			if s[index-2] != 0x00 {
				s = s[:index+2]
			} else {
				s = s[:index+3]
			}
		}
	}

	// 处理字节切片长度为奇数的情况
	if len(s)%2 != 0 {
		s = append(s, 0x00)
	}

	// 创建 UTF-16 解码器，使用 Little Endian 字节序，忽略 BOM
	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	// 创建一个 transform.Reader，用于将输入的字节数据进行解码
	transformer := transform.NewReader(bytes.NewReader(s), decoder)
	// 读取解码后的数据
	decodedData, err := io.ReadAll(transformer)
	if err != nil {
		return ""
	}

	// 处理字符串中的 \x00
	nullIndex := bytes.IndexByte(decodedData, 0x00)
	if nullIndex != -1 {
		decodedData = decodedData[:nullIndex]
	}

	return string(decodedData)
}

// ReverseSlice 反转切片
func ReverseSlice(slice []string) []string {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

// JoinReversedWithBackslash 用反斜杠连接反转后的切片
func JoinReversedWithBackslash(name []string) string {
	reversedName := ReverseSlice(name)
	return strings.Join(reversedName, "\\")
}

// partition 实现类似 Python 的 partition 方法
func Partition(path, sep string) (string, string, string) {
	// 查找分隔符的位置
	index := strings.Index(path, sep)
	if index == -1 {
		// 如果未找到分隔符，返回原字符串和两个空字符串
		return path, "", ""
	}
	// 分割字符串
	before := path[:index]
	middle := sep
	after := path[index+len(sep):]
	return before, middle, after
}
