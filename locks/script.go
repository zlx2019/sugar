/**
  @author: Zero
  @date: 2023/5/30 12:32:01
  @desc: lua脚本常量

**/

package locks

// UnlockLuaScript 用于释放锁的Lua脚本命令
// 如果`key`的value和指定的参数值相等才删除这个`key`
const UnlockLuaScript = `
	if redis.call('get',KEYS[1]) == ARGV[1] then
		return redis.call('del',KEYS[1])
	else
		return 0
	end
`

// LockExpireDelayScript 用于为锁有效期续约的Lua脚本命令
// 当`key`的value和`token`相等 && key的剩余有效期小于`triggerVal`时才进行续约
// 如果当前剩余有效期大于`triggerVal`和续约成功都返回`1`。key不存在或者token不一致则返回`0`
const LockExpireDelayScript = `
	local key = KEYS[1]
	local token = ARGV[1]
	local value = redis.call('get',key)
	if(not value or value ~= token) then
		return 0
	else
		local triggerVal = ARGV[2]
		if redis.call('pttl',key) > tonumber(triggerVal) then
			return 1
		else
			local incrVal = ARGV[3]
			return redis.call('pexpire',key,incrVal)
		end
	end
`