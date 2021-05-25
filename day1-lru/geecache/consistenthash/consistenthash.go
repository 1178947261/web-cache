package consistenthash

import "hash/crc32"

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
