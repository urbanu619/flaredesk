package app

// 引入关联包
import (
	"github.com/shopspring/decimal"
)

type MagicAssetRwRecord struct {
	Id            int64           `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt     int64           `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt     int64           `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	OrderId       string          `json:"orderId" gorm:"column:order_id;type:varchar(50);comment:订单ID-唯一识别;index;unique:uni_magic_asset_rw_record_order_id"`
	UserId        int64           `json:"userId" gorm:"column:user_id;type:bigint;comment:用户ID;index;NOT NULL"`
	Uid           string          `json:"uid" gorm:"column:uid;type:varchar(64);comment:交易所ID;index"`
	OpenId        string          `json:"openId" gorm:"column:open_id;type:varchar(128);comment:uuid"`
	Symbol        string          `json:"symbol" gorm:"column:symbol;type:varchar(20);comment:币种名"`
	Amount        decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(25,8);comment:发生可用金额;NOT NULL"`
	Ratio         decimal.Decimal `json:"ratio" gorm:"column:ratio;type:decimal(25,8);comment:手续费比例;NOT NULL"`
	BaseFee       decimal.Decimal `json:"baseFee" gorm:"column:base_fee;type:decimal(25,8);comment:单笔基础手续费;NOT NULL"`
	Fee           decimal.Decimal `json:"fee" gorm:"column:fee;type:decimal(25,8);comment:手续费数量 金额*比例+单笔基础手续费;NOT NULL"`
	RealAmount    decimal.Decimal `json:"realAmount" gorm:"column:real_amount;type:decimal(25,8);comment:实际操作数量;NOT NULL"`
	Status        int64           `json:"status" gorm:"column:status;type:bigint;comment:状态:0 待审核 1 审核通过 2 完成转出/转入 >=3 审核不通过/失败"`
	WithdrawState int64           `json:"withdrawState" gorm:"column:withdraw_state;type:bigint;comment:提现状态:0-待提现 1-已提现 此字段不允许后台修改"`
	Direction     int64           `json:"direction" gorm:"column:direction;type:bigint;comment:资金方向:1 交易所转入-充值 2 交易所转出-提现 "`
	DirectionName string          `json:"directionName" gorm:"column:direction_name;type:varchar(128);comment:类型名称"`
	Describe      string          `json:"describe" gorm:"column:describe;type:varchar(255);comment:流水备注信息"`
}

func (*MagicAssetRwRecord) TableName() string {
	return "magic_asset_rw_record"
}

func NewMagicAssetRwRecord() *MagicAssetRwRecord {
	return &MagicAssetRwRecord{}
}
