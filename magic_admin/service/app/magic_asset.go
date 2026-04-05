package app

import (
	"go_server/base/core"
	model "go_server/model/biz_modules/app"
	"go_server/model/common/response"
	"go_server/service/base"
	"strings"

	"github.com/gin-gonic/gin"
)

type MagicAssetService struct {
	base.BizCommonService
}

func (s *MagicAssetService) Get(c *gin.Context) {
	s.SetDbAlias("app")
	base.NewBaseHandler(model.NewMagicAsset()).Get(c, s.DB())
}

func (s *MagicAssetService) Find(c *gin.Context) {
	s.SetDbAlias("app")
	type request[T any] struct {
		base.ListRequest[T]
		Id        *int64  `form:"id"`
		IsExport  *bool   `form:"isExport"`
		Fields    *string `form:"fields"` // 指定返回字段 , 分割
		UserId    *int64  `form:"userId"`
		Symbol    *string `form:"symbol"`
		BeginTime *int64  `form:"beginTime"`
		EndTime   *int64  `form:"endTime"`
		Uid       *string `form:"uid"`
	}
	req := new(request[model.MagicAsset])
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

	// 导出处理
	fields := make([]string, 0)
	if req.Fields != nil {
		fields = strings.Split(*req.Fields, ",")
	}
	var url string
	colInfo := s.GetColumnCommentFromStruct(model.MagicAsset{})
	var err error
	if req.IsExport != nil && *req.IsExport {
		if len(fields) == 0 {
			for _, col := range colInfo {
				fields = append(fields, col.Field)
			}
		}
		core.Log.Infof("导出的字段:%s", fields)
		url, err = base.ExportCsv[model.MagicAsset](db, fields, colInfo)
		if err != nil {
			response.Resp(c, err.Error())
			return
		}
	}
	resp, err := base.NewQueryBaseHandler(model.NewMagicAsset()).List(db, req)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}

	// 批量查询用户绑定的 Telegram 用户名
	userIds := make([]int64, 0, len(resp.List))
	for _, item := range resp.List {
		userIds = append(userIds, item.UserId)
	}
	bindMap := make(map[int64]string) // userId -> telegramUsername
	if len(userIds) > 0 {
		var binds []model.TgUserBind
		s.DB().Where("user_id IN ? AND bind_status = 1", userIds).Find(&binds)
		for _, b := range binds {
			bindMap[b.UserId] = b.TelegramUsername
		}
	}

	// 构造带 telegramUsername 的返回列表
	type AssetWithBind struct {
		model.MagicAsset
		TelegramUsername string `json:"telegramUsername"`
	}
	listWithBind := make([]AssetWithBind, 0, len(resp.List))
	for _, item := range resp.List {
		listWithBind = append(listWithBind, AssetWithBind{
			MagicAsset:       item,
			TelegramUsername: bindMap[item.UserId],
		})
	}

	response.Resp(c, map[string]interface{}{
		"url":    url,
		"cols":   colInfo,
		"list":   listWithBind,
		"paging": resp.Paging,
	})
}

func (s *MagicAssetService) Comment(c *gin.Context) {
	s.SetDbAlias("app")
	dbs, err := s.GetColumnComment("app", model.NewMagicAsset().TableName())
	if err != nil {
		response.Resp(c, "获取失败")
		return
	} else {
		response.Resp(c, gin.H{"dbs": dbs})
		return
	}
}
