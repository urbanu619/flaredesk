package base

import (
	"fmt"
	"gorm.io/gorm"
	"math"
)

type Pagination struct {
	Current   int   `json:"current" query:"current" form:"current" `         //当前页
	PageSize  int   `json:"pageSize" query:"pageSize" form:"pageSize"  `     //每页条数
	Total     int64 `json:"total" query:"total" form:"total" `               //总条数
	Count     int   `json:"count" query:"count" form:"count" `               //总页数
	StartNums int   `json:"startNums"   query:"startNums" form:"startNums" ` //起始条数
}

func NewPagination() *Pagination {
	return &Pagination{}
}

func InitPaging(p, s int) *Pagination {
	return &Pagination{Current: p, PageSize: s}
}

func (p *Pagination) Computer() {
	if p.Current < 1 {
		p.Current = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	p.StartNums = p.PageSize * (p.Current - 1)
	count := math.Ceil(float64(p.Total) / float64(p.PageSize))
	p.Count = int(count)
}

// 获取单条记录

func GetOne[T any](db *gorm.DB, args ...interface{}) (*T, bool) {
	db = buildEqualDB(db, args...)
	r := new(T)
	if err := db.First(&r).Error; err != nil {
		return r, false
	}
	return r, true
}

func GetLast[T any](db *gorm.DB, args ...interface{}) (*T, bool) {
	db = buildEqualDB(db, args...)
	r := new(T)
	if err := db.Last(&r).Error; err != nil {
		return r, false
	}
	return r, true
}

// 统计条数

func CountRow[T any](db *gorm.DB, args ...interface{}) int64 {
	db = buildEqualDB(db, args...)
	r := new(T)
	var total int64 = 0
	db.Model(&r).Count(&total)
	return total
}

// 获取所有数据 -- 一次性获取所有 -- 对于已知数量不多的情况下使用

func GetMore[T any](db *gorm.DB, args ...interface{}) ([]*T, error) {
	db = buildEqualDB(db, args...)
	res := make([]*T, 0)
	var total int64
	err := db.Model(&res).Count(&total).Error
	if err != nil {
		return res, err
	}
	if total == 0 {
		return res, nil
	}
	if err := db.Find(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}

// 分页查询 -- 分页查询

func GetMoreWithPage[T any](db *gorm.DB, paging *Pagination, args ...interface{}) ([]*T, *Pagination, error) {
	if paging == nil {
		paging = InitPaging(1, 10)
	}
	db = buildEqualDB(db, args...)
	res := make([]*T, 0)

	err := db.Model(&res).Count(&paging.Total).Error
	if err != nil {
		return res, paging, err
	}
	paging.Computer()
	if paging.Total == 0 {
		return res, paging, nil
	}
	err = db.Model(&res).Limit(paging.PageSize).Offset(paging.StartNums).Find(&res).Error
	return res, paging, err
}

// 批量分页查询

func getRowsWithPage[T any](db *gorm.DB, paging *Pagination) ([]*T, error) {
	paging.Computer()
	res := make([]*T, 0)
	err := db.Model(&res).Limit(paging.PageSize).Offset(paging.StartNums).Find(&res).Error
	return res, err
}

const BatchPageSize = 5000

func ObtainAllRowsInBatches[T any](db *gorm.DB) ([]*T, error) {
	count := CountRow[T](db)
	// 分页循环读取用户数据
	list := make([]*T, 0, count)
	paging := &Pagination{
		PageSize: BatchPageSize,
		Total:    count,
	}
	paging.Computer()
	if count > 0 {
		for i := 0; i < paging.Count; i++ {
			items, err := getRowsWithPage[T](db, &Pagination{
				PageSize: BatchPageSize,
				Current:  i + 1,
				Total:    count,
			})

			if err != nil {
				return nil, err
			}
			list = append(list, items...) // 直接追加整个切片
		}
	}
	return list, nil
}

// 构建查询

func buildEqualDB(db *gorm.DB, args ...interface{}) *gorm.DB {
	tms := len(args) / 2
	if tms > 0 {
		for i := 0; i < tms; i++ {
			db = db.Where(fmt.Sprintf("%s = ?", args[i*2]), args[i*2+1])
		}
	}
	return db
}

const QueryPageSize int64 = 100000 // 批量查询数据最大条数
