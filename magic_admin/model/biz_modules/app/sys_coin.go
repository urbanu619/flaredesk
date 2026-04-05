package app

// 引入关联包
import (
	"github.com/shopspring/decimal"
)

type SysCoin struct {
	Id                   int64           `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt            int64           `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt            int64           `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	Symbol               string          `json:"symbol" gorm:"column:symbol;type:varchar(20);comment:币种名;index;unique:uni_sys_coin_symbol"`
	Icon                 string          `json:"icon" gorm:"column:icon;type:varchar(255);comment:图标"`
	WithdrawStatus       int8            `json:"withdrawStatus" gorm:"column:withdraw_status;type:tinyint;comment:是否允许提现(全局):0=禁止，1=允许"`
	WithdrawNoauditLimit decimal.Decimal `json:"withdrawNoauditLimit" gorm:"column:withdraw_noaudit_limit;type:decimal(25,8);comment:提现免审额度，0表示都需要审核"`
	MinWithdrawAmount    decimal.Decimal `json:"minWithdrawAmount" gorm:"column:min_withdraw_amount;type:decimal(25,8);comment:提现最小限额，0表示无限制"`
	MaxWithdrawAmount    decimal.Decimal `json:"maxWithdrawAmount" gorm:"column:max_withdraw_amount;type:decimal(25,8);comment:每笔提现最大限额，-1表示无限制，0表示不能提现"`
	WithdrawRatio        decimal.Decimal `json:"withdrawRatio" gorm:"column:withdraw_ratio;type:decimal(25,8);comment:提现手续费比例"`
	WithdrawBaseAmount   decimal.Decimal `json:"withdrawBaseAmount" gorm:"column:withdraw_base_amount;type:decimal(25,8);comment:单次提现手续费数量"`
	Enable               bool            `json:"enable" gorm:"column:enable;type:tinyint(1);comment:是否有效"`
	WithdrawDayMax       decimal.Decimal `json:"withdrawDayMax" gorm:"column:withdraw_day_max;type:decimal(25,8);comment:单日可提现最大数量"`
}

func (*SysCoin) TableName() string {
	return "sys_coin"
}

func NewSysCoin() *SysCoin {
	return &SysCoin{}
}
