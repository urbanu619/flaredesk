package config

import (
	"gorm.io/gorm/logger"
	"path/filepath"
	"strings"
)

type Mysql struct {
	GeneralDB `yaml:",inline" mapstructure:",squash"`
}

func (m *Mysql) Dsn() string {
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
}

type DsnProvider interface {
	Dsn() string
}

// GeneralDB 也被 Pgsql 和 Mysql 原样使用
type GeneralDB struct {
	Prefix       string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	Port         string `mapstructure:"port" json:"port" yaml:"port"`                               // 数据库端口
	Config       string `mapstructure:"config" json:"config" yaml:"config"`                         // 高级配置
	Dbname       string `mapstructure:"db-name" json:"db-name" yaml:"db-name"`                      // 数据库名
	Username     string `mapstructure:"username" json:"username" yaml:"username"`                   // 数据库账号
	Password     string `mapstructure:"password" json:"password" yaml:"password"`                   // 数据库密码
	Path         string `mapstructure:"path" json:"path" yaml:"path"`                               // 数据库地址
	Engine       string `mapstructure:"engine" json:"engine" yaml:"engine" default:"InnoDB"`        // 数据库引擎，默认InnoDB
	LogMode      string `mapstructure:"log-mode" json:"log-mode" yaml:"log-mode"`                   // 是否开启Gorm全局日志
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"` // 打开到数据库的最大连接数
	Singular     bool   `mapstructure:"singular" json:"singular" yaml:"singular"`                   // 是否开启全局禁用复数，true表示开启
	LogZap       bool   `mapstructure:"log-zap" json:"log-zap" yaml:"log-zap"`                      // 是否通过zap写入日志文件
	Driver       string `mapstructure:"driver" json:"driver" yaml:"driver"`                         // mysql（默认）或 sqlite
}

// UseSQLite 为 true 时使用 SQLite 单文件数据库（path/db-name 含义见 SqliteFile）。
func (g GeneralDB) UseSQLite() bool {
	return strings.EqualFold(strings.TrimSpace(g.Driver), "sqlite")
}

// SqliteFile 返回 SQLite 数据库文件路径：path 以 .db/.sqlite 结尾则视为完整路径，否则与 db-name 拼接为目录+文件名。
func (g GeneralDB) SqliteFile() string {
	p := strings.TrimSpace(g.Path)
	if p == "" {
		if g.Dbname == "" {
			return "flaredesk.db"
		}
		return g.Dbname + ".db"
	}
	lower := strings.ToLower(p)
	if strings.HasSuffix(lower, ".db") || strings.HasSuffix(lower, ".sqlite") {
		return p
	}
	if g.Dbname == "" {
		return filepath.Join(p, "flaredesk.db")
	}
	return filepath.Join(p, g.Dbname+".db")
}

func (c GeneralDB) LogLevel() logger.LogLevel {
	switch strings.ToLower(c.LogMode) {
	case "silent", "Silent":
		return logger.Silent
	case "error", "Error":
		return logger.Error
	case "warn", "Warn":
		return logger.Warn
	case "info", "Info":
		return logger.Info
	default:
		return logger.Info
	}
}
