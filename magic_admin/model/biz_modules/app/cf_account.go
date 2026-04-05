package app

// CfAccount Cloudflare 账号
type CfAccount struct {
	Id        int64  `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt int64  `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt int64  `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	Name      string `json:"name" gorm:"column:name;type:varchar(100);comment:账号名称"`
	Email     string `json:"email" gorm:"column:email;type:varchar(255);comment:CF账号邮箱"`
	ApiToken  string `json:"apiToken" gorm:"column:api_token;type:varchar(500);comment:API Token"`
	Status    int8   `json:"status" gorm:"column:status;type:tinyint;default:1;comment:状态: 1=正常, 2=禁用"`
	Remark    string `json:"remark" gorm:"column:remark;type:varchar(500);comment:备注"`
}

func (*CfAccount) TableName() string {
	return "cf_account"
}

func NewCfAccount() *CfAccount {
	return &CfAccount{}
}

func (*CfAccount) Comment() string {
	return "Cloudflare账号表"
}
