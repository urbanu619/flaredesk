package utils

import (
	"regexp"
	"strings"
)

// SQLSecurityUtils SQL安全工具类
type SQLSecurityUtils struct{}

// 危险的SQL关键字
var dangerousKeywords = []string{
	"select", "insert", "update", "delete", "drop", "create", "alter", "exec", "execute",
	"union", "script", "javascript", "vbscript", "onload", "onerror", "onclick",
	"--", "/*", "*/", ";", "'", "\"", "xp_", "sp_", "master", "truncate", "shutdown",
}

// 允许的排序字段字符（字母、数字、下划线、点号、反引号、空格、逗号、DESC、ASC）
var orderByPattern = regexp.MustCompile(`^[a-zA-Z0-9_.,\x60\s]+(\s+(ASC|DESC))?(\s*,\s*[a-zA-Z0-9_.,\x60\s]+(\s+(ASC|DESC))?)*$`)

// 允许的列名字符（字母、数字、下划线）
var columnNamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)

// 允许的表名字符（字母、数字、下划线）
var tableNamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]*$`)

// ValidateColumnName 验证列名是否安全
func (s *SQLSecurityUtils) ValidateColumnName(columnName string) bool {
	if columnName == "" {
		return false
	}

	// 检查长度
	if len(columnName) > 64 {
		return false
	}

	// 检查是否匹配安全的列名模式
	return columnNamePattern.MatchString(columnName)
}

// ValidateTableName 验证表名是否安全
func (s *SQLSecurityUtils) ValidateTableName(tableName string) bool {
	if tableName == "" {
		return false
	}

	// 检查长度
	if len(tableName) > 64 {
		return false
	}

	// 检查是否匹配安全的表名模式
	return tableNamePattern.MatchString(tableName)
}

// ValidateOrderBy 验证ORDER BY子句是否安全
func (s *SQLSecurityUtils) ValidateOrderBy(orderBy string) bool {
	if orderBy == "" {
		return true // 空的ORDER BY是安全的
	}

	// 检查长度
	if len(orderBy) > 200 {
		return false
	}

	// 转换为大写进行检查
	upperOrderBy := strings.ToUpper(orderBy)

	// 检查是否包含危险关键字
	for _, keyword := range dangerousKeywords {
		if strings.Contains(upperOrderBy, strings.ToUpper(keyword)) {
			// 允许ASC和DESC关键字
			if keyword != "desc" && keyword != "asc" {
				return false
			}
		}
	}

	// 检查是否匹配安全的ORDER BY模式
	return orderByPattern.MatchString(orderBy)
}

// SanitizeString 清理字符串，移除潜在的SQL注入字符
func (s *SQLSecurityUtils) SanitizeString(input string) string {
	if input == "" {
		return input
	}

	// 移除或转义危险字符
	input = strings.ReplaceAll(input, "'", "''")    // 转义单引号
	input = strings.ReplaceAll(input, "\"", "\\\"") // 转义双引号
	input = strings.ReplaceAll(input, "\\", "\\\\") // 转义反斜杠
	input = strings.ReplaceAll(input, "\x00", "")   // 移除NULL字符
	input = strings.ReplaceAll(input, "\n", " ")    // 替换换行符
	input = strings.ReplaceAll(input, "\r", " ")    // 替换回车符
	input = strings.ReplaceAll(input, "\t", " ")    // 替换制表符

	return strings.TrimSpace(input)
}

// ValidateOperator 验证操作符是否在允许的列表中
func (s *SQLSecurityUtils) ValidateOperator(operator string) bool {
	allowedOps := []string{"=", ">", "<", ">=", "<=", "like", "between", "in", "not in"}
	operator = strings.ToLower(strings.TrimSpace(operator))

	for _, allowedOp := range allowedOps {
		if operator == allowedOp {
			return true
		}
	}
	return false
}

// ValidateLimit 验证LIMIT值是否安全
func (s *SQLSecurityUtils) ValidateLimit(limit int) bool {
	return limit > 0 && limit <= 10000 // 限制最大查询数量
}

// ValidateOffset 验证OFFSET值是否安全
func (s *SQLSecurityUtils) ValidateOffset(offset int) bool {
	return offset >= 0 && offset <= 1000000 // 限制最大偏移量
}

// ContainsDangerousKeywords 检查输入是否包含危险的SQL关键字
func (s *SQLSecurityUtils) ContainsDangerousKeywords(input string) bool {
	upperInput := strings.ToUpper(input)

	for _, keyword := range dangerousKeywords {
		if strings.Contains(upperInput, strings.ToUpper(keyword)) {
			return true
		}
	}
	return false
}

// 创建全局实例
var SQLSecurity = &SQLSecurityUtils{}
