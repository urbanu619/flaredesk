import http from "@/api";

// 冲级活动
export const promotionUpgradeLevelList = (params) => {
  return http.get("/admin/api/app/promotion_upgrade_level/find", params);
};

export default {};
