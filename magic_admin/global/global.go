package global

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"sync"
)

var (
	AMS_DB                  *gorm.DB            // 管理后台数据库
	AMS_BIZ_DBS             map[string]*gorm.DB // 业务库列表
	AMS_BIZ_ALIAS_DB_MAP    map[string]string   // -别名 - 数据库DB 映射
	AMS_BIZ_ALIAS_PROXY_MAP map[string]string   // 别名-代理 映射
	GVA_ROUTERS             gin.RoutesInfo      // 当前http路由信息
	lock                    sync.RWMutex
	AMS_REDIS               redis.UniversalClient
)

// BizDBByAlias 通过名称获取db list中的db

func BizDBByAlias(dbAlias string) (*gorm.DB, error) {
	lock.RLock()
	defer lock.RUnlock()
	// dbname转换为alias
	for k, dbName := range AMS_BIZ_ALIAS_DB_MAP {
		if dbName == dbAlias {
			dbAlias = k
		}
	}
	v, ok := AMS_BIZ_DBS[dbAlias]
	if !ok {
		return nil, errors.New("biz db not found")
	}
	return v, nil
}

func DefaultAlias() string {
	for k, _ := range AMS_BIZ_ALIAS_DB_MAP {
		return k
	}
	return ""
}
