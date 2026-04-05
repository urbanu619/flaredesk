package base

import (
	"fmt"
	"go_server/utils"
	"gorm.io/gorm"
	"strings"
)

type BaseRequest struct {
	Id   uint           `json:"id" form:"id"`
	Data map[string]any `json:"data"`
}

type ListRequestColumn struct {
	Key   string  `json:"key"`
	Op    *string `json:"op"`
	Value any     `json:"value"`
}

// 定义范型

type ListRequestInterface[T any] interface {
	Query(db *gorm.DB) (*ListResponse[T], error)
}

// ListRequest T 为 model
type ListRequest[T any] struct {
	Current  int                 `json:"current" query:"current" form:"current" `     //当前页
	PageSize int                 `json:"pageSize" query:"pageSize" form:"pageSize"  ` //每页条数
	Column   []ListRequestColumn `json:"column"`
	Order    *string             `json:"order" form:"order" `
}

func NewListRequest[T any]() *ListRequest[T] {
	request := new(ListRequest[T])
	request.PageSize = 10
	request.Current = 1
	return request
}

type ListResponse[T any] struct {
	List   []T         `json:"list"`
	Paging *Pagination `json:"paging"`
}

func (req *ListRequest[T]) Query(db *gorm.DB) (*ListResponse[T], error) {
	list := make([]T, 0)
	query := db.Model(&list)
	if len(req.Column) > 0 {
		for _, column := range req.Column {
			// 验证列名安全性
			if !utils.SQLSecurity.ValidateColumnName(column.Key) {
				continue // 跳过不安全的列名
			}
			// 验证操作符安全性
			if column.Op != nil && !utils.SQLSecurity.ValidateOperator(*column.Op) {
				continue // 跳过不安全的操作符
			}
			if column.Op != nil && !strings.Contains(">,=,<,>=,<=,between,like", *column.Op) {
				continue
			}
			switch {
			case column.Op != nil && *column.Op == "between":
				if t, ok := column.Value.([]any); ok {
					if len(t) > 1 {
						query.Where(fmt.Sprintf("`%s` %s ? and ?", column.Key, *column.Op), t[0], t[1])
					}
				}
			case column.Op != nil && *column.Op == "like":
				if t, ok := column.Value.(string); ok {
					query.Where(fmt.Sprintf("`%s` like ?", column.Key), "%"+strings.TrimSpace(t)+"%")
				}
			case column.Op != nil:
				query.Where(fmt.Sprintf("`%s` %s ?", column.Key, *column.Op), column.Value)
			default:
				query.Where(fmt.Sprintf("`%s` = ?", column.Key), column.Value)
			}
		}
	}
	// 验证ORDER BY安全性
	if req.Order != nil && utils.SQLSecurity.ValidateOrderBy(*req.Order) {
		query = query.Order(utils.WordsToSnakeCase(*req.Order))
	}
	paging := NewPagination()
	paging.Current = req.Current
	paging.PageSize = req.PageSize
	if err := query.Count(&paging.Total).Error; err != nil {
		return nil, err
	}
	paging.Computer() //分页器计数
	query = query.Limit(paging.PageSize).Offset(paging.StartNums)
	err := query.Find(&list).Error
	if err != nil {
		return nil, err
	}
	return &ListResponse[T]{
		List:   list,
		Paging: paging,
	}, nil
}
