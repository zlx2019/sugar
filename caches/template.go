/**
  @author: Zero
  @date: 2023/4/22 16:43:03
  @desc: 缓存组件抽象模板

**/

package caches

import (
	"time"
)

// CacheTemplate 缓存组件顶级接口
type CacheTemplate interface {
	//Set 设置一个缓存
	Set(key string, value any) error
	// SetExpire 设置一个带有有效期的缓存
	SetExpire(key string, value any, expire time.Duration) error
	// Get 获取一个缓存
	Get(key string) *Reply
	// Del 删除一个或者多个缓存
	Del(keys ...string) error
	// Exists 检查一个缓存是否存在
	Exists(key string) bool
	// Keys 匹配所有符合规则的Key
	Keys(pattern string) []string
	// ExpireAdd  延长有效期
	ExpireAdd(key string, time time.Duration) bool
	// ExpireSetup 设置有效期为指定时间
	ExpireSetup(key string, time time.Time) bool
	// GetExpire 获取一个Key的剩余有效期
	GetExpire(key string) (time.Duration, error)
}
