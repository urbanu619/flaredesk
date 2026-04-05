package app

// CfZone 本地缓存的 Cloudflare Zone 信息
type CfZone struct {
	Id            int64  `json:"id" gorm:"column:id;type:bigint;primarykey;NOT NULL"`
	CreatedAt     int64  `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt     int64  `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	AccountId     int64  `json:"accountId" gorm:"column:account_id;type:bigint;index;comment:CF账号ID"`
	ZoneId        string `json:"zoneId" gorm:"column:zone_id;type:varchar(64);uniqueIndex;comment:CF Zone ID"`
	Name          string `json:"name" gorm:"column:name;type:varchar(255);index;comment:域名"`
	Status        string `json:"status" gorm:"column:status;type:varchar(32);comment:状态"`
	Paused        bool   `json:"paused" gorm:"column:paused;comment:是否暂停"`
	PlanName      string `json:"planName" gorm:"column:plan_name;type:varchar(64);comment:套餐名称"`
	NameServers   string `json:"nameServers" gorm:"column:name_servers;type:varchar(500);comment:NS服务器JSON"`
	ActivatedOn   string `json:"activatedOn" gorm:"column:activated_on;type:varchar(64);comment:激活时间"`
	Remark        string `json:"remark" gorm:"column:remark;type:varchar(500);comment:备注"`
}

func (*CfZone) TableName() string {
	return "cf_zone"
}

func NewCfZone() *CfZone {
	return &CfZone{}
}

func (*CfZone) Comment() string {
	return "Cloudflare Zone缓存表"
}
