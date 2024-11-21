
CREATE TABLE IF NOT EXISTS `playerdata` (`playerId` int NOT NULL COMMENT '玩家id',
PRIMARY KEY (`playerId`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci; 
 ALTER TABLE `playerdata` ADD COLUMN `playerName` varchar(255) NULL COMMENT '玩家名称' AFTER `playerId`; 
 ALTER TABLE `playerdata` ADD COLUMN `level` int NULL COMMENT '玩家等级' AFTER `playerName`; 
 ALTER TABLE `playerdata` ADD COLUMN `exp` int NULL COMMENT '玩家经验' AFTER `level`; 
 ALTER TABLE `playerdata` ADD COLUMN `gold` int NULL COMMENT '金币' AFTER `exp`; 
 ALTER TABLE `playerdata` ADD COLUMN `diamond` int NULL COMMENT '钻石' AFTER `gold`; 
 ALTER TABLE `playerdata` ADD COLUMN `userSetting` mediumblob NULL COMMENT '玩家设置' AFTER `diamond`; 
 ALTER TABLE `playerdata` ADD COLUMN `modules` mediumblob NULL COMMENT '各个模块数据' AFTER `userSetting`; 


CREATE TABLE IF NOT EXISTS `account` (`account_id` int NOT NULL COMMENT '账号id',
PRIMARY KEY (`account_id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci; 
 ALTER TABLE `account` ADD COLUMN `account_name` varchar(255) NULL COMMENT '账号' AFTER `account_id`; 
 ALTER TABLE `account` ADD COLUMN `create_timer` int NULL COMMENT '创建时间' AFTER `account_name`; 
 ALTER TABLE `account` ADD COLUMN `login_timer` int NULL COMMENT '登录时间' AFTER `create_timer`; 
 ALTER TABLE `account` ADD COLUMN `logout_timer` int NULL COMMENT '登出时间' AFTER `login_timer`; 
 ALTER TABLE `account` ADD COLUMN `phone` varchar(1024) NULL COMMENT 'len[1024]  手机号' AFTER `logout_timer`; 
 ALTER TABLE `account` ADD COLUMN `role_list` mediumblob NULL COMMENT '角色列表' AFTER `phone`; 

