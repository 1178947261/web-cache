package lru

import "container/list"

//Cache 是LRU Cache。并发访问是不安全的。
type Cache struct {
	maxBytes   int64                         // 可以使用的最大内存
	nBytes     int64                         //  已经使用的内存
	doubleList *list.List                    //标准库的双向链表
	cache      map[string]*list.Element      // key 就是 普通的字符串 value 就是-双向链表中对应的节点的指针
	OnEvent    func(key string, value Value) // 可选，并在清除条目时执行。 是某条记录被移除时的回调函数，可以为 nil。
}

// 键值对 entry 是双向链表节点的数据类型，在链表中仍保存每个值对应的 key 的好处在于，淘汰队首节点时，需要用 key 从字典中删除对应的映射。
// 获取-元素一起存储值的数据结构。

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

// New new是一个Cache 的构造函数
func New(maxBytes int64, OnEvent func(string, Value)) *Cache {
	return &Cache{
		maxBytes:   maxBytes,
		doubleList: list.New(),
		cache:      make(map[string]*list.Element),
		OnEvent:    OnEvent,
	}
}

// Get 查找键的值
func (c *Cache) Get(key string) (value Value, ok bool) {
	// 我们先从缓存key-value 中取出value
	element, ok := c.cache[key]
	// 如果对应的连接节点存在，那我们就将节点移动到队尾部，并且返回查找的值
	//  移动到队尾（双向链表作为队列，队首队尾是相对的，在这里约定 front 为队尾）;
	if ok {
		c.doubleList.MoveToFront(element)
	}
	return
}

// 这里的删除，实际上是缓存淘汰。即移除最近最少访问的节点（队首）

// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	// 首先我们要取到-首节点
	element := c.doubleList.Back()
	// 如果对应的节点存在那我们就将他从链表中删除
	if element != nil {
		// 删除该节点
		c.doubleList.Remove(element)
		// 获取-与此元素一起存储的值。
		kv := element.Value.(*entry)
		// 从MAP 中删除
		delete(c.cache, kv.key)
		//更新当前所用的内存 c.nBytes。
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvent != nil {
			c.OnEvent(kv.key, kv.value)
		}
	}
}

// Add 如果键存在，则更新对应节点的值，并将该节点移到队尾。
//不存在则是新增场景，首先队尾添加新节点 &entry{key, value}, 并字典中添加 key 和节点的映射关系。
//更新 c.nbytes，如果超过了设定的最大值 c.maxBytes，则移除最少访问的节点。
// Add adds a value to the cache.
func (c *Cache) Add(key string, value Value) {
	// 我们先从缓存key-value 中取出value
	element, ok := c.cache[key]
	// 如果对应的连接节点存在，那我们就将节点移动到队尾部，并且返回查找的值
	//  移动到队尾（双向链表作为队列，队首队尾是相对的，在这里约定 front 为队尾）;
	if ok {
		c.doubleList.MoveToFront(element)
		// 获取-与此元素一起存储的值。
		kv := element.Value.(*entry)
		//更新当前所用的内存 c.nBytes。
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// 往链表的末尾插入
		ele := c.doubleList.PushFront(&entry{key, value})
		// 添加到-map
		c.cache[key] = ele
		// 更新内存
		c.nBytes += int64(len(key)) + int64(value.Len())

	}
	// 更新 c.nBytes，如果超过了设定的最大值 c.maxBytes，则移除最少访问的节点。
	// 这里的删除，实际上是缓存淘汰。即移除最近最少访问的节点（队首）
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.doubleList.Len()
}
