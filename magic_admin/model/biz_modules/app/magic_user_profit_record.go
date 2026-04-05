package app

// 引入关联包
import (
	"github.com/shopspring/decimal"
)

type MagicUserProfitRecord struct { 
	Id int64 `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt int64 `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt int64 `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	RawUserId string `json:"rawUserId" gorm:"column:raw_user_id;type:varchar(64);comment:源用户UID;index"`
	RawUserUid string `json:"rawUserUid" gorm:"column:raw_user_uid;type:varchar(64);comment:源用户UID;index"`
	RewardUserId string `json:"rewardUserId" gorm:"column:reward_user_id;type:varchar(64);comment:收益用户ID;index"`
	RewardUserUid string `json:"rewardUserUid" gorm:"column:reward_user_uid;type:varchar(64);comment:收益用户UID;index"`
	RewardSymbol string `json:"rewardSymbol" gorm:"column:reward_symbol;type:varchar(64);comment:奖励币种;index"`
	SymbolUsdPrice decimal.Decimal `json:"symbolUsdPrice" gorm:"column:symbol_usd_price;type:decimal(25,18);comment:奖励币种单价"`
	Indicator decimal.Decimal `json:"indicator" gorm:"column:indicator;type:decimal(25,8);comment:奖励指标值: 等级 排名 ..."`
	CapitalUsd decimal.Decimal `json:"capitalUsd" gorm:"column:capital_usd;type:decimal(25,8);comment:奖励基数价值"`
	RewardRatio decimal.Decimal `json:"rewardRatio" gorm:"column:reward_ratio;type:decimal(25,18);comment:奖励比例"`
	RewardUsdValue decimal.Decimal `json:"rewardUsdValue" gorm:"column:reward_usd_value;type:decimal(25,8);comment:奖励USD价值"`
	RewardQuantity decimal.Decimal `json:"rewardQuantity" gorm:"column:reward_quantity;type:decimal(25,8);comment:奖励币种数量"`
	BusinessNumber int `json:"businessNumber" gorm:"column:business_number;type:int;comment:业务场景;index;NOT NULL"`
	BusinessName string `json:"businessName" gorm:"column:business_name;type:varchar(64);comment:业务场名称"`
	ContextName string `json:"contextName" gorm:"column:context_name;type:varchar(128);comment:上下文名"`
	ContextValue string `json:"contextValue" gorm:"column:context_value;type:varchar(128);comment:上下文值"`
	RewardPeriod string `json:"rewardPeriod" gorm:"column:reward_period;type:varchar(64);comment:奖励期数;index"`
	RewardDate string `json:"rewardDate" gorm:"column:reward_date;type:varchar(64);comment:奖励日期;index"`
	State string `json:"state" gorm:"column:state;type:varchar(32);comment:状态waiting 待发放 success发放成功"`
	Describe string `json:"describe" gorm:"column:describe;type:varchar(1024);comment:收益详情"`
}

func (*MagicUserProfitRecord) TableName() string {
	return "magic_user_profit_record"
}

func NewMagicUserProfitRecord() *MagicUserProfitRecord {
	return &MagicUserProfitRecord{}
}
