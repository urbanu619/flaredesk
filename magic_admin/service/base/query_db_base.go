package base

import (
	"fmt"
	"gorm.io/gorm"
)

type QueryBaseHandler[T any] struct {
}

func NewQueryBaseHandler[T any](t *T) *QueryBaseHandler[T] {
	return &QueryBaseHandler[T]{}
}

func (*QueryBaseHandler[T]) DeleteOne(db *gorm.DB) error {
	tempInfo := new(T)
	// 先通过Id获取记录'
	if err := db.First(&tempInfo).Error; err != nil {
		return err
	}
	exCmd := db.Model(&tempInfo).Delete(&tempInfo)
	if exCmd.RowsAffected != 1 {
		return fmt.Errorf("delete fail")
	}
	if exCmd.Error != nil {
		return exCmd.Error
	}
	return nil
}

func (*QueryBaseHandler[T]) Get(db *gorm.DB) (*T, error) {
	tempInfo := new(T)
	// 先通过Id获取记录'
	if err := db.First(&tempInfo).Error; err != nil {
		return tempInfo, err
	}
	return tempInfo, nil
}

func (*QueryBaseHandler[T]) List(db *gorm.DB, req ListRequestInterface[T]) (*ListResponse[T], error) {
	return req.Query(db)
}

func (*QueryBaseHandler[T]) Export(db *gorm.DB, req ListRequestInterface[T]) (*ListResponse[T], error) {
	return req.Query(db)
}
