package app

// 引入关联包
import (
	"github.com/shopspring/decimal"
)

type MagicUserQuota struct {
	CreatedAt                   int64           `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt                   int64           `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	UserId                      int64           `json:"userId" gorm:"column:user_id;type:bigint;comment:id;primarykey;NOT NULL"`
	Uid                         string          `json:"uid" gorm:"column:uid;type:varchar(64);comment:交易所ID:;index;unique:uni_magic_user_quota_uid"`
	ParentId                    int64           `json:"parentId" gorm:"column:parent_id;type:bigint;comment:上级ID"`
	ParentUid                   string          `json:"parentUid" gorm:"column:parent_uid;type:varchar(64);comment:上级交易所ID:;index"`
	ParentSetTime               int64           `json:"parentSetTime" gorm:"column:parent_set_time;type:bigint;comment:首次设置上级时间"`
	ParentUpdateTime            int64           `json:"parentUpdateTime" gorm:"column:parent_update_time;type:bigint;comment:更新上级时间"`
	LockLevel                   int8            `json:"lockLevel" gorm:"column:lock_level;type:tinyint;comment:锁定团队等级"`
	AcLevel                     int8            `json:"acLevel" gorm:"column:ac_level;type:tinyint;comment:团队业绩等级"`
	Level                       int8            `json:"level" gorm:"column:level;type:tinyint;comment:团队等级=max(lockLevel,acLevel)"`
	InviteCount                 int64           `json:"inviteCount" gorm:"column:invite_count;type:bigint;comment:直推人数"`
	TeamCount                   int64           `json:"teamCount" gorm:"column:team_count;type:bigint;comment:团队人数"`
	BuyParentNodes              int64           `json:"buyParentNodes" gorm:"column:buy_parent_nodes;type:bigint;comment:购买节点数量"`
	TeamParentNodes             int64           `json:"teamParentNodes" gorm:"column:team_parent_nodes;type:bigint;comment:团队节点数量"`
	BuyEcologicalNodes          int64           `json:"buyEcologicalNodes" gorm:"column:buy_ecological_nodes;type:bigint;comment:购买生态节点数量"`
	TeamEcologicalNodes         int64           `json:"teamEcologicalNodes" gorm:"column:team_ecological_nodes;type:bigint;comment:团队生态节点数量"`
	PersonAchievement           decimal.Decimal `json:"personAchievement" gorm:"column:person_achievement;type:decimal(25,8);comment:个人业绩USD"`
	TeamTodayAchievement        decimal.Decimal `json:"teamTodayAchievement" gorm:"column:team_today_achievement;type:decimal(25,8);comment:当日团队业绩"`
	TeamAchievement             decimal.Decimal `json:"teamAchievement" gorm:"column:team_achievement;type:decimal(25,8);comment:团队总业绩"`
	LargeRegionUserId           int64           `json:"largeRegionUserId" gorm:"column:large_region_user_id;type:bigint;comment:大区用户ID"`
	LargeRegionAchievement      decimal.Decimal `json:"largeRegionAchievement" gorm:"column:large_region_achievement;type:decimal(25,8);comment:大区业绩"`
	FewTeamAchievement          decimal.Decimal `json:"fewTeamAchievement" gorm:"column:few_team_achievement;type:decimal(25,8);comment:小团队业绩"`
	TodayFewTeamAchievement     decimal.Decimal `json:"todayFewTeamAchievement" gorm:"column:today_few_team_achievement;type:decimal(25,8);comment:今日小团队业绩"`
	PersonRecharge              decimal.Decimal `json:"personRecharge" gorm:"column:person_recharge;type:decimal(25,8);comment:个人累计充值"`
	TeamRecharge                decimal.Decimal `json:"teamRecharge" gorm:"column:team_recharge;type:decimal(25,8);comment:团队累计充值"`
	PersonWithdrawal            decimal.Decimal `json:"personWithdrawal" gorm:"column:person_withdrawal;type:decimal(25,8);comment:个人累计提现"`
	TeamWithdrawal              decimal.Decimal `json:"teamWithdrawal" gorm:"column:team_withdrawal;type:decimal(25,8);comment:团队累计提现"`
	GiftPersonAchievement       decimal.Decimal `json:"giftPersonAchievement" gorm:"column:gift_person_achievement;type:decimal(25,8);comment:赠送质押业绩"`
	GiftLargeRegionAchievement  decimal.Decimal `json:"giftLargeRegionAchievement" gorm:"column:gift_large_region_achievement;type:decimal(25,8);comment:赠送大区业绩-直通车业绩"`
	IsOld                       bool            `json:"isOld" gorm:"column:is_old;type:tinyint(1);comment:团队累计提现"`
	IsExemptionLargeAchievement bool            `json:"isExemptionLargeAchievement" gorm:"column:is_exemption_large_achievement;type:tinyint(1);comment:是否免除大区业绩"`
	InflatedFewTeamAchievement  decimal.Decimal `json:"inflatedFewTeamAchievement" gorm:"column:inflated_few_team_achievement;type:decimal(25,8);comment:虚增小区业绩"`
}

func (*MagicUserQuota) TableName() string {
	return "magic_user_quota"
}

func NewMagicUserQuota() *MagicUserQuota {
	return &MagicUserQuota{}
}
