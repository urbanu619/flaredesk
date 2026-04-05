package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/demdxx/gocast"
	"github.com/gin-gonic/gin"
	"go_server/model/biz_modules/app"
	"go_server/model/common/response"
	"go_server/service/base"
)

const cfAPIBase = "https://api.cloudflare.com/client/v4"

type CfDnsService struct {
	base.BizCommonService
}

// cfRequest 向 CF API 发送请求
func (s *CfDnsService) cfRequest(method, url, token string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// getToken 从数据库取 API Token
func (s *CfDnsService) getToken(accountId int64) (string, bool) {
	account, ok := base.GetOne[app.CfAccount](s.DB(), "id", accountId)
	if !ok || account.Status != 1 {
		return "", false
	}
	return account.ApiToken, true
}

// Zones 获取域名列表（自动翻页，返回全部）
func (s *CfDnsService) Zones(c *gin.Context) {
	s.SetDbAlias("app")
	accountId := gocast.ToInt64(c.Query("account_id"))
	if accountId == 0 {
		response.Resp(c, "account_id 必填")
		return
	}
	token, ok := s.getToken(accountId)
	if !ok {
		response.Resp(c, "账号不存在或已禁用")
		return
	}

	var allZones []interface{}
	page := 1
	for {
		url := fmt.Sprintf("%s/zones?per_page=50&page=%d", cfAPIBase, page)
		data, err := s.cfRequest("GET", url, token, nil)
		if err != nil {
			response.Resp(c, err.Error())
			return
		}
		var result map[string]interface{}
		if err := json.Unmarshal(data, &result); err != nil {
			response.Resp(c, err.Error())
			return
		}
		cfSuccess, _ := result["success"].(bool)
		if !cfSuccess {
			response.Resp(c, result)
			return
		}
		zones, _ := result["result"].([]interface{})
		allZones = append(allZones, zones...)

		resultInfo, _ := result["result_info"].(map[string]interface{})
		totalPages := gocast.ToInt64(resultInfo["total_pages"])
		if int64(page) >= totalPages {
			break
		}
		page++
	}

	response.Resp(c, map[string]interface{}{
		"result":      allZones,
		"total_count": len(allZones),
		"success":     true,
	})
}

// Records 获取 DNS 记录列表
func (s *CfDnsService) Records(c *gin.Context) {
	s.SetDbAlias("app")
	accountId := gocast.ToInt64(c.Query("account_id"))
	zoneId := c.Query("zone_id")
	if accountId == 0 || zoneId == "" {
		response.Resp(c, "account_id 和 zone_id 必填")
		return
	}
	token, ok := s.getToken(accountId)
	if !ok {
		response.Resp(c, "账号不存在或已禁用")
		return
	}

	// 支持按 type / name 过滤
	query := fmt.Sprintf("%s/zones/%s/dns_records?per_page=100", cfAPIBase, zoneId)
	if t := c.Query("type"); t != "" {
		query += "&type=" + t
	}
	if n := c.Query("name"); n != "" {
		query += "&name=" + n
	}

	data, err := s.cfRequest("GET", query, token, nil)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}

	var result map[string]interface{}
	json.Unmarshal(data, &result)
	response.Resp(c, result)
}

// CreateRecord 创建单条 DNS 记录
func (s *CfDnsService) CreateRecord(c *gin.Context) {
	s.SetDbAlias("app")
	type request struct {
		AccountId int64  `json:"accountId" binding:"required"`
		ZoneId    string `json:"zoneId" binding:"required"`
		Type      string `json:"type" binding:"required"`
		Name      string `json:"name" binding:"required"`
		Content   string `json:"content" binding:"required"`
		Ttl       int    `json:"ttl"`
		Proxied   bool   `json:"proxied"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	token, ok := s.getToken(req.AccountId)
	if !ok {
		response.Resp(c, "账号不存在或已禁用")
		return
	}

	ttl := req.Ttl
	if ttl == 0 {
		ttl = 1 // auto
	}
	payload := map[string]interface{}{
		"type":    req.Type,
		"name":    req.Name,
		"content": req.Content,
		"ttl":     ttl,
		"proxied": req.Proxied,
	}

	url := fmt.Sprintf("%s/zones/%s/dns_records", cfAPIBase, req.ZoneId)
	data, err := s.cfRequest("POST", url, token, payload)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}

	var result map[string]interface{}
	json.Unmarshal(data, &result)
	response.Resp(c, result)
}

// BatchCreateRecord 批量创建 DNS 记录
func (s *CfDnsService) BatchCreateRecord(c *gin.Context) {
	s.SetDbAlias("app")
	type dnsItem struct {
		Type    string `json:"type"`
		Name    string `json:"name"`
		Content string `json:"content"`
		Ttl     int    `json:"ttl"`
		Proxied bool   `json:"proxied"`
	}
	type request struct {
		AccountId int64     `json:"accountId" binding:"required"`
		ZoneId    string    `json:"zoneId" binding:"required"`
		Records   []dnsItem `json:"records" binding:"required"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	token, ok := s.getToken(req.AccountId)
	if !ok {
		response.Resp(c, "账号不存在或已禁用")
		return
	}

	var results []map[string]interface{}
	var failCount int
	url := fmt.Sprintf("%s/zones/%s/dns_records", cfAPIBase, req.ZoneId)

	for _, item := range req.Records {
		ttl := item.Ttl
		if ttl == 0 {
			ttl = 1
		}
		payload := map[string]interface{}{
			"type":    item.Type,
			"name":    item.Name,
			"content": item.Content,
			"ttl":     ttl,
			"proxied": item.Proxied,
		}
		data, err := s.cfRequest("POST", url, token, payload)
		if err != nil {
			failCount++
			continue
		}
		var r map[string]interface{}
		json.Unmarshal(data, &r)
		results = append(results, r)
	}

	response.Resp(c, map[string]interface{}{
		"total":   len(req.Records),
		"success": len(results),
		"fail":    failCount,
		"results": results,
	})
}

// UpdateRecord 更新 DNS 记录
func (s *CfDnsService) UpdateRecord(c *gin.Context) {
	s.SetDbAlias("app")
	type request struct {
		AccountId int64  `json:"accountId" binding:"required"`
		ZoneId    string `json:"zoneId" binding:"required"`
		RecordId  string `json:"recordId" binding:"required"`
		Type      string `json:"type"`
		Name      string `json:"name"`
		Content   string `json:"content"`
		Ttl       int    `json:"ttl"`
		Proxied   *bool  `json:"proxied"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	token, ok := s.getToken(req.AccountId)
	if !ok {
		response.Resp(c, "账号不存在或已禁用")
		return
	}

	payload := make(map[string]interface{})
	if req.Type != "" {
		payload["type"] = req.Type
	}
	if req.Name != "" {
		payload["name"] = req.Name
	}
	if req.Content != "" {
		payload["content"] = req.Content
	}
	if req.Ttl > 0 {
		payload["ttl"] = req.Ttl
	}
	if req.Proxied != nil {
		payload["proxied"] = *req.Proxied
	}

	url := fmt.Sprintf("%s/zones/%s/dns_records/%s", cfAPIBase, req.ZoneId, req.RecordId)
	data, err := s.cfRequest("PATCH", url, token, payload)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}

	var result map[string]interface{}
	json.Unmarshal(data, &result)
	response.Resp(c, result)
}

// DeleteRecord 删除 DNS 记录
func (s *CfDnsService) DeleteRecord(c *gin.Context) {
	s.SetDbAlias("app")
	accountId := gocast.ToInt64(c.Query("account_id"))
	zoneId := c.Query("zone_id")
	recordId := c.Query("record_id")
	if accountId == 0 || zoneId == "" || recordId == "" {
		response.Resp(c, "account_id、zone_id、record_id 必填")
		return
	}
	token, ok := s.getToken(accountId)
	if !ok {
		response.Resp(c, "账号不存在或已禁用")
		return
	}

	url := fmt.Sprintf("%s/zones/%s/dns_records/%s", cfAPIBase, zoneId, recordId)
	data, err := s.cfRequest("DELETE", url, token, nil)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}

	var result map[string]interface{}
	json.Unmarshal(data, &result)
	response.Resp(c, result)
}

// getZoneRecords 获取 zone 的所有 DNS 记录（按 type/name 可选过滤）
func (s *CfDnsService) getZoneRecords(zoneId, token, recType, name string) ([]map[string]interface{}, error) {
	query := fmt.Sprintf("%s/zones/%s/dns_records?per_page=100", cfAPIBase, zoneId)
	if recType != "" {
		query += "&type=" + recType
	}
	if name != "" {
		query += "&name=" + name
	}
	data, err := s.cfRequest("GET", query, token, nil)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	json.Unmarshal(data, &result)
	raw, _ := result["result"].([]interface{})
	var records []map[string]interface{}
	for _, r := range raw {
		if rec, ok := r.(map[string]interface{}); ok {
			records = append(records, rec)
		}
	}
	return records, nil
}

// CrossZoneDeleteRecords 跨域名批量删除匹配的 DNS 记录
func (s *CfDnsService) CrossZoneDeleteRecords(c *gin.Context) {
	s.SetDbAlias("app")
	type request struct {
		AccountId int64    `json:"accountId" binding:"required"`
		ZoneIds   []string `json:"zoneIds" binding:"required"`
		Type      string   `json:"type"`
		Name      string   `json:"name"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	if req.Type == "" && req.Name == "" {
		response.Resp(c, "type 和 name 至少填一个")
		return
	}
	token, ok := s.getToken(req.AccountId)
	if !ok {
		response.Resp(c, "账号不存在或已禁用")
		return
	}

	var deleted, failed int
	for _, zoneId := range req.ZoneIds {
		records, err := s.getZoneRecords(zoneId, token, req.Type, req.Name)
		if err != nil {
			failed++
			continue
		}
		for _, rec := range records {
			recordId, _ := rec["id"].(string)
			if recordId == "" {
				continue
			}
			delURL := fmt.Sprintf("%s/zones/%s/dns_records/%s", cfAPIBase, zoneId, recordId)
			if _, err := s.cfRequest("DELETE", delURL, token, nil); err != nil {
				failed++
			} else {
				deleted++
			}
		}
	}
	response.Resp(c, map[string]interface{}{
		"zones":   len(req.ZoneIds),
		"deleted": deleted,
		"fail":    failed,
	})
}

// CrossZoneToggleProxy 跨域名批量切换橙云
func (s *CfDnsService) CrossZoneToggleProxy(c *gin.Context) {
	s.SetDbAlias("app")
	type request struct {
		AccountId int64    `json:"accountId" binding:"required"`
		ZoneIds   []string `json:"zoneIds" binding:"required"`
		Type      string   `json:"type"`
		Name      string   `json:"name"`
		Proxied   bool     `json:"proxied"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	token, ok := s.getToken(req.AccountId)
	if !ok {
		response.Resp(c, "账号不存在或已禁用")
		return
	}

	var updated, failed int
	for _, zoneId := range req.ZoneIds {
		records, err := s.getZoneRecords(zoneId, token, req.Type, req.Name)
		if err != nil {
			failed++
			continue
		}
		for _, rec := range records {
			recordId, _ := rec["id"].(string)
			recType, _ := rec["type"].(string)
			if recordId == "" {
				continue
			}
			// 只有 A/AAAA/CNAME 支持 proxied
			if recType != "A" && recType != "AAAA" && recType != "CNAME" {
				continue
			}
			patchURL := fmt.Sprintf("%s/zones/%s/dns_records/%s", cfAPIBase, zoneId, recordId)
			if _, err := s.cfRequest("PATCH", patchURL, token, map[string]interface{}{"proxied": req.Proxied}); err != nil {
				failed++
			} else {
				updated++
			}
		}
	}
	response.Resp(c, map[string]interface{}{
		"zones":   len(req.ZoneIds),
		"updated": updated,
		"fail":    failed,
	})
}

// ToggleProxy 批量切换橙云（proxied）
func (s *CfDnsService) ToggleProxy(c *gin.Context) {
	s.SetDbAlias("app")
	type request struct {
		AccountId int64    `json:"accountId" binding:"required"`
		ZoneId    string   `json:"zoneId" binding:"required"`
		RecordIds []string `json:"recordIds" binding:"required"`
		Proxied   bool     `json:"proxied"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	token, ok := s.getToken(req.AccountId)
	if !ok {
		response.Resp(c, "账号不存在或已禁用")
		return
	}

	var failCount int
	for _, recordId := range req.RecordIds {
		url := fmt.Sprintf("%s/zones/%s/dns_records/%s", cfAPIBase, req.ZoneId, recordId)
		_, err := s.cfRequest("PATCH", url, token, map[string]interface{}{"proxied": req.Proxied})
		if err != nil {
			failCount++
		}
	}

	response.Resp(c, map[string]interface{}{
		"total":   len(req.RecordIds),
		"success": len(req.RecordIds) - failCount,
		"fail":    failCount,
	})
}
