package config

type Redis struct {
	Addr           string   `mapstructure:"addr" json:"addr" yaml:"addr"`                                // url
	Password       string   `mapstructure:"password" json:"password" yaml:"password"`                    // 密码
	Db             int      `mapstructure:"db" json:"db" yaml:"db"`                                      // index db
	EnabledTls     bool     `mapstructure:"enabledTls" json:"enabledTls" yaml:"enabledTls"`              // 是否 启用tls
	EnabledCluster bool     `mapstructure:"enabledCluster" json:"enabledCluster"  yaml:"enabledCluster"` // 是否集群
	Cluster        []string `mapstructure:"cluster" json:"cluster" yaml:"cluster"`
}

type RedisMq struct {
	Addr           string    `mapstructure:"addr" json:"addr" yaml:"addr"`                                // url
	Password       string    `mapstructure:"password" json:"password" yaml:"password"`                    // 密码
	Db             int       `mapstructure:"db" json:"db" yaml:"db"`                                      // index db
	EnabledTls     bool      `mapstructure:"enabledTls" json:"enabledTls" yaml:"enabledTls"`              // 是否 启用tls
	EnabledCluster bool      `mapstructure:"enabledCluster" json:"enabledCluster"  yaml:"enabledCluster"` // 是否集群
	Cluster        []string  `mapstructure:"cluster" json:"cluster" yaml:"cluster"`
	Consumer       *Consumer `mapstructure:"consumer" json:"consumer" yaml:"consumer"`
}

type Consumer struct {
	Concurrency int      `mapstructure:"concurrency" json:"concurrency" yaml:"concurrency"`
	Queues      []*Queue `mapstructure:"queues" json:"queues" yaml:"queues"`
}

type Queue struct {
	Name     string `mapstructure:"name" json:"name" yaml:"name"`
	Priority int    `mapstructure:"priority" json:"priority" yaml:"priority"`
}
