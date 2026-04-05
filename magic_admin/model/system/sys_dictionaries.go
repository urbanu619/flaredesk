package system

import (
	"go_server/model/common"
)

type Dictionaries struct {
	common.GormFullModel
	Name   string `json:"name" gorm:"column:name;type:varchar(50);comment:名称"`
	Key    string `json:"key" gorm:"column:key;;unique;type:varchar(50);comment:关键词"`
	Enable bool   `json:"enable" gorm:"column:enable;type:tinyint(1);comment:是否有效"`
	Desc   string `json:"desc" gorm:"column:desc;type:varchar(255);comment:说明"`
}

func (*Dictionaries) TableName() string {
	return common.ModelPrefix + "dictionaries"
}

func NewDictionaries() *Dictionaries {
	return &Dictionaries{}
}

func (*Dictionaries) Comment() string {
	return "系统字典表"
}

type DictionaryDetail struct {
	common.GormFullModel
	DictionaryId int64  `json:"dictionaryId" gorm:"comment:字典ID"`
	Label        string `json:"label" gorm:"type:varchar(50);comment:展示值"`
	Value        string `json:"value" gorm:"type:varchar(50);comment:字典值"`
	Extend       string `json:"extend" gorm:"type:varchar(100);comment:扩展值"`
	Enable       bool   `json:"enable" gorm:"column:enable;type:tinyint(1);comment:是否有效"`
	Sort         int64  `json:"sort" gorm:"comment:排序"`
}

func (*DictionaryDetail) TableName() string {
	return common.ModelPrefix + "dictionary_detail"
}

func NewDictionaryDetail() *DictionaryDetail {
	return &DictionaryDetail{}
}

func (*DictionaryDetail) Comment() string {
	return "系统字典详情"
}
