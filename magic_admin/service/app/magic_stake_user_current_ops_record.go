package app

import (
	"go_server/base/core"
	model "go_server/model/biz_modules/app"
	"go_server/model/common/response"
	"go_server/service/base"
	"strings"

	"github.com/gin-gonic/gin"
)

type MagicStakeUserCurrentOpsRecordService struct {
	base.BizCommonService
}

func (s *MagicStakeUserCurrentOpsRecordService) Get(c *gin.Context) {
	s.SetDbAlias("app")
	base.NewBaseHandler(model.NewMagicStakeUserCurrentOpsRecord()).Get(c, s.DB())
}

func (s *MagicStakeUserCurrentOpsRecordService) Find(c *gin.Context) {
	s.SetDbAlias("app")
	type request[T any] struct {
		base.ListRequest[T]
		Id             *int64  `form:"id"`
		IsExport       *bool   `form:"isExport"`
		Fields         *string `form:"fields"` // 指定返回字段 , 分割
		UserId         *int64  `form:"userId"`
		BeginTime      *int64  `form:"beginTime"`
		EndTime        *int64  `form:"endTime"`
		Uid            *string `form:"uid"`
		TopUid         *string `form:"topUid"` // 顶层上级UID
		BusinessNumber *int64  `form:"businessNumber"`
	}
	req := new(request[model.MagicStakeUserCurrentOpsRecord])
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
	if req.Uid != nil && *req.Uid != "" {
		db = db.Where("uid", req.Uid)
	}
	if req.BeginTime != nil && *req.BeginTime > 0 && req.EndTime != nil && *req.EndTime > 0 {
		db = db.Where("created_at between ? and ?", *req.BeginTime, *req.EndTime)
	}
	if req.BusinessNumber != nil && *req.BusinessNumber != 0 {
		db = db.Where("business_number", req.BusinessNumber)
	}

	// TopUid 筛选逻辑
	if req.TopUid != nil && *req.TopUid != "" {
		// 1. 先查出 TopUid 对应的 int64 ID
		var topUser model.MagicUser
		if err := s.DB().Model(&model.MagicUser{}).Where("uid", req.TopUid).Select("id").First(&topUser).Error; err == nil {
			// 2. 查出该用户伞下所有用户的 ID
			var userIds []int64
			// 使用 JSON_CONTAINS 查找 parent_ids 包含 topUser.Id 的用户
			if err := s.DB().Model(&model.MagicUser{}).Where("JSON_CONTAINS(parent_ids, CAST(? AS JSON), '$')", topUser.Id).Pluck("id", &userIds).Error; err == nil && len(userIds) > 0 {
				// 3. 使用 IN 查询
				db = db.Where("user_id IN ?", userIds)
			} else {
				// 查无下级或查询出错，返回空
				db = db.Where("1 = 0")
			}
		} else {
			// 如果找不到这个 TopUid，直接返回空数据
			db = db.Where("1 = 0")
		}
	}

	// 导出处理
	fields := make([]string, 0)
	if req.Fields != nil {
		fields = strings.Split(*req.Fields, ",")
	}
	var url string
	colInfo := s.GetColumnCommentFromStruct(model.MagicStakeUserCurrentOpsRecord{})
	var err error
	if req.IsExport != nil && *req.IsExport {
		if len(fields) == 0 {
			for _, col := range colInfo {
				fields = append(fields, col.Field)
			}
		}
		core.Log.Infof("导出的字段:%s", fields)
		url, err = base.ExportCsv[model.MagicStakeUserCurrentOpsRecord](db, fields, colInfo)
		if err != nil {
			response.Resp(c, err.Error())
			return
		}
	}
	resp, err := base.NewQueryBaseHandler(model.NewMagicStakeUserCurrentOpsRecord()).List(db, req)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c, map[string]interface{}{
		"url":    url,
		"cols":   colInfo,
		"list":   resp.List,
		"paging": resp.Paging,
		"symbol": "TD", // TD or USDT todo 后期支持多币种需要调整
		"businessOptions": map[string]string{
			//OrderBsMap
			"201": "活期质押",
			"202": "活期赎回",
			"203": "领取收益",
			"204": "复投收益",
			"205": "订单迁移",
		},
		"stakeQueueStateOptions": map[string]string{
			"waiting": "排队中",
			"success": "成功",
			"cancel":  "取消",
		},
	})
}

func (s *MagicStakeUserCurrentOpsRecordService) Comment(c *gin.Context) {
	s.SetDbAlias("app")
	dbs, err := s.GetColumnComment("app", model.NewMagicStakeUserCurrentOpsRecord().TableName())
	if err != nil {
		response.Resp(c, "获取失败")
		return
	} else {
		response.Resp(c, gin.H{"dbs": dbs})
		return
	}
}
