package system

import (
	"go_server/service"

	"github.com/gin-gonic/gin"
)

//------------------------------------------自动化代码路由------------------------------------------//

type AutoRouter struct {
}

func (AutoRouter) Route() string {
	return "/auto"
}

// 控制层与实现层 合二为一 让同一个业务 尽量在一个文件中实现与暴露
var autoService = service.RealizationLayer.SystemServiceGroup.AutoService

func (h AutoRouter) Register(group *gin.RouterGroup) {
	group.POST("server/code", autoService.ServerCode)
	group.POST("server/rollback", autoService.Rollback)
	group.POST("model/auto", autoService.ModelAuto)

}

//------------------------------------------数据库信息查看路由------------------------------------------//

type DbRouter struct {
}

func (DbRouter) Route() string {
	return "/db"
}

// 控制层与实现层 合二为一 让同一个业务 尽量在一个文件中实现与暴露
var sysDbService = service.RealizationLayer.SystemServiceGroup.DBService

func (h DbRouter) Register(group *gin.RouterGroup) {
	group.GET("dbs", sysDbService.Dbs)
	group.GET("tbs", sysDbService.Tbs)
	group.GET("cols", sysDbService.Cols)
}

//------------------------------------------用户管理路由------------------------------------------//

type UserRouter struct {
}

func (UserRouter) Route() string {
	return "/user"
}

// 控制层与实现层 合二为一 让同一个业务 尽量在一个文件中实现与暴露
var sysUserService = service.RealizationLayer.SystemServiceGroup.UserService

func (h UserRouter) Register(group *gin.RouterGroup) {
	// 设置本人信息
	group.GET("info", sysUserService.Info)
	group.POST("set", sysUserService.Set)
	group.GET("getGoogleKey", sysUserService.GetGoogleKey)
	group.POST("cancelGoogleKey", sysUserService.CancelGoogleKey)
	group.POST("replaceGoogleKey", sysUserService.ReplaceGoogleKey)
	// 管理
	group.GET("get", sysUserService.Get)          // 获取用户详情
	group.GET("find", sysUserService.Find)        // 用户列表
	group.POST("setUser", sysUserService.SetUser) // 设置用户信息
	group.POST("create", sysUserService.Create)   // 创建用户
	group.GET("delete", sysUserService.Del)       // 删除用户
	group.GET("logs", sysUserService.Logs)        // 操作日志

}

//------------------------------------------角色管理路由------------------------------------------//

type RoleRouter struct {
}

func (RoleRouter) Route() string {
	return "/role"
}

var roleService = service.RealizationLayer.SystemServiceGroup.RoleService

func (h RoleRouter) Register(group *gin.RouterGroup) {
	group.GET("get", roleService.Get)
	group.GET("find", roleService.Find)
	group.POST("set", roleService.Set)
	group.POST("create", roleService.Create)
	group.GET("delete", roleService.Del)
}

//------------------------------------------菜单管理路由------------------------------------------//

type MenuRouter struct {
}

func (MenuRouter) Route() string {
	return "/menu"
}

var menuService = service.RealizationLayer.SystemServiceGroup.MenuService

func (h MenuRouter) Register(group *gin.RouterGroup) {
	group.GET("tree", menuService.Tree)

	group.POST("get", menuService.Get)
	group.POST("set", menuService.Set)
	group.POST("find", menuService.Find)
	group.POST("create", menuService.Create)
	group.GET("delete", menuService.Del)
}

//------------------------------------------apis管理路由------------------------------------------//

type ApisRouter struct {
}

func (ApisRouter) Route() string {
	return "/apis"
}

var apisService = service.RealizationLayer.SystemServiceGroup.ApisService

func (h ApisRouter) Register(group *gin.RouterGroup) {
	group.GET("tree", apisService.Tree)
	group.GET("get", apisService.Get)
	group.GET("find", apisService.Find)
	group.POST("set", apisService.Set)
	group.POST("create", apisService.Create)
	group.GET("delete", apisService.Del)
}

type DictionaryRouter struct {
}

func (DictionaryRouter) Route() string {
	return "/dictionary"
}

var dictionaryService = service.RealizationLayer.SystemServiceGroup.DictionaryService

func (h DictionaryRouter) Register(group *gin.RouterGroup) {
	group.GET("tree", dictionaryService.DictionariesTree)

	group.POST("find", dictionaryService.FindDictionaries)
	group.POST("create", dictionaryService.CreateDictionary)
	group.POST("set", dictionaryService.UpdateDictionary)
	group.POST("delete", dictionaryService.DelDictionary)

	group.POST("detail/find", dictionaryService.FindDictionaryDetail)
	group.POST("detail/create", dictionaryService.CreateDictionaryDetail)
	group.POST("detail/set", dictionaryService.UpdateDictionaryDetail)
	group.POST("detail/delete", dictionaryService.DelDictionaryDetail)

}

type FileRouter struct {
}

func (FileRouter) Route() string {
	return "/file"
}

var fileRouter = service.RealizationLayer.SystemServiceGroup.FileService

func (h FileRouter) Register(group *gin.RouterGroup) {
	group.POST("/upload", fileRouter.UploadFile)
	group.GET("/delete", fileRouter.DeleteFile)
	group.GET("/oss/auth", fileRouter.OssAuth)

}

type SignRouter struct {
}

func (SignRouter) Route() string {
	return "/sign"
}

var signRouter = service.RealizationLayer.SystemServiceGroup.SignService

func (h SignRouter) Register(group *gin.RouterGroup) {
	group.POST("/set", signRouter.Set)
	group.GET("/find", signRouter.Find)
}
