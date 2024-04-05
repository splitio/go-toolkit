package cache

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

const (
	NoTTL = -1
)

// SimpleLRU is an in-memory TTL & LRU cache
type SimpleLRU[K comparable, V any] interface {
	Get(key K) (V, error)
	Set(key K, value V) error
}

// SimpleLRUImpl implements the Simple interface
type SimpleLRUImpl[K comparable, V any] struct {
	ttl    time.Duration
	maxLen int
	ttls   map[K]time.Time
	items  map[K]*list.Element
	lru    *list.List
	mutex  sync.Mutex
}

type centry[K comparable, V any] struct {
	key   K
	value V
}

// Get retrieves an item if exist, nil + an error otherwise
func (c *SimpleLRUImpl[K, V]) Get(key K) (V, error) {
	var empty V
	c.mutex.Lock()
	defer c.mutex.Unlock()
	node, ok := c.items[key]
	if !ok {
		return empty, &Miss{Where: "LOCAL", Key: key}
	}

	entry, ok := node.Value.(centry[K, V])
	if !ok {
		return empty, fmt.Errorf("Invalid data in cache for key %v", key)
	}

	if c.ttls != nil { // TTL enabled
		ttl, ok := c.ttls[key]
		if !ok {
			return empty, fmt.Errorf(
				"Missing TTL for key %v. Wrapping as expired: %w",
				key,
				&Expired{Key: key, Value: entry.value, When: ttl.Add(c.ttl)},
			)
		}

		if time.Now().UnixNano() > ttl.UnixNano() {
			return empty, &Expired{Key: key, Value: entry.value, When: ttl.Add(c.ttl)}
		}
	}

	c.lru.MoveToFront(node)
	return entry.value, nil
}

// Set adds a new item. Since the cache being full results in removing the LRU element, this method never fails.
func (c *SimpleLRUImpl[K, V]) Set(key K, value V) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if node, ok := c.items[key]; ok {
		c.lru.MoveToFront(node)
		node.Value = centry[K, V]{key: key, value: value}
	} else {
		// Drop the LRU item on the list before adding a new one.
		if c.lru.Len() == c.maxLen {
			entry, ok := c.lru.Back().Value.(centry[K, V])
			if !ok {
				return fmt.Errorf("Invalid data in list for key %v", key)
			}
			key := entry.key
			delete(c.items, key)
			if c.ttls != nil {
				delete(c.ttls, key)
			}
			c.lru.Remove(c.lru.Back())
		}

		ptr := c.lru.PushFront(centry[K, V]{key: key, value: value})
		c.items[key] = ptr
	}

	if c.ttls != nil {
		c.ttls[key] = time.Now().Add(c.ttl)
	}
	return nil
}

// NewSimple returns a new Simple instance of the specified size and TTL
func NewSimpleLRU[K comparable, V any](maxSize int, ttl time.Duration) (*SimpleLRUImpl[K, V], error) {
	if maxSize <= 0 {
		return nil, fmt.Errorf("Cache size should be > 0. Is: %d", maxSize)
	}
    
    var ttls map[K]time.Time = nil
    if ttl != NoTTL {
        ttls = make(map[K]time.Time)
    }

	return &SimpleLRUImpl[K, V]{
		maxLen: maxSize,
		ttl:    ttl,
		lru:    new(list.List),
		items:  make(map[K]*list.Element, maxSize),
		ttls:   ttls,
	}, nil
}

var _ SimpleLRU[string, int] = (*SimpleLRUImpl[string, int])(nil)
