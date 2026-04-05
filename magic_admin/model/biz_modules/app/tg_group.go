package app

// TgGroup Telegram群组信息
type TgGroup struct {
	Id          int64  `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt   int64  `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt   int64  `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	ChatId      int64  `json:"chatId" gorm:"column:chat_id;type:bigint;uniqueIndex;comment:Telegram群组ID;NOT NULL"`
	ChatType    string `json:"chatType" gorm:"column:chat_type;type:varchar(20);comment:群组类型(group/supergroup/channel)"`
	Title       string `json:"title" gorm:"column:title;type:varchar(255);comment:群组标题"`
	Username    string `json:"username" gorm:"column:username;type:varchar(100);comment:群组用户名"`
	Description string `json:"description" gorm:"column:description;type:text;comment:群组描述"`
	MemberCount int    `json:"memberCount" gorm:"column:member_count;type:int;comment:成员数量"`
	Status      int8   `json:"status" gorm:"column:status;type:tinyint;default:1;comment:状态: 1=正常, 2=已禁用, 3=已退出"`
	BotJoinedAt int64  `json:"botJoinedAt" gorm:"column:bot_joined_at;type:bigint;comment:机器人加入时间"`
	Remark      string `json:"remark" gorm:"column:remark;type:varchar(500);comment:备注"`
}

func (*TgGroup) TableName() string {
	return "magic_tg_group"
}

func NewTgGroup() *TgGroup {
	return &TgGroup{}
}

func (*TgGroup) Comment() string {
	return "Telegram群组信息表"
}
