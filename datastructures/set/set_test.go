package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	s := New[int](10)
	assert.Equal(t, int(0), s.Len())
	s.Add(1, 2, 3, 4, 5, 6, 7, 8, 9)
	assert.Equal(t, int(9), s.Len())
	for i := 1; i <= 9; i++ {
		assert.True(t, s.Contains(i))
	}
	assert.False(t, s.Contains(0))
	assert.False(t, s.Contains(10))

	s.Add(0)
	assert.True(t, s.Contains(0))
	s.Remove(0)
	assert.False(t, s.Contains(0))

	assert.Equal(t, s, s.Clone())

	s.Intersect(Define(5, 245))
	assert.Equal(t, Define(int(5)), s)

	s.Union(Define(10, 20, 30))
	assert.Equal(t, Define(5, 10, 20, 30), s)

	assert.Equal(t, Define(1, 2, 3), Define(3, 2, 1))
	assert.Equal(t, Define(1), Define(1, 1, 1))

	asSlice := s.ToSlice()
	assert.ElementsMatch(t, []int{5, 10, 20, 30}, asSlice)
	n := New[int](4)
	n.AddFromSlice(asSlice)
	assert.Equal(t, Define(5, 10, 20, 30), n)

	e := New[int](10)
	e.Add(1, 2, 3)
	assert.True(t, e.IsSubSet(Define(1, 2, 3, 4, 5, 6, 7, 8, 9)))
	assert.False(t, e.IsSubSet(Define(1, 2)))
}
