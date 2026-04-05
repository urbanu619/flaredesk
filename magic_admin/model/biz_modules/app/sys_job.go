package app

// 引入关联包
import (
	"github.com/shopspring/decimal"
)

type SysJob struct { 
	Id int64 `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt int64 `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt int64 `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	JobDate string `json:"jobDate" gorm:"column:job_date;type:varchar(64);comment:任务日期;index;unique:uni_sys_job_job_date"`
	State string `json:"state" gorm:"column:state;type:varchar(256);comment:快照任务状态: waiting-等待执行 running-进行中  finish-完成 fail-失败"`
	TodayNetPledged decimal.Decimal `json:"todayNetPledged" gorm:"column:today_net_pledged;type:decimal(25,8);comment:当日平台总质押"`
	TodayStaticProfit decimal.Decimal `json:"todayStaticProfit" gorm:"column:today_static_profit;type:decimal(25,8);comment:当日平台静态总收益"`
	TodayStaticUsdProfit decimal.Decimal `json:"todayStaticUsdProfit" gorm:"column:today_static_usd_profit;type:decimal(25,8);comment:当日平台静态总USD收益"`
	TodayFeeProfit decimal.Decimal `json:"todayFeeProfit" gorm:"column:today_fee_profit;type:decimal(25,8);comment:当日平台手续费总收益"`
	IsSendNodeProfit bool `json:"isSendNodeProfit" gorm:"column:is_send_node_profit;type:tinyint(1);comment:是否发放节点奖励"`
	IsSendAcProfit bool `json:"isSendAcProfit" gorm:"column:is_send_ac_profit;type:tinyint(1);comment:是否发放小区业绩排行奖励"`
	IsSendStakeProfit bool `json:"isSendStakeProfit" gorm:"column:is_send_stake_profit;type:tinyint(1);comment:是否发放个人质押排行奖励"`
	Desc string `json:"desc" gorm:"column:desc;type:text;comment:错误信息"`
}

func (*SysJob) TableName() string {
	return "sys_job"
}

func NewSysJob() *SysJob {
	return &SysJob{}
}
