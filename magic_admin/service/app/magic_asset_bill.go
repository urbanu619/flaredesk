package app

import (
	"go_server/base/core"
	model "go_server/model/biz_modules/app"
	"go_server/model/common/response"
	"go_server/service/base"
	"strings"

	"github.com/gin-gonic/gin"
)

type MagicAssetBillService struct {
	base.BizCommonService
}

func (s *MagicAssetBillService) Get(c *gin.Context) {
	s.SetDbAlias("app")
	base.NewBaseHandler(model.NewMagicAssetBill()).Get(c, s.DB())
}

func (s *MagicAssetBillService) Find(c *gin.Context) {
	s.SetDbAlias("app")
	type request[T any] struct {
		base.ListRequest[T]
		Id             *int64  `form:"id"`
		IsExport       *bool   `form:"isExport"`
		Fields         *string `form:"fields"` // 指定返回字段 , 分割
		UserId         *int64  `form:"userId"`
		Symbol         *string `form:"symbol"`
		BeginTime      *int64  `form:"beginTime"`
		EndTime        *int64  `form:"endTime"`
		Uid            *string `form:"uid"`
		BusinessNumber *int64  `form:"businessNumber"`
	}
	req := new(request[model.MagicAssetBill])
	if err := c.BindQuery(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	db := s.DB()
	if req.Id != nil && *req.Id != 0 {
		db = db.Where("id", req.Id)
	}
	if req.UserId != nil && *req.UserId != 0 {
		db = db.Where("user_id", req.UserId)
	}
	if req.Symbol != nil && *req.Symbol != "" {
		db = db.Where("symbol", req.Symbol)
	}
	if req.BeginTime != nil && *req.BeginTime > 0 && req.EndTime != nil && *req.EndTime > 0 {
		db = db.Where("created_at between ? and ?", *req.BeginTime, *req.EndTime)
	}
	if req.Uid != nil && *req.Uid != "" {
		db = db.Where("uid", req.Uid)
	}
	if req.BusinessNumber != nil && *req.BusinessNumber != 0 {
		db = db.Where("business_number", req.BusinessNumber)
	}
	// 导出处理
	fields := make([]string, 0)
	if req.Fields != nil {
		fields = strings.Split(*req.Fields, ",")
	}
	var url string
	colInfo := s.GetColumnCommentFromStruct(model.MagicAssetBill{})
	var err error
	if req.IsExport != nil && *req.IsExport {
		if len(fields) == 0 {
			for _, col := range colInfo {
				fields = append(fields, col.Field)
			}
		}
		core.Log.Infof("导出的字段:%s", fields)
		url, err = base.ExportCsv[model.MagicAssetBill](db, fields, colInfo)
		if err != nil {
			response.Resp(c, err.Error())
			return
		}
	}
	resp, err := base.NewQueryBaseHandler(model.NewMagicAssetBill()).List(db, req)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c, map[string]interface{}{
		"url":    url,
		"cols":   colInfo,
		"list":   resp.List,
		"paging": resp.Paging,
		"businessOptions": map[string]string{
			//MagicAssetBsMap
			"101": "充值",
			"102": "提现申请",
			"103": "提现成功扣除",
			"104": "提现失败",
			"105": "系统增减资产",
			"106": "收益提取",
			"107": "质押",
			"108": "赎回",
			"109": "升级奖励",
			"110": "资产迁移",
		},
	})
}

func (s *MagicAssetBillService) Comment(c *gin.Context) {
	s.SetDbAlias("app")
	dbs, err := s.GetColumnComment("app", model.NewMagicAssetBill().TableName())
	if err != nil {
		response.Resp(c, "获取失败")
		return
	} else {
		response.Resp(c, gin.H{"dbs": dbs})
		return
	}
}
