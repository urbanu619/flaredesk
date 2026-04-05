package core

import (
	"fmt"
	"go_server/base/config"
	bizApp "go_server/model/biz_modules/app"
	"go_server/model/system"
	"gorm.io/gorm"
)

// 数据库同步

// MigrateTable 迁移表接口
type MigrateTable interface {
	TableName() string
	Comment() string
}

var allTables = []MigrateTable{
	&system.Administrator{},
	&system.AdministratorLog{},
	&system.Role{},
	&system.Apis{},
	&system.Menus{},
	&system.Dictionaries{},
	&system.DictionaryDetail{},
	&system.FileJob{},
	&system.SysSignConfig{},
}

// bizTables 业务库需要迁移的表
var bizTables = []MigrateTable{
	&bizApp.CfAccount{},
	&bizApp.CfDnsTemplate{},
	&bizApp.CfZone{},
}

func Migrates() {
	// 系统表迁移（主库）
	db := MainDb()
	CurrentDatabase := db.Migrator().CurrentDatabase()
	Log.Info(fmt.Sprintf("当前数据库[%s]", CurrentDatabase))
	mTables := make([]MigrateTable, 0)
	mTables = append(mTables, allTables...)
	migrationTable(db, mTables) // 同步数据库结构
	sysDataInit(db)             // 数据初始化
	sysSignInfoInit(db)         // 系统交互密钥初始化

	// 业务表迁移（biz 库）
	bizDbs, _, _ := BizDbs()
	for alias, bizDb := range bizDbs {
		Log.Info(fmt.Sprintf("开始迁移业务库[%s]", alias))
		migrationTable(bizDb, bizTables)
	}

	Log.Info(fmt.Sprintf("数据库迁移完成"))
}

func sysDataInit(db *gorm.DB) {
	if err := system.NewAdministrator().DataInit(db); err != nil {
		Log.Error("error:", err.Error())
	}
	if err := system.NewApis().DataInit(db); err != nil {
		Log.Error("error:", err.Error())
	}
	if err := system.NewMenus().DataInit(db); err != nil {
		Log.Error("error:", err.Error())
	}
	if err := system.NewRole().DataInit(db); err != nil {
		Log.Error("error:", err.Error())
	}
}

// migrationTable 迁移数据表
func migrationTable(db *gorm.DB, tables []MigrateTable) {
	for _, table := range tables {
		//slog.Info(fmt.Sprintf("开始迁移[%s]表", table.TableName()))
		db = db.Set("gorm:table_options", "ENGINE=InnoDB")
		comment := fmt.Sprintf("COMMENT='%s'", table.Comment())
		db = db.Set("gorm:table_options", comment)
		err := db.Migrator().AutoMigrate(table)
		if err != nil {
			Log.Error(fmt.Sprintf("[%s]表迁移失败：%s", table.TableName(), err.Error()))
		}
		Log.Info(fmt.Sprintf("[%s]表迁移完成", table.TableName()))
	}
}

const (
	SignKey                   = "sign"
	SignSystemName            = "app"   // 本系统名称 -- 用于接口交互时候身份识别
	SignSystemAdminServerName = "admin" // 管理后台服务
)

// 初始化系统签名信息

func sysSignInfoInit(db *gorm.DB) {
	var rows []system.SysSignConfig
	// 初始化本系统密钥信息
	findRow := system.NewSysSignConfig()
	err := db.Model(&system.SysSignConfig{}).Where("sign_name", SignSystemAdminServerName).Find(&findRow).Error
	if err != nil {
		panic(err)
	}
	if findRow.ID == 0 {
		addr, pri, err := GenerateSysInfo()
		if err != nil {
			panic(err)
		}
		rows = append(rows, system.SysSignConfig{
			IsSystemSign: true,
			SignName:     SignSystemAdminServerName,
			SignAddress:  addr,
			SignPriKey:   pri,
			SignExpSec:   5,
			SysUrl:       fmt.Sprintf("http://127.0.0.1:%d", config.AdminHostPort),
		})
		//bizAddr, bizPri, err := GenerateSysInfo()
		//if err != nil {
		//	panic(err)
		//}
		//rows = append(rows, system.SysSignConfig{
		//	IsSystemSign: false,
		//	SignName:     SignSystemName,
		//	SignAddress:  bizAddr,
		//	SignPriKey:   bizPri,
		//	SignExpSec:   5,
		//	SysUrl:       fmt.Sprintf("http://127.0.0.1:%d", config.BizHostPort),
		//})
	}
	if len(rows) > 0 {
		err = db.Create(&rows).Error
		if err != nil {
			panic(err)
		}
	}
}
