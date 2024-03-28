package queuecache

import (
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
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
		assert.True(t, ok)
		assert.Equal(t, index, asInt)
	}

	offset := 5
	next5, err := myCache.Fetch(5)
	assert.Nil(t, err)
	for index, item := range next5 {
		asInt, ok := item.(int)
		assert.True(t, ok)
		assert.Equal(t, index+offset, asInt)
	}

	index = 0
	myCache = New(10, fetchMore)
	for i := 0; i < 100; i++ {
		elem, err := myCache.Fetch(1)
        assert.Nil(t, err)
		asInt, ok := elem[0].(int)
        assert.True(t, ok)
        assert.Equal(t, i, asInt)
	}

	elems, err := myCache.Fetch(1)
    assert.Nil(t, elems)
    assert.ErrorContains(t, err, "NO_MORE_DATA")

	// Set index to 0 so that refill works and restart tests.
	index = 0
	for i := 0; i < 100; i++ {
		elem, err := myCache.Fetch(1)
        assert.Nil(t, err)

		asInt, ok := elem[0].(int)
        assert.True(t, ok)
        assert.Equal(t, i, asInt)
	}
}

func TestRefillPanic(t *testing.T) {
	fetchMore := func(count int) ([]interface{}, error) {
		panic("something")
	}

	myCache := New(10, fetchMore)
	result, err := myCache.Fetch(5)
    assert.Nil(t, result)
    assert.NotNil(t, err)

	_, ok := err.(*RefillError)
    assert.True(t, ok)
}

func TestCountWorksProperly(t *testing.T) {
	cache := InMemoryQueueCacheOverlay[int]{maxSize: 100}

	cache.readCursor = 0
	cache.writeCursor = 0
    assert.Equal(t, 0, cache.Count())

	cache.readCursor = 0
	cache.writeCursor = 1
    assert.Equal(t, 1, cache.Count())

	cache.readCursor = 50
	cache.writeCursor = 99
    assert.Equal(t, 49, cache.Count())

	cache.readCursor = 50
	cache.writeCursor = 20
    assert.Equal(t, 70, cache.Count())
}
