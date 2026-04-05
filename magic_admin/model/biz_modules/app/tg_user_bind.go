package app

import "time"

// TgUserBind Telegram用户绑定信息
type TgUserBind struct {
	Id                int64     `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt         int64     `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt         int64     `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	UserId            int64     `json:"userId" gorm:"column:user_id;type:bigint;index;comment:平台用户ID;NOT NULL"`
	TelegramId        int64     `json:"telegramId" gorm:"column:telegram_id;type:bigint;uniqueIndex;comment:Telegram用户ID;NOT NULL"`
	TelegramUsername  string    `json:"telegramUsername" gorm:"column:telegram_username;type:varchar(255);comment:Telegram用户名"`
	TelegramFirstName string    `json:"telegramFirstName" gorm:"column:telegram_first_name;type:varchar(255);comment:Telegram名字"`
	BindStatus        int       `json:"bindStatus" gorm:"column:bind_status;type:tinyint;default:1;comment:绑定状态: 1=已绑定, 0=已解绑"`
	BindTime          time.Time `json:"bindTime" gorm:"column:bind_time;type:timestamp;default:CURRENT_TIMESTAMP;comment:绑定时间"`
}

func (*TgUserBind) TableName() string {
	return "magic_tg_user_bind"
}

func (*TgUserBind) Comment() string {
	return "Telegram用户绑定表"
}

func NewTgUserBind() *TgUserBind {
	return &TgUserBind{}
}
