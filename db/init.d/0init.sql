CREATE DATABASE `message`;

USE `message`;

SET NAMES utf8mb4;

DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '10000から始まるuniq id',
  `email` varchar(256) NOT NULL COMMENT 'ログインに使用される',
  `password` varchar(256) NOT NULL COMMENT 'ハッシュ化されたパスワード、ソルトなどを保持する',
  `name` varchar(256) NOT NULL COMMENT '表示名',
  `login_at` datetime NOT NULL,
  `update_at` datetime NOT NULL,
  `create_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE current_timestamp(),
  `deleted` tinyint(1) unsigned zerofill DEFAULT 0 COMMENT '1 ... 削除',
  `two_step_verification_code` varchar(6) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL COMMENT '2段階認証コード',
  `two_verificated_at` datetime DEFAULT NULL COMMENT '2段階認証成功日時',
  `two_verificated` tinyint(1) unsigned zerofill DEFAULT 0 COMMENT '2段階認証成功フラグ 1..success',
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='userの基本情報を保持する';


DROP TABLE IF EXISTS `webpush`;
CREATE TABLE `webpush` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned NOT NULL,
  `endpoint` varchar(256) DEFAULT NULL,
  `p256dh` varchar(128) DEFAULT NULL,
  `auth` text DEFAULT NULL,
  `expiration_time` datetime DEFAULT NULL,
  `user_agent` varchar(256) DEFAULT NULL,
  `create_at` timestamp NOT NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE current_timestamp(),
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `webpush_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin COMMENT='store webpush of users. hard delete';

