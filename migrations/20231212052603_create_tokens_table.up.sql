CREATE TABLE IF NOT EXISTS `tokens` (
  `id` INT(11) PRIMARY KEY NOT NULL,   
  `user_id` INT(11) UNSIGNED,       
  `name` VARCHAR(255)  NOT NULL default '',     
  `email`  VARCHAR(255) NOT NULL,
  `token_hash` VARBINARY(255) NOT NULL);

ALTER TABLE `tokens` ADD COLUMN `created_at` TIMESTAMP DEFAULT NOW();
ALTER TABLE `tokens` ADD COLUMN `updated_at` TIMESTAMP DEFAULT NOW();

