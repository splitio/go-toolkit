package cache

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimpleCache(t *testing.T) {
	cache, err := NewSimpleLRU[string, int](5, 1*time.Second)
	assert.Nil(t, err)

	for i := 1; i <= 5; i++ {
		err := cache.Set(fmt.Sprintf("someKey%d", i), i)
		assert.Nil(t, err)
	}

	for i := 1; i <= 5; i++ {
		val, err := cache.Get(fmt.Sprintf("someKey%d", i))
		assert.Nil(t, err)
		assert.Equal(t, i, val)
	}

	cache.Set("someKey6", 6)

	// Oldest item (1) should have been removed
	val, err := cache.Get("someKey1")
	assert.NotNil(t, err)
	asMiss, ok := err.(*Miss)
	assert.True(t, ok)
	assert.Equal(t, "someKey1", asMiss.Key)
	assert.Equal(t, "LOCAL", asMiss.Where)
	assert.Equal(t, 0, val)

	// 2-6 should be available
	for i := 2; i <= 6; i++ {
		val, err := cache.Get(fmt.Sprintf("someKey%d", i))
		assert.Nil(t, err)
		assert.Equal(t, i, val)
	}

	assert.Equal(t, 5, len(cache.items))
	assert.Equal(t, 5, len(cache.ttls))
	assert.Equal(t, 5, cache.lru.Len())

	time.Sleep(2 * time.Second) // Wait for all keys to expire.

	for i := 2; i <= 6; i++ {
		val, err := cache.Get(fmt.Sprintf("someKey%d", i))
		assert.Equal(t, 0, val)
		assert.NotNil(t, err)
		asExpired, ok := err.(*Expired)
		assert.True(t, ok)
		assert.Equal(t, fmt.Sprintf("someKey%d", i), asExpired.Key)
		assert.Equal(t, i, asExpired.Value)

		ttl, ok := cache.ttls[fmt.Sprintf("someKey%d", i)]
		assert.True(t, ok)
		assert.Equal(t, asExpired.When, ttl.Add(cache.ttl))

	}
}

func TestSimpleCacheHighConcurrency(t *testing.T) {

	cache, err := NewSimpleLRU[string, int](500, 1*time.Second)
	if err != nil {
		t.Error("No error should have been returned. Got: ", err)
	}

	iterations := 500000
	wg := sync.WaitGroup{}
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		r := rand.Intn(500)
		if i%2 == 0 {
			go func() {
				defer wg.Done()
				cache.Set(fmt.Sprintf("someKey%d", r), r)
			}()

		} else {
			go func() {
				defer wg.Done()
				cache.Get(fmt.Sprintf("someKey%d", r))
			}()
		}
	}
	wg.Wait()
}

func TestInt64Cache(t *testing.T) {
	c, err := NewSimpleLRU[int64, int64](5, NoTTL)
	assert.Nil(t, err)

	for i := int64(1); i <= 5; i++ {
		assert.Nil(t, c.Set(i, i))
	}

	for i := int64(1); i <= 5; i++ {
		val, err := c.Get(i)
        assert.Nil(t, err)
        assert.Equal(t, i, val)
	}

	c.Set(6, 6)

	// Oldest item (1) should have been removed
	val, err := c.Get(1)
    assert.NotNil(t, err)
	_, ok := err.(*Miss)
    assert.True(t, ok)
    assert.Equal(t, int64(0), val)

	// 2-6 should be available
	for i := int64(2); i <= 6; i++ {
		val, err := c.Get(i)
        assert.Nil(t, err)
        assert.Equal(t, i, val)
	}

    assert.Equal(t, 5, len(c.items))
}
