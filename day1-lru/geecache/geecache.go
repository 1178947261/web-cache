package geecache

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
