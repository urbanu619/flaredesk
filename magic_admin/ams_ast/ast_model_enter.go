package ams_ast

import (
	"fmt"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// 表信息结构体

type TableInfo struct {
	Module     string
	TableName  string
	StructName string
	Fields     []*FieldInfo
	HasDecimal bool
	HasTime    bool
	HasJson    bool
}

// 字段信息结构体

type FieldInfo struct {
	Name         string
	Type         string
	GormTag      string
	JsonTag      string
	Comment      string
	DefaultValue string
	IsPrimary    bool   // 是否主键
	IsNullable   bool   // 是否可为空
	IsIndex      bool   // 是否索引
	IsUnique     bool   // 是否唯一约束索引
	UniqueName   string // 约束名
}

// 获取表的列信息

type Column struct {
	ColumnName       string `gorm:"column:COLUMN_NAME"`
	DataType         string `gorm:"column:DATA_TYPE"`
	ColumnType       string `gorm:"column:COLUMN_TYPE"`
	ColumnKey        string `gorm:"column:COLUMN_KEY"`
	IsNullable       string `gorm:"column:IS_NULLABLE"`
	ColumnDefault    string `gorm:"column:COLUMN_DEFAULT"`
	ColumnComment    string `gorm:"column:COLUMN_COMMENT"`
	NumericPrecision *int64 `gorm:"column:NUMERIC_PRECISION"`
	NumericScale     *int64 `gorm:"column:NUMERIC_SCALE"`
}

// 唯一约束信息

type Constraint struct {
	ConstraintName       string `gorm:"column:CONSTRAINT_NAME"`
	ColumnName           string `gorm:"column:COLUMN_NAME"`
	ReferencedTableName  string `gorm:"column:REFERENCED_TABLE_NAME"`
	ReferencedColumnName string `gorm:"column:REFERENCED_COLUMN_NAME"`
}

// 索引信息

type IndexInfo struct {
	IndexName  string `gorm:"column:INDEX_NAME"`
	ColumnName string `gorm:"column:COLUMN_NAME"`
	NonUnique  int    `gorm:"column:NON_UNIQUE"`
}

// 模型注册器

type ModelRegister struct {
	DbAlias           string // 数据库别名
	ModelTemplateFile string `json:"modelTemplateFile"`
	TargetDir         string `json:"targetDir"`  // 文件生成目标文件夹
	TargetFile        string `json:"targetFile"` // 模型文件
	TableName         string `json:"tableName"`
	StructName        string `json:"structName"` // 生成的结构体名称
	BaseAst
}

// 数据库 表名

func BuildModelRegister(dbAlias, dbName, tableName string) *ModelRegister {
	register := &ModelRegister{
		DbAlias:    dbAlias,                                   // 数据库别名 biz
		TableName:  tableName,                                 // 表名称 demo_user
		StructName: CapitalizeOrLower(toCamelCase(tableName)), // 首字母大写驼峰写法 DemoUser
	}
	register.ModelTemplateFile = register.modelTemplateFile() // 模版文件地址
	register.TargetFile = register.modelTargetFile(tableName) // 模型目标文件
	return register
}

// 生成表结构

func (s *ModelRegister) Initialize(db *gorm.DB) error {
	var files *template.Template
	files, err := template.ParseFiles(s.ModelTemplateFile)
	if err != nil {
		return err
	}
	//slog.Infof("模型模版文件:%s", s.ModelTemplateFile)
	err = os.MkdirAll(filepath.Dir(s.TargetFile), os.ModePerm)
	if err != nil {
		return err
	}

	if FileExists(s.TargetFile) {
		err = os.Remove(s.TargetFile)
		if err != nil {
			return err
		}
	}
	var file *os.File
	file, err = os.Create(s.TargetFile)
	if err != nil {
		return err
	}
	data, err := s.getTableInfo(db)
	if err != nil {
		return err
	}
	err = files.Execute(file, data)
	_ = file.Close()
	if err != nil {
		return err
	}
	return nil
}

// 回滚

func (s *ModelRegister) RollbackModel() error {
	if FileExists(s.TargetFile) {
		err := os.Remove(s.TargetFile)
		if err != nil {
			return err
		}
	}
	return nil
}

// 模板文件地址

func (s *ModelRegister) modelTemplateFile() string {
	return filepath.Join(template_root, "model/model.go.tpl")
}

// 生成模型文件路径

func (s *ModelRegister) modelTargetFile(tableName string) string {
	return filepath.Join(s.modelTargetDir(), toSnakeCase(tableName)+".go")
}

// 模型目标文件夹 -- go_server/server/model/server_model_biz/DbAlias

func (s *ModelRegister) modelTargetDir() string {
	return filepath.Join(modelTargetPath, server_model_biz, s.DbAlias)
}

// 获取表信息

func (s *ModelRegister) getTableInfo(db *gorm.DB) (*TableInfo, error) {
	if db.Dialector.Name() == "sqlite" {
		return nil, fmt.Errorf("SQLite 暂不支持从数据库逆向生成模型，请使用 MySQL 或手写 model")
	}
	tableInfo := &TableInfo{
		Module:     s.DbAlias,
		TableName:  s.TableName,
		StructName: s.StructName,
		Fields:     make([]*FieldInfo, 0),
		HasDecimal: false,
		HasTime:    false,
		HasJson:    false,
	}

	cols, err := s.getTableCols(db)
	if err != nil {
		return nil, err
	}
	constraint, err := s.getTableConstraint(db)
	if err != nil {
		return nil, err
	}
	indexes, err := s.getTableIndexInfo(db)
	if err != nil {
		return nil, err
	}
	for _, col := range cols {
		// 是否为Primary index
		// 获取唯一约束信息
		uniqueName, isUnique := s.colUniqueInfo(constraint, col.ColumnName)
		isIndex, isPrimary := s.colIndexInfo(indexes, col.ColumnName)
		if col.DataType == "decimal" {
			tableInfo.HasDecimal = true
		}
		if col.DataType == "datetime" {
			tableInfo.HasTime = true
		}
		if col.DataType == "json" {
			tableInfo.HasJson = true
		}
		filedInfo := &FieldInfo{
			Name:         CapitalizeOrLower(toCamelCase(col.ColumnName)),
			Type:         s.sqlTypeToGoType(col.DataType, col.ColumnType),
			GormTag:      "",
			JsonTag:      fmt.Sprintf(`json:"%s"`, toCamelCase(col.ColumnName)),
			Comment:      col.ColumnComment,
			DefaultValue: col.ColumnDefault,
			IsNullable:   col.IsNullable == "YES",
			IsPrimary:    isPrimary,
			IsIndex:      isIndex,
			IsUnique:     isUnique,
			UniqueName:   uniqueName,
		}
		// 组装 GormTag JsonTag
		// `json:"userId" gorm:"uniqueIndex:idx_asset_user_coin;comment:用户ID;NOT NULL"`
		// column:rate;comment:汇率
		gormTags := make([]string, 0)
		gormTags = append(gormTags, fmt.Sprintf("column:%s", col.ColumnName))
		gormTags = append(gormTags, fmt.Sprintf("type:%s", col.ColumnType))
		gormTags = append(gormTags, fmt.Sprintf("comment:%s", col.ColumnComment))
		if isIndex {
			if isPrimary {
				gormTags = append(gormTags, "primarykey")
			} else {
				gormTags = append(gormTags, "index")
			}
		}
		if isUnique && !isPrimary {
			gormTags = append(gormTags, fmt.Sprintf("unique:%s", filedInfo.UniqueName))
		}
		if !filedInfo.IsNullable {
			gormTags = append(gormTags, "NOT NULL")
		}
		filedInfo.GormTag = fmt.Sprintf(`gorm:"%s"`, strings.Join(gormTags, ";"))
		tableInfo.Fields = append(tableInfo.Fields, filedInfo)
	}
	return tableInfo, nil
}

func (s *ModelRegister) colUniqueInfo(items []*Constraint, colName string) (string, bool) {
	for _, item := range items {
		if item.ColumnName == colName {
			return item.ConstraintName, true
		}
	}

	return "", false
}

func (s *ModelRegister) colIndexInfo(items []*IndexInfo, colName string) (isIndex bool, isPrimary bool) {
	for _, item := range items {
		if item.ColumnName == colName {
			return true, item.IndexName == "PRIMARY"
		}
	}
	return false, false
}

// 获取列信息

func (s *ModelRegister) getTableCols(db *gorm.DB) ([]*Column, error) {
	columns := make([]*Column, 0)
	if err := db.Raw(`
		SELECT 
			COLUMN_NAME,
			DATA_TYPE,
			COLUMN_TYPE,
			COLUMN_KEY,
			IS_NULLABLE,
			COLUMN_DEFAULT,
			COLUMN_COMMENT,
			NUMERIC_PRECISION,
			NUMERIC_SCALE
		FROM INFORMATION_SCHEMA.COLUMNS 
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
		ORDER BY ORDINAL_POSITION`, s.TableName).Scan(&columns).Error; err != nil {
		panic(err)
	}
	return columns, nil
}

// 获取表唯一约束信息

func (s *ModelRegister) getTableConstraint(db *gorm.DB) ([]*Constraint, error) {
	constraints := make([]*Constraint, 0)
	if err := db.Raw(`
		SELECT 
			k.CONSTRAINT_NAME,
			k.COLUMN_NAME,
			k.REFERENCED_TABLE_NAME,
			k.REFERENCED_COLUMN_NAME
		FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS t
		JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE k
		USING(CONSTRAINT_NAME,TABLE_SCHEMA,TABLE_NAME) 
		WHERE t.TABLE_SCHEMA = DATABASE() 
		AND t.TABLE_NAME = ?
		AND t.CONSTRAINT_TYPE = 'UNIQUE'
		ORDER BY k.ORDINAL_POSITION`, s.TableName).Scan(&constraints).Error; err != nil {
		return nil, err
	}
	return constraints, nil
}

// 获取表索引信息

func (s *ModelRegister) getTableIndexInfo(db *gorm.DB) ([]*IndexInfo, error) {
	// 获取索引信息
	indexes := make([]*IndexInfo, 0)
	if err := db.Raw(`
		SELECT 
			INDEX_NAME,
			COLUMN_NAME,
			NON_UNIQUE
		FROM INFORMATION_SCHEMA.STATISTICS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
		ORDER BY INDEX_NAME, SEQ_IN_INDEX`, s.TableName).Scan(&indexes).Error; err != nil {
		return nil, err
	}
	return indexes, nil
}

// SQL类型转Go类型

func (s *ModelRegister) sqlTypeToGoType(dataType, columnType string) string {
	switch strings.ToLower(dataType) {
	case "tinyint":
		if strings.Contains(columnType, "tinyint(1)") {
			return "bool"
		}
		return "int8"
	case "smallint":
		return "int16"
	case "mediumint", "int":
		return "int"
	case "bigint":
		return "int64"
	case "float":
		return "float32"
	case "double":
		return "float64"
	case "char", "varchar", "tinytext", "text", "mediumtext", "longtext", "enum", "set":
		return "string"
	case "date", "datetime", "timestamp", "time":
		return "time.Time"
	case "binary", "varbinary", "tinyblob", "blob", "mediumblob", "longblob":
		return "[]byte"
	case "json":
		return "datatypes.JSON"
	case "bit":
		return "[]uint8"
	case "decimal":
		return "decimal.Decimal"
	default:
		return "interface{}"
	}
}
