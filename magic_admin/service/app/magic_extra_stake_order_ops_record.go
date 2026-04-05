package app

import (
	model "go_server/model/biz_modules/app"
	"go_server/service/base"
	"go_server/base/core"
	"go_server/model/common/response"
	"github.com/gin-gonic/gin"
	"strings"
)

type MagicExtraStakeOrderOpsRecordService struct {
	base.BizCommonService
}

func (s *MagicExtraStakeOrderOpsRecordService) Get(c *gin.Context) {
    s.SetDbAlias("app")
    base.NewBaseHandler(model.NewMagicExtraStakeOrderOpsRecord()).Get(c,s.DB())
}

func (s *MagicExtraStakeOrderOpsRecordService) Find(c *gin.Context) {
    s.SetDbAlias("app")
    type request[T any] struct {
    	base.ListRequest[T]
    	Id        *int64  `form:"id"`
    	IsExport         *bool   `form:"isExport"`
        Fields           *string `form:"fields"` // 指定返回字段 , 分割
    	BeginTime *int64  `form:"beginTime"`
        EndTime   *int64  `form:"endTime"`
    }
    req := new(request[model.MagicExtraStakeOrderOpsRecord])
    if err := c.BindQuery(req); err != nil {
    	response.Resp(c, err.Error())
    	return
    }
    db := s.DB()
    if req.Id != nil && *req.Id != 0 {
    	db = db.Where("id", req.Id)
    }
    if req.BeginTime != nil && *req.BeginTime > 0 && req.EndTime != nil && *req.EndTime > 0{
    		db = db.Where("created_at between ? and ?", *req.BeginTime, *req.EndTime)
    }

    // 导出处理
    fields := make([]string, 0)
    if req.Fields != nil {
    	fields = strings.Split(*req.Fields, ",")
    }
    var url string
    colInfo := s.GetColumnCommentFromStruct(model.MagicExtraStakeOrderOpsRecord{})
    var err error
    if req.IsExport != nil && *req.IsExport {
    	if len(fields) == 0 {
    		for _, col := range colInfo {
    			fields = append(fields, col.Field)
    		}
    	}
    	core.Log.Infof("导出的字段:%s", fields)
    	url, err = base.ExportCsv[model.MagicExtraStakeOrderOpsRecord](db, fields, colInfo)
    	if err != nil {
    		response.Resp(c, err.Error())
    		return
    	}
    }
    resp, err := base.NewQueryBaseHandler(model.NewMagicExtraStakeOrderOpsRecord()).List(db, req)
    if err != nil {
    	response.Resp(c, err.Error())
    	return
    }
    response.Resp(c, map[string]interface{}{
    		"url":      url,
    		"cols":     colInfo,
    		"list":     resp.List,
    		"paging":   resp.Paging,
    	})
}

func (s *MagicExtraStakeOrderOpsRecordService) Comment(c *gin.Context) {
	s.SetDbAlias("app")
	dbs, err := s.GetColumnComment("app", model.NewMagicExtraStakeOrderOpsRecord().TableName())
	if err != nil {
		response.Resp(c, "获取失败")
		return
	} else {
		response.Resp(c, gin.H{"dbs": dbs})
		return
	}
}