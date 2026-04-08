# Flaredesk

**English** · [简体中文](README.zh-CN.md)

A **self-hosted web panel for Cloudflare DNS** at scale: multiple accounts, bulk edits across zones, orange-cloud toggles, and DNS templates—things the Cloudflare dashboard does not do well for many zones.

[![Website](https://img.shields.io/badge/Website-flaredesk--site.vercel.app-orange)](https://flaredesk-site.vercel.app)
[![Telegram](https://img.shields.io/badge/Telegram-FlaredeskCommunity-blue?logo=telegram)](https://t.me/FlaredeskCommunity)
[![GitHub Stars](https://img.shields.io/github/stars/urbanu619/flaredesk?style=social)](https://github.com/urbanu619/flaredesk)

**Your data stays on your own server.** API tokens never pass through a third party.

> If this project helps you, a ⭐ star helps others discover it.

**Quick try:** copy [`magic_admin/conf.d/config.sqlite.embedded.json.example`](magic_admin/conf.d/config.sqlite.embedded.json.example) to `magic_admin/conf.d/config.json` to run with **SQLite + embedded Redis** without installing MySQL or Redis separately. Then follow [Deployment](#deployment) below. **New to the CLI?** Step-by-step guides: [`docs/beginner-install.md`](docs/beginner-install.md) · [`docs/完全新手安装指南.md`](docs/完全新手安装指南.md) (Chinese).

---

## Features

- **Multiple Cloudflare accounts** in one UI
- **Zone list** — sync all zones from CF for quick search
- **DNS CRUD** for records
- **Bulk add** — many records on one zone
- **Cross-zone bulk add** — same records to many zones at once
- **Cross-zone bulk delete** — filter by type/name
- **Cross-zone orange-cloud toggle** — proxy on/off in bulk
- **DNS templates** — save record sets and apply to new zones
- **MCP server** — use Claude, Cursor, and other MCP-capable tools

---

## Requirements

- Go 1.21+
- Node.js 18+

**Database & Redis (pick one mode)**

| Mode | Notes |
|------|--------|
| **Minimal local (recommended for a quick try)** | **SQLite** + **embedded Redis** (`redis.embedded: true`). No separate MySQL/Redis. Set `mysql.driver` to `sqlite`; `path` is the data directory and combines with `db-name` as `./data/<name>.db`. Copy [magic_admin/conf.d/config.sqlite.embedded.json.example](magic_admin/conf.d/config.sqlite.embedded.json.example) to `config.json`. |
| **Traditional** | MySQL 5.7+ and Redis 6+; fill DSN in config as below. |

> With SQLite, the “generate model from DB” feature is unavailable—use MySQL or maintain models by hand.

Beginner walkthroughs (Windows / macOS): use the **`docs/beginner-install.md`** and **`docs/完全新手安装指南.md`** links in the **Quick try** paragraph at the top (not repeated here).

### Embedded frontend (single binary, optional)

The Vite build can be **embedded** into the Go binary and served on the same origin as `/admin/api` (no separate Nginx or Vite in production for that mode).

1. Edit `magic_admin/conf.d/app.json`, set **`"serve-web": true`** (default `false`; local dev still uses `npm run dev`).
2. From repo root run **`./scripts/embed-web.sh`** (syncs `magic_admin_web/dist` into `magic_admin/webdist/dist` before compile).
3. `cd magic_admin && go build -o flaredesk .`, then open **`http://127.0.0.1:<addr>/`** (`<addr>` from `app.json`, default `2022`).

Without running the embed script, the repo shows a placeholder page until the frontend is built.

---

## Deployment

### 1. Clone

```bash
git clone https://github.com/urbanu619/flaredesk.git
cd flaredesk
```

### 2. Backend config

Edit `magic_admin/conf.d/config.json`, for example:

```json
{
  "mysql": {
    "path": "127.0.0.1",
    "port": "3306",
    "db-name": "flaredesk",
    "username": "root",
    "password": "YOUR_MYSQL_PASSWORD"
  },
  "redis": {
    "addr": "127.0.0.1:6379",
    "password": ""
  },
  "jwt": {
    "signing-key": "use-a-random-string"
  }
}
```

> **MySQL:** create the database first: `CREATE DATABASE flaredesk DEFAULT CHARACTER SET utf8mb4;`
>
> **SQLite:** no server DB; ensure `path` is writable (e.g. `./data`); the `db-name.db` file is created on first run.
>
> Schema is created on first startup—no manual SQL import.

### 3. Run backend

```bash
cd magic_admin
go run main.go api
```

First run will:

- Create tables
- Seed menus and roles
- Create the default admin user

Default login: `superman` / `666666` (**change the password immediately**)

API listens on port `2022` by default.

### 4. Build frontend

```bash
cd magic_admin_web
npm install
npm run build
```

Output: `magic_admin_web/dist/` — serve with Nginx or use embedded mode above.

### 5. Nginx example

```nginx
server {
    listen 80;
    server_name your.domain.or.ip;

    root /path/to/flaredesk/magic_admin_web/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /admin/api/ {
        proxy_pass http://127.0.0.1:2022;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

## Usage

1. Open the deployed URL in a browser
2. Sign in with `superman` / `666666` (then change password)
3. **Cloudflare → Accounts** — add CF account and API token
4. **Cloudflare → Zones** — **Sync from CF** to import zones
5. Manage DNS in bulk as needed

### Cloudflare API token

1. Open [Cloudflare API Tokens](https://dash.cloudflare.com/profile/api-tokens)
2. Create a token with at least `Zone` → `DNS` → `Edit` for bulk DNS.
3. For **Origin CA bulk certificate** flows, the token also needs `Zone` → `SSL and Certificates` → `Edit`. Otherwise Cloudflare may return `Authentication error` (10000)—meaning the token lacks that scope, not that it is invalid.
4. Paste the token and Account ID in the account UI

---

## MCP Server (AI tools)

[![npm version](https://img.shields.io/npm/v/flaredesk-mcp)](https://www.npmjs.com/package/flaredesk-mcp)

Use natural language in Claude Desktop, Cursor, Windsurf, or any MCP-capable client to manage Cloudflare DNS.

**You do not need to deploy Flaredesk for MCP—only a Cloudflare API token.**

### Install

```bash
npx flaredesk-mcp
```

### Claude Desktop

Edit `~/Library/Application Support/Claude/claude_desktop_config.json` (macOS):

```json
{
  "mcpServers": {
    "flaredesk": {
      "command": "npx",
      "args": ["flaredesk-mcp"],
      "env": {
        "CF_API_TOKEN": "YOUR_CLOUDFLARE_API_TOKEN"
      }
    }
  }
}
```

### Cursor / Windsurf

Add the same MCP server block in your MCP settings.

### Tools

| Tool | Purpose |
|------|---------|
| `list_zones` | List zones under the account |
| `list_dns_records` | List DNS records for a zone |
| `create_dns_record` | Create a record |
| `update_dns_record` | Update a record |
| `delete_dns_record` | Delete a record |
| `batch_create_dns_records` | Bulk create on one zone |
| `cross_zone_create_dns_records` | Bulk create across zones |
| `cross_zone_delete_records` | Bulk delete across zones |
| `cross_zone_toggle_proxy` | Bulk orange-cloud toggle |

### Example prompts

> “Point all A records on example.com to 1.2.3.4”
>
> “List all zones in my account”
>
> “Turn on orange cloud for www on zone xxx”

---

## Security

- All data (tokens, zone metadata) stays **on your server**
- Change the default password and use HTTPS in production
- Replace `jwt.signing-key` with a random secret

---

## Built With

**Backend**
- [Gin](https://github.com/gin-gonic/gin) — HTTP framework
- [GORM](https://gorm.io) — ORM
- [Viper](https://github.com/spf13/viper) — configuration
- [zap](https://github.com/uber-go/zap) — logging

**Frontend**
- [Vue 3](https://vuejs.org)
- [Element Plus](https://element-plus.org)
- [Vite](https://vitejs.dev)
- [Pinia](https://pinia.vuejs.org)
- [Vue Router](https://router.vuejs.org)

**MCP Server**
- [Model Context Protocol SDK](https://github.com/modelcontextprotocol/typescript-sdk)
- [Zod](https://github.com/colinhacks/zod)

---

## License

MIT
