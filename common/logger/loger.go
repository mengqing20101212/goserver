package logger

import "io/fs"

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
const default_spilt_log_size = 1024 * 1024 * 1024

/*
日志文件的配置信息
*/
type loggerConfig struct {
	//日志输出的目录
	logDir string
	//日志输出的文件名
	logName string
	//是否按照格式化输出 默认 true 是格式化输出 即打印额外的时间、行号等其他信息, false 是不按照格式化输出
	isFormat bool
	//日志切割的尺寸 默认1G切割
	spiltSize int
}

type Logger struct {
	config loggerConfig
	file   fs.File
}

func (self Logger) InitByConfig(config loggerConfig) (success bool) {
	self.config = config
	return true
}
func (self Logger) Init(dir, filename *string) (success bool) {
	self.config = loggerConfig{*dir, *filename, true, default_spilt_log_size}
	println(self.config.logName)
	return true
}
