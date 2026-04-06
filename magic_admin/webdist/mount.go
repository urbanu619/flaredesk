package webdist

import (
	"io/fs"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go_server/base/core"
)

func contentTypeByExt(rel string) string {
	switch strings.ToLower(path.Ext(rel)) {
	case ".js":
		return "text/javascript; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".svg":
		return "image/svg+xml"
	case ".json":
		return "application/json"
	case ".html":
		return "text/html; charset=utf-8"
	case ".png":
		return "image/png"
	case ".ico":
		return "image/x-icon"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".map":
		return "application/json"
	default:
		return "application/octet-stream"
	}
}

// Mount 将嵌入式前端挂到未匹配路由（/admin、/health 除外）。
// /assets 使用 fs.ReadFile：不经过 gin.FileFromFS/http.FileServer，避免路径与 embed 组合时 Open 失败却仅返回「404 page not found」。
// 路径一律从 Request.URL.Path 截取 /assets/ 之后部分，避免依赖 *filepath 参数名在各版本 Gin 下的差异。
func Mount(engine *gin.Engine, enable bool) {
	if !enable {
		return
	}
	root, err := fs.Sub(Dist, "dist")
	if err != nil {
		core.Log.Warnf("webdist: 无法打开嵌入式目录: %v", err)
		return
	}

	serveAsset := func(c *gin.Context) {
		rest, ok := strings.CutPrefix(c.Request.URL.Path, "/assets/")
		if !ok {
			c.Status(http.StatusNotFound)
			return
		}
		if rest == "" || strings.Contains(rest, "..") {
			c.Status(http.StatusNotFound)
			return
		}
		rel := path.Join("assets", rest)
		if rel != "assets" && !strings.HasPrefix(rel, "assets/") {
			c.Status(http.StatusNotFound)
			return
		}
		data, err := fs.ReadFile(root, rel)
		if err != nil {
			core.Log.Warnf("webdist: embed 无此文件 %q（请先 ./scripts/embed-web.sh 后重新 go build）: %v", rel, err)
			c.Status(http.StatusNotFound)
			return
		}
		ct := contentTypeByExt(rel)
		c.Header("Cache-Control", "public, max-age=31536000, immutable")
		c.Header("X-Content-Type-Options", "nosniff")
		if c.Request.Method == http.MethodHead {
			c.Header("Content-Type", ct)
			c.Header("Content-Length", strconv.Itoa(len(data)))
			c.Status(http.StatusOK)
			return
		}
		c.Data(http.StatusOK, ct, data)
	}
	engine.GET("/assets/*filepath", serveAsset)
	engine.HEAD("/assets/*filepath", serveAsset)

	engine.GET("/logo.svg", func(c *gin.Context) {
		c.FileFromFS("logo.svg", http.FS(root))
	})

	engine.NoRoute(func(c *gin.Context) {
		p := c.Request.URL.Path
		if strings.HasPrefix(p, "/admin") || strings.HasPrefix(p, "/health") {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "not found"})
			return
		}
		// 根路径、显式 index、以及无后缀的 history 路由：直接读字节，避免与目录索引 302 冲突。
		if p == "/" || p == "/index.html" || !strings.Contains(path.Base(p), ".") {
			data, err := fs.ReadFile(root, "index.html")
			if err != nil {
				c.String(http.StatusInternalServerError, "webdist: 读取 index.html 失败: %v", err)
				return
			}
			c.Data(http.StatusOK, "text/html; charset=utf-8", data)
			return
		}
		// dist 根目录下其它带后缀文件（如未来增加的 favicon.ico）
		name := strings.TrimPrefix(p, "/")
		if !strings.Contains(name, "/") {
			if b, err := fs.ReadFile(root, name); err == nil {
				c.Data(http.StatusOK, http.DetectContentType(b), b)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "msg": "not found"})
	})
	core.Log.Info("webdist: 已启用嵌入式前端，浏览器访问本机监听端口根路径即可（与 API 同源 /admin）")
}
