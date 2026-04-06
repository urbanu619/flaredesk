#!/usr/bin/env bash
# 将 magic_admin_web 构建产物复制到 magic_admin/webdist/dist，供 Go embed 打入二进制（桌面式 / 单进程交付）。
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT/magic_admin_web"
[ -d node_modules ] || npm install
npm run build
rm -rf "$ROOT/magic_admin/webdist/dist"
mkdir -p "$ROOT/magic_admin/webdist"
cp -R "$ROOT/magic_admin_web/dist" "$ROOT/magic_admin/webdist/"
echo "OK: magic_admin_web/dist -> magic_admin/webdist/dist （请再执行 go build）"
