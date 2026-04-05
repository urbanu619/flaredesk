<template>
  <el-dropdown trigger="click">
    <div class="avatar">
      <div class="iconView" v-if="userStore.userInfo.avatar">
        <img :src="userStore.userInfo.avatar" alt="avatar"/>
      </div>

      <div class="iconView" v-else>
        <img src="@/assets/images/avatar.gif" alt="avatar"/>
      </div>
    </div>

    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item @click="openDialog('passwordRef')">
          <el-icon>
            <Edit/>
          </el-icon>
          修改信息
        </el-dropdown-item>
        <el-dropdown-item @click="openDialog('infoRef')">
          <el-icon>
            <Setting />
          </el-icon>
          设置谷歌密钥
        </el-dropdown-item>
        <el-dropdown-item @click="openDialog('delAuthRef')">
          <el-icon>
            <Delete />
          </el-icon>
          取消谷歌密钥
        </el-dropdown-item>
        <el-dropdown-item divided @click="logout">
          <el-icon>
            <SwitchButton/>
          </el-icon>
          {{ $t("header.logout") }}
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
  <!-- infoDialog -->
  <InfoDialog ref="infoRef"></InfoDialog>
  <!-- passwordDialog -->
  <PasswordDialog ref="passwordRef"></PasswordDialog>
</template>

<script setup lang="ts">
import {ref} from "vue";
import {LOGIN_URL} from "@/config";
import {useRouter} from "vue-router";
import {logoutApi} from "@/api/modules/login";
import {useUserStore} from "@/stores/modules/user";
import {useAuthStore} from "@/stores/modules/auth"
import {ElMessageBox, ElMessage} from "element-plus";
import InfoDialog from "./InfoDialog.vue";
import PasswordDialog from "./PasswordDialog.vue";


const router = useRouter();
const userStore = useUserStore();
const useAuth = useAuthStore();

// 退出登录
const logout = () => {
  ElMessageBox.confirm("您是否确认退出登录?", "温馨提示", {
    confirmButtonText: "确定",
    cancelButtonText: "取消",
    type: "warning"
  }).then(async () => {
    // 1.执行退出登录接口
    // await logoutApi();

    // 2.清除 Token
    userStore.setToken("");
    userStore.setUserInfo({avatar: "", nickname: ""});
    useAuth.clearMenusList()
    localStorage.clear()
    // 3.重定向到登陆页
    router.replace(LOGIN_URL);
    ElMessage.success("退出登录成功！");
  });
};

// 打开修改密码和个人信息弹窗
const infoRef = ref<InstanceType<typeof InfoDialog> | null>(null);
const passwordRef = ref<InstanceType<typeof PasswordDialog> | null>(null);
const openDialog = (ref: string) => {
  if (ref == "infoRef") infoRef.value?.openDialog(); 
  if (ref == "delAuthRef") infoRef.value?.closeGoogleAuthDialog(); 
  if (ref == "passwordRef") passwordRef.value?.openDialog();
  
};
</script>

<style scoped lang="scss">


.avatar {
  .iconView {
    width: 40px;
    height: 40px;
    overflow: hidden;
    cursor: pointer;
    border-radius: 50%;
    border: 1px solid #eee;
    overflow: hidden;

    img {
      width: 40px;
      height: 40px;
      border-radius: 50%;
    }
  }


}
</style>
