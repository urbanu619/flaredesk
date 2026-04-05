package common

import "gorm.io/gorm"

const (
	ModelPrefix     = "ams_"
	SplittingSymbol = ","
)

// 自增长ID变更启始值 db.Exec("ALTER TABLE your_table_name AUTO_INCREMENT = 1000")

type GormIdModel struct {
	ID int64 `json:"id" gorm:"primarykey;comment:id"`
}

type GormTimeModel struct {
	CreatedAt int64 `json:"created_at,omitempty" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt int64 `json:"updated_at,omitempty" gorm:"autoUpdateTime;comment:更新时间"`
}

type GormDeleteModel struct {
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间"`
}

type GormBaseModel struct {
	GormIdModel
	GormTimeModel
}

type GormAllTimeModel struct {
	GormTimeModel
	GormDeleteModel
}
type GormFullModel struct {
	GormIdModel
	GormTimeModel
	GormDeleteModel
}
