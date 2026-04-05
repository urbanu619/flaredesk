package app

import (
	"github.com/gin-gonic/gin"
	"go_server/global"
)

// All routes to be registered

var (
	allRouters = []global.ContextInterface{MagicAssetRouter{}, MagicAssetBillRouter{}, MagicAssetRwCallbackLogRouter{}, MagicAssetRwRecordRouter{}, MagicStakeProductRouter{}, MagicStakeQueueInfoRouter{}, MagicStakeUserCurrentOpsRecordRouter{}, MagicStakeUserCurrentOrderRouter{}, MagicUserRouter{}, MagicUserActionLogRouter{}, MagicUserLogsRouter{}, MagicUserProfitRouter{}, MagicUserProfitRecordRouter{}, MagicUserQuotaRouter{}, NodeBannerRouter{}, NodeInfoRouter{}, NodeOrderRouter{}, NodeOrderPaymentsRouter{}, PromotionUpgradeLevelRouter{}, SysCoinRouter{}, SysI18nRouter{}, SysJobRouter{}, SysLevelConfigRouter{}, SysSignConfigRouter{}, SysStakePeriodJobRouter{}, MagicExtraStakeOrderOpsRecordRouter{}, MagicExtraStakeOrderRouter{}, TgRedPacketConfigRouter{}, TgRedPacketSendRouter{}, TgRedPacketRecordRouter{}, TgGroupRouter{}, TgUserBindRouter{}, CfAccountRouter{}, CfDnsRouter{}, CfDnsTemplateRouter{}, CfZoneRouter{}}
)

type RouterGroup struct {
}

func (RouterGroup) Route() string {
	return "/app"
}

func (h RouterGroup) Register(group *gin.RouterGroup) {
	for _, item := range allRouters {
		global.RegisterRouter(group, item)
	}
}
