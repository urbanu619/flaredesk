package app

// 引入关联包

type MagicUserActionLog struct { 
	Id int64 `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt int64 `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt int64 `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	UserId int64 `json:"userId" gorm:"column:user_id;type:bigint;comment:用户ID"`
	Uid string `json:"uid" gorm:"column:uid;type:varchar(64);comment:用户UID"`
	Source string `json:"source" gorm:"column:source;type:text;comment:原信息"`
	Target string `json:"target" gorm:"column:target;type:text;comment:目标信息"`
	Action string `json:"action" gorm:"column:action;type:varchar(42);comment:变更类型:修改地址 重置上级 修改上级"`
}

func (*MagicUserActionLog) TableName() string {
	return "magic_user_action_log"
}

func NewMagicUserActionLog() *MagicUserActionLog {
	return &MagicUserActionLog{}
}
