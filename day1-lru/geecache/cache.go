package geecache

import (
	lru2 "geecahe/lru"
	"sync"
)

/**
cache.go 的实现非常简单，实例化 lru，封装 get 和 add 方法，并添加互斥锁 mu。
在 add 方法中，判断了 c.lru 是否为 nil，如果等于 nil 再创建实例。这种方法称之为延迟初始化(Lazy Initialization)，
一个对象的延迟初始化意味着该对象的创建将会延迟至第一次使用该对象时。主要用于提高性能，并减少程序内存要求
*/

type Cache struct {
	mu         sync.Mutex
	lru        *lru2.Cache
	cacheBytes int64
}

func (c *Cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru2.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)

}
