package system

import (
	"go_server/model/common"
)

type SysSignConfig struct {
	common.GormBaseModel
	IsSystemSign bool   `json:"isSystemSign" gorm:"comment:'是否本系统密钥信息'"`
	SignName     string `json:"signName" gorm:"type:varchar(45);unique;comment:'签名系统名称-前缀:FOMO-PRO'"`
	SignAddress  string `json:"signAddress" gorm:"type:varchar(42);comment:'系统地址'"`
	SignPriKey   string `json:"signPriKey" gorm:"type:varchar(512);comment:'密文'"`
	SignExpSec   int64  `json:"signExpSec" gorm:"comment:'超时时间S'"`
	SysUrl       string `json:"sysUrl" gorm:"type:varchar(512);comment:外部系统请求链接"`
}

func (*SysSignConfig) TableName() string {
	return common.ModelPrefix + "sign_config"
}

func (*SysSignConfig) Comment() string {
	return "签名系统配置表"
}

func NewSysSignConfig() *SysSignConfig {
	return &SysSignConfig{}
}
