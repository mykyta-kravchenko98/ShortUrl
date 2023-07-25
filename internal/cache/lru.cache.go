package cache

import "sync"

type lruCache struct {
	Capacity int
	Cache    map[string]*cacheNode
	Head     *cacheNode
	Teil     *cacheNode
	m        *sync.RWMutex
}

type cacheNode struct {
	Key   string
	Value string
	Next  *cacheNode
	Prev  *cacheNode
}

//LRUCache its interface for comunication with lru cache
type LRUCache interface {
	Get(key string) string
	Put(key string, value string)
	Length() int
}

// InitLRUCache its method for creating instance of lruCache and return LRUCache interface
func InitLRUCache(capacity int) LRUCache {
	cache := &lruCache{
		Capacity: capacity,
		Cache:    make(map[string]*cacheNode),
		m:        &sync.RWMutex{},
	}
	return cache
}

func (с *lruCache) Length() int {
	return len(с.Cache)
}

func (с *lruCache) Get(key string) string {
	с.m.Lock()
	defer с.m.Unlock()
	if value, ok := с.Cache[key]; ok {
		с.MoveFront(value)
		return value.Value
	}

	return ""
}

func (с *lruCache) Put(key string, value string) {
	с.m.Lock()
	defer с.m.Unlock()
	if result, found := с.Cache[key]; found {
		result.Value = value
		с.MoveFront(result)
	} else {
		newNode := &cacheNode{Key: key, Value: value}
		if len(с.Cache) >= с.Capacity {
			delete(с.Cache, с.Teil.Key)
			с.RemoveTail()
		}
		с.Cache[key] = newNode
		с.AddNode(newNode)
	}
}

func (с *lruCache) MoveFront(node *cacheNode) {
	if node == с.Head {
		return
	}

	с.RemoveNode(node)
	с.AddNode(node)
}

func (с *lruCache) RemoveTail() {
	с.RemoveNode(с.Teil)
}

func (с *lruCache) RemoveNode(node *cacheNode) {
	if node == с.Head {
		с.Head = node.Next
	}

	if node == с.Teil {
		с.Teil = node.Prev
	}

	if node.Prev != nil {
		node.Prev.Next = node.Next
	}

	if node.Next != nil {
		node.Next.Prev = node.Prev
	}

	node.Prev = nil
	node.Next = nil
}

func (с *lruCache) AddNode(node *cacheNode) {
	if с.Head == nil {
		с.Head = node
		с.Teil = node
	} else {
		с.Head.Prev = node
		node.Next = с.Head
		с.Head = node
	}
}
