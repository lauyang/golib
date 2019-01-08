package memcache

import (
	"fmt"
	"testing"
	"time"

	"net/http"
	_ "net/http/pprof"
)

func init() {
	go func() {
		http.ListenAndServe(":1100", nil)
	}()
}

type myNode struct {
	node MemNode

	val  string
	data [1024 * 1024]byte
}

func (self *myNode) GetMemNode() *MemNode {
	return &self.node
}

func (self *myNode) GetInterface() interface{} {
	return self
}

// 打印全部数据
func printAll(mem *MemCache) {
	fmt.Println("==========map==========")
	for key, cache := range mem.storeMap {
		node, _ := cache.(*myNode)
		fmt.Println(key, "=>", node.val)
	}

	fmt.Println("=========list=========")
	first := mem.storeFirst
	for ; first != nil; first = first.next {
		fmt.Println(first.key, "=>", first)
	}
}

func _Test_Memcache(t *testing.T) {
	mem := MemCache{}
	mem2 := MemCache{}

	mem.Init(5, 1000)
	mem2.Init(1000, 3)
	// 把数据加入
	var data [10]*myNode
	for i := 0; i < 10; i++ {
		data[i] = &myNode{val: fmt.Sprint("value-", i+1)}
	}

	fmt.Println("test timeout >>>>>>>>>>")
	for index, node := range data {
		err := mem.Set(fmt.Sprint("key-", index+1), node)
		if nil != err {
			fmt.Println("add:", err)
		}

		printAll(&mem)
		time.Sleep(time.Second * 1)

		node, ok := mem.Get("key-1").(*myNode)
		if false == ok {
			t.Error("not find key-1")
			return
		} else {
			fmt.Println(node.node.key, "=>", node.val)
		}

		node, ok = mem.Get("key-11").(*myNode)
		fmt.Println(ok, node)
	}

	fmt.Println("test max >>>>>>>>>>")
	for index, node := range data {
		err := mem2.Set(fmt.Sprint("key-", index+1), node)
		if nil != err {
			fmt.Println("add:", err)
		}

		printAll(&mem2)
	}
}

//func Test_Memcache_2(t *testing.T) {
//	mem := &MemCache{}

//	mem.Init(5, 1024)
//	for n := 0; true; n++ {
//		for i := 0; i < 2048; i++ {
//			node := &myNode{}
//			copy(node.data[:], []byte("start"))

//			node.data[1024*1024-1] = 'd'
//			node.data[1024*1024-2] = 'n'
//			node.data[1024*1024-3] = 'e'

//			mem.Set(fmt.Sprint("key-", i+1), node)
//		}

//		fmt.Println(n, ". will sleep 6 second.")
//		time.Sleep(time.Second * 6)
//	}
//}
