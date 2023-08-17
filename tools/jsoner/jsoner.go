package jsoner

import (
	"errors"
	"reflect"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

var jsonConf = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {
	// 解析时容忍空数组作为对象
	extra.RegisterFuzzyDecoders()
}
func Marshal(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, errors.New("invalid memory address or nil pointer dereference")
	}

	return jsonConf.Marshal(v)
}

func MarshalToString(v interface{}) (string, error) {
	byteTmp, err := Marshal(v)
	return Bytes2String(byteTmp), err
}

func Unmarshal(data string, v interface{}) error {
	return jsonConf.Unmarshal(String2Bytes(data), v)
}

func UnmarshalByte(data []byte, v interface{}) error {
	return jsonConf.Unmarshal(data, v)
}

// Bytes2String 高效将bytes转成string, 减少拷贝性能消耗
func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// String2Bytes 高效将string转成bytes, 减少拷贝性能消耗
func String2Bytes(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := *(*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return b
}
