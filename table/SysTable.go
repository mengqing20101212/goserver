package table

import (
	"database/sql"
	"errors"
	"fmt"
	"goserver/common/db"
	"goserver/common/logger"
	"unsafe"
)

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
	status      bool
}

type SysTableSqlOptional struct {
}

func (self *SysTableSqlOptional) OnQuerySuccess(flag bool, rows *sql.Rows) any {
	list := make([]SysTable, 0)
	for rows.Next() {
		data := SysTable{}
		rows.Scan(&data.id, &data.description, &data.role_name, &data.status)
		list = append(list, data)
	}
	return list
}

func (self *SysTableSqlOptional) selectSql(sql string, params any) ([]SysTable, error) {
	manger := db.GetDataBaseManger()
	if manger == nil || !manger.IsConnectFlag() {
		return nil, errors.New("not found DataBaseManger or DataBaseManger not connect")
	}

	res, list := db.GetDataBaseManger().Query(sql, params, self)
	if res {
		return list.([]SysTable), nil
	}
	return nil, errors.New(fmt.Sprintf("not found data by sql:%s", sql))
}

func (self *SysTableSqlOptional) SelectAll() ([]SysTable, error) {
	sql := "select * from crm_role"
	return self.selectSql(sql, nil)
}
func (self *SysTableSqlOptional) Save(data *SysTable) (bool, error) {
	sql := "INSERT INTO crm_role (description,role_name,status) values (?,?,?) ON DUPLICATE KEY UPDATE description=?, status=?;"
	manger := db.GetDataBaseManger()
	if manger == nil || !manger.IsConnectFlag() {
		return false, errors.New("not found DataBaseManger or DataBaseManger not connect")
	}
	result, err := manger.GetDB().Exec(sql, data.description, data.role_name, data.status, data.description, data.status)
	if err != nil {
		return false, errors.New(fmt.Sprintf("save sql error table:SysTable, sql:%s, data:%v", sql, data))
	}
	id, err := result.LastInsertId()
	if err != nil {
		return false, errors.New(fmt.Sprintf(" get insert id error sql:%s, data:%v", sql, data))
	}
	data.id = uint64(id)
	return true, nil
}

type SysTableProxy struct {
	SysTable
	SysTableSqlOptional
	changeFlag bool
}

func NewSysTable() *SysTableProxy {
	table := &SysTableProxy{changeFlag: false}
	db.GetCacheService().AddCacheFunc(unsafe.Pointer(table), table.autoSaveData)
	return table
}

func FindOneSysTable(id uint64) *SysTableProxy {
	sql := "select * from crm_role where id=?;"
	opt := SysTableSqlOptional{}
	list, err := opt.selectSql(sql, id)
	if err != nil {
		return nil
	}
	if len(list) > 0 {
		data := &SysTableProxy{
			list[0],
			opt,
			false,
		}
		return data
	}
	return nil
}

func (self *SysTableProxy) SetId(id uint64) *SysTableProxy {
	self.id = id
	self.changeFlag = true
	return self
}
func (self *SysTableProxy) GetId() uint64 {
	return self.id
}
func (self *SysTableProxy) SetDescription(description string) *SysTableProxy {
	self.description = description
	self.changeFlag = true
	return self
}
func (self *SysTableProxy) GetDescription() string {
	return self.description
}
func (self *SysTableProxy) SetRole_name(roleName string) *SysTableProxy {
	self.role_name = roleName
	self.changeFlag = true
	return self
}
func (self *SysTableProxy) GetRole_name() string {
	return self.role_name
}
func (self *SysTableProxy) SetStatus(status bool) *SysTableProxy {
	self.status = status
	self.changeFlag = true
	return self
}
func (self *SysTableProxy) GetStatus() bool {
	return self.status
}
func (self *SysTableProxy) AddStatus(addStatus bool) *SysTableProxy {
	self.status = addStatus
	self.changeFlag = true
	return self
}

func (self *SysTableProxy) GetSysTable() *SysTable {
	return &self.SysTable
}
func (self *SysTableProxy) autoSaveData() {
	if self.changeFlag {
		self.changeFlag = false
		saveFlag, err := self.Save(&self.SysTable)
		if err != nil {
			self.changeFlag = true
			logger.Error(fmt.Sprintf("autoSaveData save error :%s, data:%v ", err, self.SysTable))
			return
		}
		if saveFlag {
			self.changeFlag = false
		} else {
			self.changeFlag = true
		}
	}
}

// *****begin****//

func (self *SysTableProxy) Test1() {
	self.changeFlag = true
}
func (self *SysTableProxy) Test2() {
	self.changeFlag = true
}

//CREATE TABLE `account` (
//`account_id` int NOT NULL,
//PRIMARY KEY (`account_id`)
//) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
//ALTER TABLE `account` ADD COLUMN `account_name` varchar(255) NULL AFTER `account_id`;
//*****end****//
