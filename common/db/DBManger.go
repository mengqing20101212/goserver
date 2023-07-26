package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"goserver/common/logger"
)

type DBManger struct {
	dbUrl       string
	db          *sqlx.DB
	connectFlag bool
}

func (self *DBManger) Execute(sql string, params any) bool {
	_, err := self.db.Exec(sql, params)
	if err != nil {
		self.connectFlag = false
		logger.Error(fmt.Sprintf("Execute sql error:%s, sql:%s, params:%s", err, sql, params))
		return false
	}
	logger.Info(fmt.Sprintf("sql:%s, params:%s", sql, params))
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
func (self *DBManger) Query(sqlStr string, params any, sqlOpt TableInterface) (bool, any) {
	var rows *sql.Rows
	var err error
	if params == nil {
		rows, err = self.db.Query(sqlStr)
	} else {
		rows, err = self.db.Query(sqlStr, params)
	}
	if err != nil {
		logger.Error(fmt.Sprintf("Query data error sqlStr:%s, params:%s, error:%s", sqlStr, params, err))
		return false, nil
	}
	defer rows.Close()
	return true, sqlOpt.OnQuerySuccess(true, rows)
}

func (self *DBManger) Insert(sql string, params any) bool {
	return self.Execute(sql, params)
}

func (self *DBManger) Update(sql string, params any) bool {
	return self.Execute(sql, params)

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

func GetDataBaseManger() *DBManger {
	return &DbManger
}

type TableInterface interface {
	OnQuerySuccess(flag bool, rows *sql.Rows) any
}

/*
*
CREATE TABLE `crm_role` (
`id` bigint NOT NULL AUTO_INCREMENT,
`description` varchar(512) DEFAULT NULL,
`role_name` varchar(30) NOT NULL,
`status` bit(1) DEFAULT NULL,
PRIMARY KEY (`id`),
UNIQUE KEY `UK_r0jsnwb00o0n376ghyuahuqfg` (`role_name`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb3;
*
*/
type SysTable struct {
	id          uint64
	description string
	role_name   string
	status      int32
}

type SysTableSqlOpt struct {
}

func (self *SysTableSqlOpt) OnQuerySuccess(flag bool, rows *sql.Rows) any {
	list := make([]SysTable, 1)
	for rows.Next() {
		data := SysTable{}
		rows.Scan(&data.id, &data.description, &data.role_name, &data.status)
		list = append(list, data)
	}
	return list
}

func (self *SysTableSqlOpt) selectSql(sql string, params any) ([]SysTable, error) {
	manger := GetDataBaseManger()
	if manger == nil || !manger.connectFlag {
		return nil, errors.New("not found DataBaseManger or DataBaseManger not connect")
	}

	res, list := GetDataBaseManger().Query(sql, params, self)
	if res {
		fmt.Println(list)
		return list.([]SysTable), nil
	}
	return nil, errors.New(fmt.Sprintf("not found data by sql:%s, urlDb:%s", sql, manger.dbUrl))
}

func (self *SysTableSqlOpt) SelectAll() ([]SysTable, error) {
	sql := "select * from crm_role"
	return self.selectSql(sql, nil)
}
