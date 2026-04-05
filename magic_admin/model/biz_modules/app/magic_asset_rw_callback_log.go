package app

// 引入关联包

type MagicAssetRwCallbackLog struct { 
	Id int64 `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt int64 `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt int64 `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	Currency string `json:"currency" gorm:"column:currency;type:varchar(20);comment:币种缩写"`
	Amount string `json:"amount" gorm:"column:amount;type:varchar(50);comment:金额"`
	OrderId string `json:"orderId" gorm:"column:order_id;type:varchar(50);comment:订单ID-唯一;index;unique:uni_magic_asset_rw_callback_log_order_id"`
	OpenId string `json:"openId" gorm:"column:open_id;type:varchar(50);comment:openId"`
	CreateTime string `json:"createTime" gorm:"column:create_time;type:varchar(20);comment:时间戳"`
	Type string `json:"type" gorm:"column:type;type:varchar(20);comment:类型 user_to_application:充值 application_to_user:提现"`
	AppId string `json:"appId" gorm:"column:app_id;type:varchar(50);comment:应用ID"`
	CustomParam string `json:"customParam" gorm:"column:custom_param;type:varchar(100);comment:透传信息"`
	Sign string `json:"sign" gorm:"column:sign;type:varchar(100);comment:签名信息"`
	State int8 `json:"state" gorm:"column:state;type:tinyint;comment:处理状态 0-待处理 1-成功 2-失败"`
	ErrDesc string `json:"errDesc" gorm:"column:err_desc;type:varchar(255);comment:错误信息"`
	Status int64 `json:"status" gorm:"column:status;type:bigint;comment:充值、提现回调状态码 0 / 1 成功 2 失败"`
}

func (*MagicAssetRwCallbackLog) TableName() string {
	return "magic_asset_rw_callback_log"
}

func NewMagicAssetRwCallbackLog() *MagicAssetRwCallbackLog {
	return &MagicAssetRwCallbackLog{}
}
