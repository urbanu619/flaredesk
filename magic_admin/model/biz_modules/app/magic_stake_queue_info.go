package app

// 引入关联包
import (
	"github.com/shopspring/decimal"
)

type MagicStakeQueueInfo struct { 
	Id int64 `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt int64 `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt int64 `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	RewardDate string `json:"rewardDate" gorm:"column:reward_date;type:varchar(64);comment:奖励日期;index"`
	SuccessOrders int64 `json:"successOrders" gorm:"column:success_orders;type:bigint;comment:排队成功订单数量"`
	UpperUsdAmount decimal.Decimal `json:"upperUsdAmount" gorm:"column:upper_usd_amount;type:decimal(25,8);comment:每日可通过USD数量"`
	PassedUsdAmount decimal.Decimal `json:"passedUsdAmount" gorm:"column:passed_usd_amount;type:decimal(25,8);comment:已通过USD数量"`
}

func (*MagicStakeQueueInfo) TableName() string {
	return "magic_stake_queue_info"
}

func NewMagicStakeQueueInfo() *MagicStakeQueueInfo {
	return &MagicStakeQueueInfo{}
}
