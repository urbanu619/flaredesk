import http from "@/api";


/* 产品管理  */

// 产品列表
export const apoProduct = (params) => {
  return http.get("/admin/api/juapp/apo_product/find", params);
};

// 用户活期产品
export const currentProduct = (params) => {
  return http.get("/admin/api/juapp/apo_user_current_product/find", params);
};

//  活期操作日志
export const currentProductLog = (params) => {
  return http.get("/admin/api/juapp/apo_user_current_wealth_log/find", params);
};

//  用户定期产品
export const regularProduct = (params) => {
  return http.get("/admin/api/juapp/apo_user_regula_product/find", params);
};


//  用户定期订单
export const regularProductOrder = (params) => {
  return http.get("/admin/api/juapp/apo_user_regula_order/find", params);
};


// 定期操作日志
export const regularProductLog = (params) => {
  return http.get("/admin/api/juapp/apo_user_regula_wealth_log/find", params);
};
