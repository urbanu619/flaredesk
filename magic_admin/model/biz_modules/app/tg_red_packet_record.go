package app

import (
	"github.com/shopspring/decimal"
	"time"
)

// TgRedPacketRecord 抢红包记录
type TgRedPacketRecord struct {
	Id               int64           `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt        int64           `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt        int64           `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	PacketId         int64           `json:"packetId" gorm:"column:packet_id;type:bigint;index;comment:红包ID"`
	PacketNo         string          `json:"packetNo" gorm:"column:packet_no;type:varchar(64);index;comment:红包编号"`
	UserId           int64           `json:"userId" gorm:"column:user_id;type:bigint;index;comment:用户ID"`
	TelegramId       int64           `json:"telegramId" gorm:"column:telegram_id;type:bigint;index;comment:Telegram用户ID"`
	TelegramUsername string          `json:"telegramUsername" gorm:"column:telegram_username;type:varchar(255);comment:Telegram用户名"`
	Amount           decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(25,8);comment:抢到金额"`
	IsBest           int             `json:"isBest" gorm:"column:is_best;type:tinyint;default:0;comment:是否手气最佳"`
	Sequence         int64           `json:"sequence" gorm:"column:sequence;type:bigint;comment:抢红包序号"`
	GrabbedAt        time.Time       `json:"grabbedAt" gorm:"column:grabbed_at;type:timestamp;comment:抢红包时间"`
}

func (*TgRedPacketRecord) TableName() string {
	return "magic_tg_red_packet_record"
}

func NewTgRedPacketRecord() *TgRedPacketRecord {
	return &TgRedPacketRecord{}
}

func (*TgRedPacketRecord) Comment() string {
	return "TG红包领取记录表"
}
