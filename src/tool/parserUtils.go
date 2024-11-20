package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ScanOutFileExtCode(file, begin, end string) *string {
	scanStr := ""
	fs, err := os.OpenFile(file, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("ScanOutFileExtCode error :", err, "name:", file)
		return &scanStr
	}
	defer fs.Close()
	scanner := bufio.NewScanner(fs)
	beginFlag, endFlag := false, false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, begin) {
			beginFlag = true
		} else if strings.Contains(line, end) {
			endFlag = true
		} else if beginFlag && !endFlag {
			scanStr += "\r\n" + line
		}
	}
	return &scanStr
}
