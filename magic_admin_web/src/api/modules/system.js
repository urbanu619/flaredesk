import http from "@/api";

/**
 * @name 系统管理
 */

// 菜单树查看
export const getTreeList = () => {
  return http.get("/admin/api/sys/menu/tree", {});
};


// 获取管理员列表
export const getUserList = (params) => {
  return http.get("/admin/api/sys/user/find", params);
};

// 创建管理员
export const sysCreateUser = (params) => {
  return http.post("/admin/api/sys/user/create", params);
};

// 修改管理员
export const sysEditUser = (params) => {
  return http.post("/admin/api/sys/user/setUser", params);
};


//删除管理员

export const sysDelUser = (params) => {
  return http.get("/admin/api/sys/user/delete", params);
};


//获取角色列表
export const getSysRoleList = (params) => {
  return http.get("/admin/api/sys/role/find", params);
};

//创建菜单
export const sysMenuCreate = (params) => {
  return http.post("/admin/api/sys/menu/create", params);
};
//修改菜单
export const sysMenuSet = (params) => {
  return http.post("/admin/api/sys/menu/set", params);
};

//删除菜单
export const sysMenuDel = (params) => {
  return http.get("/admin/api/sys/menu/delete", params);
};

// 创建角色
export const sysRoleCreate = (params) => {
  return http.post("/admin/api/sys/role/create", params);
};
// 修改角色
export const sysRoleSet = (params) => {
  return http.post("/admin/api/sys/role/set", params);
};

// 删除角色
export const sysRoleDel = (params) => {
  return http.get("/admin/api/sys/role/delete", params);
};

// apis信息列表
export const apisFind = (params) => {
  return http.get("/admin/api/sys/apis/find", params);
};

// 创建api组
export const apisCreate = (params) => {
  return http.post("/admin/api/sys/apis/create", params);
};

// apis修改
export const apisSet = (params) => {
  return http.post("/admin/api/sys/apis/set", params);
};

// apis删除
export const apisDelete = (params) => {
  return http.get("/admin/api/sys/apis/delete", params);
};

// 取消谷歌密钥
export const cloearsysGoogle = (params) => {
  return http.post("/admin/api/sys/user/cancelGoogleKey", params);
};

// 设置/更改谷歌密钥
export const setsysGoogle = (params) => {
  return http.post("/admin/api/sys/user/replaceGoogleKey", params);
};

// 获取谷歌密钥
export const userGetGoogleKey = (params) => {
  return http.get("/admin/api/sys/user/getGoogleKey", params);
};