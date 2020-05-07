package common

import (
	"fmt"
	"testing"
)

func TestPartition(t *testing.T) {
	maxItems := 10000
	keys := make([]string, 0)
	partition := Partition(keys, maxItems)
	if len(partition) != 0 {
		t.Error("Unexpected quantity of partition")
	}

	keys2 := make([]string, 0, 10000)
	for i := 0; i < 300; i++ {
		keys2 = append(keys2, fmt.Sprintf("test_%d", i))
	}

	partition2 := Partition(keys2, maxItems)
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

	partition3 := Partition(keys3, maxItems)
	if len(partition3) != 2 {
		t.Error("Unexpected quantity of partition")
	}
	if len(partition3[0]) != 10000 || len(partition3[1]) != 5000 {
		t.Error("Unexpected items per chunk")
	}
}
