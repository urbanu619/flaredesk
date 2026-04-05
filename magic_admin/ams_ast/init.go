package ams_ast

import (
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

// 基础配置

var (
	root, _                 = filepath.Abs("./") // 项目根路径 -- 该路径仅用于开发阶段
	server_mod              = "go_server"        // go服务模型名
	server_enter            = "enter.go"         // go服务 入口文件名称
	server_router           = "router"           // 路由服务路径
	server_model            = "model"            // 模型模块根路径
	server_model_biz        = "biz_modules"      // 业务模块存放路径
	server_service          = "service"          // 实现层路径
	server_biz_base_service = "base"             // 业务基础模块
	// server
	// 路由相关
	serverRootRouterPath      = filepath.Join(root, server_router)                // 跟路由文件夹路径
	serverRootRouterEnterFile = filepath.Join(serverRootRouterPath, server_enter) // 根路由入口文件地址  任务:往文件中注册子路由组
	// 模型相关
	modelTargetPath = filepath.Join(root, server_model) // 模型入口
	// template 模版使用相关规范
	// tips: 模板路径与服务配置路径保持一致 /根路径root/服务名/模块名/文件名
	// demo:
	//		子路由入口文件 model - 模型名称 命名规范 小写蛇形
	//		模板所在路径 ./resource/package/go_server/router/model/enter.go.tpl
	// 		生成路径: root/go_server/router/model/enter.go
	//
	//		子路由
	// 		模板所在路径: ./resource/package/go_server/router/model/enter.go.tpl
	//		文件生成路径: root/go_server/router/model/model.go
	template_root = filepath.Join(root, "/ams_ast/resource/package") // 模版存放根路径
)

// snakeString 驼峰转蛇形 XxYy to xx_yy , XxYY to xx_y_y

func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		//判断如果字母为大写的A-Z就在前面拼接一个_
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	//ToLower把大写字母统一转小写
	return strings.ToLower(string(data[:]))
}

// 下划线转驼峰 xx_yy to xxYy
func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		if i > 0 {
			parts[i] = strings.Title(parts[i])
		}
	}
	return strings.Join(parts, "")
}

// 驼峰转下划线 XxYy to xx_yy , XxYY to xx_y_y
func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

// 字符串首字母大写/小写

func CapitalizeOrLower(str string, toLower ...bool) string {
	if str == "" {
		return str
	}
	r, size := utf8.DecodeRuneInString(str)
	if len(toLower) > 0 {
		return string(unicode.ToLower(r)) + str[size:]
	}
	return string(unicode.ToUpper(r)) + str[size:]
}

// 文件是否存在

func FileExists(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
