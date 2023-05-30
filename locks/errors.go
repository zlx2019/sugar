/**
  @author: Zero
  @date: 2023/5/30 12:32:39
  @desc: 错误列表

**/

package locks

import "errors"

var (
	// LockBlockingTimeOutErr 阻塞模式获取锁超时错误
	LockBlockingTimeOutErr = errors.New("lock failed blocking timeout")
	// LockNotRetryErr 阻塞模式自旋重试加锁次数已用尽错误
	LockNotRetryErr = errors.New("lock failed retry used up")

	// UnlockWithoutOwnershipErr 释放一把对自己无所有权的锁从而产生的错误
	UnlockWithoutOwnershipErr = errors.New("unlock failed without ownership")
	// DelayLockWithoutOwnershipErr 对一个没有所有权的锁续约从而产生的错误
	DelayLockWithoutOwnershipErr = errors.New("delay lock failed without ownership")
)