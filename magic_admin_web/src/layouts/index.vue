<!-- 💥 这里是一次性加载 LayoutComponents -->
<template>

  <el-watermark id="watermark" :font="font" :content="watermark ? ['Admin', 'Happy Working'] : ''">
    <component :is="LayoutVertical"/>
    <ThemeDrawer/>
  </el-watermark>
</template>

<script setup lang="ts" name="layout">
import {computed, reactive, watch, onMounted} from "vue";
import {ElMessage} from "element-plus";
import {useGlobalStore} from "@/stores/modules/global";
import ThemeDrawer from "./components/ThemeDrawer/index.vue";
import LayoutVertical from "./LayoutVertical/index.vue";
import { jwtDecode } from "jwt-decode"
import { useRouter } from "vue-router"
import dayjs from "dayjs"
import { LOGIN_URL } from "@/config";
import {useUserStore} from "@/stores/modules/user";
import {useAuthStore} from "@/stores/modules/auth"

const router = useRouter()
const userStore = useUserStore();
const useAuth = useAuthStore();

const globalStore = useGlobalStore();
const isDark = computed(() => globalStore.isDark);
const watermark = computed(() => globalStore.watermark);

const font = reactive({color: "rgba(0, 0, 0, .15)"});

onMounted(() => {
  const str = localStorage.getItem("userInfo")
  if(str == null) {
    router.push("/login")
  }
  const userInfo = JSON.parse(str)

  const token = jwtDecode(userInfo.token)
  if(dayjs.unix(token.exp!) < dayjs()) {
    ElMessage.error("登录超时")
    userStore.setToken("");
    userStore.setUserInfo({avatar: "", nickname: ""});
    useAuth.clearMenusList()
    localStorage.clear()
    // 3.重定向到登陆页
    router.replace(LOGIN_URL);
    router.push({path: '/login', replace: true})
  }
})

watch(isDark, () => (font.color = isDark.value ? "rgba(255, 255, 255, .15)" : "rgba(0, 0, 0, .15)"), {
  immediate: true
});
</script>

<style scoped lang="scss">
.layout {
  min-width: 600px;
}
</style>
