package app

// CfDnsTemplate DNS 记录模板
type CfDnsTemplate struct {
	Id        int64  `json:"id" gorm:"column:id;type:bigint;primarykey;NOT NULL"`
	CreatedAt int64  `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt int64  `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	Name      string `json:"name" gorm:"column:name;type:varchar(100);comment:模板名称"`
	Remark    string `json:"remark" gorm:"column:remark;type:varchar(500);comment:备注"`
	Records   string `json:"records" gorm:"column:records;type:text;comment:DNS记录JSON"`
}

func (*CfDnsTemplate) TableName() string {
	return "cf_dns_template"
}

func NewCfDnsTemplate() *CfDnsTemplate {
	return &CfDnsTemplate{}
}

func (*CfDnsTemplate) Comment() string {
	return "DNS记录模板表"
}
