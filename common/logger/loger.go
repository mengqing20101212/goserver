package logger

import (
	"fmt"
	"goserver/common"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

/*
日志模块，负责打印日志
*/
/**
日志打印级别
*/
const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
)

// 默认1G切割一个文件
const defaultSpiltLogSize = 1024 * 1024 * 1024

/*
日志文件的配置信息
*/
type LoggerConfig struct {
	//日志输出的目录
	logDir string
	//日志输出的文件名
	logName string
	//是否按照格式化输出 默认 true 是格式化输出 即打印额外的时间、行号等其他信息, false 是不按照格式化输出
	isFormat bool
	//日志切割的尺寸 默认1G切割
	spiltSize int64
}

type Logger struct {
	config      LoggerConfig
	level       int
	initFlag    bool
	receiverMsg chan string
	file        *os.File
}

var log Logger

func InitByConfig(config LoggerConfig) (success bool) {
	if log.initFlag {
		return false
	}
	return initBase(config)
}

func initBase(config LoggerConfig) (success bool) {
	log.config = config
	log.initFlag = true
	log.level = DEBUG
	log.receiverMsg = make(chan string, 1024*3)
	return tryOpenFile()
}
func Init(dir, filename string) (success bool) {
	if log.initFlag {
		return false
	}
	config := LoggerConfig{dir, filename, true, defaultSpiltLogSize}
	return initBase(config)
}

func tryOpenFile() bool {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get current working directory:", err)
		return false
	}

	relativePath := path.Join(log.config.logDir, log.config.logName)
	absolutePath := filepath.Join(currentDir, relativePath)
	// 创建目录及其父目录
	err = os.MkdirAll(filepath.Dir(absolutePath), 0755)
	if err != nil {
		fmt.Println("Failed to create directory:", err)
		return false
	}
	fs, err := os.OpenFile(absolutePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("open file error file:", absolutePath, " , err: ", err)
		return false
	}

	log.file = fs
	Info("init file success path:" + absolutePath + ", level:" + strconv.Itoa(log.level))
	go writeLog()
	return true
}

func writeLog() {
	for {
		str, ok := <-log.receiverMsg
		if !ok {
			return
		}
		_, err := fmt.Fprintln(os.Stdout, str)
		if err != nil {
			os.Exit(10)
		}
		_, err = fmt.Fprintln(log.file, str)
		if err != nil {
			os.Exit(11)
		}
		checkSpiltFile()
	}
}

func checkSpiltFile() {
	// 获取文件信息
	fileInfo, err := log.file.Stat()
	if err != nil {
		return
	}
	if fileInfo.Size() >= log.config.spiltSize {
		oldFileName := log.file.Name()
		dir, fileName := path.Split(strings.ReplaceAll(log.file.Name(), "\\", "/"))
		log.file.Close()
		newFIleName := path.Join(dir, fmt.Sprintf("%s_%s", common.GetYYYY_MM_DD_HH_mm_ss(), fileName))
		os.Rename(log.file.Name(), newFIleName)
		// 打开新的文件
		newFile, err := os.OpenFile(oldFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
		if err != nil {
			return
		}
		log.file = newFile
	}
}

func Info(msg string) {
	if log.level > INFO {
		return
	}
	if log.config.isFormat {
		log.receiverMsg <- format(msg, "INFO")
	} else {
		log.receiverMsg <- msg
	}
}

func format(msg string, level string) string {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	id := uint64(0)
	fmt.Sscanf(string(buf[:n]), "goroutine %d", &id)
	return fmt.Sprintf(" [coid:%d] %s %s %s,", id, time.Now().Format("2006-01-02 15:04:05.000"), level, msg)
}
func Debug(msg string) {
	if log.level > DEBUG {
		return
	}
	if log.config.isFormat {
		log.receiverMsg <- format(msg, "DEBUG")
	} else {
		log.receiverMsg <- msg
	}
}
func Warn(msg string) {
	if log.level > WARNING {
		return
	}
	if log.config.isFormat {
		log.receiverMsg <- format(msg, "WARN")
	} else {
		log.receiverMsg <- msg
	}
}

func Error(msg string) {
	if log.level > ERROR {
		return
	}
	if log.config.isFormat {
		log.receiverMsg <- format(msg, "ERROR")
	} else {
		log.receiverMsg <- msg
	}
}
