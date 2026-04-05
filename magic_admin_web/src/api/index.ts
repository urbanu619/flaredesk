import axios, {AxiosInstance} from "axios";
import {showFullScreenLoading, tryHideFullScreenLoading} from "@/components/Loading/fullScreen";
import {LOGIN_URL} from "@/config";
import {ElMessage} from "element-plus";

import {ResultEnum} from "@/enums/httpEnum";
import {checkStatus} from "./helper/checkStatus";
import {AxiosCanceler} from "./helper/axiosCancel";
import {useUserStore} from "@/stores/modules/user";
import router from "@/routers";

const headers_json: any = {
  "Content-Type": "application/json; charset=utf-8",
  Authorization: ""
}
const headers_form: any = {
  'Content-Type': 'multipart/form-data',
}
let config = {
  // 默认地址请求地址，可在 .env.** 文件中修改
  // baseURL: import.meta.env.VITE_API_URL,
  baseURL: '/',
  // 设置超时时间
  timeout: ResultEnum.TIMEOUT,
  // 跨域时候允许携带凭证
  // withCredentials: true
  cancel: true,
  loading: false,
  headers: {
    "Accept": "application/json",
    'Content-Type': "application/json; charset=utf-8",
    Authorization: ""
  },
  isJsonn: null,
};


const axiosCanceler = new AxiosCanceler();

// 辅助函数：移除对象中值为 null 或 undefined 的属性
function cleanParams(obj: any): any {
  if (obj instanceof FormData) return obj;
  if (Array.isArray(obj)) {
    return obj
      .map(item => cleanParams(item))
      .filter(item => item !== null && item !== undefined);
  } else if (typeof obj === 'object' && obj !== null) {
    const cleanedObj: any = {};
    Object.keys(obj).forEach(key => {
      const value = obj[key];
      if (value !== null && value !== undefined) {
        cleanedObj[key] = cleanParams(value);
      }
    });
    return cleanedObj;
  }
  return obj;
}

class RequestHttp {
  service: AxiosInstance;

  public constructor(config) {

    this.service = axios.create(config);
    //请求拦截器
    this.service.interceptors.request.use((config) => {
        const userStore = useUserStore();
        axiosCanceler.addPending(config);
        showFullScreenLoading();
        config.headers['Content-Type'] = config.isJsonn ? "application/json; charset=utf-8" : 'multipart/form-data'
        if (userStore.token) {
          config.headers.Authorization = userStore.token;
        }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      });

    /**
     * @description 响应拦截器
     *  服务器换返回信息 -> [拦截统一处理] -> 客户端JS获取到信息
     */
    this.service.interceptors.response.use((response) => {


        const {data, config} = response;

        const userStore = useUserStore();
        axiosCanceler.removePending(config);
        // config.loading &&
        tryHideFullScreenLoading();
        // 登录失效
        if (data.code == ResultEnum.OVERDUE) {
          userStore.setToken("");
          router.replace(LOGIN_URL);
          ElMessage.error(data.msg);
          return Promise.reject(data);
        }
        // 全局错误信息拦截（防止下载文件的时候返回数据流，没有 code 直接报错）
        if (data.code && data.code !== ResultEnum.SUCCESS) {
          ElMessage.error(data.msg);
          return Promise.reject(data);
        }
        // 成功请求（在页面上除非特殊情况，否则不用处理失败逻辑）
        return data;
      },
      async (error) => {
        console.error(error)
        const {response} = error;
        tryHideFullScreenLoading();
        // 请求超时 && 网络错误单独判断，没有 response
        if (error.message.indexOf("timeout") !== -1) ElMessage.error("请求超时！请您稍后重试");
        if (error.message.indexOf("Network Error") !== -1) ElMessage.error("网络错误！请您稍后重试");
        // 根据服务器响应的错误状态码，做不同的处理
        if (response) checkStatus(response.status);
        // 服务器结果都没有返回(可能服务器错误可能客户端断网)，断网处理:可以跳转到断网页面
        if (!window.navigator.onLine) router.replace("/500");
        return Promise.reject(error);
      }
    );
  }

  /**
   * @description 常用请求方法封装
   */
  get(url, params, _object) {
    _object = {
      isJsonn: true
    }
    //return this.service.get(url, { params, ..._object });
    return this.service.get(url, { params: cleanParams(params), ..._object });
  }

  post(url, params, _object) {
    _object = {
      isJsonn: true
    }
    // return this.service.post(url, params, { ..._object });
    return this.service.post(url, cleanParams(params), { ..._object });
  }

  put(url, params, _object) {
    _object = {
      isJsonn: true
    }
    // return this.service.put(url, params, _object);
    return this.service.put(url, cleanParams(params), _object);
  }

  delete(url, params, _object) {
    _object = {
      isJsonn: true
    }
    // return this.service.delete(url, { params, ..._object });
    return this.service.delete(url, { params: cleanParams(params), ..._object });
  }

  upload(url, params, _object) {
    _object = {
      isJsonn: false
    }
    return this.service.post(url, params, {..._object});
  }

  download(url, params, _object) {
    _object = {
      isJsonn: true
    }
    // return this.service.post(url, params, { ..._object, responseType: "blob" });
    return this.service.post(url, cleanParams(params), { ..._object, responseType: "blob" });
  }
}

export default new RequestHttp(config);
