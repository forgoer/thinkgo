package helper

import (
	"reflect"
	"fmt"
	"strconv"
)

func GetString(value interface{}) (string,error)  {
	if value == nil {
		return "", nil
	}

	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64,reflect.Float32, reflect.Float64:
		return NumericalToString(value), nil
	case reflect.String:
		return value.(string),nil
	default:
		return "", fmt.Errorf("Convert: type %T Cannot be converted to a string value %#v", t, value)
	}
}

func NumericalToString(value interface{}) string {
	var val string

	switch value.(type) {
	default:
		return ""
	case int:
		intVal, _ := value.(int)
		val = strconv.FormatInt(int64(intVal), 10)
	case int8:
		intVal, _ := value.(int8)
		val = strconv.FormatInt(int64(intVal), 10)
	case int16:
		intVal, _ := value.(int16)
		val = strconv.FormatInt(int64(intVal), 10)
	case int32:
		intVal, _ := value.(int32)
		val = strconv.FormatInt(int64(intVal), 10)
	case int64:
		intVal, _ := value.(int64)
		val = strconv.FormatInt(int64(intVal), 10)
	case uint:
		intVal, _ := value.(uint)
		val = strconv.FormatUint(uint64(intVal), 10)
	case uint8:
		intVal, _ := value.(uint8)
		val = strconv.FormatUint(uint64(intVal), 10)
	case uint16:
		intVal, _ := value.(uint16)
		val = strconv.FormatUint(uint64(intVal), 10)
	case uint32:
		intVal, _ := value.(uint32)
		val = strconv.FormatUint(uint64(intVal), 10)
	case uint64:
		intVal, _ := value.(uint64)
		val = strconv.FormatUint(uint64(intVal), 10)
	case float32:
		floatVal, _ := value.(float32)
		val = strconv.FormatFloat(float64(floatVal), 'f', -1, 32)
	case float64:
		floatVal, _ := value.(float64)
		val = strconv.FormatFloat(float64(floatVal), 'f', -1, 64)
	}
	return val
}
