package geecache

/**
cache.go 的实现非常简单，实例化 lru，封装 get 和 add 方法，并添加互斥锁 mu。
在 add 方法中，判断了 c.lru 是否为 nil，如果等于 nil 再创建实例。这种方法称之为延迟初始化(Lazy Initialization)，
一个对象的延迟初始化意味着该对象的创建将会延迟至第一次使用该对象时。主要用于提高性能，并减少程序内存要求
*/