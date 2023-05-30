/**
  @author: Zero
  @date: 2023/5/29 22:39:38
  @desc: Redis分布式锁单元测试

**/

package locks

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zlx2019/sugar/caches"
	"testing"
	"time"
)

// 测试分布式锁,非阻塞模式
func TestRedisDistributedLock_Lock_NotBlocking(t *testing.T) {
	// 创建一个Redis分布式锁-非阻塞模式
	lock := NewRedisDistributedLock("lock-key", caches.NewDefaultRedisTemplate())
	// 加锁
	err := lock.Lock(context.Background())
	if err == nil {
		//加锁成功
		err = lock.Unlock(context.Background())
		if err == nil {
			//解锁成功
		}
	}
}

// 测试Redis分布式锁,看门狗模式1 使用自定义有效时间,主动开启看门狗模式
func TestRedisDistributedLock_Lock_WatchDog1(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	lock := NewRedisDistributedLock("dog-lock-key", caches.NewDefaultRedisTemplate(),WithExpire(time.Second * 5),WithWatchDog())
	err := lock.Lock(context.Background())
	is.NoError(err)
	time.Sleep(time.Second * 10)
	is.NoError(lock.Unlock(context.Background()))

	time.Sleep(time.Second * 5)

	is.NoError(lock.Lock(context.Background()))
	time.Sleep(time.Second * 10)
	is.NoError(lock.Unlock(context.Background()))
}

// 测试Redis分布式锁,看门狗模式,使用默认有效时间10s
func TestRedisDistributedLock_Lock_WatchDog2(t *testing.T) {
	t.Parallel()
	is := assert.New(t)
	lock := NewRedisDistributedLock("dog-lock-key", caches.NewDefaultRedisTemplate())
	err := lock.Lock(context.Background())
	is.NoError(err)
	time.Sleep(time.Second * 5)
	is.NoError(lock.Unlock(context.Background()))
	time.Sleep(time.Second * 5)

	is.NoError(lock.Lock(context.Background()))
	time.Sleep(time.Second * 5)
	is.NoError(lock.Unlock(context.Background()))
	time.Sleep(time.Second * 5)
}


// 测试Redis分布式锁 阻塞模式
func TestRedisDistributedLock_Lock_Blocking1(t *testing.T) {
	t.Parallel()
	is := assert.New(t)
	// 创建Redis分布式锁,阻塞式
	var lock DistributedLock = NewRedisDistributedLock("lock2", caches.NewDefaultRedisTemplate(), WithBlocking(),WithExpire(time.Second * 10))
	// 加锁
	err := lock.Lock(context.Background())
	is.NoError(err)

	// 再尝试加锁
	is.Error(lock.Lock(context.Background()))

	// 解锁
	is.NoError(lock.Unlock(context.Background()))
	// 再加锁
	is.NoError(lock.Lock(context.Background()))
	is.NoError(lock.Unlock(context.Background()))
}

// 测试Redis分布式锁 阻塞模式, 使用重试次数+间隔时间
func TestRedisDistributedLock_Lock2(t *testing.T) {
	t.Parallel()
	is := assert.New(t)
	lock := NewRedisDistributedLock("lock3", caches.NewDefaultRedisTemplate(), WithBlocking(), WithExpire(time.Second*10), WithBlockingWaitTime(time.Second))
	err := lock.Lock(context.Background())
	is.NoError(err)
	err = lock.Lock(context.Background())
	if err != nil {
		fmt.Println(err.Error())
	}
}

func TestTimeProportion(t *testing.T) {
	//unit := time.Second
	//l := unit / 4
	//fmt.Println(l.Milliseconds())

	//s := time.Duration(float64(time.Second * 100) * 0.3)
	//fmt.Println(s.Milliseconds())

	fmt.Println(time.Duration(float64(time.Second * 1) * 0.25))
	fmt.Println(time.Duration(float64(time.Second * 1) * 0.3))
	fmt.Println(time.Duration(float64(time.Second * 1) * 0.75))
}