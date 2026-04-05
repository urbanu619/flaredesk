package core

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

// 限流器

type RateLimiter struct {
	client redis.UniversalClient
}

func NewRateLimiter(client redis.UniversalClient) *RateLimiter {
	return &RateLimiter{client: client}
}

func (r *RateLimiter) IsAllowedSimple(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	// 先检查 key 的类型，如果是错误类型则删除
	keyType, err := r.client.Type(ctx, key).Result()
	if err != nil {
		return false, err
	}
	// 如果 key 存在且不是 string 类型，删除它
	if keyType != "none" && keyType != "string" {
		r.client.Del(ctx, key)
	}
	// 使用 INCR 和 EXPIRE 的简化方案
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if count == 1 {
		// 第一次设置时添加过期时间
		r.client.Expire(ctx, key, window)
	}
	return count <= int64(limit), nil
}

// 检查是否可以执行方法

func (r *RateLimiter) CanExecuteMethod(ctx context.Context, methodName string, info string) (bool, error) {
	key := fmt.Sprintf("rate_limit:%s:%s", methodName, info)
	return r.IsAllowedSimple(ctx, key, 10, 10*time.Minute)
}

// ClearLimit 清除指定 key 的限流缓存
func (r *RateLimiter) ClearLimit(ctx context.Context, methodName string, info string) error {
	key := fmt.Sprintf("rate_limit:%s:%s", methodName, info)
	return r.client.Del(ctx, key).Err()
}

// ClearLimitByPattern 根据模式匹配清除多个限流缓存
func (r *RateLimiter) ClearLimitByPattern(ctx context.Context, pattern string) error {
	// 使用 SCAN 命令避免在大量 key 时阻塞 Redis
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}
	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}
	return nil
}

func (r *RateLimiter) ClearMethodLimit(ctx context.Context) error {
	key := fmt.Sprintf("rate_limit:*")
	return r.ClearLimitByPattern(ctx, key)
}
