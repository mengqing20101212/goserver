package table

import (
	"database/sql"
	"db"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"logger"
	"unsafe"
)

type AccountTableSqlOptional struct {
}

func (self *AccountTableSqlOptional) OnQuerySuccess(flag bool, rows *sql.Rows) any {
	list := make([]Account, 0)
	for rows.Next() {
		data := Account{}
		bs7 := make([]byte, 1024)

		rows.Scan(&data.AccountId, &data.AccountName, &data.CreateTimer, &data.LoginTimer, &data.LogoutTimer, &data.Phone, &bs7)

		role_list := RoleShowList{}
		err := proto.Unmarshal(bs7, &role_list)
		if err != nil {
			logger.DbLogger.Error(fmt.Sprintf("Unmarshal bs:%v, data:%v", bs7, data))
			continue
		}
		data.RoleList = &role_list

		list = append(list, data)
	}
	return list
}

func (self *AccountTableSqlOptional) selectSql(sql string, params ...any) ([]Account, error) {
	manger := db.GetDataBaseManger()
	if manger == nil || !manger.IsConnectFlag() {
		return nil, errors.New(" AccountTableSqlOptional not found DataBaseManger or DataBaseManger not connect")
	}

	res, list := db.GetDataBaseManger().Query(sql, params, self)
	if res {
		return list.([]Account), nil
	}
	return nil, errors.New(fmt.Sprintf("not found data by sql:%s", sql))
}
func (self *AccountTableSqlOptional) Save(data *Account) (bool, error) {
	RoleListBytes := make([]byte, 1024)
	proto.UnmarshalMerge(RoleListBytes, data.GetRoleList())
	sql := "INSERT INTO account (account_id, account_name, create_timer, login_timer, logout_timer, phone, role_list) values (?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE   account_name=?, create_timer=?, login_timer=?, logout_timer=?, phone=?, role_list=?;"
	manger := db.GetDataBaseManger()
	if manger == nil || !manger.IsConnectFlag() {
		return false, errors.New(" AccountTableSqlOptional not found DataBaseManger or DataBaseManger not connect")
	}
	_, err := manger.ExecuteSqlResult(sql, &data.AccountId, &data.AccountName, &data.CreateTimer, &data.LoginTimer, &data.LogoutTimer, &data.Phone, &RoleListBytes)
	if err != nil {
		return false, errors.New(fmt.Sprintf("save sql error table:Account, sql:%s, data:%v", sql, data))
	}
	return true, nil
}

func (self *AccountTableSqlOptional) SelectAll() ([]Account, error) {
	return self.selectSql("select * from account limit 1000", nil)
}

type AccountTableProxy struct {
	Account
	AccountTableSqlOptional
	changeFlag bool
}

func (self *AccountTableProxy) GetAccountId() uint64 {
	return self.AccountId
}
func (self *AccountTableProxy) SetAccountId(account_id uint64) *AccountTableProxy {
	self.changeFlag = true
	self.AccountId = account_id
	return self
}
func (self *AccountTableProxy) GetAccountName() string {
	return self.AccountName
}
func (self *AccountTableProxy) SetAccountName(account_name string) *AccountTableProxy {
	self.changeFlag = true
	self.AccountName = account_name
	return self
}
func (self *AccountTableProxy) GetCreateTimer() uint32 {
	return self.CreateTimer
}
func (self *AccountTableProxy) SetCreateTimer(create_timer uint32) *AccountTableProxy {
	self.changeFlag = true
	self.CreateTimer = create_timer
	return self
}
func (self *AccountTableProxy) GetLoginTimer() uint32 {
	return self.LoginTimer
}
func (self *AccountTableProxy) SetLoginTimer(login_timer uint32) *AccountTableProxy {
	self.changeFlag = true
	self.LoginTimer = login_timer
	return self
}
func (self *AccountTableProxy) GetLogoutTimer() uint32 {
	return self.LogoutTimer
}
func (self *AccountTableProxy) SetLogoutTimer(logout_timer uint32) *AccountTableProxy {
	self.changeFlag = true
	self.LogoutTimer = logout_timer
	return self
}
func (self *AccountTableProxy) GetPhone() string {
	return self.Phone
}
func (self *AccountTableProxy) SetPhone(phone string) *AccountTableProxy {
	self.changeFlag = true
	self.Phone = phone
	return self
}
func (self *AccountTableProxy) GetRoleList() *RoleShowList {
	return self.RoleList
}
func (self *AccountTableProxy) SetRoleList(role_list *RoleShowList) *AccountTableProxy {
	self.changeFlag = true
	self.RoleList = role_list
	return self
}

// NewAccountTable isCache 初始化该数据库Account表的时候，是否给 TableCacheService 托管， true则数据的入库操作由
// TableCacheService 托管默认 10秒钟扫描该实体，保存入库。 false 则需要自己保存.
func NewAccountTable(isCache bool) *AccountTableProxy {
	table := &AccountTableProxy{changeFlag: false}
	if isCache {
		db.GetCacheService().AddCacheFunc(unsafe.Pointer(table), table.autoSaveData)
	}
	return table
}

// Destroy 如果 Account表的数据由 TableCacheService 托管时，在该数据生命周期结束卸载时候 需要手动调用一下该方法，从TableCacheService
// 中删除
func (self *AccountTableProxy) Destroy() {
	db.GetCacheService().DelCacheFunc(unsafe.Pointer(self))
	if self.changeFlag {
		self.autoSaveData()
	}
}

func FindOneAccountTable(id uint64, isCache bool) *AccountTableProxy {
	sql := "select * from account where id=?;"
	opt := AccountTableSqlOptional{}
	list, err := opt.selectSql(sql, id)
	if err != nil {
		return nil
	}
	if len(list) > 0 {
		data := NewAccountTable(isCache)
		data.Account = list[0]
		data.AccountTableSqlOptional = opt
		return data
	}
	return nil
}

func (self *AccountTableProxy) autoSaveData() {
	if self.changeFlag {
		self.changeFlag = false
		saveFlag, err := self.Save(&self.Account)
		if err != nil {
			self.changeFlag = true
			logger.DbLogger.Error(fmt.Sprintf("autoSaveData save error :%s, data:%v ", err, self.Account))
			return
		}
		if saveFlag {
			self.changeFlag = false
		} else {
			self.changeFlag = true
		}
	}
}

//***** 自定义代码区 begin ****

func (receiver AccountTableProxy) test1() {
	fmt.Println("test1")
}
func (receiver AccountTableProxy) test2() {
	fmt.Println("test1")
}

//***** 自定义代码区 end ****
