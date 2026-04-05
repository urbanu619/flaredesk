package app

// 引入关联包
import (
	"github.com/shopspring/decimal"
)

type SysStakePeriodJob struct { 
	Id int64 `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt int64 `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt int64 `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	PeriodNo string `json:"periodNo" gorm:"column:period_no;type:varchar(64);comment:期编号;index;unique:idx_sys_stake_period_job_period_no"`
	JobState string `json:"jobState" gorm:"column:job_state;type:varchar(32);comment:状态:running 进行中 success 成功  failure 失败"`
	Price decimal.Decimal `json:"price" gorm:"column:price;type:decimal(25,8);comment:当期单价"`
	StaticSymbolProfit decimal.Decimal `json:"staticSymbolProfit" gorm:"column:static_symbol_profit;type:decimal(25,8);comment:期静态币总收益"`
	StaticUsdProfit decimal.Decimal `json:"staticUsdProfit" gorm:"column:static_usd_profit;type:decimal(25,8);comment:期USD总收益"`
	LevelProfit decimal.Decimal `json:"levelProfit" gorm:"column:level_profit;type:decimal(25,8);comment:动态等级奖励"`
	LevelProfitIsCal bool `json:"levelProfitIsCal" gorm:"column:level_profit_is_cal;type:tinyint(1);comment:动态等级奖励是否计算完成"`
	Desc string `json:"desc" gorm:"column:desc;type:text;comment:错误信息"`
}

func (*SysStakePeriodJob) TableName() string {
	return "sys_stake_period_job"
}

func NewSysStakePeriodJob() *SysStakePeriodJob {
	return &SysStakePeriodJob{}
}
