<template>
  <div :class="['editor-box', self_disabled ? 'editor-disabled' : '']">
    <Toolbar v-if="!hideToolBar" class="editor-toolbar" :editor="editorRef" :default-config="toolbarConfig" :mode="mode"/>
    <Editor
      v-model="valueHtml"
      class="editor-content"
      :style="{ height }"
      :mode="mode"
      :default-config="editorConfig"
      @on-created="handleCreated"
      @on-blur="handleBlur"
    />
  </div>
</template>

<script setup lang="ts" name="WangEditor">
import {nextTick, computed, inject, shallowRef, onBeforeUnmount, ref} from "vue";
import {IToolbarConfig, IEditorConfig} from "@wangeditor/editor";
import {Editor, Toolbar} from "@wangeditor/editor-for-vue";
import {UploadFiles} from "../../api/uploadFile.js";

import "@wangeditor/editor/dist/css/style.css";
import {ElNotification, formContextKey, formItemContextKey} from "element-plus";

// 富文本 DOM 元素
const editorRef = shallowRef();

// 实列化编辑器
const handleCreated = (editor: any) => {
  editorRef.value = editor;
};

// 接收父组件参数，并设置默认值
interface RichEditorProps {
  value: string; // 富文本值 ==> 必传
  toolbarConfig?: Partial<IToolbarConfig>; // 工具栏配置 ==> 非必传（默认为空）
  editorConfig?: Partial<IEditorConfig>; // 编辑器配置 ==> 非必传（默认为空）
  height?: string; // 富文本高度 ==> 非必传（默认为 500px）
  mode?: "default" | "simple"; // 富文本模式 ==> 非必传（默认为 default）
  hideToolBar?: boolean; // 是否隐藏工具栏 ==> 非必传（默认为false）
  disabled?: boolean; // 是否禁用编辑器 ==> 非必传（默认为false）
}

const props = withDefaults(defineProps<RichEditorProps>(), {
  toolbarConfig: () => {
    return {
      excludeKeys: []
    };
  },
  editorConfig: () => {
    return {
      placeholder: "请输入内容...",
      MENU_CONF: {}
    };
  },
  height: "500px",
  mode: "default",
  hideToolBar: false,
  disabled: false
});

// 获取 el-form 组件上下文
const formContext = inject(formContextKey, void 0);
// 获取 el-form-item 组件上下文
const formItemContext = inject(formItemContextKey, void 0);
// 判断是否禁用上传和删除
const self_disabled = computed(() => {
  return props.disabled || formContext?.disabled;
});

// 判断当前富文本编辑器是否禁用
if (self_disabled.value) nextTick(() => editorRef.value.disable());

// 富文本的内容监听，触发父组件改变，实现双向数据绑定
const emit = defineEmits<{
  "update:value": [value: string];
  "check-validate": [];
}>();
const valueHtml = computed({
  get() {
    return props.value;
  },
  set(val: string) {
    // 防止富文本内容为空时，校验失败
    if (editorRef.value.isEmpty()) val = "";
    emit("update:value", val);
  }
});

/**
 * @description 图片自定义上传
 * @param file 上传的文件
 * @param insertFn 上传成功后的回调函数（插入到富文本编辑器中）
 * */
type InsertFnTypeImg = (url: string, alt?: string, href?: string) => void;
props.editorConfig.MENU_CONF!["uploadImage"] = {
  async customUpload(file: File, insertFn: InsertFnTypeImg) {
    if (!uploadImgValidate(file)) return;
    try {
      const res = await UploadFiles(file);
      if (res?.res?.requestUrls?.length) {
        insertFn(res.res.requestUrls[0]);
      } else {

      }
      
    } catch (error) {
      console.log(error);
    }
  }
};

// 图片上传前判断

const imagesType = ref(["image/jpeg", "image/png", "image/gif", "image/svg+xml"])  // 文件类型
const imagesSize = ref(5) // 文件大小

const uploadImgValidate = (file: File): boolean => {
  console.log("rawFile", file)
  const imgSize = file.size < imagesSize.value * 1024 * 1024;
  const imgType = imagesType.value.includes(file.type);

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
        message: "上传图片大小不能超过" + 5 + "M！",
        type: "warning"
      });
    }, 0);
    return false
  }

  return true;
};

/**
 * @description 视频自定义上传
 * @param file 上传的文件
 * @param insertFn 上传成功后的回调函数（插入到富文本编辑器中）
 * */
type InsertFnTypeVideo = (url: string, poster?: string) => void;
props.editorConfig.MENU_CONF!["uploadVideo"] = {
  async customUpload(file: File, insertFn: InsertFnTypeVideo) {
    if (!uploadVideoValidate(file)) return;
    let formData = new FormData();
    formData.append("file", file);
    try {
      //const {data} = await UploadFiles(file);
      //insertFn(data.fileUrl);
      const res = await UploadFiles(file);
      console.log(res, 333)
      if (res?.res?.requestUrls?.length) {
        insertFn(res.res.requestUrls[0]);
      } else {

      }
    } catch (error) {
      console.log(error);
    }
  }
};

// 视频上传前判断
const videoType = ref(["video/mp4"])
const uploadVideoValidate = (file: File): boolean => {
  console.log(file);
  const mp4Type = videoType.value.includes(file.type);

  console.log("videoType", mp4Type)

  if (!mp4Type) {
    ElNotification.warning("请上传mp4格式的视频文件！");
    return false
  } else if (file.size > 10 * 1024 * 1024) {
    ElNotification.error('视频大小不能超过10M')
    return false
  }
  return true;
};

// 编辑框失去焦点时触发
const handleBlur = () => {
  formItemContext?.prop && formContext?.validateField([formItemContext.prop as string]);
  console.log(valueHtml.value)
  emit("update:value", valueHtml.value);
};

// 组件销毁时，也及时销毁编辑器
onBeforeUnmount(() => {
  if (!editorRef.value) return;
  editorRef.value.destroy();
});

defineExpose({
  editor: editorRef
});
</script>

<style scoped lang="scss">
@user "./index.scss";
</style>
