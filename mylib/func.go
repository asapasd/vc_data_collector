package mylib

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"time"
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

// DoLog 函数
func DoLog(content string) {
	logfile, err := os.OpenFile("./log/"+time.Now().Format("20060102")+".txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic("cannnot open log file:" + err.Error())
	}
	defer logfile.Close()

	// io.MultiWriteで、
	// 標準出力とファイルの両方を束ねて、
	// logの出力先に設定する
	//log.SetOutput(io.MultiWriter(logfile, os.Stdout))

	// Only logfileに出力
	log.SetOutput(logfile)
	// microsecond resolution: 01:23:23.123123.  assumes Ltime.
	log.SetFlags(log.Lmicroseconds)
	// 内容を書き込み
	log.Println(content)
}
