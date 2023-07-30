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

var TablePbDir = "G:\\WORK\\me\\goserver\\proto\\table"

var tablePbFiles = make([]string, 10)
var tableFileProtoMap = make(map[string]TableMsgProto)

type TableFiledProto struct {
	FiledType          string
	FiledName          string
	Value              int
	CamelCaseFiledName string
	IsBaseType         bool
}

type TableMsgProto struct {
	PbName    string
	FiledList []TableFiledProto
	FileName  string
}

func main() {
	fmt.Println("start parse table buffer")
	args := os.Args
	if len(args) > 1 {
		TablePbDir = args[1]
	}
	loadTableProtoFiles(TablePbDir)
	fmt.Println(tableFileProtoMap)
	createGoTableFile()
	createSqlFile()
	fmt.Println("end parse table buffer")
}

func createSqlFile() {

}

func createGoTableFile() {

	// Relative path to the template file
	relativePath := "Table.tmpl"
	// Resolve the absolute path using filepath.Join
	absPath := filepath.Join(TablePbDir, "../", relativePath)
	tmpl, err := template.New("Table.tmpl").Funcs(template.FuncMap{
		"TrimSuffix":     strings.TrimSuffix,
		"scanVal":        ScanVal1,
		"makeScanBytes":  makeScanBytes,
		"protoUnmarshal": protoUnmarshal,
		"saveSqlData":    saveSqlData,
		"createSetGet":   createSetGet,
		"ToLower":        strings.ToLower,
	}).ParseFiles(absPath)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}

	var buf bytes.Buffer
	for fileName, proto := range tableFileProtoMap {
		if fileName != proto.FileName {
			continue
		}
		err = tmpl.Execute(&buf, proto)
		if err != nil {
			fmt.Println("mapTemplate Execute error :", err)
		}
		outFile := "../../table"
		outFile = filepath.Join(TablePbDir, outFile, fileName+"Table.go")
		strBegin := "//*****begin****//"
		strEnd := "//*****end****//"
		scanExtCode := scanOutFileExtCode(outFile, strBegin, strEnd)
		os.Remove(outFile)
		fs, err := os.OpenFile(outFile, os.O_RDWR|os.O_CREATE, 0755)
		defer fs.Close()
		if err != nil {
			fmt.Println(" err:", err)
			return
		}
		fmt.Println("outFile:", outFile, ", scanExtCode: ", scanExtCode)
		fs.Write(buf.Bytes())
		fs.WriteString(strBegin)
		fs.WriteString("\r\n")
		fs.WriteString(scanExtCode)
		fs.WriteString("\r\n")
		fs.WriteString(strEnd)
		fmt.Println(outFile)
	}
}

func ScanVal1(list []TableFiledProto) string {
	result := ""
	for i, filed := range list {
		if IsBaseType(filed.FiledType) {
			if i == 0 {
				result += fmt.Sprintf("&data.%s", filed.CamelCaseFiledName)
			} else {
				result += fmt.Sprintf(", &data.%s", filed.CamelCaseFiledName)
			}
		} else {
			if i == 0 {
				result += fmt.Sprintf("&bs%d", filed.Value)
			} else {
				result += fmt.Sprintf(", &bs%d", filed.Value)
			}
		}
	}
	return result
}
func makeScanBytes(list []TableFiledProto) string {
	result := ""
	for _, filed := range list {
		if !IsBaseType(filed.FiledType) {
			result += fmt.Sprintf("bs%d := make([]byte, 1024)\n", filed.Value)
		}
	}
	result += "\r\n"
	return result
}
func protoUnmarshal(list []TableFiledProto) string {

	result := ""
	for _, filed := range list {
		if !IsBaseType(filed.FiledType) {
			result += "\r\n"
			result += fmt.Sprintf("	   	%s := %s{}\n", filed.FiledName, filed.FiledType)
			result += fmt.Sprintf("		err := proto.Unmarshal(bs%d, &%s)\n", filed.Value, filed.FiledName)
			result += fmt.Sprintf("		if err != nil {\n")
			result += fmt.Sprintf("		logger.Error(fmt.Sprintf(\"Unmarshal bs:%%v, data:%%v\", bs%d, data)) \n", filed.Value)
			result += "		  continue\n"
			result += "		}\n"
			result += fmt.Sprintf("		data.%s =&%s\n", filed.CamelCaseFiledName, filed.FiledName)
		}
	}
	result += "\r\n"
	return result
}
func saveSqlData(data TableMsgProto) string {
	result := ""
	sqlBegin := "INSERT INTO " + strings.ToLower(data.FileName) + " ("
	sqlValues := " values ("
	sqlUpdate := " ON DUPLICATE KEY UPDATE  "
	params := ""
	for i, filed := range data.FiledList {
		if IsBaseType(filed.FiledType) {
			params += fmt.Sprintf(", &data.%s", filed.CamelCaseFiledName)
		} else {
			params += fmt.Sprintf(", &%sBytes", filed.CamelCaseFiledName)
			result += fmt.Sprintf("%sBytes := make([]byte, 1024)\n", filed.CamelCaseFiledName)
			result += fmt.Sprintf("proto.UnmarshalMerge(%sBytes, data.Get%s())\n", filed.CamelCaseFiledName, filed.CamelCaseFiledName)
		}
		if i == 0 {
			sqlBegin += filed.FiledName
			sqlValues += "?"
		} else {
			sqlBegin += ", " + filed.FiledName
			sqlValues += ", ?"
			sqlUpdate += fmt.Sprintf(" %s=?,", filed.FiledName)
		}
	}
	result += "sql := \"" + sqlBegin + ")" + sqlValues + ")" + sqlUpdate[:len(sqlUpdate)-1] + ";\"\n"
	result += "manger := db.GetDataBaseManger()\n"
	result += "if manger == nil || !manger.IsConnectFlag() {\n"
	result += fmt.Sprintf("return false, errors.New(\" %sTableSqlOptional not found DataBaseManger or DataBaseManger not connect\")", data.FileName)
	result += "}\n"
	result += fmt.Sprintf("_, err := manger.GetDB().Exec(sql %s)\n", strings.TrimSuffix(params, ","))
	result += "if err != nil {\n"
	result += fmt.Sprintf("return false, errors.New(fmt.Sprintf(\"save sql error table:%s, sql:%%s, data:%%v\", sql, data))\n", data.FileName)
	result += "}"
	return result
}

//	func (self *AccountTableProxy) GetAccountId() uint64 {
//		return self.AccountId
//	}
//
//	func (self *AccountTableProxy) SetAccountId(accountId uint64) *AccountTableProxy {
//		self.changeFlag = true
//		self.AccountId = accountId
//		return self
//	}
func createSetGet(data TableMsgProto) string {
	result := ""
	for _, filed := range data.FiledList {
		if IsBaseType(filed.FiledType) {
			result += fmt.Sprintf("func (self *%sTableProxy) Get%s() %s {\n", data.FileName, filed.CamelCaseFiledName, filed.FiledType)
			result += fmt.Sprintf("return self.%s\n", filed.CamelCaseFiledName)
			result += "}\r\n"

			result += fmt.Sprintf("func (self *%sTableProxy) Set%s(%s %s) *%sTableProxy {\n", data.FileName, filed.CamelCaseFiledName, filed.FiledName, filed.FiledType, data.FileName)
			result += "self.changeFlag = true\n"
			result += fmt.Sprintf("self.%s = %s\n", filed.CamelCaseFiledName, filed.FiledName)
			result += " return self\n"
			result += "}\n"
		} else {
			result += fmt.Sprintf("func (self *%sTableProxy) Get%s() *%s {\n", data.FileName, filed.CamelCaseFiledName, filed.FiledType)
			result += fmt.Sprintf("return self.%s\n", filed.CamelCaseFiledName)
			result += "}\r\n"

			result += fmt.Sprintf("func (self *%sTableProxy) Set%s(%s* %s) *%sTableProxy {\n", data.FileName, filed.CamelCaseFiledName, filed.FiledName, filed.FiledType, data.FileName)
			result += "self.changeFlag = true\n"
			result += fmt.Sprintf("self.%s = %s\n", filed.CamelCaseFiledName, filed.FiledName)
			result += " return  self\n"
			result += "}\n"
		}
	}
	return result
}

func scanOutFileExtCode(file string, begin string, end string) string {
	fs, err := os.OpenFile(file, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("scanOutFileExtCode error :", err, "name:", file)
		return ""
	}
	defer fs.Close()
	scanner := bufio.NewScanner(fs)
	scanStr := ""
	beginFlag, endFlag := false, false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, begin) {
			beginFlag = true
		} else if strings.Contains(line, end) {
			endFlag = true

		} else if beginFlag && !endFlag {
			fmt.Println(line)
			scanStr += "\r\n" + line
		}
	}
	return scanStr
}

func loadTableProtoFiles(pbDir string) {
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
			loadTableProtoFiles(file)
		} else {
			fileType := path.Ext(fileStat.Name())
			if fileType != ".proto" {
				continue
			}
			waitNum++
			fileList = append(fileList, file)
		}
	}
	wg.Add(waitNum)
	for _, file := range fileList {
		go parseNormalTableFile(file, &wg)
	}
	wg.Wait()
}

func parseNormalTableFile(name string, wg *sync.WaitGroup) {
	fs, err := os.OpenFile(name, os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("parseCmdFile error :", err, "name:", name)
		os.Exit(1)
	}
	defer fs.Close()
	defer wg.Done()
	fileNameWithExtension := filepath.Base(name)
	fmt.Println("File name with extension:", fileNameWithExtension)

	// 去除后缀，得到不带后缀的文件名
	fileNameWithoutExtension := strings.TrimSuffix(fileNameWithExtension, filepath.Ext(fileNameWithExtension))
	scanner := bufio.NewScanner(fs)
	// 逐行读取文件内容
	begin := false
	msgProto := TableMsgProto{}
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
			msgProto.PbName = ss[1]
			msgProto.FileName = fileNameWithoutExtension
			continue
		} else if !begin && len(msgProto.PbName) > 0 && line == "{" {
			begin = true
			continue
		} else if strings.Contains(line, "}") {
			tableFileProtoMap[msgProto.PbName] = msgProto
			msgProto = TableMsgProto{}
			begin = false
		} else if begin {
			endIndex := strings.Index(line, ";")
			if endIndex < 0 {
				continue
			}
			line = line[:endIndex]
			line = strings.ReplaceAll(line, "  ", " ")
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

			iVal, err := strconv.Atoi(strings.TrimSpace(ss[1]))
			if err != nil {
				fmt.Println(fmt.Sprintf("line:%s err :%s", line, err))
				continue
			}
			filed := TableFiledProto{
				FiledType:          fileType,
				FiledName:          fileName,
				Value:              iVal,
				CamelCaseFiledName: CamelCase(fileName),
				IsBaseType:         IsBaseType(fileType),
			}
			msgProto.FiledList = append(msgProto.FiledList, filed)
		}
	}
}
func CamelCase(str string) string {
	parts := strings.Split(str, "_")
	for i, part := range parts {
		parts[i] = strings.Title(part)
	}
	return strings.Join(parts, "")
}

func IsBaseType(filedType string) bool {
	filedType = strings.TrimSpace(filedType)
	return filedType == "uint64" || filedType == "int64" || filedType == "uint32" || filedType == "int32" || filedType == "uint16" || filedType == "int16" || filedType == "bool" || filedType == "string"
}
