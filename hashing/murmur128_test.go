package hashing

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMurmur128(t *testing.T) {
	raw, err := os.ReadFile("../testdata/murmur3_64_uuids.csv")
	assert.Nil(t, err)

	lines := strings.Split(string(raw), "\n")
	for _, line := range lines {
		if line == "" || line == "\n" {
			continue
		}
		fields := strings.Split(line, ",")
		seed, _ := strconv.ParseUint(fields[1], 10, 64)
		expected, _ := strconv.ParseInt(fields[2], 10, 64)

		h1, _ := Sum128WithSeed([]byte(fields[0]), uint32(seed))
		assert.Equal(t, expected, int64(h1))
	}
}
