import http from "@/api";


/* APP设置模块  */

// AIC众筹管理 - 各期众筹记录
export const appCrowdfundRecord = (params) => {
  return http.get("/admin/api/juapp/app_crowdfund_record/find", params);
};

// AIC众筹管理 - 地址众筹日志
export const appCrowdfundRecordLogs = (params) => {
  return http.get("/admin/api/juapp/app_crowdfund_record_logs/find", params);
};

// AIC众筹管理 - 众筹轮信息
export const appCrowdfundRounds = (params) => {
  return http.get("/admin/api/juapp/app_crowdfund_rounds/find", params);
};

// AIC众筹管理 - 众筹期信息
export const appCrowdfundperiod = (params) => {
  return http.get("/admin/api/juapp/app_crowdfund_period/find", params);
};

// AIC众筹管理 - 每日任务
export const appCrowdfundJob= (params) => {
  return http.get("/admin/api/juapp/app_crowdfund_job/find", params);
};

// FMC矿场管理 - 矿场产出任务
export const miningJob= (params) => {
  return http.get("/admin/api/juapp/fmc_mining_job/find", params);
};

// FMC矿场管理 - FMC矿机信息
export const miningMachine= (params) => {
  return http.get("/admin/api/juapp/fmc_mining_machine/find", params);
};

// FMC矿场管理 - 其他矿机
export const miningMachineOther= (params) => {
  return http.get("/admin/api/juapp/fmc_mining_compensation_machine/find", params);
};

// FMC矿场管理 - FMC矿产记录
export const miningMachineRecord= (params) => {
  return http.get("/admin/api/juapp/fmc_mining_record/find", params);
};

// APP设置 - 产品信息
export const appSysConfig= (params) => {
  return http.get("/admin/api/juapp/earn_product/find", params);
};

// APP设置 - 等级配置
export const appLevelConfig= (params) => {
  return http.get("/admin/api/juapp/sys_level_config/find", params);
};

// APP设置 - 产品配置
export const appSysInfo= (params) => {
  return http.get("/admin/api/juapp/sys_config/find", params);
};

// APP设置 - 产品配置 - 设置修改系统配置
export const setSysConfig = (params) => {
  return http.post("/admin/api/proxy/biz/set/sys/config", params);
};


// APP设置 - 设置修改等级配置
export const setLevelConfig = (params) => {
  return http.post("/admin/api/proxy/mng/set/level/config", params);
};

// APP设置 - 设置产品配置
export const setSysInfo = (params) => {
  return http.post("/admin/api/proxy/admin/set/sys/config", params);
};


// APP设置 - 日志
export const sysLog = (params) => {
  return http.get("/admin/api/juapp/sys_job/find", params);
};

// APP设置 - 产品信息 - 更新产品
export const updateProduct = (params) => {
  return http.post("/admin/api/proxy/mng/update/product", params);
};

// 系统币种配置
export const sysCoinList = (params) => {
  return http.get("/admin/api/app/sys_coin/find", params);
};

// 创建币种
export const sysCoinUpdate = (params) => {
  return http.post("/admin/api/proxy/v1/adi/asset/set/coin", params);
};

// 等级配置（appLevelConfig 已存在，增加统一命名）
export const sysLevelConfigList = (params) => {
  return http.get("/admin/api/app/sys_level_config/find", params);
};

// 签名配置
export const sysSignConfigList = (params) => {
  return http.get("/admin/api/app/sys_sign_config/find", params);
};

// 创建等级配置
export const sysLevelConfigCreate = (params) => {
  return http.post("/admin/api/proxy/v1/adi/sys/level/config", params);
};

// 产品配置列表
export const magicStakeProductList = (params) => {
  return http.get("/admin/api/app/magic_stake_product/find", params);
};

// 更新产品配置
export const magicStakeProductUpdate = (params) => {
  return http.post("/admin/api/proxy/v1/adi/magic/update/product", params);
};