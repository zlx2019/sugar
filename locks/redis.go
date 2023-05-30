/**
  @author: Zero
  @date: 2023/5/29 20:54:41
  @desc: Redis分布式实现

**/

package locks

import (
	"context"
	"fmt"
	"github.com/zlx2019/toys/randoms"
	"pluto/caches"
	"sync/atomic"
	"time"
)


// RedisDistributedLock 基于Redis的分布式锁
type RedisDistributedLock struct {
	// 锁的功能配置
	LockOptions
	// Redis客户端
	template caches.RedisTemplate
	// 锁的Key
	key string
	// 锁的身份标识,用于防止锁被他人释放
	// 使用获取锁者的 进程ID + 协程ID作为标识 或者使用UUID或其他随机值都可以
	token string
}


// NewRedisDistributedLock 创建一个Redis分布式锁
func NewRedisDistributedLock(key string, template caches.RedisTemplate, opts ...LockOption) *RedisDistributedLock {
	// 创建锁
	lock := RedisDistributedLock{
		key:      key,
		token:    randoms.RandomString(15),
		template: template,
	}
	// 设置锁的配置选项
	for _, opt := range opts {
		opt(&lock.LockOptions)
	}
	// 检查选项,填补默认参数
	optionWithDefault(&lock.LockOptions)
	return &lock
}

// Lock 加锁
func (lock *RedisDistributedLock) Lock(ctx context.Context)(err error) {
	// 锁的续约处理
	defer func() {
		// 如果加锁失败或者没有开启看门狗 直接退出
		if err != nil || !lock.enabled {
			return
		}
		// 获取锁后,初始化看门狗状态等
		lock.doWatchDog(ctx)
	}()

	// 无论阻塞与非阻塞模式,都要先加一次锁
	err = lock.tryLock(ctx)
	if err == nil {
		// 加锁成功
		return nil
	}
	// ======加锁失败处理=======

	// 非阻塞模式直接返回error
	if !lock.blocking {
		return err
	}
	// 阻塞模式继续尝试加锁(自旋+重试)
	return lock.loopTryLock(ctx)
}

// 尝试加锁,如果加锁失败则返回error
func (lock *RedisDistributedLock) tryLock(ctx context.Context) error {
	// 加锁
	return lock.template.SetNEX(ctx, lock.key, lock.token, lock.expire)
}

// 循环尝试加锁,直到阻塞时长用尽、可重试次数用尽、context中断。
func (lock *RedisDistributedLock) loopTryLock(ctx context.Context) error {
	// 超时通知器  如果超过锁的 `blockingTime`时长还未抢抢到锁,则表示获取锁超时
	timeOutChan := time.After(lock.blockingTime)
	// 轮询定时器 每隔锁的`retryWaitingTime`时长尝试加锁一次,直到`retry`次数用尽
	loopTicker := time.NewTicker(lock.retryWaitingTime)
	defer loopTicker.Stop()

	// 开始循环获取锁
	for range loopTicker.C {
		select {
		case <-ctx.Done():
			// 整个上下文终止
			return fmt.Errorf("lock failed ctx timeout, err: %w", ctx.Err())
		case <-timeOutChan:
			// 阻塞等待到达上限时间
			return LockBlockingTimeOutErr
		default:
			//放行,继续尝试加锁
		}
		err := lock.tryLock(ctx)
		if err == nil {
			// 加锁成功
			return nil
		}
		// 递减可重试次数
		if lock.retry -= 1; lock.retry <= 0 {
			// 已经没有可重试的次数
			return LockNotRetryErr
		}
	}
	// 不可达
	return nil
}

// 初始化看门狗运行状态
func (lock *RedisDistributedLock) doWatchDog(ctx context.Context) {
	// 更新状态: 将status从`0`更改为`1`;从静止状态更新为运行状态
	// 这里主要是确保之前开启的看门狗已经被停止了
	for !atomic.CompareAndSwapInt32(&lock.status, stop, running) {
	}
	fmt.Println("watch dog status stop to running")
	// 获取看门狗的停止函数
	ctx, lock.cancelFn = context.WithCancel(ctx)
	// 启动看门狗异步任务
	go func() {
		// 在任务结束之前,将状态恢复到停止状态
		defer func() {
			atomic.StoreInt32(&lock.status,stop)
			fmt.Println("watch dog status to stop")
		}()
		lock.runningWatchDog(ctx)
	}()
}

// 运行续约异步任务
func (lock *RedisDistributedLock)runningWatchDog(ctx context.Context)  {
	// 轮询间隔时间为 锁的有效期时长比例的%25
	intervalTime := time.Duration(float64(lock.expire) * 0.25)
	// 有效期不足%30比例时则续约
	triggerTime := time.Duration(float64(lock.expire) * 0.3)
	// 续约原有的%75的时间比例
	incrTime := time.Duration(float64(lock.expire) * 0.75)
	fmt.Printf("watch dog  intervalTime: %v, triggerTime: %v, incrTime: %v  \n",intervalTime.Milliseconds(),triggerTime.Milliseconds(),incrTime.Milliseconds())
	// 根据间隔时间创建一个定时器
	loop := time.NewTicker(intervalTime)
	defer loop.Stop()
	for range loop.C{
		select {
		case <-ctx.Done():
			// 已释放锁,停止任务
			return
		default:
		}
		// 执行续约
		_ = lock.delayExpire(ctx, triggerTime, incrTime)
	}
}

// 为锁的有效期执行续约操作,前提是剩余有效期不足%30,并且确保该锁属于自己
// triggerTime: 触发续约的阈值毫秒数,当有效期低于该数值才会续约
// incrTime: 要续约的毫秒数
func (lock *RedisDistributedLock) delayExpire(ctx context.Context,triggerTime,incrTime time.Duration) error{
	// 在Redis中 最小时间单位为毫秒
	trigger := triggerTime.Milliseconds()
	incr := incrTime.Milliseconds()
	// 通过lua脚本实现原子性续约
	// 返回`1`表示续约成功或者剩余期限还很多不需要续约
	val, err := lock.template.Eval(ctx, LockExpireDelayScript, []string{lock.key}, []any{lock.token, trigger, incr})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if v,ok :=val.(int64); ok && v != 1{
		// 不能对非自己的锁进行续约
		return DelayLockWithoutOwnershipErr
	}
	return nil
}


// Unlock 释放锁
// 释放锁的同时需要确认释放者的身份,所以基于lua脚本来实现操作的原子性
func (lock *RedisDistributedLock) Unlock(ctx context.Context) (err error) {
	// 看门狗处理
	defer func() {
		// 解锁失败或者没有开启看门狗 不作处理
		if err != nil || !lock.enabled {
			return
		}
		// 停止看门狗
		if lock.cancelFn != nil {
			lock.cancelFn()
		}
	}()
	// 释放锁
	val, err := lock.template.Eval(ctx, UnlockLuaScript, []string{lock.key}, []any{lock.token})
	if err != nil {
		return err
	}
	if v,ok :=val.(int64);ok && v != 1 {
		// 如果返回值不为`1`也视为失败。可能没有锁的释放权\锁已经提前失效
		return UnlockWithoutOwnershipErr
	}
	return nil
}
