package queuecache

import (
	"errors"
	"math"
	"testing"
)

func TestCacheBasicUsage(t *testing.T) {

	data := make([]int, 100)
	for index := range data {
		data[index] = index
	}
	index := 0

	fetchMore := func(count int) ([]interface{}, error) {
		if index == 100 {
			return nil, errors.New("NO_MORE_DATA")
		}
		numberOfElementsToReturn := int(math.Min(float64(count), float64(100-index)))
		toReturn := make([]interface{}, numberOfElementsToReturn)
		for localIndex := range toReturn {
			toReturn[localIndex] = data[index]
			index++
		}

		return toReturn, nil
	}

	myCache := New(10, fetchMore)
	first5, err := myCache.Fetch(5)
	if err != nil {
		t.Error(err)
	}

	for index, item := range first5 {
		asInt, ok := item.(int)
		if !ok {
			t.Error("Item should be stored as int and isn't")
		}

		if asInt != index {
			t.Error("Each number should be equal to its index")
		}
	}

	offset := 5
	next5, err := myCache.Fetch(5)
	if err != nil {
		t.Error(err)
	}
	for index, item := range next5 {
		asInt, ok := item.(int)
		if !ok {
			t.Error("Item should be stored as int and isn't")
		}

		if asInt != index+offset {
			t.Error("Each number should be equal to its index")
		}
	}

	index = 0
	myCache = New(10, fetchMore)
	for i := 0; i < 100; i++ {
		elem, err := myCache.Fetch(1)
		if err != nil {
			t.Error(err)
		}

		asInt, ok := elem[0].(int)
		if !ok {
			t.Error("Item should be stored as int and isn't")
		}

		if asInt != i {
			t.Error("Each number should be equal to its index")
			t.Error("asInt", asInt)
			t.Error("index", i)
		}
	}

	elems, err := myCache.Fetch(1)
	if elems != nil {
		t.Error("Elem should be nil and is: ", elems)
	}

	if err == nil || err.Error() != "NO_MORE_DATA" {
		t.Error("Error should be NO_MORE_DATA and is: ", err.Error())
	}

	// Set index to 0 so that refill works and restart tests.
	index = 0
	for i := 0; i < 100; i++ {
		elem, err := myCache.Fetch(1)
		if err != nil {
			t.Error(err)
		}

		asInt, ok := elem[0].(int)
		if !ok {
			t.Error("Item should be stored as int and isn't")
		}

		if asInt != i {
			t.Error("Each number should be equal to its index")
			t.Error("asInt", asInt)
			t.Error("index", i)
		}
	}
}

func TestRefillPanic(t *testing.T) {
	fetchMore := func(count int) ([]interface{}, error) {
		panic("something")
	}

	myCache := New(10, fetchMore)
	result, err := myCache.Fetch(5)

	if result != nil {
		t.Error("Result should have been nil and is: ", result)
	}
	if err == nil {
		t.Error("Error should not have been nil")
	}

	_, ok := err.(*RefillError)
	if !ok {
		t.Error("Returned error should have been a RefillError")
	}
}

func TestCountWorksProperly(t *testing.T) {
	cache := InMemoryQueueCacheOverlay{maxSize: 100}

	cache.readCursor = 0
	cache.writeCursor = 0
	if cache.Count() != 0 {
		t.Error("Count should be 0 and is: ", cache.Count())
	}

	cache.readCursor = 0
	cache.writeCursor = 1
	if cache.Count() != 1 {
		t.Error("Count should be 1 and is: ", cache.Count())
	}

	cache.readCursor = 50
	cache.writeCursor = 99
	if cache.Count() != 49 {
		t.Error("Count should be 49 and is: ", cache.Count())
	}

	cache.readCursor = 50
	cache.writeCursor = 20
	if cache.Count() != 70 {
		t.Error("Count should be 69 and is: ", cache.Count())
	}
}
