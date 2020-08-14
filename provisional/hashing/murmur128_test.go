package hashing

import (
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
)

func TestMurmur128(t *testing.T) {
	raw, err := ioutil.ReadFile("../../testfiles/murmur3_64_uuids.csv")
	if err != nil {
		t.Error("error reading murmur128 test cases files: ", err.Error())
	}

	lines := strings.Split(string(raw), "\n")
	for _, line := range lines {
		if line == "" || line == "\n" {
			continue
		}
		fields := strings.Split(line, ",")
		seed, _ := strconv.ParseUint(fields[1], 10, 64)
		expected, _ := strconv.ParseInt(fields[2], 10, 64)

		h1, _ := Sum128WithSeed([]byte(fields[0]), uint32(seed))
		if int64(h1) != expected {
			t.Errorf("Hashes don't match. Expected: %d, actual: %d", expected, uint64(h1))
		}
	}
}
