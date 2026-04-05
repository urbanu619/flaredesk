package scheduler

import (
	"fmt"
	model "go_server/model/biz_modules/app"
	"go_server/service/app"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"go_server/base/core"
	"gorm.io/gorm"
)

// 秒级 cron 解析器（与 cron.WithSeconds() 一致）
var secondParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

// RedPacketScheduler 红包定时任务调度器
type RedPacketScheduler struct {
	cron      *cron.Cron
	db        *gorm.DB
	jobs      map[int64]cron.EntryID // configId -> entryID
	jobCrons  map[int64]string       // configId -> cronExpr（用于检测变化）
	jobsMutex sync.RWMutex
	service   *app.TgRedPacketSendService
}

var (
	scheduler *RedPacketScheduler
	once      sync.Once
)

// GetScheduler 获取调度器单例
func GetScheduler(db *gorm.DB) *RedPacketScheduler {
	once.Do(func() {
		scheduler = &RedPacketScheduler{
			cron:     cron.New(cron.WithSeconds()), // 支持秒级
			db:       db,
			jobs:     make(map[int64]cron.EntryID),
			jobCrons: make(map[int64]string),
			service:  &app.TgRedPacketSendService{},
		}
	})
	return scheduler
}

// Start 启动调度器
func (s *RedPacketScheduler) Start() error {
	core.Log.Info("🕐 启动红包定时任务调度器...")

	// 加载所有启用的定时配置
	if err := s.LoadAllConfigs(); err != nil {
		return fmt.Errorf("加载定时配置失败: %v", err)
	}

	// 启动 cron
	s.cron.Start()

	// 启动配置监听器（每分钟检查一次配置变化）
	go s.watchConfigChanges()

	core.Log.Info("✅ 红包定时任务调度器启动成功")
	return nil
}

// Stop 停止调度器
func (s *RedPacketScheduler) Stop() {
	core.Log.Info("🛑 停止红包定时任务调度器...")
	s.cron.Stop()
}

// LoadAllConfigs 加载所有启用的定时配置
func (s *RedPacketScheduler) LoadAllConfigs() error {
	var configs []model.TgRedPacketConfig

	// 查询所有启用的定时红包配置
	err := s.db.Where("config_type = ? AND status = ?", 1, 1).Find(&configs).Error
	if err != nil {
		return err
	}

	core.Log.Infof("📋 加载到 %d 个定时红包配置", len(configs))

	// 添加到调度器
	for _, config := range configs {
		if err := s.AddJob(&config); err != nil {
			core.Log.Errorf("添加定时任务失败 [ID:%d]: %v", config.Id, err)
			continue
		}
	}

	return nil
}

// AddJob 添加定时任务
func (s *RedPacketScheduler) AddJob(config *model.TgRedPacketConfig) error {
	s.jobsMutex.Lock()
	defer s.jobsMutex.Unlock()

	return s.addJobLocked(config)
}

// addJobLocked 添加定时任务（需要提前加锁）
func (s *RedPacketScheduler) addJobLocked(config *model.TgRedPacketConfig) error {
	// 检查是否已存在
	if _, exists := s.jobs[config.Id]; exists {
		s.removeJobLocked(config.Id)
	}

	// 验证 Cron 表达式
	if config.CronExpr == "" {
		return fmt.Errorf("Cron 表达式为空")
	}

	// 添加任务
	entryID, err := s.cron.AddFunc(config.CronExpr, func() {
		s.executeJob(config.Id)
	})

	if err != nil {
		return fmt.Errorf("添加 Cron 任务失败: %v", err)
	}

	s.jobs[config.Id] = entryID
	s.jobCrons[config.Id] = config.CronExpr

	core.Log.Infof("✅ 添加定时任务 [ID:%d, Name:%s, Cron:%s]",
		config.Id, config.ConfigName, config.CronExpr)

	return nil
}

// RemoveJob 移除定时任务
func (s *RedPacketScheduler) RemoveJob(configId int64) {
	s.jobsMutex.Lock()
	defer s.jobsMutex.Unlock()
	s.removeJobLocked(configId)
}

// removeJobLocked 移除定时任务（需要提前加锁）
func (s *RedPacketScheduler) removeJobLocked(configId int64) {
	if entryID, exists := s.jobs[configId]; exists {
		s.cron.Remove(entryID)
		delete(s.jobs, configId)
		delete(s.jobCrons, configId)
		core.Log.Infof("🗑️ 移除定时任务 [ID:%d]", configId)
	}
}

// executeJob 执行定时任务
func (s *RedPacketScheduler) executeJob(configId int64) {
	core.Log.Infof("⏰ 执行定时红包任务 [ID:%d]", configId)

	// 查询配置
	var config model.TgRedPacketConfig
	if err := s.db.Where("id = ?", configId).First(&config).Error; err != nil {
		core.Log.Errorf("查询配置失败 [ID:%d]: %v", configId, err)
		return
	}

	// 检查状态
	if config.Status != 1 {
		core.Log.Warnf("配置已禁用，跳过执行 [ID:%d]", configId)
		return
	}

	// 使用配置的时区做时间比较
	now := time.Now()
	if config.TimeZone != "" {
		if loc, err := time.LoadLocation(config.TimeZone); err == nil {
			now = now.In(loc)
		}
	}

	if config.StartDate != nil && !config.StartDate.IsZero() && now.Before(*config.StartDate) {
		core.Log.Infof("未到开始时间，跳过执行 [ID:%d, StartDate:%v, Now:%v]", configId, config.StartDate, now)
		return
	}
	if config.EndDate != nil && !config.EndDate.IsZero() && now.After(*config.EndDate) {
		core.Log.Infof("已过结束时间，停止任务 [ID:%d, EndDate:%v, Now:%v]", configId, config.EndDate, now)
		s.RemoveJob(configId)
		return
	}

	// 调用发送服务
	s.service.SetDbAlias("app")
	packetNo, err := s.service.ExecuteSendRedPacket(&config)
	if err != nil {
		core.Log.Errorf("发送红包失败 [ID:%d]: %v", configId, err)
		return
	}

	// 更新执行统计
	updates := map[string]interface{}{
		"last_exec_time": time.Now(),
		"exec_count":     config.ExecCount + 1,
	}

	// 计算下次执行时间（使用秒级解析器，与 cron.WithSeconds() 一致）
	schedule, err := secondParser.Parse(config.CronExpr)
	if err == nil {
		nextTime := schedule.Next(time.Now())
		updates["next_exec_time"] = nextTime
	} else {
		core.Log.Warnf("解析 Cron 表达式失败 [ID:%d, Cron:%s]: %v", configId, config.CronExpr, err)
	}

	s.db.Model(&model.TgRedPacketConfig{}).Where("id = ?", configId).Updates(updates)

	core.Log.Infof("✅ 定时红包发送成功 [ID:%d, PacketNo:%s]", configId, packetNo)
}

// watchConfigChanges 监听配置变化
func (s *RedPacketScheduler) watchConfigChanges() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.reloadConfigs()
	}
}

// reloadConfigs 重新加载配置
func (s *RedPacketScheduler) reloadConfigs() {
	var configs []model.TgRedPacketConfig

	// 查询所有启用的定时配置
	err := s.db.Where("config_type = ? AND status = ?", 1, 1).Find(&configs).Error
	if err != nil {
		core.Log.Errorf("重新加载配置失败: %v", err)
		return
	}

	s.jobsMutex.Lock()
	defer s.jobsMutex.Unlock()

	// 构建当前应该存在的任务 ID 集合
	activeIds := make(map[int64]bool)
	for _, config := range configs {
		activeIds[config.Id] = true

		existingCron, exists := s.jobCrons[config.Id]
		if !exists {
			// 任务不存在，添加
			if err := s.addJobLocked(&config); err != nil {
				core.Log.Errorf("重新加载-添加定时任务失败 [ID:%d]: %v", config.Id, err)
			}
		} else if existingCron != config.CronExpr {
			// Cron 表达式变了，重新注册
			core.Log.Infof("🔄 检测到 Cron 变化 [ID:%d]: %s -> %s", config.Id, existingCron, config.CronExpr)
			if err := s.addJobLocked(&config); err != nil {
				core.Log.Errorf("重新加载-更新定时任务失败 [ID:%d]: %v", config.Id, err)
			}
		}
	}

	// 移除不再需要的任务
	for id := range s.jobs {
		if !activeIds[id] {
			s.removeJobLocked(id)
		}
	}
}
