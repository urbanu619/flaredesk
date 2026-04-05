package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type EnvConfig struct {
	Mysql   Mysql     `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	BizList []BizInfo `mapstructure:"biz-list" json:"biz-list" yaml:"biz-list"`
	JWT     JWT       `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Zap     *Zap      `mapstructure:"zap" json:"zap" yaml:"zap"`
	Redis   *Redis    `mapstructure:"redis" json:"redis" yaml:"redis"`
	File    *File     `mapstructure:"file" json:"file" yaml:"file"`
	StsOss  *OssSts   `mapstructure:"sts-oss" json:"sts-oss" yaml:"sts-oss"`
}

var envConfig *EnvConfig

func EnvConf() *EnvConfig {
	if envConfig == nil {
		appEnvConfigInit()
	}
	return envConfig
}

// 环境配置初始化

func appEnvConfigInit() {
	configFile := GetConfigFileNameByMod(AppConf().Mod)
	vp := viper.New()
	vp.SetConfigFile(configFile)
	vp.SetConfigType("json")
	setEnvDefaultConf(vp)
	handVpToFile(vp)
	if err := vp.Unmarshal(&envConfig); err != nil {
		panic(err)
	}
}

func setEnvDefaultConf(v *viper.Viper) {
	v.SetDefault("file", &File{
		OssType:   OssTypeLocal,
		Path:      fmt.Sprintf("http://127.0.0.1:%d", AdminHostPort),
		ProxyPath: "/admin/api/static",
		StorePath: "./static",
		OriginConf: map[string]interface{}{
			"endpoint":          "",
			"access-key-id":     "",
			"access-key-secret": "",
			"bucket-name":       "",
			"bucket-url":        "",
			"base-path":         "",
		},
	})
	v.SetDefault("mysql", Mysql{
		GeneralDB: GeneralDB{
			Prefix:       "",
			Port:         "3306",
			Config:       `charset=utf8mb4&parseTime=True&loc=Local`,
			Dbname:       AdminDbName,
			Username:     "root",
			Password:     "123456",
			Path:         "127.0.0.1",
			Engine:       "",
			LogMode:      "error",
			MaxIdleConns: 10,
			MaxOpenConns: 10,
			Singular:     false,
			LogZap:       false,
		}})
	v.SetDefault("biz-list", []BizInfo{
		{
			AliasName:  BizAliasName,
			ProxyUrl:   fmt.Sprintf("http://127.0.0.1:%d", BizHostPort),
			ProxyAlias: BizAliasName,
			GeneralDB: GeneralDB{
				Prefix:       "",
				Port:         "3306",
				Config:       `charset=utf8mb4&parseTime=True&loc=Local`,
				Dbname:       BizDbName,
				Username:     "root",
				Password:     "123456",
				Path:         "127.0.0.1",
				Engine:       "",
				LogMode:      "error",
				MaxIdleConns: 10,
				MaxOpenConns: 10,
				Singular:     false,
				LogZap:       false,
			},
		},
	})

	v.SetDefault("zap", &Zap{
		TagName:       "app",
		Level:         "debug",
		Prefix:        "AMS_",
		Format:        "_json",
		Director:      "logs",
		EncodeLevel:   "LowercaseLevelEncoder",
		StacktraceKey: "error",
		MaxAge:        3,
		RotationSize:  100,
		RotationCount: 3,
		ShowLine:      true,
		LogInConsole:  true,
		LogOutputFile: true,
	})

	v.SetDefault("redis", &Redis{
		Addr:           "127.0.0.1:6379",
		Password:       "",
		Db:             0,
		EnabledTls:     false,
		EnabledCluster: false,
		Cluster:        []string{"127.0.0.1:6379"},
	})

	v.SetDefault("jwt", &JWT{
		SigningKey:  fmt.Sprintf("%s.!@#$1234", BizAliasName),
		ExpiresTime: 24,
		Issuer:      "appTemp",
	})
	v.SetDefault("sts-oss", &OssSts{
		AccessKeyId:        "",
		AccessKeySecret:    "",
		BucketName:         "",
		BucketUrl:          "",
		BasePath:           "",
		StsEndpoint:        "",
		StsDurationSeconds: 1200,
		StsRoleSessionName: "",
		StsRoleArn:         "",
		StsRegion:          "",
	})
}
