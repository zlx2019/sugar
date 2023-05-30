/**
  @author: Zero
  @date: 2023/4/22 16:43:30
  @desc: 缓存客户端

**/

package caches

// CacheClient 缓存组件客户端
type CacheClient struct {
	// 缓存模板实例
	CacheTemplate
}

// NewCacheClient Create Cache Client With Template
func NewCacheClient(template CacheTemplate) *CacheClient {
	return &CacheClient{
		template,
	}
}
