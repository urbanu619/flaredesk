package app

// 引入关联包
import (
	"github.com/shopspring/decimal"
)

type MagicStakeProduct struct { 
	Id int64 `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt int64 `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt int64 `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	Name string `json:"name" gorm:"column:name;type:varchar(256);comment: 产品名称"`
	Symbol string `json:"symbol" gorm:"column:symbol;type:varchar(64);comment:支付币种"`
	CycleDay int64 `json:"cycleDay" gorm:"column:cycle_day;type:bigint;comment:质押周期天0:活期"`
	ProfitCycleHours int64 `json:"profitCycleHours" gorm:"column:profit_cycle_hours;type:bigint;comment:收益周期"`
	PledgeMode string `json:"pledgeMode" gorm:"column:pledge_mode;type:varchar(256);comment:质押类型:current活期 regularly 定期"`
	PeriodProfitRatio decimal.Decimal `json:"periodProfitRatio" gorm:"column:period_profit_ratio;type:decimal(25,8);comment:期收益比例"`
	IsQueue bool `json:"isQueue" gorm:"column:is_queue;type:tinyint(1);comment:产品是否开启排队"`
	QuotaPerMinutes string `json:"quotaPerMinutes" gorm:"column:quota_per_minutes;type:varchar(256);comment:每分钟释放生效额度"`
	CurrentAvailableQuota string `json:"currentAvailableQuota" gorm:"column:current_available_quota;type:varchar(256);comment:当前累计可用生效额度"`
	CollectionFeeRatio decimal.Decimal `json:"collectionFeeRatio" gorm:"column:collection_fee_ratio;type:decimal(25,8);comment:领取奖励手续费比例"`
	ClaimFeeRatio decimal.Decimal `json:"claimFeeRatio" gorm:"column:claim_fee_ratio;type:decimal(25,8);comment:收益提取手续费比例"`
	RedeemFeeRatio decimal.Decimal `json:"redeemFeeRatio" gorm:"column:redeem_fee_ratio;type:decimal(25,8);comment:赎回手续费比例"`
	MinQuantity decimal.Decimal `json:"minQuantity" gorm:"column:min_quantity;type:decimal(25,8);comment:最小质押数量"`
	OverflowPeriods int64 `json:"overflowPeriods" gorm:"column:overflow_periods;type:bigint;comment:奖励溢出期数"`
	Sort int64 `json:"sort" gorm:"column:sort;type:bigint;comment:排序"`
	Enable int8 `json:"enable" gorm:"column:enable;type:tinyint;comment:是否有效 1=可质押 0=暂停质押"`
}

func (*MagicStakeProduct) TableName() string {
	return "magic_stake_product"
}

func NewMagicStakeProduct() *MagicStakeProduct {
	return &MagicStakeProduct{}
}
