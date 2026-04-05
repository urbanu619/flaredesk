import http from "@/api";

// 文件上传模块
// 上传获取授权
export const uploadOssAuth = () => {
  return http.get("/admin/api/sys/file/oss/auth", {}, {})
};

// 文件上传
export const uploadUrl = (params) => {
  return http.upload("/admin/api/sys/file/upload", params, {})
};

// 文件上传
export const uploadImg = (params) => {
  return http.upload("/admin/api/sys/file/upload", params, {})
};

// 用户资产导出


