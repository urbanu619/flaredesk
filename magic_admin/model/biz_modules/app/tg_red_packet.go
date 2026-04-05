package app

import "time"

// TgRedPacket 红包主记录
type TgRedPacket struct {
	Id            int64     `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt     int64     `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt     int64     `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	PacketNo      string    `json:"packetNo" gorm:"column:packet_no;type:varchar(64);uniqueIndex;comment:红包编号;NOT NULL"`
	UserId        int64     `json:"userId" gorm:"column:user_id;type:bigint;index;comment:发送者用户ID;NOT NULL"`
	GroupId       string    `json:"groupId" gorm:"column:group_id;type:varchar(128);index;comment:群组ID;NOT NULL"`
	PacketType    int8      `json:"packetType" gorm:"column:packet_type;type:tinyint;comment:红包类型:1=普通红包,2=手气红包;NOT NULL"`
	TotalAmount   float64   `json:"totalAmount" gorm:"column:total_amount;type:decimal(25,8);comment:红包总金额;NOT NULL"`
	TotalCount    int64     `json:"totalCount" gorm:"column:total_count;type:bigint;comment:红包总个数;NOT NULL"`
	GrabbedCount  int64     `json:"grabbedCount" gorm:"column:grabbed_count;type:bigint;default:0;comment:已抢个数"`
	GrabbedAmount float64   `json:"grabbedAmount" gorm:"column:grabbed_amount;type:decimal(25,8);default:0;comment:已抢金额"`
	RemainCount   int64     `json:"remainCount" gorm:"column:remain_count;type:bigint;comment:剩余个数;NOT NULL"`
	RemainAmount  float64   `json:"remainAmount" gorm:"column:remain_amount;type:decimal(25,8);comment:剩余金额;NOT NULL"`
	Symbol        string    `json:"symbol" gorm:"column:symbol;type:varchar(20);default:VND;comment:币种"`
	MessageId     int64     `json:"messageId" gorm:"column:message_id;type:bigint;comment:Telegram消息ID"`
	BlessingWords string    `json:"blessingWords" gorm:"column:blessing_words;type:varchar(255);comment:祝福语"`
	Status        int8      `json:"status" gorm:"column:status;type:tinyint;default:1;index;comment:状态:1=进行中,2=已抢完,3=已过期"`
	ExpireAt      time.Time `json:"expireAt" gorm:"column:expire_at;type:timestamp;comment:过期时间"`
	CompletedAt   time.Time `json:"completedAt" gorm:"column:completed_at;type:timestamp;comment:完成时间"`
}

func (*TgRedPacket) TableName() string {
	return "magic_tg_red_packet"
}

func NewTgRedPacket() *TgRedPacket {
	return &TgRedPacket{}
}

func (*TgRedPacket) Comment() string {
	return "TG红包记录表"
}
