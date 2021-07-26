package cache

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

// LocalCache is an in-memory TTL & LRU cache
type LocalCache interface {
	Get(key string) (interface{}, error)
	Set(key string, value interface{}) error
}

// LocalCacheImpl implements the LocalCache interface
type LocalCacheImpl struct {
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

// Get retrieves an item if exist, nil + an error otherwise
func (c *LocalCacheImpl) Get(key string) (interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	node, ok := c.items[key]
	if !ok {
		return nil, &Miss{Where: "LOCAL", Key: key}
	}

	entry, ok := node.Value.(entry)
	if !ok {
		return nil, fmt.Errorf("Invalid data in cache for key %s", key)
	}

	ttl, ok := c.ttls[key]
	if !ok {
		return nil, fmt.Errorf(
			"Missing TTL for key %s. Wrapping as expired: %w",
			key,
			&Expired{Key: key, Value: entry.value, When: ttl.Add(c.ttl)},
		)
	}

	if time.Now().UnixNano() > ttl.UnixNano() {
		return nil, &Expired{Key: key, Value: entry.value, When: ttl.Add(c.ttl)}
	}

	c.lru.MoveToFront(node)
	return entry.value, nil
}

// Set adds a new item. Since the cache being full results in removing the LRU element, this method never fails.
func (c *LocalCacheImpl) Set(key string, value interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if node, ok := c.items[key]; ok {
		c.lru.MoveToFront(node)
		node.Value = entry{key: key, value: value}
	} else {
		// Drop the LRU item on the list before adding a new one.
		if c.lru.Len() == c.maxLen {
			entry, ok := c.lru.Back().Value.(entry)
			if !ok {
				return fmt.Errorf("Invalid data in list for key %s", key)
			}
			key := entry.key
			delete(c.items, key)
			delete(c.ttls, key)
			c.lru.Remove(c.lru.Back())
		}

		ptr := c.lru.PushFront(entry{key: key, value: value})
		c.items[key] = ptr
	}
	c.ttls[key] = time.Now().Add(c.ttl)
	return nil
}

// NewLocalCache returns a new LocalCache instance of the specified size and TTL
func NewLocalCache(maxSize int, ttl time.Duration) (*LocalCacheImpl, error) {
	if maxSize <= 0 {
		return nil, fmt.Errorf("Cache size should be > 0. Is: %d", maxSize)
	}

	return &LocalCacheImpl{
		maxLen: maxSize,
		ttl:    ttl,
		lru:    new(list.List),
		items:  make(map[string]*list.Element, maxSize),
		ttls:   make(map[string]time.Time),
	}, nil
}
