<template>
  <div class="upload-box">
    <el-upload
      action="#"
      :class="['upload']"
      :multiple="false"
      :show-file-list="false"
      :http-request="handleHttpUpload"
      :before-upload="beforeUpload"
      :drag="true"
    >
      <template v-if="imageUrl">
        <img :src="imageUrl" class="upload-image"/>
        <div class="upload-handle" @click.stop>
          <div class="handle-icon" @click="deleteImg">
            <el-icon>
              <Delete/>
            </el-icon>
            <span>删除</span>
          </div>
        </div>
      </template>
      <template v-else>
        <div class="upload-empty">
          <slot name="empty">
            <el-icon>
              <Plus/>
            </el-icon>
            <!-- <span>请上传图片</span> -->
          </slot>
        </div>
      </template>
    </el-upload>
    <div class="el-upload__tip">
      <slot name="tip"></slot>
    </div>
    <el-image-viewer v-if="imgViewVisible" :url-list="[imageUrl]" @close="imgViewVisible = false"/>
  </div>
</template>

<script setup name="UploadImg">
import {ref, computed, inject, onMounted} from "vue";
import OSS from 'ali-oss';

import {uploadImg} from "@/api/modules/upload.js";
import {ElNotification, formContextKey, formItemContextKey} from "element-plus";
import {uploadOssAuth} from "../../api/modules/upload.js";
import {UploadFiles} from "@/api/uploadFile.js";
// import {UploadProps, UploadRequestOptions} from "element-plus";

// 查看图片
const imgViewVisible = ref(false);

const imagesType = ref(["image/jpeg", "image/png", "image/gif", "image/svg+xml"])  // 文件类型
const imagesSize = ref(2) // 文件大小
// 接受父组件参数
const props = defineProps(["avatar"]);
const emit = defineEmits(["upImgSuccess", "upImgError", "upDel"]);

const imageUrl = ref("")

onMounted(() => {
  if (props.avatar) {
    imageUrl.value = props.avatar
  } else {
    imageUrl.value = ''
  }
})

const getUploadOssAuth = () => {
  return new Promise((reslove, reject) => {
    uploadOssAuth().then(res => {
      if (res.code == 200) {
        reslove(res.data)
      } else {
        ElNotification.error(res.msg)
        reject({})

      }
    })

  })


}
const handleHttpUpload = async (file) => {
  let res = await UploadFiles(file.file)
  if (res) {
    imageUrl.value = res?.res?.requestUrls || []
    uploadSuccess( imageUrl.value[0] )
  } else {

  }
  // let resp = await getUploadOssAuth()
  // console.log("options", file.file)
  // let {bucket, bucketUrl, credentials, region, filePath} = resp
  // let {AccessKeyId, AccessKeySecret, Expiration, SecurityToken} = credentials
  //
  // try {
  //   // 从后端获取STS临时凭证（这里仅为示例，实际应用中应从后端API获取）
  //   const ossClient = new OSS({
  //     region, // 设置OSS的区域，例如'oss-cn-hangzhou'
  //     accessKeyId: AccessKeyId,
  //     accessKeySecret: AccessKeySecret,
  //     stsToken: SecurityToken, // 使用STS的临时凭证进行授权上传
  //     bucket: bucket, // 设置你的bucket名称
  //
  //     refreshSTSToken: async () => {
  //       // Token 过期时自动调用
  //       const newToken = await this.getUploadOssAuth();
  //       if (newToken) {
  //         return {
  //           accessKeyId: newToken.credentials.accessKeyId,
  //           accessKeySecret: newToken.credentials.accessKeySecret,
  //           stsToken: newToken.credentials.stsToken
  //         };
  //       }
  //     },
  //     refreshSTSTokenInterval: 300000 // 5分钟检查一次(单位毫秒)
  //   });
  //   // 上传文件到OSS，这里使用了默认的上传方法，你也可以使用分片上传等方法根据需求调整。
  //   const result = await ossClient.put(filePath + "/" + file.file.name, file.file); // 上传文件到OSS，文件名与本地文件名相同。你也可以自定义文件名。
  //   console.log(result)
  //   imageUrl.value = result
  //   uploadSuccess()
  //   console.log('Upload success:', result); // 上传成功后的回调处理。例如显示上传结果等。
  // } catch (error) {
  //   uploadError()
  //   console.error('Upload error:', error); // 上传失败的处理逻辑。例如显示错误信息等。
  // }


  //
  //
  //
  //
  // const res = await uploadImg(formData);
  // if (res.code == 200) {
  //   imageUrl.value = res.data.path
  // } else {
  //   ElNotification.error(res.msg)
  // }

};


//删除图片
const deleteImg = () => {
  imageUrl.value = ''
  emit("upDel")
};

/**
 * @description 文件上传之前判断
 * @param rawFile 选择的文件
 * */
const beforeUpload = (rawFile) => {
  console.log("rawFile", rawFile)
  const imgSize = rawFile.size < imagesSize.value * 1024 * 1024;
  const imgType = imagesType.value.includes(rawFile.type);

  console.log(imgSize, imgType)

  if (!imgType) {
    ElNotification({
      title: "温馨提示",
      message: "上传图片不符合所需的格式！",
      type: "warning"
    });
    return false
  } else if (!imgSize) {
    setTimeout(() => {
      ElNotification({
        title: "温馨提示",
        message: "上传图片大小不能超过" + fileSize.value + "M！",
        type: "warning"
      });
    }, 0);
    return false
  }


  return true;
};

/**
 * @description 图片上传成功
 * */
const uploadSuccess = (img) => {
  ElNotification({
    title: "温馨提示",
    message: "图片上传成功！",
    type: "success"
  });
  emit("upImgSuccess", img)
};

/**
 * @description 图片上传错误
 * */
const uploadError = () => {
  ElNotification({
    title: "温馨提示",
    message: "图片上传失败，请您重新上传！",
    type: "error"
  });
};
</script>

<style scoped lang="scss">
.is-error {
  .upload {
    width: 100px;
    height: 100px;

    :deep(.el-upload),
    :deep(.el-upload-dragger) {
      border: 1px dashed var(--el-color-danger) !important;

      &:hover {
        border-color: var(--el-color-primary) !important;
      }
    }
  }
}

:deep(.disabled) {
  .el-upload,
  .el-upload-dragger {
    cursor: not-allowed !important;
    background: var(--el-disabled-bg-color);
    border: 1px dashed var(--el-border-color-darker) !important;

    &:hover {
      border: 1px dashed var(--el-border-color-darker) !important;
    }
  }
}

.upload-box {
  width: 100px;
  height: 100px;

  .no-border {
    :deep(.el-upload) {
      border: none !important;
    }
  }

  :deep(.upload) {
    width: 100% !important;
    height: 100% !important;

    .el-upload {
      width: 100% !important;
      height: 100% !important;
      position: relative;
      display: flex;
      align-items: center;
      justify-content: center;
      overflow: hidden;
      border: 1px dashed var(--el-border-color-darker);
      border-radius: v-bind(borderRadius);
      transition: var(--el-transition-duration-fast);

      &:hover {
        border-color: var(--el-color-primary);

        .upload-handle {
          opacity: 1;
        }
      }

      .el-upload-dragger {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 100%;
        height: 100%;
        padding: 0;
        overflow: hidden;
        background-color: transparent;
        border: 1px dashed var(--el-border-color-darker);
        border-radius: v-bind(borderRadius);

        &:hover {
          border: 1px dashed var(--el-color-primary);
        }
      }

      .el-upload-dragger.is-dragover {
        background-color: var(--el-color-primary-light-9);
        border: 2px dashed var(--el-color-primary) !important;
      }

      .upload-image {
        width: 100%;
        height: 100%;
        object-fit: contain;
      }

      .upload-empty {
        position: relative;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        font-size: 12px;
        line-height: 30px;
        color: var(--el-color-info);

        .el-icon {
          font-size: 28px;
          color: var(--el-text-color-secondary);
        }
      }

      .upload-handle {
        position: absolute;
        top: 0;
        right: 0;
        box-sizing: border-box;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 100%;
        height: 100%;
        cursor: pointer;
        background: rgb(0 0 0 / 60%);
        opacity: 0;
        transition: var(--el-transition-duration-fast);

        .handle-icon {
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          padding: 0 6%;
          color: aliceblue;

          .el-icon {
            margin-bottom: 40%;
            font-size: 130%;
            line-height: 130%;
          }

          span {
            font-size: 85%;
            line-height: 85%;
          }
        }
      }
    }
  }

  .el-upload__tip {
    line-height: 18px;
    text-align: center;
  }
}
</style>
