/**
  @author: Zero
  @date: 2023/5/29 21:39:55
  @desc: 锁的配置选项

**/

package locks

import (
	"context"
	"time"
)

const (
	// 锁的默认有效时长为3秒
	defaultExpire = time.Second * 10
	// 默认加锁阻塞等待时长为3秒
	defaultBlockingTime = time.Second * 3
	// 默认加锁重试次数为5次
	defaultRetry = 5

	// 看门狗状态
	stop = 0 // 停止
	running = 1 //运行
)

// LockOptions 分布式锁功能配置选项
type LockOptions struct {
	// 锁的有效时长,默认为3s
	expire time.Duration
	// 看门狗
	WatchDog

	// 锁类型是否为阻塞式模式(自旋+重试)
	// 默认为非阻塞模式(成功或失败直接返回)
	blocking bool
	// 阻塞模式下阻塞的时长,默认为3s
	blockingTime time.Duration

	// 阻塞模式下重试获取锁的次数,默认为5
	retry int
	// 每次重试间隔等待时间,默认为 阻塞时长/重试次数 (blockingTime / retry)
	retryWaitingTime time.Duration
}

// WatchDog 看门狗,为锁的有效期自动续约
// 规则: 每%25的时间比例检查一次是否续约,不足%30则续约
// 每次续约%75的时间比例
type WatchDog struct {
	// 是否启用看门狗: 如果用户没有显示为锁添加有效期,那么就启动看门狗
	enabled bool
	// 运行状态,用于标识看门狗是否正在运行
	status int32
	// 用于关闭看门狗的函数
	cancelFn context.CancelFunc
}

// LockOption 选项闭包
type LockOption func(options *LockOptions)

// WithBlocking 设置锁为阻塞模式
func WithBlocking() LockOption {
	return func(options *LockOptions) {
		options.blocking = true
	}
}

// WithBlockingWaitTime 设置锁的阻塞等待时长
func WithBlockingWaitTime(blockingTime time.Duration) LockOption {
	return func(options *LockOptions) {
		options.blockingTime = blockingTime
	}
}

// WithExpire 设置锁的有效时长
func WithExpire(expire time.Duration) LockOption {
	return func(options *LockOptions) {
		options.expire = expire
	}
}

// WithRetry 设置锁自旋的次数
func WithRetry(retry int) LockOption {
	return func(options *LockOptions) {
		options.retry = retry
	}
}

// WithRetryWaitingTime 设置锁自旋间隔时间
func WithRetryWaitingTime(waitingTime time.Duration) LockOption {
	return func(options *LockOptions) {
		options.retryWaitingTime = waitingTime
	}
}

// WithWatchDog 启用锁的有效期自动续约
func WithWatchDog()LockOption {
	return func(options *LockOptions) {
		options.enabled = true
	}
}

// 设置默认选项参数
// 如果没有设置某项参数,则使用默认参数
func optionWithDefault(options *LockOptions) {
	// 设置锁的默认有效期,并且启用自动续约(看门狗)
	if options.expire <= 0 {
		options.expire = defaultExpire
		options.enabled = true
	}
	// 阻塞模式,设置加锁的默认阻塞等待时长
	if options.blocking && options.blockingTime <= 0 {
		options.blockingTime = defaultBlockingTime
	}
	// 阻塞模式,设置加锁的默认重试次数
	if options.blocking && options.retry <= 0 {
		options.retry = defaultRetry
	}
	// 阻塞模式,设置默认的重试间隔时间
	if options.blocking && options.retryWaitingTime <= 0 {
		// 间隔时间 = 总时长 / 重试次数
		sleep := options.blockingTime.Milliseconds() / int64(options.retry)
		// 设置间隔时间,以毫秒为单位
		options.retryWaitingTime = time.Millisecond * time.Duration(sleep)
	}
}
