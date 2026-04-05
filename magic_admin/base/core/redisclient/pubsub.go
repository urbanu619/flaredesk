package redisclient

import (
	"context"
	"encoding/json"
	"go_server/base/core"
	"time"
)

// Publish 发布消息
func Publish(channel string, message interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := DefaultClient().Publish(ctx, channel, ObjToJson(message)).Result()
	return err
}

func ObjToJson(obj interface{}) string {
	if obj == nil {
		return ""
	}
	b, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(b)
}

// Subscribe 订阅消息（支持 Redis 断线重连）
func Subscribe(channel string, handler func(data string) error) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		pubsub := DefaultClient().Subscribe(ctx, channel)
		ch := pubsub.Channel()

		core.Log.Info("Redis 订阅成功, channel: ", channel)

		// 监听消息
		for {
			select {
			case <-ctx.Done(): // 监听外部取消信号，安全退出
				_ = pubsub.Close() // 确保关闭 Redis 订阅连接
				core.Log.Warn("Redis 订阅退出, channel: ", channel)
				return nil

			case msg, ok := <-ch: // 监听 Redis 消息
				if !ok {
					core.Log.Warn("Redis 订阅通道关闭, 5秒后重连, channel: ", channel)
					_ = pubsub.Close()          // 关闭当前连接，防止遗留连接
					time.Sleep(5 * time.Second) // 休眠后重新订阅
					break                       // 退出当前 for 循环，进入重连
				}
				// 处理收到的消息
				if err := handler(msg.Payload); err != nil {
					core.Log.Error("消息处理失败, channel: ", channel, ", error: ", err)
				} else {
					core.Log.Debug("收到 Redis 消息, channel: ", channel, ", message: ", msg.Payload)
				}
			}
		}
	}
}
