# sugar
this is A library of caching components and support for distributed locks
<hr>

目前支持的缓存组件有:
- Redis
<hr>

## 分布式锁组件
- Redis
- Etcd

<hr>

### Redis分布式锁
**1. 非阻塞模式**
```go
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
```
**2. 阻塞模式**
```go
lock := NewRedisDistributedLock("lock-key", caches.NewDefaultRedisTemplate(),WithBlocking())
lock.Lock(ctx)
defer lock.Unlock(ctx)
```
**3. 看门狗模式**
```go
lock := NewRedisDistributedLock("dog-lock-key", caches.NewDefaultRedisTemplate(),WithExpire(time.Second * 5),WithWatchDog())
```