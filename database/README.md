# 数据库说明

## 数据库信息

- **数据库名**: app_platform
- **字符集**: utf8mb4
- **排序规则**: utf8mb4_general_ci
- **引擎**: InnoDB

## 数据表列表

系统共包含 **16张数据表**：

| 表名 | 说明 | 主要字段 |
|------|------|---------|
| users | 用户表 | id, username, password, email, phone, status |
| apps | APP应用表 | id, name, app_id, app_secret, description, owner_id |
| module_templates | 模块模板表 | id, name, code, description, category |
| app_modules | APP模块关联表 | id, app_id, module_id, status, config |
| configs | 配置表 | id, app_id, key, value, env, version, status |
| versions | 版本管理表 | id, app_id, version_code, platform, download_url |
| files | 文件存储表 | id, app_id, filename, file_path, file_size |
| messages | 消息表 | id, app_id, user_id, title, content, status |
| logs | 日志表 | id, app_id, level, message, context |
| push_records | 推送记录表 | id, app_id, title, content, status, sent_count |
| analytics_events | 数据埋点表 | id, app_id, event_name, event_data, user_id |
| monitor_metrics | 监控指标表 | id, app_id, metric_name, metric_value, tags |
| alert_rules | 告警规则表 | id, app_id, name, metric_name, condition, threshold |
| alert_records | 告警记录表 | id, rule_id, app_id, message, level, status |
| audit_logs | 审计日志表 | id, user_id, app_id, action, resource, ip_address |
| app_users | APP用户表 | id, app_id, user_id, username, phone, email |

## 初始化数据

### 管理员账号

- **用户名**: admin
- **密码**: admin123
- **邮箱**: admin@example.com

⚠️ **重要**: 生产环境请务必修改默认密码！

### 模块模板

系统预置了11个功能模块：

1. **配置管理** (config_management) - 多环境配置管理
2. **埋点服务** (analytics) - 事件上报、数据统计
3. **文件存储** (file_storage) - 文件上传下载
4. **日志服务** (log_service) - 日志查询统计
5. **消息中心** (message_center) - 站内消息
6. **监控服务** (monitor_service) - 性能监控告警
7. **推送服务** (push_service) - APP推送
8. **用户管理** (user_management) - 用户列表管理
9. **版本管理** (version_management) - 版本发布
10. **WebSocket服务** (websocket_service) - 实时通信
11. **审计日志** (audit_log) - 操作审计

## 使用方法

### 1. 创建数据库

```sql
CREATE DATABASE IF NOT EXISTS app_platform CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;
USE app_platform;
```

### 2. 导入表结构

```bash
mysql -h your-host -u your-username -p app_platform < schema.sql
```

或在MySQL客户端中：

```sql
SOURCE /path/to/schema.sql;
```

### 3. 验证导入

```sql
-- 查看所有表
SHOW TABLES;

-- 查看管理员账号
SELECT * FROM users WHERE username='admin';

-- 查看模块模板
SELECT * FROM module_templates;
```

## 数据库配置

后端配置文件位置：`backend/configs/config.yaml`

```yaml
database:
  driver: mysql
  host: your-database-host
  port: 3306
  username: your-username
  password: your-password
  database: app_platform
  charset: utf8mb4
  parse_time: true
  loc: Local
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 3600
```

## 索引说明

所有表都包含以下索引优化：

- **主键索引**: 所有表的 `id` 字段
- **唯一索引**: 用户名、APP ID、模块代码等唯一标识
- **普通索引**: 外键、状态字段、时间字段
- **软删除索引**: `deleted_at` 字段（GORM软删除支持）

## 注意事项

1. **软删除**: 大部分表使用 `deleted_at` 字段实现软删除，删除的记录不会真正从数据库中移除
2. **时间戳**: 使用 `created_at` 和 `updated_at` 自动维护创建和更新时间
3. **字符集**: 使用 utf8mb4 支持完整的 Unicode 字符（包括 Emoji）
4. **外键**: 使用逻辑外键而非物理外键，提高性能和灵活性

## 备份建议

1. **定期备份**: 建议每天自动备份数据库
2. **保留周期**: 至少保留30天的备份
3. **备份命令**:

```bash
mysqldump -h host -u username -p app_platform > backup_$(date +%Y%m%d).sql
```

4. **恢复命令**:

```bash
mysql -h host -u username -p app_platform < backup_20260111.sql
```

## 性能优化建议

1. **日志表清理**: 定期清理 `logs` 和 `analytics_events` 表的历史数据
2. **索引优化**: 根据实际查询情况添加复合索引
3. **分表策略**: 数据量大时考虑按时间分表（日志、埋点等表）
4. **读写分离**: 生产环境建议配置主从复制
5. **连接池**: 合理配置数据库连接池参数

---

**最后更新**: 2026年1月11日
