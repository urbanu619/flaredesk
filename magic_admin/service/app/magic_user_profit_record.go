package app

import (
	"go_server/base/core"
	model "go_server/model/biz_modules/app"
	"go_server/model/common/response"
	"go_server/service/base"
	"strings"

	"github.com/gin-gonic/gin"
)

type MagicUserProfitRecordService struct {
	base.BizCommonService
}

func (s *MagicUserProfitRecordService) Get(c *gin.Context) {
	s.SetDbAlias("app")
	base.NewBaseHandler(model.NewMagicUserProfitRecord()).Get(c, s.DB())
}

func (s *MagicUserProfitRecordService) Find(c *gin.Context) {
	s.SetDbAlias("app")
	type request[T any] struct {
		base.ListRequest[T]
		Id             *int64  `form:"id"`
		IsExport       *bool   `form:"isExport"`
		Fields         *string `form:"fields"` // 指定返回字段 , 分割
		State          *string `form:"state"`
		rewardUserId   *string `form:"rewardUserId"`
		RewardUserUid  *string `form:"rewardUserUid"`
		BeginTime      *int64  `form:"beginTime"`
		EndTime        *int64  `form:"endTime"`
		BusinessNumber *int64  `form:"businessNumber"`
	}
	req := new(request[model.MagicUserProfitRecord])
	if err := c.BindQuery(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	db := s.DB()
	if req.Id != nil && *req.Id != 0 {
		db = db.Where("id", req.Id)
	}
	if req.State != nil && *req.State != "" {
		db = db.Where("state", req.State)
	}
	if req.rewardUserId != nil && *req.rewardUserId != "" {
		db = db.Where("reward_user_id", req.rewardUserId)
	}
	if req.RewardUserUid != nil && *req.RewardUserUid != "" {
		db = db.Where("reward_user_uid", req.RewardUserUid)
	}
	if req.BusinessNumber != nil && *req.BusinessNumber != 0 {
		db = db.Where("business_number", req.BusinessNumber)
	}
	if req.BeginTime != nil && *req.BeginTime > 0 && req.EndTime != nil && *req.EndTime > 0 {
		db = db.Where("created_at between ? and ?", *req.BeginTime, *req.EndTime)
	}

	// 导出处理
	fields := make([]string, 0)
	if req.Fields != nil {
		fields = strings.Split(*req.Fields, ",")
	}
	var url string
	colInfo := s.GetColumnCommentFromStruct(model.MagicUserProfitRecord{})
	var err error
	if req.IsExport != nil && *req.IsExport {
		if len(fields) == 0 {
			for _, col := range colInfo {
				fields = append(fields, col.Field)
			}
		}
		core.Log.Infof("导出的字段:%s", fields)
		url, err = base.ExportCsv[model.MagicUserProfitRecord](db, fields, colInfo)
		if err != nil {
			response.Resp(c, err.Error())
			return
		}
	}
	resp, err := base.NewQueryBaseHandler(model.NewMagicUserProfitRecord()).List(db, req)
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
			//NodeAssetBsMap
			"151": "合伙人节点分红",
			"152": "生态节点分红",
			"153": "合伙人节点返佣",
			"154": "生态节点返佣",
			//ProfitBsMap
			"301": "活期静态收益",
			"302": "动态均分收益",
			"303": "动态加权收益",
			"304": "小区业绩排行榜收益",
			"305": "个人质押排行榜收益",
			"306": "领取收益",
			"307": "收益迁移",
		},
	})
}

func (s *MagicUserProfitRecordService) Comment(c *gin.Context) {
	s.SetDbAlias("app")
	dbs, err := s.GetColumnComment("app", model.NewMagicUserProfitRecord().TableName())
	if err != nil {
		response.Resp(c, "获取失败")
		return
	} else {
		response.Resp(c, gin.H{"dbs": dbs})
		return
	}
}
