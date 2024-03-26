package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRef(t *testing.T) {
	str1 := "hello"
	assert.Equal(t, &str1, Ref[string]("hello"))
	i64 := int64(123456)
	assert.Equal(t, &i64, Ref[int64](123456))
}

func TestRefOrNil(t *testing.T) {
	str1 := "hello"
	assert.Equal(t, &str1, RefOrNil("hello"))
	assert.Equal(t, (*string)(nil), RefOrNil(""))
}

func TestPartitionSliceByLength(t *testing.T) {
	maxItems := 10000
	keys := make([]string, 0)
	partition := PartitionSliceByLength(keys, maxItems)
	if len(partition) != 0 {
		t.Error("Unexpected quantity of partition")
	}

	keys2 := make([]string, 0, 10000)
	for i := 0; i < 300; i++ {
		keys2 = append(keys2, fmt.Sprintf("test_%d", i))
	}

	partition2 := PartitionSliceByLength(keys2, maxItems)
	if len(partition2) != 1 {
		t.Error("Unexpected quantity of partition")
	}
	if len(partition2[0]) != 300 {
		t.Error("Unexpected items per chunk")
	}

	keys3 := make([]string, 0, 15000)
	for i := 0; i < 15000; i++ {
		keys3 = append(keys3, fmt.Sprintf("test_%d", i))
	}

	partition3 := PartitionSliceByLength(keys3, maxItems)
	if len(partition3) != 2 {
		t.Error("Unexpected quantity of partition")
	}
	if len(partition3[0]) != 10000 || len(partition3[1]) != 5000 {
		t.Error("Unexpected items per chunk")
	}
}

func TestDedupeInNewSlice(t *testing.T) {
	assert.ElementsMatch(t, []int{1, 2, 3}, UnorderedDedupedCopy([]int{3, 2, 2, 3, 1}))
	assert.ElementsMatch(t, []int{1, 2, 3}, UnorderedDedupedCopy([]int{1, 2, 3}))
	assert.ElementsMatch(t, []int{}, UnorderedDedupedCopy([]int{}))
	assert.ElementsMatch(t, []string{"a", "c"}, UnorderedDedupedCopy([]string{"c", "c", "a"}))
}

func TestValueOr(t *testing.T) {
	assert.Equal(t, int64(3), ValueOr[int64](0, 3))
	assert.Equal(t, (*int)(nil), ValueOr[*int](nil, nil))
	assert.Equal(t, Ref(int(3)), ValueOr[*int](nil, Ref(int(3))))
	assert.Equal(t, Ref(int(4)), ValueOr[*int](Ref(int(4)), Ref(int(3))))
}

func TestMax(t *testing.T) {
	assert.Equal(t, int(5), Max[int](1, 2, 3, 4, 5))
	assert.Equal(t, int(5), Max[int](5))
}

func TestMin(t *testing.T) {
	assert.Equal(t, int(1), Min[int](1, 2, 3, 4, 5))
	assert.Equal(t, int(1), Min[int](1))
}
