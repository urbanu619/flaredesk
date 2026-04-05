package app

// 引入关联包
import (
	"github.com/shopspring/decimal"
)

type MagicStakeUserCurrentOrder struct { 
	Id int64 `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt int64 `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt int64 `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	ProductId int64 `json:"productId" gorm:"column:product_id;type:bigint;comment:产品ID"`
	ProductName string `json:"productName" gorm:"column:product_name;type:varchar(64);comment:产品名称"`
	UserId int64 `json:"userId" gorm:"column:user_id;type:bigint;comment:用户ID;index"`
	Uid string `json:"uid" gorm:"column:uid;type:varchar(64);comment:交易所ID;index;unique:idx_user_uid_merchant_id"`
	Symbol string `json:"symbol" gorm:"column:symbol;type:varchar(32);comment:质押币种"`
	CumQuantity decimal.Decimal `json:"cumQuantity" gorm:"column:cum_quantity;type:decimal(25,8);comment:累计质押数量"`
	Quantity decimal.Decimal `json:"quantity" gorm:"column:quantity;type:decimal(25,8);comment:当前质押数量"`
	CumUsdAmount decimal.Decimal `json:"cumUsdAmount" gorm:"column:cum_usd_amount;type:decimal(25,8);comment:累计质押USD价值"`
	UsdAmount decimal.Decimal `json:"usdAmount" gorm:"column:usd_amount;type:decimal(25,8);comment:当前质押USD价值"`
	FirstPledgeDate string `json:"firstPledgeDate" gorm:"column:first_pledge_date;type:varchar(20);comment:首次质押日"`
	RedeemedQuantity decimal.Decimal `json:"redeemedQuantity" gorm:"column:redeemed_quantity;type:decimal(25,8);comment:累计释放股票数量"`
	RedeemedUsdAmount decimal.Decimal `json:"redeemedUsdAmount" gorm:"column:redeemed_usd_amount;type:decimal(25,8);comment:累计释放USD价值"`
	LastProfitPeriod string `json:"lastProfitPeriod" gorm:"column:last_profit_period;type:varchar(32);comment:最近产出期号"`
	AvailableQuantity decimal.Decimal `json:"availableQuantity" gorm:"column:available_quantity;type:decimal(25,8);comment:待领取数量"`
	AvailableUsdAmount decimal.Decimal `json:"availableUsdAmount" gorm:"column:available_usd_amount;type:decimal(25,8);comment:待领取价值"`
	CumClaimQuantity decimal.Decimal `json:"cumClaimQuantity" gorm:"column:cum_claim_quantity;type:decimal(25,8);comment:累计领取数量"`
	CumClaimUsdAmount decimal.Decimal `json:"cumClaimUsdAmount" gorm:"column:cum_claim_usd_amount;type:decimal(25,8);comment:累计领取价值"`
	LastClaimPeriod string `json:"lastClaimPeriod" gorm:"column:last_claim_period;type:varchar(32);comment:最近领取期号"`
	Version int64 `json:"version" gorm:"column:version;type:bigint;comment:事务版本"`
	ErrDesc string `json:"errDesc" gorm:"column:err_desc;type:text;comment:最近一次错误信息备注"`
}

func (*MagicStakeUserCurrentOrder) TableName() string {
	return "magic_stake_user_current_order"
}

func NewMagicStakeUserCurrentOrder() *MagicStakeUserCurrentOrder {
	return &MagicStakeUserCurrentOrder{}
}
