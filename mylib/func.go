package mylib

import (
	"io/ioutil"
	"log"
	"reflect"
)

// CheckError 函数
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

// ArrayContains sliceの中に特定の文字列が含まれるかを返す
func ArrayContains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}

// MapContains mapの中に特定の文字列が含まれるかを返す
func MapContains(arr map[string]string, str string) bool {
	if _, ok := arr[str]; ok {
		return true
	}
	return false
}

// CreateJSONfile 函数
func CreateJSONfile(fp string, s []byte) {
	err := ioutil.WriteFile(fp, s, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// StructToMap 函数 struct 转换成 map
func StructToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	elem := reflect.ValueOf(data).Elem()
	size := elem.NumField()

	for i := 0; i < size; i++ {
		field := elem.Type().Field(i).Name
		value := elem.Field(i).Interface()
		result[field] = value
	}

	return result
}
