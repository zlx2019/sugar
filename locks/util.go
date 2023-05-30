/**
  @author: Zero
  @date: 2023/5/30 11:44:21
  @desc: 工具库

**/

package locks

import (
	"fmt"
	"github.com/zlx2019/toys/system"
)

// GetLockToken 获取锁的Token标识,使用进程ID + 协程ID
func GetLockToken() string{
	return fmt.Sprintf("%s_%s",system.GetCurrentProcessID(),system.GetCurrentGoroutineID())
}
