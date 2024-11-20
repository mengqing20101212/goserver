

CREATE TABLE IF NOT EXISTS `account` (`account_id` int NOT NULL COMMENT '账号id',
PRIMARY KEY (`account_id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci; 
 ALTER TABLE `account` ADD COLUMN `account_name` varchar(255) NULL COMMENT '账号' AFTER `account_id`; 
 ALTER TABLE `account` ADD COLUMN `create_timer` int NULL COMMENT '创建时间' AFTER `account_name`; 
 ALTER TABLE `account` ADD COLUMN `login_timer` int NULL COMMENT '登录时间' AFTER `create_timer`; 
 ALTER TABLE `account` ADD COLUMN `logout_timer` int NULL COMMENT '登出时间' AFTER `login_timer`; 
 ALTER TABLE `account` ADD COLUMN `phone` varchar(1024) NULL COMMENT 'len[1024]  手机号' AFTER `logout_timer`; 
 ALTER TABLE `account` ADD COLUMN `role_list` mediumblob NULL COMMENT '角色列表' AFTER `phone`; 
