import {defineConfig, loadEnv, ConfigEnv, UserConfig} from "vite";
import {resolve} from "path";
import {wrapperEnv} from "./build/getEnv";

import {createVitePlugins} from "./build/plugins";
import pkg from "./package.json";
import dayjs from "dayjs";
import javascriptObfuscator from "vite-plugin-javascript-obfuscator";

// 扩展代理配置类型
interface CustomProxyOptions {
  target: string
  changeOrigin?: boolean
  secure?: boolean
  pathRewrite?: Record<string, string>
}

const {dependencies, devDependencies, name, version} = pkg;
const __APP_INFO__ = {
  pkg: {dependencies, devDependencies, name, version},
  lastBuildTime: dayjs().format("YYYY-MM-DD HH:mm:ss")
};

const timeStamp = new Date().getTime();

// @see: https://vitejs.dev/config/
export default defineConfig(({mode}: ConfigEnv): UserConfig => {
  const root = process.cwd();
  const env = loadEnv(mode, root);
  const viteEnv = wrapperEnv(env);

  return {
    base: "/",// viteEnv.VITE_PUBLIC_PATH,
    root,
    resolve: {
      alias: {
        "@": resolve(__dirname, "./src"),
        "vue-i18n": "vue-i18n/dist/vue-i18n.cjs.js"
      }
    },
    define: {
      __APP_INFO__: JSON.stringify(__APP_INFO__)
    },
    css: {
      preprocessorOptions: {
        scss: {
          additionalData: `@use "@/styles/var.scss";`,
          quietDeps: true  // 忽略依赖警告
        }
      }
    },
    // 应用代理配置
    server: {
      port: 8080,
      host: '0.0.0.0',
      proxy: {
        "/admin": {
          target: "http://127.0.0.1:2022",
          changeOrigin: true,
          secure: false
        } as CustomProxyOptions // 类型断言
      }
    },
    // plugins: createVitePlugins(viteEnv),
    plugins: [
      ...createVitePlugins(viteEnv),
      javascriptObfuscator({
        include: ['src/api/**/*.js', 'src/api/**/*.ts'],
        options: {
          compact: true,
          controlFlowFlattening: true,
          deadCodeInjection: false,
          rotateStringArray: true,
          stringArray: true,
          stringArrayEncoding: ['rc4'], // rc4 或 'base64'
          stringArrayThreshold: 1
        }
      })
    ],
    esbuild: {
      // pure: viteEnv.VITE_DROP_CONSOLE ? ["console.log", "debugger"] : []
    },
    build: {
      outDir: "dist",
      minify: "esbuild",
      // esbuild 打包更快，但是不能去除 console.log，terser打包慢，但能去除 console.log
      // minify: "terser",
      // terserOptions: {
      // 	compress: {
      // 		drop_console: viteEnv.VITE_DROP_CONSOLE,
      // 		drop_debugger: true
      // 	}
      // },
      sourcemap: false,
      // 禁用 gzip 压缩大小报告，可略微减少打包时间
      reportCompressedSize: false,
      // 规定触发警告的 chunk 大小
      chunkSizeWarningLimit: 5000,
      rollupOptions: {
        output: {
          // Static resource classification and packaging
          chunkFileNames: `assets/js/[name]-[hash].${timeStamp}.js`,
          entryFileNames: `assets/js/[name]-[hash].${timeStamp}.js`,
          assetFileNames: `assets/[ext]/[name]-[hash].${timeStamp}.[ext]`
        }
      }
    }
  };
});
