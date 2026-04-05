# Flaredesk

**Cloudflare DNS 批量管理面板** — 为多站站长、出海团队、大量域名运维场景而生。

[![Telegram](https://img.shields.io/badge/Telegram-FlaredeskCommunity-blue?logo=telegram)](https://t.me/FlaredeskCommunity)
[![GitHub Stars](https://img.shields.io/github/stars/urbanu619/flaredesk?style=social)](https://github.com/urbanu619/flaredesk)

做 CF 控制台做不到的事：多账号统一管理、跨域名批量操作、一键橙云代理切换。

**你的数据只在你自己的服务器上。** API Token 不经过任何第三方。

> 如果这个项目对你有帮助，欢迎点一个 ⭐ Star，让更多人发现它。

---

## 功能

- **多账号管理** — 添加多个 CF 账号，统一面板操作
- **Zone 列表** — 一键从 CF 同步所有域名到本地，快速检索
- **DNS 管理** — 查看、新增、编辑、删除 DNS 记录
- **批量新增** — 同一域名批量添加多条记录
- **跨域名批量新增** — 多个域名同时添加相同记录
- **跨域名批量删除** — 按类型/名称批量清理
- **跨域名橙云切换** — 一键开启/关闭多个域名的 CF 代理
- **DNS 模板** — 保存常用记录组合，一键应用到多个域名
- **MCP Server** — 支持在 Claude / Cursor 等 AI 工具中直接调用

---

## 环境要求

- Go 1.21+
- Node.js 18+
- MySQL 5.7+
- Redis 6+

---

## 部署步骤

### 1. 克隆代码

```bash
git clone https://github.com/urbanu619/flaredesk.git
cd flaredesk
```

### 2. 配置后端

编辑 `magic_admin/conf.d/config.json`，修改以下字段：

```json
{
  "mysql": {
    "path": "127.0.0.1",
    "port": "3306",
    "db-name": "flaredesk",
    "username": "root",
    "password": "你的MySQL密码"
  },
  "redis": {
    "addr": "127.0.0.1:6379",
    "password": ""
  },
  "jwt": {
    "signing-key": "改成一个随机字符串"
  }
}
```

> MySQL 需要提前创建数据库：`CREATE DATABASE flaredesk DEFAULT CHARACTER SET utf8mb4;`
>
> 表结构会在首次启动时自动创建，无需手动执行 SQL。

### 3. 启动后端

```bash
cd magic_admin
go run main.go api
```

首次启动会自动：
- 创建所有数据表
- 初始化菜单和角色
- 创建默认管理员账号

默认账号：`superman` / `666666`（**请登录后立即修改密码**）

后端监听端口：`2022`

### 4. 构建前端

```bash
cd magic_admin_web
npm install
npm run build
```

构建产物在 `magic_admin_web/dist/`，用 Nginx 托管即可。

### 5. 配置 Nginx

```nginx
server {
    listen 80;
    server_name 你的域名或IP;

    # 前端静态文件
    root /path/to/flaredesk/magic_admin_web/dist;
    index index.html;

    # 前端路由
    location / {
        try_files $uri $uri/ /index.html;
    }

    # 后端 API 代理
    location /admin/api/ {
        proxy_pass http://127.0.0.1:2022;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

## 使用

1. 打开浏览器访问部署地址
2. 用 `superman` / `666666` 登录
3. 进入「Cloudflare → 账号管理」，添加 CF 账号和 API Token
4. 进入「Cloudflare → Zone 列表」，点击「从 CF 同步」导入域名
5. 开始批量管理 DNS 记录

### 如何获取 CF API Token

1. 登录 [Cloudflare Dashboard](https://dash.cloudflare.com/profile/api-tokens)
2. 创建 Token，权限选择：`Zone:DNS:Edit`
3. 将 Token 和 Account ID 填入账号管理页面

---

## MCP Server（AI 工具集成）

[![npm version](https://img.shields.io/npm/v/flaredesk-mcp)](https://www.npmjs.com/package/flaredesk-mcp)

在 Claude Desktop / Cursor / Windsurf 等支持 MCP 的 AI 工具中，用自然语言直接管理 Cloudflare DNS。

**无需部署 flaredesk，只需要 CF API Token。**

### 安装

```bash
npx flaredesk-mcp
```

### 配置 Claude Desktop

编辑 `~/Library/Application Support/Claude/claude_desktop_config.json`（Mac）：

```json
{
  "mcpServers": {
    "flaredesk": {
      "command": "npx",
      "args": ["flaredesk-mcp"],
      "env": {
        "CF_API_TOKEN": "你的 Cloudflare API Token"
      }
    }
  }
}
```

### 配置 Cursor / Windsurf

在 MCP 设置中添加同样的配置即可。

### 支持的操作

| 工具 | 说明 |
|------|------|
| `list_zones` | 列出账号下所有域名 |
| `list_dns_records` | 查看域名的 DNS 记录 |
| `create_dns_record` | 新增 DNS 记录 |
| `update_dns_record` | 修改 DNS 记录 |
| `delete_dns_record` | 删除 DNS 记录 |
| `batch_create_dns_records` | 批量新增记录（单域名） |
| `cross_zone_create_dns_records` | 跨域名批量新增 |
| `cross_zone_delete_records` | 跨域名批量删除 |
| `cross_zone_toggle_proxy` | 跨域名橙云代理批量切换 |

### 使用示例

> "帮我把 example.com 的所有 A 记录指向 1.2.3.4"
>
> "列出我账号下所有域名"
>
> "把 zone id xxx 下的 www 记录的橙云代理开启"

---

## 安全说明

- 所有数据（API Token、域名信息）存储在**你自己的服务器**上
- 建议部署后修改默认密码，并配置 HTTPS
- JWT 密钥请替换为随机字符串

---

## Built With

**后端**
- [Gin](https://github.com/gin-gonic/gin) — Go HTTP 框架
- [GORM](https://gorm.io) — Go ORM
- [Viper](https://github.com/spf13/viper) — 配置管理
- [zap](https://github.com/uber-go/zap) — 高性能日志

**前端**
- [Vue 3](https://vuejs.org) — 前端框架
- [Element Plus](https://element-plus.org) — UI 组件库
- [Vite](https://vitejs.dev) — 构建工具
- [Pinia](https://pinia.vuejs.org) — 状态管理
- [Vue Router](https://router.vuejs.org) — 路由

**MCP Server**
- [Model Context Protocol SDK](https://github.com/modelcontextprotocol/typescript-sdk) — MCP TypeScript SDK
- [Zod](https://github.com/colinhacks/zod) — 参数校验

---

## License

MIT
