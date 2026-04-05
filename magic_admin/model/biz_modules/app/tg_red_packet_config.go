package app

import (
	"github.com/shopspring/decimal"
	"time"
)

type TgRedPacketConfig struct {
	Id            int64           `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt     int64           `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt     int64           `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	ConfigName    string          `json:"configName" gorm:"column:config_name;type:varchar(128);comment:配置名称;NOT NULL"`
	ConfigType    int8            `json:"configType" gorm:"column:config_type;type:tinyint;comment:配置类型: 1=定时红包, 2=手动触发;NOT NULL"`
	GroupId       string          `json:"groupId" gorm:"column:group_id;type:varchar(128);comment:Telegram群组ID;index;NOT NULL"`
	GroupName     string          `json:"groupName" gorm:"column:group_name;type:varchar(255);comment:群组名称"`
	PacketType    int8            `json:"packetType" gorm:"column:packet_type;type:tinyint;comment:红包类型: 1=普通红包, 2=手气红包;NOT NULL"`
	TotalAmount   decimal.Decimal `json:"totalAmount" gorm:"column:total_amount;type:decimal(25,8);comment:红包总金额;NOT NULL"`
	TotalCount    int             `json:"totalCount" gorm:"column:total_count;type:int;comment:红包总个数;NOT NULL"`
	Symbol        string          `json:"symbol" gorm:"column:symbol;type:varchar(20);comment:币种"`
	ExpireMinutes int             `json:"expireMinutes" gorm:"column:expire_minutes;type:int;default:10;comment:红包过期时间(分钟)"`
	MaxGrabAmount decimal.Decimal `json:"maxGrabAmount" gorm:"column:max_grab_amount;type:decimal(25,8);default:0;comment:单人最大可领取金额(0=不限制)"`
	Lang          string          `json:"lang" gorm:"column:lang;type:varchar(10);default:'vi';comment:消息语言: vi=越南语, id=印尼语, en=英语, zh=中文"`
	BlessingWords string          `json:"blessingWords" gorm:"column:blessing_words;type:varchar(255);comment:祝福语"`
	CronExpr      string          `json:"cronExpr" gorm:"column:cron_expr;type:varchar(128);comment:Cron表达式（如: 0 0 12 * * *）"`
	TimeZone      string          `json:"timeZone" gorm:"column:time_zone;type:varchar(64);comment:时区"`
	StartDate     *time.Time      `json:"startDate" gorm:"column:start_date;type:timestamp;default:null;comment:开始日期"`
	EndDate       *time.Time      `json:"endDate" gorm:"column:end_date;type:timestamp;default:null;comment:结束日期"`
	Status        int8            `json:"status" gorm:"column:status;type:tinyint;comment:状态: 1=启用, 2=禁用, 3=已删除;index"`
	LastExecTime  *time.Time      `json:"lastExecTime" gorm:"column:last_exec_time;type:timestamp;default:null;comment:上次执行时间"`
	NextExecTime  *time.Time      `json:"nextExecTime" gorm:"column:next_exec_time;type:timestamp;default:null;comment:下次执行时间"`
	ExecCount     int             `json:"execCount" gorm:"column:exec_count;type:int;comment:已执行次数"`
	CreatorId     int64           `json:"creatorId" gorm:"column:creator_id;type:bigint;comment:创建者ID"`
	CreatorName   string          `json:"creatorName" gorm:"column:creator_name;type:varchar(128);comment:创建者名称"`
	Remark        string          `json:"remark" gorm:"column:remark;type:varchar(500);comment:备注"`
}

func (*TgRedPacketConfig) TableName() string {
	return "magic_tg_red_packet_config"
}

func NewTgRedPacketConfig() *TgRedPacketConfig {
	return &TgRedPacketConfig{}
}
