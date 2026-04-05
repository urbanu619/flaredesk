package config

import (
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	stsClient "github.com/alibabacloud-go/sts-20150401/v2/client"
	stsService "github.com/alibabacloud-go/tea-utils/v2/service"
)

type OssSts struct {
	AccessKeyId        string `mapstructure:"ACCESS_KEY_ID" json:"ACCESS_KEY_ID" yaml:"ACCESS_KEY_ID"`
	AccessKeySecret    string `mapstructure:"ACCESS_KEY_SECRET" json:"ACCESS_KEY_SECRET" yaml:"ACCESS_KEY_SECRET"`
	BucketName         string `mapstructure:"BUCKET_NAME" json:"BUCKET_NAME" yaml:"BUCKET_NAME"`
	BucketUrl          string `mapstructure:"BUCKET_URL" json:"BUCKET_URL" yaml:"BUCKET_URL"`
	BasePath           string `mapstructure:"BASE_PATH" json:"BASE_PATH" yaml:"BASE_PATH"`
	StsEndpoint        string `mapstructure:"STS_ENDPOINT" json:"STS_ENDPOINT" yaml:"STS_ENDPOINT"`
	StsDurationSeconds int64  `mapstructure:"STS_DURATION_SECONDS" json:"STS_DURATION_SECONDS" yaml:"STS_DURATION_SECONDS"` // 仅支持15分钟:900 / 60分钟:3600
	StsRoleSessionName string `mapstructure:"STS_ROLE_SESSION_NAME" json:"STS_ROLE_SESSION_NAME" yaml:"STS_ROLE_SESSION_NAME"`
	StsRoleArn         string `mapstructure:"STS_ROLE_ARN" json:"STS_ROLE_ARN" yaml:"STS_ROLE_ARN"`
	StsRegion          string `mapstructure:"STS_REGION" json:"STS_REGION" yaml:"STS_REGION"`
}

var stsAliEngine *StsAliOss

func StsAliEngin() *StsAliOss {
	if stsAliEngine == nil {
		configSts := &openapi.Config{
			AccessKeyId:     &EnvConf().StsOss.AccessKeyId,
			AccessKeySecret: &EnvConf().StsOss.AccessKeySecret,
			Endpoint:        &EnvConf().StsOss.StsEndpoint,
		}
		// docs https://api.alibabacloud.com/product/Sts
		stsCli, err := stsClient.NewClient(configSts)
		if err != nil {
			panic(err)
		}
		stsAliEngine = &StsAliOss{
			stsCli,
			&EnvConf().StsOss.StsRoleArn,
			&EnvConf().StsOss.StsDurationSeconds,
			EnvConf().StsOss.BucketName,
			EnvConf().StsOss.BucketUrl,
			EnvConf().StsOss.BasePath}
	}
	return stsAliEngine
}

type StsAliOss struct {
	stsClient       *stsClient.Client
	roleArn         *string
	durationSeconds *int64
	bucketName      string
	bucketUrl       string // 访问地址
	basePath        string //上传路径
}

// 阿里STS授权上传

func (s *StsAliOss) AliSTS(roleName string) (*stsClient.AssumeRoleResponse, error) {
	assumeRoleRequest := &stsClient.AssumeRoleRequest{
		DurationSeconds: s.durationSeconds,
		RoleArn:         s.roleArn,
		RoleSessionName: &roleName,
	}
	runtime := &stsService.RuntimeOptions{}
	return s.stsClient.AssumeRoleWithOptions(assumeRoleRequest, runtime)
}
