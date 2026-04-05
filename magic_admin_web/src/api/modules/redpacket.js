import http from "@/api";

const BASE = "/admin/api/app/tg_red_packet";

// ==================== 红包配置管理 ====================

/**
 * 获取红包配置列表
 * @param {Object} params - 查询参数
 */
export const getRedPacketConfigList = params => {
  return http.get(`${BASE}_config/find`, params);
};

/**
 * 获取单个红包配置
 * @param {Number} id - 配置ID
 */
export const getRedPacketConfig = id => {
  return http.get(`${BASE}_config/get`, { id });
};

/**
 * 创建红包配置
 * @param {Object} data - 配置数据
 */
export const createRedPacketConfig = data => {
  return http.post(`${BASE}_config/create`, data);
};

/**
 * 更新红包配置
 * @param {Object} data - 配置数据
 */
export const updateRedPacketConfig = data => {
  return http.post(`${BASE}_config/update`, data);
};

/**
 * 删除红包配置
 * @param {Number} id - 配置ID
 */
export const deleteRedPacketConfig = id => {
  return http.post(`${BASE}_config/delete`, { id });
};

/**
 * 切换配置状态（启用/禁用）
 * @param {Number} id - 配置ID
 * @param {Number} status - 状态：1=启用, 2=禁用
 */
export const toggleRedPacketConfigStatus = (id, status) => {
  return http.post(`${BASE}_config/toggleStatus`, { id, status });
};

/**
 * 获取表字段注释
 */
export const getRedPacketConfigComment = () => {
  return http.get(`${BASE}_config/comment`);
};

// ==================== 红包发送 ====================

/**
 * 基于配置ID手动发送红包
 * @param {Number} configId - 配置ID
 */
export const sendRedPacketManual = configId => {
  return http.post(`${BASE}_send/sendManual`, { configId });
};

/**
 * 直接发送红包（无需预先配置）
 * @param {Object} data - 红包数据
 */
export const sendRedPacketDirect = data => {
  return http.post(`${BASE}_send/sendDirect`, data);
};

/**
 * 获取Telegram群组列表
 */
export const getTelegramGroups = () => {
  return http.get(`${BASE}_send/getGroups`);
};

// ==================== 红包记录 ====================

/**
 * 获取红包记录列表
 * @param {Object} params - 查询参数
 */
export const getRedPacketRecordList = params => {
  return http.get(`/admin/api/app/tg_red_packet/find`, params);
};

/**
 * 获取抢红包记录列表（领取明细）
 * @param {Object} params - 查询参数
 */
export const getGrabRecordList = params => {
  return http.get(`/admin/api/app/tg_red_packet/findGrabRecords`, params);
};

// ==================== 群组管理 ====================

/**
 * 获取群组列表
 * @param {Object} params - 查询参数
 */
export const getGroupList = params => {
  return http.get(`/admin/api/app/tg_group/find`, params);
};

/**
 * 获取群组详情
 * @param {Number} id - 群组ID
 */
export const getGroup = id => {
  return http.get(`/admin/api/app/tg_group/get`, { id });
};

/**
 * 创建群组
 * @param {Object} data - 群组数据
 */
export const createGroup = data => {
  return http.post(`/admin/api/app/tg_group/create`, data);
};

/**
 * 更新群组
 * @param {Object} data - 群组数据
 */
export const updateGroup = data => {
  return http.post(`/admin/api/app/tg_group/update`, data);
};

/**
 * 删除群组
 * @param {Number} id - 群组ID
 */
export const deleteGroup = id => {
  return http.get(`/admin/api/app/tg_group/delete`, { id });
};

/**
 * 从Bot同步群组信息
 */
export const syncGroupsFromBot = () => {
  return http.post(`/admin/api/app/tg_group/syncFromBot`);
};

// ==================== 用户绑定管理 ====================

/**
 * 获取用户绑定列表
 * @param {Object} params - 查询参数
 */
export const getUserBindList = params => {
  return http.get(`/admin/api/app/tg_user_bind/find`, params);
};

/**
 * 获取用户绑定详情
 * @param {Number} id - 绑定ID
 */
export const getUserBind = id => {
  return http.get(`/admin/api/app/tg_user_bind/get`, { id });
};

/**
 * 创建用户绑定
 * @param {Object} data - 绑定数据
 */
export const createUserBind = data => {
  return http.post(`/admin/api/app/tg_user_bind/create`, data);
};

/**
 * 更新用户绑定
 * @param {Object} data - 绑定数据
 */
export const updateUserBind = data => {
  return http.post(`/admin/api/app/tg_user_bind/update`, data);
};

/**
 * 解绑用户（软删除）
 * @param {Number} id - 绑定ID
 */
export const deleteUserBind = id => {
  return http.get(`/admin/api/app/tg_user_bind/delete`, { id });
};

// ==================== 用户资产管理 ====================

/**
 * 获取用户资产列表
 * @param {Object} params - 查询参数
 */
export const getUserAssetList = params => {
  return http.get(`/admin/api/app/magic_asset/find`, params);
};

/**
 * 获取资产流水明细
 * @param {Object} params - 查询参数（userId, symbol 等）
 */
export const getUserAssetBillList = params => {
  return http.get(`/admin/api/app/magic_asset_bill/find`, params);
};
