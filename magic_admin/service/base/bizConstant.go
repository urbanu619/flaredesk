package base

type BusinessType struct {
	BusinessNumber int    // 业务场景编号
	BusinessName   string // 场景名称
	ContextName    string //上下文名
	ContextValue   string // 上下文值
	Desc           string // 备注信息
}

func BsById(i int) *BusinessType {
	if v, ok := businessNumberMap[i]; ok {
		return &BusinessType{
			BusinessNumber: i,
			BusinessName:   v,
		}
	}
	return &BusinessType{
		BusinessNumber: UnknownNumber,
		BusinessName:   businessNumberMap[UnknownNumber],
	}
}

var BillBsMap = map[int]string{
	BsDeposit:                "充值",
	BsWithdraw:               "提现",
	BsStaticProfit:           "静态奖励",
	BsDirectProfit:           "直推奖励",
	BsInviteProfit:           "间推奖励",
	BsEqualProfit:            "平级奖励",
	BsTeamProfit:             "团队奖励",
	BsDaoProfit:              "DAO收益",
	BsFundOrganizationProfit: "基金会收益",
	BsLeaderProfit:           "DAO团队长收益",
	BsAgentProfit:            "DAO代理收益",
	BsPartnerProfit:          "DAO合伙人收益",
}

var businessNumberMap = map[int]string{
	UnknownNumber:  "未知",
	BsDeposit:      "充值",
	BsWithdraw:     "提现",
	BsDistribution: "分发",
	BsSysDeduction: "系统划扣",

	BsJoinGroup:   "参团",
	BsCancelGroup: "撤团",

	BsStaticProfit:           "静态奖励",
	BsDirectProfit:           "直推奖励",
	BsInviteProfit:           "间推奖励",
	BsEqualProfit:            "平级奖励",
	BsTeamProfit:             "团队奖励",
	BsDaoPurchase:            "DAO产品购买",
	BsDaoProfit:              "DAO收益",
	BsFundOrganizationProfit: "基金会收益",

	BsLeaderProfit:  "DAO团队长收益",
	BsAgentProfit:   "DAO代理收益",
	BsPartnerProfit: "DAO合伙人收益",
	BsFmcProducer:   "FMC产矿",
}

// 资金相关

const (
	UnknownNumber  int = iota + 100
	BsDeposit          // 充值
	BsWithdraw         // 提现
	BsDistribution     // 分发
	BsSysDeduction     // 系统划扣
)

// 参团

const (
	UnknownJoinGroupNumber int = iota + 200
	BsJoinGroup                // 参团
	BsCancelGroup              // 撤团
	BsGroupReturnCost          // 退回本金
)

// 参团

const (
	_                        int = iota + 300
	BsStaticProfit               // 静态奖励
	BsDirectProfit               // 直推奖励
	BsInviteProfit               // 间推奖励
	BsEqualProfit                // 平级奖励
	BsTeamProfit                 // 团队奖励
	BsDaoPurchase                // DAO产品购买
	BsDaoProfit                  // DAO收益
	BsFmcProducer                // FMC收益
	BsFundOrganizationProfit     // 基金会收益
	BsLeaderProfit
	BsAgentProfit
	BsPartnerProfit
)

var ProfitRecordBsMap = map[int]string{
	BsDirectProfit: "直推奖励",
	BsInviteProfit: "间推奖励",
	BsEqualProfit:  "平级奖励",
	BsTeamProfit:   "团队奖励",
}
