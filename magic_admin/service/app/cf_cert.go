package app

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go_server/model/common/response"
	"go_server/service/base"
)

// zoneIDParam 兼容前端传入字符串（CF zone id）或数字（本地库主键），避免 JSON 反序列化失败。
type zoneIDParam string

func (z *zoneIDParam) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) > 0 && data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*z = zoneIDParam(s)
		return nil
	}
	var n json.Number
	if err := json.Unmarshal(data, &n); err == nil {
		*z = zoneIDParam(n.String())
		return nil
	}
	return fmt.Errorf("zoneId 须为字符串或数字")
}

// cfOriginCAErrDetail 解析 CF 返回的 errors，并对常见认证/权限问题附加说明。
// 参见：https://developers.cloudflare.com/ssl/origin-configuration/origin-ca/#api-calls
func cfOriginCAErrDetail(cfResp map[string]interface{}) string {
	errBytes, _ := json.Marshal(cfResp["errors"])
	msg := "CF API 错误: " + string(errBytes)
	low := bytes.ToLower(errBytes)
	if bytes.Contains(low, []byte("authentication")) || bytes.Contains(errBytes, []byte("10000")) {
		msg += " — 多为权限问题：当前账号里保存的 API Token 未授予 Origin CA 所需权限。" +
			"请在 Cloudflare 为该 Token 勾选「Zone → SSL and Certificates → 编辑」，保存后把新 Token 填回本系统「Cloudflare 账号管理」。"
	}
	return msg
}

type CfCertService struct {
	base.BizCommonService
}

// getTokenByCertService 获取 API Token（复用 CfDnsService 逻辑）
func (s *CfCertService) getTokenByCertService(accountId int64) (string, bool) {
	dns := &CfDnsService{}
	dns.SetDbAlias("app")
	return dns.getToken(accountId)
}

// generateDomainCSR 为域名生成 RSA-2048 私钥和 CSR
func generateDomainCSR(domain string) (privateKeyPEM []byte, csrPEM []byte, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("生成私钥失败: %w", err)
	}

	privBuf := &bytes.Buffer{}
	if err = pem.Encode(privBuf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}); err != nil {
		return nil, nil, fmt.Errorf("编码私钥失败: %w", err)
	}

	csrTemplate := &x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   domain,
			Organization: []string{"Cloudflare Origin CA"},
		},
		DNSNames: []string{domain, "*." + domain},
	}
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, csrTemplate, privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("生成 CSR 失败: %w", err)
	}

	csrBuf := &bytes.Buffer{}
	if err = pem.Encode(csrBuf, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrDER}); err != nil {
		return nil, nil, fmt.Errorf("编码 CSR 失败: %w", err)
	}

	return privBuf.Bytes(), csrBuf.Bytes(), nil
}

// BatchGenerateCerts 批量生成 CF Origin 证书，打包为 zip 下载
func (s *CfCertService) BatchGenerateCerts(c *gin.Context) {
	s.SetDbAlias("app")

	type zoneItem struct {
		ZoneId zoneIDParam `json:"zoneId" binding:"required"`
		Domain string      `json:"domain" binding:"required"`
	}
	type reqBody struct {
		AccountId    int64      `json:"accountId" binding:"required"`
		Zones        []zoneItem `json:"zones" binding:"required"`
		ValidityDays int        `json:"validityDays"` // 365 或 5475（15年），默认 5475
		RequestType  string     `json:"requestType"`  // "origin-rsa"（默认）或 "origin-ecc"
	}
	type itemResult struct {
		Domain  string `json:"domain"`
		Success bool   `json:"success"`
		Error   string `json:"error,omitempty"`
	}

	req := new(reqBody)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}

	token, ok := s.getTokenByCertService(req.AccountId)
	if !ok {
		response.Resp(c, "账号不存在或已禁用")
		return
	}

	if req.ValidityDays == 0 {
		req.ValidityDays = 5475
	}
	if req.RequestType == "" {
		req.RequestType = "origin-rsa"
	}

	dns := &CfDnsService{}

	zipBuf := &bytes.Buffer{}
	zipWriter := zip.NewWriter(zipBuf)

	var allResults []itemResult
	successCount := 0

	for _, zone := range req.Zones {
		item := itemResult{Domain: zone.Domain}

		// 1. 生成私钥 + CSR
		privKeyPEM, csrPEM, err := generateDomainCSR(zone.Domain)
		if err != nil {
			item.Error = err.Error()
			allResults = append(allResults, item)
			continue
		}

		// 2. 调用 CF Origin CA API
		payload := map[string]interface{}{
			"hostnames":          []string{zone.Domain, "*." + zone.Domain},
			"requested_validity": req.ValidityDays,
			"request_type":       req.RequestType,
			"csr":                string(csrPEM),
		}

		apiURL := fmt.Sprintf("%s/certificates", cfAPIBase)
		data, err := dns.cfRequest("POST", apiURL, token, payload)
		if err != nil {
			item.Error = "调用 CF API 失败: " + err.Error()
			allResults = append(allResults, item)
			continue
		}

		var cfResp map[string]interface{}
		if err := json.Unmarshal(data, &cfResp); err != nil {
			item.Error = "解析 CF 响应失败: " + err.Error()
			allResults = append(allResults, item)
			continue
		}
		if ok2, _ := cfResp["success"].(bool); !ok2 {
			item.Error = cfOriginCAErrDetail(cfResp)
			allResults = append(allResults, item)
			continue
		}

		// 3. 提取证书 PEM
		cfResult, _ := cfResp["result"].(map[string]interface{})
		certPEM, _ := cfResult["certificate"].(string)
		if certPEM == "" {
			item.Error = "CF 返回证书为空"
			allResults = append(allResults, item)
			continue
		}

		// 4. 写入 zip（目录：域名/）
		dirPrefix := zone.Domain + "/"

		privW, err := zipWriter.Create(dirPrefix + "private.key")
		if err == nil {
			privW.Write(privKeyPEM)
		}
		certW, err := zipWriter.Create(dirPrefix + "cert.pem")
		if err == nil {
			certW.Write([]byte(certPEM))
		}

		item.Success = true
		successCount++
		allResults = append(allResults, item)
	}

	zipWriter.Close()

	// 全部失败时返回 JSON
	if successCount == 0 {
		response.Resp(c, map[string]interface{}{
			"total":   len(req.Zones),
			"success": 0,
			"fail":    len(req.Zones),
			"results": allResults,
		})
		return
	}

	// 下载 zip
	fileName := fmt.Sprintf("cf_certs_%s.zip", time.Now().Format("20060102_150405"))
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/zip")
	c.Data(http.StatusOK, "application/zip", zipBuf.Bytes())
}

// ListCerts 查询账号下已签发的 CF Origin 证书列表
func (s *CfCertService) ListCerts(c *gin.Context) {
	s.SetDbAlias("app")

	accountIdStr := c.Query("account_id")
	if accountIdStr == "" {
		response.Resp(c, "account_id 必填")
		return
	}

	accountId, err := strconv.ParseInt(accountIdStr, 10, 64)
	if err != nil || accountId == 0 {
		response.Resp(c, "account_id 格式错误")
		return
	}

	token, ok := s.getTokenByCertService(accountId)
	if !ok {
		response.Resp(c, "账号不存在或已禁用")
		return
	}

	dns := &CfDnsService{}
	apiURL := fmt.Sprintf("%s/certificates?per_page=100", cfAPIBase)
	data, err2 := dns.cfRequest("GET", apiURL, token, nil)
	if err2 != nil {
		response.Resp(c, err2.Error())
		return
	}

	var result map[string]interface{}
	json.Unmarshal(data, &result)
	response.Resp(c, result)
}
