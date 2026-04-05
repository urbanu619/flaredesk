package app

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/demdxx/gocast"
	"github.com/gin-gonic/gin"
	"go_server/model/biz_modules/app"
	"go_server/model/common/response"
	"go_server/service/base"
)

type CfZoneService struct {
	base.BizCommonService
}

// Sync 从 CF API 同步域名到本地数据库
func (s *CfZoneService) Sync(c *gin.Context) {
	s.SetDbAlias("app")
	accountId := gocast.ToInt64(c.Query("account_id"))
	if accountId == 0 {
		response.Resp(c, "account_id 必填")
		return
	}

	// 复用 CfDnsService 的 token 获取逻辑
	account, ok := base.GetOne[app.CfAccount](s.DB(), "id", accountId)
	if !ok || account.Status != 1 {
		response.Resp(c, "账号不存在或已禁用")
		return
	}

	dns := &CfDnsService{}
	dns.SetDbAlias("app")

	// 翻页拉取所有 Zone
	var allZones []map[string]interface{}
	page := 1
	for {
		url := fmt.Sprintf("%s/zones?per_page=50&page=%d", cfAPIBase, page)
		data, err := dns.cfRequest("GET", url, account.ApiToken, nil)
		if err != nil {
			response.Resp(c, err.Error())
			return
		}
		var result map[string]interface{}
		json.Unmarshal(data, &result)
		cfSuccess, _ := result["success"].(bool)
		if !cfSuccess {
			response.Resp(c, result)
			return
		}
		raw, _ := result["result"].([]interface{})
		for _, r := range raw {
			if z, ok := r.(map[string]interface{}); ok {
				allZones = append(allZones, z)
			}
		}
		resultInfo, _ := result["result_info"].(map[string]interface{})
		if int64(page) >= gocast.ToInt64(resultInfo["total_pages"]) {
			break
		}
		page++
	}

	now := time.Now().Unix()
	var created, updated int

	for _, z := range allZones {
		zoneId, _ := z["id"].(string)
		name, _ := z["name"].(string)
		status, _ := z["status"].(string)
		paused, _ := z["paused"].(bool)
		activatedOn, _ := z["activated_on"].(string)

		planName := ""
		if plan, ok := z["plan"].(map[string]interface{}); ok {
			planName, _ = plan["name"].(string)
		}

		nsBytes, _ := json.Marshal(z["name_servers"])

		existing, exists := base.GetOne[app.CfZone](s.DB(), "zone_id", zoneId)
		if exists {
			s.DB().Model(existing).Updates(map[string]interface{}{
				"name":         name,
				"status":       status,
				"paused":       paused,
				"plan_name":    planName,
				"name_servers": string(nsBytes),
				"activated_on": activatedOn,
				"updated_at":   now,
			})
			updated++
		} else {
			s.DB().Create(&app.CfZone{
				AccountId:   accountId,
				ZoneId:      zoneId,
				Name:        name,
				Status:      status,
				Paused:      paused,
				PlanName:    planName,
				NameServers: string(nsBytes),
				ActivatedOn: activatedOn,
				CreatedAt:   now,
				UpdatedAt:   now,
			})
			created++
		}
	}

	response.Resp(c, map[string]interface{}{
		"total":   len(allZones),
		"created": created,
		"updated": updated,
	})
}

// Find 查询本地 Zone 列表
func (s *CfZoneService) Find(c *gin.Context) {
	s.SetDbAlias("app")
	type request[T any] struct {
		base.ListRequest[T]
		AccountId int64  `form:"accountId"`
		Name      string `form:"name"`
		Status    string `form:"status"`
	}
	req := new(request[app.CfZone])
	if err := c.BindQuery(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	db := s.DB()
	if req.AccountId > 0 {
		db = db.Where("account_id = ?", req.AccountId)
	}
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}
	db = db.Order("name ASC")
	resp, err := base.NewQueryBaseHandler(app.NewCfZone()).List(db, req)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c, resp)
}

// UpdateRemark 更新备注
func (s *CfZoneService) UpdateRemark(c *gin.Context) {
	s.SetDbAlias("app")
	type request struct {
		Id     int64  `json:"id" binding:"required"`
		Remark string `json:"remark"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	zone, ok := base.GetOne[app.CfZone](s.DB(), "id", req.Id)
	if !ok {
		response.Resp(c, "记录不存在")
		return
	}
	s.DB().Model(zone).Updates(map[string]interface{}{
		"remark":     req.Remark,
		"updated_at": time.Now().Unix(),
	})
	response.Resp(c)
}

// All 返回某账号下所有 Zone（不分页，供批量操作使用）
func (s *CfZoneService) All(c *gin.Context) {
	s.SetDbAlias("app")
	accountId := gocast.ToInt64(c.Query("account_id"))
	if accountId == 0 {
		response.Resp(c, "account_id 必填")
		return
	}
	var zones []app.CfZone
	s.DB().Where("account_id = ?", accountId).Order("name ASC").Find(&zones)
	response.Resp(c, zones)
}
