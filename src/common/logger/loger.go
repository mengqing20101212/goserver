package logger

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

/*
日志模块，负责打印日志
日志打印级别
*/

type LoggerLevel int

const (
	DEBUG LoggerLevel = iota
	INFO
	WARNING
	ERROR
)

func (level LoggerLevel) string() string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	}
	return "NONE"
}

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
	isStdout  bool
}

type Logger struct {
	config      LoggerConfig
	level       LoggerLevel
	initFlag    bool
	receiverMsg chan *string
	file        *os.File
}

type LoggerFactory struct {
	loggerMap map[string]*Logger
	lock      sync.Mutex
}

var loggerFactory = LoggerFactory{make(map[string]*Logger), sync.Mutex{}}

func (self *LoggerConfig) getLogFullName() string {
	return path.Join(self.logDir, self.logName)
}
func (self *LoggerFactory) InitByConfig(config LoggerConfig) *Logger {
	self.lock.Lock()
	defer self.lock.Unlock()
	log, ok := self.loggerMap[config.getLogFullName()]
	if ok {
		return log
	}
	log = initBase(config)
	if log != nil {
		self.loggerMap[config.getLogFullName()] = log
	}
	return log
}

func initBase(config LoggerConfig) *Logger {
	log := new(Logger)
	log.config = config
	log.initFlag = true
	log.level = DEBUG
	log.receiverMsg = make(chan *string, 1024*3)
	if log.tryOpenFile() {
		return log
	}
	return nil
}
func InitNull() *Logger {
	currentDir, _ := os.Getwd()
	dir := path.Join(currentDir, "logs")
	fileName := "system.log"
	return Init(dir, fileName)
}
func Init(dir, filename string) *Logger {
	config := LoggerConfig{dir, filename, true, defaultSpiltLogSize, true}
	return loggerFactory.InitByConfig(config)
}

func (self *Logger) tryOpenFile() bool {
	if self.file != nil {
		return true
	}
	absolutePath := self.config.getLogFullName()
	// 创建目录及其父目录
	err := os.MkdirAll(filepath.Dir(absolutePath), 0755)
	if err != nil {
		fmt.Println("Failed to create directory:", err)
		return false
	}
	fs, err := os.OpenFile(absolutePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		fmt.Println("open file error file:", absolutePath, " , err: ", err)
		return false
	}

	self.file = fs
	self.Info("init file success path:" + absolutePath + ", level:" + self.level.string())
	go self.writeLog()
	return true
}

func (self *Logger) writeLog() {
	for {
		str, ok := <-self.receiverMsg
		if !ok {
			return
		}
		if self.config.isStdout {
			_, err := fmt.Fprintln(os.Stdout, *str)
			if err != nil {
				os.Exit(10)
			}
		}
		_, err := fmt.Fprintln(self.file, *str)
		if err != nil {
			os.Exit(11)
		}
		self.checkSpiltFile()
	}
}

func (self *Logger) checkSpiltFile() {
	// 获取文件信息
	fileInfo, err := self.file.Stat()
	if err != nil {
		return
	}
	if fileInfo.Size() >= self.config.spiltSize {
		oldFileName := self.file.Name()
		dir, fileName := path.Split(strings.ReplaceAll(self.file.Name(), "\\", "/"))
		_ = self.file.Close()
		newFIleName := path.Join(dir, fmt.Sprintf("%s_%s", time.Now().Format("2006_01_01_15_04_05"), fileName))
		_ = os.Rename(self.file.Name(), newFIleName)
		self.file.Close()
		// 打开新的文件
		newFile, err := os.OpenFile(oldFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
		if err != nil {
			return
		}
		self.file = newFile
	}
}

func (log *Logger) Info(msg string) {
	if log.level > INFO {
		return
	}
	if log.config.isFormat {
		s := formatMsg(&msg, INFO.string())
		log.receiverMsg <- s
	} else {
		log.receiverMsg <- &msg
	}
}

func formatMsg(msg *string, level string) *string {
	_, fileName, line, ok := runtime.Caller(2)
	if !ok {
		s := fmt.Sprintf("formatMsg error level:%s  msg:%s ", level, *msg)
		return &s
	}
	s := fmt.Sprintf("%s [%s:%d] [%s] %s", time.Now().Format("2006-01-02 15:04:05.000"), fileName, line, level, *msg)
	return &s
}

func format(msg string, level string) string {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	id := uint64(0)
	_, _ = fmt.Sscanf(string(buf[:n]), "goroutine %d", &id)
	pc := make([]uintptr, 10) // 假设最多获取 10 层调用栈
	n = runtime.Callers(1, pc)
	frames := runtime.CallersFrames(pc[:n])
	fileName := ""
	line := 0
	funcName := ""
	for i := 0; i < 3; i++ {
		frame, more := frames.Next()
		if !more {
			break
		}
		_, fileName = path.Split(frame.File)
		line = frame.Line
		funcName = strings.Split(frame.Function, ".")[1]
	}
	return fmt.Sprintf(" [coid:%d] %s [%s:%d] [%s] %s %s", id, time.Now().Format("2006-01-02 15:04:05.000"), fileName, line, funcName, level, msg)
}
func (log *Logger) Debug(msg string) {
	if log.level > DEBUG {
		return
	}
	if log.config.isFormat {
		s := formatMsg(&msg, log.level.string())
		log.receiverMsg <- s
	} else {
		log.receiverMsg <- &msg
	}
}
func (log *Logger) Warn(msg string) {
	if log.level > WARNING {
		return
	}
	if log.config.isFormat {
		s := formatMsg(&msg, log.level.string())
		log.receiverMsg <- s
	} else {
		log.receiverMsg <- &msg
	}
}

func (log *Logger) Error(msg string) {
	if log.level > ERROR {
		return
	}
	if log.config.isFormat {

		log.receiverMsg <- formatMsg(&msg, log.level.string())
	} else {
		log.receiverMsg <- &msg
	}
}

func (log *Logger) IsDebug() bool {
	return log.level == DEBUG
}
func (log *Logger) IsInfo() bool {
	return log.level == INFO
}
func (log *Logger) IsWARNING() bool {
	return log.level == WARNING
}
func (log *Logger) IsERROR() bool {
	return log.level == ERROR
}
