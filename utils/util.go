package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"reflect"
)

func Round(v float32) int {
	if v < 0 {
		return int(v - 0.5)
	} else {
		return int(v + 0.4999999)
	}
}

func Sha1Sign(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return hex.EncodeToString(t.Sum(nil))
}

//判断一个数组中是否有重复元素
func IsRepeated(data interface{}) (bool, interface{}) {
	dv := reflect.ValueOf(data)
	dt := reflect.TypeOf(data)
	if dt.Kind() == reflect.Ptr {
		if dv.IsNil() {
			return false, nil
		}
		dv = dv.Elem()
		dt = dv.Type()
	}
	switch dt.Kind() {
	case reflect.Array, reflect.Slice:
		var m map[interface{}]int = make(map[interface{}]int)
		for i := 0; i < dv.Len(); i++ {
			elem := dv.Index(i).Interface()
			if _, exist := m[elem]; exist {
				return true, elem
			}
			m[elem] = 1
		}
	default:
		return false, nil
	}
	return false, nil
}
