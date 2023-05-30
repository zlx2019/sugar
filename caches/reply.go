/**
  @author: Zero
  @date: 2023/4/23 13:22:13
  @desc: 操作响应结果体,通常用于获取缓存

**/

package caches

import (
	"github.com/redis/go-redis/v9"
	"github.com/zlx2019/toys/converts"
)

// Reply 操作响应
type Reply struct {
	// Redis操作响应对象
	cmd *redis.StringCmd
	err error //错误响应
	ok  bool  //操作是否成功

}

func NewReply(cmd *redis.StringCmd) *Reply {
	return &Reply{
		cmd: cmd,
		err: cmd.Err(),
		ok:  cmd.Err() == nil,
	}
}

// Err 获取响应错误
func (reply *Reply) Err() error {
	return reply.err
}

// Ok 获取响应状态
func (reply *Reply) Ok() bool {
	return reply.ok
}

// GetValue  获取响应结果
func (reply *Reply) GetValue() interface{} {
	return reply.cmd.Val()
}

// GetString 以String类型获取响应结果
func (reply *Reply) GetString() string {
	return converts.ToString(reply.cmd.Val())
}

// GetBytes 以[]byte类型获取响应结果
func (reply *Reply) GetBytes() []byte {
	bytes, _ := reply.cmd.Bytes()
	return bytes
}

// ToAny 将结果以Json字节数组形式 写入到一个目标对象中
func (reply *Reply) ToAny(target any) error {
	return converts.ReadJsonBytesToAny(reply.GetBytes(), target)
}
