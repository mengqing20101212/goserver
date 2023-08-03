package main

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var excelDir = "G:\\WORK\\me\\goserver\\excel"

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
	wg.Add(len(files) * 2)
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

func createServerConfig(e *ExcelData, s *sync.WaitGroup) {
	defer s.Done()
}

var configMap map[string]ExcelData

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

func readConfig(wg *sync.WaitGroup, fs string, baseFs string, fsIndex int) {
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
				if j == 0 {
					if colCell == "#" {
						ed.readCellList = append(ed.readCellList, true)
					} else {
						ed.readCellList = append(ed.readCellList, false)
					}
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
						fmt.Println("read file:", fs, ", row:", i, ", cell:", j, ", error:", err)
						os.Exit(0)
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
							fmt.Println("data not is int read file:", fs, ", row:", i, ", cell:", j, ", error:", err)
							os.Exit(0)
						}
					}
					list := ed.recordList[i-6]
					list[j] = colCell
				}
			}
		}
	}

	configMap[ed.fileBaseName] = ed
}

func readServerConfig(wg *sync.WaitGroup, fs string, fsIndex int) {
	defer wg.Done()
}
