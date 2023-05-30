/**
  @author: Zero
  @date: 2023/4/22 20:15:29
  @desc: Redis 单元测试

**/

package caches

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/zlx2019/toys/converts"
	"testing"
	"time"
)

// 测试设置缓存
func TestSet(t *testing.T) {
	t.Parallel()
	is := assert.New(t)
	template := NewDefaultRedisTemplate()
	var num = 123
	err := template.Set("testKey1", num)
	is.NoError(err)
	reply := template.Get("testKey1")
	err = reply.Err()
	is.NoError(err)
	i, _ := reply.cmd.Int()
	is.Equal(i, num)
}

type Student struct {
	Name string
	Sex  bool
	Age  int
}

// 测试设置并且读取缓存
func TestGet(t *testing.T) {
	t.Parallel()
	is := assert.New(t)

	s1 := Student{
		Name: "小明",
		Sex:  true,
		Age:  18,
	}
	// 获取连接
	template := NewDefaultRedisTemplate()

	// 设置缓存
	err := template.Set("s1", s1)
	is.NoError(err)

	// 获取缓存
	reply := template.Get("s1")
	is.True(reply.Ok())
	is.NoError(reply.Err())

	// 断言1
	// 将结果转为[]byte,再通过json反序列为对象
	s2, err := converts.ReadJsonBytes[Student](reply.GetBytes())
	is.NoError(err)
	// 验证Name是否一致
	is.Equal(s1.Name, s2.Name)

	// 断言2
	// 将结果转换为String,再通过json反序列化为对象
	s3, err := converts.ReadJson[Student](reply.GetString())
	fmt.Println(s3)
	is.NoError(err)
	is.Equal(s1.Age, s3.Age)

	// 断言3
	var s4 Student
	err = reply.ToAny(&s4)
	is.NoError(err)
	is.Equal(s1.Name, s4.Name)
	is.Equal(s1.Sex, s4.Sex)
}

func TestSetExpire(t *testing.T) {
	t.Parallel()
	is := assert.New(t)
	template := NewDefaultRedisTemplate()
	// 设置有效期3秒
	err := template.SetExpire("kk1", "vv1", time.Second*3)
	is.NoError(err)

	// 2秒后取一次
	time.Sleep(time.Second * 2)
	reply := template.Get("kk1")
	is.NoError(reply.Err())

	// 2秒后再取一次,不存在
	time.Sleep(time.Second * 2)
	reply = template.Get("kk1")
	is.Error(reply.Err())
	is.Equal(reply.GetValue(), "")
}

func TestDel(t *testing.T) {
	t.Parallel()
	is := assert.New(t)
	template := NewDefaultRedisTemplate()
	err := template.Del("ns1")
	is.NoError(err)
}

func TestExists(t *testing.T) {
	t.Parallel()
	is := assert.New(t)
	template := NewDefaultRedisTemplate()
	// 已存在断言
	is.NoError(template.Set("l1", "v1"))
	is.True(template.Exists("l1"))
	// 不存在断言
	is.False(template.Exists("dwadwadaw"))
}

func TestExpireAdd(t *testing.T) {
	t.Parallel()
	is := assert.New(t)
	template := NewDefaultRedisTemplate()
	// 设置缓存
	err := template.Set("ok1", "ok1")
	is.NoError(err)
	// 为缓存设置有效期 3s
	ok := template.ExpireAdd("ok1", time.Second*3)
	is.True(ok)
	// 获取缓存
	reply := template.Get("ok1")
	is.NoError(reply.Err())
	is.Equal(reply.GetString(), "ok1")

	// 等待缓存失效
	time.Sleep(time.Second * 4)
	// 再次获取缓存
	reply = template.Get("ok1")
	// redis.Nil表示Key已不存在的错误
	is.Equal(reply.Err(), redis.Nil)
}

func TestGetExpire(t *testing.T) {
	t.Parallel()
	is := assert.New(t)
	template := NewDefaultRedisTemplate()
	err := template.SetExpire("l1", true, time.Second*60)
	is.NoError(err)
	expire, err := template.GetExpire("l1")
	is.NoError(err)
	fmt.Println(expire)
}

func TestRedisTemplate_SetNEX(t *testing.T) {
	t.Parallel()
	is := assert.New(t)
	template := NewDefaultRedisTemplate()
	err := template.SetNEX(context.Background(), "lock1", "xxxx", time.Second*10)
	is.NoError(err)

	err = template.SetNEX(context.Background(), "lock1", "dddd", time.Second*10)
	is.Error(err)

}
