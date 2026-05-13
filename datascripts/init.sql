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
-- 每 10 分钟插入一条增量记录，总访问量 = SUM(count)
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `counters` (
  `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `created_at` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3)                                COMMENT '创建时间',
  `updated_at` DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间',

  `url`        VARCHAR(256)    NOT NULL                               COMMENT '被统计的请求路径，如 /counter',
  `count`      BIGINT          NOT NULL DEFAULT 0                     COMMENT '本时间段内的访问增量',

  PRIMARY KEY (`id`),
  KEY `idx_url` (`url`)          -- 普通索引，允许多行，用于 SUM 查询
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='游客访问计数器表（每10分钟一条增量记录）';

-- 查询总访问量示例：SELECT SUM(count) FROM counters WHERE url = '/counter';

-- ------------------------------------------------------------
-- 测试结果表
-- 存储用户提交的测试答案和分数，scores/answer_details 为 JSON 类型
-- ------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `test_results` (
  `id`              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
  `created_at`      DATETIME(3)     NOT NULL DEFAULT CURRENT_TIMESTAMP(3)                                COMMENT '创建时间',
  
  `user_id`         BIGINT UNSIGNED NOT NULL                               COMMENT '用户 ID（关联 users.id）',
  `scores`          JSON            NOT NULL                               COMMENT '各类型分数，如 {"1":18,"2":12}',
  `answer_details`  JSON            NOT NULL                               COMMENT '答题详情数组',
  `total_questions` INT             NOT NULL                               COMMENT '总题数',
  `answered_count`  INT             NOT NULL                               COMMENT '已答题数',
  `timestamp`       BIGINT          NOT NULL                               COMMENT '前端提交时间戳（毫秒）',

  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB
  DEFAULT CHARSET=utf8mb4
  COLLATE=utf8mb4_unicode_ci
  COMMENT='测试结果表';
