package webdist

import (
	"embed"
)

// Dist 为前端构建目录（Vite `npm run build` 产出）。
// 占位仅含 index.html，完整资源请执行仓库根目录 scripts/embed-web.sh 再编译后端。
//
//go:embed all:dist
var Dist embed.FS
