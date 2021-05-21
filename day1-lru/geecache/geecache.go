package geecache

import (
	"fmt"
	"log"
	"sync"
)

//Getter 会加载密钥的数据。
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc 通过函数实现Getter。
type GetterFunc func(key string) ([]byte, error)

/**
定义接口 Getter 和 回调函数 Get(key string)([]byte, error)，参数是 key，返回值是 []byte。
定义函数类型 GetterFunc，并实现 Getter 接口的 Get 方法。
函数类型实现某一个接口，称之为接口型函数，方便使用者在调用时既能够传入函数作为参数，也能够传入实现了该接口的结构体作为参数。
*/

// Get 实现Getter接口功能
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

/**
	一个 Group 可以认为是一个缓存的命名空间，每个 Group 拥有一个唯一的名称 name。比如可以创建三个 Group，缓存学生的成绩命名为 scores，
缓存学生信息的命名为 info，缓存学生课程的命名为 courses。
*/

// Group 是一个缓存名称空间，相关的加载数据分布在 Group 是 GeeCache 最核心的数据结构，负责与用户的交互，并且控制缓存值存储和获取的流程。
type Group struct {
	name string //一个 Group 可以认为是一个缓存的命名空间，每个 Group 拥有一个唯一的名称 name。比如可以创建三个 Group，
	// 缓存学生的成绩命名为 scores，缓存学生信息的命名为 info，缓存学生课程的命名为 courses
	getter    Getter //第二个属性是 getter Getter，即缓存未命中时获取源数据的回调(callback)。
	mainCache Cache  //第三个属性是 mainCache cache，即一开始实现的并发缓存。
}

// 读写锁

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

/**
接收 key --> 检查是否被缓存 -----> 返回缓存值 ⑴
                |  否                         是
                |-----> 是否应当从远程节点获取 -----> 与远程节点交互 --> 返回缓存值 ⑵
                            |  否
                            |-----> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值 ⑶
*/

// NewGroup 构建函数 NewGroup 用来实例化 Group，并且将 group 存储在全局变量 groups 中
func NewGroup(name string, cacheByte int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		mainCache: Cache{cacheBytes: cacheByte},
		getter:    getter,
	}
	groups[name] = g //将 group 存储在全局变量 groups 中。
	return g
}

// GetGroup 返回先前使用NewGroup创建的命名组，或者如果没有这样的组，则为nil。

func GetGroup(key string) *Group {
	mu.RLock()
	defer mu.Unlock()
	return groups[key]
}

/**
Get 方法实现了上述所说的流程 ⑴ 和 ⑶。
流程 ⑴ ：从 mainCache 中查找缓存，如果存在则返回缓存值。
流程 ⑶ ：缓存不存在，则调用 load 方法，load 调用 getLocally（分布式场景下会调用 getFromPeer 从其他节点获取），
getLocally 调用用户回调函数 g.getter.Get() 获取源数据，并且将源数据添加到缓存 mainCache 中（通过 populateCache 方法）
*/

func (g *Group) Get(key string) (ByteView, error) {

	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	// 从缓存中取出数据  ---
	if v, ok := g.mainCache.get(key); ok {

		log.Println("[GeeCache] hit")
		return v, nil
	}
	// 如果-缓存失效-即缓存未命中时获取源数据的回调(callback)。
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {

	// 调用 未命中缓存情况下的--的源数据的回调(callback)。
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err

	}
	//  创建存入的缓存的数据类型
	value := ByteView{b: cloneBytes(bytes)}
	// 写入缓存中
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {

	// 调用- Cache.go 把数据添加到-缓存
	g.mainCache.add(key, value)
}
