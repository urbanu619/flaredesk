# Flaredesk — Beginner install guide (Windows / macOS)

**English** · [简体中文](完全新手安装指南.md)

For readers who are **new to the command line** (no coding experience required).  
This guide uses plain language: what to do each step and what you should **normally** see. If something differs, copy the **full error text** when asking for help—it makes debugging much easier.

---

## Pick a path first

| What you want | Difficulty | Notes |
|----------------|------------|--------|
| **Manage Cloudflare DNS in an AI app** (e.g. Cursor / Claude) with natural language | Easiest | **No** need to install this repo or a database. Install **Node** and add one MCP config. |
| **Use the full web panel** (multi-account, bulk DNS, etc.) | Harder | Install **Go** and **Node** on your machine. You can use **SQLite** (no MySQL) and **embedded Redis** (no Redis server). MySQL + separate Redis adds more steps. |

**Suggestion:** If this is your first time, finish path **A (MCP)** first to confirm your Cloudflare token works, then decide whether to run the full panel.

The guide has two parts: **A. MCP only** and **B. Full web UI (this repo)**.

---

## Quick glossary

- **Terminal / shell**: A window where you type commands. On Windows use **PowerShell** or **cmd**; on macOS **Terminal**.
- **Path / folder**: The same folders you see in Explorer / Finder. `cd some/path` means “go into that folder before running the next commands”.
- **Paste**: In terminals, **right‑click** to paste on Windows, or **Cmd + V** on macOS Terminal (usually).

---

## A. MCP only (no Flaredesk web app)

**Result:** In a supported AI client you can ask things like “list my zones” or “add an A record for …”.

### A1. Install Node.js (both OSes)

1. Open https://nodejs.org/  
2. Download the **LTS** build and install with the defaults.  
3. **Verify:**

**Windows:** Open **PowerShell** (search in Start), run:

```text
node -v
```

Press Enter.  
**Expected:** A line like `v20.x.x`.

**macOS:** Open **Terminal** (Spotlight: `Terminal`), same command `node -v`.  
**Expected:** Same as above.

If you see “not recognized” or `command not found`, Node is missing or not on `PATH`—close the terminal, reopen, or reinstall Node.

### A2. Configure MCP in your AI app

Follow the **MCP Server** section in the repo root **[README.md](../README.md)** (default English); Chinese mirror: **[README.zh-CN.md](../README.zh-CN.md)** under **「MCP Server（AI 工具集成）」**.

You need a **Cloudflare API token** with permission to edit **Zone → DNS** (at minimum).

**Expected:** After saving config and restarting the client, asking it to run something like `list_zones` returns your zone list.

---

## B. Full web panel (run backend + frontend locally)

**Result:** Log into Flaredesk in the browser and use multi-account and bulk DNS features.  
**Trade-off:** More steps; you need **Go** and **Node**. Optional **SQLite + embedded Redis** avoids installing MySQL/Redis as services.

### B1. Get the code on your machine

**Option A (if you use Git):** clone the repo.  
**Option B (no Git):** On GitHub use **Code → Download ZIP** and extract to e.g. a `flaredesk` folder on your Desktop.

**Expected:** You see folders such as `magic_admin` and `magic_admin_web`.

### B1a. (Recommended) Minimal setup: SQLite + embedded Redis

**You skip:** Installing MySQL, installing/running Redis as a service, and manual `CREATE DATABASE`.

1. In Explorer / Finder open `magic_admin/conf.d/`.  
2. Copy **`config.sqlite.embedded.json.example`** and rename the copy to **`config.json`** (back up any existing `config.json` first, or merge fields).  
3. In JSON, ensure **`mysql`** has `"driver": "sqlite"`, **`path`** is a data directory (e.g. `"./data"`), **`db-name`** becomes `./data/<name>.db`, and **`redis`** has `"embedded": true`.

**Expected:** `magic_admin/conf.d/config.json` exists and `jwt.signing-key` is a random string you chose (see **[README.md](../README.md)**).

If you finished **B1a**, jump to **B5 Start the backend** (skip **B2**, **B3**; **B4** is optional—just verify `jwt`).

---

### B2. Install tooling (only for MySQL + separate Redis)

After each install, **close and reopen** the terminal, then run the checks.

| Tool | Role | Notes |
|------|------|--------|
| **Go** | Backend | https://go.dev/dl/ |
| **Node.js** | Frontend build | https://nodejs.org/ LTS |
| **MySQL** | Database | Windows: MySQL Installer; macOS: installer or Homebrew |
| **Redis** | Cache | Windows: official or WSL; macOS: `brew install redis`, etc. |

**Check commands:**

```text
go version
node -v
mysql --version
redis-cli ping
```

**Expected:**  
- `go version` shows Go **1.21+**.  
- `node -v` prints a version.  
- MySQL/Redis reachable (`redis-cli ping` → `PONG`).

Fix any failing tool before continuing.

> **If you use B1a (SQLite):** skip B2 and B3.

### B3. Create the database (MySQL only)

```sql
CREATE DATABASE flaredesk DEFAULT CHARACTER SET utf8mb4;
```

**Expected:** Command succeeds with no error.

> **If you use B1a (SQLite):** do not run this SQL.

### B4. Edit config

Open with Notepad (Windows) or TextEdit (macOS):

`magic_admin/conf.d/config.json`

- **After B1a:** mainly confirm **`jwt.signing-key`** is random.  
- **For MySQL:** fill **`mysql`**, **`redis`** (or `redis.embedded: true` to skip a Redis server), **`jwt`** per **[README.md](../README.md)** → **Deployment → Configure backend**; use your local MySQL password.

**Expected:** File saves; JSON is valid (matching quotes/braces, no trailing comma after the last field).

### B5. Start the backend

**Windows (PowerShell):**

```powershell
cd path\to\flaredesk\magic_admin
go run main.go api
```

**macOS (Terminal):**

```bash
cd ~/path/to/flaredesk/magic_admin
go run main.go api
```

**Expected:** Logs keep printing; the process stays running; near the end you should see port **2022**.  
**On error:** Copy **everything** from `go run` through the error message.

First successful start creates tables and the default admin (see **[README.md](../README.md)**).

### B6. Build the frontend

**Open a second terminal** (leave the backend running).

**Windows:**

```powershell
cd path\to\flaredesk\magic_admin_web
npm install
npm run build
```

**macOS:**

```bash
cd ~/path/to/flaredesk/magic_admin_web
npm install
npm run build
```

**Expected:** Build finishes successfully; folder `magic_admin_web/dist` exists.

### B7. Open in the browser (quick local test without Nginx)

If you don’t have Nginx yet, you can **temporarily** serve the built files (local testing only):

```bash
cd magic_admin_web/dist
npx --yes serve -s . -l 3000
```

Then visit `http://127.0.0.1:3000`.

**Important:** The UI calls `/admin/api/`. You must **proxy** `/admin/api/` to `http://127.0.0.1:2022` with **Nginx** or similar—see the Nginx example in **[README.md](../README.md)**. Without that, login fails. If proxies are new to you, this is the most common beginner blocker—use the README sample or ask someone to set the proxy once.

**Expected:** Page loads; with a correct proxy, sign in with the default account (**change the password immediately**).

---

## How to ask for help effectively

1. **OS:** Windows 11 or which macOS version.  
2. **Path A (MCP) or B (full panel).**  
3. **Last command you ran** (copy/paste).  
4. **Full terminal output** including the error.  
5. **What you expected vs what happened** (one sentence).

---

## Relation to “one-click installers”

The easiest experience for non-technical users is an **official GUI installer** (or single `.exe` / `.app`) with bundled DB and browser—product work. This repo already supports **SQLite + embedded Redis** to skip two external services; shipping prebuilt binaries could shorten steps further.

Until then, this guide splits **MCP first, panel second** and spells out **expected output** to reduce trial-and-error. If you document your own walkthrough with screenshots, it will help other beginners too.
