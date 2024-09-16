package logger

import "common"

var logDir = common.Context.Config.LogDir
var SystemLogger = Init(logDir, "system")
var NetLogger = Init(logDir, "net")
var DbLogger = Init(logDir, "db")
var RpcLogger = Init(logDir, "rpc")
