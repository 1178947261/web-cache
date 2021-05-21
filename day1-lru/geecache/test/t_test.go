package test

import (
	"fmt"
	geecache "geecahe"
	"reflect"
	"sync"
	"testing"
	"time"
)

var m sync.Mutex

var set = make(map[int]bool, 0)

func printOnce(num int) {
	m.Lock()
	defer m.Unlock()
	if _, exist := set[num]; !exist {
		fmt.Println(num)
	}
	set[num] = true

}

func TestGet(t *testing.T) {
	for i := 0; i < 10; i++ {
		go printOnce(100)
	}
	time.Sleep(time.Second)
}
func TestGetter(t *testing.T) {
	var f geecache.Getter = geecache.GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")

	}
	fmt.Println(f.Get("key"))
}
