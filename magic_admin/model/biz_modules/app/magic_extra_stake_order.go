package app

// 引入关联包

type MagicExtraStakeOrder struct { 
}

func (*MagicExtraStakeOrder) TableName() string {
	return "magic_extra_stake_order"
}

func NewMagicExtraStakeOrder() *MagicExtraStakeOrder {
	return &MagicExtraStakeOrder{}
}
