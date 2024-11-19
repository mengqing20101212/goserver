
CREATE TABLE IF NOT EXISTS `account` (`account_id` int NOT NULL,
PRIMARY KEY (`account_id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci; 
 ALTER TABLE `account` ADD COLUMN `account_name` varchar(255) NULL COMMENT ' 2' AFTER `account_id`; 
 ALTER TABLE `account` ADD COLUMN `create_timer` int NULL AFTER `account_name`; 
 ALTER TABLE `account` ADD COLUMN `login_timer` int NULL AFTER `create_timer`; 
 ALTER TABLE `account` ADD COLUMN `logout_timer` int NULL AFTER `login_timer`; 
 ALTER TABLE `account` ADD COLUMN `phone` varchar(255) NULL COMMENT ' 6' AFTER `logout_timer`; 
 ALTER TABLE `account` ADD COLUMN `role_list` mediumblob NULL AFTER `phone`; 

