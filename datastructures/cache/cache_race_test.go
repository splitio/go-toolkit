// +build !race

package cache

import (
	"math/rand"
	"sync"
	"testing"
)

func TestLocalCacheHighConcurrency(t *testing.T) {

	c, err := NewSimpleLRU[int64, int64](500, NoTTL)
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
