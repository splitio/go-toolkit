package hashing

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

)

func TestMurmurHashOnAlphanumericData(t *testing.T) {
	inFile, err := os.Open("../testdata/murmur3-sample-data-v2.csv")
    assert.Nil(t, err)
	defer inFile.Close()

	reader := csv.NewReader(bufio.NewReader(inFile))

	var arr []string
	line := 0
	for err != io.EOF {
		line++
		arr, err = reader.Read()
		if len(arr) < 4 {
			continue // Skip empty lines
		}
		seed, _ := strconv.ParseInt(arr[0], 10, 32)
		str := arr[1]
		digest, _ := strconv.ParseUint(arr[2], 10, 32)

		calculated := NewMurmur332Hasher(uint32(seed)).Hash([]byte(str))
        assert.Equal(t, calculated, uint32(digest))
	}
}

func TestMurmurHashOnNonAlphanumericData(t *testing.T) {
	inFile, err := os.Open("../testdata/murmur3-sample-data-non-alpha-numeric-v2.csv")
    assert.Nil(t, err)
	defer inFile.Close()

	reader := csv.NewReader(bufio.NewReader(inFile))

	var arr []string
	line := 0
	for err != io.EOF {
		line++
		arr, err = reader.Read()
		if len(arr) < 4 {
			continue // Skip empty lines
		}
		seed, _ := strconv.ParseInt(arr[0], 10, 32)
		str := arr[1]
		digest, _ := strconv.ParseUint(arr[2], 10, 32)

		calculated := NewMurmur332Hasher(uint32(seed)).Hash([]byte(str))
        assert.Equal(t, calculated, uint32(digest))
	}
}
