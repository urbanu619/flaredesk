package app

// 引入关联包

type NodeBanner struct { 
	Id int64 `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt int64 `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt int64 `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	Name string `json:"name" gorm:"column:name;type:varchar(256);comment: banner名称"`
	Url string `json:"url" gorm:"column:url;type:text;comment:bannerURL"`
	IsDisplay int8 `json:"isDisplay" gorm:"column:is_display;type:tinyint;comment:是否有效 1:开放展示 0:不开放展示"`
	Sort int64 `json:"sort" gorm:"column:sort;type:bigint;comment:排序"`
	Remark string `json:"remark" gorm:"column:remark;type:varchar(200);comment:备注信息"`
}

func (*NodeBanner) TableName() string {
	return "node_banner"
}

func NewNodeBanner() *NodeBanner {
	return &NodeBanner{}
}
