package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

//用map填充结构
func MapToStruct(data map[string]interface{}, obj interface{}) error {
	structType := reflect.TypeOf(obj).Elem()
	structValue := reflect.ValueOf(obj).Elem()
	for i := 0; i < structType.NumField(); i++ {
		if value, ok := data[structType.Field(i).Tag.Get("json")]; ok {
			structFieldValue := structValue.Field(i)
			if !structFieldValue.CanSet() {
				return fmt.Errorf("Cannot set %s field value", data[structType.Field(i).Tag.Get("json")])
			}
			val := reflect.ValueOf(value)
			structFieldType := structFieldValue.Type()

			var err error
			if structFieldType != val.Type() {
				val, err = TypeConversion(fmt.Sprintf("%v", value), structFieldValue.Type().Name())
				if err != nil {
					return err
				}
			}

			structFieldValue.Set(val)
		}
	}
	return nil
}

func SliceMapToSlice(datas []map[string]interface{}, objs []interface{}) {
	//TODO
}

//类型转换
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "uint" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(uint(i)), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	//else if .......增加其他一些类型的转换 (bool,string不需要)

	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}
