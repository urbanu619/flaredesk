package system

import (
	"github.com/gin-gonic/gin"
	"go_server/base/config"
	"go_server/model/common/response"
	"go_server/service/base"
	"go_server/utils"
)

type FileService struct {
	base.SysCommonService
}

func (s *FileService) UploadFile(c *gin.Context) {
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		response.Resp(c, "接收文件失败")
		return
	}
	path, file, err := config.EnvConf().File.UploadFile(header) // 文件上传后拿到文件路径
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c, map[string]interface{}{
		"path": path,
		"file": file,
	})
}

func (s *FileService) DeleteFile(c *gin.Context) {
	filename, ok := c.GetQuery("filename")
	if !ok {
		response.Resp(c, "接收文件失败")
		return
	}
	err := config.EnvConf().File.DeleteFile(filename) // 文件上传后拿到文件路径
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	response.Resp(c)
}

func (s *FileService) OssAuth(c *gin.Context) {
	resp, err := config.StsAliEngin().AliSTS("fomo-sts")
	if err != nil {
		response.Resp(c, err.Error())
		return
	}
	savePath := utils.ToPath(config.EnvConf().StsOss.BasePath)
	response.Resp(c, map[string]any{
		"region":      config.EnvConf().StsOss.StsRegion,
		"bucket":      config.EnvConf().StsOss.BucketName,
		"bucketUrl":   config.EnvConf().StsOss.BucketUrl,
		"credentials": resp.Body.Credentials,
		"filePath":    savePath,
	})

}
