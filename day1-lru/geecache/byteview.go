package geecache

/**
ByteView 只有一个数据成员，b []byte，b 将会存储真实的缓存值。选择 byte 类型是为了能够支持任意的数据类型的存储，例如字符串、图片等。
实现 Len() int 方法，我们在 lru.Cache 的实现中，要求被缓存对象必须实现 Value 接口，即 Len() int 方法，返回其所占的内存大小。
b 是只读的，使用 ByteSlice() 方法返回一个拷贝，防止缓存值被外部程序修改。
*/

type ByteView struct {
	b []byte
}

// Len 返回其所占的内存大小。
func (b ByteView) len() int {
	return len(b.b)

}

// ByteSlice  方法返回一个拷贝
func (b ByteView) ByteSlice() []byte {
	return cloneBytes(b.b)
}

// 字符串以字符串形式返回数据，并在必要时进行复制。
func (b ByteView) String() string {
	return string(b.b)
}

// 克隆
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}