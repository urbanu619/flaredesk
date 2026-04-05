import http from "@/api";


/* app设置  */

// 语言清单
export const proxyBizLang = (params) => {
  return http.get("/admin/api/proxy/biz/lang", params);
};
// 问题列表
export const qaFind = (params) => {
  return http.get("/admin/api/juapp/app_question_answer/find", params);
};
// 问题列表 -增加
export const qaAdd = (params) => {
  return http.post("/admin/api/proxy/biz/qa/add", params);
};

// 问题列表 -修改
export const qaEdit = (params) => {
  return http.post("/admin/api/proxy/biz/qa/edit", params);
};
// 问题列表 -删除
export const qaDel = (params) => {
  return http.post("/admin/api/proxy/biz/qa/del", params);
};

// 众筹仓配置
export const crowdfundWarehouses = (params) => {
  return http.get("/admin/api/juapp/app_crowdfund_warehouses/find", params);
};
