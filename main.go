/*
###################################################
V1.0 2019-10-07 新規作成（呉）
###################################################
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/asapasd/vc_data_collector/mylib"
)

/*#################### 变量以及结构体 ###################*/
// 全局变量
var (
	wg sync.WaitGroup //wg用来等待程序完成
	c  Config         //自定义数据类型
)

// Config 文件读取结构体定义
type Config struct {
	SetProcsNum  int
	SetURL       string
	SetPort      string
	FileDir      string
	ModelCode    string
	SiteCode     string
	FactoryCode  string
	LineCode     []string
	AssyProcess  map[string][]string
	SerialLength int
}

// ReusultDetails JSON数据结构定义
type rJSON struct {
	ModelCode   string              `json:"model_cd"`
	SiteCode    string              `json:"site_cd"`
	FactoryCode string              `json:"factory_cd"`
	LineCode    string              `json:"line_cd"`
	ProcessCode string              `json:"process_cd"`
	SerialCode  string              `json:"serial_cd"`
	DatatypeID  string              `json:"datatype_id"`
	LotCode     string              `json:"lot_cd"`
	MoCode      string              `json:"mo_cd"`
	Inspects    []map[string]string `json:"inspect_items"`
	testAttr    map[string]string   `json:"test_attributes"`
	childInfo   map[string]string   `json:"child_serial_info"`
	jigInfo     map[string]string   `json:"jig_info"`
	macInfo     map[string]string   `json:"machine_info"`
	opInfo      map[string]string   `json:"operator_info"`
	mcInfo      map[string]string   `json:"mcset_info"`
}

// MachineSettings JSON数据结构定义
type mJSON struct {
	ModelCode   string            `json:"model_cd"`
	SiteCode    string            `json:"site_cd"`
	FactoryCode string            `json:"factory_cd"`
	LineCode    string            `json:"line_cd"`
	ProcessCode string            `json:"process_cd"`
	SerialCode  string            `json:"serial_cd"`
	DatatypeID  string            `json:"datatype_id"`
	LotCode     string            `json:"lot_cd"`
	MoCode      string            `json:"mo_cd"`
	Settings    map[string]string `json:"settings"`
	testAttr    map[string]string `json:"test_attributes"`
}

/*#################### 子函数 ###################*/
// checkJSON 子函数，使用具名返回值judge
func checkJSON(c Config, key string, d interface{}) (result string, fname string, message string) {
	var judge = make(map[string]string)

	// interface{}をJSON文字列に変換
	bytes, _ := json.Marshal(d)
	switch key {
	case "sendResultDetails":
		// rJSON结构体实例化
		var jj rJSON
		if err := json.Unmarshal(bytes, &jj); err != nil {
			// JSON解码失败
			judge["result"] = "NG"
			judge["message"] = "sendResultDetails JSON decode fail"
			return judge["result"], "N/A", judge["message"]
		}
		// model_cd不一致
		if jj.ModelCode != c.ModelCode {
			judge["result"] = "NG"
			judge["message"] = "model_cd unmatch"
			return judge["result"], jj.SerialCode, judge["message"]
		}
		// site_cd不一致
		if jj.SiteCode != c.SiteCode {
			judge["result"] = "NG"
			judge["message"] = "site_cd unmatch"
			return judge["result"], jj.SerialCode, judge["message"]
		}
		// factory_cd不一致
		if jj.FactoryCode != c.FactoryCode {
			judge["result"] = "NG"
			judge["message"] = "factory_cd unmatch"
			return judge["result"], jj.SerialCode, judge["message"]
		}
		// line_cd设定不包含
		if !mylib.ArrayContains(c.LineCode, jj.LineCode) {
			judge["result"] = "NG"
			judge["message"] = "line_cd unmatch"
			return judge["result"], jj.SerialCode, judge["message"]
		}
		// datatype_id & process_cd
		if !mylib.ArrayContains(c.AssyProcess[jj.DatatypeID], jj.ProcessCode) {
			judge["result"] = "NG"
			judge["message"] = "datatype_id & process_cd unmatch"
			return judge["result"], jj.SerialCode, judge["message"]
		}
		// serial_cd
		if len(jj.SerialCode) < c.SerialLength {
			judge["result"] = "NG"
			judge["message"] = "serial_cd error"
			return judge["result"], jj.SerialCode, judge["message"]
		}
		// 最终OK
		judge["result"] = "OK"
		judge["message"] = "file created"
		judge["fname"] = key + "_" + jj.ModelCode
		judge["fname"] += jj.SiteCode
		judge["fname"] += jj.FactoryCode
		judge["fname"] += jj.LineCode
		judge["fname"] += jj.ProcessCode + "_"
		judge["fname"] += jj.SerialCode + "_"
		judge["fname"] += time.Now().Format("20060102150405.000000") + ".json"
		return judge["result"], judge["fname"], judge["message"]
	case "sendMachineSettings":
		// mJSON结构体实例化
		var jj mJSON
		if err := json.Unmarshal(bytes, &jj); err != nil {
			// JSON解码失败
			judge["result"] = "NG"
			judge["message"] = "sendMachineSettings JSON decode fail"
			return judge["result"], "N/A", judge["message"]
		}
		// model_cd不一致
		if jj.ModelCode != c.ModelCode {
			judge["result"] = "NG"
			judge["message"] = "model_cd unmatch"
			return judge["result"], jj.SerialCode, judge["message"]
		}
		// site_cd不一致
		if jj.SiteCode != c.SiteCode {
			judge["result"] = "NG"
			judge["message"] = "site_cd unmatch"
			return judge["result"], jj.SerialCode, judge["message"]
		}
		// factory_cd不一致
		if jj.FactoryCode != c.FactoryCode {
			judge["result"] = "NG"
			judge["message"] = "factory_cd unmatch"
			return judge["result"], jj.SerialCode, judge["message"]
		}
		// line_cd设定不包含
		if !mylib.ArrayContains(c.LineCode, jj.LineCode) {
			judge["result"] = "NG"
			judge["message"] = "line_cd unmatch"
			return judge["result"], jj.SerialCode, judge["message"]
		}
		// datatype_id & process_cd
		if !mylib.ArrayContains(c.AssyProcess[jj.DatatypeID], jj.ProcessCode) {
			judge["result"] = "NG"
			judge["message"] = "datatype_id & process_cd unmatch"
			return judge["result"], jj.SerialCode, judge["message"]
		}
		// serial_cd
		if len(jj.SerialCode) < c.SerialLength {
			judge["result"] = "NG"
			judge["message"] = "serial_cd error"
			return judge["result"], jj.SerialCode, judge["message"]
		}
		// 最终OK
		judge["result"] = "OK"
		judge["message"] = "file created"
		judge["fname"] = key + "_" + jj.ModelCode
		judge["fname"] += jj.SiteCode
		judge["fname"] += jj.FactoryCode
		judge["fname"] += jj.LineCode
		judge["fname"] += jj.ProcessCode + "_"
		judge["fname"] += jj.SerialCode + "_"
		judge["fname"] += time.Now().Format("20060102150405.000000") + ".json"
		return judge["result"], judge["fname"], judge["message"]
	default:
		judge["result"] = "NG"
		judge["message"] = "send JSON type unmatch (sendResultDetails or sendMachineSettings)"
		return judge["result"], "N/A", judge["message"]
	}
}

// helloWorld 子函数
func helloWorld(w http.ResponseWriter, r *http.Request) {
	var content string
	
	//if r.URL.Path != ""
	content = "client=" + r.RemoteAddr + ", time=" + time.Now().Format("2006-01-02 15:04:05") + ", accept=" + r.Header.Get("Accept") + ", method=" + r.Method
	//fmt.Println(content)
	mylib.LogAccess(content)

	// 响应不同类型的内容
	switch r.Header.Get("Accept") {
	case "application/json":
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	default:
		// 501
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
		return
	}

	// 响应不同类型的请求
	if r.Method == "POST" {
		// 获取POST请求发来的数据
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		jsonByteData := []byte(reqBody)
		var j map[string]interface{}
		if err := json.Unmarshal(jsonByteData, &j); err != nil {
			// JSON解码失败
			//log.Fatal(err)
			// 500
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			return
		}
		// JSON解码后，数据确认
		for key, val := range j {
			// Check JSON
			chkResult, chkFname, chkMessage := checkJSON(c, key, val)
			//fmt.Printf("%v=%v, %v=%v\n", key, chkResult, chkFname, chkMessage)
			content = `{"status":"` + chkResult + `", "message": "` + chkFname + "=" + chkMessage + `"}`
			if chkResult != "OK" {
				// POST返信
				w.Write([]byte(content))
				mylib.LogRefuse(content)
				return
			}
			// JSON文件生成
			fPath := filepath.Join(c.FileDir, chkFname)
			mylib.CreateJSONfile(fPath, jsonByteData)
			// POST返信
			w.Write([]byte(content))
			mylib.LogJSON(content)
		}
	} else {
		// 405
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
	}
}

/*#################### 主函数 ###################*/
func main() {
	// OS info
	osInfo := runtime.GOOS
	arch := runtime.GOARCH
	fmt.Printf("# OS=%s / architecture=%s\n", osInfo, arch)
	// 获取当前工作目录
	dir, _ := os.Getwd()
	fmt.Println("# Current dir=" + dir)

	// 读取配置文件
	var ConfigFilePath string
	if osInfo == "windows" {
		ConfigFilePath = dir + "\\config.toml"
	}
	if osInfo == "linux" {
		ConfigFilePath = dir + "/config.toml"
	}
	// 解析配置文件，存至全局变量，结构体c
	_, err := toml.DecodeFile(ConfigFilePath, &c)
	if err != nil {
		// エラーメッセージを表示してプログラムを終了
		log.Fatal(err)
	}

	// 设定使用几个逻辑处理器
	MaxProcsCan := runtime.NumCPU()
	MaxProcsSet := c.SetProcsNum
	runtime.GOMAXPROCS(MaxProcsSet)
	fmt.Println("# Use " + strconv.Itoa(MaxProcsSet) + "/" + strconv.Itoa(MaxProcsCan) + " logic processors")

	// JSON文件夹存在确认
	if _, err := os.Stat(c.FileDir); os.IsNotExist(err) {
		// path does not exist
		if err := os.MkdirAll(c.FileDir, 0777); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("# JSON Output dir=" + c.FileDir)

	// LOG文件夹存在确认
	if _, err := os.Stat("./log"); os.IsNotExist(err) {
		// path does not exist
		if err := os.MkdirAll("./log", 0777); err != nil {
			log.Fatal(err)
		}
	}

	// 显示现在的时间
	st := time.Now()
	fmt.Printf("--> Start at %s, port=%s, url=%s\n\n", st.Format("2006-01-02 15:04:05"), c.SetPort, c.SetURL)

	// 启动web svr
	http.HandleFunc(c.SetURL, helloWorld)
	http.ListenAndServe(c.SetPort, nil)
}
