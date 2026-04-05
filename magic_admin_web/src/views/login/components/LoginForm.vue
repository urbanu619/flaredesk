<template>
  <el-form ref="loginFormRef" :model="loginForm" :rules="loginRules" size="large">
    <el-form-item prop="username">
      <el-input v-model="loginForm.username" placeholder="用户名" @blur="validateUsername">
        <template #prefix>
          <el-icon class="el-input__icon">
            <user/>
          </el-icon>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item prop="password">
      <el-input v-model="loginForm.password" type="password" placeholder="密码"
                @keyup.native.enter="login(loginFormRef)"
                show-password autocomplete="new-password">
        <template #prefix>
          <el-icon class="el-input__icon">
            <lock/>
          </el-icon>
        </template>
      </el-input>
    </el-form-item>
    <el-form-item prop="code" v-if="loginForm.isNeedCaptcha">
      <div class="verify_code_cls">
        <el-input class="input" v-model="loginForm.code" type="text" autocomplete="off" maxlength="6"
                  @keyup.native.enter="login(loginFormRef)" placeholder="验证码" clearable/>
        <el-image class="img" @click="getCaptcha()" :src="captcha"/>
      </div>
    </el-form-item>
    <el-form-item prop="googleCode" v-if="loginForm.isSetGoogleAuth" @keyup.native.enter="login(loginFormRef)">
      <el-input maxlength="6" v-model="loginForm.googleCode" type="text" autocomplete="off"
                placeholder="Google Code" clearable/>
    </el-form-item>
  </el-form>
  <div class="login-btn">
    <el-button :icon="CircleClose" round size="large" @click="resetForm(loginFormRef)"> 重置</el-button>
    <el-button :icon="UserFilled" round size="large" type="primary" :loading="loading" @click="login(loginFormRef)">
      登录
    </el-button>
  </div>
</template>

<script setup>
import {ref, reactive, onMounted, onBeforeUnmount} from "vue";
import {useRouter} from "vue-router";
import {HOME_URL} from "@/config";


import {ElMessage} from "element-plus";
import {loginApi, captchaGet, checkPost} from "@/api/modules/login";


import {useUserStore} from "@/stores/modules/user";
import {useTabsStore} from "@/stores/modules/tabs";
import {useKeepAliveStore} from "@/stores/modules/keepAlive";

import {CircleClose, UserFilled} from "@element-plus/icons-vue";


import {useAuthStore} from "@/stores/modules/auth.js"


const useAuth = useAuthStore()


const router = useRouter();
const userStore = useUserStore();
const tabsStore = useTabsStore();
const keepAliveStore = useKeepAliveStore();

const loginFormRef = ref();


const loginRules = reactive({
  username: [{required: true, message: "请输入用户名", trigger: "blur"}],
  password: [{required: true, message: "请输入密码", trigger: "blur"}]
});

const loading = ref(false);
const loginForm = reactive({
  code: "",
  cid: "",  // 图形验证ID
  googleCode: "", // 谷歌验证码
  username: "",  // 用户名
  password: "",  // 密码
  isNeedCaptcha: false,  // 是否需要图形验证
  isSetGoogleAuth: false, // 是否需要谷歌验证
});
const captcha = ref("")


onMounted(() => {
  if (userStore.token) {
    router.push(HOME_URL);
    return
  }
  // 监听 enter 事件（调用登录）
  // document.onkeydown = (e) => {
  //   if (e.code === "Enter" || e.code === "enter" || e.code === "NumpadEnter") {
  //     if (loading.value) return;
  //     login(loginFormRef.value);
  //   }
  // };
  getCaptcha()
})
;
// 获取验证码
const getCaptcha = async () => {
  let resp = await captchaGet()
  if (resp.code === 200) {
    await setCapcha(resp.data)
  } else {
    ElMessage.error(resp.msg)
  }
}
const validateUsername = async () => {

  if(!loginForm.username){
    return
  }
  
  let params = {
    username: loginForm.username,
  }
  let resp = await checkPost(params)
  console.log("账号信息检查返回信息 resp:", resp)
  if (resp.code === 200) {
    let isSetGoogleAuth = resp.data?.isSetGoogleAuth ?? false
    let isNeedCaptcha = resp.data?.isNeedCaptcha ?? false
    setVerifyData(isSetGoogleAuth, isNeedCaptcha)

    if(isSetGoogleAuth || isNeedCaptcha){
      // 需要输入谷歌验证码
      return resp.data
    }

    if (!resp.data.userIsExit) {
      ElMessage.error("fail times:" + resp.data.t)
    }

  } else {
    ElMessage.error(resp.msg)
  }
}

const setVerifyData = (isSetGoogleAuth, isNeedCaptcha) => {
  loginForm.isSetGoogleAuth = isSetGoogleAuth
  loginForm.isNeedCaptcha = isNeedCaptcha
}
const setCapcha = async (data) => {
  loginForm.cid = data.cid
  captcha.value = data.captcha
}


// login
const login = async (formEl) => {

  // 检测账号
  if(!loginForm.isNeedCaptcha){
    const res = await validateUsername()
    console.log(res, 'res')
    // 需要输入谷歌验证码或者验证码，但是没有输入，不继续往下走了
    if(res && (res.isNeedCaptcha && !loginForm.code)){
      ElMessage.error("请输入验证码")
      return 
    }
    if(res && (res.isSetGoogleAuth && !loginForm.googleCode)){
      ElMessage.error("请输入谷歌验证码")
      return 
    }
  }
  
  if (!formEl) return;
  formEl.validate(async valid => {
    if (!valid) return;
    loading.value = true;
    try {
      // 1.执行登录接口
      let params = {
        ...loginForm
      }
      const res = await loginApi(params);
      loading.value = false;
      if (res.code === 200) {
        let userinfo = {
          avatar: res.data?.avatar,
          isSuper: res.data?.isSuper,
          isSetGoogleAuth: res.data?.isSetGoogleAuth,
          nickname: res.data?.nickname
        }
        userStore.setUserInfo(userinfo)
        userStore.setToken(res.data?.token);

        await useAuth.getAuthMenuList()

        // let resp = await getTree()
        // if (resp) {
        //   // 2.添加动态路由
        // // 3.清空 tabs、keepAlive 数据
        await tabsStore.setTabs([]);
        await keepAliveStore.setKeepAliveName([]);
        // 4.跳转到首页
        await router.push("/home/index");
        // }
      } else {
        ElMessage.error(res.msg)
      }
    } finally {
      loading.value = false;
    }
  });
};


// resetForm
const resetForm = (formEl) => {
  if (!formEl) return;
  formEl.resetFields();
};

onBeforeUnmount(() => {
  document.onkeydown = null;
});
</script>

<style scoped lang="scss">
@use "../index.scss";
</style>
