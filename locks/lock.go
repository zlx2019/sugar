/**
  @author: Zero
  @date: 2023/5/29 20:52:39
  @desc: 分布式锁抽象接口

**/

package locks

import "context"

// DistributedLock 顶级分布式锁接口
type DistributedLock interface {
	// Lock 加锁
	Lock(ctx context.Context) error
	// Unlock 释放锁
	Unlock(ctx context.Context) error
}
