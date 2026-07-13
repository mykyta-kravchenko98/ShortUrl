package cache

import (
	"context"
	"log/slog"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const meterName = "shorturl"

type lruCache struct {
	Capacity int
	Cache    map[string]*cacheNode
	Head     *cacheNode
	Teil     *cacheNode
	m        *sync.RWMutex

	requests  metric.Int64Counter
	evictions metric.Int64Counter
}

type cacheNode struct {
	Key   string
	Value string
	Next  *cacheNode
	Prev  *cacheNode
}

// LRUCache its interface for communication with lru cache
type LRUCache interface {
	Get(ctx context.Context, key string) string
	Put(ctx context.Context, key string, value string)
	Length() int
}

// InitLRUCache its method for creating instance of lruCache and return LRUCache interface
func InitLRUCache(capacity int) LRUCache {
	meter := otel.Meter(meterName)

	requests, err := meter.Int64Counter("cache.requests",
		metric.WithDescription("LRU cache lookups, tagged by outcome (hit/miss)"))
	if err != nil {
		panic(err)
	}
	evictions, err := meter.Int64Counter("cache.evictions",
		metric.WithDescription("Entries evicted from the LRU cache to stay under capacity"))
	if err != nil {
		panic(err)
	}

	cache := &lruCache{
		Capacity:  capacity,
		Cache:     make(map[string]*cacheNode),
		m:         &sync.RWMutex{},
		requests:  requests,
		evictions: evictions,
	}
	return cache
}

func (c *lruCache) Length() int {
	return len(c.Cache)
}

func (c *lruCache) Get(ctx context.Context, key string) string {
	c.m.Lock()
	defer c.m.Unlock()
	if value, ok := c.Cache[key]; ok {
		c.MoveFront(value)
		c.requests.Add(ctx, 1, metric.WithAttributes(attribute.String("result", "hit")))
		return value.Value
	}

	c.requests.Add(ctx, 1, metric.WithAttributes(attribute.String("result", "miss")))
	return ""
}

func (c *lruCache) Put(ctx context.Context, key string, value string) {
	c.m.Lock()
	defer c.m.Unlock()
	if result, found := c.Cache[key]; found {
		result.Value = value
		c.MoveFront(result)
	} else {
		newNode := &cacheNode{Key: key, Value: value}
		if len(c.Cache) >= c.Capacity {
			evictedKey := c.Teil.Key
			delete(c.Cache, evictedKey)
			c.RemoveTail()
			c.evictions.Add(ctx, 1)
			slog.DebugContext(ctx, "cache eviction", "evictedKey", evictedKey, "capacity", c.Capacity)
		}
		c.Cache[key] = newNode
		c.AddNode(newNode)
	}
}

func (c *lruCache) MoveFront(node *cacheNode) {
	if node == c.Head {
		return
	}

	c.RemoveNode(node)
	c.AddNode(node)
}

func (c *lruCache) RemoveTail() {
	c.RemoveNode(c.Teil)
}

func (c *lruCache) RemoveNode(node *cacheNode) {
	if node == c.Head {
		c.Head = node.Next
	}

	if node == c.Teil {
		c.Teil = node.Prev
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

func (c *lruCache) AddNode(node *cacheNode) {
	if c.Head == nil {
		c.Head = node
		c.Teil = node
	} else {
		c.Head.Prev = node
		node.Next = c.Head
		c.Head = node
	}
}
