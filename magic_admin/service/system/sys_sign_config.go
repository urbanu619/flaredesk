package system

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_server/model/common/response"
	"go_server/model/system"
	"go_server/service/base"
)

type SignService struct {
	base.SysCommonService
}

func (s *SignService) Find(c *gin.Context) {
	db := s.DB()
	items, err := base.GetMore[system.SysSignConfig](db)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	type SignInfo struct {
		ID           int64  `json:"id" gorm:"primarykey;comment:id"`
		IsSystemSign bool   `json:"isSystemSign" gorm:"comment:'是否本系统密钥信息'"`
		SignName     string `json:"signName" gorm:"type:varchar(45);unique;comment:'签名系统名称-前缀:FOMO-PRO'"`
		SignAddress  string `json:"signAddress" gorm:"type:varchar(42);comment:'系统地址'"`
		SignExpSec   int64  `json:"signExpSec" gorm:"comment:'超时时间S'"`
		SysUrl       string `json:"sysUrl" gorm:"type:varchar(512);comment:外部系统请求链接"`
	}
	list := make([]*SignInfo, 0)
	for _, item := range items {
		list = append(list, &SignInfo{
			ID:           item.ID,
			IsSystemSign: item.IsSystemSign,
			SignName:     item.SignName,
			SignAddress:  item.SignAddress,
			SignExpSec:   item.SignExpSec,
			SysUrl:       item.SysUrl,
		})
	}
	response.Resp(c, map[string]interface{}{
		"list": list,
	})
}

func (s *SignService) Set(c *gin.Context) {
	// 限制用户增加必须管理员才可以操作
	userId := c.GetInt64("userId")
	if userId != system.AdminId {
		response.Resp(c, "不允许操作")
		return
	}
	type request struct {
		Id          interface{} `json:"id"`
		SignName    string      `json:"signName" gorm:"type:varchar(45);unique;comment:'签名系统名称-前缀:FOMO-PRO'"`
		SignAddress string      `json:"signAddress" gorm:"type:varchar(42);comment:'系统地址'"`
		SignExpSec  int64       `json:"signExpSec" gorm:"comment:'超时时间S'"`
		SysUrl      string      `json:"sysUrl" gorm:"type:varchar(512);comment:外部系统请求链接"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	row, ok := base.GetOne[system.SysSignConfig](s.DB(), "id", req.Id)
	if !ok {
		response.Resp(c, fmt.Sprintf("配置ID:%d不存在", req.Id))
		return
	}
	if row.IsSystemSign {
		response.Resp(c, fmt.Sprintf("本系统地址不允许修改"))
		return
	}
	row.SignAddress = req.SignAddress
	row.SignExpSec = req.SignExpSec
	row.SysUrl = req.SysUrl
	if err := s.DB().Save(&row).Error; err != nil {
		response.Resp(c, fmt.Sprintf("修改失败:%s", err.Error()))
		return
	}
	response.Resp(c, nil)
}
