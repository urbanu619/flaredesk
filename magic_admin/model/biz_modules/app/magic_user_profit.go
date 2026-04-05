package app

// 引入关联包
import (
	"github.com/shopspring/decimal"
)

type MagicUserProfit struct {
	CreatedAt                     int64           `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt                     int64           `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	UserId                        int64           `json:"userId" gorm:"column:user_id;type:bigint;comment:用户ID;primarykey;NOT NULL"`
	Uid                           string          `json:"uid" gorm:"column:uid;type:varchar(64);comment:交易所ID;index;unique:uni_magic_user_profit_uid"`
	ProfitBalance                 decimal.Decimal `json:"profitBalance" gorm:"column:profit_balance;type:decimal(25,8);comment:收益余额"`
	CumClaimLevel                 int64           `json:"cumClaimLevel" gorm:"column:cum_claim_level;type:bigint;comment:已领取等级"`
	CumClaimLevelUpgradeQuantity  decimal.Decimal `json:"cumClaimLevelUpgradeQuantity" gorm:"column:cum_claim_level_upgrade_quantity;type:decimal(25,8);comment:累计领取等级升级数量"`
	CumClaimLevelUpgradeUsdAmount decimal.Decimal `json:"cumClaimLevelUpgradeUsdAmount" gorm:"column:cum_claim_level_upgrade_usd_amount;type:decimal(25,8);comment:累计领取等级升级USD价值"`
	TodayNodeProfitQuantity       decimal.Decimal `json:"todayNodeProfitQuantity" gorm:"column:today_node_profit_quantity;type:decimal(25,8);comment:当日节点奖励产出系统币"`
	NodeProfitQuantity            decimal.Decimal `json:"nodeProfitQuantity" gorm:"column:node_profit_quantity;type:decimal(25,8);comment:累计节点产出系统币"`
	NodeProfitUsdAmount           decimal.Decimal `json:"nodeProfitUsdAmount" gorm:"column:node_profit_usd_amount;type:decimal(25,8);comment:累计节点产出USD价值"`
	PersonPledgeRankingProfit     decimal.Decimal `json:"personPledgeRankingProfit" gorm:"column:person_pledge_ranking_profit;type:decimal(25,8);comment:个人质押排行榜收益"`
	PersonPledgeRankingUsdProfit  decimal.Decimal `json:"personPledgeRankingUsdProfit" gorm:"column:person_pledge_ranking_usd_profit;type:decimal(25,8);comment:个人质押排行榜收益U"`
	FewAcRankingProfit            decimal.Decimal `json:"fewAcRankingProfit" gorm:"column:few_ac_ranking_profit;type:decimal(25,8);comment:小区业绩排行榜收益"`
	FewAcRankingUsdProfit         decimal.Decimal `json:"fewAcRankingUsdProfit" gorm:"column:few_ac_ranking_usd_profit;type:decimal(25,8);comment:小区业绩排行榜收益U"`
	TodayDynamicProfitQuantity    decimal.Decimal `json:"todayDynamicProfitQuantity" gorm:"column:today_dynamic_profit_quantity;type:decimal(25,8);comment:当日动态奖励产出系统币"`
	DynamicProfitQuantity         decimal.Decimal `json:"dynamicProfitQuantity" gorm:"column:dynamic_profit_quantity;type:decimal(25,8);comment:累计动态奖励系统币"`
	DynamicProfitUsdAmount        decimal.Decimal `json:"dynamicProfitUsdAmount" gorm:"column:dynamic_profit_usd_amount;type:decimal(25,8);comment:累计动态奖励USD价值"`
	StaticProfitQuantity          decimal.Decimal `json:"staticProfitQuantity" gorm:"column:static_profit_quantity;type:decimal(25,8);comment:累计静态产出系统币"`
	StaticProfitUsdAmount         decimal.Decimal `json:"staticProfitUsdAmount" gorm:"column:static_profit_usd_amount;type:decimal(25,8);comment:累计静态产出USD价值"`
	TotalProfitQuantity           decimal.Decimal `json:"totalProfitQuantity" gorm:"column:total_profit_quantity;type:decimal(25,8);comment:总收益币"`
	TotalProfitUsd                decimal.Decimal `json:"totalProfitUsd" gorm:"column:total_profit_usd;type:decimal(25,8);comment:总收益USD价值"`
	Version                       int64           `json:"version" gorm:"column:version;type:bigint;comment:事务版本"`
}

func (*MagicUserProfit) TableName() string {
	return "magic_user_profit"
}

func NewMagicUserProfit() *MagicUserProfit {
	return &MagicUserProfit{}
}
