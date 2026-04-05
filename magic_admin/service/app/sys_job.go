package app

import (
	"go_server/base/core"
	model "go_server/model/biz_modules/app"
	"go_server/model/common/response"
	"go_server/service/base"
	"strings"

	"github.com/gin-gonic/gin"
)

type SysJobService struct {
	base.BizCommonService
}

func (s *SysJobService) Get(c *gin.Context) {
	s.SetDbAlias("app")
	base.NewBaseHandler(model.NewSysJob()).Get(c, s.DB())
}

func (s *SysJobService) Find(c *gin.Context) {
	s.SetDbAlias("app")
	type request[T any] struct {
		base.ListRequest[T]
		Id        *int64  `form:"id"`
		IsExport  *bool   `form:"isExport"`
		Fields    *string `form:"fields"` // 指定返回字段 , 分割
		BeginTime *int64  `form:"beginTime"`
		EndTime   *int64  `form:"endTime"`
		JobDate   *string `form:"jobDate"`
		State     *string `form:"state"`
	}
	req := new(request[model.SysJob])
	if err := c.BindQuery(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	db := s.DB()
	if req.Id != nil && *req.Id != 0 {
		db = db.Where("id", req.Id)
	}
	if req.BeginTime != nil && *req.BeginTime > 0 && req.EndTime != nil && *req.EndTime > 0 {
		db = db.Where("created_at between ? and ?", *req.BeginTime, *req.EndTime)
	}
	if req.JobDate != nil && *req.JobDate != "" {
		db = db.Where("job_date", req.JobDate)
	}
	if req.State != nil && *req.State != "" {
		db = db.Where("state", req.State)
	}
	// 导出处理
	fields := make([]string, 0)
	if req.Fields != nil {
		fields = strings.Split(*req.Fields, ",")
	}
	var url string
	colInfo := s.GetColumnCommentFromStruct(model.SysJob{})
	var err error
	if req.IsExport != nil && *req.IsExport {
		if len(fields) == 0 {
			for _, col := range colInfo {
				fields = append(fields, col.Field)
			}
		}
		core.Log.Infof("导出的字段:%s", fields)
		url, err = base.ExportCsv[model.SysJob](db, fields, colInfo)
		if err != nil {
			response.Resp(c, err.Error())
			return
		}
	}
	resp, err := base.NewQueryBaseHandler(model.NewSysJob()).List(db, req)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c, map[string]interface{}{
		"url":    url,
		"cols":   colInfo,
		"list":   resp.List,
		"paging": resp.Paging,
	})
}

func (s *SysJobService) Comment(c *gin.Context) {
	s.SetDbAlias("app")
	dbs, err := s.GetColumnComment("app", model.NewSysJob().TableName())
	if err != nil {
		response.Resp(c, "获取失败")
		return
	} else {
		response.Resp(c, gin.H{"dbs": dbs})
		return
	}
}
