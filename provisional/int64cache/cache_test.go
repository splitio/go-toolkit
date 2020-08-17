package int64cache

import (
	"math/rand"
	"sync"
	"testing"
)

func TestInt64Cache(t *testing.T) {
	c, err := NewInt64Cache(5)
	if err != nil {
		t.Error("No error should have been returned. Got: ", err)
	}

	for i := int64(1); i <= 5; i++ {
		err := c.Set(i, i)
		if err != nil {
			t.Errorf("Setting value '%d', should not have raised an error. Got: %s", i, err)
		}
	}

	for i := int64(1); i <= 5; i++ {
		val, err := c.Get(i)
		if err != nil {
			t.Errorf("Getting value '%d', should not have raised an error. Got: %s", i, err)
		}
		if val != i {
			t.Errorf("Value for key '%d' should be %d. Is %d", i, i, val)
		}
	}

	c.Set(6, 6)

	// Oldest item (1) should have been removed
	val, err := c.Get(1)
	if err == nil {
		t.Errorf("Getting value 'someKey1', should not have raised an error. Got: %s", err)
	}

	_, ok := err.(*Miss)
	if !ok {
		t.Errorf("Error should be of type Miss. Is %T", err)
	}

	if val != 0 {
		t.Errorf("Value for key 'someKey1' should be nil. Is %d", val)
	}

	// 2-6 should be available
	for i := int64(2); i <= 6; i++ {
		val, err := c.Get(i)
		if err != nil {
			t.Errorf("Getting value '%d', should not have raised an error. Got: %s", i, err)
		}
		if val != i {
			t.Errorf("Value for key '%d' should be %d. Is %d", i, i, val)
		}
	}

	if len(c.items) != 5 {
		t.Error("Items size should be 5. is: ", len(c.items))
	}
}

func TestLocalCacheHighConcurrency(t *testing.T) {

	c, err := NewInt64Cache(500)
	if err != nil {
		t.Error("No error should have been returned. Got: ", err)
	}

	iterations := int64(500000)
	wg := sync.WaitGroup{}
	wg.Add(int(iterations))
	for i := int64(0); i < iterations; i++ {
		r := int64(rand.Intn(500))
		if i%2 == 0 {
			go func() {
				defer wg.Done()
				c.Set(r, r)
			}()

		} else {
			go func() {
				defer wg.Done()
				c.Get(r)
			}()
		}
	}
	wg.Wait()
}
