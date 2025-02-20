package main

import (
	"fmt"

	"github.com/OblivionTime/go-registry/registry"
)

// 遍历所有项
func ergodic() {
	reg := registry.NewRegistry("SAM")
	rec(reg.Root(), 0)
}
func rec(key *registry.RegistryKey, depth int) {
	// 初始化一个空字符串用于存储制表符
	tabs := ""
	// 循环 depth 次，每次添加一个制表符到 tabs 字符串中
	for i := 0; i < depth; i++ {
		tabs += "\t"
	}
	fmt.Printf("%s%s\n", tabs, key.Path())
	for _, k := range key.Subkeys() {
		rec(k, depth+1)
	}
}

// 查找键并打印所有字符串值
func findKeyAndPrintValues(keyPath string) {
	reg := registry.NewRegistry("SAM")
	key := reg.Open(keyPath)
	if key == nil {
		fmt.Printf("未找到键: %s\n", keyPath)
		return
	}
	values := key.Values()
	for _, value := range values {
		fmt.Printf("key的值为:%s\t类型为:%s\t值为:%v\n", value.Name(), value.Value_type(), value.Value(0))
	}
}

// 获取某个键的值
func getKeyAndPrintValues(keyPath string) {
	reg := registry.NewRegistry("SAM")
	key := reg.Open(keyPath)
	if key == nil {
		fmt.Printf("未找到键: %s\n", keyPath)
		return
	}
	fmt.Println(key.GetStringValue("V"))
	fmt.Println(key.GetBinaryValue("V"))
}

func main() {
	//遍历所有项
	// ergodic()
	// 查找键并打印所有字符串值
	// findKeyAndPrintValues("SAM\\Domains\\Account\\Users\\000003E9")
	// 获取某个键的值
	getKeyAndPrintValues("SAM\\Domains\\Account\\Users\\000003E9")
}
