import http from "@/api";


// 修改用户信息
export const userSetUser = (params) => {
  return http.post(`/admin/api/proxy/v1/adi/user/set/info`, params);
};

// 重置用户上级
export const userSetParent = (params) => {
  return http.post("/admin/api/proxy/v1/adi/user/set/parent", params);
};

// 获取用户资产
export const getUserAsset = (params) => {
  return http.get("/admin/api/juapp/apo_asset/find", params);
};

// 用户资产 - 导出
export const getUserAssetExport = (params) => {
  return http.get("/admin/api/juapp/apo_asset/export", params);
};

// 系统划转
export const handlerAsset = (params) => {
  return http.post(`/admin/api/proxy/mng/handler/asset`, params);
};


/* 用户管理模块  */


/**
 * 获取用户信息
 */
export const userBaseInfoList = (params) => {
  return http.get("/admin/api/app/macmagic_user/find", params);
};

// magic_user 列表查询
export const magicUserList = (params) => {
  return http.get("/admin/api/app/magic_user/find", params);
};


// magic_user_quota 列表查询
export const magicUserQuotaList = (params) => {
  return http.get("/admin/api/app/magic_user_quota/find", params);
};

// 更新用户配额
export const magicUserQuotaUpdate = (params) => {
  return http.post("/admin/api/proxy/v1/adi/user/set/quota", params);
};

// magic_user_profit 列表查询
export const magicUserProfitList = (params) => {
  return http.get("/admin/api/app/magic_user_profit/find", params);
};

// magic_user_profit_record 列表查询
export const magicUserProfitRecordList = (params) => {
  return http.get("/admin/api/app/magic_user_profit_record/find", params);
};

// magic_user_logs 列表查询
export const magicUserLogsList = (params) => {
  return http.get("/admin/api/app/magic_user_logs/find", params);
};

// magic_user_action_log 列表查询
export const magicUserActionLogList = (params) => {
  return http.get("/admin/api/app/magic_user_action_log/find", params);
};

