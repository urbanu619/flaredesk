package app

// 引入关联包
import (
	"github.com/shopspring/decimal"
)

type SysLevelConfig struct { 
	Level int64 `json:"level" gorm:"column:level;type:bigint;comment:用户等级;index;unique:uni_sys_level_config_level"`
	Icon string `json:"icon" gorm:"column:icon;type:varchar(255);comment:等级图标"`
	LevelName string `json:"levelName" gorm:"column:level_name;type:varchar(256);comment:等级名称"`
	PersonAchievement decimal.Decimal `json:"personAchievement" gorm:"column:person_achievement;type:decimal(25,8);comment:个人业绩要求"`
	TeamAchievement decimal.Decimal `json:"teamAchievement" gorm:"column:team_achievement;type:decimal(25,8);comment:团队业绩要求"`
	FewTeamAchievement decimal.Decimal `json:"fewTeamAchievement" gorm:"column:few_team_achievement;type:decimal(25,8);comment:小团队业绩要求"`
	StaticRatio decimal.Decimal `json:"staticRatio" gorm:"column:static_ratio;type:decimal(25,8);comment:总静态收益占比"`
	AvgRatio decimal.Decimal `json:"avgRatio" gorm:"column:avg_ratio;type:decimal(25,8);comment:均分占比"`
	WeightedRatio decimal.Decimal `json:"weightedRatio" gorm:"column:weighted_ratio;type:decimal(25,8);comment:加权占比"`
	GiftLargeRegionAchievement decimal.Decimal `json:"giftLargeRegionAchievement" gorm:"column:gift_large_region_achievement;type:decimal(25,8);comment:赠送大区业绩"`
	LevelUpgradeUsdProfit decimal.Decimal `json:"levelUpgradeUsdProfit" gorm:"column:level_upgrade_usd_profit;type:decimal(25,8);comment:升级奖励U"`
}

func (*SysLevelConfig) TableName() string {
	return "sys_level_config"
}

func NewSysLevelConfig() *SysLevelConfig {
	return &SysLevelConfig{}
}
