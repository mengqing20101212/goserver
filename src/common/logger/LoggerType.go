package logger

var logDir string

var SystemLogger *Logger

var NetLogger *Logger

var DbLogger *Logger

var RpcLogger *Logger

func InitType(dir string) {
	logDir = dir
	SystemLogger = Init(logDir, "system")
	NetLogger = Init(logDir, "net")
	DbLogger = Init(logDir, "db")
	RpcLogger = Init(logDir, "rpc")
}
