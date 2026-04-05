package app

// 引入关联包
import (
	"github.com/shopspring/decimal"
)

type NodeOrderPayments struct { 
	Id int64 `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt int64 `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt int64 `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	PaymentRecordId string `json:"paymentRecordId" gorm:"column:payment_record_id;type:varchar(64);comment:支付订单ID;index;unique:uni_node_order_payments_payment_record_id"`
	NodeOrderRecordId string `json:"nodeOrderRecordId" gorm:"column:node_order_record_id;type:varchar(64);comment:节点购买记录ID;index"`
	UserId int64 `json:"userId" gorm:"column:user_id;type:bigint;comment:用户ID;index"`
	Uid string `json:"uid" gorm:"column:uid;type:varchar(64);comment:用户交易所ID:;index"`
	OrderUsdAmount decimal.Decimal `json:"orderUsdAmount" gorm:"column:order_usd_amount;type:decimal(25,8);comment:订单总价值"`
	PaySymbol string `json:"paySymbol" gorm:"column:pay_symbol;type:varchar(64);comment:支付币种"`
	PaySymbolUsdPrice decimal.Decimal `json:"paySymbolUsdPrice" gorm:"column:pay_symbol_usd_price;type:decimal(25,8);comment:支付币种单价"`
	PayQuantity decimal.Decimal `json:"payQuantity" gorm:"column:pay_quantity;type:decimal(25,8);comment:支付数量"`
	PayUsdAmount decimal.Decimal `json:"payUsdAmount" gorm:"column:pay_usd_amount;type:decimal(25,8);comment:支付总价值"`
}

func (*NodeOrderPayments) TableName() string {
	return "node_order_payments"
}

func NewNodeOrderPayments() *NodeOrderPayments {
	return &NodeOrderPayments{}
}
