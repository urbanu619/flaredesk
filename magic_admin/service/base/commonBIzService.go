package base

import (
	"fmt"
	"github.com/samber/lo"
	"go_server/base/config"
	"go_server/base/core"
	"go_server/global"
	"go_server/utils"
	"gorm.io/gorm"
	"reflect"
	"strings"
	"sync"
	"time"
)

// 业务库公共服务

type BizCommonService struct {
	DbAlias string
}

var lock sync.RWMutex

func (s *BizCommonService) mustAlias() string {
	lock.RLock()
	defer lock.RUnlock()
	if s.DbAlias == "" {
		return global.DefaultAlias()
	}
	return s.DbAlias
}

func (s *BizCommonService) SetDbAlias(alias string) {
	lock.RLock()
	defer lock.RUnlock()
	s.DbAlias = alias
}

func (s *BizCommonService) DB() *gorm.DB {
	db, _ := global.BizDBByAlias(s.mustAlias())
	return db
}

func (s *BizCommonService) GetColumnComment(dbAlias string, tableName string) (data []Column, err error) {
	var entities []Column
	sql := `
	SELECT 
    c.COLUMN_NAME column_name,
    c.DATA_TYPE data_type,
    CASE c.DATA_TYPE
        WHEN 'longtext' THEN c.CHARACTER_MAXIMUM_LENGTH
        WHEN 'varchar' THEN c.CHARACTER_MAXIMUM_LENGTH
        WHEN 'double' THEN CONCAT_WS(',', c.NUMERIC_PRECISION, c.NUMERIC_SCALE)
        WHEN 'decimal' THEN CONCAT_WS(',', c.NUMERIC_PRECISION, c.NUMERIC_SCALE)
        WHEN 'int' THEN c.NUMERIC_PRECISION
        WHEN 'bigint' THEN c.NUMERIC_PRECISION
        ELSE '' 
    END AS data_type_long,
    c.COLUMN_COMMENT column_comment,
    CASE WHEN kcu.COLUMN_NAME IS NOT NULL THEN 1 ELSE 0 END AS primary_key,
    c.ORDINAL_POSITION
FROM 
    INFORMATION_SCHEMA.COLUMNS c
LEFT JOIN 
    INFORMATION_SCHEMA.KEY_COLUMN_USAGE kcu 
ON 
    c.TABLE_SCHEMA = kcu.TABLE_SCHEMA 
    AND c.TABLE_NAME = kcu.TABLE_NAME 
    AND c.COLUMN_NAME = kcu.COLUMN_NAME 
    AND kcu.CONSTRAINT_NAME = 'PRIMARY'
WHERE 
    c.TABLE_NAME = ? 
    AND c.TABLE_SCHEMA = ?
ORDER BY 
    c.ORDINAL_POSITION;`
	db, err := global.BizDBByAlias(s.mustAlias())
	if err != nil {
		return entities, err
	}
	core.Log.Infof("dbName:%s tableName:%s", db.Migrator().CurrentDatabase(), tableName)
	err = db.Raw(sql, tableName, db.Migrator().CurrentDatabase()).Scan(&entities).Error
	return entities, err
}

// 反射获得结构体字段与备注信息 获取json gorm.comment

type ClmInfo struct {
	Field   string `json:"field"`
	Json    string `json:"json"`
	Comment string `json:"comment"`
}

func (s *BizCommonService) GetColumnCommentFromStruct(anyStruct any) []ClmInfo {
	t := reflect.TypeOf(anyStruct)
	var fieldInfoList []ClmInfo
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		// 处理嵌套结构体
		//if field.Anonymous && field.Type.Kind() == reflect.Struct {
		//	// 递归处理嵌套结构体的字段
		//	nestedFields := s.GetColumnCommentFromStruct(reflect.New(field.Type).Elem().Interface())
		//	fieldInfoList = append(fieldInfoList, nestedFields...)
		//	continue
		//}
		// 处理嵌套结构体（包括指针类型）
		if field.Anonymous {
			fieldType := field.Type
			// 如果是指针类型，获取其指向的类型
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}
			// 如果是结构体类型，递归处理
			if fieldType.Kind() == reflect.Struct {
				// 创建该类型的零值实例进行递归
				nestedFields := s.GetColumnCommentFromStruct(reflect.New(fieldType).Elem().Interface())
				fieldInfoList = append(fieldInfoList, nestedFields...)
				continue
			}
		}
		// 获取json标签值
		jsonTag := field.Tag.Get("json")
		// 获取gorm标签中的comment值
		gormTag := field.Tag.Get("gorm")
		comment := s.extractCommentFromGormTag(gormTag)
		// 添加到结果列表
		fieldInfoList = append(fieldInfoList, ClmInfo{
			Field:   field.Name,
			Json:    jsonTag,
			Comment: comment,
		})
	}
	return fieldInfoList
}

func (s *BizCommonService) extractCommentFromGormTag(gormTag string) string {
	pairs := strings.Split(gormTag, ";")
	for _, pair := range pairs {
		if strings.HasPrefix(pair, "comment:") {
			return strings.TrimPrefix(pair, "comment:")
		}
	}
	return ""
}

// 通用CSV导出方法
// 表名 传入需要导出字段名 表字段信息

func ExportCsv[T any](db *gorm.DB, fields []string, colInfo []ClmInfo) (string, error) {
	core.Log.Infof("申请导出字段:%+v", fields)
	allRecords, err := GetMore[T](db)
	if err != nil {
		return "", err
	}
	core.Log.Infof("导出数据总量:%d", len(allRecords))
	if len(allRecords) == 0 {
		return "", nil
	}
	var csvHeard []string
	var csvCols []string
	for _, cls := range colInfo {
		if lo.Contains(fields, cls.Field) {
			csvHeard = append(csvHeard, cls.Comment)
			csvCols = append(csvCols, cls.Field)
		}
	}
	core.Log.Infof("导出表头:%+v", csvHeard)
	core.Log.Infof("可导出csvCols字段:%+v", csvCols)
	var csvData [][]string
	for _, item := range allRecords {
		var row []string
		val := reflect.ValueOf(item)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		for _, fieldName := range csvCols {
			fieldVal := val.FieldByName(fieldName)
			if !fieldVal.IsValid() {
				// 如果直接字段找不到，尝试从嵌套结构体查找
				for i := 0; i < val.NumField(); i++ {
					field := val.Type().Field(i)
					if field.Anonymous {
						nestedVal := val.Field(i)
						// 处理嵌套结构体指针
						if nestedVal.Kind() == reflect.Ptr {
							if nestedVal.IsNil() {
								continue
							}
							nestedVal = nestedVal.Elem()
						}
						if nestedVal.IsValid() {
							nestedField := nestedVal.FieldByName(fieldName)
							if nestedField.IsValid() {
								fieldVal = nestedField
								break
							}
						}
					}
				}
			}
			var strVal string
			if fieldVal.IsValid() {
				strVal = fmt.Sprintf("%v", fieldVal.Interface())
			}
			row = append(row, strVal)
		}
		csvData = append(csvData, row)
	}

	filename := fmt.Sprintf("export_%s.csv", time.Now().Format("0102150405"))
	viewPath := utils.ToPath(config.EnvConf().File.Path, config.EnvConf().File.ProxyPath, filename)
	storePath := utils.ToPath(config.EnvConf().File.StorePath, filename)

	_, err = config.ExportCsv(csvHeard, csvData, storePath)
	if err != nil {
		return "", err
	}
	return viewPath, nil
}
