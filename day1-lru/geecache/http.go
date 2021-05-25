package geecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// 提供被其他节点访问的能力(基于http)
const defaultBasePath = "/_geecache/"

// HTTPPool 作为承载节点间 HTTP 通信的核心数据结构（包括服务端和客户端，今天只实现服务端）。
type HTTPPool struct {
	// this peer's base URL, e.g. "https://example.net:8000"
	self     string
	basePath string
}

// NewHTTPPool initializes an HTTP pool of peers.
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Log 服务器名称的信息
func (h HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", h.self, fmt.Sprintf(format, v...))
}

//ServeHTTP 处理所有http请求
func (h HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// 检查URL 是否相等
	if !strings.HasPrefix(r.URL.Path, h.basePath) {
		panic("HTTPPool  服务意外路径: " + r.URL.Path)
	}
	h.Log("%s %s", r.Method, r.URL.Path)

	// 按照斜杠 分割为-切片
	parts := strings.SplitN(r.URL.Path[len(h.basePath):], "/", 2)
	// 判断是否包含了-缓存分组和缓存的KEY
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]
	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "没有找到-该缓存分组："+groupName, http.StatusNotFound)
		return
	}
	view, err := group.Get(key)

	if err != nil {
		http.Error(w, "没有找到-该缓存："+err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}
