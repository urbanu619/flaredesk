import http from "@/api";

// 用户节点订单
export const nodeOrderList = (params) => {
  return http.get("/admin/api/app/node_order/find", params);
};

// 创建用户节点订单
export const nodeOrderCreate = (params) => {
  return http.post("/admin/api/proxy/v1/adi/node/add/order", params);
};

// 用户节点支付信息
export const nodeOrderPaymentsList = (params) => {
  return http.get("/admin/api/app/node_order_payments/find", params);
};

// 节点Banner
export const nodeBannerList = (params) => {
  return http.get("/admin/api/app/node_banner/find", params);
};

// 节点配置
export const nodeInfoList = (params) => {
  return http.get("/admin/api/app/node_info/find", params);
};

// 更新节点配置
export const nodeInfoUpdate = (params) => {
  return http.post("/admin/api/proxy/v1/adi/node/info/update", params);
};

// 创建节点Banner
export const nodeBannerCreate = (params) => {
  return http.post("/admin/api/proxy/v1/adi/node/banner/create", params);
};

// 更新节点Banner
export const nodeBannerUpdate = (params) => {
  return http.post("/admin/api/proxy/v1/adi/node/banner/update", params);
};

export default {};
