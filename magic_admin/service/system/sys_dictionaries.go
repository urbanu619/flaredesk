package system

import (
	"fmt"
	"github.com/demdxx/gocast"
	"github.com/gin-gonic/gin"
	"go_server/model/common/response"
	system2 "go_server/model/system"
	"go_server/service/base"
)

type DictionaryService struct {
	base.SysCommonService
}

// 字典服务 增 删 改 查
// 所有字典信息

func (s *DictionaryService) DictionariesTree(c *gin.Context) {
	type Detail struct {
		DictionaryId int64  `json:"dictionaryId" gorm:"comment:字典ID"`
		Label        string `json:"label" gorm:"type:varchar(50);comment:展示值"`
		Value        string `json:"value" gorm:"type:varchar(50);comment:字典值"`
		Extend       string `json:"extend" gorm:"type:varchar(100);comment:扩展值"`
		Sort         int64  `json:"sort" gorm:"comment:排序"`
	}
	type DictionaryNode struct {
		Key     string    `json:"key"`
		Name    string    `json:"name"`
		Desc    string    `json:"desc"`
		Details []*Detail `json:"details"`
	}
	dictionaryMap := make(map[int64]*DictionaryNode)
	dictionaries, err := base.GetMore[system2.Dictionaries](s.DB(), "enable", true)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	for _, item := range dictionaries {
		dictionaryMap[item.ID] = &DictionaryNode{
			Key:     item.Key,
			Name:    item.Name,
			Desc:    item.Desc,
			Details: make([]*Detail, 0),
		}
	}
	dictionaryDetails, err := base.GetMore[system2.DictionaryDetail](s.DB().Order("sort"), "enable", true)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	for _, item := range dictionaryDetails {
		if _, ok := dictionaryMap[item.DictionaryId]; ok {
			dictionaryMap[item.DictionaryId].Details = append(dictionaryMap[item.DictionaryId].Details, &Detail{
				DictionaryId: item.DictionaryId,
				Label:        item.Label,
				Value:        item.Value,
				Extend:       item.Extend,
				Sort:         item.Sort,
			})
		}
	}
	res := make(map[string]interface{})
	res["dictionaryMap"] = dictionaryMap
	response.Resp(c, res)
	return
}

// --finish

func (s *DictionaryService) FindDictionaries(c *gin.Context) {
	type request[T any] struct {
		base.ListRequest[T]
		Id     interface{} `form:"id"`
		Name   string      `json:"name" gorm:"column:name;type:varchar(50);comment:名称"`
		Key    string      `json:"key" gorm:"column:key;;unique;type:varchar(50);comment:关键词"`
		Enable interface{} `json:"enable" gorm:"column:enable;type:tinyint(1);comment:是否有效"`
	}
	req := new(request[system2.Dictionaries])
	if err := c.BindQuery(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	db := s.DB()
	if req.Id != nil && gocast.ToInt64(req.Id) != 0 {
		db = db.Where("id", req.Id)
	}
	if req.Name != "" {
		db = db.Where("label LIKE ?", "%"+req.Name+"%")
	}
	if req.Key != "" {
		db = db.Where("`key` LIKE ?", "%"+req.Key+"%")
	}
	if req.Enable != nil {
		db = db.Where("enable = ?", gocast.ToBool(req.Enable))
	}
	resp, err := base.NewQueryBaseHandler(system2.NewDictionaries()).List(db, req)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c, resp)
}

func (s *DictionaryService) CreateDictionary(c *gin.Context) {
	userId := c.GetInt64("userId")
	// 限制用户增加必须管理员才可以操作
	if userId != system2.AdminId {
		response.Resp(c, "不允许操作")
		return
	}
	type request struct {
		Name string `json:"name" gorm:"column:name;type:varchar(50);comment:名称"`
		Key  string `json:"key" gorm:"column:key;;unique;type:varchar(50);comment:关键词"`
		Desc string `json:"desc" gorm:"column:desc;type:varchar(255);comment:说明"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	if base.CountRow[system2.Dictionaries](s.DB(), "key", req.Key) > 0 {
		response.Resp(c, "key 不可重复")
		return
	}
	row := &system2.Dictionaries{
		Name:   req.Name,
		Key:    req.Key,
		Enable: true,
		Desc:   req.Desc,
	}
	if err := s.DB().Create(&row).Error; err != nil {
		response.Resp(c, "添加失败")
		return
	}
	response.Resp(c)
	return
}

func (s *DictionaryService) UpdateDictionary(c *gin.Context) {
	userId := c.GetInt64("userId")
	// 限制用户增加必须管理员才可以操作
	if userId != system2.AdminId {
		response.Resp(c, "不允许操作")
		return
	}
	type request struct {
		Name   string      `json:"name" gorm:"column:name;type:varchar(50);comment:名称"`
		Key    string      `json:"key" gorm:"column:key;;unique;type:varchar(50);comment:关键词"`
		Enable interface{} `json:"enable"`
		Desc   string      `json:"desc" gorm:"column:desc;type:varchar(255);comment:说明"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	if req.Key == "" || req.Name == "" {
		response.Resp(c, "key, name 必填")
		return
	}
	row, ok := base.GetOne[system2.Dictionaries](s.DB(), "key", req.Key)
	if !ok {
		response.Resp(c, fmt.Sprintf("字典:%s不存在", req.Key))
		return
	}
	row.Name = req.Name
	row.Desc = req.Desc
	row.Enable = gocast.ToBool(req.Enable)
	if err := s.DB().Save(&row).Error; err != nil {
		response.Resp(c, fmt.Sprintf("修改失败:%s", err.Error()))
		return
	}
	response.Resp(c)
	return
}

func (s *DictionaryService) DelDictionary(c *gin.Context) {
	userId := c.GetInt64("userId")
	// 限制用户增加必须管理员才可以操作
	if userId != system2.AdminId {
		response.Resp(c, "不允许操作")
		return
	}
	type request struct {
		Key string `json:"key" gorm:"column:key;;unique;type:varchar(50);comment:关键词"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	row, ok := base.GetOne[system2.Dictionaries](s.DB(), "key", req.Key)
	if !ok {
		response.Resp(c, fmt.Sprintf("字典:%s不存在", req.Key))
		return
	}
	txDb := s.DB().Begin()
	if err := txDb.Delete(&row).Error; err != nil {
		response.Resp(c, "删除失败")
		txDb.Rollback()
		return
	}
	if err := txDb.
		Where("dictionary_id", row.ID).
		Delete(&system2.DictionaryDetail{}).Error; err != nil {
		response.Resp(c, "删除失败")
		txDb.Rollback()
		return
	}
	txDb.Commit()
	response.Resp(c)
	return
}

// --finish

func (s *DictionaryService) FindDictionaryDetail(c *gin.Context) {
	type request[T any] struct {
		base.ListRequest[T]
		Id     interface{} `form:"id"`
		Label  string      `json:"label" gorm:"type:varchar(50);comment:展示值"`
		Value  string      `json:"value" gorm:"type:varchar(50);comment:字典值"`
		Enable interface{} `json:"enable" gorm:"column:enable;type:tinyint(1);comment:是否有效"`
	}
	req := new(request[system2.DictionaryDetail])
	if err := c.BindQuery(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	db := s.DB()
	if req.Id != nil && gocast.ToInt64(req.Id) != 0 {
		db = db.Where("id", req.Id)
	}
	if req.Label != "" {
		db = db.Where("label LIKE ?", "%"+req.Label+"%")
	}
	if req.Value != "" {
		db = db.Where("value LIKE ?", "%"+req.Value+"%")
	}
	if req.Enable != nil {
		db = db.Where("enable = ?", gocast.ToBool(req.Enable))
	}
	resp, err := base.NewQueryBaseHandler(system2.NewDictionaryDetail()).List(db, req)
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c, resp)
}

func (s *DictionaryService) CreateDictionaryDetail(c *gin.Context) {
	userId := c.GetInt64("userId")
	// 限制用户增加必须管理员才可以操作
	if userId != system2.AdminId {
		response.Resp(c, "不允许操作")
		return
	}
	type request struct {
		DictionaryId interface{} `json:"dictionaryId" gorm:"comment:字典ID"`
		Label        string      `json:"label" gorm:"type:varchar(50);comment:展示值"`
		Value        string      `json:"value" gorm:"type:varchar(50);comment:字典值"`
		Extend       string      `json:"extend" gorm:"type:varchar(100);comment:扩展值"`
		Enable       interface{} `json:"enable" gorm:"column:enable;type:tinyint(1);comment:是否有效"`
		Sort         interface{} `json:"sort" gorm:"comment:排序"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	if req.DictionaryId == nil || req.Label == "" || req.Value == "" {
		response.Resp(c, "DictionaryId,Label,Value 必填")
		return
	}
	if base.CountRow[system2.Dictionaries](s.DB(), "id", req.DictionaryId) == 0 {
		response.Resp(c, "DictionaryId不存在")
		return
	}
	row := &system2.DictionaryDetail{
		DictionaryId: gocast.ToInt64(req.DictionaryId),
		Label:        req.Label,
		Value:        req.Value,
		Extend:       req.Extend,
		Enable:       gocast.ToBool(req.Enable),
		Sort:         gocast.ToInt64(req.Sort),
	}
	if err := s.DB().Create(&row).Error; err != nil {
		response.Resp(c, "添加失败")
		return
	}
	response.Resp(c)
	return
}

func (s *DictionaryService) UpdateDictionaryDetail(c *gin.Context) {
	userId := c.GetInt64("userId")
	// 限制用户增加必须管理员才可以操作
	if err := s.CheckIsAdmin(userId); err != nil {
		response.Resp(c, err.Error())
		return
	}
	type request struct {
		Id     interface{} `json:"id"`
		Label  string      `json:"label" gorm:"type:varchar(50);comment:展示值"`
		Value  string      `json:"value" gorm:"type:varchar(50);comment:字典值"`
		Extend string      `json:"extend" gorm:"type:varchar(100);comment:扩展值"`
		Enable interface{} `json:"enable" gorm:"column:enable;type:tinyint(1);comment:是否有效"`
		Sort   interface{} `json:"sort" gorm:"comment:排序"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	if req.Id == nil || req.Id == 0 || req.Label == "" || req.Value == "" {
		response.Resp(c, "Id、label、value必填")
		return
	}
	row, ok := base.GetOne[system2.DictionaryDetail](s.DB(), "id", req.Id)
	if !ok {
		response.Resp(c, fmt.Sprintf("字典值:%s不存在", req.Id))
		return
	}
	row.Label = req.Label
	row.Value = req.Value
	row.Extend = req.Extend
	row.Sort = gocast.ToInt64(req.Sort)
	row.Enable = gocast.ToBool(req.Enable)
	if err := s.DB().Save(&row).Error; err != nil {
		response.Resp(c, fmt.Sprintf("修改失败:%s", err.Error()))
		return
	}
	response.Resp(c)
	return
}

func (s *DictionaryService) DelDictionaryDetail(c *gin.Context) {
	userId := c.GetInt64("userId")
	if err := s.CheckIsAdmin(userId); err != nil {
		response.Resp(c, err.Error())
		return
	}
	type request struct {
		Id interface{} `json:"id"`
	}
	req := new(request)
	if err := c.BindJSON(req); err != nil {
		response.Resp(c, err.Error())
		return
	}
	if req.Id == nil || gocast.ToInt64(req.Id) == 0 {
		response.Resp(c, "Id必填")
		return
	}
	row, ok := base.GetOne[system2.DictionaryDetail](s.DB(), "id", req.Id)
	if !ok {
		response.Resp(c, fmt.Sprintf("字典值:%s不存在", req.Id))
		return
	}
	if err := s.DB().Delete(&row).Error; err != nil {
		response.Resp(c, fmt.Sprintf("修改失败:%s", err.Error()))
		return
	}
	response.Resp(c)
}
