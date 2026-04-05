package core

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

// 获取文件名

func GetCsvFileName(dir string) string {
	filename := fmt.Sprintf("%s_%d.xlsx", dir, time.Now().Unix())
	return filename
}

// 获取csv文件的路径

func GetCsvPath(filename string, dir string) string {
	execPath, err := os.Executable()
	if err != nil {
		return ""
	}
	// 获取可执行文件目录并构造目标路径
	baseDir := filepath.Dir(execPath)
	excelPath := filepath.Join(baseDir, "uploads", dir, filename)

	// 创建目录（如果不存在）
	if err := os.MkdirAll(filepath.Dir(excelPath), os.ModePerm); err != nil {
		return ""
	}

	return excelPath
}

// 导出表格

func ExportExcel[T any](headers []string, data []T, savePath string) (string, error) {
	if savePath == "" {
		return "", fmt.Errorf("获取文件路径失败")
	}
	f := excelize.NewFile()
	var err error
	sheetName := "Sheet1"
	err = f.SetSheetName("Sheet1", sheetName)
	if err != nil {
		return "", err
	}
	// 写入表头
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		err = f.SetCellValue(sheetName, cell, header)
		if err != nil {
			return "", err
		}
	}
	// 使用反射获取结构体字段值
	for rowIndex, row := range data {
		rowValue := reflect.ValueOf(row)
		for colIndex := 0; colIndex < rowValue.NumField(); colIndex++ {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
			field := rowValue.Field(colIndex)
			if field.CanInterface() {
				err = f.SetCellValue(sheetName, cell, field.Interface())
				if err != nil {
					return "", err
				}
			}
		}
	}
	// 保存 Excel 文件
	if err := f.SaveAs(savePath); err != nil {
		return "", fmt.Errorf("保存文件失败: %v", err)
	}
	return savePath, nil
}
