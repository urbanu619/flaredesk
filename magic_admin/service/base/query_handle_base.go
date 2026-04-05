package base

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go_server/model/common/response"
	"gorm.io/gorm"
)

// gin 封装模型curd处理
// 使用范型 较为鸡肋 不建议使用

type ModelBaseHandler[T any] struct {
}

func NewBaseHandler[T any](t *T) *ModelBaseHandler[T] {
	return &ModelBaseHandler[T]{}
}

func (*ModelBaseHandler[T]) DeleteOne(c *gin.Context, db *gorm.DB) {
	req := new(BaseRequest)
	if err := c.BindQuery(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	if req.Id == 0 {
		response.Resp(c, "id is zero")
		return
	}
	tempInfo := new(T)
	// 先通过Id获取记录'
	if err := db.First(&tempInfo, req.Id).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}
	exCmd := db.Model(&tempInfo).Delete(&tempInfo)
	if exCmd.RowsAffected != 1 {
		response.Resp(c, "delete fail")
		return
	}
	if exCmd.Error != nil {
		response.Resp(c, exCmd.Error.Error())
		return
	}
	response.Resp(c)
	return
}

func (*ModelBaseHandler[T]) UpdateOne(c *gin.Context, db *gorm.DB) {
	req := new(BaseRequest)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	if req.Id == 0 {
		response.Resp(c, "id is zero")
		return
	}
	tempInfo := new(T)
	// 先通过Id获取记录'
	if err := db.First(&tempInfo, req.Id).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}
	dbEx := db.Model(&tempInfo).Where("id", req.Id).
		Updates(req.Data)
	if dbEx.Error != nil {
		response.Resp(c, dbEx.Error.Error())
		return
	}
	if dbEx.RowsAffected == 0 {
		response.Resp(c, "update fail")
		return
	}
	if err := db.First(&tempInfo, req.Id).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c, map[string]any{
		"info":         tempInfo,
		"rowsAffected": dbEx.RowsAffected,
	})
	return
}

func (*ModelBaseHandler[T]) Get(c *gin.Context, db *gorm.DB) {
	req := new(BaseRequest)
	if err := c.BindQuery(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	if req.Id == 0 {
		response.Resp(c, "id is zero")
		return
	}
	tempInfo := new(T)
	// 先通过Id获取记录'
	if err := db.First(&tempInfo, req.Id).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c, map[string]any{
		"info": tempInfo,
	})
	return
}

func (*ModelBaseHandler[T]) Create(c *gin.Context, db *gorm.DB) {
	req := new(BaseRequest)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	info := new(T)
	body, _ := json.Marshal(req.Data)
	if err := json.Unmarshal(body, info); err != nil {
		response.Resp(c, err.Error())
		return
	}
	if err := db.Create(info).Error; err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c, map[string]any{
		"info": info,
	})
}

func (*ModelBaseHandler[T]) List(c *gin.Context, db *gorm.DB) {
	req := NewListRequest[T]()
	if err := c.BindQuery(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	resp, err := req.Query(db)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c, resp)
	return
}
