package app

// 引入关联包

type MagicExtraStakeOrderOpsRecord struct { 
}

func (*MagicExtraStakeOrderOpsRecord) TableName() string {
	return "magic_extra_stake_order_ops_record"
}

func NewMagicExtraStakeOrderOpsRecord() *MagicExtraStakeOrderOpsRecord {
	return &MagicExtraStakeOrderOpsRecord{}
}
