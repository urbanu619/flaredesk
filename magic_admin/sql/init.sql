-- Flaredesk 业务库初始化
-- 执行前请先创建数据库: CREATE DATABASE flaredesk DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE IF NOT EXISTS `cf_account` (
  `id`         bigint       NOT NULL AUTO_INCREMENT COMMENT 'id' PRIMARY KEY,
  `created_at` bigint       DEFAULT NULL COMMENT '创建时间',
  `updated_at` bigint       DEFAULT NULL COMMENT '更新时间',
  `name`       varchar(100) DEFAULT NULL COMMENT '账号名称',
  `email`      varchar(255) DEFAULT NULL COMMENT 'CF账号邮箱',
  `api_token`  varchar(500) DEFAULT NULL COMMENT 'API Token',
  `status`     tinyint      DEFAULT 1    COMMENT '状态: 1=正常, 2=禁用',
  `remark`     varchar(500) DEFAULT NULL COMMENT '备注'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Cloudflare账号表';
