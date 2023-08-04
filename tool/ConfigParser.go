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

var excelDir = "G:\\WORK\\me\\goserver\\excel"
var configMap = make(map[string]ExcelData)

func main() {
	if len(os.Args) > 1 {
		excelDir = os.Args[1]
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
	for _, data := range configMap {
		if data.fileBaseName == "serverConfig" {
			go createServerConfig(&data, &wg)
		} else {
			go createConfig(&data, &wg)
		}
	}
	wg.Wait()
	fmt.Println(123)
}

func createConfig(e *ExcelData, s *sync.WaitGroup) {
	defer s.Done()

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

func createServerConfig(e *ExcelData, s *sync.WaitGroup) {
	defer s.Done()
	cfgList := make([]serverConfig, 0)
	for _, line := range e.recordList {
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
		fileName := fmt.Sprintf("%sConfig_%s_%s.yaml", config.Name, config.ZoneId, config.PartId)
		relativePath := "serverConfigYaml.tmpl"

		absPath := filepath.Join(excelDir, "../tool", relativePath)
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
		fmt.Println(buf.String())
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
	fileName        string
	fileBaseName    string
	fileIndex       int
	writeFlagList   []int
	decList         []string
	recordList      [][]string
	titleServerList []string
	titleClientList []string
	readCellList    []bool
	typeStringList  []string
}

func recoverFunction() {
	if r := recover(); r != nil {
		fmt.Println("Recovered from panic:", r)
	}
}
func readConfig(wg *sync.WaitGroup, fs string, baseFs string, fsIndex int) {
	defer recoverFunction()
	defer wg.Done()
	ed := ExcelData{fileName: fs, fileBaseName: baseFs, fileIndex: fsIndex}
	f, err := excelize.OpenFile(ed.fileName)
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

	for i, row := range rows {
		for j, colCell := range row {
			colCell = strings.ToLower(strings.TrimSpace(colCell))
			if i == 0 { //第一行#	#
				if colCell == "#" {
					ed.readCellList = append(ed.readCellList, true)
				} else {
					ed.readCellList = append(ed.readCellList, false)
				}
			} else if i == 1 { // server_flag	id	name
				if ed.readCellList[j] {
					ed.titleServerList = append(ed.titleServerList, colCell)
				}
			} else if i == 2 { // client_flag	id	name
				if ed.readCellList[j] {
					ed.titleClientList = append(ed.titleClientList, colCell)
				}
			} else if i == 3 { //STRING	INT	STRING
				if ed.readCellList[j] {
					ed.typeStringList = append(ed.typeStringList, colCell)
				}
			} else if i == 4 { //转表标记	服务器id	服务器名称
				if ed.readCellList[j] {
					ed.decList = append(ed.decList, colCell)
				}
			} else if i == 5 { //0	100	100
				if ed.readCellList[j] {
					flag, err := strconv.Atoi(colCell)
					if err != nil {
						panic(fmt.Sprintf("read file:", fs, ", row:", i, ", cell:", j, ", error:", err))
					}
					ed.writeFlagList = append(ed.writeFlagList, flag)
				}
			} else {
				if j == 0 {
					if colCell != "#" {
						continue
					} else {
						list := make([]string, len(ed.decList))
						ed.recordList = append(ed.recordList, list)
					}
				}
				if ed.readCellList[j] {
					if ed.typeStringList[j] == "int" {
						_, err := strconv.Atoi(colCell)
						if err != nil {
							continue
							//panic(fmt.Sprintf("data not is int read file:", fs, ", row:", i, ", cell:", j, ", error:", err))
						}
					}
					list := ed.recordList[i-6]
					list[j] = colCell
				}
			}
		}
	}
	configMap[ed.fileBaseName] = ed
	fmt.Println(configMap)
}
