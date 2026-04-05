package ams_ast

import (
	"testing"
)

// 测试根路径

func TestRoot(t *testing.T) {
	t.Logf("go服务 服务根路由 入口文件地址:%s", ServerRootRouterEnterFile())
	t.Logf("go服务 子路由 入口文件地址:%s", ServerSubRouterEnterFile("biz"))
	t.Logf("go服务 子路由 注册文件地址:%s", ServerSubRouterRegisterFile("biz", "HelloWorld"))
}
