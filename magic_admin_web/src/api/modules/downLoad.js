import http from "@/api";

// 文件上传模块
// 用户资产导出
export const userBillExport = (params) => {
  return http.get("/admin/api/fmc/bill/export", params, {})
};
