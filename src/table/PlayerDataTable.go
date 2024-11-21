package table

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"goserver/common/db"
	"goserver/common/logger"
	"unsafe"
)

type PlayerDataTableSqlOptional struct {
}

func (self *PlayerDataTableSqlOptional) OnQuerySuccess(flag bool, rows *sql.Rows) any {
	list := make([]PlayerData, 0)
	for rows.Next() {
		data := PlayerData{}
		bs7 := make([]byte, 1024)
		bs8 := make([]byte, 1024)

		rows.Scan(&data.PlayerId, &data.PlayerName, &data.Level, &data.Exp, &data.Gold, &data.Diamond, &bs7, &bs8)

		userSetting := UserSetting{}
		err := proto.Unmarshal(bs7, &userSetting)
		if err != nil {
			logger.DbLogger.Error(fmt.Sprintf("Unmarshal bs:%v, data:%v", bs7, data))
			continue
		}
		data.UserSetting = &userSetting

		modules := Modules{}
		err = proto.Unmarshal(bs8, &modules)
		if err != nil {
			logger.DbLogger.Error(fmt.Sprintf("Unmarshal bs:%v, data:%v", bs8, data))
			continue
		}
		data.Modules = &modules

		list = append(list, data)
	}
	return list
}

func (self *PlayerDataTableSqlOptional) selectSql(sql string, params ...any) ([]PlayerData, error) {
	manger := db.GetDataBaseManger()
	if manger == nil || !manger.IsConnectFlag() {
		return nil, errors.New(" PlayerDataTableSqlOptional not found DataBaseManger or DataBaseManger not connect")
	}

	res, list := db.GetDataBaseManger().Query(sql, params, self)
	if res {
		return list.([]PlayerData), nil
	}
	return nil, errors.New(fmt.Sprintf("not found data by sql:%s", sql))
}
func (self *PlayerDataTableSqlOptional) Save(data *PlayerData) (bool, error) {
	UserSettingBytes := make([]byte, 1024)
	proto.UnmarshalMerge(UserSettingBytes, data.GetUserSetting())
	ModulesBytes := make([]byte, 1024)
	proto.UnmarshalMerge(ModulesBytes, data.GetModules())
	sql := "INSERT INTO playerdata (playerId, playerName, level, exp, gold, diamond, userSetting, modules) values (?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE   playerName=?, level=?, exp=?, gold=?, diamond=?, userSetting=?, modules=?;"
	manger := db.GetDataBaseManger()
	if manger == nil || !manger.IsConnectFlag() {
		return false, errors.New(" PlayerDataTableSqlOptional not found DataBaseManger or DataBaseManger not connect")
	}
	_, err := manger.ExecuteSqlResult(sql, &data.PlayerId, &data.PlayerName, &data.Level, &data.Exp, &data.Gold, &data.Diamond, &UserSettingBytes, &ModulesBytes)
	if err != nil {
		return false, errors.New(fmt.Sprintf("save sql error table:PlayerData, sql:%s, data:%v", sql, data))
	}
	return true, nil
}

func (self *PlayerDataTableSqlOptional) SelectAll() ([]PlayerData, error) {
	return self.selectSql("select * from account limit 1000", nil)
}

type PlayerDataTableProxy struct {
	PlayerData
	PlayerDataTableSqlOptional
	changeFlag bool
}

func (self *PlayerDataTableProxy) GetPlayerId() uint64 {
	return self.PlayerId
}
func (self *PlayerDataTableProxy) SetPlayerId(playerId uint64) *PlayerDataTableProxy {
	self.changeFlag = true
	self.PlayerId = playerId
	return self
}
func (self *PlayerDataTableProxy) GetPlayerName() string {
	return self.PlayerName
}
func (self *PlayerDataTableProxy) SetPlayerName(playerName string) *PlayerDataTableProxy {
	self.changeFlag = true
	self.PlayerName = playerName
	return self
}
func (self *PlayerDataTableProxy) GetLevel() int32 {
	return self.Level
}
func (self *PlayerDataTableProxy) SetLevel(level int32) *PlayerDataTableProxy {
	self.changeFlag = true
	self.Level = level
	return self
}
func (self *PlayerDataTableProxy) GetExp() uint32 {
	return self.Exp
}
func (self *PlayerDataTableProxy) SetExp(exp uint32) *PlayerDataTableProxy {
	self.changeFlag = true
	self.Exp = exp
	return self
}
func (self *PlayerDataTableProxy) GetGold() uint64 {
	return self.Gold
}
func (self *PlayerDataTableProxy) SetGold(gold uint64) *PlayerDataTableProxy {
	self.changeFlag = true
	self.Gold = gold
	return self
}
func (self *PlayerDataTableProxy) GetDiamond() uint64 {
	return self.Diamond
}
func (self *PlayerDataTableProxy) SetDiamond(diamond uint64) *PlayerDataTableProxy {
	self.changeFlag = true
	self.Diamond = diamond
	return self
}
func (self *PlayerDataTableProxy) GetUserSetting() *UserSetting {
	return self.UserSetting
}
func (self *PlayerDataTableProxy) SetUserSetting(userSetting *UserSetting) *PlayerDataTableProxy {
	self.changeFlag = true
	self.UserSetting = userSetting
	return self
}
func (self *PlayerDataTableProxy) GetModules() *Modules {
	return self.Modules
}
func (self *PlayerDataTableProxy) SetModules(modules *Modules) *PlayerDataTableProxy {
	self.changeFlag = true
	self.Modules = modules
	return self
}

// NewPlayerDataTable isCache 初始化该数据库PlayerData表的时候，是否给 TableCacheService 托管， true则数据的入库操作由
// TableCacheService 托管默认 10秒钟扫描该实体，保存入库。 false 则需要自己保存.
func NewPlayerDataTable(isCache bool) *PlayerDataTableProxy {
	table := &PlayerDataTableProxy{changeFlag: false}
	if isCache {
		db.GetCacheService().AddCacheFunc(unsafe.Pointer(table), table.autoSaveData)
	}
	return table
}

// Destroy 如果 PlayerData表的数据由 TableCacheService 托管时，在该数据生命周期结束卸载时候 需要手动调用一下该方法，从TableCacheService
// 中删除
func (self *PlayerDataTableProxy) Destroy() {
	db.GetCacheService().DelCacheFunc(unsafe.Pointer(self))
	if self.changeFlag {
		self.autoSaveData()
	}
}

func FindOnePlayerDataTable(id uint64, isCache bool) *PlayerDataTableProxy {
	sql := "select * from playerdata where id=?;"
	opt := PlayerDataTableSqlOptional{}
	list, err := opt.selectSql(sql, id)
	if err != nil {
		return nil
	}
	if len(list) > 0 {
		data := NewPlayerDataTable(isCache)
		data.PlayerData = list[0]
		data.PlayerDataTableSqlOptional = opt
		return data
	}
	return nil
}

func (self *PlayerDataTableProxy) autoSaveData() {
	if self.changeFlag {
		self.changeFlag = false
		saveFlag, err := self.Save(&self.PlayerData)
		if err != nil {
			self.changeFlag = true
			logger.DbLogger.Error(fmt.Sprintf("autoSaveData save error :%s, data:%v ", err, self.PlayerData))
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

//***** 自定义代码区 end ****
