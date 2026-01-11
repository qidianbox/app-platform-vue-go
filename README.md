# APP中台管理系统

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org)
[![Vue Version](https://img.shields.io/badge/Vue-3.5+-4FC08D?logo=vue.js)](https://vuejs.org)

统一的企业级应用管理平台，提供APP管理、用户管理、配置管理、监控告警等完整功能。

## ✨ 特性

- 🚀 **现代技术栈**: Vue 3 + Element Plus + Go + Gin
- 📦 **模块化设计**: 11个独立功能模块，按需启用
- 🔐 **安全可靠**: JWT认证、权限管理、操作审计
- 📊 **数据分析**: 埋点统计、监控告警、日志分析
- 🎨 **企业级UI**: 响应式设计，支持PC和移动端
- 🔧 **易于部署**: 支持Docker、K8s等多种部署方式

## 📋 功能模块

### 核心功能
- ✅ **认证授权** - JWT登录、权限管理、操作审计
- ✅ **APP管理** - 应用创建、密钥管理、模块配置
- ✅ **用户管理** - 用户列表、状态管理、搜索筛选
- ✅ **配置中心** - 多环境配置、版本管理、配置发布
- ✅ **版本管理** - 版本发布、灰度发布、强制更新

### 11个业务模块

| 模块 | 功能 | API数量 |
|------|------|---------|
| 配置管理 | 配置列表、创建/更新/发布、配置历史 | 5个 |
| 埋点服务 | 事件上报、批量上报、事件统计、漏斗分析 | 6个 |
| 文件存储 | 文件上传/下载、文件列表、存储统计 | 5个 |
| 日志服务 | 日志查询、上报、统计、导出、清理 | 5个 |
| 消息中心 | 发送消息、消息列表、消息模板、批量发送 | 6个 |
| 监控服务 | 上报指标、监控指标、告警管理、健康检查 | 5个 |
| 推送服务 | 创建推送、发送推送、推送统计、推送模板 | 6个 |
| 用户管理 | 用户列表、用户详情、状态管理、用户统计 | 4个 |
| 版本管理 | 版本列表、创建/发布/下线、更新检查 | 6个 |
| WebSocket服务 | 实时连接、数据推送、告警推送 | 3个 |
| 审计日志 | 日志列表、审计统计、导出审计日志 | 3个 |

**总计**: 54个API接口

## 🚀 快速开始

### 环境要求

- **Node.js**: 18.0+
- **Go**: 1.21+
- **MySQL**: 8.0+
- **Redis**: 6.0+ (可选)

### 1. 克隆项目

```bash
git clone https://github.com/your-username/app-platform.git
cd app-platform
```

### 2. 初始化数据库

```bash
# 创建数据库
mysql -u root -p -e "CREATE DATABASE app_platform CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;"

# 导入表结构
mysql -u root -p app_platform < database/schema.sql
```

详细说明请查看 [database/README.md](database/README.md)

### 3. 配置后端

编辑 `backend/configs/config.yaml`:

```yaml
database:
  host: localhost
  port: 3306
  username: root
  password: your-password
  database: app_platform
```

### 4. 启动后端

```bash
cd backend
go mod download
go run cmd/server/main.go
```

后端将运行在 `http://localhost:8080`

### 5. 配置前端

创建 `frontend/.env`:

```
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

### 6. 启动前端

```bash
cd frontend
pnpm install
pnpm dev
```

前端将运行在 `http://localhost:5173`

### 7. 登录系统

- **地址**: http://localhost:5173
- **用户名**: admin
- **密码**: admin123

⚠️ **首次登录后请立即修改密码！**

## 📁 项目结构

```
app-platform/
├── frontend/              # Vue 3 前端
│   ├── src/
│   │   ├── views/        # 页面组件
│   │   ├── components/   # 公共组件
│   │   ├── router/       # 路由配置
│   │   ├── store/        # 状态管理
│   │   ├── api/          # API接口
│   │   └── utils/        # 工具函数
│   ├── public/           # 静态资源
│   └── package.json
├── backend/               # Go 后端
│   ├── cmd/              # 主程序入口
│   │   └── server/
│   ├── core/             # 核心模块
│   │   └── module/       # 模块系统
│   ├── models/           # 数据模型
│   ├── modules/          # 业务模块
│   │   ├── analytics/    # 埋点服务
│   │   ├── config/       # 配置管理
│   │   ├── file/         # 文件存储
│   │   ├── log/          # 日志服务
│   │   ├── message/      # 消息中心
│   │   ├── monitor/      # 监控服务
│   │   ├── push/         # 推送服务
│   │   └── ...
│   ├── middleware/       # 中间件
│   ├── utils/            # 工具函数
│   ├── configs/          # 配置文件
│   └── go.mod
├── database/              # 数据库
│   ├── schema.sql        # 表结构
│   └── README.md         # 数据库说明
├── docs/                  # 文档
├── deploy/                # 部署配置
└── README.md              # 本文档
```

## 🔧 开发指南

### 前端开发

```bash
cd frontend

# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev

# 构建生产版本
pnpm build

# 代码检查
pnpm lint
```

### 后端开发

```bash
cd backend

# 下载依赖
go mod download

# 运行开发服务器
go run cmd/server/main.go

# 构建生产版本
go build -o server cmd/server/main.go

# 运行测试
go test ./...
```

### 添加新模块

1. 在 `backend/modules/` 创建模块目录
2. 实现 `Module` 接口（参考 `core/module/module.go`）
3. 在 `modules/loader.go` 注册模块
4. 在前端添加对应的页面和API调用

## 📦 部署

### Docker部署

```bash
# 构建镜像
docker-compose build

# 启动服务
docker-compose up -d
```

### 生产部署

详细部署文档请查看 [deploy/README.md](deploy/README.md)

支持的部署方式：
- Docker + Docker Compose
- Kubernetes (K8s)
- 阿里云SAE（Serverless应用引擎）
- 阿里云ACK（容器服务）
- 传统服务器部署

## 🔒 安全建议

1. **修改默认密码**: 首次登录后立即修改admin密码
2. **配置HTTPS**: 生产环境必须启用HTTPS
3. **数据库安全**: 
   - 使用强密码
   - 限制访问IP
   - 启用SSL连接
4. **环境变量**: 敏感信息使用环境变量而非配置文件
5. **定期更新**: 及时更新依赖包修复安全漏洞
6. **备份策略**: 定期备份数据库和配置文件

## 📊 技术栈

### 前端
- **框架**: Vue 3.5 (Composition API)
- **UI组件**: Element Plus 2.13
- **状态管理**: Pinia
- **路由**: Vue Router 4
- **HTTP客户端**: Axios
- **图表**: ECharts 5
- **构建工具**: Vite 5

### 后端
- **语言**: Go 1.21
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL 8.0
- **缓存**: Redis (可选)
- **认证**: JWT

## 📝 API文档

API文档位于 `docs/api/` 目录，或访问运行中的服务：

```
http://localhost:8080/swagger/index.html
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

[MIT License](LICENSE)

## 📧 联系方式

- **项目负责人**: [待填写]
- **技术支持**: [待填写]
- **问题反馈**: [GitHub Issues](https://github.com/your-username/app-platform/issues)

---

**最后更新**: 2026年1月11日
