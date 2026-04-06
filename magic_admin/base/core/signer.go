package core

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"go_server/model/system"
)

// GenerateSysInfo 生成系统内部签名用的随机地址和私钥（纯随机 hex）
func GenerateSysInfo() (addr string, priKey string, err error) {
	addrBytes := make([]byte, 20)
	if _, err = rand.Read(addrBytes); err != nil {
		return "", "", fmt.Errorf("生成地址失败: %w", err)
	}

	priBytes := make([]byte, 32)
	if _, err = rand.Read(priBytes); err != nil {
		return "", "", fmt.Errorf("生成私钥失败: %w", err)
	}

	addr = "0x" + hex.EncodeToString(addrBytes)
	priKey = hex.EncodeToString(priBytes)
	return addr, priKey, nil
}

// BuildSignMessage 生成内部代理请求签名（HMAC-SHA256）
// 格式：address:timestamp:hmac_sig
func BuildSignMessage() (string, error) {
	db := MainDb()
	cfg := system.NewSysSignConfig()
	if err := db.Where("sign_name = ? AND is_system_sign = ?", SignSystemAdminServerName, true).
		First(cfg).Error; err != nil {
		return "", fmt.Errorf("读取签名配置失败: %w", err)
	}

	ts := strconv.FormatInt(time.Now().Unix(), 10)
	msg := cfg.SignAddress + ":" + ts

	priKeyBytes, err := hex.DecodeString(cfg.SignPriKey)
	if err != nil {
		// 私钥格式不是 hex（兼容旧数据），直接用原始字符串作为 key
		priKeyBytes = []byte(cfg.SignPriKey)
	}

	mac := hmac.New(sha256.New, priKeyBytes)
	mac.Write([]byte(msg))
	sig := hex.EncodeToString(mac.Sum(nil))

	return msg + ":" + sig, nil
}
