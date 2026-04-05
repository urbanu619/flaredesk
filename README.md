# Flaredesk

**Cloudflare DNS 批量管理面板** — 为多站站长、出海团队、博彩/游戏行业运维而生。

做 CF 控制台做不到的事：多账号统一管理、跨域名批量操作、一键橙云代理切换。

**你的数据只在你自己的服务器上。** API Token 不经过任何第三方。

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

在 Claude Desktop / Cursor 等支持 MCP 的 AI 工具中直接管理 DNS：

```bash
cd flaredesk-mcp
npm install
npm run build
```

在 AI 工具的 MCP 配置中添加：

```json
{
  "mcpServers": {
    "flaredesk": {
      "command": "node",
      "args": ["/path/to/flaredesk/flaredesk-mcp/dist/index.js"],
      "env": {
        "CF_API_TOKEN": "你的CF API Token"
      }
    }
  }
}
```

支持的操作：列出域名、查询/新增/删除 DNS 记录、批量操作、跨域名橙云切换。

---

## 安全说明

- 所有数据（API Token、域名信息）存储在**你自己的服务器**上
- 建议部署后修改默认密码，并配置 HTTPS
- JWT 密钥请替换为随机字符串

---

## License

MIT
