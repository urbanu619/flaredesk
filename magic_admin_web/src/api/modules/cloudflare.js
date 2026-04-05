import http from "@/api";

const ACCOUNT = "/admin/api/app/cf_account";
const DNS = "/admin/api/app/cf_dns";

// ==================== 账号管理 ====================

export const getCfAccountList = params => http.get(`${ACCOUNT}/find`, params);
export const getCfAccount = id => http.get(`${ACCOUNT}/get`, { id });
export const createCfAccount = data => http.post(`${ACCOUNT}/create`, data);
export const updateCfAccount = data => http.post(`${ACCOUNT}/update`, data);
export const deleteCfAccount = id => http.get(`${ACCOUNT}/delete`, { id });

// ==================== DNS 管理 ====================

export const getCfZones = accountId => http.get(`${DNS}/zones`, { account_id: accountId });
export const getCfRecords = params => http.get(`${DNS}/records`, params);
export const createCfRecord = data => http.post(`${DNS}/record/create`, data);
export const batchCreateCfRecord = data => http.post(`${DNS}/record/batch_create`, data);
export const updateCfRecord = data => http.post(`${DNS}/record/update`, data);
export const deleteCfRecord = params => http.get(`${DNS}/record/delete`, params);
export const toggleCfProxy = data => http.post(`${DNS}/record/toggle_proxy`, data);
export const crossZoneDeleteRecords = data => http.post(`${DNS}/cross_zone/delete`, data);
export const crossZoneToggleProxy = data => http.post(`${DNS}/cross_zone/toggle_proxy`, data);

// ==================== Zone 本地缓存 ====================
const ZONE = "/admin/api/app/cf_zone";
export const syncCfZones = accountId => http.get(`${ZONE}/sync`, { account_id: accountId });
export const getLocalZoneList = params => http.get(`${ZONE}/find`, params);
export const getLocalZoneAll = accountId => http.get(`${ZONE}/all`, { account_id: accountId });
export const updateZoneRemark = data => http.post(`${ZONE}/update_remark`, data);

// ==================== DNS 模板 ====================
const TEMPLATE = "/admin/api/app/cf_dns_template";
export const getDnsTemplateList = params => http.get(`${TEMPLATE}/find`, params);
export const getDnsTemplate = id => http.get(`${TEMPLATE}/get`, { id });
export const createDnsTemplate = data => http.post(`${TEMPLATE}/create`, data);
export const updateDnsTemplate = data => http.post(`${TEMPLATE}/update`, data);
export const deleteDnsTemplate = id => http.get(`${TEMPLATE}/delete`, { id });
