package app

// 引入关联包
import (
	"github.com/shopspring/decimal"
)

type MagicStakeUserCurrentOpsRecord struct { 
	Id int64 `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt int64 `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt int64 `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	RecordId string `json:"recordId" gorm:"column:record_id;type:varchar(64);comment:记录ID;index"`
	UserId int64 `json:"userId" gorm:"column:user_id;type:bigint;comment:用户ID;index"`
	Uid string `json:"uid" gorm:"column:uid;type:varchar(64);comment:交易所ID;index"`
	Price decimal.Decimal `json:"price" gorm:"column:price;type:decimal(25,8);comment:单价(U)"`
	BusinessNumber int `json:"businessNumber" gorm:"column:business_number;type:int;comment:业务场景;index;NOT NULL"`
	BusinessName string `json:"businessName" gorm:"column:business_name;type:varchar(64);comment:业务场名称"`
	ContextName string `json:"contextName" gorm:"column:context_name;type:varchar(128);comment:上下文名"`
	ContextValue string `json:"contextValue" gorm:"column:context_value;type:varchar(128);comment:上下文值"`
	BeforeQuantity decimal.Decimal `json:"beforeQuantity" gorm:"column:before_quantity;type:decimal(25,8);comment:操作前质押数量"`
	BeforeUsdAmount decimal.Decimal `json:"beforeUsdAmount" gorm:"column:before_usd_amount;type:decimal(25,8);comment:操作前USD价值"`
	Quantity decimal.Decimal `json:"quantity" gorm:"column:quantity;type:decimal(25,8);comment:质押数量"`
	UsdValue decimal.Decimal `json:"usdValue" gorm:"column:usd_value;type:decimal(25,8);comment:U价值"`
	AfterUsdAmount decimal.Decimal `json:"afterUsdAmount" gorm:"column:after_usd_amount;type:decimal(25,8);comment:操作后USD价值"`
	AfterQuantity decimal.Decimal `json:"afterQuantity" gorm:"column:after_quantity;type:decimal(25,8);comment:操作后质押数量"`
	QueueState string `json:"queueState" gorm:"column:queue_state;type:varchar(32);comment:订单状态:waiting 排队中 success 排队成功  cancel 已取消"`
	AdminId int64 `json:"adminId" gorm:"column:admin_id;type:bigint;comment:来源:0 用户 大于0表示管理后台"`
	EffTime int64 `json:"effTime" gorm:"column:eff_time;type:bigint;comment:生效时间戳"`
	Remark string `json:"remark" gorm:"column:remark;type:text;comment:备注"`
}

func (*MagicStakeUserCurrentOpsRecord) TableName() string {
	return "magic_stake_user_current_ops_record"
}

func NewMagicStakeUserCurrentOpsRecord() *MagicStakeUserCurrentOpsRecord {
	return &MagicStakeUserCurrentOpsRecord{}
}
