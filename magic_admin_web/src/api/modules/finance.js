import http from "@/api";


/* 资产信息模块  */

// 转账清单
export const appTransferRecord = (params) => {
  return http.get("/admin/api/juapp/apo_transfer_record/find", params);
};

// 转账清单 导出
export const appTransferRecordExport = (params) => {
  return http.get("/admin/api/juapp/apo_transfer_record/export", params);
};


// 资产流水
export const billRecord = (params) => {
  return http.get("/admin/api/juapp/apo_asset_bill/find", params);
};


// 资产流水 - 导出
export const billRecordExport = (params) => {
  return http.get("/admin/api/juapp/apo_bill/export", params);

};

// 收益领取记录
export const apoUserClaimLogs = (params) => {
  return http.get("/admin/api/juapp/apo_user_claim_logs/find", params);
};

// 基金释放记录
export const apoUserFundReleaseLog = (params) => {
  return http.get("/admin/api/juapp/apo_user_fund_release_log/find", params);
};

// 奖励记录
export const apoUserRewardLog = (params) => {
  return http.get("/admin/api/juapp/apo_user_reward_record_views/find", params);
};


// 流水详情- 导出
export const billRecordExportDetail = (params) => {
  return http.get("/admin/api/juapp/apo_bill/export/detail", params);
};

// 分发日志
export const distributeTaskLog = (params) => {
  return http.get("/admin/api/juapp/app_transfer_distribute_task/find", params);
};

// 提现审核
export const handlerTransfer = (params) => {
  return http.post("/admin/api/proxy/v1/adi/asset/withdraw/verify", params);
};

// ju兑换记录
export const earnExchangeRecord = (params) => {
  return http.get("/admin/api/juapp/apo_exchange_record/find", params);
};

// 获取ju算力列表
export const exchangeJuRecord = (params) => {
  return http.post("/admin/api/proxy/mng/stage/exchange/ju/data", params);
};

// 兑换ju算力
export const exchangeJuPower = (params) => {
  return http.post("/admin/api/proxy/mng/exchange/ju/power", params);
};

// // 兑换ju算力
// export const exchangeJuPower = (params) => {
//   return http.post("/admin/api/proxy/mng/exchange/ju/power", params);
// };

// magic_asset 列表查询
export const magicAssetList = (params) => {
  return http.get("/admin/api/app/magic_asset/find", params);
};

// magic_asset_bill 列表查询
export const magicAssetBillList = (params) => {
  return http.get("/admin/api/app/magic_asset_bill/find", params);
};

// magic_asset_rw_record 列表查询
export const magicAssetRwRecordList = (params) => {
  return http.get("/admin/api/app/magic_asset_rw_record/find", params);
};

// magic_asset_config 列表查询
export const magicAssetConfigList = (params) => {
  return http.get("/admin/api/app/magic_asset_config/find", params);
};