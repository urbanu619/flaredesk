import http from "@/api";

// 期任务
export const sysStakePeriodJobList = (params) => {
  return http.get("/admin/api/app/sys_stake_period_job/find", params);
};

// 日任务
export const sysJobList = (params) => {
  return http.get("/admin/api/app/sys_job/find", params);
};

export default {};
