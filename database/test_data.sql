-- 测试数据脚本
-- 用于开发和测试环境

-- 0. 清理旧数据（如果存在）
DELETE FROM audit_logs WHERE user_id = '1' AND user_name = 'admin' AND action IN ('登录系统', '查看应用列表', '创建应用');
DELETE FROM users WHERE username LIKE 'testuser%';
DELETE FROM apps WHERE app_id LIKE 'test_app_%';

-- 1. 添加测试应用
INSERT INTO apps (name, app_id, app_secret, description, icon, status, created_at, updated_at) VALUES
('测试APP-电商平台', 'test_app_001', 'secret_key_001_abcdefghijklmnopqrstuvwxyz', '这是一个电商平台测试应用，用于测试商品管理、订单处理等功能', 'https://via.placeholder.com/128/FF6B6B/FFFFFF?text=E-Shop', 1, NOW(), NOW()),
('测试APP-社交网络', 'test_app_002', 'secret_key_002_zyxwvutsrqponmlkjihgfedcba', '社交网络测试应用，包含用户关系、动态发布、消息推送等功能', 'https://via.placeholder.com/128/4ECDC4/FFFFFF?text=Social', 1, NOW(), NOW()),
('测试APP-在线教育', 'test_app_003', 'secret_key_003_1234567890abcdefghijklmnop', '在线教育平台，提供课程管理、作业提交、在线考试等功能', 'https://via.placeholder.com/128/95E1D3/FFFFFF?text=Edu', 1, NOW(), NOW()),
('测试APP-物流管理', 'test_app_004', 'secret_key_004_qwertyuiopasdfghjklzxcvbnm', '物流管理系统，支持订单跟踪、配送管理、仓储管理', 'https://via.placeholder.com/128/F38181/FFFFFF?text=Logistics', 1, NOW(), NOW()),
('测试APP-金融服务', 'test_app_005', 'secret_key_005_mnbvcxzlkjhgfdsapoiuytrewq', '金融服务应用，包含账户管理、交易记录、风险控制等模块', 'https://via.placeholder.com/128/AA96DA/FFFFFF?text=Finance', 1, NOW(), NOW());

-- 2. 添加测试用户（密码统一为 test123）
-- 密码哈希: $2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIr.oaLQKO
INSERT INTO users (username, password, email, phone, status, created_at, updated_at) VALUES
('testuser01', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIr.oaLQKO', 'testuser01@example.com', '13800138001', 1, NOW(), NOW()),
('testuser02', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIr.oaLQKO', 'testuser02@example.com', '13800138002', 1, NOW(), NOW()),
('testuser03', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIr.oaLQKO', 'testuser03@example.com', '13800138003', 1, NOW(), NOW()),
('testuser04', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIr.oaLQKO', 'testuser04@example.com', '13800138004', 1, NOW(), NOW()),
('testuser05', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIr.oaLQKO', 'testuser05@example.com', '13800138005', 1, NOW(), NOW()),
('testuser06', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIr.oaLQKO', 'testuser06@example.com', '13800138006', 1, NOW(), NOW()),
('testuser07', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIr.oaLQKO', 'testuser07@example.com', '13800138007', 1, NOW(), NOW()),
('testuser08', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIr.oaLQKO', 'testuser08@example.com', '13800138008', 1, NOW(), NOW()),
('testuser09', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIr.oaLQKO', 'testuser09@example.com', '13800138009', 1, NOW(), NOW()),
('testuser10', '$2b$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewY5GyYIr.oaLQKO', 'testuser10@example.com', '13800138010', 1, NOW(), NOW());

-- 3. 添加审计日志测试数据
INSERT INTO audit_logs (user_id, user_name, action, resource, resource_id, request_method, request_path, ip_address, user_agent, status_code, created_at) VALUES
('1', 'admin', '登录系统', '认证', '0', 'POST', '/api/v1/admin/login', '127.0.0.1', 'Mozilla/5.0', 200, DATE_SUB(NOW(), INTERVAL 1 DAY)),
('1', 'admin', '查看应用列表', '应用', '0', 'GET', '/api/v1/apps', '127.0.0.1', 'Mozilla/5.0', 200, DATE_SUB(NOW(), INTERVAL 1 DAY)),
('1', 'admin', '创建应用', '应用', '1', 'POST', '/api/v1/apps', '127.0.0.1', 'Mozilla/5.0', 201, DATE_SUB(NOW(), INTERVAL 23 HOUR)),
('1', 'admin', '更新应用', '应用', '1', 'PUT', '/api/v1/apps/1', '127.0.0.1', 'Mozilla/5.0', 200, DATE_SUB(NOW(), INTERVAL 20 HOUR)),
('1', 'admin', '查看用户列表', '用户', '0', 'GET', '/api/v1/users', '127.0.0.1', 'Mozilla/5.0', 200, DATE_SUB(NOW(), INTERVAL 18 HOUR)),
('1', 'admin', '创建用户', '用户', '2', 'POST', '/api/v1/users', '127.0.0.1', 'Mozilla/5.0', 201, DATE_SUB(NOW(), INTERVAL 15 HOUR)),
('1', 'admin', '查看模块列表', '模块', '0', 'GET', '/api/v1/modules', '127.0.0.1', 'Mozilla/5.0', 200, DATE_SUB(NOW(), INTERVAL 12 HOUR)),
('1', 'admin', '更新模块配置', '模块', '1', 'PUT', '/api/v1/modules/1', '127.0.0.1', 'Mozilla/5.0', 200, DATE_SUB(NOW(), INTERVAL 10 HOUR)),
('1', 'admin', '删除应用', '应用', '2', 'DELETE', '/api/v1/apps/2', '127.0.0.1', 'Mozilla/5.0', 200, DATE_SUB(NOW(), INTERVAL 8 HOUR)),
('1', 'admin', '查看审计日志', '审计', '0', 'GET', '/api/v1/audit-logs', '127.0.0.1', 'Mozilla/5.0', 200, DATE_SUB(NOW(), INTERVAL 5 HOUR)),
('1', 'admin', '导出数据', '系统', '0', 'POST', '/api/v1/system/export', '127.0.0.1', 'Mozilla/5.0', 200, DATE_SUB(NOW(), INTERVAL 3 HOUR)),
('1', 'admin', '修改密码', '认证', '1', 'PUT', '/api/v1/admin/password', '127.0.0.1', 'Mozilla/5.0', 200, DATE_SUB(NOW(), INTERVAL 2 HOUR)),
('1', 'admin', '登录系统', '认证', '0', 'POST', '/api/v1/admin/login', '127.0.0.1', 'Mozilla/5.0', 200, DATE_SUB(NOW(), INTERVAL 1 HOUR)),
('1', 'admin', '查看仪表盘', '仪表盘', '0', 'GET', '/api/v1/dashboard', '127.0.0.1', 'Mozilla/5.0', 200, DATE_SUB(NOW(), INTERVAL 30 MINUTE)),
('1', 'admin', '登录系统', '认证', '0', 'POST', '/api/v1/admin/login', '127.0.0.1', 'Mozilla/5.0', 200, NOW());



-- 查询统计信息
SELECT '=== 测试数据统计 ===' as info;
SELECT 'Apps' as table_name, COUNT(*) as count FROM apps;
SELECT 'Users' as table_name, COUNT(*) as count FROM users;
SELECT 'Audit Logs' as table_name, COUNT(*) as count FROM audit_logs;
