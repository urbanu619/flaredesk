package app

// 引入关联包
import (
	"github.com/shopspring/decimal"
)

type MagicAsset struct { 
	Id int64 `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt int64 `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt int64 `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	UserId int64 `json:"userId" gorm:"column:user_id;type:bigint;comment:用户ID;index;unique:idx_address_coin;NOT NULL"`
	Uid string `json:"uid" gorm:"column:uid;type:varchar(64);comment:交易所ID:;index"`
	Symbol string `json:"symbol" gorm:"column:symbol;type:varchar(20);comment:币种名;index;unique:idx_address_coin"`
	Balance decimal.Decimal `json:"balance" gorm:"column:balance;type:decimal(25,8);comment:净资产"`
	Frozen decimal.Decimal `json:"frozen" gorm:"column:frozen;type:decimal(25,8);comment:冻结资产"`
	TotalAmount decimal.Decimal `json:"totalAmount" gorm:"column:total_amount;type:decimal(25,8);comment:总资产"`
	Version int64 `json:"version" gorm:"column:version;type:bigint;comment:事务版本"`
}

func (*MagicAsset) TableName() string {
	return "magic_asset"
}

func NewMagicAsset() *MagicAsset {
	return &MagicAsset{}
}
