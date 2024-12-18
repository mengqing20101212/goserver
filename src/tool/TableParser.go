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
	FiledType          string //Mysql 字段
	FiledTypeLen       string //MySQL 字段的长度 有可能是 text
	FileShowDesc       string //显示的注释
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

func main1() {
	fmt.Println("start parse table buffer")
	args := os.Args
	if len(args) > 1 {
		TablePbDir = args[1]
	} else {
		TablePbDir, _ = os.Getwd()
	}
	loadTableProtoFiles(TablePbDir)
	//fmt.Println(tableFileProtoMap)
	createGoTableFile()
	createSqlFile()
	fmt.Println("end parse table buffer")
}

func createSqlFile() {
	tempOutFile := filepath.Join(TablePbDir, "../script/sql", "tempInit.sql")
	outFile := filepath.Join(TablePbDir, "../script/sql", "Init.sql")
	os.Remove(tempOutFile)
	sqlList := make([]string, 0)
	for _, data := range tableFileProtoMap {
		if strings.ToLower(data.PbName) != strings.ToLower(data.FileName) {
			continue
		}
		sql := ""
		fmt.Sprintf("createSqlFile:" + data.FileName)
		lastFiled := TableFiledProto{}
		for i, filed := range data.FiledList {
			if i == 0 {
				sql += "\r\nCREATE TABLE IF NOT EXISTS `" + strings.ToLower(data.FileName) + "` ("
				sql += "`" + filed.FiledName + "` " + getDBType(filed.FiledType, filed.FiledTypeLen) + " NOT NULL"
				if len(filed.FileShowDesc) > 0 {
					sql += " COMMENT '" + filed.FileShowDesc + "'"
				}
				sql += ",\r\n"
				sql += "PRIMARY KEY (`" + filed.FiledName + "`)"
				sql += ") ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci; \n"
			} else {
				sql += " ALTER TABLE `" + strings.ToLower(data.FileName) + "` ADD COLUMN `" + filed.FiledName + "` " + getDBType(filed.FiledType, filed.FiledTypeLen) + " NULL"
				if len(filed.FileShowDesc) > 0 {
					sql += " COMMENT '" + filed.FileShowDesc + "'"
				}
				sql += " AFTER `" + lastFiled.FiledName + "`; \n"
			}
			lastFiled = filed
		}
		sqlList = append(sqlList, sql)
		fmt.Println("createSql:" + sql)
	}

	fs, err := os.OpenFile(tempOutFile, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("open file tempOutFile:", tempOutFile, ", err:", err)
		return
	}
	fs.WriteString(strings.Join(sqlList, "\r\n"))
	os.Remove(outFile)
	fs.Close()
	os.Rename(tempOutFile, outFile)
}

func getDBType(filedType, fileTypeLen string) string {
	if filedType == "uint64" || filedType == "int64" || filedType == "uint32" || filedType == "int32" {
		return "int"

	} else if filedType == "bool" {
		return "TINYINT(1)"
	} else if filedType == "string" {
		if fileTypeLen == "text" {
			return "text"
		} else {
			return "varchar(" + fileTypeLen + ")"
		}
	} else {
		return "mediumblob"
	}

}

func createGoTableFile() {

	relativePath := "Table.tmpl"
	absPath := filepath.Join(TablePbDir, relativePath)
	for fileName, proto := range tableFileProtoMap {
		var buf bytes.Buffer
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
		if len(fileName) < 2 {
			continue
		}
		if fileName != proto.FileName {
			continue
		}
		err = tmpl.Execute(&buf, proto)
		if err != nil {
			fmt.Println("mapTemplate Execute error :", err)
		}
		outFile := "../src/table"
		outFile = filepath.Join(TablePbDir, outFile, fileName+"Table.go")
		strBegin := "//***** 自定义代码区 begin ****"
		strEnd := "//***** 自定义代码区 end ****"
		scanExtCode := ScanOutFileExtCode(outFile, strBegin[2:], strEnd)
		os.Remove(outFile)
		fs, err := os.OpenFile(outFile, os.O_RDWR|os.O_CREATE, 0755)
		defer fs.Close()
		if err != nil {
			fmt.Println(" err:", err)
			return
		}
		fmt.Println("outFile:", outFile, ", scanExtCode: ", *scanExtCode)
		fs.Write(buf.Bytes())
		fs.WriteString(strBegin)
		fs.WriteString("\r\n")
		fs.WriteString(*scanExtCode)
		fs.WriteString("\r\n")
		fs.WriteString(strEnd)
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
	index := 0
	for _, filed := range list {
		if !IsBaseType(filed.FiledType) {
			result += "\r\n"
			result += fmt.Sprintf("	   	%s := %s{}\n", filed.FiledName, filed.FiledType)
			if index == 0 {
				result += fmt.Sprintf("		err := proto.Unmarshal(bs%d, &%s)\n", filed.Value, filed.FiledName)
			} else {
				result += fmt.Sprintf("		err = proto.Unmarshal(bs%d, &%s)\n", filed.Value, filed.FiledName)
			}
			index++
			result += fmt.Sprintf("		if err != nil {\n")
			result += fmt.Sprintf("		logger.DbLogger.Error(fmt.Sprintf(\"Unmarshal bs:%%v, data:%%v\", bs%d, data)) \n", filed.Value)
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
	result += fmt.Sprintf("_, err := manger.ExecuteSqlResult(sql %s)\n", strings.TrimSuffix(params, ","))
	result += "if err != nil {\n"
	result += fmt.Sprintf("return false, errors.New(fmt.Sprintf(\"save sql error table:%s, sql:%%s, data:%%v\", sql, data))\n", data.FileName)
	result += "}"
	return result
}

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
	fmt.Println("解析proto文件:", name)

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
		oldLine := line
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
			fileType := leftss[0] //字段类型
			fileName := leftss[1] //字段名称
			fileTypeLen := "255"  //字段类型长度 默认是 0
			showDescArr := strings.Split(oldLine, "//")
			showDesc := ""
			if len(showDescArr) > 1 {
				showDesc = showDescArr[1]
			}
			if fileType == "string" {
				rss := ""
				if len(showDescArr) > 1 {
					rss = showDescArr[1]
					e := 0
					if strings.Contains(rss, "len[") {
						b := strings.Index(rss, "len[")
						e := strings.Index(rss, "]")
						fileTypeLen = rss[b+4 : e]
					}
					showDesc = rss[e:]
				}
			}

			iVal, err := strconv.Atoi(strings.TrimSpace(ss[1]))
			if err != nil {
				fmt.Println(fmt.Sprintf("line:%s err :%s", line, err))
				continue
			}
			filed := TableFiledProto{
				FiledType:          fileType,
				FiledTypeLen:       fileTypeLen,
				FileShowDesc:       showDesc,
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
