<template>
  <el-dialog v-model="dialogVisible" title="修改信息" width="500px" draggable>
    <el-form :model="form_user" label-width="100">
      <el-form-item label="头像">
        <upImage :avatar="form_user.avatar" @upImgSuccess="upSucess" @upDel="upDel"></upImage>
      </el-form-item>
      <el-form-item label="管理员昵称">
        <el-input v-model="form_user.nickname"/>
      </el-form-item>
      <el-form-item label="密码">
        <el-input v-model="form_user.password" type="password"/>
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="closePopup">取消</el-button>
        <el-button type="primary" @click="modifyUserConfirm">确认</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from "vue";
import { useUserStore } from "@/stores/modules/user";
import upImage from "@/components/Upload/Img.vue"
import { ElMessage } from "element-plus";
import { userSetUser } from "@/api/modules/user.js";

const userStore = useUserStore();
const dialogVisible = ref(false);

const form_user = ref({
  avatar: "",//  string  头像
  nickname: "",  //string  昵称 
  password: "",  //string:"" 密码
})

onMounted(() => {
  form_user.value = userStore.userInfo
  form_user.value.avatar = userStore.userInfo.avatar
})

watch(() => userStore.userInfo, (res) => {
  form_user.value = res
  form_user.value.avatar = userStore.userInfo.avatar
}, {deep: true})

// 头像上传完成回调
const upSucess = (data) => {
  form_user.value.avatar = data
}
// 头像删除回调
const upDel = () => {
  form_user.value.avatar = ""
}

// 关闭弹窗
const closePopup = () => { 
  dialogVisible.value = false
}

// 修改资料提交
const modifyUserConfirm = () => {
  const { avatar, nickname, password } = form_user.value
  const params = {
    avatar, nickname, password
  }
  userSetUser(params).then(res => {
    if (res.code == 200) {
      let obj = {
        avatar: res.data.avatar,
        nickname: res.data.nickname,
      }
      userStore.setUserInfo(obj)
      ElMessage.success(res.msg)
      closePopup()
    } else {
      ElMessage.error(res.msg)
    }
  })
}

// 打开弹窗
const openDialog = () => {
  dialogVisible.value = true;
};

// 暴露方法
defineExpose({ openDialog });
</script>
