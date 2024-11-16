package logger

var logDir string

var SystemLogger *Logger

var NetLogger *Logger

var DbLogger *Logger

var RpcLogger *Logger

func InitType(dir string) {
	logDir = dir
	SystemLogger = Init(logDir, "system.log")
	NetLogger = Init(logDir, "net.log")
	DbLogger = Init(logDir, "db.log")
	RpcLogger = Init(logDir, "rpc.log")
}
