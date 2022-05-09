package utils

// 简单封装了json的Marshal功能
import (
	"bytes"
	"encoding/json"
)

// JsonEncode pingpp.JsonEncode(param1)
func JsonEncode(v interface{}) ([]byte, error) {
	return json.Marshal(&v)
}

// JsonDecode 简单封装了json的UnMarshal功能
// Example pingpp.JsonDecode(param1, param2)
// param1：需要转换成结构体的json数据
// param2：转换后数据容器
func JsonDecode(p []byte, v interface{}) error {
	obj := json.NewDecoder(bytes.NewBuffer(p))
	obj.UseNumber()
	return obj.Decode(&v)
}
