package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

// var pbDir = "D:\\WORK\\me\\gameProject\\proto"
var pbDir = ""

var fileProtoMap = make(map[string]MsgProto)
var lock = sync.Mutex{}

type FiledProto struct {
	filedType string // 数据类型
	filedName string // 字段名
	index     string // 字段的排序
}

type MsgProto struct {
	pbName    string
	filedList []FiledProto
}

const (
	enum = iota
)

func main2() {
	fmt.Println("start parse proto buffer")
	if pbDir == "" {
		pbDir, _ = os.Getwd()
	}
	loadProtoFiles(pbDir)
	createGoHandlerFile()
	fmt.Println("end parse proto buffer")
}

func createGoHandlerFile() {
	cmdMap := make(map[string]int)
	for key, val := range fileProtoMap {
		if key == "CMD" {
			for _, filed := range val.filedList {
				intVal, err := strconv.Atoi(strings.TrimSpace(filed.index))
				if err != nil {
					fmt.Println(err)
				}
				cmdMap[strings.TrimSpace(filed.filedName)] = intVal
			}
		}
	}

	type handleType map[int]string
	cmdCsHandler := make(handleType)
	cmdScHandler := make(handleType)
	handleData := [2]handleType{cmdScHandler, cmdScHandler}
	for key, _ := range fileProtoMap {
		if key == "" || len(key) <= 2 {
			continue
		}
		if key[:2] == "cs" {
			handler := key[2:]
			cmd := cmdMap[handler]
			if cmd != 0 {
				cmdCsHandler[cmd] = key[2:]
			}
		} else if key[:2] == "sc" {
			handler := key[2:]
			cmd := cmdMap[handler]
			if cmd != 0 {
				cmdScHandler[cmd] = key[2:]
			}
		}
	}

	if len(cmdCsHandler) <= 0 {
		fmt.Println("len(cmdCsHandler) <= 0 fileProtoMap：", fileProtoMap)
		return
	}

	// Relative path to the template file
	relativePath := "MsgCreateFactory.tmpl"

	// Resolve the absolute path using filepath.Join
	absPath := filepath.Join(pbDir, relativePath)
	tmpl, err := template.New("MsgCreateFactory.tmpl").ParseFiles(absPath)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, handleData)
	if err != nil {
		fmt.Println("mapTemplate Execute error :", err)
	}
	//fmt.Println(buf.String())
	outFile := "../src/protobufMsg/MsgCreateFactory.go"
	outFile = filepath.Join(pbDir, outFile)
	fmt.Println("create MsgCreateFactory :", outFile)
	os.Remove(outFile)
	fs, err := os.OpenFile(outFile, os.O_RDWR|os.O_CREATE, 0755)
	defer fs.Close()
	if err != nil {
		fmt.Println(" err:", err)
		return
	}
	fs.Write(buf.Bytes())
	fmt.Println("create MsgCreateFactory.go 文件")
}

func loadProtoFiles(pbDir string) {
	var files []string

	err := filepath.Walk(pbDir, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	i := 0
	var wg sync.WaitGroup

	waitNum := 0
	fileList := make([]string, 0)
	cmdFileName := ""
	for _, file := range files {
		fileStat, err := os.Stat(file)
		i++
		if i == 1 {
			continue
		}
		if err != nil {
			continue
		}
		if fileStat.IsDir() {
			loadProtoFiles(file)
		} else {
			fileType := path.Ext(fileStat.Name())
			if fileType != ".proto" {
				continue
			}
			fileNameOnly := strings.TrimSuffix(fileStat.Name(), fileType)
			waitNum++
			if fileNameOnly == "Cmd" {
				cmdFileName = file
			} else {
				fileList = append(fileList, file)
			}
		}
	}
	wg.Add(waitNum)
	go parseCmdFile(cmdFileName, &wg)
	for _, file := range fileList {
		go parseNormalFile(file, &wg)
	}
	wg.Wait()
}

func parseNormalFile(name string, wg *sync.WaitGroup) {

	fmt.Println("parseNormalFile:", name)
	fs, err := os.OpenFile(name, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("parseCmdFile error :", err, "name:", name)
		return
	}
	defer fs.Close()
	defer wg.Done()
	scanner := bufio.NewScanner(fs)

	// 逐行读取文件内容
	msgProto := MsgProto{}
	begin := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) <= 0 {
			continue
		}
		index := strings.Index(line, "//")
		if index > 0 {
			line = line[:index]
		}
		if len(line) > 7 && line[:7] == "message" {
			ss := strings.Split(line, " ")
			if len(ss) != 2 {
				fmt.Println("parse error message line:", line)
				continue
			}

			msgProto.pbName = ss[1]
			continue
		} else if !begin && len(msgProto.pbName) > 0 && line == "{" {
			begin = true
			continue
		} else if strings.Contains(line, "}") {
			lock.Lock()
			fileProtoMap[msgProto.pbName] = msgProto
			lock.Unlock()
			begin = false
		} else if begin {
			endIndex := strings.Index(line, ";")
			if endIndex < 0 {
				continue
			}
			line = line[:endIndex]
			ss := strings.Split(line, "=")
			if len(ss) != 2 {
				fmt.Println("this line:", line, ", ss.lne != 2 ")
				continue
			}
			leftss := strings.Split(strings.TrimSpace(ss[0]), " ")
			if len(leftss) != 2 {
				continue
			}
			fileType := leftss[0]
			fileName := leftss[1]

			filed := FiledProto{
				filedType: fileType,
				filedName: fileName,
				index:     ss[1],
			}
			msgProto.filedList = append(msgProto.filedList, filed)
		}
	}
}

func parseCmdFile(name string, wg *sync.WaitGroup) {
	fmt.Println("parseCmdFile:", name)
	fs, err := os.OpenFile(name, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("parseCmdFile error :", err, "name:", name)
		return
	}
	defer fs.Close()
	defer wg.Done()
	// 创建一个新的 bufio.Scanner 来读取文件内容
	scanner := bufio.NewScanner(fs)

	// 逐行读取文件内容
	msgProto := MsgProto{}
	begin := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) <= 0 {
			continue
		}
		index := strings.Index(line, "//")
		if index > 0 {
			line = line[:index]
		}
		if strings.Contains(line, "enum CMD") {
			msgProto.pbName = "CMD"
			continue
		} else if !begin && len(msgProto.pbName) > 0 && line == "{" {
			begin = true
			continue
		} else if strings.Contains(line, "}") {
			lock.Lock()
			fileProtoMap[msgProto.pbName] = msgProto
			lock.Unlock()
			return
		} else if begin {
			endIndex := strings.Index(line, ";")
			if endIndex < 0 {
				continue
			}
			line = line[:endIndex]
			ss := strings.Split(line, "=")
			if len(ss) != 2 {
				fmt.Println("error line: ", line, " ss.len != 2")
				continue
			}
			/*nameIndex := strings.Index(ss[0], "cmd_")
			if nameIndex < 0 {
				fmt.Println("this cmd not found  prefix cmd_ ")
				os.Exit(1)
			}*/
			filed := FiledProto{
				filedType: "enum",
				filedName: ss[0],
				index:     ss[1],
			}
			msgProto.filedList = append(msgProto.filedList, filed)
		}

	}
}
