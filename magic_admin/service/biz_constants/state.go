package biz_constants

// 通用状态
// 使用场景 充值、提现
// 订单 任务状态 ===
const (
	StateUnknown     = "unknown"   // 未知状态
	StateWaiting     = "waiting"   // 待处理
	StateRunning     = "running"   // 进行中
	StateSuccess     = "success"   // 成功
	StateFailure     = "failure"   // 失败
	StateFinishQuota = "quota"     // 完成指标计算
	StateCancel      = "cancel"    // 取消
	StateSettled     = "settled"   // 结算完成
	StateEffective   = "effective" // 申购股票质押生效中
)

// 充提状态

var ClaimDepositStateMap = map[string]string{
	StateWaiting: "待处理",
	StateRunning: "处理中",
	StateSuccess: "成功",
	StateFailure: "失败",
}
