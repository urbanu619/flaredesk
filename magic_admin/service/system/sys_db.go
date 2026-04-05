package system

import (
	"github.com/gin-gonic/gin"
	"go_server/model/common/response"
	"go_server/service/base"
)

type DBService struct {
	base.SysCommonService
}

func (s *DBService) Dbs(c *gin.Context) {
	dbs, err := s.GetBizDbs()
	if err != nil {
		response.Resp(c, "获取失败")
	} else {
		response.Resp(c, gin.H{"alias": dbs})
	}
}

func (s *DBService) Tbs(c *gin.Context) {
	alias := c.Query("alias")
	dbs, err := s.GetTables(alias)
	if err != nil {
		response.Resp(c, "获取失败")
	} else {
		response.Resp(c, gin.H{"dbs": dbs})
	}
}

func (s *DBService) Cols(c *gin.Context) {
	alias := c.Query("alias")
	tableName := c.Query("tableName")
	if alias == "" || tableName == "" {
		response.Resp(c, "参数不足")
	}
	dbs, err := s.GetColumn(alias, tableName)
	if err != nil {
		response.Resp(c, "获取失败")
	} else {
		response.Resp(c, gin.H{"dbs": dbs})
	}
	response.Resp(c)
}
