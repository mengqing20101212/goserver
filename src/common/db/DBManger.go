package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"logger"
)

type DBManger struct {
	dbUrl       string
	db          *sqlx.DB
	connectFlag bool
}

var DbLogger *logger.Logger

func (self *DBManger) IsConnectFlag() bool {
	return self.connectFlag
}

func (self *DBManger) Execute(sql string, params ...any) bool {
	_, err := self.db.Exec(sql, params)
	if err != nil {
		self.connectFlag = false
		DbLogger.Error(fmt.Sprintf("Execute sql error:%s, sql:%s, params:%s", err, sql, params))
		return false
	}
	DbLogger.Info(fmt.Sprintf("sql:%s, params:%s", sql, params))
	return true
}

func (this *DBManger) ExecuteSqlResult(sql string, data ...any) (sql.Result, error) {
	return this.db.Exec(sql, data)
}
func (self *DBManger) ExecuteSql(sql string) bool {
	_, err := self.db.Exec(sql)
	if err != nil {
		self.connectFlag = false
		DbLogger.Error(fmt.Sprintf("Execute sql error:%s, sql:%s", err, sql))
		return false
	}
	DbLogger.Info(fmt.Sprintf("sql:%s", sql))
	return true
}
func (self *DBManger) Query(sqlStr string, params any, sqlOpt TableInterface) (bool, any) {
	var rows *sql.Rows
	var err error
	if params == nil {
		rows, err = self.db.Query(sqlStr)
	} else {
		rows, err = self.db.Query(sqlStr, params)
	}
	if err != nil {
		DbLogger.Error(fmt.Sprintf("Query data error sqlStr:%s, params:%s, error:%s", sqlStr, params, err))
		return false, nil
	}
	defer rows.Close()
	resultList := sqlOpt.OnQuerySuccess(true, rows)
	return resultList != nil, resultList
}

func (self *DBManger) Insert(sql string, params ...any) bool {
	return self.Execute(sql, params)
}

func (self *DBManger) Update(sql string, params ...any) bool {
	return self.Execute(sql, params)
}
func (self *DBManger) GetDB() *sqlx.DB {
	return self.db
}

func (self *DBManger) CloseDBConnect() {
	err := self.db.Close()
	if err != nil {
		return
	}

}

var DbManger DBManger

func InitDefaultDataBase(userName, passWord, ip, databases string, port int32) bool {
	return InitDataBase(&DbManger, userName, passWord, ip, databases, port)
}

func InitDataBase(manger *DBManger, userName, passWord, ip, databases string, port int32) bool {
	if DbLogger == nil {
		DbLogger = logger.DbLogger
	}
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", userName, passWord, ip, port, databases)
	database, err := sqlx.Open("mysql", dbUrl)
	if err != nil {
		//logger.Error(fmt.Sprintf(" InitDataBase init error  err :%s, userName:%s, passWord:%s, ip:%s, databases:%s, port:%d, dbUrl:%s", err, userName, passWord, ip, database, port, dbUrl))
		return false
	}
	err = database.Ping()
	if err != nil {
		database.Close()
		DbLogger.Error(fmt.Sprintf(" InitDataBase Ping databases error err :%s, userName:%s, passWord:%s, ip:%s, databases:%s, port:%d, dbUrl:%s", err, userName, passWord, ip, database, port, dbUrl))
		return false
	}
	database.SetMaxIdleConns(10)
	database.SetMaxOpenConns(10)
	manger.dbUrl = dbUrl
	manger.db = database
	manger.connectFlag = true
	DbLogger.Info(fmt.Sprintf("InitDataBase success dbUrl:%s", dbUrl))
	return true
}

func GetDataBaseManger() *DBManger {
	return &DbManger
}

type TableInterface interface {
	OnQuerySuccess(flag bool, rows *sql.Rows) any
}
