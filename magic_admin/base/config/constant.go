package config

import "os"

var autoCreateConfig bool = true

func init() {
	if autoCreateConfig {
		if _, err := os.Stat(ConfigDir); os.IsNotExist(err) {
			err := os.Mkdir(ConfigDir, 0755) // 设置文件权限为rwxr-xr-x（默认值）
			if err != nil {
				panic(err)
			}
		}
	}
}

// 是否自动创建配置文件

const (
	VerifyCode        = "111111"
	ConfigDir         = "./conf.d/" //
	DefaultMod        = "default"
	ModEnvDev         = "dev"
	ModEnvTest        = "test"
	ModEnvProd        = "prod"
	ConfigAppFile     = ConfigDir + "app.json"
	ConfigDefaultFile = ConfigDir + "config.json"
	ConfigDevFile     = ConfigDir + "config.dev.json"
	ConfigTestFile    = ConfigDir + "config.test.json"
	ConfigProdFile    = ConfigDir + "config.prod.json"
)

var configMap = map[string]string{
	DefaultMod: ConfigDefaultFile,
	ModEnvDev:  ConfigDevFile,
	ModEnvTest: ConfigTestFile,
	ModEnvProd: ConfigProdFile,
}

func GetConfigFileNameByMod(mod string) string {
	if v, ok := configMap[mod]; ok {
		return v
	}
	return ConfigDefaultFile
}
