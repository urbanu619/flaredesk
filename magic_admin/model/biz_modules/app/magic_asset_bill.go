package app

// 引入关联包
import (
	"github.com/shopspring/decimal"
)

type MagicAssetBill struct { 
	Id int64 `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt int64 `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt int64 `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	AssetId int64 `json:"assetId" gorm:"column:asset_id;type:bigint;comment:资产表ID"`
	UserId int64 `json:"userId" gorm:"column:user_id;type:bigint;comment:用户ID;index;NOT NULL"`
	Uid string `json:"uid" gorm:"column:uid;type:varchar(64);comment:用户UID;index;NOT NULL"`
	Symbol string `json:"symbol" gorm:"column:symbol;type:varchar(20);comment:币种名"`
	BusinessNumber int `json:"businessNumber" gorm:"column:business_number;type:int;comment:业务场景;index;NOT NULL"`
	BusinessName string `json:"businessName" gorm:"column:business_name;type:varchar(50);comment:业务场名称;index"`
	BeforeAmount decimal.Decimal `json:"beforeAmount" gorm:"column:before_amount;type:decimal(25,8);comment:交易前可用余额;NOT NULL"`
	Balance decimal.Decimal `json:"balance" gorm:"column:balance;type:decimal(25,8);comment:净资产"`
	Frozen decimal.Decimal `json:"frozen" gorm:"column:frozen;type:decimal(25,8);comment:冻结资产"`
	Amount decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(25,8);comment:发生可用金额;NOT NULL"`
	AfterAmount decimal.Decimal `json:"afterAmount" gorm:"column:after_amount;type:decimal(25,8);comment:交易后可用余额;NOT NULL"`
	ContextName string `json:"contextName" gorm:"column:context_name;type:varchar(128);comment:上下文名"`
	ContextValue string `json:"contextValue" gorm:"column:context_value;type:varchar(128);comment:上下文值"`
	Describe string `json:"describe" gorm:"column:describe;type:text;comment:流水备注信息"`
}

func (*MagicAssetBill) TableName() string {
	return "magic_asset_bill"
}

func NewMagicAssetBill() *MagicAssetBill {
	return &MagicAssetBill{}
}
