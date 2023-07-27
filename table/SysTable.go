package table

import (
	"database/sql"
	"errors"
	"fmt"
	"goserver/common/db"
	"goserver/common/logger"
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
	status      int32
}

type SysTableSqlOptional struct {
}

func (self *SysTableSqlOptional) OnQuerySuccess(flag bool, rows *sql.Rows) any {
	list := make([]SysTable, 1)
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
		fmt.Println(list)
		return list.([]SysTable), nil
	}
	return nil, errors.New(fmt.Sprintf("not found data by sql:%s", sql))
}

func (self *SysTableSqlOptional) SelectAll() ([]SysTable, error) {
	sql := "select * from crm_role"
	return self.selectSql(sql, nil)
}
func (self *SysTableSqlOptional) Save(data *SysTable) (bool, error) {
	sql := "insert into crm_role (id,description,role_name,status) values (?,?,?,?);"
	manger := db.GetDataBaseManger()
	if manger == nil || !manger.IsConnectFlag() {
		return false, errors.New("not found DataBaseManger or DataBaseManger not connect")
	}
	result, err := manger.GetDB().Exec(sql, data.id, data.description, data.role_name, data.status)
	if err != nil {
		return false, errors.New(fmt.Sprintf("save sql error table:SysTable, sql:%s, data:%s", sql, data))
	}
	id, err := result.LastInsertId()
	if err != nil {
		return false, errors.New(fmt.Sprintf(" get insert id error sql:%s, data:%s", sql, data))
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
	return &SysTableProxy{changeFlag: false}
}

func (self *SysTableProxy) SetId(id uint64) {
	self.id = id
	self.changeFlag = true
}
func (self *SysTableProxy) GetId() uint64 {
	return self.id
}
func (self *SysTableProxy) SetDescription(description string) {
	self.description = description
	self.changeFlag = true
}
func (self *SysTableProxy) GetDescription() string {
	return self.description
}
func (self *SysTableProxy) SetRole_name(roleName string) {
	self.role_name = roleName
	self.changeFlag = true
}
func (self *SysTableProxy) GetRole_name() string {
	return self.role_name
}
func (self *SysTableProxy) SetStatus(status int32) {
	self.status = status
	self.changeFlag = true
}
func (self *SysTableProxy) GetStatus() int32 {
	return self.status
}
func (self *SysTableProxy) AddStatus(addStatus int32) {
	self.status += addStatus
	self.changeFlag = true
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
			logger.Error(fmt.Sprintf("autoSaveData save error :%s, data:%s ", err, self.SysTable))
			return
		}
		if saveFlag {
			self.changeFlag = false
		} else {
			self.changeFlag = true
		}
	}
}

var cacheList = make([]func(), 10)

func addFun(f func()) {
	cacheList = append(cacheList, f)
	for _, autoSave := range cacheList {
		autoSave()
	}
}

func init() {
	bean := NewSysTable()
	addFun(bean.autoSaveData)
}
