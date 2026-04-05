package config

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
	"go_server/utils"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"
)

type File struct {
	OssType    string                 `mapstructure:"oss-type" json:"oss-type" yaml:"oss-type"`          // 文件存储类型
	Path       string                 `mapstructure:"path" json:"path" yaml:"path"`                      // 本地文件访问路径
	ProxyPath  string                 `mapstructure:"proxy-path" json:"proxy-path" yaml:"proxy-path"`    // 系统代理路径
	StorePath  string                 `mapstructure:"store-path" json:"store-path" yaml:"store-path"`    // 本地文件存储路径
	OriginConf map[string]interface{} `mapstructure:"origin-conf" json:"origin-conf" yaml:"origin-conf"` // 远程文件存储配置
}

// 导出csv表

func ExportCsv(headers []string, data [][]string, filename string) (string, error) {
	if filename == "" || len(data) == 0 {
		return "", fmt.Errorf("数据错误")
	}
	// exportCSV 导出CSV文件
	// 生成文件名，包含当前时间戳
	// 创建文件
	file, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()
	file.WriteString("\xEF\xBB\xBF")
	// 创建CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()
	// 写入表头
	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("写入表头失败: %v", err)
	}
	// 写入数据
	for _, record := range data {
		if err := writer.Write(record); err != nil {
			return "", fmt.Errorf("写入数据失败: %v", err)
		}
	}
	return filename, nil
}

func ExportExcelFile[T any](headers []string, data []T, savePath string) error {
	if savePath == "" || len(data) == 0 {
		return fmt.Errorf("无目标数据")
	}
	f := excelize.NewFile()
	sheetName := "Sheet1"
	if err := f.SetSheetName("Sheet1", sheetName); err != nil {
		return err

	}

	// 写入表头
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		if err := f.SetCellValue(sheetName, cell, header); err != nil {
			return err

		}
	}

	// 使用反射获取结构体字段值
	for rowIndex, row := range data {
		rowValue := reflect.ValueOf(row)
		for colIndex := 0; colIndex < rowValue.NumField(); colIndex++ {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
			field := rowValue.Field(colIndex)
			if field.CanInterface() {
				if err := f.SetCellValue(sheetName, cell, field.Interface()); err != nil {
					return err

				}
			}
		}
	}
	// 保存 Excel 文件
	if err := f.SaveAs(savePath); err != nil {
		return err
	}
	return nil
}

func (a *File) UploadFile(file *multipart.FileHeader) (viewPath, savePath string, err error) {
	switch a.OssType {
	case OssTypeLocal:
		viewPath, savePath, err = a.uploadFile(file)
		return
	case OssTypeAws:
		engine := new(AwsS3)
		byteData, _ := json.Marshal(a.OriginConf)
		err = json.Unmarshal(byteData, &engine)
		if err != nil {
			return
		}
		viewPath, savePath, err = engine.UploadFile(file)
	case OssTypeAli:
		engine := new(AliOSS)
		byteData, _ := json.Marshal(a.OriginConf)
		err = json.Unmarshal(byteData, &engine)
		if err != nil {
			return
		}
		viewPath, savePath, err = engine.UploadFile(file)
		return
	default:
		err = fmt.Errorf("oss type not support")
		break
	}
	return
}

func (a *File) DeleteFile(key string) (err error) {
	switch a.OssType {
	case OssTypeLocal:
		err = a.deleteFile(key)
		return
	case OssTypeAli:
		engine := new(AwsS3)
		byteData, _ := json.Marshal(a.OriginConf)
		err = json.Unmarshal(byteData, &engine)
		if err != nil {
			return
		}
		err = engine.DeleteFile(key)
		return
	case OssTypeAws:
		engine := new(AliOSS)
		byteData, _ := json.Marshal(a.OriginConf)
		err = json.Unmarshal(byteData, &engine)
		if err != nil {
			return
		}
		err = engine.DeleteFile(key)
		return
	default:
		err = fmt.Errorf("oss type not support")
		break
	}
	return
}

func (a *File) uploadFile(file *multipart.FileHeader) (viewPath, filename string, err error) {
	// 读取文件后缀
	ext := filepath.Ext(file.Filename)
	// 读取文件名并md5
	name := strings.TrimSuffix(file.Filename, ext)
	name = utils.Md5ByteEncode([]byte(name))
	// 拼接新文件名
	filename = name + "_" + time.Now().Format("20060102150405") + ext
	// 尝试创建此路径
	mkdirErr := os.MkdirAll(a.StorePath, os.ModePerm)
	if mkdirErr != nil {
		return "", "", errors.New("function os.MkdirAll() failed, err:" + mkdirErr.Error())
	}
	// 拼接路径和文件名
	p := utils.ToPath(a.StorePath, filename) // 存储文件
	viewPath = utils.ToPath(a.Path, a.ProxyPath, filename)

	f, openError := file.Open() // 读取文件
	if openError != nil {
		return "", "", errors.New("function file.Open() failed, err:" + openError.Error())
	}
	defer func() {
		_ = f.Close()
	}() // 创建文件 defer 关闭

	out, createErr := os.Create(p)
	if createErr != nil {
		return "", "", errors.New("function os.Create() failed, err:" + createErr.Error())
	}
	defer func() {
		_ = out.Close()
	}()
	// 创建文件 defer 关闭
	_, copyErr := io.Copy(out, f) // 传输（拷贝）文件
	if copyErr != nil {
		return "", "", errors.New("function io.Copy() failed, err:" + copyErr.Error())
	}
	return
}

var mu sync.Mutex

func (a *File) deleteFile(key string) (err error) {
	// 检查 key 是否为空
	if key == "" {
		return errors.New("key不能为空")
	}
	// 验证 key 是否包含非法字符或尝试访问存储路径之外的文件
	if strings.Contains(key, "..") || strings.ContainsAny(key, `\/:*?"<>|`) {
		return errors.New("非法的key")
	}
	p := filepath.Join(a.StorePath, key)
	// 检查文件是否存在
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return errors.New("文件不存在")
	}
	// 使用文件锁防止并发删除
	mu.Lock()
	defer mu.Unlock()
	err = os.Remove(p)
	if err != nil {
		return errors.New("文件删除失败: " + err.Error())
	}
	return
}

const (
	FileJobRunning = 1 // 进行中
	FileJobFinish  = 2 // 已完成
	FileJobFail    = 3 // 失败

	OssTypeLocal = "local"
	OssTypeAli   = "ali"
	OssTypeAws   = "aws"
)

type OriginFile interface {
	UploadFile(*multipart.FileHeader) (string, string, error)
	DeleteFile(key string) error
}

type AliOSS struct {
	Endpoint        string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	AccessKeyId     string `mapstructure:"access-key-id" json:"access-key-id" yaml:"access-key-id"`
	AccessKeySecret string `mapstructure:"access-key-secret" json:"access-key-secret" yaml:"access-key-secret"`
	BucketName      string `mapstructure:"bucket-name" json:"bucket-name" yaml:"bucket-name"`
	BucketUrl       string `mapstructure:"bucket-url" json:"bucket-url" yaml:"bucket-url"`
	BasePath        string `mapstructure:"base-path" json:"base-path" yaml:"base-path"`
}

func (a *AliOSS) UploadFile(file *multipart.FileHeader) (string, string, error) {
	bucket, err := a.aliBucket()
	if err != nil {
		return "", "", errors.New("aliBucket Failed, err:" + err.Error())
	}
	// 读取本地文件。
	f, openError := file.Open()
	if openError != nil {
		return "", "", errors.New("function file.Open() Failed, err:" + openError.Error())
	}
	defer func() {
		_ = f.Close()
	}() // 创建文件 defer 关闭
	// 上传阿里云路径 文件名格式 自己可以改 建议保证唯一性
	yunFileTmpPath := a.BasePath + "/" + "uploads" + "/" + time.Now().Format("2006-01-02") + "/" + file.Filename
	// 上传文件流。
	err = bucket.PutObject(yunFileTmpPath, f)
	if err != nil {
		return "", "", errors.New("PutObject Failed, err:" + err.Error())
	}
	return a.BucketUrl + "/" + yunFileTmpPath, yunFileTmpPath, nil
}

func (a *AliOSS) DeleteFile(key string) error {
	bucket, err := a.aliBucket()
	if err != nil {
		return errors.New("aliBucket err:" + err.Error())
	}
	// 删除单个文件。objectName表示删除OSS文件时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	// 如需删除文件夹，请将objectName设置为对应的文件夹名称。如果文件夹非空，则需要将文件夹下的所有object删除后才能删除该文件夹。
	err = bucket.DeleteObject(key)
	if err != nil {
		return errors.New("bucket.DeleteObject, err:" + err.Error())
	}
	return nil
}

func (a *AliOSS) aliBucket() (*oss.Bucket, error) {
	var err error
	var client *oss.Client
	var buket *oss.Bucket
	// 创建OSSClient实例。
	client, err = oss.New(a.Endpoint, a.AccessKeyId, a.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	// 获取存储空间。
	buket, err = client.Bucket(a.BucketName)
	if err != nil {
		return nil, err
	}
	return buket, err
}

type AwsS3 struct {
	Bucket           string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`
	Region           string `mapstructure:"region" json:"region" yaml:"region"`
	Endpoint         string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	SecretID         string `mapstructure:"secret-id" json:"secret-id" yaml:"secret-id"`
	SecretKey        string `mapstructure:"secret-key" json:"secret-key" yaml:"secret-key"`
	BaseURL          string `mapstructure:"base-url" json:"base-url" yaml:"base-url"`
	PathPrefix       string `mapstructure:"path-prefix" json:"path-prefix" yaml:"path-prefix"`
	S3ForcePathStyle bool   `mapstructure:"s3-force-path-style" json:"s3-force-path-style" yaml:"s3-force-path-style"`
	DisableSSL       bool   `mapstructure:"disable-ssl" json:"disable-ssl" yaml:"disable-ssl"`
	_session         *session.Session
}

func (a *AwsS3) UploadFile(file *multipart.FileHeader) (string, string, error) {
	_session := a.newSession()
	uploader := s3manager.NewUploader(_session)
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename)
	filename := a.PathPrefix + "/" + fileKey
	f, openError := file.Open()
	if openError != nil {
		return "", "", errors.New("function file.Open() failed, err:" + openError.Error())
	}
	defer func() {
		_ = f.Close()
	}() // 创建文件 defer 关闭
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(a.Bucket),
		Key:    aws.String(filename),
		Body:   f,
	})
	if err != nil {
		return "", "", err
	}
	return a.BaseURL + "/" + filename, fileKey, nil
}

func (a *AwsS3) DeleteFile(key string) error {
	_session := a.newSession()
	svc := s3.New(_session)
	filename := a.PathPrefix + "/" + key
	bucket := a.Bucket
	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return errors.New("function svc.DeleteObject() failed, err:" + err.Error())
	}
	_ = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	return nil
}

// newSession Create S3 session

func (a *AwsS3) newSession() *session.Session {
	sess, _ := session.NewSession(&aws.Config{
		Region:           aws.String(a.Region),
		Endpoint:         aws.String(a.Endpoint), //minio在这里设置地址,可以兼容
		S3ForcePathStyle: aws.Bool(a.S3ForcePathStyle),
		DisableSSL:       aws.Bool(a.DisableSSL),
		Credentials: credentials.NewStaticCredentials(
			a.SecretID,
			a.SecretKey,
			"",
		),
	})
	return sess
}
