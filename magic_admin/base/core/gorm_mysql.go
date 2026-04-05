package core

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	config2 "go_server/base/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

func MainDb() *gorm.DB {
	return gormMysqlByConfig(config2.EnvConf().Mysql)
}

func BizDbs() (map[string]*gorm.DB, map[string]string, map[string]string) {
	dbMap := make(map[string]*gorm.DB)
	aliasDBMap := make(map[string]string)
	aliasProxyMap := make(map[string]string)
	for _, info := range config2.EnvConf().BizList {
		// 通过别名获取 db
		dbMap[info.AliasName] = gormMysqlByConfig(config2.Mysql{GeneralDB: info.GeneralDB})
		aliasDBMap[info.AliasName] = info.Dbname       // 通过db名 获取别名
		aliasProxyMap[info.ProxyAlias] = info.ProxyUrl // 别名 - 代理映射
		// 别名/db名 获取对应的业务请求地址
	}
	return dbMap, aliasDBMap, aliasProxyMap
}

func gormMysqlByConfig(m config2.Mysql) *gorm.DB {
	dialer := mysql.New(mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         256,     // string 类型字段的默认长度
		DisableDatetimePrecision:  true,    // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,    // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,    // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	})
	db, err := gorm.Open(dialer, gormConfig(m.Prefix, m.Singular))
	if err != nil {
		panic(err)
	}
	db.InstanceSet("gorm:table_options", "ENGINE=InnoDB")
	sqlDB, _ := db.DB()
	sqlDB.SetConnMaxIdleTime(20 * time.Second)
	sqlDB.SetConnMaxLifetime(20 * time.Second)
	sqlDB.SetMaxIdleConns(m.MaxIdleConns)
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)
	err = sqlDB.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

// Config gorm 自定义配置
func gormConfig(prefix string, singular bool) *gorm.Config {
	var general config2.GeneralDB
	general = config2.EnvConf().Mysql.GeneralDB
	return &gorm.Config{
		Logger: logger.New(NewWriter(general), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      general.LogLevel(),
			Colorful:      true,
		}),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,
			SingularTable: singular,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}
}

type Writer struct {
	config config2.GeneralDB
	writer logger.Writer
}

func NewWriter(config config2.GeneralDB) *Writer {
	return &Writer{config: config}
}

// Printf 格式化打印日志
func (c *Writer) Printf(message string, data ...any) {
	// 当有日志时候均需要输出到控制台
	fmt.Printf(message, data...)
	// 当开启了zap的情况，会打印到日志记录
	if c.config.LogZap {
		switch c.config.LogLevel() {
		case logger.Silent:
			Log.Debug(fmt.Sprintf(message, data...))
		case logger.Error:
			Log.Error(fmt.Sprintf(message, data...))
		case logger.Warn:
			Log.Warn(fmt.Sprintf(message, data...))
		case logger.Info:
			Log.Info(fmt.Sprintf(message, data...))
		default:
			Log.Info(fmt.Sprintf(message, data...))
		}
		return
	}
}
