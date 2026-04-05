package app

// 引入关联包

type MagicUserLogs struct { 
	Id int64 `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt int64 `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt int64 `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	UserId int64 `json:"userId" gorm:"column:user_id;type:bigint;comment:用户Id"`
	UserAgent string `json:"userAgent" gorm:"column:user_agent;type:varchar(1024);comment:客户端信息"`
	Ip string `json:"ip" gorm:"column:ip;type:varchar(64);comment:ip"`
	Action string `json:"action" gorm:"column:action;type:varchar(20);comment:请求方法"`
	Query string `json:"query" gorm:"column:query;type:text;comment:query"`
	Path string `json:"path" gorm:"column:path;type:varchar(200);comment:path"`
	Status int64 `json:"status" gorm:"column:status;type:bigint;comment:状态"`
	Request string `json:"request" gorm:"column:request;type:text;comment:请求内容消息"`
	Response string `json:"response" gorm:"column:response;type:text;comment:请求返回信息"`
	Elapsed string `json:"elapsed" gorm:"column:elapsed;type:varchar(50);comment:耗时"`
}

func (*MagicUserLogs) TableName() string {
	return "magic_user_logs"
}

func NewMagicUserLogs() *MagicUserLogs {
	return &MagicUserLogs{}
}
