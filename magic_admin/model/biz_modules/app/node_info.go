package app

// 引入关联包
import (
	"github.com/shopspring/decimal"
)

type NodeInfo struct { 
	Id int64 `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt int64 `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt int64 `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	NodeName string `json:"nodeName" gorm:"column:node_name;type:varchar(36);comment:节点名称;index;unique:uni_node_info_node_name"`
	SupportPaymentSymbols string `json:"supportPaymentSymbols" gorm:"column:support_payment_symbols;type:varchar(64);comment:支持的支付币种,分割"`
	UsdPrice decimal.Decimal `json:"usdPrice" gorm:"column:usd_price;type:decimal(25,8);comment:单价"`
	SoldQuantity decimal.Decimal `json:"soldQuantity" gorm:"column:sold_quantity;type:decimal(25,2);comment:已售出数量"`
	UpperQuantityLimit decimal.Decimal `json:"upperQuantityLimit" gorm:"column:upper_quantity_limit;type:decimal(25,2);comment:可售出上限"`
	LimitTime int64 `json:"limitTime" gorm:"column:limit_time;type:bigint;comment:结束时间"`
	Sort int64 `json:"sort" gorm:"column:sort;type:bigint;comment:排序"`
	MagicProfitRatio decimal.Decimal `json:"magicProfitRatio" gorm:"column:magic_profit_ratio;type:decimal(25,2);comment:平台静态产出分红比例"`
	PersonBuyLimit int64 `json:"personBuyLimit" gorm:"column:person_buy_limit;type:bigint;comment:产品用户购买有效上限"`
	RebateRatio decimal.Decimal `json:"rebateRatio" gorm:"column:rebate_ratio;type:decimal(25,2);comment:最近一层返佣比例"`
	DividendRatio decimal.Decimal `json:"dividendRatio" gorm:"column:dividend_ratio;type:decimal(25,8);comment:分红比例"`
	IsDisplay int8 `json:"isDisplay" gorm:"column:is_display;type:tinyint;comment:是否有效 1:开放展示 0:不开放展示"`
	Enable int8 `json:"enable" gorm:"column:enable;type:tinyint;comment:是否有效 1:开放购买 0:未开放购买"`
}

func (*NodeInfo) TableName() string {
	return "node_info"
}

func NewNodeInfo() *NodeInfo {
	return &NodeInfo{}
}
