package base

// 基础信息 映射
//
//var userMap map[int64]*model.ProUser
//var userMapExp int64
//
//func (s *BizCommonService) CacheUserById(id int64) (*model.ProUser, error) {
//	if userMapExp < time.Now().Unix() {
//		if err := s.initUserMap(); err != nil {
//			return nil, err
//		}
//	}
//	v, ok := userMap[id]
//	if ok {
//		return v, nil
//	}
//	return nil, fmt.Errorf("not found")
//}
//
//func (s *BizCommonService) initUserMap() error {
//	userMap = make(map[int64]*model.ProUser)
//	allUser, err := GetMore[model.ProUser](s.DB().Where("1=1"))
//	if err != nil {
//		return err
//	}
//	core.Log.Infof("用户总数 初始化 :%d", len(allUser))
//	for _, user := range allUser {
//		userMap[user.Id] = user
//	}
//	userMapExp = time.Now().Unix() + 60
//	return nil
//}
//
//// 用户指标信息 映射
//var userDaoQuotaMap map[int64]*model.ProUserDaoQuota
//var userDaoQuotaMapExp int64
//
//func (s *BizCommonService) CacheUserDaoQuotaById(id int64) (*model.ProUserDaoQuota, error) {
//	if userDaoQuotaMapExp < time.Now().Unix() {
//		if err := s.initUserDaoQuotaMap(); err != nil {
//			return nil, err
//		}
//	}
//	v, ok := userDaoQuotaMap[id]
//	if ok {
//		return v, nil
//	}
//	return nil, fmt.Errorf("not found")
//}
//
//func (s *BizCommonService) initUserDaoQuotaMap() error {
//	userDaoQuotaMap = make(map[int64]*model.ProUserDaoQuota)
//	items, err := GetMore[model.ProUserDaoQuota](s.DB().Where("1=1"))
//	if err != nil {
//		return err
//	}
//	for _, item := range items {
//		userDaoQuotaMap[item.UserId] = item
//	}
//	userDaoQuotaMapExp = time.Now().Unix() + 60
//	return nil
//}
