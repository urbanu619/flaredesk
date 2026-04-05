package app

// 引入关联包

type SysI18n struct { 
	Id int64 `json:"id" gorm:"column:id;type:bigint;comment:;primarykey;NOT NULL"`
	Lang string `json:"lang" gorm:"column:lang;type:varchar(45);comment:'语言, en, cn, hk'"`
	Name string `json:"name" gorm:"column:name;type:varchar(256);comment:'名称'"`
	Text string `json:"text" gorm:"column:text;type:text;comment:'文本'"`
}

func (*SysI18n) TableName() string {
	return "sys_i18n"
}

func NewSysI18n() *SysI18n {
	return &SysI18n{}
}
