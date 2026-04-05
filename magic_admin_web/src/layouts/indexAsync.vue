<!-- 💥 这里是异步加载 LayoutComponents -->
<template>
  <el-watermark id="watermark" :font="font" :content="watermark ? ['Admin', 'Happy Working'] : ''">
    <suspense>
      <template #default>
        <component :is="LayoutVertical" />
      </template>
      <template #fallback>
        <Loading />
      </template>
    </suspense>
    <ThemeDrawer />
  </el-watermark>
</template>

<script setup lang="ts" name="layoutAsync">
import { computed, defineAsyncComponent, reactive, watch, type Component, onMounted } from "vue";
import { LayoutType } from "@/stores/interface";
import { useGlobalStore } from "@/stores/modules/global";
import Loading from "@/components/Loading/index.vue";
import ThemeDrawer from "./components/ThemeDrawer/index.vue";
import LayoutVertical from "@/layouts/LayoutVertical/index.vue";
import {useUserStore} from "@/stores/modules/user";
import {useAuthStore} from "@/stores/modules/auth"
import { jwtDecode } from "jwt-decode";
import dayjs from "dayjs";
import { ElMessage } from "element-plus";
import { useRouter } from "vue-router";
import { LOGIN_URL } from "@/config";

const router = useRouter()
const userStore = useUserStore();
const useAuth = useAuthStore();

const globalStore = useGlobalStore();

const isDark = computed(() => globalStore.isDark);
const layout = computed(() => globalStore.layout);
const watermark = computed(() => globalStore.watermark);

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

const font = reactive({ color: "rgba(0, 0, 0, .15)" });
watch(isDark, () => (font.color = isDark.value ? "rgba(255, 255, 255, .15)" : "rgba(0, 0, 0, .15)"), {
  immediate: true
});
</script>

<style scoped lang="scss">
.layout {
  min-width: 600px;
}
</style>
