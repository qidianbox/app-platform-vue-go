-- APP中台管理系统数据库表结构
-- 数据库: app_platform
-- 字符集: utf8mb4

-- 1. 用户表
CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL COMMENT '用户名',
  `password` varchar(255) NOT NULL COMMENT '密码hash',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `status` tinyint(4) DEFAULT '1' COMMENT '状态 1:正常 2:禁用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  KEY `idx_email` (`email`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 2. APP应用表
CREATE TABLE IF NOT EXISTS `apps` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT 'APP名称',
  `app_id` varchar(50) NOT NULL COMMENT 'APP唯一标识',
  `app_secret` varchar(255) NOT NULL COMMENT 'APP密钥',
  `description` text COMMENT 'APP描述',
  `icon` varchar(255) DEFAULT NULL COMMENT 'APP图标URL',
  `status` tinyint(4) DEFAULT '1' COMMENT '状态 1:正常 2:禁用',
  `owner_id` bigint(20) unsigned DEFAULT NULL COMMENT '所有者ID',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_app_id` (`app_id`),
  KEY `idx_owner_id` (`owner_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='APP应用表';

-- 3. 模块模板表
CREATE TABLE IF NOT EXISTS `module_templates` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '模块名称',
  `code` varchar(50) NOT NULL COMMENT '模块代码',
  `description` text COMMENT '模块描述',
  `category` varchar(50) DEFAULT NULL COMMENT '模块分类',
  `icon` varchar(255) DEFAULT NULL COMMENT '模块图标',
  `sort_order` int(11) DEFAULT '0' COMMENT '排序',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code` (`code`),
  KEY `idx_category` (`category`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='模块模板表';

-- 4. APP模块关联表
CREATE TABLE IF NOT EXISTS `app_modules` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `app_id` bigint(20) unsigned NOT NULL COMMENT 'APP ID',
  `module_id` bigint(20) unsigned NOT NULL COMMENT '模块ID',
  `status` tinyint(4) DEFAULT '1' COMMENT '状态 1:启用 2:禁用',
  `config` text COMMENT '模块配置JSON',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_app_module` (`app_id`,`module_id`),
  KEY `idx_module_id` (`module_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='APP模块关联表';

-- 5. 配置表
CREATE TABLE IF NOT EXISTS `configs` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `app_id` bigint(20) unsigned NOT NULL COMMENT 'APP ID',
  `key` varchar(100) NOT NULL COMMENT '配置键',
  `value` text COMMENT '配置值',
  `env` varchar(20) DEFAULT 'prod' COMMENT '环境 dev/test/prod',
  `version` int(11) DEFAULT '1' COMMENT '版本号',
  `status` tinyint(4) DEFAULT '1' COMMENT '状态 1:草稿 2:已发布',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_app_key_env` (`app_id`,`key`,`env`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='配置表';

-- 6. 版本管理表
CREATE TABLE IF NOT EXISTS `versions` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `app_id` bigint(20) unsigned NOT NULL COMMENT 'APP ID',
  `version_code` varchar(50) NOT NULL COMMENT '版本号',
  `version_name` varchar(100) DEFAULT NULL COMMENT '版本名称',
  `platform` varchar(20) DEFAULT NULL COMMENT '平台 ios/android/web',
  `download_url` varchar(500) DEFAULT NULL COMMENT '下载地址',
  `update_log` text COMMENT '更新日志',
  `force_update` tinyint(4) DEFAULT '0' COMMENT '是否强制更新',
  `status` tinyint(4) DEFAULT '1' COMMENT '状态 1:测试 2:灰度 3:正式 4:下线',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_app_id` (`app_id`),
  KEY `idx_platform` (`platform`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='版本管理表';

-- 7. 文件存储表
CREATE TABLE IF NOT EXISTS `files` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `app_id` bigint(20) unsigned NOT NULL COMMENT 'APP ID',
  `filename` varchar(255) NOT NULL COMMENT '文件名',
  `file_path` varchar(500) NOT NULL COMMENT '文件路径',
  `file_size` bigint(20) DEFAULT NULL COMMENT '文件大小(字节)',
  `mime_type` varchar(100) DEFAULT NULL COMMENT 'MIME类型',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_app_id` (`app_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文件存储表';

-- 8. 消息表
CREATE TABLE IF NOT EXISTS `messages` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `app_id` bigint(20) unsigned NOT NULL COMMENT 'APP ID',
  `user_id` bigint(20) unsigned DEFAULT NULL COMMENT '接收用户ID',
  `title` varchar(200) NOT NULL COMMENT '消息标题',
  `content` text COMMENT '消息内容',
  `type` varchar(20) DEFAULT 'system' COMMENT '消息类型',
  `status` tinyint(4) DEFAULT '0' COMMENT '状态 0:未读 1:已读',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `read_at` timestamp NULL DEFAULT NULL COMMENT '阅读时间',
  PRIMARY KEY (`id`),
  KEY `idx_app_id` (`app_id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='消息表';

-- 9. 日志表
CREATE TABLE IF NOT EXISTS `logs` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `app_id` bigint(20) unsigned NOT NULL COMMENT 'APP ID',
  `level` varchar(20) DEFAULT 'info' COMMENT '日志级别',
  `message` text COMMENT '日志消息',
  `context` text COMMENT '上下文JSON',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_app_id` (`app_id`),
  KEY `idx_level` (`level`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='日志表';

-- 10. 推送记录表
CREATE TABLE IF NOT EXISTS `push_records` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `app_id` bigint(20) unsigned NOT NULL COMMENT 'APP ID',
  `title` varchar(200) NOT NULL COMMENT '推送标题',
  `content` text COMMENT '推送内容',
  `target_type` varchar(20) DEFAULT 'all' COMMENT '推送目标类型',
  `target_ids` text COMMENT '目标ID列表JSON',
  `status` tinyint(4) DEFAULT '0' COMMENT '状态 0:待发送 1:发送中 2:已完成 3:失败',
  `sent_count` int(11) DEFAULT '0' COMMENT '已发送数量',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `sent_at` timestamp NULL DEFAULT NULL COMMENT '发送时间',
  PRIMARY KEY (`id`),
  KEY `idx_app_id` (`app_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='推送记录表';

-- 11. 数据埋点表
CREATE TABLE IF NOT EXISTS `analytics_events` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `app_id` bigint(20) unsigned NOT NULL COMMENT 'APP ID',
  `event_name` varchar(100) NOT NULL COMMENT '事件名称',
  `event_data` text COMMENT '事件数据JSON',
  `user_id` bigint(20) unsigned DEFAULT NULL COMMENT '用户ID',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_app_id` (`app_id`),
  KEY `idx_event_name` (`event_name`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据埋点表';

-- 12. 监控指标表
CREATE TABLE IF NOT EXISTS `monitor_metrics` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `app_id` bigint(20) unsigned NOT NULL COMMENT 'APP ID',
  `metric_name` varchar(100) NOT NULL COMMENT '指标名称',
  `metric_value` decimal(20,4) DEFAULT NULL COMMENT '指标值',
  `tags` text COMMENT '标签JSON',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_app_id` (`app_id`),
  KEY `idx_metric_name` (`metric_name`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='监控指标表';

-- 13. 告警规则表
CREATE TABLE IF NOT EXISTS `alert_rules` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `app_id` bigint(20) unsigned NOT NULL COMMENT 'APP ID',
  `name` varchar(100) NOT NULL COMMENT '规则名称',
  `metric_name` varchar(100) NOT NULL COMMENT '监控指标',
  `condition` varchar(50) DEFAULT NULL COMMENT '条件 gt/lt/eq',
  `threshold` decimal(20,4) DEFAULT NULL COMMENT '阈值',
  `status` tinyint(4) DEFAULT '1' COMMENT '状态 1:启用 2:禁用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_app_id` (`app_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='告警规则表';

-- 14. 告警记录表
CREATE TABLE IF NOT EXISTS `alert_records` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `rule_id` bigint(20) unsigned NOT NULL COMMENT '规则ID',
  `app_id` bigint(20) unsigned NOT NULL COMMENT 'APP ID',
  `message` text COMMENT '告警消息',
  `level` varchar(20) DEFAULT 'warning' COMMENT '告警级别',
  `status` tinyint(4) DEFAULT '0' COMMENT '状态 0:未处理 1:已处理',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `handled_at` timestamp NULL DEFAULT NULL COMMENT '处理时间',
  PRIMARY KEY (`id`),
  KEY `idx_rule_id` (`rule_id`),
  KEY `idx_app_id` (`app_id`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='告警记录表';

-- 15. 审计日志表
CREATE TABLE IF NOT EXISTS `audit_logs` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned DEFAULT NULL COMMENT '操作用户ID',
  `app_id` bigint(20) unsigned DEFAULT NULL COMMENT 'APP ID',
  `action` varchar(100) NOT NULL COMMENT '操作动作',
  `resource` varchar(100) DEFAULT NULL COMMENT '操作资源',
  `resource_id` bigint(20) unsigned DEFAULT NULL COMMENT '资源ID',
  `ip_address` varchar(50) DEFAULT NULL COMMENT 'IP地址',
  `user_agent` varchar(500) DEFAULT NULL COMMENT 'User Agent',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_app_id` (`app_id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='审计日志表';

-- 16. APP用户关联表
CREATE TABLE IF NOT EXISTS `app_users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `app_id` bigint(20) unsigned NOT NULL COMMENT 'APP ID',
  `user_id` varchar(100) NOT NULL COMMENT '用户ID（APP内部）',
  `username` varchar(100) DEFAULT NULL COMMENT '用户名',
  `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
  `email` varchar(100) DEFAULT NULL COMMENT '邮箱',
  `status` tinyint(4) DEFAULT '1' COMMENT '状态 1:正常 2:禁用',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_app_user` (`app_id`,`user_id`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='APP用户表';

-- 初始化管理员账号
INSERT INTO `users` (`username`, `password`, `email`, `status`) VALUES
('admin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iAt6Z5EH', 'admin@example.com', 1);
-- 密码: admin123

-- 初始化模块模板
INSERT INTO `module_templates` (`name`, `code`, `description`, `category`, `sort_order`) VALUES
('配置管理', 'config_management', '多环境配置管理、配置发布、配置历史', '基础服务', 1),
('埋点服务', 'analytics', '事件上报、数据统计、漏斗分析', '数据分析', 2),
('文件存储', 'file_storage', '文件上传、下载、管理', '基础服务', 3),
('日志服务', 'log_service', '日志查询、统计、导出', '运维监控', 4),
('消息中心', 'message_center', '站内消息、消息推送', '消息通知', 5),
('监控服务', 'monitor_service', '性能监控、告警管理', '运维监控', 6),
('推送服务', 'push_service', 'APP推送、批量推送', '消息通知', 7),
('用户管理', 'user_management', '用户列表、用户状态管理', '用户权限', 8),
('版本管理', 'version_management', '版本发布、灰度发布、强制更新', '基础服务', 9),
('WebSocket服务', 'websocket_service', '实时通信、数据推送', '基础服务', 10),
('审计日志', 'audit_log', '操作审计、日志查询', '安全审计', 11);
