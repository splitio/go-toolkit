package cache

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type LocalCache struct {
	ttl    time.Duration
	maxLen int
	ttls   map[string]time.Time
	items  map[string]*list.Element
	lru    *list.List
	mutex  sync.Mutex
}

type entry struct {
	key   string
	value interface{}
}

func (c *LocalCache) Get(key string) (interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	node, ok := c.items[key]
	if !ok {
		return nil, &Miss{Where: "LOCAL", Key: key}
	}

	val := node.Value.(entry).value
	ttl, ok := c.ttls[key]
	if !ok || time.Now().Sub(ttl) > c.ttl {
		return nil, &Expired{Key: key, Value: val, When: ttl.Add(c.ttl)}
	}

	c.lru.MoveToFront(node)
	return val, nil
}

func (c *LocalCache) Set(key string, value interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if node, ok := c.items[key]; ok {
		c.lru.MoveToFront(node)
		node.Value = entry{key: key, value: value}
	} else {
		// Drop the LRU item on the list before adding a new one.
		if c.lru.Len() == c.maxLen {
			key := c.lru.Back().Value.(entry).key
			delete(c.items, key)
			c.lru.Remove(c.lru.Back())
		}

		ptr := c.lru.PushFront(entry{key: key, value: value})
		c.items[key] = ptr
		c.ttls[key] = time.Now()
	}

	return nil
}

func NewLocalCache(maxSize int, ttl time.Duration) (*LocalCache, error) {
	if maxSize <= 0 {
		return nil, fmt.Errorf("Cache size should be > 0. Is: %d", maxSize)
	}

	return &LocalCache{
		maxLen: maxSize,
		ttl:    ttl,
		lru:    new(list.List),
		items:  make(map[string]*list.Element, maxSize),
		ttls:   make(map[string]time.Time),
	}, nil
}
