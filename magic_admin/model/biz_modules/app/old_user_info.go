package app

// 引入关联包
import (
	"github.com/shopspring/decimal"
)

type OldUserInfo struct {
	Id                            int64           `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt                     int64           `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt                     int64           `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	Uid                           string          `json:"uid" gorm:"column:uid;type:varchar(64);comment:交易所ID;index;unique:idx_user_uid_merchant_id"`
	ParentUid                     string          `json:"parentUid" gorm:"column:parent_uid;type:varchar(64);comment:上级交易所ID;index"`
	TdBalance                     decimal.Decimal `json:"tdBalance" gorm:"column:td_balance;type:decimal(25,8);comment:td余额"`
	UsdtBalance                   decimal.Decimal `json:"usdtBalance" gorm:"column:usdt_balance;type:decimal(25,8);comment:usdt余额"`
	ParentNodes                   decimal.Decimal `json:"parentNodes" gorm:"column:parent_nodes;type:decimal(25,8);comment:节点数量??"`
	ParentAmount                  decimal.Decimal `json:"parentAmount" gorm:"column:parent_amount;type:decimal(25,8);comment:节点总价值"`
	CumClaimLevel                 int64           `json:"cumClaimLevel" gorm:"column:cum_claim_level;type:bigint;comment:已领取等级"`
	CumClaimLevelUpgradeQuantity  decimal.Decimal `json:"cumClaimLevelUpgradeQuantity" gorm:"column:cum_claim_level_upgrade_quantity;type:decimal(25,8);comment:累计领取等级升级数量"`
	CumClaimLevelUpgradeUsdAmount decimal.Decimal `json:"cumClaimLevelUpgradeUsdAmount" gorm:"column:cum_claim_level_upgrade_usd_amount;type:decimal(25,8);comment:累计领取等级升级USD价值"`
	StakeQuantity                 decimal.Decimal `json:"stakeQuantity" gorm:"column:stake_quantity;type:decimal(25,8);comment:当前质押数量"`
	StakeUsdAmount                decimal.Decimal `json:"stakeUsdAmount" gorm:"column:stake_usd_amount;type:decimal(25,8);comment:当前质押USD价值"`
	StakeCollectionBalance        decimal.Decimal `json:"stakeCollectionBalance" gorm:"column:stake_collection_balance;type:decimal(25,8);comment:可领取td静态奖励数量"`
	ImportStakeQuantity           decimal.Decimal `json:"importStakeQuantity" gorm:"column:import_stake_quantity;type:decimal(25,8);comment:当前质押数量"`
	ImportStakeUsdAmount          decimal.Decimal `json:"importStakeUsdAmount" gorm:"column:import_stake_usd_amount;type:decimal(25,8);comment:当前质押USD价值"`
	ImportStakeCollectionBalance  decimal.Decimal `json:"importStakeCollectionBalance" gorm:"column:import_stake_collection_balance;type:decimal(25,8);comment:可领取td静态奖励数量"`
	ProfitBalance                 decimal.Decimal `json:"profitBalance" gorm:"column:profit_balance;type:decimal(25,8);comment:动静态可提取收益数量"`
	IsHandler                     bool            `json:"isHandler" gorm:"column:is_handler;type:tinyint(1);comment:当前数据是否处理完成"`
	Enable                        bool            `json:"enable" gorm:"column:enable;type:tinyint(1);comment:账户是否有效 1=有效 0=无效"`
	IsRobot                       bool            `json:"isRobot" gorm:"column:is_robot;type:tinyint(1);comment:"`
	WithoutStakeRanking           bool            `json:"withoutStakeRanking" gorm:"column:without_stake_ranking;type:tinyint(1);comment:0 参与 1不参与个人排行"`
	WithoutFewRegionRanking       bool            `json:"withoutFewRegionRanking" gorm:"column:without_few_region_ranking;type:tinyint(1);comment:0 参与 1不参与小区业绩排行"`
	LockStakeProfit               bool            `json:"lockStakeProfit" gorm:"column:lock_stake_profit;type:tinyint(1);comment:0 正常 1锁定不产出收益"`
	IsNotStakeLimit               bool            `json:"isNotStakeLimit" gorm:"column:is_not_stake_limit;type:tinyint(1);comment:是否无质押限制"`
	IsNotCountStakeAc             bool            `json:"isNotCountStakeAc" gorm:"column:is_not_count_stake_ac;type:tinyint(1);comment:是否不计算全网业绩统计"`
	IsZero                        bool            `json:"isZero" gorm:"column:is_zero;type:tinyint(1);comment:是否零号线"`
	LockWithdraw                  bool            `json:"lockWithdraw" gorm:"column:lock_withdraw;type:tinyint(1);comment:锁定提现  0=未锁定  1=锁定"`
	SubsLockWithdraw              bool            `json:"subsLockWithdraw" gorm:"column:subs_lock_withdraw;type:tinyint(1);comment:锁定伞下提现  0=未锁定  1=锁定"`
	IsOld                         bool            `json:"isOld" gorm:"column:is_old;type:tinyint(1);comment:是否老规则用户"`
	IsExemptionLargeAchievement   int64           `json:"isExemptionLargeAchievement" gorm:"column:is_exemption_large_achievement;type:bigint;comment:是否免除大区业绩"`
	GiftPersonAchievement         decimal.Decimal `json:"giftPersonAchievement" gorm:"column:gift_person_achievement;type:decimal(25,8);comment:赠送质押业绩"`
	InflatedFewTeamAchievement    decimal.Decimal `json:"inflatedFewTeamAchievement" gorm:"column:inflated_few_team_achievement;type:decimal(25,8);comment:虚增小区业绩"`
	LockRedeem                    int8            `json:"lockRedeem" gorm:"column:lock_redeem;type:tinyint;comment:赎回锁定 0 正常 1 锁定"`
	UnCollection                  int8            `json:"unCollection" gorm:"column:un_collection;type:tinyint;comment:0 正常 1 自动领取静态收益"`
	PerformanceRobot              int8            `json:"performanceRobot" gorm:"column:performance_robot;type:tinyint;comment:0 正常 1 自动领取静态收益"`
	LevelRobot                    int8            `json:"levelRobot" gorm:"column:level_robot;type:tinyint;comment:0 正常 1 自动领取静态收益"`
	InternalAccount               int8            `json:"internalAccount" gorm:"column:internal_account;type:tinyint;comment:0 正常 1 自动领取静态收益"`
	FireRobot                     int8            `json:"fireRobot" gorm:"column:fire_robot;type:tinyint;comment:0 正常 1 自动领取静态收益"`
	IsHandlerOrder                bool            `json:"isHandlerOrder" gorm:"column:is_handler_order;type:tinyint(1);comment:订单数据是否处理完成"`
	IsHandlerExtraOrder           bool            `json:"isHandlerExtraOrder" gorm:"column:is_handler_extra_order;type:tinyint(1);comment:extra订单数据是否处理完成"`
}

func (*OldUserInfo) TableName() string {
	return "old_user_info"
}

func NewOldUserInfo() *OldUserInfo {
	return &OldUserInfo{}
}
