package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

const pbDir = "G:\\WORK\\me\\goserver\\proto"

var pbFiles = make([]string, 10)
var fileProtoMap = make(map[string]MsgProto)

type FiledProto struct {
	filedType string
	filedName string
	value     string
}

type MsgProto struct {
	pbName    string
	filedList []FiledProto
}

const (
	enum = iota
)

func main() {
	fmt.Println("start parse proto buffer")
	loadProtoFiles(pbDir)
	fmt.Println(fileProtoMap)
	fmt.Println("end parse proto buffer")
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

	wg.Add(2)

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
			if fileNameOnly == "Cmd" {
				go parseCmdFile(file, &wg)
			} else {
				go parseNormalFile(file, &wg)
			}
		}
		fmt.Println(file)
	}
	wg.Wait()
}

func parseNormalFile(name string, wg *sync.WaitGroup) {
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
			fileProtoMap[msgProto.pbName] = msgProto
			begin = false
		} else if begin {
			endIndex := strings.Index(line, ";")
			if endIndex < 0 {
				continue
			}
			line = line[:endIndex]
			ss := strings.Split(line, "=")
			if len(ss) != 2 {
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
				value:     ss[1],
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
			fileProtoMap[msgProto.pbName] = msgProto
			return
		} else if begin {
			endIndex := strings.Index(line, ";")
			if endIndex < 0 {
				continue
			}
			line = line[:endIndex]
			ss := strings.Split(line, "=")
			if len(ss) != 2 {
				continue
			}
			nameIndex := strings.Index(ss[0], "cmd_")
			filed := FiledProto{
				filedType: "enum",
				filedName: ss[0][nameIndex+4:],
				value:     ss[1],
			}
			msgProto.filedList = append(msgProto.filedList, filed)
		}

	}
}
