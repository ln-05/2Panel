-- 创建 sys_databases 表
USE 123;

CREATE TABLE IF NOT EXISTS `sys_databases` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '数据库名称',
  `type` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '数据库类型',
  `host` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主机地址',
  `port` bigint NOT NULL COMMENT '端口号',
  `username` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名',
  `password` varchar(191) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码',
  `database` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '数据库名',
  `description` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '描述信息',
  `status` varchar(191) COLLATE utf8mb4_unicode_ci DEFAULT 'unknown' COMMENT '连接状态',
  `last_test_at` datetime(3) DEFAULT NULL COMMENT '最后测试时间',
  PRIMARY KEY (`id`),
  KEY `idx_sys_databases_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='数据库连接配置表';

-- 插入一些示例数据
INSERT INTO `sys_databases` (`name`, `type`, `host`, `port`, `username`, `password`, `database`, `description`, `status`, `created_at`, `updated_at`) VALUES
('本地MySQL', 'mysql', 'localhost', 3306, 'root', 'password123', 'test_db', '本地测试数据库', 'unknown', NOW(), NOW()),
('生产MySQL', 'mysql', '192.168.1.100', 3306, 'prod_user', 'prod_pass', 'production', '生产环境数据库', 'unknown', NOW(), NOW());