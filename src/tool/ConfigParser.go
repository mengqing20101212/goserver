package main

import (
	"bytes"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

var excelDir = ""
var configMap = make(map[string]ExcelData)

func main() {
	if len(os.Args) > 1 {
		excelDir = os.Args[1]
	} else {
		excelDir, _ = os.Getwd()
	}
	fmt.Println("baseUrl:", excelDir)
	var files []string

	err := filepath.Walk(excelDir, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		fmt.Println("find excel dir: ", excelDir, ", err:", err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(files))
	for i, fs := range files {
		fileStat, err := os.Stat(fs)
		if err != nil {
			fmt.Println("fs error:", err, ", fs name:", fs)
			wg.Done()
			continue
		}
		fileType := path.Ext(fileStat.Name())
		if fileType != ".xlsx" {
			wg.Done()
			continue
		}
		fileNameOnly := strings.TrimSuffix(fileStat.Name(), fileType)
		go readConfig(&wg, fs, fileNameOnly, i)
	}
	wg.Wait()
	wg = sync.WaitGroup{}
	wg.Add(len(configMap))
	serverList := make([]string, 0)
	for _, data := range configMap {

		if data.FileBaseName == "src" {
			go createServerConfig(data, &wg)
		} else {
			go createConfig(data, &wg)
			serverList = append(serverList, strings.Title(data.FileBaseName))
		}
	}
	wg.Wait()

	createConfigManger(serverList)

}

func createConfigManger(list []string) {
	outFile := filepath.Join(excelDir, "../src/config", "ConfigManger.go")
	relativePath := "ConfigManger.tmpl"
	absPath := filepath.Join(excelDir, relativePath)
	tmpl, err := template.New("ConfigManger.tmpl").ParseFiles(absPath)
	if err != nil {
		panic(fmt.Sprintf("Error parsing template:%s", err))
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, list)
	if err != nil {
		fmt.Println("mapTemplate Execute error :", err)
	}
	os.Remove(outFile)
	fs, err := os.OpenFile(outFile, os.O_CREATE|os.O_RDWR, 755)
	if err != nil {
		panic(fmt.Sprintf("createConfigManger   error :%s", err))
	}
	fs.WriteString(buf.String())
	fs.Close()
}

func createConfig(e ExcelData, s *sync.WaitGroup) {
	defer s.Done()
	createJsonFile(e)
	relativePath := "Config.tmpl"
	absPath := filepath.Join(excelDir, relativePath)
	tmpl, err := template.New("Config.tmpl").Funcs(template.FuncMap{
		"TrimSuffix": strings.TrimSuffix,
		"makeJson":   makeJson,
		"ToLower":    strings.ToLower,
	}).ParseFiles(absPath)
	if err != nil {
		panic(fmt.Sprintf("Error parsing template:%s", err))
	}

	var buf bytes.Buffer
	copyE := e
	copyE.FileBaseNameOld = copyE.FileBaseName
	copyE.FileBaseName = strings.Title(copyE.FileBaseName)
	err = tmpl.Execute(&buf, copyE)
	if err != nil {
		fmt.Println("mapTemplate Execute error :", err)
	}
	outFile := filepath.Join(excelDir, "../src/config", e.FileBaseName+"Config.go")
	strBegin := "//***** 自定义代码区 begin ****"
	strEnd := "//***** 自定义代码区 end ****"
	scanExtCode := *ScanOutFileExtCode(outFile, strBegin[2:], strEnd)
	strBeginFiled := "//***** 自定义代码区 filed begin ****"
	strEndFiled := "//***** 自定义代码区 filed end ****"
	scanExtFiled := *ScanOutFileExtCode(outFile, strBeginFiled[2:], strEndFiled)
	scanExtFiled += "\r\n"
	os.Remove(outFile)
	fs, err := os.OpenFile(outFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println(" err:", err)
		return
	}
	defer fs.Close()
	if len(scanExtCode) < 10 {
		scanExtCode = "func (self *" + strings.Title(e.FileBaseName) + "ConfigSwitch) afterLoad(config* " + strings.Title(e.FileBaseName) + "ConfigPtr) {\n}"

	}
	outStr := buf.String()
	strs := strings.Split(outStr, strBeginFiled)
	outStr = strs[0]
	outStr += "\r\n" + strBeginFiled
	outStr += scanExtFiled
	if len(strs) > 1 {
		outStr += strs[1]
	}

	fs.WriteString(outStr)
	fs.WriteString(strBegin)
	fs.WriteString("\r\n")
	fs.WriteString(scanExtCode)
	fs.WriteString("\r\n")
	fs.WriteString(strEnd)
	fmt.Println("创建config文件:" + outFile)
}

func makeJson(e ExcelData) string {
	str := ""
	for i, title := range e.TitleServerList {
		if i == 0 {
			continue
		}
		if strings.ToLower(e.TypeStringList[i]) == "int" {
			str += fmt.Sprintf("%s    %s  `json:\"%s\"`    //%s \n", strings.Title(title), "int32", strings.Title(title), e.DecList[i])
		} else {

			str += fmt.Sprintf("%s  %s `json:\"%s\"`  //%s \n", strings.Title(title), "string", strings.Title(title), e.DecList[i])
		}
	}
	return str
}

func createJsonFile(e ExcelData) {
	list := make([]string, 0)
	for _, data := range e.RecordList {
		str := "{"
		for j, col := range data {
			if j == 0 {
				continue
			}
			str += "\"" + strings.Title(e.TitleServerList[j]) + "\": " + getVal(e.TypeStringList[j], col) + ","
		}
		str = str[:len(str)-1]
		str += "} "
		list = append(list, str)
	}
	json := strings.Join(list, "\t\r\n")
	outFile := filepath.Join(excelDir, "../conf", e.FileBaseName+".txt")

	os.Remove(outFile)
	fs, err := os.OpenFile(outFile, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(fmt.Sprintf("not create file:%s, err:%s", outFile, err))
	}
	defer fs.Close()
	fs.WriteString(json)
}

func getVal(strType string, col string) string {
	if strings.ToLower(strType) == "int" {
		return strings.TrimSpace(col)
	} else if strings.ToLower(strType) == "string" {
		return "\"" + strings.TrimSpace(col) + "\""
	} else {
		return strings.TrimSpace(col)
	}
}

type serverConfig struct {
	Id            string
	Name          string
	ServerType    string
	Port          string
	GroupId       string
	GroupName     string
	MysqlIp       string
	MysqlPort     string
	MysqlUsername string
	MysqlPassword string
	MysqlDbname   string
	RedisIp       string
	RedisPort     string
	RedisUsername string
	RedisPassword string
	ZoneId        string
	PartId        string
	InputSize     string
	OutPutSize    string
	MaxConnectNum string
}

func createServerConfig(e ExcelData, s *sync.WaitGroup) {
	defer s.Done()
	cfgList := make([]serverConfig, 0)
	for _, line := range e.RecordList {
		cfg := serverConfig{}
		for j, col := range line {
			if j == 1 {
				cfg.Id = col
			} else if j == 2 {
				cfg.Name = col
			} else if j == 3 {
				cfg.ServerType = col
			} else if j == 4 {
				cfg.Port = col
			} else if j == 5 {
				cfg.GroupId = col
			} else if j == 6 {
				cfg.GroupName = col
			} else if j == 7 {
				cfg.MysqlIp = col
			} else if j == 8 {
				cfg.MysqlPort = col
			} else if j == 9 {
				cfg.MysqlUsername = col
			} else if j == 10 {
				cfg.MysqlDbname = col
			} else if j == 11 {
				cfg.MysqlPassword = col
			} else if j == 12 {
				cfg.RedisIp = col
			} else if j == 13 {
				cfg.RedisPort = col
			} else if j == 14 {
				cfg.RedisPassword = col
			} else if j == 15 {
				cfg.RedisUsername = col
			} else if j == 16 {
				cfg.MaxConnectNum = col
			} else if j == 17 {
				cfg.InputSize = getByteSize(col)
			} else if j == 18 {
				cfg.OutPutSize = getByteSize(col)
			} else if j == 19 {
				cfg.ZoneId = col
			} else if j == 20 {
				cfg.PartId = col
			}
		}
		cfgList = append(cfgList, cfg)
	}
	for _, config := range cfgList {
		fileName := fmt.Sprintf("%sConfig_%s_%s_%s.yaml", config.Name, config.ZoneId, config.PartId, config.Port)
		relativePath := "serverConfigYaml.tmpl"

		absPath := filepath.Join(excelDir, "../src/tool", relativePath)
		tmpl, err := template.New("serverConfigYaml.tmpl").ParseFiles(absPath)
		if err != nil {
			fmt.Println("Error parsing template:", err)
			return
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, config)
		if err != nil {
			fmt.Println("mapTemplate Execute error :", err)
		}
		outFile := filepath.Join(excelDir, "../script/config", fileName)
		os.Remove(outFile)
		fs, err := os.OpenFile(outFile, os.O_RDWR|os.O_CREATE, 0755)
		defer fs.Close()
		if err != nil {
			fmt.Println(" err:", err)
			return
		}
		fs.WriteString(buf.String())
	}
}

func getByteSize(col string) string {
	if col == "" {
		return "0"
	}
	val, err := strconv.Atoi(col[:len(col)-1])
	if err != nil {
		panic(fmt.Sprintf("col:%s, err:%s", col, err))
	}
	t := col[len(col)-1:]
	K := 1204
	M := 1024 * K

	if strings.ToLower(t) == "k" {
		val *= K
	} else if strings.ToLower(t) == "k" {
		val *= M
	}
	return strconv.Itoa(val)
}

type ExcelData struct {
	FileName        string
	FileBaseName    string
	FileBaseNameOld string
	FileIndex       int
	WriteFlagList   []int
	DecList         []string
	RecordList      [][]string
	TitleServerList []string
	TitleClientList []string
	ReadCellList    []bool
	TypeStringList  []string
}

func recoverFunction() {
	if r := recover(); r != nil {
		fmt.Println("Recovered from panic:", r)
	}
}
func readConfig(wg *sync.WaitGroup, fs string, baseFs string, fsIndex int) {
	defer recoverFunction()
	defer wg.Done()
	ed := ExcelData{FileName: fs, FileBaseName: baseFs, FileIndex: fsIndex}
	f, err := excelize.OpenFile(ed.FileName)
	if err != nil {
		fmt.Println("open file", fs, ", err:", err)
		return
	}
	defer f.Close()
	// 获取 Sheet1 上所有单元格
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println("read rows sheet1 fail ", fs, ", err:", err)
		return
	}

	skipNum := 6

	for i, row := range rows {
		if row == nil {
			skipNum++
		}
		for j, colCell := range row {
			colCell = strings.ToLower(strings.TrimSpace(colCell))
			if i == 0 { //第一行#	#
				if colCell == "#" {
					ed.ReadCellList = append(ed.ReadCellList, true)
				} else {
					ed.ReadCellList = append(ed.ReadCellList, false)
				}
			} else if i == 1 { // server_flag	id	name
				if ed.ReadCellList[j] {
					ed.TitleServerList = append(ed.TitleServerList, colCell)
				}
			} else if i == 2 { // client_flag	id	name
				if ed.ReadCellList[j] {
					ed.TitleClientList = append(ed.TitleClientList, colCell)
				}
			} else if i == 3 { //STRING	INT	STRING
				if ed.ReadCellList[j] {
					ed.TypeStringList = append(ed.TypeStringList, colCell)
				}
			} else if i == 4 { //转表标记	服务器id	服务器名称
				if ed.ReadCellList[j] {
					ed.DecList = append(ed.DecList, colCell)
				}
			} else if i == 5 { //0	100	100
				if ed.ReadCellList[j] {
					flag, err := strconv.Atoi(colCell)
					if err != nil {
						panic(fmt.Sprintf("read file:", fs, ", row:", i, ", cell:", j, ", error:", err))
					}
					ed.WriteFlagList = append(ed.WriteFlagList, flag)
				}
			} else {
				if j == 0 {
					if colCell != "#" {
						skipNum++
						break
					} else {
						list := make([]string, len(ed.DecList))
						ed.RecordList = append(ed.RecordList, list)
						continue
					}
				}
				if ed.ReadCellList[j] {
					if ed.TypeStringList[j] == "int" {
						if colCell == "" {
							colCell = "0"
						}
						_, err := strconv.Atoi(colCell)
						if err != nil {
							continue
							//panic(fmt.Sprintf("data not is int read file:", fs, ", row:", i, ", cell:", j, ", error:", err))
						}
					}
					list := ed.RecordList[i-skipNum]
					list[j] = colCell
				}
			}
		}
	}
	configMap[ed.FileBaseName] = ed
}
