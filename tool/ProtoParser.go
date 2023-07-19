package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const pbDir = "."

var pbFiles = make([]string, 10)
var fileProtoMap = make(map[string]MsgProto)

type FiledProto struct {
	filedType string
	filedName string
}

type MsgProto struct {
	pbName    string
	filedList []FiledProto
}

func main() {
	fmt.Println("start parse proto buffer")
	loadProtoFiles(pbDir)
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
	for _, file := range files {
		fileStat, err := os.Stat(file)
		if err != nil {
			continue
		}
		if fileStat.IsDir() {
			loadProtoFiles(file)
		} else {
			fileType := path.Ext(fileStat.Name())
			if fileType != "proto" {
				continue
			}
			fileNameOnly := strings.TrimSuffix(fileType, fileType)
			if fileNameOnly == "cmd" {
				go parseCmdFile(fileStat.Name())
			} else {
				go parseNormalFile(fileStat.Name())
			}
		}
		fmt.Println(file)
	}
}

func parseNormalFile(name string) {

}

func parseCmdFile(name string) {

}
