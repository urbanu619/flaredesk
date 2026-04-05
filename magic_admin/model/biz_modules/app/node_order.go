package app

// 引入关联包
import (
	"github.com/shopspring/decimal"
)

type NodeOrder struct { 
	Id int64 `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt int64 `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt int64 `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	RecordId string `json:"recordId" gorm:"column:record_id;type:varchar(64);comment:奖励记录ID;index;unique:uni_node_order_record_id"`
	UserId int64 `json:"userId" gorm:"column:user_id;type:bigint;comment:用户ID;index"`
	Uid string `json:"uid" gorm:"column:uid;type:varchar(64);comment:用户交易所ID:;index"`
	ProductId int64 `json:"productId" gorm:"column:product_id;type:bigint;comment:产品ID"`
	ProductName string `json:"productName" gorm:"column:product_name;type:varchar(64);comment:产品名称"`
	UsdPrice decimal.Decimal `json:"usdPrice" gorm:"column:usd_price;type:decimal(25,8);comment:单价"`
	Quantity decimal.Decimal `json:"quantity" gorm:"column:quantity;type:decimal(25,2);comment:数量"`
	UsdAmount decimal.Decimal `json:"usdAmount" gorm:"column:usd_amount;type:decimal(25,8);comment:总价值"`
	PayState int64 `json:"payState" gorm:"column:pay_state;type:bigint;comment:支付状态:0:待支付 1: 支付成功 2: 支付失败 "`
	State int64 `json:"state" gorm:"column:state;type:bigint;comment:订单状态:0:未知 1:有效 2:过期失效 3:其他失效 "`
	Source int64 `json:"source" gorm:"column:source;type:bigint;comment:来源:0 用户购买 1: 2: 3: 4:后台购买"`
}

func (*NodeOrder) TableName() string {
	return "node_order"
}

func NewNodeOrder() *NodeOrder {
	return &NodeOrder{}
}
