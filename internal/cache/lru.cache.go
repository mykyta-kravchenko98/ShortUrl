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

type LRUCache interface {
	Get(key string) string
	Put(key string, value string)
	Length() int
}

func InitLRUCache(capacity int) LRUCache {
	cache := &lruCache{
		Capacity: capacity,
		Cache:    make(map[string]*cacheNode),
		m:        &sync.RWMutex{},
	}
	return cache
}

func (this *lruCache) Length() int {
	return len(this.Cache)
}

func (this *lruCache) Get(key string) string {
	this.m.Lock()
	defer this.m.Unlock()
	if value, ok := this.Cache[key]; ok {
		this.MoveFront(value)
		return value.Value
	}

	return ""
}

func (this *lruCache) Put(key string, value string) {
	this.m.Lock()
	defer this.m.Unlock()
	if result, found := this.Cache[key]; found {
		result.Value = value
		this.MoveFront(result)
	} else {
		newNode := &cacheNode{Key: key, Value: value}
		if len(this.Cache) >= this.Capacity {
			delete(this.Cache, this.Teil.Key)
			this.RemoveTail()
		}
		this.Cache[key] = newNode
		this.AddNode(newNode)
	}
}

func (this *lruCache) MoveFront(node *cacheNode) {
	if node == this.Head {
		return
	}

	this.RemoveNode(node)
	this.AddNode(node)
}

func (this *lruCache) RemoveTail() {
	this.RemoveNode(this.Teil)
}

func (this *lruCache) RemoveNode(node *cacheNode) {
	if node == this.Head {
		this.Head = node.Next
	}

	if node == this.Teil {
		this.Teil = node.Prev
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

func (this *lruCache) AddNode(node *cacheNode) {
	if this.Head == nil {
		this.Head = node
		this.Teil = node
	} else {
		this.Head.Prev = node
		node.Next = this.Head
		this.Head = node
	}
}
