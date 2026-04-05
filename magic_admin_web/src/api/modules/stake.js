import http from "@/api";

// 活期用户信息
export const magicStakeUserCurrentOrderList = (params) => {
  return http.get("/admin/api/app/magic_stake_user_current_order/find", params);
};

// 活期质押记录
export const magicStakeUserCurrentOpsRecordList = (params) => {
  return http.get("/admin/api/app/magic_stake_user_current_ops_record/find", params);
};

export default {};
