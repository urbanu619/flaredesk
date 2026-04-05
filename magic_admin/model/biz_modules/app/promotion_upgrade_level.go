package app

// 引入关联包
import (
	"github.com/shopspring/decimal"
)

type PromotionUpgradeLevel struct {
	Id                int64           `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt         int64           `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt         int64           `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	UserId            int64           `json:"userId" gorm:"column:user_id;type:bigint;comment:用户ID;index"`
	Uid               string          `json:"uid" gorm:"column:uid;type:varchar(64);comment:交易所ID;index;unique:idx_user_uid_merchant_id"`
	RealAcLevel       int8            `json:"realAcLevel" gorm:"column:real_ac_level;type:tinyint;comment:奖励业绩等级;index"`
	Quantity          decimal.Decimal `json:"quantity" gorm:"column:quantity;type:decimal(25,8);comment:本次领取奖励数量"`
	Price             decimal.Decimal `json:"price" gorm:"column:price;type:decimal(25,8);comment:单价"`
	UsdAmount         decimal.Decimal `json:"usdAmount" gorm:"column:usd_amount;type:decimal(25,8);comment:本次领取奖励USD价值"`
	CumClaimQuantity  decimal.Decimal `json:"cumClaimQuantity" gorm:"column:cum_claim_quantity;type:decimal(25,8);comment:累计领取数量"`
	CumClaimUsdAmount decimal.Decimal `json:"cumClaimUsdAmount" gorm:"column:cum_claim_usd_amount;type:decimal(25,8);comment:累计领取USD价值"`
}

func (*PromotionUpgradeLevel) TableName() string {
	return "promotion_upgrade_level"
}

func NewPromotionUpgradeLevel() *PromotionUpgradeLevel {
	return &PromotionUpgradeLevel{}
}
