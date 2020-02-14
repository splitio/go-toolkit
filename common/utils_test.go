package common

import (
	"fmt"
	"testing"
)

func TestGenerateChunks(t *testing.T) {
	maxItems := 10000
	keys := make([]string, 0)
	chunks := GenerateChunks(keys, maxItems)
	if len(chunks) != 0 {
		t.Error("Unexpected quantity of chunks")
	}

	keys2 := make([]string, 0, 10000)
	for i := 0; i < 300; i++ {
		keys2 = append(keys2, fmt.Sprintf("test_%d", i))
	}

	chunks2 := GenerateChunks(keys2, maxItems)
	if len(chunks2) != 1 {
		t.Error("Unexpected quantity of chunks")
	}
	if len(chunks2[0]) != 300 {
		t.Error("Unexpected items per chunk")
	}

	keys3 := make([]string, 0, 15000)
	for i := 0; i < 15000; i++ {
		keys3 = append(keys3, fmt.Sprintf("test_%d", i))
	}

	chunks3 := GenerateChunks(keys3, maxItems)
	if len(chunks3) != 2 {
		t.Error("Unexpected quantity of chunks")
	}
	if len(chunks3[0]) != 10000 || len(chunks3[1]) != 5000 {
		t.Error("Unexpected items per chunk")
	}
}
