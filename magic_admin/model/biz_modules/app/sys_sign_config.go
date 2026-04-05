package app

// 引入关联包

type SysSignConfig struct { 
	Id int64 `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt int64 `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt int64 `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	IsSystemSign bool `json:"isSystemSign" gorm:"column:is_system_sign;type:tinyint(1);comment:'是否本系统密钥信息'"`
	SignName string `json:"signName" gorm:"column:sign_name;type:varchar(45);comment:'签名系统名称-前缀:FOMO-PRO';index;unique:uni_sys_sign_config_sign_name"`
	SignAddress string `json:"signAddress" gorm:"column:sign_address;type:varchar(42);comment:'系统地址'"`
	SignPriKey string `json:"signPriKey" gorm:"column:sign_pri_key;type:varchar(512);comment:'密文'"`
	SignExpSec int64 `json:"signExpSec" gorm:"column:sign_exp_sec;type:bigint;comment:'超时时间S'"`
	SysUrl string `json:"sysUrl" gorm:"column:sys_url;type:varchar(512);comment:外部系统请求链接"`
	ExchangeUrl string `json:"exchangeUrl" gorm:"column:exchange_url;type:varchar(512);comment:交易所URL请求链接"`
	ExchangeAppId string `json:"exchangeAppId" gorm:"column:exchange_app_id;type:varchar(512);comment:appID"`
	ExchangeAppSecret string `json:"exchangeAppSecret" gorm:"column:exchange_app_secret;type:varchar(512);comment:appSecret"`
}

func (*SysSignConfig) TableName() string {
	return "sys_sign_config"
}

func NewSysSignConfig() *SysSignConfig {
	return &SysSignConfig{}
}
