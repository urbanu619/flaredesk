// ? 系统全局字典

/**
 * @description：用户性别
 */
export const genderType = [
  { label: "男", value: 1 },
  { label: "女", value: 2 }
];

/**
 * @description：用户状态
 */
export const userStatus = [
  { label: "启用", value: 1, tagType: "success" },
  { label: "禁用", value: 0, tagType: "danger" }
];

/**
 * @description：提现/锁定状态（提现锁定示例）
 */
export const lockWithdrawOption = [
  { label: "未锁定", value: 0 },
  { label: "锁定", value: 1 }
];

/**
 * @description：白名单标记
 */
export const whiteListOption = [
  { label: "正常用户", value: 0 },
  { label: "白名单", value: 1 }
];

/**
 * @description：是否复投（布尔映射）
 */
export const booleanOption = [
  { label: "否", value: false },
  { label: "是", value: true }
];


/**
 * @description： 用户等级
 *
 */
export const levelListOption = [
  { label: "VIP0", value: 0 },
  { label: "VIP1", value: 1 },
  { label: "VIP2", value: 2 },
  { label: "VIP3", value: 3 },
  { label: "VIP4", value: 4 },
  { label: "VIP5", value: 5 },
  { label: "VIP6", value: 6 },
  { label: "VIP7", value: 7 },
];

/**
 * @description：币种下拉 TODO
 */
export const symbolList = [
  { label: "AIC", value: "AIC" },
  // { label: "FMC", value: "FMC" }
];

/**
 * @description：币种划转方式
 */
export const transferMethodList = [
  { label: "增加", value: "+" },
  { label: "减少", value: "-" }
];

/**
 * @description：交易类型下拉
 * user_to_system:充值 
 * system_to_user:提现
 */
export const transactionTypeList = [
  { label: "充值", value: "user_to_system" },
  { label: "提现", value: "system_to_user" }
];

/**
 * @description：交易状态 
 *  " 提现状态:0 无状态 1 申请冻结 2 提现完成 3 提现失败退回",
 */
export const transferStatusList = [
  { label: "无状态", value: 0 },
  { label: "申请冻结", value: 1 }, 
  { label: "提现完成", value: 2}, 
  { label: "提现失败退回", value: 3}
];


/**
 * @description：分发日志-业务场景下拉
 */
export const businessNameList = [
  { label: "购买定期", value: 1 },
  { label: "购买活期", value: 2 },
  { label: "赎回定期", value: 3 },
  { label: "赎回活期", value: 4 },
  { label: "收益复投", value: 5 },
  { label: "加速释放", value: 6 },
  { label: "兑换JU算力", value: 7 },
  { label: "充值", value: 101 },
  { label: "提现", value: 102 },
  { label: "系统转入", value: 105 },
  { label: "系统转出", value: 106 },
  { label: "系统迁移", value: 107 },
  { label: "定期收益", value: 201 },
  { label: "活期收益", value: 202 },
  { label: "直推收益", value: 203 },
  { label: "团队收益", value: 204 },
  { label: "平级收益", value: 205 },
];

/**
 * @description：众筹状态 
 * 众筹状态 0 待开启 1 进行中 2 结束
 */
export const warehouseStatusList = [
  { label: "待开启", value: 0 },
  { label: "进行中", value: 1 },
  { label: "结束", value: 2 },
];



/**
 * @description：ju记录状态
 * ExchangeRecordStatus        = 0 // 待处理
// ExchangeRecordStatusing     = 1 // 处理中
// ExchangeRecordStatusSuccess = 2 // 处理成功
// ExchangeRecordStatusFail    = 3 // 处理失败
 */
export const juStatusOption = [
  { label: "待处理", value: 0 },
  { label: "处理中", value: 1 },
  { label: "成功", value: 2 },
  { label: "失败", value: 3 },
];

/**
 * @description：兑换阶段
 */
export const juExchangeStepOption = [
  { label: "2024-12-31前", value: '1' },
  { label: "2025-01-01至04-30间", value: '2' },
  { label: "2025-05-01后", value: '3' },
]


/**
 * @description：期众筹状态 
 * 0 待开启 1 进行中 2 众筹成功 3 众筹失败 
 */
export const roundStatusList = [
  { label: "待开启", value: 0 },
  { label: "进行中", value: 1 },
  { label: "众筹成功", value: 2 },
  { label: "众筹失败", value: 3 },
];

/**
 * @description：期众筹结算状态 
 * 状态: 0 待结算 1 结算中 2 正常结算完成 3 爆仓结算完成
 */
export const warehouseSettlementStatusList = [
  { label: "待结算", value: 0 },
  { label: "结算中", value: 1 },
  { label: "正常结算完成", value: 2 },
  { label: "爆仓结算完成", value: 3 },
];

/**
 * @description：期众筹结算状态 
 * 分发状态:0 待分发 1 已结算-分发中 2 分发完成 3 爆仓退回
 */
export const periodDistributeStatusList = [
  { label: "待分发", value: 0 },
  { label: "已结算-分发中", value: 1 },
  { label: "分发完成", value: 2 },
  { label: "爆仓退回", value: 3 },
];

/**
 * @description：众筹分发类型 
 *  分发 1  结算 2  奖励 3 瓜分 4
 */
export const distributeTypeOptions = [
  { label: "分发", value: 1 },
  { label: "结算", value: 2 },
  { label: "奖励", value: 3 },
  { label: "瓜分", value: 4 },
];

/**
 * @description：任务状态
 *  1 进行中 2 已完成 3 失败
 */
export const taskStateOption = [
  { label: "进行中", value: 1 },
  { label: "已完成", value: 2 },
  { label: "失败", value: 3 },
];

/**
 * @description：sys_job - 快照任务状态
 * waiting: 等待执行
 * running: 进行中
 * finish: 完成
 * fail: 失败
 */
export const sysJobStateOption = [
  { label: "等待执行", value: "waiting" },
  { label: "进行中", value: "running" },
  { label: "完成", value: "finish" },
  { label: "失败", value: "fail" },
];

/**
 * @description：sys_stake_period_job - 期任务状态
 * running: 进行中
 * success: 成功
 * failure: 失败
 */
export const stakePeriodJobStateOption = [
  { label: "等待执行", value: "waiting" },

  { label: "进行中", value: "running" },
  { label: "成功", value: "success" },
  { label: "失败", value: "failure" },
];

/**
 * @description：分发状态
 *  分发状态: 0 待发放 1 已发放 2 发放失败
 */
export const distributeStateOptions = [
  { label: "待发放", value: 0 },
  { label: "已发放", value: 1 },
  { label: "发放失败", value: 2 },
];

/**
 * @description：资产信息-分发日志-链上分发状态
 * 0 待申请 1 申请失败 2 已申请成功 3 任务执行成功 4 任务执行失败
 */
export const distributeWithdrawStateOptions = [
  { label: "待申请", value: 0 },
  { label: "申请失败", value: 1 },
  { label: "已申请成功", value: 2 },
  { label: "任务执行成功", value: 3 },
  { label: "任务执行失败", value: 4 },
];


/**
 * @description： 矿机补偿状态
 * 0 待补偿 1补偿完成 2补偿失败
 */
export const compensationStatusList = [
  { label: "待补偿", value: 0 },
  { label: "补偿完成", value: 1 },
  { label: "补偿失败", value: 2 },
];


/**
 * @description： FMC矿产记录分发状态
 * 0 待发放 1 已发放 2 发放失败
 */
export const distributeStatusList = [
  { label: "待发放", value: 0 },
  { label: "已发放", value: 1 },
  { label: "发放失败", value: 2 },
];

/**
 * @description： 节点等级
 * 1 创世节点 2 超级节点 3 社区节点 4 个人节点
 */
export const nodeLevelList = [
  { label: "无节点", value: 0 },
  { label: "创世节点", value: 1 },
  { label: "超级节点", value: 2 },
  { label: "社区节点", value: 3 },
  { label: "个人节点", value: 4 },
];

/**
 * @description： 用户等级
 * 1 创世节点 2 超级节点 3 社区节点 4 个人节点
 */
export const levelList = [
  { label: "VIP0", value: 0 },
  { label: "VIP1", value: 1 },
  { label: "VIP2", value: 2 },
  { label: "VIP3", value: 3 },
  { label: "VIP4", value: 4 },
  { label: "VIP5", value: 5 },
  { label: "VIP6", value: 6 },
  { label: "VIP7", value: 7 },
  { label: "VIP8", value: 8 },
  { label: "VIP9", value: 9 },
  { label: "一星 ⭐️", value: 10 },
  { label: "二星 ⭐️⭐️", value: 11 },
  { label: "三星 ⭐️⭐️⭐️", value: 12 },
];

/**
 * @description： 变更类型
 *  :修改地址 重置上级 修改上级
 */
export const actionlList = [
  { label: "修改地址", value: '修改地址' },
  { label: "重置上级", value: '重置上级' },
  { label: "修改上级", value: '修改上级' }
];

/**
 * @description： 区块提现状态
 *  0 创建 1 挂起 2 成功 3 失败 4 跳过
 */
export const blockWithdrawStateOption = [
  { label: "创建", value: 0},
  { label: "pending", value: 1},
  { label: "成功", value: 2},
  { label: "失败", value: 3},
  { label: "跳过", value: 4},
];

/**
 * @description： 合约操作表状态
 *  0 创建 1 已操作 2 已完成
 */
export const blockActionOrderOption = [
  { label: "创建", value: 0},
  { label: "已操作", value: 1},
  { label: "已完成", value: 2},
  { label: "失败", value: 3},
];

// 产品类型
// 类型 0 活期 1 定期
export const productTypeOption = [
  { label: "活期理财", value: 0},
  { label: "定期理财", value: 1},
];


// 业务场景名称
//1： 购买定期 2 购买活期 3 赎回定期 4 赎回活期 5 收益复投 6 加速释放 7 兑换ju算力
export const businessNumberOptions = [
  { label: "购买定期", value: 1},
  { label: "购买活期", value: 2},
  { label: "赎回定期", value: 3},
  { label: "赎回活期", value: 4},
  { label: "收益复投", value: 5},
  { label: "加速释放", value: 6},
  { label: "兑换ju算力", value: 7},
];

// 任务状态
// :0-等待快照 1-发放中 2-发放完成 3-发放失败
export const syaLogStateOption = [
  { label: "等待快照", value: 0},
  { label: "发放中", value: 1},
  { label: "发放完成", value: 2},
  { label: "发放失败", value: 3},
];

// 转账清单-状态
// 状态:0 待审核 1 审核通过 2 完成转出/转入 3 审核不通过
export const ordertateOption = [
  { label: "待审核", value: 0},
  { label: "审核通过", value: 1},
  { label: "完成", value: 2},
  { label: "审核不通过", value: 3},
  { label: "失败", value: 4},
  { label: "失败", value: 5},
  { label: "失败", value: 6},
];


// 交易状态
// 转出状态:0-待转出 1-已转出 此字段不允许后台修改
export const transfertateOption = [
  { label: "待转出", value: 0},
  { label: "已转出", value: 1},
];

// 资金方向
//1 交易所转入 2 交易所转出 3 系统转入 4 系统转出 5 用户互转转入 6  用户互转转出
// 1 交易所转入(充值) 2 交易所转出（提现） 3 系统转入(分发) 4 系统转出（扣回） 5 用户互转转入(**你取个名吧) 6  用户互转转出(**你取个名吧)
export const directionOption = [
  { label: "充值", value: 1},
  { label: "提现", value: 2},
  { label: "分发", value: 3},
  { label: "扣回", value: 4},
  { label: "互转转入", value: 5},
  { label: "互转转出", value: 6},
];


// 状态
//  0 确认中 1-有效
export const currentProductState = [
  { label: "确认中", value: 0},
  { label: "有效", value: 1},
];

// 活期理财 - 状态
// 0：确认中 1：有效，2：已经期  3：已复投 4：已赎回 5：系统赎回 6：转移
export const currentState = [
  { label: "确认中", value: 0},
  { label: "有效", value: 1},
  { label: "已经期", value: 2},
  { label: "已复投", value: 3},
  { label: "已赎回", value: 4},
  { label: "系统赎回", value: 5},
  { label: "转移", value: 6},
];

// 定期理财 - 状态
// "1-有效 2 已过期待赎回 3 已赎回
export const regularState = [
  { label: "有效", value: 1},
  { label: "已过期待赎回", value: 2},
  { label: "已赎回", value: 3},
];
// 定期理财 类型
// 
export const regularProductOption = [
  { label: "30天理财", value: '30天理财'},
  { label: "180天理财", value: '180天理财'},
  { label: "360天理财", value: '360天理财'},
  { label: "360天保本理财", value: '360天保本理财'},
];

// 
// 0 用户手动购买 1 系统收益复投
export const buySystemOption = [
  { label: "用户手动购买", value: 0},
  { label: "系统收益复投", value: 1},
];


// 收益类型
// 1 静态 2 团队动态 3 平级 4 活期 
export const rewardTypeOption = [
  { label: "静态", value: 1},
  { label: "团队动态", value: 2},
  { label: "平级", value: 3},
  { label: "活期", value: 4},
];


// 发放类型
// 0 待发放 1 发放中 2 已完成
export const sendStateOption = [
  { label: "待发放", value: 0},
  { label: "发放中", value: 1},
  { label: "已完成", value: 2},
];
