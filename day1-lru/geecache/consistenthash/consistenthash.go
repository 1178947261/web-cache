package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
type Hash func(data []byte) uint32

// Map 是一致性哈希算法的主数据结构，包含 4 个成员变量：Hash 函数 hash；虚拟节点倍数 replicas；
//  哈希环 keys；虚拟节点与真实节点的映射表 hashMap，键是虚拟节点的哈希值，值是真实节点的名称。
// 构造函数 New() 允许自定义虚拟节点倍数和 Hash 函数。

// Map 包含所有哈希键
type Map struct {
	hash     Hash
	replicas int
	keys     []int // Sorted
	hashMap  map[int]string
}

// New creates a Map instance
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,             //虚拟节点倍数
		hash:     fn,                   // 自定义的hash
		hashMap:  make(map[int]string), // 拟节点与真实节点的映射表 hashMap，键是虚拟节点的哈希值，值是真实节点的名称。
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

/**

Add 函数允许传入 0 或 多个真实节点的名称。
对每一个真实节点 key，对应创建 m.replicas 个虚拟节点，虚拟节点的名称是：strconv.Iota(i) + key，即通过添加编号的方式区分不同虚拟节点。
使用 m.hash() 计算虚拟节点的哈希值，使用 append(m.keys, hash) 添加到环上。
在 hashMap 中增加虚拟节点和真实节点的映射关系。
最后一步，环上的哈希值排序。
*/

// Add adds some keys to the hash.
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// 转为-为了拼接而不是-计算
			// m/hash 可以是自定义传入的方法-
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash) // 虚拟节点的key
			m.hashMap[hash] = key         //  虚拟节点和真实节点做了映射-
		}
	}
	// 排序
	sort.Ints(m.keys)
}

func (m Map) Get(key string) string {

	// 先判断 key
	if len(key) == 0 {
		return ""
	}
	// 计算哈希值-
	hash := int(m.hash([]byte(key)))

	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]

}

// Remove use to remove a key and its virtual keys on the ring and map
func (m *Map) Remove(key string) {
	for i := 0; i < m.replicas; i++ {
		hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
		idx := sort.SearchInts(m.keys, hash)
		m.keys = append(m.keys[:idx], m.keys[idx+1:]...)
		delete(m.hashMap, hash)
	}
}
