package system

import (
	"github.com/jameskeane/bcrypt"
	"go_server/model/common"
	"time"
)

type Administrator struct {
	ID int64 `json:"id" gorm:"primarykey;comment:id"`
	common.GormAllTimeModel
	Nickname      string     `json:"nickname" gorm:"column:nickname;type:char(20);comment:昵称"`
	Username      string     `json:"username" gorm:"column:username;type:varchar(50);unique;comment:用户名"`
	Password      string     `json:"-" gorm:"column:password;type:varchar(100);comment:密码"`
	Salt          string     `json:"-" gorm:"column:salt;varchar(50);comment:延码"`
	RoleId        int64      `json:"role_id" gorm:"column:role_id;comment:角色Id"`
	Avatar        string     `json:"avatar" gorm:"type:varchar(100);comment:头像地址"`
	Lock          bool       `json:"lock" gorm:"column:lock;default:false;comment:是否锁定"`
	Token         string     `json:"-" gorm:";column:token;type:text;comment:token"`
	LastLoginIp   string     `json:"last_login_ip" gorm:"column:last_login_ip;type:varchar(20);comment:最后登录ip"`
	LastLoginTime *time.Time `json:"last_login_time" gorm:"column:last_login_time;comment:最后登录时间"`
	GoogleKey     string     `json:"-" gorm:"column:google_key;type:varchar(100);comment:谷歌密钥"`
}

func (*Administrator) TableName() string {
	return common.ModelPrefix + "administrator"
}

func NewAdministrator() *Administrator {
	return &Administrator{}
}

func (*Administrator) Comment() string {
	return "管理员"
}

func (r *Administrator) EncodePassword(passWord string) (enPassWord, salt string) {
	salt, err := bcrypt.Salt(bcrypt.DefaultRounds)
	if err != nil {
		panic(err)
	}
	enPassWord, err = bcrypt.Hash(passWord, salt)
	if err != nil {
		panic(err)
	}
	return
}

func (r *Administrator) CheckPassWord(checkPassWord string) bool {
	checkPassWordEn, _ := bcrypt.Hash(checkPassWord, r.Salt)
	if checkPassWordEn == r.Password {
		return true
	}
	return false
}

type AdministratorLog struct {
	common.GormFullModel
	AdminId   int64  `json:"adminId" gorm:"comment:管理员Id"`
	Method    string `json:"method" gorm:"type:varchar(20);comment:请求方法"`
	Path      string `json:"path" gorm:"type:varchar(200);comment:请求方法"`
	Ip        string `json:"ip" gorm:"type:varchar(20);comment:ip"`
	UserAgent string `json:"userAgent" gorm:"type:text;comment:用户请求信息"`
}

func (*AdministratorLog) TableName() string {
	return common.ModelPrefix + "administrator_log"
}

func NewAdministratorLog() *AdministratorLog {
	return &AdministratorLog{}
}

func (*AdministratorLog) Comment() string {
	return "管理员日志"
}
