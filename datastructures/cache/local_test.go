package cache

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestLocalCache(t *testing.T) {
	cache, err := NewLocalCache(5, 1*time.Second)
	if err != nil {
		t.Error("No error should have been returned. Got: ", err)
	}

	for i := 1; i <= 5; i++ {
		err := cache.Set(fmt.Sprintf("someKey%d", i), i)
		if err != nil {
			t.Errorf("Setting value 'someKey%d', should not have raised an error. Got: %s", i, err)
		}
	}

	for i := 1; i <= 5; i++ {
		val, err := cache.Get(fmt.Sprintf("someKey%d", i))
		if err != nil {
			t.Errorf("Getting value 'someKey%d', should not have raised an error. Got: %s", i, err)
		}
		if val != i {
			t.Errorf("Value for key 'someKey%d' should be %d. Is %d", i, i, val)
		}
	}

	cache.Set("someKey6", 6)

	// Oldest item (1) should have been removed
	val, err := cache.Get("someKey1")
	if err == nil {
		t.Errorf("Getting value 'someKey1', should not have raised an error. Got: %s", err)
	}

	asMiss, ok := err.(*Miss)
	if !ok {
		t.Errorf("Error should be of type Miss. Is %T", err)
	}

	if asMiss.Key != "someKey1" || asMiss.Where != "LOCAL" {
		t.Errorf("Incorrect data within the Miss error. Got: %+v", asMiss)
	}

	if val != nil {
		t.Errorf("Value for key 'someKey1' should be nil. Is %d", val)
	}

	// 2-6 should be available
	for i := 2; i <= 6; i++ {
		val, err := cache.Get(fmt.Sprintf("someKey%d", i))
		if err != nil {
			t.Errorf("Getting value 'someKey%d', should not have raised an error. Got: %s", i, err)
		}
		if val != i {
			t.Errorf("Value for key 'someKey%d' should be %d. Is %d", i, i, val)
		}
	}

	if len(cache.items) != 5 {
		t.Error("Items size should be 5. is: ", len(cache.items))
	}

	if cache.lru.Len() != 5 {
		t.Error("LRU size should be 5. is: ", cache.lru.Len())
	}

	time.Sleep(2 * time.Second) // Wait for all keys to expire.
	for i := 2; i <= 6; i++ {
		val, err := cache.Get(fmt.Sprintf("someKey%d", i))
		if val != nil {
			t.Errorf("No value should have been returned for expired key 'someKey%d'.", i)
		}

		if err == nil {
			t.Errorf("Getting value 'someKey%d', should have raised an 'Expired' error. Got nil", i)
			continue
		}

		asExpiredErr, ok := err.(*Expired)
		if !ok {
			t.Errorf("Returned error should be of 'Expired' type. Is %T", err)
			continue
		}

		if asExpiredErr.Key != fmt.Sprintf("someKey%d", i) {
			t.Errorf("Key in Expired error should be 'someKey%d'. Is: '%s'", i, asExpiredErr.Key)
		}

		if asExpiredErr.Value != i {
			t.Errorf("Value in Expired error should be %d. Is %+v", i, asExpiredErr.Value)
		}

		ttl, ok := cache.ttls[fmt.Sprintf("someKey%d", i)]
		if !ok {
			t.Errorf("A ttl entry should exist for key 'someKey%d'", i)
			continue
		}

		if asExpiredErr.When != ttl.Add(cache.ttl) {
			t.Errorf("Key 'someKey%d' should have expired at %+v. It did at %+v", i, ttl.Add(cache.ttl), asExpiredErr.When)
		}
	}
}

func TestLocalCacheHighConcurrency(t *testing.T) {

	cache, err := NewLocalCache(500, 1*time.Second)
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
