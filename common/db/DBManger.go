package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"goserver/common/logger"
)

type TableInterface interface {
	ToSql() string
	GetParams() []any
	TableName() string
}

type DBManger struct {
	dbUrl       string
	db          *sqlx.DB
	connectFlag bool
}

func (self *DBManger) Execute(opt TableInterface) bool {
	_, err := self.db.Exec(opt.ToSql(), opt.GetParams())
	if err != nil {
		self.connectFlag = false
		logger.Error(fmt.Sprintf("Execute sql error:%s, sql:%s, params:%s", err, opt.ToSql(), opt.GetParams()))
		return false
	}
	logger.Info(fmt.Sprintf("sql:%s, params:%s", opt.ToSql(), opt.GetParams()))
	return true
}
func (self *DBManger) ExecuteSql(sql string) bool {
	_, err := self.db.Exec(sql)
	if err != nil {
		self.connectFlag = false
		logger.Error(fmt.Sprintf("Execute sql error:%s, sql:%s", err, sql))
		return false
	}
	logger.Info(fmt.Sprintf("sql:%s", sql))
	return true
}

var DbManger DBManger

func InitDefaultDataBase(userName, passWord, ip, databases string, port int32) bool {
	return InitDataBase(&DbManger, userName, passWord, ip, databases, port)
}

func InitDataBase(manger *DBManger, userName, passWord, ip, databases string, port int32) bool {
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", userName, passWord, ip, port, databases)
	database, err := sqlx.Open("mysql", dbUrl)
	if err != nil {
		//logger.Error(fmt.Sprintf(" InitDataBase init error  err :%s, userName:%s, passWord:%s, ip:%s, databases:%s, port:%d, dbUrl:%s", err, userName, passWord, ip, database, port, dbUrl))
		return false
	}
	err = database.Ping()
	if err != nil {
		database.Close()
		//logger.Error(fmt.Sprintf(" InitDataBase Ping databases error err :%s, userName:%s, passWord:%s, ip:%s, databases:%s, port:%d, dbUrl:%s", err, userName, passWord, ip, database, port, dbUrl))
		return false
	}
	manger.dbUrl = dbUrl
	manger.db = database
	manger.connectFlag = true
	//logger.Info(fmt.Sprintf("InitDataBase success dbUrl:%s", dbUrl))
	return true
}
