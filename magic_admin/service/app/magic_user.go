package app

import (
	"go_server/base/core"
	model "go_server/model/biz_modules/app"
	"go_server/model/common/response"
	"go_server/service/base"
	"strings"

	"github.com/gin-gonic/gin"
)

type MagicUserService struct {
	base.BizCommonService
}

func (s *MagicUserService) Get(c *gin.Context) {
	s.SetDbAlias("app")
	base.NewBaseHandler(model.NewMagicUser()).Get(c, s.DB())
}

func (s *MagicUserService) Find(c *gin.Context) {
	s.SetDbAlias("app")
	type request[T any] struct {
		base.ListRequest[T]
		Id                      *int64  `form:"id"`
		IsExport                *bool   `form:"isExport"`
		Fields                  *string `form:"fields"` // 指定返回字段 , 分割
		Uid                     *string `form:"uid"`
		Code                    *string `form:"code"`
		Enable                  *int8   `form:"enable"`
		BeginTime               *int64  `form:"beginTime"`
		EndTime                 *int64  `form:"endTime"`
		SubsLockWithdraw        *int8   `form:"subsLockWithdraw"`
		LockRedeem              *int8   `form:"lockRedeem"`
		WithoutStakeRanking     *int8   `form:"withoutStakeRanking"`
		WithoutFewRegionRanking *int8   `form:"withoutFewRegionRanking"`
		LockStake               *int8   `form:"lockStake"`
		LockStakeProfit         *int8   `form:"lockStakeProfit"`
		UnCollection            *int8   `form:"unCollection"`
	}
	req := new(request[model.MagicUser])
	if err := c.BindQuery(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	db := s.DB()
	if req.Id != nil && *req.Id != 0 {
		db = db.Where("id", req.Id)
	}

	if req.Uid != nil && *req.Uid != "" {
		db = db.Where("uid", req.Uid)
	}
	if req.Code != nil && *req.Code != "" {
		db = db.Where("code", req.Code)
	}
	if req.Enable != nil {
		db = db.Where("enable", req.Enable)
	}

	if req.SubsLockWithdraw != nil {
		db = db.Where("subs_lock_withdraw", req.SubsLockWithdraw)
	}
	if req.LockRedeem != nil {
		db = db.Where("lock_redeem", req.LockRedeem)
	}
	if req.WithoutStakeRanking != nil {
		db = db.Where("without_stake_ranking", req.WithoutStakeRanking)
	}
	if req.WithoutFewRegionRanking != nil {
		db = db.Where("without_few_region_ranking", req.WithoutFewRegionRanking)
	}
	if req.LockStake != nil {
		db = db.Where("lock_stake", req.LockStake)
	}
	if req.LockStakeProfit != nil {
		db = db.Where("lock_stake_profit", req.LockStakeProfit)
	}
	if req.UnCollection != nil {
		db = db.Where("un_collection", req.UnCollection)
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
	colInfo := s.GetColumnCommentFromStruct(model.MagicUser{})
	var err error
	if req.IsExport != nil && *req.IsExport {
		if len(fields) == 0 {
			for _, col := range colInfo {
				fields = append(fields, col.Field)
			}
		}
		core.Log.Infof("导出的字段:%s", fields)
		url, err = base.ExportCsv[model.MagicUser](db, fields, colInfo)
		if err != nil {
			response.Resp(c, err.Error())
			return
		}
	}
	resp, err := base.NewQueryBaseHandler(model.NewMagicUser()).List(db, req)
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

func (s *MagicUserService) Comment(c *gin.Context) {
	s.SetDbAlias("app")
	dbs, err := s.GetColumnComment("app", model.NewMagicUser().TableName())
	if err != nil {
		response.Resp(c, "获取失败")
		return
	} else {
		response.Resp(c, gin.H{"dbs": dbs})
		return
	}
}
