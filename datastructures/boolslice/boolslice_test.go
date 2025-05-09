package boolslice

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolSlice(t *testing.T) {
	_, err := NewBoolSlice(12)
    assert.NotNil(t, err)

	b, err := NewBoolSlice(int(math.Pow(2, 15)))
    assert.Nil(t, err)

	i1 := 12
	i2 := 20
	i3 := 123
	i4 := 2000
	i5 := 8192

    assert.Equal(t, ErrorOutOfBounds, b.Set(int(math.Pow(2, 15)) + 1))
	assert.Nil(t, b.Set(i1))
	assert.Nil(t, b.Set(i2))
	assert.Nil(t, b.Set(i3))
	assert.Nil(t, b.Set(i4))
	assert.Nil(t, b.Set(i5))

    set, err := b.Get(int(math.Pow(2, 15)) + 1)
    assert.False(t, set)
    assert.Equal(t, ErrorOutOfBounds, err)

    for _, i := range []int{i1, i2, i3, i4, i5} {
        res, err := b.Get(i)
        assert.Nil(t, err)
        assert.True(t, res)
    }

    for _, i := range []int{200, 500} {
        res, err := b.Get(i)
        assert.Nil(t, err)
        assert.False(t, res)
    }

    assert.Equal(t, math.Pow(2, 15)/8, float64(len(b.Bytes())))
    assert.Equal(t, ErrorOutOfBounds, b.Clear(int(math.Pow(2, 15)) + 1))
    assert.Nil(t, b.Clear(i1))

    v, err := b.Get(i1)
    assert.Nil(t, err)
    assert.False(t, v)

    res, err := Rebuild(1, nil)
    assert.Nil(t, res)
    assert.NotNil(t, err)

    res, err = Rebuild(8, nil)
    assert.Nil(t, res)
    assert.NotNil(t, err)

	rebuilt, err := Rebuild(int(math.Pow(2, 15)), b.Bytes())
    assert.Nil(t, err)

    for _, i := range []int{i2, i3, i4, i5} {
        res, err := rebuilt.Get(i)
        assert.Nil(t, err)
        assert.True(t, res)
    }

    for _, i := range []int{200, 5000} {
        res, err := rebuilt.Get(i)
        assert.Nil(t, err)
        assert.False(t, res)
    }
}
