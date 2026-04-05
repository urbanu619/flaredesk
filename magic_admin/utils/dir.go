package utils

import (
	"github.com/pkg/errors"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: PathExists
//@description: 文件目录是否存在
//@param: path string
//@return: bool, error

func PathExists(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err == nil {
		if fi.IsDir() {
			return true, nil
		}
		return false, errors.New("存在同名文件")
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateDir
//@description: 批量创建文件夹
//@param: dirs ...string
//@return: err error

func CreateDir(dirs ...string) (err error) {
	for _, v := range dirs {
		exist, err := PathExists(v)
		if err != nil {
			return err
		}
		if !exist {
			if err := os.MkdirAll(v, os.ModePerm); err != nil {
				return err
			}
		}
	}
	return err
}

//@author: [songzhibin97](https://github.com/songzhibin97)
//@function: FileMove
//@description: 文件移动供外部调用
//@param: src string, dst string(src: 源位置,绝对路径or相对路径, dst: 目标位置,绝对路径or相对路径,必须为文件夹)
//@return: err error

func FileMove(src string, dst string) (err error) {
	if dst == "" {
		return nil
	}
	src, err = filepath.Abs(src)
	if err != nil {
		return err
	}
	dst, err = filepath.Abs(dst)
	if err != nil {
		return err
	}
	revoke := false
	dir := filepath.Dir(dst)
Redirect:
	_, err = os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 0o755)
		if err != nil {
			return err
		}
		if !revoke {
			revoke = true
			goto Redirect
		}
	}
	return os.Rename(src, dst)
}

func DeLFile(filePath string) error {
	return os.RemoveAll(filePath)
}

//@author: [songzhibin97](https://github.com/songzhibin97)
//@function: TrimSpace
//@description: 去除结构体空格
//@param: target interface (target: 目标结构体,传入必须是指针类型)
//@return: null

func TrimSpace(target interface{}) {
	t := reflect.TypeOf(target)
	if t.Kind() != reflect.Ptr {
		return
	}
	t = t.Elem()
	v := reflect.ValueOf(target).Elem()
	for i := 0; i < t.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.String:
			v.Field(i).SetString(strings.TrimSpace(v.Field(i).String()))
		}
	}
}

// FileExist 判断文件是否存在
func FileExist(path string) bool {
	fi, err := os.Lstat(path)
	if err == nil {
		return !fi.IsDir()
	}
	return !os.IsNotExist(err)
}

// AbsPath 绝对路径生成
func AbsPath(args ...string) string {
	absPath, err := os.Getwd()
	if err != nil {
		return ""
	}
	savePath := ToPath(absPath, ToPath(args...))
	err = CheckPathIsExist(savePath)
	if err != nil {
		return ""
	}
	return savePath
}

// ToPath 转化成路径
func ToPath(args ...string) string {
	if len(args) == 0 {
		return ""
	}
	strList := make([]string, 0)
	for i, p := range args {
		if strings.TrimSpace(p) == "" {
			continue
		}
		if i == 0 {
			strList = append(strList, strings.TrimRight(strings.TrimSpace(p), "/"))
			continue
		}

		strList = append(strList, strings.Trim(strings.TrimSpace(p), "/"))
	}
	return strings.Join(strList, "/")
}

// CheckPathIsExist 检查路径并创建
func CheckPathIsExist(filePath string) error {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(filePath, os.ModePerm) //生成多级目录
			if err != nil {
				return err
			}
			return err
		}
	}
	return err
}

// GetFileContentType 获取文件类型
func GetFileContentType(out *os.File) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}
	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)
	return contentType, nil
}
