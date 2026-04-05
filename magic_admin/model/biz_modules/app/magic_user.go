package app

// 引入关联包
import (
	"time"
)

type MagicUser struct {
	Id                int64     `json:"id" gorm:"column:id;type:bigint;comment:id;primarykey;NOT NULL"`
	CreatedAt         int64     `json:"createdAt" gorm:"column:created_at;type:bigint;comment:创建时间"`
	UpdatedAt         int64     `json:"updatedAt" gorm:"column:updated_at;type:bigint;comment:更新时间"`
	DeletedAt         time.Time `json:"deletedAt" gorm:"column:deleted_at;type:datetime;comment:删除时间;index"`
	Uid               string    `json:"uid" gorm:"column:uid;type:varchar(64);comment:交易所ID:;index"`
	OpenId            string    `json:"openId" gorm:"column:open_id;type:varchar(64);comment:平台OPEN_ID"`
	Code              string    `json:"code" gorm:"column:code;type:varchar(32);comment:平台邀请码"`
	RefCode           string    `json:"refCode" gorm:"column:ref_code;type:varchar(32);comment:平台上级邀请码;index"`
	ParentId          int64     `json:"parentId" gorm:"column:parent_id;type:bigint;comment:上级ID"`
	ParentUid         string    `json:"parentUid" gorm:"column:parent_uid;type:varchar(64);comment:上级交易所ID;index"`
	ParentIds         []int64   `json:"parentIds" gorm:"serializer:json;type:json;column:parent_ids;comment:上级ID路径，使用json查询"`
	LastLoginTime     int64     `json:"lastLoginTime" gorm:"column:last_login_time;type:bigint;comment:最近一次登陆时间"`
	LastLoginIp       string    `json:"lastLoginIp" gorm:"column:last_login_ip;type:varchar(200);comment:最后登陆IP"`
	LastLoginRemoteIp string    `json:"lastLoginRemoteIp" gorm:"column:last_login_remote_ip;type:varchar(200);comment:最后登陆RemoteIP"`
	IsReinvestment    bool      `json:"isReinvestment" gorm:"column:is_reinvestment;type:tinyint(1);comment:是否开启复投"`
	LockWithdraw      int8      `json:"lockWithdraw" gorm:"column:lock_withdraw;type:tinyint;comment:锁定提现  0=未锁定  1=锁定"`
	IsWhite           int8      `json:"isWhite" gorm:"column:is_white;type:tinyint;comment:  1=白名单 0=正常用户"`
	IsRoot            bool      `json:"isRoot" gorm:"column:is_root;type:tinyint(1);comment:是否根账号"`
	IsZero            bool      `json:"isZero" gorm:"column:is_zero;type:tinyint(1);comment:是否零号线"`
	Enable            int8      `json:"enable" gorm:"column:enable;type:tinyint;comment:账户是否有效 1=有效 0=无效"`
	LockRedeem        int8      `json:"lockRedeem" gorm:"column:lock_redeem;type:tinyint;comment:赎回锁定"`
	LockStakeProfit   int8      `json:"lockStakeProfit" gorm:"column:lock_stake_profit;type:tinyint;comment:0 正常 1锁定不产出收益"`
	SubsLockWithdraw  int8      `json:"subsLockWithdraw" gorm:"column:subs_lock_withdraw;type:tinyint;comment:锁定伞下提现  0=未锁定  1=锁定"`
	PerformanceRobot  int8      `json:"performanceRobot" gorm:"column:performance_robot;type:tinyint;comment:0 正常 1 业绩机器人"`
	LevelRobot        int8      `json:"levelRobot" gorm:"column:level_robot;type:tinyint;comment:0 正常 1 等级机器人"`
	InternalAccount   int8      `json:"internalAccount" gorm:"column:internal_account;type:tinyint;comment:0 正常 1 内部账号"`
	FireRobot         int8      `json:"fireRobot" gorm:"column:fire_robot;type:tinyint;comment:0 正常 1 点火机器人"`
	Remark            string    `json:"remark" gorm:"column:remark;type:varchar(200);comment:备注信息"`
}

func (*MagicUser) TableName() string {
	return "magic_user"
}

func NewMagicUser() *MagicUser {
	return &MagicUser{}
}
