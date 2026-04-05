import {createApp} from "vue";
import App from "./App.vue";
// reset style sheet
import "@/styles/reset.scss";
// CSS common style sheet
import "@/styles/common.scss";
// iconfont css
import "@/assets/iconfont/iconfont.scss";
// font css
import "@/assets/fonts/font.scss";
// element css
import "element-plus/dist/index.css";
// element dark css
import "element-plus/theme-chalk/dark/css-vars.css";
// custom element dark css
import "@/styles/element-dark.scss";
// custom element css
import "@/styles/element.scss";
// svg icons
import "virtual:svg-icons-register";
// element plus
import ElementPlus from "element-plus";
// element icons
import * as Icons from "@element-plus/icons-vue";


// vue Router
import router from "@/routers";
// vue i18n
import I18n from "@/languages/index";
// pinia store
import pinia from "@/stores/index.js";
// errorHandler
import errorHandler from "@/utils/errorHandler";

// 自定义指令
import addressFormat from './directives/addressFormat'

// import {config} from "./contract/client.js";
import {WagmiPlugin} from '@wagmi/vue'
import {QueryClient, VueQueryPlugin} from "@tanstack/vue-query";

const queryClient = new QueryClient()

const app = createApp(App);

// app.config.errorHandler = errorHandler;

app.directive('address-format', addressFormat)

// register the element Icons component
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(ElementPlus).use(router).use(I18n).use(pinia).
  // use(WagmiPlugin, {config: config}).
  use(VueQueryPlugin, {queryClient}).
  mount("#app");
