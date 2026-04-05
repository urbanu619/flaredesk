package config

type BizInfo struct {
	AliasName  string `mapstructure:"alias-name" json:"alias-name" yaml:"alias-name"`    // 数据库别名 -- 一旦生成代码 别名不变
	ProxyUrl   string `mapstructure:"proxy-url" json:"proxy-url" yaml:"proxy-url"`       // 业务操作URL
	ProxyAlias string `mapstructure:"proxy-alias" json:"proxy-alias" yaml:"proxy-alias"` // mng 代理服务别名
	GeneralDB  `yaml:",inline" mapstructure:",squash" json:",inline"`
}
