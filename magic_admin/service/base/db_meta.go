package base

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// MetaListTables 列出用户表（MySQL / SQLite）。
func MetaListTables(db *gorm.DB) ([]Table, error) {
	var entities []Table
	switch db.Dialector.Name() {
	case "sqlite":
		err := db.Raw(`SELECT name AS table_name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name`).Scan(&entities).Error
		return entities, err
	default:
		sql := `SELECT table_name AS table_name FROM information_schema.tables WHERE table_schema = ?`
		err := db.Raw(sql, db.Migrator().CurrentDatabase()).Scan(&entities).Error
		return entities, err
	}
}

// MetaListColumns 列出表字段（MySQL / SQLite）。
func MetaListColumns(db *gorm.DB, tableName string) ([]Column, error) {
	switch db.Dialector.Name() {
	case "sqlite":
		return metaListColumnsSQLite(db, tableName)
	default:
		return metaListColumnsMySQL(db, tableName)
	}
}

func metaListColumnsMySQL(db *gorm.DB, tableName string) ([]Column, error) {
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
    c.ORDINAL_POSITION`
	err := db.Raw(sql, tableName, db.Migrator().CurrentDatabase()).Scan(&entities).Error
	return entities, err
}

func metaListColumnsSQLite(db *gorm.DB, tableName string) ([]Column, error) {
	tn := strings.ReplaceAll(tableName, `"`, `""`)
	var rows []struct {
		Name string `gorm:"column:name"`
		Type string `gorm:"column:type"`
		Pk   int    `gorm:"column:pk"`
	}
	q := fmt.Sprintf(`PRAGMA table_info("%s")`, tn)
	if err := db.Raw(q).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]Column, 0, len(rows))
	for _, r := range rows {
		out = append(out, Column{
			DataType:      r.Type,
			ColumnName:    r.Name,
			DataTypeLong:  r.Type,
			ColumnComment: "",
			PrimaryKey:    r.Pk == 1,
		})
	}
	return out, nil
}
