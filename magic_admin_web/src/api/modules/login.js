import authMenuList from "@/assets/json/authMenuList.json";
import authButtonList from "@/assets/json/authButtonList.json";
import http from "@/api";

/**
 * @name 登录模块
 */

// 获取验证码
export const captchaGet = async () => {
  return await http.get("/admin/api/app/login/generateCaptcha", {}, {loading: false});
};

export const checkPost = async (data) => {
  return await http.post("/admin/api/app/login/check", data, {loading: true});
};


// 用户登录
export const loginApi = (data) => {
  return http.post("/admin/api/app/login/in", data, {loading: true}); // 正常 post json 请求  ==>  application/json
};

// 获取菜单列表
export const getAuthMenuListApi = () => {
  // 如果想让菜单变为本地数据，注释上一行代码，并引入本地 authMenuList.json 数据
  return authMenuList;
};

// 获取按钮权限
export const getAuthButtonListApi = () => {
  // 如果想让按钮权限变为本地数据，注释上一行代码，并引入本地 authButtonList.json 数据
  return authButtonList;
};

// 用户退出登录
export const logoutApi = () => {
  return http.post("/admin/api/app/login/out");
};



