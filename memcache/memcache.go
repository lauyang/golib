package memcache

import (
	"errors"
	"sync"
	"time"
)

type MemNode struct {
	activeTime int64
	key        string

	prev *MemNode
	next *MemNode
}

type MemCache struct {
	activeMax int64
	countMax  int64

	storeMap    map[string]MemNodeInterface
	storeFirst  *MemNode
	storeLast   *MemNode
	storeLocker sync.Mutex
}

type MemNodeInterface interface {
	GetMemNode() *MemNode
	GetInterface() interface{}
}

// 初始化
// activeTime: 存活时间，以秒为单位
// max: 最大记录数, 超过最大记录数，将活跃时间最久的那一个自动删除
func (self *MemCache) Init(activeMax int64, countMax int64) {
	self.activeMax = activeMax
	self.countMax = countMax

	self.storeMap = make(map[string]MemNodeInterface)
	self.storeFirst = nil
	self.storeLast = nil
}

// 添加新项
func (self *MemCache) Set(key string, mem MemNodeInterface) error {
	// 判断是否为空对像
	if nil == mem {
		return errors.New("not found interface MemNodeInterface")
	}

	// 判断是否为要设置的内存
	newNode := mem.GetMemNode()
	if nil == newNode {
		return errors.New("not found interface MemNode")
	}

	newNode.activeTime = time.Now().Unix() + self.activeMax
	newNode.key = key
	// 加锁操作
	self.storeLocker.Lock()
	defer self.storeLocker.Unlock()

	// 要找key是否存在
	old, ok := self.storeMap[key]
	if ok && nil != old {
		// 替换旧的
		node := old.GetMemNode()
		if nil != node {
			// 删除旧的
			self.deleteNode(node)
			// 添加新的
			self.storeMap[key] = mem
			self.pushList(newNode)
			return nil
		}
	}

	nowTime := time.Now().Unix()
	// 删除过期的节点
	for len(self.storeMap) > 0 {
		node := self.storeLast
		// 判断链表是否还有值
		if nil == node {
			break
		}

		if node.activeTime < nowTime {
			self.deleteNode(node)
		} else {
			// 没有这句会成为死循环
			break
		}
	}

	// 判断是否大于最大值
	if int64(len(self.storeMap)) >= self.countMax {
		// 删除最不活跃的那一个
		self.deleteNode(self.storeLast)
	}

	// 加入
	self.storeMap[key] = mem
	self.pushList(newNode)

	return nil
}

// 获取节点
func (self *MemCache) Get(key string) interface{} {
	// 加锁
	self.storeLocker.Lock()
	defer self.storeLocker.Unlock()

	// 从map中查找
	mem, ok := self.storeMap[key]
	if false == ok || nil == mem {
		return nil
	}

	node := mem.GetMemNode()
	if false == ok || nil == node {
		return nil
	}

	// 更改活跃时间
	self.popList(node)
	node.activeTime = time.Now().Unix() + self.activeMax
	self.pushList(node)

	return mem.GetInterface()
}

func (self *MemCache) Size() int {
	return len(self.storeMap)
}

// 删除节点
func (self *MemCache) deleteNode(node *MemNode) {
	if nil == node {
		return
	}

	// 从map中移除
	delete(self.storeMap, node.key)
	// 从链表中移除
	self.popList(node)
}

// 加入链表
func (self *MemCache) pushList(node *MemNode) {
	node.prev = nil
	node.next = self.storeFirst

	if nil != self.storeFirst {
		self.storeFirst.prev = node
	}

	self.storeFirst = node
	// 判断链表是否为空
	if nil == self.storeLast {
		self.storeLast = node
	}
}

// 从链表中移出
func (self *MemCache) popList(node *MemNode) {
	if nil == node {
		return
	}

	prev := node.prev
	next := node.next

	node.prev = nil
	node.next = nil

	if nil != prev {
		prev.next = next
	}

	if nil != next {
		next.prev = prev
	}

	if node == self.storeLast {
		self.storeLast = prev
	}

	if node == self.storeFirst {
		self.storeFirst = next
	}
}
