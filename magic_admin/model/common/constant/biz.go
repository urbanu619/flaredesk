package constant

// 资金变化场景

// 公共场景 -- 必填

type BusinessType struct {
	BusinessNumber int    // 业务场景编号
	BusinessName   string // 场景名称
	ContextName    string //上下文名
	ContextValue   string // 上下文值
	Desc           string // 备注信息
}

var defaultBs = newBusinessType(UnknownNumber, BusinessNumberMap[UnknownNumber])

func BsById(i int) *BusinessType {
	if v, ok := BusinessNumberMap[i]; ok {
		return &BusinessType{
			BusinessNumber: i,
			BusinessName:   v,
		}
	}
	return defaultBs
}

func newBusinessType(businessNumber int, businessName string) *BusinessType {
	return &BusinessType{
		BusinessNumber: businessNumber,
		BusinessName:   businessName,
	}
}

// 产品业务场景

const (
	ProductNumber     int = iota
	BsProductPledge       // 质押
	BsProductReply        // 一键复投
	BsProductTransfer     // 产品转移
)

// 转账业务场景编号

const (
	TransferNumber            int = iota + 100
	BsTransferExchangeCashIn      // 交易所转入
	BsTransferExchangeCashOut     // 交易所转出
	BsTransferUserCashIn          // 用户转入
	BsTransferUserCashOut         // 用户转出
	BsTransferSystemCashIn        // 系统转入
	BsTransferSystemCashOut       // 系统转出
)

// 奖励业务场景

const (
	RewardNumber                   int = iota + 200
	BsRewardPersonProfit               // 个人收益
	BsRewardDirectProfit               // 直推收益
	BsRewardTeamInviteProfit           // 团队收益
	BsRewardCommunityLevelProfit       // 社区收益
	BsRewardWeightedDividendProfit     // 加权分红
	BsRewardPromoterProfit             // 加速释放
	BsRewardActiveProfit               // 活动奖励
	BsRewardLotteryProfit              // 乐透奖励

)

func GetBizName(i int) string {
	if v, ok := BusinessNumberMap[i]; ok {
		return v
	}
	return BusinessNumberMap[UnknownNumber]
}

const UnknownNumber = 0

var BusinessNumberMap = map[int]string{
	UnknownNumber:                  "未知",
	BsProductPledge:                "质押",   // 质押
	BsProductReply:                 "一键复投", // 一键复投
	BsProductTransfer:              "遗矿转移",
	BsTransferExchangeCashIn:       "交易所转入", // 交易所转入
	BsTransferExchangeCashOut:      "交易所转出", // 交易所转出
	BsTransferUserCashIn:           "用户转入",  // 用户转入
	BsTransferUserCashOut:          "用户转出",  // 用户转出
	BsTransferSystemCashIn:         "系统转入",  // 系统转入
	BsTransferSystemCashOut:        "系统转出",  // 系统转出
	BsRewardPersonProfit:           "个人收益",  //   个人收益
	BsRewardDirectProfit:           "直推奖励",  // 直推奖励
	BsRewardTeamInviteProfit:       "团队收益",  //   团队推荐收益
	BsRewardCommunityLevelProfit:   "社区收益",  //   社区收益
	BsRewardWeightedDividendProfit: "加权分红",  // 加权分红
	BsRewardPromoterProfit:         "加速释放",  // 加速释放
	BsRewardActiveProfit:           "活动奖励",
	BsRewardLotteryProfit:          "聚乐透奖金",
}
