package redisclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	config2 "go_server/base/config"
	"go_server/base/core"
	"strings"
	"time"
)

var (
	defaultClient redis.UniversalClient
	_redisSync    *redsync.Redsync
	_redisPool    redsyncredis.Pool
)

func DefaultClient() redis.UniversalClient {
	if defaultClient == nil {
		redisClientInit()
	}
	return defaultClient
}

func redisSync() *redsync.Redsync {
	if _redisSync == nil {
		redisClientInit()
	}
	return _redisSync
}

func redisPool() redsyncredis.Pool {
	if _redisPool == nil {
		redisClientInit()
	}
	return _redisPool
}

var exConf *config2.Redis

func conf() *config2.Redis {
	if exConf == nil {
		exConf = config2.EnvConf().Redis
	}
	return exConf
}

// LockOptions 定义分布式锁的配置选项
type LockOptions struct {
	Expiration int32         // 锁的过期时间（秒）
	MaxRetries int           // 最大重试次数
	RetryDelay time.Duration // 重试间隔
}

// DefaultLockOptions 提供锁配置的默认值
var DefaultLockOptions = LockOptions{
	Expiration: 3,                      // 默认过期时间 3 秒
	MaxRetries: 5,                      // 默认最大重试 5 次
	RetryDelay: 200 * time.Millisecond, // 默认重试间隔 200 毫秒
}

func redisClientInit() {
	addr := conf().Addr             //viper.GetString("redis.addr")
	password := conf().Password     //viper.GetString("redis.password")
	db := conf().Db                 //viper.GetInt("redis.db")
	enabledTls := conf().EnabledTls //viper.GetBool("redis.enabledTls")
	if conf().EnabledCluster {      //viper.GetBool("redis.enabledCluster")
		addrs := conf().Cluster
		opts := &redis.ClusterOptions{
			Addrs:    addrs,
			Password: password,
		}
		if enabledTls {
			opts.TLSConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		defaultClient = redis.NewClusterClient(opts)
		//if err := defaultClient.Ping(context.TODO()).Err(); err != nil {
		//	core.Log.Fatal("集群 redis 连接失败, addr: " + strings.Join(addrs, ",") + ", error: " + err.Error())
		//}
		core.Log.Info("集群 redis 连接成功, addrs: ", strings.Join(addrs, ","))
	} else {
		opts := &redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}
		if enabledTls {
			opts.TLSConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
		defaultClient = redis.NewClient(opts)
		//if err := defaultClient.Ping(context.Background()).Err(); err != nil {
		//	core.Log.Fatal("单体 redis 连接失败, addr: " + addr + ", error: " + err.Error())
		//}
		core.Log.Info(" 单体 redis 连接成功, addr: ", addr)
	}
	// 初始化 redsync
	_redisPool = goredis.NewPool(defaultClient)
	_redisSync = redsync.New(_redisPool)
}

var redisKeyNonce = "nonce:%s"

func MustCode(address string) (string, error) {
	// 如果存在 nonce 则返回
	key := fmt.Sprintf(redisKeyNonce, address)
	nonce, err := DefaultClient().Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		return "", fmt.Errorf("GetNonce error")
	}
	if nonce != "" {
		return nonce, nil
	}
	// 保存到 redis
	result := DefaultClient().Set(context.Background(), key, nonce, time.Duration(0))
	if result.Err() != nil {
		return "", fmt.Errorf("GetNonce error")
	}
	return nonce, nil
}
