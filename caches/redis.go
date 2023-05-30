/**
  @author: Zero
  @date: 2023/4/22 17:05:42
  @desc: 缓存组件Redis API封装

**/

package caches

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/zlx2019/toys/converts"
	"log"
	"reflect"
	"time"
)

// 默认上下文,不设置超时时间
var defaultCtx = context.Background()

// RedisTemplate Redis缓存组件模板实现
type RedisTemplate struct {
	Client *redis.Client
}

// NewDefaultRedisTemplate 创建一个默认Redis客户端
func NewDefaultRedisTemplate() RedisTemplate {
	return RedisTemplate{
		Client: redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "root1234",
			DB:       0,
		}),
	}
}

// Set 设置一个缓存
func (template *RedisTemplate) Set(key string, value any) error {
	return template.SetExpire(key, value, 0)
}

// SetExpire 设置一个带有有效时间的缓存
func (template *RedisTemplate) SetExpire(key string, value any, expire time.Duration) error {
	var body any
	var err error
	// 通过反射断言类型
	switch reflect.TypeOf(value).Kind() {
	// Struct、Slice、Map等复杂结构自定义序列化为[]byte,避免没有实现BinaryMarshaler()而发生错误
	// 后续有更好的方案再度优化
	case reflect.Struct, reflect.Slice, reflect.Map:
		body, err = converts.ToBytes(value)
	default:
		body = value
	}
	if err != nil {
		return err
	}
	status := template.Client.Set(defaultCtx, key, body, expire)
	if status.Err() != nil {
		return status.Err()
	}
	return nil
}

// Get 根据Key读取一个缓存
func (template *RedisTemplate) Get(key string) *Reply {
	cmd := template.Client.Get(defaultCtx, key)
	//TODO  Nil表示Key不存在,严格来说并不是一种错误,暂不处理
	//if err := cmd.Err(); err != nil && err == redis.Nil {
	//
	//}
	return NewReply(cmd)
}

// Del 删除一个或多个缓存
func (template *RedisTemplate) Del(keys ...string) error {
	return template.Client.Del(defaultCtx, keys...).Err()
}

// Exists 检查一个缓存是否存在
func (template *RedisTemplate) Exists(key string) bool {
	// 当n为1时表示存在,表示不存在。
	n, err := template.Client.Exists(defaultCtx, key).Result()
	if err != nil {
		log.Printf("exists command error %s", err)
		return false
	}
	if n < 1 {
		return false
	}
	return true
}

// Keys 匹配所有符合规则的Key
func (template *RedisTemplate) Keys(pattern string) []string {
	keys, err := template.Client.Keys(defaultCtx, pattern).Result()
	if err != nil {
		return []string{}
	}
	return keys
}

// ExpireAdd 延长一个缓存的有效期
func (template *RedisTemplate) ExpireAdd(key string, time time.Duration) bool {
	ok, err := template.Client.Expire(defaultCtx, key, time).Result()
	if err != nil {
		return false
	}
	return ok
}

// ExpireSetup 设置有效期为指定时间
func (template *RedisTemplate) ExpireSetup(key string, time time.Time) bool {
	ok, err := template.Client.ExpireAt(defaultCtx, key, time).Result()
	if err != nil {
		return false
	}
	return ok
}

// GetExpire 获取一个Key的剩余有效期
func (template *RedisTemplate) GetExpire(key string) (time.Duration, error) {
	return template.Client.TTL(defaultCtx, key).Result()
}

// SetNEX 设置一个缓存,带有有效期,如果key已存在,则设置失败.
// 该方法持有原子性,通常用来做分布式锁
func (template *RedisTemplate) SetNEX(ctx context.Context, key, value string, expire time.Duration) error {
	if key == "" || value == "" {
		return errors.New("redis set key or value can't be empty")
	}
	// 执行Set命令
	// EX 表示使用秒作为时间单位 || PX 表示使用毫秒作为时间单位
	// NX 表示如果key已存在则设置失败
	cmd := template.Client.Do(ctx, "SET", key, value, "PX", expire.Milliseconds(), "NX")
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

// Eval 执行Lua脚本
// script Lua脚本内容
// keys 脚本Key列表
// args 脚本参数列表
func (template *RedisTemplate) Eval(ctx context.Context, script string, keys []string, args []any) (interface{},error) {
	return template.Client.Eval(ctx, script, keys, args...).Result()
}
