package system

import (
	"go_server/model/common"
)

type FileJob struct {
	common.GormFullModel
	Path     string `json:"path" gorm:"column:path;type:varchar(255);comment:数据文件路径"`
	Total    int64  `json:"total" gorm:"comment:数据量"`
	Size     string `json:"size" gorm:"type:varchar(80);comment:数据大小"`
	Status   int    `json:"status" gorm:"comment:任务状态:1=执行中 2=执行完成 3=执行失败"`
	Desc     string `json:"desc" gorm:"type:varchar(512);comment:错误信息"`
	CreateBy string `json:"createBy" gorm:"type:varchar(80);comment:任务创建者"`
}

func (*FileJob) TableName() string {
	return common.ModelPrefix + "file_job"
}

func (*FileJob) Comment() string {
	return "文件导出任务"
}

func NewFileJob() *FileJob {
	return &FileJob{}
}
