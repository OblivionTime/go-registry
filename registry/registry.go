package registry

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/OblivionTime/go-registry/utils"
)

type Registry struct {
	Buffers []byte
	Regf    *REGFBlock
}

func NewRegistry(filePath string) *Registry {
	buf, _ := os.ReadFile(filePath)
	return &Registry{
		Buffers: buf,
		Regf:    NewREGFBlock(buf, 0, nil),
	}
}

func (r *Registry) Root() *RegistryKey {
	return NewRegistryKey(r.Regf.FirstKey())
}
func (r *Registry) Open(p string) *RegistryKey {
	return r.Root().FindKey(p)
}

type RegistryKey struct {
	Nkrecord *NKRecord
}

func NewRegistryKey(nkrecord *NKRecord) *RegistryKey {
	return &RegistryKey{
		Nkrecord: nkrecord,
	}
}
func (r *RegistryKey) Subkeys() []*RegistryKey {
	result := make([]*RegistryKey, 0)
	if r.Nkrecord.Subkey_number() == 0 {
		return result
	}
	l := r.Nkrecord.Subkey_List()
	for _, v := range l.Keys() {
		result = append(result, NewRegistryKey(v))
	}
	return result

}
func (r *RegistryKey) SubKey(name string) *RegistryKey {
	if r.Nkrecord.Subkey_number() == 0 {
		return nil
	}
	for _, k := range r.Nkrecord.Subkey_List().Keys() {
		if strings.EqualFold(k.name(), name) {
			return NewRegistryKey(k)
		}
	}
	return nil
}
func (r *RegistryKey) Path() string {
	return r.Nkrecord.Path()
}

func (r *RegistryKey) FindKey(p string) *RegistryKey {
	if p == "" {
		return r
	}
	immediate, _, future := utils.Partition(p, "\\")
	return r.SubKey(immediate).FindKey(future)
}

func (r *RegistryKey) Values() []*RegistryValue {
	result := make([]*RegistryValue, 0)
	list := r.Nkrecord.Values_list()
	if list == nil {
		return result
	}
	for _, v := range r.Nkrecord.Values_list().Values() {
		result = append(result, NewRegistryValue(v))
	}
	return result
}
func (r *RegistryKey) Value(name string) *RegistryValue {
	if name == "(default)" {
		name = ""
	}
	for _, v := range r.Nkrecord.Values_list().Values() {
		if strings.EqualFold(v.Name(), name) {
			return NewRegistryValue(v)
		}
	}
	return nil
}
func (r *RegistryKey) GetStringValue(name string) (string, error) {
	if name == "(default)" {
		name = ""
	}
	value := r.Value(name)
	if value == nil {
		return "", errors.New("未找到指定的注册表项")
	}
	data_type := value.Value_type_ori()
	if slices.Contains(stringTypes, data_type) {
		return value.Value(0).(string), nil
	}
	return "", fmt.Errorf("指定的注册表项的类型为:%s,而不是字符串类型", value.Value_type())
}
func (r *RegistryKey) GetBinaryValue(name string) ([]byte, error) {
	if name == "(default)" {
		name = ""
	}
	value := r.Value(name)
	if value == nil {
		return nil, errors.New("未找到指定的注册表项")
	}
	data_type := value.Value_type_ori()
	if slices.Contains(byteArrayTypes, data_type) {
		return value.Value(0).([]byte), nil
	}
	return nil, fmt.Errorf("指定的注册表项的类型为:%s,而不是字节数组类型", value.Value_type())
}
func (r *RegistryKey) GetInt32Value(name string) (uint32, error) {
	if name == "(default)" {
		name = ""
	}
	value := r.Value(name)
	if value == nil {
		return 0, errors.New("未找到指定的注册表项")
	}
	data_type := value.Value_type_ori()
	if slices.Contains(int32Types, data_type) {
		return value.Value(0).(uint32), nil
	}
	return 0, fmt.Errorf("指定的注册表项的类型为:%s,而不是int32类型", value.Value_type())
}
func (r *RegistryKey) GetInt64Value(name string) (uint64, error) {
	if name == "(default)" {
		name = ""
	}
	value := r.Value(name)
	if value == nil {
		return 0, errors.New("未找到指定的注册表项")
	}
	data_type := value.Value_type_ori()
	if slices.Contains(int64Types, data_type) {
		return value.Value(0).(uint64), nil
	}
	return 0, fmt.Errorf("指定的注册表项的类型为:%s,而不是int64类型", value.Value_type())
}

type RegistryValue struct {
	Vkrecord *VKRecord
}

func NewRegistryValue(vkrecord *VKRecord) *RegistryValue {
	return &RegistryValue{
		Vkrecord: vkrecord,
	}
}
func (r *RegistryValue) Name() string {
	if r.Vkrecord.Has_name() {
		return r.Vkrecord.Name()
	} else {
		return "(default)"
	}
}
func (r *RegistryValue) Value_type() string {
	return r.Vkrecord.Data_type_str()
}
func (r *RegistryValue) Value_type_ori() int {
	return r.Vkrecord.Data_type_ori()
}
func (r *RegistryValue) Value(overrun int) interface{} {
	return r.Vkrecord.Data(overrun)
}
