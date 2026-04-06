package base

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go_server/base/config"
	"go_server/base/core"
	"go_server/global"
	system2 "go_server/model/system"
	"gorm.io/gorm"
	"time"
)

var sysCommonService *SysCommonService

// 业务库公共服务

type SysCommonService struct {
}

func GetSysService() *SysCommonService {
	if sysCommonService == nil {
		sysCommonService = &SysCommonService{}
	}
	return sysCommonService
}

func (*SysCommonService) DB() *gorm.DB {
	return global.AMS_DB
}

func (*SysCommonService) CheckIsAdmin(userId int64) error {
	// 限制用户增加必须管理员才可以操作
	if userId != system2.AdminId {
		return fmt.Errorf("不允许操作")
	}
	return nil
}

func (*SysCommonService) Redis() redis.UniversalClient {
	return global.AMS_REDIS
}

func (s *SysCommonService) unique(a string) string {
	return fmt.Sprintf("ADMIN:UNIKEY:%s", a)
}

// 接口重复调用限制

func (s *SysCommonService) RepeatFilter(serviceName string, timeSt time.Duration) bool {
	if err := s.Redis().Ping(context.Background()).Err(); err != nil {
		return true
	}
	// 尝试设置KEY
	nxBool, err := s.Redis().SetNX(context.Background(), s.unique(serviceName), true, timeSt).Result()
	if err != nil {
		return false
	}
	if !nxBool {
		return false
	}
	// 获取过期时间 -- 可用于验证码获取上 暂五使用场景
	_, err = s.Redis().TTL(context.Background(), s.unique(serviceName)).Result()
	if err != nil {
		return false
	}
	return true
}

// GetDB 获取数据库的所有数据库名

func (s *SysCommonService) GetBizDbs() (alias map[string]string, err error) {
	return global.AMS_BIZ_ALIAS_DB_MAP, err
}

// GetTables 获取数据库的所有表名

type Table struct {
	TableName string `json:"tableName" gorm:"column:table_name"`
}

func (s *SysCommonService) GetTables(dbAlias string) ([]Table, error) {
	db, err := global.BizDBByAlias(dbAlias)
	if err != nil {
		return nil, err
	}
	core.Log.Infof("dbName:%s", db.Migrator().CurrentDatabase())
	return MetaListTables(db)
}

// GetColumn 获取指定数据库和指定数据表的所有字段名,类型值等

type Column struct {
	DataType      string `json:"dataType" gorm:"column:data_type"`
	ColumnName    string `json:"columnName" gorm:"column:column_name"`
	DataTypeLong  string `json:"dataTypeLong" gorm:"column:data_type_long"`
	ColumnComment string `json:"columnComment" gorm:"column:column_comment"`
	PrimaryKey    bool   `json:"primaryKey" gorm:"column:primary_key"`
}

func (s *SysCommonService) GetColumn(dbAlias string, tableName string) (data []Column, err error) {
	db, err := global.BizDBByAlias(dbAlias)
	if err != nil {
		return nil, err
	}
	core.Log.Infof("dbName:%s tableName:%s", db.Migrator().CurrentDatabase(), tableName)
	return MetaListColumns(db, tableName)
}

func (s *SysCommonService) CheckFileJob() bool {
	row, ok := GetLast[system2.FileJob](s.DB())
	if !ok {
		return true
	}
	return row.Status != config.FileJobRunning
}

// 新建一个运行中文件任务

func (s *SysCommonService) NewFileJob(userId int64) (*system2.FileJob, error) {
	row, ok := GetOne[system2.Administrator](s.DB(), "id", userId)
	if !ok {
		return nil, fmt.Errorf("admin is not found")
	}
	// 创建文件生成任务
	newJob := &system2.FileJob{
		Path:     "",
		Total:    0,
		Size:     "",
		Status:   config.FileJobRunning,
		Desc:     "",
		CreateBy: fmt.Sprintf("%s:%d", row.Nickname, row.ID),
	}
	err := s.DB().Create(&newJob).Error
	return newJob, err
}

// 保存文件任务

func (s *SysCommonService) SaveJob(job *system2.FileJob) error {
	_, ok := GetOne[system2.FileJob](s.DB(), "id", job.ID)
	if !ok {
		return fmt.Errorf("file job is not found")
	}
	err := s.DB().Save(&job).Error
	return err
}
