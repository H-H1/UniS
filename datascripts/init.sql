-- ============================================================
-- uniS 数据库初始化脚本
-- 数据库: unis
-- 字符集: utf8mb4 / utf8mb4_unicode_ci
-- ============================================================

CREATE DATABASE IF NOT EXISTS `unis`
  DEFAULT CHARACTER SET utf8mb4
  DEFAULT COLLATE utf8mb4_unicode_ci;

USE `unis`;

-- ------------------------------------------------------------
-- 微信用户表
-- gorm.Model 对应: id, created_at, updated_at, deleted_at
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `users` (
  `id`          BIGINT UNSIGNED  NOT NULL AUTO_INCREMENT COMMENT '主键',
  `created_at`  DATETIME(3)      NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
  `updated_at`  DATETIME(3)      NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',
  `deleted_at`  DATETIME(3)               DEFAULT NULL COMMENT '软删除时间',

  `open_id`     VARCHAR(64)      NOT NULL COMMENT '微信 OpenID',
  `union_id`    VARCHAR(64)               DEFAULT '' COMMENT '微信 UnionID（开放平台）',
  `nick_name`   VARCHAR(64)               DEFAULT '' COMMENT '昵称',
  `avatar_url`  VARCHAR(512)              DEFAULT '' COMMENT '头像地址',
  `gender`      TINYINT          NOT NULL DEFAULT 0  COMMENT '性别 0未知 1男 2女',
  `country`     VARCHAR(64)               DEFAULT '' COMMENT '国家',
  `province`    VARCHAR(64)               DEFAULT '' COMMENT '省份',
  `city`        VARCHAR(64)               DEFAULT '' COMMENT '城市',
  `session_key` VARCHAR(128)              DEFAULT '' COMMENT '微信 SessionKey（不对外暴露）',

  PRIMARY KEY (`id`),
  UNIQUE KEY `udx_open_id`   (`open_id`),
  KEY        `idx_union_id`  (`union_id`),
  KEY        `idx_deleted_at`(`deleted_at`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='微信用户表';

-- ------------------------------------------------------------
-- 游客访问计数器表
-- 每个 URL 独立一行，记录累计访问次数
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `counters` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `created_at` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3)                       COMMENT '创建时间',
  `updated_at` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',

  `url`        VARCHAR(256)    NOT NULL                               COMMENT '被统计的请求路径，如 /counter',
  `count`      BIGINT UNSIGNED NOT NULL DEFAULT 0                     COMMENT '累计访问次数',

  PRIMARY KEY (`id`),
  UNIQUE KEY `udx_url` (`url`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='游客访问计数器表';

-- 预插入 /counter 路径的初始行，避免首次查询返回空
INSERT IGNORE INTO `counters` (`url`, `count`) VALUES ('/counter', 0);
