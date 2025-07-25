-- 创建数据库管理表
CREATE DATABASE IF NOT EXISTS myapp CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE myapp;

-- 数据库连接配置表
CREATE TABLE IF NOT EXISTS sys_databases (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '数据库名称',
    type ENUM('mysql', 'postgresql', 'redis', 'mongodb') NOT NULL COMMENT '数据库类型',
    host VARCHAR(255) NOT NULL COMMENT '主机地址',
    port INT NOT NULL COMMENT '端口号',
    database_name VARCHAR(100) COMMENT '数据库名',
    username VARCHAR(100) NOT NULL COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码(加密存储)',
    description TEXT COMMENT '描述信息',
    status ENUM('online', 'offline', 'unknown') DEFAULT 'unknown' COMMENT '连接状态',
    last_test_time DATETIME COMMENT '最后测试时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME COMMENT '软删除时间',
    
    INDEX idx_type (type),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at),
    INDEX idx_name (name)
) COMMENT='数据库连接配置表';

-- 插入一些示例数据
INSERT INTO sys_databases (name, type, host, port, username, password, database_name, description, status) VALUES
('本地MySQL', 'mysql', 'localhost', 3306, 'root', 'password123', 'test_db', '本地测试数据库', 'online'),
('生产MySQL', 'mysql', '192.168.1.100', 3306, 'prod_user', 'prod_pass', 'production', '生产环境数据库', 'online'),
('Redis缓存', 'redis', '192.168.1.101', 6379, 'redis_user', 'redis_pass', '0', 'Redis缓存服务器', 'online');

-- 创建用户权限表（可选）
CREATE TABLE IF NOT EXISTS sys_users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100),
    role ENUM('admin', 'user') DEFAULT 'user',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) COMMENT='用户表';

-- 插入默认管理员用户
INSERT INTO sys_users (username, password, email, role) VALUES
('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKXISwKlSRnbAFAKdP5MzvFaWhfG', 'admin@example.com', 'admin');
-- 默认密码: admin123