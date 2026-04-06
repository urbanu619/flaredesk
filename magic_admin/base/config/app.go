package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type App struct {
	Addr         int    `mapstructure:"addr" json:"addr" yaml:"addr"`                            // 端口值
	Mod          string `mapstructure:"mod" json:"mod" yaml:"mod"`                               // 运行模式
	RouterPrefix string `mapstructure:"router-prefix" json:"router-prefix" yaml:"router-prefix"` //  路由前缀
	ProxyUrl     string `mapstructure:"proxy-url" json:"proxy-url" yaml:"proxy-url"`             // 业务默认代理地址
	ServeWeb     bool   `mapstructure:"serve-web" json:"serve-web" yaml:"serve-web"`             // true：由 Go 嵌入并提供前端静态资源（webdist/dist），与 API 同源
}

func AppConf() *App {
	if appConf == nil {
		appConfInit()
	}
	return appConf
}

var appConf *App

// app配置初始化

func appConfInit() {
	vp := viper.New()
	vp.SetConfigFile(ConfigAppFile) // app.json
	vp.SetConfigType("json")
	vp.SetDefault("addr", AdminHostPort)
	vp.SetDefault("mod", DefaultMod)
	vp.SetDefault("router-prefix", "admin")
	vp.SetDefault("proxy-url", fmt.Sprintf("http://127.0.0.1:%d", BizHostPort))
	vp.SetDefault("serve-web", false)
	handVpToFile(vp) // 文件读取 如果读取不到则生成(dev环境)
	if err := vp.Unmarshal(&appConf); err != nil {
		panic(err)
	}
}

func handVpToFile(vp *viper.Viper) {
	err := vp.ReadInConfig()
	if err != nil {
		if err := vp.WriteConfig(); err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
	}
}
