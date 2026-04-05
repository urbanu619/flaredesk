import { uploadOssAuth } from "./modules/upload.js";
import { ElNotification } from "element-plus";
import OSS from "ali-oss";
// 上传文件

// 获取授权
const getUploadOssAuth = () => {
  return new Promise((reslove, reject) => {
    uploadOssAuth().then(res => {
      if (res.code == 200) {
        reslove(res.data);
      } else {
        ElNotification.error(res.msg);
        reject({});
      }
    });
  });
};

export const UploadFiles = async file => {
  console.log("file", file);

  let { code, data, msg } = await uploadOssAuth();
  if (code !== 200) {
    ElNotification.error(msg);
    return false;
  }
  //

  let { bucket, bucketUrl, credentials, region, filePath } = data;
  let { AccessKeyId, AccessKeySecret, Expiration, SecurityToken } = credentials;

  try {
    // 从后端获取STS临时凭证（这里仅为示例，实际应用中应从后端API获取）
    console.log(region, AccessKeyId, AccessKeySecret, SecurityToken, bucket);
    const ossClient = new OSS({
      region: region, // 设置OSS的区域，例如'oss-cn-hangzhou' "oss-" +
      accessKeyId: AccessKeyId,
      accessKeySecret: AccessKeySecret,
      stsToken: SecurityToken, // 使用STS的临时凭证进行授权上传
      bucket: bucket, // 设置你的bucket名称
      secure: true, // 强制使用HTTPS
      // 添加跨域相关配置
      crossDomain: true,
      headers: {
        "Access-Control-Allow-Origin": "*",
        "Access-Control-Allow-Methods": "PUT, GET, POST, DELETE, OPTIONS",
        "Access-Control-Allow-Headers": "*"
      }
      // refreshSTSToken: async () => {
      //   // Token 过期时自动调用
      //   const newToken = await getUploadOssAuth();
      //   if (newToken) {
      //     return {
      //       accessKeyId: newToken.credentials.accessKeyId,
      //       accessKeySecret: newToken.credentials.accessKeySecret,
      //       stsToken: newToken.credentials.stsToken
      //     };
      //   }
      // },
      // refreshSTSTokenInterval: 300000 // 5分钟检查一次(单位毫秒)
    });
    // 上传文件到OSS，这里使用了默认的上传方法，你也可以使用分片上传等方法根据需求调整。
    // const result = await ossClient.put(filePath + "/" + file.name, file); // 上传文件到OSS，文件名与本地文件名相同。你也可以自定义文件名。
    // console.log("Upload success:", result); // 上传成功后的回调处理。例如显示上传结果等。

    // 推荐使用分片上传
    try {
      // const result = await ossClient.multipartUpload(
      //   `${filePath}/${file.name}`,
      //   file,
      //   { partSize: 1024 * 1024, retryCount: 3 }
      // );
      const result = await ossClient.put(filePath + "/" + file.name, file); // 上传文件到OSS，文件名与本地文件名相同。你也可以自定义文件名。
      console.log("Upload success:", result); // 上传成功后的回调处理。例如显示上传结果等。

      return result;
    } catch (error) {
      console.error("Upload error:", error); // 上传失败的处理逻辑。例如显示错误信息等。
      ElNotification.error("上传失败");
    } 
  } catch (error) {
    console.error("Upload error:", error); // 上传失败的处理逻辑。例如显示错误信息等。
    ElNotification.error("上传失败");
  }
};
